import Foundation

// MARK: - Deep Link

enum DeepLink {
    case dashboard
    case alertDetail(ActivityEvent)
    case panicView
    case group
}

// MARK: - App State

/// Singleton observable used to drive cross-cutting navigation from push
/// notification taps. Singleton so AppDelegate can reach it without a
/// SwiftUI environment.
final class AppState: ObservableObject {
    static let shared = AppState()
    @Published var deepLink: DeepLink?
    private init() {}

    func navigate(to link: DeepLink) {
        DispatchQueue.main.async { self.deepLink = link }
    }
}
