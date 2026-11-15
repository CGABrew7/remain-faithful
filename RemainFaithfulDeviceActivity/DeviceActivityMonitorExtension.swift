import DeviceActivity
import FamilyControls
import Foundation
import ManagedSettings
import os

private let logger = Logger(
    subsystem: "com.remainfaithful.app.DeviceActivityMonitor",
    category: "monitor"
)

// The system launches this process on the daily schedule.
// It applies or clears shields on this phone. It does not open a network connection.
class DeviceActivityMonitorExtension: DeviceActivityMonitor {

    private let appGroupID = "group.com.remainfaithful.app"
    private lazy var defaults = UserDefaults(suiteName: appGroupID)
    private let store = ManagedSettingsStore()

    override func intervalDidStart(for activity: DeviceActivityName) {
        super.intervalDidStart(for: activity)
        logger.info("Interval started: \(activity.rawValue)")
        syncShields()
    }

    override func intervalDidEnd(for activity: DeviceActivityName) {
        super.intervalDidEnd(for: activity)
        logger.info("Interval ended: \(activity.rawValue)")
        syncShields()
    }

    override func eventDidReachThreshold(
        _ event: DeviceActivityEvent.Name,
        activity: DeviceActivityName
    ) {
        super.eventDidReachThreshold(event, activity: activity)
        syncShields()
    }

    private func syncShields() {
        let now = Date().timeIntervalSince1970
        let session = defaults?.bool(forKey: SessionShield.enabledKey) ?? false
        let flag = defaults?.bool(forKey: SessionShield.broadcastingKey) ?? false
        let stamp = defaults?.double(forKey: SessionShield.stampKey) ?? 0
        let blocking = defaults?.bool(forKey: "screenTimeShieldingEnabled") ?? false

        if session {
            if SessionShield.shouldClearShields(
                sessionShieldOn: true,
                broadcastingFlag: flag,
                stamp: stamp,
                now: now
            ) {
                clearShields()
                logger.info("Session Shield: broadcast confirmed, shields lifted")
            } else {
                applyShields()
                logger.info("Session Shield: broadcast not confirmed, shields applied")
            }
            return
        }
        if blocking {
            applyShields()
        }
    }

    private func applyShields() {
        guard let selection = loadSelection() else { return }
        let apps = selection.applicationTokens
        let cats = selection.categoryTokens
        store.shield.applications = apps.isEmpty ? nil : apps
        store.shield.applicationCategories = cats.isEmpty ? nil : .specific(cats)
        store.webContent.blockedByFilter = .auto()
    }

    private func clearShields() {
        store.shield.applications = nil
        store.shield.applicationCategories = nil
        store.webContent.blockedByFilter = nil
    }

    private func loadSelection() -> FamilyActivitySelection? {
        guard let data = defaults?.data(forKey: "screenTimeSelection") else { return nil }
        return try? JSONDecoder().decode(FamilyActivitySelection.self, from: data)
    }
}

extension DeviceActivityName {
    static let daily = DeviceActivityName("com.remainfaithful.daily")
}
