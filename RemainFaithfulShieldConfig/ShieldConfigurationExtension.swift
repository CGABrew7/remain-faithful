import ManagedSettings
import ManagedSettingsUI
import UIKit

// Provides custom branding for the iOS Screen Time block screen.
// Replaces Apple's generic "This app is not available" UI with Remain Faithful
// branding whenever ManagedSettingsStore shields an app or category.
//
// When App Lockout (Feature 1) is active the subtitle instructs the user to
// open Remain Faithful and start Deep Scan to lift the shield. This configuration
// is read directly from the shared app group UserDefaults — no entitlement change
// is needed because the ShieldConfig extension already has family-controls.
class ShieldConfigurationExtension: ShieldConfigurationDataSource {

    override func configuration(shielding application: Application) -> ShieldConfiguration {
        makeConfiguration()
    }

    override func configuration(shielding application: Application,
                                in category: ActivityCategory) -> ShieldConfiguration {
        makeConfiguration()
    }

    override func configuration(shielding webDomain: WebDomain) -> ShieldConfiguration {
        makeConfiguration()
    }

    override func configuration(shielding webDomain: WebDomain,
                                in category: ActivityCategory) -> ShieldConfiguration {
        makeConfiguration()
    }

    // MARK: - Private

    private func makeConfiguration() -> ShieldConfiguration {
        noteBlockedAppOpened()
        let navy  = UIColor(red: 0.04, green: 0.09, blue: 0.16, alpha: 1.0)
        let gold  = UIColor(red: 0.87, green: 0.62, blue: 0.24, alpha: 1.0)
        let dim   = UIColor(white: 0.65, alpha: 1.0)

        let lockoutEnabled = UserDefaults(suiteName: "group.com.remainfaithful.app")?
            .bool(forKey: "appLockoutEnabled") ?? false

        let subtitleText = lockoutEnabled
            ? "Open Remain Faithful and start Deep Scan to unlock."
            : "Your partner is standing with you."

        return ShieldConfiguration(
            backgroundBlurStyle: .dark,
            backgroundColor: navy,
            icon: UIImage(named: "AppIcon-1024"),
            title: ShieldConfiguration.Label(
                text: "Remain Faithful",
                color: .white
            ),
            subtitle: ShieldConfiguration.Label(
                text: subtitleText,
                color: dim
            ),
            primaryButtonLabel: ShieldConfiguration.Label(
                text: "I need help",
                color: navy
            ),
            primaryButtonBackgroundColor: gold
        )
    }

    /// Opening a blocked app tells partners through the same alert path as other app activity.
    /// The note is metadata only. It does not include the app name.
    private func noteBlockedAppOpened() {
        let defaults = UserDefaults(suiteName: "group.com.remainfaithful.app")
        let now = Date().timeIntervalSince1970
        let last = defaults?.double(forKey: "lastBlockedOpenAlert") ?? 0
        if last > 0, now - last < 120 { return }
        defaults?.set(now, forKey: "lastBlockedOpenAlert")

        let id = UUID().uuidString
        let account = defaults?.integer(forKey: "authUserID") ?? 0
        var row: [String: Any] = [
            "id": id,
            "timestamp": now,
            "category": "app_usage",
            "severity": "concerning",
            "summary": "Blocked app opened",
            "tier": 4,
            "confidence": 1.0,
            "processed": false,
        ]
        if account > 0 { row["accountID"] = account }
        var queue = (defaults?.array(forKey: "pendingRawEvents") as? [[String: Any]]) ?? []
        queue.append(row)
        if queue.count > 50 { queue = Array(queue.suffix(50)) }
        defaults?.set(queue, forKey: "pendingRawEvents")
        Task { await postBlockedOpen(id: id, timestamp: now, defaults: defaults) }
    }

    private func postBlockedOpen(id: String, timestamp: TimeInterval, defaults: UserDefaults?) async {
        guard let apiBase = defaults?.string(forKey: "apiBaseURL"),
              let token = defaults?.string(forKey: "authToken"),
              let url = URL(string: apiBase + "/events") else { return }
        let iso = ISO8601DateFormatter()
        iso.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        var req = URLRequest(url: url, timeoutInterval: 10)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        req.httpBody = try? JSONSerialization.data(withJSONObject: [
            "category": "app_usage",
            "severity": "concerning",
            "summary": "Blocked app opened",
            "timestamp": iso.string(from: Date(timeIntervalSince1970: timestamp)),
        ])
        guard req.httpBody != nil else { return }
        if let (_, resp) = try? await URLSession.shared.data(for: req),
           let http = resp as? HTTPURLResponse, (200...299).contains(http.statusCode) {
            guard var queue = defaults?.array(forKey: "pendingRawEvents") as? [[String: Any]] else { return }
            for i in queue.indices where (queue[i]["id"] as? String) == id {
                queue[i]["processed"] = true
            }
            defaults?.set(queue, forKey: "pendingRawEvents")
        }
    }
}
