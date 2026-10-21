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
        XCTAssertNil(on.actionTitle)

        let scanning = MonitoringPresentation.make(
            authorizationApproved: true,
            selectedCount: 2,
            alertsOn: true,
            blockingOn: true,
            deepScanOn: true,
            extensionSignedOut: false
        )
        XCTAssertNil(scanning.actionTitle)
        XCTAssertEqual(
            MonitoringPresentation.make(
                authorizationApproved: false,
                selectedCount: 0,
                alertsOn: false,
                blockingOn: false,
                deepScanOn: false,
                extensionSignedOut: false
            ).actionTitle,
            "Allow Screen Time"
        )
    }

    func testDeepScanIsExplainedOnceInOnboardingCopy() {
        let copy = ([
            DeepScanOnboarding.title,
            DeepScanOnboarding.subtitle,
            DeepScanOnboarding.summary,
            DeepScanOnboarding.limits,
        ] + DeepScanOnboarding.steps).joined(separator: " ")
        XCTAssertFalse(copy.contains("—"))
        XCTAssertFalse(copy.contains("–"))
        XCTAssertFalse(copy.localizedCaseInsensitiveContains("How Deep Scan works"))
        XCTAssertTrue(copy.contains("Control Center"))
        XCTAssertTrue(copy.contains("Screen Recording"))
        XCTAssertTrue(copy.contains("Settings"))
        XCTAssertTrue(copy.contains("Netflix"))
        XCTAssertTrue(copy.contains("notifies your partners"))
        XCTAssertTrue(copy.contains("does not replace Alerts or Blocking"))
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
        XCTAssertEqual(EventCategory(apiString: "selection_cleared"), .tamper)
        XCTAssertEqual(EventCategory(apiString: "deep_scan_stopped"), .tamper)
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

    func testPinButtonWaitsForAnAcceptedPartnership() {
        XCTAssertFalse(PartnerRow.showsPin(isPartner: true, status: "pending"))
        XCTAssertFalse(PartnerRow.showsPin(isPartner: false, status: "accepted"))
        XCTAssertTrue(PartnerRow.showsPin(isPartner: true, status: "accepted"))
        XCTAssertEqual(
            PINSaveMessage.text(for: APIError.server("not authorized to set pin on this relationship", code: nil)),
            "You can set a PIN after this partnership is accepted."
        )
        XCTAssertEqual(
            PINSaveMessage.text(for: APIError.server(
                "Your partner request isn't accepted yet. Confirm your email, then accept the invite.",
                code: "relationship_not_accepted"
            )),
            "Your partner request isn't accepted yet. Confirm your email, then accept the invite."
        )
    }

    func testUnconfirmedEmailExplainsTheMissingAcceptButton() {
        XCTAssertTrue(PartnerRow.needsEmailConfirm(awaitingMe: true, inviteID: nil, emailVerified: false, listCode: nil))
        XCTAssertTrue(PartnerRow.needsEmailConfirm(awaitingMe: true, inviteID: nil, emailVerified: nil, listCode: "email_not_verified"))
        XCTAssertFalse(PartnerRow.needsEmailConfirm(awaitingMe: true, inviteID: 4, emailVerified: false, listCode: nil))
        XCTAssertFalse(PartnerRow.needsEmailConfirm(awaitingMe: true, inviteID: nil, emailVerified: nil, listCode: nil))
        XCTAssertFalse(PartnerRow.needsEmailConfirm(awaitingMe: true, inviteID: nil, emailVerified: false, listCode: "email_required"))
        XCTAssertTrue(PartnerRow.needsAddedEmail(listCode: "email_required"))
        XCTAssertEqual(
            UserFacingError.message(for: APIError.server("HTTP 403", code: "email_not_verified")),
            "Confirm your email to accept."
        )
        XCTAssertEqual(
            UserFacingError.message(for: APIError.server("HTTP 409", code: "email_required")),
            "Add an email in your profile and confirm it."
        )
        XCTAssertEqual(
            UserFacingError.message(for: APIError.server("Too many confirmation emails. Try again in an hour.", code: "rate_limited")),
            "Too many confirmation emails. Try again in an hour."
        )
        XCTAssertEqual(MailboxNotice.resendStatus("sent"), "Check your email for the confirmation link.")
        XCTAssertEqual(MailboxNotice.resendStatus("already_verified"), "Your email is already confirmed.")
        XCTAssertTrue(AppleMailbox.isStandIn(email: "abc.123@privaterelay.appleid.com", appleUserID: "abc.123"))
        XCTAssertFalse(AppleMailbox.isStandIn(email: "hide.me@privaterelay.appleid.com", appleUserID: "abc.123"))
        XCTAssertEqual(
            UserFacingError.message(for: APIError.server("Confirm your email to accept this invite.", code: "email_not_verified")),
            "Confirm your email to accept this invite."
        )
    }

    func testQueuedAlertsStayWithTheAccountThatMadeThem() {
        XCTAssertFalse(QueuedEventGate.shouldUpload(ownerAccountID: 4, signedInAccountID: 9))
        XCTAssertFalse(QueuedEventGate.shouldUpload(ownerAccountID: 4, signedInAccountID: 0))
        XCTAssertTrue(QueuedEventGate.shouldUpload(ownerAccountID: 4, signedInAccountID: 4))
        XCTAssertTrue(QueuedEventGate.shouldUpload(ownerAccountID: nil, signedInAccountID: 4))
    }

    func testScreenTimeRevocationWaitsForTheNextSignIn() {
        let signedOut = ScreenTimeRevocation.plan(wasApproved: true, nowApproved: false, authenticated: false, alreadyPending: false)
        XCTAssertFalse(signedOut.send)
        XCTAssertTrue(signedOut.pending)
        let signedIn = ScreenTimeRevocation.plan(wasApproved: false, nowApproved: false, authenticated: true, alreadyPending: true)
        XCTAssertTrue(signedIn.send)
        XCTAssertFalse(signedIn.pending)
        let restored = ScreenTimeRevocation.plan(wasApproved: false, nowApproved: true, authenticated: true, alreadyPending: true)
        XCTAssertFalse(restored.send)
        XCTAssertFalse(restored.pending)
    }

    func testSetupStatusAndHomeCopyStayPlain() {
        let lines = SetupStatusCopy.lines(screenTime: true, alerts: true, blocking: false, deepScan: false, emailInvite: true, linkShared: false)
        XCTAssertEqual(lines.map(\.text), [
            "Screen Time is on",
            "Alerts are on",
            "Blocking is off",
            "Deep Scan is off",
            "Email invite is on",
            "Invite link is off",
        ])
        XCTAssertEqual(DeepScanOnboarding.later, "Set up later")
        XCTAssertEqual(HomeGreeting.line(hour: 9), "Good morning")
        XCTAssertEqual(HomeGreeting.welcome, "Glad you are here.")
        XCTAssertEqual(FinishSetupCopy.title, "Finish setup")
        let copy = FinishSetupCopy.detail(kind: .needsPermission)
            + FinishSetupCopy.detail(kind: .noApps)
            + FinishSetupCopy.detail(kind: .off)
            + HomeGreeting.welcome
        XCTAssertFalse(copy.contains("—"))
        XCTAssertFalse(copy.contains("–"))
        let on = MonitoringPresentation.make(
            authorizationApproved: true,
            selectedCount: 1,
            alertsOn: true,
            blockingOn: false,
            deepScanOn: false,
            extensionSignedOut: false
        )
        XCTAssertTrue(on.detail.contains("You are not doing this alone."))
    }

    func testSupportTapsCountOnThisPhone() {
        let id = 880011
        UserDefaults.standard.removeObject(forKey: "supportTapCount.\(id)")
        defer { UserDefaults.standard.removeObject(forKey: "supportTapCount.\(id)") }
        XCTAssertEqual(PhoneSignals.supportTaps(userID: id), 0)
        PhoneSignals.recordSupportTap(userID: id)
        PhoneSignals.recordSupportTap(userID: id)
        XCTAssertEqual(PhoneSignals.supportTaps(userID: id), 2)
        XCTAssertEqual(PhoneSignals.supportTaps(userID: 0), 0)
    }
}
