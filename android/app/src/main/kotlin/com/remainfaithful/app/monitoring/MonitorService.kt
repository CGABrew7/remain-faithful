package com.remainfaithful.app.monitoring

import android.app.AppOpsManager
import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.app.usage.UsageEvents
import android.app.usage.UsageStatsManager
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.pm.ServiceInfo
import android.os.Build
import android.os.IBinder
import androidx.core.app.NotificationCompat
import com.remainfaithful.app.MainActivity
import com.remainfaithful.app.R
import com.remainfaithful.app.data.ApiException
import com.remainfaithful.app.data.Repository
import com.remainfaithful.app.theme.RfColors
import com.remainfaithful.core.Copy
import com.remainfaithful.core.DevicePolicy
import com.remainfaithful.core.EventPolicy
import com.remainfaithful.core.MonitorInput
import com.remainfaithful.core.MonitoringPolicy
import com.remainfaithful.core.Payloads
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.delay
import kotlinx.coroutines.isActive
import kotlinx.coroutines.launch
import java.util.Calendar

class MonitorService : Service() {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.Default)
    private var loop: Job? = null
    private var lastCheck = System.currentTimeMillis()
    private var lastHeartbeat = 0L
    private val recentBlocked = mutableMapOf<String, Long>()

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val notification = ongoing("Alerts and Blocking stay as you set them.")
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.UPSIDE_DOWN_CAKE) {
            startForeground(NOTE_ID, notification, ServiceInfo.FOREGROUND_SERVICE_TYPE_SPECIAL_USE)
        } else {
            startForeground(NOTE_ID, notification)
        }
        if (loop == null) {
            loop = scope.launch {
                while (isActive) {
                    runCatching { tick() }
                    delay(20_000)
                }
            }
        }
        return START_STICKY
    }

    override fun onDestroy() {
        loop?.cancel()
        ShieldWindow.hide(this)
        super.onDestroy()
    }

    private fun tick() {
        val repo = Repository.get(this)
        if (!repo.alertsOn && !repo.blockingOn) {
            stopSelf()
            return
        }
        if (!UsageAccess.granted(this)) {
            if (!repo.revocationSent) {
                repo.revocationSent = true
                repo.alertsOn = false
                repo.blockingOn = false
                runCatching {
                    repo.protection(
                        "family_controls_revoked",
                        "Usage access was turned off. Alerts and Blocking are off.",
                    )
                }
            }
            stopSelf()
            return
        }
        repo.revocationSent = false
        val now = System.currentTimeMillis()
        val snapshot = UsageAccess.snapshot(this, lastCheck)
        lastCheck = now
        val opened = snapshot.opened.filter { pkg ->
            val at = recentBlocked[pkg] ?: 0L
            if (now - at < 5 * 60_000) false else {
                recentBlocked[pkg] = now
                true
            }
        }.toSet()
        val decision = MonitoringPolicy.decide(
            MonitorInput(
                selected = repo.chosenPackages,
                alertsOn = repo.alertsOn,
                blockingOn = repo.blockingOn,
                lockoutOn = repo.lockoutOn,
                deepScanActive = repo.deepScanRunning,
                alreadyAlertedToday = repo.alertedToday(),
                foregroundPackage = snapshot.foreground,
                foregroundMsToday = snapshot.msToday,
                openedSinceLastCheck = if (repo.blockingOn) opened else emptySet(),
                ownPackage = packageName,
            ),
        )
        decision.uploads.forEach { payload ->
            runCatching { repo.api().createEvent(payload) }
        }
        repo.markAlerted(decision.alertedPackages)
        if (decision.showShield && SettingsCompat.canDraw(this)) {
            ShieldWindow.show(this, repo.lockoutOn, com.remainfaithful.app.BuildConfig.DEEP_SCAN_ENABLED)
        } else {
            ShieldWindow.hide(this)
            if (decision.showShield) {
                Notify.alert(this, "Remain Faithful", "Your partner is standing with you.", "CONTENT_FLAGGED")
            }
        }
        if (now - lastHeartbeat > 120_000) {
            lastHeartbeat = now
            val error = runCatching { repo.sendHeartbeat() }.exceptionOrNull()
            if (error is ApiException && DevicePolicy.removed(error.status, error.code)) {
                Notify.alert(this, "Remain Faithful", Copy.deviceRemoved, "DEVICE_REMOVED")
                ShieldWindow.hide(this)
                stopSelf()
                return
            }
        }
    }

    private fun ongoing(text: String): Notification {
        val pending = PendingIntent.getActivity(
            this, 1, Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
        )
        return NotificationCompat.Builder(this, CHANNEL_MONITOR)
            .setSmallIcon(R.drawable.ic_stat)
            .setColor(RfColors.goldArgb)
            .setContentTitle("Remain Faithful")
            .setContentText(text)
            .setOngoing(true)
            .setContentIntent(pending)
            .build()
    }

    companion object {
        const val NOTE_ID = 41
        const val CHANNEL_MONITOR = "monitoring"

        fun start(context: Context) {
            val intent = Intent(context, MonitorService::class.java)
            androidx.core.content.ContextCompat.startForegroundService(context, intent)
        }

        fun stop(context: Context) {
            context.stopService(Intent(context, MonitorService::class.java))
        }
    }
}

object UsageAccess {
    fun granted(context: Context): Boolean {
        val appOps = context.getSystemService(AppOpsManager::class.java) ?: return false
        val mode = appOps.unsafeCheckOpNoThrow(
            AppOpsManager.OPSTR_GET_USAGE_STATS,
            android.os.Process.myUid(),
            context.packageName,
        )
        return mode == AppOpsManager.MODE_ALLOWED
    }

    fun snapshot(context: Context, since: Long): Snapshot {
        val usm = context.getSystemService(UsageStatsManager::class.java)
            ?: return Snapshot(null, emptyMap(), emptySet())
        val start = Calendar.getInstance().apply {
            set(Calendar.HOUR_OF_DAY, 0)
            set(Calendar.MINUTE, 0)
            set(Calendar.SECOND, 0)
            set(Calendar.MILLISECOND, 0)
        }.timeInMillis
        val ms = mutableMapOf<String, Long>()
        usm.queryUsageStats(UsageStatsManager.INTERVAL_DAILY, start, System.currentTimeMillis())
            ?.forEach { stat ->
                if (stat.totalTimeInForeground > 0) ms[stat.packageName] = stat.totalTimeInForeground
            }
        val opened = mutableSetOf<String>()
        var foreground: String? = null
        var foregroundAt = 0L
        val events = usm.queryEvents(since.coerceAtLeast(System.currentTimeMillis() - 60_000), System.currentTimeMillis())
        val event = UsageEvents.Event()
        while (events.hasNextEvent()) {
            events.getNextEvent(event)
            if (event.eventType == UsageEvents.Event.ACTIVITY_RESUMED) {
                opened += event.packageName
                if (event.timeStamp >= foregroundAt) {
                    foregroundAt = event.timeStamp
                    foreground = event.packageName
                }
            }
        }
        return Snapshot(foreground, ms, opened)
    }

    data class Snapshot(val foreground: String?, val msToday: Map<String, Long>, val opened: Set<String>)
}

object SettingsCompat {
    fun canDraw(context: Context): Boolean = android.provider.Settings.canDrawOverlays(context)
}

object Notify {
    fun alert(context: Context, title: String, body: String, type: String) {
        val channel = if (type == "PANIC_ALERT") CHANNEL_URGENT else CHANNEL_ALERTS
        val pending = PendingIntent.getActivity(
            context,
            type.hashCode(),
            Intent(context, MainActivity::class.java).putExtra(MainActivity.EXTRA_PUSH_TYPE, type),
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
        )
        val note = NotificationCompat.Builder(context, channel)
            .setSmallIcon(R.drawable.ic_stat)
            .setColor(RfColors.goldArgb)
            .setContentTitle(title)
            .setContentText(body)
            .setAutoCancel(true)
            .setContentIntent(pending)
            .setPriority(NotificationCompat.PRIORITY_HIGH)
            .build()
        val manager = context.getSystemService(NotificationManager::class.java)
        manager.notify((System.currentTimeMillis() % Int.MAX_VALUE).toInt(), note)
    }

    const val CHANNEL_ALERTS = "alerts"
    const val CHANNEL_URGENT = "urgent"
}

fun createChannels(context: Context) {
    val manager = context.getSystemService(NotificationManager::class.java)
    manager.createNotificationChannel(NotificationChannel(MonitorService.CHANNEL_MONITOR, context.getString(R.string.monitoring_channel), NotificationManager.IMPORTANCE_LOW))
    manager.createNotificationChannel(NotificationChannel(Notify.CHANNEL_ALERTS, context.getString(R.string.alerts_channel), NotificationManager.IMPORTANCE_HIGH))
    manager.createNotificationChannel(NotificationChannel(Notify.CHANNEL_URGENT, context.getString(R.string.urgent_channel), NotificationManager.IMPORTANCE_HIGH))
}

class BootReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent?) {
        val repo = Repository.get(context)
        if (!repo.alertsOn && !repo.blockingOn) return
        Notify.alert(
            context,
            "Remain Faithful",
            "Open the app to keep Alerts and Blocking on.",
            "MONITORING_STOPPED",
        )
    }
}

object EventUpload {
    fun safe(summary: String): Boolean = EventPolicy.isAllowedSummary(summary) && !EventPolicy.leaks(summary)
    fun protectionOk(type: String): Boolean = type in Payloads.protectionTypes
}
