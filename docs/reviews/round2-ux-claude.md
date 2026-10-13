# Remain Faithful — UX / UI / Friction Review (Round 2, Claude)

**Reviewer:** Claude (Opus 5)
**Date:** 2026-10-13
**Commit reviewed:** `784998f`
**Perspective:** Two personas walked end to end — **Dave**, a 41-year-old non-technical first-time user who heard about the app at men's group, and **Marcus**, the friend Dave invites as his accountability partner.
**Method:** Read-only source walkthrough of all SwiftUI views and website pages. Contrast ratios computed from the actual palette constants. No product code was modified.

> Round 1 covered security, privacy, and cost. This review is about whether a real person can get the thing working. Where a UX dead end has an underlying defect already filed in round 1, it is cross-referenced as **[R1-Fxx]** rather than re-litigated.

---

## The short version

Dave can install the app, complete onboarding, see a dashboard that says **"Monitoring Paused"**, and have **nothing whatsoever being monitored** — no apps selected, both monitoring toggles off — with no screen anywhere telling him that. The flow that actually turns monitoring on lives three taps deep in Settings and is never linked from onboarding.

Meanwhile Marcus gets an invite email, clicks it, and lands on a page that tells him to go search the App Store by hand — because **the website has no App Store link anywhere**, even though the app is approved and live.

Those two things, plus the category mismatch in finding **U-03**, mean the product's core loop does not complete for a new pair of users. Everything else in this document is secondary to those.

---

## Tap counts

Counted as discrete user actions (taps, keyboard Next/Return, system-dialog buttons). Gestures noted separately.

| Flow | Taps | Notes |
|---|---:|---|
| **Sign up → onboarding complete** (invite a partner) | **12** | 1 welcome, 4 account, 3 invite, 4 permissions |
| **Sign up → onboarding complete** (skip invite) | **8** | "Skip for now" is one tap |
| **…then actually enable monitoring** | **+8 → 20 total** | Settings → App Restrictions → picker → 2 toggles. **Not linked from onboarding.** |
| **Start Deep Scan** (Control Center already set up) | **3** + 1 swipe + 3 s wait | Leaves the app entirely |
| **Start Deep Scan** (Control Center *not* set up) | **+5 in the Settings app** | App never explains this |
| **Invite a partner from Settings** | **5** | Settings → Manage Partners → field → type → Send |
| **Create a group** | **~9** | 4-step sheet: name, covenant, invites, confirm |
| **Join a group** | **0 in-app** | No join flow exists — you are added by email [R1-F05] |
| **Partner reads an alert** (adult\_content / gambling) | **2** | Tap push → detail view |
| **Partner reads an alert** (app\_usage / self\_harm / violence) | **1 → dead end** | Push taps to an empty dashboard — **U-03** |
| **Partner reads a tamper alert** ("X turned off monitoring") | **1 → dead end** | Routes to dashboard; nothing is stored [R1-F03] |
| **Set a partner PIN** (partner side) | **6** | Settings → Manage Partners → member → Set PIN → 4 digits |
| **Send feedback** | **6** | Settings → Feedback → type → summary → details → Send |
| **Sign out** | **2** | Settings → Sign Out. **No confirmation. No PIN gate.** — U-06 |
| **Delete account** | **4 + typing "DELETE"** | Correctly heavy friction |

The two numbers that matter: **20 taps before anything is monitored**, and **2 taps to turn it all off**.

---

## Top 5 places a new user quits

### 1. The dashboard says "Monitoring Paused" when monitoring is on — and he never turned anything on anyway

Dave finishes onboarding, grants Screen Time, lands on Home, and reads a grey card: **"Monitoring Paused — Activity recording is paused · Session ended."**

`StatusCard` is driven **solely by `isBroadcasting`** — the ReplayKit Deep Scan state. The file says so outright:

```swift
// DashboardView.swift:454-456
// Driven solely by isBroadcasting (the real broadcast state from EventProcessor).
```
```swift
// DashboardView.swift:487-490
Text(isBroadcasting ? "Monitoring Active" : "Monitoring Paused")
Text(isBroadcasting ? "Your screen activity is being monitored"
                    : "Activity recording is paused")
```

Deep Scan is the *optional* layer. The always-on Family Controls filtering and DeviceActivity monitoring — what the website calls the "credible core" — are invisible to this card. So the headline answer to the only question that matters is wrong in both directions: on when the core is off, off when the core is on.

And in Dave's case nothing *is* on, because `isMonitoringEnabled` and `isShieldingEnabled` both default to `false` (`ActivitySelectionManager.swift:31-36`) and no apps are selected (`selectionSummary` returns `"Nothing selected"`, `:220`).

**Smallest fix.** Make the card reflect a composite state, computed from what is actually true, and never say "Paused" when the core is running:

| Condition | Card |
|---|---|
| Family Controls not approved | **"Setup needed"** → tap to grant |
| Approved, 0 apps selected | **"No apps selected"** → tap to pick apps |
| Apps selected, both toggles off | **"Monitoring off"** → tap to turn on |
| Monitoring and/or shielding on | **"Monitoring on"** + "Deep Scan: off" as a secondary line |
| Broadcast running | **"Monitoring on · Deep Scan active"** |

Keep Deep Scan as a *subtitle*, never the headline. **Effort: low** (one computed property + copy). **Impact: highest in this document.**

---

### 2. Nothing ever tells him to pick apps

`ScreenTimeMonitoringView` — the only place apps and categories can be chosen — is referenced in exactly one spot:

```
$ grep -rn "ScreenTimeMonitoringView" RemainFaithful
SettingsView.swift:94:  NavigationStack { ScreenTimeMonitoringView() }
```

Onboarding's final step (`MonitoringSetupStep`) requests the Family Controls *permission* and then ends. Permission is not selection. Dave has authorized Screen Time and selected nothing, and the app presents this as done — the button literally reads **"Get Started"** (`OnboardingView.swift:490`).

The app tour does point at it ("Choose what to block", `AppTour.swift:61-62`) — but it is **step 7 of 8**, and one tap on "Skip" dismisses the tour permanently (`ContentView.swift:55` → `finishTour()` sets `tourCompletedVersion`).

**Smallest fix.** Add a fifth onboarding step that presents the `FamilyActivityPicker` directly after permission is granted, with a "Choose apps" primary button and a genuinely secondary "I'll do this later". If zero apps are selected when onboarding ends, show a persistent dashboard banner: *"No apps are being monitored yet — choose apps"* linking straight to the picker. **Effort: medium** (reuse the existing view). **Impact: very high.**

---

### 3. "Continue Without Monitoring" looks exactly like the real button

On the final onboarding step, when permission has not been granted, **two identical full-width gold buttons** are stacked:

```swift
// OnboardingView.swift:479-491
RFButton(title: didRequestScreenTime ? "Waiting for permission…" : "Enable App Monitoring", ...)
RFButton(title: fcManager.authorizationStatus == .approved ? "Get Started"
                                                           : "Continue Without Monitoring", ...)
```

`RFButton` has a single enabled style — gold fill, navy text, 56 pt (`:560-579`). There is no secondary variant. So the escape hatch from the app's entire purpose carries the same visual weight as the primary action, and sits directly beneath it where a fast-tapping user lands.

**Smallest fix.** Give `RFButton` a `.secondary` style (text-only or outline) and apply it to "Continue Without Monitoring". Reword to **"Skip for now — monitoring stays off"** so the consequence is in the label. **Effort: trivial.** **Impact: high.**

---

### 4. Marcus clicks the invite link and is told to go search the App Store

Every invite email lands on `/invite` or `/accept-invite`, both rendering `InviteLanding`. Step 1 reads:

```ts
// website/src/components/InviteLanding.tsx:14, 33, 52
body: 'Remain Faithful is available for iPhone (iOS 17+). Open the App Store on your iPhone and search for Remain Faithful.',
```

There is **no App Store link on this page**. The two buttons are "How it works" and "What partners can see" — both marketing pages (`:103-110`).

A sitewide search confirms the link does not exist anywhere:

```
$ grep -rn "apps.apple.com\|itunes.apple" website/src
(no matches)
```

This is the highest-intent moment in the entire funnel — a man who just got an email from a friend asking him to help — and the page sends him to search a store by name.

**Smallest fix.** Add the App Store URL as the primary button on `InviteLanding`, and as step 1's inline link. **Effort: trivial.** **Impact: very high.**

---

### 5. Forgot Password shows a green checkmark for a flow that cannot work

```swift
// LoginView.swift:462-474
if sent {
    Image(systemName: "checkmark.circle.fill")   // green
    Text("Check your email")
    Text("If an account exists for \(email),\nyou'll receive a reset link shortly.")
}
```

Per **[R1-F06]**, that link resolves to a dead domain or a 404 — no `/reset-password` route exists on any site. A returning user who forgot his password is shown a confident success state and then hits nothing. For this audience, "email support about my porn-accountability account" is a step many will not take; the realistic outcome is abandonment.

**Smallest fix.** Ship the `/reset-password` page and set `APP_BASE_URL` (both in R1-F06). Until then, change the success copy to name the sender and the next step, and add a "Didn't get it? Contact support" link with a working mailto. **Effort: low for the copy, medium for the page.** **Impact: high for returning users.**

---

## Findings by flow

Ranked within each section. **Impact** = how many users hit it × how badly. **Effort** = eng time for the smallest fix.

### Sign up & onboarding

**U-01 · Terms and Privacy are not links** — P1 · impact high · effort trivial
`OnboardingView.swift:275-280` renders *"By continuing you agree to our Terms of Service and Privacy Policy."* as plain `Text` at 11 pt, white at 30 % opacity. Nothing is tappable. The user agrees to documents he cannot open, and the line is the lowest-contrast text on the screen (**2.64:1**, see the contrast table).
**Fix:** make both phrases `Link`s to `remainfaithful.com/terms` and `/privacy`; raise opacity to 0.6.

**U-02 · The disabled button never says why** — P1 · impact high · effort low
`canSubmit` requires `password.count >= 8` (`OnboardingView.swift:192-197`), and `LoginView.swift:30-32` does the same on *sign-in*. A user who types a 6-character password sees a dead grey button and no message. Worse, the rule lives only in the placeholder — `"Password (8+ characters)"` — which **disappears the moment he starts typing**.
**Fix:** persistent helper text under the field that turns from grey to green as each rule is met. On sign-in, drop the length requirement entirely — let the server reject it, so a legacy short password is not locked out by the client.

**U-03 · Four of six event categories are invisible in the app** — P0 · impact critical · effort low

The backend accepts six categories:

```go
// backend/internal/handler/events.go:43-47
"adult_content": true, "gambling": true,
"violence": true, "self_harm": true, "clean": true,
"app_usage": true,
```

The iOS enum models four, two of which the backend never sends:

```swift
// DashboardView.swift:5-9
case adultContent = "Adult Content"
case gambling     = "Gambling"
case socialMedia  = "Social Media"   // backend never sends this
case gaming       = "Gaming"         // backend never sends this
```

`EventCategory(apiString:)` returns `nil` for anything else (`:28`), and every call site discards nils:

```swift
DashboardView.swift:372    events = remote.compactMap { ActivityEvent.from(remote: $0) }
SettingsView.swift:1967    events = remote.compactMap { ActivityEvent.from(remote: $0) }
GroupView.swift:820        .compactMap { ActivityEvent.from(remote: $0.event) }
```

**Consequences.** `app_usage` is the category the always-on DeviceActivity layer emits — `DeviceActivityMonitorExtension.swift:65` and `:96` both post `"category": "app_usage"`. So **every event produced by the primary monitoring layer is silently dropped by the app**: it never appears on Dave's dashboard, never appears in Marcus's alert list, and never appears in the group view. The partner gets a push, taps it, and `routeNotification` falls through to `.dashboard` because the payload will not map (`NotificationService.swift:125-129`) — landing on a screen showing nothing about it.

The same is true of `self_harm`, whose push copy reads *"needs you right now — their device flagged self-harm content. Reach out immediately."* (`events.go:264`). Marcus taps that, and the app shows him an ordinary dashboard.

**Fix:** add `appUsage`, `violence`, `selfHarm`, and `clean` cases with icons, tints, and conversation starters; map them in `init?(apiString:)`. Delete `socialMedia` and `gaming` or add them to the server allowlist. Add a `default` case rather than returning `nil`, so an unknown future category degrades to a generic card instead of vanishing.

**U-04 · "Two layers" contradicts the website's five** — P3 · impact low · effort trivial
`OnboardingView.swift:457` — *"Two layers work together"* — lists only App Monitoring and Deep Scan. The site describes always-on filtering, DeviceActivity, time-window shielding, Deep Scan, and App Lockout.
**Fix:** say "Two things to set up now" and link to the full explanation.

**U-05 · The Deep Scan instruction is the least readable text on the screen** — P2 · impact high · effort trivial
`OnboardingView.swift:500-505`: *"Deep Scan starts from Control Center → Screen Recording → Remain Faithful"*, 11 pt, white at 30 % opacity (**2.64:1 — fails WCAG AA**). This is the only time the app explains how to start Deep Scan during onboarding, and it assumes Screen Recording is already in Control Center. If it is not, the user needs five taps in the iOS Settings app that nothing mentions.
**Fix:** promote to a proper card at 13 pt / 0.6 opacity with numbered steps, plus a line covering "don't see Screen Recording? Settings → Control Center → add it."

### Inviting a partner

**U-06 · Sign Out is the easiest way to switch everything off, and nobody is told** — P0 · impact critical · effort low

```swift
// SettingsView.swift:403-407
Button {
    AuthState.shared.clearSession()
    userName  = ""
    userEmail = ""
} label: { Text("Sign Out") }
```

No confirmation dialog. No PIN gate. No partner notification.

`clearSession()` removes the token from the App Group store the extensions read (`AuthState.swift:87`), so event upload and heartbeats stop instantly (`SampleHandler.swift:543,700`; `DeviceActivityMonitorExtension.swift:81`). `protection.go:58-67` has no alert type for sign-out, so nothing is sent.

The inconsistency is stark: **"App Restrictions" is PIN-gated** (`SettingsView.swift:207`, `withPINGate`), **App Lockout is PIN-gated** (`:307`), and **Delete Account requires typing "DELETE"** (`:564`). The one action that silently defeats all of it is two taps with no friction at all.

**Fix:** put Sign Out behind `withPINGate` like the other protection controls, add a confirmation sheet reading *"Signing out stops monitoring. Your partners will be told."*, and add a `signed_out` type to the protection-alert allowlist so the partner is notified.

**U-07 · Invite errors are replaced with a guess** — P2 · impact medium · effort trivial

```swift
// OnboardingView.swift:424-429
} catch {
    errorMsg = "Couldn't send invite. Check the email and try again."
}
```

Every failure becomes "check the email", discarding the server's actual message. Invite yourself and the backend returns *"cannot add yourself as a partner"* (`invites.go:57`); the user is told to check an address that is perfectly valid.
**Fix:** surface `APIError.server(msg)` directly — the backend messages are already user-appropriate — and keep the generic line only for transport errors.

**U-08 · "Invitation Sent!" is secretly the Continue button** — P3 · impact low · effort trivial
`OnboardingView.swift:389-394`: after a successful send the button's label becomes "Invitation Sent!" and its action becomes `onContinue()`. Nothing signals that tapping it advances.
**Fix:** show a separate confirmation row and change the button to "Continue".

**U-09 · Inviting an existing partner silently does nothing** — P0 · impact critical · effort low · **[R1-F04]**
The API returns `{"status":"connected"}` and the UI shows success, but the relationship is stuck `pending` forever and that partner receives no alerts, ever. From Dave's side the partner appears in his list and looks fine. Purely a UX note here: **until R1-F04 is fixed, the app is confidently reporting a protection that does not exist.** Once fixed, show pending partners with an explicit *"Waiting for Marcus to accept"* state rather than rendering them identically to accepted ones.

### Groups

**U-10 · The invite result never says what happened** — P2 · impact medium · effort trivial

```swift
// GroupView.swift:1327-1330
try await APIClient.shared.groupEmailInvite(groupID: gid, email: email)
await MainActor.run { withAnimation(...) { isSent = true } }
```

The backend distinguishes three outcomes — `added` (already a user, joined now), `already_member` (409), and `invited` (emailed, must sign up). The client collapses all of them into one "sent" state. Dave never learns whether his friend is in the group right now or needs to go install an app.
**Fix:** decode `status` and show three distinct messages.

**U-11 · "No Group Yet" offers only half the paths** — P3 · impact low · effort trivial
`GroupView.swift:244-260` reads *"You are not in any accountability groups. Create one, or ask a member for an invite link."* But there is **no join-by-link flow** — `/join/[id]` on the website explicitly tells the visitor to ask to be invited by email (`InviteLanding.tsx:55-57`), because membership is email-keyed. The empty state advertises a path that does not exist.
**Fix:** *"Create a group, or ask a member to invite the email address on your account."*

**U-12 · Raw error strings reach the user** — P3 · impact low · effort trivial
`GroupView.swift:1333` assigns `error.localizedDescription`, which for `.decoding` renders as *"Unexpected response format: …"* and for `.network` as *"Network error: …"* (`APIClient.swift:145-146`).
**Fix:** map to plain language; log the detail instead.

### Alerts

**U-13 · The streak grid labels the wrong days** — P2 · impact medium · effort trivial
`StreakCard` hardcodes `["S","M","T","W","T","F","S"]` (`DashboardView.swift:550`), but the backend returns a **rolling** window — *"week\[0] = 6 days ago, week\[6] = today"* (`users.go:255-256`). Unless today happens to be Saturday, every check and cross sits under the wrong letter.
**Fix:** generate labels from today's weekday, walking back six days.

**U-14 · A brand-new user sees six red X marks on day one** — P2 · impact high · effort trivial
`users.go:319` sets `week[i] = !day.Before(createdDay) && !flagged[dayStr]`, so days *before the account existed* come back `false`. The card renders `false` as a red X (`DashboardView.swift:586-589`). Dave's very first screen tells him he failed six days running — before he had the app.
**Fix:** return a tri-state (`clean` / `flagged` / `before-signup`) and render the third as a dimmed empty circle.

**U-15 · Tamper alerts dead-end on the dashboard** — P1 · impact high · effort medium · **[R1-F03]**
`NotificationService.swift:134-139` routes `PROTECTION_ALERT` to `.dashboard`, and those alerts are never persisted server-side. Marcus receives *"Dave turned off activity monitoring"*, taps it, and lands on a screen with no trace of the message. Fixing R1-F03 (persisting them) is the prerequisite; the UX fix is a **Protection** section in the alerts list plus a dashboard banner while a tamper state is unresolved.

**U-16 · The in-app alert banner auto-dismisses in 6 seconds** — P3 · impact low · effort trivial
`DashboardView.swift:328-334`. A foreground alert about a partner's relapse vanishes on a timer with no history entry if missed.
**Fix:** require an explicit dismiss for `severe`, or leave a persistent row behind.

### Deep Scan, Lockout, PIN

**U-17 · "App Lockout (Deep Scan)" / "Off — optional hardening"** — P2 · impact medium · effort trivial
`SettingsView.swift:256-262`. Three pieces of jargon in one row; "optional hardening" is engineering vocabulary.
**Fix:** *"Keep apps blocked unless Deep Scan is running"* / *"Off — apps stay blocked only by your normal settings."*

**U-18 · Six names for two concepts** — P2 · impact medium · effort low
"App Monitoring" (onboarding), "App Restrictions" (nav title + Settings row), "Monitor App Usage" (toggle), "Block Selected Apps" (toggle), "Choose Apps to Restrict" (button), "Screen Time Access" (permission), plus Settings sections headed **MONITORING**, **SCREEN TIME**, and **PROTECTION** that all govern overlapping things.
**Fix:** settle on two user-facing nouns — **Alerts** (partner gets told) and **Blocking** (app will not open) — and use them everywhere.

**U-19 · A wrong PIN is indistinguishable from no PIN** — P2 · impact medium · effort low · relates to **[R1-F24]**
`VerifyPIN` returns `{"success": false}` both when the PIN is wrong and when none is configured (`protection.go:289-293`).
**Fix:** distinguish the two server-side and tell the user *"No PIN is set for this partnership"* rather than implying he typed it wrong.

### Tour, feedback, settings

**U-20 · One tap kills the tour forever** — P2 · impact medium · effort trivial
`ContentView.swift:55` — `onSkip: { finishTour() }` writes `tourCompletedVersion` and it never returns. The replay path exists (`.replayAppTour`) but lives in Settings.
**Fix:** make Skip mean "not now" and re-offer once on the next launch; keep permanent dismissal for reaching the final step.

**U-21 · Feedback truncates silently while typing** — P2 · impact low · effort trivial
`FeedbackSheet.swift:196-197` and `:215-216` hard-truncate at 120 and 4000 characters with no counter anywhere. Characters simply stop appearing.
**Fix:** show "n/120" once past ~80 %, and stop silently discarding keystrokes.

**U-22 · `support@remainfaithful.com` may not be a real mailbox** — P2 · impact medium · effort trivial · *suspected*
`FeedbackSheet.swift:243,272` tell users to email `support@remainfaithful.com`, but the backend routes contact and feedback mail to `jeff@hanokventures.co` (`contact.go:44`, `feedback.go:57`) and SendGrid sends from `noreply@remainfaithful.app` (`email.go:27`) — a **third** domain, which is also the dead one from R1-F06. I cannot verify mailbox existence from the repo.
**Fix:** confirm the address receives mail, or print the address that actually does. Align the sending domain with the live site.

**U-23 · The tour's first step confirms the status card is wrong** — P2 · impact medium · effort trivial
`AppTour.swift:40-41` names the card **"Deep Scan status"**. The tour is accurate; the card's own copy ("Monitoring Paused") is what overclaims. Fixing U-01's status model should bring this label along.

---

## Accessibility

### Dynamic Type — effectively unsupported

```
fixed-point fonts  .font(.system(size: N))   361
Dynamic Type styles .font(.body/.headline/…)  13
scaled custom fonts .font(..., relativeTo:)    2
```

Roughly **96 % of text in the app does not scale.** A user at Larger Accessibility Sizes sees the same 11 pt disclaimers as everyone else. Only `TourOverlay` does it properly (`AppTour.swift:118`, `@ScaledMetric(relativeTo: .body)`) — that is the pattern to copy.

Compounding it, fixed-height containers will clip rather than grow: `RFButton` is `.frame(height: 56)` (`OnboardingView.swift:571`), field rows are `.frame(height: 56)` (`:585`), and the Forgot Password sheet is pinned to `.presentationDetents([.height(340)])` (`LoginView.swift:80`).

**Fix, in priority order:** (1) convert the Forgot Password detent to `.medium`; (2) change the ~20 fixed heights to `minHeight`; (3) migrate text to `.font(.system(.body, design: .serif))`-style APIs or `@ScaledMetric`, starting with onboarding, the status card, and alert detail. **Effort: medium-high but incremental.**

### Contrast — low-opacity text fails AA

Computed from `Color.rfNavy` `(18,33,76)`, `rfNavyMid` `(28,48,102)`, `rfGold` `(209,171,76)` (`OnboardingView.swift:6-8`), compositing white at the opacity used:

| Element | Opacity | vs rfNavy | vs rfNavyMid | Verdict |
|---|---:|---:|---:|---|
| Terms / Privacy line (`Onb:277`) | 0.30 | 2.64 | 2.50 | **FAIL** |
| Deep Scan hint (`Onb:502`) | 0.30 | 2.64 | 2.50 | **FAIL** |
| Disabled button label (`Onb:569`) | 0.35 | 2.46 | — | **FAIL** |
| Field placeholders (`Onb:217`) | 0.38 | 3.41 | 3.15 | large text only |
| Scripture line (`Onb:146`) | 0.38 | 3.41 | 3.15 | large text only |
| "Skip for now" (`Onb:403`) | 0.40 | 3.63 | 3.33 | large text only |
| tierRow detail (`Onb:529`) | 0.50 | 4.88 | 4.34 | borderline on gradient |
| Subtitles (`Onb:556`) | 0.55 | 5.61 | 4.93 | PASS |
| Gold button text (navy on gold) | — | **7.13** | — | PASS |

The gold-on-navy brand pairing is genuinely strong at 7.13:1. The failures all come from stacking low opacity on small type.

**Fix:** floor body text at opacity **0.60** and secondary text at **0.55**; never pair opacity < 0.5 with sizes under 13 pt. Note that the 0.50 rows fail against the *lighter* end of the gradient, so test against `rfNavyMid`, not `rfNavy`.

### Tap targets — sheet close buttons are ~30 pt

The dismiss control repeats across at least eight sheets as a 14 pt glyph with 8 pt padding — about **30 × 30 pt**, below the 44 × 44 pt HIG minimum:

```swift
// ScreenTimeMonitoringView.swift:41-46 (same shape at LoginView.swift:449,
// FeedbackSheet.swift:150, SettingsView.swift:1253/1645/1811/1897, DashboardView.swift:990)
Image(systemName: "xmark")
    .font(.system(size: 14, weight: .semibold))
    .padding(8)
    .background(Circle().fill(Color.white.opacity(0.08)))
```

**Fix:** add `.frame(width: 44, height: 44)` to the button label — the circle can stay visually 30 pt while the hit area grows. One shared `CloseButton` view fixes all eight.

### VoiceOver — very thin

Across 9,868 lines: **13** `accessibilityLabel`, **2** `accessibilityHint`, **0** `accessibilityValue`, **3** `accessibilityHidden`.

Specific gaps worth fixing first:
- **StatusCard** conveys state through colour and a pulsing dot with no combined label. Add `.accessibilityElement(children: .combine)` and a value of "Monitoring on" / "off".
- **StreakCard's 7 dots** are bare `checkmark`/`xmark` images (`DashboardView.swift:586-589`) — VoiceOver reads fourteen meaningless glyph names. Label each "Monday, clean" / "Tuesday, flagged".
- **Decorative halo circles** in onboarding and login are not hidden; mark them `.accessibilityHidden(true)`.
- **Severity is colour-only** in `FlagSeverity.color` (`:90-96`). Add a text or shape cue so red/amber/green is not the sole signal.

---

## Website

**W-01 · No App Store link anywhere** — P0 · impact critical · effort trivial
Covered as quit point #4. The homepage's primary gold CTA is `<a href="#download">Get started</a>` (`page.tsx:104`), and the `#download` section it scrolls to (`:139-168`) contains **no download link** — only "How it works" and "Android notify". A live, approved iPhone app has no way to be installed from its own marketing site.
**Fix:** App Store badge as the primary button in the hero, in `#download`, and on `InviteLanding`.

**W-02 · The invite page's instructions cannot work for existing users** — P1 · impact high · effort low · **[R1-F04]**
`/accept-invite` is the link sent to people who *already have accounts* (`invites.go:73`), yet `InviteLanding`'s partner copy leads with *"Create your account with this email"* and the token is deliberately unused — *"The token is never read or sent anywhere from the website"* (`invite/page.tsx:12-13`). For an existing user there is currently **no path at all** to accept.
**Fix:** once R1-F04 persists the token, make this page deep-link into the app to redeem it. Until then, lead the `/accept-invite` variant with the existing-account instructions instead of burying them in step 3.

**W-03 · Waitlist inputs have no labels** — P2 · impact medium · effort trivial
All three `WaitlistForm` variants use `placeholder` as the only label (`WaitlistForm.tsx:57-61, 80-84, 110-122`) — no `<label>`, no `aria-label`. Placeholder-as-label is a WCAG 3.3.2 failure and the text vanishes on input.
**Fix:** add visually-hidden `<label htmlFor>` to each.

**W-04 · Errors are rendered in the brand gold** — P3 · impact low · effort trivial
`WaitlistForm.tsx:72,93` style failures with `text-wax`, the same gold used for CTAs and accents, so an error does not read as one.
**Fix:** a distinct error colour plus `role="alert"`.

**W-05 · Seven pages lack canonical URLs** — P3 · carried from R1-F41.

**Website credit where due.** The accessibility fundamentals here are notably better than the app's: a real skip link (`layout.tsx:106` + `globals.css:134`), a global `:focus-visible` gold outline (`:98-100`), a proper focus replacement on inputs rather than a bare `outline: none` (`:157-161`), exactly one `<h1>` per page across all 24 pages, `min-h-11` touch targets on nav links, and `aria-hidden` on decorative ornaments. The privacy copy is also unusually careful and accurate.

---

## Fix list, ranked by impact ÷ effort

### Do first — trivial effort, large payoff

| # | Fix | Where |
|---|---|---|
| 1 | Add the App Store link (hero, `#download`, `InviteLanding`) | W-01 |
| 2 | Make "Continue Without Monitoring" visually secondary and name the consequence | U-03 |
| 3 | Make Terms / Privacy real links and raise contrast to 0.6 | U-01 |
| 4 | Generate streak day letters from today instead of hardcoding Sun–Sat | U-13 |
| 5 | Render pre-signup days as dimmed circles, not red X marks | U-14 |
| 6 | Add `.frame(width: 44, height: 44)` to the shared close button | a11y |
| 7 | Surface the server's actual invite error text | U-07 |
| 8 | Raise the Deep Scan instruction to 13 pt / 0.6 opacity | U-05 |

### Do next — low effort, high impact

| # | Fix | Where |
|---|---|---|
| 9 | **Rewrite `StatusCard` to a composite monitoring state** | U-01 / quit #1 |
| 10 | **Add the four missing `EventCategory` cases and a `default`** | U-03 |
| 11 | **PIN-gate Sign Out, confirm it, and notify the partner** | U-06 |
| 12 | Persistent "No apps selected" dashboard banner → picker | quit #2 |
| 13 | Password rules as live helper text; drop the length gate on sign-in | U-02 |
| 14 | Decode and display the three group-invite outcomes | U-10 |
| 15 | Labels on waitlist inputs | W-03 |
| 16 | Forgot Password detent `.medium`; fixed heights → `minHeight` | a11y |

### Then — medium effort

| # | Fix | Where |
|---|---|---|
| 17 | Add an app-picker step to onboarding | quit #2 |
| 18 | Settle the vocabulary on **Alerts** and **Blocking** throughout | U-18 |
| 19 | Protection-alert history in the alerts list (needs R1-F03) | U-15 |
| 20 | VoiceOver labels for StatusCard, StreakCard, severity | a11y |
| 21 | Ship `/reset-password`; fix the success copy | quit #5 |
| 22 | Migrate text to Dynamic Type, onboarding first | a11y |

---

## What already works well

- **The app tour** is the best-written surface in the product: short, concrete, plain-spoken ("Type an email. They join when they sign up with that email."), correctly skips the invite step when there is no group, and is the only view that handles Dynamic Type and Reduce Motion properly (`AppTour.swift:118`, `DashboardView.swift:386-392`).
- **Delete Account** applies friction exactly where it belongs — type "DELETE" to confirm (`SettingsView.swift:564`). The problem is that Sign Out has none.
- **Conversation starters** in `AlertDetailView` are genuinely thoughtful and non-shaming — *"What was going on in your heart before this happened?"* — and are the clearest expression of the product's philosophy anywhere in the codebase.
- **`PauseNotice`** tells the truth about whether pausing notifies partners, conditioned on the actual lockout state (`AppTour.swift:70-77`). That is the honesty the rest of the status UI should match.
- **The `StatusCard` comment** records why a previous local-only toggle was removed because it contradicted the real broadcast state — the right instinct, applied to the wrong source of truth.
- **Website accessibility fundamentals** (skip link, focus rings, heading structure, touch targets) are solid and ahead of the app.
- **Privacy copy** on `/privacy-architecture` and `/how-it-works` is precise, and round 1 confirmed it matches the code.

---

## Appendix — the two walkthroughs

**Dave, first run.** Welcome → "Begin Your Journey" → name/email/password → "Get Started" → invite Marcus → "Skip for now" → "Enable App Monitoring" → Apple's Screen Time dialog (no priming beyond two short rows) → "Get Started" → notification permission → Home. **8 taps.** The card reads "Monitoring Paused". The streak card shows **six red X marks**. The tour may appear; if he taps Skip it never returns. Nothing is monitored and nothing says so. To fix it he must find Settings → App Restrictions → Choose Apps → picker → two toggles: **8 more taps, none signposted**.

**Marcus, invited.** Email arrives: *"Dave invited you to Remain Faithful"* → taps "Accept Invitation" → lands on `/invite` → told to open the App Store and search by hand, with no link → if he perseveres, installs, and signs up **with the exact invited address**, the invite redeems (`auth.go:71`). If he already had an account, nothing happens at all [R1-F04]. Once connected, Dave's `app_usage` alerts are dropped by his app [U-03], and tamper alerts tap through to an empty dashboard [R1-F03].

Both journeys fail at a different point, for different reasons, and neither failure is visible to the person experiencing it. That invisibility — not any single bug — is the core finding of this review.
