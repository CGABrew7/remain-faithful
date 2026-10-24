package com.remainfaithful.app.data

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.remainfaithful.core.HeartbeatPolicy
import com.remainfaithful.core.Payloads
import java.time.LocalDate

class Repository private constructor(context: Context) {
    private val app = context.applicationContext
    private val prefs: SharedPreferences = openPrefs(app)
    private val api = ApiClient(
        baseUrl = com.remainfaithful.app.BuildConfig.API_BASE_URL,
        tokenProvider = { token },
        onToken = { saveToken(it) },
    )

    var token: String?
        get() = prefs.getString(KEY_TOKEN, null)
        private set(value) {
            prefs.edit().putString(KEY_TOKEN, value).apply()
        }

    val isSignedIn: Boolean get() = !token.isNullOrBlank()

    fun saveSession(session: AuthSession) {
        token = session.token
        prefs.edit()
            .putLong(KEY_USER_ID, session.user.id)
            .putString(KEY_NAME, session.user.name)
            .putString(KEY_EMAIL, session.user.email)
            .putBoolean(KEY_VERIFIED, session.user.emailVerified != false)
            .putString(KEY_PENDING, session.user.pendingEmail)
            .putString(KEY_CREATED, session.user.createdAt)
            .apply()
    }

    fun saveUser(user: RemoteUser) {
        prefs.edit()
            .putLong(KEY_USER_ID, user.id)
            .putString(KEY_NAME, user.name)
            .putString(KEY_EMAIL, user.email)
            .putBoolean(KEY_VERIFIED, user.emailVerified != false)
            .putString(KEY_PENDING, user.pendingEmail)
            .putString(KEY_CREATED, user.createdAt)
            .apply()
    }

    fun saveToken(value: String) {
        token = value
    }

    val userId: Long get() = prefs.getLong(KEY_USER_ID, 0L)
    val userName: String get() = prefs.getString(KEY_NAME, "").orEmpty()
    val userEmail: String get() = prefs.getString(KEY_EMAIL, "").orEmpty()
    val emailVerified: Boolean get() = prefs.getBoolean(KEY_VERIFIED, true)
    val pendingEmail: String? get() = prefs.getString(KEY_PENDING, null)

    var onboardingDone: Boolean
        get() = prefs.getBoolean(KEY_ONBOARDING, false)
        set(value) { prefs.edit().putBoolean(KEY_ONBOARDING, value).apply() }

    var tourVersion: Int
        get() = prefs.getInt(KEY_TOUR, 0)
        set(value) { prefs.edit().putInt(KEY_TOUR, value).apply() }

    var alertsOn: Boolean
        get() = prefs.getBoolean(KEY_ALERTS, false)
        set(value) { prefs.edit().putBoolean(KEY_ALERTS, value).apply() }

    var blockingOn: Boolean
        get() = prefs.getBoolean(KEY_BLOCKING, false)
        set(value) { prefs.edit().putBoolean(KEY_BLOCKING, value).apply() }

    var lockoutOn: Boolean
        get() = prefs.getBoolean(KEY_LOCKOUT, false)
        set(value) { prefs.edit().putBoolean(KEY_LOCKOUT, value).apply() }

    var deepScanRunning: Boolean
        get() = prefs.getBoolean(KEY_SCAN, false)
        set(value) { prefs.edit().putBoolean(KEY_SCAN, value).apply() }

    var lastFrameAt: Long
        get() = prefs.getLong(KEY_FRAME, 0L)
        set(value) { prefs.edit().putLong(KEY_FRAME, value).apply() }

    var primaryGroupId: Long
        get() = prefs.getLong(KEY_GROUP, 0L)
        set(value) { prefs.edit().putLong(KEY_GROUP, value).apply() }

    var supportTaps: Int
        get() = prefs.getInt(KEY_TAPS, 0)
        set(value) { prefs.edit().putInt(KEY_TAPS, value).apply() }

    var revocationSent: Boolean
        get() = prefs.getBoolean(KEY_REVOKED, false)
        set(value) { prefs.edit().putBoolean(KEY_REVOKED, value).apply() }

    var chosenPackages: Set<String>
        get() = prefs.getStringSet(KEY_PACKAGES, emptySet()).orEmpty()
        set(value) { prefs.edit().putStringSet(KEY_PACKAGES, value.toSet()).apply() }

    fun alertedToday(): Set<String> {
        val today = LocalDate.now().toString()
        if (prefs.getString(KEY_ALERT_DAY, "") != today) return emptySet()
        return prefs.getStringSet(KEY_ALERTED, emptySet()).orEmpty()
    }

    fun markAlerted(packages: Set<String>) {
        if (packages.isEmpty()) return
        val today = LocalDate.now().toString()
        val current = if (prefs.getString(KEY_ALERT_DAY, "") == today) alertedToday() else emptySet()
        prefs.edit()
            .putString(KEY_ALERT_DAY, today)
            .putStringSet(KEY_ALERTED, current + packages)
            .apply()
    }

    fun clearSession() {
        alertsOn = false
        blockingOn = false
        deepScanRunning = false
        token = null
        prefs.edit()
            .remove(KEY_USER_ID)
            .remove(KEY_NAME)
            .remove(KEY_EMAIL)
            .apply()
    }

    fun api(): ApiClient = api

    fun sendHeartbeat() {
        if (!isSignedIn) return
        val since = if (lastFrameAt == 0L) null else System.currentTimeMillis() - lastFrameAt
        val screen = HeartbeatPolicy.screen(alertsOn, blockingOn, deepScanRunning, since) ?: return
        api.heartbeat(HeartbeatPolicy.body(screen).getValue("screen"))
    }

    fun sendSignals(usageGranted: Boolean) {
        if (!isSignedIn) return
        api.signals(
            Payloads.signals(
                usageGranted = usageGranted,
                monitoredCount = chosenPackages.size,
                appVersion = com.remainfaithful.app.BuildConfig.VERSION_NAME,
                supportTaps = supportTaps,
            ),
        )
    }

    fun registerPush(pushToken: String) {
        if (!isSignedIn || pushToken.isBlank()) return
        api.registerDevice(Payloads.deviceToken(pushToken))
    }

    fun protection(type: String, detail: String) {
        if (!isSignedIn) return
        api.protection(type, detail)
    }

    companion object {
        @Volatile private var instance: Repository? = null

        fun get(context: Context): Repository =
            instance ?: synchronized(this) {
                instance ?: Repository(context.applicationContext).also { instance = it }
            }

        private const val KEY_TOKEN = "auth_token"
        private const val KEY_USER_ID = "user_id"
        private const val KEY_NAME = "user_name"
        private const val KEY_EMAIL = "user_email"
        private const val KEY_VERIFIED = "email_verified"
        private const val KEY_PENDING = "pending_email"
        private const val KEY_CREATED = "created_at"
        private const val KEY_ONBOARDING = "onboarding_done"
        private const val KEY_TOUR = "tour_version"
        private const val KEY_ALERTS = "alerts_on"
        private const val KEY_BLOCKING = "blocking_on"
        private const val KEY_LOCKOUT = "lockout_on"
        private const val KEY_SCAN = "deep_scan"
        private const val KEY_FRAME = "last_frame"
        private const val KEY_GROUP = "primary_group"
        private const val KEY_TAPS = "support_taps"
        private const val KEY_REVOKED = "revocation_sent"
        private const val KEY_PACKAGES = "chosen_packages"
        private const val KEY_ALERT_DAY = "alert_day"
        private const val KEY_ALERTED = "alerted_packages"
    }
}

private fun openPrefs(context: Context): SharedPreferences {
    return try {
        val masterKey = MasterKey.Builder(context)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build()
        EncryptedSharedPreferences.create(
            context,
            "rf_secure",
            masterKey,
            EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
            EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM,
        )
    } catch (_: Exception) {
        context.getSharedPreferences("rf_secure_fallback", Context.MODE_PRIVATE)
    }
}
