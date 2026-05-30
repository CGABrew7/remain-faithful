import Foundation

// MARK: - Deep Link

enum DeepLink {
    case dashboard
    case alertDetail(ActivityEvent)
    case panicView
    case group
}

// MARK: - Partner Flag Alert (shown as in-app banner)

struct PartnerFlagAlert: Identifiable {
    let id = UUID()
    let senderName: String
    let category: String
    let event: ActivityEvent?

    var displayCategory: String {
        switch category {
        case "adult_content":  return "adult content"
        case "gambling":       return "gambling"
        case "violence":       return "violent content"
        case "self_harm":      return "self-harm content"
        default:               return category.replacingOccurrences(of: "_", with: " ")
        }
    }
}

// MARK: - App State

/// Singleton observable that drives cross-cutting navigation and real-time
/// UI state (notification taps, partner alert banners, unread badge count).
/// Singleton so AppDelegate can reach it without a SwiftUI environment.
final class AppState: ObservableObject {
    static let shared = AppState()

    @Published var deepLink: DeepLink?
    @Published var partnerFlagAlert: PartnerFlagAlert?
    @Published var unreadAlertCount: Int = 0

    private init() {}

    func navigate(to link: DeepLink) {
        DispatchQueue.main.async { self.deepLink = link }
    }

    func showPartnerAlert(senderName: String, category: String, event: ActivityEvent?) {
        DispatchQueue.main.async {
            self.partnerFlagAlert = PartnerFlagAlert(
                senderName: senderName,
                category: category,
                event: event
            )
            self.unreadAlertCount += 1
        }
    }

    func clearPartnerAlert() {
        DispatchQueue.main.async { self.partnerFlagAlert = nil }
    }

    func resetUnreadCount() {
        DispatchQueue.main.async { self.unreadAlertCount = 0 }
    }
}
