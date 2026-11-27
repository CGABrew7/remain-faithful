import Foundation

/// On-device check-in rules. No network, no app names, no location.
enum CoachCopy {
    static let alertType = "coach_nudge"
    static let checkInTitle = "A good time to check in"
    static let checkInBody = "Take a moment with Remain Faithful. A short check-in can keep you steady."
    static let primaryButton = "Check in now"
    static let secondaryButton = "Not now"
    static let partnerTitle = "Ask your partner to reach out"
    static let partnerBody = "If a hard stretch is starting, Remain Faithful can ask your partner to send you a word of support. They will not see what was on your screen, and they will not see which app you opened."
    static let partnerButton = "Turn on partner check-ins"
    static let partnerSentence = "Your partner would welcome a word from you today."
    static let awayTitle = "I'm away"
    static let awayBody = "Turn this on while your usual routine is off. Remain Faithful will check in with you a little more carefully until you turn it off."
    static let notificationType = "COACH_CHECK_IN"
    static let writingHelper = "Warm the check-in words on this iPhone"

    static var userFacing: [String] {
        [
            checkInTitle, checkInBody, primaryButton, secondaryButton,
            partnerTitle, partnerBody, partnerButton, partnerSentence,
            awayTitle, awayBody, writingHelper,
        ]
    }
}

enum CoachKeys {
    static let enabled = "coachEnabled"
    static let partnerNudge = "coachPartnerNudgeEnabled"
    static let lateStart = "coachLateNightStart"
    static let lateEnd = "coachLateNightEnd"
    static let away = "imAway"
    static let lastPrompt = "coachLastLocalPromptAt"
    static let lastNudge = "coachLastPartnerNudgeAt"
    static let threshold = "coachThresholdReachedAt"
    static let blockedCount = "blockedOpenCountToday"
    static let blockedDay = "blockedOpenDay"
    static let blockedStamp = "blockedOpenCountedStamp"
    static let checkedInDay = "coachCheckedInDay"
    static let owner = "coachOwnerUserID"
    static let writing = "coachWritingHelper"

    static let all = [
        enabled, partnerNudge, lateStart, lateEnd, away, lastPrompt, lastNudge,
        threshold, blockedCount, blockedDay, blockedStamp, checkedInDay, owner, writing,
    ]
}

enum CoachReply: Equatable {
    /// The prompt may appear. Nothing is sent to a partner.
    case pending
    case notNow
    case checkIn
}

struct CoachClock: Equatable {
    var year: Int
    var month: Int
    var day: Int
    var hour: Int
    var minute: Int
    /// 1 is Sunday, matching Calendar's Gregorian weekday.
    var weekday: Int
    var epoch: TimeInterval
}

struct CoachSignals: Equatable {
    var enabled: Bool
    var partnerNudgeEnabled: Bool
    var imAway: Bool
    var lateStartHour: Int
    var lateEndHour: Int
    var missedCheckIn: Bool
    var shieldOpensToday: Int
    var thresholdFresh: Bool
    var recentSession: Bool
    var lastPromptAt: TimeInterval?
    var lastNudge: CoachClock?
}

struct CoachDecision: Equatable {
    var showPrompt: Bool
    var sendPartnerNudge: Bool
}

enum CoachEngine {
    static let quietSeconds: TimeInterval = 6 * 60 * 60
    static let shieldOpensForPrompt = 2
    static let defaultLateStart = 23
    static let defaultLateEnd = 5
    static let thresholdFreshSeconds: TimeInterval = 12 * 60 * 60
    static let sessionRecentSeconds: TimeInterval = 2 * 60 * 60

    static func clock(from date: Date, timeZone: TimeZone) -> CoachClock {
        var calendar = Calendar(identifier: .gregorian)
        calendar.timeZone = timeZone
        let parts = calendar.dateComponents([.year, .month, .day, .hour, .minute, .weekday], from: date)
        return CoachClock(
            year: parts.year ?? 0,
            month: parts.month ?? 1,
            day: parts.day ?? 1,
            hour: parts.hour ?? 0,
            minute: parts.minute ?? 0,
            weekday: parts.weekday ?? 1,
            epoch: date.timeIntervalSince1970
        )
    }

    static func evaluate(clock: CoachClock, signals: CoachSignals, reply: CoachReply) -> CoachDecision {
        let open = windowOpen(clock: clock, signals: signals)
        let quiet = inQuietPeriod(clock: clock, signals: signals)
        let show = open && !quiet && reply == .pending
        let send = open
            && reply == .checkIn
            && signals.partnerNudgeEnabled
            && !nudgedToday(clock: clock, signals: signals)
        return CoachDecision(showPrompt: show, sendPartnerNudge: send)
    }

    static func windowOpen(clock: CoachClock, signals: CoachSignals) -> Bool {
        guard signals.enabled else { return false }
        if lateNight(hour: clock.hour, start: signals.lateStartHour, end: signals.lateEndHour) { return true }
        if clock.weekday == 1 && clock.hour >= 17 && clock.hour < 23 { return true }
        if signals.missedCheckIn && clock.hour >= 18 { return true }
        if signals.shieldOpensToday >= shieldOpensForPrompt { return true }
        if signals.imAway && clock.hour >= 8 && clock.hour < 22 { return true }
        if signals.thresholdFresh { return true }
        if signals.recentSession { return true }
        return false
    }

    /// Start is inclusive and end is exclusive. Equal hours mean the window is off.
    static func lateNight(hour: Int, start: Int, end: Int) -> Bool {
        let startHour = normalizeHour(start)
        let endHour = normalizeHour(end)
        let now = normalizeHour(hour)
        if startHour == endHour { return false }
        if startHour < endHour {
            return now >= startHour && now < endHour
        }
        return now >= startHour || now < endHour
    }

    static func inQuietPeriod(clock: CoachClock, signals: CoachSignals) -> Bool {
        guard let last = signals.lastPromptAt else { return false }
        return clock.epoch - last < quietSeconds
    }

    static func nudgedToday(clock: CoachClock, signals: CoachSignals) -> Bool {
        guard let last = signals.lastNudge else { return false }
        return last.year == clock.year && last.month == clock.month && last.day == clock.day
    }

    static func hourLabel(_ hour: Int) -> String {
        let normalized = normalizeHour(hour)
        let suffix = normalized < 12 ? "AM" : "PM"
        let hour12 = normalized % 12 == 0 ? 12 : normalized % 12
        return "\(hour12) \(suffix)"
    }

    static func dayKey(_ clock: CoachClock) -> String {
        String(format: "%04d-%02d-%02d", clock.year, clock.month, clock.day)
    }

    private static func normalizeHour(_ hour: Int) -> Int {
        let mod = hour % 24
        return mod >= 0 ? mod : mod + 24
    }
}

/// App-group values for the coach. Wiped on sign-out and when the account changes.
enum CoachStore {
    private static var defaults: UserDefaults? {
        UserDefaults(suiteName: SessionShield.appGroupID)
    }

    static func wipe(_ defaults: UserDefaults) {
        for key in CoachKeys.all {
            defaults.removeObject(forKey: key)
        }
    }

    static func wipe() {
        guard let defaults else { return }
        wipe(defaults)
    }

    static func prepareForUser(_ userID: Int) {
        guard let defaults else { return }
        let owner = defaults.integer(forKey: CoachKeys.owner)
        if owner != 0 && owner != userID {
            wipe(defaults)
        }
        defaults.set(userID, forKey: CoachKeys.owner)
    }

    static func signals(now: Date = Date(), timeZone: TimeZone = .current) -> CoachSignals {
        let defaults = defaults
        let clock = CoachEngine.clock(from: now, timeZone: timeZone)
        let day = CoachEngine.dayKey(clock)
        let enabled = defaults?.object(forKey: CoachKeys.enabled) as? Bool ?? true
        let start = defaults?.object(forKey: CoachKeys.lateStart) as? Int ?? CoachEngine.defaultLateStart
        let end = defaults?.object(forKey: CoachKeys.lateEnd) as? Int ?? CoachEngine.defaultLateEnd
        let prompt = defaults?.double(forKey: CoachKeys.lastPrompt) ?? 0
        let nudge = defaults?.double(forKey: CoachKeys.lastNudge) ?? 0
        let threshold = defaults?.double(forKey: CoachKeys.threshold) ?? 0
        let opens = (defaults?.string(forKey: CoachKeys.blockedDay) == day)
            ? (defaults?.integer(forKey: CoachKeys.blockedCount) ?? 0)
            : 0
        let checked = defaults?.string(forKey: CoachKeys.checkedInDay) == day
        let stamp = UserDefaults(suiteName: SessionShield.appGroupID)?.double(forKey: SessionShield.stampKey) ?? 0
        let broadcasting = UserDefaults(suiteName: SessionShield.appGroupID)?.bool(forKey: SessionShield.broadcastingKey) ?? false
        let recent = broadcasting && stamp > 0 && now.timeIntervalSince1970 - stamp <= CoachEngine.sessionRecentSeconds
        return CoachSignals(
            enabled: enabled,
            partnerNudgeEnabled: defaults?.bool(forKey: CoachKeys.partnerNudge) ?? false,
            imAway: defaults?.bool(forKey: CoachKeys.away) ?? false,
            lateStartHour: start,
            lateEndHour: end,
            missedCheckIn: !checked,
            shieldOpensToday: opens,
            thresholdFresh: threshold > 0 && now.timeIntervalSince1970 - threshold <= CoachEngine.thresholdFreshSeconds,
            recentSession: recent,
            lastPromptAt: prompt > 0 ? prompt : nil,
            lastNudge: nudge > 0 ? CoachEngine.clock(from: Date(timeIntervalSince1970: nudge), timeZone: timeZone) : nil
        )
    }

    static func markPrompt(at date: Date = Date()) {
        defaults?.set(date.timeIntervalSince1970, forKey: CoachKeys.lastPrompt)
    }

    static func markNudge(at date: Date = Date()) {
        defaults?.set(date.timeIntervalSince1970, forKey: CoachKeys.lastNudge)
    }

    static func markCheckedIn(at date: Date = Date(), timeZone: TimeZone = .current) {
        let clock = CoachEngine.clock(from: date, timeZone: timeZone)
        defaults?.set(CoachEngine.dayKey(clock), forKey: CoachKeys.checkedInDay)
    }

    static func setEnabled(_ on: Bool) { defaults?.set(on, forKey: CoachKeys.enabled) }
    static func setPartnerNudge(_ on: Bool) { defaults?.set(on, forKey: CoachKeys.partnerNudge) }
    static func setAway(_ on: Bool) { defaults?.set(on, forKey: CoachKeys.away) }
    static func setLateStart(_ hour: Int) { defaults?.set(hour, forKey: CoachKeys.lateStart) }
    static func setLateEnd(_ hour: Int) { defaults?.set(hour, forKey: CoachKeys.lateEnd) }
    static func setWritingHelper(_ on: Bool) { defaults?.set(on, forKey: CoachKeys.writing) }

    static var writingHelperOn: Bool {
        defaults?.object(forKey: CoachKeys.writing) as? Bool ?? true
    }

    /// One count per shield-open stamp, for the local day only.
    static func noteShieldOpen(stamp: TimeInterval, now: Date = Date(), timeZone: TimeZone = .current) {
        guard let defaults, stamp > 0 else { return }
        if defaults.double(forKey: CoachKeys.blockedStamp) == stamp { return }
        let clock = CoachEngine.clock(from: now, timeZone: timeZone)
        let day = CoachEngine.dayKey(clock)
        let count = defaults.string(forKey: CoachKeys.blockedDay) == day
            ? defaults.integer(forKey: CoachKeys.blockedCount)
            : 0
        defaults.set(day, forKey: CoachKeys.blockedDay)
        defaults.set(count + 1, forKey: CoachKeys.blockedCount)
        defaults.set(stamp, forKey: CoachKeys.blockedStamp)
    }
}

enum AccountQueue {
    static func withoutUser(_ userID: Int, items: [ProtectionAlerts.Item]) -> [ProtectionAlerts.Item] {
        items.filter { $0.userID != userID }
    }
}
