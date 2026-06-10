import Foundation
import Combine
import FamilyControls
import ManagedSettings

// Manages the FamilyActivitySelection (apps and categories to block) and the
// ManagedSettingsStore shield. Selection is persisted to the shared app group
// so extensions can read it. The store persists independently of app lifecycle —
// shields remain active across app restarts until explicitly cleared.
final class ActivitySelectionManager: ObservableObject {
    static let shared = ActivitySelectionManager()

    private let store    = ManagedSettingsStore()
    private let defaults = UserDefaults(suiteName: "group.com.remainfaithful.app")

    @Published var selection        = FamilyActivitySelection()
    @Published var isShieldingEnabled: Bool = false

    private init() {
        isShieldingEnabled = defaults?.bool(forKey: "screenTimeShieldingEnabled") ?? false

        if let data  = defaults?.data(forKey: "screenTimeSelection"),
           let saved = try? JSONDecoder().decode(FamilyActivitySelection.self, from: data) {
            selection = saved
        }

        // Re-apply shield on launch — store persists independently but this
        // ensures we're in sync after an OS update or unexpected clear.
        if isShieldingEnabled { applyShielding() }
    }

    // Call after the FamilyActivityPicker updates the selection.
    func selectionDidChange() {
        if let data = try? JSONEncoder().encode(selection) {
            defaults?.set(data, forKey: "screenTimeSelection")
        }
        if isShieldingEnabled { applyShielding() }
    }

    func setShieldingEnabled(_ enabled: Bool) {
        isShieldingEnabled = enabled
        defaults?.set(enabled, forKey: "screenTimeShieldingEnabled")
        if enabled {
            applyShielding()
        } else {
            clearShielding()
        }
    }

    func applyShielding() {
        let apps       = selection.applicationTokens
        let categories = selection.categoryTokens

        store.shield.applications        = apps.isEmpty ? nil : apps
        store.shield.applicationCategories = categories.isEmpty
            ? nil
            : .specific(categories)
    }

    func clearShielding() {
        store.shield.applications          = nil
        store.shield.applicationCategories = nil
    }

    // Human-readable summary of the current selection.
    var selectionSummary: String {
        let apps = selection.applicationTokens.count
        let cats = selection.categoryTokens.count
        guard apps > 0 || cats > 0 else { return "Nothing selected" }
        var parts: [String] = []
        if apps > 0 { parts.append("\(apps) app\(apps == 1 ? "" : "s")") }
        if cats > 0 { parts.append("\(cats) categor\(cats == 1 ? "y" : "ies")") }
        return parts.joined(separator: ", ")
    }
}
