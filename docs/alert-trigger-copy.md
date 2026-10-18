# Partner alert trigger

Paste this sentence into the App Store listing and the app. Do not paraphrase it.

> If Alerts is on, your partner gets one alert that day after a chosen app or category is used for a minute. If Blocking is on, opening a chosen app shows Apple's Screen Time shield and does not itself send an alert.

## What the code does

Alerts and Blocking are separate switches in `RemainFaithful/ScreenTimeMonitoringView.swift`. They do not do the same thing.

**Alerts (monitored apps and categories).** `ActivitySelectionManager.startMonitoring` registers a daily DeviceActivity schedule (`00:00` to `23:59`, repeating). Each selected set has a threshold of `DateComponents(minute: 1)`: one event for chosen apps (`com.remainfaithful.appusage`) and, if categories are selected, a second event for chosen categories (`com.remainfaithful.categoryusage`). `DeviceActivityMonitorExtension.eventDidReachThreshold` runs when that cumulative usage crosses one minute. The extension comments that this fires once per event per daily interval. It posts `POST /events` with category `app_usage`, severity `concerning`, and summary `Monitored app used for 1+ minute` or `Monitored app category used for 1+ minute`. `backend/internal/handler/events.go` `CreateEvent` stores that event and inserts a partner alert. The push body for `app_usage` says the person spent 1+ minute today in an app they are monitoring.

If both apps and categories are selected, each set can send its own alert that day. The sentence says "one alert that day" for a chosen app or category, which is the once-per-day rule on each set.

**Blocking (shielded apps).** `ActivitySelectionManager.applyShielding` puts the chosen apps and categories on `ManagedSettingsStore.shield`. `RemainFaithfulShieldConfig/ShieldConfigurationExtension.swift` only draws the Screen Time shield (title, subtitle, and an "I need help" button). It does not post an event. There is no ShieldAction extension. Opening a blocked app does not create a partner alert. A partner is told if Blocking is turned off (`shielding_disabled` in `backend/internal/handler/protection.go`), which is a protection notice, not an "app opened" alert.

Optional Deep Scan can store a separate content alert when it flags a frame. That path is not this sentence. Deep Scan cannot see DRM streaming.

## Where this sentence lives on the site

`website/src/lib/product-copy.ts` exports `ALERT_TRIGGER` with this exact wording. Pages that describe the trigger use that constant.
