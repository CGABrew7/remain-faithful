import SwiftUI
import FamilyControls

// The navy and gold palette now lives in RFTheme.swift alongside the shared
// surface, button and toggle styles. Values are unchanged.

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

    // Welcome is step 0 and is hidden from the dots.
    // 1 account, 2 invite, 3 Screen Time, 4 Deep Scan, 5 apps.
    private let totalSteps = 6

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
                    case 3:
                        MonitoringSetupStep(onContinue: advance, onSkip: advance)
                            .transition(.push(from: .trailing))
                    case 4:
                        DeepScanStep(onContinue: advance)
                            .transition(.push(from: .trailing))
                    default:
                        AppPickerStep {
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
                RFCloseButton { dismiss() }
                    .padding(8)
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
            .accessibilityHidden(true)
            .padding(.bottom, 44)

            Text("Remain Faithful")
                .font(.largeTitle.bold())
                .fontDesign(.serif)
                .foregroundStyle(.white)
                .multilineTextAlignment(.center)
                .accessibilityAddTraits(.isHeader)

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
                .font(.footnote)
                .fontDesign(.serif)
                .foregroundStyle(Color.white.opacity(0.72))
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
    @State private var isAdult = false
    @State private var showPassword = false

    private enum Field { case name, email, password }

    private var canSubmit: Bool {
        !name.trimmingCharacters(in: .whitespaces).isEmpty
            && email.contains("@") && email.contains(".")
            && PasswordRules.isMet(password)
            && isAdult
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
                        .foregroundStyle(focused == .password ? Color.rfGold : Color.white.opacity(0.55))
                        .frame(width: 22)
                    Group {
                        if showPassword {
                            TextField("", text: $password,
                                      prompt: Text("Password").foregroundColor(Color.white.opacity(0.55)))
                        } else {
                            SecureField("", text: $password,
                                        prompt: Text("Password").foregroundColor(Color.white.opacity(0.55)))
                        }
                    }
                    .textContentType(.newPassword)
                    .foregroundStyle(.white)
                    .focused($focused, equals: .password)
                    .submitLabel(.done)
                    .accessibilityIdentifier("password-field")
                    Button {
                        showPassword.toggle()
                    } label: {
                        Image(systemName: showPassword ? "eye.slash" : "eye")
                            .frame(width: 44, height: 44)
                    }
                    .foregroundStyle(Color.white.opacity(0.72))
                    .accessibilityLabel(showPassword ? "Hide password" : "Show password")
                }
                .fieldStyle(isFocused: focused == .password)

                Text(PasswordRules.helper(password))
                    .font(.footnote)
                    .foregroundStyle(PasswordRules.isMet(password) ? Color(red: 0.55, green: 0.85, blue: 0.55) : Color.white.opacity(0.72))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .accessibilityLabel(PasswordRules.helper(password))

                Button {
                    isAdult.toggle()
                } label: {
                    HStack(alignment: .top, spacing: 12) {
                        Image(systemName: isAdult ? "checkmark.square.fill" : "square")
                            .font(.title3)
                            .foregroundStyle(isAdult ? Color.rfGold : Color.white.opacity(0.72))
                            .frame(width: 44, height: 44)
                        Text("I am 18 or older, and I am setting this up for myself.")
                            .font(.subheadline)
                            .foregroundStyle(Color.white.opacity(0.85))
                            .multilineTextAlignment(.leading)
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .padding(.top, 10)
                    }
                }
                .buttonStyle(.plain)
                .accessibilityLabel("I am 18 or older, and I am setting this up for myself")
                .accessibilityValue(isAdult ? "Checked" : "Unchecked")
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

            VStack(spacing: 4) {
                Text("By continuing you agree to our")
                    .foregroundStyle(Color.white.opacity(0.72))
                HStack(spacing: 6) {
                    Link("Terms of Service", destination: URL(string: "https://remainfaithful.com/terms")!)
                    Text("and")
                        .foregroundStyle(Color.white.opacity(0.72))
                    Link("Privacy Policy", destination: URL(string: "https://remainfaithful.com/privacy")!)
                }
            }
            .font(.footnote)
            .tint(Color.rfGold)
            .multilineTextAlignment(.center)
            .padding(.top, 16)

            Spacer()

            RFButton(
                title: isLoading ? "Creating account…" : "Get Started",
                isEnabled: canSubmit,
                isLoading: isLoading,
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
                    errorMsg = UserFacingError.message(for: e)
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
    @State private var resultMessage: String?
    @State private var errorMsg: String?

    private let green = Color(red: 0.20, green: 0.78, blue: 0.45)

    private var canSend: Bool {
        partnerEmail.contains("@") && partnerEmail.contains(".") && !isSending && resultMessage == nil
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
                subtitle: "If someone already invited you, accept it in Settings after you sign in. You don't have to invite them back."
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
                                  .foregroundColor(Color.white.opacity(0.72)))
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

                if let resultMessage {
                    Text(resultMessage)
                        .font(.subheadline)
                        .foregroundStyle(Color.white.opacity(0.85))
                        .multilineTextAlignment(.center)
                        .fixedSize(horizontal: false, vertical: true)
                    RFButton(title: "Continue", action: onContinue)
                } else {
                    RFButton(
                        title: isSending ? "Sending…" : "Send Invitation",
                        isEnabled: canSend,
                        isLoading: isSending,
                        action: sendInvite
                    )
                }
            }

            Spacer()

            Button(action: onSkip) {
                Text("Skip for now")
                    .font(.body)
                    .foregroundStyle(Color.white.opacity(0.72))
                    .frame(minHeight: 44)
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
                let result = try await APIClient.shared.invitePartner(email: addr)
                await MainActor.run {
                    isSending = false
                    resultMessage = result.message
                }
            } catch {
                await MainActor.run {
                    isSending = false
                    errorMsg = UserFacingError.message(for: error)
                }
            }
        }
    }
}

// MARK: - Step 4: Monitoring setup

private struct MonitoringSetupStep: View {
    let onContinue: () -> Void
    let onSkip: () -> Void

    @ObservedObject private var fcManager = FamilyControlsManager.shared
    @State private var didRequestScreenTime = false

    var body: some View {
        ScrollView(showsIndicators: false) {
            VStack(spacing: 0) {
                ZStack {
                    Circle().fill(Color.rfGold.opacity(0.08)).frame(width: 110, height: 110)
                    Circle().fill(Color.rfGold.opacity(0.13)).frame(width: 78, height: 78)
                    Image(systemName: "shield.fill")
                        .font(.system(size: 36))
                        .foregroundStyle(Color.rfGold)
                }
                .accessibilityHidden(true)
                .padding(.top, 12)
                .padding(.bottom, 24)

                stepHeader(
                    title: "Allow Screen Time",
                    subtitle: "Two things to set up now: Alerts and Blocking."
                )
                .padding(.bottom, 20)

                Text("Remain Faithful needs Apple’s Screen Time permission so Alerts and Blocking can use the apps you choose. Your partner cannot browse your phone. Only an alert category and a heartbeat leave this device — not photos, messages, or the screen itself. If you skip this, monitoring stays off. You can turn it on later in Settings.")
                    .font(.subheadline)
                    .foregroundStyle(Color.white.opacity(0.78))
                    .multilineTextAlignment(.leading)
                    .fixedSize(horizontal: false, vertical: true)
                    .padding(.bottom, 16)

                VStack(spacing: 12) {
                    tierRow(
                        icon: "bell.fill",
                        tint: Color(red: 0.28, green: 0.56, blue: 0.95),
                        title: "Alerts",
                        detail: "Your partner is told after about one minute of use in an app you choose, once a day."
                    )
                    tierRow(
                        icon: "hand.raised.fill",
                        tint: Color.rfGold,
                        title: "Blocking",
                        detail: "Chosen apps show Apple’s Screen Time shield. You can still ask Apple for an override."
                    )
                }
                .padding(.bottom, 20)

                if fcManager.authorizationStatus != .approved {
                    RFButton(
                        title: didRequestScreenTime ? "Waiting for permission…" : "Allow Screen Time",
                        isLoading: didRequestScreenTime
                    ) {
                        didRequestScreenTime = true
                        Task {
                            await fcManager.requestAuthorization()
                            await MainActor.run {
                                didRequestScreenTime = false
                                if fcManager.authorizationStatus == .approved {
                                    onContinue()
                                }
                            }
                        }
                    }
                    .padding(.bottom, 12)

                    RFButton(
                        title: "Skip for now — monitoring stays off",
                        role: .secondary,
                        action: onSkip
                    )
                } else {
                    RFButton(title: "Continue", action: onContinue)
                }
            }
            .padding(.horizontal, 24)
            .padding(.bottom, 28)
        }
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
                    .font(.subheadline)
                    .foregroundStyle(Color.white.opacity(0.72))
                    .fixedSize(horizontal: false, vertical: true)
            }
            Spacer()
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 14)
        .background(
            RoundedRectangle(cornerRadius: 14)
                .fill(Color.white.opacity(0.06))
                .overlay(RoundedRectangle(cornerRadius: 14)
                    .strokeBorder(RFStroke.lit, lineWidth: 1))
        )
    }
}

// MARK: - Step 5: Deep Scan

private struct DeepScanStep: View {
    let onContinue: () -> Void

    var body: some View {
        ScrollView(showsIndicators: false) {
            VStack(alignment: .leading, spacing: 16) {
                stepHeader(
                    title: DeepScanOnboarding.title,
                    subtitle: DeepScanOnboarding.subtitle
                )
                .frame(maxWidth: .infinity)

                Text(DeepScanOnboarding.summary)
                    .font(.body)
                    .foregroundStyle(Color.white.opacity(0.82))
                    .fixedSize(horizontal: false, vertical: true)

                VStack(alignment: .leading, spacing: 12) {
                    ForEach(Array(DeepScanOnboarding.steps.enumerated()), id: \.offset) { index, line in
                        HStack(alignment: .firstTextBaseline, spacing: 10) {
                            Text("\(index + 1).")
                                .font(.body.weight(.semibold))
                                .foregroundStyle(Color.rfGold)
                                .accessibilityHidden(true)
                            Text(line)
                                .font(.body)
                                .foregroundStyle(Color.white.opacity(0.88))
                                .fixedSize(horizontal: false, vertical: true)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                        .accessibilityElement(children: .combine)
                    }
                }

                Text(DeepScanOnboarding.limits)
                    .font(.body)
                    .foregroundStyle(Color.white.opacity(0.82))
                    .fixedSize(horizontal: false, vertical: true)

                RFButton(title: "Continue", action: onContinue)
                    .padding(.top, 8)
            }
            .padding(.horizontal, 24)
            .padding(.vertical, 24)
        }
    }
}

// MARK: - Shared helpers

private struct AppPickerStep: View {
    let onComplete: () -> Void

    @ObservedObject private var fcManager = FamilyControlsManager.shared
    @ObservedObject private var selection = ActivitySelectionManager.shared
    @State private var showingPicker = false

    private var selectedCount: Int {
        selection.selection.applicationTokens.count + selection.selection.categoryTokens.count
    }

    var body: some View {
        ScrollView(showsIndicators: false) {
            VStack(alignment: .leading, spacing: 16) {
                stepHeader(
                    title: "Choose apps",
                    subtitle: "Select the apps you want included. You can change this later."
                )
                .frame(maxWidth: .infinity)

                Text(selectedCount == 0
                     ? "No apps selected yet."
                     : selection.selectionSummary)
                    .font(.subheadline.weight(.semibold))
                    .foregroundStyle(Color.white.opacity(0.85))
                    .frame(maxWidth: .infinity, alignment: .center)

                if fcManager.authorizationStatus == .approved {
                    RFButton(title: "Choose apps") { showingPicker = true }
                    RFButton(
                        title: selectedCount == 0 ? "Save selected apps" : "Turn on Alerts",
                        isEnabled: selectedCount > 0
                    ) {
                        selection.selectionDidChange()
                        selection.setMonitoringEnabled(true)
                        NotificationService.shared.requestPermission(then: onComplete)
                    }
                } else {
                    Text("Screen Time access is still off, so the app list cannot open. Monitoring stays off until you allow it.")
                        .font(.subheadline)
                        .foregroundStyle(Color.white.opacity(0.78))
                    RFButton(title: "Allow Screen Time") {
                        Task { await fcManager.requestAuthorization() }
                    }
                }

                RFButton(title: "I'll do this later", role: .secondary) {
                    NotificationService.shared.requestPermission(then: onComplete)
                }
            }
            .padding(.horizontal, 24)
            .padding(.vertical, 24)
        }
        .familyActivityPicker(isPresented: $showingPicker, selection: $selection.selection)
        .onChange(of: selection.selection) { _, _ in
            selection.selectionDidChange()
        }
    }
}

private func stepHeader(title: String, subtitle: String) -> some View {
    VStack(spacing: 10) {
        Text(title)
            .font(.largeTitle.bold())
            .fontDesign(.serif)
            .foregroundStyle(.white)
            .multilineTextAlignment(.center)
            .lineSpacing(2)
        Text(subtitle)
            .font(.body)
            .foregroundStyle(Color.white.opacity(0.72))
            .multilineTextAlignment(.center)
    }
}

private extension View {
    func fieldStyle(isFocused: Bool) -> some View {
        self
            .padding(.horizontal, 18)
            .frame(minHeight: 44)
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
