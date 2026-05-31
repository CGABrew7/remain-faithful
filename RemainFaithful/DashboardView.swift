import SwiftUI

// MARK: - Models

enum MonitoringStatus { case active, paused }

enum EventCategory: String {
    case adultContent = "Adult Content"
    case gambling     = "Gambling"
    case socialMedia  = "Social Media"
    case gaming       = "Gaming"

    var icon: String {
        switch self {
        case .adultContent: "eye.slash.fill"
        case .gambling:     "dollarsign.circle.fill"
        case .socialMedia:  "bubble.left.and.bubble.right.fill"
        case .gaming:       "gamecontroller.fill"
        }
    }

    // Maps the string stored in the API back to an enum case.
    // Accepts both the rawValue ("Adult Content") and snake_case ("adult_content").
    init?(apiString: String) {
        switch apiString.lowercased().replacingOccurrences(of: " ", with: "_") {
        case "adult_content": self = .adultContent
        case "gambling":      self = .gambling
        case "social_media":  self = .socialMedia
        case "gaming":        self = .gaming
        default:              return nil
        }
    }

    var tint: Color {
        switch self {
        case .adultContent: Color(red: 0.90, green: 0.25, blue: 0.30)
        case .gambling:     Color(red: 0.95, green: 0.58, blue: 0.12)
        case .socialMedia:  Color(red: 0.28, green: 0.56, blue: 0.95)
        case .gaming:       Color(red: 0.60, green: 0.30, blue: 0.90)
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
        }
    }
}

enum FlagSeverity: String {
    case low    = "Low"
    case medium = "Medium"
    case high   = "High"

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
    let category:    EventCategory
    let severity:    FlagSeverity
    let description: String
    let minutesAgo:  Int

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

    var extendedDescription: String {
        switch category {
        case .adultContent:
            if description.contains("search") {
                return "Search terms associated with adult content were detected during this session. The monitoring system logged the query for accountability review."
            }
            return "Adult content was detected during a web browsing session. The monitoring system flagged one or more pages containing explicit material."
        case .gambling:
            return "A gambling or sports-betting website was visited during this session. The monitoring system flagged the domain as a gambling-related destination."
        case .socialMedia:
            return "Extended time was spent on social media platforms during this session. Usage exceeded typical daily thresholds and was flagged for review."
        case .gaming:
            return "A gaming session was detected that ran late into the night. Extended play during these hours is flagged as a potential accountability concern."
        }
    }
}

// MARK: - API conversion

extension ActivityEvent {
    /// Converts a RemoteEvent from the API into a local ActivityEvent.
    /// Returns nil if the category or severity value is unrecognised.
    static func from(remote r: RemoteEvent) -> ActivityEvent? {
        // Category: try rawValue first ("Adult Content"), then snake_case ("adult_content")
        let cat = EventCategory(rawValue: r.category) ?? EventCategory(apiString: r.category)
        guard let category = cat else { return nil }

        // Severity: API stores lowercase ("low"), enum rawValues are title-cased ("Low")
        let severity = FlagSeverity(rawValue: r.severity.capitalized) ?? .low

        // Parse ISO-8601 timestamp; try with and without fractional seconds
        let fmt = ISO8601DateFormatter()
        fmt.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        let date = fmt.date(from: r.timestamp) ?? {
            let f2 = ISO8601DateFormatter()
            f2.formatOptions = [.withInternetDateTime]
            return f2.date(from: r.timestamp) ?? Date()
        }()
        let minutesAgo = max(0, Int(Date().timeIntervalSince(date) / 60))

        return ActivityEvent(category: category, severity: severity,
                             description: r.summary, minutesAgo: minutesAgo)
    }
}

// MARK: - Notification payload conversion

extension ActivityEvent {
    /// Reconstruct an ActivityEvent from an APNs notification userInfo dict so
    /// tapping a CONTENT_FLAGGED push can open AlertDetailView directly.
    static func from(notificationPayload p: [AnyHashable: Any]) -> ActivityEvent? {
        guard let catStr = p["category"] as? String else { return nil }
        let category = EventCategory(rawValue: catStr) ?? EventCategory(apiString: catStr)
        guard let category else { return nil }

        let severityStr = p["severity"] as? String ?? "low"
        let severity    = FlagSeverity(rawValue: severityStr.capitalized) ?? .low
        let description = p["summary"]  as? String ?? category.rawValue

        var minutesAgo = 0
        if let ts = p["timestamp"] as? String {
            let fmt = ISO8601DateFormatter()
            fmt.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
            let date = fmt.date(from: ts) ?? {
                let f2 = ISO8601DateFormatter()
                f2.formatOptions = [.withInternetDateTime]
                return f2.date(from: ts) ?? Date()
            }()
            minutesAgo = max(0, Int(Date().timeIntervalSince(date) / 60))
        }
        return ActivityEvent(category: category, severity: severity,
                             description: description, minutesAgo: minutesAgo)
    }
}

// MARK: - Placeholder data
// PLACEHOLDER: sampleEvents and weekClean are shown until real data loads from GET /events.
// streak data (days:14, best:21) is also hardcoded until a streak API endpoint is built.

private let sampleEvents: [ActivityEvent] = [
    .init(category: .adultContent, severity: .high,   description: "Flagged during browsing session",  minutesAgo: 23),
    .init(category: .socialMedia,  severity: .low,    description: "Extended social media usage",      minutesAgo: 145),
    .init(category: .gambling,     severity: .medium, description: "Sports betting site visited",      minutesAgo: 360),
    .init(category: .adultContent, severity: .medium, description: "Flagged search terms detected",    minutesAgo: 1440),
    .init(category: .gaming,       severity: .low,    description: "Late-night gaming session",        minutesAgo: 2160),
]

// true = clean day, false = flagged day  (Sun → Sat)
private let weekClean = [true, true, false, true, true, true, false]

// MARK: - Dashboard root

struct DashboardView: View {
    @AppStorage("userName")                  private var userName             = ""
    @AppStorage("hasDonated")                private var hasDonated           = false
    @AppStorage("donateBannerLastDismissed") private var donateBannerLastDismissed: Double = 0
    @EnvironmentObject private var appState: AppState
    @Binding var showPanic: Bool

    @State private var status:         MonitoringStatus = .active
    @State private var events:         [ActivityEvent]  = sampleEvents
    @State private var showDonation:   Bool             = false
    @State private var isLoadingEvents = false

    private var shouldShowDonateBanner: Bool {
        guard !hasDonated else { return false }
        guard donateBannerLastDismissed > 0 else { return true }
        let days = (Date().timeIntervalSince1970 - donateBannerLastDismissed) / 86_400
        return days >= 7
    }

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

            ScrollView(showsIndicators: false) {
                VStack(spacing: 18) {
                    headerRow
                    StatusCard(status: $status)
                    VerseCard()
                    StreakCard(days: 14, best: 21, week: weekClean)
                    ActivitySection(events: events, isLoading: isLoadingEvents)
                    panicButton
                    if shouldShowDonateBanner {
                        DonateBanner(
                            onDonate:  { showDonation = true;
                                         donateBannerLastDismissed = Date().timeIntervalSince1970 },
                            onDismiss: { donateBannerLastDismissed = Date().timeIntervalSince1970 }
                        )
                        .transition(.opacity.combined(with: .scale(scale: 0.97, anchor: .top)))
                    }
                }
                .padding(.horizontal, 20)
                .padding(.top, 8)
                .padding(.bottom, 32)
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
                // Auto-dismiss after 6 seconds; reset if a new alert replaces this one
                .task(id: alert.id) {
                    try? await Task.sleep(for: .seconds(6))
                    withAnimation(.spring(response: 0.35, dampingFraction: 0.75)) {
                        appState.clearPartnerAlert()
                    }
                }
            }
        }
        .toolbar(.hidden, for: .navigationBar)
        .sheet(isPresented: $showDonation) { DonationView() }
        .task { await loadEvents() }
    }

    @MainActor
    private func loadEvents() async {
        guard APIClient.shared.isAuthenticated else { return }
        isLoadingEvents = true
        defer { isLoadingEvents = false }
        do {
            let remote = try await APIClient.shared.listEvents()
            let converted = remote.compactMap { ActivityEvent.from(remote: $0) }
            if !converted.isEmpty { events = converted }
            // If empty, placeholder sampleEvents remain visible so the UI isn't blank.
        } catch {
            // Keep placeholder data visible if the server is unreachable.
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

private struct StatusCard: View {
    @Binding var status: MonitoringStatus
    @State private var pulse: CGFloat = 1.0

    private var isActive: Bool { status == .active }

    private static let green   = Color(red: 0.18, green: 0.82, blue: 0.48)
    private static let bgOn    = LinearGradient(colors: [Color(red: 0.07, green: 0.40, blue: 0.26), Color(red: 0.05, green: 0.26, blue: 0.18)], startPoint: .topLeading, endPoint: .bottomTrailing)
    private static let bgOff   = LinearGradient(colors: [Color(red: 0.16, green: 0.20, blue: 0.32), Color(red: 0.11, green: 0.15, blue: 0.26)], startPoint: .topLeading, endPoint: .bottomTrailing)

    var body: some View {
        VStack(spacing: 0) {
            // Top row
            HStack(alignment: .top, spacing: 14) {
                // Pulse ring + dot
                ZStack {
                    if isActive {
                        Circle()
                            .fill(Self.green.opacity(0.22))
                            .frame(width: 44, height: 44)
                            .scaleEffect(pulse)
                    }
                    Circle()
                        .fill(isActive ? Self.green : Color(.systemGray3))
                        .frame(width: 11, height: 11)
                }
                .frame(width: 44, height: 44)

                VStack(alignment: .leading, spacing: 5) {
                    Text(isActive ? "Monitoring Active" : "Monitoring Paused")
                        .font(.system(size: 18, weight: .bold))
                        .foregroundStyle(.white)
                    Text(isActive
                         ? "Your screen activity is being monitored"
                         : "Activity recording is paused")
                        .font(.system(size: 13))
                        .foregroundStyle(Color.white.opacity(0.6))
                        .fixedSize(horizontal: false, vertical: true)
                    if isActive {
                        Text("Pausing will notify your accountability group")
                            .font(.system(size: 12, weight: .medium))
                            .foregroundStyle(Color(red: 0.95, green: 0.72, blue: 0.22))
                            .fixedSize(horizontal: false, vertical: true)
                    }
                }

                Spacer(minLength: 0)
            }
            .padding(.bottom, 18)

            Divider().overlay(Color.white.opacity(0.10))
                .padding(.bottom, 14)

            // Bottom row
            HStack {
                HStack(spacing: 7) {
                    Image(systemName: isActive ? "shield.fill" : "shield.slash.fill")
                        .font(.system(size: 12))
                        .foregroundStyle(isActive ? Self.green : Color(.systemGray3))
                    Text(isActive ? "Active session" : "Session ended")
                        .font(.system(size: 12))
                        .foregroundStyle(Color.white.opacity(0.45))
                }

                Spacer()

                Button {
                    withAnimation(.spring(response: 0.3, dampingFraction: 0.75)) {
                        status = isActive ? .paused : .active
                    }
                } label: {
                    Text(isActive ? "Pause" : "Resume")
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundStyle(isActive ? Self.green : Color.rfGold)
                        .padding(.horizontal, 16)
                        .padding(.vertical, 7)
                        .background(Capsule().fill(
                            isActive ? Self.green.opacity(0.15) : Color.rfGold.opacity(0.15)
                        ))
                }
            }
        }
        .padding(20)
        .background(
            RoundedRectangle(cornerRadius: 20)
                .fill(isActive ? Self.bgOn : Self.bgOff)
                .overlay(
                    RoundedRectangle(cornerRadius: 20)
                        .stroke(
                            isActive ? Self.green.opacity(0.25) : Color.white.opacity(0.07),
                            lineWidth: 1
                        )
                )
        )
        .animation(.easeInOut(duration: 0.28), value: isActive)
        .onAppear {
            withAnimation(.easeInOut(duration: 1.1).repeatForever(autoreverses: true)) {
                pulse = 1.55
            }
        }
    }
}

// MARK: - Streak card

private struct StreakCard: View {
    let days: Int
    let best: Int
    let week: [Bool]

    private let dayLabels = ["S", "M", "T", "W", "T", "F", "S"]
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
                ForEach(0..<7, id: \.self) { i in
                    VStack(spacing: 6) {
                        ZStack {
                            Circle()
                                .fill(week[i]
                                      ? Color.rfGold.opacity(0.18)
                                      : flagRed.opacity(0.18))
                                .frame(width: 28, height: 28)
                            Image(systemName: week[i] ? "checkmark" : "xmark")
                                .font(.system(size: 10, weight: .bold))
                                .foregroundStyle(week[i] ? Color.rfGold : flagRed)
                        }
                        Text(dayLabels[i])
                            .font(.system(size: 10, weight: .medium))
                            .foregroundStyle(Color.white.opacity(0.32))
                    }
                    .frame(maxWidth: .infinity)
                }
            }
        }
        .padding(20)
        .background(card)
    }
}

// MARK: - Verse card

private struct VerseCard: View {
    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            Text("VERSE OF THE DAY")
                .font(.system(size: 10, weight: .bold))
                .foregroundStyle(Color.rfGold.opacity(0.75))
                .kerning(1.4)

            Text("\"Watch, stand fast in the faith, be brave, be strong.\"")
                .font(.system(size: 14, weight: .medium, design: .serif))
                .foregroundStyle(.white)
                .lineSpacing(4)
                .fixedSize(horizontal: false, vertical: true)

            Text("— 1 Corinthians 16:13")
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
    let events: [ActivityEvent]
    var isLoading: Bool = false

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Recent Flags")
                    .font(.system(size: 16, weight: .semibold))
                    .foregroundStyle(.white)
                Spacer()
                if isLoading {
                    ProgressView().tint(Color.rfGold).scaleEffect(0.75)
                } else {
                    Text("See all")
                        .font(.system(size: 13, weight: .medium))
                        .foregroundStyle(Color.rfGold)
                }
            }

            if events.isEmpty && !isLoading {
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

    private var emptyState: some View {
        VStack(spacing: 12) {
            Image(systemName: "checkmark.shield.fill")
                .font(.system(size: 34))
                .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45))
            Text("No flags in the last 7 days")
                .font(.system(size: 15))
                .foregroundStyle(Color.white.opacity(0.45))
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
                    .font(.system(size: 10, weight: .bold))
                    .kerning(0.5)
                    .foregroundStyle(event.severity.color)
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

// MARK: - Donate Banner

private struct DonateBanner: View {
    let onDonate:  () -> Void
    let onDismiss: () -> Void

    var body: some View {
        HStack(spacing: 12) {
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.14))
                    .frame(width: 36, height: 36)
                Image(systemName: "heart.fill")
                    .font(.system(size: 14))
                    .foregroundStyle(Color.rfGold)
            }

            VStack(alignment: .leading, spacing: 2) {
                Text("Keep RF Free — Donate")
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundStyle(.white)
                Text("This app runs on generous support.")
                    .font(.system(size: 11))
                    .foregroundStyle(Color.white.opacity(0.45))
            }

            Spacer(minLength: 4)

            Button(action: onDonate) {
                Text("Give")
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundStyle(Color(red: 0.05, green: 0.09, blue: 0.22))
                    .padding(.horizontal, 12)
                    .padding(.vertical, 6)
                    .background(Capsule().fill(Color.rfGold))
            }

            Button(action: onDismiss) {
                Image(systemName: "xmark")
                    .font(.system(size: 10, weight: .bold))
                    .foregroundStyle(Color.white.opacity(0.35))
                    .padding(6)
                    .background(Circle().fill(Color.white.opacity(0.07)))
            }
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 10)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(Color.white.opacity(0.04))
                .overlay(RoundedRectangle(cornerRadius: 12)
                    .stroke(Color.rfGold.opacity(0.22), lineWidth: 1))
        )
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

                Button(action: onDismiss) {
                    Image(systemName: "xmark")
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(Color.white.opacity(0.45))
                        .padding(8)
                        .background(Circle().fill(Color.white.opacity(0.08)))
                }
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

private var card: some View {
    RoundedRectangle(cornerRadius: 16)
        .fill(Color.white.opacity(0.055))
        .overlay(
            RoundedRectangle(cornerRadius: 16)
                .stroke(Color.white.opacity(0.08), lineWidth: 1)
        )
}

#Preview {
    @Previewable @State var showPanic = false
    NavigationStack {
        DashboardView(showPanic: $showPanic)
            .navigationDestination(for: ActivityEvent.self) { AlertDetailView(event: $0) }
    }
    .environmentObject(AppState.shared)
}
