import SwiftUI

// MARK: - Anchors

enum TourAnchor: Hashable {
    case homeStatus
    case homeStreak
    case homeFlags
    case homeSupport
    case groupMembers
    case groupCreate
    case groupInvite
    case settingsPartners
    case settingsRestrictions
    case settingsFeedback
}

struct TourStep: Equatable {
    let tab: Int
    let anchor: TourAnchor
    let title: String
    let body: String
}

enum TourContent {
    /// Bump this when new steps are worth seeing again. A stored value below
    /// this shows the tour once more on this install.
    static let version = 1

    static func shouldAutoPresent(completedVersion: Int, isDemoMode: Bool, forceTour: Bool) -> Bool {
        if forceTour { return true }
        if isDemoMode { return false }
        return completedVersion < version
    }

    /// `inGroup` is true when this phone already has a group to point at.
    /// Without one, the invite step is left out because that button is not on screen.
    static func steps(inGroup: Bool) -> [TourStep] {
        var steps: [TourStep] = [
            TourStep(tab: 0, anchor: .homeStatus, title: "Deep Scan status",
                     body: "This shows if Deep Scan is running. Tap it for the four steps to start it from Control Center."),
            TourStep(tab: 0, anchor: .homeStreak, title: "Your clean streak",
                     body: "Days in a row without a flag. The circles are this week."),
            TourStep(tab: 0, anchor: .homeFlags, title: "Recent flags",
                     body: "Your latest alerts. Tap one to see what happened and a few questions to talk through with your partner."),
            TourStep(tab: 0, anchor: .homeSupport, title: "When it gets hard",
                     body: "One tap sends an urgent notice to your partners and your group."),
        ]
        if inGroup {
            steps.append(TourStep(tab: 1, anchor: .groupMembers, title: "Your group",
                                   body: "Everyone's streak and how they're doing. Groups hold up to 12 people."))
            steps.append(TourStep(tab: 1, anchor: .groupInvite, title: "Invite someone",
                                   body: "Type an email. They join when they sign up with that email."))
        } else {
            steps.append(TourStep(tab: 1, anchor: .groupCreate, title: "Your group",
                                   body: "Start a group here. It takes four short steps."))
        }
        steps.append(contentsOf: [
            TourStep(tab: 2, anchor: .settingsPartners, title: "Your partners",
                     body: "Add partners and star your primary partner."),
            TourStep(tab: 2, anchor: .settingsRestrictions, title: "Choose what to block",
                     body: "Pick the apps that tempt you. Block them or alert your partner after one minute a day."),
            TourStep(tab: 2, anchor: .settingsFeedback, title: "Tell us what you think",
                     body: "Ideas, fixes, problems. We read every one."),
        ])
        return steps
    }
}

enum PauseNotice {
    /// A partner notice goes out only when App Lockout is on and Deep Scan stops.
    static func text(lockoutEnabled: Bool) -> String {
        lockoutEnabled
            ? "Pausing will notify your accountability partners"
            : "Pausing does not notify your partners"
    }
}

extension Notification.Name {
    static let replayAppTour = Notification.Name("replayAppTour")
    static let notificationPermissionResolved = Notification.Name("notificationPermissionResolved")
}

final class TourController: ObservableObject {
    static let shared = TourController()
    @Published var activeAnchor: TourAnchor?
    private init() {}
}

struct TourAnchorKey: PreferenceKey {
    static var defaultValue: [TourAnchor: Anchor<CGRect>] = [:]
    static func reduce(value: inout [TourAnchor: Anchor<CGRect>],
                        nextValue: () -> [TourAnchor: Anchor<CGRect>]) {
        value.merge(nextValue(), uniquingKeysWith: { _, new in new })
    }
}

extension View {
    func tourAnchor(_ id: TourAnchor) -> some View {
        anchorPreference(key: TourAnchorKey.self, value: .bounds) { [id: $0] }
            .id(id)
    }
}

// MARK: - Overlay

struct TourOverlay: View {
    let step: TourStep
    let index: Int
    let count: Int
    let hole: Anchor<CGRect>?
    let proxy: GeometryProxy
    let onBack: () -> Void
    let onNext: () -> Void
    let onSkip: () -> Void

    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @ScaledMetric(relativeTo: .body) private var bodySize: CGFloat = 15

    private var isLast: Bool { index + 1 >= count }
    private var holeRect: CGRect? { hole.map { proxy[$0] } }

    var body: some View {
        let rect = holeRect
        ZStack(alignment: .top) {
            spotlight(rect)
                .accessibilityHidden(true)
            card
                .padding(.top, cardTop(rect))
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .accessibilityElement(children: .contain)
        .accessibilityAddTraits(.isModal)
        .animation(reduceMotion ? .easeOut(duration: 0.2) : .easeInOut(duration: 0.25), value: index)
    }

    private func spotlight(_ rect: CGRect?) -> some View {
        Canvas { context, size in
            var path = Path()
            path.addRect(CGRect(origin: .zero, size: size))
            if let rect {
                let padded = rect.insetBy(dx: -8, dy: -8)
                path.addRoundedRect(in: padded, cornerSize: CGSize(width: 14, height: 14))
            }
            context.fill(path, with: .color(Color.rfNavy.opacity(0.78)), style: FillStyle(eoFill: true))
        }
        .ignoresSafeArea()
        .contentShape(Rectangle())
        .accessibilityHidden(true)
    }

    private var card: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(alignment: .firstTextBaseline, spacing: 12) {
                Text(step.title)
                    .font(.system(.title3, design: .serif).weight(.bold))
                    .foregroundStyle(Color.rfGold)
                    .fixedSize(horizontal: false, vertical: true)
                    .accessibilityAddTraits(.isHeader)
                    .accessibilitySortPriority(6)
                Spacer(minLength: 8)
                Button(action: onSkip) {
                    Text("Skip tour")
                        .font(.body)
                        .foregroundStyle(Color.white.opacity(0.9))
                }
                .accessibilityLabel("Skip tour")
                .accessibilitySortPriority(1)
            }

            Text(step.body)
                .font(.system(size: bodySize))
                .foregroundStyle(.white)
                .fixedSize(horizontal: false, vertical: true)
                .accessibilitySortPriority(5)

            Text("\(index + 1) of \(count)")
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(Color.rfGold.opacity(0.9))
                .accessibilityLabel("Step \(index + 1) of \(count)")
                .accessibilitySortPriority(4)

            HStack {
                if index > 0 {
                    Button("Back", action: onBack)
                        .font(.body)
                        .foregroundStyle(Color.white.opacity(0.9))
                        .accessibilityLabel("Back")
                        .accessibilitySortPriority(2)
                }
                Spacer(minLength: 0)
                Button(action: onNext) {
                    Text(isLast ? "Done" : "Next")
                        .font(.body.weight(.semibold))
                        .foregroundStyle(Color.rfNavy)
                        .padding(.horizontal, 22)
                        .padding(.vertical, 10)
                        .background(Capsule().fill(Color.rfGold))
                }
                .accessibilityLabel(isLast ? "Done" : "Next")
                .accessibilitySortPriority(3)
            }
        }
        .padding(18)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(Color.rfNavy)
                .overlay(
                    RoundedRectangle(cornerRadius: 16)
                        .stroke(Color.rfGold.opacity(0.45), lineWidth: 1)
                )
        )
        .padding(.horizontal, 20)
        .accessibilityElement(children: .contain)
    }

    private func cardTop(_ rect: CGRect?) -> CGFloat {
        let container = proxy.size.height
        let estimated: CGFloat = 240
        guard let rect else { return max(72, (container - estimated) / 2) }
        let below = rect.maxY + 16
        if below + estimated < container - 12 {
            return max(12, below)
        }
        let above = rect.minY - 16 - estimated
        if above > 56 { return above }
        return 64
    }
}
