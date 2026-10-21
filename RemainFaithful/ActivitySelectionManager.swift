import Foundation
import Combine
import FamilyControls
import ManagedSettings
import DeviceActivity
import os

private let logger = Logger(subsystem: "com.remainfaithful.app",
                            category: "ActivitySelectionManager")

// Typed names that must match DeviceActivityMonitorExtension.swift exactly.
extension DeviceActivityName {
    static let daily = DeviceActivityName("com.remainfaithful.daily")
}
extension DeviceActivityEvent.Name {
    static let appUsage      = DeviceActivityEvent.Name("com.remainfaithful.appusage")
    static let categoryUsage = DeviceActivityEvent.Name("com.remainfaithful.categoryusage")
}

// Manages the FamilyActivitySelection (apps and categories to block/monitor),
// the ManagedSettingsStore shield, and the DeviceActivityCenter schedule.
// Selection is persisted to the shared app group so extensions can read it.
final class ActivitySelectionManager: ObservableObject {
    static let shared = ActivitySelectionManager()

    private let store          = ManagedSettingsStore()
    private let activityCenter = DeviceActivityCenter()
    private let defaults       = UserDefaults(suiteName: "group.com.remainfaithful.app")

    @Published var selection         = FamilyActivitySelection()
    @Published var isShieldingEnabled:  Bool = false
    @Published var isMonitoringEnabled: Bool = false

    private init() {
        isShieldingEnabled  = defaults?.bool(forKey: "screenTimeShieldingEnabled")  ?? false
        isMonitoringEnabled = defaults?.bool(forKey: "deviceActivityMonitoringEnabled") ?? false

        if let data  = defaults?.data(forKey: "screenTimeSelection"),
           let saved = try? JSONDecoder().decode(FamilyActivitySelection.self, from: data) {
            selection = saved
        }

        if isShieldingEnabled  { applyShielding() }
        if isMonitoringEnabled { startMonitoring() }
    }

    // MARK: - Selection

    func selectionDidChange() {
        let nowEmpty = selection.applicationTokens.isEmpty && selection.categoryTokens.isEmpty
        let wasEmpty = persistedSelectionIsEmpty()
        if let data = try? JSONEncoder().encode(selection) {
            defaults?.set(data, forKey: "screenTimeSelection")
        }
        // Clearing every chosen app is a reduction, even if a toggle was already off.
        if !wasEmpty && nowEmpty {
            ProtectionAlerts.send(
                type: "selection_cleared",
                detail: "All chosen apps were removed."
            )
            if isMonitoringEnabled { stopMonitoring() }
            if isShieldingEnabled { clearShielding() }
            return
        }
        if isShieldingEnabled  { applyShielding() }
        // Restart monitoring with updated token set
        if isMonitoringEnabled { startMonitoring() }
    }

    private func persistedSelectionIsEmpty() -> Bool {
        guard let data = defaults?.data(forKey: "screenTimeSelection"),
              let saved = try? JSONDecoder().decode(FamilyActivitySelection.self, from: data) else {
            return true
        }
        return saved.applicationTokens.isEmpty && saved.categoryTokens.isEmpty
    }

    // MARK: - Shielding (ManagedSettings)

    func setShieldingEnabled(_ enabled: Bool) {
        isShieldingEnabled = enabled
        defaults?.set(enabled, forKey: "screenTimeShieldingEnabled")
        if enabled {
            applyShielding()
        } else {
            clearShielding()
            ProtectionAlerts.send(
                type: "shielding_disabled",
                detail: "Blocking was turned off."
            )
        }
    }

    func applyShielding() {
        let apps = selection.applicationTokens
        let cats = selection.categoryTokens
        store.shield.applications        = apps.isEmpty ? nil : apps
        store.shield.applicationCategories = cats.isEmpty ? nil : .specific(cats)
    }

    func clearShielding() {
        store.shield.applications          = nil
        store.shield.applicationCategories = nil
    }

    // MARK: - Activity monitoring (DeviceActivity)
    //
    // DeviceActivityCenter fires DeviceActivityMonitorExtension.eventDidReachThreshold
    // when cumulative usage of the monitored set crosses 1 minute within the daily
    // schedule window (00:00–23:59, repeating). The minimum threshold enforced by
    // the API is 1 minute — sub-minute values are rejected at startMonitoring time.

    func setMonitoringEnabled(_ enabled: Bool) {
        if enabled {
            startMonitoring()
        } else {
            let wasOn = isMonitoringEnabled
            stopMonitoring()
            if wasOn {
                ProtectionAlerts.send(
                    type: "monitoring_disabled",
                    detail: "Alerts were turned off."
                )
            }
        }
    }

    func startMonitoring() {
        guard !selection.applicationTokens.isEmpty || !selection.categoryTokens.isEmpty else {
            logger.info("startMonitoring: no apps/categories selected — skipping")
            let wasOn = isMonitoringEnabled
            isMonitoringEnabled = false
            defaults?.set(false, forKey: "deviceActivityMonitoringEnabled")
            if wasOn {
                ProtectionAlerts.send(
                    type: "monitoring_disabled",
                    detail: "Alerts stopped because no apps are selected."
                )
            }
            return
        }

        // Always-on daily window; threshold resets at midnight each day.
        let schedule = DeviceActivitySchedule(
            intervalStart: DateComponents(hour: 0, minute: 0, second: 0),
            intervalEnd:   DateComponents(hour: 23, minute: 59, second: 59),
            repeats: true
        )

        // Build event dictionary — separate events for apps vs. categories
        // so the extension can distinguish them in its summarize() method.
        var events: [DeviceActivityEvent.Name: DeviceActivityEvent] = [:]

        if !selection.applicationTokens.isEmpty {
            events[.appUsage] = DeviceActivityEvent(
                applications: selection.applicationTokens,
                threshold: DateComponents(minute: 1)
            )
        }
        if !selection.categoryTokens.isEmpty {
            events[.categoryUsage] = DeviceActivityEvent(
                categories: selection.categoryTokens,
                threshold: DateComponents(minute: 1)
            )
        }

        do {
            // stopMonitoring first to update the token set on re-starts
            activityCenter.stopMonitoring([.daily])
            try activityCenter.startMonitoring(.daily, during: schedule, events: events)
            isMonitoringEnabled = true
            defaults?.set(true, forKey: "deviceActivityMonitoringEnabled")
            logger.info("DeviceActivity monitoring started (\(events.count) event(s))")
        } catch {
            isMonitoringEnabled = false
            defaults?.set(false, forKey: "deviceActivityMonitoringEnabled")
            logger.error("DeviceActivity startMonitoring failed: \(error.localizedDescription)")
            ProtectionAlerts.send(
                type: "monitoring_disabled",
                detail: "Alerts were turned off."
            )
        }
    }

    func stopMonitoring() {
        activityCenter.stopMonitoring([.daily])
        isMonitoringEnabled = false
        defaults?.set(false, forKey: "deviceActivityMonitoringEnabled")
        logger.info("DeviceActivity monitoring stopped")
    }

    /// Stops Alerts and Blocking and drops this account's queued events.
    /// Does not send a second protection notice; sign-out sends `signed_out`.
    func resetForSignOut() {
        activityCenter.stopMonitoring([.daily])
        isMonitoringEnabled = false
        isShieldingEnabled = false
        defaults?.set(false, forKey: "deviceActivityMonitoringEnabled")
        defaults?.set(false, forKey: "screenTimeShieldingEnabled")
        clearShielding()
        selection = FamilyActivitySelection()
        defaults?.removeObject(forKey: "screenTimeSelection")
        defaults?.removeObject(forKey: "pendingRawEvents")
        defaults?.removeObject(forKey: "pendingEvents")
        defaults?.removeObject(forKey: "lastAppUsageTimestamp")
    }

    // MARK: - Pending raw event drain
    //
    // The DeviceActivityMonitor extension writes to "pendingRawEvents" in the shared
    // app group when it fires. EventProcessor calls this on every app foreground to
    // upload any queued events that the extension couldn't post directly.
    func drainPendingRawEvents() async {
        guard let queue = defaults?.array(forKey: "pendingRawEvents") as? [[String: Any]],
              !queue.isEmpty else { return }

        guard let apiBase = defaults?.string(forKey: "apiBaseURL"),
              let token   = defaults?.string(forKey: "authToken") else { return }

        let iso = ISO8601DateFormatter()
        iso.formatOptions = [.withInternetDateTime, .withFractionalSeconds]

        let signedIn = defaults?.integer(forKey: "authUserID") ?? 0
        var updated: [[String: Any]] = []
        var stop = false
        for var row in queue {
            if (row["processed"] as? Bool) == true {
                updated.append(row)
                continue
            }
            let owner = (row["accountID"] as? Int) ?? (row["accountID"] as? NSNumber)?.intValue
            if !QueuedEventGate.shouldUpload(ownerAccountID: owner, signedInAccountID: signedIn) {
                continue
            }
            if stop {
                updated.append(row)
                continue
            }
            let ts  = row["timestamp"] as? TimeInterval ?? Date().timeIntervalSince1970
            let cat = row["category"]  as? String ?? "app_usage"
            let sev = row["severity"]  as? String ?? "concerning"
            let sum = row["summary"]   as? String ?? "App activity detected"

            guard let url  = URL(string: apiBase + "/events") else {
                updated.append(row)
                continue
            }
            var req = URLRequest(url: url, timeoutInterval: 10)
            req.httpMethod = "POST"
            req.setValue("application/json", forHTTPHeaderField: "Content-Type")
            req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
            guard let body = try? JSONSerialization.data(withJSONObject: [
                "category": cat, "severity": sev, "summary": sum,
                "timestamp": iso.string(from: Date(timeIntervalSince1970: ts)),
            ]) else {
                updated.append(row)
                continue
            }
            req.httpBody = body

            if let (_, resp) = try? await URLSession.shared.data(for: req),
               let http = resp as? HTTPURLResponse, (200...299).contains(http.statusCode) {
                row["processed"] = true
                updated.append(row)
                logger.info("Drained raw event: HTTP \(http.statusCode)")
            } else {
                updated.append(row)
                stop = true
            }
        }
        defaults?.set(updated, forKey: "pendingRawEvents")
    }

    // MARK: - App usage prompt

    // Returns the timestamp of the most recent DeviceActivity threshold crossing
    // if it occurred within the last 4 hours, nil otherwise.
    var recentAppUsageDate: Date? {
        let ts = defaults?.double(forKey: "lastAppUsageTimestamp") ?? 0
        guard ts > 0 else { return nil }
        let date = Date(timeIntervalSince1970: ts)
        return Date().timeIntervalSince(date) < 4 * 3600 ? date : nil
    }

    // Called after the user acts on the prompt (dismisses or taps Start Monitoring)
    // so it doesn't reappear until the next DeviceActivity threshold crossing.
    func clearAppUsagePrompt() {
        defaults?.removeObject(forKey: "lastAppUsageTimestamp")
    }

    // MARK: - Selection summary

    var selectionSummary: String {
        let apps = selection.applicationTokens.count
        let cats = selection.categoryTokens.count
        guard apps > 0 || cats > 0 else { return "No apps selected" }
        var parts: [String] = []
        if apps > 0 { parts.append("\(apps) app\(apps == 1 ? "" : "s")") }
        if cats > 0 { parts.append("\(cats) categor\(cats == 1 ? "y" : "ies")") }
        return parts.joined(separator: ", ")
    }
}
