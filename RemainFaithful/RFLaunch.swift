import SwiftUI

// Launch splash for the entry screens.
//
// A thin yellow line draws the shield mark, holds, glows, then the mark
// eases into the header while the first screen rises. The draw eases in and
// out over 1.55s. Reduce Motion shows the finished mark and fades the screen.
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
    @State private var elapsed: TimeInterval = 0
    @State private var lastTick: Date?
    @State private var linkMark: Bool

    private let drawDuration: TimeInterval = 1.55
    private let holdDuration: TimeInterval = 0.45

    init(enabled: Bool, skipsIntro: Bool, @ViewBuilder content: @escaping () -> Content) {
        self.enabled = enabled
        self.skipsIntro = skipsIntro
        self.content = content
        let done = skipsIntro || !enabled
        _settled = State(initialValue: done)
        _revealed = State(initialValue: done)
        _draw = State(initialValue: done ? 1 : 0)
        _glow = State(initialValue: 0)
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
        .animation(.easeInOut(duration: 0.85), value: settled)
        .onChange(of: scenePhase) { _, phase in
            if phase != .active {
                lastTick = nil
            }
        }
        .onAppear {
            guard enabled, !skipsIntro, reduceMotion else { return }
            draw = 1
            glow = 0
            settled = true
            linkMark = false
            withAnimation(.easeInOut(duration: 0.45)) {
                revealed = true
            }
        }
        .task(id: settled) {
            guard settled, linkMark else { return }
            try? await Task.sleep(nanoseconds: 950_000_000)
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
            let held = min(1, (elapsed - drawDuration) / 0.4)
            glow = Self.easeInOut(held)
        }
        if elapsed >= drawDuration + holdDuration {
            withAnimation(.easeInOut(duration: 0.85)) {
                settled = true
                revealed = true
                glow = 0.28
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
    var namespace: Namespace.ID?
    var linked: Bool

    private var titleOpacity: Double {
        if compact { return 1 }
        return Double(min(1, max(0, (draw - 0.62) / 0.28)))
    }

    var body: some View {
        let lockup = VStack(spacing: compact ? 12 : 22) {
            RFDrawnMark(draw: max(draw, compact ? 1 : 0), glow: glow, lineWidth: 1.5)
                .frame(width: compact ? 44 : 92, height: compact ? 50 : 104)
            Text("Remain Faithful")
                .font(compact ? .largeTitle.bold() : .system(size: 40, weight: .bold, design: .serif))
                .fontDesign(.serif)
                .foregroundStyle(.white)
                .multilineTextAlignment(.center)
                .minimumScaleFactor(0.7)
                .opacity(titleOpacity)
                .accessibilityAddTraits(.isHeader)
        }

        if linked, let namespace {
            lockup.matchedGeometryEffect(id: RFLaunchID.mark, in: namespace)
        } else {
            lockup
        }
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
                    let maxR = max(size.width, size.height) * 0.78
                    let period = 18.0
                    let count = 3
                    for i in 0..<count {
                        let shift = Double(i) / Double(count)
                        let u = (t / period + shift).truncatingRemainder(dividingBy: 1)
                        let radius = max(1, maxR * u)
                        let alpha = (1 - u) * 0.10
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
