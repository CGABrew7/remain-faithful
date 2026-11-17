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
        if isMonitoringEnabled { startMonitoringIfCurrent() }
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
        // Apple's built-in adult content filter for Safari and other apps that
        // honor Screen Time web rules. It follows the Blocking switch.
        store.webContent.blockedByFilter = .auto()
    }

    func clearShielding() {
        store.shield.applications          = nil
        store.shield.applicationCategories = nil
        store.webContent.blockedByFilter   = nil
    }

    // MARK: - Activity monitoring (DeviceActivity)
    //
    // The daily schedule wakes the monitor extension so it can re-apply shields.
    // It does not register usage thresholds and it does not upload usage.

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

        let schedule = dailySchedule()

        do {
            activityCenter.stopMonitoring([.daily])
            try activityCenter.startMonitoring(.daily, during: schedule, events: [:])
            isMonitoringEnabled = true
            defaults?.set(true, forKey: "deviceActivityMonitoringEnabled")
            defaults?.set(scheduleSignature(), forKey: "deviceActivityScheduleSignature")
            logger.info("DeviceActivity schedule started")
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

    /// Launch restarts the schedule when the signature changed.
    /// There is no usage total to preserve.
    func startMonitoringIfCurrent() {
        let signature = scheduleSignature()
        let running = activityCenter.activities.contains(.daily)
        let saved = defaults?.string(forKey: "deviceActivityScheduleSignature")
        if running && saved == signature {
            isMonitoringEnabled = true
            defaults?.set(true, forKey: "deviceActivityMonitoringEnabled")
            logger.info("DeviceActivity monitoring already active")
            return
        }
        startMonitoring()
    }

    private func scheduleSignature() -> String {
        let data = (try? JSONEncoder().encode(selection)) ?? Data()
        var hash: UInt64 = 1469598103934665603
        for byte in data {
            hash ^= UInt64(byte)
            hash &*= 1099511628211
        }
        return "session-shield-v1:\(hash)"
    }

    private func dailySchedule() -> DeviceActivitySchedule {
        DeviceActivitySchedule(
            intervalStart: DateComponents(hour: 0, minute: 0, second: 0),
            intervalEnd: DateComponents(hour: 23, minute: 59, second: 59),
            repeats: true
        )
    }

    /// Keeps the monitor extension scheduled while Session Shield is on,
    /// without turning the Alerts switch on.
    func startShieldSchedule() {
        guard !selection.applicationTokens.isEmpty || !selection.categoryTokens.isEmpty else { return }
        if isMonitoringEnabled {
            startMonitoring()
            return
        }
        do {
            activityCenter.stopMonitoring([.daily])
            try activityCenter.startMonitoring(.daily, during: dailySchedule(), events: [:])
            logger.info("Session Shield schedule started")
        } catch {
            logger.error("Session Shield schedule failed: \(error.localizedDescription)")
        }
    }

    func stopMonitoring() {
        if !AppLockoutManager.shared.isEnabled {
            activityCenter.stopMonitoring([.daily])
            defaults?.removeObject(forKey: "deviceActivityScheduleSignature")
        }
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
        defaults?.removeObject(forKey: "deviceActivityScheduleSignature")
        defaults?.set(false, forKey: "screenTimeShieldingEnabled")
        clearShielding()
        selection = FamilyActivitySelection()
        defaults?.removeObject(forKey: "screenTimeSelection")
        defaults?.removeObject(forKey: "pendingRawEvents")
        defaults?.removeObject(forKey: "pendingEvents")
        defaults?.removeObject(forKey: "lastAppUsageTimestamp")
    }

    /// Drops any usage queue left by an older build. Those rows are not uploaded.
    func drainPendingRawEvents() async {
        defaults?.removeObject(forKey: "pendingRawEvents")
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
