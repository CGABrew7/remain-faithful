package com.remainfaithful.app.push

import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import com.remainfaithful.app.data.Repository
import com.remainfaithful.app.monitoring.Notify
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch

class RfFirebaseService : FirebaseMessagingService() {
    override fun onNewToken(token: String) {
        CoroutineScope(Dispatchers.IO).launch {
            runCatching { Repository.get(this@RfFirebaseService).registerPush(token) }
        }
    }

    override fun onMessageReceived(message: RemoteMessage) {
        val data = message.data
        val title = message.notification?.title ?: data["title"] ?: "Remain Faithful"
        val body = message.notification?.body ?: data["body"].orEmpty()
        val type = data["notification_type"].orEmpty()
        if (title.isNotBlank() || body.isNotBlank()) {
            Notify.alert(this, title, body, type.ifBlank { "CONTENT_FLAGGED" })
        }
    }
}
