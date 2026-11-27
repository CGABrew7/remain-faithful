package com.remainfaithful.core

enum class CoachReply { Pending, NotNow, CheckIn }

data class CoachClock(
    val year: Int,
    val month: Int,
    val day: Int,
    val hour: Int,
    val minute: Int,
    val weekday: Int,
    val epochMillis: Long,
)

data class CoachSignals(
    val enabled: Boolean,
    val partnerNudgeEnabled: Boolean,
    val imAway: Boolean,
    val lateStartHour: Int,
    val lateEndHour: Int,
    val missedCheckIn: Boolean,
    val shieldOpensToday: Int,
    val thresholdFresh: Boolean,
    val recentSession: Boolean,
    val lastPromptAt: Long?,
    val lastNudge: CoachClock?,
)

data class CoachDecision(
    val showPrompt: Boolean,
    val sendPartnerNudge: Boolean,
)

object CoachPolicy {
    const val quietMillis = 6L * 60 * 60 * 1000
    const val shieldOpensForPrompt = 2
    const val defaultLateStart = 23
    const val defaultLateEnd = 5
    const val alertType = "coach_nudge"
    const val notificationType = "COACH_CHECK_IN"

    fun evaluate(clock: CoachClock, signals: CoachSignals, reply: CoachReply): CoachDecision {
        val open = windowOpen(clock, signals)
        val quiet = inQuietPeriod(clock, signals)
        return CoachDecision(
            showPrompt = open && !quiet && reply == CoachReply.Pending,
            sendPartnerNudge = open &&
                reply == CoachReply.CheckIn &&
                signals.partnerNudgeEnabled &&
                !nudgedToday(clock, signals),
        )
    }

    fun windowOpen(clock: CoachClock, signals: CoachSignals): Boolean {
        if (!signals.enabled) return false
        if (lateNight(clock.hour, signals.lateStartHour, signals.lateEndHour)) return true
        if (clock.weekday == 1 && clock.hour in 17..22) return true
        if (signals.missedCheckIn && clock.hour >= 18) return true
        if (signals.shieldOpensToday >= shieldOpensForPrompt) return true
        if (signals.imAway && clock.hour in 8..21) return true
        if (signals.thresholdFresh || signals.recentSession) return true
        return false
    }

    fun lateNight(hour: Int, start: Int, end: Int): Boolean {
        val startHour = normalize(start)
        val endHour = normalize(end)
        val now = normalize(hour)
        if (startHour == endHour) return false
        return if (startHour < endHour) now in startHour until endHour else now >= startHour || now < endHour
    }

    fun inQuietPeriod(clock: CoachClock, signals: CoachSignals): Boolean {
        val last = signals.lastPromptAt ?: return false
        return clock.epochMillis - last < quietMillis
    }

    fun nudgedToday(clock: CoachClock, signals: CoachSignals): Boolean {
        val last = signals.lastNudge ?: return false
        return last.year == clock.year && last.month == clock.month && last.day == clock.day
    }

    fun hourLabel(hour: Int): String {
        val normalized = normalize(hour)
        val suffix = if (normalized < 12) "AM" else "PM"
        val hour12 = if (normalized % 12 == 0) 12 else normalized % 12
        return "$hour12 $suffix"
    }

    fun dayKey(clock: CoachClock): String = "%04d-%02d-%02d".format(clock.year, clock.month, clock.day)

    fun clock(now: java.time.ZonedDateTime): CoachClock = CoachClock(
        year = now.year,
        month = now.monthValue,
        day = now.dayOfMonth,
        hour = now.hour,
        minute = now.minute,
        weekday = now.dayOfWeek.value % 7 + 1,
        epochMillis = now.toInstant().toEpochMilli(),
    )

    private fun normalize(hour: Int): Int {
        val mod = hour % 24
        return if (mod >= 0) mod else mod + 24
    }
}

object CoachKeys {
    val all = listOf(
        "coachEnabled",
        "coachPartnerNudgeEnabled",
        "coachLateNightStart",
        "coachLateNightEnd",
        "imAway",
        "coachLastLocalPromptAt",
        "coachLastPartnerNudgeAt",
        "coachThresholdReachedAt",
        "blockedOpenCountToday",
        "blockedOpenDay",
        "coachCheckedInDay",
        "coachOwnerUserID",
        "coachPendingNudgeUser",
    )

    fun wipe(raw: MutableMap<String, Any?>) {
        all.forEach { raw.remove(it) }
    }
}
