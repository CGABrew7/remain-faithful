import XCTest

// MARK: - Classifier under test
// Copied from RemainFaithfulBroadcast/SampleHandler.swift (see FIX 13 comment).
// Keep in sync with the broadcast extension copy.

private enum TestContentCategory: String {
    case explicitSexual = "adult_content"
    case gambling       = "gambling"
    case violence       = "violence"
    case selfHarm       = "self_harm"
    case clean          = "clean"
}

private struct TestTextClassification {
    let category: TestContentCategory
    let confidence: Float
}

private final class TestFallbackTextClassifier {

    private let trainingData: [TestContentCategory: [(term: String, weight: Float)]] = [

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

    func classify(_ text: String) -> TestTextClassification {
        let lower = text.lowercased()
        var scores = [TestContentCategory: Float]()

        for (category, terms) in trainingData {
            var score: Float = 0
            for (term, weight) in terms where lower.contains(term) {
                score += weight
            }
            if score > 0 {
                let normalized = min(1.0, 0.35 + score * 0.15)
                scores[category] = normalized
            }
        }

        if let top = scores.max(by: { $0.value < $1.value }) {
            return TestTextClassification(category: top.key, confidence: top.value)
        }
        return TestTextClassification(category: .clean, confidence: 0.95)
    }
}

// MARK: - Tests

final class FallbackTextClassifierTests: XCTestCase {

    let classifier = TestFallbackTextClassifier()

    func testAdultContent_highConfidence() {
        let result = classifier.classify("pornhub.com free videos")
        XCTAssertEqual(result.category, .explicitSexual)
        XCTAssertGreaterThanOrEqual(result.confidence, 0.5)
    }

    func testGambling_highConfidence() {
        let result = classifier.classify("draftkings sportsbook parlay odds calculator")
        XCTAssertEqual(result.category, .gambling)
        XCTAssertGreaterThanOrEqual(result.confidence, 0.5)
    }

    func testSelfHarm_highConfidence() {
        let result = classifier.classify("how to kill myself method")
        XCTAssertEqual(result.category, .selfHarm)
        XCTAssertGreaterThanOrEqual(result.confidence, 0.5)
    }

    func testViolence_highConfidence() {
        let result = classifier.classify("bestgore execution video graphic death")
        XCTAssertEqual(result.category, .violence)
        XCTAssertGreaterThanOrEqual(result.confidence, 0.5)
    }

    func testCleanText_classifiesAsClean() {
        let result = classifier.classify("quarterly earnings report revenue growth shareholder value")
        XCTAssertEqual(result.category, .clean)
    }

    func testCleanText_confidenceAtLeast90Percent() {
        let result = classifier.classify("quarterly earnings report revenue growth shareholder value")
        XCTAssertGreaterThanOrEqual(result.confidence, 0.9)
    }
}
