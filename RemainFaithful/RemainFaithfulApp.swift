import SwiftUI
import UserNotifications

@main
struct RemainFaithfulApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @StateObject private var appState = AppState.shared

    var body: some Scene {
        WindowGroup {
            Group {
                if hasCompletedOnboarding {
                    ContentView()
                } else {
                    OnboardingView()
                }
            }
            .environmentObject(appState)
        }
    }
}

// MARK: - App Delegate

final class AppDelegate: NSObject, UIApplicationDelegate, UNUserNotificationCenterDelegate {

    func application(
        _ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
    ) -> Bool {
        UNUserNotificationCenter.current().delegate = self
        // Enable background fetch for silent alert polling.
        application.setMinimumBackgroundFetchInterval(
            UIApplication.backgroundFetchIntervalMinimum
        )
        // Handle notification that cold-launched the app.
        if let payload = launchOptions?[.remoteNotification] as? [AnyHashable: Any] {
            handlePayload(payload)
        }
        return true
    }

    // MARK: - Remote notification registration

    func application(
        _ application: UIApplication,
        didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data
    ) {
        NotificationService.shared.handleDeviceToken(deviceToken)
    }

    func application(
        _ application: UIApplication,
        didFailToRegisterForRemoteNotificationsWithError error: Error
    ) {
        print("[APNs] registration failed: \(error)")
    }

    // MARK: - Background fetch

    func application(
        _ application: UIApplication,
        performFetchWithCompletionHandler completionHandler: @escaping (UIBackgroundFetchResult) -> Void
    ) {
        NotificationService.shared.performBackgroundFetch(completionHandler: completionHandler)
    }

    // MARK: - UNUserNotificationCenterDelegate

    /// Show banners even when the app is in the foreground.
    func userNotificationCenter(
        _ center: UNUserNotificationCenter,
        willPresent notification: UNNotification,
        withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void
    ) {
        completionHandler([.banner, .sound, .badge])
    }

    /// Route taps whether the app was in background or was cold-launched.
    func userNotificationCenter(
        _ center: UNUserNotificationCenter,
        didReceive response: UNNotificationResponse,
        withCompletionHandler completionHandler: @escaping () -> Void
    ) {
        NotificationService.shared.handleTap(response)
        completionHandler()
    }

    // MARK: - Private

    private func handlePayload(_ payload: [AnyHashable: Any]) {
        // Defer navigation until after the root view appears by posting to the
        // AppState singleton (which will be picked up by ContentView.onReceive).
        let type = payload["notification_type"] as? String ?? ""
        switch type {
        case NotificationType.contentFlagged:
            if let event = ActivityEvent.from(notificationPayload: payload) {
                AppState.shared.navigate(to: .alertDetail(event))
            }
        case NotificationType.panicAlert:
            AppState.shared.navigate(to: .panicView)
        default:
            break
        }
    }
}
