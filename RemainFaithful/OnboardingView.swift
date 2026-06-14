import SwiftUI

// MARK: - Color palette

extension Color {
    static let rfNavy    = Color(red: 0.07, green: 0.13, blue: 0.30)
    static let rfNavyMid = Color(red: 0.11, green: 0.19, blue: 0.40)
    static let rfGold    = Color(red: 0.82, green: 0.67, blue: 0.30)
}

// MARK: - Container

struct OnboardingView: View {
    var showDismissButton: Bool = false
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @State private var showLogin = false
    @Environment(\.dismiss) private var dismiss
    @AppStorage("userName")  private var storedName  = ""
    @AppStorage("userEmail") private var storedEmail = ""
    @State private var step     = 0
    @State private var name     = ""
    @State private var email    = ""
    @State private var password = ""

    // Total visible steps (step 0 = welcome, hidden from dots)
    private let totalSteps = 4

    var body: some View {
        ZStack {
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

                ZStack {
                    switch step {
                    case 0:
                        WelcomeStep(onContinue: advance, onSignIn: { showLogin = true })
                            .transition(.push(from: .trailing))
                    case 1:
                        CreateAccountStep(name: $name, email: $email, password: $password) {
                            storedName  = name
                            storedEmail = email
                            advance()
                        }
                        .transition(.push(from: .trailing))
                    case 2:
                        InvitePartnerStep(onContinue: advance, onSkip: advance)
                            .transition(.push(from: .trailing))
                    default:
                        MonitoringSetupStep {
                            hasCompletedOnboarding = true
                        }
                        .transition(.push(from: .trailing))
                    }
                }
                .animation(.easeInOut(duration: 0.35), value: step)
            }
        }
        .overlay(alignment: .topLeading) {
            if showDismissButton {
                Button { dismiss() } label: {
                    Image(systemName: "xmark.circle.fill")
                        .font(.system(size: 28))
                        .foregroundStyle(Color.white.opacity(0.55))
                        .padding(20)
                }
            }
        }
        .fullScreenCover(isPresented: $showLogin) {
            LoginView(showDismissButton: true)
        }
    }

    private var stepDots: some View {
        HStack(spacing: 7) {
            ForEach(1..<totalSteps) { i in
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
    var onSignIn: (() -> Void)? = nil

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
                Image(systemName: "checkmark.shield.fill")
                    .font(.system(size: 54, weight: .medium))
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

            Text("\"But you must remain faithful to the things you have been taught.\"  —  2 Timothy 3:14")
                .font(.system(size: 12, weight: .regular, design: .serif))
                .foregroundStyle(Color.white.opacity(0.38))
                .multilineTextAlignment(.center)
                .padding(.top, 24)

            Spacer()
            Spacer()

            RFButton(title: "Begin Your Journey", action: onContinue)
                .padding(.bottom, 16)

            Button { onSignIn?() } label: {
                HStack(spacing: 4) {
                    Text("Already have an account?")
                        .foregroundStyle(Color.white.opacity(0.50))
                    Text("Sign In")
                        .foregroundStyle(Color.rfGold)
                        .fontWeight(.semibold)
                }
                .font(.system(size: 15))
            }
            .padding(.bottom, 44)
        }
        .padding(.horizontal, 32)
    }

    private var line: some View {
        Rectangle()
            .fill(Color.rfGold.opacity(0.35))
            .frame(height: 1)
    }
}

// MARK: - Step 2: Create Account

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
                        .accessibilityIdentifier("name-field")
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
                        .accessibilityIdentifier("email-field")
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
                        .textContentType(.password)
                        .foregroundStyle(.white)
                        .focused($focused, equals: .password)
                        .submitLabel(.done)
                        .accessibilityIdentifier("password-field")
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
            do {
                _ = try await APIClient.shared.register(name: name, email: email, password: password)
                _ = try await APIClient.shared.login(email: email, password: password)
                await MainActor.run { onComplete() }
            } catch let e as APIError {
                await MainActor.run {
                    if case .server(let msg) = e {
                        errorMsg = msg
                    } else {
                        errorMsg = e.localizedDescription
                    }
                    isLoading = false
                }
            } catch {
                await MainActor.run {
                    errorMsg = "Something went wrong. Please try again."
                    isLoading = false
                }
            }
        }
    }
}

// MARK: - Step 3: Invite Partner

private struct InvitePartnerStep: View {
    let onContinue: () -> Void
    let onSkip:     () -> Void

    @FocusState private var emailFocused: Bool
    @State private var partnerEmail = ""
    @State private var isSending    = false
    @State private var didSend      = false
    @State private var errorMsg: String?

    private let green = Color(red: 0.20, green: 0.78, blue: 0.45)

    private var canSend: Bool {
        partnerEmail.contains("@") && partnerEmail.contains(".") && !isSending && !didSend
    }

    var body: some View {
        VStack(spacing: 0) {
            Spacer()

            // Icon
            ZStack {
                Circle().fill(Color.rfGold.opacity(0.08)).frame(width: 130, height: 130)
                Circle().fill(Color.rfGold.opacity(0.13)).frame(width: 90, height: 90)
                Image(systemName: "person.2.fill")
                    .font(.system(size: 40))
                    .foregroundStyle(Color.rfGold)
            }
            .padding(.bottom, 36)

            stepHeader(
                title: "Invite your\naccountability partner",
                subtitle: "Accountability works best with a trusted friend"
            )
            .padding(.bottom, 32)

            VStack(spacing: 14) {
                HStack(spacing: 12) {
                    Image(systemName: "envelope.fill")
                        .font(.system(size: 15))
                        .foregroundStyle(emailFocused ? Color.rfGold : Color.white.opacity(0.45))
                        .frame(width: 22)
                    TextField("", text: $partnerEmail,
                              prompt: Text("Partner's email address")
                                  .foregroundColor(Color.white.opacity(0.38)))
                        .keyboardType(.emailAddress)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .foregroundStyle(.white)
                        .focused($emailFocused)
                        .submitLabel(.send)
                        .onSubmit { sendInvite() }
                }
                .fieldStyle(isFocused: emailFocused)

                if let err = errorMsg {
                    Text(err)
                        .font(.system(size: 13))
                        .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                        .transition(.opacity)
                }

                RFButton(
                    title: didSend ? "Invitation Sent!" : (isSending ? "Sending…" : "Send Invitation"),
                    isEnabled: didSend || canSend
                ) {
                    didSend ? onContinue() : sendInvite()
                }
                .animation(.easeInOut(duration: 0.2), value: didSend)
            }

            Spacer()

            Button(action: onSkip) {
                Text("Skip for now")
                    .font(.system(size: 15))
                    .foregroundStyle(Color.white.opacity(0.40))
            }
            .padding(.bottom, 48)
        }
        .padding(.horizontal, 24)
        .animation(.easeInOut(duration: 0.2), value: errorMsg)
    }

    private func sendInvite() {
        guard canSend else { return }
        emailFocused = false
        isSending = true
        errorMsg  = nil
        let addr = partnerEmail.trimmingCharacters(in: .whitespaces).lowercased()
        Task {
            do {
                try await APIClient.shared.invitePartner(email: addr)
                await MainActor.run {
                    isSending = false
                    didSend   = true
                }
            } catch {
                await MainActor.run {
                    isSending = false
                    errorMsg  = "Couldn't send invite. Check the email and try again."
                }
            }
        }
    }
}

// MARK: - Step 4: Monitoring setup

private struct MonitoringSetupStep: View {
    let onComplete: () -> Void

    @ObservedObject private var fcManager = FamilyControlsManager.shared
    @State private var didRequestScreenTime = false

    var body: some View {
        VStack(spacing: 0) {
            Spacer()

            ZStack {
                Circle().fill(Color.rfGold.opacity(0.08)).frame(width: 130, height: 130)
                Circle().fill(Color.rfGold.opacity(0.13)).frame(width: 90, height: 90)
                Image(systemName: "shield.fill")
                    .font(.system(size: 40))
                    .foregroundStyle(Color.rfGold)
            }
            .padding(.bottom, 36)

            stepHeader(
                title: "Turn on\nmonitoring",
                subtitle: "Two layers work together to keep you accountable"
            )
            .padding(.bottom, 32)

            VStack(spacing: 12) {
                tierRow(
                    icon: "apps.iphone",
                    tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                    title: "App Monitoring",
                    detail: "Alerts your partner after 1 minute of use in any app you choose."
                )
                tierRow(
                    icon: "eye.fill",
                    tint: Color.rfGold,
                    title: "Deep Scan",
                    detail: "Visual monitoring during screen broadcasts detects sensitive content in real time."
                )
            }
            .padding(.horizontal, 4)
            .padding(.bottom, 28)

            if fcManager.authorizationStatus != .approved {
                RFButton(
                    title: didRequestScreenTime ? "Waiting for permission…" : "Enable App Monitoring",
                    isEnabled: !didRequestScreenTime
                ) {
                    didRequestScreenTime = true
                    Task { await fcManager.requestAuthorization() }
                }
                .padding(.bottom, 12)
            }

            RFButton(
                title: fcManager.authorizationStatus == .approved ? "Get Started" : "Continue Without Monitoring",
                isEnabled: true
            ) {
                // Pass onComplete as completion so the root-view swap only
                // happens after the permission dialog fully resolves. Calling
                // onComplete() synchronously before the dialog dismisses causes
                // iOS to drop the dialog (the presenting view hierarchy is gone).
                NotificationService.shared.requestPermission(then: onComplete)
            }

            Text("Deep Scan starts from Control Center → Screen Recording → Remain Faithful")
                .font(.system(size: 11))
                .foregroundStyle(Color.white.opacity(0.30))
                .multilineTextAlignment(.center)
                .lineSpacing(3)
                .padding(.top, 14)

            Spacer()
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 24)
    }

    private func tierRow(icon: String, tint: Color, title: String, detail: String) -> some View {
        HStack(alignment: .top, spacing: 16) {
            ZStack {
                RoundedRectangle(cornerRadius: 10)
                    .fill(tint.opacity(0.14))
                    .frame(width: 44, height: 44)
                Image(systemName: icon)
                    .font(.system(size: 18))
                    .foregroundStyle(tint)
            }
            VStack(alignment: .leading, spacing: 4) {
                Text(title)
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundStyle(.white)
                Text(detail)
                    .font(.system(size: 13))
                    .foregroundStyle(Color.white.opacity(0.50))
                    .fixedSize(horizontal: false, vertical: true)
            }
            Spacer()
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 14)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(Color.white.opacity(0.055))
                .overlay(RoundedRectangle(cornerRadius: 14)
                    .stroke(Color.white.opacity(0.08), lineWidth: 1))
        )
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
