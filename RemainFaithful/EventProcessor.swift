import Foundation
import os

// MARK: - Shared event model (mirrors DetectedEvent in the broadcast extension)

// IMPORTANT: This struct must remain byte-for-byte identical to DetectedEvent
// in RemainFaithful/EventProcessor.swift. Both targets share App Group UserDefaults
// for IPC. A future refactor should move this to a shared Swift package.
struct DetectedEvent: Codable {
    let id: String
    let timestamp: Double
    let category: String
    let severity: String
    let summary: String
    let tier: Int
    let confidence: Float
    var processed: Bool
}

// MARK: - EventProcessor

/// Reads events written by the broadcast extension into the shared App Group
/// UserDefaults container and POSTs them to the backend.
final class EventProcessor {
    static let shared = EventProcessor()
    private let appGroupID = "group.com.remainfaithful.app"
    private let logger = Logger(subsystem: "com.remainfaithful.app", category: "EventProcessor")
    private lazy var sharedDefaults: UserDefaults? = UserDefaults(suiteName: appGroupID)
    private init() {}

    // MARK: - Called on foreground / background fetch

    func processPendingEvents() async {
        guard APIClient.shared.isAuthenticated else { return }
        guard let defaults = sharedDefaults else { return }

        // Sync the API base URL into the shared container so the extension can read it.
        let apiBase = ProcessInfo.processInfo.environment["API_BASE_URL"] ?? "https://remain-faithful-api.fly.dev"
        defaults.set(apiBase, forKey: "apiBaseURL")

        guard let data = defaults.data(forKey: "pendingEvents"),
              var events = try? JSONDecoder().decode([DetectedEvent].self, from: data) else { return }

        let pending = events.indices.filter { !events[$0].processed }
        guard !pending.isEmpty else { return }

        logger.info("Processing \(pending.count) pending events from broadcast extension")

        var changed = false
        for i in pending {
            do {
                _ = try await APIClient.shared.createEvent(
                    category: events[i].category,
                    severity: events[i].severity,
                    summary: events[i].summary
                )
                events[i].processed = true
                changed = true
                logger.info("Posted event \(events[i].id) (\(events[i].category))")
            } catch {
                // Stop on first network failure to avoid piling up failures
                logger.error("Failed to post event \(events[i].id): \(error.localizedDescription)")
                break
            }
        }

        if changed, let newData = try? JSONEncoder().encode(events) {
            defaults.set(newData, forKey: "pendingEvents")
        }
    }

    // MARK: - Debugging helpers

    func pendingEventCount() -> Int {
        guard let defaults = sharedDefaults,
              let data = defaults.data(forKey: "pendingEvents"),
              let events = try? JSONDecoder().decode([DetectedEvent].self, from: data) else { return 0 }
        return events.filter { !$0.processed }.count
    }

    func isBroadcasting() -> Bool {
        sharedDefaults?.bool(forKey: "isBroadcasting") ?? false
    }
}
