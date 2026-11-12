# Partner alert trigger

Paste this sentence into the App Store listing and the app. Do not paraphrase it.

> If Alerts is on, your partner is told when a chosen app or category is used for 1 minute, then again at 5, 15, and 30 minutes that day. If Blocking is on, opening a chosen app shows Apple's Screen Time shield and does not itself send an alert.

## What the code does

Alerts and Blocking are separate switches in `RemainFaithful/ScreenTimeMonitoringView.swift`. They do not do the same thing.

**Alerts (monitored apps and categories).** `ActivitySelectionManager.startMonitoring` registers a daily DeviceActivity schedule (`00:00` to `23:59`, repeating). Each selected set has four thresholds, at 1, 5, 15, and 30 minutes. Chosen apps use `com.remainfaithful.appusage` plus `.5`, `.15`, and `.30`. Categories use `com.remainfaithful.categoryusage` with the same suffixes. `DeviceActivityMonitorExtension.eventDidReachThreshold` runs when cumulative usage crosses that event's threshold. Each event fires once per daily interval. It posts `POST /events` with category `app_usage`, severity `concerning`, and a fixed summary such as `Monitored app used for 1+ minute` or `Monitored app category used for 15+ minutes`. `backend/internal/handler/events.go` accepts those summaries and no others. The push names the threshold that was crossed.

If both apps and categories are selected, each set sends its own alerts. Desktop and Android use the same minute marks and the same summaries.

**Blocking (shielded apps).** `ActivitySelectionManager.applyShielding` puts the chosen apps and categories on `ManagedSettingsStore.shield`. `RemainFaithfulShieldConfig/ShieldConfigurationExtension.swift` only draws the Screen Time shield (title, subtitle, and an "I need help" button). It does not post an event. There is no ShieldAction extension. Opening a blocked app does not create a partner alert. A partner is told if Blocking is turned off (`shielding_disabled` in `backend/internal/handler/protection.go`), which is a protection notice, not an "app opened" alert.

Optional Deep Scan can store a separate content alert when it flags a frame. That path is not this sentence. Deep Scan cannot see DRM streaming.

## Where this sentence lives on the site

`website/src/lib/product-copy.ts` exports `ALERT_TRIGGER` with this exact wording. Pages that describe the trigger use that constant.
