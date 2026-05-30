import SwiftUI

struct ContentView: View {
    var body: some View {
        TabView {
            NavigationStack { DashboardView() }
                .tabItem { Label("Home",     systemImage: "house.fill") }

            NavigationStack { GroupView() }
                .tabItem { Label("Group",    systemImage: "person.2.fill") }

            NavigationStack { SettingsView() }
                .tabItem { Label("Settings", systemImage: "gearshape.fill") }
        }
        .tint(Color.rfGold)
        .toolbarBackground(Color(red: 0.05, green: 0.09, blue: 0.22), for: .tabBar)
        .toolbarBackground(.visible, for: .tabBar)
    }
}


#Preview {
    ContentView()
}
