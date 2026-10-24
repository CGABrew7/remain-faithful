package com.remainfaithful.app.data

import com.remainfaithful.core.PinCheck
import com.remainfaithful.core.friendlyServerMessage
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONArray
import org.json.JSONObject
import java.util.concurrent.TimeUnit

class ApiException(val status: Int, message: String, val code: String?) : Exception(message)

data class RemoteUser(
    val id: Long,
    val name: String,
    val email: String,
    val createdAt: String?,
    val emailVerified: Boolean?,
    val pendingEmail: String?,
)

data class AuthSession(val token: String, val user: RemoteUser)

data class RemoteEvent(
    val id: Long,
    val category: String,
    val severity: String,
    val summary: String,
    val timestamp: String,
)

data class RemoteAlert(
    val id: Long,
    val seen: Boolean,
    val createdAt: String,
    val kind: String?,
    val protectionType: String?,
    val summary: String?,
    val event: RemoteEvent?,
    val discussed: Boolean,
)

data class RemoteStats(val streakDays: Int, val bestStreak: Int, val weekStates: List<String>)

data class RemoteRel(
    val id: Long,
    val userId: Long,
    val partnerId: Long,
    val status: String,
    val isPrimary: Boolean,
    val pinSet: Boolean,
    val partnerName: String,
    val partnerEmail: String,
)

data class PartnerInvite(val id: Long, val inviterName: String, val createdAt: String)
data class GroupInvite(val id: Long, val groupId: Long, val groupName: String, val inviterName: String)

data class GroupMember(
    val userId: Long,
    val role: String,
    val name: String,
    val flagsLast30: Int,
    val streakDays: Int,
)

data class RemoteGroup(
    val id: Long,
    val name: String,
    val covenant: String,
    val members: List<GroupMember>,
)

class ApiClient(
    private val baseUrl: String,
    private val tokenProvider: () -> String?,
    private val onToken: (String) -> Unit,
) {
    private val http = OkHttpClient.Builder()
        .callTimeout(20, TimeUnit.SECONDS)
        .build()
    private val jsonType = "application/json; charset=utf-8".toMediaType()

    fun register(name: String, email: String, password: String): RemoteUser {
        val body = post("/auth/register", mapOf("name" to name, "email" to email, "password" to password), auth = false)
        return body.toUser()
    }

    fun login(email: String, password: String): AuthSession {
        val body = post("/auth/login", mapOf("email" to email, "password" to password), auth = false)
        val session = AuthSession(body.getString("token"), body.getJSONObject("user").toUser())
        onToken(session.token)
        return session
    }

    fun googleSignIn(idToken: String): AuthSession {
        val body = post("/auth/google", com.remainfaithful.core.Payloads.googleIdToken(idToken), auth = false)
        val session = AuthSession(body.getString("token"), body.getJSONObject("user").toUser())
        onToken(session.token)
        return session
    }

    fun forgotPassword(email: String) {
        post("/auth/forgot-password", mapOf("email" to email), auth = false)
    }

    fun refresh() {
        val body = post("/auth/refresh", emptyMap())
        val token = body.optString("token")
        if (token.isNotBlank()) onToken(token)
    }

    fun resendVerification(): String {
        val body = post("/auth/resend-verification", emptyMap())
        return when (body.optString("status")) {
            "already_verified" -> "Your email is already confirmed."
            "sent" -> "Check your inbox for a confirmation email."
            else -> "Check your inbox."
        }
    }

    fun me(): RemoteUser = get("/users/me").toUser()

    fun updateMe(name: String, email: String): RemoteUser =
        put("/users/me", mapOf("name" to name, "email" to email)).toUser()

    fun deleteMe() {
        delete("/users/me")
    }

    fun stats(): RemoteStats {
        val body = get("/users/me/stats")
        val states = when {
            body.has("week_states") -> body.getJSONArray("week_states").toStates()
            body.has("week") -> body.getJSONArray("week").toStates()
            else -> emptyList()
        }
        return RemoteStats(body.optInt("streak_days"), body.optInt("best_streak"), states)
    }

    fun events(): List<RemoteEvent> = getArray("/events").toEvents()

    fun alerts(): List<RemoteAlert> = getArray("/alerts").toAlerts()

    fun markAlertsSeen() {
        post("/alerts/mark-seen", emptyMap())
    }

    fun markAlertSeen(id: Long) {
        try {
            post("/alerts/$id/seen", emptyMap())
        } catch (e: ApiException) {
            if (e.status != 404) throw e
        }
    }

    fun markDiscussed(id: Long) {
        send("PATCH", "/alerts/$id/discussed", null)
    }

    fun createEvent(fields: Map<String, String>): RemoteEvent = post("/events", fields).toEvent()

    fun heartbeat(screen: String) {
        post("/heartbeat", mapOf("screen" to screen))
    }

    fun signals(body: Map<String, Any>) {
        try {
            postAny("/users/me/signals", body)
        } catch (_: ApiException) {
            // Older servers ignore this.
        }
    }

    fun registerDevice(body: Map<String, String>) {
        post("/users/device-token", body)
    }

    fun protection(type: String, detail: String) {
        post("/protection/alerts", mapOf("type" to type, "detail" to detail))
    }

    fun pinSet(): Boolean = get("/protection/pin").optBoolean("is_set")

    fun verifyPin(pin: String): PinCheck {
        val raw = send("POST", "/protection/pin/verify", mapOf("pin" to pin), acceptErrorBody = true)
        val body = JSONObject(raw.ifBlank { "{}" })
        if (body.has("success") || body.has("reason")) {
            return PinCheck.from(body.optBoolean("success"), body.optString("reason").ifBlank { null })
        }
        throw ApiException(400, friendlyServerMessage(body.optString("error")), body.optString("code").ifBlank { null })
    }

    fun setPin(relationshipId: Long, pin: String) {
        post("/relationships/$relationshipId/pin", mapOf("pin" to pin))
    }

    fun removePin(relationshipId: Long) {
        delete("/relationships/$relationshipId/pin")
    }

    fun relationships(): List<RemoteRel> = getArray("/relationships").toRels()

    fun invitePartner(email: String): String =
        post("/relationships/invite", mapOf("email" to email)).optString("status")

    fun acceptPartnerInvite(id: Long) {
        post("/relationships/invites/$id/accept", emptyMap())
    }

    fun declinePartnerInvite(id: Long) {
        post("/relationships/invites/$id/decline", emptyMap())
    }

    fun acceptPartnerToken(token: String) {
        post("/relationships/accept-invite", mapOf("token" to token))
    }

    fun pendingPartners(): List<PartnerInvite> = getArray("/relationships/invites/pending").let { arr ->
        (0 until arr.length()).map { i ->
            val o = arr.getJSONObject(i)
            PartnerInvite(o.optLong("id"), o.optString("inviter_name"), o.optString("created_at"))
        }
    }

    fun deleteRelationship(id: Long) {
        delete("/relationships/$id")
    }

    fun setPrimary(id: Long) {
        put("/relationships/$id/primary", emptyMap())
    }

    fun groups(): List<RemoteGroup> = getArray("/groups").toGroups(withMembers = false)

    fun group(id: Long): RemoteGroup = get("/groups/$id").toGroup()

    fun createGroup(name: String, covenant: String): RemoteGroup =
        post("/groups", mapOf("name" to name, "covenant" to covenant)).toGroup()

    fun updateGroup(id: Long, name: String?, covenant: String?) {
        val fields = linkedMapOf<String, String>()
        if (name != null) fields["name"] = name
        if (covenant != null) fields["covenant"] = covenant
        send("PATCH", "/groups/$id", fields)
    }

    fun emailInvite(groupId: Long, email: String): String =
        post("/groups/$groupId/email-invite", mapOf("email" to email)).optString("status")

    fun leaveGroup(id: Long) {
        delete("/groups/$id/members/me")
    }

    fun leaveAll() {
        post("/groups/leave-all", emptyMap())
    }

    fun encourage(groupId: Long, userId: Long) {
        postAny("/groups/$groupId/encouragement", mapOf("user_id" to userId))
    }

    fun pendingGroups(): List<GroupInvite> = getArray("/groups/invites/pending").let { arr ->
        (0 until arr.length()).map { i ->
            val o = arr.getJSONObject(i)
            GroupInvite(o.optLong("id"), o.optLong("group_id"), o.optString("group_name"), o.optString("inviter_name"))
        }
    }

    fun acceptGroupInvite(id: Long) {
        post("/groups/invites/$id/accept", emptyMap())
    }

    fun declineGroupInvite(id: Long) {
        post("/groups/invites/$id/decline", emptyMap())
    }

    fun acceptGroupToken(token: String) {
        post("/groups/accept-invite", mapOf("token" to token))
    }

    fun panic() {
        post("/panic", emptyMap())
    }

    fun feedback(fields: Map<String, String>) {
        post("/feedback", fields)
    }

    private fun get(path: String): JSONObject = JSONObject(send("GET", path, null))
    private fun getArray(path: String): JSONArray = JSONArray(send("GET", path, null))
    private fun post(path: String, body: Map<String, String>, auth: Boolean = true): JSONObject =
        JSONObject(send("POST", path, body, auth).ifBlank { "{}" })

    private fun postAny(path: String, body: Map<String, Any>): JSONObject =
        JSONObject(sendAny("POST", path, body).ifBlank { "{}" })

    private fun put(path: String, body: Map<String, String>): JSONObject =
        JSONObject(send("PUT", path, body).ifBlank { "{}" })

    private fun delete(path: String) {
        send("DELETE", path, null)
    }

    private fun send(method: String, path: String, body: Map<String, String>?, auth: Boolean = true, acceptErrorBody: Boolean = false): String {
        maybeRefresh(auth, path)
        val builder = Request.Builder().url(baseUrl.trimEnd('/') + path)
        if (auth) {
            val token = tokenProvider() ?: throw ApiException(401, "Not signed in. Please log in again.", null)
            builder.header("Authorization", "Bearer $token")
        }
        val payload = if (body == null) null else jsonOf(body).toRequestBody(jsonType)
        builder.method(method, if (method == "GET" || method == "DELETE") null else payload ?: "{}".toRequestBody(jsonType))
        return execute(builder.build(), acceptErrorBody)
    }

    private fun sendAny(method: String, path: String, body: Map<String, Any>): String {
        maybeRefresh(true, path)
        val token = tokenProvider() ?: throw ApiException(401, "Not signed in. Please log in again.", null)
        val json = JSONObject()
        body.forEach { (k, v) ->
            when (v) {
                is Boolean -> json.put(k, v)
                is Int -> json.put(k, v)
                is Long -> json.put(k, v)
                else -> json.put(k, v.toString())
            }
        }
        val request = Request.Builder()
            .url(baseUrl.trimEnd('/') + path)
            .header("Authorization", "Bearer $token")
            .method(method, json.toString().toRequestBody(jsonType))
            .build()
        return execute(request, false)
    }

    private fun maybeRefresh(auth: Boolean, path: String) {
        if (!auth || path.endsWith("/auth/refresh")) return
        val token = tokenProvider() ?: return
        if (!com.remainfaithful.core.JwtClock.shouldRefresh(token, System.currentTimeMillis() / 1000)) return
        try {
            val request = Request.Builder()
                .url(baseUrl.trimEnd('/') + "/auth/refresh")
                .header("Authorization", "Bearer $token")
                .post("{}".toRequestBody(jsonType))
                .build()
            val raw = execute(request, false)
            val next = JSONObject(raw).optString("token")
            if (next.isNotBlank()) onToken(next)
        } catch (_: Exception) {
            // A failed refresh leaves the current token in place.
        }
    }

    private fun execute(request: Request, acceptErrorBody: Boolean): String {
        try {
            http.newCall(request).execute().use { response ->
                val raw = response.body?.string().orEmpty()
                if (response.isSuccessful) return raw
                val parsed = runCatching { JSONObject(raw) }.getOrNull()
                if (acceptErrorBody && parsed != null && (parsed.has("success") || parsed.has("reason"))) return raw
                val message = friendlyServerMessage(parsed?.optString("error") ?: "HTTP ${response.code}")
                throw ApiException(response.code, message, parsed?.optString("code")?.ifBlank { null })
            }
        } catch (e: ApiException) {
            throw e
        } catch (e: Exception) {
            throw ApiException(0, com.remainfaithful.core.Copy.networkError, null)
        }
    }
}

private fun jsonOf(body: Map<String, String>): String {
    val o = JSONObject()
    body.forEach { (k, v) -> o.put(k, v) }
    return o.toString()
}

private fun JSONObject.toUser() = RemoteUser(
    id = optLong("id"),
    name = optString("name"),
    email = optString("email"),
    createdAt = optString("created_at").ifBlank { null },
    emailVerified = if (has("email_verified")) optBoolean("email_verified") else null,
    pendingEmail = optString("pending_email").ifBlank { null },
)

private fun JSONObject.toEvent() = RemoteEvent(
    id = optLong("id"),
    category = optString("category"),
    severity = optString("severity"),
    summary = optString("summary"),
    timestamp = optString("timestamp"),
)

private fun JSONArray.toEvents() = (0 until length()).map { getJSONObject(it).toEvent() }

private fun JSONArray.toAlerts(): List<RemoteAlert> = (0 until length()).map { i ->
    val o = getJSONObject(i)
    val type = listOf(o.optString("alert_type"), o.optString("protection_type"), o.optString("type"))
        .firstOrNull { it.isNotBlank() && it != "null" }
    RemoteAlert(
        id = o.optLong("id"),
        seen = o.optBoolean("seen"),
        createdAt = o.optString("created_at"),
        kind = o.optString("kind").ifBlank { null },
        protectionType = type,
        summary = o.optString("summary").ifBlank { null },
        event = if (o.has("event") && !o.isNull("event")) o.getJSONObject("event").toEvent() else null,
        discussed = o.optBoolean("discussed"),
    )
}

private fun JSONArray.toStates(): List<String> = (0 until length()).map { i ->
    when (val v = get(i)) {
        is Boolean -> if (v) "clean" else "flagged"
        else -> v.toString()
    }
}

private fun JSONArray.toRels(): List<RemoteRel> = (0 until length()).map { i ->
    val o = getJSONObject(i)
    val partner = o.optJSONObject("partner")
    RemoteRel(
        id = o.optLong("id"),
        userId = o.optLong("user_id"),
        partnerId = o.optLong("partner_id"),
        status = o.optString("status"),
        isPrimary = o.optBoolean("is_primary"),
        pinSet = o.optBoolean("pin_set"),
        partnerName = partner?.optString("name").orEmpty(),
        partnerEmail = partner?.optString("email").orEmpty(),
    )
}

private fun JSONArray.toGroups(withMembers: Boolean) = (0 until length()).map { getJSONObject(it).toGroup() }

private fun JSONObject.toGroup(): RemoteGroup {
    val members = if (has("members") && !isNull("members")) {
        val arr = getJSONArray("members")
        (0 until arr.length()).map { i ->
            val m = arr.getJSONObject(i)
            val user = m.optJSONObject("user")
            GroupMember(
                userId = m.optLong("user_id"),
                role = m.optString("role"),
                name = user?.optString("name").orEmpty(),
                flagsLast30 = m.optInt("flags_last_30"),
                streakDays = m.optInt("streak_days"),
            )
        }
    } else {
        emptyList()
    }
    return RemoteGroup(optLong("id"), optString("name"), optString("covenant"), members)
}
