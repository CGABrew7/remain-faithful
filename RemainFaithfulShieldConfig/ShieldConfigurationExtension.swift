import ManagedSettings
import ManagedSettingsUI
import UIKit

// Returns the shield's colors, title, subtitle, and buttons.
// This extension does not open a network connection and does not write
// an app name, bundle id, or domain anywhere.
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

    private func makeConfiguration() -> ShieldConfiguration {
        let navy = UIColor(red: 0.04, green: 0.09, blue: 0.16, alpha: 1.0)
        let gold = UIColor(red: 240.0/255.0, green: 208.0/255.0, blue: 48.0/255.0, alpha: 1.0)
        let dim = UIColor(white: 0.65, alpha: 1.0)

        let sessionShield = UserDefaults(suiteName: "group.com.remainfaithful.app")?
            .bool(forKey: "appLockoutEnabled") ?? false

        let subtitleText = sessionShield
            ? "Start Deep Scan to use this app for a while."
            : "Your partner is standing with you."

        return ShieldConfiguration(
            backgroundBlurStyle: .dark,
            backgroundColor: navy,
            icon: UIImage(named: "AppIcon-1024"),
            title: ShieldConfiguration.Label(text: "Remain Faithful", color: .white),
            subtitle: ShieldConfiguration.Label(text: subtitleText, color: dim),
            primaryButtonLabel: ShieldConfiguration.Label(text: "Open Remain Faithful", color: navy),
            primaryButtonBackgroundColor: gold
        )
    }
}
