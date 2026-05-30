import SwiftUI

struct AlertDetailView: View {
    let event: ActivityEvent

    @Environment(\.dismiss) private var dismiss
    @State private var isDiscussed = false

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            VStack(spacing: 0) {
                navBar

                ScrollView(showsIndicators: false) {
                    VStack(spacing: 22) {
                        heroHeader
                        detectedSection
                        conversationSection
                    }
                    .padding(.horizontal, 20)
                    .padding(.top, 4)
                    .padding(.bottom, 32)
                }

                markDiscussedButton
                    .padding(.horizontal, 20)
                    .padding(.top, 12)
                    .padding(.bottom, 40)
            }
        }
        .navigationBarHidden(true)
    }

    // MARK: - Nav bar

    private var navBar: some View {
        HStack {
            Button { dismiss() } label: {
                HStack(spacing: 6) {
                    Image(systemName: "chevron.left")
                        .font(.system(size: 15, weight: .semibold))
                    Text("Recent Flags")
                        .font(.system(size: 15, weight: .medium))
                }
                .foregroundStyle(Color.rfGold)
            }
            Spacer()
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 14)
    }

    // MARK: - Hero header

    private var heroHeader: some View {
        VStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(event.category.tint.opacity(0.13))
                    .frame(width: 88, height: 88)
                Image(systemName: event.category.icon)
                    .font(.system(size: 34))
                    .foregroundStyle(event.category.tint)
            }

            VStack(spacing: 10) {
                Text(event.category.rawValue)
                    .font(.system(size: 28, weight: .bold, design: .serif))
                    .foregroundStyle(.white)

                Text(event.severity.rawValue.uppercased())
                    .font(.system(size: 11, weight: .bold))
                    .kerning(0.8)
                    .foregroundStyle(event.severity.color)
                    .padding(.horizontal, 14)
                    .padding(.vertical, 5)
                    .background(Capsule().fill(event.severity.color.opacity(0.14)))
            }
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 6)
    }

    // MARK: - Detected section

    private var detectedSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            sectionLabel("DETECTED")

            VStack(alignment: .leading, spacing: 14) {
                HStack(spacing: 10) {
                    Image(systemName: "clock.fill")
                        .font(.system(size: 13))
                        .foregroundStyle(Color.white.opacity(0.38))
                    Text(event.fullTimestamp)
                        .font(.system(size: 14))
                        .foregroundStyle(Color.white.opacity(0.65))
                    Spacer()
                    Text(event.timeLabel)
                        .font(.system(size: 12))
                        .foregroundStyle(Color.white.opacity(0.32))
                }

                Divider().overlay(Color.white.opacity(0.07))

                Text(event.extendedDescription)
                    .font(.system(size: 15))
                    .foregroundStyle(Color.white.opacity(0.85))
                    .lineSpacing(5)
                    .fixedSize(horizontal: false, vertical: true)
            }
            .padding(16)
            .background(detailCard)
        }
    }

    // MARK: - Conversation starter section

    private var conversationSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            sectionLabel("CONVERSATION STARTER")

            VStack(alignment: .leading, spacing: 16) {
                HStack(spacing: 10) {
                    Image(systemName: "bubble.left.and.bubble.right.fill")
                        .font(.system(size: 13))
                        .foregroundStyle(Color.rfGold.opacity(0.70))
                    Text("Suggested talking points for your accountability partner")
                        .font(.system(size: 13))
                        .foregroundStyle(Color.white.opacity(0.48))
                        .fixedSize(horizontal: false, vertical: true)
                }

                Divider().overlay(Color.white.opacity(0.07))

                VStack(alignment: .leading, spacing: 14) {
                    ForEach(Array(event.category.conversationStarters.enumerated()), id: \.offset) { _, starter in
                        HStack(alignment: .top, spacing: 12) {
                            Circle()
                                .fill(Color.rfGold.opacity(0.55))
                                .frame(width: 5, height: 5)
                                .padding(.top, 8)
                            Text(starter)
                                .font(.system(size: 14))
                                .foregroundStyle(Color.white.opacity(0.82))
                                .lineSpacing(4)
                                .fixedSize(horizontal: false, vertical: true)
                        }
                    }
                }
            }
            .padding(16)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Color.rfGold.opacity(0.06))
                    .overlay(RoundedRectangle(cornerRadius: 16)
                        .stroke(Color.rfGold.opacity(0.18), lineWidth: 1))
            )
        }
    }

    // MARK: - Mark as discussed button

    private var markDiscussedButton: some View {
        let green = Color(red: 0.20, green: 0.78, blue: 0.45)
        let tint  = isDiscussed ? green : Color.rfGold

        return Button {
            withAnimation(.spring(response: 0.35, dampingFraction: 0.7)) {
                isDiscussed = true
            }
        } label: {
            HStack(spacing: 10) {
                Image(systemName: isDiscussed ? "checkmark.circle.fill" : "checkmark.circle")
                    .font(.system(size: 18))
                Text(isDiscussed ? "Marked as Discussed" : "Mark as Discussed")
                    .font(.system(size: 16, weight: .semibold))
            }
            .foregroundStyle(tint)
            .frame(maxWidth: .infinity)
            .frame(height: 54)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(tint.opacity(0.11))
                    .overlay(RoundedRectangle(cornerRadius: 14)
                        .stroke(tint.opacity(0.32), lineWidth: 1.5))
            )
        }
        .disabled(isDiscussed)
        .animation(.easeInOut(duration: 0.25), value: isDiscussed)
    }

    // MARK: - Helpers

    private func sectionLabel(_ text: String) -> some View {
        Text(text)
            .font(.system(size: 10, weight: .bold))
            .foregroundStyle(Color.rfGold.opacity(0.75))
            .kerning(1.4)
    }

    private var detailCard: some View {
        RoundedRectangle(cornerRadius: 16)
            .fill(Color.white.opacity(0.055))
            .overlay(RoundedRectangle(cornerRadius: 16)
                .stroke(Color.white.opacity(0.08), lineWidth: 1))
    }
}

#Preview {
    NavigationStack {
        AlertDetailView(event: ActivityEvent(
            category: .adultContent,
            severity: .high,
            description: "Flagged during browsing session",
            minutesAgo: 23
        ))
    }
}
