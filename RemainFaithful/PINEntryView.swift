import SwiftUI

// 4-digit PIN entry sheet used when a partner-set PIN gates a security action.
//
// Usage:
//   PINEntryView(title:..., subtitle:...) { pin in
//       let ok = await PartnerPINManager.shared.verifyPIN(pin)
//       if ok { /* dismiss & execute action */ }
//       return ok
//   } onCancel: { /* dismiss without action */ }
//
// The caller owns presentation state. PINEntryView does NOT dismiss itself —
// the onVerify and onCancel closures are where the caller should set
// their showPINGate flag to false (after executing or cancelling).
struct PINEntryView: View {
    let title:    String
    let subtitle: String
    let onVerify: (String) async -> Bool
    let onCancel: () -> Void

    @State private var digits:         [Int]    = []
    @State private var isVerifying:    Bool     = false
    @State private var errorMessage:   String?  = nil
    @State private var indicatorOffset: CGFloat = 0

    var body: some View {
        ZStack {
            Color.rfNavy.ignoresSafeArea()
            VStack(spacing: 0) {
                dragPill
                lockIcon
                titleBlock
                Spacer().frame(height: 28)
                indicators
                errorLabel
                Spacer().frame(height: 28)
                numpad
                cancelButton
            }
        }
        .onChange(of: digits.count) { _, count in
            guard count == 4, !isVerifying else { return }
            submit()
        }
    }

    // MARK: - Sub-views

    private var dragPill: some View {
        RoundedRectangle(cornerRadius: 3)
            .fill(Color.white.opacity(0.18))
            .frame(width: 40, height: 4)
            .padding(.top, 12)
            .padding(.bottom, 28)
    }

    private var lockIcon: some View {
        ZStack {
            Circle()
                .fill(Color.rfGold.opacity(0.12))
                .frame(width: 72, height: 72)
            Image(systemName: "lock.shield.fill")
                .font(.system(size: 32))
                .foregroundStyle(Color.rfGold)
        }
        .padding(.bottom, 20)
    }

    private var titleBlock: some View {
        VStack(spacing: 8) {
            Text(title)
                .font(.system(size: 22, weight: .bold, design: .serif))
                .foregroundStyle(.white)
                .multilineTextAlignment(.center)
                .padding(.horizontal, 24)
            Text(subtitle)
                .font(.system(size: 14))
                .foregroundStyle(Color.white.opacity(0.55))
                .multilineTextAlignment(.center)
                .padding(.horizontal, 36)
                .lineSpacing(3)
        }
    }

    private var indicators: some View {
        HStack(spacing: 20) {
            ForEach(0..<4, id: \.self) { index in
                Circle()
                    .fill(index < digits.count ? Color.rfGold : Color.white.opacity(0.14))
                    .frame(width: 16, height: 16)
                    .overlay(Circle().stroke(Color.white.opacity(0.22), lineWidth: 1))
            }
        }
        .offset(x: indicatorOffset)
        .animation(.interpolatingSpring(stiffness: 500, damping: 8), value: indicatorOffset)
        .padding(.bottom, errorMessage == nil ? 0 : 0)
    }

    @ViewBuilder
    private var errorLabel: some View {
        if let error = errorMessage {
            Text(error)
                .font(.system(size: 13))
                .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                .multilineTextAlignment(.center)
                .padding(.horizontal, 36)
                .padding(.top, 12)
                .transition(.opacity)
        } else {
            Color.clear.frame(height: 0)
        }
    }

    private var numpad: some View {
        let rows: [[String]] = [
            ["1", "2", "3"],
            ["4", "5", "6"],
            ["7", "8", "9"],
            ["",  "0", "⌫"],
        ]
        return VStack(spacing: 10) {
            ForEach(rows, id: \.self) { row in
                HStack(spacing: 10) {
                    ForEach(row, id: \.self) { key in
                        numpadKey(key)
                    }
                }
            }
        }
        .padding(.horizontal, 44)
        .padding(.bottom, 12)
    }

    private func numpadKey(_ key: String) -> some View {
        let isDelete  = key == "⌫"
        let isEmpty   = key.isEmpty
        let disabled  = isEmpty || (digits.count == 4 && !isDelete) || isVerifying

        return Button {
            handleKey(key)
        } label: {
            Group {
                if isDelete {
                    Image(systemName: "delete.left")
                        .font(.system(size: 18, weight: .medium))
                        .foregroundStyle(Color.white.opacity(0.75))
                } else if isEmpty {
                    Color.clear
                } else {
                    Text(key)
                        .font(.system(size: 22, weight: .semibold))
                        .foregroundStyle(.white)
                }
            }
            .frame(maxWidth: .infinity)
            .frame(height: 62)
            .background {
                if !isEmpty {
                    RoundedRectangle(cornerRadius: 14)
                        .fill(Color.white.opacity(0.07))
                        .overlay(RoundedRectangle(cornerRadius: 14)
                            .stroke(Color.white.opacity(0.06), lineWidth: 1))
                }
            }
        }
        .disabled(disabled)
    }

    private var cancelButton: some View {
        Button(action: onCancel) {
            Text("Cancel")
                .font(.system(size: 15, weight: .semibold))
                .foregroundStyle(Color.white.opacity(0.50))
                .frame(maxWidth: .infinity)
                .frame(height: 44)
        }
        .padding(.horizontal, 44)
        .padding(.bottom, 24)
    }

    // MARK: - Logic

    private func handleKey(_ key: String) {
        withAnimation(.easeInOut(duration: 0.08)) {
            errorMessage = nil
            if key == "⌫" {
                if !digits.isEmpty { digits.removeLast() }
            } else if let d = Int(key), digits.count < 4 {
                digits.append(d)
            }
        }
    }

    private func submit() {
        let pin = digits.map(String.init).joined()
        isVerifying = true
        Task {
            let success = await onVerify(pin)
            await MainActor.run {
                isVerifying = false
                if !success {
                    shake()
                    withAnimation { errorMessage = "Incorrect PIN. Contact your accountability partner." }
                    digits = []
                }
                // On success the caller's onVerify closure closes the sheet.
            }
        }
    }

    private func shake() {
        indicatorOffset = 10
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.15) {
            indicatorOffset = -10
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.15) {
                indicatorOffset = 6
                DispatchQueue.main.asyncAfter(deadline: .now() + 0.12) {
                    indicatorOffset = 0
                }
            }
        }
    }
}

// MARK: - Set-PIN sheet (partner side)

// Shown in ManagePartnersView for relationships where the current user is
// the monitoring partner. Lets them set or change the 4-digit protection PIN.
struct SetPINSheet: View {
    let relationshipID: Int
    let partnerName:    String
    let hasExistingPIN: Bool
    let onDone: () -> Void

    @Environment(\.dismiss) private var dismiss

    @State private var digits:        [Int]    = []
    @State private var confirmDigits: [Int]    = []
    @State private var phase:         SetPINPhase = .enter
    @State private var isSaving:      Bool     = false
    @State private var errorMessage:  String?  = nil
    @State private var saveError:     String?  = nil

    enum SetPINPhase { case enter, confirm, done }

    private let gold = Color.rfGold
    private let navy = Color.rfNavy

    var body: some View {
        ZStack {
            Color(red: 0.07, green: 0.11, blue: 0.24).ignoresSafeArea()
            VStack(spacing: 0) {
                RoundedRectangle(cornerRadius: 3)
                    .fill(Color.white.opacity(0.18))
                    .frame(width: 40, height: 4)
                    .padding(.top, 12)
                    .padding(.bottom, 24)

                ZStack {
                    Circle()
                        .fill(gold.opacity(0.12))
                        .frame(width: 64, height: 64)
                    Image(systemName: "lock.fill")
                        .font(.system(size: 28))
                        .foregroundStyle(gold)
                }
                .padding(.bottom, 16)

                Text(phase == .confirm ? "Confirm PIN" : (hasExistingPIN ? "Change PIN" : "Set Protection PIN"))
                    .font(.system(size: 20, weight: .bold, design: .serif))
                    .foregroundStyle(.white)
                    .padding(.bottom, 8)

                Text(phase == .confirm
                     ? "Enter the same PIN again to confirm."
                     : "Set a 4-digit PIN for \(partnerName). Share it verbally — it will be required to change protection settings.")
                    .font(.system(size: 13))
                    .foregroundStyle(Color.white.opacity(0.55))
                    .multilineTextAlignment(.center)
                    .padding(.horizontal, 28)
                    .lineSpacing(3)
                    .padding(.bottom, 24)

                // Indicators
                HStack(spacing: 20) {
                    let current = phase == .confirm ? confirmDigits : digits
                    ForEach(0..<4, id: \.self) { i in
                        Circle()
                            .fill(i < current.count ? gold : Color.white.opacity(0.14))
                            .frame(width: 16, height: 16)
                            .overlay(Circle().stroke(Color.white.opacity(0.22), lineWidth: 1))
                    }
                }
                .padding(.bottom, 8)

                if let err = errorMessage {
                    Text(err)
                        .font(.system(size: 12))
                        .foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                        .padding(.top, 6)
                        .padding(.horizontal, 28)
                }

                Spacer().frame(height: 24)

                // Numpad
                let rows: [[String]] = [["1","2","3"],["4","5","6"],["7","8","9"],["","0","⌫"]]
                VStack(spacing: 10) {
                    ForEach(rows, id: \.self) { row in
                        HStack(spacing: 10) {
                            ForEach(row, id: \.self) { key in
                                setPINKey(key)
                            }
                        }
                    }
                }
                .padding(.horizontal, 44)
                .padding(.bottom, 12)

                if let err = saveError {
                    Text(err).font(.system(size: 12)).foregroundStyle(Color(red: 0.90, green: 0.35, blue: 0.35))
                        .padding(.horizontal, 28).padding(.bottom, 8)
                }

                Button { dismiss() } label: {
                    Text("Cancel")
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(Color.white.opacity(0.50))
                        .frame(maxWidth: .infinity).frame(height: 44)
                }
                .padding(.horizontal, 44).padding(.bottom, 24)
            }
        }
        .onChange(of: digits.count) { _, count in
            guard count == 4, phase == .enter else { return }
            phase = .confirm
        }
        .onChange(of: confirmDigits.count) { _, count in
            guard count == 4, phase == .confirm else { return }
            if confirmDigits == digits {
                savePIN()
            } else {
                withAnimation { errorMessage = "PINs don't match — try again." }
                confirmDigits = []
            }
        }
    }

    private func setPINKey(_ key: String) -> some View {
        let current = phase == .confirm ? confirmDigits : digits
        let disabled = key.isEmpty || (current.count == 4 && key != "⌫") || isSaving
        return Button {
            errorMessage = nil
            if key == "⌫" {
                if phase == .confirm {
                    if !confirmDigits.isEmpty { confirmDigits.removeLast() }
                } else {
                    if !digits.isEmpty { digits.removeLast() }
                }
            } else if let d = Int(key), current.count < 4 {
                if phase == .confirm { confirmDigits.append(d) }
                else { digits.append(d) }
            }
        } label: {
            Group {
                if key == "⌫" {
                    Image(systemName: "delete.left").font(.system(size: 18, weight: .medium)).foregroundStyle(Color.white.opacity(0.75))
                } else if key.isEmpty {
                    Color.clear
                } else {
                    Text(key).font(.system(size: 22, weight: .semibold)).foregroundStyle(.white)
                }
            }
            .frame(maxWidth: .infinity).frame(height: 62)
            .background {
            if !key.isEmpty {
                RoundedRectangle(cornerRadius: 14).fill(Color.white.opacity(0.07))
            }
        }
        }
        .disabled(disabled)
    }

    private func savePIN() {
        isSaving = true
        let pin        = digits.map(String.init).joined()
        let wasChanged = hasExistingPIN
        Task {
            do {
                try await APIClient.shared.setRelationshipPIN(relationshipID: relationshipID, pin: pin)
                // Notify the monitored user that their PIN protection was changed.
                if wasChanged {
                    try? await APIClient.shared.sendProtectionAlert(
                        type: "pin_changed",
                        detail: "Protection PIN was changed by your accountability partner."
                    )
                }
                await MainActor.run {
                    isSaving = false
                    onDone()
                    dismiss()
                }
            } catch {
                await MainActor.run {
                    isSaving = false
                    saveError = "Couldn't save PIN. Try again."
                    confirmDigits = []
                    phase = .enter
                    digits = []
                }
            }
        }
    }
}
