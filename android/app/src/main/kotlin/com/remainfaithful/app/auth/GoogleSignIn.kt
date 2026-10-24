package com.remainfaithful.app.auth

import android.content.Context
import androidx.credentials.CredentialManager
import androidx.credentials.GetCredentialRequest
import androidx.credentials.exceptions.GetCredentialCancellationException
import androidx.credentials.exceptions.NoCredentialException
import com.google.android.libraries.identity.googleid.GetSignInWithGoogleOption
import com.google.android.libraries.identity.googleid.GoogleIdTokenCredential

class GoogleSignInException(message: String) : Exception(message)

object GoogleSignIn {
    /**
     * Opens the Google account picker and returns an ID token for POST /auth/google.
     * A null result means the person closed the picker.
     */
    suspend fun requestIdToken(context: Context, serverClientId: String): String? {
        if (serverClientId.isBlank()) throw GoogleSignInException("unconfigured")
        val request = GetCredentialRequest.Builder()
            .addCredentialOption(GetSignInWithGoogleOption.Builder(serverClientId).build())
            .build()
        try {
            val result = CredentialManager.create(context).getCredential(context, request)
            val credential = result.credential
            if (credential.type != GoogleIdTokenCredential.TYPE_GOOGLE_ID_TOKEN_CREDENTIAL) {
                throw GoogleSignInException("missing")
            }
            return GoogleIdTokenCredential.createFrom(credential.data).idToken
        } catch (_: GetCredentialCancellationException) {
            return null
        } catch (_: NoCredentialException) {
            throw GoogleSignInException("none")
        }
    }
}
