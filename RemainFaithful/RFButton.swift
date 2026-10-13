import SwiftUI

// Premium navy-and-gold buttons. Primary is a gold gradient with navy text.
// Secondary is a navy fill with a thin gold border. Destructive stays red.

enum RFButtonRole {
    case primary
    case secondary
    case destructive
}

struct RFButton: View {
    let title: String
    var role: RFButtonRole = .primary
    var isEnabled: Bool = true
    /// Set when motion would be distracting (a static control, or Reduce Motion).
    var staticPress: Bool = false
    let action: () -> Void

    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        Button(action: {
            if isEnabled {
                UIImpactFeedbackGenerator(style: .light).impactOccurred()
            }
            action()
        }) {
            Text(title)
                .font(.body.weight(.semibold))
                .foregroundStyle(foreground)
                .multilineTextAlignment(.center)
                .lineLimit(3)
                .minimumScaleFactor(0.85)
                .frame(maxWidth: .infinity)
                .frame(minHeight: 44)
                .padding(.vertical, 12)
                .padding(.horizontal, 16)
                .background(background)
                .overlay(border)
                .shadow(color: shadowColor, radius: role == .primary && isEnabled ? 10 : 0, x: 0, y: 6)
        }
        .buttonStyle(RFPressStyle(enabled: isEnabled && !staticPress && !reduceMotion))
        .disabled(!isEnabled)
        .accessibilityAddTraits(.isButton)
    }

    private var foreground: Color {
        guard isEnabled else { return Color.white.opacity(0.55) }
        switch role {
        case .primary: return Color.rfNavy
        case .secondary: return Color.white
        case .destructive: return Color.white
        }
    }

    @ViewBuilder
    private var background: some View {
        let shape = RoundedRectangle(cornerRadius: 14, style: .continuous)
        if !isEnabled {
            shape.fill(Color.white.opacity(0.10))
        } else {
            switch role {
            case .primary:
                shape.fill(
                    LinearGradient(
                        colors: [
                            Color(red: 0.95, green: 0.84, blue: 0.52),
                            Color(red: 0.72, green: 0.55, blue: 0.18)
                        ],
                        startPoint: .top,
                        endPoint: .bottom
                    )
                )
                .overlay(
                    LinearGradient(
                        colors: [Color.white.opacity(0.42), Color.white.opacity(0)],
                        startPoint: .top,
                        endPoint: .center
                    )
                    .clipShape(shape)
                    .allowsHitTesting(false)
                )
            case .secondary:
                shape.fill(Color.rfNavy)
            case .destructive:
                shape.fill(Color(red: 0.62, green: 0.14, blue: 0.16))
            }
        }
    }

    @ViewBuilder
    private var border: some View {
        let shape = RoundedRectangle(cornerRadius: 14, style: .continuous)
        switch role {
        case .secondary where isEnabled:
            shape.stroke(Color.rfGold, lineWidth: 1.5)
        case .destructive where isEnabled:
            shape.stroke(Color.white.opacity(0.18), lineWidth: 1)
        default:
            shape.stroke(Color.clear, lineWidth: 0)
        }
    }

    private var shadowColor: Color {
        guard isEnabled, role == .primary else { return .clear }
        return Color(red: 0.72, green: 0.55, blue: 0.18).opacity(0.38)
    }
}

private struct RFPressStyle: ButtonStyle {
    var enabled: Bool

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(enabled && configuration.isPressed ? 0.96 : 1)
            .animation(.easeOut(duration: 0.12), value: configuration.isPressed)
    }
}

/// Visual glyph stays about 30 pt. The hit target is 44 pt.
struct RFCloseButton: View {
    var label: String = "Close"
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            Image(systemName: "xmark")
                .font(.system(size: 14, weight: .semibold))
                .foregroundStyle(Color.white.opacity(0.75))
                .frame(width: 30, height: 30)
                .background(Circle().fill(Color.white.opacity(0.10)))
                .frame(width: 44, height: 44)
                .contentShape(Rectangle())
        }
        .buttonStyle(.plain)
        .accessibilityLabel(label)
    }
}
