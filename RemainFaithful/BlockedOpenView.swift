import SwiftUI

/// Shown after the shield button. Nothing is sent until Tell my partner.
struct BlockedOpenView: View {
    let openedAt: Date
    let onClose: () -> Void

    @State private var sending = false
    @State private var sent = false
    @State private var errorText: String?
    @State private var voicePlayed = false
    @State private var voiceHidden = false
    @StateObject private var player = FutureSelfPlayer()

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()
            VStack(alignment: .leading, spacing: 20) {
                Text("A blocked app was opened")
                    .font(.system(size: 28, weight: .semibold))
                    .foregroundStyle(.white)
                    .fixedSize(horizontal: false, vertical: true)

                Text("Tell my partner sends the time below and the words \"A blocked app was opened.\"")
                    .font(.body)
                    .foregroundStyle(Color.white.opacity(0.78))
                    .fixedSize(horizontal: false, vertical: true)

                Text(openedAt.formatted(date: .abbreviated, time: .shortened))
                    .font(.system(size: 17, weight: .medium))
                    .foregroundStyle(Color.rfGold)

                if FutureSelfStore.exists(.why), !voiceHidden {
                    Text(FutureSelfCopy.playTitle)
                        .font(.system(size: 17, weight: .semibold))
                        .foregroundStyle(.white)
                    HStack(spacing: 12) {
                        RFButton(title: FutureSelfCopy.playAgain, role: .secondary) {
                            player.play(.why)
                        }
                        RFButton(title: FutureSelfCopy.skip, role: .secondary) {
                            player.stop()
                            voiceHidden = true
                        }
                        .accessibilityIdentifier("future-self-skip-playback")
                    }
                }

                if let errorText {
                    Text(errorText)
                        .font(.footnote)
                        .foregroundStyle(Color(red: 0.95, green: 0.45, blue: 0.40))
                }
                if sent {
                    Text("Your partner has this notice.")
                        .font(.body)
                        .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45))
                }

                Spacer()

                RFButton(title: sent ? "Done" : "Tell my partner", isLoading: sending) {
                    if sent {
                        finish()
                    } else {
                        Task { await send() }
                    }
                }
                .disabled(sending)

                RFButton(title: "Not now", role: .secondary) {
                    finish()
                }
                .disabled(sending)
            }
            .padding(24)
        }
        .onAppear {
            guard FutureSelfPlayback.shouldAutoPlay(hasClip: FutureSelfStore.exists(.why), alreadyPlayedThisOpen: voicePlayed) else { return }
            voicePlayed = true
            player.play(.why)
        }
        .onDisappear { player.stop() }
    }

    private func send() async {
        sending = true
        errorText = nil
        do {
            try await APIClient.shared.sendBlockedOpen(at: openedAt)
            BlockedOpenRequest.clear()
            sent = true
        } catch {
            errorText = "The notice is still on this phone. Try again when you have a connection."
        }
        sending = false
    }

    private func finish() {
        BlockedOpenRequest.clear()
        onClose()
    }
}
