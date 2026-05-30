import Foundation

// MARK: - Stored user (light snapshot kept in Keychain alongside the token)

struct StoredUser: Codable {
    let id: Int
    let name: String
    let email: String
}

// MARK: - AuthState

/// Singleton observable that owns the session token stored in Keychain.
/// SwiftUI views observe `isAuthenticated`; all token reads/writes go through here.
final class AuthState: ObservableObject {
    static let shared = AuthState()

    @Published private(set) var isAuthenticated: Bool
    @Published private(set) var currentUser: StoredUser?

    private init() {
        isAuthenticated = KeychainHelper.shared.get("authToken") != nil
        if let data = KeychainHelper.shared.getData("currentUser"),
           let user = try? JSONDecoder().decode(StoredUser.self, from: data) {
            currentUser = user
        }
    }

    var token: String? { KeychainHelper.shared.get("authToken") }

    func setSession(token: String, user: RemoteUser) {
        KeychainHelper.shared.set(token, for: "authToken")
        let stored = StoredUser(id: user.id, name: user.name, email: user.email)
        if let data = try? JSONEncoder().encode(stored) {
            KeychainHelper.shared.setData(data, for: "currentUser")
        }
        DispatchQueue.main.async {
            self.currentUser = stored
            self.isAuthenticated = true
        }
    }

    func clearSession() {
        KeychainHelper.shared.delete("authToken")
        KeychainHelper.shared.delete("currentUser")
        DispatchQueue.main.async {
            self.currentUser = nil
            self.isAuthenticated = false
        }
    }

    /// Decode the JWT expiry without a dependency — used for auto-refresh.
    var tokenExpiresAt: Date? {
        guard let t = token else { return nil }
        let parts = t.components(separatedBy: ".")
        guard parts.count == 3 else { return nil }
        var b64 = parts[1]
        let rem = b64.count % 4
        if rem != 0 { b64 += String(repeating: "=", count: 4 - rem) }
        b64 = b64.replacingOccurrences(of: "-", with: "+")
                 .replacingOccurrences(of: "_", with: "/")
        guard let data = Data(base64Encoded: b64),
              let obj  = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              let exp  = obj["exp"] as? TimeInterval else { return nil }
        return Date(timeIntervalSince1970: exp)
    }
}
