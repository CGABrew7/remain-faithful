import Foundation

/// Optional on-phone wording. The partner sentence is never rewritten.
enum CoachWarmth {
    static func checkInBody() async -> String {
        guard CoachStore.writingHelperOn else { return CoachCopy.checkInBody }
        #if canImport(FoundationModels)
        if #available(iOS 26.0, *) {
            return await OnPhoneLine.rewrite(CoachCopy.checkInBody)
        }
        #endif
        return CoachCopy.checkInBody
    }
}

#if canImport(FoundationModels)
import FoundationModels

@available(iOS 26.0, *)
private enum OnPhoneLine {
    static func rewrite(_ fallback: String) async -> String {
        let model = SystemLanguageModel.default
        guard case .available = model.availability else { return fallback }
        let session = LanguageModelSession(instructions: """
        Rewrite a short check-in invitation for the person using this phone. \
        Keep it warm and under 140 characters. Do not mention apps, places, or what they were doing.
        """)
        do {
            let response = try await session.respond(to: "Rewrite this invitation: \(fallback)")
            let text = response.content.trimmingCharacters(in: .whitespacesAndNewlines)
            if text.isEmpty || text.count > 180 { return fallback }
            if text.contains("\u{2014}") || text.contains("\u{2013}") { return fallback }
            return text
        } catch {
            return fallback
        }
    }
}
#endif
