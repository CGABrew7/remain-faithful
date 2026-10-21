import Foundation

// MARK: - Stored user (light snapshot kept in Keychain alongside the token)

struct StoredUser: Codable {
    let id: Int
    let name: String
    let email: String
    let createdAt: String?
    var emailVerified: Bool? = nil
}

// MARK: - AuthState

/// Singleton observable that owns the session token stored in Keychain.
/// SwiftUI views observe `isAuthenticated`; all token reads/writes go through here.
final class AuthState: ObservableObject {
    static let shared = AuthState(keychain: .shared)

    @Published private(set) var isAuthenticated: Bool
    @Published private(set) var currentUser: StoredUser?

    let keychain: KeychainHelper
    private var inMemoryToken: String?
    private let sharedDefaults = UserDefaults(suiteName: "group.com.remainfaithful.app")

    /// UserDefaults key for the Group tab's selected group (`@AppStorage` in GroupView
    /// and SettingsView). It belongs to the signed-in account, so it is cleared
    /// whenever the session ends or a different account signs in.
    static let primaryGroupIDKey = "primaryGroupID"

    init(keychain: KeychainHelper = .shared) {
        self.keychain = keychain
        inMemoryToken = keychain.get("authToken") ?? sharedDefaults?.string(forKey: "authToken")
        isAuthenticated = inMemoryToken != nil
        if let data = keychain.getData("currentUser"),
           let user = try? JSONDecoder().decode(StoredUser.self, from: data) {
            currentUser = user
            sharedDefaults?.set(user.id, forKey: "authUserID")
            emailRequired = UserDefaults.standard.bool(forKey: Self.emailRequiredKey(userID: user.id))
            pendingEmail = UserDefaults.standard.string(forKey: Self.pendingEmailKey(userID: user.id))
        }
    }

    @Published private(set) var emailRequired = false
    @Published private(set) var pendingEmail: String?

    private static func emailRequiredKey(userID: Int) -> String { "emailRequired.\(userID)" }
    private static func pendingEmailKey(userID: Int) -> String { "pendingEmail.\(userID)" }

    var token: String? { inMemoryToken ?? keychain.get("authToken") ?? sharedDefaults?.string(forKey: "authToken") }

    func setSession(token: String, user: RemoteUser) {
        if currentUser?.id != user.id {
            UserDefaults.standard.removeObject(forKey: Self.primaryGroupIDKey)
            clearQueuedEvents()
        }
        inMemoryToken = token
        keychain.set(token, for: "authToken")
        sharedDefaults?.set(token, forKey: "authToken")
        sharedDefaults?.set(user.id, forKey: "authUserID")
        sharedDefaults?.set(false, forKey: "extensionAuthFailed")
        let stored = StoredUser(id: user.id, name: user.name, email: user.email, createdAt: user.createdAt, emailVerified: user.emailVerified)
        if let data = try? JSONEncoder().encode(stored) {
            keychain.setData(data, for: "currentUser")
        }
        let required = UserDefaults.standard.bool(forKey: Self.emailRequiredKey(userID: user.id))
        let pending = UserDefaults.standard.string(forKey: Self.pendingEmailKey(userID: user.id))
        let verified = user.emailVerified == true
        DispatchQueue.main.async {
            self.currentUser = stored
            self.isAuthenticated = true
            if verified {
                self.clearMailboxFlags(userID: user.id)
            } else {
                self.emailRequired = required
                self.pendingEmail = pending
            }
        }
    }

    /// Apple sign-in with no shared email asks the person to add one in profile.
    func noteAppleSignIn(userID: Int, email: String, appleUserID: String) {
        let missing = AppleMailbox.isStandIn(email: email, appleUserID: appleUserID)
        UserDefaults.standard.set(missing, forKey: Self.emailRequiredKey(userID: userID))
        guard missing else { return }
        DispatchQueue.main.async {
            if self.currentUser?.id == userID || self.currentUser == nil {
                self.emailRequired = true
            }
        }
    }

    func noteEmailRequired() {
        guard let id = currentUser?.id else { return }
        UserDefaults.standard.set(true, forKey: Self.emailRequiredKey(userID: id))
        DispatchQueue.main.async { self.emailRequired = true }
    }

    func notePendingEmail(_ email: String) {
        guard let id = currentUser?.id else { return }
        let trimmed = email.trimmingCharacters(in: .whitespacesAndNewlines)
        guard trimmed.contains("@") else { return }
        UserDefaults.standard.set(trimmed, forKey: Self.pendingEmailKey(userID: id))
        UserDefaults.standard.set(false, forKey: Self.emailRequiredKey(userID: id))
        DispatchQueue.main.async {
            self.pendingEmail = trimmed
            self.emailRequired = false
        }
    }

    private func clearMailboxFlags(userID: Int) {
        UserDefaults.standard.removeObject(forKey: Self.emailRequiredKey(userID: userID))
        UserDefaults.standard.removeObject(forKey: Self.pendingEmailKey(userID: userID))
        emailRequired = false
        pendingEmail = nil
    }

    /// Refresh writes the same token to memory, Keychain, and the app group
    /// the broadcast extension reads. A Keychain-only write leaves Deep Scan
    /// on a token that dies after about a day.
    func updateToken(_ token: String) {
        inMemoryToken = token
        keychain.set(token, for: "authToken")
        sharedDefaults?.set(token, forKey: "authToken")
        sharedDefaults?.set(false, forKey: "extensionAuthFailed")
    }

    func updateProfile(name: String, email: String) {
        guard let user = currentUser else { return }
        let updated = StoredUser(id: user.id, name: name, email: email, createdAt: user.createdAt, emailVerified: user.emailVerified)
        if let data = try? JSONEncoder().encode(updated) {
            keychain.setData(data, for: "currentUser")
        }
        UserDefaults.standard.set(name,  forKey: "userName")
        UserDefaults.standard.set(email, forKey: "userEmail")
        DispatchQueue.main.async { self.currentUser = updated }
    }

    func refreshFromAPI() async {
        guard isAuthenticated else { return }
        guard let user = try? await APIClient.shared.me() else { return }
        let updated = StoredUser(id: user.id, name: user.name, email: user.email, createdAt: user.createdAt, emailVerified: user.emailVerified)
        if let data = try? JSONEncoder().encode(updated) {
            keychain.setData(data, for: "currentUser")
        }
        UserDefaults.standard.set(user.name,  forKey: "userName")
        UserDefaults.standard.set(user.email, forKey: "userEmail")
        sharedDefaults?.set(user.id, forKey: "authUserID")
        let verified = user.emailVerified == true
        DispatchQueue.main.async {
            self.currentUser = updated
            if verified {
                self.clearMailboxFlags(userID: user.id)
            }
        }
    }

    func markEmailVerified() {
        guard let user = currentUser else { return }
        let updated = StoredUser(id: user.id, name: user.name, email: user.email, createdAt: user.createdAt, emailVerified: true)
        if let data = try? JSONEncoder().encode(updated) {
            keychain.setData(data, for: "currentUser")
        }
        DispatchQueue.main.async {
            self.currentUser = updated
            self.clearMailboxFlags(userID: user.id)
        }
    }

    /// Unsent alerts belong to the person who was signed in. An ended session clears them.
    func clearQueuedEvents() {
        sharedDefaults?.removeObject(forKey: "pendingEvents")
        sharedDefaults?.removeObject(forKey: "pendingRawEvents")
    }

    func clearSession() {
        inMemoryToken = nil
        keychain.delete("authToken")
        keychain.delete("currentUser")
        sharedDefaults?.removeObject(forKey: "authToken")
        sharedDefaults?.removeObject(forKey: "authUserID")
        clearQueuedEvents()
        UserDefaults.standard.removeObject(forKey: Self.primaryGroupIDKey)
        DispatchQueue.main.async {
            self.currentUser = nil
            self.isAuthenticated = false
            self.emailRequired = false
            self.pendingEmail = nil
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
