import XCTest
@testable import RemainFaithful

final class Build10Tests: XCTestCase {

    // MARK: - Monitoring status

    func testMonitoringOnNeverSaysPaused() {
        let on = MonitoringPresentation.make(
            authorizationApproved: true,
            selectedCount: 2,
            alertsOn: true,
            blockingOn: false,
            deepScanOn: false,
            extensionSignedOut: false
        )
        XCTAssertEqual(on.kind, .on)
        XCTAssertTrue(on.coreIsOn)
        XCTAssertEqual(on.title, "Monitoring on")
        XCTAssertFalse(on.title.localizedCaseInsensitiveContains("paused"))
        XCTAssertFalse(on.detail.localizedCaseInsensitiveContains("paused"))
        XCTAssertTrue(on.detail.contains("Alerts on"))
        XCTAssertTrue(on.detail.contains("Blocking off"))
        XCTAssertTrue(on.deepScanLine.contains("optional"))
    }

    func testNoAppsBeatsALeftoverToggle() {
        let card = MonitoringPresentation.make(
            authorizationApproved: true,
            selectedCount: 0,
            alertsOn: true,
            blockingOn: true,
            deepScanOn: true,
            extensionSignedOut: false
        )
        XCTAssertEqual(card.kind, .noApps)
        XCTAssertEqual(card.title, "No apps selected")
        XCTAssertFalse(card.coreIsOn)
    }

    func testDeepScanSignOutIsASubtitleNotTheHeadline() {
        let card = MonitoringPresentation.make(
            authorizationApproved: true,
            selectedCount: 1,
            alertsOn: true,
            blockingOn: true,
            deepScanOn: false,
            extensionSignedOut: true
        )
        XCTAssertEqual(card.title, "Monitoring on")
        XCTAssertTrue(card.deepScanLine.contains("sign in again"))
        XCTAssertTrue(card.deepScanLine.localizedCaseInsensitiveContains("Deep Scan"))
    }

    // MARK: - Categories

    func testServerCategoriesAreKept() {
        XCTAssertEqual(EventCategory(apiString: "app_usage"), .appUsage)
        XCTAssertEqual(EventCategory(apiString: "self_harm"), .selfHarm)
        XCTAssertEqual(EventCategory(apiString: "violence"), .violence)
        XCTAssertEqual(EventCategory(apiString: "gambling"), .gambling)
        XCTAssertEqual(EventCategory(apiString: "adult_content"), .adultContent)
        XCTAssertNotEqual(EventCategory(apiString: "gambling").rawValue, EventCategory.adultContent.rawValue)
        XCTAssertNotEqual(EventCategory(apiString: "violence").rawValue, EventCategory.adultContent.rawValue)
        XCTAssertNotEqual(EventCategory(apiString: "self_harm").rawValue, EventCategory.adultContent.rawValue)
    }

    func testUnknownCategoryIsShownInsteadOfDropped() {
        XCTAssertEqual(EventCategory(apiString: "not_a_real_category"), .unknown)
        XCTAssertEqual(EventCategory(apiString: "signed_out"), .tamper)
        XCTAssertEqual(EventCategory(apiString: "shielding_disabled"), .tamper)
        let event = ActivityEvent.from(notificationPayload: [
            "notification_type": "CONTENT_FLAGGED",
            "category": "app_usage",
            "summary": "Safari was used",
        ])
        XCTAssertEqual(event?.category, .appUsage)
        XCTAssertEqual(event?.description, "Safari was used")
    }

    func testContentFlagWithoutCategoryStillHasNothingToOpen() {
        let event = ActivityEvent.from(notificationPayload: [
            "notification_type": "CONTENT_FLAGGED",
        ])
        XCTAssertNil(event)
    }

    func testProtectionAlertOpensATamperRow() {
        let event = ActivityEvent.from(notificationPayload: [
            "notification_type": "PROTECTION_ALERT",
            "protection_type": "signed_out",
            "detail": "Signed out on this phone.",
        ])
        XCTAssertEqual(event?.category, .tamper)
        XCTAssertEqual(event?.description, "Signed out on this phone.")
    }

    // MARK: - Streak

    func testStreakLettersEndOnToday() {
        var calendar = Calendar(identifier: .gregorian)
        calendar.locale = Locale(identifier: "en_US_POSIX")
        calendar.timeZone = TimeZone(secondsFromGMT: 0)!
        let today = calendar.date(from: DateComponents(year: 2026, month: 10, day: 13))!
        let letters = StreakWeek.letters(count: 7, endingOn: today, calendar: calendar)
        XCTAssertEqual(letters.count, 7)
        XCTAssertEqual(letters.last, "T")
        XCTAssertEqual(
            StreakWeek.dayName(index: 6, count: 7, endingOn: today, calendar: calendar),
            "Tuesday"
        )
        XCTAssertEqual(letters, ["W", "T", "F", "S", "S", "M", "T"])
    }

    func testDaysBeforeSignupAreDimmed() {
        var calendar = Calendar(identifier: .gregorian)
        calendar.locale = Locale(identifier: "en_US_POSIX")
        calendar.timeZone = TimeZone(secondsFromGMT: 0)!
        let today = calendar.date(from: DateComponents(year: 2026, month: 10, day: 13))!
        let created = calendar.date(from: DateComponents(year: 2026, month: 10, day: 12))!
        let raw = Array(repeating: StreakDayInput.flagged, count: 7)
        let marks = StreakWeek.marks(raw: raw, accountCreated: created, today: today, calendar: calendar)
        XCTAssertEqual(marks.filter { $0 == .beforeSignup }.count, 5)
        XCTAssertEqual(marks[5], .flagged)
        XCTAssertEqual(marks[6], .flagged)
        XCTAssertNotEqual(marks[0], .flagged)
    }

    func testWeekDecodesBoolsAndBeforeSignup() throws {
        struct Box: Decodable { let week: [StreakDayInput] }
        let data = Data(#"{"week":[true,false,"before_signup"]}"#.utf8)
        let box = try JSONDecoder().decode(Box.self, from: data)
        XCTAssertEqual(box.week, [.clean, .flagged, .beforeSignup])
        let states = Data(#"{"week":["before_account","clean","flagged"]}"#.utf8)
        let decoded = try JSONDecoder().decode(Box.self, from: states)
        XCTAssertEqual(decoded.week, [.beforeSignup, .clean, .flagged])
    }

    // MARK: - Invites and password

    func testPartnerInviteStatuses() {
        XCTAssertEqual(PartnerInviteResult.from(status: "connected"), .waitingToAccept)
        XCTAssertEqual(PartnerInviteResult.from(status: "pending"), .waitingToAccept)
        XCTAssertEqual(PartnerInviteResult.from(status: "invited"), .invitedToSignUp)
        XCTAssertEqual(PartnerInviteResult.from(status: "accepted"), .accepted)
        XCTAssertTrue(PartnerInviteResult.from(status: "connected").message.contains("accept"))
    }

    func testGroupInviteOutcomes() {
        XCTAssertEqual(GroupInviteOutcome.from(status: "added"), .added)
        XCTAssertEqual(GroupInviteOutcome.from(status: "already_member"), .alreadyMember)
        XCTAssertEqual(GroupInviteOutcome.from(status: "invited"), .invited)
        XCTAssertTrue(GroupInviteOutcome.added.message.contains("in the group"))
        XCTAssertTrue(GroupInviteOutcome.alreadyMember.message.contains("already"))
        XCTAssertTrue(GroupInviteOutcome.invited.message.contains("emailed"))
    }

    func testPasswordHelperStatesTheRule() {
        XCTAssertEqual(PasswordRules.helper(""), "At least 8 characters.")
        XCTAssertEqual(PasswordRules.helper("abc"), "3 of 8 characters.")
        XCTAssertFalse(PasswordRules.isMet("abc"))
        XCTAssertTrue(PasswordRules.isMet("12345678"))
        XCTAssertEqual(PasswordRules.helper("12345678"), "8 or more characters.")
    }

    func testPlainErrors() {
        XCTAssertEqual(
            UserFacingError.friendlyServer("email already registered"),
            "That email already has an account. Sign in."
        )
        XCTAssertEqual(
            UserFacingError.friendlyServer("groups have a maximum of 12 members"),
            "Groups hold up to 12 people."
        )
        XCTAssertEqual(PINCheck.from(success: false, reason: "no_pin"), .noPIN)
        XCTAssertEqual(PINCheck.from(success: false, reason: "wrong_pin"), .wrong)
        XCTAssertEqual(PINCheck.from(success: false, reason: "rate_limited"), .rateLimited)
        XCTAssertTrue(PINCheck.rateLimited.message.contains("15"))
        XCTAssertEqual(PINCheck.from(success: false, reason: nil), .wrong)
        XCTAssertEqual(PINCheck.from(success: true, reason: nil), .success)
    }
}
