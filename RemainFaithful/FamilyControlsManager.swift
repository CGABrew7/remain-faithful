import Foundation
import FamilyControls
import Combine

// Manages FamilyControls authorization for self-monitoring (individual mode).
// Individual mode lets an adult monitor their own device — no parental pairing needed.
// Authorization persists until the user revokes it in iOS Settings > Screen Time.
final class FamilyControlsManager: ObservableObject {
    static let shared = FamilyControlsManager()

    @Published private(set) var authorizationStatus: AuthorizationStatus = .notDetermined

    private var cancellables = Set<AnyCancellable>()

    private init() {
        authorizationStatus = AuthorizationCenter.shared.authorizationStatus
        // Watch for revocation (e.g., user turns off Screen Time in iOS Settings).
        AuthorizationCenter.shared.objectWillChange
            .receive(on: DispatchQueue.main)
            .sink { [weak self] in
                self?.authorizationStatus = AuthorizationCenter.shared.authorizationStatus
            }
            .store(in: &cancellables)
    }

    /// If Screen Time was allowed and is no longer, tell the partner.
    /// A change noticed while signed out is sent the next time someone signs in.
    func noteRevocationIfNeeded() {
        let approvedKey = "lastKnownFamilyControlsApproved"
        let pendingKey = "familyControlsRevocationUnsent"
        let nowApproved = authorizationStatus == .approved
        let wasApproved = UserDefaults.standard.bool(forKey: approvedKey)
        let pending = UserDefaults.standard.bool(forKey: pendingKey)
        let plan = ScreenTimeRevocation.plan(
            wasApproved: wasApproved,
            nowApproved: nowApproved,
            authenticated: AuthState.shared.isAuthenticated,
            alreadyPending: pending
        )
        UserDefaults.standard.set(plan.pending, forKey: pendingKey)
        UserDefaults.standard.set(plan.storeApproved, forKey: approvedKey)
        if plan.send {
            ProtectionAlerts.send(
                type: "family_controls_revoked",
                detail: "Screen Time access was turned off. Alerts and Blocking are off."
            )
        }
    }

    func requestAuthorization() async {
        do {
            try await AuthorizationCenter.shared.requestAuthorization(for: .individual)
        } catch {
            // Denial or restriction; status is captured below regardless.
        }
        await MainActor.run {
            authorizationStatus = AuthorizationCenter.shared.authorizationStatus
        }
    }
}
