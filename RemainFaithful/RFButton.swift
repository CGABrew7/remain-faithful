import SwiftUI

// Flat buttons. Primary is solid #F0D030 with navy text.
// Secondary is a navy fill with a thin yellow border. Destructive stays red.
// Press darkens slightly and scales to 0.98. No gradient, shine, or glow.
//
// Contrast (measured, see RFTheme):
//   primary    navy on #F0D030   10.22:1
//   secondary  white on navy  12.4:1, gold border 5.7:1
//   destructive white on red   7.7:1
//   disabled   white on tint   6.5:1

enum RFButtonRole {
    case primary
    case secondary
    case destructive
}

struct RFButton: View {
    let title: String
    var role: RFButtonRole = .primary
    var isEnabled: Bool = true
    /// Shows a spinner, keeps the title, and blocks input without greying out.
    var isLoading: Bool = false
    /// Set when motion would be distracting (a static control, or Reduce Motion).
    var staticPress: Bool = false
    let action: () -> Void

    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    /// Loading is interactive-looking but not tappable.
    private var isInteractive: Bool { isEnabled && !isLoading }

    /// While loading the button keeps its full colour rather than greying out,
    /// so a brief save does not make the control look broken. Call sites whose
    /// `isEnabled` already drops to false during a request still render as the
    /// live button for as long as `isLoading` is set.
    private var appearsEnabled: Bool { isEnabled || isLoading }

    var body: some View {
        Button(action: {
            guard isInteractive else { return }
            haptic()
            action()
        }) {
            label
                .frame(maxWidth: .infinity)
                .frame(minHeight: RFMetrics.minTap)
                .padding(.vertical, 12)
                .padding(.horizontal, 16)
                .background(face)
                .overlay(border)
                .clipShape(shape)
                .contentShape(shape)
        }
        .buttonStyle(
            RFPressStyle(
                animates: isInteractive && !staticPress && !reduceMotion,
                liftsShadow: appearsEnabled,
                shadowTint: shadowColor
            )
        )
        .disabled(!isInteractive)
        .accessibilityAddTraits(.isButton)
        .accessibilityLabel(Text(title))
        .accessibilityValue(isLoading ? Text("Working") : Text(""))
        .accessibilityHint(isEnabled || isLoading ? Text("") : Text("Unavailable"))
    }

    private var shape: RoundedRectangle {
        RoundedRectangle(cornerRadius: RFMetrics.controlRadius, style: .continuous)
    }

    // MARK: - Label

    private var label: some View {
        HStack(spacing: 10) {
            if isLoading {
                ProgressView()
                    .progressViewStyle(.circular)
                    .controlSize(.small)
                    .tint(foreground)
                    .accessibilityHidden(true)
            }
            Text(title)
                .font(.body.weight(.semibold))
                .foregroundStyle(foreground)
                .multilineTextAlignment(.center)
                .lineLimit(3)
                .minimumScaleFactor(0.85)
        }
    }

    private var foreground: Color {
        // Disabled text measures 6.5:1 on the disabled fill.
        guard appearsEnabled else { return Color.white.opacity(0.70) }
        switch role {
        case .primary:     return Color.rfNavy
        case .secondary:   return Color.white
        case .destructive: return Color.white
        }
    }

    // MARK: - Face

    @ViewBuilder
    private var face: some View {
        if !appearsEnabled {
            shape.fill(Color.white.opacity(0.08))
        } else {
            switch role {
            case .primary:
                shape.fill(Color.rfGold)
            case .secondary:
                shape
                    .fill(Color.rfNavy)
                    .overlay(shape.fill(Color.white.opacity(0.04)))
            case .destructive:
                shape.fill(Color.rfDanger)
            }
        }
    }

    // MARK: - Border

    @ViewBuilder
    private var border: some View {
        if !appearsEnabled {
            shape.strokeBorder(Color.white.opacity(0.10), lineWidth: 1)
        } else {
            switch role {
            case .primary:
                shape.strokeBorder(Color.clear, lineWidth: 0)
            case .secondary:
                shape.strokeBorder(Color.rfGold, lineWidth: 1.5)
            case .destructive:
                shape.strokeBorder(Color.white.opacity(0.20), lineWidth: 1)
            }
        }
    }

    private var shadowColor: Color {
        guard appearsEnabled else { return .clear }
        switch role {
        case .primary:     return .clear
        case .destructive: return Color.rfDanger.opacity(0.42)
        case .secondary:   return Color.black.opacity(0.30)
        }
    }

    private func haptic() {
        switch role {
        case .destructive:
            UIImpactFeedbackGenerator(style: .rigid).impactOccurred()
        default:
            UIImpactFeedbackGenerator(style: .light).impactOccurred()
        }
    }
}

/// Press feel: a slight darken and a scale of 0.98. No bounce.
private struct RFPressStyle: ButtonStyle {
    var animates: Bool
    var liftsShadow: Bool
    var shadowTint: Color

    func makeBody(configuration: Configuration) -> some View {
        let pressed = animates && configuration.isPressed
        _ = (liftsShadow, shadowTint)
        return configuration.label
            .scaleEffect(pressed ? 0.98 : 1)
            .brightness(pressed ? -0.06 : 0)
            .animation(.easeInOut(duration: 0.16), value: configuration.isPressed)
    }
}

/// Visual glyph stays about 30 pt. The hit target is 44 pt.
struct RFCloseButton: View {
    var label: String = "Close"
    let action: () -> Void

    var body: some View {
        Button(action: {
            UIImpactFeedbackGenerator(style: .light).impactOccurred()
            action()
        }) {
            Image(systemName: "xmark")
                .font(.system(size: 14, weight: .semibold))
                .foregroundStyle(Color.white.opacity(0.80))
                .frame(width: 30, height: 30)
                .background(
                    Circle()
                        .fill(Color.white.opacity(0.10))
                        .overlay(Circle().strokeBorder(Color.white.opacity(0.14), lineWidth: 1))
                )
                .frame(width: RFMetrics.minTap, height: RFMetrics.minTap)
                .contentShape(Rectangle())
        }
        .buttonStyle(RFPressStyle(animates: true, liftsShadow: false, shadowTint: .clear))
        .accessibilityLabel(label)
    }
}
