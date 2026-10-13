import SwiftUI

// MARK: - Models

enum EventCategory: String {
    case adultContent = "Adult Content"
    case gambling     = "Gambling"
    case violence     = "Violence"
    case selfHarm     = "Self-Harm"
    case appUsage     = "App Use"
    case clean        = "All Clear"
    case tamper       = "Monitoring Change"
    case socialMedia  = "Social Media"
    case gaming       = "Gaming"
    case unknown      = "Alert"

    var icon: String {
        switch self {
        case .adultContent: "eye.slash.fill"
        case .gambling:     "dollarsign.circle.fill"
        case .violence:     "exclamationmark.triangle.fill"
        case .selfHarm:     "heart.text.square.fill"
        case .appUsage:     "apps.iphone"
        case .clean:        "checkmark.shield.fill"
        case .tamper:       "lock.slash.fill"
        case .socialMedia:  "bubble.left.and.bubble.right.fill"
        case .gaming:       "gamecontroller.fill"
        case .unknown:      "bell.fill"
        }
    }

    /// Never returns nil. An unrecognised server category becomes `.unknown`
    /// so the row still shows.
    init(apiString: String) {
        switch apiString.lowercased().replacingOccurrences(of: " ", with: "_") {
        case "adult_content": self = .adultContent
        case "gambling":      self = .gambling
        case "violence":      self = .violence
        case "self_harm":     self = .selfHarm
        case "app_usage":     self = .appUsage
        case "clean":         self = .clean
        case "social_media":  self = .socialMedia
        case "gaming":        self = .gaming
        case "monitoring_disabled", "shielding_disabled", "signed_out",
             "family_controls_revoked", "deep_scan_stopped", "lockout_disabled",
             "selection_cleared", "heartbeat_silence", "wrong_pin_attempt",
             "pin_removed", "pin_changed", "protection", "tamper":
            self = .tamper
        default:
            self = .unknown
        }
    }

    var tint: Color {
        switch self {
        case .adultContent: Color(red: 0.90, green: 0.25, blue: 0.30)
        case .gambling:     Color(red: 0.95, green: 0.58, blue: 0.12)
        case .violence:     Color(red: 0.85, green: 0.28, blue: 0.22)
        case .selfHarm:     Color(red: 0.95, green: 0.45, blue: 0.35)
        case .appUsage:     Color(red: 0.28, green: 0.56, blue: 0.95)
        case .clean:        Color(red: 0.20, green: 0.78, blue: 0.45)
        case .tamper:       Color(red: 0.95, green: 0.62, blue: 0.10)
        case .socialMedia:  Color(red: 0.28, green: 0.56, blue: 0.95)
        case .gaming:       Color(red: 0.60, green: 0.30, blue: 0.90)
        case .unknown:      Color.rfGold
        }
    }

    var conversationStarters: [String] {
        switch self {
        case .adultContent:
            return [
                "What was going on in your heart before this happened — were you stressed, lonely, or bored?",
                "Was this something you sought out, or did you stumble onto it? What happened in that moment?",
                "What does your relationship with purity look like right now? What would help most this week?",
                "How can I be praying for you specifically in this area?"
            ]
        case .gambling:
            return [
                "Are there financial pressures or anxieties you've been carrying lately that we should talk about?",
                "What draws you to this — the thrill, a sense of control, or something else entirely?",
                "Have you spoken with a pastor or counselor about this pattern? Would that be helpful?",
                "What would one concrete boundary look like here, and how can I help hold you to it?"
            ]
        case .socialMedia:
            return [
                "How much of your time and mental energy is social media taking up lately?",
                "Are you using it to connect, or to avoid something? What might you be escaping?",
                "How do you feel after time on these platforms — refreshed, or drained and restless?",
                "What's one limit we could agree on together around your phone use this week?"
            ]
        case .gaming:
            return [
                "What were you feeling right before the session started — stressed, restless, or wanting to check out?",
                "Do you think gaming is helping you unwind, or helping you avoid something harder?",
                "Late nights often point to something deeper. How's your sleep and your emotional state been?",
                "What would a healthy rhythm with gaming look like for you?"
            ]
        case .violence:
            return [
                "What were you looking at, and how are you doing right now?",
                "Was this something you went looking for, or did it show up in front of you?",
                "What would help you step away from this the next time it appears?",
                "How can I check in with you later today?"
            ]
        case .selfHarm:
            return [
                "Are you safe right now? I want to hear that from you.",
                "You don't have to explain everything. What do you need in the next hour?",
                "Who else can be with you — me, a pastor, or a counselor?",
                "Can we stay in touch until this feels smaller?"
            ]
        case .appUsage:
            return [
                "You opened an app you asked to be watched. What was going on just before?",
                "Was this a slip, or a choice you had already talked about?",
                "What would make the next hour easier?",
                "Do you want Blocking on for this app, or is an Alert enough?"
            ]
        case .clean:
            return [
                "What's been helping this week?",
                "Is there a pressure point we should name before it grows?",
                "How can I pray for you specifically?"
            ]
        case .tamper:
            return [
                "Monitoring changed on the phone. What happened?",
                "Was this intentional, or did something else turn it off?",
                "What do you want back on — Alerts, Blocking, or both?",
                "When can we confirm it's on again?"
            ]
        case .unknown:
            return [
                "What does this alert mean on your end?",
                "What was happening when it arrived?",
                "What would a helpful next step look like today?"
            ]
        }
    }
}

enum FlagSeverity: String {
    case low    = "Low"
    case medium = "Medium"
    case high   = "High"

    // Maps backend severity strings (informational/concerning/severe) to UI cases.
    init?(apiString: String) {
        switch apiString.lowercased() {
        case "informational": self = .low
        case "concerning":    self = .medium
        case "severe":        self = .high
        default:              return nil
        }
    }

    var color: Color {
        switch self {
        case .low:    Color(red: 0.20, green: 0.78, blue: 0.45)
        case .medium: Color(red: 0.95, green: 0.62, blue: 0.10)
        case .high:   Color(red: 0.90, green: 0.25, blue: 0.30)
        }
    }
}

struct ActivityEvent: Identifiable, Hashable {
    let id = UUID()
    let alertID:     Int?        // backend alert ID; nil for the user's own events
    let category:    EventCategory
    let severity:    FlagSeverity
    let description: String
    let minutesAgo:  Int
    var discussed:   Bool        // initial state loaded from API; toggled locally after PATCH

    init(alertID: Int? = nil, category: EventCategory, severity: FlagSeverity,
         description: String, minutesAgo: Int, discussed: Bool = false) {
        self.alertID     = alertID
        self.category    = category
        self.severity    = severity
        self.description = description
        self.minutesAgo  = minutesAgo
        self.discussed   = discussed
    }

    static func == (lhs: ActivityEvent, rhs: ActivityEvent) -> Bool { lhs.id == rhs.id }
    func hash(into hasher: inout Hasher) { hasher.combine(id) }

    var timeLabel: String {
        switch minutesAgo {
        case 0..<60:    return "\(minutesAgo)m ago"
        case 60..<1440: return "\(minutesAgo / 60)h ago"
        default:        return "\(minutesAgo / 1440)d ago"
        }
    }

    var fullTimestamp: String {
        let eventDate = Date().addingTimeInterval(-Double(minutesAgo) * 60)
        let cal = Calendar.current
        let fmt = DateFormatter()
        fmt.timeStyle = .short
        if cal.isDateInToday(eventDate)     { return "Today at \(fmt.string(from: eventDate))" }
        if cal.isDateInYesterday(eventDate) { return "Yesterday at \(fmt.string(from: eventDate))" }
        fmt.dateStyle = .medium
        fmt.timeStyle = .short
        return fmt.string(from: eventDate)
    }

    /// The server summary, plus one honest line. This does not invent a browsing
    /// story, and it does not claim Deep Scan can see protected video.
    var extendedDescription: String {
        let summary = description.trimmingCharacters(in: .whitespacesAndNewlines)
        let note = category.honestNote
        if summary.isEmpty { return note }
        if summary == note { return summary }
        return summary + "\n\n" + note
    }
}

private extension EventCategory {
    var honestNote: String {
        switch self {
        case .adultContent:
            return "This is the category your phone sent. It is not a copy of what was on screen. Deep Scan cannot see protected video such as Netflix."
        case .gambling:
            return "This is a gambling alert. The picture and the page stay on the phone."
        case .violence:
            return "This is a violence alert. Reach out. The screen itself was not uploaded."
        case .selfHarm:
            return "This is a self-harm alert. Check on them. The screen itself was not uploaded."
        case .appUsage:
            return "An app you chose was used for about a minute today. This is once a day, not every time the app opens."
        case .clean:
            return "Nothing in this check needed an alert."
        case .tamper:
            return "Alerts or Blocking changed on the phone. This is a notice, not screen content."
        case .socialMedia, .gaming:
            return "This is an older sample category. It is not a live screen capture."
        case .unknown:
            return "The phone sent an alert the app did not recognise. It is shown so it is not lost."
        }
    }
}

// MARK: - API conversion

extension ActivityEvent {
    /// Converts a RemoteEvent from the API into a local ActivityEvent.
    /// Unknown categories become `.unknown` instead of being dropped.
    static func from(remote r: RemoteEvent, alertID: Int? = nil, discussed: Bool = false) -> ActivityEvent {
        let category = EventCategory(rawValue: r.category) ?? EventCategory(apiString: r.category)
        let severity = FlagSeverity(rawValue: r.severity.capitalized) ?? FlagSeverity(apiString: r.severity) ?? .low
        let minutesAgo = Self.minutesAgo(from: r.timestamp)
        return ActivityEvent(alertID: alertID, category: category, severity: severity,
                             description: r.summary, minutesAgo: minutesAgo, discussed: discussed)
    }

    static func from(alert: RemoteAlert) -> ActivityEvent {
        if alert.kind == "protection",
           let raw = alert.protectionType, !raw.isEmpty {
            let category = EventCategory(apiString: raw)
            let summary = (alert.event?.summary ?? alert.summary ?? "")
                .trimmingCharacters(in: .whitespacesAndNewlines)
            let description = summary.isEmpty ? category.rawValue : summary
            let stamp: String
            if let event = alert.event, !event.timestamp.isEmpty {
                stamp = event.timestamp
            } else {
                stamp = alert.createdAt
            }
            let severity: FlagSeverity
            if let event = alert.event {
                severity = FlagSeverity(rawValue: event.severity.capitalized)
                    ?? FlagSeverity(apiString: event.severity)
                    ?? .medium
            } else {
                severity = (raw == "self_harm" || category == .selfHarm) ? .high : .medium
            }
            return ActivityEvent(alertID: alert.id, category: category, severity: severity,
                                 description: description, minutesAgo: minutesAgo(from: stamp))
        }
        if let event = alert.event {
            return from(remote: event, alertID: alert.id, discussed: false)
        }
        let raw = alert.protectionType ?? alert.kind ?? "protection"
        let category = EventCategory(apiString: raw)
        let summary = alert.summary?.trimmingCharacters(in: .whitespacesAndNewlines)
        let description = (summary?.isEmpty == false) ? summary! : category.rawValue
        let minutesAgo = Self.minutesAgo(from: alert.createdAt)
        let severity: FlagSeverity = (raw == "self_harm" || category == .selfHarm) ? .high : .medium
        return ActivityEvent(alertID: alert.id, category: category, severity: severity,
                             description: description, minutesAgo: minutesAgo)
    }

    static func minutesAgo(from timestamp: String) -> Int {
        guard !timestamp.isEmpty else { return 0 }
        let date = AccountDate.parse(timestamp) ?? Date()
        return max(0, Int(Date().timeIntervalSince(date) / 60))
    }
}

// MARK: - Notification payload conversion

extension ActivityEvent {
    /// Reconstruct an ActivityEvent from an APNs notification userInfo dict so
    /// tapping a CONTENT_FLAGGED push can open AlertDetailView directly.
    static func from(notificationPayload p: [AnyHashable: Any]) -> ActivityEvent? {
        let type = p["notification_type"] as? String ?? ""
        let catStr = (p["category"] as? String)
            ?? (p["protection_type"] as? String)
            ?? (p["type"] as? String)
        if catStr == nil && type != NotificationType.contentFlagged && type != NotificationType.protectionAlert {
            return nil
        }
        // A content push with no category still has nothing to open.
        if type == NotificationType.contentFlagged && catStr == nil { return nil }

        let resolved = catStr ?? "tamper"
        let category = EventCategory(rawValue: resolved) ?? EventCategory(apiString: resolved)
        let severityStr = p["severity"] as? String ?? (category == .selfHarm ? "severe" : "concerning")
        let severity = FlagSeverity(rawValue: severityStr.capitalized) ?? FlagSeverity(apiString: severityStr) ?? .low
        let description = (p["summary"] as? String)
            ?? (p["detail"] as? String)
            ?? category.rawValue
        let stamp = (p["timestamp"] as? String) ?? ""
        let minutesAgo = stamp.isEmpty ? 0 : minutesAgo(from: stamp)
        let alertID = (p["alert_id"] as? Int) ?? Int(p["alert_id"] as? String ?? "")
        return ActivityEvent(alertID: alertID, category: category, severity: severity,
                             description: description, minutesAgo: minutesAgo)
    }
}

// MARK: - Placeholder data (used in demo mode only)

private let sampleEvents: [ActivityEvent] = [
    .init(category: .adultContent, severity: .high,   description: "Flagged during browsing session",  minutesAgo: 23),
    .init(category: .socialMedia,  severity: .low,    description: "Extended social media usage",      minutesAgo: 145),
    .init(category: .gambling,     severity: .medium, description: "Sports betting site visited",      minutesAgo: 360),
    .init(category: .adultContent, severity: .medium, description: "Flagged search terms detected",    minutesAgo: 1440),
    .init(category: .gaming,       severity: .low,    description: "Late-night gaming session",        minutesAgo: 2160),
]

// MARK: - Dashboard root

struct DashboardView: View {
    @AppStorage("userName")                  private var userName             = ""
    @EnvironmentObject private var appState: AppState
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @ObservedObject private var tour = TourController.shared
    @Binding var showPanic: Bool

    @State private var events:           [ActivityEvent]  = []
    @State private var partnerAlerts:    [ActivityEvent]  = []
    @State private var streakDays:       Int              = 0
    @State private var streakBest:       Int              = 0
    @State private var streakMarks:      [StreakMark]     = Array(repeating: .beforeSignup, count: 7)
    @State private var isLoadingEvents   = false
    @State private var eventsLoadError:  String?
    @State private var partnerLoadError: String?
    @State private var isBroadcasting       = false
    @State private var showAppUsagePrompt   = false
    @State private var showBlockingSetup    = false
    @State private var showDeepScanHelp     = false
    @ObservedObject private var selection = ActivitySelectionManager.shared
    @ObservedObject private var familyControls = FamilyControlsManager.shared

    private var greeting: String {
        let h = Calendar.current.component(.hour, from: Date())
        if h < 12 { return "Good morning" }
        if h < 17 { return "Good afternoon" }
        return "Good evening"
    }

    private var firstName: String {
        let first = userName.components(separatedBy: " ").first ?? ""
        return first.isEmpty ? "there" : first
    }

    var body: some View {
        ZStack(alignment: .top) {
            Color.rfNavy.ignoresSafeArea()

            ScrollViewReader { proxy in
            ScrollView(showsIndicators: false) {
                VStack(spacing: 18) {
                    headerRow
                    StatusCard(
                        presentation: monitoring,
                        onAction: { handleStatusAction() }
                    )
                    .tourAnchor(.homeStatus)
                    if monitoring.kind == .noApps {
                        NoAppsBanner { showBlockingSetup = true }
                    }
                    if showAppUsagePrompt && !isBroadcasting {
                        AppUsagePromptBanner(
                            onStart:   {
                                ActivitySelectionManager.shared.clearAppUsagePrompt()
                                showAppUsagePrompt = false
                            },
                            onDismiss: {
                                ActivitySelectionManager.shared.clearAppUsagePrompt()
                                withAnimation(.easeOut(duration: 0.2)) { showAppUsagePrompt = false }
                            }
                        )
                        .transition(.move(edge: .top).combined(with: .opacity))
                    }
                    VerseCard()
                    StreakCard(days: streakDays, best: streakBest, marks: streakMarks)
                        .tourAnchor(.homeStreak)
                    ActivitySection(
                        title: "Alerts",
                        events: partnerAlerts,
                        isLoading: isLoadingEvents,
                        error: partnerLoadError,
                        emptyTitle: partnerEmptyTitle,
                        onRetry: { Task { await loadEvents() } }
                    )
                    ActivitySection(
                        title: "Your activity",
                        events: events,
                        isLoading: isLoadingEvents,
                        error: eventsLoadError,
                        emptyTitle: ownEmptyTitle,
                        onRetry: { Task { await loadEvents() } }
                    )
                    .tourAnchor(.homeFlags)
                    panicButton
                        .tourAnchor(.homeSupport)
                }
                .padding(.horizontal, 20)
                .padding(.top, 8)
                .padding(.bottom, 32)
            }
            .onAppear { scrollTour(proxy) }
            .onChange(of: tour.activeAnchor) { _, _ in scrollTour(proxy) }
            }

            // Partner flag banner — slides in from top when a foreground push arrives
            if let alert = appState.partnerFlagAlert {
                PartnerFlagBanner(
                    alert: alert,
                    onTap: {
                        withAnimation(.spring(response: 0.35, dampingFraction: 0.75)) {
                            appState.clearPartnerAlert()
                        }
                        if let event = alert.event {
                            AppState.shared.navigate(to: .alertDetail(event))
                        }
                    },
                    onDismiss: {
                        withAnimation(.spring(response: 0.35, dampingFraction: 0.75)) {
                            appState.clearPartnerAlert()
                        }
                    }
                )
                .transition(.move(edge: .top).combined(with: .opacity))
                .zIndex(10)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                // Low and medium banners leave after 6 seconds. A severe alert stays until it is opened or dismissed.
                .task(id: alert.id) {
                    if alert.event?.severity == .high { return }
                    try? await Task.sleep(for: .seconds(6))
                    withAnimation(.spring(response: 0.35, dampingFraction: 0.75)) {
                        appState.clearPartnerAlert()
                    }
                }
            }
        }
        .toolbar(.hidden, for: .navigationBar)
        .sheet(isPresented: $showBlockingSetup) {
            NavigationStack { ScreenTimeMonitoringView() }
        }
        .sheet(isPresented: $showDeepScanHelp) {
            BroadcastInstructionsSheet()
        }
        .task { await loadEvents() }
        .task {
            // Poll broadcast status every 5 seconds while the view is active.
            while !Task.isCancelled {
                isBroadcasting = EventProcessor.shared.isBroadcasting()
                try? await Task.sleep(for: .seconds(5))
            }
        }
        .onReceive(NotificationCenter.default.publisher(for: UIApplication.willEnterForegroundNotification)) { _ in
            isBroadcasting = EventProcessor.shared.isBroadcasting()
            if ActivitySelectionManager.shared.recentAppUsageDate != nil && !isBroadcasting {
                withAnimation(.spring(response: 0.4, dampingFraction: 0.8)) {
                    showAppUsagePrompt = true
                }
            }
        }
    }

    @MainActor
    private func loadEvents() async {
        if appState.isDemoMode {
            events     = sampleEvents
            partnerAlerts = [
                ActivityEvent(alertID: 1, category: .appUsage, severity: .medium,
                              description: "Chosen app used for about a minute", minutesAgo: 40),
                ActivityEvent(alertID: 2, category: .selfHarm, severity: .high,
                              description: "Self-harm content was flagged", minutesAgo: 90),
            ]
            streakDays = 14
            streakBest = 21
            streakMarks = [.beforeSignup, .beforeSignup, .clean, .clean, .flagged, .clean, .clean]
            return
        }
        guard APIClient.shared.isAuthenticated else { return }
        isLoadingEvents = true
        eventsLoadError = nil
        partnerLoadError = nil
        defer { isLoadingEvents = false }
        do {
            async let remoteEvents = APIClient.shared.listEvents()
            async let stats        = APIClient.shared.getUserStats()
            let (remote, userStats) = try await (remoteEvents, stats)
            events     = remote.map { ActivityEvent.from(remote: $0) }
            streakDays = userStats.streakDays
            streakBest = userStats.bestStreak
            let created = AccountDate.parse(AuthState.shared.currentUser?.createdAt)
            let raw: [StreakDayInput]
            if !userStats.weekStates.isEmpty {
                raw = userStats.weekStates
            } else if !userStats.week.isEmpty {
                raw = userStats.week
            } else {
                raw = Array(repeating: StreakDayInput.flagged, count: 7)
            }
            streakMarks = StreakWeek.marks(raw: raw, accountCreated: created, today: Date())
        } catch {
            eventsLoadError = "Couldn't load activity. Tap to retry."
        }
        do {
            let alerts = try await APIClient.shared.listAlerts()
            partnerAlerts = alerts.map { ActivityEvent.from(alert: $0) }
        } catch {
            partnerLoadError = "Couldn't load alerts. Tap to retry."
        }
    }

    private var monitoring: MonitoringPresentation {
        let selected = selection.selection.applicationTokens.count
            + selection.selection.categoryTokens.count
        let signedOut = UserDefaults(suiteName: "group.com.remainfaithful.app")?
            .bool(forKey: "extensionAuthFailed") ?? false
        return MonitoringPresentation.make(
            authorizationApproved: familyControls.authorizationStatus == .approved,
            selectedCount: selected,
            alertsOn: selection.isMonitoringEnabled,
            blockingOn: selection.isShieldingEnabled,
            deepScanOn: isBroadcasting,
            extensionSignedOut: signedOut
        )
    }

    private var setupUnfinished: Bool {
        monitoring.kind == .needsPermission || monitoring.kind == .noApps || monitoring.kind == .off
    }

    private var ownEmptyTitle: String {
        if setupUnfinished {
            return "No alerts yet. Choose apps under Blocking if you have not finished setup."
        }
        return "No alerts yet."
    }

    private var partnerEmptyTitle: String {
        "No alerts yet. This does not mean their Blocking is on."
    }

    private func handleStatusAction() {
        switch monitoring.kind {
        case .needsPermission:
            Task { await familyControls.requestAuthorization() }
        case .noApps, .off:
            showBlockingSetup = true
        case .on:
            showDeepScanHelp = true
        }
    }

    private func scrollTour(_ proxy: ScrollViewProxy) {
        guard let anchor = tour.activeAnchor else { return }
        switch anchor {
        case .homeStatus, .homeStreak, .homeFlags, .homeSupport:
            if reduceMotion {
                proxy.scrollTo(anchor, anchor: .center)
            } else {
                withAnimation(.easeInOut(duration: 0.35)) {
                    proxy.scrollTo(anchor, anchor: .center)
                }
            }
        default:
            break
        }
    }

    private var panicButton: some View {
        Button { showPanic = true } label: {
            HStack(spacing: 12) {
                Image(systemName: "heart.fill")
                    .font(.system(size: 15))
                Text("I Need Support Right Now")
                    .font(.system(size: 15, weight: .semibold))
            }
            .foregroundStyle(.white)
            .frame(maxWidth: .infinity)
            .frame(height: 54)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(
                        LinearGradient(
                            colors: [
                                Color(red: 0.45, green: 0.20, blue: 0.70),
                                Color(red: 0.28, green: 0.12, blue: 0.50)
                            ],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .overlay(
                        RoundedRectangle(cornerRadius: 14)
                            .stroke(Color.white.opacity(0.12), lineWidth: 1)
                    )
            )
        }
    }

    private var headerRow: some View {
        HStack(alignment: .center) {
            VStack(alignment: .leading, spacing: 3) {
                Text(greeting)
                    .font(.system(size: 13, weight: .medium))
                    .foregroundStyle(Color.white.opacity(0.45))
                Text(firstName)
                    .font(.system(size: 26, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
            }
            Spacer()
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.14))
                    .frame(width: 44, height: 44)
                Text(String(firstName.prefix(1)).uppercased())
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundStyle(Color.rfGold)
            }
        }
        .padding(.top, 4)
    }
}

// MARK: - Status card
//
// Composite of Screen Time permission, chosen apps, Alerts, and Blocking.
// Deep Scan is a subtitle. The card never says monitoring is paused while
// Alerts or Blocking are on.

private struct StatusCard: View {
    let presentation: MonitoringPresentation
    let onAction: () -> Void
    @State private var pulse: CGFloat = 1.0
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    private static let green   = Color(red: 0.18, green: 0.82, blue: 0.48)
    private static let amber   = Color(red: 0.95, green: 0.72, blue: 0.22)
    private static let bgOn    = LinearGradient(colors: [Color(red: 0.07, green: 0.40, blue: 0.26), Color(red: 0.05, green: 0.26, blue: 0.18)], startPoint: .topLeading, endPoint: .bottomTrailing)
    private static let bgOff   = LinearGradient(colors: [Color(red: 0.16, green: 0.20, blue: 0.32), Color(red: 0.11, green: 0.15, blue: 0.26)], startPoint: .topLeading, endPoint: .bottomTrailing)

    private var accent: Color { presentation.coreIsOn ? Self.green : Self.amber }

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            HStack(alignment: .top, spacing: 14) {
                ZStack {
                    if presentation.coreIsOn {
                        Circle()
                            .fill(accent.opacity(0.22))
                            .frame(width: 44, height: 44)
                            .scaleEffect(pulse)
                    }
                    Circle()
                        .fill(accent)
                        .frame(width: 11, height: 11)
                        .shadow(color: accent.opacity(0.9), radius: 5)
                }
                .frame(width: 44, height: 44)
                .accessibilityHidden(true)

                VStack(alignment: .leading, spacing: 5) {
                    Text(presentation.title)
                        .font(.title3.weight(.bold))
                        .foregroundStyle(.white)
                    Text(presentation.detail)
                        .font(.subheadline)
                        .foregroundStyle(Color.white.opacity(0.72))
                        .fixedSize(horizontal: false, vertical: true)
                    Text(presentation.deepScanLine)
                        .font(.footnote)
                        .foregroundStyle(Color.white.opacity(0.66))
                        .fixedSize(horizontal: false, vertical: true)
                    Text(PauseNotice.text)
                        .font(.footnote.weight(.medium))
                        .foregroundStyle(Self.amber)
                        .fixedSize(horizontal: false, vertical: true)
                }
                Spacer(minLength: 0)
            }
            .padding(.bottom, 16)

            RFButton(title: presentation.actionTitle, action: onAction)
        }
        .padding(20)
        .background(
            RoundedRectangle(cornerRadius: 20, style: .continuous)
                .fill(presentation.coreIsOn ? Self.bgOn : Self.bgOff)
                // Lit top edge over the accent rim, matching RFSurface.
                .overlay(
                    RoundedRectangle(cornerRadius: 20, style: .continuous)
                        .strokeBorder(
                            LinearGradient(
                                colors: [accent.opacity(0.55), accent.opacity(0.18)],
                                startPoint: .top, endPoint: .bottom
                            ),
                            lineWidth: 1
                        )
                )
                .compositingGroup()
                .shadow(color: .black.opacity(0.26), radius: 3,  x: 0, y: 2)
                .shadow(color: .black.opacity(0.20), radius: 18, x: 0, y: 10)
        )
        .clipShape(RoundedRectangle(cornerRadius: 20, style: .continuous))
        .accessibilityElement(children: .combine)
        .accessibilityLabel("\(presentation.title). \(presentation.detail). \(presentation.deepScanLine)")
        .accessibilityHint(presentation.actionTitle)
        .accessibilityAddTraits(.isButton)
        .onAppear {
            // Reduce Motion keeps the ring static rather than pulsing forever.
            guard !reduceMotion else { return }
            withAnimation(.easeInOut(duration: 1.1).repeatForever(autoreverses: true)) {
                pulse = 1.55
            }
        }
    }
}

private struct NoAppsBanner: View {
    let onChoose: () -> Void

    var body: some View {
        Button(action: onChoose) {
            HStack(spacing: 12) {
                Image(systemName: "apps.iphone")
                    .font(.body.weight(.semibold))
                    .foregroundStyle(Color.rfGold)
                    .frame(width: 44, height: 44)
                VStack(alignment: .leading, spacing: 2) {
                    Text("No apps selected")
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(.white)
                    Text("Choose apps to turn Alerts and Blocking on.")
                        .font(.footnote)
                        .foregroundStyle(Color.white.opacity(0.72))
                        .fixedSize(horizontal: false, vertical: true)
                }
                Spacer(minLength: 0)
                Image(systemName: "chevron.right")
                    .foregroundStyle(Color.white.opacity(0.55))
            }
            .padding(12)
            .frame(minHeight: 44)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color.rfGold.opacity(0.12))
                    .overlay(RoundedRectangle(cornerRadius: 14).stroke(Color.rfGold.opacity(0.45), lineWidth: 1))
            )
        }
        .buttonStyle(.plain)
        .accessibilityLabel("No apps selected. Choose apps.")
    }
}

// MARK: - Streak card

private struct StreakCard: View {
    let days: Int
    let best: Int
    let marks: [StreakMark]

    private var dayLabels: [String] { StreakWeek.letters(count: marks.count, endingOn: Date()) }
    private let flagRed   = Color(red: 0.88, green: 0.25, blue: 0.28)

    var body: some View {
        VStack(spacing: 0) {
            HStack(alignment: .top) {
                VStack(alignment: .leading, spacing: 2) {
                    HStack(alignment: .lastTextBaseline, spacing: 6) {
                        Text("\(days)")
                            .font(.system(size: 58, weight: .bold, design: .rounded))
                            .foregroundStyle(Color.rfGold)
                        Text("days")
                            .font(.system(size: 17, weight: .medium))
                            .foregroundStyle(Color.white.opacity(0.45))
                            .padding(.bottom, 8)
                    }
                    Text("clean streak")
                        .font(.system(size: 13, weight: .medium))
                        .foregroundStyle(Color.white.opacity(0.45))
                }

                Spacer()

                VStack(alignment: .trailing, spacing: 8) {
                    Text("Best: \(best)d")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundStyle(Color.white.opacity(0.38))
                }
                .padding(.top, 4)
            }
            .padding(.bottom, 18)

            Divider().overlay(Color.white.opacity(0.08))
                .padding(.bottom, 14)

            // 7-day dots
            HStack(spacing: 0) {
                ForEach(marks.indices, id: \.self) { i in
                    let mark = marks[i]
                    let letter = i < dayLabels.count ? dayLabels[i] : ""
                    VStack(spacing: 6) {
                        ZStack {
                            Circle()
                                .fill(markFill(mark))
                                .frame(width: 28, height: 28)
                            Image(systemName: markSymbol(mark))
                                .font(.system(size: 10, weight: .bold))
                                .foregroundStyle(markColor(mark))
                        }
                        Text(letter)
                            .font(.caption2.weight(.medium))
                            .foregroundStyle(Color.white.opacity(0.66))
                    }
                    .frame(maxWidth: .infinity)
                    .frame(minHeight: 44)
                    .accessibilityElement(children: .ignore)
                    .accessibilityLabel("\(StreakWeek.dayName(index: i, count: marks.count, endingOn: Date())), \(mark.accessibilityName)")
                }
            }
        }
        .padding(20)
        .background(card)
    }

    private func markFill(_ mark: StreakMark) -> Color {
        switch mark {
        case .clean: return Color.rfGold.opacity(0.18)
        case .flagged: return flagRed.opacity(0.18)
        case .beforeSignup: return Color.white.opacity(0.06)
        }
    }

    private func markColor(_ mark: StreakMark) -> Color {
        switch mark {
        case .clean: return Color.rfGold
        case .flagged: return flagRed
        case .beforeSignup: return Color.white.opacity(0.35)
        }
    }

    private func markSymbol(_ mark: StreakMark) -> String {
        switch mark {
        case .clean: return "checkmark"
        case .flagged: return "xmark"
        case .beforeSignup: return "circle"
        }
    }
}

// MARK: - Verse card

private struct VerseCard: View {
    private static let verses: [(text: String, ref: String)] = [
        ("Watch, stand fast in the faith, be brave, be strong.",
         "1 Corinthians 16:13"),
        ("Flee from sexual immorality. Every other sin a person commits is outside the body, but the sexually immoral person sins against his own body.",
         "1 Corinthians 6:18"),
        ("Submit yourselves therefore to God. Resist the devil, and he will flee from you.",
         "James 4:7"),
        ("No temptation has overtaken you that is not common to man. God is faithful, and he will not let you be tempted beyond your ability.",
         "1 Corinthians 10:13"),
        ("I have stored up your word in my heart, that I might not sin against you.",
         "Psalm 119:11"),
        ("Blessed is the man who remains steadfast under trial, for when he has stood the test he will receive the crown of life.",
         "James 1:12"),
        ("Iron sharpens iron, and one man sharpens another.",
         "Proverbs 27:17"),
        ("Create in me a clean heart, O God, and renew a right spirit within me.",
         "Psalm 51:10"),
        ("Walk by the Spirit, and you will not gratify the desires of the flesh.",
         "Galatians 5:16"),
        ("Put on the whole armor of God, that you may be able to stand against the schemes of the devil.",
         "Ephesians 6:11"),
        ("I can do all things through him who strengthens me.",
         "Philippians 4:13"),
        ("For everyone who has been born of God overcomes the world. And this is the victory that has overcome the world — our faith.",
         "1 John 5:4"),
        ("Do you not know that your body is a temple of the Holy Spirit within you, whom you have from God? You are not your own.",
         "1 Corinthians 6:19"),
        ("Set your minds on things that are above, not on things that are on earth.",
         "Colossians 3:2"),
        ("Finally, brothers, whatever is true, whatever is honorable, whatever is just, whatever is pure — think about these things.",
         "Philippians 4:8"),
        ("If we confess our sins, he is faithful and just to forgive us our sins and to cleanse us from all unrighteousness.",
         "1 John 1:9"),
        ("How can a young man keep his way pure? By guarding it according to your word.",
         "Psalm 119:9"),
        ("For the grace of God has appeared, bringing salvation for all people, training us to renounce ungodliness and worldly passions.",
         "Titus 2:11–12"),
        ("Confess your sins to one another and pray for one another, that you may be healed.",
         "James 5:16"),
        ("He who conceals his transgressions will not prosper, but he who confesses and forsakes them will obtain mercy.",
         "Proverbs 28:13"),
        ("Be sober-minded; be watchful. Your adversary the devil prowls around like a roaring lion, seeking someone to devour.",
         "1 Peter 5:8"),
        ("Therefore, since we are surrounded by so great a cloud of witnesses, let us lay aside every weight, and sin which clings so closely.",
         "Hebrews 12:1"),
        ("The Lord will fight for you; you need only to be still.",
         "Exodus 14:14"),
        ("For we do not have a high priest who is unable to sympathize with our weaknesses, but one who in every respect has been tempted as we are, yet without sin.",
         "Hebrews 4:15"),
        ("But I say, walk by the Spirit, and you will not gratify the desires of the flesh.",
         "Galatians 5:16"),
        ("Or do you not know that the unrighteous will not inherit the kingdom of God? Do not be deceived.",
         "1 Corinthians 6:9"),
        ("The name of the Lord is a strong tower; the righteous man runs into it and is safe.",
         "Proverbs 18:10"),
        ("So flee youthful passions and pursue righteousness, faith, love, and peace, along with those who call on the Lord from a pure heart.",
         "2 Timothy 2:22"),
        ("For the weapons of our warfare are not of the flesh but have divine power to destroy strongholds.",
         "2 Corinthians 10:4"),
        ("And we all, with unveiled face, beholding the glory of the Lord, are being transformed into the same image from one degree of glory to another.",
         "2 Corinthians 3:18"),
    ]

    private static var todaysVerse: (text: String, ref: String) {
        let dayOfYear = Calendar.current.ordinality(of: .day, in: .year, for: Date()) ?? 1
        return verses[(dayOfYear - 1) % verses.count]
    }

    var body: some View {
        let verse = Self.todaysVerse
        VStack(alignment: .leading, spacing: 10) {
            Text("VERSE OF THE DAY")
                .font(.system(size: 10, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.4)

            Text("\u{201C}\(verse.text)\u{201D}")
                .font(.system(size: 14, weight: .medium, design: .serif))
                .foregroundStyle(.white)
                .lineSpacing(4)
                .fixedSize(horizontal: false, vertical: true)

            Text("— \(verse.ref)")
                .font(.system(size: 12))
                .foregroundStyle(Color.white.opacity(0.45))
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.rfGold.opacity(0.07))
                .overlay(
                    RoundedRectangle(cornerRadius: 16)
                        .stroke(Color.rfGold.opacity(0.18), lineWidth: 1)
                )
        )
    }
}

// MARK: - Activity feed

private struct ActivitySection: View {
    var title: String = "Alerts"
    let events:    [ActivityEvent]
    var isLoading: Bool   = false
    var error:     String? = nil
    var emptyTitle: String = "No alerts yet."
    var onRetry:   (() -> Void)? = nil

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text(title)
                    .font(.headline)
                    .foregroundStyle(.white)
                Spacer()
                if isLoading {
                    ProgressView().tint(Color.rfGold).scaleEffect(0.75)
                }
            }

            if isLoading && events.isEmpty {
                loadingState
            } else if let err = error {
                errorState(err)
            } else if events.isEmpty {
                emptyState
            } else {
                VStack(spacing: 10) {
                    ForEach(events) { event in
                        NavigationLink(value: event) {
                            ActivityRow(event: event)
                        }
                        .buttonStyle(.plain)
                    }
                }
            }
        }
    }

    private var loadingState: some View {
        HStack {
            Spacer()
            ProgressView()
                .tint(Color.rfGold)
                .padding(.vertical, 36)
            Spacer()
        }
    }

    private func errorState(_ message: String) -> some View {
        Button(action: { onRetry?() }) {
            VStack(spacing: 12) {
                Image(systemName: "exclamationmark.triangle.fill")
                    .font(.system(size: 28))
                    .foregroundStyle(Color(red: 0.95, green: 0.62, blue: 0.10))
                Text(message)
                    .font(.system(size: 14))
                    .foregroundStyle(Color.white.opacity(0.60))
                    .multilineTextAlignment(.center)
                Text("Tap to retry")
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundStyle(Color.rfGold)
                    .padding(.horizontal, 16)
                    .padding(.vertical, 6)
                    .background(Capsule().fill(Color.rfGold.opacity(0.14)))
            }
            .frame(maxWidth: .infinity)
            .padding(.vertical, 28)
        }
        .buttonStyle(.plain)
    }

    private var emptyState: some View {
        VStack(spacing: 12) {
            Image(systemName: "checkmark.shield.fill")
                .font(.system(size: 34))
                .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45))
            Text(emptyTitle)
                .font(.subheadline)
                .foregroundStyle(Color.white.opacity(0.72))
                .multilineTextAlignment(.center)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 36)
    }
}

private struct ActivityRow: View {
    let event: ActivityEvent

    var body: some View {
        HStack(spacing: 14) {
            ZStack {
                Circle()
                    .fill(event.category.tint.opacity(0.14))
                    .frame(width: 44, height: 44)
                Image(systemName: event.category.icon)
                    .font(.system(size: 17))
                    .foregroundStyle(event.category.tint)
            }

            VStack(alignment: .leading, spacing: 3) {
                Text(event.category.rawValue)
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(.white)
                Text(event.description)
                    .font(.system(size: 12))
                    .foregroundStyle(Color.white.opacity(0.48))
            }

            Spacer(minLength: 8)

            VStack(alignment: .trailing, spacing: 5) {
                Text(event.severity.rawValue.uppercased())
                    .font(.caption2.weight(.bold))
                    .kerning(0.5)
                    .foregroundStyle(event.severity.color)
                    .accessibilityLabel("Severity \(event.severity.rawValue)")
                    .padding(.horizontal, 8)
                    .padding(.vertical, 3)
                    .background(Capsule().fill(event.severity.color.opacity(0.14)))

                Text(event.timeLabel)
                    .font(.system(size: 11))
                    .foregroundStyle(Color.white.opacity(0.32))
            }
        }
        .padding(14)
        .background(card)
    }
}

// MARK: - Optional Deep Scan prompt

private struct AppUsagePromptBanner: View {
    let onStart:   () -> Void
    let onDismiss: () -> Void

    @State private var showInstructions = false

    private let blue   = Color(red: 0.28, green: 0.56, blue: 0.95)
    private let bgFill = Color(red: 0.07, green: 0.13, blue: 0.28)

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(spacing: 12) {
                ZStack {
                    Circle()
                        .fill(blue.opacity(0.18))
                        .frame(width: 38, height: 38)
                    Image(systemName: "apps.iphone")
                        .font(.system(size: 16))
                        .foregroundStyle(blue)
                }
                VStack(alignment: .leading, spacing: 2) {
                    Text("An app you chose was used")
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(.white)
                    Text("Alerts and Blocking stay as you set them. Deep Scan is optional.")
                        .font(.footnote)
                        .foregroundStyle(Color.white.opacity(0.75))
                        .fixedSize(horizontal: false, vertical: true)
                }
                Spacer(minLength: 0)
                RFCloseButton(action: onDismiss)
            }

            RFButton(title: "How Deep Scan works", role: .secondary) {
                onStart()
                showInstructions = true
            }
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 12)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(bgFill)
                .overlay(
                    RoundedRectangle(cornerRadius: 14)
                        .stroke(blue.opacity(0.35), lineWidth: 1.5)
                )
        )
        .sheet(isPresented: $showInstructions) {
            BroadcastInstructionsSheet()
        }
    }
}

private struct BroadcastInstructionsSheet: View {
    @Environment(\.dismiss) private var dismiss

    private let steps: [(String, String)] = [
        ("1", "Swipe down from the top-right corner to open Control Center."),
        ("2", "Tap and hold the Screen Recording button (circle inside a circle)."),
        ("3", "Select \u{201C}Remain Faithful\u{201D} from the list of broadcast receivers."),
        ("4", "Tap “Start Broadcast”. Deep Scan starts. Alerts and Blocking do not change."),
    ]

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
                        Text("Deep Scan")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("Optional. It does not replace Alerts or Blocking.")
                            .font(.subheadline)
                            .foregroundStyle(Color.white.opacity(0.75))
                            .fixedSize(horizontal: false, vertical: true)
                    }
                    Spacer()
                    RFCloseButton { dismiss() }
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 24)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.bottom, 24)

                VStack(spacing: 16) {
                    ForEach(steps, id: \.0) { step in
                        HStack(alignment: .top, spacing: 16) {
                            ZStack {
                                Circle()
                                    .fill(Color.rfGold.opacity(0.14))
                                    .frame(width: 32, height: 32)
                                Text(step.0)
                                    .font(.system(size: 14, weight: .bold))
                                    .foregroundStyle(Color.rfGold)
                            }
                            Text(step.1)
                                .font(.system(size: 14))
                                .foregroundStyle(Color.white.opacity(0.80))
                                .lineSpacing(3)
                                .fixedSize(horizontal: false, vertical: true)
                            Spacer(minLength: 0)
                        }
                    }
                }
                .padding(.horizontal, 24)

                Text("Protected video such as Netflix is not visible. Only an alert category and a heartbeat leave this phone.")
                    .font(.footnote)
                    .foregroundStyle(Color.white.opacity(0.75))
                    .fixedSize(horizontal: false, vertical: true)
                    .padding(.horizontal, 24)
                    .padding(.top, 16)

                Spacer()

                RFButton(title: "Got it") { dismiss() }
                    .padding(.horizontal, 24)
                    .padding(.bottom, 36)
            }
        }
        .presentationDetents([.medium, .large])
    }
}

// MARK: - Partner Flag Banner

private struct PartnerFlagBanner: View {
    let alert: PartnerFlagAlert
    let onTap: () -> Void
    let onDismiss: () -> Void

    private let orange = Color(red: 0.95, green: 0.58, blue: 0.12)

    var body: some View {
        Button(action: onTap) {
            HStack(spacing: 12) {
                ZStack {
                    Circle()
                        .fill(orange.opacity(0.18))
                        .frame(width: 38, height: 38)
                    Image(systemName: "exclamationmark.triangle.fill")
                        .font(.system(size: 16))
                        .foregroundStyle(orange)
                }

                VStack(alignment: .leading, spacing: 2) {
                    Text("\(alert.senderName) flagged \(alert.displayCategory)")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundStyle(.white)
                        .lineLimit(1)
                    Text("Tap to review the alert")
                        .font(.system(size: 12))
                        .foregroundStyle(Color.white.opacity(0.50))
                }

                Spacer(minLength: 4)

                RFCloseButton(label: "Dismiss alert", action: onDismiss)
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 12)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color(red: 0.14, green: 0.18, blue: 0.34))
                    .overlay(
                        RoundedRectangle(cornerRadius: 14)
                            .stroke(orange.opacity(0.40), lineWidth: 1.5)
                    )
                    .shadow(color: .black.opacity(0.35), radius: 12, y: 4)
            )
        }
        .buttonStyle(.plain)
    }
}

// MARK: - Shared card background

private var card: some View { RFSurface() }

#Preview {
    @Previewable @State var showPanic = false
    NavigationStack {
        DashboardView(showPanic: $showPanic)
            .navigationDestination(for: ActivityEvent.self) { AlertDetailView(event: $0) }
    }
    .environmentObject(AppState.shared)
}
