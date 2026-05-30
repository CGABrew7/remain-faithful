import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var appState: AppState
    @Environment(\.scenePhase) private var scenePhase
    @State private var selectedTab  = 0
    @State private var homePath     = NavigationPath()
    @State private var showPanic    = false

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
        .fullScreenCover(isPresented: $showPanic) { PanicView() }
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
            guard phase == .active, APIClient.shared.isAuthenticated else { return }
            Task {
                let count = (try? await APIClient.shared.alertUnreadCount()) ?? 0
                appState.unreadAlertCount = count
                // Auto-clear badge if user is already on Home tab
                if selectedTab == 0 && count > 0 { markSeenAndClearBadge() }
            }
        }
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
