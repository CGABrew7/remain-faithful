import ReplayKit
import SensitiveContentAnalysis
import Vision
import CoreImage
import NaturalLanguage
import os

// MARK: - Shared event model (mirrors DetectedEvent in the main app target)

// IMPORTANT: This struct must remain byte-for-byte identical to DetectedEvent
// in RemainFaithful/EventProcessor.swift. Both targets share App Group UserDefaults
// for IPC. A future refactor should move this to a shared Swift package.
struct DetectedEvent: Codable {
    let id: String
    let timestamp: Double
    let category: String  // "adult_content" | "gambling" | "violence" | "self_harm" | "clean"
    let severity: String  // "informational" | "concerning" | "severe"
    let summary: String
    let tier: Int         // 1, 2, or 3
    let confidence: Float
    var processed: Bool
}

// MARK: - Content category

enum ContentCategory: String {
    case explicitSexual = "adult_content"
    case gambling       = "gambling"
    case violence       = "violence"
    case selfHarm       = "self_harm"
    case clean          = "clean"

    func severity(for confidence: Float) -> String {
        guard self != .clean else { return "informational" }
        return confidence >= 0.75 ? "severe" : "concerning"
    }
}

// MARK: - Text classification result

struct TextClassification {
    let category: ContentCategory
    let confidence: Float
}

// MARK: - Fallback text classifier
//
// Implements the same interface a CoreML NLModel would provide.
// Uses the 20+ inline training examples per category as weighted keyword anchors.
// A production build would swap this for a bundled .mlmodel trained with CreateML
// on the same example sets.

final class FallbackTextClassifier {

    // 20+ representative training examples per category (used as keyword anchors).
    private let trainingData: [ContentCategory: [(term: String, weight: Float)]] = [

        .explicitSexual: [
            ("pornhub", 1.0), ("xvideos", 1.0), ("xhamster", 1.0), ("onlyfans", 0.9),
            ("chaturbate", 1.0), ("livejasmin", 0.9), ("cam4.com", 1.0),
            ("adult content", 0.9), ("xxx", 0.9), ("nude", 0.7), ("naked", 0.7),
            ("pornography", 1.0), ("erotic", 0.7), ("sex video", 0.9),
            ("watch porn", 1.0), ("free porn", 1.0), ("hot girls", 0.6),
            ("barely legal", 0.9), ("nsfw", 0.7), ("cam girl", 0.85),
            ("strip club", 0.8), ("escort", 0.75), ("masturbat", 0.9),
            ("explicit", 0.5), ("orgasm", 0.85),
        ],

        .gambling: [
            ("draftkings", 0.95), ("fanduel", 0.95), ("betmgm", 0.95),
            ("bovada", 0.9), ("pokerstars", 0.9), ("bet365", 0.9),
            ("sports betting", 0.9), ("place a bet", 0.85), ("casino", 0.75),
            ("poker chips", 0.8), ("blackjack", 0.75), ("roulette", 0.75),
            ("slot machine", 0.8), ("jackpot", 0.7), ("wager", 0.75),
            ("sportsbook", 0.9), ("free bet", 0.8), ("live betting", 0.85),
            ("parlay", 0.8), ("moneyline", 0.8), ("point spread", 0.8),
            ("gambling", 0.9), ("odds calculator", 0.8),
        ],

        .violence: [
            ("gore", 1.0), ("graphic violence", 1.0), ("execution", 0.8),
            ("beheading", 1.0), ("massacre", 0.9), ("snuff", 0.95),
            ("murder video", 1.0), ("death footage", 0.9), ("brutal attack", 0.85),
            ("shooting video", 0.8), ("war crimes", 0.8), ("torture", 0.85),
            ("gunshot wound", 0.85), ("dismember", 0.9), ("decapitat", 0.95),
            ("slaughter", 0.8), ("carnage", 0.8), ("graphic death", 0.9),
            ("liveleak", 0.85), ("bestgore", 1.0), ("shock site", 0.85),
            ("violent video", 0.85),
        ],

        .selfHarm: [
            ("suicide method", 1.0), ("how to kill myself", 1.0),
            ("self harm", 0.85), ("self-harm", 0.85), ("cutting myself", 0.9),
            ("overdose on", 0.9), ("end my life", 0.95), ("suicidal", 0.8),
            ("razor blade cut", 0.9), ("painless death", 0.9),
            ("no reason to live", 0.9), ("goodbye note", 0.85),
            ("ending it all", 0.85), ("don't want to live", 0.9),
            ("ways to die", 0.85), ("pills to die", 0.95),
            ("hang myself", 0.95), ("jump off bridge", 0.9),
            ("pro-ana", 0.85), ("starvation tips", 0.85), ("body check", 0.75),
            ("trigger warning self", 0.7),
        ],
    ]

    func classify(_ text: String) -> TextClassification {
        let lower = text.lowercased()
        var scores = [ContentCategory: Float]()

        for (category, terms) in trainingData {
            var score: Float = 0
            for (term, weight) in terms where lower.contains(term) {
                score += weight
            }
            if score > 0 {
                // Non-linear normalization: first match is impactful, additional matches add less
                let normalized = min(1.0, 0.35 + score * 0.15)
                scores[category] = normalized
            }
        }

        if let top = scores.max(by: { $0.value < $1.value }) {
            return TextClassification(category: top.key, confidence: top.value)
        }
        return TextClassification(category: .clean, confidence: 0.95)
    }
}

// MARK: - Tier 1 rules

private struct Tier1Rules {

    // Known adult/gambling/harmful domains.
    static let urlBlocklist: Set<String> = [
        "pornhub.com", "xvideos.com", "xhamster.com", "xnxx.com", "redtube.com",
        "youporn.com", "tube8.com", "spankbang.com", "ixxx.com", "beeg.com",
        "brazzers.com", "realitykings.com", "bangbros.com", "onlyfans.com",
        "chaturbate.com", "livejasmin.com", "cam4.com", "stripchat.com",
        "adultfriendfinder.com", "ashleymadison.com", "benaughty.com",
        "seeking.com", "seekingarrangement.com",
        "bovada.lv", "bet365.com", "draftkings.com", "fanduel.com",
        "pokerstars.com", "betmgm.com", "barstoolsportsbook.com",
        "bestgore.com", "liveleak.com",
    ]

    // Explicit keyword patterns tested against full OCR text.
    static let keywordPatterns: [NSRegularExpression] = [
        "\\b(porn|pornography|xxx|adult.?content)\\b",
        "\\b(nude|naked|erotic|nsfw)\\b",
        "\\b(sex.?video|free.?porn|watch.?porn|hot.?girls|cam.?girl)\\b",
        "\\b(onlyfans|chaturbate|livejasmin|pornhub|xvideos|xhamster)\\b",
        "\\b(tinder|grindr|ashley.?madison|seeking.?arrangement|adult.?friend.?finder)\\b",
        "\\b(casino|sports.?betting|place.?a.?bet|draftkings|fanduel|betmgm|sportsbook)\\b",
        "\\b(gore|beheading|snuff|bestgore|liveleak|graphic.?violence|execution.?video)\\b",
        "\\b(suicide.?method|how.?to.?kill.?myself|self.?harm|end.?my.?life|suicidal)\\b",
    ].compactMap { try? NSRegularExpression(pattern: $0, options: .caseInsensitive) }

    static func matchesURL(_ domain: String) -> Bool {
        let clean = domain.hasPrefix("www.") ? String(domain.dropFirst(4)) : domain
        return urlBlocklist.contains { clean == $0 || clean.hasSuffix("." + $0) }
    }

    static func matchesKeyword(in text: String) -> Bool {
        let lower = text.lowercased()
        let range = NSRange(lower.startIndex..., in: lower)
        return keywordPatterns.contains { $0.firstMatch(in: lower, options: [], range: range) != nil }
    }
}

// MARK: - Perceptual hash (dHash — 64-bit difference hash)

private func dHash(from ciImage: CIImage, context: CIContext) -> UInt64 {
    // Scale to 9×8 pixels for 64-bit hash
    let target = CGRect(x: 0, y: 0, width: 9, height: 8)
    let scaled = ciImage.transformed(by: CGAffineTransform(
        scaleX: 9.0 / ciImage.extent.width,
        y: 8.0 / ciImage.extent.height
    )).cropped(to: target)

    // Desaturate
    let gray = scaled.applyingFilter("CIPhotoEffectNoir")

    // Render to 9×8 RGBA8 bitmap
    var pixels = [UInt8](repeating: 0, count: 9 * 8 * 4)
    context.render(gray,
                   toBitmap: &pixels,
                   rowBytes: 9 * 4,
                   bounds: target,
                   format: .RGBA8,
                   colorSpace: nil)

    // Build hash: compare adjacent horizontal pixels (red channel)
    var hash: UInt64 = 0
    for row in 0..<8 {
        for col in 0..<8 {
            if pixels[(row * 9 + col) * 4] < pixels[(row * 9 + col + 1) * 4] {
                hash |= 1 << UInt64(row * 8 + col)
            }
        }
    }
    return hash
}

private func hammingDistance(_ a: UInt64, _ b: UInt64) -> Int {
    (a ^ b).nonzeroBitCount
}

// MARK: - Tier 3 cloud fallback

private func sendToClassify(text: String, apiBase: String, defaults: UserDefaults?) async -> TextClassification? {
    guard !text.isEmpty else { return nil }
    guard let url = URL(string: apiBase + "/classify") else { return nil }
    var req = URLRequest(url: url, timeoutInterval: 10)
    req.httpMethod = "POST"
    req.setValue("application/json", forHTTPHeaderField: "Content-Type")
    guard let body = try? JSONSerialization.data(withJSONObject: ["text": String(text.prefix(500))]) else {
        return nil
    }
    if let secret = defaults?.string(forKey: "classifySecret"), !secret.isEmpty {
        req.setValue(secret, forHTTPHeaderField: "X-Classify-Secret")
    }
    req.httpBody = body
    guard let (data, _) = try? await URLSession.shared.data(for: req),
          let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
          let cat = json["category"] as? String,
          let conf = json["confidence"] as? Double,
          let category = ContentCategory(rawValue: cat) else { return nil }
    return TextClassification(category: category, confidence: Float(conf))
}

// MARK: - Dedup store
//
// Thread-safe via Swift actor. Tracks 5-minute dedup windows per category|severity key.
// record() returns true (suppress) if an identical event fired within the current window.
// drainExpired() returns and clears windows that have elapsed with unsent suppressions
// so a "continued activity" summary can be posted.

private actor DedupStore {
    struct Entry {
        var windowStart: Date
        var suppressedCount: Int
    }

    private let windowSeconds: TimeInterval
    private var state: [String: Entry] = [:]

    init(windowSeconds: TimeInterval) {
        self.windowSeconds = windowSeconds
    }

    func record(key: String, now: Date = Date()) -> Bool {
        if let entry = state[key], now.timeIntervalSince(entry.windowStart) < windowSeconds {
            state[key]!.suppressedCount += 1
            return true  // suppress
        }
        state[key] = Entry(windowStart: now, suppressedCount: 0)
        return false  // new window — post immediately
    }

    func drainExpired(now: Date = Date()) -> [(category: String, severity: String, count: Int)] {
        var result: [(String, String, Int)] = []
        for (key, entry) in state where
                now.timeIntervalSince(entry.windowStart) >= windowSeconds &&
                entry.suppressedCount > 0 {
            let parts = key.components(separatedBy: "|")
            result.append((parts[0], parts.count > 1 ? parts[1] : "concerning", entry.suppressedCount))
            state.removeValue(forKey: key)
        }
        return result
    }
}

// MARK: - SampleHandler

private let logger = Logger(subsystem: "com.remainfaithful.app.broadcast",
                            category: "classification")

class SampleHandler: RPBroadcastSampleHandler {

    private let appGroupID = "group.com.remainfaithful.app"

    private lazy var sharedDefaults: UserDefaults? = UserDefaults(suiteName: appGroupID)

    private let scaAnalyzer = SCSensitivityAnalyzer()
    private let textClassifier = FallbackTextClassifier()
    private let ciContext = CIContext(options: [.useSoftwareRenderer: false])

    // 3-second intervals balance detection accuracy with battery efficiency.
    private var lastAnalysisTime: Date = Date()

    // Rate-limit tier 1a hash detections before dedup takes over (avoids spawning 30 tasks/sec).
    private var lastHashMatchTime: Date = .distantPast

    // Heartbeat: track last frame to report active vs idle every 2 minutes.
    private var _lastFrameTime: Date = .distantPast
    private let frameTimeLock = NSLock()
    private var lastFrameTime: Date {
        get { frameTimeLock.withLock { _lastFrameTime } }
        set { frameTimeLock.withLock { _lastFrameTime = newValue } }
    }
    private var heartbeatTask: Task<Void, Never>?
    private var retryTask:     Task<Void, Never>?

    // Hash blocklist — populated from server-side list at broadcast start.
    private var hashBlocklist: Set<UInt64> = []

    // Dedup: identical category|severity events within a 5-minute window are suppressed.
    // At most one "continued activity" summary is sent per expired window.
    private let dedupStore = DedupStore(windowSeconds: 300)

    // MARK: - Lifecycle

    override func broadcastStarted(withSetupInfo setupInfo: [String: NSObject]?) {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
        sharedDefaults?.set(Date().timeIntervalSince1970, forKey: "broadcastStartTime")
        loadHashBlocklist()
        let policy = scaAnalyzer.analysisPolicy
        logger.info("Broadcast started — SCA policy: \(String(describing: policy))")
        startHeartbeatLoop()
        startRetryLoop()
    }

    override func broadcastPaused() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
    }

    override func broadcastResumed() {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
    }

    override func broadcastFinished() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
        heartbeatTask?.cancel()
        retryTask?.cancel()
        logger.info("Broadcast finished")
    }

    // MARK: - Frame processing

    override func processSampleBuffer(_ sampleBuffer: CMSampleBuffer,
                                      with sampleBufferType: RPSampleBufferType) {
        guard sampleBufferType == .video,
              let pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else { return }

        let ciImage = CIImage(cvPixelBuffer: pixelBuffer)
        lastFrameTime = Date()

        // — Tier 1a: perceptual hash check (every frame, ~0.1 ms) —
        logger.debug("Frame received")
        let hash = dHash(from: ciImage, context: ciContext)
        let now = Date()
        if hashBlocklist.contains(where: { hammingDistance($0, hash) < 10 }) {
            // Rate-limit to one task per 3s; the dedup store handles the 5-min window.
            if now.timeIntervalSince(lastHashMatchTime) >= 3.0 {
                lastHashMatchTime = now
                logger.warning("Tier 1 hash match")
                let event = makeEvent(
                    category: .explicitSexual, tier: 1, confidence: 1.0,
                    summary: "Explicit image detected (perceptual hash match)"
                )
                Task.detached(priority: .utility) { [weak self] in
                    await self?.handleEvent(event)
                }
            }
        }

        // — Full Tier 2 processing every 3 seconds —
        guard now.timeIntervalSince(lastAnalysisTime) >= 3.0 else { return }
        lastAnalysisTime = now
        logger.info("Tier 2 analysis triggered")

        // Render CGImage once; reused by OCR, SCA, and hash update
        guard let cgImage = ciContext.createCGImage(ciImage, from: ciImage.extent) else { return }

        let capturedDefaults = sharedDefaults
        let capturedAnalyzer = scaAnalyzer
        let capturedClassifier = textClassifier
        let apiBase = capturedDefaults?.string(forKey: "apiBaseURL") ?? "http://localhost:8080"

        Task.detached(priority: .utility) { [self] in
            await self.analyzeFrame(
                cgImage: cgImage, ciImage: ciImage,
                scaAnalyzer: capturedAnalyzer,
                textClassifier: capturedClassifier,
                defaults: capturedDefaults,
                apiBase: apiBase
            )
        }
    }

    // MARK: - Full frame analysis (async, off sample-buffer thread)

    private func analyzeFrame(
        cgImage: CGImage,
        ciImage: CIImage,
        scaAnalyzer: SCSensitivityAnalyzer,
        textClassifier: FallbackTextClassifier,
        defaults: UserDefaults?,
        apiBase: String
    ) async {
        // ——— TIER 2a: Vision OCR ———
        let ocrText = await extractText(from: cgImage)
        logger.info("OCR: \(ocrText.count) chars — \(String(ocrText.prefix(120)).debugDescription)")

        // ——— TIER 1b: URL + keyword check on extracted text ———
        var tier1Triggered = false
        let domains = extractDomains(from: ocrText)
        for domain in domains where Tier1Rules.matchesURL(domain) {
            logger.warning("Tier 1 URL match: \(domain)")
            let event = makeEvent(
                category: .explicitSexual, tier: 1, confidence: 1.0,
                summary: "Blocked domain detected"
            )
            await handleEvent(event)
            tier1Triggered = true
        }
        if !tier1Triggered && Tier1Rules.matchesKeyword(in: ocrText) {
            logger.warning("Tier 1 keyword match: \(String(ocrText.prefix(80)))")
            let event = makeEvent(
                category: .explicitSexual, tier: 1, confidence: 0.95,
                summary: "Explicit keyword detected in screen text"
            )
            await handleEvent(event)
            tier1Triggered = true
        }

        // ——— TIER 2b: SensitiveContentAnalysis (image) ———
        var scaIsSensitive = false
        if scaAnalyzer.analysisPolicy != .disabled {
            let tmpURL = FileManager.default.temporaryDirectory
                .appendingPathComponent("rf_sca_\(UUID().uuidString).jpg")
            defer { try? FileManager.default.removeItem(at: tmpURL) }
            if let jpeg = cgImageToJPEG(cgImage) {
                try? jpeg.write(to: tmpURL, options: .atomic)
                if let analysis = try? await scaAnalyzer.analyzeImage(at: tmpURL) {
                    scaIsSensitive = analysis.isSensitive
                    logger.info("SCA: sensitive=\(analysis.isSensitive)")
                }
            }
        } else {
            logger.debug("SCA: policy=disabled, skipping")
        }

        // ——— TIER 2c: Text classifier ———
        let textResult = textClassifier.classify(ocrText)
        logger.info("Text classifier: \(textResult.category.rawValue) @ \(String(format: "%.2f", textResult.confidence))")

        // ——— Composite score ———
        let scaScore: Float = scaIsSensitive ? 0.85 : 0.0
        var compositeCategory = textResult.category
        var compositeConfidence: Float

        if textResult.category != .clean {
            compositeConfidence = max(scaScore, textResult.confidence)
        } else if scaIsSensitive {
            compositeCategory  = .explicitSexual
            compositeConfidence = scaScore
        } else {
            compositeConfidence = 0.0
        }

        logger.info("Composite: \(compositeCategory.rawValue) @ \(String(format: "%.2f", compositeConfidence))")

        // ——— Route by confidence ———
        if compositeConfidence >= 0.3 && compositeCategory != .clean {
            // ——— TIER 3: cloud fallback for ambiguous zone ———
            var finalCategory   = compositeCategory
            var finalConfidence = compositeConfidence

            if compositeConfidence >= 0.5 && compositeConfidence < 0.7 {
                logger.info("Sending to Tier 3 /classify — confidence \(String(format: "%.2f", compositeConfidence))")
                if let cloud = await sendToClassify(text: ocrText, apiBase: apiBase, defaults: defaults) {
                    finalCategory   = cloud.category
                    finalConfidence = cloud.confidence
                    logger.info("Tier 3 result: \(cloud.category.rawValue) @ \(String(format: "%.2f", cloud.confidence))")
                }
            }

            if finalCategory != .clean && finalConfidence >= 0.3 {
                let tier = compositeConfidence >= 0.7 ? 2 : (finalConfidence >= 0.7 ? 3 : 2)
                let event = makeEvent(
                    category: finalCategory,
                    tier: tier,
                    confidence: finalConfidence,
                    summary: buildSummary(category: finalCategory, confidence: finalConfidence)
                )
                await handleEvent(event)
            }
        }
    }

    // MARK: - Event handling with dedup

    // Check the 5-minute dedup window. If the same category|severity fired within the
    // current window, suppress. Otherwise post immediately and start a new window.
    private func handleEvent(_ event: DetectedEvent) async {
        let key = "\(event.category)|\(event.severity)"
        let suppressed = await dedupStore.record(key: key)
        if suppressed {
            logger.info("Dedup suppressed \(key) within 5-min window")
            return
        }
        logger.info("Dedup new window for \(key) — posting to backend")
        await postEvent(event)
    }

    // MARK: - Direct upload to backend

    private func uploadEvent(_ event: DetectedEvent) async -> Bool {
        guard let defaults = sharedDefaults,
              let apiBase  = defaults.string(forKey: "apiBaseURL"),
              let token    = defaults.string(forKey: "authToken"),
              let url      = URL(string: apiBase + "/events") else {
            logger.warning("uploadEvent: missing apiBaseURL or authToken")
            return false
        }
        let isoFmt = ISO8601DateFormatter()
        isoFmt.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        let tsStr = isoFmt.string(from: Date(timeIntervalSince1970: event.timestamp))

        var req = URLRequest(url: url, timeoutInterval: 10)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        let body: [String: Any] = [
            "category":  event.category,
            "severity":  event.severity,
            "summary":   event.summary,
            "timestamp": tsStr,
        ]
        guard let data = try? JSONSerialization.data(withJSONObject: body) else { return false }
        req.httpBody = data

        do {
            let (_, resp) = try await URLSession.shared.data(for: req)
            if let http = resp as? HTTPURLResponse, (200...299).contains(http.statusCode) {
                logger.info("Event uploaded: \(event.category) \(event.severity)")
                return true
            }
            let code = (resp as? HTTPURLResponse)?.statusCode ?? 0
            logger.warning("Event upload HTTP \(code)")
            return false
        } catch {
            logger.warning("Event upload network error: \(error.localizedDescription)")
            return false
        }
    }

    private func postEvent(_ event: DetectedEvent) async {
        let ok = await uploadEvent(event)
        if !ok {
            logger.warning("Event upload failed — writing to retry queue")
            writeEventToRetryQueue(event)
        }
    }

    // MARK: - Retry queue (network failures only)
    //
    // EventProcessor in the main app also drains this on foreground / BGAppRefreshTask,
    // providing a second-chance drain if the extension is killed before retrying.

    private func writeEventToRetryQueue(_ event: DetectedEvent) {
        guard let defaults = sharedDefaults else { return }
        var events: [DetectedEvent] = []
        if let data = defaults.data(forKey: "pendingEvents"),
           let existing = try? JSONDecoder().decode([DetectedEvent].self, from: data) {
            events = existing.filter { !$0.processed }
        }
        events.append(event)
        if events.count > 50 { events = Array(events.suffix(50)) }
        if let data = try? JSONEncoder().encode(events) {
            defaults.set(data, forKey: "pendingEvents")
            defaults.set(event.severity,  forKey: "lastEventSeverity")
            defaults.set(event.summary,   forKey: "lastEventSummary")
            defaults.set(event.timestamp, forKey: "lastEventTimestamp")
        }
    }

    // MARK: - Retry + continued-activity loop (60 s)

    private func startRetryLoop() {
        retryTask = Task.detached(priority: .utility) { [weak self] in
            while !Task.isCancelled {
                try? await Task.sleep(nanoseconds: 60 * 1_000_000_000)
                guard !Task.isCancelled else { break }
                await self?.drainRetryQueue()
                await self?.flushContinuedActivity()
            }
        }
    }

    private func drainRetryQueue() async {
        guard let defaults = sharedDefaults,
              let data   = defaults.data(forKey: "pendingEvents"),
              var events = try? JSONDecoder().decode([DetectedEvent].self, from: data) else { return }
        let pending = events.indices.filter { !events[$0].processed }
        guard !pending.isEmpty else { return }
        logger.info("Retry queue: draining \(pending.count) event(s)")
        var changed = false
        for i in pending {
            let ok = await uploadEvent(events[i])
            if ok {
                events[i].processed = true
                changed = true
            } else {
                break  // stop on first failure; try again next 60 s cycle
            }
        }
        if changed, let newData = try? JSONEncoder().encode(events) {
            defaults.set(newData, forKey: "pendingEvents")
        }
    }

    // For each 5-minute window that has elapsed with suppressed detections,
    // send exactly one "continued activity" summary event.
    private func flushContinuedActivity() async {
        let expired = await dedupStore.drainExpired()
        for (cat, _, count) in expired {
            let category = ContentCategory(rawValue: cat) ?? .explicitSexual
            let summary  = continuedActivitySummary(category: category, count: count + 1)
            let event    = makeEvent(category: category, tier: 2, confidence: 0.85, summary: summary)
            logger.info("Sending continued-activity summary: \(summary)")
            await postEvent(event)
        }
    }

    private func continuedActivitySummary(category: ContentCategory, count: Int) -> String {
        let base: String
        switch category {
        case .explicitSexual: base = "Explicit content"
        case .gambling:       base = "Gambling content"
        case .violence:       base = "Violent content"
        case .selfHarm:       base = "Self-harm content"
        case .clean:          base = "Activity"
        }
        return "\(base) — detected \(count) time\(count == 1 ? "" : "s") in 5 min"
    }

    // MARK: - Vision OCR

    private func extractText(from cgImage: CGImage) async -> String {
        await withCheckedContinuation { continuation in
            let request = VNRecognizeTextRequest { req, _ in
                let text = (req.results as? [VNRecognizedTextObservation] ?? [])
                    .compactMap { $0.topCandidates(1).first?.string }
                    .joined(separator: " ")
                continuation.resume(returning: text)
            }
            request.recognitionLevel = .fast
            request.usesLanguageCorrection = false
            let handler = VNImageRequestHandler(cgImage: cgImage, options: [:])
            try? handler.perform([request])
        }
    }

    private func extractDomains(from text: String) -> [String] {
        guard let detector = try? NSDataDetector(
            types: NSTextCheckingResult.CheckingType.link.rawValue
        ) else { return [] }
        let range = NSRange(text.startIndex..., in: text)
        return detector
            .matches(in: text, options: [], range: range)
            .compactMap { $0.url?.host?.lowercased().replacingOccurrences(of: "www.", with: "") }
    }

    // MARK: - Helpers

    private func makeEvent(category: ContentCategory, tier: Int,
                           confidence: Float, summary: String) -> DetectedEvent {
        DetectedEvent(
            id: UUID().uuidString,
            timestamp: Date().timeIntervalSince1970,
            category: category.rawValue,
            severity: category.severity(for: confidence),
            summary: summary,
            tier: tier,
            confidence: confidence,
            processed: false
        )
    }

    private func buildSummary(category: ContentCategory, confidence: Float) -> String {
        switch category {
        case .explicitSexual: return "Explicit content detected"
        case .gambling:       return "Gambling content detected"
        case .violence:       return "Violent content detected"
        case .selfHarm:       return "Self-harm content detected"
        case .clean:          return "Content reviewed — no concerns"
        }
    }

    private func cgImageToJPEG(_ cgImage: CGImage) -> Data? {
        UIImage(cgImage: cgImage).jpegData(compressionQuality: 0.70)
    }

    // MARK: - Heartbeat

    private func startHeartbeatLoop() {
        heartbeatTask = Task.detached(priority: .utility) { [weak self] in
            while !Task.isCancelled {
                try? await Task.sleep(nanoseconds: 2 * 60 * 1_000_000_000)
                guard !Task.isCancelled else { break }
                await self?.sendHeartbeat()
            }
        }
    }

    private func sendHeartbeat() async {
        let screen = Date().timeIntervalSince(lastFrameTime) < 120 ? "active" : "idle"
        let hasToken = sharedDefaults?.string(forKey: "authToken") != nil
        logger.info("Heartbeat: screen=\(screen) hasToken=\(hasToken)")
        guard let apiBase = sharedDefaults?.string(forKey: "apiBaseURL"),
              let token   = sharedDefaults?.string(forKey: "authToken"),
              let url     = URL(string: apiBase + "/heartbeat") else {
            logger.error("Heartbeat skipped — missing apiBaseURL or authToken in sharedDefaults")
            return
        }
        var req = URLRequest(url: url, timeoutInterval: 10)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        req.httpBody = try? JSONSerialization.data(withJSONObject: ["screen": screen])
        if let (_, resp) = try? await URLSession.shared.data(for: req),
           let http = resp as? HTTPURLResponse {
            logger.info("Heartbeat response: \(http.statusCode)")
        }
    }

    private func loadHashBlocklist() {
        // Main app can push known-bad hashes via shared defaults as a JSON [UInt64] array.
        guard let defaults = sharedDefaults,
              let data = defaults.data(forKey: "hashBlocklist"),
              let hashes = try? JSONDecoder().decode([UInt64].self, from: data) else { return }
        hashBlocklist = Set(hashes)
        logger.info("Loaded \(hashes.count) perceptual hashes into blocklist")
    }
}
