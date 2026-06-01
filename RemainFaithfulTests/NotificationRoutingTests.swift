import XCTest
import Combine
@testable import RemainFaithful

final class NotificationRoutingTests: XCTestCase {

    var cancellables = Set<AnyCancellable>()

    override func setUp() {
        super.setUp()
        cancellables = []
        // Reset deepLink on the main thread before each test
        let cleared = expectation(description: "deepLink cleared")
        DispatchQueue.main.async {
            AppState.shared.deepLink = nil
            cleared.fulfill()
        }
        wait(for: [cleared], timeout: 1)
    }

    override func tearDown() {
        cancellables = []
        super.tearDown()
    }

    // MARK: - CONTENT_FLAGGED

    func testContentFlagged_navigatesToAlertDetail() {
        let info: [AnyHashable: Any] = [
            "notification_type": "CONTENT_FLAGGED",
            "category": "adult_content",
            "severity": "High",
            "summary": "Explicit content detected",
        ]

        let exp = expectation(description: "deepLink set to alertDetail")
        AppState.shared.$deepLink
            .compactMap { $0 }
            .first()
            .sink { link in
                if case .alertDetail = link { exp.fulfill() }
            }
            .store(in: &cancellables)

        NotificationService.shared.routeNotification(userInfo: info)
        wait(for: [exp], timeout: 1)

        guard case .alertDetail = AppState.shared.deepLink else {
            XCTFail("Expected .alertDetail, got \(String(describing: AppState.shared.deepLink))")
            return
        }
    }

    func testContentFlagged_withUnparsablePayload_fallsBackToDashboard() {
        // Missing "category" key — ActivityEvent.from returns nil, should route to .dashboard
        let info: [AnyHashable: Any] = [
            "notification_type": "CONTENT_FLAGGED",
        ]

        let exp = expectation(description: "deepLink set to dashboard")
        AppState.shared.$deepLink
            .compactMap { $0 }
            .first()
            .sink { link in
                if case .dashboard = link { exp.fulfill() }
            }
            .store(in: &cancellables)

        NotificationService.shared.routeNotification(userInfo: info)
        wait(for: [exp], timeout: 1)
    }

    // MARK: - PANIC_ALERT

    func testPanicAlert_navigatesToPanicView() {
        let info: [AnyHashable: Any] = ["notification_type": "PANIC_ALERT"]

        let exp = expectation(description: "deepLink set to panicView")
        AppState.shared.$deepLink
            .compactMap { $0 }
            .first()
            .sink { link in
                if case .panicView = link { exp.fulfill() }
            }
            .store(in: &cancellables)

        NotificationService.shared.routeNotification(userInfo: info)
        wait(for: [exp], timeout: 1)
    }

    // MARK: - Unknown type

    func testUnknownType_navigatesToDashboard() {
        let info: [AnyHashable: Any] = ["notification_type": "SOME_UNKNOWN_TYPE_XYZ"]

        let exp = expectation(description: "deepLink set to dashboard")
        AppState.shared.$deepLink
            .compactMap { $0 }
            .first()
            .sink { link in
                if case .dashboard = link { exp.fulfill() }
            }
            .store(in: &cancellables)

        NotificationService.shared.routeNotification(userInfo: info)
        wait(for: [exp], timeout: 1)

        guard case .dashboard = AppState.shared.deepLink else {
            XCTFail("Expected .dashboard")
            return
        }
    }
}
