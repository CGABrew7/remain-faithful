import SwiftUI

// Launch splash for the entry screens.
//
// A thin yellow line draws the shield, a soft bloom opens, and the serif
// wordmark fades in one letter at a time. The mark then eases into the header.
// Draw is 2.4s, ease in and out. Reduce Motion shows the finished mark and fades.
// Rings are a 12 fps canvas, paused when the app is not active, and absent
// when Reduce Motion is on.

enum RFLaunchID {
    static let mark = "rf.launch.mark"
}

struct RFLaunchContext {
    var settled: Bool
    var revealed: Bool
    var reduceMotion: Bool
    var draw: CGFloat
    var glow: CGFloat
    var letters: CGFloat
    var namespace: Namespace.ID
    var linkMark: Bool
}

private struct RFLaunchKey: EnvironmentKey {
    static let defaultValue: RFLaunchContext? = nil
}

extension EnvironmentValues {
    var rfLaunch: RFLaunchContext? {
        get { self[RFLaunchKey.self] }
        set { self[RFLaunchKey.self] = newValue }
    }
}

struct RFLaunchHost<Content: View>: View {
    var enabled: Bool
    var skipsIntro: Bool
    @ViewBuilder var content: () -> Content

    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @Environment(\.scenePhase) private var scenePhase
    @Namespace private var markSpace

    @State private var settled: Bool
    @State private var revealed: Bool
    @State private var draw: CGFloat
    @State private var glow: CGFloat
    @State private var letters: CGFloat
    @State private var elapsed: TimeInterval = 0
    @State private var lastTick: Date?
    @State private var linkMark: Bool

    private let drawDuration: TimeInterval = 2.4
    private let letterLead: TimeInterval = 2.15
    private let letterDuration: TimeInterval = 1.4
    private let holdDuration: TimeInterval = 0.7

    init(enabled: Bool, skipsIntro: Bool, @ViewBuilder content: @escaping () -> Content) {
        self.enabled = enabled
        self.skipsIntro = skipsIntro
        self.content = content
        let done = skipsIntro || !enabled
        _settled = State(initialValue: done)
        _revealed = State(initialValue: done)
        _draw = State(initialValue: done ? 1 : 0)
        _glow = State(initialValue: 0)
        _letters = State(initialValue: done ? 1 : 0)
        _linkMark = State(initialValue: !done)
    }

    private var showSplashMark: Bool {
        enabled && !skipsIntro && !reduceMotion && !settled
    }

    private var context: RFLaunchContext {
        RFLaunchContext(
            settled: skipsIntro || !enabled || reduceMotion || settled,
            revealed: skipsIntro || !enabled || revealed,
            reduceMotion: reduceMotion,
            draw: (reduceMotion || skipsIntro || !enabled) ? 1 : draw,
            glow: glow,
            letters: (reduceMotion || skipsIntro || !enabled) ? 1 : letters,
            namespace: markSpace,
            linkMark: linkMark && enabled && !reduceMotion && !skipsIntro
        )
    }

    var body: some View {
        ZStack {
            content()
                .environment(\.rfLaunch, context)

            if showSplashMark {
                RFLaunchMark(
                    compact: false,
                    draw: draw,
                    glow: glow,
                    letters: letters,
                    namespace: markSpace,
                    linked: linkMark
                )
                .allowsHitTesting(false)
            }

            if showSplashMark && scenePhase == .active {
                TimelineView(.animation(minimumInterval: 1.0 / 30.0, paused: false)) { timeline in
                    Color.clear
                        .onChange(of: timeline.date) { _, date in
                            advance(to: date)
                        }
                }
                .frame(width: 0, height: 0)
                .accessibilityHidden(true)
                .allowsHitTesting(false)
            }
        }
        .animation(.easeInOut(duration: 1.2), value: settled)
        .onChange(of: scenePhase) { _, phase in
            if phase != .active {
                lastTick = nil
            }
        }
        .onAppear {
            guard enabled, !skipsIntro, reduceMotion else { return }
            draw = 1
            glow = 0
            letters = 1
            settled = true
            linkMark = false
            withAnimation(.easeInOut(duration: 0.45)) {
                revealed = true
            }
        }
        .task(id: settled) {
            guard settled, linkMark else { return }
            try? await Task.sleep(nanoseconds: 1_300_000_000)
            if !Task.isCancelled {
                linkMark = false
            }
        }
    }

    private func advance(to date: Date) {
        guard showSplashMark, scenePhase == .active else {
            lastTick = nil
            return
        }
        defer { lastTick = date }
        guard let lastTick else { return }
        let dt = min(0.05, max(0, date.timeIntervalSince(lastTick)))
        guard dt > 0 else { return }
        elapsed += dt
        let linear = min(1, elapsed / drawDuration)
        draw = Self.easeInOut(linear)
        if elapsed > drawDuration {
            let bloom = min(1, (elapsed - drawDuration) / 0.9)
            glow = Self.easeInOut(bloom) * (1 - min(1, max(0, (elapsed - letterLead - 0.4) / 1.2)) * 0.65)
        }
        if elapsed > letterLead {
            letters = Self.easeInOut(min(1, (elapsed - letterLead) / letterDuration))
        }
        if elapsed >= letterLead + letterDuration + holdDuration {
            withAnimation(.easeInOut(duration: 1.2)) {
                settled = true
                revealed = true
                glow = 0.16
                letters = 1
            }
        }
    }

    private static func easeInOut(_ t: CGFloat) -> CGFloat {
        let x = min(1, max(0, t))
        if x < 0.5 {
            return 2 * x * x
        }
        let f = -2 * x + 2
        return 1 - (f * f) / 2
    }
}

// MARK: - Mark

struct RFShieldOutline: Shape {
    func path(in rect: CGRect) -> Path {
        let sx = rect.width / 32
        let sy = rect.height / 36
        func p(_ x: CGFloat, _ y: CGFloat) -> CGPoint {
            CGPoint(x: rect.minX + x * sx, y: rect.minY + y * sy)
        }
        var path = Path()
        path.move(to: p(16, 0))
        path.addLine(to: p(2, 6))
        path.addLine(to: p(2, 18))
        path.addCurve(to: p(16, 36), control1: p(2, 26.284), control2: p(8.268, 33.916))
        path.addCurve(to: p(30, 18), control1: p(23.732, 33.916), control2: p(30, 26.284))
        path.addLine(to: p(30, 6))
        path.closeSubpath()
        return path
    }
}

struct RFShieldCheck: Shape {
    func path(in rect: CGRect) -> Path {
        let sx = rect.width / 32
        let sy = rect.height / 36
        func p(_ x: CGFloat, _ y: CGFloat) -> CGPoint {
            CGPoint(x: rect.minX + x * sx, y: rect.minY + y * sy)
        }
        var path = Path()
        path.move(to: p(11, 18))
        path.addLine(to: p(14.5, 21.5))
        path.addLine(to: p(21, 14))
        return path
    }
}

struct RFDrawnMark: View {
    var draw: CGFloat
    var glow: CGFloat
    var lineWidth: CGFloat = 1.5

    private var shieldTrim: CGFloat { min(1, max(0, draw / 0.78)) }
    private var checkTrim: CGFloat { min(1, max(0, (draw - 0.72) / 0.28)) }

    var body: some View {
        ZStack {
            stroke(RFShieldOutline(), trim: shieldTrim)
            stroke(RFShieldCheck(), trim: checkTrim)
        }
        .accessibilityHidden(true)
    }

    private func stroke<S: Shape>(_ shape: S, trim: CGFloat) -> some View {
        ZStack {
            shape
                .trim(from: 0, to: trim)
                .stroke(Color.rfGold.opacity(0.28 * glow), style: StrokeStyle(lineWidth: lineWidth + 5, lineCap: .round, lineJoin: .round))
            shape
                .trim(from: 0, to: trim)
                .stroke(Color.rfGold, style: StrokeStyle(lineWidth: lineWidth, lineCap: .round, lineJoin: .round))
        }
    }
}

struct RFLaunchMark: View {
    var compact: Bool
    var draw: CGFloat
    var glow: CGFloat
    var letters: CGFloat = 1
    var namespace: Namespace.ID?
    var linked: Bool

    private let word = Array("Remain Faithful")

    var body: some View {
        let lockup = VStack(spacing: compact ? 12 : 26) {
            ZStack {
                Circle()
                    .fill(Color.rfGold.opacity(0.18 * glow))
                    .frame(width: compact ? 64 : 190, height: compact ? 64 : 190)
                    .blur(radius: compact ? 8 : 30)
                RFDrawnMark(draw: max(draw, compact ? 1 : 0), glow: glow, lineWidth: 1.5)
                    .frame(width: compact ? 44 : 92, height: compact ? 50 : 104)
            }
            wordmark
        }

        if linked, let namespace {
            lockup.matchedGeometryEffect(id: RFLaunchID.mark, in: namespace)
        } else {
            lockup
        }
    }

    private var wordmark: some View {
        HStack(spacing: 0) {
            ForEach(Array(word.enumerated()), id: \.offset) { index, character in
                let span = 1 / CGFloat(max(word.count, 1))
                let start = CGFloat(index) * span * 0.84
                let shown = compact ? 1 : min(1, max(0, (letters - start) / max(span, 0.08)))
                Text(String(character))
                    .opacity(Double(shown))
            }
        }
        .font(compact ? .largeTitle.bold() : .system(size: 40, weight: .bold, design: .serif))
        .fontDesign(.serif)
        .foregroundStyle(.white)
        .minimumScaleFactor(0.7)
        .accessibilityElement(children: .ignore)
        .accessibilityLabel("Remain Faithful")
        .accessibilityAddTraits(.isHeader)
    }
}

struct RFRise: ViewModifier {
    var shown: Bool
    var index: Int

    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    func body(content: Content) -> some View {
        content
            .opacity(shown ? 1 : 0)
            .offset(y: shown || reduceMotion ? 0 : 18)
            .allowsHitTesting(shown)
            .accessibilityHidden(!shown)
            .animation(
                reduceMotion
                    ? .easeInOut(duration: 0.4)
                    : .easeInOut(duration: 0.7).delay(0.08 * Double(index)),
                value: shown
            )
    }
}

extension View {
    func rfRise(shown: Bool, index: Int) -> some View {
        modifier(RFRise(shown: shown, index: index))
    }
}

struct RFAmbientRings: View {
    @Environment(\.scenePhase) private var scenePhase
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        if reduceMotion {
            Color.clear
                .allowsHitTesting(false)
                .accessibilityHidden(true)
        } else {
            TimelineView(.animation(minimumInterval: 1.0 / 12.0, paused: scenePhase != .active)) { timeline in
                Canvas { context, size in
                    let t = timeline.date.timeIntervalSinceReferenceDate
                    let center = CGPoint(x: size.width / 2, y: size.height * 0.42)
                    let maxR = max(size.width, size.height) * 0.62
                    let count = 3
                    for i in 0..<count {
                        let shift = Double(i) / Double(count)
                        let breathe = 0.5 + 0.5 * sin(t * 0.32 + shift * .pi * 2)
                        let radius = max(1, maxR * (0.34 + 0.22 * Double(i)) * (0.94 + 0.07 * breathe))
                        let alpha = 0.035 + 0.045 * breathe
                        var path = Path()
                        path.addEllipse(in: CGRect(
                            x: center.x - radius,
                            y: center.y - radius,
                            width: radius * 2,
                            height: radius * 2
                        ))
                        context.stroke(
                            path,
                            with: .color(Color.rfGold.opacity(alpha)),
                            lineWidth: 1
                        )
                    }
                }
            }
            .allowsHitTesting(false)
            .accessibilityHidden(true)
        }
    }
}
