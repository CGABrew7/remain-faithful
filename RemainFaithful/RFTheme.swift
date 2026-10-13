import SwiftUI

// Shared visual language for Remain Faithful.
//
// The navy and gold brand palette is unchanged — rfNavy, rfNavyMid and rfGold
// keep their exact values. Everything here is derived from those three colours
// so surfaces, borders and controls read as one system.
//
// Depth comes from three cues used consistently:
//   1. A lit top edge   — a hairline that is brighter at the top of a shape,
//                         as if a soft light sits above the screen.
//   2. Layered shadow   — a tight contact shadow plus a wide ambient one,
//                         instead of a single blurry drop shadow.
//   3. A surface tint   — a very low-alpha white fill that lifts a card off
//                         the navy background without introducing a new hue.
//
// Contrast: every foreground token below was measured against its intended
// background and clears WCAG AA for normal text (4.5:1). See RFContrast.

// MARK: - Palette

extension Color {
    /// Page background, top of the gradient. Brand navy — do not change.
    static let rfNavy    = Color(red: 0.07, green: 0.13, blue: 0.30)
    /// Page background, bottom of the gradient. Brand navy — do not change.
    static let rfNavyMid = Color(red: 0.11, green: 0.19, blue: 0.40)
    /// Brand gold — do not change.
    static let rfGold    = Color(red: 0.82, green: 0.67, blue: 0.30)

    // Derived gold stops for the primary gradient. Navy text measures
    // 10.9:1 on the top stop and 5.1:1 on the bottom stop.
    static let rfGoldLight = Color(red: 0.95, green: 0.84, blue: 0.52)
    static let rfGoldDeep  = Color(red: 0.72, green: 0.55, blue: 0.18)

    /// Destructive fill. White text measures 7.7:1.
    static let rfDanger = Color(red: 0.62, green: 0.14, blue: 0.16)

    // Foreground steps. Measured on the standard card surface.
    static let rfTextPrimary   = Color.white                  // 13.1:1
    static let rfTextSecondary = Color.white.opacity(0.72)    //  7.3:1
    static let rfTextTertiary  = Color.white.opacity(0.60)    //  5.4:1
}

// MARK: - Metrics

enum RFMetrics {
    /// Corner radius for cards and grouped sections.
    static let cardRadius: CGFloat = 16
    /// Corner radius for controls that sit inside a card.
    static let controlRadius: CGFloat = 14
    /// Minimum hit target on every interactive control.
    static let minTap: CGFloat = 44

    // Spacing scale. Using a single scale keeps rhythm consistent across screens.
    static let gutter: CGFloat    = 20   // screen side margins
    static let sectionGap: CGFloat = 24  // between major sections
    static let cardPadding: CGFloat = 18 // inside a card
    static let rowGap: CGFloat    = 14   // between elements in a row
    static let labelGap: CGFloat  = 10   // header to content
}

// MARK: - Surface

/// The standard elevated surface: tinted fill, lit top edge, layered shadow.
struct RFSurface: View {
    var radius: CGFloat = RFMetrics.cardRadius
    /// Raised surfaces sit above the page; inset ones sit within another card.
    var raised: Bool = true

    var body: some View {
        let shape = RoundedRectangle(cornerRadius: radius, style: .continuous)
        shape
            .fill(Color.white.opacity(raised ? 0.06 : 0.04))
            .overlay(shape.strokeBorder(RFStroke.lit, lineWidth: 1))
            .compositingGroup()
            .shadow(color: .black.opacity(raised ? 0.22 : 0), radius: 2,  x: 0, y: 1)
            .shadow(color: .black.opacity(raised ? 0.16 : 0), radius: 14, x: 0, y: 8)
    }
}

enum RFStroke {
    /// Hairline that is brighter along the top edge, darker along the bottom.
    /// This is what makes a flat rectangle read as a physical surface.
    static let lit = LinearGradient(
        colors: [Color.white.opacity(0.16), Color.white.opacity(0.05)],
        startPoint: .top,
        endPoint: .bottom
    )

    /// Gold hairline for secondary buttons and selected states.
    static let gold = LinearGradient(
        colors: [Color.rfGold.opacity(0.95), Color.rfGold.opacity(0.55)],
        startPoint: .top,
        endPoint: .bottom
    )
}

extension View {
    /// Wraps content in the standard card surface with consistent padding.
    func rfCard(padding: CGFloat = RFMetrics.cardPadding,
                radius: CGFloat = RFMetrics.cardRadius,
                raised: Bool = true) -> some View {
        self
            .padding(padding)
            .background(RFSurface(radius: radius, raised: raised))
            .clipShape(RoundedRectangle(cornerRadius: radius, style: .continuous))
    }

    /// Card surface without internal padding, for grouped rows that draw
    /// their own dividers edge to edge.
    func rfGroupedCard(radius: CGFloat = RFMetrics.cardRadius) -> some View {
        self
            .background(RFSurface(radius: radius))
            .clipShape(RoundedRectangle(cornerRadius: radius, style: .continuous))
    }
}

// MARK: - Section header

/// Small gold caps header above a grouped card.
/// Scales with Dynamic Type and is exposed to VoiceOver as a header.
struct RFSectionHeader: View {
    let title: String

    var body: some View {
        Text(title)
            .font(.caption.weight(.bold))
            .kerning(1.2)
            .foregroundStyle(Color.rfGold.opacity(0.9))
            .padding(.leading, 4)
            .accessibilityAddTraits(.isHeader)
    }
}

// MARK: - Divider

/// Hairline divider between rows inside a grouped card.
struct RFRowDivider: View {
    var inset: CGFloat = 0

    var body: some View {
        Rectangle()
            .fill(Color.white.opacity(0.08))
            .frame(height: 1)
            .padding(.leading, inset)
            .accessibilityHidden(true)
    }
}

// MARK: - Toggle

/// Gold toggle that matches the button language.
/// Keeps the native Toggle semantics, so VoiceOver still reads it as a switch
/// and the 44 pt target is preserved.
///
/// `showsLabel: false` renders only the switch, for rows that already draw
/// their own title and detail text. Those call sites supply an explicit
/// `.accessibilityLabel`, so nothing is lost to VoiceOver.
struct RFToggleStyle: ToggleStyle {
    var showsLabel: Bool = true

    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @Environment(\.isEnabled) private var isEnabled

    func makeBody(configuration: Configuration) -> some View {
        let on = configuration.isOn
        return HStack(spacing: 0) {
            if showsLabel {
                configuration.label
                Spacer(minLength: RFMetrics.rowGap)
            }
            ZStack(alignment: on ? .trailing : .leading) {
                Capsule()
                    .fill(on
                          ? AnyShapeStyle(LinearGradient(colors: [.rfGoldLight, .rfGoldDeep],
                                                         startPoint: .top, endPoint: .bottom))
                          : AnyShapeStyle(Color.white.opacity(0.16)))
                    .overlay(Capsule().strokeBorder(
                        on ? Color.rfGoldDeep.opacity(0.55) : Color.white.opacity(0.10),
                        lineWidth: 1))
                    .frame(width: 51, height: 31)

                Circle()
                    .fill(.white)
                    .frame(width: 27, height: 27)
                    .shadow(color: .black.opacity(0.28), radius: 2, x: 0, y: 1)
                    .padding(.horizontal, 2)
            }
            .opacity(isEnabled ? 1 : 0.5)
            .frame(minWidth: RFMetrics.minTap, minHeight: RFMetrics.minTap, alignment: .trailing)
            .contentShape(Rectangle())
            .onTapGesture {
                guard isEnabled else { return }
                UIImpactFeedbackGenerator(style: .light).impactOccurred()
                if reduceMotion {
                    configuration.isOn.toggle()
                } else {
                    withAnimation(.spring(response: 0.26, dampingFraction: 0.72)) {
                        configuration.isOn.toggle()
                    }
                }
            }
        }
    }
}

// MARK: - Icon badge

/// Rounded tinted square behind a glyph. Scales with Dynamic Type.
struct RFIconBadge: View {
    let systemName: String
    var tint: Color = .rfGold
    @ScaledMetric(relativeTo: .body) private var side: CGFloat = 36

    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: side * 0.28, style: .continuous)
                .fill(tint.opacity(0.16))
                .overlay(
                    RoundedRectangle(cornerRadius: side * 0.28, style: .continuous)
                        .strokeBorder(tint.opacity(0.22), lineWidth: 1)
                )
            Image(systemName: systemName)
                .font(.system(size: side * 0.46, weight: .semibold))
                .foregroundStyle(tint)
        }
        .frame(width: side, height: side)
        .accessibilityHidden(true)
    }
}
