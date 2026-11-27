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
import com.remainfaithful.core.BlurPolicy
import com.remainfaithful.core.DevicePolicy
import com.remainfaithful.core.FrameHygiene
import com.remainfaithful.app.theme.RfColors
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.delay
import kotlinx.coroutines.isActive
import kotlinx.coroutines.launch

/**
 * User-started screen session. Each frame is closed in memory. Frames are never
 * written, logged, or uploaded. A blur overlay can follow an on-phone score.
 * No weight file ships in this build, so the overlay stays down until one is chosen.
 */
class DeepScanService : Service() {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.Default)
    private var loop: Job? = null
    private var projection: MediaProjection? = null
    private var reader: ImageReader? = null
    private var display: android.hardware.display.VirtualDisplay? = null
    private var notifiedStop = false
    private var ending = false
    private var lastSampleAt = 0L

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
        val size = captureSize()
        reader = ImageReader.newInstance(size.first, size.second, PixelFormat.RGBA_8888, 2)
        reader?.setOnImageAvailableListener({ imageReader ->
            val image = imageReader.acquireLatestImage() ?: return@setOnImageAvailableListener
            val now = System.currentTimeMillis()
            try {
                if (!FrameHygiene.closeBeforeUse) return@setOnImageAvailableListener
                if (FrameHygiene.outputPath() != null) return@setOnImageAvailableListener
                Repository.get(this).lastFrameAt = now
                if (!BlurPolicy.weightsReady() || !BlurPolicy.shouldSample(now, lastSampleAt)) return@setOnImageAvailableListener
                lastSampleAt = now
                val sensitivity = Repository.get(this).blurSensitivity
                val command = BlurPolicy.command(null, sensitivity)
                Handler(Looper.getMainLooper()).post {
                    BlurOverlay.apply(this, command, null)
                }
            } finally {
                image.close()
            }
        }, Handler(Looper.getMainLooper()))
        display = projection?.createVirtualDisplay(
            "rf-deep-scan",
            size.first,
            size.second,
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
        BlurOverlay.hide(this)
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

    private fun captureSize(): Pair<Int, Int> {
        if (!BlurPolicy.weightsReady()) return 32 to 32
        val metrics = resources.displayMetrics
        return BlurPolicy.workingSize(metrics.widthPixels, metrics.heightPixels)
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
