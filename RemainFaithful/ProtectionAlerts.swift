import Foundation

// Queues protection notices when the network fails, and retries them for the
// same account. A later sign-in must not upload the previous person's queue.

enum DeliveryResult {
    case sent
    case ended
    case failed
}

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

    /// Asks the server to record one notice and waits for that request.
    /// A network failure stays on this phone so the caller can retry.
    static func confirmDelivery(type: String, detail: String) async -> DeliveryResult {
        guard APIClient.shared.isAuthenticated else { return .ended }
        do {
            try await APIClient.shared.sendProtectionAlert(type: type, detail: detail)
            return .sent
        } catch let error as APIError {
            if case .unauthenticated = error { return .ended }
            return .failed
        } catch {
            return .failed
        }
    }

    @discardableResult
    static func deliver(type: String, detail: String) async -> Bool {
        guard APIClient.shared.isAuthenticated else {
            enqueue(type: type, detail: detail)
            return false
        }
        let attempts = fallbackTypes(for: type)
        var lastWasNetwork = false
        for (index, candidate) in attempts.enumerated() {
            do {
                try await APIClient.shared.sendProtectionAlert(type: candidate, detail: detail)
                if candidate == CoachCopy.alertType { CoachStore.markNudge() }
                return true
            } catch let error as APIError {
                switch error {
                case .network:
                    lastWasNetwork = true
                case .server(let message, _) where message.contains("already sent today"):
                    CoachStore.markNudge()
                    return false
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
        return false
    }

    static func flush() async {
        guard APIClient.shared.isAuthenticated else { return }
        let mine = queued().filter { $0.userID == currentUserID && shouldFlush($0) }
        var remaining = queued().filter { $0.userID != currentUserID }
        for item in mine {
            do {
                try await APIClient.shared.sendProtectionAlert(type: item.type, detail: item.detail)
                if item.type == CoachCopy.alertType { CoachStore.markNudge() }
            } catch let error as APIError {
                if case .server(let message, _) = error, message.contains("already sent today") {
                    CoachStore.markNudge()
                } else if case .network = error {
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
        save(AccountQueue.withoutUser(currentUserID, items: queued()))
    }

    /// A queued partner check-in leaves only while the person still opted in,
    /// the daily cap is open, and the sentence is the fixed one.
    private static func shouldFlush(_ item: Item) -> Bool {
        guard item.type == CoachCopy.alertType else { return true }
        guard item.detail == CoachCopy.partnerSentence else { return false }
        let now = Date()
        let signals = CoachStore.signals(now: now)
        if !signals.partnerNudgeEnabled { return false }
        if CoachEngine.nudgedToday(clock: CoachEngine.clock(from: now, timeZone: .current), signals: signals) {
            return false
        }
        return true
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
