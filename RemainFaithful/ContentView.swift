import SwiftUI
import UserNotifications

struct ContentView: View {
    @EnvironmentObject private var appState: AppState
    @Environment(\.scenePhase) private var scenePhase
    @AppStorage("tourCompletedVersion") private var tourCompletedVersion = 0
    @AppStorage("primaryGroupID") private var primaryGroupID = 0
    @State private var selectedTab  = 0
    @State private var homePath     = NavigationPath()
    @State private var showPanic    = false
    @State private var tourVisible  = false
    @State private var tourStep     = 0
    @State private var tourSteps: [TourStep] = []

    var body: some View {
        TabView(selection: $selectedTab) {
            NavigationStack(path: $homePath) {
                DashboardView(showPanic: $showPanic)
                    .navigationDestination(for: ActivityEvent.self) { event in
                        AlertDetailView(event: event)
                    }
            }
            .tabItem { Label("Home", systemImage: "house.fill") }
            .badge(appState.unreadAlertCount)
            .tag(0)

            NavigationStack { GroupView() }
                .tabItem { Label("Group",    systemImage: "person.2.fill") }
                .tag(1)

            NavigationStack { SettingsView() }
                .tabItem { Label("Settings", systemImage: "gearshape.fill") }
                .tag(2)
        }
        .tint(Color.rfGold)
        .toolbarBackground(Color(red: 0.05, green: 0.09, blue: 0.22), for: .tabBar)
        .toolbarBackground(.visible, for: .tabBar)
        .accessibilityHidden(tourVisible)
        .fullScreenCover(isPresented: $showPanic) { PanicView() }
        .overlayPreferenceValue(TourAnchorKey.self) { anchors in
            if tourVisible, tourSteps.indices.contains(tourStep) {
                GeometryReader { proxy in
                    TourOverlay(
                        step: tourSteps[tourStep],
                        index: tourStep,
                        count: tourSteps.count,
                        hole: anchors[tourSteps[tourStep].anchor],
                        proxy: proxy,
                        onBack: { if tourStep > 0 { tourStep -= 1 } },
                        onNext: {
                            if tourStep + 1 >= tourSteps.count { finishTour() }
                            else { tourStep += 1 }
                        },
                        onSkip: { finishTour() }
                    )
                }
                .ignoresSafeArea()
            }
        }
        #if targetEnvironment(simulator)
        .safeAreaInset(edge: .top, spacing: 0) {
            if appState.isDemoMode {
                HStack(spacing: 5) {
                    Image(systemName: "play.rectangle.fill")
                        .font(.system(size: 10, weight: .bold))
                    Text("DEMO MODE — sample data only")
                        .font(.system(size: 11, weight: .bold))
                }
                .foregroundStyle(Color(red: 0.10, green: 0.10, blue: 0.10))
                .frame(maxWidth: .infinity)
                .padding(.vertical, 6)
                .background(Color(red: 0.98, green: 0.82, blue: 0.25))
            }
        }
        #endif
        // Deep-link handling from notification taps
        .onReceive(appState.$deepLink.compactMap { $0 }) { link in
            handleDeepLink(link)
        }
        // Mark alerts seen and clear badge when user lands on Home
        .onChange(of: selectedTab) { _, newTab in
            if newTab == 0 { markSeenAndClearBadge() }
        }
        // Refresh unread count each time the app comes to the foreground
        .onChange(of: scenePhase) { _, phase in
            if phase == .active {
                Task { await presentTourIfNeeded(replay: false) }
            }
            guard phase == .active, APIClient.shared.isAuthenticated else { return }
            Task {
                let count = (try? await APIClient.shared.alertUnreadCount()) ?? 0
                appState.unreadAlertCount = count
                // Auto-clear badge if user is already on Home tab
                if selectedTab == 0 && count > 0 { markSeenAndClearBadge() }
            }
        }
        .onAppear {
            Task { await presentTourIfNeeded(replay: false) }
        }
        .onChange(of: tourStep) { _, _ in focusTourStep() }
        .onChange(of: showPanic) { _, isOpen in
            if isOpen { tourVisible = false }
        }
        .onReceive(NotificationCenter.default.publisher(for: .replayAppTour)) { _ in
            tourCompletedVersion = 0
            tourVisible = false
            Task { await presentTourIfNeeded(replay: true) }
        }
        .onReceive(NotificationCenter.default.publisher(for: .notificationPermissionResolved)) { _ in
            Task { await presentTourIfNeeded(replay: false) }
        }
    }

    private func finishTour() {
        tourCompletedVersion = TourContent.version
        tourVisible = false
        TourController.shared.activeAnchor = nil
    }

    private func focusTourStep() {
        guard tourVisible, tourSteps.indices.contains(tourStep) else { return }
        let step = tourSteps[tourStep]
        selectedTab = step.tab
        TourController.shared.activeAnchor = step.anchor
    }

    @MainActor
    private func presentTourIfNeeded(replay: Bool) async {
        let force = CommandLine.arguments.contains("-forceTour")
        let wants = replay || TourContent.shouldAutoPresent(
            completedVersion: tourCompletedVersion,
            isDemoMode: appState.isDemoMode,
            forceTour: force
        )
        guard wants else { return }
        guard scenePhase == .active || replay else { return }
        guard !showPanic else { return }
        if tourVisible && !replay { return }

        // Stay out of the way of the notification permission prompt.
        if !force {
            let status = await UNUserNotificationCenter.current().notificationSettings().authorizationStatus
            if status == .notDetermined { return }
            if tourVisible && !replay { return }
        }

        if !appState.isDemoMode, APIClient.shared.isAuthenticated {
            if let groups = try? await APIClient.shared.listMyGroups() {
                let resolved = GroupSelection.resolvePrimaryGroupID(
                    stored: primaryGroupID,
                    memberGroupIDs: groups.map(\.id)
                )
                if resolved != primaryGroupID { primaryGroupID = resolved }
            }
            if tourVisible && !replay { return }
        }

        let inGroup = appState.isDemoMode || primaryGroupID > 0
        tourSteps = TourContent.steps(inGroup: inGroup)
        tourStep = 0
        tourVisible = true
        focusTourStep()
    }

    private func markSeenAndClearBadge() {
        appState.resetUnreadCount()
        Task { try? await APIClient.shared.markAlertsSeen() }
    }

    private func handleDeepLink(_ link: DeepLink) {
        appState.deepLink = nil
        switch link {
        case .alertDetail(let event):
            selectedTab = 0
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
                homePath.append(event)
            }
        case .panicView:
            selectedTab = 0
            showPanic   = true
        case .group:
            selectedTab = 1
        case .dashboard:
            selectedTab = 0
        }
    }
}

#Preview {
    ContentView()
        .environmentObject(AppState.shared)
}
