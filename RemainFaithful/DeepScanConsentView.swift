import SwiftUI
import ReplayKit

/// Names the recipient and the exact fields before the first Deep Scan or Session Shield use.
struct DeepScanConsentView: View {
    let onAccept: () -> Void
    let onDecline: () -> Void

    @State private var recipientLine = SharingRecipients.line(partners: [], groups: [])

    var body: some View {
        NavigationStack {
            ZStack {
                Color.rfNavy.ignoresSafeArea()
                ScrollView {
                    VStack(alignment: .leading, spacing: 16) {
                        Text("Deep Scan sharing")
                            .font(.system(size: 28, weight: .semibold))
                            .foregroundStyle(.white)

                        Text("While you run Deep Scan, \(recipientLine) can receive:")
                            .font(.body)
                            .foregroundStyle(Color.white.opacity(0.82))
                            .fixedSize(horizontal: false, vertical: true)

                        VStack(alignment: .leading, spacing: 10) {
                            field("Category: adult content, gambling, violence, or self-harm")
                            field("Severity")
                            field("Confidence: low, medium, or high")
                            field("Time")
                            field("A notice when Deep Scan stops")
                        }

                        Text("You start the broadcast from the system picker. iOS shows the recording indicator for the whole session. Session Shield keeps your chosen apps shielded until that broadcast is running, and shields them again when it stops.")
                            .font(.subheadline)
                            .foregroundStyle(Color.white.opacity(0.75))
                            .fixedSize(horizontal: false, vertical: true)

                        RFButton(title: "I agree") {
                            DeepScanConsent.accept()
                            onAccept()
                        }
                        RFButton(title: "Not now", role: .secondary, action: onDecline)
                    }
                    .padding(24)
                }
            }
            .task { await loadRecipients() }
        }
    }

    private func field(_ text: String) -> some View {
        HStack(alignment: .top, spacing: 10) {
            Image(systemName: "circle.fill")
                .font(.system(size: 6))
                .foregroundStyle(Color.rfGold)
                .padding(.top, 7)
            Text(text)
                .font(.body)
                .foregroundStyle(.white)
                .fixedSize(horizontal: false, vertical: true)
        }
    }

    private func loadRecipients() async {
        var partners: [String] = []
        var groups: [String] = []
        if let rels = try? await APIClient.shared.listRelationships() {
            partners = rels.filter { $0.status == "accepted" }.map { $0.partner.name }
        }
        if let mine = try? await APIClient.shared.listMyGroups() {
            groups = mine.map(\.name)
        }
        recipientLine = SharingRecipients.line(partners: partners, groups: groups)
    }
}

struct SystemBroadcastPicker: UIViewRepresentable {
    func makeUIView(context: Context) -> RPSystemBroadcastPickerView {
        let picker = RPSystemBroadcastPickerView(frame: CGRect(x: 0, y: 0, width: 52, height: 52))
        picker.preferredExtension = "com.remainfaithful.app.broadcast"
        picker.showsMicrophoneButton = false
        return picker
    }

    func updateUIView(_ uiView: RPSystemBroadcastPickerView, context: Context) {}
}
