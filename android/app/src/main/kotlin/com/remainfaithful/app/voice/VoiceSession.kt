package com.remainfaithful.app.voice

import android.content.Context
import android.media.MediaRecorder
import android.os.Build
import android.os.SystemClock
import com.remainfaithful.core.Copy
import com.remainfaithful.core.FutureSelfFiles
import com.remainfaithful.core.FutureSelfRules
import java.io.File

class VoiceSession(private val context: Context) {
    var recording: Boolean = false
        private set

    private var recorder: MediaRecorder? = null
    private var started = 0L
    private var clipName = FutureSelfFiles.whyName

    fun start(name: String) {
        stop(save = false)
        clipName = name
        val file = File(context.cacheDir, "take.m4a")
        if (file.exists()) file.delete()
        val next = if (Build.VERSION.SDK_INT >= 31) {
            MediaRecorder(context)
        } else {
            @Suppress("DEPRECATION")
            MediaRecorder()
        }
        next.setAudioSource(MediaRecorder.AudioSource.MIC)
        next.setOutputFormat(MediaRecorder.OutputFormat.MPEG_4)
        next.setAudioEncoder(MediaRecorder.AudioEncoder.AAC)
        next.setOutputFile(file.absolutePath)
        next.setMaxDuration(FutureSelfRules.maxMs.toInt())
        next.prepare()
        next.start()
        recorder = next
        started = SystemClock.elapsedRealtime()
        recording = true
    }

    fun stop(save: Boolean = true): String? {
        val elapsed = if (started == 0L) 0L else SystemClock.elapsedRealtime() - started
        recording = false
        started = 0L
        runCatching { recorder?.stop() }
        runCatching { recorder?.release() }
        recorder = null
        val file = File(context.cacheDir, "take.m4a")
        if (!save || !file.exists()) {
            file.delete()
            return null
        }
        val kept = FutureSelfFiles.replace(context.filesDir, clipName, file, elapsed)
        return if (kept) Copy.voiceSaved else Copy.voiceShort
    }
}
