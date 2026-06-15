import Foundation
import Combine
import os

private let lockoutLogger = Logger(subsystem: "com.remainfaithful.app", category: "AppLockout")

// Feature 1 — App Lockout tied to Deep Scan (optional, default OFF).
//
// When enabled, the selected apps remain shielded via ManagedSettings except
// while the ReplayKit broadcast (Deep Scan) is active. When Deep Scan stops,
// the shield re-applies and the accountability partner is notified.
//
// Cross-process IPC: the broadcast extension posts a Darwin notification whenever
// isBroadcasting changes. This manager observes that notification and immediately
// syncs the shield state via ActivitySelectionManager. If the main app is suspended
// at that moment, syncShieldState() catches up on the next scene-active transition.
//
// No new entitlement is required. ManagedSettingsStore lives in the main app target
// (which already holds com.apple.developer.family-controls). The broadcast extension
// only writes isBroadcasting to shared UserDefaults; it never calls ManagedSettings.
final class AppLockoutManager: ObservableObject {
    static let shared = AppLockoutManager()

    private let appGroupID = "group.com.remainfaithful.app"
    private let defaults: UserDefaults?

    @Published private(set) var isEnabled: Bool = false

    private init() {
        defaults = UserDefaults(suiteName: appGroupID)
        isEnabled = defaults?.bool(forKey: "appLockoutEnabled") ?? false
        registerDarwinObserver()
    }

    // MARK: - Enable / disable

    func setEnabled(_ enabled: Bool) {
        isEnabled = enabled
        defaults?.set(enabled, forKey: "appLockoutEnabled")
        if enabled {
            syncShieldState()
        } else {
            // Restore normal (user-configured) shield state when lockout is turned off.
            let asm = ActivitySelectionManager.shared
            if asm.isShieldingEnabled {
                asm.applyShielding()
            } else {
                asm.clearShielding()
            }
        }
    }

    // MARK: - Shield sync

    // Call on scene becoming active and on Darwin notification from broadcast extension.
    func syncShieldState() {
        guard isEnabled else { return }
        let isBroadcasting = defaults?.bool(forKey: "isBroadcasting") ?? false
        let asm = ActivitySelectionManager.shared
        if isBroadcasting {
            asm.clearShielding()
            lockoutLogger.info("Lockout: broadcast active — shield lifted")
        } else {
            asm.applyShielding()
            lockoutLogger.info("Lockout: broadcast inactive — shield applied")
        }
    }

    // MARK: - Darwin notification observer

    private func registerDarwinObserver() {
        // The closure below references only AppLockoutManager.shared, a global
        // static — no local capture — so it is implicitly @convention(c) compatible.
        CFNotificationCenterAddObserver(
            CFNotificationCenterGetDarwinNotifyCenter(),
            nil,
            { _, _, _, _, _ in AppLockoutManager.shared.handleBroadcastStateChange() },
            "com.remainfaithful.broadcastStateChanged" as CFString,
            nil,
            .deliverImmediately
        )
    }

    func handleBroadcastStateChange() {
        let isBroadcasting = defaults?.bool(forKey: "isBroadcasting") ?? false
        DispatchQueue.main.async {
            guard self.isEnabled else { return }
            if isBroadcasting {
                ActivitySelectionManager.shared.clearShielding()
                lockoutLogger.info("Lockout (Darwin): broadcast started — shield lifted")
            } else {
                ActivitySelectionManager.shared.applyShielding()
                lockoutLogger.info("Lockout (Darwin): broadcast ended — shield re-applied")
                Task {
                    try? await APIClient.shared.sendProtectionAlert(
                        type: "deep_scan_stopped",
                        detail: "Deep Scan ended — app lockout shields re-applied."
                    )
                }
            }
        }
    }
}
