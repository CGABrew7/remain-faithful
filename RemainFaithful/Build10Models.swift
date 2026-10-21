import Foundation

// MARK: - Monitoring status
//
// The home card is a composite of what is actually true. Deep Scan is a subtitle.
// "Monitoring Paused" is never used while Alerts or Blocking are on.

enum MonitoringKind: Equatable {
    case needsPermission
    case noApps
    case off
    case on
}

struct MonitoringPresentation: Equatable {
    let kind: MonitoringKind
    let title: String
    let detail: String
    let deepScanLine: String
    /// Setup actions only. When Alerts or Blocking are already on, Home has no Deep Scan button.
    let actionTitle: String?
    /// True when Alerts or Blocking is on. Deep Scan alone does not count.
    let coreIsOn: Bool

    static func make(
        authorizationApproved: Bool,
        selectedCount: Int,
        alertsOn: Bool,
        blockingOn: Bool,
        deepScanOn: Bool,
        extensionSignedOut: Bool
    ) -> MonitoringPresentation {
        let deepScan = deepScanLine(on: deepScanOn, extensionSignedOut: extensionSignedOut)

        if !authorizationApproved {
            return MonitoringPresentation(
                kind: .needsPermission,
                title: "Setup needed",
                detail: "Screen Time access is off. Alerts and Blocking stay off until you allow it.",
                deepScanLine: deepScan,
                actionTitle: "Allow Screen Time",
                coreIsOn: false
            )
        }
        if selectedCount == 0 {
            return MonitoringPresentation(
                kind: .noApps,
                title: "No apps selected",
                detail: "Nothing is being watched yet. Choose the apps you want help with. You can take your time.",
                deepScanLine: deepScan,
                actionTitle: "Choose apps",
                coreIsOn: false
            )
        }
        if !alertsOn && !blockingOn {
            return MonitoringPresentation(
                kind: .off,
                title: "Monitoring off",
                detail: "Apps are chosen. Alerts and Blocking are both off. Turn them on when you are ready.",
                deepScanLine: deepScan,
                actionTitle: "Turn monitoring on",
                coreIsOn: false
            )
        }

        var parts: [String] = []
        parts.append(alertsOn ? "Alerts on" : "Alerts off")
        parts.append(blockingOn ? "Blocking on" : "Blocking off")
        return MonitoringPresentation(
            kind: .on,
            title: "Monitoring on",
            detail: parts.joined(separator: ". ") + ". You are not doing this alone.",
            deepScanLine: deepScan,
            actionTitle: nil,
            coreIsOn: true
        )
    }

    private static func deepScanLine(on: Bool, extensionSignedOut: Bool) -> String {
        if extensionSignedOut {
            return "Deep Scan stopped. Sign in again on this phone."
        }
        if on {
            return "Deep Scan is on. Protected video such as Netflix is not visible."
        }
        return "Deep Scan is off. It is optional and separate from Alerts and Blocking."
    }
}

// MARK: - Deep Scan onboarding
//
// Explained once during setup. Home does not repeat a "How Deep Scan works" button.

enum DeepScanOnboarding {
    static let title = "Deep Scan"
    static let subtitle = "Optional. It does not replace Alerts or Blocking."
    static let summary = "Deep Scan looks at this phone only while you start it from Control Center. Alerts and Blocking stay as you set them."
    static let steps = [
        "Open Control Center. If you do not see Screen Recording, add it in Settings.",
        "Touch and hold Screen Recording, choose Remain Faithful, then start it.",
    ]
    static let limits = "Protected video such as Netflix is not visible. Stopping or pausing Deep Scan notifies your partners."
    static let later = "Set up later"
}

// MARK: - Streak week
//
// The server window is rolling: index 0 is six days ago, the last index is today.
// Days before the account existed are dimmed, never a red X.

enum StreakDayInput: Decodable, Equatable {
    case clean
    case flagged
    case beforeSignup

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let flag = try? container.decode(Bool.self) {
            self = flag ? .clean : .flagged
            return
        }
        let raw = (try? container.decode(String.self))?.lowercased() ?? ""
        switch raw {
        case "clean", "true", "1":
            self = .clean
        case "before_signup", "before-signup", "before_account", "not_applicable", "na", "none":
            self = .beforeSignup
        default:
            self = .flagged
        }
    }
}

enum StreakMark: Equatable {
    case clean
    case flagged
    case beforeSignup

    var accessibilityName: String {
        switch self {
        case .clean: return "clean"
        case .flagged: return "flagged"
        case .beforeSignup: return "before you signed up"
        }
    }
}

enum StreakWeek {
    /// One letter per day, oldest first, ending on `today`.
    static func letters(count: Int = 7, endingOn today: Date, calendar: Calendar = .current) -> [String] {
        let total = max(count, 1)
        let start = calendar.startOfDay(for: today)
        return (0..<total).map { index in
            let daysAgo = total - 1 - index
            let date = calendar.date(byAdding: .day, value: -daysAgo, to: start) ?? start
            let weekday = calendar.component(.weekday, from: date)
            let symbols = calendar.veryShortWeekdaySymbols
            let symbol = symbols.indices.contains(weekday - 1) ? symbols[weekday - 1] : "?"
            return String(symbol.prefix(1)).uppercased()
        }
    }

    static func dayName(index: Int, count: Int = 7, endingOn today: Date, calendar: Calendar = .current) -> String {
        let total = max(count, 1)
        let start = calendar.startOfDay(for: today)
        let daysAgo = total - 1 - index
        let date = calendar.date(byAdding: .day, value: -daysAgo, to: start) ?? start
        let formatter = DateFormatter()
        formatter.calendar = calendar
        // Pin the formatter to the calendar's zone and locale. Without this it
        // falls back to the device's, and a day built in one zone is named in
        // another — VoiceOver then reads the wrong weekday for a streak dot.
        formatter.timeZone = calendar.timeZone
        formatter.locale = calendar.locale ?? .current
        formatter.dateFormat = "EEEE"
        return formatter.string(from: date)
    }

    static func marks(
        raw: [StreakDayInput],
        accountCreated: Date?,
        today: Date,
        calendar: Calendar = .current
    ) -> [StreakMark] {
        let total = raw.count
        let start = calendar.startOfDay(for: today)
        let createdDay = accountCreated.map { calendar.startOfDay(for: $0) }
        return raw.enumerated().map { index, value in
            let daysAgo = total - 1 - index
            let day = calendar.date(byAdding: .day, value: -daysAgo, to: start) ?? start
            if let createdDay, day < createdDay {
                return .beforeSignup
            }
            switch value {
            case .clean: return .clean
            case .beforeSignup: return .beforeSignup
            case .flagged: return .flagged
            }
        }
    }
}

// MARK: - Invites

enum PartnerInviteResult: Equatable {
    case waitingToAccept
    case invitedToSignUp
    case accepted
    case unknown

    /// `connected` is the legacy response for a row that is still pending.
    /// Treat it as waiting until the server says `accepted`.
    static func from(status: String) -> PartnerInviteResult {
        switch status.lowercased().trimmingCharacters(in: .whitespaces) {
        case "pending", "connected", "invited_existing":
            return .waitingToAccept
        case "invited":
            return .invitedToSignUp
        case "accepted":
            return .accepted
        default:
            return .unknown
        }
    }

    var message: String {
        switch self {
        case .waitingToAccept:
            return "Invite sent. Waiting for them to accept. Alerts start after they accept."
        case .invitedToSignUp:
            return "Invite sent. They can accept after they create an account with this email. They don't need an account before you invite them."
        case .accepted:
            return "They're connected. Alerts can start."
        case .unknown:
            return "Invite sent. They'll need to accept before alerts start."
        }
    }
}

enum GroupInviteOutcome: Equatable {
    case added
    case alreadyMember
    case invited

    static func from(status: String) -> GroupInviteOutcome? {
        switch status.lowercased() {
        case "added": return .added
        case "already_member", "conflict": return .alreadyMember
        case "invited": return .invited
        default: return nil
        }
    }

    var message: String {
        switch self {
        case .added:
            return "They're in the group now."
        case .alreadyMember:
            return "They're already in this group."
        case .invited:
            return "We emailed them. They join after they accept. If they don't have an account yet, they can accept after they create one with that email."
        }
    }
}

// MARK: - Password rules

enum PasswordRules {
    static let minimumCount = 8

    static func isMet(_ password: String) -> Bool {
        password.count >= minimumCount
    }

    /// Live helper under the sign-up field. Empty still states the rule.
    static func helper(_ password: String) -> String {
        if password.isEmpty {
            return "At least 8 characters."
        }
        if isMet(password) {
            return "8 or more characters."
        }
        return "\(password.count) of 8 characters."
    }
}

// MARK: - Plain-language errors

enum UserFacingError {
    static func message(for error: Error) -> String {
        guard let api = error as? APIError else {
            return "Something went wrong. Please try again."
        }
        switch api {
        case .server(let text, let code):
            if let coded = MailboxNotice.message(code: code, text: text) {
                return coded
            }
            return friendlyServer(text)
        case .network:
            return "Couldn't reach Remain Faithful. Check your connection and try again."
        case .decoding:
            return "We got an unexpected reply. Try again in a moment."
        case .unauthenticated:
            return "Please sign in again."
        case .badURL:
            return "Something went wrong. Please try again."
        }
    }

    static func friendlyServer(_ text: String) -> String {
        let lower = text.lowercased()
        if lower.contains("email already registered") || lower.contains("already registered") {
            return "That email already has an account. Sign in."
        }
        if lower.contains("cannot add yourself") {
            return "You can't invite your own email."
        }
        if lower.contains("maximum of 12") || lower.contains("maximum of twelve") {
            return "Groups hold up to 12 people."
        }
        if lower.contains("relationship not found") {
            return "Couldn't remove them. Try again."
        }
        if lower == "not authorized to set pin on this relationship" {
            return "You can set a PIN after this partnership is accepted."
        }
        return text
    }
}

// MARK: - Partner row

enum PartnerRow {
    static func showsPin(isPartner: Bool, status: String) -> Bool {
        isPartner && status.lowercased() == "accepted"
    }

    /// A pending invite with no accept button while the mailbox is still unconfirmed.
    static func needsEmailConfirm(awaitingMe: Bool, inviteID: Int?, emailVerified: Bool?, listCode: String?) -> Bool {
        guard awaitingMe, inviteID == nil else { return false }
        if listCode == "email_required" { return false }
        if listCode == "email_not_verified" { return true }
        return emailVerified == false
    }

    /// Apple sign-in did not share an email, so the invite cannot be accepted yet.
    static func needsAddedEmail(listCode: String?) -> Bool {
        listCode == "email_required"
    }
}

enum MailboxNotice {
    static let confirmFallback = "Confirm your email to accept."
    static let addEmailFallback = "Add an email in your profile and confirm it."
    static let pinPendingFallback = "You can set a PIN after this partnership is accepted."
    static let resendLimitFallback = "Too many confirmation emails. Try again in an hour."
    static let resendSent = "Check your email for the confirmation link."
    static let alreadyConfirmed = "Your email is already confirmed."

    static func message(code: String?, text: String) -> String? {
        let trimmed = text.trimmingCharacters(in: .whitespacesAndNewlines)
        let usable = !trimmed.isEmpty && !trimmed.hasPrefix("HTTP ")
        switch code {
        case "email_not_verified":
            return usable ? trimmed : confirmFallback
        case "email_required":
            return usable ? trimmed : addEmailFallback
        case "relationship_not_accepted":
            return usable ? trimmed : pinPendingFallback
        case "rate_limited":
            return usable ? trimmed : resendLimitFallback
        default:
            return nil
        }
    }

    static func resendStatus(_ status: String?) -> String {
        if status == "already_verified" { return alreadyConfirmed }
        return resendSent
    }
}

/// True when Apple sign-in did not share an email address.
enum AppleMailbox {
    static func isStandIn(email: String, appleUserID: String) -> Bool {
        let sub = appleUserID.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !sub.isEmpty else { return false }
        let expected = sub.lowercased() + "@privaterelay.appleid.com"
        return email.trimmingCharacters(in: .whitespacesAndNewlines).lowercased() == expected
    }
}

enum PINSaveMessage {
    static func text(for error: Error) -> String {
        let message = UserFacingError.message(for: error).trimmingCharacters(in: .whitespacesAndNewlines)
        if message.isEmpty { return "Couldn't save PIN. Try again." }
        return message
    }
}

// MARK: - Queued events stay with the account that created them

enum QueuedEventGate {
    static func shouldUpload(ownerAccountID: Int?, signedInAccountID: Int) -> Bool {
        let owner = ownerAccountID ?? 0
        if owner > 0 && signedInAccountID > 0 && owner != signedInAccountID {
            return false
        }
        if owner > 0 && signedInAccountID == 0 {
            return false
        }
        return true
    }
}

// MARK: - Screen Time turned off while signed out

struct ScreenTimeRevocationPlan: Equatable {
    var send: Bool
    var pending: Bool
    var storeApproved: Bool
}

enum ScreenTimeRevocation {
    static func plan(wasApproved: Bool, nowApproved: Bool, authenticated: Bool, alreadyPending: Bool) -> ScreenTimeRevocationPlan {
        var pending = alreadyPending
        if wasApproved && !nowApproved {
            pending = true
        }
        if nowApproved {
            pending = false
        }
        let send = authenticated && pending && !nowApproved
        if send {
            pending = false
        }
        return ScreenTimeRevocationPlan(send: send, pending: pending, storeApproved: nowApproved)
    }
}

// MARK: - End of setup

struct SetupLine: Equatable {
    let label: String
    let isOn: Bool
    var text: String {
        if label == "Alerts" {
            return isOn ? "Alerts are on" : "Alerts are off"
        }
        return "\(label) is \(isOn ? "on" : "off")"
    }
}

enum SetupStatusCopy {
    static func lines(screenTime: Bool, alerts: Bool, blocking: Bool, deepScan: Bool, emailInvite: Bool, linkShared: Bool) -> [SetupLine] {
        [
            SetupLine(label: "Screen Time", isOn: screenTime),
            SetupLine(label: "Alerts", isOn: alerts),
            SetupLine(label: "Blocking", isOn: blocking),
            SetupLine(label: "Deep Scan", isOn: deepScan),
            SetupLine(label: "Email invite", isOn: emailInvite),
            SetupLine(label: "Invite link", isOn: linkShared),
        ]
    }
}

enum FinishSetupCopy {
    static let title = "Finish setup"
    static func detail(kind: MonitoringKind) -> String {
        switch kind {
        case .needsPermission:
            return "Screen Time is still off. Allow it so Alerts and Blocking can help."
        case .noApps:
            return "Choose the apps you want help with. It takes a moment."
        case .off:
            return "Alerts and Blocking are off. Turn them on when you are ready."
        case .on:
            return ""
        }
    }
}

enum HomeGreeting {
    static func line(hour: Int) -> String {
        if hour < 12 { return "Good morning" }
        if hour < 17 { return "Good afternoon" }
        return "Good evening"
    }
    static let welcome = "Glad you are here."
}

enum PhoneSignals {
    private static func key(userID: Int) -> String { "supportTapCount.\(userID)" }

    static func supportTaps(userID: Int) -> Int {
        guard userID > 0 else { return 0 }
        return UserDefaults.standard.integer(forKey: key(userID: userID))
    }

    static func recordSupportTap(userID: Int) {
        guard userID > 0 else { return }
        UserDefaults.standard.set(supportTaps(userID: userID) + 1, forKey: key(userID: userID))
    }
}

// MARK: - PIN check
//
// Contract: POST /protection/pin/verify returns
// {"success": false, "reason": "no_pin" | "wrong_pin"}
// or 429 {"success": false, "reason": "rate_limited"}.
// Older servers omit `reason`. A false with no reason stays "wrong PIN".

enum PINCheck: Equatable {
    case success
    case wrong
    case noPIN
    case rateLimited
    case unavailable

    var succeeded: Bool {
        if case .success = self { return true }
        return false
    }

    var message: String {
        switch self {
        case .success:
            return ""
        case .wrong:
            return "Incorrect PIN. Your partner is the one who knows it."
        case .noPIN:
            return "No PIN is set for this partnership."
        case .rateLimited:
            return "Too many tries. Wait 15 minutes."
        case .unavailable:
            return "Couldn't check the PIN. Check your connection and try again."
        }
    }

    static func from(success: Bool, reason: String?) -> PINCheck {
        if success { return .success }
        switch reason?.lowercased() {
        case "no_pin", "not_set", "none", "missing":
            return .noPIN
        case "rate_limited":
            return .rateLimited
        default:
            return .wrong
        }
    }
}

// MARK: - Account dates

enum AccountDate {
    static func parse(_ raw: String?) -> Date? {
        guard let raw, !raw.isEmpty else { return nil }
        let withFraction = ISO8601DateFormatter()
        withFraction.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = withFraction.date(from: raw) { return date }
        let plain = ISO8601DateFormatter()
        plain.formatOptions = [.withInternetDateTime]
        return plain.date(from: raw)
    }
}
