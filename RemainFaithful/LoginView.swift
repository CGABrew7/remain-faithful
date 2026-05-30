import SwiftUI
import AuthenticationServices
import GoogleSignIn

// MARK: - LoginView

struct LoginView: View {
    @AppStorage("hasCompletedOnboarding") private var hasCompletedOnboarding = false
    @AppStorage("userName")  private var storedName  = ""
    @AppStorage("userEmail") private var storedEmail = ""

    @State private var email    = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMsg: String?
    @State private var showForgotPassword = false

    @FocusState private var focused: Field?
    private enum Field { case email, password }

    private var canSubmit: Bool {
        email.contains("@") && email.contains(".") && password.count >= 8 && !isLoading
    }

    var body: some View {
        ZStack {
            LinearGradient(colors: [Color.rfNavy, Color.rfNavyMid],
                           startPoint: .topLeading, endPoint: .bottomTrailing)
                .ignoresSafeArea()

            ScrollView(showsIndicators: false) {
                VStack(spacing: 0) {
                    logoSection
                        .padding(.top, 60)
                        .padding(.bottom, 44)

                    formSection
                        .padding(.horizontal, 24)

                    if let msg = errorMsg {
                        Text(msg)
                            .font(.system(size: 13))
                            .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                            .multilineTextAlignment(.center)
                            .padding(.top, 12)
                            .padding(.horizontal, 24)
                            .transition(.opacity)
                    }

                    signInButton
                        .padding(.horizontal, 24)
                        .padding(.top, 20)

                    forgotPasswordButton
                        .padding(.top, 14)

                    orDivider
                        .padding(.horizontal, 24)
                        .padding(.vertical, 24)

                    socialButtons
                        .padding(.horizontal, 24)

                    Spacer(minLength: 48)
                }
            }
        }
        .animation(.easeInOut(duration: 0.2), value: errorMsg)
        .sheet(isPresented: $showForgotPassword) {
            ForgotPasswordSheet()
                .presentationDetents([.height(340)])
                .presentationDragIndicator(.hidden)
        }
    }

    // MARK: - Logo

    private var logoSection: some View {
        VStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.08))
                    .frame(width: 100, height: 100)
                Circle()
                    .fill(Color.rfGold.opacity(0.12))
                    .frame(width: 72, height: 72)
                Image(systemName: "cross.fill")
                    .font(.system(size: 30, weight: .medium))
                    .foregroundStyle(Color.rfGold)
            }
            VStack(spacing: 6) {
                Text("Welcome Back")
                    .font(.system(size: 30, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                Text("Sign in to continue your journey")
                    .font(.system(size: 15))
                    .foregroundStyle(Color.white.opacity(0.55))
            }
        }
    }

    // MARK: - Form

    private var formSection: some View {
        VStack(spacing: 14) {
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
            .loginFieldStyle(isFocused: focused == .email)

            // Password
            HStack(spacing: 12) {
                Image(systemName: "lock.fill")
                    .font(.system(size: 15))
                    .foregroundStyle(focused == .password ? Color.rfGold : Color.white.opacity(0.45))
                    .frame(width: 22)
                SecureField("", text: $password,
                            prompt: Text("Password").foregroundColor(Color.white.opacity(0.38)))
                    .textContentType(.password)
                    .foregroundStyle(.white)
                    .focused($focused, equals: .password)
                    .submitLabel(.done)
                    .onSubmit { if canSubmit { signIn() } }
            }
            .loginFieldStyle(isFocused: focused == .password)
        }
    }

    // MARK: - Sign In button

    private var signInButton: some View {
        Button(action: signIn) {
            Group {
                if isLoading {
                    ProgressView().tint(Color.rfNavy)
                } else {
                    Text("Sign In")
                        .font(.system(size: 17, weight: .semibold))
                        .foregroundStyle(canSubmit ? Color.rfNavy : Color.white.opacity(0.35))
                }
            }
            .frame(maxWidth: .infinity)
            .frame(height: 56)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(canSubmit ? Color.rfGold : Color.white.opacity(0.08))
            )
        }
        .disabled(!canSubmit)
    }

    // MARK: - Forgot password

    private var forgotPasswordButton: some View {
        Button { showForgotPassword = true } label: {
            Text("Forgot Password?")
                .font(.system(size: 14))
                .foregroundStyle(Color.rfGold.opacity(0.80))
        }
    }

    // MARK: - Or divider

    private var orDivider: some View {
        HStack(spacing: 14) {
            Rectangle()
                .fill(Color.white.opacity(0.12))
                .frame(height: 1)
            Text("or")
                .font(.system(size: 13))
                .foregroundStyle(Color.white.opacity(0.35))
            Rectangle()
                .fill(Color.white.opacity(0.12))
                .frame(height: 1)
        }
    }

    // MARK: - Social buttons

    private var socialButtons: some View {
        VStack(spacing: 12) {
            // Apple Sign In
            SignInWithAppleButton(.signIn,
                onRequest: { req in
                    req.requestedScopes = [.fullName, .email]
                },
                onCompletion: handleAppleCompletion
            )
            .signInWithAppleButtonStyle(.white)
            .frame(height: 54)
            .cornerRadius(14)

            // Google Sign In
            Button(action: handleGoogleSignIn) {
                HStack(spacing: 12) {
                    // Google "G" logo approximation
                    ZStack {
                        Text("G")
                            .font(.system(size: 17, weight: .bold))
                            .foregroundStyle(Color(red: 0.26, green: 0.52, blue: 0.96))
                    }
                    .frame(width: 22)
                    Text("Sign in with Google")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundStyle(Color(red: 0.10, green: 0.10, blue: 0.10))
                    Spacer()
                }
                .padding(.horizontal, 16)
                .frame(maxWidth: .infinity)
                .frame(height: 54)
                .background(RoundedRectangle(cornerRadius: 14).fill(.white))
                .overlay(RoundedRectangle(cornerRadius: 14)
                    .stroke(Color(white: 0.82), lineWidth: 1))
            }
        }
    }

    // MARK: - Actions

    private func signIn() {
        guard canSubmit else { return }
        focused = nil
        isLoading = true
        errorMsg  = nil
        Task {
            do {
                let resp = try await APIClient.shared.login(email: email, password: password)
                await MainActor.run {
                    storedName  = resp.user.name
                    storedEmail = resp.user.email
                    AuthState.shared.setSession(token: resp.token, user: resp.user)
                }
            } catch let e as APIError {
                await MainActor.run {
                    errorMsg  = e.errorDescription
                    isLoading = false
                }
            } catch {
                await MainActor.run {
                    errorMsg  = "Network error — check your connection"
                    isLoading = false
                }
            }
        }
    }

    private func handleAppleCompletion(_ result: Result<ASAuthorization, Error>) {
        switch result {
        case .failure(let err):
            if (err as? ASAuthorizationError)?.code == .canceled { return }
            errorMsg = err.localizedDescription
        case .success(let auth):
            guard let cred = auth.credential as? ASAuthorizationAppleIDCredential,
                  let tokenData = cred.identityToken,
                  let identityToken = String(data: tokenData, encoding: .utf8),
                  let codeData = cred.authorizationCode,
                  let authCode = String(data: codeData, encoding: .utf8) else {
                errorMsg = "Apple Sign In failed — missing credentials"
                return
            }
            let firstName = cred.fullName?.givenName ?? ""
            let lastName  = cred.fullName?.familyName ?? ""
            isLoading = true
            errorMsg  = nil
            Task {
                do {
                    let resp = try await APIClient.shared.appleSignIn(
                        identityToken: identityToken,
                        authorizationCode: authCode,
                        firstName: firstName,
                        lastName: lastName
                    )
                    await MainActor.run {
                        storedName  = resp.user.name
                        storedEmail = resp.user.email
                        hasCompletedOnboarding = true
                        AuthState.shared.setSession(token: resp.token, user: resp.user)
                    }
                } catch {
                    await MainActor.run {
                        errorMsg  = "Apple Sign In failed — \(error.localizedDescription)"
                        isLoading = false
                    }
                }
            }
        }
    }

    private func handleGoogleSignIn() {
        guard GIDSignIn.sharedInstance.configuration != nil else {
            errorMsg = "Google Sign In is not available"
            return
        }
        guard let scene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
              let rootVC = scene.windows.first(where: { $0.isKeyWindow })?.rootViewController else { return }
        isLoading = true
        errorMsg  = nil
        GIDSignIn.sharedInstance.signIn(withPresenting: rootVC) { result, error in
            if let error {
                DispatchQueue.main.async {
                    errorMsg  = "Google Sign In failed — \(error.localizedDescription)"
                    isLoading = false
                }
                return
            }
            guard let idToken = result?.user.idToken?.tokenString else {
                DispatchQueue.main.async {
                    errorMsg  = "Google Sign In failed — missing ID token"
                    isLoading = false
                }
                return
            }
            Task {
                do {
                    let resp = try await APIClient.shared.googleSignIn(idToken: idToken)
                    await MainActor.run {
                        storedName  = resp.user.name
                        storedEmail = resp.user.email
                        hasCompletedOnboarding = true
                        AuthState.shared.setSession(token: resp.token, user: resp.user)
                    }
                } catch {
                    await MainActor.run {
                        errorMsg  = "Google Sign In failed — \(error.localizedDescription)"
                        isLoading = false
                    }
                }
            }
        }
    }
}

// MARK: - Forgot Password sheet

private struct ForgotPasswordSheet: View {
    @Environment(\.dismiss) private var dismiss
    @State private var email = ""
    @State private var sent  = false
    @State private var isLoading = false
    @State private var errorMsg: String?
    @FocusState private var focused: Bool

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 20)

                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Forgot Password")
                            .font(.system(size: 20, weight: .bold, design: .serif))
                            .foregroundStyle(.white)
                        Text("We'll send a reset link to your email")
                            .font(.system(size: 13))
                            .foregroundStyle(Color.white.opacity(0.45))
                    }
                    Spacer()
                    Button { dismiss() } label: {
                        Image(systemName: "xmark")
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundStyle(Color.white.opacity(0.55))
                            .padding(10)
                            .background(Circle().fill(Color.white.opacity(0.09)))
                    }
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 24)

                Divider().overlay(Color.white.opacity(0.08))
                    .padding(.bottom, 24)

                if sent {
                    VStack(spacing: 12) {
                        Image(systemName: "checkmark.circle.fill")
                            .font(.system(size: 40))
                            .foregroundStyle(Color(red: 0.20, green: 0.78, blue: 0.45))
                        Text("Check your email")
                            .font(.system(size: 17, weight: .semibold))
                            .foregroundStyle(.white)
                        Text("If an account exists for \(email),\nyou'll receive a reset link shortly.")
                            .font(.system(size: 14))
                            .foregroundStyle(Color.white.opacity(0.55))
                            .multilineTextAlignment(.center)
                    }
                    .padding(.horizontal, 24)
                } else {
                    VStack(spacing: 14) {
                        HStack(spacing: 12) {
                            Image(systemName: "envelope.fill")
                                .font(.system(size: 15))
                                .foregroundStyle(focused ? Color.rfGold : Color.white.opacity(0.45))
                                .frame(width: 22)
                            TextField("", text: $email,
                                      prompt: Text("Email address").foregroundColor(Color.white.opacity(0.38)))
                                .keyboardType(.emailAddress)
                                .textContentType(.emailAddress)
                                .textInputAutocapitalization(.never)
                                .autocorrectionDisabled()
                                .foregroundStyle(.white)
                                .focused($focused)
                                .submitLabel(.send)
                                .onSubmit { sendReset() }
                        }
                        .loginFieldStyle(isFocused: focused)

                        if let msg = errorMsg {
                            Text(msg)
                                .font(.system(size: 13))
                                .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                        }

                        Button(action: sendReset) {
                            Text(isLoading ? "Sending…" : "Send Reset Link")
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundStyle(canSend ? Color.rfNavy : Color.white.opacity(0.35))
                                .frame(maxWidth: .infinity)
                                .frame(height: 54)
                                .background(
                                    RoundedRectangle(cornerRadius: 14)
                                        .fill(canSend ? Color.rfGold : Color.white.opacity(0.08))
                                )
                        }
                        .disabled(!canSend)
                    }
                    .padding(.horizontal, 24)
                }
                Spacer()
            }
        }
    }

    private var canSend: Bool {
        email.contains("@") && email.contains(".") && !isLoading
    }

    private func sendReset() {
        guard canSend else { return }
        focused = false
        isLoading = true
        errorMsg  = nil
        Task {
            do {
                try await APIClient.shared.forgotPassword(email: email)
                await MainActor.run { sent = true }
            } catch {
                await MainActor.run {
                    errorMsg  = "Failed to send reset link — check your connection"
                    isLoading = false
                }
            }
        }
    }
}

// MARK: - Field style helper (local to LoginView)

private extension View {
    func loginFieldStyle(isFocused: Bool) -> some View {
        self
            .padding(.horizontal, 18)
            .frame(height: 56)
            .background(
                RoundedRectangle(cornerRadius: 14)
                    .fill(Color.white.opacity(isFocused ? 0.11 : 0.07))
                    .overlay(
                        RoundedRectangle(cornerRadius: 14)
                            .stroke(isFocused ? Color.rfGold : Color.white.opacity(0.10),
                                    lineWidth: 1.5)
                    )
            )
            .animation(.easeInOut(duration: 0.18), value: isFocused)
    }
}

#Preview {
    LoginView()
}
