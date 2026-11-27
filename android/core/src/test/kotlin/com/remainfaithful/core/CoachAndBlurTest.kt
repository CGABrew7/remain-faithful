package com.remainfaithful.core

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test
import java.io.File

class CoachAndBlurTest {
    private fun clock(hour: Int, weekday: Int = 4, day: Int = 10) = CoachClock(
        year = 2026, month = 11, day = day, hour = hour, minute = 0, weekday = weekday,
        epochMillis = 1_700_000_000_000L + hour * 3_600_000L,
    )

    private fun signals(
        enabled: Boolean = true,
        partner: Boolean = false,
        away: Boolean = false,
        missed: Boolean = false,
        opens: Int = 0,
        threshold: Boolean = false,
        session: Boolean = false,
        lastPrompt: Long? = null,
        lastNudge: CoachClock? = null,
        lateStart: Int = 23,
        lateEnd: Int = 5,
    ) = CoachSignals(
        enabled, partner, away, lateStart, lateEnd, missed, opens, threshold, session, lastPrompt, lastNudge,
    )

    @Test
    fun lateNightSundayMissedAwayAndQuietPeriod() {
        assertTrue(CoachPolicy.lateNight(23, 23, 5))
        assertTrue(CoachPolicy.lateNight(0, 23, 5))
        assertTrue(CoachPolicy.lateNight(4, 23, 5))
        assertFalse(CoachPolicy.lateNight(5, 23, 5))
        assertFalse(CoachPolicy.lateNight(3, 3, 3))
        val night = CoachPolicy.evaluate(clock(23), signals(), CoachReply.Pending)
        assertTrue(night.showPrompt)
        assertFalse(night.sendPartnerNudge)
        assertFalse(CoachPolicy.evaluate(clock(12), signals(), CoachReply.Pending).showPrompt)
        assertTrue(CoachPolicy.evaluate(clock(18, weekday = 1), signals(), CoachReply.Pending).showPrompt)
        assertFalse(CoachPolicy.evaluate(clock(18, weekday = 2), signals(), CoachReply.Pending).showPrompt)
        assertFalse(CoachPolicy.evaluate(clock(17), signals(missed = true), CoachReply.Pending).showPrompt)
        assertTrue(CoachPolicy.evaluate(clock(18), signals(missed = true), CoachReply.Pending).showPrompt)
        val away = signals(away = true, lateStart = 0, lateEnd = 0)
        assertTrue(CoachPolicy.evaluate(clock(10), away, CoachReply.Pending).showPrompt)
        assertFalse(CoachPolicy.evaluate(clock(3), away, CoachReply.Pending).showPrompt)
        val now = clock(23)
        val quiet = signals(lastPrompt = now.epochMillis - 1_000)
        assertFalse(CoachPolicy.evaluate(now, quiet, CoachReply.Pending).showPrompt)
        val later = now.copy(epochMillis = quiet.lastPromptAt!! + CoachPolicy.quietMillis)
        assertTrue(CoachPolicy.evaluate(later, quiet, CoachReply.Pending).showPrompt)
    }

    @Test
    fun notNowAndOptInAndDailyCap() {
        val open = signals(partner = true)
        val now = clock(23)
        val declined = CoachPolicy.evaluate(now, open, CoachReply.NotNow)
        assertFalse(declined.showPrompt)
        assertFalse(declined.sendPartnerNudge)
        assertFalse(CoachPolicy.evaluate(now, open, CoachReply.Pending).sendPartnerNudge)
        assertFalse(CoachPolicy.evaluate(now, signals(partner = false), CoachReply.CheckIn).sendPartnerNudge)
        assertTrue(CoachPolicy.evaluate(now, open, CoachReply.CheckIn).sendPartnerNudge)
        assertFalse(CoachPolicy.evaluate(now, signals(partner = true, lastNudge = now), CoachReply.CheckIn).sendPartnerNudge)
        assertTrue(CoachPolicy.evaluate(clock(23, day = 11), signals(partner = true, lastNudge = now), CoachReply.CheckIn).sendPartnerNudge)
        assertFalse(CoachPolicy.evaluate(now, signals(enabled = false, partner = true), CoachReply.CheckIn).sendPartnerNudge)
        assertFalse(CoachPolicy.evaluate(clock(12), signals(opens = 1), CoachReply.Pending).showPrompt)
        assertTrue(CoachPolicy.evaluate(clock(12), signals(opens = 2), CoachReply.Pending).showPrompt)
        assertTrue(CoachPolicy.evaluate(clock(12), signals(threshold = true), CoachReply.Pending).showPrompt)
    }

    @Test
    fun partnerSentenceAndWipeStayFixed() {
        assertEquals(Copy.partnerSentence, "Your partner would welcome a word from you today.")
        assertTrue(Payloads.protectionTypes.contains(CoachPolicy.alertType))
        val body = Payloads.protection(CoachPolicy.alertType, Copy.partnerSentence)
        assertEquals(setOf("type", "detail"), body.keys)
        assertEquals(Copy.partnerSentence, body.getValue("detail"))
        assertFalse(body.values.any { it.contains("http") || it.contains("com.") })
        val bag = mutableMapOf<String, Any?>("imAway" to true, "coachPartnerNudgeEnabled" to true, "other" to 1)
        CoachKeys.wipe(bag)
        assertFalse(bag.containsKey("imAway"))
        assertEquals(1, bag.getValue("other"))
        assertEquals("11 PM", CoachPolicy.hourLabel(23))
        assertEquals("5 AM", CoachPolicy.hourLabel(5))
    }

    @Test
    fun blurThresholdRegionAndNoFrameFile() {
        val wide = BlurPolicy.workingSize(1080, 2400)
        assertTrue(wide.first <= BlurPolicy.longEdge && wide.second <= BlurPolicy.longEdge)
        assertTrue(maxOf(wide.first, wide.second) == BlurPolicy.longEdge || maxOf(wide.first, wide.second) <= BlurPolicy.longEdge)
        assertTrue(BlurPolicy.threshold(1f) < BlurPolicy.threshold(0f))
        assertEquals(BlurCommand.Keep, BlurPolicy.command(null, 0.55f))
        assertEquals(BlurCommand.Hide, BlurPolicy.command(FrameSample(0.1f, null), 0.2f))
        assertEquals(BlurCommand.Full, BlurPolicy.command(FrameSample(0.99f, null), 0.8f))
        val box = BlurBox(0.1f, 0.2f, 0.4f, 0.6f)
        assertEquals(BlurCommand.Region, BlurPolicy.command(FrameSample(0.99f, box), 0.8f))
        assertEquals(BlurCommand.Full, BlurPolicy.command(FrameSample(0.99f, BlurBox(0.2f, 0.2f, 0.1f, 0.3f)), 0.8f))
        assertFalse(BlurPolicy.shouldSample(1000, 800))
        assertTrue(BlurPolicy.shouldSample(1000, 1000 - BlurPolicy.minIntervalMs))
        assertNull(FrameHygiene.outputPath())
        assertTrue(FrameHygiene.closeBeforeUse)
        assertFalse(BlurPolicy.weightsReady())
    }

    @Test
    fun voiceCapSkipAndWipe() {
        assertFalse(FutureSelfRules.keep(0))
        assertTrue(FutureSelfRules.keep(1_000))
        assertTrue(FutureSelfRules.keep(45_000))
        assertFalse(FutureSelfRules.keep(46_000))
        assertFalse(FutureSelfRules.skipSendsPartnerNotice)
        assertFalse(FutureSelfRules.shouldAutoPlay(false, false))
        assertTrue(FutureSelfRules.shouldAutoPlay(true, false))
        val root = File(System.getProperty("java.io.tmpdir"), "rf-voice-${System.nanoTime()}")
        root.mkdirs()
        val source = File(root, "take.m4a")
        source.writeText("first")
        assertTrue(FutureSelfFiles.replace(root, FutureSelfFiles.whyName, source, 12_000))
        val again = File(root, "take2.m4a")
        again.writeText("second")
        assertTrue(FutureSelfFiles.replace(root, FutureSelfFiles.whyName, again, 20_000))
        assertEquals("second", File(FutureSelfFiles.directory(root), FutureSelfFiles.whyName).readText())
        assertEquals(1, FutureSelfFiles.count(root))
        val tooLong = File(root, "long.m4a")
        tooLong.writeText("nope")
        assertFalse(FutureSelfFiles.replace(root, FutureSelfFiles.rememberName, tooLong, 46_000))
        FutureSelfFiles.wipe(root)
        assertFalse(FutureSelfFiles.directory(root).exists())
        assertEquals(0, FutureSelfFiles.count(root))
    }
}
