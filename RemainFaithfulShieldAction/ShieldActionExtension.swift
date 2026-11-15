import ManagedSettings
import UserNotifications

/// Shield Action returns a response. It does not read an app name, a bundle id, or a domain.
/// The system does not provide those values here. The only app-group write is a timestamp.
class ShieldActionExtension: ShieldActionDelegate {

    override func handle(
        action: ShieldAction,
        for application: ApplicationToken,
        completionHandler: @escaping (ShieldActionResponse) -> Void
    ) {
        respond(to: action, completionHandler: completionHandler)
    }

    override func handle(
        action: ShieldAction,
        for webDomain: WebDomainToken,
        completionHandler: @escaping (ShieldActionResponse) -> Void
    ) {
        respond(to: action, completionHandler: completionHandler)
    }

    override func handle(
        action: ShieldAction,
        for category: ActivityCategoryToken,
        completionHandler: @escaping (ShieldActionResponse) -> Void
    ) {
        respond(to: action, completionHandler: completionHandler)
    }

    private func respond(
        to action: ShieldAction,
        completionHandler: @escaping (ShieldActionResponse) -> Void
    ) {
        guard action == .primaryButtonPressed else {
            completionHandler(.close)
            return
        }
        BlockedOpenRequest.record()
        if #available(iOS 26.5, *) {
            completionHandler(.openParentalControlsApp)
        } else {
            postOpenNotification()
            completionHandler(.close)
        }
    }

    private func postOpenNotification() {
        let content = UNMutableNotificationContent()
        content.title = "Remain Faithful"
        content.body = "A blocked app was opened"
        content.sound = .default
        content.userInfo = ["notification_type": BlockedOpenRequest.notificationType]
        let request = UNNotificationRequest(
            identifier: "blocked-open-\(UUID().uuidString)",
            content: content,
            trigger: nil
        )
        UNUserNotificationCenter.current().add(request, withCompletionHandler: nil)
    }
}
