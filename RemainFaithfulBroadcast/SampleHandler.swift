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

final class FallbackTextClassifier {

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
            for (term, weight) in terms where lower.contains(term) { score += weight }
            if score > 0 {
                scores[category] = min(1.0, 0.35 + score * 0.15)
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
    let target = CGRect(x: 0, y: 0, width: 9, height: 8)
    let scaled = ciImage.transformed(by: CGAffineTransform(
        scaleX: 9.0 / ciImage.extent.width,
        y: 8.0 / ciImage.extent.height
    )).cropped(to: target)
    let gray = scaled.applyingFilter("CIPhotoEffectNoir")
    var pixels = [UInt8](repeating: 0, count: 9 * 8 * 4)
    context.render(gray, toBitmap: &pixels, rowBytes: 9 * 4,
                   bounds: target, format: .RGBA8, colorSpace: nil)
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

private func hammingDistance(_ a: UInt64, _ b: UInt64) -> Int { (a ^ b).nonzeroBitCount }

// MARK: - Tier 3 cloud fallback

private func sendToClassify(text: String, apiBase: String, defaults: UserDefaults?) async -> TextClassification? {
    guard !text.isEmpty, let url = URL(string: apiBase + "/classify") else { return nil }
    var req = URLRequest(url: url, timeoutInterval: 10)
    req.httpMethod = "POST"
    req.setValue("application/json", forHTTPHeaderField: "Content-Type")
    guard let body = try? JSONSerialization.data(withJSONObject: ["text": String(text.prefix(500))]) else { return nil }
    if let secret = defaults?.string(forKey: "classifySecret"), !secret.isEmpty {
        req.setValue(secret, forHTTPHeaderField: "X-Classify-Secret")
    }
    req.httpBody = body
    guard let (data, _) = try? await URLSession.shared.data(for: req),
          let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
          let cat  = json["category"] as? String,
          let conf = json["confidence"] as? Double,
          let category = ContentCategory(rawValue: cat) else { return nil }
    return TextClassification(category: category, confidence: Float(conf))
}

// MARK: - Dedup store

private actor DedupStore {
    struct Entry { var windowStart: Date; var suppressedCount: Int }
    private let windowSeconds: TimeInterval
    private var state: [String: Entry] = [:]

    init(windowSeconds: TimeInterval) { self.windowSeconds = windowSeconds }

    func record(key: String, now: Date = Date()) -> Bool {
        if let e = state[key], now.timeIntervalSince(e.windowStart) < windowSeconds {
            state[key]!.suppressedCount += 1
            return true
        }
        state[key] = Entry(windowStart: now, suppressedCount: 0)
        return false
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

// MARK: - Memory helpers
//
// os_proc_available_memory() reports bytes the process can still allocate before
// the jetsam limit (extension limit ≈ 50 MB). Negative values are not possible;
// zero means the next allocation will likely cause an OOM kill.

private func availableMemoryMB() -> Int {
    Int(os_proc_available_memory() / 1_048_576)
}

// MARK: - SampleHandler
//
// Broadcast continuation note: RPBroadcastSampleHandler runs in a separate
// process from the host app (com.remainfaithful.app). iOS continues delivering
// sample buffers through host-app backgrounding, screen lock, and screen unlock.
// The broadcast only stops when the user explicitly ends it via the iOS system
// UI or when we call finishBroadcastWithError(_:). We never call that method
// outside of a genuine error, so monitoring is never tied to app lifecycle.

private let logger = Logger(subsystem: "com.remainfaithful.app.broadcast",
                            category: "classification")

class SampleHandler: RPBroadcastSampleHandler {

    private let appGroupID = "group.com.remainfaithful.app"
    private lazy var sharedDefaults: UserDefaults? = UserDefaults(suiteName: appGroupID)

    // var so they can be recreated during the periodic resource reset.
    // Peak memory before reset: ~25–40 MB (full-res CGImage + Vision + SCA)
    // Peak memory after this change: ~14–23 MB (960px CGImage + autoreleased Vision/SCA)
    private var scaAnalyzer    = SCSensitivityAnalyzer()
    private var textClassifier = FallbackTextClassifier()
    private var ciContext      = CIContext(options: [.useSoftwareRenderer: false])

    // Maximum pixel dimension for per-frame analysis.
    // 960 px is sufficient for Vision OCR (fast mode) and SCA; reduces CGImage
    // allocation from ~8 MB (1920×1080) to ~2 MB (960×540).
    private let analysisMaxDim: CGFloat = 960

    private var lastAnalysisTime: Date  = Date()
    private var lastHashMatchTime: Date = .distantPast
    private var lastResourceReset: Date = Date()
    private var lastMemoryLogTime: Date = .distantPast

    private var _lastFrameTime: Date = .distantPast
    private let frameTimeLock = NSLock()
    private var lastFrameTime: Date {
        get { frameTimeLock.withLock { _lastFrameTime } }
        set { frameTimeLock.withLock { _lastFrameTime = newValue } }
    }

    private var heartbeatTask: Task<Void, Never>?
    private var retryTask:     Task<Void, Never>?
    private var hashBlocklist: Set<UInt64> = []
    private let dedupStore = DedupStore(windowSeconds: 300)

    // MARK: - Lifecycle

    override func broadcastStarted(withSetupInfo setupInfo: [String: NSObject]?) {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
        sharedDefaults?.set(Date().timeIntervalSince1970, forKey: "broadcastStartTime")
        loadHashBlocklist()
        let policy = scaAnalyzer.analysisPolicy
        logger.info("Broadcast started — SCA policy: \(String(describing: policy)) — available: \(availableMemoryMB()) MB")
        startHeartbeatLoop()
        startRetryLoop()
    }

    override func broadcastPaused()   { sharedDefaults?.set(false, forKey: "isBroadcasting") }
    override func broadcastResumed()  { sharedDefaults?.set(true,  forKey: "isBroadcasting") }

    override func broadcastFinished() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
        heartbeatTask?.cancel()
        retryTask?.cancel()
        logger.info("Broadcast finished — available: \(availableMemoryMB()) MB")
    }

    // MARK: - Frame processing

    override func processSampleBuffer(_ sampleBuffer: CMSampleBuffer,
                                      with sampleBufferType: RPSampleBufferType) {
        guard sampleBufferType == .video,
              let pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else { return }

        let now = Date()
        lastFrameTime = now

        // Periodic resource reset every 10 minutes — flush Core Image texture cache,
        // SCA model memory, and any ML fragmentation. Does NOT end or pause the broadcast.
        if now.timeIntervalSince(lastResourceReset) >= 600 {
            lastResourceReset = now
            performResourceReset()
        }

        // Wrap all per-frame work in an autorelease pool to flush Obj-C objects promptly.
        autoreleasepool {
            let ciImage = CIImage(cvPixelBuffer: pixelBuffer)

            // — Tier 1a: perceptual hash check (every frame, ~0.1 ms) —
            logger.debug("Frame received")
            let hash = dHash(from: ciImage, context: ciContext)
            if hashBlocklist.contains(where: { hammingDistance($0, hash) < 10 }) {
                if now.timeIntervalSince(lastHashMatchTime) >= 3.0 {
                    lastHashMatchTime = now
                    logger.warning("Tier 1 hash match — available: \(availableMemoryMB()) MB")
                    let event = makeEvent(category: .explicitSexual, tier: 1, confidence: 1.0,
                                         summary: "Explicit image detected (perceptual hash match)")
                    Task.detached(priority: .utility) { [weak self] in await self?.handleEvent(event) }
                }
            }

            // — Full Tier 2 processing every 3 seconds —
            guard now.timeIntervalSince(lastAnalysisTime) >= 3.0 else { return }
            lastAnalysisTime = now

            // Downscale to analysisMaxDim on the longer edge before rendering.
            // 1920×1080 → 960×540: CGImage drops from ~8 MB to ~2 MB (BGRA8).
            let extent = ciImage.extent
            let longerEdge = max(extent.width, extent.height)
            let scale = longerEdge > analysisMaxDim ? analysisMaxDim / longerEdge : 1.0
            let renderImage = scale < 1.0
                ? ciImage.transformed(by: CGAffineTransform(scaleX: scale, y: scale))
                : ciImage

            // Render inside nested autorelease pool; ciContext may produce
            // intermediate Obj-C objects during rasterization.
            guard let cgImage: CGImage = autoreleasepool(invoking: {
                ciContext.createCGImage(renderImage, from: renderImage.extent)
            }) else { return }

            // ciImage (pixel buffer wrapper) is no longer needed — do not capture it.
            let preMemMB = availableMemoryMB()
            logger.info("Tier 2 analysis — pre: \(preMemMB) MB available (frame \(Int(extent.width))×\(Int(extent.height)) → \(Int(renderImage.extent.width))×\(Int(renderImage.extent.height)))")

            let capturedDefaults    = sharedDefaults
            let capturedAnalyzer    = scaAnalyzer
            let capturedClassifier  = textClassifier
            let apiBase = capturedDefaults?.string(forKey: "apiBaseURL") ?? "http://localhost:8080"

            Task.detached(priority: .utility) { [self] in
                await self.analyzeFrame(
                    cgImage: cgImage,       // downscaled, ~2 MB
                    scaAnalyzer: capturedAnalyzer,
                    textClassifier: capturedClassifier,
                    defaults: capturedDefaults,
                    apiBase: apiBase,
                    preMemMB: preMemMB
                )
            }
            // cgImage is now owned exclusively by the detached task. ciImage is
            // released here as the autorelease pool drains at the end of this block.
        }
    }

    // MARK: - Periodic resource reset (every 10 minutes)
    //
    // Recreates ciContext, scaAnalyzer, and textClassifier to flush accumulated
    // Core Image texture caches, ML model memory fragmentation, and stale buffers.
    // Running tasks hold captured references to the old objects and complete normally.
    // The broadcast stream is never interrupted.

    private func performResourceReset() {
        let before = availableMemoryMB()
        ciContext      = CIContext(options: [.useSoftwareRenderer: false])
        scaAnalyzer    = SCSensitivityAnalyzer()
        textClassifier = FallbackTextClassifier()
        let after = availableMemoryMB()
        logger.info("Resource reset — available: \(before) MB → \(after) MB (Δ\(after - before) MB)")
    }

    // MARK: - Full frame analysis (async, off sample-buffer thread)

    private func analyzeFrame(
        cgImage: CGImage,            // downscaled to ≤960 px, ~2 MB
        scaAnalyzer: SCSensitivityAnalyzer,
        textClassifier: FallbackTextClassifier,
        defaults: UserDefaults?,
        apiBase: String,
        preMemMB: Int
    ) async {
        // ——— TIER 2a: Vision OCR ———
        // Wrap synchronous Vision work in autoreleasepool so VNImageRequestHandler,
        // VNRecognizeTextRequest, and VNRecognizedTextObservation are released immediately.
        let ocrText: String = autoreleasepool {
            var result = ""
            let request = VNRecognizeTextRequest { req, _ in
                result = (req.results as? [VNRecognizedTextObservation] ?? [])
                    .compactMap { $0.topCandidates(1).first?.string }
                    .joined(separator: " ")
            }
            request.recognitionLevel = .fast
            request.usesLanguageCorrection = false
            let handler = VNImageRequestHandler(cgImage: cgImage, options: [:])
            try? handler.perform([request])
            // handler, request, and all VN objects go out of scope here and are released
            return result
        }
        logger.info("OCR: \(ocrText.count) chars — \(String(ocrText.prefix(80)).debugDescription)")

        // ——— TIER 1b: URL + keyword check on extracted text ———
        var tier1Triggered = false
        let domains = extractDomains(from: ocrText)
        for domain in domains where Tier1Rules.matchesURL(domain) {
            logger.warning("Tier 1 URL match: \(domain)")
            await handleEvent(makeEvent(category: .explicitSexual, tier: 1, confidence: 1.0,
                                        summary: "Blocked domain detected"))
            tier1Triggered = true
        }
        if !tier1Triggered && Tier1Rules.matchesKeyword(in: ocrText) {
            logger.warning("Tier 1 keyword match: \(String(ocrText.prefix(80)))")
            await handleEvent(makeEvent(category: .explicitSexual, tier: 1, confidence: 0.95,
                                        summary: "Explicit keyword detected in screen text"))
            tier1Triggered = true
        }

        // ——— TIER 2b: SensitiveContentAnalysis (image) ———
        // Write JPEG inside autoreleasepool to release UIImage and Data promptly;
        // only the file URL crosses the autorelease boundary into the async SCA call.
        var scaIsSensitive = false
        if scaAnalyzer.analysisPolicy != .disabled {
            let tmpURL = FileManager.default.temporaryDirectory
                .appendingPathComponent("rf_sca_\(UUID().uuidString).jpg")
            defer { try? FileManager.default.removeItem(at: tmpURL) }
            let didWrite: Bool = autoreleasepool {
                guard let jpeg = UIImage(cgImage: cgImage).jpegData(compressionQuality: 0.70) else { return false }
                return (try? jpeg.write(to: tmpURL, options: .atomic)) != nil
                // jpeg Data and UIImage are released here
            }
            if didWrite, let analysis = try? await scaAnalyzer.analyzeImage(at: tmpURL) {
                scaIsSensitive = analysis.isSensitive
                logger.info("SCA: sensitive=\(analysis.isSensitive)")
            }
        } else {
            logger.debug("SCA: policy=disabled, skipping")
        }

        // ——— TIER 2c: Text classifier ———
        let textResult = textClassifier.classify(ocrText)
        logger.info("Text classifier: \(textResult.category.rawValue) @ \(String(format: "%.2f", textResult.confidence))")

        // ——— Composite score ———
        let scaScore: Float = scaIsSensitive ? 0.85 : 0.0
        var compositeCategory   = textResult.category
        var compositeConfidence: Float

        if textResult.category != .clean {
            compositeConfidence = max(scaScore, textResult.confidence)
        } else if scaIsSensitive {
            compositeCategory   = .explicitSexual
            compositeConfidence = scaScore
        } else {
            compositeConfidence = 0.0
        }

        let postMemMB = availableMemoryMB()
        logger.info("Composite: \(compositeCategory.rawValue) @ \(String(format: "%.2f", compositeConfidence)) — post: \(postMemMB) MB available (Δ\(preMemMB - postMemMB) MB used)")

        // ——— Route by confidence ———
        if compositeConfidence >= 0.3 && compositeCategory != .clean {
            var finalCategory   = compositeCategory
            var finalConfidence = compositeConfidence

            if compositeConfidence >= 0.5 && compositeConfidence < 0.7 {
                logger.info("Tier 3 /classify — confidence \(String(format: "%.2f", compositeConfidence))")
                if let cloud = await sendToClassify(text: ocrText, apiBase: apiBase, defaults: defaults) {
                    finalCategory   = cloud.category
                    finalConfidence = cloud.confidence
                    logger.info("Tier 3 result: \(cloud.category.rawValue) @ \(String(format: "%.2f", cloud.confidence))")
                }
            }

            if finalCategory != .clean && finalConfidence >= 0.3 {
                let tier = compositeConfidence >= 0.7 ? 2 : (finalConfidence >= 0.7 ? 3 : 2)
                await handleEvent(makeEvent(category: finalCategory, tier: tier,
                                            confidence: finalConfidence,
                                            summary: buildSummary(category: finalCategory,
                                                                   confidence: finalConfidence)))
            }
        }
        // cgImage is released here as analyzeFrame returns — no cross-frame retention.
    }

    // MARK: - Event handling with dedup

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
        let body: [String: Any] = ["category": event.category, "severity": event.severity,
                                   "summary": event.summary, "timestamp": tsStr]
        guard let data = try? JSONSerialization.data(withJSONObject: body) else { return false }
        req.httpBody = data
        do {
            let (_, resp) = try await URLSession.shared.data(for: req)
            if let http = resp as? HTTPURLResponse, (200...299).contains(http.statusCode) {
                logger.info("Event uploaded: \(event.category) \(event.severity)")
                return true
            }
            logger.warning("Event upload HTTP \((resp as? HTTPURLResponse)?.statusCode ?? 0)")
            return false
        } catch {
            logger.warning("Event upload network error: \(error.localizedDescription)")
            return false
        }
    }

    private func postEvent(_ event: DetectedEvent) async {
        if !(await uploadEvent(event)) {
            logger.warning("Event upload failed — writing to retry queue")
            writeEventToRetryQueue(event)
        }
    }

    // MARK: - Retry queue (network failures only)

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
                logger.info("Retry loop tick — available: \(availableMemoryMB()) MB")
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
            if await uploadEvent(events[i]) {
                events[i].processed = true; changed = true
            } else { break }
        }
        if changed, let newData = try? JSONEncoder().encode(events) {
            defaults.set(newData, forKey: "pendingEvents")
        }
    }

    private func flushContinuedActivity() async {
        let expired = await dedupStore.drainExpired()
        for (cat, _, count) in expired {
            let category = ContentCategory(rawValue: cat) ?? .explicitSexual
            let summary  = continuedActivitySummary(category: category, count: count + 1)
            logger.info("Continued-activity: \(summary)")
            await postEvent(makeEvent(category: category, tier: 2, confidence: 0.85, summary: summary))
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

    // MARK: - Vision OCR (synchronous via perform — safe to wrap in autoreleasepool)

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
        DetectedEvent(id: UUID().uuidString, timestamp: Date().timeIntervalSince1970,
                      category: category.rawValue, severity: category.severity(for: confidence),
                      summary: summary, tier: tier, confidence: confidence, processed: false)
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
        logger.info("Heartbeat: screen=\(screen) available=\(availableMemoryMB()) MB")
        guard let apiBase = sharedDefaults?.string(forKey: "apiBaseURL"),
              let token   = sharedDefaults?.string(forKey: "authToken"),
              let url     = URL(string: apiBase + "/heartbeat") else {
            logger.error("Heartbeat skipped — missing apiBaseURL or authToken")
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
        guard let defaults = sharedDefaults,
              let data   = defaults.data(forKey: "hashBlocklist"),
              let hashes = try? JSONDecoder().decode([UInt64].self, from: data) else { return }
        hashBlocklist = Set(hashes)
        logger.info("Loaded \(hashes.count) perceptual hashes into blocklist")
    }
}
