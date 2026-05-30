import ReplayKit
import SensitiveContentAnalysis
import Vision
import CoreImage
import NaturalLanguage
import os

// MARK: - Shared event model (mirrors DetectedEvent in the main app target)

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

// MARK: - Event persistence (App Group shared container)

private func writeEvent(_ event: DetectedEvent, to defaults: UserDefaults?) {
    guard let defaults else { return }
    var events: [DetectedEvent] = []
    if let data = defaults.data(forKey: "pendingEvents"),
       let existing = try? JSONDecoder().decode([DetectedEvent].self, from: data) {
        events = existing
    }
    events.append(event)
    // Cap at 200 events; drop oldest unprocessed first
    if events.count > 200 { events = Array(events.suffix(200)) }
    if let data = try? JSONEncoder().encode(events) {
        defaults.set(data, forKey: "pendingEvents")
        defaults.set(event.severity, forKey: "lastEventSeverity")
        defaults.set(event.summary, forKey: "lastEventSummary")
        defaults.set(event.timestamp, forKey: "lastEventTimestamp")
    }
}

// MARK: - Tier 3 cloud fallback

private func sendToClassify(text: String, apiBase: String) async -> TextClassification? {
    guard !text.isEmpty else { return nil }
    guard let url = URL(string: apiBase + "/classify") else { return nil }
    var req = URLRequest(url: url, timeoutInterval: 10)
    req.httpMethod = "POST"
    req.setValue("application/json", forHTTPHeaderField: "Content-Type")
    guard let body = try? JSONSerialization.data(withJSONObject: ["text": String(text.prefix(500))]) else {
        return nil
    }
    req.httpBody = body
    guard let (data, _) = try? await URLSession.shared.data(for: req),
          let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
          let cat = json["category"] as? String,
          let conf = json["confidence"] as? Double,
          let category = ContentCategory(rawValue: cat) else { return nil }
    return TextClassification(category: category, confidence: Float(conf))
}

// MARK: - SampleHandler

private let logger = Logger(subsystem: "com.remainfaithful.app.broadcast",
                            category: "classification")

class SampleHandler: RPBroadcastSampleHandler {

    private let appGroupID = "group.com.remainfaithful.app"

    private var sharedDefaults: UserDefaults? { UserDefaults(suiteName: appGroupID) }

    private let scaAnalyzer = SCSensitivityAnalyzer()
    private let textClassifier = FallbackTextClassifier()
    private let ciContext = CIContext(options: [.useSoftwareRenderer: false])

    private var frameCount = 0
    private let sampleInterval = 30

    // Hash blocklist — populated from server-side list at broadcast start
    private var hashBlocklist: Set<UInt64> = []

    // MARK: - Lifecycle

    override func broadcastStarted(withSetupInfo setupInfo: [String: NSObject]?) {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
        sharedDefaults?.set(Date().timeIntervalSince1970, forKey: "broadcastStartTime")
        loadHashBlocklist()
        let policy = scaAnalyzer.analysisPolicy
        logger.info("Broadcast started — SCA policy: \(String(describing: policy))")
    }

    override func broadcastPaused() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
    }

    override func broadcastResumed() {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
    }

    override func broadcastFinished() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
        logger.info("Broadcast finished — frames processed: \(self.frameCount)")
    }

    // MARK: - Frame processing

    override func processSampleBuffer(_ sampleBuffer: CMSampleBuffer,
                                      with sampleBufferType: RPSampleBufferType) {
        guard sampleBufferType == .video,
              let pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else { return }

        frameCount += 1
        let ciImage = CIImage(cvPixelBuffer: pixelBuffer)

        // — Tier 1a: perceptual hash check (every frame, ~0.1 ms) —
        let hash = dHash(from: ciImage, context: ciContext)
        if hashBlocklist.contains(where: { hammingDistance($0, hash) < 10 }) {
            logger.warning("Tier 1 hash match on frame \(self.frameCount)")
            let event = makeEvent(
                category: .explicitSexual, tier: 1, confidence: 1.0,
                summary: "Explicit image detected (perceptual hash match)"
            )
            writeEvent(event, to: sharedDefaults)
        }

        // — Full Tier 2 processing on every Nth frame —
        guard frameCount % sampleInterval == 0 else { return }

        // Render CGImage once; reused by OCR, SCA, and hash update
        guard let cgImage = ciContext.createCGImage(ciImage, from: ciImage.extent) else { return }

        let capturedDefaults = sharedDefaults
        let capturedFrame    = frameCount
        let capturedAnalyzer = scaAnalyzer
        let capturedClassifier = textClassifier
        let apiBase = capturedDefaults?.string(forKey: "apiBaseURL") ?? "http://localhost:8080"

        Task.detached(priority: .utility) { [self] in
            await self.analyzeFrame(
                cgImage: cgImage, ciImage: ciImage,
                frameNumber: capturedFrame,
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
        frameNumber: Int,
        scaAnalyzer: SCSensitivityAnalyzer,
        textClassifier: FallbackTextClassifier,
        defaults: UserDefaults?,
        apiBase: String
    ) async {
        // ——— TIER 2a: Vision OCR ———
        let ocrText = await extractText(from: cgImage)
        logger.info("OCR frame \(frameNumber): \(ocrText.count) chars — \(String(ocrText.prefix(120)).debugDescription)")

        defaults?.set(ocrText.prefix(500).description, forKey: "lastOCRText")
        defaults?.set(Date().timeIntervalSince1970, forKey: "lastOCRTimestamp")

        // ——— TIER 1b: URL + keyword check on extracted text ———
        var tier1Triggered = false
        let domains = extractDomains(from: ocrText)
        for domain in domains where Tier1Rules.matchesURL(domain) {
            logger.warning("Tier 1 URL match: \(domain) on frame \(frameNumber)")
            let event = makeEvent(
                category: .explicitSexual, tier: 1, confidence: 1.0,
                summary: "Blocked domain detected: \(domain)"
            )
            writeEvent(event, to: defaults)
            tier1Triggered = true
        }
        if !tier1Triggered && Tier1Rules.matchesKeyword(in: ocrText) {
            logger.warning("Tier 1 keyword match on frame \(frameNumber): \(String(ocrText.prefix(80)))")
            let event = makeEvent(
                category: .explicitSexual, tier: 1, confidence: 0.95,
                summary: "Explicit keyword detected in screen text"
            )
            writeEvent(event, to: defaults)
            tier1Triggered = true
        }

        // ——— TIER 2b: SensitiveContentAnalysis (image) ———
        var scaIsSensitive = false
        if scaAnalyzer.analysisPolicy != .disabled {
            let tmpURL = FileManager.default.temporaryDirectory
                .appendingPathComponent("rf_sca_\(frameNumber).jpg")
            defer { try? FileManager.default.removeItem(at: tmpURL) }
            if let jpeg = cgImageToJPEG(cgImage) {
                try? jpeg.write(to: tmpURL, options: .atomic)
                if let analysis = try? await scaAnalyzer.analyzeImage(at: tmpURL) {
                    scaIsSensitive = analysis.isSensitive
                    logger.info("SCA frame \(frameNumber): sensitive=\(analysis.isSensitive)")
                }
            }
        } else {
            logger.debug("SCA frame \(frameNumber): policy=disabled, skipping")
        }

        // ——— TIER 2c: Text classifier ———
        let textResult = textClassifier.classify(ocrText)
        logger.info("Text classifier frame \(frameNumber): \(textResult.category.rawValue) @ \(String(format: "%.2f", textResult.confidence))")

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

        logger.info("Composite frame \(frameNumber): \(compositeCategory.rawValue) @ \(String(format: "%.2f", compositeConfidence))")

        // ——— Route by confidence ———
        if compositeConfidence >= 0.3 && compositeCategory != .clean {
            // ——— TIER 3: cloud fallback for ambiguous zone ———
            var finalCategory   = compositeCategory
            var finalConfidence = compositeConfidence

            if compositeConfidence >= 0.5 && compositeConfidence < 0.7 {
                logger.info("Sending to Tier 3 /classify — confidence \(String(format: "%.2f", compositeConfidence))")
                if let cloud = await sendToClassify(text: ocrText, apiBase: apiBase) {
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
                    summary: buildSummary(category: finalCategory,
                                         confidence: finalConfidence,
                                         textSnippet: ocrText)
                )
                writeEvent(event, to: defaults)
            }
        }
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

    private func buildSummary(category: ContentCategory, confidence: Float,
                              textSnippet: String) -> String {
        let pct = Int(confidence * 100)
        let snippet = textSnippet.prefix(60).trimmingCharacters(in: .whitespacesAndNewlines)
        switch category {
        case .explicitSexual:
            return "Explicit content detected (\(pct)% confidence)\(snippet.isEmpty ? "" : ": \"\(snippet)\"")"
        case .gambling:
            return "Gambling content detected (\(pct)% confidence)\(snippet.isEmpty ? "" : ": \"\(snippet)\"")"
        case .violence:
            return "Violent content detected (\(pct)% confidence)\(snippet.isEmpty ? "" : ": \"\(snippet)\"")"
        case .selfHarm:
            return "Self-harm content detected (\(pct)% confidence) — support resources available"
        case .clean:
            return "Content reviewed — no concerns"
        }
    }

    private func cgImageToJPEG(_ cgImage: CGImage) -> Data? {
        UIImage(cgImage: cgImage).jpegData(compressionQuality: 0.70)
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
