package com.remainfaithful.core

/**
 * Event metadata that is allowed to leave the phone.
 * Summaries match the Go API allowlist. App names, package names, and URLs
 * are never part of a payload.
 */
object EventPolicy {
    const val CATEGORY_APP = "app_usage"
    const val SEVERITY_CONCERNING = "concerning"
    const val SUMMARY_MINUTE = "Monitored app used for 1+ minute"
    const val SUMMARY_BLOCKED = "Blocked app opened"

    private val dash = "\u2014"

    private val categories = setOf(
        "adult_content", "gambling", "violence", "self_harm", "clean", "app_usage",
    )
    private val severities = setOf("informational", "concerning", "severe")

    private val allowed = setOf(
        "Explicit content detected",
        "Gambling content detected",
        "Violent content detected",
        "Self-harm content detected",
        "Content reviewed $dash no concerns",
        "Explicit image detected (perceptual hash match)",
        "Blocked domain detected",
        "Explicit keyword detected in screen text",
        SUMMARY_MINUTE,
        "Monitored app category used for 1+ minute",
        "Monitored app activity detected",
        "App activity detected",
        SUMMARY_BLOCKED,
        "Adult content flagged on screen",
        "Gambling flagged on screen",
        "Violence flagged on screen",
        "Self-harm content flagged on screen",
        "Content reviewed",
    )

    private val continued = Regex(
        "^(Explicit content|Gambling content|Violent content|Self-harm content|Activity) $dash detected [1-9][0-9]{0,3} times? in 5 min$",
    )

    fun isAllowedSummary(summary: String): Boolean {
        if (summary.isEmpty() || summary.length > 120) return false
        if (summary in allowed) return true
        return continued.matches(summary)
    }

    fun leaks(text: String): Boolean {
        val lower = text.lowercase()
        return "http://" in lower || "https://" in lower || "www." in lower || "://" in lower
    }

    fun metadata(
        category: String,
        severity: String,
        summary: String,
        timestamp: String? = null,
    ): Map<String, String> {
        require(category in categories) { "invalid category" }
        require(severity in severities) { "invalid severity" }
        require(isAllowedSummary(summary)) { "invalid summary" }
        require(!leaks(summary)) { "summary leaks" }
        val body = linkedMapOf(
            "category" to category,
            "severity" to severity,
            "summary" to summary,
        )
        if (!timestamp.isNullOrBlank()) body["timestamp"] = timestamp
        return body
    }

    fun minuteEvent(): Map<String, String> =
        metadata(CATEGORY_APP, SEVERITY_CONCERNING, SUMMARY_MINUTE)

    fun blockedEvent(): Map<String, String> =
        metadata(CATEGORY_APP, SEVERITY_CONCERNING, SUMMARY_BLOCKED)
}

data class MonitorInput(
    val selected: Set<String>,
    val alertsOn: Boolean,
    val blockingOn: Boolean,
    val lockoutOn: Boolean,
    val deepScanActive: Boolean,
    val alreadyAlertedToday: Set<String>,
    val foregroundPackage: String?,
    val foregroundMsToday: Map<String, Long>,
    val openedSinceLastCheck: Set<String>,
    val ownPackage: String,
)

data class MonitorResult(
    val uploads: List<Map<String, String>>,
    val alertedPackages: Set<String>,
    val showShield: Boolean,
)

object MonitoringPolicy {
    const val minuteMs = 60_000L

    fun decide(input: MonitorInput): MonitorResult {
        val uploads = mutableListOf<Map<String, String>>()
        val alerted = mutableSetOf<String>()
        if (input.alertsOn) {
            for ((pkg, ms) in input.foregroundMsToday) {
                if (pkg.isBlank() || pkg == input.ownPackage || pkg !in input.selected) continue
                if (pkg in input.alreadyAlertedToday) continue
                if (ms >= minuteMs) {
                    uploads += EventPolicy.minuteEvent()
                    alerted += pkg
                }
            }
        }
        if (input.blockingOn) {
            for (pkg in input.openedSinceLastCheck) {
                if (pkg.isBlank() || pkg == input.ownPackage || pkg !in input.selected) continue
                uploads += EventPolicy.blockedEvent()
            }
        }
        val watched = input.foregroundPackage != null &&
            input.foregroundPackage != input.ownPackage &&
            input.foregroundPackage in input.selected
        val shield = input.blockingOn && watched && (!input.lockoutOn || !input.deepScanActive)
        return MonitorResult(uploads, alerted, shield)
    }
}

object HeartbeatPolicy {
    const val idleAfterMs = 120_000L

    fun screen(
        alertsOn: Boolean,
        blockingOn: Boolean,
        deepScanRunning: Boolean,
        millisSinceFrame: Long?,
    ): String? {
        if (deepScanRunning) {
            return if (millisSinceFrame != null && millisSinceFrame <= idleAfterMs) "active" else "idle"
        }
        if (alertsOn || blockingOn) return "monitoring"
        return null
    }

    fun body(screen: String): Map<String, String> {
        require(screen == "active" || screen == "idle" || screen == "monitoring")
        return mapOf("screen" to screen)
    }
}

object GroupPolicy {
    const val maxMembers = 12

    fun canInvite(memberCount: Int): Boolean = memberCount < maxMembers

    fun fullMessage(memberCount: Int): String? =
        if (memberCount >= maxMembers) "Groups hold up to 12 people." else null
}

object AuthRules {
    fun passwordOk(password: String): Boolean = password.length >= 8

    fun canRegister(adult: Boolean, name: String, email: String, password: String): Boolean =
        adult && name.isNotBlank() && email.contains("@") && email.contains(".") && passwordOk(password)

    fun passwordHint(password: String): String = when {
        password.isEmpty() -> "At least 8 characters."
        passwordOk(password) -> "8 or more characters."
        else -> "${password.length} of 8 characters."
    }
}

data class PinCheck(val success: Boolean, val reason: String?) {
    val message: String
        get() = when {
            success -> ""
            reason == "rate_limited" -> "Too many PIN attempts. Try again in 15 minutes."
            reason == "no_pin" -> "No Partner PIN is set."
            reason == "wrong_pin" -> "Incorrect PIN. Your partner is the one who knows it."
            else -> "Could not check the PIN."
        }

    companion object {
        fun from(success: Boolean, reason: String?): PinCheck = PinCheck(success, reason)
    }
}

object FeedbackRules {
    val types = listOf("Idea", "Improvement", "Problem")

    fun error(type: String, summary: String, details: String): String? {
        if (type !in types) return "type must be Idea, Improvement, or Problem"
        if (summary.isBlank() || summary.length > 120) return "summary is required (120 characters max)"
        if (details.isBlank() || details.length > 4000) return "details are required (4000 characters max)"
        return null
    }
}

object Payloads {
    val protectionTypes = setOf(
        "monitoring_disabled",
        "shielding_disabled",
        "lockout_disabled",
        "deep_scan_stopped",
        "wrong_pin_attempt",
        "pin_removed",
        "pin_changed",
        "family_controls_revoked",
        "signed_out",
        "selection_cleared",
    )

    fun deviceToken(token: String): Map<String, String> = mapOf(
        "token" to token,
        "platform" to "android",
        "environment" to "production",
    )

    fun signals(
        usageGranted: Boolean,
        monitoredCount: Int,
        appVersion: String,
        supportTaps: Int,
    ): Map<String, Any> = linkedMapOf(
        "screen_time_granted" to usageGranted,
        "monitored_app_count" to monitoredCount,
        "app_version" to appVersion,
        "support_taps" to supportTaps,
    )

    fun protection(type: String, detail: String): Map<String, String> {
        require(type in protectionTypes)
        require(!EventPolicy.leaks(detail))
        return mapOf("type" to type, "detail" to detail)
    }

    fun pin(pin: String): Map<String, String> {
        require(pin.length == 4 && pin.all { it.isDigit() })
        return mapOf("pin" to pin)
    }
}

object JwtClock {
    fun expiresAtEpochSeconds(token: String): Long? {
        val parts = token.split('.')
        if (parts.size < 2) return null
        val payload = try {
            java.util.Base64.getUrlDecoder().decode(pad(parts[1])).decodeToString()
        } catch (_: IllegalArgumentException) {
            return null
        }
        val match = Regex("\"exp\"\\s*:\\s*(\\d+)").find(payload) ?: return null
        return match.groupValues[1].toLongOrNull()
    }

    fun shouldRefresh(token: String, nowEpochSeconds: Long): Boolean {
        val exp = expiresAtEpochSeconds(token) ?: return false
        return exp - nowEpochSeconds < 3600
    }

    private fun pad(value: String): String {
        val mod = value.length % 4
        return if (mod == 0) value else value + "=".repeat(4 - mod)
    }
}

enum class StreakMark { Clean, Flagged, BeforeSignup }

object StreakWeek {
    fun mark(raw: String): StreakMark = when (raw.lowercase()) {
        "clean", "true", "1" -> StreakMark.Clean
        "before_signup", "before-signup", "before_account", "not_applicable", "na", "none" ->
            StreakMark.BeforeSignup
        else -> StreakMark.Flagged
    }

    fun marks(raw: List<String>): List<StreakMark> = raw.map { mark(it) }
}

enum class MonitoringKind { NeedsPermission, NoApps, Off, On }

data class MonitoringPresentation(
    val kind: MonitoringKind,
    val title: String,
    val detail: String,
    val deepScanLine: String,
    val actionTitle: String?,
    val coreIsOn: Boolean,
)

object MonitoringStatus {
    fun make(
        usageGranted: Boolean,
        selectedCount: Int,
        alertsOn: Boolean,
        blockingOn: Boolean,
        deepScanOn: Boolean,
    ): MonitoringPresentation {
        val deep = if (deepScanOn) {
            "Deep Scan is on. Protected video such as Netflix is not visible."
        } else {
            "Deep Scan is off. It is optional and separate from Alerts and Blocking."
        }
        if (!usageGranted) {
            return MonitoringPresentation(
                MonitoringKind.NeedsPermission,
                "Setup needed",
                "Usage access is off. Alerts and Blocking stay off until you allow it.",
                deep,
                "Allow usage access",
                false,
            )
        }
        if (selectedCount == 0) {
            return MonitoringPresentation(
                MonitoringKind.NoApps,
                "No apps selected",
                "Nothing is being watched yet. Choose the apps you want help with. You can take your time.",
                deep,
                "Choose apps",
                false,
            )
        }
        if (!alertsOn && !blockingOn) {
            return MonitoringPresentation(
                MonitoringKind.Off,
                "Monitoring off",
                "Apps are chosen. Alerts and Blocking are both off. Turn them on when you are ready.",
                deep,
                "Turn monitoring on",
                false,
            )
        }
        val parts = listOf(
            if (alertsOn) "Alerts on" else "Alerts off",
            if (blockingOn) "Blocking on" else "Blocking off",
        )
        return MonitoringPresentation(
            MonitoringKind.On,
            "Monitoring on",
            parts.joinToString(". ") + ". You are not doing this alone.",
            deep,
            null,
            true,
        )
    }
}

fun friendlyServerMessage(raw: String?): String {
    val cleaned = raw.orEmpty()
        .replace("\u2014", ",")
        .replace("\u2013", "-")
        .trim()
    if (cleaned.contains("maximum of 12")) return "Groups hold up to 12 people."
    if (cleaned.isEmpty()) return "Something went wrong. Try again."
    return cleaned
}

data class TourStep(val tab: Int, val title: String, val body: String)

object TourContent {
    const val version = 1

    fun shouldAutoPresent(completedVersion: Int, force: Boolean): Boolean {
        if (force) return true
        return completedVersion < version
    }

    fun steps(inGroup: Boolean): List<TourStep> {
        val steps = mutableListOf(
            TourStep(0, "Monitoring status", "This shows whether Alerts and Blocking are on. Deep Scan is the smaller line. It is optional."),
            TourStep(0, "Your clean streak", "Days in a row without a flag. The circles are this week."),
            TourStep(0, "Recent flags", "Your latest alerts. Tap one to see what happened and a few questions to talk through with your partner."),
            TourStep(0, "When it gets hard", "One tap sends an urgent notice to your partners and your group."),
        )
        if (inGroup) {
            steps += TourStep(1, "Your group", "Everyone's streak and how they're doing. Groups hold up to 12 people.")
            steps += TourStep(1, "Invite someone", "Type an email. They join when they accept, or when they sign up with that email.")
        } else {
            steps += TourStep(1, "Your group", "Start a group here. A name is enough. Groups hold up to 12 people.")
        }
        steps += TourStep(2, "Your partners", "Add partners and star your primary partner.")
        steps += TourStep(2, "Choose apps", "Pick the apps that tempt you. Blocking covers them with a shield you asked for. Alerts tell your partner after one minute a day.")
        steps += TourStep(2, "Tell us what you think", "Ideas, fixes, problems. We read every one.")
        return steps
    }
}
