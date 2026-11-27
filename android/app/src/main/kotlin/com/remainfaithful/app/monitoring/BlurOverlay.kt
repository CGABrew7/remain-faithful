package com.remainfaithful.app.monitoring

import android.content.Context
import android.graphics.PixelFormat
import android.os.Build
import android.view.Gravity
import android.view.View
import android.view.WindowManager
import android.widget.FrameLayout
import com.remainfaithful.core.BlurBox
import com.remainfaithful.core.BlurCommand

/**
 * On-phone blur during a user-started Deep Scan session.
 * The overlay never writes a frame and never reads a file.
 */
object BlurOverlay {
    private var root: FrameLayout? = null

    fun apply(context: Context, command: BlurCommand, box: BlurBox?) {
        when (command) {
            BlurCommand.Keep -> Unit
            BlurCommand.Hide -> hide(context)
            BlurCommand.Full -> show(context, null)
            BlurCommand.Region -> show(context, box)
        }
    }

    fun hide(context: Context) {
        val current = root ?: return
        val window = context.applicationContext.getSystemService(WindowManager::class.java)
        runCatching { window.removeView(current) }
        root = null
    }

    private fun show(context: Context, box: BlurBox?) {
        if (!SettingsCompat.canDraw(context)) return
        hide(context)
        val app = context.applicationContext
        val window = app.getSystemService(WindowManager::class.java)
        val layout = FrameLayout(app)
        val veil = View(app).apply {
            setBackgroundColor(0xE60B1733.toInt())
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
                setRenderEffect(
                    android.graphics.RenderEffect.createBlurEffect(
                        48f,
                        48f,
                        android.graphics.Shader.TileMode.CLAMP,
                    ),
                )
            }
        }
        if (box == null || !box.valid()) {
            layout.addView(
                veil,
                FrameLayout.LayoutParams(
                    FrameLayout.LayoutParams.MATCH_PARENT,
                    FrameLayout.LayoutParams.MATCH_PARENT,
                ),
            )
        } else {
            layout.addView(veil, regionParams(app, box))
        }
        val params = WindowManager.LayoutParams(
            WindowManager.LayoutParams.MATCH_PARENT,
            WindowManager.LayoutParams.MATCH_PARENT,
            WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
            WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE or
                WindowManager.LayoutParams.FLAG_NOT_TOUCHABLE or
                WindowManager.LayoutParams.FLAG_LAYOUT_IN_SCREEN,
            PixelFormat.TRANSLUCENT,
        )
        window.addView(layout, params)
        root = layout
    }

    private fun regionParams(context: Context, box: BlurBox): FrameLayout.LayoutParams {
        val metrics = context.resources.displayMetrics
        val width = ((box.right - box.left) * metrics.widthPixels).toInt().coerceAtLeast(1)
        val height = ((box.bottom - box.top) * metrics.heightPixels).toInt().coerceAtLeast(1)
        return FrameLayout.LayoutParams(width, height, Gravity.TOP or Gravity.START).apply {
            leftMargin = (box.left * metrics.widthPixels).toInt()
            topMargin = (box.top * metrics.heightPixels).toInt()
        }
    }
}
