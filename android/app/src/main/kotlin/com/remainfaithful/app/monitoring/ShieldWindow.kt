package com.remainfaithful.app.monitoring

import android.content.Context
import android.content.Intent
import android.graphics.PixelFormat
import android.view.Gravity
import android.view.WindowManager
import android.widget.Button
import android.widget.LinearLayout
import android.widget.TextView
import com.remainfaithful.app.MainActivity
import com.remainfaithful.app.theme.RfColors
import com.remainfaithful.core.Copy

object ShieldWindow {
    private var view: LinearLayout? = null

    fun show(context: Context, lockout: Boolean) {
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
            text = if (lockout) Copy.shieldLockout else Copy.shieldOpen
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
        val current = view ?: return
        val window = context.applicationContext.getSystemService(WindowManager::class.java)
        runCatching { window.removeView(current) }
        view = null
    }
}
