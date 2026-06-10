import DeviceActivity
import Foundation
import os

private let logger = Logger(
    subsystem: "com.remainfaithful.app.DeviceActivityMonitor",
    category: "monitor"
)

// Registered as the extension's principal class via NSExtensionPrincipalClass.
// The system launches this process when a DeviceActivityEvent threshold is crossed.
// It does NOT run inside the host app — it survives app termination independently.
class DeviceActivityMonitorExtension: DeviceActivityMonitor {

    private let appGroupID = "group.com.remainfaithful.app"
    private lazy var defaults = UserDefaults(suiteName: appGroupID)

    // MARK: - Schedule callbacks

    override func intervalDidStart(for activity: DeviceActivityName) {
        super.intervalDidStart(for: activity)
        logger.info("Interval started: \(activity.rawValue)")
    }

    override func intervalDidEnd(for activity: DeviceActivityName) {
        super.intervalDidEnd(for: activity)
        logger.info("Interval ended: \(activity.rawValue)")
    }

    // MARK: - Threshold callback

    // Fires when cumulative usage of a monitored app set crosses the event threshold
    // (minimum 1 minute per DeviceActivity API contract). One fire per event per
    // schedule interval — the daily interval resets at midnight when repeats = true.
    override func eventDidReachThreshold(_ event: DeviceActivityEvent.Name,
                                          activity: DeviceActivityName) {
        super.eventDidReachThreshold(event, activity: activity)
        logger.info("Usage threshold reached: event=\(event.rawValue) activity=\(activity.rawValue)")

        let id        = UUID().uuidString
        let timestamp = Date().timeIntervalSince1970
        let summary   = summarize(event: event)

        // Write to pending queue so EventProcessor drains it on next app foreground,
        // guaranteeing delivery even if the direct POST below fails.
        enqueue(id: id, timestamp: timestamp, summary: summary)

        // Stamp the alert time so the main app can surface a broadcast prompt on foreground.
        defaults?.set(timestamp, forKey: "lastAppUsageTimestamp")

        // Attempt direct backend POST — same pattern as the broadcast extension.
        Task { await post(id: id, timestamp: timestamp, summary: summary) }
    }

    // MARK: - Queue (app group UserDefaults)

    // Stored as [[String: Any]] under "pendingRawEvents" — a separate key from
    // "pendingEvents" (which stores DetectedEvent-encoded Data). EventProcessor
    // checks both keys on foreground.
    private func enqueue(id: String, timestamp: TimeInterval, summary: String) {
        var queue = (defaults?.array(forKey: "pendingRawEvents") as? [[String: Any]]) ?? []
        queue.append([
            "id":         id,
            "timestamp":  timestamp,
            "category":   "app_usage",
            "severity":   "concerning",
            "summary":    summary,
            "tier":       4,
            "confidence": 1.0,
            "processed":  false,
        ])
        if queue.count > 50 { queue = Array(queue.suffix(50)) }
        defaults?.set(queue, forKey: "pendingRawEvents")
        logger.info("Queued raw event (depth: \(queue.count))")
    }

    // MARK: - Direct backend POST

    private func post(id: String, timestamp: TimeInterval, summary: String) async {
        guard let apiBase = defaults?.string(forKey: "apiBaseURL"),
              let token   = defaults?.string(forKey: "authToken"),
              let url     = URL(string: apiBase + "/events") else {
            logger.warning("post: missing apiBaseURL or authToken in shared defaults")
            return
        }

        let iso = ISO8601DateFormatter()
        iso.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        let tsStr = iso.string(from: Date(timeIntervalSince1970: timestamp))

        var req = URLRequest(url: url, timeoutInterval: 10)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        guard let body = try? JSONSerialization.data(withJSONObject: [
            "category":  "app_usage",
            "severity":  "concerning",
            "summary":   summary,
            "timestamp": tsStr,
        ]) else { return }
        req.httpBody = body

        do {
            let (_, resp) = try await URLSession.shared.data(for: req)
            let code = (resp as? HTTPURLResponse)?.statusCode ?? 0
            if (200...299).contains(code) {
                logger.info("Event posted: HTTP \(code)")
                // Mark as processed so EventProcessor skips it
                markProcessed(id: id)
            } else {
                logger.warning("Event post HTTP \(code) — left in retry queue")
            }
        } catch {
            logger.warning("Event post failed: \(error.localizedDescription) — left in retry queue")
        }
    }

    private func markProcessed(id: String) {
        guard var queue = defaults?.array(forKey: "pendingRawEvents") as? [[String: Any]] else { return }
        for i in queue.indices where (queue[i]["id"] as? String) == id {
            queue[i]["processed"] = true
        }
        defaults?.set(queue, forKey: "pendingRawEvents")
    }

    // MARK: - Summary text

    private func summarize(event: DeviceActivityEvent.Name) -> String {
        switch event {
        case .appUsage:      return "Monitored app used for 1+ minute"
        case .categoryUsage: return "Monitored app category used for 1+ minute"
        default:             return "Monitored app activity detected"
        }
    }
}

// MARK: - Typed activity and event names

extension DeviceActivityName {
    // The daily 00:00–23:59 repeating schedule window.
    static let daily = DeviceActivityName("com.remainfaithful.daily")
}

extension DeviceActivityEvent.Name {
    // Threshold event for individually selected apps.
    static let appUsage = DeviceActivityEvent.Name("com.remainfaithful.appusage")
    // Threshold event for selected app categories.
    static let categoryUsage = DeviceActivityEvent.Name("com.remainfaithful.categoryusage")
}
