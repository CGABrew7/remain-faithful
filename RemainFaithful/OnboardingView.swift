import SwiftUI

// MARK: - Color palette

extension Color {
    static let rfNavy    = Color(red: 0.07, green: 0.13, blue: 0.30)
    static let rfNavyMid = Color(red: 0.11, green: 0.19, blue: 0.40)
    static let rfGold    = Color(red: 0.82, green: 0.67, blue: 0.30)
}

// MARK: - Data

enum AccountabilityType {
    case oneToOne, smallGroup
}

// MARK: - Container

struct OnboardingView: View {
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @AppStorage("userName")  private var storedName  = ""
    @AppStorage("userEmail") private var storedEmail = ""
    @State private var step = 0
    @State private var selectedType: AccountabilityType?
    @State private var name     = ""
    @State private var email    = ""
    @State private var password = ""

    var body: some View {
        ZStack {
            // Background gradient — shared across all steps
            LinearGradient(
                colors: [Color.rfNavy, Color.rfNavyMid],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            VStack(spacing: 0) {
                stepDots
                    .padding(.top, 16)
                    .opacity(step == 0 ? 0 : 1)

                // Swap step views with a push transition
                ZStack {
                    if step == 0 {
                        WelcomeStep(onContinue: advance)
                            .transition(.push(from: .trailing))
                    } else if step == 1 {
                        AccountabilityTypeStep(
                            selected: $selectedType,
                            onContinue: advance
                        )
                        .transition(.push(from: .trailing))
                    } else {
                        CreateAccountStep(name: $name, email: $email, password: $password) {
                            storedName  = name
                            storedEmail = email
                            hasCompletedOnboarding = true
                            NotificationService.shared.requestPermission()
                        }
                        .transition(.push(from: .trailing))
                    }
                }
                .animation(.easeInOut(duration: 0.35), value: step)
            }
        }
    }

    private var stepDots: some View {
        HStack(spacing: 7) {
            ForEach(0..<3) { i in
                Capsule()
                    .fill(i == step ? Color.rfGold : Color.white.opacity(0.25))
                    .frame(width: i == step ? 20 : 7, height: 7)
                    .animation(.spring(response: 0.35, dampingFraction: 0.7), value: step)
            }
        }
    }

    private func advance() {
        withAnimation { step += 1 }
    }
}

// MARK: - Step 1: Welcome

private struct WelcomeStep: View {
    let onContinue: () -> Void

    var body: some View {
        VStack(spacing: 0) {
            Spacer()

            // Halo + cross icon
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.08))
                    .frame(width: 180, height: 180)
                Circle()
                    .fill(Color.rfGold.opacity(0.12))
                    .frame(width: 130, height: 130)
                Image(systemName: "cross.fill")
                    .font(.system(size: 52, weight: .medium))
                    .foregroundStyle(Color.rfGold)
            }
            .padding(.bottom, 44)

            Text("Remain Faithful")
                .font(.system(size: 36, weight: .bold, design: .serif))
                .foregroundStyle(.white)
                .multilineTextAlignment(.center)

            // Decorative rule
            HStack(spacing: 12) {
                line
                Image(systemName: "leaf.fill")
                    .font(.system(size: 10))
                    .foregroundStyle(Color.rfGold.opacity(0.7))
                line
            }
            .padding(.vertical, 20)
            .padding(.horizontal, 60)

            Text("Accountability between\ntrusted friends")
                .font(.system(size: 19, weight: .regular, design: .serif))
                .foregroundStyle(Color.white.opacity(0.75))
                .multilineTextAlignment(.center)
                .lineSpacing(5)

            Text("\"Iron sharpens iron.\"  —  Proverbs 27:17")
                .font(.system(size: 12, weight: .regular, design: .serif))
                .foregroundStyle(Color.white.opacity(0.38))
                .multilineTextAlignment(.center)
                .padding(.top, 24)

            Spacer()
            Spacer()

            RFButton(title: "Begin Your Journey", action: onContinue)
                .padding(.bottom, 52)
        }
        .padding(.horizontal, 32)
    }

    private var line: some View {
        Rectangle()
            .fill(Color.rfGold.opacity(0.35))
            .frame(height: 1)
    }
}

// MARK: - Step 2: Accountability Type

private struct AccountabilityTypeStep: View {
    @Binding var selected: AccountabilityType?
    let onContinue: () -> Void

    var body: some View {
        VStack(spacing: 0) {
            Spacer()

            stepHeader(
                title: "How will you\nwalk together?",
                subtitle: "Choose your accountability style"
            )
            .padding(.bottom, 36)

            VStack(spacing: 14) {
                TypeCard(
                    icon: "person.2.fill",
                    title: "One-to-One Partner",
                    description: "Walk closely with a single trusted friend who knows your journey.",
                    isSelected: selected == .oneToOne
                ) { selected = .oneToOne }

                TypeCard(
                    icon: "person.3.fill",
                    title: "Small Group",
                    description: "Share the journey with a small circle of trusted friends.",
                    isSelected: selected == .smallGroup
                ) { selected = .smallGroup }
            }

            Spacer()

            RFButton(
                title: "Continue",
                isEnabled: selected != nil,
                action: onContinue
            )
            .padding(.bottom, 52)
            .animation(.easeInOut(duration: 0.2), value: selected != nil)
        }
        .padding(.horizontal, 24)
    }
}

private struct TypeCard: View {
    let icon: String
    let title: String
    let description: String
    let isSelected: Bool
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            HStack(alignment: .center, spacing: 16) {
                ZStack {
                    Circle()
                        .fill(isSelected ? Color.rfGold : Color.white.opacity(0.10))
                        .frame(width: 52, height: 52)
                    Image(systemName: icon)
                        .font(.system(size: 21))
                        .foregroundStyle(isSelected ? Color.rfNavy : .white)
                }

                VStack(alignment: .leading, spacing: 4) {
                    Text(title)
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundStyle(.white)
                    Text(description)
                        .font(.system(size: 13))
                        .foregroundStyle(Color.white.opacity(0.55))
                        .fixedSize(horizontal: false, vertical: true)
                        .lineSpacing(2)
                }

                Spacer(minLength: 8)

                Image(systemName: isSelected ? "checkmark.circle.fill" : "circle")
                    .font(.system(size: 22))
                    .foregroundStyle(isSelected ? Color.rfGold : Color.white.opacity(0.25))
            }
            .padding(18)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Color.white.opacity(isSelected ? 0.11 : 0.06))
                    .overlay(
                        RoundedRectangle(cornerRadius: 16)
                            .stroke(
                                isSelected ? Color.rfGold : Color.white.opacity(0.10),
                                lineWidth: isSelected ? 1.5 : 1
                            )
                    )
            )
        }
        .buttonStyle(.plain)
        .animation(.easeInOut(duration: 0.18), value: isSelected)
    }
}

// MARK: - Step 3: Create Account

private struct CreateAccountStep: View {
    @Binding var name:     String
    @Binding var email:    String
    @Binding var password: String
    let onComplete: () -> Void

    @FocusState private var focused: Field?
    @State private var isLoading  = false
    @State private var errorMsg: String?

    private enum Field { case name, email, password }

    private var canSubmit: Bool {
        !name.trimmingCharacters(in: .whitespaces).isEmpty
            && email.contains("@") && email.contains(".")
            && password.count >= 8
            && !isLoading
    }

    var body: some View {
        VStack(spacing: 0) {
            Spacer()

            stepHeader(
                title: "Create your\naccount",
                subtitle: "Your journey begins here"
            )
            .padding(.bottom, 36)

            VStack(spacing: 14) {
                // Name
                HStack(spacing: 12) {
                    Image(systemName: "person.fill")
                        .font(.system(size: 15))
                        .foregroundStyle(focused == .name ? Color.rfGold : Color.white.opacity(0.45))
                        .frame(width: 22)
                    TextField("", text: $name,
                              prompt: Text("Your name").foregroundColor(Color.white.opacity(0.38)))
                        .textContentType(.name)
                        .textInputAutocapitalization(.words)
                        .foregroundStyle(.white)
                        .focused($focused, equals: .name)
                        .submitLabel(.next)
                        .onSubmit { focused = .email }
                }
                .fieldStyle(isFocused: focused == .name)

                // Email
                HStack(spacing: 12) {
                    Image(systemName: "envelope.fill")
                        .font(.system(size: 15))
                        .foregroundStyle(focused == .email ? Color.rfGold : Color.white.opacity(0.45))
                        .frame(width: 22)
                    TextField("", text: $email,
                              prompt: Text("Email address").foregroundColor(Color.white.opacity(0.38)))
                        .keyboardType(.emailAddress)
                        .textContentType(.emailAddress)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .foregroundStyle(.white)
                        .focused($focused, equals: .email)
                        .submitLabel(.next)
                        .onSubmit { focused = .password }
                }
                .fieldStyle(isFocused: focused == .email)

                // Password
                HStack(spacing: 12) {
                    Image(systemName: "lock.fill")
                        .font(.system(size: 15))
                        .foregroundStyle(focused == .password ? Color.rfGold : Color.white.opacity(0.45))
                        .frame(width: 22)
                    SecureField("", text: $password,
                                prompt: Text("Password (8+ characters)").foregroundColor(Color.white.opacity(0.38)))
                        .textContentType(.newPassword)
                        .foregroundStyle(.white)
                        .focused($focused, equals: .password)
                        .submitLabel(.done)
                        .onSubmit { if canSubmit { submit() } }
                }
                .fieldStyle(isFocused: focused == .password)
            }

            // Inline API error
            if let msg = errorMsg {
                Text(msg)
                    .font(.system(size: 13))
                    .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                    .multilineTextAlignment(.center)
                    .padding(.top, 12)
                    .transition(.opacity)
            }

            Text("By continuing you agree to our Terms of Service\nand Privacy Policy.")
                .font(.system(size: 11))
                .foregroundStyle(Color.white.opacity(0.30))
                .multilineTextAlignment(.center)
                .lineSpacing(3)
                .padding(.top, 16)

            Spacer()

            RFButton(
                title: isLoading ? "Creating account…" : "Get Started",
                isEnabled: canSubmit,
                action: submit
            )
            .padding(.bottom, 52)
            .animation(.easeInOut(duration: 0.2), value: canSubmit)
        }
        .padding(.horizontal, 24)
        .animation(.easeInOut(duration: 0.2), value: errorMsg)
    }

    private func submit() {
        guard canSubmit else { return }
        focused = nil
        isLoading = true
        errorMsg  = nil
        Task {
            // Try the backend; fall through locally if unavailable so the
            // simulator works without a running server.
            do {
                _ = try await APIClient.shared.register(name: name, email: email, password: password)
                _ = try await APIClient.shared.login(email: email, password: password)
            } catch let e as APIError {
                // Surface server-side errors (e.g. "email already taken") but
                // treat network failures as offline mode and proceed anyway.
                if case .server(let msg) = e {
                    await MainActor.run {
                        errorMsg  = msg
                        isLoading = false
                    }
                    return
                }
            } catch { }
            await MainActor.run { onComplete() }
        }
    }
}

// MARK: - Shared helpers

private func stepHeader(title: String, subtitle: String) -> some View {
    VStack(spacing: 10) {
        Text(title)
            .font(.system(size: 32, weight: .bold, design: .serif))
            .foregroundStyle(.white)
            .multilineTextAlignment(.center)
            .lineSpacing(2)
        Text(subtitle)
            .font(.system(size: 15))
            .foregroundStyle(Color.white.opacity(0.55))
    }
}

private struct RFButton: View {
    let title: String
    var isEnabled: Bool = true
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            Text(title)
                .font(.system(size: 17, weight: .semibold))
                .foregroundStyle(isEnabled ? Color.rfNavy : Color.white.opacity(0.35))
                .frame(maxWidth: .infinity)
                .frame(height: 56)
                .background(
                    RoundedRectangle(cornerRadius: 14)
                        .fill(isEnabled ? Color.rfGold : Color.white.opacity(0.08))
                )
        }
        .disabled(!isEnabled)
    }
}

private extension View {
    func fieldStyle(isFocused: Bool) -> some View {
        self
            .padding(.horizontal, 18)
            .frame(height: 56)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color.white.opacity(isFocused ? 0.11 : 0.07))
                    .overlay(
                        RoundedRectangle(cornerRadius: 14)
                            .stroke(
                                isFocused ? Color.rfGold : Color.white.opacity(0.10),
                                lineWidth: 1.5
                            )
                    )
            )
            .animation(.easeInOut(duration: 0.18), value: isFocused)
    }
}

#Preview {
    OnboardingView()
}
