import Foundation

/// Shared rules for Session Shield and for the fields a Deep Scan notice may carry.
/// Extensions compile this file too. It stays free of UIKit and of any network call.
enum SessionShield {
    static let appGroupID = "group.com.remainfaithful.app"
    /// Existing storage key. The shield reads the same flag.
    static let enabledKey = "appLockoutEnabled"
    static let broadcastingKey = "isBroadcasting"
    /// Written by the broadcast extension while frames are arriving. Readers treat a missing or stale stamp as not running.
    static let stampKey = "broadcastStateStamp"
    static let freshSeconds: TimeInterval = 45

    static func broadcastIsRunning(flag: Bool, stamp: TimeInterval, now: TimeInterval) -> Bool {
        guard flag, stamp > 0 else { return false }
        return now - stamp <= freshSeconds
    }

    /// Session Shield lifts shields only while a fresh broadcast is confirmed.
    /// A missing flag, a false flag, or a stale stamp stays shielded.
    static func shouldClearShields(
        sessionShieldOn: Bool,
        broadcastingFlag: Bool,
        stamp: TimeInterval,
        now: TimeInterval
    ) -> Bool {
        guard sessionShieldOn else { return false }
        return broadcastIsRunning(flag: broadcastingFlag, stamp: stamp, now: now)
    }
}

enum ConfidenceBucket {
    /// low, medium, or high. The raw score stays on the phone.
    static func name(for confidence: Float) -> String {
        if confidence >= 0.75 { return "high" }
        if confidence >= 0.50 { return "medium" }
        return "low"
    }
}

enum BlockedOpenRequest {
    static let pendingKey = "pendingBlockedOpenAt"
    static let notificationType = "BLOCKED_APP_OPENED"
    static let summary = "Blocked app opened"

    static func record(now: TimeInterval = Date().timeIntervalSince1970) {
        UserDefaults(suiteName: SessionShield.appGroupID)?.set(now, forKey: pendingKey)
    }

    static func pendingDate(now: TimeInterval = Date().timeIntervalSince1970) -> Date? {
        let defaults = UserDefaults(suiteName: SessionShield.appGroupID)
        let stamp = defaults?.double(forKey: pendingKey) ?? 0
        guard stamp > 0 else { return nil }
        // Ignore a stale tap from a previous day.
        guard now - stamp < 6 * 60 * 60 else { return nil }
        return Date(timeIntervalSince1970: stamp)
    }

    static func clear() {
        UserDefaults(suiteName: SessionShield.appGroupID)?.removeObject(forKey: pendingKey)
    }
}

enum CheckInReminder {
    static let notificationType = "DAILY_CHECK_IN"
    static let enabledKey = "dailyCheckInReminderOn"
    static let identifier = "com.remainfaithful.dailyCheckIn"
}

enum DeepScanConsent {
    static let acceptedKey = "deepScanSharingConsentAt"

    static var hasAccepted: Bool {
        UserDefaults.standard.double(forKey: acceptedKey) > 0
    }

    static func accept() {
        UserDefaults.standard.set(Date().timeIntervalSince1970, forKey: acceptedKey)
    }
}

enum SharingRecipients {
    static func line(partners: [String], groups: [String]) -> String {
        let people = partners.map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }.filter { !$0.isEmpty }
        let rooms = groups.map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }.filter { !$0.isEmpty }
        if people.isEmpty && rooms.isEmpty {
            return "Your accountability partners and your covenant group"
        }
        var parts: [String] = []
        if !people.isEmpty { parts.append(people.joined(separator: ", ")) }
        if !rooms.isEmpty { parts.append(rooms.joined(separator: ", ")) }
        return parts.joined(separator: " and ")
    }
}
