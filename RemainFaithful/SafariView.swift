import SafariServices
import SwiftUI

/// Thin UIViewControllerRepresentable wrapper around SFSafariViewController.
/// Used to present Stripe Checkout and other external URLs inside the app.
struct SafariView: UIViewControllerRepresentable {
    let url: URL

    func makeUIViewController(context: Context) -> SFSafariViewController {
        let config = SFSafariViewController.Configuration()
        config.entersReaderIfAvailable = false
        let vc = SFSafariViewController(url: url, configuration: config)
        vc.preferredBarTintColor     = UIColor(red: 0.05, green: 0.09, blue: 0.22, alpha: 1)
        vc.preferredControlTintColor = UIColor(red: 0.83, green: 0.68, blue: 0.28, alpha: 1)
        return vc
    }

    func updateUIViewController(_ vc: SFSafariViewController, context: Context) {}
}
