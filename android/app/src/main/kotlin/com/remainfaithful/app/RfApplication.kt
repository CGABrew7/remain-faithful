package com.remainfaithful.app

import android.app.Application
import com.remainfaithful.app.monitoring.createChannels

class RfApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        createChannels(this)
    }
}
