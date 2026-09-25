# Remain Faithful — App Store Submission Package

**Version:** 1.0.0  
**Bundle ID:** `com.remainfaithful.app`  
**Target:** iOS 17+ iPhone  
**Account:** Jeff Brewer / CGABrew7 Apple Developer Team  
**EOW goal:** TestFlight internal → Submit for Review

This file is the paste-ready App Store Connect package plus the steps only Jeff can perform. Agents cannot click Submit in App Store Connect, cannot approve the Family Controls distribution entitlement, and cannot archive from Xcode.

---

## Honest status (2026-09-22)

| Item | Status |
|---|---|
| Website `https://www.remainfaithful.com` | Live |
| API `https://remain-faithful-api.fly.dev/health` | Live `{"status":"ok"}` |
| Protection stack (Family Controls + DeviceActivity + Shield + Deep Scan + Partner PIN + App Lockout) | In source |
| Heartbeat silence sweep (every 5 min, default 30 min) | Wired in `main.go` |
| Woodfield Stripe Payment Links (8 amounts) | Live on `donate.stripe.com` (PR #19) |
| App Store listing | **Not submitted** |
| Family Controls **distribution** entitlement | **Jeff must confirm Apple approved it for this Team ID** |
| TestFlight build | **Jeff must archive + upload from Xcode** |

---

## What only Jeff can do (do these in order)

1. **Apple Developer → Certificates, Identifiers & Profiles**  
   Confirm App ID `com.remainfaithful.app` has:
   - Push Notifications
   - Sign in with Apple
   - App Groups `group.com.remainfaithful.app`
   - **Family Controls** (distribution, not just development)
   If Family Controls shows “Development only,” request the distribution entitlement from Apple *today*. Without it, TestFlight external and App Store will fail signing.

2. **App Store Connect → New App**  
   Name: Remain Faithful  
   Primary language: English (US)  
   Bundle ID: `com.remainfaithful.app`  
   SKU: `remainfaithful-ios-001`  
   User access: Full Access for Jeff

3. **Archive in Xcode**  
   Scheme: `RemainFaithful`  
   Destination: Any iOS Device (arm64)  
   Product → Archive → Distribute App → App Store Connect → Upload  
   Include all extensions: Broadcast, DeviceActivity, ShieldConfig.

4. **TestFlight internal**  
   Add Jeff + 1–2 partners as internal testers. Walk the protection checklist below on a real iPhone before submitting.

5. **Confirm Woodfield Stripe account name**  
   Open any Payment Link. The Stripe Checkout header must show the Woodfield Foundation legal name, not a personal name. Receipts that will be called tax-deductible must show the 501(c)(3) legal name + EIN.

6. **Paste listing copy** from this file into App Store Connect (name, subtitle, description, keywords, what’s new, review notes, privacy nutrition).

7. **Screenshots**  
   Required: 6.7" (iPhone 16 Pro Max / 15 Pro Max) and 6.1" (iPhone 16 / 15). Optional 5.5".  
   Capture from TestFlight: Onboarding, Dashboard (protection on), Group / partner invite, App Restrictions (FamilyActivityPicker after auth), Alert detail (metadata only), Donate.

8. **Submit for Review**  
   Only after TestFlight protection checklist is green and Family Controls distribution entitlement is approved.

---

## App Store Connect listing (paste-ready)

### Name
Remain Faithful

### Subtitle (30 characters max)
Free peer accountability

### Promotional text (170 characters, optional)
Always-on app blocking. Partners get a category and time — never screenshots. Free forever, sustained by tax-deductible gifts.

### Description

Remain Faithful is free peer accountability for adults who want to stay faithful — spouses, friends, mentors, and church groups.

Most accountability apps either send screenshots of your screen to another person, or they depend on you telling the truth at the worst moment. Remain Faithful does neither.

HOW IT WORKS

1. Invite a trusted partner or a small group (up to 12).
2. Turn on always-on filtering with Apple Screen Time. Chosen apps and categories stay blocked even if you force-quit the app or restart the phone.
3. If a blocked app is opened, your partner gets a discreet alert: category, severity, time, and a short system summary. Not the app name. Not a screenshot. Not the page you were on.

PRIVACY, BY DESIGN

• Screenshots, OCR text, and raw screen content never leave the device.
• Optional Deep Scan uses Apple Vision and Sensitive Content Analysis on-device only.
• You choose your partners. You can pause or remove access instantly.
• Source is public: github.com/CGABrew7/remain-faithful

PROTECTION THAT HOLDS

• Always-on Family Controls blocking
• DeviceActivity usage signals
• Optional time-window shielding
• Optional Deep Scan for high-risk windows
• Partner PIN so settings cannot be quietly changed
• Alerts if Screen Time permission is revoked or the device goes silent

FREE FOREVER

No subscription. No premium tier. No ads. Optional donations are processed by Woodfield Foundation and help keep servers and outreach running.

Remain Faithful is for adults 18+. It is a covenant tool, not surveillance, and not a substitute for pastoral care or counseling.

### Keywords (100 characters max, comma-separated, no spaces after commas if tight)
accountability,purity,faithful,christian,screen time,blocker,partner,church,marriage,recovery

Count check: accountability,purity,faithful,christian,screen time,blocker,partner,church,marriage,recovery = 99 characters.

### What’s New (1.0.0)
First public release. Always-on Screen Time blocking, partner and group alerts, on-device Deep Scan, Partner PIN, and tax-deductible donations to keep the app free.

### Support URL
https://www.remainfaithful.com

### Marketing URL
https://www.remainfaithful.com

### Privacy Policy URL
https://www.remainfaithful.com/privacy

### Category
Primary: Lifestyle  
Secondary: Health & Fitness

### Age rating
18+  
Reason: frequent/intense mature/suggestive themes discussed in the context of avoiding them. No pornography is shown in the app.

### Pricing
Free. No In-App Purchases required. Donations are optional and processed outside the app via Stripe (Safari).

---

## Privacy nutrition labels (App Store Connect)

Match `RemainFaithful/PrivacyInfo.xcprivacy`:

| Data type | Linked to identity | Used for tracking | Purpose |
|---|---|---|---|
| Name | Yes | No | App Functionality |
| Email Address | Yes | No | App Functionality |
| Device ID (APNs token) | Yes | No | App Functionality |
| Other Usage Data (alert category, severity, timestamp, short summary) | Yes | No | App Functionality |

Not collected: precise location, contacts, photos, browsing history, screenshots, product interaction for advertising, tracking.

NSPrivacyTracking = false.

---

## Review notes (paste into App Store Connect)

Remain Faithful is an adult self-accountability app. It is not parental control for children.

Screen Time / Family Controls
We request the Family Controls entitlement in individual (self) mode so an adult can block apps and categories on their own iPhone. Authorization is requested in-app via AuthorizationCenter.requestAuthorization(for: .individual). Shields use ManagedSettingsStore. Usage signals use DeviceActivityMonitor. Revoking Screen Time in iOS Settings clears shields and notifies the user’s chosen partners that protection is off.

Broadcast / Deep Scan
The ReplayKit Broadcast Upload extension (“Remain Faithful Deep Scan”) is optional and user-started every session. Frames are classified on-device with Vision OCR and SensitiveContentAnalysis. Screenshots, OCR text, and raw frames are never uploaded. Partners receive only category, severity, timestamp, and a short system-generated summary. DRM video frames are black and are not classified.

Demo account
Email: review@remainfaithful.com
Password: (Jeff: set a review account before submit and paste it here)

How to test the core flow
1. Sign in with the demo account.
2. Settings → App Restrictions → Enable Screen Time Access → approve individual Family Controls.
3. Choose a harmless category or app, enable Block Selected Apps.
4. Confirm the iOS shield appears if that app is opened.
5. Invite flow: Group tab → New Group or partner invite. Alerts show metadata only.
6. Donate opens Stripe Checkout in Safari. No StoreKit purchase is required to use the app.

Contact for review questions: jeff@hanokventures.co or support@remainfaithful.com

---

## Protection checklist (TestFlight, real device)

Run this on a physical iPhone before Submit.

- [ ] Family Controls individual authorization succeeds
- [ ] Selected apps show the system Screen Time shield
- [ ] Force-quit Remain Faithful — shield still holds
- [ ] Reboot — shield still holds
- [ ] Turning blocking off sends `shielding_disabled` to partner
- [ ] Revoking Screen Time in iOS Settings sends `family_controls_revoked`
- [ ] Deep Scan starts only after system broadcast permission
- [ ] Stopping Deep Scan with App Lockout on re-applies shield + `deep_scan_stopped`
- [ ] Partner PIN required to change protected settings; wrong PIN alerts partner
- [ ] Heartbeat posts while Deep Scan / monitoring is active
- [ ] After HEARTBEAT_SILENCE_MINUTES (default 30) of no heartbeat, partner receives `heartbeat_silence`
- [ ] Alert detail in partner UI shows category / severity / time / summary — never a screenshot
- [ ] Donate preset $10 opens Woodfield Stripe Checkout
- [ ] Sign in with Apple works
- [ ] Delete Account in Settings cascades server-side

---

## Screenshot shot list (6.7" and 6.1")

1. Hero — Onboarding covenant line + “Free forever”
2. Dashboard — protection ON, streak, partner count
3. App Restrictions — Screen Time approved, apps selected, Block enabled
4. Group — invite code + member list (3–12)
5. Alert — metadata card, no screenshot, “Start a conversation” prompt
6. Donate — $10, Woodfield / tax-deductible line visible

Do not show explicit content, blocked-site pages, or any real person’s alert about adult content in screenshots. Use a generic “Restricted app opened” fixture if needed.

---

## Entitlements inventory

Main app (`RemainFaithful.entitlements`)
- `aps-environment` production
- Sign in with Apple
- App Group `group.com.remainfaithful.app`
- `com.apple.developer.family-controls`

DeviceActivity + ShieldConfig extensions
- Family Controls + App Group

Broadcast extension
- App Group only (no Family Controls; writes `isBroadcasting` to the group)

---

## Known review risks (do not hide these)

1. Family Controls distribution entitlement missing → signing / upload failure.
2. Reviewer treats Deep Scan as surveillance → rejection under 5.1.1 / 2.5.14. Mitigate with review notes above and the privacy site.
3. Donations called tax-deductible if Stripe account legal name / EIN is wrong → legal risk, not just review.
4. 18+ rating required. Do not market as a kids’ Screen Time app.

---

## After upload

Typical first-review window: 24–48 hours, longer if Family Controls is new on the account. If Apple asks for a screen recording of Deep Scan, record a session on Settings or Safari home — never on adult content — and show the partner alert metadata card.
