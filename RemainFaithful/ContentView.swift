import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var appState: AppState
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
            .tabItem { Label("Home",     systemImage: "house.fill") }
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
        // Deep-link handling from notification taps
        .onReceive(appState.$deepLink.compactMap { $0 }) { link in
            handleDeepLink(link)
        }
    }

    private func handleDeepLink(_ link: DeepLink) {
        appState.deepLink = nil
        switch link {
        case .alertDetail(let event):
            selectedTab = 0
            // Small delay lets the tab switch settle before pushing.
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
