import SwiftUI
import StoreKit
import FamilyControls

// MARK: - SettingsView

struct SettingsView: View {
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @AppStorage("userName")             private var userName             = ""
    @AppStorage("userEmail")            private var userEmail            = ""
    @AppStorage("monitoringActive")     private var monitoringActive     = true
    @AppStorage("notificationsEnabled") private var notificationsEnabled = true
    @AppStorage("dataRetentionDays")    private var dataRetentionDays    = 30

    @Environment(\.openURL)       private var openURL
    @Environment(\.requestReview) private var requestReview

    @ObservedObject private var fcManager = FamilyControlsManager.shared

    @State private var showScreenTime       = false
    @State private var showRetentionPicker  = false
    @State private var showCovenant         = false
    @State private var showActivityLog      = false
    @State private var showLeaveConfirm     = false
    @State private var showDeleteConfirm    = false
    @State private var showDonation         = false
    @State private var showManagePartners   = false
    @State private var showManageGroups     = false
    @State private var showHowItWorks       = false
    @State private var showEditProfile      = false

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
                    screenTimeSection
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
        .sheet(isPresented: $showScreenTime) {
            NavigationStack { ScreenTimeMonitoringView() }
        }
        .sheet(isPresented: $showDonation)        { DonationView() }
        .sheet(isPresented: $showCovenant)        { SettingsCovenantSheet() }
        .sheet(isPresented: $showActivityLog)     { ActivityLogSheet() }
        .sheet(isPresented: $showManagePartners)  { ManagePartnersView() }
        .sheet(isPresented: $showManageGroups)    { ManageGroupsView() }
        .sheet(isPresented: $showHowItWorks)      { HowItWorksView() }
        .sheet(isPresented: $showEditProfile)    { EditProfileSheet() }
        .sheet(isPresented: $showLeaveConfirm) {
            DestructiveConfirmSheet(
                title: "Leave All Groups",
                message: "You will be removed from all accountability groups. This cannot be undone.",
                warning: "All members of your accountability groups will be notified that you are leaving.",
                confirmLabel: "Leave All Groups"
            ) {
                Task {
                    try? await APIClient.shared.leaveAllGroups()
                    UserDefaults.standard.removeObject(forKey: "primaryGroupID")
                }
            }
        }
        .sheet(isPresented: $showDeleteConfirm) {
            DestructiveConfirmSheet(
                title: "Delete Account",
                message: "All your data will be permanently deleted. This cannot be undone.",
                warning: "All members of your accountability groups will be notified that your account has been deleted.",
                confirmLabel: "Delete Account"
            ) {
                AuthState.shared.clearSession()
                userName  = ""
                userEmail = ""
                hasCompletedOnboarding = false
            }
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

            Button { showEditProfile = true } label: {
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
            SRow(icon: "person.2.fill", tint: Color.rfGold, label: "Manage Partners") { showManagePartners = true }
            rowDivider
            SRow(icon: "person.3.fill", tint: Color.rfGold, label: "Manage Groups")   { showManageGroups = true }
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
        }
    }

    // MARK: - Screen Time

    private var screenTimeSection: some View {
        SettingsSection(title: "SCREEN TIME") {
            Button { showScreenTime = true } label: {
                HStack(spacing: 14) {
                    iconBadge("hand.raised.fill",
                              tint: Color(red: 0.52, green: 0.36, blue: 0.92))
                    VStack(alignment: .leading, spacing: 3) {
                        Text("App Restrictions")
                            .font(.system(size: 15, weight: .medium))
                            .foregroundStyle(.white)
                        Text(screenTimeStatusLabel)
                            .font(.system(size: 12))
                            .foregroundStyle(screenTimeStatusColor)
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

    private var screenTimeStatusLabel: String {
        switch fcManager.authorizationStatus {
        case .approved:      return "Authorized"
        case .denied:        return "Denied — tap to review"
        case .notDetermined: return "Tap to set up"
        @unknown default:    return "Unknown"
        }
    }

    private var screenTimeStatusColor: Color {
        switch fcManager.authorizationStatus {
        case .approved:      return Color(red: 0.20, green: 0.78, blue: 0.45)
        case .denied:        return Color(red: 0.90, green: 0.30, blue: 0.30)
        case .notDetermined: return Color.rfGold.opacity(0.75)
        @unknown default:    return Color.white.opacity(0.40)
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
                 label: "How It Works") { showHowItWorks = true }
            rowDivider
            SRow(icon: "bubble.left.fill",
                 tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                 label: "Contact Support") {
                if let url = URL(string: "mailto:support@remainfaithful.com?subject=Remain%20Faithful%20Support%20Request") {
                    UIApplication.shared.open(url)
                }
            }
            rowDivider
            SRow(icon: "star.fill",
                 tint: Color(red: 0.95, green: 0.62, blue: 0.10),
                 label: "Rate the App") { requestReview() }
            rowDivider
            SRow(icon: "lightbulb.fill",
                 tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                 label: "Suggest an Improvement") {
                if let url = URL(string: "mailto:support@remainfaithful.com?subject=Remain%20Faithful%20Feature%20Suggestion") {
                    UIApplication.shared.open(url)
                }
            }
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
                     label: "Leave All Groups") { showLeaveConfirm  = true }
                rowDivider
                DRow(icon: "trash.fill",
                     label: "Delete Account")   { showDeleteConfirm = true }
            }
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Self.red.opacity(0.06))
                    .overlay(RoundedRectangle(cornerRadius: 16)
                        .stroke(Self.red.opacity(0.20), lineWidth: 1))
            )
            .clipShape(RoundedRectangle(cornerRadius: 16))

            Button {
                AuthState.shared.clearSession()
                userName  = ""
                userEmail = ""
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
            .padding(.leading, 64)
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

// MARK: - Standard settings row

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

// MARK: - Destructive confirm sheet

private struct DestructiveConfirmSheet: View {
    let title:        String
    let message:      String
    let warning:      String
    let confirmLabel: String
    let onConfirm:    () -> Void
    @Environment(\.dismiss) private var dismiss

    private let red = Color(red: 0.90, green: 0.30, blue: 0.30)

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 28)

                ZStack {
                    Circle()
                        .fill(red.opacity(0.12))
                        .frame(width: 64, height: 64)
                    Image(systemName: "exclamationmark.triangle.fill")
                        .font(.system(size: 28))
                        .foregroundStyle(red)
                }
                .padding(.bottom, 20)

                Text(title)
                    .font(.system(size: 20, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                    .multilineTextAlignment(.center)
                    .padding(.horizontal, 24)
                    .padding(.bottom, 12)

                Text(message)
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.65))
                    .multilineTextAlignment(.center)
                    .padding(.horizontal, 24)
                    .lineSpacing(3)
                    .padding(.bottom, 14)

                Text(warning)
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundStyle(red)
                    .multilineTextAlignment(.center)
                    .padding(.horizontal, 24)
                    .lineSpacing(3)
                    .padding(.bottom, 28)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.bottom, 20)

                VStack(spacing: 12) {
                    Button {
                        onConfirm()
                        dismiss()
                    } label: {
                        Text(confirmLabel)
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundStyle(.white)
                            .frame(maxWidth: .infinity)
                            .frame(height: 52)
                            .background(RoundedRectangle(cornerRadius: 14).fill(red))
                    }

                    Button { dismiss() } label: {
                        Text("Cancel")
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundStyle(Color.white.opacity(0.70))
                            .frame(maxWidth: .infinity)
                            .frame(height: 52)
                            .background(
                                RoundedRectangle(cornerRadius: 14)
                                    .fill(Color.white.opacity(0.08))
                                    .overlay(RoundedRectangle(cornerRadius: 14)
                                        .stroke(Color.white.opacity(0.12), lineWidth: 1))
                            )
                    }
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 36)
            }
        }
        .presentationDetents([.medium])
    }
}

// MARK: - Manage Partners

private struct ManagePartnersView: View {
    @Environment(\.dismiss) private var dismiss
    @FocusState private var emailFocused: Bool
    @State private var inviteEmail   = ""
    @State private var isSending     = false
    @State private var sentSuccess   = false
    @State private var errorMsg: String?
    @State private var partnerToRemove: PartnerItem? = nil
    @State private var showRemoveConfirm = false

    private let green = Color(red: 0.20, green: 0.78, blue: 0.45)
    private let red   = Color(red: 0.90, green: 0.30, blue: 0.30)

    @State private var partners:         [PartnerItem] = []
    @State private var isLoadingPartners = false
    @State private var isSettingPrimary  = false

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                dragIndicator()
                sheetHeader(title: "Manage Partners", subtitle: "Your one-to-one accountability partners")
                Divider().overlay(Color.white.opacity(0.08))
                ScrollView(showsIndicators: false) {
                    VStack(spacing: 20) {
                        partnersListSection
                        inviteSection
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 20)
                }
            }
        }
        .task { await loadPartners() }
        .alert("Remove Partner", isPresented: $showRemoveConfirm) {
            Button("Remove", role: .destructive) {
                guard let p = partnerToRemove else { return }
                Task {
                    try? await APIClient.shared.deleteRelationship(id: p.relationshipID)
                    await MainActor.run { partners.removeAll { $0.id == p.id } }
                }
            }
            Button("Cancel", role: .cancel) { }
        } message: {
            Text("This will end your accountability partnership. \(partnerToRemove?.name ?? "") will be notified.")
        }
    }

    private var partnersListSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("CURRENT PARTNERS")
                    .font(.system(size: 11, weight: .bold))
                    .foregroundStyle(Color.rfGold.opacity(0.75))
                    .kerning(1.2)
                if isLoadingPartners {
                    ProgressView().tint(Color.rfGold).scaleEffect(0.7).padding(.leading, 4)
                }
            }

            if !isLoadingPartners && partners.isEmpty {
                Text("No partners yet. Invite someone below.")
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.45))
                    .padding(.vertical, 12)
            } else {
                VStack(spacing: 0) {
                    ForEach(partners) { partner in
                        HStack(spacing: 14) {
                            ZStack {
                                Circle()
                                    .fill(Color.rfGold.opacity(0.14))
                                    .frame(width: 40, height: 40)
                                Text(partner.initials)
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundStyle(Color.rfGold)
                            }
                            VStack(alignment: .leading, spacing: 3) {
                                HStack(spacing: 6) {
                                    Text(partner.name)
                                        .font(.system(size: 15, weight: .semibold))
                                        .foregroundStyle(.white)
                                    if partner.isPrimary {
                                        Text("Primary")
                                            .font(.system(size: 10, weight: .semibold))
                                            .foregroundStyle(Color.rfGold)
                                            .padding(.horizontal, 6)
                                            .padding(.vertical, 2)
                                            .background(Capsule().fill(Color.rfGold.opacity(0.14)))
                                    }
                                }
                                Text(partner.email)
                                    .font(.system(size: 12))
                                    .foregroundStyle(Color.white.opacity(0.42))
                            }
                            Spacer()
                            Button {
                                setPrimary(partner: partner)
                            } label: {
                                Image(systemName: partner.isPrimary ? "star.fill" : "star")
                                    .font(.system(size: 15))
                                    .foregroundStyle(partner.isPrimary ? Color.rfGold : Color.white.opacity(0.30))
                            }
                            .padding(.trailing, 4)
                            Button {
                                partnerToRemove = partner
                                showRemoveConfirm = true
                            } label: {
                                Text("Remove")
                                    .font(.system(size: 12, weight: .semibold))
                                    .foregroundStyle(red)
                                    .padding(.horizontal, 10)
                                    .padding(.vertical, 5)
                                    .background(Capsule().fill(red.opacity(0.12)))
                            }
                        }
                        .padding(.vertical, 12)
                        .padding(.horizontal, 16)
                        if partner.id != partners.last?.id {
                            Divider().overlay(Color.white.opacity(0.07)).padding(.leading, 70)
                        }
                    }
                }
                .background(
                    RoundedRectangle(cornerRadius: 16)
                        .fill(Color.white.opacity(0.055))
                        .overlay(RoundedRectangle(cornerRadius: 16).stroke(Color.white.opacity(0.08), lineWidth: 1))
                )
                .clipShape(RoundedRectangle(cornerRadius: 16))
            }
        }
    }

    private var inviteSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("ADD NEW PARTNER")
                .font(.system(size: 11, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.2)

            HStack(spacing: 12) {
                Image(systemName: "envelope.fill")
                    .font(.system(size: 15))
                    .foregroundStyle(emailFocused ? Color.rfGold : Color.white.opacity(0.45))
                    .frame(width: 22)
                TextField("", text: $inviteEmail,
                          prompt: Text("Partner's email address").foregroundColor(Color.white.opacity(0.38)))
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
                    .overlay(RoundedRectangle(cornerRadius: 14)
                        .stroke(emailFocused ? Color.rfGold : Color.white.opacity(0.10), lineWidth: 1.5))
            )
            .animation(.easeInOut(duration: 0.18), value: emailFocused)

            if let err = errorMsg {
                Text(err)
                    .font(.system(size: 13))
                    .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                    .transition(.opacity)
            }

            Button(action: sendInvite) {
                HStack(spacing: 10) {
                    Image(systemName: sentSuccess ? "checkmark.circle.fill" : "paperplane.fill")
                        .font(.system(size: 16))
                    Text(sentSuccess ? "Invitation Sent!" : "Send Invitation")
                        .font(.system(size: 16, weight: .semibold))
                }
                .foregroundStyle(sentSuccess ? green : Color.rfNavy)
                .frame(maxWidth: .infinity)
                .frame(height: 52)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(sentSuccess ? green.opacity(0.15) : (canSendInvite ? Color.rfGold : Color.rfGold.opacity(0.35)))
                )
            }
            .disabled(!canSendInvite)
            .animation(.easeInOut(duration: 0.2), value: sentSuccess)
        }
    }

    private var canSendInvite: Bool {
        inviteEmail.contains("@") && inviteEmail.contains(".") && !sentSuccess && !isSending
    }

    @MainActor
    private func loadPartners() async {
        guard APIClient.shared.isAuthenticated else { return }
        isLoadingPartners = true
        defer { isLoadingPartners = false }
        if let rels = try? await APIClient.shared.listRelationships() {
            partners = rels.map { r in
                PartnerItem(relationshipID: r.id, name: r.partner.name,
                            email: r.partner.email, status: r.status, isPrimary: r.isPrimary)
            }
        }
    }

    private func setPrimary(partner: PartnerItem) {
        guard !partner.isPrimary else { return }
        Task {
            guard APIClient.shared.isAuthenticated else { return }
            try? await APIClient.shared.setPrimaryPartner(relationshipID: partner.relationshipID)
            await MainActor.run {
                for i in partners.indices {
                    partners[i].isPrimary = (partners[i].id == partner.id)
                }
            }
        }
    }

    private func sendInvite() {
        guard canSendInvite else { return }
        emailFocused = false
        errorMsg = nil
        isSending = true
        let email = inviteEmail
        Task {
            do {
                try await APIClient.shared.invitePartner(email: email)
                await MainActor.run {
                    withAnimation { sentSuccess = true }
                    isSending = false
                    inviteEmail = ""
                }
                // Reload partner list to include the new connection if they had an account.
                await loadPartners()
            } catch {
                await MainActor.run {
                    errorMsg = error.localizedDescription
                    isSending = false
                }
            }
        }
    }
}

private struct GroupItem: Identifiable {
    let id:   Int
    let name: String
}

// MARK: - Group Invite Sheet

private struct GroupInviteSheet: View {
    let group: GroupItem
    @Environment(\.dismiss) private var dismiss
    @FocusState private var emailFocused: Bool
    @State private var inviteEmail  = ""
    @State private var isSending    = false
    @State private var sentSuccess  = false
    @State private var errorMsg: String?

    private let green = Color(red: 0.20, green: 0.78, blue: 0.45)

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                dragIndicator()
                sheetHeader(title: "Invite Member", subtitle: "Add someone to \(group.name)")
                Divider().overlay(Color.white.opacity(0.08))
                VStack(spacing: 16) {
                    HStack(spacing: 12) {
                        Image(systemName: "envelope.fill")
                            .font(.system(size: 15))
                            .foregroundStyle(emailFocused ? Color.rfGold : Color.white.opacity(0.45))
                            .frame(width: 22)
                        TextField("", text: $inviteEmail,
                                  prompt: Text("Member's email address").foregroundColor(Color.white.opacity(0.38)))
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
                            .overlay(RoundedRectangle(cornerRadius: 14)
                                .stroke(emailFocused ? Color.rfGold : Color.white.opacity(0.10), lineWidth: 1.5))
                    )
                    .animation(.easeInOut(duration: 0.18), value: emailFocused)

                    if let err = errorMsg {
                        Text(err).font(.system(size: 13)).foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                    }

                    Button(action: sendInvite) {
                        HStack(spacing: 10) {
                            Image(systemName: sentSuccess ? "checkmark.circle.fill" : "paperplane.fill")
                                .font(.system(size: 16))
                            Text(sentSuccess ? "Invitation Sent!" : "Send Invitation")
                                .font(.system(size: 16, weight: .semibold))
                        }
                        .foregroundStyle(sentSuccess ? green : Color.rfNavy)
                        .frame(maxWidth: .infinity)
                        .frame(height: 52)
                        .background(RoundedRectangle(cornerRadius: 14)
                            .fill(sentSuccess ? green.opacity(0.15) : (canSend ? Color.rfGold : Color.rfGold.opacity(0.35))))
                    }
                    .disabled(!canSend)
                    .animation(.easeInOut(duration: 0.2), value: sentSuccess)
                }
                .padding(24)
                Spacer()
            }
        }
    }

    private var canSend: Bool {
        inviteEmail.contains("@") && inviteEmail.contains(".") && !sentSuccess && !isSending
    }

    private func sendInvite() {
        guard canSend else { return }
        emailFocused = false
        errorMsg = nil
        isSending = true
        let email = inviteEmail
        Task {
            do {
                try await APIClient.shared.groupEmailInvite(groupID: group.id, email: email)
                await MainActor.run {
                    withAnimation { sentSuccess = true }
                    isSending = false
                    inviteEmail = ""
                }
            } catch {
                await MainActor.run {
                    errorMsg = error.localizedDescription
                    isSending = false
                }
            }
        }
    }
}

// MARK: - Edit Profile Sheet

private struct EditProfileSheet: View {
    @Environment(\.dismiss) private var dismiss
    @FocusState private var nameFocused:  Bool
    @FocusState private var emailFocused: Bool
    @State private var name:      String
    @State private var email:     String
    @State private var isSaving   = false
    @State private var saveError: String?

    init() {
        let user = AuthState.shared.currentUser
        _name  = State(initialValue: user?.name  ?? UserDefaults.standard.string(forKey: "userName")  ?? "")
        _email = State(initialValue: user?.email ?? UserDefaults.standard.string(forKey: "userEmail") ?? "")
    }

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                dragIndicator()
                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Edit Profile")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Update your name or email")
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

                VStack(spacing: 14) {
                    inputField(icon: "person.fill", placeholder: "Full name",
                               text: $name, focused: $nameFocused,
                               keyboardType: .default, caps: .words)
                    inputField(icon: "envelope.fill", placeholder: "Email address",
                               text: $email, focused: $emailFocused,
                               keyboardType: .emailAddress, caps: .never)

                    if let err = saveError {
                        Text(err).font(.system(size: 13)).foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                    }

                    Button(action: save) {
                        HStack(spacing: 10) {
                            if isSaving { ProgressView().tint(Color.rfNavy).scaleEffect(0.9) }
                            Text(isSaving ? "Saving…" : "Save Changes")
                                .font(.system(size: 16, weight: .semibold))
                        }
                        .foregroundStyle(canSave ? Color.rfNavy : Color.white.opacity(0.35))
                        .frame(maxWidth: .infinity)
                        .frame(height: 52)
                        .background(RoundedRectangle(cornerRadius: 14)
                            .fill(canSave ? Color.rfGold : Color.white.opacity(0.08)))
                    }
                    .disabled(!canSave)
                }
                .padding(24)
                Spacer()
            }
        }
        .presentationDetents([.medium])
    }

    private var canSave: Bool {
        !name.trimmingCharacters(in: .whitespaces).isEmpty &&
        email.contains("@") && email.contains(".") && !isSaving
    }

    private func save() {
        guard canSave else { return }
        nameFocused  = false
        emailFocused = false
        saveError = nil
        isSaving = true
        let n = name.trimmingCharacters(in: .whitespaces)
        let e = email.trimmingCharacters(in: .whitespaces).lowercased()
        Task {
            do {
                _ = try await APIClient.shared.updateMe(name: n, email: e)
                AuthState.shared.updateProfile(name: n, email: e)
                await MainActor.run { isSaving = false }
                dismiss()
            } catch {
                await MainActor.run {
                    saveError = error.localizedDescription
                    isSaving = false
                }
            }
        }
    }

    private func inputField(icon: String, placeholder: String,
                            text: Binding<String>,
                            focused: FocusState<Bool>.Binding,
                            keyboardType: UIKeyboardType,
                            caps: TextInputAutocapitalization) -> some View {
        HStack(spacing: 12) {
            Image(systemName: icon)
                .font(.system(size: 15))
                .foregroundStyle(focused.wrappedValue ? Color.rfGold : Color.white.opacity(0.45))
                .frame(width: 22)
            TextField("", text: text,
                      prompt: Text(placeholder).foregroundColor(Color.white.opacity(0.38)))
                .keyboardType(keyboardType)
                .textInputAutocapitalization(caps)
                .autocorrectionDisabled()
                .foregroundStyle(.white)
                .focused(focused)
        }
        .padding(.horizontal, 18)
        .frame(height: 54)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(Color.white.opacity(focused.wrappedValue ? 0.11 : 0.07))
                .overlay(RoundedRectangle(cornerRadius: 14)
                    .stroke(focused.wrappedValue ? Color.rfGold : Color.white.opacity(0.10), lineWidth: 1.5))
        )
        .animation(.easeInOut(duration: 0.18), value: focused.wrappedValue)
    }
}

private struct PartnerItem: Identifiable {
    let id             = UUID()
    let relationshipID: Int
    let name:           String
    let email:          String
    var status:         String = "accepted"
    var isPrimary:      Bool   = false
    var initials: String {
        name.components(separatedBy: " ")
            .compactMap(\.first).prefix(2).map(String.init).joined().uppercased()
    }
}

// MARK: - Manage Groups

private struct ManageGroupsView: View {
    @AppStorage("primaryGroupID") private var primaryGroupID = 0
    @Environment(\.dismiss) private var dismiss
    @State private var newGroupName      = ""
    @State private var isCreating        = false
    @State private var createError: String?
    @State private var groupToLeave: GroupItem?
    @State private var showLeaveConfirm  = false
    @State private var groupForInvite: GroupItem?
    @State private var showInviteSheet   = false
    @State private var groups: [GroupItem] = []
    @FocusState private var groupNameFocused: Bool

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                dragIndicator()
                sheetHeader(title: "Manage Groups", subtitle: "Groups you belong to")
                Divider().overlay(Color.white.opacity(0.08))
                ScrollView(showsIndicators: false) {
                    VStack(spacing: 20) {
                        groupsListSection
                        createGroupSection
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 20)
                }
            }
        }
        .alert("Leave Group", isPresented: $showLeaveConfirm) {
            Button("Leave", role: .destructive) {
                guard let g = groupToLeave else { return }
                Task {
                    try? await APIClient.shared.leaveGroup(groupID: g.id)
                    await MainActor.run {
                        groups.removeAll { $0.id == g.id }
                        if primaryGroupID == g.id {
                            UserDefaults.standard.removeObject(forKey: "primaryGroupID")
                        }
                    }
                }
            }
            Button("Cancel", role: .cancel) { }
        } message: {
            Text("You will be removed from \(groupToLeave?.name ?? "this group") and all members will be notified.")
        }
        .sheet(isPresented: $showInviteSheet) {
            if let g = groupForInvite {
                GroupInviteSheet(group: g)
                    .presentationDetents([.medium])
            }
        }
        .task {
            guard let remoteGroups = try? await APIClient.shared.listMyGroups(),
                  !remoteGroups.isEmpty else { return }
            groups = remoteGroups.map { GroupItem(id: $0.id, name: $0.name) }
            if primaryGroupID == 0 { primaryGroupID = groups[0].id }
        }
    }

    private var groupsListSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("YOUR GROUPS")
                .font(.system(size: 11, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.2)

            if groups.isEmpty {
                Text("You're not in any groups yet.")
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.45))
                    .padding(.vertical, 12)
            } else {
                VStack(spacing: 0) {
                    ForEach(groups) { group in
                        VStack(spacing: 0) {
                            HStack(spacing: 14) {
                                ZStack {
                                    RoundedRectangle(cornerRadius: 10)
                                        .fill(Color.rfGold.opacity(0.14))
                                        .frame(width: 40, height: 40)
                                    Image(systemName: "person.3.fill")
                                        .font(.system(size: 16))
                                        .foregroundStyle(Color.rfGold)
                                }
                                Text(group.name)
                                    .font(.system(size: 15, weight: .semibold))
                                    .foregroundStyle(.white)
                                Spacer()
                            }
                            .padding(.vertical, 14)
                            .padding(.horizontal, 16)

                            HStack(spacing: 10) {
                                actionChip(label: "Invite Member", icon: "person.badge.plus") {
                                    groupForInvite = group
                                    showInviteSheet = true
                                }
                                actionChip(label: "Leave Group", icon: "rectangle.portrait.and.arrow.right", isDestructive: true) {
                                    groupToLeave = group
                                    showLeaveConfirm = true
                                }
                            }
                            .padding(.horizontal, 16)
                            .padding(.bottom, 14)
                        }
                    }
                }
                .background(
                    RoundedRectangle(cornerRadius: 16)
                        .fill(Color.white.opacity(0.055))
                        .overlay(RoundedRectangle(cornerRadius: 16).stroke(Color.white.opacity(0.08), lineWidth: 1))
                )
                .clipShape(RoundedRectangle(cornerRadius: 16))
            }
        }
    }

    private var createGroupSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("CREATE NEW GROUP")
                .font(.system(size: 11, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.2)

            HStack(spacing: 12) {
                Image(systemName: "person.3.fill")
                    .font(.system(size: 15))
                    .foregroundStyle(groupNameFocused ? Color.rfGold : Color.white.opacity(0.45))
                    .frame(width: 22)
                TextField("", text: $newGroupName,
                          prompt: Text("Group name (e.g. Iron Brotherhood)").foregroundColor(Color.white.opacity(0.38)))
                    .textInputAutocapitalization(.words)
                    .foregroundStyle(.white)
                    .focused($groupNameFocused)
                    .submitLabel(.done)
                    .onSubmit { createGroup() }
            }
            .padding(.horizontal, 18)
            .frame(height: 54)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color.white.opacity(groupNameFocused ? 0.11 : 0.07))
                    .overlay(RoundedRectangle(cornerRadius: 14)
                        .stroke(groupNameFocused ? Color.rfGold : Color.white.opacity(0.10), lineWidth: 1.5))
            )
            .animation(.easeInOut(duration: 0.18), value: groupNameFocused)

            if let err = createError {
                Text(err)
                    .font(.system(size: 13))
                    .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                    .transition(.opacity)
            }

            Button(action: createGroup) {
                HStack(spacing: 10) {
                    if isCreating {
                        ProgressView().tint(Color.rfNavy).scaleEffect(0.9)
                    } else {
                        Image(systemName: "plus.circle.fill")
                            .font(.system(size: 16))
                    }
                    Text(isCreating ? "Creating…" : "Create Group")
                        .font(.system(size: 16, weight: .semibold))
                }
                .foregroundStyle(canCreate ? Color.rfNavy : Color.white.opacity(0.35))
                .frame(maxWidth: .infinity)
                .frame(height: 52)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(canCreate ? Color.rfGold : Color.white.opacity(0.08))
                )
            }
            .disabled(!canCreate)
        }
    }

    private var canCreate: Bool {
        !newGroupName.trimmingCharacters(in: .whitespaces).isEmpty && !isCreating
    }

    private func createGroup() {
        guard canCreate else { return }
        groupNameFocused = false
        createError = nil
        isCreating = true
        let name = newGroupName.trimmingCharacters(in: .whitespaces)
        Task {
            do {
                let group = try await APIClient.shared.createGroup(name: name)
                await MainActor.run {
                    primaryGroupID = group.id
                    groups.append(GroupItem(id: group.id, name: group.name))
                    newGroupName = ""
                    isCreating = false
                }
            } catch {
                await MainActor.run {
                    createError = error.localizedDescription
                    isCreating = false
                }
            }
        }
    }

    private func actionChip(label: String, icon: String, isDestructive: Bool = false, action: @escaping () -> Void) -> some View {
        let tint = isDestructive ? Color(red: 0.90, green: 0.30, blue: 0.30) : Color.rfGold
        return Button(action: action) {
            HStack(spacing: 5) {
                Image(systemName: icon)
                    .font(.system(size: 11, weight: .semibold))
                Text(label)
                    .font(.system(size: 12, weight: .semibold))
            }
            .foregroundStyle(tint)
            .padding(.horizontal, 10)
            .padding(.vertical, 6)
            .background(Capsule().fill(tint.opacity(0.12)))
        }
    }
}

// MARK: - How It Works

private struct HowItWorksView: View {
    @Environment(\.dismiss) private var dismiss

    private let steps: [(icon: String, title: String, body: String)] = [
        ("person.2.fill",
         "Choose Your Partners",
         "Invite a trusted friend or join a small group of men who will walk honestly alongside you. Accountability works because you choose people who genuinely care."),
        ("shield.fill",
         "Start Monitoring",
         "Remain Faithful quietly monitors your screen activity in the background. When something concerning is detected, your partners are notified — no hiding, no excuses."),
        ("checkmark.shield.fill",
         "Stay Accountable",
         "Review your activity together, have honest conversations, and spur one another toward freedom. Iron sharpens iron — and so one man sharpens another."),
    ]

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                dragIndicator()

                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("How It Works")
                            .font(.system(size: 24, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Three steps to freedom")
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
                .padding(.bottom, 24)

                Divider().overlay(Color.white.opacity(0.08))

                ScrollView(showsIndicators: false) {
                    VStack(spacing: 20) {
                        ForEach(Array(steps.enumerated()), id: \.offset) { index, step in
                            stepCard(number: index + 1, icon: step.icon, title: step.title, body: step.body)
                        }

                        Text("\"Iron sharpens iron, and one man sharpens another.\"")
                            .font(.system(size: 14, weight: .medium, design: .serif))
                            .foregroundStyle(Color.white.opacity(0.55))
                            .multilineTextAlignment(.center)
                            .lineSpacing(4)
                            .padding(.top, 8)

                        Text("— Proverbs 27:17")
                            .font(.system(size: 12))
                            .foregroundStyle(Color.rfGold.opacity(0.65))
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 24)
                }
            }
        }
    }

    private func stepCard(number: Int, icon: String, title: String, body: String) -> some View {
        HStack(alignment: .top, spacing: 16) {
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.13))
                    .frame(width: 52, height: 52)
                Image(systemName: icon)
                    .font(.system(size: 20))
                    .foregroundStyle(Color.rfGold)
            }

            VStack(alignment: .leading, spacing: 8) {
                HStack(spacing: 8) {
                    Text("0\(number)")
                        .font(.system(size: 11, weight: .bold))
                        .foregroundStyle(Color.rfGold.opacity(0.70))
                        .kerning(1.0)
                    Text(title)
                        .font(.system(size: 16, weight: .bold, design: .serif))
                        .foregroundStyle(.white)
                }
                Text(body)
                    .font(.system(size: 13))
                    .foregroundStyle(Color.white.opacity(0.60))
                    .lineSpacing(4)
                    .fixedSize(horizontal: false, vertical: true)
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.white.opacity(0.055))
                .overlay(RoundedRectangle(cornerRadius: 16)
                    .stroke(Color.rfGold.opacity(0.14), lineWidth: 1))
        )
    }
}

// MARK: - Shared sheet helpers

private func dragIndicator() -> some View {
    RoundedRectangle(cornerRadius: 3)
        .fill(Color.white.opacity(0.18))
        .frame(width: 40, height: 4)
        .padding(.top, 12)
        .padding(.bottom, 20)
}

private func sheetHeader(title: String, subtitle: String) -> some View {
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
    }
    .padding(.horizontal, 24)
    .padding(.bottom, 20)
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
    @State private var events:    [ActivityEvent] = []
    @State private var isLoading = false
    @State private var loadError: String?

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

                if isLoading {
                    Spacer()
                    ProgressView().tint(Color.rfGold).scaleEffect(1.2)
                    Spacer()
                } else if let err = loadError {
                    Spacer()
                    VStack(spacing: 12) {
                        Image(systemName: "exclamationmark.triangle")
                            .font(.system(size: 32))
                            .foregroundStyle(Color.rfGold.opacity(0.55))
                        Text(err)
                            .font(.system(size: 14))
                            .foregroundStyle(Color.white.opacity(0.5))
                            .multilineTextAlignment(.center)
                        Button("Try Again") { Task { await load() } }
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundStyle(Color.rfGold)
                    }
                    .padding(.horizontal, 40)
                    Spacer()
                } else if events.isEmpty {
                    Spacer()
                    VStack(spacing: 12) {
                        Image(systemName: "checkmark.shield.fill")
                            .font(.system(size: 38))
                            .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45).opacity(0.55))
                        Text("No flagged activity")
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundStyle(.white)
                        Text("Your monitoring history will appear here.")
                            .font(.system(size: 14))
                            .foregroundStyle(Color.white.opacity(0.45))
                            .multilineTextAlignment(.center)
                    }
                    Spacer()
                } else {
                    ScrollView(showsIndicators: false) {
                        LazyVStack(spacing: 8) {
                            ForEach(events) { event in
                                LogRow(event: event)
                            }
                        }
                        .padding(.horizontal, 20)
                        .padding(.vertical, 16)
                    }
                }
            }
        }
        .task { await load() }
    }

    @MainActor
    private func load() async {
        guard APIClient.shared.isAuthenticated else { return }
        isLoading = true
        loadError = nil
        defer { isLoading = false }
        do {
            let remote = try await APIClient.shared.listEvents()
            events = remote.compactMap { ActivityEvent.from(remote: $0) }
        } catch {
            loadError = "Couldn't load history. Tap to retry."
        }
    }
}

private struct LogRow: View {
    let event: ActivityEvent

    var body: some View {
        HStack(spacing: 14) {
            ZStack {
                Circle()
                    .fill(event.category.tint.opacity(0.14))
                    .frame(width: 42, height: 42)
                Image(systemName: event.category.icon)
                    .font(.system(size: 16))
                    .foregroundStyle(event.category.tint)
            }

            VStack(alignment: .leading, spacing: 3) {
                Text(event.category.rawValue)
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(.white)
                Text(event.description)
                    .font(.system(size: 12))
                    .foregroundStyle(Color.white.opacity(0.48))
                    .lineLimit(1)
            }

            Spacer(minLength: 8)

            VStack(alignment: .trailing, spacing: 5) {
                Text(event.severity.rawValue.uppercased())
                    .font(.system(size: 10, weight: .bold))
                    .kerning(0.5)
                    .foregroundStyle(event.severity.color)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 3)
                    .background(Capsule().fill(event.severity.color.opacity(0.14)))
                Text(event.fullTimestamp)
                    .font(.system(size: 11))
                    .foregroundStyle(Color.white.opacity(0.32))
            }
        }
        .padding(14)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(Color.white.opacity(0.055))
                .overlay(RoundedRectangle(cornerRadius: 14)
                    .stroke(Color.white.opacity(0.07), lineWidth: 1))
        )
    }
}

#Preview {
    SettingsView()
}
