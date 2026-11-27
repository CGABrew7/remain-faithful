package com.remainfaithful.core

data class BlurBox(
    val left: Float,
    val top: Float,
    val right: Float,
    val bottom: Float,
) {
    fun valid(): Boolean =
        right > left && bottom > top && left >= 0f && top >= 0f && right <= 1f && bottom <= 1f
}

data class FrameSample(val score: Float, val box: BlurBox?)

enum class BlurCommand { Keep, Hide, Full, Region }

object BlurPolicy {
    const val longEdge = 480
    const val minIntervalMs = 333L

    fun weightsReady(): Boolean = false

    fun workingSize(width: Int, height: Int, edge: Int = longEdge): Pair<Int, Int> {
        val longSide = maxOf(width, height).coerceAtLeast(1)
        val scale = edge.toFloat() / longSide.toFloat()
        val w = (width * scale).toInt().coerceIn(32, edge)
        val h = (height * scale).toInt().coerceIn(32, edge)
        return w to h
    }

    fun threshold(sensitivity: Float): Float {
        val level = sensitivity.coerceIn(0f, 1f)
        return 0.90f - level * 0.55f
    }

    fun command(sample: FrameSample?, sensitivity: Float): BlurCommand {
        if (sample == null) return BlurCommand.Keep
        if (sample.score < threshold(sensitivity)) return BlurCommand.Hide
        val box = sample.box
        return if (box != null && box.valid()) BlurCommand.Region else BlurCommand.Full
    }

    fun shouldSample(now: Long, last: Long): Boolean = now - last >= minIntervalMs
}

object FrameHygiene {
    const val closeBeforeUse = true

    fun outputPath(): String? = null
}

object FutureSelfFiles {
    const val whyName = "why.m4a"
    const val rememberName = "remember.m4a"

    fun directory(base: java.io.File): java.io.File = java.io.File(base, "future-self")

    fun count(base: java.io.File): Int {
        val folder = directory(base)
        if (!folder.isDirectory) return 0
        return listOf(whyName, rememberName).count { java.io.File(folder, it).isFile }
    }

    fun replace(base: java.io.File, name: String, source: java.io.File, durationMs: Long): Boolean {
        if (!FutureSelfRules.keep(durationMs)) {
            source.delete()
            return false
        }
        val folder = directory(base)
        folder.mkdirs()
        val destination = java.io.File(folder, name)
        if (destination.exists()) destination.delete()
        if (!source.renameTo(destination)) {
            source.copyTo(destination, overwrite = true)
            source.delete()
        }
        return destination.isFile
    }

    fun wipe(base: java.io.File) {
        val folder = directory(base)
        folder.listFiles()?.forEach { it.delete() }
        folder.delete()
    }
}

object FutureSelfRules {
    const val minMs = 1_000L
    const val maxMs = 45_000L
    const val skipSendsPartnerNotice = false

    fun keep(durationMs: Long): Boolean = durationMs in minMs..maxMs

    fun shouldAutoPlay(hasClip: Boolean, alreadyPlayedThisOpen: Boolean): Boolean =
        hasClip && !alreadyPlayedThisOpen
}
