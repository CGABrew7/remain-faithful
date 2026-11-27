# Alert trigger copy (1.0.4)

> If Alerts is on, an on-device schedule stays running and your partner is told when you turn it off. If Blocking is on, Apple's adult web filter turns on and opening a chosen app shows Apple's Screen Time shield. You can open Remain Faithful from that shield and tap Tell my partner to send the time and the words a blocked app was opened.

**Alerts.** `ActivitySelectionManager.startMonitoring` registers a daily DeviceActivity schedule with no usage thresholds. The monitor extension applies or clears shields locally. It does not call the API.

**Blocking.** Managed Settings shields the selection and turns on Apple's adult web filter. The shield button opens Remain Faithful. `POST /events` with summary `Blocked app opened` runs only after Tell my partner, from the main app.

**Session Shield.** Optional. Chosen apps stay shielded unless a fresh Deep Scan broadcast flag is set. Stopping the broadcast re-applies the shields.

**Deep Scan.** User-started ReplayKit broadcast. `POST /events` carries category, severity, a fixed summary, timestamp, tier, and confidence_bucket.

**Check-in.** `POST /check-ins` after the user taps Send.

**Coach nudge.** `POST /protection/alerts` with type `coach_nudge` and the fixed detail `Your partner would welcome a word from you today.` The main app sends it after the person opts in and taps Check in now. Not now does not send it. One accepted notice per local day. Screen Time extensions do not send it.

**Android.** The Android build can still send the minute summaries and `Blocked app opened`. Those summaries stay accepted by the API.
