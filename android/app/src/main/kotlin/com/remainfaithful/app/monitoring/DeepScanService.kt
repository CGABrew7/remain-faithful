package com.remainfaithful.app.monitoring

import android.app.Activity
import android.app.Notification
import android.app.PendingIntent
import android.app.Service
import android.content.Context
import android.content.Intent
import android.content.pm.ServiceInfo
import android.graphics.PixelFormat
import android.hardware.display.DisplayManager
import android.media.ImageReader
import android.media.projection.MediaProjection
import android.media.projection.MediaProjectionManager
import android.os.Build
import android.os.Handler
import android.os.IBinder
import android.os.Looper
import androidx.core.app.NotificationCompat
import com.remainfaithful.app.MainActivity
import com.remainfaithful.app.R
import com.remainfaithful.app.data.ApiException
import com.remainfaithful.app.data.Repository
import com.remainfaithful.core.DevicePolicy
import com.remainfaithful.app.theme.RfColors
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.delay
import kotlinx.coroutines.isActive
import kotlinx.coroutines.launch

/**
 * User-started screen session. Frames are closed immediately and are never
 * written, logged, or uploaded. Android has no equivalent of Apple's on-device
 * sensitive-content API, so this session does not invent a category from pixels.
 */
class DeepScanService : Service() {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.Default)
    private var loop: Job? = null
    private var projection: MediaProjection? = null
    private var reader: ImageReader? = null
    private var display: android.hardware.display.VirtualDisplay? = null
    private var notifiedStop = false
    private var ending = false

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val note = notification()
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            startForeground(NOTE_ID, note, ServiceInfo.FOREGROUND_SERVICE_TYPE_MEDIA_PROJECTION)
        } else {
            startForeground(NOTE_ID, note)
        }
        val code = intent?.getIntExtra(EXTRA_CODE, Activity.RESULT_CANCELED) ?: Activity.RESULT_CANCELED
        val data = intent?.let { readData(it) }
        if (code != Activity.RESULT_OK || data == null) {
            stopSelf()
            return START_NOT_STICKY
        }
        val repo = Repository.get(this)
        repo.deepScanRunning = true
        val manager = getSystemService(MediaProjectionManager::class.java)
        projection = manager.getMediaProjection(code, data)
        projection?.registerCallback(object : MediaProjection.Callback() {
            override fun onStop() {
                end(notify = true)
            }
        }, Handler(Looper.getMainLooper()))
        reader = ImageReader.newInstance(32, 32, PixelFormat.RGBA_8888, 2)
        reader?.setOnImageAvailableListener({ imageReader ->
            val image = imageReader.acquireLatestImage() ?: return@setOnImageAvailableListener
            image.close()
            Repository.get(this).lastFrameAt = System.currentTimeMillis()
        }, Handler(Looper.getMainLooper()))
        display = projection?.createVirtualDisplay(
            "rf-deep-scan",
            32,
            32,
            resources.displayMetrics.densityDpi,
            DisplayManager.VIRTUAL_DISPLAY_FLAG_AUTO_MIRROR,
            reader?.surface,
            null,
            null,
        )
        loop = scope.launch {
            while (isActive) {
                val error = runCatching { Repository.get(this@DeepScanService).sendHeartbeat() }.exceptionOrNull()
                if (error is ApiException && DevicePolicy.removed(error.status, error.code)) {
                    end(notify = false)
                    stopSelf()
                    break
                }
                delay(120_000)
            }
        }
        return START_NOT_STICKY
    }

    override fun onDestroy() {
        end(notify = true)
        super.onDestroy()
    }

    private fun end(notify: Boolean) {
        if (ending) return
        ending = true
        loop?.cancel()
        display?.release()
        reader?.close()
        projection?.stop()
        display = null
        reader = null
        projection = null
        val repo = Repository.get(this)
        val wasRunning = repo.deepScanRunning
        repo.deepScanRunning = false
        if (notify && wasRunning && !notifiedStop) {
            notifiedStop = true
            runCatching { repo.protection("deep_scan_stopped", "Deep Scan stopped.") }
        }
    }

    private fun notification(): Notification {
        val pending = PendingIntent.getActivity(
            this, 3, Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
        )
        return NotificationCompat.Builder(this, MonitorService.CHANNEL_MONITOR)
            .setSmallIcon(R.drawable.ic_stat)
            .setColor(RfColors.goldArgb)
            .setContentTitle("Deep Scan")
            .setContentText("Running on this phone. The screen is not uploaded.")
            .setOngoing(true)
            .setContentIntent(pending)
            .build()
    }

    companion object {
        const val NOTE_ID = 42
        const val EXTRA_CODE = "code"
        const val EXTRA_DATA = "data"

        fun start(context: Context, code: Int, data: Intent) {
            val intent = Intent(context, DeepScanService::class.java)
                .putExtra(EXTRA_CODE, code)
                .putExtra(EXTRA_DATA, data)
            androidx.core.content.ContextCompat.startForegroundService(context, intent)
        }

        private fun readData(intent: Intent): Intent? {
            return if (Build.VERSION.SDK_INT >= 33) {
                intent.getParcelableExtra(EXTRA_DATA, Intent::class.java)
            } else {
                @Suppress("DEPRECATION")
                intent.getParcelableExtra(EXTRA_DATA)
            }
        }
    }
}
