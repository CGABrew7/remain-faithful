import Darwin
import SwiftUI
import UIKit

// MARK: - What the form sends

enum FeedbackKind: String, CaseIterable, Identifiable {
    case idea = "Idea"
    case improvement = "Improvement"
    case problem = "Problem"
    var id: String { rawValue }
}

struct FeedbackSubmission: Equatable {
    var type: FeedbackKind
    var summary: String
    var details: String
    var includeDevice: Bool
    var appVersion: String
    var deviceModel: String
    var osVersion: String

    var trimmedSummary: String {
        summary.trimmingCharacters(in: .whitespacesAndNewlines)
    }
    var trimmedDetails: String {
        details.trimmingCharacters(in: .whitespacesAndNewlines)
    }

    var canSend: Bool {
        !trimmedSummary.isEmpty && trimmedSummary.count <= 120
            && !trimmedDetails.isEmpty && trimmedDetails.count <= 4000
    }

    /// Exactly the device line shown under the toggle, and only sent when it is on.
    var deviceCaption: String {
        "\(appVersion), \(deviceModel), \(osVersion)"
    }

    func apiFields() -> [String: String] {
        var fields = [
            "type": type.rawValue,
            "summary": trimmedSummary,
            "details": trimmedDetails,
        ]
        if includeDevice {
            fields["app_version"] = appVersion
            fields["device_model"] = deviceModel
            fields["os_version"] = osVersion
        }
        return fields
    }

    func mailtoURL() -> URL? {
        var body = trimmedDetails
        if includeDevice {
            body += "\n\n" + deviceCaption
        }
        var components = URLComponents()
        components.scheme = "mailto"
        components.path = "support@remainfaithful.com"
        components.queryItems = [
            URLQueryItem(name: "subject", value: "[\(type.rawValue)] \(trimmedSummary)"),
            URLQueryItem(name: "body", value: body),
        ]
        return components.url
    }
}

enum AppDeviceLabel {
    static var appVersion: String {
        let version = Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "1.0"
        let build = Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") as? String ?? "8"
        return "Remain Faithful \(version) (\(build))"
    }

    static var modelIdentifier: String {
        var system = utsname()
        uname(&system)
        return withUnsafePointer(to: &system.machine) {
            $0.withMemoryRebound(to: CChar.self, capacity: 1) {
                String(cString: $0)
            }
        }
    }

    static var osVersion: String {
        "iOS \(UIDevice.current.systemVersion)"
    }
}

// MARK: - Sheet

struct FeedbackSheet: View {
    let replyEmail: String
    @Environment(\.dismiss) private var dismiss
    @Environment(\.openURL) private var openURL

    @State private var type: FeedbackKind = .idea
    @State private var summary = ""
    @State private var details = ""
    @State private var includeDevice = false
    @State private var isSending = false
    @State private var errorText: String?
    @State private var didSend = false
    @ScaledMetric(relativeTo: .body) private var bodySize: CGFloat = 15

    private var submission: FeedbackSubmission {
        FeedbackSubmission(
            type: type,
            summary: summary,
            details: details,
            includeDevice: includeDevice,
            appVersion: AppDeviceLabel.appVersion,
            deviceModel: AppDeviceLabel.modelIdentifier,
            osVersion: AppDeviceLabel.osVersion
        )
    }

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                header
                Divider().overlay(Color.white.opacity(0.08))
                if didSend {
                    thankYou
                } else {
                    form
                }
            }
        }
    }

    private var header: some View {
        HStack(alignment: .top) {
            VStack(alignment: .leading, spacing: 4) {
                Text("Send ideas or report a problem")
                    .font(.system(.title3, design: .serif).weight(.bold))
                    .foregroundStyle(.white)
                    .fixedSize(horizontal: false, vertical: true)
                Text(replyLine)
                    .font(.system(size: bodySize - 2))
                    .foregroundStyle(Color.white.opacity(0.55))
                    .fixedSize(horizontal: false, vertical: true)
            }
            Spacer(minLength: 12)
            Button { dismiss() } label: {
                Image(systemName: "xmark")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(Color.white.opacity(0.55))
                    .padding(10)
                    .background(Circle().fill(Color.white.opacity(0.09)))
            }
            .accessibilityLabel("Close")
        }
        .padding(.horizontal, 24)
        .padding(.top, 22)
        .padding(.bottom, 16)
    }

    private var replyLine: String {
        let email = replyEmail.trimmingCharacters(in: .whitespacesAndNewlines)
        if email.isEmpty { return "We'll reply to your account email." }
        return "We'll reply to \(email)."
    }

    private var form: some View {
        ScrollView(showsIndicators: false) {
            VStack(alignment: .leading, spacing: 18) {
                VStack(alignment: .leading, spacing: 8) {
                    Text("What kind?")
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(Color.rfGold)
                    Picker("What kind?", selection: $type) {
                        ForEach(FeedbackKind.allCases) { kind in
                            Text(kind.rawValue).tag(kind)
                        }
                    }
                    .pickerStyle(.segmented)
                    .accessibilityLabel("What kind of feedback?")
                }

                VStack(alignment: .leading, spacing: 8) {
                    Text("In a few words")
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(.white)
                    TextField("In a few words", text: $summary, axis: .vertical)
                        .lineLimit(1...3)
                        .textInputAutocapitalization(.sentences)
                        .foregroundStyle(.white)
                        .padding(14)
                        .background(fieldBackground)
                        .accessibilityLabel("In a few words")
                        .onChange(of: summary) { _, newValue in
                            if newValue.count > 120 {
                                summary = String(newValue.prefix(120))
                            }
                        }
                }

                VStack(alignment: .leading, spacing: 8) {
                    Text("Tell us more")
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(.white)
                    TextEditor(text: $details)
                        .scrollContentBackground(.hidden)
                        .frame(minHeight: 140)
                        .foregroundStyle(.white)
                        .font(.system(size: bodySize))
                        .padding(10)
                        .background(fieldBackground)
                        .accessibilityLabel("Tell us more")
                        .onChange(of: details) { _, newValue in
                            if newValue.count > 4000 {
                                details = String(newValue.prefix(4000))
                            }
                        }
                }

                Toggle(isOn: $includeDevice) {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Include app version and iPhone model")
                            .font(.subheadline.weight(.semibold))
                            .foregroundStyle(.white)
                        Text(submission.deviceCaption)
                            .font(.footnote)
                            .foregroundStyle(Color.white.opacity(0.55))
                            .fixedSize(horizontal: false, vertical: true)
                    }
                }
                .tint(Color.rfGold)
                .accessibilityHint("Off by default. Turns on only the version, model, and iOS version shown here.")

                if let errorText {
                    Text(errorText)
                        .font(.subheadline)
                        .foregroundStyle(Color(red: 0.95, green: 0.55, blue: 0.45))
                        .fixedSize(horizontal: false, vertical: true)
                    Button {
                        if let url = submission.mailtoURL() { openURL(url) }
                    } label: {
                        Text("Email support@remainfaithful.com")
                            .font(.subheadline.weight(.semibold))
                            .foregroundStyle(Color.rfGold)
                    }
                    .accessibilityLabel("Email support")
                }

                Button(action: send) {
                    HStack {
                        if isSending { ProgressView().tint(Color.rfNavy) }
                        Text(isSending ? "Sending…" : "Send")
                            .font(.body.weight(.semibold))
                    }
                    .foregroundStyle(submission.canSend ? Color.rfNavy : Color.white.opacity(0.35))
                    .frame(maxWidth: .infinity)
                    .frame(height: 52)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(submission.canSend ? Color.rfGold : Color.white.opacity(0.08))
                    )
                }
                .disabled(!submission.canSend || isSending)
                .accessibilityLabel("Send")

                Button {
                    if let url = URL(string: "mailto:support@remainfaithful.com") {
                        openURL(url)
                    }
                } label: {
                    Text("Prefer email? support@remainfaithful.com")
                        .font(.footnote)
                        .foregroundStyle(Color.white.opacity(0.55))
                        .frame(maxWidth: .infinity)
                }
                .accessibilityLabel("Prefer email? support at remain faithful dot com")
            }
            .padding(24)
        }
    }

    private var thankYou: some View {
        VStack(spacing: 12) {
            Spacer()
            Text("Thank you. We got it.")
                .font(.system(.title3, design: .serif).weight(.bold))
                .foregroundStyle(Color.rfGold)
                .multilineTextAlignment(.center)
            Text(replyLine)
                .font(.system(size: bodySize))
                .foregroundStyle(Color.white.opacity(0.7))
                .multilineTextAlignment(.center)
            Spacer()
        }
        .padding(24)
        .accessibilityElement(children: .combine)
    }

    private var fieldBackground: some View {
        RoundedRectangle(cornerRadius: 14)
            .fill(Color.white.opacity(0.07))
            .overlay(
                RoundedRectangle(cornerRadius: 14)
                    .stroke(Color.white.opacity(0.10), lineWidth: 1)
            )
    }

    private func send() {
        let draft = submission
        guard draft.canSend, !isSending else { return }
        isSending = true
        errorText = nil
        Task {
            do {
                try await APIClient.shared.sendFeedback(draft.apiFields())
                await MainActor.run {
                    isSending = false
                    didSend = true
                }
                try? await Task.sleep(for: .seconds(1.2))
                await MainActor.run { dismiss() }
            } catch {
                await MainActor.run {
                    isSending = false
                    errorText = "Couldn't send. Try again or email support@remainfaithful.com"
                }
            }
        }
    }
}

struct JSONStringMap: Encodable {
    let values: [String: String]

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: DynamicKey.self)
        for (key, value) in values {
            try container.encode(value, forKey: DynamicKey(key))
        }
    }

    struct DynamicKey: CodingKey {
        var stringValue: String
        var intValue: Int? { nil }
        init(_ value: String) { stringValue = value }
        init?(stringValue: String) { self.stringValue = stringValue }
        init?(intValue: Int) { nil }
    }
}
