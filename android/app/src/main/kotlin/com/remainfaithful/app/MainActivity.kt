package com.remainfaithful.app

import android.content.Intent
import android.net.Uri
import android.os.Build
import android.os.Bundle
import android.provider.Settings
import androidx.activity.ComponentActivity
import androidx.activity.SystemBarStyle
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.runtime.LaunchedEffect
import androidx.core.view.WindowCompat
import androidx.lifecycle.ViewModelProvider
import com.remainfaithful.app.monitoring.DeepScanService
import com.remainfaithful.app.theme.RfColors
import com.remainfaithful.app.theme.RfTheme
import com.remainfaithful.app.ui.AppModel
import com.remainfaithful.app.ui.RootScreen
import com.remainfaithful.app.ui.Sheet

class MainActivity : ComponentActivity() {
    private lateinit var model: AppModel

    private val notifications = registerForActivityResult(ActivityResultContracts.RequestPermission()) { }

    private val capture = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
        val data = result.data
        if (result.resultCode == RESULT_OK && data != null) {
            DeepScanService.start(this, result.resultCode, data)
            if (::model.isInitialized) model.deepScanOn = true
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge(statusBarStyle = SystemBarStyle.dark(RfColors.navyArgb))
        super.onCreate(savedInstanceState)
        WindowCompat.getInsetsController(window, window.decorView).isAppearanceLightStatusBars = false
        model = ViewModelProvider(this)[AppModel::class.java]
        handleIntent(intent)
        setContent {
            RfTheme {
                RootScreen(model)
                val nonce = model.systemNonce
                LaunchedEffect(nonce) {
                    when (model.systemRequest) {
                        "usage" -> startActivity(Intent(Settings.ACTION_USAGE_ACCESS_SETTINGS))
                        "overlay" -> startActivity(
                            Intent(Settings.ACTION_MANAGE_OVERLAY_PERMISSION, Uri.parse("package:$packageName")),
                        )
                        "notifications" -> if (Build.VERSION.SDK_INT >= 33) {
                            notifications.launch(android.Manifest.permission.POST_NOTIFICATIONS)
                        }
                        "capture" -> {
                            val manager = getSystemService(android.media.projection.MediaProjectionManager::class.java)
                            capture.launch(manager.createScreenCaptureIntent())
                        }
                    }
                    model.systemRequest = null
                }
            }
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        setIntent(intent)
        if (::model.isInitialized) handleIntent(intent)
    }

    override fun onResume() {
        super.onResume()
        if (::model.isInitialized) model.onResume()
    }

    private fun handleIntent(intent: Intent?) {
        val uri = intent?.data
        val token = uri?.getQueryParameter("token")
        if (!token.isNullOrBlank()) {
            val type = uri.getQueryParameter("type") ?: if (uri.path.orEmpty().contains("group")) "group" else "partner"
            model.acceptLink(token, type)
        }
        when (intent?.getStringExtra(EXTRA_OPEN)) {
            OPEN_PANIC -> model.sheet = Sheet.Panic
            OPEN_DEEP_SCAN -> model.askSystem("capture")
            OPEN_PIN -> {
                model.pendingAfterPin = { }
                model.sheet = Sheet.Pin
            }
        }
        when (intent?.getStringExtra(EXTRA_PUSH_TYPE)) {
            "PANIC_ALERT" -> model.sheet = Sheet.Panic
            "CONTENT_FLAGGED", "PROTECTION_ALERT" -> model.tab = 0
        }
    }

    companion object {
        const val EXTRA_OPEN = "open"
        const val EXTRA_PUSH_TYPE = "push_type"
        const val OPEN_PANIC = "panic"
        const val OPEN_DEEP_SCAN = "deep_scan"
        const val OPEN_PIN = "pin"
    }
}
