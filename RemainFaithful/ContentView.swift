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

// MARK: - Settings tab (placeholder)

struct SettingsView: View {
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @AppStorage("userName")  private var userName  = ""
    @AppStorage("userEmail") private var userEmail = ""

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            VStack(spacing: 0) {
                // Profile block
                VStack(spacing: 12) {
                    ZStack {
                        Circle()
                            .fill(Color.rfGold.opacity(0.14))
                            .frame(width: 72, height: 72)
                        Text(String(userName.prefix(1)).uppercased())
                            .font(.system(size: 30, weight: .semibold))
                            .foregroundStyle(Color.rfGold)
                    }
                    Text(userName.isEmpty ? "Your Name" : userName)
                        .font(.system(size: 18, weight: .semibold))
                        .foregroundStyle(.white)
                    Text(userEmail.isEmpty ? "your@email.com" : userEmail)
                        .font(.system(size: 13))
                        .foregroundStyle(Color.white.opacity(0.45))
                }
                .padding(.top, 60)
                .padding(.bottom, 40)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.horizontal, 24)

                Spacer()

                // Sign out / reset
                Button {
                    userName  = ""
                    userEmail = ""
                    hasCompletedOnboarding = false
                } label: {
                    Text("Sign Out")
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(Color(red: 0.90, green: 0.30, blue: 0.30))
                        .frame(maxWidth: .infinity)
                        .frame(height: 50)
                        .background(
                            RoundedRectangle(cornerRadius: 14)
                                .fill(Color(red: 0.90, green: 0.30, blue: 0.30).opacity(0.10))
                                .overlay(
                                    RoundedRectangle(cornerRadius: 14)
                                        .stroke(Color(red: 0.90, green: 0.30, blue: 0.30).opacity(0.20), lineWidth: 1)
                                )
                        )
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 48)
            }
        }
        .toolbar(.hidden, for: .navigationBar)
    }
}

#Preview {
    ContentView()
}
