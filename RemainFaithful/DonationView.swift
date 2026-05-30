import SwiftUI

struct DonationView: View {
    @Environment(\.dismiss) private var dismiss

    @State private var isMonthly     = false
    @State private var selectedPreset: Int? = 10
    @State private var customAmount  = ""
    @State private var isLoading     = false
    @State private var checkoutURL:  URL?   = nil
    @State private var showSafari    = false
    @State private var errorMessage: String? = nil
    @State private var whyExpanded   = false

    private let presets = [5, 10, 25, 50]
    private let pink    = Color(red: 0.90, green: 0.25, blue: 0.48)

    private var effectiveAmount: Int? {
        if let p = selectedPreset { return p }
        let v = Int(customAmount.trimmingCharacters(in: .whitespaces))
        return v.flatMap { $0 >= 1 ? $0 : nil }
    }

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()

            ScrollView(showsIndicators: false) {
                VStack(spacing: 24) {
                    closeButton
                    headerSection
                    frequencyToggle
                    amountGrid
                    if selectedPreset == nil {
                        customAmountField
                    }
                    donateButton
                    whyDonateSection
                    Spacer(minLength: 40)
                }
                .padding(.horizontal, 20)
                .padding(.top, 8)
            }
        }
        .sheet(isPresented: $showSafari) {
            if let url = checkoutURL {
                SafariView(url: url).ignoresSafeArea()
            }
        }
    }

    // MARK: - Close

    private var closeButton: some View {
        HStack {
            Spacer()
            Button { dismiss() } label: {
                Image(systemName: "xmark")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(Color.white.opacity(0.55))
                    .padding(10)
                    .background(Circle().fill(Color.white.opacity(0.09)))
            }
        }
    }

    // MARK: - Header

    private var headerSection: some View {
        VStack(spacing: 14) {
            ZStack {
                Circle()
                    .fill(pink.opacity(0.14))
                    .frame(width: 72, height: 72)
                Image(systemName: "heart.fill")
                    .font(.system(size: 28))
                    .foregroundStyle(pink)
            }

            Text("Support Remain Faithful")
                .font(.system(size: 24, weight: .bold, design: .serif))
                .foregroundStyle(.white)
                .multilineTextAlignment(.center)

            Text("Remain Faithful is completely free — no subscriptions, no paywalls. It exists because people like you believe in this mission.")
                .font(.system(size: 14))
                .foregroundStyle(Color.white.opacity(0.60))
                .multilineTextAlignment(.center)
                .lineSpacing(4)
        }
    }

    // MARK: - Monthly toggle

    private var frequencyToggle: some View {
        HStack(spacing: 0) {
            FreqTab(label: "One-time", selected: !isMonthly) {
                withAnimation(.spring(response: 0.25, dampingFraction: 0.75)) { isMonthly = false }
            }
            FreqTab(label: "Monthly",  selected: isMonthly) {
                withAnimation(.spring(response: 0.25, dampingFraction: 0.75)) { isMonthly = true  }
            }
        }
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(Color.white.opacity(0.06))
                .overlay(RoundedRectangle(cornerRadius: 12)
                    .stroke(Color.white.opacity(0.09), lineWidth: 1))
        )
        .clipShape(RoundedRectangle(cornerRadius: 12))
    }

    // MARK: - Amount grid

    private var amountGrid: some View {
        LazyVGrid(
            columns: Array(repeating: GridItem(.flexible(), spacing: 12), count: 4),
            spacing: 12
        ) {
            ForEach(presets, id: \.self) { amount in
                AmountBtn(label: "$\(amount)", selected: selectedPreset == amount) {
                    withAnimation(.spring(response: 0.22, dampingFraction: 0.7)) {
                        selectedPreset = amount
                        customAmount   = ""
                    }
                }
            }
            AmountBtn(label: "Custom", selected: selectedPreset == nil) {
                withAnimation(.spring(response: 0.22, dampingFraction: 0.7)) {
                    selectedPreset = nil
                }
            }
        }
    }

    // MARK: - Custom amount

    private var customAmountField: some View {
        HStack(spacing: 8) {
            Text("$")
                .font(.system(size: 24, weight: .semibold))
                .foregroundStyle(Color.rfGold)
            TextField("Enter amount", text: $customAmount)
                .keyboardType(.numberPad)
                .font(.system(size: 24, weight: .semibold))
                .foregroundStyle(.white)
        }
        .padding(16)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(Color.white.opacity(0.07))
                .overlay(RoundedRectangle(cornerRadius: 12)
                    .stroke(Color.rfGold.opacity(0.40), lineWidth: 1.5))
        )
    }

    // MARK: - Donate button

    private var donateButton: some View {
        VStack(spacing: 12) {
            let amount = effectiveAmount ?? 0
            Button {
                guard amount >= 1 else { return }
                Task { await startCheckout(amount: amount) }
            } label: {
                Group {
                    if isLoading {
                        ProgressView().tint(.white)
                    } else {
                        let suffix = isMonthly ? "/mo" : ""
                        Text(amount > 0 ? "Donate $\(amount)\(suffix)" : "Select an amount")
                            .font(.system(size: 17, weight: .bold))
                            .foregroundStyle(amount > 0 ? .white : Color.white.opacity(0.40))
                    }
                }
                .frame(maxWidth: .infinity)
                .frame(height: 56)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(
                            LinearGradient(
                                colors: amount > 0
                                    ? [pink, Color(red: 0.65, green: 0.15, blue: 0.35)]
                                    : [Color.white.opacity(0.08), Color.white.opacity(0.06)],
                                startPoint: .topLeading,
                                endPoint: .bottomTrailing
                            )
                        )
                )
            }
            .disabled(isLoading || amount < 1)

            if let msg = errorMessage {
                Text(msg)
                    .font(.system(size: 13))
                    .foregroundStyle(Color(red: 0.90, green: 0.30, blue: 0.30))
                    .multilineTextAlignment(.center)
            }

            Text("Powered by Stripe · Secure checkout")
                .font(.system(size: 12))
                .foregroundStyle(Color.white.opacity(0.30))
        }
    }

    // MARK: - Why donate?

    private var whyDonateSection: some View {
        VStack(spacing: 0) {
            Button {
                withAnimation(.spring(response: 0.3, dampingFraction: 0.8)) {
                    whyExpanded.toggle()
                }
            } label: {
                HStack {
                    Text("Why donate?")
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(.white)
                    Spacer()
                    Image(systemName: whyExpanded ? "chevron.up" : "chevron.down")
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(Color.white.opacity(0.40))
                }
                .padding(16)
            }

            if whyExpanded {
                VStack(alignment: .leading, spacing: 14) {
                    WhyRow(icon: "server.rack",
                           text: "Server costs — AI content analysis and push notifications aren't free. Your gift keeps them running.")
                    WhyRow(icon: "lock.shield.fill",
                           text: "Privacy infrastructure — keeping your data off ad networks and behind proper encryption costs money.")
                    WhyRow(icon: "heart.fill",
                           text: "Ministry mission — we believe every man deserves this tool, regardless of income.")
                    WhyRow(icon: "cross.fill",
                           text: "100% mission-driven — no VC money, no ads, no data selling. Just people helping people stay free.")
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 18)
            }
        }
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.white.opacity(0.055))
                .overlay(RoundedRectangle(cornerRadius: 16)
                    .stroke(Color.white.opacity(0.08), lineWidth: 1))
        )
        .clipShape(RoundedRectangle(cornerRadius: 16))
    }

    // MARK: - Checkout

    @MainActor
    private func startCheckout(amount: Int) async {
        isLoading    = true
        errorMessage = nil
        defer { isLoading = false }
        do {
            let url = try await APIClient.shared.createCheckoutSession(
                amountDollars: amount, monthly: isMonthly)
            checkoutURL = url
            showSafari  = true
        } catch {
            errorMessage = "Couldn't start checkout. Please try again."
        }
    }
}

// MARK: - Sub-views

private struct FreqTab: View {
    let label:    String
    let selected: Bool
    let action:   () -> Void

    var body: some View {
        Button(action: action) {
            Text(label)
                .font(.system(size: 14, weight: .semibold))
                .foregroundStyle(selected ? Color(red: 0.05, green: 0.09, blue: 0.22) : Color.white.opacity(0.55))
                .frame(maxWidth: .infinity)
                .frame(height: 44)
                .background(selected ? Color.rfGold : Color.clear)
                .clipShape(RoundedRectangle(cornerRadius: 11))
                .padding(2)
        }
    }
}

private struct AmountBtn: View {
    let label:    String
    let selected: Bool
    let action:   () -> Void

    var body: some View {
        Button(action: action) {
            Text(label)
                .font(.system(size: 15, weight: .semibold))
                .foregroundStyle(selected ? Color(red: 0.05, green: 0.09, blue: 0.22) : .white)
                .frame(maxWidth: .infinity)
                .frame(height: 48)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(selected ? Color.rfGold : Color.white.opacity(0.07))
                        .overlay(RoundedRectangle(cornerRadius: 12)
                            .stroke(selected ? Color.clear : Color.white.opacity(0.10), lineWidth: 1))
                )
        }
    }
}

private struct WhyRow: View {
    let icon: String
    let text: String

    var body: some View {
        HStack(alignment: .top, spacing: 10) {
            Image(systemName: icon)
                .font(.system(size: 13))
                .foregroundStyle(Color.rfGold.opacity(0.80))
                .frame(width: 18, alignment: .center)
                .padding(.top, 1)
            Text(text)
                .font(.system(size: 13))
                .foregroundStyle(Color.white.opacity(0.65))
                .lineSpacing(3)
                .fixedSize(horizontal: false, vertical: true)
        }
    }
}

#Preview {
    DonationView()
}
