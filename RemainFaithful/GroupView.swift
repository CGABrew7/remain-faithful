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
    let userId: Int
    let name: String
    let streak: Int
    let health: AccountabilityHealth
    var phone: String? = nil

    init(userId: Int = 0, name: String, streak: Int, health: AccountabilityHealth, phone: String? = nil) {
        self.userId = userId
        self.name   = name
        self.streak = streak
        self.health = health
        self.phone  = phone
    }
}

// MARK: - Placeholder data

private let groupName = "Iron Brotherhood"

private let sampleMembers: [GroupMember] = [
    .init(name: "James Bishop",  streak: 21, health: .strong,     phone: "5555550101"),
    .init(name: "Marcus Cole",   streak: 7,  health: .watchful,   phone: "5555550102"),
    .init(name: "David Torres",  streak: 3,  health: .struggling, phone: "5555550103"),
    .init(name: "Nathan Wells",  streak: 14, health: .strong,     phone: "5555550104"),
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
    @AppStorage("primaryGroupID")    private var primaryGroupID    = 0
    @AppStorage("customCovenantText") private var customCovenantText = ""
    @EnvironmentObject private var appState: AppState

    @State private var showCovenant       = false
    @State private var showInvite         = false
    @State private var showRenameGroup    = false
    @State private var renameGroupText    = ""
    @State private var selectedMember: GroupMember? = nil
    @State private var showEditCovenant   = false
    @State private var showCovenantAlert  = false
    @State private var showCreateGroup    = false
    @State private var liveGroupName      = ""
    @State private var liveMembers:  [GroupMember] = []
    @State private var isLoading          = false
    @State private var loadError: String?

    private var displayName: String {
        if !liveGroupName.isEmpty { return liveGroupName }
        return appState.isDemoMode ? groupName : ""
    }
    private var displayMembers: [GroupMember] {
        if !liveMembers.isEmpty { return liveMembers }
        return appState.isDemoMode ? sampleMembers : []
    }
    private var displayCovenant: String       { customCovenantText.isEmpty ? covenantText : customCovenantText }

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            if !appState.isDemoMode && primaryGroupID == 0 {
                noGroupEmptyState
            } else {
                VStack(spacing: 0) {
                    ScrollView(showsIndicators: false) {
                        VStack(spacing: 18) {
                            GroupHeaderCard(name: displayName, memberCount: displayMembers.count) {
                                renameGroupText = displayName
                                showRenameGroup = true
                            }
                            if isLoading {
                                ProgressView()
                                    .tint(Color.rfGold)
                                    .padding(.vertical, 12)
                            } else if let err = loadError {
                                Text(err)
                                    .font(.system(size: 13))
                                    .foregroundStyle(Color.white.opacity(0.45))
                                    .multilineTextAlignment(.center)
                                    .padding(.vertical, 8)
                            }
                            if !displayMembers.isEmpty {
                                membersSection
                            }
                            CovenantButton { showCovenant = true }
                        }
                        .padding(.horizontal, 20)
                        .padding(.top, 8)
                        .padding(.bottom, 24)
                    }

                    inviteBar
                }
            }
        }
        .toolbar(.hidden, for: .navigationBar)
        .sheet(isPresented: $showCovenant) {
            CovenantSheet(covenantText: displayCovenant) { showEditCovenant = true }
        }
        .sheet(isPresented: $showInvite) {
            InviteSheet(groupName: displayName, groupID: primaryGroupID)
        }
        .sheet(isPresented: $showRenameGroup) {
            RenameGroupSheet(name: $renameGroupText) { newName in
                liveGroupName = newName
            }
        }
        .sheet(item: $selectedMember) { member in
            MemberDetailView(member: member)
        }
        .sheet(isPresented: $showEditCovenant) {
            EditCovenantSheet(text: $customCovenantText) {
                showCovenantAlert = true
            }
        }
        .sheet(isPresented: $showCreateGroup) {
            CreateGroupSheet { createdName in
                liveGroupName = createdName
            }
        }
        .alert("Covenant Updated", isPresented: $showCovenantAlert) {
            Button("OK", role: .cancel) { }
        } message: {
            Text("All members have been notified that the group covenant has been updated and will be prompted to re-accept it.")
        }
        .task { await discoverGroup() }
        .task(id: primaryGroupID) { await loadGroup() }
    }

    private var noGroupEmptyState: some View {
        VStack(spacing: 28) {
            Spacer()
            Image(systemName: "person.3.fill")
                .font(.system(size: 52))
                .foregroundStyle(Color.rfGold.opacity(0.35))
            VStack(spacing: 10) {
                Text("No Group Yet")
                    .font(.system(size: 22, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                Text("You are not in any accountability groups.\nCreate one, or ask a member for an invite link.")
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.50))
                    .multilineTextAlignment(.center)
                    .lineSpacing(3)
            }
            .padding(.horizontal, 32)
            Button {
                showCreateGroup = true
            } label: {
                HStack(spacing: 10) {
                    Image(systemName: "plus.circle.fill")
                        .font(.system(size: 16))
                    Text("Create a Group")
                        .font(.system(size: 16, weight: .semibold))
                }
                .foregroundStyle(Color.rfNavy)
                .frame(maxWidth: .infinity)
                .frame(height: 54)
                .background(RoundedRectangle(cornerRadius: 14).fill(Color.rfGold))
            }
            .padding(.horizontal, 40)
            Spacer()
        }
    }

    @MainActor
    private func discoverGroup() async {
        guard !appState.isDemoMode, primaryGroupID == 0, APIClient.shared.isAuthenticated else { return }
        guard let groups = try? await APIClient.shared.listMyGroups(), let first = groups.first else { return }
        primaryGroupID = first.id
    }

    @MainActor
    private func loadGroup() async {
        guard !appState.isDemoMode else { return }
        guard primaryGroupID > 0, APIClient.shared.isAuthenticated else { return }
        isLoading = true
        loadError = nil
        defer { isLoading = false }
        do {
            let group = try await APIClient.shared.getGroup(id: primaryGroupID)
            liveGroupName = group.name
            liveMembers = (group.members ?? []).map { m in
                let health: AccountabilityHealth
                switch m.flagsLast30 {
                case 0:    health = .strong
                case 1, 2: health = .watchful
                default:   health = .struggling
                }
                return GroupMember(userId: m.userId, name: m.user.name,
                                   streak: m.streakDays, health: health)
            }
        } catch {
            loadError = error.localizedDescription
        }
    }

    private var membersSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Members")
                .font(.system(size: 16, weight: .semibold))
                .foregroundStyle(.white)

            VStack(spacing: 10) {
                ForEach(displayMembers) { member in
                    Button { selectedMember = member } label: {
                        MemberRow(member: member)
                    }
                    .buttonStyle(.plain)
                }
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
    let onRename: () -> Void

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

            Button(action: onRename) {
                Image(systemName: "pencil")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(Color.rfGold)
                    .padding(9)
                    .background(Circle().fill(Color.rfGold.opacity(0.13)))
            }
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

// MARK: - Rename group sheet

private struct RenameGroupSheet: View {
    @Binding var name: String
    let onSave: (String) -> Void
    @Environment(\.dismiss) private var dismiss
    @FocusState private var focused: Bool

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 24)

                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Rename Group")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Choose a new name for this group")
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
                    .padding(.bottom, 24)

                HStack(spacing: 12) {
                    Image(systemName: "person.3.fill")
                        .font(.system(size: 15))
                        .foregroundStyle(focused ? Color.rfGold : Color.white.opacity(0.45))
                        .frame(width: 22)
                    TextField("", text: $name,
                              prompt: Text("Group name").foregroundColor(Color.white.opacity(0.38)))
                        .textInputAutocapitalization(.words)
                        .foregroundStyle(.white)
                        .focused($focused)
                        .submitLabel(.done)
                        .onSubmit { saveAndDismiss() }
                }
                .padding(.horizontal, 18)
                .frame(height: 54)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(Color.white.opacity(focused ? 0.11 : 0.07))
                        .overlay(RoundedRectangle(cornerRadius: 14)
                            .stroke(focused ? Color.rfGold : Color.white.opacity(0.10), lineWidth: 1.5))
                )
                .animation(.easeInOut(duration: 0.18), value: focused)
                .padding(.horizontal, 24)
                .padding(.bottom, 24)

                Button(action: saveAndDismiss) {
                    Text("Save")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundStyle(name.trimmingCharacters(in: .whitespaces).isEmpty ? Color.white.opacity(0.35) : Color.rfNavy)
                        .frame(maxWidth: .infinity)
                        .frame(height: 52)
                        .background(
                            RoundedRectangle(cornerRadius: 14)
                                .fill(name.trimmingCharacters(in: .whitespaces).isEmpty ? Color.white.opacity(0.08) : Color.rfGold)
                        )
                }
                .disabled(name.trimmingCharacters(in: .whitespaces).isEmpty)
                .padding(.horizontal, 24)
                .padding(.bottom, 32)
            }
        }
        .presentationDetents([.height(320)])
        .onAppear { focused = true }
    }

    private func saveAndDismiss() {
        let trimmed = name.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty else { return }
        onSave(trimmed)
        dismiss()
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

            Image(systemName: "chevron.right")
                .font(.system(size: 11, weight: .semibold))
                .foregroundStyle(Color.white.opacity(0.20))
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

// MARK: - Member detail view

private struct MemberDetailView: View {
    let member: GroupMember
    @Environment(\.dismiss)  private var dismiss
    @Environment(\.openURL)  private var openURL
    @State private var encouragementSent = false
    @State private var memberAlerts: [ActivityEvent] = []
    @State private var alertsLoading = true

    private var firstName: String {
        member.name.components(separatedBy: " ").first ?? member.name
    }

    private var initials: String {
        member.name.components(separatedBy: " ")
            .compactMap(\.first).prefix(2).map(String.init).joined().uppercased()
    }

    var body: some View {
        ZStack {
            Color(red: 0.04, green: 0.07, blue: 0.18).ignoresSafeArea()

            VStack(spacing: 0) {
                dismissBar

                ScrollView(showsIndicators: false) {
                    VStack(spacing: 20) {
                        // Avatar + name + health
                        VStack(spacing: 12) {
                            ZStack {
                                Circle()
                                    .fill(Color.rfGold.opacity(0.14))
                                    .frame(width: 80, height: 80)
                                Text(initials)
                                    .font(.system(size: 28, weight: .semibold))
                                    .foregroundStyle(Color.rfGold)
                            }
                            .overlay(alignment: .bottomTrailing) {
                                Circle()
                                    .fill(member.health.color)
                                    .frame(width: 18, height: 18)
                                    .overlay(Circle().stroke(Color(red: 0.04, green: 0.07, blue: 0.18), lineWidth: 2.5))
                                    .offset(x: 3, y: 3)
                            }

                            Text(member.name)
                                .font(.system(size: 22, weight: .bold, design: .serif))
                                .foregroundStyle(.white)

                            HStack(spacing: 6) {
                                Circle()
                                    .fill(member.health.color)
                                    .frame(width: 8, height: 8)
                                Text(member.health.label)
                                    .font(.system(size: 13, weight: .medium))
                                    .foregroundStyle(member.health.color)
                            }
                        }
                        .padding(.top, 8)

                        // Streak card
                        HStack(spacing: 0) {
                            VStack(spacing: 4) {
                                Text("\(member.streak)")
                                    .font(.system(size: 40, weight: .bold, design: .rounded))
                                    .foregroundStyle(Color.rfGold)
                                Text("day streak")
                                    .font(.system(size: 12))
                                    .foregroundStyle(Color.white.opacity(0.45))
                            }
                            .frame(maxWidth: .infinity)
                        }
                        .padding(20)
                        .background(
                            RoundedRectangle(cornerRadius: 16)
                                .fill(Color.white.opacity(0.055))
                                .overlay(RoundedRectangle(cornerRadius: 16)
                                    .stroke(Color.white.opacity(0.08), lineWidth: 1))
                        )

                        // Recent flags
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Recent Flags")
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundStyle(.white)

                            if alertsLoading {
                                HStack {
                                    Spacer()
                                    ProgressView().tint(Color.rfGold)
                                    Spacer()
                                }
                                .padding(.vertical, 12)
                            } else if memberAlerts.isEmpty {
                                Text("No recent flags")
                                    .font(.system(size: 13))
                                    .foregroundStyle(Color.white.opacity(0.35))
                                    .frame(maxWidth: .infinity, alignment: .center)
                                    .padding(.vertical, 12)
                            } else {
                                VStack(spacing: 10) {
                                    ForEach(memberAlerts) { event in
                                        memberEventRow(event)
                                    }
                                }
                            }
                        }

                        // Action buttons
                        VStack(spacing: 12) {
                            if let phone = member.phone {
                                Button {
                                    if let url = URL(string: "tel://\(phone)") {
                                        openURL(url)
                                    }
                                } label: {
                                    HStack(spacing: 10) {
                                        Image(systemName: "phone.fill")
                                            .font(.system(size: 16))
                                        Text("Call \(firstName)")
                                            .font(.system(size: 16, weight: .semibold))
                                    }
                                    .foregroundStyle(Color.rfNavy)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 54)
                                    .background(RoundedRectangle(cornerRadius: 14).fill(Color.rfGold))
                                }
                            }

                            Button {
                                withAnimation(.spring(response: 0.35, dampingFraction: 0.7)) {
                                    encouragementSent = true
                                }
                            } label: {
                                HStack(spacing: 10) {
                                    Image(systemName: encouragementSent ? "checkmark.circle.fill" : "hand.raised.fill")
                                        .font(.system(size: 16))
                                    Text(encouragementSent ? "Encouragement Sent!" : "Send Encouragement")
                                        .font(.system(size: 16, weight: .semibold))
                                }
                                .foregroundStyle(encouragementSent
                                    ? Color(red: 0.20, green: 0.78, blue: 0.45)
                                    : Color(red: 0.28, green: 0.56, blue: 0.95))
                                .frame(maxWidth: .infinity)
                                .frame(height: 54)
                                .background(
                                    RoundedRectangle(cornerRadius: 14)
                                        .fill((encouragementSent
                                               ? Color(red: 0.20, green: 0.78, blue: 0.45)
                                               : Color(red: 0.28, green: 0.56, blue: 0.95)).opacity(0.13))
                                        .overlay(RoundedRectangle(cornerRadius: 14)
                                            .stroke((encouragementSent
                                                     ? Color(red: 0.20, green: 0.78, blue: 0.45)
                                                     : Color(red: 0.28, green: 0.56, blue: 0.95)).opacity(0.35), lineWidth: 1.5))
                                )
                            }
                            .disabled(encouragementSent)
                            .animation(.easeInOut(duration: 0.25), value: encouragementSent)
                        }

                        Button { dismiss() } label: {
                            Text("Close")
                                .font(.system(size: 15, weight: .semibold))
                                .foregroundStyle(Color.white.opacity(0.60))
                                .frame(maxWidth: .infinity)
                                .frame(height: 50)
                                .background(
                                    RoundedRectangle(cornerRadius: 14)
                                        .fill(Color.white.opacity(0.07))
                                        .overlay(RoundedRectangle(cornerRadius: 14)
                                            .stroke(Color.white.opacity(0.12), lineWidth: 1))
                                )
                        }
                    }
                    .padding(.horizontal, 24)
                    .padding(.bottom, 32)
                }
            }
        }
        .task {
            defer { alertsLoading = false }
            guard member.userId != 0 else { return }
            if let alerts = try? await APIClient.shared.listAlerts() {
                memberAlerts = alerts
                    .filter { $0.event.userId == member.userId }
                    .prefix(10)
                    .compactMap { ActivityEvent.from(remote: $0.event) }
            }
        }
    }

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

    private func memberEventRow(_ event: ActivityEvent) -> some View {
        HStack(spacing: 14) {
            ZStack {
                Circle()
                    .fill(event.category.tint.opacity(0.14))
                    .frame(width: 38, height: 38)
                Image(systemName: event.category.icon)
                    .font(.system(size: 14))
                    .foregroundStyle(event.category.tint)
            }
            VStack(alignment: .leading, spacing: 3) {
                Text(event.category.rawValue)
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundStyle(.white)
                Text(event.description)
                    .font(.system(size: 11))
                    .foregroundStyle(Color.white.opacity(0.48))
            }
            Spacer(minLength: 4)
            VStack(alignment: .trailing, spacing: 4) {
                Text(event.severity.rawValue.uppercased())
                    .font(.system(size: 9, weight: .bold))
                    .kerning(0.5)
                    .foregroundStyle(event.severity.color)
                    .padding(.horizontal, 7)
                    .padding(.vertical, 3)
                    .background(Capsule().fill(event.severity.color.opacity(0.14)))
                Text(event.timeLabel)
                    .font(.system(size: 10))
                    .foregroundStyle(Color.white.opacity(0.32))
            }
        }
        .padding(12)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(Color.white.opacity(0.055))
                .overlay(RoundedRectangle(cornerRadius: 14).stroke(Color.white.opacity(0.08), lineWidth: 1))
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
    let covenantText: String
    let onEdit: () -> Void
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

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.bottom, 16)

                Button {
                    dismiss()
                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.35) { onEdit() }
                } label: {
                    HStack(spacing: 8) {
                        Image(systemName: "pencil")
                            .font(.system(size: 14, weight: .semibold))
                        Text("Edit Covenant")
                            .font(.system(size: 15, weight: .semibold))
                    }
                    .foregroundStyle(Color.rfGold)
                    .frame(maxWidth: .infinity)
                    .frame(height: 50)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(Color.rfGold.opacity(0.10))
                            .overlay(RoundedRectangle(cornerRadius: 14)
                                .stroke(Color.rfGold.opacity(0.22), lineWidth: 1))
                    )
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 32)
            }
        }
    }
}

// MARK: - Edit covenant sheet

private struct EditCovenantSheet: View {
    @Binding var text: String
    let onSave: () -> Void
    @Environment(\.dismiss) private var dismiss
    @State private var draft = ""

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
                        Text("Edit Covenant")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Changes will notify all members")
                            .font(.system(size: 13))
                            .foregroundStyle(Color(red: 0.95, green: 0.72, blue: 0.22))
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
                .padding(.bottom, 16)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.bottom, 16)

                TextEditor(text: $draft)
                    .scrollContentBackground(.hidden)
                    .background(Color.clear)
                    .foregroundStyle(Color.white.opacity(0.85))
                    .font(.system(size: 14, design: .serif))
                    .lineSpacing(4)
                    .padding(.horizontal, 20)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.top, 16)
                    .padding(.bottom, 16)

                Button {
                    text = draft
                    dismiss()
                    onSave()
                } label: {
                    Text("Save Covenant")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundStyle(Color.rfNavy)
                        .frame(maxWidth: .infinity)
                        .frame(height: 52)
                        .background(RoundedRectangle(cornerRadius: 14).fill(Color.rfGold))
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 32)
            }
        }
        .onAppear {
            draft = text.isEmpty ? covenantText : text
        }
    }
}

// MARK: - Invite sheet

private struct InviteSheet: View {
    let groupName: String
    let groupID:   Int
    @Environment(\.dismiss) private var dismiss
    @FocusState private var emailFocused: Bool
    @State private var inviteEmail  = ""
    @State private var isSent       = false
    @State private var copied       = false
    @State private var inviteError: String?

    private var inviteCode: String { "remainfaithful.app/join/\(groupID)" }
    private var inviteURL: URL { URL(string: "https://\(inviteCode)") ?? URL(string: "https://remainfaithful.app")! }
    private let green = Color(red: 0.20, green: 0.78, blue: 0.45)

    private var canSend: Bool {
        inviteEmail.contains("@") && inviteEmail.contains(".") && !isSent
    }

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
                        Text("Add someone to \(groupName)")
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
                    .padding(.bottom, 28)

                VStack(alignment: .leading, spacing: 12) {
                    Text("INVITE BY EMAIL")
                        .font(.system(size: 11, weight: .bold))
                        .foregroundStyle(Color.rfGold.opacity(0.75))
                        .kerning(1.2)
                        .padding(.horizontal, 24)

                    HStack(spacing: 12) {
                        Image(systemName: "envelope.fill")
                            .font(.system(size: 15))
                            .foregroundStyle(emailFocused ? Color.rfGold : Color.white.opacity(0.45))
                            .frame(width: 22)
                        TextField("", text: $inviteEmail,
                                  prompt: Text("Email address").foregroundColor(Color.white.opacity(0.38)))
                            .keyboardType(.emailAddress)
                            .textInputAutocapitalization(.never)
                            .autocorrectionDisabled()
                            .foregroundStyle(.white)
                            .focused($emailFocused)
                            .submitLabel(.send)
                            .onSubmit { sendInvite() }
                    }
                    .padding(.horizontal, 18)
                    .frame(height: 54)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(Color.white.opacity(emailFocused ? 0.11 : 0.07))
                            .overlay(
                                RoundedRectangle(cornerRadius: 14)
                                    .stroke(emailFocused ? Color.rfGold : Color.white.opacity(0.10),
                                            lineWidth: 1.5)
                            )
                    )
                    .animation(.easeInOut(duration: 0.18), value: emailFocused)
                    .padding(.horizontal, 24)

                    Button { sendInvite() } label: {
                        HStack(spacing: 10) {
                            Image(systemName: isSent ? "checkmark.circle.fill" : "paperplane.fill")
                                .font(.system(size: 16))
                            Text(isSent ? "Invite Sent!" : "Send Invite")
                                .font(.system(size: 16, weight: .semibold))
                        }
                        .foregroundStyle(isSent ? green : Color.rfNavy)
                        .frame(maxWidth: .infinity)
                        .frame(height: 52)
                        .background(
                            RoundedRectangle(cornerRadius: 14)
                                .fill(isSent ? green.opacity(0.15)
                                      : (canSend ? Color.rfGold : Color.rfGold.opacity(0.35)))
                        )
                    }
                    .disabled(!canSend)
                    .animation(.easeInOut(duration: 0.2), value: isSent)
                    .padding(.horizontal, 24)

                    if let err = inviteError {
                        Text(err)
                            .font(.system(size: 13))
                            .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                            .multilineTextAlignment(.center)
                            .padding(.horizontal, 24)
                            .transition(.opacity)
                    }
                }

                HStack {
                    Rectangle().fill(Color.white.opacity(0.08)).frame(height: 1)
                    Text("or share a link")
                        .font(.system(size: 11))
                        .foregroundStyle(Color.white.opacity(0.30))
                        .padding(.horizontal, 12)
                        .fixedSize()
                    Rectangle().fill(Color.white.opacity(0.08)).frame(height: 1)
                }
                .padding(.horizontal, 24)
                .padding(.vertical, 24)

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
                            .foregroundStyle(copied ? green : Color.rfGold)
                            .padding(.horizontal, 12)
                            .padding(.vertical, 6)
                            .background(Capsule().fill(copied ? green.opacity(0.14) : Color.rfGold.opacity(0.14)))
                    }
                }
                .padding(14)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(Color.white.opacity(0.055))
                        .overlay(RoundedRectangle(cornerRadius: 12).stroke(Color.white.opacity(0.08), lineWidth: 1))
                )
                .padding(.horizontal, 24)
                .padding(.bottom, 16)

                ShareLink(item: inviteURL) {
                    HStack(spacing: 10) {
                        Image(systemName: "square.and.arrow.up")
                            .font(.system(size: 16, weight: .semibold))
                        Text("Share Invite Link")
                            .font(.system(size: 16, weight: .semibold))
                    }
                    .foregroundStyle(Color.rfNavy)
                    .frame(maxWidth: .infinity)
                    .frame(height: 52)
                    .background(RoundedRectangle(cornerRadius: 14).fill(Color.rfGold.opacity(0.6)))
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 8)
            }
        }
    }

    private func sendInvite() {
        guard canSend else { return }
        emailFocused = false
        inviteError  = nil
        let email = inviteEmail
        let gid   = groupID
        Task {
            do {
                try await APIClient.shared.inviteMember(groupID: gid, email: email)
                await MainActor.run {
                    withAnimation(.spring(response: 0.35, dampingFraction: 0.7)) { isSent = true }
                }
            } catch {
                await MainActor.run {
                    inviteError = error.localizedDescription
                }
            }
        }
    }
}

// MARK: - Create Group sheet (multi-step)

private let defaultCovenant = """
We, the members of this accountability group, covenant together before God and one another to:

1. PURSUE PURITY — We commit to actively pursue sexual purity in thought, word, and deed.

2. WALK IN HONESTY — We will be completely truthful about our struggles, failures, and victories.

3. PRAY FAITHFULLY — We commit to pray regularly for each member of this group.

4. RESPOND QUICKLY — When a brother sends an alert, we will respond promptly with encouragement and prayer.

5. GUARD ONE ANOTHER — We take responsibility for one another's spiritual health.

6. HOLD CONFIDENCE — What is shared in this group stays in this group.

7. PRESS FORWARD — We refuse to be defined by our failures. We spur one another toward freedom in Christ.

Signed and agreed upon this day, before God and this brotherhood.
"""

private struct CreateGroupSheet: View {
    let onCreate: (String) -> Void
    @AppStorage("primaryGroupID") private var primaryGroupID = 0
    @Environment(\.dismiss) private var dismiss

    @State private var step         = 0          // 0=name, 1=covenant, 2=invite, 3=review
    @State private var groupName    = ""
    @State private var covenant     = defaultCovenant
    @State private var inviteEmails: [String] = [""]  // at least one field
    @State private var isCreating   = false
    @State private var createError: String?

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                // Drag indicator + step dots
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 16)

                HStack(spacing: 7) {
                    ForEach(0..<4) { i in
                        Capsule()
                            .fill(i == step ? Color.rfGold : Color.white.opacity(0.20))
                            .frame(width: i == step ? 20 : 7, height: 7)
                            .animation(.spring(response: 0.35, dampingFraction: 0.7), value: step)
                    }
                }
                .padding(.bottom, 20)

                // Step content
                ZStack {
                    if step == 0 { stepName }
                    else if step == 1 { stepCovenant }
                    else if step == 2 { stepInvite }
                    else { stepReview }
                }
                .animation(.easeInOut(duration: 0.3), value: step)
            }
        }
        .presentationDetents([.large])
    }

    // MARK: Step 0 — Name

    @FocusState private var nameFocused: Bool

    private var stepName: some View {
        VStack(spacing: 0) {
            sheetHeaderBlock(
                title: "Name your group",
                subtitle: "Choose a name that captures your brotherhood"
            )

            HStack(spacing: 12) {
                Image(systemName: "person.3.fill")
                    .font(.system(size: 15))
                    .foregroundStyle(nameFocused ? Color.rfGold : Color.white.opacity(0.45))
                    .frame(width: 22)
                TextField("", text: $groupName,
                          prompt: Text("e.g. Iron Brotherhood").foregroundColor(Color.white.opacity(0.38)))
                    .textInputAutocapitalization(.words)
                    .foregroundStyle(.white)
                    .focused($nameFocused)
                    .submitLabel(.next)
                    .onSubmit { if !groupName.trimmingCharacters(in: .whitespaces).isEmpty { step = 1 } }
            }
            .padding(.horizontal, 18)
            .frame(height: 54)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color.white.opacity(nameFocused ? 0.11 : 0.07))
                    .overlay(RoundedRectangle(cornerRadius: 14)
                        .stroke(nameFocused ? Color.rfGold : Color.white.opacity(0.10), lineWidth: 1.5))
            )
            .animation(.easeInOut(duration: 0.18), value: nameFocused)
            .padding(.horizontal, 24)
            .padding(.bottom, 24)

            navButtons(
                nextLabel: "Set Covenant",
                nextEnabled: !groupName.trimmingCharacters(in: .whitespaces).isEmpty,
                onNext: { step = 1 }
            )
        }
        .onAppear { nameFocused = true }
    }

    // MARK: Step 1 — Covenant

    private var stepCovenant: some View {
        VStack(spacing: 0) {
            sheetHeaderBlock(
                title: "Group Covenant",
                subtitle: "Edit the agreement your members will sign"
            )

            TextEditor(text: $covenant)
                .scrollContentBackground(.hidden)
                .background(Color.white.opacity(0.07))
                .foregroundStyle(Color.white.opacity(0.85))
                .font(.system(size: 13, design: .serif))
                .lineSpacing(3)
                .clipShape(RoundedRectangle(cornerRadius: 14))
                .overlay(RoundedRectangle(cornerRadius: 14).stroke(Color.white.opacity(0.10), lineWidth: 1))
                .frame(maxHeight: .infinity)
                .padding(.horizontal, 24)
                .padding(.bottom, 16)

            navButtons(
                backLabel: "Back",
                nextLabel: "Invite Members",
                nextEnabled: !covenant.trimmingCharacters(in: .whitespaces).isEmpty,
                onBack: { step = 0 },
                onNext: { step = 2 }
            )
        }
    }

    // MARK: Step 2 — Invite

    private var stepInvite: some View {
        VStack(spacing: 0) {
            sheetHeaderBlock(
                title: "Invite Members",
                subtitle: "Enter emails — they'll receive an invitation link"
            )

            ScrollView(showsIndicators: false) {
                VStack(spacing: 10) {
                    ForEach(inviteEmails.indices, id: \.self) { i in
                        HStack(spacing: 10) {
                            Image(systemName: "envelope.fill")
                                .font(.system(size: 14))
                                .foregroundStyle(Color.white.opacity(0.40))
                                .frame(width: 20)
                            TextField("", text: $inviteEmails[i],
                                      prompt: Text("Email address").foregroundColor(Color.white.opacity(0.32)))
                                .keyboardType(.emailAddress)
                                .textInputAutocapitalization(.never)
                                .autocorrectionDisabled()
                                .foregroundStyle(.white)
                                .font(.system(size: 14))
                            if inviteEmails.count > 1 {
                                Button { inviteEmails.remove(at: i) } label: {
                                    Image(systemName: "minus.circle.fill")
                                        .foregroundStyle(Color(red: 0.90, green: 0.30, blue: 0.30).opacity(0.70))
                                }
                            }
                        }
                        .padding(.horizontal, 16)
                        .frame(height: 50)
                        .background(
                            RoundedRectangle(cornerRadius: 12)
                                .fill(Color.white.opacity(0.07))
                                .overlay(RoundedRectangle(cornerRadius: 12)
                                    .stroke(Color.white.opacity(0.09), lineWidth: 1))
                        )
                    }

                    Button {
                        inviteEmails.append("")
                    } label: {
                        HStack(spacing: 6) {
                            Image(systemName: "plus.circle")
                                .font(.system(size: 14))
                            Text("Add another email")
                                .font(.system(size: 14))
                        }
                        .foregroundStyle(Color.rfGold.opacity(0.75))
                    }
                    .padding(.top, 4)
                }
                .padding(.horizontal, 24)
            }
            .frame(maxHeight: .infinity)
            .padding(.bottom, 8)

            navButtons(
                backLabel: "Back",
                nextLabel: "Review & Create",
                nextEnabled: true,
                onBack: { step = 1 },
                onNext: { step = 3 }
            )
        }
    }

    // MARK: Step 3 — Review

    private var stepReview: some View {
        VStack(spacing: 0) {
            sheetHeaderBlock(
                title: "Review",
                subtitle: "Confirm and create your group"
            )

            ScrollView(showsIndicators: false) {
                VStack(alignment: .leading, spacing: 16) {
                    reviewRow(icon: "person.3.fill", label: "Name", value: groupName)
                    Divider().overlay(Color.white.opacity(0.08))
                    reviewRow(icon: "doc.text.fill", label: "Covenant", value: "Custom group covenant set")

                    let validEmails = inviteEmails.filter { $0.contains("@") }
                    if !validEmails.isEmpty {
                        Divider().overlay(Color.white.opacity(0.08))
                        VStack(alignment: .leading, spacing: 8) {
                            HStack(spacing: 8) {
                                Image(systemName: "envelope.fill")
                                    .font(.system(size: 14))
                                    .foregroundStyle(Color.rfGold.opacity(0.7))
                                Text("Inviting \(validEmails.count) member\(validEmails.count == 1 ? "" : "s")")
                                    .font(.system(size: 13, weight: .semibold))
                                    .foregroundStyle(.white)
                            }
                            ForEach(validEmails, id: \.self) { email in
                                Text("• " + email)
                                    .font(.system(size: 12))
                                    .foregroundStyle(Color.white.opacity(0.55))
                            }
                        }
                    }

                    if let err = createError {
                        Text(err)
                            .font(.system(size: 13))
                            .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                            .transition(.opacity)
                    }
                }
                .padding(20)
                .background(
                    RoundedRectangle(cornerRadius: 16)
                        .fill(Color.white.opacity(0.055))
                        .overlay(RoundedRectangle(cornerRadius: 16)
                            .stroke(Color.white.opacity(0.08), lineWidth: 1))
                )
                .padding(.horizontal, 24)
            }
            .frame(maxHeight: .infinity)
            .padding(.bottom, 8)

            navButtons(
                backLabel: "Back",
                nextLabel: isCreating ? "Creating…" : "Create Group",
                nextEnabled: !isCreating,
                onBack: { step = 2 },
                onNext: { Task { await createGroup() } }
            )
        }
    }

    // MARK: Create

    @MainActor
    private func createGroup() async {
        guard !isCreating else { return }
        isCreating = true
        createError = nil
        let trimmed = groupName.trimmingCharacters(in: .whitespaces)
        do {
            let group = try await APIClient.shared.createGroup(name: trimmed, covenant: covenant)
            primaryGroupID = group.id
            // Send email invites for all valid addresses.
            let validEmails = inviteEmails.filter { $0.contains("@") && $0.contains(".") }
            for email in validEmails {
                try? await APIClient.shared.groupEmailInvite(groupID: group.id, email: email)
            }
            onCreate(group.name)
            dismiss()
        } catch {
            createError = error.localizedDescription
        }
        isCreating = false
    }

    // MARK: Shared subviews

    private func sheetHeaderBlock(title: String, subtitle: String) -> some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                Text(title)
                    .font(.system(size: 20, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                Text(subtitle)
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
        .padding(.bottom, 20)
    }

    private func reviewRow(icon: String, label: String, value: String) -> some View {
        HStack(spacing: 12) {
            Image(systemName: icon)
                .font(.system(size: 14))
                .foregroundStyle(Color.rfGold.opacity(0.7))
                .frame(width: 20)
            VStack(alignment: .leading, spacing: 2) {
                Text(label)
                    .font(.system(size: 11, weight: .bold))
                    .foregroundStyle(Color.rfGold.opacity(0.65))
                    .kerning(0.8)
                Text(value)
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(.white)
            }
        }
    }

    private func navButtons(
        backLabel: String? = nil,
        nextLabel: String,
        nextEnabled: Bool,
        onBack: (() -> Void)? = nil,
        onNext: @escaping () -> Void
    ) -> some View {
        HStack(spacing: 12) {
            if let backLabel, let onBack {
                Button(action: onBack) {
                    Text(backLabel)
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(Color.white.opacity(0.65))
                        .frame(maxWidth: .infinity)
                        .frame(height: 50)
                        .background(
                            RoundedRectangle(cornerRadius: 14)
                                .fill(Color.white.opacity(0.08))
                                .overlay(RoundedRectangle(cornerRadius: 14)
                                    .stroke(Color.white.opacity(0.12), lineWidth: 1))
                        )
                }
            }
            Button(action: onNext) {
                Text(nextLabel)
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundStyle(nextEnabled ? Color.rfNavy : Color.white.opacity(0.35))
                    .frame(maxWidth: .infinity)
                    .frame(height: 50)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(nextEnabled ? Color.rfGold : Color.white.opacity(0.08))
                    )
            }
            .disabled(!nextEnabled)
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 32)
    }
}

#Preview {
    GroupView()
        .environmentObject(AppState.shared)
}
