import SwiftUI

struct ContentView: View {
    var body: some View {
        NavigationStack {
            VStack(spacing: 32) {
                Spacer()

                Image(systemName: "cross.fill")
                    .font(.system(size: 44))
                    .foregroundStyle(Color.rfGold)

                Text("Remain Faithful")
                    .font(.system(size: 28, weight: .bold, design: .serif))

                Text("You're all set. More features coming soon.")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
                    .multilineTextAlignment(.center)

                Spacer()
            }
            .padding(.horizontal, 32)
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}

#Preview {
    ContentView()
}
