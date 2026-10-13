import Foundation

// Queues protection notices when the network fails, and retries them for the
// same account. A later sign-in must not upload the previous person's queue.

enum ProtectionAlerts {
    private static let queueKey = "pendingProtectionAlerts"
    private static let defaults = UserDefaults(suiteName: "group.com.remainfaithful.app")

    struct Item: Codable {
        let userID: Int
        let type: String
        let detail: String
    }

    static func send(type: String, detail: String) {
        Task { await deliver(type: type, detail: detail) }
    }

    static func deliver(type: String, detail: String) async {
        guard APIClient.shared.isAuthenticated else {
            enqueue(type: type, detail: detail)
            return
        }
        let attempts = fallbackTypes(for: type)
        var lastWasNetwork = false
        for (index, candidate) in attempts.enumerated() {
            do {
                try await APIClient.shared.sendProtectionAlert(type: candidate, detail: detail)
                return
            } catch let error as APIError {
                switch error {
                case .network:
                    lastWasNetwork = true
                case .server where index < attempts.count - 1:
                    continue
                default:
                    lastWasNetwork = false
                }
            } catch {
                lastWasNetwork = true
            }
            if lastWasNetwork { break }
        }
        if lastWasNetwork {
            enqueue(type: type, detail: detail)
        }
    }

    static func flush() async {
        guard APIClient.shared.isAuthenticated else { return }
        let mine = queued().filter { $0.userID == currentUserID }
        var remaining = queued().filter { $0.userID != currentUserID }
        for item in mine {
            do {
                try await APIClient.shared.sendProtectionAlert(type: item.type, detail: item.detail)
            } catch let error as APIError {
                if case .network = error {
                    remaining.append(item)
                } else if item.type == "signed_out" {
                    try? await APIClient.shared.sendProtectionAlert(
                        type: "monitoring_disabled",
                        detail: item.detail
                    )
                }
            } catch {
                remaining.append(item)
            }
        }
        save(remaining)
    }

    /// Drop this account's unsent notices after sign-out so the next person
    /// cannot upload them.
    static func dropCurrentUser() {
        let kept = queued().filter { $0.userID != currentUserID }
        save(kept)
    }

    private static func fallbackTypes(for type: String) -> [String] {
        if type == "signed_out" {
            return ["signed_out", "monitoring_disabled"]
        }
        return [type]
    }

    private static var currentUserID: Int {
        AuthState.shared.currentUser?.id ?? 0
    }

    private static func enqueue(type: String, detail: String) {
        var items = queued()
        items.append(Item(userID: currentUserID, type: type, detail: detail))
        if items.count > 40 { items = Array(items.suffix(40)) }
        save(items)
    }

    private static func queued() -> [Item] {
        guard let data = defaults?.data(forKey: queueKey) else { return [] }
        return (try? JSONDecoder().decode([Item].self, from: data)) ?? []
    }

    private static func save(_ items: [Item]) {
        if let data = try? JSONEncoder().encode(items) {
            defaults?.set(data, forKey: queueKey)
        }
    }
}
