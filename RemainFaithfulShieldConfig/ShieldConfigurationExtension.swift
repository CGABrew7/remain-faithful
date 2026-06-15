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
}
