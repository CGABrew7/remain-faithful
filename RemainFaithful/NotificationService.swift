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
        let info = response.notification.request.content.userInfo
        let type = info["notification_type"] as? String ?? ""

        switch type {
        case NotificationType.contentFlagged:
            if let event = ActivityEvent.from(notificationPayload: info) {
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

    /// Used by `application(_:performFetchWithCompletionHandler:)`.
    func performBackgroundFetch(completionHandler: @escaping (UIBackgroundFetchResult) -> Void) {
        guard APIClient.shared.isAuthenticated else {
            completionHandler(.noData)
            return
        }
        Task {
            do {
                let events = try await APIClient.shared.listEvents()
                completionHandler(events.isEmpty ? .noData : .newData)
            } catch {
                completionHandler(.failed)
            }
        }
    }
}
