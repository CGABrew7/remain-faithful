import SwiftUI
import FamilyControls

// Phase 1 Screen Time monitoring: FamilyActivityPicker + ManagedSettings shielding.
// The existing ReplayKit broadcast/deep-scan feature is a separate, independent tier
// and is not referenced or modified here.
struct ScreenTimeMonitoringView: View {
    @ObservedObject private var fcManager  = FamilyControlsManager.shared
    @ObservedObject private var selManager = ActivitySelectionManager.shared

    @State private var showingPicker = false
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            ScrollView(showsIndicators: false) {
                VStack(spacing: 24) {
                    authorizationCard

                    if fcManager.authorizationStatus == .approved {
                        selectionCard
                        blockingCard
                        if selManager.isMonitoringEnabled || selManager.isShieldingEnabled {
                            shieldingNoteCard
                        }
                    }
                }
                .padding(.horizontal, 20)
                .padding(.top, 16)
                .padding(.bottom, 48)
            }
        }
        .navigationBarBackButtonHidden(true)
        .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
                RFCloseButton { dismiss() }
            }
            ToolbarItem(placement: .principal) {
                Text("Blocking")
                    .font(.system(size: 17, weight: .semibold))
                    .foregroundStyle(.white)
            }
        }
        // FamilyActivityPicker is presented as a full-screen system sheet.
        .familyActivityPicker(isPresented: $showingPicker, selection: $selManager.selection)
        .onChange(of: selManager.selection) { _, _ in
            selManager.selectionDidChange()
        }
    }

    // MARK: - Authorization card

    private var authorizationCard: some View {
        STSection(title: "SCREEN TIME ACCESS") {
            HStack(spacing: 14) {
                stBadge("lock.shield.fill", tint: authTint)

                VStack(alignment: .leading, spacing: 3) {
                    Text("Screen Time Permission")
                        .font(.system(size: 15, weight: .medium))
                        .foregroundStyle(.white)
                    Text(authStatusLabel)
                        .font(.system(size: 12))
                        .foregroundStyle(authTint)
                }
                Spacer()

                if fcManager.authorizationStatus == .approved {
                    Image(systemName: "checkmark.circle.fill")
                        .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45))
                        .font(.system(size: 20))
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 14)

            if fcManager.authorizationStatus != .approved {
                STDivider()

                Button {
                    Task { await fcManager.requestAuthorization() }
                } label: {
                    Text(fcManager.authorizationStatus == .denied
                         ? "Open Screen Time Settings"
                         : "Enable Screen Time Access")
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(Color.rfGold)
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 14)
                }
                .simultaneousGesture(TapGesture().onEnded {
                    // If denied, iOS will not re-show the prompt; open Settings instead.
                    if fcManager.authorizationStatus == .denied {
                        if let url = URL(string: UIApplication.openSettingsURLString) {
                            UIApplication.shared.open(url)
                        }
                    }
                })

                Text(authHintText)
                    .font(.system(size: 12))
                    .foregroundStyle(Color.white.opacity(0.40))
                    .multilineTextAlignment(.center)
                    .padding(.horizontal, 16)
                    .padding(.bottom, 14)
            }
        }
    }

    // MARK: - Selection card

    private var selectionCard: some View {
        STSection(title: "APPS & CATEGORIES") {
            Button {
                showingPicker = true
            } label: {
                HStack(spacing: 14) {
                    stBadge("square.grid.2x2.fill",
                            tint: Color(red: 0.52, green: 0.36, blue: 0.92))

                    VStack(alignment: .leading, spacing: 3) {
                        Text("Choose apps")
                            .font(.system(size: 15, weight: .medium))
                            .foregroundStyle(.white)
                        Text(selManager.selectionSummary)
                            .font(.system(size: 12))
                            .foregroundStyle(Color.white.opacity(0.45))
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

    // MARK: - Monitoring + blocking card

    private var blockingCard: some View {
        STSection(title: "ACTIONS") {
            // Monitor toggle — DeviceActivity fires after 1 min of app usage per day
            HStack(spacing: 14) {
                stBadge("chart.bar.fill",
                        tint: selManager.isMonitoringEnabled
                            ? Color(red: 0.28, green: 0.56, blue: 0.95)
                            : Color.white.opacity(0.35))
                VStack(alignment: .leading, spacing: 3) {
                    Text("Alerts")
                        .font(.system(size: 15, weight: .medium))
                        .foregroundStyle(.white)
                    Text("Tell your partner after 1 minute of use, once a day")
                        .font(.system(size: 12))
                        .foregroundStyle(Color.white.opacity(0.45))
                }
                Spacer()
                Toggle("Alerts", isOn: Binding(
                    get: { selManager.isMonitoringEnabled },
                    set: { selManager.setMonitoringEnabled($0) }
                ))
                .labelsHidden()
                .tint(Color(red: 0.28, green: 0.56, blue: 0.95))
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)

            Divider().overlay(Color.white.opacity(0.07)).padding(.leading, 64)

            // Shield toggle — ManagedSettings blocks the app immediately
            HStack(spacing: 14) {
                stBadge("hand.raised.fill",
                        tint: selManager.isShieldingEnabled
                            ? Color(red: 0.90, green: 0.30, blue: 0.30)
                            : Color.white.opacity(0.35))

                VStack(alignment: .leading, spacing: 3) {
                    Text("Blocking")
                        .font(.system(size: 15, weight: .medium))
                        .foregroundStyle(.white)
                    Text("Shows Apple's Screen Time shield when the app opens")
                        .font(.system(size: 12))
                        .foregroundStyle(Color.white.opacity(0.45))
                }
                Spacer()

                Toggle("Blocking", isOn: Binding(
                    get:  { selManager.isShieldingEnabled },
                    set:  { selManager.setShieldingEnabled($0) }
                ))
                .labelsHidden()
                .tint(Color(red: 0.90, green: 0.30, blue: 0.30))
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
        }
    }

    // MARK: - Active shielding note

    private var shieldingNoteCard: some View {
        HStack(alignment: .top, spacing: 10) {
            Image(systemName: "info.circle.fill")
                .font(.system(size: 14))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .padding(.top, 1)

            Text("Alerts tell your partner after one minute of use, once a day. Blocking shows Apple's Screen Time shield. You can still ask Apple for an override inside the blocked app. Deep Scan does not see protected video such as Netflix.")
                .font(.system(size: 12))
                .foregroundStyle(Color.white.opacity(0.45))
                .fixedSize(horizontal: false, vertical: true)
        }
        .padding(.horizontal, 4)
    }

    // MARK: - Helpers

    private var authTint: Color {
        switch fcManager.authorizationStatus {
        case .approved:      return Color(red: 0.20, green: 0.78, blue: 0.45)
        case .denied:        return Color(red: 0.90, green: 0.30, blue: 0.30)
        case .notDetermined: return Color.rfGold
        @unknown default:    return Color.white
        }
    }

    private var authStatusLabel: String {
        switch fcManager.authorizationStatus {
        case .approved:      return "Allowed — choose apps next"
        case .denied:        return "Denied"
        case .notDetermined: return "Not yet enabled"
        @unknown default:    return "Unknown"
        }
    }

    private var authHintText: String {
        switch fcManager.authorizationStatus {
        case .denied:
            return "Screen Time access was denied. Tap above to open iOS Settings and re-enable it under Screen Time > Remain Faithful."
        default:
            return "Remain Faithful needs Screen Time access for Alerts and Blocking. Your partner cannot browse your phone. If you skip it, monitoring stays off."
        }
    }
}

// MARK: - Local section + badge components (match SettingsView style)

private struct STSection<Content: View>: View {
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

private struct STDivider: View {
    var body: some View {
        Divider()
            .overlay(Color.white.opacity(0.07))
            .padding(.leading, 64)
    }
}

private func stBadge(_ name: String, tint: Color) -> some View {
    ZStack {
        RoundedRectangle(cornerRadius: 8)
            .fill(tint.opacity(0.15))
            .frame(width: 34, height: 34)
        Image(systemName: name)
            .font(.system(size: 15))
            .foregroundStyle(tint)
    }
}

#Preview {
    NavigationStack {
        ScreenTimeMonitoringView()
    }
}
