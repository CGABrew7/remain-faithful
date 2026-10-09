import XCTest
@testable import RemainFaithful

final class HonestyAndTourTests: XCTestCase {

    func testTourDoesNotAutoPresentInDemoUnlessForced() {
        XCTAssertFalse(TourContent.shouldAutoPresent(completedVersion: 0, isDemoMode: true, forceTour: false))
        XCTAssertTrue(TourContent.shouldAutoPresent(completedVersion: 0, isDemoMode: true, forceTour: true))
    }

    func testTourShowsOnceUntilTheStoredVersionCatchesUp() {
        XCTAssertTrue(TourContent.shouldAutoPresent(completedVersion: 0, isDemoMode: false, forceTour: false))
        XCTAssertFalse(TourContent.shouldAutoPresent(
            completedVersion: TourContent.version, isDemoMode: false, forceTour: false))
        XCTAssertTrue(TourContent.shouldAutoPresent(
            completedVersion: TourContent.version, isDemoMode: false, forceTour: true))
    }

    func testInGroupTourPointsAtMembersThenInvite() {
        let steps = TourContent.steps(inGroup: true)
        XCTAssertEqual(steps.count, 9)
        XCTAssertEqual(steps[4].anchor, .groupMembers)
        XCTAssertEqual(steps[5].anchor, .groupInvite)
        XCTAssertTrue(steps.contains { $0.title == "Invite someone" })
        XCTAssertEqual(steps.last?.anchor, .settingsFeedback)
    }

    func testNoGroupTourOmitsTheInviteStep() {
        let steps = TourContent.steps(inGroup: false)
        XCTAssertEqual(steps.count, 8)
        XCTAssertEqual(steps[4].anchor, .groupCreate)
        XCTAssertFalse(steps.contains { $0.title == "Invite someone" })
        XCTAssertFalse(steps.contains { $0.anchor == .groupInvite })
    }

    func testPauseNoticeMatchesLockout() {
        XCTAssertEqual(
            PauseNotice.text(lockoutEnabled: true),
            "Pausing will notify your accountability partners")
        XCTAssertEqual(
            PauseNotice.text(lockoutEnabled: false),
            "Pausing does not notify your partners")
    }

    func testFeedbackLeavesDeviceFieldsOutUnlessAsked() {
        let off = sample(includeDevice: false, summary: "  Button did nothing  ", details: " I tapped it. ")
        XCTAssertTrue(off.canSend)
        let fields = off.apiFields()
        XCTAssertEqual(fields["type"], "Problem")
        XCTAssertEqual(fields["summary"], "Button did nothing")
        XCTAssertEqual(fields["details"], "I tapped it.")
        XCTAssertNil(fields["app_version"])
        XCTAssertNil(fields["device_model"])
        XCTAssertNil(fields["os_version"])

        let mailto = off.mailtoURL()?.absoluteString ?? ""
        XCTAssertTrue(mailto.contains("remainfaithful.com"))
        XCTAssertFalse(mailto.contains("iPhone15"))
    }

    func testFeedbackIncludesDeviceFieldsWhenAsked() {
        let on = sample(includeDevice: true, summary: "A thought", details: "More detail")
        let fields = on.apiFields()
        XCTAssertEqual(fields["app_version"], "Remain Faithful 1.0 (8)")
        XCTAssertEqual(fields["device_model"], "iPhone15,2")
        XCTAssertEqual(fields["os_version"], "iOS 18.1")
        let mailto = on.mailtoURL()?.absoluteString ?? ""
        XCTAssertTrue(mailto.contains("iPhone15"))
    }

    func testFeedbackRejectsBlankText() {
        XCTAssertFalse(sample(includeDevice: false, summary: "   ", details: "").canSend)
        XCTAssertFalse(sample(includeDevice: false, summary: "Hi", details: "   ").canSend)
    }

    private func sample(includeDevice: Bool, summary: String, details: String) -> FeedbackSubmission {
        FeedbackSubmission(
            type: .problem,
            summary: summary,
            details: details,
            includeDevice: includeDevice,
            appVersion: "Remain Faithful 1.0 (8)",
            deviceModel: "iPhone15,2",
            osVersion: "iOS 18.1"
        )
    }
}
