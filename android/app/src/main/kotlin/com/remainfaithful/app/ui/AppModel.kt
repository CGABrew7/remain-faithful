package com.remainfaithful.app.ui

import android.app.Application
import android.content.Intent
import android.content.pm.PackageManager
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.remainfaithful.app.BuildConfig
import com.remainfaithful.app.auth.GoogleSignIn
import com.remainfaithful.app.auth.GoogleSignInException
import com.remainfaithful.app.data.ApiException
import com.remainfaithful.app.data.GroupInvite
import com.remainfaithful.app.data.GroupMember
import com.remainfaithful.app.data.PartnerInvite
import com.remainfaithful.app.data.RemoteAlert
import com.remainfaithful.app.data.RemoteEvent
import com.remainfaithful.app.data.RemoteGroup
import com.remainfaithful.app.data.RemoteRel
import com.remainfaithful.app.data.Repository
import com.remainfaithful.app.monitoring.MonitorService
import com.remainfaithful.app.monitoring.UsageAccess
import com.remainfaithful.app.data.RemoteDevice
import com.remainfaithful.core.AuthRules
import com.remainfaithful.core.DevicePolicy
import com.remainfaithful.core.CategoryLabels
import com.remainfaithful.core.Copy
import com.remainfaithful.core.GroupPolicy
import com.remainfaithful.core.MonitoringStatus
import com.remainfaithful.core.PinCheck
import com.remainfaithful.core.StreakWeek
import com.remainfaithful.core.TourContent
import com.google.firebase.messaging.FirebaseMessaging
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.time.Duration
import java.time.Instant
import java.time.LocalDate
import java.time.OffsetDateTime
import java.time.ZoneId

enum class Phase { Loading, Onboarding, Login, App }
enum class Sheet { Panic, Feedback, Pin, Activity, How, Delete, Partners, Monitoring, DeepScan, Alert, Profile, Forgot, RemoveDevice }

data class DeviceRow(
    val id: String,
    val label: String,
    val monitoring: Boolean,
    val lastCheckIn: String,
    val current: Boolean,
)

data class LocalApp(val packageName: String, val label: String)

data class UiEvent(
    val key: String,
    val alertId: Long?,
    val category: String,
    val severity: String,
    val summary: String,
    val timeLabel: String,
    val discussed: Boolean,
)

class AppModel(app: Application) : AndroidViewModel(app) {
    private val repo = Repository.get(app)

    var phase by mutableStateOf(Phase.Loading)
    var step by mutableIntStateOf(0)
    var tab by mutableIntStateOf(0)
    var banner by mutableStateOf<String?>(null)
    var busy by mutableStateOf(false)
    var sheet by mutableStateOf<Sheet?>(null)
    var systemRequest by mutableStateOf<String?>(null)
    var systemNonce by mutableIntStateOf(0)

    var name by mutableStateOf("")
    var email by mutableStateOf("")
    var verified by mutableStateOf(true)
    var pendingEmail by mutableStateOf<String?>(null)
    var userId by mutableStateOf(0L)

    var streak by mutableIntStateOf(0)
    var best by mutableIntStateOf(0)
    var week by mutableStateOf(listOf<com.remainfaithful.core.StreakMark>())
    var mine by mutableStateOf(listOf<UiEvent>())
    var alerts by mutableStateOf(listOf<UiEvent>())
    var selectedAlert by mutableStateOf<UiEvent?>(null)

    var alertsOn by mutableStateOf(false)
    var blockingOn by mutableStateOf(false)
    var lockoutOn by mutableStateOf(false)
    var deepScanOn by mutableStateOf(false)
    var usageGranted by mutableStateOf(false)
    var overlayGranted by mutableStateOf(false)
    var apps by mutableStateOf(listOf<LocalApp>())
    var chosen by mutableStateOf(setOf<String>())
    var pinSet by mutableStateOf(false)

    var rels by mutableStateOf(listOf<RemoteRel>())
    var partnerInvites by mutableStateOf(listOf<PartnerInvite>())
    var groups by mutableStateOf(listOf<RemoteGroup>())
    var group by mutableStateOf<RemoteGroup?>(null)
    var groupInvites by mutableStateOf(listOf<GroupInvite>())

    var devices by mutableStateOf(listOf<DeviceRow>())
    var devicesUnavailable by mutableStateOf(false)
    var pendingDevice by mutableStateOf<DeviceRow?>(null)

    var tourOn by mutableStateOf(false)
    var tourIndex by mutableIntStateOf(0)
    var tourForced by mutableStateOf(false)
    var pendingAfterPin: (() -> Unit)? = null

    var adult by mutableStateOf(false)
    var formName by mutableStateOf("")
    var formEmail by mutableStateOf("")
    var formPassword by mutableStateOf("")

    val status get() = MonitoringStatus.make(
        usageGranted,
        chosen.size,
        alertsOn,
        blockingOn,
        deepScanOn,
        deepScanAvailable = BuildConfig.DEEP_SCAN_ENABLED,
    )
    val tourSteps get() = TourContent.steps(group != null || groups.isNotEmpty(), BuildConfig.DEEP_SCAN_ENABLED)
    val inGroup get() = (group?.members?.size ?: 0) > 0 || groups.isNotEmpty()

    init {
        viewModelScope.launch { boot() }
    }

    private suspend fun boot() {
        val context = getApplication<Application>()
        usageGranted = UsageAccess.granted(context)
        overlayGranted = com.remainfaithful.app.monitoring.SettingsCompat.canDraw(context)
        alertsOn = repo.alertsOn
        blockingOn = repo.blockingOn
        lockoutOn = repo.lockoutOn
        deepScanOn = repo.deepScanRunning
        chosen = repo.chosenPackages
        name = repo.userName
        email = repo.userEmail
        verified = repo.emailVerified
        pendingEmail = repo.pendingEmail
        userId = repo.userId
        phase = when {
            repo.isSignedIn -> Phase.App
            repo.onboardingDone -> Phase.Login
            else -> Phase.Onboarding
        }
        if (phase == Phase.App) {
            refresh()
            syncMonitoring()
            adoptRemovedIfNeeded()
            maybeTour(false)
            refreshPush()
        }
    }

    fun askSystem(kind: String) {
        systemRequest = kind
        systemNonce++
    }

    fun onResume() {
        val context = getApplication<Application>()
        usageGranted = UsageAccess.granted(context)
        overlayGranted = com.remainfaithful.app.monitoring.SettingsCompat.canDraw(context)
        deepScanOn = repo.deepScanRunning
        if (phase == Phase.App) {
            adoptRemovedIfNeeded()
            syncMonitoring()
            viewModelScope.launch {
                val error = runCatching { withContext(Dispatchers.IO) { repo.sendHeartbeat() } }.exceptionOrNull()
                if (error is ApiException && DevicePolicy.removed(error.status, error.code)) {
                    adoptRemoved()
                    banner = Copy.deviceRemoved
                } else {
                    runCatching { withContext(Dispatchers.IO) { reloadDevices() } }
                }
                runCatching { withContext(Dispatchers.IO) { repo.sendSignals(usageGranted) } }
            }
        }
    }

    fun syncMonitoring() {
        val context = getApplication<Application>()
        if ((alertsOn || blockingOn) && usageGranted) {
            MonitorService.start(context)
        } else if (!alertsOn && !blockingOn) {
            MonitorService.stop(context)
        }
    }

    fun refresh() {
        viewModelScope.launch {
            try {
                val me = withContext(Dispatchers.IO) { repo.api().me() }
                repo.saveUser(me)
                name = me.name
                email = me.email
                verified = me.emailVerified != false
                pendingEmail = me.pendingEmail
                userId = me.id
                val stats = withContext(Dispatchers.IO) { repo.api().stats() }
                streak = stats.streakDays
                best = stats.bestStreak
                week = StreakWeek.marks(stats.weekStates)
                mine = withContext(Dispatchers.IO) { repo.api().events() }.map { it.toUi() }
                alerts = withContext(Dispatchers.IO) { repo.api().alerts() }.map { it.toUi() }
                pinSet = withContext(Dispatchers.IO) { repo.api().pinSet() }
                rels = withContext(Dispatchers.IO) { repo.api().relationships() }
                partnerInvites = withContext(Dispatchers.IO) { repo.api().pendingPartners() }
                groups = withContext(Dispatchers.IO) { repo.api().groups() }
                groupInvites = withContext(Dispatchers.IO) { repo.api().pendingGroups() }
                val gid = repo.primaryGroupId.takeIf { id -> groups.any { it.id == id } } ?: groups.firstOrNull()?.id
                if (gid != null) {
                    repo.primaryGroupId = gid
                    group = withContext(Dispatchers.IO) { repo.api().group(gid) }
                } else {
                    group = null
                }
                withContext(Dispatchers.IO) { repo.sendSignals(usageGranted) }
                runCatching { reloadDevices() }
            } catch (e: Exception) {
                banner = message(e)
                if (e is ApiException && e.status == 401) signOutLocal()
            }
        }
    }

    fun register() {
        if (!AuthRules.canRegister(adult, formName, formEmail, formPassword)) {
            banner = if (!adult) Copy.ageGate else AuthRules.passwordHint(formPassword)
            return
        }
        launchIo {
            val user = repo.api().register(formName.trim(), formEmail.trim(), formPassword)
            val session = repo.api().login(formEmail.trim(), formPassword)
            repo.saveSession(session.copy(user = session.user.copy(emailVerified = user.emailVerified)))
            name = session.user.name
            email = session.user.email
            userId = session.user.id
            verified = user.emailVerified != false
            step = 2
            refreshPush()
        }
    }

    fun login() {
        launchIo {
            val session = repo.api().login(formEmail.trim(), formPassword)
            applySession(session, enterApp = true)
        }
    }

    fun continueWithGoogle(activity: android.content.Context, fromCreate: Boolean) {
        if (fromCreate && !adult) {
            banner = Copy.ageGate
            return
        }
        if (BuildConfig.GOOGLE_WEB_CLIENT_ID.isBlank()) {
            banner = Copy.googleUnavailable
            return
        }
        viewModelScope.launch {
            busy = true
            banner = null
            try {
                val idToken = GoogleSignIn.requestIdToken(activity, BuildConfig.GOOGLE_WEB_CLIENT_ID)
                if (idToken != null) {
                    val session = withContext(Dispatchers.IO) { repo.api().googleSignIn(idToken) }
                    applySession(session, enterApp = !fromCreate)
                    if (fromCreate) step = 2
                }
            } catch (e: GoogleSignInException) {
                banner = when (e.message) {
                    "none" -> Copy.googleNone
                    "unconfigured" -> Copy.googleUnavailable
                    else -> Copy.googleFailed
                }
            } catch (e: Exception) {
                banner = message(e)
                if (e is ApiException && e.status == 401) signOutLocal()
            } finally {
                busy = false
            }
        }
    }

    private fun applySession(session: com.remainfaithful.app.data.AuthSession, enterApp: Boolean) {
        repo.saveSession(session)
        name = session.user.name
        email = session.user.email
        userId = session.user.id
        verified = session.user.emailVerified != false
        refreshPush()
        if (enterApp) {
            repo.onboardingDone = true
            phase = Phase.App
            refresh()
            maybeTour(false)
        }
    }

    fun forgot(address: String) {
        launchIo {
            repo.api().forgotPassword(address.trim())
            banner = Copy.forgotSent(address.trim())
            sheet = null
        }
    }

    fun resend() {
        launchIo {
            val note = repo.api().resendVerification()
            if (note.contains("already")) verified = true
            banner = note
        }
    }

    fun finishOnboarding() {
        repo.onboardingDone = true
        phase = Phase.App
        refresh()
        syncMonitoring()
        maybeTour(false)
    }

    fun signOut() {
        gate {
            launchIo {
                runCatching { repo.protection("signed_out", "Signed out. Monitoring on this phone has stopped.") }
                runCatching { repo.reportMonitoring(false) }
                signOutLocal()
            }
        }
    }

    private fun signOutLocal() {
        val context = getApplication<Application>()
        MonitorService.stop(context)
        repo.clearSession()
        phase = Phase.Login
        sheet = null
        tourOn = false
    }

    fun deleteAccount(confirm: String) {
        if (confirm != "DELETE") {
            banner = "Type DELETE to confirm"
            return
        }
        launchIo {
            repo.api().deleteMe()
            signOutLocal()
        }
    }

    fun updateProfile(newName: String, newEmail: String) {
        launchIo {
            val user = repo.api().updateMe(newName.trim(), newEmail.trim())
            repo.saveUser(user)
            name = user.name
            email = user.email
            pendingEmail = user.pendingEmail
            sheet = null
        }
    }

    fun setAlerts(on: Boolean) {
        val wasWatching = alertsOn || blockingOn
        if (!on && alertsOn) {
            alertsOn = false
            repo.alertsOn = false
            launchIo {
                repo.protection("monitoring_disabled", "Alerts were turned off.")
                repo.reportMonitoring(repo.alertsOn || repo.blockingOn)
            }
        } else if (on) {
            if (!usageGranted) {
                sheet = Sheet.Monitoring
                return
            }
            alertsOn = true
            repo.alertsOn = true
            publishDevice(turnedOn = !wasWatching)
        }
        syncMonitoring()
    }

    fun setBlocking(on: Boolean) {
        val wasWatching = alertsOn || blockingOn
        if (!on && blockingOn) {
            blockingOn = false
            repo.blockingOn = false
            launchIo {
                repo.protection("shielding_disabled", "Blocking was turned off.")
                repo.reportMonitoring(repo.alertsOn || repo.blockingOn)
            }
        } else if (on) {
            if (!usageGranted) {
                sheet = Sheet.Monitoring
                return
            }
            blockingOn = true
            repo.blockingOn = true
            if (!overlayGranted) askSystem("overlay")
            publishDevice(turnedOn = !wasWatching)
        }
        syncMonitoring()
    }

    fun setLockout(on: Boolean) {
        if (!on && lockoutOn) {
            gate {
                lockoutOn = false
                repo.lockoutOn = false
                launchIo { repo.protection("lockout_disabled", "App Lockout was turned off.") }
            }
        } else if (on) {
            lockoutOn = true
            repo.lockoutOn = true
        }
    }

    fun saveChosen(next: Set<String>) {
        val cleared = chosen.isNotEmpty() && next.isEmpty()
        val wasWatching = alertsOn || blockingOn
        chosen = next
        repo.chosenPackages = next
        val stopped = next.isEmpty() && wasWatching
        if (stopped) {
            alertsOn = false
            blockingOn = false
            repo.alertsOn = false
            repo.blockingOn = false
            syncMonitoring()
        }
        if (cleared || stopped) {
            launchIo {
                if (cleared) repo.protection("selection_cleared", "All chosen apps were removed.")
                if (stopped) {
                    repo.protection("monitoring_disabled", "Alerts stopped because no apps are selected.")
                    repo.reportMonitoring(false)
                }
            }
        }
    }

    fun loadApps() {
        viewModelScope.launch {
            val context = getApplication<Application>()
            apps = withContext(Dispatchers.Default) {
                val pm = context.packageManager
                val intent = Intent(Intent.ACTION_MAIN).addCategory(Intent.CATEGORY_LAUNCHER)
                pm.queryIntentActivities(intent, PackageManager.MATCH_ALL)
                    .map { LocalApp(it.activityInfo.packageName, it.loadLabel(pm).toString()) }
                    .filter { it.packageName != context.packageName }
                    .distinctBy { it.packageName }
                    .sortedBy { it.label.lowercase() }
            }
        }
    }

    fun openMonitoring() {
        gate { sheet = Sheet.Monitoring }
    }

    fun gate(action: () -> Unit) {
        if (pinSet) {
            pendingAfterPin = action
            sheet = Sheet.Pin
        } else {
            action()
        }
    }

    fun submitPin(pin: String, onDone: (PinCheck) -> Unit) {
        launchIo {
            val check = repo.api().verifyPin(pin)
            if (!check.success && check.reason == "wrong_pin") {
                runCatching { repo.protection("wrong_pin_attempt", "$name entered a wrong Protection PIN.") }
            }
            if (check.success) {
                val next = pendingAfterPin
                pendingAfterPin = null
                sheet = null
                next?.invoke()
            }
            onDone(check)
        }
    }

    fun setPartnerPin(relationshipId: Long, pin: String) {
        launchIo {
            repo.api().setPin(relationshipId, pin)
            repo.protection("pin_changed", "Protection PIN was changed by your accountability partner.")
            refresh()
        }
    }

    fun removePartnerPin(relationshipId: Long) {
        launchIo {
            repo.api().removePin(relationshipId)
            repo.protection("pin_removed", "Protection PIN removed by accountability partner.")
            refresh()
        }
    }

    fun invitePartner(address: String) {
        launchIo {
            val status = repo.api().invitePartner(address.trim())
            banner = when (status) {
                "accepted" -> "You are connected."
                "invited" -> "Invite sent. They can join when they create an account."
                else -> "Invite sent. They accept from Settings."
            }
            refresh()
        }
    }

    fun acceptPartner(id: Long) = launchIo { repo.api().acceptPartnerInvite(id); refresh() }
    fun declinePartner(id: Long) = launchIo { repo.api().declinePartnerInvite(id); refresh() }
    fun acceptGroup(id: Long) = launchIo { repo.api().acceptGroupInvite(id); refresh() }
    fun declineGroup(id: Long) = launchIo { repo.api().declineGroupInvite(id); refresh() }

    fun acceptLink(token: String, type: String?) {
        launchIo {
            if (type == "group") repo.api().acceptGroupToken(token) else repo.api().acceptPartnerToken(token)
            banner = "Invite accepted."
            refresh()
        }
    }

    fun createGroup(groupName: String) {
        launchIo {
            val created = repo.api().createGroup(groupName.trim(), "")
            repo.primaryGroupId = created.id
            refresh()
        }
    }

    fun inviteMember(address: String) {
        val current = group ?: return
        val count = current.members.size
        if (!GroupPolicy.canInvite(count)) {
            banner = GroupPolicy.fullMessage(count)
            return
        }
        launchIo {
            val status = repo.api().emailInvite(current.id, address.trim())
            banner = when (status) {
                "added" -> "Added."
                "already_member" -> "Already in the group."
                else -> "Invite sent."
            }
            refresh()
        }
    }

    fun updateCovenant(text: String) {
        val current = group ?: return
        launchIo {
            repo.api().updateGroup(current.id, null, text)
            refresh()
        }
    }

    fun encourage(member: GroupMember) {
        val current = group ?: return
        launchIo {
            repo.api().encourage(current.id, member.userId)
            banner = "Encouragement sent."
        }
    }

    fun leaveGroup() {
        val current = group ?: return
        launchIo {
            repo.api().leaveGroup(current.id)
            repo.primaryGroupId = 0
            refresh()
        }
    }

    fun leaveAll() {
        launchIo {
            repo.api().leaveAll()
            repo.primaryGroupId = 0
            refresh()
        }
    }

    fun removePartner(rel: RemoteRel) {
        launchIo { repo.api().deleteRelationship(rel.id); refresh() }
    }

    fun makePrimary(rel: RemoteRel) {
        launchIo { repo.api().setPrimary(rel.id); refresh() }
    }

    fun panic() {
        launchIo {
            repo.supportTaps = repo.supportTaps + 1
            repo.api().panic()
            repo.sendSignals(usageGranted)
            banner = "Your partners have been told."
            sheet = null
        }
    }

    fun feedback(type: String, summary: String, details: String, includeDevice: Boolean) {
        val error = com.remainfaithful.core.FeedbackRules.error(type, summary, details)
        if (error != null) {
            banner = error
            return
        }
        launchIo {
            val fields = linkedMapOf("type" to type, "summary" to summary.trim(), "details" to details.trim())
            if (includeDevice) {
                fields["app_version"] = com.remainfaithful.app.BuildConfig.VERSION_NAME
                fields["os_version"] = android.os.Build.VERSION.RELEASE ?: ""
                fields["device_model"] = android.os.Build.MODEL ?: ""
            }
            repo.api().feedback(fields)
            banner = "Sent. Thank you."
            sheet = null
        }
    }

    fun openAlert(event: UiEvent) {
        selectedAlert = event
        sheet = Sheet.Alert
        val id = event.alertId ?: return
        viewModelScope.launch { runCatching { withContext(Dispatchers.IO) { repo.api().markAlertSeen(id) } } }
    }

    fun discuss(event: UiEvent) {
        val id = event.alertId ?: return
        launchIo {
            repo.api().markDiscussed(id)
            selectedAlert = event.copy(discussed = true)
            refresh()
        }
    }

    fun markSeen() {
        viewModelScope.launch { runCatching { withContext(Dispatchers.IO) { repo.api().markAlertsSeen() } } }
    }

    fun replayTour() {
        tourForced = true
        maybeTour(true)
    }

    fun skipTour() {
        tourOn = false
        repo.tourVersion = TourContent.version
    }

    fun nextTour() {
        if (tourIndex + 1 >= tourSteps.size) skipTour() else tourIndex++
    }

    fun backTour() {
        if (tourIndex > 0) tourIndex--
    }

    private fun maybeTour(force: Boolean) {
        if (TourContent.shouldAutoPresent(repo.tourVersion, force)) {
            tourIndex = 0
            tourOn = true
            tab = 0
        }
    }

    private fun refreshPush() {
        runCatching {
            FirebaseMessaging.getInstance().token.addOnCompleteListener { task ->
                if (!task.isSuccessful) return@addOnCompleteListener
                val value = task.result ?: return@addOnCompleteListener
                viewModelScope.launch(Dispatchers.IO) {
                    runCatching { repo.registerPush(value) }
                }
            }
        }
    }

    private fun launchIo(block: suspend () -> Unit) {
        viewModelScope.launch {
            busy = true
            banner = null
            try {
                withContext(Dispatchers.IO) { block() }
            } catch (e: Exception) {
                if (e is ApiException && DevicePolicy.removed(e.status, e.code)) adoptRemoved()
                banner = message(e)
                if (e is ApiException && e.status == 401) signOutLocal()
            } finally {
                busy = false
            }
        }
    }

    private fun message(e: Exception): String = when {
        e is ApiException && DevicePolicy.removed(e.status, e.code) -> Copy.deviceRemoved
        e is ApiException -> e.message ?: Copy.networkError
        else -> Copy.networkError
    }

    private fun publishDevice(turnedOn: Boolean) {
        launchIo {
            if (turnedOn) repo.reclaim()
            val watching = repo.alertsOn || repo.blockingOn
            repo.reportMonitoring(watching)
            if (watching) repo.sendHeartbeat()
            runCatching { reloadDevices() }
        }
    }

    private fun adoptRemoved() {
        alertsOn = false
        blockingOn = false
        deepScanOn = false
        repo.alertsOn = false
        repo.blockingOn = false
        repo.deepScanRunning = false
        repo.removedNotice = false
        syncMonitoring()
    }

    private fun adoptRemovedIfNeeded() {
        if (!repo.removedNotice) return
        adoptRemoved()
        banner = Copy.deviceRemoved
    }

    private suspend fun reloadDevices() {
        val remote = try {
            withContext(Dispatchers.IO) { repo.listDevices() }
        } catch (e: ApiException) {
            if (e.status == 404) devicesUnavailable = true
            return
        } catch (_: Exception) {
            return
        }
        devices = deviceRows(remote)
        devicesUnavailable = false
    }

    private fun deviceRows(remote: List<RemoteDevice>): List<DeviceRow> {
        val mine = repo.deviceId
        return remote.mapNotNull { item ->
            if (!DevicePolicy.validId(item.id)) return@mapNotNull null
            val label = DevicePolicy.label(item.kind) ?: return@mapNotNull null
            DeviceRow(
                id = item.id,
                label = label,
                monitoring = item.monitoring,
                lastCheckIn = DevicePolicy.formatCheckIn(item.lastCheckIn).ifBlank { Copy.noCheckIn },
                current = item.id == mine,
            )
        }
    }

    fun askRemoveDevice(row: DeviceRow) {
        pendingDevice = row
        gate { sheet = Sheet.RemoveDevice }
    }

    fun confirmRemoveDevice() {
        val row = pendingDevice ?: return
        val wasWatching = row.current && (alertsOn || blockingOn)
        if (row.current) {
            alertsOn = false
            blockingOn = false
            repo.alertsOn = false
            repo.blockingOn = false
            syncMonitoring()
        }
        launchIo {
            try {
                repo.removeDevice(row.id)
            } catch (e: Exception) {
                if (wasWatching) {
                    runCatching { repo.protection("monitoring_disabled", "Alerts were turned off.") }
                }
                throw e
            }
            pendingDevice = null
            sheet = null
            reloadDevices()
        }
    }

    fun greeting(): String {
        val hour = java.time.LocalTime.now().hour
        return when {
            hour < 12 -> Copy.morning
            hour < 17 -> Copy.afternoon
            else -> Copy.evening
        }
    }

    fun verse() = com.remainfaithful.core.Verses.forDayOfYear(LocalDate.now().dayOfYear)

    fun inviteLink(): String {
        val id = group?.id ?: repo.primaryGroupId
        return "https://www.remainfaithful.com/join/$id"
    }
}

private fun RemoteEvent.toUi() = UiEvent(
    key = "e$id",
    alertId = null,
    category = category,
    severity = severity,
    summary = summary,
    timeLabel = timeAgo(timestamp),
    discussed = false,
)

private fun RemoteAlert.toUi(): UiEvent {
    val category = event?.category ?: protectionType ?: kind ?: "protection"
    val summary = event?.summary ?: summary ?: CategoryLabels.label(category)
    return UiEvent(
        key = "a$id",
        alertId = id,
        category = category,
        severity = event?.severity ?: "concerning",
        summary = summary,
        timeLabel = timeAgo(event?.timestamp ?: createdAt),
        discussed = discussed,
    )
}

private fun timeAgo(raw: String): String {
    val instant = parseTime(raw) ?: return raw
    val mins = Duration.between(instant, Instant.now()).toMinutes().coerceAtLeast(0)
    return when {
        mins < 60 -> "${mins}m ago"
        mins < 1440 -> "${mins / 60}h ago"
        else -> "${mins / 1440}d ago"
    }
}

private fun parseTime(raw: String): Instant? {
    runCatching { return Instant.parse(raw) }
    runCatching { return OffsetDateTime.parse(raw).toInstant() }
    runCatching { return java.time.LocalDateTime.parse(raw.replace(" ", "T")).atZone(ZoneId.systemDefault()).toInstant() }
    return null
}
