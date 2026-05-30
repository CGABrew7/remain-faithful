import SwiftUI
import UserNotifications
import GoogleSignIn

@main
struct RemainFaithfulApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @StateObject private var appState  = AppState.shared
    @StateObject private var authState = AuthState.shared

    init() {
        // Clear stale onboarding state when no account exists in Keychain.
        // UserDefaults (AppStorage) persists across simulator installs;
        // Keychain does not — so Keychain is the authoritative "has account" signal.
        if AuthState.shared.currentUser == nil && !AuthState.shared.isAuthenticated {
            UserDefaults.standard.removeObject(forKey: "hasCompletedOnboarding")
        }
    }

    // In DEBUG builds, `xcrun simctl spawn <id> defaults write com.remainfaithful.app
    // debugShowDashboard -bool true` bypasses auth so the Dashboard can be tested
    // in the simulator without a running backend.
    private var showDashboard: Bool {
        guard !authState.isAuthenticated else { return true }
        #if DEBUG
        return UserDefaults.standard.bool(forKey: "debugShowDashboard")
        #else
        return false
        #endif
    }

    var body: some Scene {
        WindowGroup {
            Group {
                if showDashboard {
                    ContentView()
                } else if hasCompletedOnboarding || authState.currentUser != nil {
                    LoginView()
                } else {
                    OnboardingView()
                }
            }
            .environmentObject(appState)
            .environmentObject(authState)
            .task {
                try? await APIClient.shared.refreshTokenIfNeeded()
                // Drain any events queued by the broadcast extension while the app was closed.
                await EventProcessor.shared.processPendingEvents()
            }
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
        configureGoogleSignIn()
        // TODO: migrate to BGAppRefreshTask (deprecated in iOS 13).
        // Add "com.remainfaithful.app.fetch" to BGTaskSchedulerPermittedIdentifiers in
        // Xcode → Target → Info, then replace with BGTaskScheduler.shared.register + submit.
        application.setMinimumBackgroundFetchInterval(
            UIApplication.backgroundFetchIntervalMinimum
        )
        if let payload = launchOptions?[.remoteNotification] as? [AnyHashable: Any] {
            handlePayload(payload)
        }
        return true
    }

    private func configureGoogleSignIn() {
        guard GIDSignIn.sharedInstance.configuration == nil else { return }
        if let path = Bundle.main.path(forResource: "GoogleService-Info", ofType: "plist"),
           let plist = NSDictionary(contentsOfFile: path),
           let clientID = plist["CLIENT_ID"] as? String {
            GIDSignIn.sharedInstance.configuration = GIDConfiguration(clientID: clientID)
        }
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

    /// Intercept foreground notifications.
    /// CONTENT_FLAGGED → show our own in-app banner; all others → system banner.
    func userNotificationCenter(
        _ center: UNUserNotificationCenter,
        willPresent notification: UNNotification,
        withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void
    ) {
        let userInfo = notification.request.content.userInfo
        let type = userInfo["notification_type"] as? String

        if type == NotificationType.contentFlagged {
            let senderName = userInfo["sender_name"] as? String ?? "Your partner"
            let category   = userInfo["category"]    as? String ?? "content"
            let event      = ActivityEvent.from(notificationPayload: userInfo)
            AppState.shared.showPartnerAlert(senderName: senderName,
                                             category: category,
                                             event: event)
            completionHandler([.sound])
        } else if type == "DONATION_THANKS" {
            UserDefaults.standard.set(true, forKey: "hasDonated")
            completionHandler([.banner, .sound])
        } else {
            completionHandler([.banner, .sound, .badge])
        }
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
