import SwiftUI
import UserNotifications

enum CoachScheduler {
    static func evaluate(now: Date = Date()) {
        guard AuthState.shared.isAuthenticated else { return }
        if let id = AuthState.shared.currentUser?.id {
            CoachStore.prepareForUser(id)
        }
        let decision = CoachEngine.evaluate(
            clock: CoachEngine.clock(from: now, timeZone: .current),
            signals: CoachStore.signals(now: now),
            reply: .pending
        )
        guard decision.showPrompt else { return }
        CoachStore.markPrompt(at: now)
        let content = UNMutableNotificationContent()
        content.title = CoachCopy.checkInTitle
        content.body = CoachCopy.checkInBody
        content.sound = .default
        content.userInfo = ["notification_type": CoachCopy.notificationType]
        let request = UNNotificationRequest(
            identifier: "coach-check-in",
            content: content,
            trigger: nil
        )
        UNUserNotificationCenter.current().add(request)
        AppState.shared.presentCoachCheckIn()
    }
}

/// Local check-in. Not now never asks a partner to reach out.
struct CoachCheckInView: View {
    var onCheckIn: () -> Void
    @Environment(\.dismiss) private var dismiss
    @State private var bodyText = CoachCopy.checkInBody
    @State private var sending = false

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()
            VStack(alignment: .leading, spacing: 20) {
                Text(CoachCopy.checkInTitle)
                    .font(.system(size: 28, weight: .semibold))
                    .foregroundStyle(.white)
                    .fixedSize(horizontal: false, vertical: true)
                    .accessibilityIdentifier("coach-check-in-title")

                Text(bodyText)
                    .font(.body)
                    .foregroundStyle(Color.white.opacity(0.78))
                    .fixedSize(horizontal: false, vertical: true)

                Spacer()

                RFButton(title: CoachCopy.primaryButton, isLoading: sending) {
                    Task { await confirm() }
                }
                .disabled(sending)
                .accessibilityIdentifier("coach-check-in")

                RFButton(title: CoachCopy.secondaryButton, role: .secondary) {
                    decline()
                }
                .disabled(sending)
                .accessibilityIdentifier("coach-not-now")
            }
            .padding(24)
        }
        .task {
            bodyText = await CoachWarmth.checkInBody()
        }
    }

    private func confirm() async {
        sending = true
        let now = Date()
        let decision = CoachEngine.evaluate(
            clock: CoachEngine.clock(from: now, timeZone: .current),
            signals: CoachStore.signals(now: now),
            reply: .checkIn
        )
        CoachStore.markPrompt(at: now)
        CoachStore.markCheckedIn(at: now)
        if decision.sendPartnerNudge {
            await ProtectionAlerts.deliver(type: CoachCopy.alertType, detail: CoachCopy.partnerSentence)
        }
        sending = false
        onCheckIn()
        dismiss()
    }

    private func decline() {
        CoachStore.markPrompt()
        dismiss()
    }
}

struct FutureSelfRecordView: View {
    var showsSkip: Bool = false
    var onDone: () -> Void
    @StateObject private var recorder = FutureSelfRecorder()
    @StateObject private var player = FutureSelfPlayer()
    @State private var partnerOn = false
    @Environment(\.scenePhase) private var scenePhase

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text(FutureSelfCopy.title)
                .font(.system(size: 28, weight: .semibold))
                .foregroundStyle(.white)
                .fixedSize(horizontal: false, vertical: true)
            Text(FutureSelfCopy.body)
                .font(.body)
                .foregroundStyle(Color.white.opacity(0.78))
                .fixedSize(horizontal: false, vertical: true)

            if let message = recorder.message {
                Text(message)
                    .font(.footnote)
                    .foregroundStyle(Color.rfGold)
            }

            RFButton(
                title: recorder.recording && recorder.active == .why ? FutureSelfCopy.stopButton : FutureSelfCopy.recordButton
            ) {
                recorder.toggle(.why)
            }
            .accessibilityIdentifier("future-self-record")

            if FutureSelfStore.exists(.why) {
                RFButton(title: FutureSelfCopy.preview, role: .secondary) {
                    player.play(.why)
                }
            }

            Text(FutureSelfCopy.rememberTitle)
                .font(.system(size: 17, weight: .semibold))
                .foregroundStyle(.white)
                .padding(.top, 8)
            RFButton(
                title: recorder.recording && recorder.active == .remember ? FutureSelfCopy.stopButton : FutureSelfCopy.recordButton,
                role: .secondary
            ) {
                recorder.toggle(.remember)
            }

            if FutureSelfStore.exists(.why) || FutureSelfStore.exists(.remember) {
                RFButton(title: FutureSelfCopy.delete, role: .destructive) {
                    player.stop()
                    FutureSelfStore.wipe()
                    recorder.message = nil
                    recorder.savedClip = nil
                }
            }

            VStack(alignment: .leading, spacing: 8) {
                Text(CoachCopy.partnerTitle)
                    .font(.system(size: 17, weight: .semibold))
                    .foregroundStyle(.white)
                Text(CoachCopy.partnerBody)
                    .font(.footnote)
                    .foregroundStyle(Color.white.opacity(0.75))
                    .fixedSize(horizontal: false, vertical: true)
                RFButton(
                    title: partnerOn ? "Partner check-ins are on" : CoachCopy.partnerButton,
                    role: .secondary
                ) {
                    partnerOn = true
                    CoachStore.setPartnerNudge(true)
                }
                .accessibilityIdentifier("coach-partner-opt-in")
            }
            .padding(.top, 8)

            Spacer(minLength: 12)

            RFButton(title: showsSkip ? "Continue" : "Done") {
                player.stop()
                if recorder.recording { recorder.stop() }
                onDone()
            }
            if showsSkip {
                RFButton(title: "Skip for now", role: .secondary) {
                    player.stop()
                    if recorder.recording { recorder.stopForBackground() }
                    onDone()
                }
                .accessibilityIdentifier("future-self-skip")
            }
        }
        .onAppear {
            partnerOn = CoachStore.signals().partnerNudgeEnabled
        }
        .onChange(of: scenePhase) { _, phase in
            if phase == .background {
                recorder.stopForBackground()
            }
        }
        .onDisappear {
            player.stop()
        }
    }
}
