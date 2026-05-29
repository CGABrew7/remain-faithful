import SwiftUI

// MARK: - Models

enum AccountabilityHealth {
    case strong
    case watchful
    case struggling

    var color: Color {
        switch self {
        case .strong:     Color(red: 0.20, green: 0.78, blue: 0.45)
        case .watchful:   Color(red: 0.95, green: 0.62, blue: 0.10)
        case .struggling: Color(red: 0.90, green: 0.25, blue: 0.30)
        }
    }

    var label: String {
        switch self {
        case .strong:     "Strong"
        case .watchful:   "Watchful"
        case .struggling: "Struggling"
        }
    }
}

struct GroupMember: Identifiable {
    let id = UUID()
    let name: String
    let streak: Int
    let health: AccountabilityHealth
}

// MARK: - Mock data

private let groupName = "Iron Brotherhood"

private let sampleMembers: [GroupMember] = [
    .init(name: "James Bishop",  streak: 21, health: .strong),
    .init(name: "Marcus Cole",   streak: 7,  health: .watchful),
    .init(name: "David Torres",  streak: 3,  health: .struggling),
    .init(name: "Nathan Wells",  streak: 14, health: .strong),
    .init(name: "Chris Hammond", streak: 0,  health: .struggling),
]

private let covenantText = """
We, the members of this accountability group, covenant together before God and one another to:

1. PURSUE PURITY — We commit to actively pursue sexual purity in thought, word, and deed, recognizing that our bodies are temples of the Holy Spirit.

2. WALK IN HONESTY — We will be completely truthful about our struggles, failures, and victories, refusing to hide in shame or deceive those who stand with us.

3. PRAY FAITHFULLY — We commit to pray regularly for each member of this group, lifting one another up in times of temptation and weakness.

4. RESPOND QUICKLY — When a brother sends an alert, we will respond promptly with encouragement, prayer, and accountability — no judgment, only grace.

5. GUARD ONE ANOTHER — We take responsibility for one another's spiritual health. We will reach out when we see a brother struggling, not waiting for them to ask.

6. HOLD CONFIDENCE — What is shared in this group stays in this group. We maintain trust by protecting each other's vulnerabilities.

7. PRESS FORWARD — We refuse to be defined by our failures. We will spur one another toward the grace of Jesus Christ and the freedom He provides.

Signed and agreed upon this day, before God and this brotherhood.
"""

// MARK: - GroupView

struct GroupView: View {
    @State private var showCovenant = false
    @State private var showInvite   = false

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            VStack(spacing: 0) {
                ScrollView(showsIndicators: false) {
                    VStack(spacing: 18) {
                        GroupHeaderCard(name: groupName, memberCount: sampleMembers.count)
                        membersSection
                        CovenantButton { showCovenant = true }
                    }
                    .padding(.horizontal, 20)
                    .padding(.top, 8)
                    .padding(.bottom, 24)
                }

                inviteBar
            }
        }
        .toolbar(.hidden, for: .navigationBar)
        .sheet(isPresented: $showCovenant) { CovenantSheet() }
        .sheet(isPresented: $showInvite)   { InviteSheet(groupName: groupName) }
    }

    private var membersSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Members")
                .font(.system(size: 16, weight: .semibold))
                .foregroundStyle(.white)

            VStack(spacing: 10) {
                ForEach(sampleMembers) { MemberRow(member: $0) }
            }
        }
    }

    private var inviteBar: some View {
        VStack(spacing: 0) {
            Divider().overlay(Color.white.opacity(0.07))
            Button {
                showInvite = true
            } label: {
                HStack(spacing: 10) {
                    Image(systemName: "person.badge.plus")
                        .font(.system(size: 16, weight: .semibold))
                    Text("Invite Member")
                        .font(.system(size: 16, weight: .semibold))
                }
                .foregroundStyle(Color.rfNavy)
                .frame(maxWidth: .infinity)
                .frame(height: 52)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(Color.rfGold)
                )
            }
            .padding(.horizontal, 20)
            .padding(.top, 14)
            .padding(.bottom, 28)
            .background(Color.rfNavy)
        }
    }
}

// MARK: - Group header card

private struct GroupHeaderCard: View {
    let name: String
    let memberCount: Int

    var body: some View {
        HStack(spacing: 16) {
            ZStack {
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color.rfGold.opacity(0.12))
                    .frame(width: 54, height: 54)
                Image(systemName: "person.3.fill")
                    .font(.system(size: 22))
                    .foregroundStyle(Color.rfGold)
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(name)
                    .font(.system(size: 22, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                Text("\(memberCount) members")
                    .font(.system(size: 13))
                    .foregroundStyle(Color.white.opacity(0.45))
            }

            Spacer()
        }
        .padding(20)
        .background(
            RoundedRectangle(cornerRadius: 20)
                .fill(
                    LinearGradient(
                        colors: [
                            Color(red: 0.14, green: 0.20, blue: 0.40),
                            Color(red: 0.09, green: 0.14, blue: 0.32),
                        ],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                )
                .overlay(
                    RoundedRectangle(cornerRadius: 20)
                        .stroke(Color.rfGold.opacity(0.18), lineWidth: 1)
                )
        )
    }
}

// MARK: - Member row

private struct MemberRow: View {
    let member: GroupMember

    private var initials: String {
        member.name
            .components(separatedBy: " ")
            .compactMap(\.first)
            .prefix(2)
            .map(String.init)
            .joined()
            .uppercased()
    }

    var body: some View {
        HStack(spacing: 14) {
            // Avatar with status dot
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.14))
                    .frame(width: 44, height: 44)
                Text(initials)
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundStyle(Color.rfGold)
            }
            .overlay(alignment: .bottomTrailing) {
                Circle()
                    .fill(member.health.color)
                    .frame(width: 12, height: 12)
                    .overlay(Circle().stroke(Color.rfNavy, lineWidth: 2))
                    .offset(x: 2, y: 2)
            }

            VStack(alignment: .leading, spacing: 3) {
                Text(member.name)
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundStyle(.white)
                Text(member.health.label)
                    .font(.system(size: 12))
                    .foregroundStyle(member.health.color.opacity(0.85))
            }

            Spacer(minLength: 8)

            VStack(alignment: .trailing, spacing: 2) {
                Text("\(member.streak)")
                    .font(.system(size: 22, weight: .bold, design: .rounded))
                    .foregroundStyle(member.streak > 0 ? Color.rfGold : Color.white.opacity(0.25))
                Text("day streak")
                    .font(.system(size: 10))
                    .foregroundStyle(Color.white.opacity(0.35))
            }
        }
        .padding(14)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.white.opacity(0.055))
                .overlay(
                    RoundedRectangle(cornerRadius: 16)
                        .stroke(Color.white.opacity(0.08), lineWidth: 1)
                )
        )
    }
}

// MARK: - Covenant button

private struct CovenantButton: View {
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            HStack(spacing: 14) {
                ZStack {
                    RoundedRectangle(cornerRadius: 10)
                        .fill(Color.rfGold.opacity(0.14))
                        .frame(width: 40, height: 40)
                    Image(systemName: "doc.text.fill")
                        .font(.system(size: 16))
                        .foregroundStyle(Color.rfGold)
                }

                VStack(alignment: .leading, spacing: 3) {
                    Text("Group Covenant")
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(.white)
                    Text("The agreement all members signed")
                        .font(.system(size: 12))
                        .foregroundStyle(Color.white.opacity(0.45))
                }

                Spacer()

                Image(systemName: "chevron.right")
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundStyle(Color.white.opacity(0.25))
            }
            .padding(16)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Color.white.opacity(0.055))
                    .overlay(
                        RoundedRectangle(cornerRadius: 16)
                            .stroke(Color.rfGold.opacity(0.14), lineWidth: 1)
                    )
            )
        }
    }
}

// MARK: - Covenant sheet

private struct CovenantSheet: View {
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()

            VStack(spacing: 0) {
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 20)

                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Group Covenant")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text(groupName)
                            .font(.system(size: 13))
                            .foregroundStyle(Color.rfGold.opacity(0.75))
                    }
                    Spacer()
                    Button { dismiss() } label: {
                        Image(systemName: "xmark")
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundStyle(Color.white.opacity(0.55))
                            .padding(10)
                            .background(Circle().fill(Color.white.opacity(0.09)))
                    }
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 20)

                Divider().overlay(Color.white.opacity(0.08))

                ScrollView(showsIndicators: false) {
                    Text(covenantText)
                        .font(.system(size: 14, design: .serif))
                        .foregroundStyle(Color.white.opacity(0.85))
                        .lineSpacing(6)
                        .padding(24)
                }
            }
        }
    }
}

// MARK: - Invite sheet

private struct InviteSheet: View {
    let groupName: String
    @Environment(\.dismiss) private var dismiss
    @State private var copied = false

    private let inviteCode = "remainfaithful.app/join/ib-4x9k"

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()

            VStack(spacing: 0) {
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 20)

                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Invite a Member")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Share the link to \(groupName)")
                            .font(.system(size: 13))
                            .foregroundStyle(Color.white.opacity(0.45))
                    }
                    Spacer()
                    Button { dismiss() } label: {
                        Image(systemName: "xmark")
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundStyle(Color.white.opacity(0.55))
                            .padding(10)
                            .background(Circle().fill(Color.white.opacity(0.09)))
                    }
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 24)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.bottom, 36)

                ZStack {
                    Circle()
                        .fill(Color.rfGold.opacity(0.12))
                        .frame(width: 80, height: 80)
                    Image(systemName: "person.badge.plus")
                        .font(.system(size: 34))
                        .foregroundStyle(Color.rfGold)
                }
                .padding(.bottom, 16)

                Text("Add a brother to your group")
                    .font(.system(size: 17, weight: .semibold))
                    .foregroundStyle(.white)
                    .padding(.bottom, 8)

                Text("Invite someone you trust to walk\nthis journey alongside you.")
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.45))
                    .multilineTextAlignment(.center)
                    .padding(.bottom, 36)

                // Invite link row
                HStack(spacing: 12) {
                    Image(systemName: "link")
                        .font(.system(size: 14))
                        .foregroundStyle(Color.rfGold.opacity(0.7))
                    Text(inviteCode)
                        .font(.system(size: 12, weight: .medium, design: .monospaced))
                        .foregroundStyle(Color.white.opacity(0.65))
                        .lineLimit(1)
                    Spacer()
                    Button {
                        UIPasteboard.general.string = inviteCode
                        withAnimation(.spring(response: 0.3, dampingFraction: 0.7)) { copied = true }
                        DispatchQueue.main.asyncAfter(deadline: .now() + 2) {
                            withAnimation { copied = false }
                        }
                    } label: {
                        Text(copied ? "Copied!" : "Copy")
                            .font(.system(size: 12, weight: .semibold))
                            .foregroundStyle(copied ? Color(red: 0.20, green: 0.78, blue: 0.45) : Color.rfGold)
                            .padding(.horizontal, 12)
                            .padding(.vertical, 6)
                            .background(
                                Capsule().fill(
                                    copied
                                    ? Color(red: 0.20, green: 0.78, blue: 0.45).opacity(0.14)
                                    : Color.rfGold.opacity(0.14)
                                )
                            )
                    }
                }
                .padding(14)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(Color.white.opacity(0.055))
                        .overlay(
                            RoundedRectangle(cornerRadius: 12)
                                .stroke(Color.white.opacity(0.08), lineWidth: 1)
                        )
                )
                .padding(.horizontal, 24)
                .padding(.bottom, 16)

                // Share button
                ShareLink(item: URL(string: "https://\(inviteCode)")!) {
                    HStack(spacing: 10) {
                        Image(systemName: "square.and.arrow.up")
                            .font(.system(size: 16, weight: .semibold))
                        Text("Share Invite Link")
                            .font(.system(size: 16, weight: .semibold))
                    }
                    .foregroundStyle(Color.rfNavy)
                    .frame(maxWidth: .infinity)
                    .frame(height: 52)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(Color.rfGold)
                    )
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 8)
            }
        }
    }
}

#Preview {
    GroupView()
}
