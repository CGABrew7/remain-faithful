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

// MARK: - EventSending

protocol EventSending {
    var isAuthenticated: Bool { get }
    func createEvent(category: String, severity: String, summary: String,
                     timestamp: String?) async throws -> RemoteEvent
}

// MARK: - EventProcessor

/// Reads events written by the broadcast extension into the shared App Group
/// UserDefaults container and POSTs them to the backend.
final class EventProcessor {
    static let shared = EventProcessor()
    private let appGroupID: String
    private let eventSender: EventSending
    private let logger = Logger(subsystem: "com.remainfaithful.app", category: "EventProcessor")
    private lazy var sharedDefaults: UserDefaults? = UserDefaults(suiteName: appGroupID)

    private init() {
        appGroupID  = "group.com.remainfaithful.app"
        eventSender = APIClient.shared
    }

    init(appGroupID: String, eventSender: EventSending) {
        self.appGroupID  = appGroupID
        self.eventSender = eventSender
    }

    // MARK: - Called on foreground / background fetch

    func processPendingEvents() async {
        guard eventSender.isAuthenticated else { return }
        guard let defaults = sharedDefaults else { return }

        // Sync configuration into the shared container so the extension can read it.
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
                _ = try await eventSender.createEvent(
                    category: events[i].category,
                    severity: events[i].severity,
                    summary: events[i].summary,
                    timestamp: nil
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

        // Drain events queued by the DeviceActivityMonitor extension (raw JSON dicts,
        // separate key because the extension can't encode DetectedEvent directly).
        await ActivitySelectionManager.shared.drainPendingRawEvents()
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
