# Partner alert trigger

Paste this sentence into the App Store listing and the app. Do not paraphrase it.

> If Alerts is on, your partner is told when a chosen app or category is used for 1 minute, then again at 5, 15, and 30 minutes that day. If Blocking is on, opening a chosen app shows Apple's Screen Time shield, Apple's adult web filter is on, and your partner gets an alert that a blocked app was opened. That alert does not include the app name.

## What the code does

Alerts and Blocking are separate switches in `RemainFaithful/ScreenTimeMonitoringView.swift`. They do not do the same thing.

**Alerts (monitored apps and categories).** `ActivitySelectionManager.startMonitoring` registers a daily DeviceActivity schedule (`00:00` to `23:59`, repeating). Each selected set has four thresholds, at 1, 5, 15, and 30 minutes. Chosen apps use `com.remainfaithful.appusage` plus `.5`, `.15`, and `.30`. Categories use `com.remainfaithful.categoryusage` with the same suffixes. `DeviceActivityMonitorExtension.eventDidReachThreshold` runs when cumulative usage crosses that event's threshold. Each event fires once per daily interval. It posts `POST /events` with category `app_usage`. Severity is concerning at 1 and 5 minutes and severe at 15 and 30. The summary is a fixed string such as `Monitored app used for 1+ minute` or `Monitored app category used for 15+ minutes`. `backend/internal/handler/events.go` accepts those summaries and no others. The push names the threshold that was crossed.

If both apps and categories are selected, each set sends its own alerts. Desktop and Android use the same minute marks and the same summaries.

**Blocking (shielded apps).** `ActivitySelectionManager.applyShielding` puts the chosen apps and categories on `ManagedSettingsStore.shield` and turns on Apple's adult web filter. `RemainFaithfulShieldConfig/ShieldConfigurationExtension.swift` draws the Screen Time shield and posts `Blocked app opened`. That alert does not include the app name. A partner is also told if Blocking is turned off (`shielding_disabled`), which is a protection notice.

Optional Deep Scan can store a separate content alert when it flags a frame. That path is not this sentence. Deep Scan cannot see DRM streaming.

## Where this sentence lives on the site

`website/src/lib/product-copy.ts` exports `ALERT_TRIGGER` with this exact wording. Pages that describe the trigger use that constant.
