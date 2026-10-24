package com.remainfaithful.core

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test

class CoreLogicTest {
    @Test
    fun uploadsAreMetadataOnly() {
        val input = MonitorInput(
            selected = setOf("com.example.browser", "com.example.game"),
            alertsOn = true,
            blockingOn = true,
            lockoutOn = true,
            deepScanActive = false,
            alreadyAlertedToday = setOf("com.example.game"),
            foregroundPackage = "com.example.browser",
            foregroundMsToday = mapOf(
                "com.example.browser" to 90_000L,
                "com.example.game" to 120_000L,
                "com.remainfaithful.app" to 500_000L,
            ),
            openedSinceLastCheck = setOf("com.example.browser", "com.remainfaithful.app"),
            ownPackage = "com.remainfaithful.app",
        )
        val result = MonitoringPolicy.decide(input)
        assertEquals(2, result.uploads.size)
        assertEquals(setOf("com.example.browser"), result.alertedPackages)
        assertTrue(result.showShield)
        result.uploads.forEach { payload ->
            assertEquals(setOf("category", "severity", "summary"), payload.keys)
            val blob = payload.values.joinToString(" ")
            assertFalse(blob.contains("com.example"))
            assertFalse(blob.contains("http"))
            assertFalse(EventPolicy.leaks(blob))
            assertTrue(EventPolicy.isAllowedSummary(payload.getValue("summary")))
        }
        assertEquals(EventPolicy.SUMMARY_MINUTE, result.uploads[0].getValue("summary"))
        assertEquals(EventPolicy.SUMMARY_BLOCKED, result.uploads[1].getValue("summary"))
    }

    @Test
    fun deepScanLiftsLockoutShield() {
        val shielded = baseInput().copy(blockingOn = true, lockoutOn = true, deepScanActive = false)
        val lifted = baseInput().copy(blockingOn = true, lockoutOn = true, deepScanActive = true)
        assertTrue(MonitoringPolicy.decide(shielded).showShield)
        assertFalse(MonitoringPolicy.decide(lifted).showShield)
    }

    @Test
    fun minuteAlertFiresOncePerPackage() {
        val first = MonitoringPolicy.decide(baseInput())
        assertEquals(1, first.uploads.size)
        val second = MonitoringPolicy.decide(
            baseInput().copy(alreadyAlertedToday = first.alertedPackages, openedSinceLastCheck = emptySet()),
        )
        assertTrue(second.uploads.isEmpty())
    }

    @Test
    fun heartbeatMatchesServerValues() {
        assertEquals("monitoring", HeartbeatPolicy.screen(true, false, false, null))
        assertEquals("active", HeartbeatPolicy.screen(true, true, true, 10_000))
        assertEquals("idle", HeartbeatPolicy.screen(true, true, true, 121_000))
        assertNull(HeartbeatPolicy.screen(false, false, false, null))
        assertEquals(setOf("screen"), HeartbeatPolicy.body("monitoring").keys)
    }

    @Test
    fun groupCapIsTwelve() {
        assertTrue(GroupPolicy.canInvite(11))
        assertFalse(GroupPolicy.canInvite(12))
        assertEquals("Groups hold up to 12 people.", GroupPolicy.fullMessage(12))
        assertEquals("Groups hold up to 12 people.", friendlyServerMessage("group has reached the maximum of 12 members"))
    }

    @Test
    fun ageGateAndPassword() {
        assertFalse(AuthRules.canRegister(false, "Ada", "ada@example.com", "longpassword"))
        assertFalse(AuthRules.canRegister(true, "Ada", "ada@example.com", "short"))
        assertTrue(AuthRules.canRegister(true, "Ada", "ada@example.com", "longpassword"))
        assertEquals("At least 8 characters.", AuthRules.passwordHint(""))
    }

    @Test
    fun pinMessages() {
        assertEquals("Incorrect PIN. Your partner is the one who knows it.", PinCheck(false, "wrong_pin").message)
        assertEquals("Too many PIN attempts. Try again in 15 minutes.", PinCheck(false, "rate_limited").message)
        assertEquals(mapOf("pin" to "1234"), Payloads.pin("1234"))
    }

    @Test
    fun googleSignInSendsOnlyTheIdToken() {
        val body = Payloads.googleIdToken("header.payload.sig")
        assertEquals(mapOf("id_token" to "header.payload.sig"), body)
    }

    fun deviceTokenAndSignalsDoNotCarryAppNames() {
        val token = Payloads.deviceToken("fcm-token")
        assertEquals(setOf("token", "platform", "environment"), token.keys)
        assertEquals("android", token.getValue("platform"))
        val signals = Payloads.signals(true, 3, "1.0.2", 1)
        assertEquals(
            setOf("screen_time_granted", "monitored_app_count", "app_version", "support_taps"),
            signals.keys,
        )
        assertEquals(3, signals.getValue("monitored_app_count"))
    }

    @Test
    fun protectionTypesMatchServerAllowlist() {
        val body = Payloads.protection("monitoring_disabled", "Alerts were turned off.")
        assertEquals("monitoring_disabled", body.getValue("type"))
        assertFalse(EventPolicy.leaks(body.getValue("detail")))
    }

    @Test
    fun tourIsSkippableAndCapsGroups() {
        assertTrue(TourContent.shouldAutoPresent(0, false))
        assertFalse(TourContent.shouldAutoPresent(1, false))
        assertTrue(TourContent.shouldAutoPresent(1, true))
        val withGroup = TourContent.steps(true)
        assertTrue(withGroup.any { it.body.contains("Groups hold up to 12 people.") })
        assertEquals(0, withGroup.first().tab)
        assertEquals(2, withGroup.last().tab)
    }

    @Test
    fun copyHasNoDashesOrBannedFraming() {
        val banned = listOf("\u2014", "\u2013", "\u2015")
        Copy.allStrings().forEach { line ->
            banned.forEach { mark ->
                assertFalse("dash in: $line", line.contains(mark))
            }
            val framing = "hon" + "est"
            assertFalse("framing in: $line", line.lowercase().contains(framing))
        }
    }

    @Test
    fun jwtRefreshWindow() {
        val token = "aaa." + java.util.Base64.getUrlEncoder().withoutPadding()
            .encodeToString("""{"exp":10000}""".toByteArray()) + ".sig"
        assertEquals(10000L, JwtClock.expiresAtEpochSeconds(token))
        assertTrue(JwtClock.shouldRefresh(token, 9000))
        assertFalse(JwtClock.shouldRefresh(token, 5000))
    }

    @Test
    fun streakAndStatus() {
        assertEquals(StreakMark.BeforeSignup, StreakWeek.mark("before_account"))
        assertEquals(StreakMark.Clean, StreakWeek.mark("true"))
        val status = MonitoringStatus.make(true, 2, true, false, false)
        assertEquals(MonitoringKind.On, status.kind)
        assertTrue(status.detail.contains("Alerts on"))
        assertTrue(status.deepScanLine.contains("Deep Scan is off"))
        val hidden = MonitoringStatus.make(true, 2, true, false, false, deepScanAvailable = false)
        assertEquals("", hidden.deepScanLine)
        assertFalse(TourContent.steps(false, deepScan = false).any { it.body.contains("Deep Scan") })
    }

    @Test
    fun feedbackRules() {
        assertNull(FeedbackRules.error("Idea", "Shorter labels", "The home card is dense."))
        assertEquals("summary is required (120 characters max)", FeedbackRules.error("Problem", "", "details"))
    }

    private fun baseInput() = MonitorInput(
        selected = setOf("com.example.browser"),
        alertsOn = true,
        blockingOn = false,
        lockoutOn = false,
        deepScanActive = false,
        alreadyAlertedToday = emptySet(),
        foregroundPackage = "com.example.browser",
        foregroundMsToday = mapOf("com.example.browser" to 60_000L),
        openedSinceLastCheck = emptySet(),
        ownPackage = "com.remainfaithful.app",
    )
}
