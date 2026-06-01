import UserNotifications
import UIKit

// MARK: - Notification type constants

enum NotificationType {
    static let contentFlagged    = "CONTENT_FLAGGED"
    static let monitoringPaused  = "MONITORING_PAUSED"
    static let monitoringResumed = "MONITORING_RESUMED"
    static let appDeleted        = "APP_DELETED"
    static let panicAlert        = "PANIC_ALERT"
    static let streakMilestone   = "STREAK_MILESTONE"
}

// MARK: - Notification Service

/// Handles the full APNs lifecycle: permission request, device token upload,
/// foreground presentation, background fetch, and notification-tap routing.
final class NotificationService: NSObject {
    static let shared = NotificationService()
    private override init() {}

    // MARK: - Permission

    /// Call once after onboarding completes. Requests authorisation then
    /// registers with APNs if granted.
    func requestPermission() {
        UNUserNotificationCenter.current().requestAuthorization(
            options: [.alert, .badge, .sound]
        ) { granted, error in
            if let error { print("[APNs] permission error: \(error)") }
            guard granted else { return }
            DispatchQueue.main.async {
                UIApplication.shared.registerForRemoteNotifications()
            }
        }
    }

    /// Call on every authenticated app launch to refresh the device token with APNs.
    /// If permission has already been granted, this re-registers and updates the
    /// token on the backend. If permission is not determined, it does nothing
    /// (the onboarding flow will call requestPermission() for new users).
    func ensureRegisteredIfAuthorized() {
        Task {
            let center   = UNUserNotificationCenter.current()
            let settings = await center.notificationSettings()
            guard settings.authorizationStatus == .authorized else { return }
            await MainActor.run {
                UIApplication.shared.registerForRemoteNotifications()
            }
        }
    }

    /// Publish whether the user has denied notification permission.
    /// Used to show an in-app prompt explaining why notifications matter.
    func checkNotificationPermission() async -> UNAuthorizationStatus {
        await UNUserNotificationCenter.current().notificationSettings().authorizationStatus
    }

    // MARK: - Device Token

    /// Convert the raw token data to a hex string and upload it to the backend.
    func handleDeviceToken(_ data: Data) {
        let token = data.map { String(format: "%02x", $0) }.joined()
        UserDefaults.standard.set(token, forKey: "apnsDeviceToken")
        Task {
            try? await APIClient.shared.registerDeviceToken(token)
        }
    }

    // MARK: - Notification Tap Routing

    /// Called by the UNUserNotificationCenter delegate when the user taps a
    /// notification. Routes to the appropriate screen via AppState.
    func handleTap(_ response: UNNotificationResponse) {
        routeNotification(userInfo: response.notification.request.content.userInfo)
    }

    func routeNotification(userInfo: [AnyHashable: Any]) {
        let type = userInfo["notification_type"] as? String ?? ""

        switch type {
        case NotificationType.contentFlagged:
            if let event = ActivityEvent.from(notificationPayload: userInfo) {
                AppState.shared.navigate(to: .alertDetail(event))
            } else {
                AppState.shared.navigate(to: .dashboard)
            }

        case NotificationType.panicAlert:
            AppState.shared.navigate(to: .panicView)

        case NotificationType.appDeleted,
             NotificationType.monitoringPaused,
             NotificationType.monitoringResumed,
             NotificationType.streakMilestone:
            AppState.shared.navigate(to: .dashboard)

        default:
            AppState.shared.navigate(to: .dashboard)
        }
    }

    // MARK: - Background Fetch

    /// Called by `application(_:performFetchWithCompletionHandler:)`.
    /// Processes events from the broadcast extension AND polls the server for
    /// new alerts so the dashboard can refresh.
    func performBackgroundFetch(completionHandler: @escaping (UIBackgroundFetchResult) -> Void) {
        guard APIClient.shared.isAuthenticated else {
            completionHandler(.noData)
            return
        }
        Task {
            // Drain any events the broadcast extension queued in shared defaults.
            await EventProcessor.shared.processPendingEvents()
            do {
                let events = try await APIClient.shared.listEvents()
                completionHandler(events.isEmpty ? .noData : .newData)
            } catch {
                completionHandler(.failed)
            }
        }
    }
}
