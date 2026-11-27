package com.remainfaithful.app.monitoring

import android.content.Context
import android.content.Intent
import android.graphics.PixelFormat
import android.media.AudioAttributes
import android.media.MediaPlayer
import android.view.Gravity
import android.view.WindowManager
import android.widget.Button
import android.widget.LinearLayout
import android.widget.TextView
import com.remainfaithful.app.MainActivity
import com.remainfaithful.app.theme.RfColors
import com.remainfaithful.core.Copy
import com.remainfaithful.core.FutureSelfFiles
import com.remainfaithful.core.FutureSelfRules
import java.io.File

object ShieldWindow {
    private var view: LinearLayout? = null
    private var player: MediaPlayer? = null

    fun show(context: Context, lockout: Boolean, deepScanEnabled: Boolean) {
        if (view != null) return
        if (!SettingsCompat.canDraw(context)) return
        val app = context.applicationContext
        val window = app.getSystemService(WindowManager::class.java)
        val layout = LinearLayout(app).apply {
            orientation = LinearLayout.VERTICAL
            gravity = Gravity.CENTER
            setBackgroundColor(RfColors.navyArgb)
            setPadding(64, 80, 64, 80)
        }
        layout.addView(TextView(app).apply {
            text = Copy.appName
            setTextColor(RfColors.goldArgb)
            textSize = 28f
            gravity = Gravity.CENTER
        })
        layout.addView(TextView(app).apply {
            text = when {
                lockout && deepScanEnabled -> Copy.shieldLockout
                lockout -> Copy.shieldLockoutPin
                else -> Copy.shieldOpen
            }
            setTextColor(0xFFFFFFFF.toInt())
            textSize = 16f
            gravity = Gravity.CENTER
            setPadding(0, 24, 0, 36)
        })
        layout.addView(Button(app).apply {
            text = Copy.needHelp
            setOnClickListener {
                hide(app)
                app.startActivity(Intent(app, MainActivity::class.java).apply {
                    addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                    putExtra(MainActivity.EXTRA_OPEN, MainActivity.OPEN_PANIC)
                })
            }
        })
        if (lockout) {
            if (deepScanEnabled) {
                layout.addView(Button(app).apply {
                    text = "Start Deep Scan"
                    setOnClickListener {
                        hide(app)
                        app.startActivity(Intent(app, MainActivity::class.java).apply {
                            addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                            putExtra(MainActivity.EXTRA_OPEN, MainActivity.OPEN_DEEP_SCAN)
                        })
                    }
                })
            }
            layout.addView(Button(app).apply {
                text = "Enter Partner PIN"
                setOnClickListener {
                    hide(app)
                    app.startActivity(Intent(app, MainActivity::class.java).apply {
                        addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                        putExtra(MainActivity.EXTRA_OPEN, MainActivity.OPEN_PIN)
                    })
                }
            })
        } else {
            layout.addView(Button(app).apply {
                text = "Continue"
                setOnClickListener { hide(app) }
            })
        }
        val clip = File(FutureSelfFiles.directory(app.filesDir), FutureSelfFiles.whyName)
        if (clip.isFile) {
            layout.addView(TextView(app).apply {
                text = Copy.voicePlayTitle
                setTextColor(RfColors.goldArgb)
                textSize = 16f
                gravity = Gravity.CENTER
                setPadding(0, 28, 0, 12)
            })
            layout.addView(Button(app).apply {
                text = Copy.voicePlayAgain
                setOnClickListener { play(app) }
            })
            layout.addView(Button(app).apply {
                text = Copy.voiceSkip
                setOnClickListener { stopVoice() }
            })
            if (FutureSelfRules.shouldAutoPlay(hasClip = true, alreadyPlayedThisOpen = false)) {
                play(app)
            }
        }
        val params = WindowManager.LayoutParams(
            WindowManager.LayoutParams.MATCH_PARENT,
            WindowManager.LayoutParams.MATCH_PARENT,
            WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
            WindowManager.LayoutParams.FLAG_LAYOUT_IN_SCREEN,
            PixelFormat.OPAQUE,
        )
        window.addView(layout, params)
        view = layout
    }

    fun hide(context: Context) {
        stopVoice()
        val current = view ?: return
        val window = context.applicationContext.getSystemService(WindowManager::class.java)
        runCatching { window.removeView(current) }
        view = null
    }

    private fun play(context: Context) {
        val clip = File(FutureSelfFiles.directory(context.applicationContext.filesDir), FutureSelfFiles.whyName)
        if (!clip.isFile) return
        stopVoice()
        val next = MediaPlayer()
        next.setAudioAttributes(
            AudioAttributes.Builder()
                .setUsage(AudioAttributes.USAGE_MEDIA)
                .setContentType(AudioAttributes.CONTENT_TYPE_SPEECH)
                .build(),
        )
        runCatching {
            next.setDataSource(clip.absolutePath)
            next.setOnPreparedListener { it.start() }
            next.prepareAsync()
            player = next
        }.onFailure {
            next.release()
        }
    }

    private fun stopVoice() {
        val current = player ?: return
        runCatching { current.stop() }
        runCatching { current.release() }
        player = null
    }
}
