import SwiftUI

// MARK: - SettingsView

struct SettingsView: View {
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @AppStorage("userName")             private var userName             = ""
    @AppStorage("userEmail")            private var userEmail            = ""
    @AppStorage("monitoringActive")     private var monitoringActive     = true
    @AppStorage("notificationsEnabled") private var notificationsEnabled = true
    @AppStorage("weeklyDigestEnabled")  private var weeklyDigestEnabled  = false
    @AppStorage("dataRetentionDays")    private var dataRetentionDays    = 30

    @State private var showRetentionPicker = false
    @State private var showCovenant        = false
    @State private var showActivityLog     = false
    @State private var showLeaveAlert      = false
    @State private var showDeleteAlert     = false
    @State private var showDonation        = false

    private var initials: String {
        let parts = userName.components(separatedBy: " ").filter { !$0.isEmpty }
        let f = parts.first.flatMap(\.first).map(String.init) ?? ""
        let l = parts.count > 1 ? (parts.last.flatMap(\.first).map(String.init) ?? "") : ""
        return (f + l).uppercased()
    }

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            ScrollView(showsIndicators: false) {
                VStack(spacing: 24) {
                    profileCard
                    accountabilitySection
                    monitoringSection
                    privacySection
                    supportUsSection
                    supportSection
                    dangerZone
                }
                .padding(.horizontal, 20)
                .padding(.top, 8)
                .padding(.bottom, 48)
            }
        }
        .toolbar(.hidden, for: .navigationBar)
        .sheet(isPresented: $showRetentionPicker) {
            RetentionPickerSheet(days: $dataRetentionDays)
                .presentationDetents([.height(320)])
                .presentationDragIndicator(.hidden)
        }
        .sheet(isPresented: $showDonation)    { DonationView() }
        .sheet(isPresented: $showCovenant)    { SettingsCovenantSheet() }
        .sheet(isPresented: $showActivityLog) { ActivityLogSheet() }
        .alert("Leave All Groups", isPresented: $showLeaveAlert) {
            Button("Leave All Groups", role: .destructive) {
                // PLACEHOLDER: call DELETE /groups/:id or a dedicated leave endpoint once built.
                // For now, clear the stored group ID so the view stops showing stale data.
                UserDefaults.standard.removeObject(forKey: "primaryGroupID")
            }
            Button("Cancel", role: .cancel) { }
        } message: {
            Text("You will be removed from all accountability groups. This cannot be undone.")
        }
        .alert("Delete Account", isPresented: $showDeleteAlert) {
            Button("Delete Account", role: .destructive) {
                AuthState.shared.clearSession()
                userName  = ""
                userEmail = ""
                hasCompletedOnboarding = false
            }
            Button("Cancel", role: .cancel) { }
        } message: {
            Text("All your data will be permanently deleted. This cannot be undone.")
        }
    }

    // MARK: - Profile card

    private var profileCard: some View {
        HStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.14))
                    .frame(width: 64, height: 64)
                Text(initials.isEmpty ? "?" : initials)
                    .font(.system(size: 22, weight: .semibold))
                    .foregroundStyle(Color.rfGold)
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(userName.isEmpty ? "Your Name" : userName)
                    .font(.system(size: 17, weight: .semibold))
                    .foregroundStyle(.white)
                Text(userEmail.isEmpty ? "your@email.com" : userEmail)
                    .font(.system(size: 13))
                    .foregroundStyle(Color.white.opacity(0.45))
            }

            Spacer()

            Button { } label: {
                Text("Edit")
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundStyle(Color.rfGold)
                    .padding(.horizontal, 14)
                    .padding(.vertical, 7)
                    .background(Capsule().fill(Color.rfGold.opacity(0.12)))
            }
        }
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.white.opacity(0.055))
                .overlay(RoundedRectangle(cornerRadius: 16)
                    .stroke(Color.rfGold.opacity(0.14), lineWidth: 1))
        )
    }

    // MARK: - Accountability

    private var accountabilitySection: some View {
        SettingsSection(title: "ACCOUNTABILITY") {
            SRow(icon: "person.2.fill", tint: Color.rfGold, label: "Manage Partners") { }
            rowDivider
            SRow(icon: "person.3.fill", tint: Color.rfGold, label: "Manage Groups")   { }
            rowDivider
            SRow(icon: "doc.text.fill", tint: Color.rfGold, label: "View Covenant")   { showCovenant = true }
        }
    }

    // MARK: - Monitoring

    private var monitoringSection: some View {
        SettingsSection(title: "MONITORING") {
            TRow(icon: "shield.fill",
                 tint: Color(red: 0.20, green: 0.78, blue: 0.45),
                 label: "Monitoring Active",
                 isOn: $monitoringActive)
            rowDivider
            TRow(icon: "bell.fill",
                 tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                 label: "Notifications",
                 isOn: $notificationsEnabled)
            rowDivider
            TRow(icon: "envelope.fill",
                 tint: Color(red: 0.65, green: 0.45, blue: 0.90),
                 label: "Weekly Digest Email",
                 isOn: $weeklyDigestEnabled)
        }
    }

    // MARK: - Privacy

    private var privacySection: some View {
        SettingsSection(title: "PRIVACY") {
            SRow(icon: "calendar",
                 tint: Color.rfGold,
                 label: "Data Retention",
                 value: "\(dataRetentionDays) days") { showRetentionPicker = true }
            rowDivider
            SRow(icon: "list.bullet.rectangle",
                 tint: Color.rfGold,
                 label: "View My Activity Log") { showActivityLog = true }
            rowDivider
            SRow(icon: "square.and.arrow.up",
                 tint: Color.rfGold,
                 label: "Export My Data") { }
        }
    }

    // MARK: - Support Us (donations)

    private var supportUsSection: some View {
        SettingsSection(title: "SUPPORT US") {
            Button { showDonation = true } label: {
                HStack(spacing: 14) {
                    iconBadge("heart.fill", tint: Color(red: 0.90, green: 0.25, blue: 0.48))
                    VStack(alignment: .leading, spacing: 2) {
                        Text("Support Remain Faithful")
                            .font(.system(size: 15, weight: .medium))
                            .foregroundStyle(.white)
                        Text("Keep RF free — one-time or monthly")
                            .font(.system(size: 12))
                            .foregroundStyle(Color.white.opacity(0.38))
                    }
                    Spacer()
                    Image(systemName: "chevron.right")
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(Color.white.opacity(0.20))
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 14)
            }
        }
    }

    // MARK: - Help & Support

    private var supportSection: some View {
        SettingsSection(title: "SUPPORT") {
            SRow(icon: "questionmark.circle.fill",
                 tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                 label: "How It Works") { }
            rowDivider
            SRow(icon: "bubble.left.fill",
                 tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                 label: "Contact Support") { }
            rowDivider
            SRow(icon: "star.fill",
                 tint: Color(red: 0.95, green: 0.62, blue: 0.10),
                 label: "Rate the App") { }
        }
    }

    // MARK: - Danger zone

    private static let red = Color(red: 0.90, green: 0.30, blue: 0.30)

    private var dangerZone: some View {
        VStack(spacing: 12) {
            Text("DANGER ZONE")
                .font(.system(size: 11, weight: .bold))
                .foregroundStyle(Self.red.opacity(0.75))
                .kerning(1.2)
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.leading, 4)

            VStack(spacing: 0) {
                DRow(icon: "rectangle.portrait.and.arrow.right.fill",
                     label: "Leave All Groups") { showLeaveAlert  = true }
                rowDivider
                DRow(icon: "trash.fill",
                     label: "Delete Account")   { showDeleteAlert = true }
            }
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Self.red.opacity(0.06))
                    .overlay(RoundedRectangle(cornerRadius: 16)
                        .stroke(Self.red.opacity(0.20), lineWidth: 1))
            )
            .clipShape(RoundedRectangle(cornerRadius: 16))

            // Sign out sits below the danger zone — less severe than deletion
            Button {
                AuthState.shared.clearSession()
                userName  = ""
                userEmail = ""
                // Keep hasCompletedOnboarding = true so LoginView appears on next launch
            } label: {
                Text("Sign Out")
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundStyle(Self.red)
                    .frame(maxWidth: .infinity)
                    .frame(height: 50)
                    .background(
                        RoundedRectangle(cornerRadius: 14)
                            .fill(Self.red.opacity(0.10))
                            .overlay(RoundedRectangle(cornerRadius: 14)
                                .stroke(Self.red.opacity(0.20), lineWidth: 1))
                    )
            }
        }
    }

    // MARK: - Helpers

    private var rowDivider: some View {
        Divider()
            .overlay(Color.white.opacity(0.07))
            .padding(.leading, 64) // icon(34) + hPad(16) + gap(14)
    }
}

// MARK: - Section wrapper

private struct SettingsSection<Content: View>: View {
    let title:   String
    let content: Content

    init(title: String, @ViewBuilder content: () -> Content) {
        self.title   = title
        self.content = content()
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            Text(title)
                .font(.system(size: 11, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.2)
                .padding(.leading, 4)

            VStack(spacing: 0) { content }
                .background(
                    RoundedRectangle(cornerRadius: 16)
                        .fill(Color.white.opacity(0.055))
                        .overlay(RoundedRectangle(cornerRadius: 16)
                            .stroke(Color.white.opacity(0.08), lineWidth: 1))
                )
                .clipShape(RoundedRectangle(cornerRadius: 16))
        }
    }
}

// MARK: - Standard settings row (navigation / disclosure)

private struct SRow: View {
    let icon:   String
    let tint:   Color
    let label:  String
    var value:  String? = nil
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            HStack(spacing: 14) {
                iconBadge(icon, tint: tint)
                Text(label)
                    .font(.system(size: 15, weight: .medium))
                    .foregroundStyle(.white)
                Spacer()
                if let value {
                    Text(value)
                        .font(.system(size: 14))
                        .foregroundStyle(Color.white.opacity(0.40))
                }
                Image(systemName: "chevron.right")
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundStyle(Color.white.opacity(0.20))
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 14)
        }
    }
}

// MARK: - Toggle row

private struct TRow: View {
    let icon:  String
    let tint:  Color
    let label: String
    @Binding var isOn: Bool

    var body: some View {
        HStack(spacing: 14) {
            iconBadge(icon, tint: tint)
            Text(label)
                .font(.system(size: 15, weight: .medium))
                .foregroundStyle(.white)
            Spacer()
            Toggle("", isOn: $isOn)
                .labelsHidden()
                .tint(Color.rfGold)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
    }
}

// MARK: - Danger row

private struct DRow: View {
    let icon:   String
    let label:  String
    let action: () -> Void

    private let red = Color(red: 0.90, green: 0.30, blue: 0.30)

    var body: some View {
        Button(action: action) {
            HStack(spacing: 14) {
                iconBadge(icon, tint: red)
                Text(label)
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundStyle(red)
                Spacer()
                Image(systemName: "chevron.right")
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundStyle(red.opacity(0.35))
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 14)
        }
    }
}

// MARK: - Shared icon badge

private func iconBadge(_ name: String, tint: Color) -> some View {
    ZStack {
        RoundedRectangle(cornerRadius: 8)
            .fill(tint.opacity(0.15))
            .frame(width: 34, height: 34)
        Image(systemName: name)
            .font(.system(size: 15))
            .foregroundStyle(tint)
    }
}

// MARK: - Data Retention picker sheet

private struct RetentionPickerSheet: View {
    @Binding var days: Int
    @Environment(\.dismiss) private var dismiss

    private let options = [30, 60, 90]

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
                        Text("Data Retention")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("How long to keep your activity logs")
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

                VStack(spacing: 0) {
                    ForEach(options, id: \.self) { option in
                        Button {
                            days = option
                            dismiss()
                        } label: {
                            HStack {
                                Text("\(option) days")
                                    .font(.system(size: 16, weight: .medium))
                                    .foregroundStyle(.white)
                                Spacer()
                                if days == option {
                                    Image(systemName: "checkmark")
                                        .font(.system(size: 14, weight: .semibold))
                                        .foregroundStyle(Color.rfGold)
                                }
                            }
                            .padding(.horizontal, 24)
                            .padding(.vertical, 18)
                        }
                        if option != options.last {
                            Divider()
                                .overlay(Color.white.opacity(0.07))
                                .padding(.horizontal, 24)
                        }
                    }
                }
                Spacer()
            }
        }
    }
}

// MARK: - Covenant sheet

private let settingsCovenantText = """
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

private struct SettingsCovenantSheet: View {
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
                        Text("The agreement all members signed")
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
                    Text(settingsCovenantText)
                        .font(.system(size: 14, design: .serif))
                        .foregroundStyle(Color.white.opacity(0.85))
                        .lineSpacing(6)
                        .padding(24)
                }
            }
        }
    }
}

// MARK: - Activity log sheet

private struct ActivityLogSheet: View {
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
                        Text("My Activity Log")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Your complete monitoring history")
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

                Spacer()
                VStack(spacing: 14) {
                    Image(systemName: "list.bullet.rectangle")
                        .font(.system(size: 38))
                        .foregroundStyle(Color.rfGold.opacity(0.45))
                    Text("Full history coming soon")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundStyle(.white)
                    Text("Your complete activity log will be\navailable in a future update.")
                        .font(.system(size: 14))
                        .foregroundStyle(Color.white.opacity(0.45))
                        .multilineTextAlignment(.center)
                }
                Spacer()
            }
        }
    }
}

#Preview {
    SettingsView()
}
