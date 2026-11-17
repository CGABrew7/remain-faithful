import SwiftUI
import UserNotifications

struct CheckInDraft: Equatable {
    var status: String
    var note: String
    var shieldingOn: Bool
    var sessionShieldOn: Bool
    var scheduleOn: Bool
    var deepScanOn: Bool

    static let statuses: [(id: String, label: String)] = [
        ("steady", "Steady"),
        ("need_support", "Asking for support"),
        ("slipped", "Had a slip"),
    ]

    static func current() -> CheckInDraft {
        let selection = ActivitySelectionManager.shared
        let defaults = UserDefaults(suiteName: SessionShield.appGroupID)
        let broadcasting = defaults?.bool(forKey: SessionShield.broadcastingKey) ?? false
        return CheckInDraft(
            status: "steady",
            note: "",
            shieldingOn: selection.isShieldingEnabled,
            sessionShieldOn: AppLockoutManager.shared.isEnabled,
            scheduleOn: selection.isMonitoringEnabled,
            deepScanOn: broadcasting
        )
    }

    var statusLabel: String {
        Self.statuses.first { $0.id == status }?.label ?? "Steady"
    }
}

/// The user writes the note and taps Send. The report view is never read.
struct CheckInView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var draft = CheckInDraft.current()
    @State private var reminderOn = UserDefaults.standard.bool(forKey: CheckInReminder.enabledKey)
    @State private var sending = false
    @State private var sent = false
    @State private var errorText: String?

    var body: some View {
        NavigationStack {
            ZStack {
                Color.rfNavy.ignoresSafeArea()
                ScrollView {
                    VStack(alignment: .leading, spacing: 18) {
                        Text("Your partner receives the status you pick, the note you write, whether Blocking and Session Shield are on, whether the on-device schedule is on, whether Deep Scan is running, and the time.")
                            .font(.subheadline)
                            .foregroundStyle(Color.white.opacity(0.75))
                            .fixedSize(horizontal: false, vertical: true)

                        VStack(spacing: 0) {
                            ForEach(CheckInDraft.statuses, id: \.id) { item in
                                Button {
                                    draft.status = item.id
                                } label: {
                                    HStack {
                                        Text(item.label)
                                            .foregroundStyle(.white)
                                        Spacer()
                                        if draft.status == item.id {
                                            Image(systemName: "checkmark")
                                                .foregroundStyle(Color.rfGold)
                                        }
                                    }
                                    .padding(.horizontal, 16)
                                    .padding(.vertical, 14)
                                }
                                if item.id != CheckInDraft.statuses.last?.id {
                                    RFRowDivider(inset: 16)
                                }
                            }
                        }
                        .rfGroupedCard()

                        TextField("A short note", text: $draft.note, axis: .vertical)
                            .lineLimit(3...6)
                            .padding(14)
                            .background(Color.white.opacity(0.06))
                            .clipShape(RoundedRectangle(cornerRadius: 12))
                            .foregroundStyle(.white)

                        Toggle("Remind me once a day", isOn: $reminderOn)
                            .tint(Color.rfGold)
                            .foregroundStyle(.white)
                            .onChange(of: reminderOn) { _, on in
                                Task { await CheckInScheduler.setEnabled(on) }
                            }

                        if let errorText {
                            Text(errorText)
                                .font(.footnote)
                                .foregroundStyle(Color(red: 0.95, green: 0.45, blue: 0.40))
                        }
                        if sent {
                            Text("Sent.")
                                .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45))
                        }

                        RFButton(title: sent ? "Close" : "Send check-in", isLoading: sending) {
                            if sent { dismiss() } else { Task { await send() } }
                        }
                        .disabled(sending)
                    }
                    .padding(20)
                }
            }
            .navigationTitle("Check-in")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Close") { dismiss() }
                        .foregroundStyle(Color.rfGold)
                }
            }
            .task {
                if reminderOn {
                    await CheckInScheduler.setEnabled(true)
                }
            }
        }
    }

    private func send() async {
        sending = true
        errorText = nil
        draft.note = String(draft.note.prefix(140))
        do {
            try await APIClient.shared.sendCheckIn(draft)
            sent = true
        } catch {
            errorText = "Still on this phone. Try again when you have a connection."
        }
        sending = false
    }
}

enum CheckInScheduler {
    static func setEnabled(_ on: Bool) async {
        UserDefaults.standard.set(on, forKey: CheckInReminder.enabledKey)
        let center = UNUserNotificationCenter.current()
        center.removePendingNotificationRequests(withIdentifiers: [CheckInReminder.identifier])
        guard on else { return }
        let settings = await center.notificationSettings()
        if settings.authorizationStatus == .notDetermined {
            NotificationService.shared.requestPermission()
        }
        var components = DateComponents()
        components.hour = 20
        components.minute = 0
        let content = UNMutableNotificationContent()
        content.title = "Daily check-in"
        content.body = "Open Remain Faithful to send today's check-in."
        content.sound = .default
        content.userInfo = ["notification_type": CheckInReminder.notificationType]
        let trigger = UNCalendarNotificationTrigger(dateMatching: components, repeats: true)
        let request = UNNotificationRequest(
            identifier: CheckInReminder.identifier,
            content: content,
            trigger: trigger
        )
        try? await center.add(request)
    }
}
