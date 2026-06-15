import Foundation
import os

private let pinLogger = Logger(subsystem: "com.remainfaithful.app", category: "PartnerPIN")

// Feature 2 — Partner PIN (optional, default OFF).
//
// A 4-digit PIN that gates security-reducing actions. The PIN is set and stored
// server-side by the accountability partner — the monitored user cannot view or
// reset it. To gate an action, call isPINSet to decide whether to require the
// PIN entry sheet, then call verifyPIN(_:) with the entered digits.
//
// Set/share flow:
//   1. The monitoring partner opens Manage Partners → taps the lock icon on
//      the partner row where they are the monitoring side.
//   2. They enter and confirm a 4-digit PIN in SetPINSheet.
//   3. The PIN is POSTed to /relationships/{id}/pin and stored as a bcrypt
//      hash server-side — it never travels back to the client in plaintext.
//   4. The partner shares the PIN with the monitored user out-of-band
//      (verbally, text message, etc.).
//   5. Subsequent gated actions on the monitored user's device require the PIN.
//
// The monitored user cannot reset the PIN, change it, or bypass the gate —
// they can only enter it. Wrong PIN attempts notify the partner via APNs.
@MainActor
final class PartnerPINManager: ObservableObject {
    static let shared = PartnerPINManager()

    @Published private(set) var isPINSet: Bool = false
    @Published private(set) var isChecking: Bool = false

    private init() {}

    // MARK: - Status

    func refreshStatus() async {
        guard APIClient.shared.isAuthenticated else { return }
        isChecking = true
        defer { isChecking = false }
        do {
            isPINSet = try await APIClient.shared.getPartnerPINStatus()
            pinLogger.info("Partner PIN status: isSet=\(self.isPINSet)")
        } catch {
            pinLogger.info("Could not fetch partner PIN status: \(error.localizedDescription)")
        }
    }

    // MARK: - Verification

    func verifyPIN(_ pin: String) async -> Bool {
        guard APIClient.shared.isAuthenticated else { return false }
        do {
            let success = try await APIClient.shared.verifyPartnerPIN(pin: pin)
            if !success { pinLogger.notice("Wrong PIN attempt") }
            return success
        } catch {
            pinLogger.error("PIN verification error: \(error.localizedDescription)")
            return false
        }
    }

    // MARK: - Wrong attempt notification

    func notifyWrongAttempt(userName: String) async {
        try? await APIClient.shared.sendProtectionAlert(
            type: "wrong_pin_attempt",
            detail: "\(userName) entered a wrong Protection PIN."
        )
    }
}
