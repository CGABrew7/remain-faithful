# Partner alert trigger

Paste this sentence into the App Store listing and the app. Do not paraphrase it.

> If Alerts is on, your partner gets one alert that day after a chosen app or category is used for a minute. If Blocking is on and you open an app you asked to block, your partner gets an alert that you opened an app you asked to block, and Apple's Screen Time shield appears. The app name is never sent.

## What the code does

Alerts and Blocking are separate switches in `RemainFaithful/ScreenTimeMonitoringView.swift`. They do not do the same thing.

**Alerts (monitored apps and categories).** `ActivitySelectionManager.startMonitoring` registers a daily DeviceActivity schedule (`00:00` to `23:59`, repeating). Each selected set has a threshold of `DateComponents(minute: 1)`: one event for chosen apps (`com.remainfaithful.appusage`) and, if categories are selected, a second event for chosen categories (`com.remainfaithful.categoryusage`). `DeviceActivityMonitorExtension.eventDidReachThreshold` runs when that cumulative usage crosses one minute. The extension comments that this fires once per event per daily interval. It posts `POST /events` with category `app_usage`, severity `concerning`, and summary `Monitored app used for 1+ minute` or `Monitored app category used for 1+ minute`. `backend/internal/handler/events.go` `CreateEvent` stores that event and inserts a partner alert. The push body for `app_usage` says the person spent 1+ minute today in an app they are monitoring.

If both apps and categories are selected, each set can send its own alert that day. The sentence says "one alert that day" for a chosen app or category, which is the once-per-day rule on each set.

**Blocking (shielded apps).** `ActivitySelectionManager.applyShielding` puts the chosen apps and categories on `ManagedSettingsStore.shield`. When Blocking is on and the user opens an app they asked to block, iOS shows that shield. `RemainFaithfulShieldConfig/ShieldConfigurationExtension.swift` draws it (title, subtitle, and an "I need help" button) and calls `noteBlockedAppOpened`. That method posts `POST /events` with category `app_usage`, severity `concerning`, and summary `Blocked app opened`. The request does not include the app name or a bundle id. Another open inside 120 seconds does not post again. `backend/internal/handler/events.go` `CreateEvent` stores the event and inserts a partner alert. The push body is the person's name plus `opened an app they asked to block.` If that post does not succeed, the same row stays in the shared queue so a later send can deliver it. Turning Blocking off still sends `shielding_disabled` from `backend/internal/handler/protection.go`. That notice is separate from the open alert.

Optional Deep Scan can store a separate content alert when it flags a frame. That path is not this sentence. Deep Scan cannot see DRM streaming.

## Where this sentence lives on the site

`website/src/lib/product-copy.ts` exports `ALERT_TRIGGER` with this exact wording. Pages that describe the trigger use that constant. The App Guide Block Selected Apps line and the software schema feature list repeat the blocking meaning in their own words.
