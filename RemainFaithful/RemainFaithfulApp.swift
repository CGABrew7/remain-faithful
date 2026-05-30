import SwiftUI
import UserNotifications
import GoogleSignIn

@main
struct RemainFaithfulApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @StateObject private var appState  = AppState.shared
    @StateObject private var authState = AuthState.shared

    var body: some Scene {
        WindowGroup {
            Group {
                if authState.isAuthenticated {
                    ContentView()
                } else if hasCompletedOnboarding {
                    LoginView()
                } else {
                    OnboardingView()
                }
            }
            .environmentObject(appState)
            .environmentObject(authState)
            .task { try? await APIClient.shared.refreshTokenIfNeeded() }
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
        application.setMinimumBackgroundFetchInterval(
            UIApplication.backgroundFetchIntervalMinimum
        )
        if let payload = launchOptions?[.remoteNotification] as? [AnyHashable: Any] {
            handlePayload(payload)
        }
        return true
    }

    // MARK: - Google Sign In URL handling

    func application(
        _ app: UIApplication,
        open url: URL,
        options: [UIApplication.OpenURLOptionsKey: Any] = [:]
    ) -> Bool {
        GIDSignIn.sharedInstance.handle(url)
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

    func userNotificationCenter(
        _ center: UNUserNotificationCenter,
        willPresent notification: UNNotification,
        withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void
    ) {
        completionHandler([.banner, .sound, .badge])
    }

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
