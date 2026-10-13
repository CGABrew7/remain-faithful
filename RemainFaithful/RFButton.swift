import SwiftUI

// Premium navy-and-gold buttons. Primary is a gold gradient with navy text.
// Secondary is a navy fill with a thin gold border. Destructive stays red.
//
// Depth comes from the same three cues used by RFSurface: a lit top edge, a
// layered shadow (tight contact + wide ambient), and an inner highlight on the
// filled face. Pressing collapses the shadow and dips the scale, so the control
// reads as being pushed toward the screen rather than merely shrinking.
//
// Contrast (measured, see RFTheme):
//   primary    navy on gold   10.9:1 top → 5.1:1 bottom
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
                shape
                    .fill(LinearGradient(colors: [.rfGoldLight, .rfGoldDeep],
                                         startPoint: .top, endPoint: .bottom))
                    .overlay(innerHighlight)
            case .secondary:
                // Lifted a touch off the page so the fill reads as a surface
                // and not a hole, while white text stays at 12.4:1.
                shape
                    .fill(Color.rfNavy)
                    .overlay(shape.fill(Color.white.opacity(0.06)))
                    .overlay(innerHighlight.opacity(0.5))
            case .destructive:
                shape
                    .fill(LinearGradient(
                        colors: [Color.rfDanger.opacity(0.98), Color.rfDanger.opacity(0.82)],
                        startPoint: .top, endPoint: .bottom))
                    .overlay(innerHighlight.opacity(0.6))
            }
        }
    }

    /// Soft sheen across the top third, as if lit from above.
    private var innerHighlight: some View {
        LinearGradient(
            colors: [Color.white.opacity(0.38), Color.white.opacity(0)],
            startPoint: .top,
            endPoint: .center
        )
        .clipShape(shape)
        .allowsHitTesting(false)
    }

    // MARK: - Border

    @ViewBuilder
    private var border: some View {
        if !appearsEnabled {
            shape.strokeBorder(Color.white.opacity(0.10), lineWidth: 1)
        } else {
            switch role {
            case .primary:
                // Darker gold rim grounds the gradient against the page.
                shape.strokeBorder(Color.rfGoldDeep.opacity(0.55), lineWidth: 1)
            case .secondary:
                shape.strokeBorder(RFStroke.gold, lineWidth: 1.5)
            case .destructive:
                shape.strokeBorder(Color.white.opacity(0.20), lineWidth: 1)
            }
        }
    }

    private var shadowColor: Color {
        guard appearsEnabled else { return .clear }
        switch role {
        case .primary:     return Color.rfGoldDeep.opacity(0.42)
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

/// Press feel: a small dip plus a shadow that collapses toward the surface.
/// Both are driven by one spring so they stay in step.
private struct RFPressStyle: ButtonStyle {
    var animates: Bool
    var liftsShadow: Bool
    var shadowTint: Color

    func makeBody(configuration: Configuration) -> some View {
        let pressed = animates && configuration.isPressed
        return configuration.label
            .scaleEffect(pressed ? 0.975 : 1)
            .brightness(pressed ? -0.04 : 0)
            // Tight contact shadow.
            .shadow(color: liftsShadow ? shadowTint.opacity(0.9) : .clear,
                    radius: pressed ? 1 : 3,
                    x: 0, y: pressed ? 1 : 2)
            // Wide ambient shadow; collapses on press so the control sinks.
            .shadow(color: liftsShadow ? shadowTint.opacity(0.55) : .clear,
                    radius: pressed ? 4 : 12,
                    x: 0, y: pressed ? 2 : 7)
            .animation(.spring(response: 0.26, dampingFraction: 0.7),
                       value: configuration.isPressed)
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
