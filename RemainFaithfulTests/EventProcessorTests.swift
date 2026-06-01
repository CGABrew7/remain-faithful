import XCTest
@testable import RemainFaithful

final class EventProcessorTests: XCTestCase {

    let testGroupID = "group.com.remainfaithful.test"
    var testDefaults: UserDefaults!

    override func setUp() {
        super.setUp()
        testDefaults = UserDefaults(suiteName: testGroupID)
        testDefaults.removeObject(forKey: "pendingEvents")
        testDefaults.removeObject(forKey: "apiBaseURL")
        testDefaults.removeObject(forKey: "classifySecret")
    }

    override func tearDown() {
        testDefaults.removePersistentDomain(forName: testGroupID)
        testDefaults = nil
        super.tearDown()
    }

    // MARK: - processPendingEvents

    func testProcessPendingEvents_doesNothingWhenNotAuthenticated() async {
        let sender = MockEventSender(authenticated: false)
        let sut = EventProcessor(appGroupID: testGroupID, eventSender: sender)
        writeEvents([makeEvent(processed: false)], to: testDefaults)

        await sut.processPendingEvents()

        XCTAssertEqual(sender.createCallCount, 0, "should not call API when unauthenticated")
    }

    func testProcessPendingEvents_marksEventsProcessedOnSuccess() async {
        let sender = MockEventSender(authenticated: true)
        let sut = EventProcessor(appGroupID: testGroupID, eventSender: sender)
        writeEvents([makeEvent(processed: false), makeEvent(processed: false)], to: testDefaults)

        await sut.processPendingEvents()

        XCTAssertEqual(sender.createCallCount, 2)
        let stored = readEvents(from: testDefaults)
        XCTAssertTrue(stored.allSatisfy(\.processed), "all events should be marked processed")
    }

    // MARK: - pendingEventCount

    func testPendingEventCount_zeroWhenNoData() {
        let sut = EventProcessor(appGroupID: testGroupID, eventSender: MockEventSender(authenticated: false))
        XCTAssertEqual(sut.pendingEventCount(), 0)
    }

    func testPendingEventCount_countsOnlyUnprocessed() {
        let sut = EventProcessor(appGroupID: testGroupID, eventSender: MockEventSender(authenticated: false))
        writeEvents([
            makeEvent(processed: false),
            makeEvent(processed: true),
            makeEvent(processed: false),
        ], to: testDefaults)

        XCTAssertEqual(sut.pendingEventCount(), 2)
    }

    // MARK: - Helpers

    private func makeEvent(processed: Bool) -> DetectedEvent {
        DetectedEvent(
            id: UUID().uuidString,
            timestamp: Date().timeIntervalSince1970,
            category: "adult_content",
            severity: "severe",
            summary: "test event",
            tier: 2,
            confidence: 0.85,
            processed: processed
        )
    }

    private func writeEvents(_ events: [DetectedEvent], to defaults: UserDefaults) {
        defaults.set(try! JSONEncoder().encode(events), forKey: "pendingEvents")
    }

    private func readEvents(from defaults: UserDefaults) -> [DetectedEvent] {
        guard let data = defaults.data(forKey: "pendingEvents") else { return [] }
        return (try? JSONDecoder().decode([DetectedEvent].self, from: data)) ?? []
    }
}

// MARK: - Mock

final class MockEventSender: EventSending {
    var isAuthenticated: Bool
    var createCallCount = 0
    var shouldThrow = false

    init(authenticated: Bool) { self.isAuthenticated = authenticated }

    func createEvent(category: String, severity: String, summary: String,
                     timestamp: String?) async throws -> RemoteEvent {
        if shouldThrow { throw URLError(.notConnectedToInternet) }
        createCallCount += 1
        return RemoteEvent(id: createCallCount, userId: 1,
                           category: category, severity: severity,
                           summary: summary, timestamp: "2025-01-01T00:00:00Z")
    }
}
