import SwiftUI

struct PanicView: View {
    @Environment(\.dismiss)  private var dismiss
    @Environment(\.openURL)  private var openURL

    @State private var alertSent      = false
    @State private var pulse: CGFloat = 1.0
    @State private var partner:        RemoteRelationship? = nil
    @State private var partnerLoaded   = false

    var body: some View {
        ZStack {
            Color(red: 0.04, green: 0.07, blue: 0.18).ignoresSafeArea()

            VStack(spacing: 0) {
                dismissBar
                ScrollView(showsIndicators: false) {
                    VStack(spacing: 28) {
                        heroSection
                        verseSection
                        partnerSection
                        groupAlertButton
                    }
                    .padding(.horizontal, 24)
                    .padding(.top, 8)
                    .padding(.bottom, 32)
                }
                closeButton
                    .padding(.horizontal, 24)
                    .padding(.top, 12)
                    .padding(.bottom, 40)
            }
        }
        .task { await loadPartner() }
    }

    @MainActor
    private func loadPartner() async {
        guard APIClient.shared.isAuthenticated else {
            partnerLoaded = true
            return
        }
        if let rels = try? await APIClient.shared.listRelationships() {
            let accepted = rels.filter { $0.status == "accepted" }
            partner = accepted.first(where: { $0.isPrimary }) ?? accepted.first
        }
        partnerLoaded = true
    }

    // MARK: - Dismiss bar

    private var dismissBar: some View {
        HStack {
            Spacer()
            Button { dismiss() } label: {
                Image(systemName: "xmark")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(Color.white.opacity(0.40))
                    .padding(10)
                    .background(Circle().fill(Color.white.opacity(0.08)))
            }
        }
        .padding(.horizontal, 24)
        .padding(.top, 16)
        .padding(.bottom, 8)
    }

    // MARK: - Hero

    private var heroSection: some View {
        VStack(spacing: 22) {
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.07))
                    .frame(width: 140, height: 140)
                    .scaleEffect(pulse)
                Circle()
                    .fill(Color.rfGold.opacity(0.11))
                    .frame(width: 100, height: 100)
                    .scaleEffect(pulse * 0.92)
                Image(systemName: "cross.fill")
                    .font(.system(size: 40, weight: .medium))
                    .foregroundStyle(Color.rfGold)
            }
            .onAppear {
                withAnimation(.easeInOut(duration: 2.0).repeatForever(autoreverses: true)) {
                    pulse = 1.20
                }
            }

            VStack(spacing: 10) {
                Text("You Are Not Alone")
                    .font(.system(size: 32, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                    .multilineTextAlignment(.center)

                Text("Your accountability partners want to hear from you right now. Reaching out is a sign of strength, not weakness.")
                    .font(.system(size: 15))
                    .foregroundStyle(Color.white.opacity(0.58))
                    .multilineTextAlignment(.center)
                    .lineSpacing(5)
            }
        }
        .padding(.top, 8)
    }

    // MARK: - Scripture

    private var verseSection: some View {
        VStack(alignment: .leading, spacing: 10) {
            Text("HOLD ON TO THIS")
                .font(.system(size: 10, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.4)

            Text("""
"No temptation has overtaken you except what is common to mankind. \
And God is faithful; he will not let you be tempted beyond what you can bear. \
But when you are tempted, he will also provide a way out so that you can endure it."
""")
                .font(.system(size: 14, weight: .medium, design: .serif))
                .foregroundStyle(.white)
                .lineSpacing(6)
                .fixedSize(horizontal: false, vertical: true)

            Text("— 1 Corinthians 10:13")
                .font(.system(size: 12))
                .foregroundStyle(Color.white.opacity(0.45))
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.rfGold.opacity(0.07))
                .overlay(RoundedRectangle(cornerRadius: 16)
                    .stroke(Color.rfGold.opacity(0.20), lineWidth: 1))
        )
    }

    // MARK: - Partner

    private var partnerSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("YOUR ACCOUNTABILITY PARTNER")
                .font(.system(size: 10, weight: .bold))
                .foregroundStyle(Color.white.opacity(0.38))
                .kerning(1.2)

            if let p = partner {
                let name = p.partner.name
                let initials = name.components(separatedBy: " ")
                    .compactMap(\.first).prefix(2).map(String.init).joined().uppercased()

                HStack(spacing: 16) {
                    ZStack {
                        Circle()
                            .fill(Color.rfGold.opacity(0.14))
                            .frame(width: 52, height: 52)
                        Text(initials)
                            .font(.system(size: 18, weight: .semibold))
                            .foregroundStyle(Color.rfGold)
                    }

                    VStack(alignment: .leading, spacing: 4) {
                        Text(name)
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundStyle(.white)
                        Text("Accountability Partner")
                            .font(.system(size: 12))
                            .foregroundStyle(Color.white.opacity(0.42))
                    }

                    Spacer()
                }
                .padding(16)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(Color.white.opacity(0.055))
                        .overlay(RoundedRectangle(cornerRadius: 14)
                            .stroke(Color.white.opacity(0.08), lineWidth: 1))
                )
            } else if partnerLoaded {
                Text("Add an accountability partner to use this feature")
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.55))
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: .infinity)
                    .padding(16)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(Color.white.opacity(0.055))
                            .overlay(RoundedRectangle(cornerRadius: 14)
                                .stroke(Color.white.opacity(0.08), lineWidth: 1))
                    )
            } else {
                ProgressView()
                    .tint(Color.rfGold)
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 16)
            }
        }
    }

    // MARK: - Group alert

    private var groupAlertButton: some View {
        let sent = alertSent
        let blue  = Color(red: 0.28, green: 0.56, blue: 0.95)
        let green = Color(red: 0.20, green: 0.78, blue: 0.45)
        let tint  = sent ? green : blue

        return Button {
            withAnimation(.spring(response: 0.35, dampingFraction: 0.7)) {
                alertSent = true
            }
            Task { try? await APIClient.shared.sendPanicAlert() }
        } label: {
            HStack(spacing: 10) {
                Image(systemName: sent ? "checkmark.circle.fill" : "bell.badge.fill")
                    .font(.system(size: 16))
                Text(sent ? "Alert Sent to Group" : "Send Alert to Group")
                    .font(.system(size: 16, weight: .semibold))
            }
            .foregroundStyle(tint)
            .frame(maxWidth: .infinity)
            .frame(height: 54)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(tint.opacity(0.12))
                    .overlay(RoundedRectangle(cornerRadius: 14)
                        .stroke(tint.opacity(0.35), lineWidth: 1.5))
            )
        }
        .disabled(sent)
        .animation(.easeInOut(duration: 0.25), value: sent)
    }

    // MARK: - Close

    private var closeButton: some View {
        Button { dismiss() } label: {
            Text("I'm okay, close this")
                .font(.system(size: 15, weight: .semibold))
                .foregroundStyle(Color.white.opacity(0.60))
                .frame(maxWidth: .infinity)
                .frame(height: 52)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(Color.white.opacity(0.07))
                        .overlay(RoundedRectangle(cornerRadius: 14)
                            .stroke(Color.white.opacity(0.12), lineWidth: 1))
                )
        }
    }
}

#Preview {
    PanicView()
}
