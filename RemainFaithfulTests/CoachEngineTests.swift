import XCTest
@testable import RemainFaithful

final class CoachEngineTests: XCTestCase {

    private func clock(hour: Int, weekday: Int = 4, day: Int = 10) -> CoachClock {
        CoachClock(year: 2026, month: 11, day: day, hour: hour, minute: 0, weekday: weekday, epoch: 1_700_000_000 + Double(hour * 3600))
    }

    private func signals(
        enabled: Bool = true,
        partner: Bool = false,
        away: Bool = false,
        missed: Bool = false,
        opens: Int = 0,
        threshold: Bool = false,
        session: Bool = false,
        lastPrompt: TimeInterval? = nil,
        lastNudge: CoachClock? = nil,
        lateStart: Int = 23,
        lateEnd: Int = 5
    ) -> CoachSignals {
        CoachSignals(
            enabled: enabled,
            partnerNudgeEnabled: partner,
            imAway: away,
            lateStartHour: lateStart,
            lateEndHour: lateEnd,
            missedCheckIn: missed,
            shieldOpensToday: opens,
            thresholdFresh: threshold,
            recentSession: session,
            lastPromptAt: lastPrompt,
            lastNudge: lastNudge
        )
    }

    func testLateNightWindowWrapsPastMidnight() {
        let night = signals()
        XCTAssertTrue(CoachEngine.lateNight(hour: 23, start: 23, end: 5))
        XCTAssertTrue(CoachEngine.lateNight(hour: 0, start: 23, end: 5))
        XCTAssertTrue(CoachEngine.lateNight(hour: 4, start: 23, end: 5))
        XCTAssertFalse(CoachEngine.lateNight(hour: 5, start: 23, end: 5))
        XCTAssertFalse(CoachEngine.lateNight(hour: 12, start: 23, end: 5))
        XCTAssertFalse(CoachEngine.lateNight(hour: 22, start: 23, end: 5))
        let shown = CoachEngine.evaluate(clock: clock(hour: 23), signals: night, reply: .pending)
        XCTAssertTrue(shown.showPrompt)
        XCTAssertFalse(shown.sendPartnerNudge)
        XCTAssertFalse(CoachEngine.evaluate(clock: clock(hour: 12), signals: night, reply: .pending).showPrompt)
    }

    func testEqualLateHoursTurnTheWindowOff() {
        XCTAssertFalse(CoachEngine.lateNight(hour: 3, start: 3, end: 3))
    }

    func testSundayEveningOpensAndMondayDoesNot() {
        let base = signals()
        let sunday = CoachEngine.evaluate(clock: clock(hour: 18, weekday: 1), signals: base, reply: .pending)
        XCTAssertTrue(sunday.showPrompt)
        let monday = CoachEngine.evaluate(clock: clock(hour: 18, weekday: 2), signals: base, reply: .pending)
        XCTAssertFalse(monday.showPrompt)
    }

    func testMissedCheckInAfterEvening() {
        let missed = signals(missed: true)
        XCTAssertFalse(CoachEngine.evaluate(clock: clock(hour: 17), signals: missed, reply: .pending).showPrompt)
        XCTAssertTrue(CoachEngine.evaluate(clock: clock(hour: 18), signals: missed, reply: .pending).showPrompt)
    }

    func testAwayOpensTheDaytimeWindow() {
        let away = signals(away: true, lateStart: 0, lateEnd: 0)
        XCTAssertTrue(CoachEngine.evaluate(clock: clock(hour: 10), signals: away, reply: .pending).showPrompt)
        XCTAssertFalse(CoachEngine.evaluate(clock: clock(hour: 3), signals: away, reply: .pending).showPrompt)
    }

    func testShieldOpensAndThresholdAndSession() {
        XCTAssertFalse(CoachEngine.evaluate(clock: clock(hour: 12), signals: signals(opens: 1), reply: .pending).showPrompt)
        XCTAssertTrue(CoachEngine.evaluate(clock: clock(hour: 12), signals: signals(opens: 2), reply: .pending).showPrompt)
        XCTAssertTrue(CoachEngine.evaluate(clock: clock(hour: 12), signals: signals(threshold: true), reply: .pending).showPrompt)
        XCTAssertTrue(CoachEngine.evaluate(clock: clock(hour: 12), signals: signals(session: true), reply: .pending).showPrompt)
    }

    func testQuietPeriodAndALaterClock() {
        let now = clock(hour: 23)
        let recent = signals(lastPrompt: now.epoch - 60)
        XCTAssertFalse(CoachEngine.evaluate(clock: now, signals: recent, reply: .pending).showPrompt)
        let later = CoachClock(
            year: now.year, month: now.month, day: now.day, hour: 23, minute: 0, weekday: now.weekday,
            epoch: (recent.lastPromptAt ?? 0) + CoachEngine.quietSeconds
        )
        XCTAssertTrue(CoachEngine.evaluate(clock: later, signals: recent, reply: .pending).showPrompt)
    }

    func testNotNowNeverNotifiesAndOptInIsRequired() {
        let open = signals(partner: true)
        let now = clock(hour: 23)
        let declined = CoachEngine.evaluate(clock: now, signals: open, reply: .notNow)
        XCTAssertFalse(declined.showPrompt)
        XCTAssertFalse(declined.sendPartnerNudge)
        let pending = CoachEngine.evaluate(clock: now, signals: open, reply: .pending)
        XCTAssertTrue(pending.showPrompt)
        XCTAssertFalse(pending.sendPartnerNudge)
        let optedOut = CoachEngine.evaluate(clock: now, signals: signals(partner: false), reply: .checkIn)
        XCTAssertFalse(optedOut.sendPartnerNudge)
        let optedIn = CoachEngine.evaluate(clock: now, signals: open, reply: .checkIn)
        XCTAssertTrue(optedIn.sendPartnerNudge)
    }

    func testDailyCapUsesTheClockDayOnly() {
        let today = clock(hour: 23)
        var nudged = signals(partner: true, lastNudge: today)
        XCTAssertFalse(CoachEngine.evaluate(clock: today, signals: nudged, reply: .checkIn).sendPartnerNudge)
        let tomorrow = clock(hour: 23, day: 11)
        XCTAssertTrue(CoachEngine.evaluate(clock: tomorrow, signals: nudged, reply: .checkIn).sendPartnerNudge)
        nudged.lastNudge = clock(hour: 1, day: 10)
        let evening = clock(hour: 22, weekday: 2)
        XCTAssertFalse(CoachEngine.windowOpen(clock: evening, signals: signals()))
        let late = clock(hour: 23)
        XCTAssertTrue(CoachEngine.windowOpen(clock: late, signals: signals()))
    }

    func testCoachOffStaysQuiet() {
        let decision = CoachEngine.evaluate(clock: clock(hour: 23), signals: signals(enabled: false, partner: true), reply: .checkIn)
        XCTAssertFalse(decision.showPrompt)
        XCTAssertFalse(decision.sendPartnerNudge)
    }

    func testPartnerSentenceIsFixedAndCopyHasNoDashes() {
        XCTAssertEqual(CoachCopy.partnerSentence, "Your partner would welcome a word from you today.")
        let banned = ["\u{2014}", "\u{2013}", "honest"]
        for line in CoachCopy.userFacing + FutureSelfCopy.userFacing {
            for mark in banned {
                XCTAssertFalse(line.lowercased().contains(mark), line)
            }
        }
    }

    func testWipeRemovesCoachKeys() {
        let defaults = UserDefaults(suiteName: "group.remainfaithful.coach.tests")!
        defaults.set(true, forKey: CoachKeys.away)
        defaults.set(true, forKey: CoachKeys.partnerNudge)
        defaults.set("queued", forKey: CoachKeys.owner)
        CoachStore.wipe(defaults)
        for key in CoachKeys.all {
            XCTAssertNil(defaults.object(forKey: key), key)
        }
    }

    func testQueueDoesNotCrossAccounts() {
        let items = [
            ProtectionAlerts.Item(userID: 7, type: CoachCopy.alertType, detail: CoachCopy.partnerSentence),
            ProtectionAlerts.Item(userID: 8, type: "signed_out", detail: "Signed out."),
        ]
        let kept = AccountQueue.withoutUser(7, items: items)
        XCTAssertEqual(kept.map(\.userID), [8])
        XCTAssertFalse(kept.contains { $0.type == CoachCopy.alertType })
    }

    func testHourLabelsStayPlain() {
        XCTAssertEqual(CoachEngine.hourLabel(23), "11 PM")
        XCTAssertEqual(CoachEngine.hourLabel(0), "12 AM")
        XCTAssertEqual(CoachEngine.hourLabel(5), "5 AM")
        XCTAssertFalse(CoachEngine.hourLabel(23).contains("-"))
    }
}
