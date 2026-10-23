# Round 2 — UX, UI, and friction

Read-only walkthrough of the SwiftUI app and the website at commit `784998f`. No product code was changed. No iPhone was available, so tap counts are the steps in the source, not a stopwatch on a device. Anything that depends on VoiceOver, Dynamic Type, or a real Screen Time dialog is marked **suspected**.

The person in this walkthrough is a non-technical adult, and then their partner. The job of the app is simple: pick the apps that tempt you, tell a person you trust, and make it obvious whether that is actually on.

---

## The first ten minutes, as the code runs them

1. Welcome. One button: "Begin Your Journey". A smaller "Sign In".
2. Name, email, password. "Get Started". The terms line is not a link.
3. Partner email, or "Skip for now".
4. "Turn on monitoring". The button asks for Family Controls. It does not ask which apps. "Continue Without Monitoring" is the same size as the finish button.
5. iOS asks about notifications. The tour then covers Home, Group, and Settings.
6. Home says **Monitoring is paused** and **No flags yet — great work**.

That last screen is the problem. The user did what the app asked. The app then says monitoring is off, and also that they are doing great. Both can be true only because "monitoring" on Home means Deep Scan, and "great work" means "no rows came back". Neither means the shield is on.

A partner who already has an account gets an email, opens a web page that tells them to sign in, signs in, and is not connected. A partner with no account must create the account with the same email. That path can work, and the app never says "you are now their partner".

---

## Flow by flow

Tap counts ignore typing. System dialogs count as a screen.

### Download

Website path: search the App Store for Remain Faithful (`InviteLanding.tsx` step 1). There is no single "Download on the App Store" button on the invite page. The partner is on a computer or in Mail, reading three steps, then has to switch to the phone.

Smallest fix: one App Store button on the invite page, and the same button on the homepage hero. The email should say "Open this on your iPhone".

Impact high, effort low. Website only.

### Sign up

Screens: welcome, create account, invite, monitoring. Then the notification dialog, then the tour.

Taps if they skip the partner and skip monitoring: Begin (1), Get Started (2), Skip (3), Continue Without Monitoring (4), then Allow or Don't Allow, then the tour's Next buttons (about 7 more) or one Skip.

Taps if they try to do it right: the same, plus Enable App Monitoring and the system Allow. They still have not chosen an app.

Problems:

- "By continuing you agree to our Terms of Service and Privacy Policy" is plain text (`OnboardingView.swift:275-276`). It is not tappable.
- No "I am 18 or older" check. The policy says 18+. The screen does not.
- Password errors from the server are shown raw (`:306-309`). "email already registered" does not say "tap Sign In".
- A network failure says "Something went wrong. Please try again." That is fine. It does not say whether the account was created.

Smallest fix: make the legal names real links, add an 18+ checkbox that blocks Get Started until it is on, and if the email exists, say "That email already has an account. Sign in."

Impact medium, effort low. App Store build.

### Family Controls and choosing apps

Onboarding screen 4 (`OnboardingView.swift:436-510`):

- Title: "Turn on monitoring".
- Two blurbs: App Monitoring ("alerts your partner after 1 minute") and Deep Scan ("visual monitoring during screen broadcasts").
- One permission button. No picker.
- Footnote, in 11-point type at 30% white: Deep Scan starts from Control Center.

Later, Settings → App Restrictions (`SettingsView.swift:205-228`):

- The row says "Authorized" in green when the permission exists. It does not say whether any app is blocked.
- Opening it can demand the partner PIN (`:207`). The toggles inside are not asked again (master review M6).
- Inside: Choose Apps (system picker), Monitor App Usage, Block Selected Apps (`ScreenTimeMonitoringView.swift:152-206`).
- The note at the bottom is honest: one minute a day, and the shield can be overridden (`:219`).

Taps from a finished signup to a real block: Settings tab (1), App Restrictions (2), Choose Apps (3), pick apps in the system sheet and done (4), Block Selected Apps toggle (5). Monitor is a sixth tap if they also want the daily alert. That is the actual setup, and onboarding never mentions it.

"Authorized" is the wrong word for "you have not picked anything".

Smallest fix: after the system permission is allowed, present the app picker before "Get Started" is enabled. On the Settings row, replace "Authorized" with "No apps selected", "Blocking 3 apps", or "Alerts on, blocking off".

Impact highest of the setup fixes. Effort medium. App Store build.

### Inviting a partner

From onboarding or Settings → Manage Partners: type email, Send Invitation. Success copy is "Invitation Sent!" (`OnboardingView.swift:390`, `SettingsView.swift:1025`). The list underneath is "CURRENT PARTNERS" with no Pending or Connected label (`SettingsView.swift:882-923`). The status field is loaded (`:1057`) and never shown.

If that email has no account yet, the server stores an invite and signup with that email accepts it. The inviter is not told "they have to create an account". The app guide does say that (`guide-content.ts:131`). The button does not.

If that email already has an account, the server returns connected and the row stays pending (master M2). Both people think it worked.

Smallest fix: the success line should be "Invite sent. They are not connected until they accept." The row should say Waiting. Do not put a pending person under CURRENT PARTNERS.

Impact high, effort low in the app once the server tells the truth. App Store build plus the server fix.

### A partner with no account

Email link opens `/invite?token=...&type=partner`. The page does not read the token (`invite/page.tsx:11-16`). Three steps: get the app, create an account with this email, sign in. That matches the server, which accepts on signup by email, not by token.

Taps on the phone after install: Begin, fill the form, Get Started, then they are in the new-user flow including "invite your partner", which is the wrong ask for someone who was just invited. Nothing says "You are now paired with Alex".

They land on Home. Their own flag list is empty. Alex's future flags do not appear there (master M5).

Smallest fix: if signup consumed an invite, the next screen is "You're paired with {name}. You'll be notified here." Skip the "invite a partner" step, or offer it as optional and second. Add the inbox from M5.

Impact high, effort medium. App Store build. The website steps can stay, with one line: "Use this exact email."

### A partner who already has an account

The email links to `/accept-invite?token=...&type=partner`. The page ignores the token (`accept-invite/page.tsx:12-14`). Step 3 says "Already have an account with this email? Sign in." Signing in does not call accept. There is no in-app accept button. The relationship stays pending, so they will not get alerts.

This is a dead end dressed as instructions.

Smallest fix: an Accept screen in the app, fed by a push or by a banner on Home: "{Name} asked you to be their partner." One button. Until that ships, the website must not say that signing in accepts the invite.

Impact high, effort medium. App Store build, and a one-line website correction now.

### Creating or joining a group

No group: one screen, "No Group Yet", one button, "Create a Group" (`GroupView.swift:237-271`). The empty state is clear. Good.

Create is four steps (`CreateGroupSheet` around `:1362`): name, a long covenant editor, invite emails, then create. The tour calls this "four short steps" (`AppTour.swift:56`). The covenant is not short. A person who only wanted a group with their spouse has to edit a multi-paragraph oath before the group exists.

Joining from a shared link: `www.remainfaithful.com/join/{id}` tells the visitor to ask for an email invite (`InviteLanding` variant `link`). The id is not an authorization. That is honest, and it is a dead end if the sharer thought the link added people. The page should say that in the first sentence.

An existing account invited by email is added immediately (master M3). The website tells them to sign in and open the Group tab, which matches the bad behavior. After the consent fix, this page must say they have to accept.

Smallest fix: create the group from the name alone. Covenant and invites come after, from the group screen, which already has both. Change the join page's first line to "This link does not add you. Ask them to invite your email."

Impact medium, effort low. App Store for the shorter create. Website for the join sentence.

### Getting and reading an alert

The monitored person's Home list is their own events (`DashboardView.swift:369-372`). A partner's flags show up as a banner for about six seconds (`:305-333`) and, if the category is one of four known ones, as a detail screen when the push is tapped (`NotificationService.swift:124-128`). Violence, self-harm, and app usage fall through to Home.

Opening Home marks every alert seen (`ContentView.swift:81-95`, `:166-168`), including ones they never opened.

A one-to-one partner has no list of the other person's alerts. The only history is inside a group member sheet, last 10, same category filter (`GroupView.swift:813-821`). Taps to see a group member's flags: Group tab, tap the person, scroll to Recent Flags. "No recent flags" is the empty state there (`:721`). It does not say whether that person has monitoring on.

The detail screen rewrites the summary into a browsing story (`DashboardView.swift:141-154`). A one-minute app-usage event can be described as a website visit if it were mapped, and gambling always becomes "a gambling or sports-betting website was visited."

Smallest fix: a "Their alerts" section on Home when you are someone's partner, using the summary the server stored. Do not clear the badge until that list is opened. Empty copy: "No alerts yet" plus "You will only see alerts after they turn monitoring on and you have accepted."

Impact high, effort medium. App Store build.

### Deep Scan

Home, when a broadcast is not running, is a button: "Monitoring is paused" / "Tap to resume accountability" (`DashboardView.swift:928-931`). Tap opens a sheet, "Resume Monitoring", four Control Center steps (`:961-966`). That is the whole feature. There is no in-app start, because iOS does not allow one. The label is the bug.

The tour calls this card "Deep Scan status" (`AppTour.swift:40-41`), which is the right name. The card itself does not use it.

If a monitored app was used, a second banner offers "Start Broadcast Scan" and then the same instructions (`:867-908`). That button does what it says. Good.

Stopping Deep Scan notifies a partner only when App Lockout is on (`SampleHandler.swift:291-304`). Lockout defaults off (`AppLockoutManager.swift:6-7`, `:31`). The privacy page says this. The home card does not.

Smallest fix: rename the card to "Deep Scan is off". Subtitle: "App blocking is separate. Deep Scan is optional and starts in Control Center." If blocking is on, say so on the same card.

Impact high, effort low. App Store build.

### Lockout and PIN

Settings → App Lockout (Deep Scan). Default off. The subtitle when off is "Off — optional hardening" (`SettingsView.swift:260-266`). "Hardening" will not mean anything to this user.

Turning it off asks for the PIN if one is set (`:299-317`). Turning it on does not.

Partner PIN status, on the monitored phone: "Not set — ask your partner to enable in Manage Partners" or "Active — your partner controls gated actions" (`:284-286`). "Gated actions" is vague. The partner sets the PIN from the lock icon on a row, only when they are the watcher (`:935-951`). The icon has no label. **Suspected:** VoiceOver says "button" or nothing useful.

The PIN pad auto-submits at 4 digits (`PINEntryView.swift:41-44`). Cancel is there. A wrong PIN shakes and shows an error string from the verify call. There is no "what happens if I forget this" because the partner set it, not the user. That should be one sentence on the pad: "Only your partner knows this PIN. Ask them if you're stuck."

Taps to turn lockout on: Settings, toggle. Two. Taps for a partner to set a PIN: Settings, Manage Partners, lock icon, enter 4 digits, save. The sheet is a detents sheet (`SettingsView.swift:864-875`).

Smallest fix: "While this is on, your chosen apps stay blocked unless Deep Scan is running. Turning it off tells your partner." Name the lock button "Set PIN" / "PIN on".

Impact medium, effort low. App Store build.

### The tour

It starts by itself once notification permission is decided (`ContentView.swift:141-163`). It hides the whole tab bar from VoiceOver while it is up (`:39`). Steps: Deep Scan, streak, flags, the support button, then either create-a-group or the group, then partners, restrictions, feedback (`AppTour.swift:38-66`). Skip exists (`AppTour.swift:162`).

It runs before the user has picked apps. The restrictions step is the one that matters, and it is near the end. "Your clean streak" on day one is an empty streak explained as success.

Smallest fix: do not auto-start the tour until they have picked apps, or make step 1 "Pick the apps to block" and switch to that screen. Keep Skip.

Impact high for first-run quit, effort low. App Store build.

### Feedback

Settings → Send Ideas or Report a Problem. This sheet is the best screen in the app. It has a kind, a short title, details, an off-by-default device line, Send, and a mailto fallback. Accessibility labels are actually set (`FeedbackSheet.swift:155-277`).

Taps: Settings, the row, pick a kind, two fields, Send. Errors depend on `sendFeedback`. **Suspected** that a server error is a raw string. The mailto path is the escape hatch, and it is labeled.

Website `/feedback` is a separate form. Fine if it posts to the same inbox. Do not make the in-app form link out to a donate page.

Smallest fix: if send fails, say "Not sent. You can email support@remainfaithful.com" and keep the mailto button visible. It is already there. Make sure the failure state does not clear the text.

Impact low, effort low. App Store only if the error state wipes the draft (not confirmed; **suspected**).

### Settings and sign out

Settings is a long scroll: profile, partners, groups, covenant, notifications, app restrictions, lockout, PIN, activity log, how it works, replay tour, contact, rate, feedback, leave all groups, delete account, sign out.

Sign out is one tap, no confirm (`SettingsView.swift:403-407`). Delete account has a confirm. Leave all groups has a confirm. Sign out is the dangerous mis-tap, and it also leaves the previous person's shields in place (master M21).

"How It Works" and the website app guide are a second explanation of the same product. They will drift. The guide already says the reset link works.

Activity log is the user's own flags again, same category filter. It is not a partner view.

Smallest fix: "Sign out?" with Cancel and Sign Out. Under it, one line: "App blocking on this iPhone stays until you turn it off."

Impact medium, effort low. App Store build.

### Website, as a partner and as a new user

Pages a partner actually hits: invite, accept-invite, join, how-it-works, privacy.

| They read | The app does |
|---|---|
| Partners are notified when a blocked category is attempted (home, about, how-it-works, privacy). | A shield has no "they tried to open it" alert. Usage is one event after one minute in the day. |
| Time windows you set (how-it-works). | One schedule, all day. |
| Sign in to accept (accept-invite, existing account). | Sign-in does not accept. |
| Create the account with this email (new account). | This one matches. |
| Join link adds you. | It does not. The page eventually says to ask for an email. The title sounds like it worked. |
| Privacy: never disclosed to any third party, then Apple, Stripe, and Google Analytics. | The second list is the true one. |
| About: about $4 per person, and about $80 a month. | Those are not the same claim. |
| App guide: reset link works for one hour. | There is no reset page. |

The invite pages do say 18+ and "not parental controls" (`InviteLanding.tsx:98-100`). Keep that.

Contact form labels are not tied to fields with `htmlFor`. The FAQ sets `aria-expanded` and leaves the answer in the page. A skip link exists (`layout.tsx:106`). Good.

Smallest fix, in order: correct the blocked-attempt and time-window sentences, fix the accept-invite step that says sign-in is enough, fix the reset-link sentence when the reset page ships, and soften the privacy opener so it matches the third-party list on the same page.

Impact high for trust, effort low. Website only.

---

## Where the words will confuse someone

| Words on screen | What a person hears | What is true |
|---|---|---|
| Monitoring is paused | The product is off | Deep Scan is off. Blocking might be on. |
| No flags yet — great work | I am protected and clean | No events were returned. Setup may never have finished. |
| Authorized | I'm done | Permission only. |
| Continue Without Monitoring | A reasonable choice | Skips the only always-on layer. |
| Invitation Sent! / CURRENT PARTNERS | They will hear about me | Maybe pending forever. |
| Off — optional hardening | Ignore this | It is the switch that makes stopping Deep Scan visible. |
| Alert partner after 1 min of daily use | They hear every time I open it | Once a day, after one minute total. |
| Tap to resume accountability | One tap resumes it | Four steps in Control Center, and it is optional. |
| Add an accountability partner to use this feature | The support button will not work | A group is enough for the server. |

---

## Accessibility

Confirmed in source:

- Almost every screen sets `.font(.system(size:))` to a fixed size. Dynamic Type will not grow those labels. Login, onboarding, home, settings, and the PIN pad all do this.
- `accessibilityLabel` shows up on the feedback sheet, the tour, one settings hint, and the app-guide link. The home banners, partner rows, lock icon, shield toggles, and PIN keys do not set one.
- Toggles use `.labelsHidden()` (`ScreenTimeMonitoringView.swift:173`, `:203`, `SettingsView.swift:271`). The words sit beside the switch as separate text. **Suspected:** VoiceOver reads an unlabeled switch. That should be checked on a phone before calling it certain.
- The tour sets `.accessibilityHidden(true)` on the tab view while the tour is visible (`ContentView.swift:39`), so a VoiceOver user cannot reach the screen underneath until they skip or finish.
- White text at 30% to 45% opacity on navy (`Color.rfNavy` is about rgb 18, 33, 77) is used for subtitles, footnotes, and the Deep Scan instructions. **Suspected** contrast failure for the 11-point and 12-point lines. Not measured with a contrast tool.
- Buttons that are full-width and 52–56 points tall meet a 44-point target. The partner-row star, lock, and dismiss X are smaller (`SettingsView.swift:927-951`, home dismiss targets around 8 points of padding). **Suspected** tight targets.
- The PIN pad is a custom grid, not a text field. **Suspected** that it has no digit labels for VoiceOver.

Smallest fix: turn the three toggles into `Toggle("Block selected apps", isOn:)` so the label is the switch. Add labels to the lock and star buttons. Replace the smallest fixed sizes (11 and 12 at low opacity) with at least 15 points and a higher opacity, or use a text style that scales. Do not hide the app from VoiceOver during the tour; hide only the dimmed backdrop.

Impact medium, effort medium. App Store build. Do this in the same build as the home-screen copy, not as a separate project.

---

## Empty states and errors that do not say what to do

| Place | What it says | What it should say |
|---|---|---|
| Home, no events | No flags yet — great work | No flags yet. If you have not chosen apps, go to Settings → App Restrictions. |
| Home, load failed | Couldn't load activity. Tap to retry. | Good. Keep it. |
| Partners, none | No partners yet. Invite someone below. | Good. Add "They must accept before alerts start." |
| Group, none | No Group Yet. Create one, or ask for an invite link. | Good. Say the link does not add them by itself. |
| Group member, no flags | No recent flags | No alerts yet. This does not mean their shield is on. |
| Support, no partner | Add an accountability partner to use this feature | Name the real requirement: a partner who has accepted, or a group. |
| Signup, email taken | Server text "email already registered" | That email already has an account. Sign in. |
| Invite send failed | Couldn't send invite. Check the email and try again. | Good enough. |
| Forgot password, offline | Failed to send reset link — check your connection | Good. After the page exists, say "Check your email on this phone." |
| Remove partner, when you are the watcher | The server error "relationship not found" | You can leave. If it fails, say "Couldn't remove them. Try again." |

---

## Top 5 moments a new user quits

1. **Home says monitoring is paused, and also that they are doing great.** They just finished setup. They do not know there is a second, real setup under Settings. They close the app and tell their partner it does not work.
   Fix: one status line with three facts — apps not chosen yet, blocking on or off, Deep Scan on or off — and a single button that opens the app picker when nothing is chosen.

2. **They allow Screen Time and think they are done.** The row turns green and says Authorized. No apps are blocked. Their partner will never get a usage alert.
   Fix: do not leave onboarding until they pick at least one app, or show a persistent "Finish setup" card that does not look like success.

3. **They invite their spouse, who already uses the app, and nothing happens.** The button says sent. The spouse signs in because the email said to. No alerts ever arrive. They decide the product is theater.
   Fix: the pending state and the Accept button from master M2. This is the partners' quit, and it is the same bug as the security hole.

4. **They forget the password.** The email arrives. The link does not open a reset page. Asking for help means emailing someone about this account. Many will not.
   Fix: the reset page from master M4. Website only, so it can ship before the next binary.

5. **The tour and "Continue Without Monitoring" offer an easy way out before the one step that matters.** Skip is good. Skipping should not look like the setup succeeded.
   Fix: put "Choose apps" first. Let them skip Deep Scan without skipping the shield. Move the tour to after that choice.

---

## Smallest fixes, ranked by impact versus effort

| Rank | Fix | Impact | Effort | Ship |
|---|---|---|---|---|
| 1 | Home status: say what is actually on. Stop "Monitoring is paused" and "great work" when setup is unfinished. | Highest. This is the first screen every day. | Low. Copy and one status card. | App Store |
| 2 | After Family Controls permission, open the app picker before setup can finish. Settings row shows "No apps selected" or "Blocking N apps", not "Authorized". | Highest. Without this the product never starts. | Medium. One more onboarding step. | App Store |
| 3 | Partner row says Waiting until they accept. Success text stops saying they are connected. Partner with an account gets an Accept button. | Highest for a couple. Also the P0 bug. | Medium. Needs the server to tell the truth. | Server, then App Store |
| 4 | A real reset-password page, and the app guide stops promising a link that 404s. | High. Recovers locked-out people without a binary. | Low. | Website and server |
| 5 | Website: delete "notified when a blocked category is attempted" and "time windows you set" until those exist. Invite page: sign-in does not accept. | High trust. Stops a partner from expecting an alert the app cannot send. | Low. | Website |
| 6 | A simple inbox for the partner, using the server's summary, including self-harm and violence. Badge clears when they open it. | High. The partner currently has nowhere to look. | Medium. | App Store |
| 7 | Do not auto-play the tour until apps are chosen. First tour step is the restrictions screen. | High for first-run. | Low. | App Store |
| 8 | Create a group from a name. Covenant later. Join page's first sentence: this link does not add you. | Medium. Removes a long form. | Low. | App Store and website |
| 9 | Sign out asks first. One line that blocking on this phone is separate from the account. Clear the previous person's queue and selection on sign-out (master M21). | Medium. Prevents a shared-phone mixup during friends-and-family. | Low. | App Store |
| 10 | Toggles carry their own labels. Lock and star buttons have names. Footnotes at 11 points and 30% white get bigger and brighter. Do not hide the app from VoiceOver during the tour. | Medium for anyone using larger text or VoiceOver. | Medium. Do it with the home-screen build. | App Store |
| 11 | 18+ checkbox. Terms and privacy are links. "Email already registered" offers Sign In. | Medium. Honest, and fewer stuck signups. | Low. | App Store |
| 12 | Lockout subtitle in plain language. PIN pad says the partner is the one who knows it. Support screen mentions a group, not only a partner. | Lower. Removes jargon at a stressed moment. | Low. | App Store |
| 13 | New partner's first screen: "You're paired with {name}." Skip asking them to invite someone else. | Medium for the partner's first launch. | Low once invite redemption returns the inviter's name. | App Store and server |
| 14 | Privacy opener matches the Apple / Stripe / Analytics list. About page stops saying $4 per person. | Medium trust for donors and the App Store privacy review. | Low. | Website |

Items 1, 2, 3, 6, and 7 belong in one app build. Item 4 and item 5 should not wait for that build.

---

## What is already easy

Worth keeping:

- Welcome is one idea and one button.
- The group empty state says what to do.
- The feedback sheet is labeled, optional about device info, and has an email fallback.
- The shield footnote tells the truth about the Screen Time override.
- The invite pages say 18+ and not parental controls.
- Skip exists on the tour and on the partner step.
- Full-width gold buttons are large enough to hit.
- "I Need Support Right Now" is plain language. The tour explains it as an urgent notice. Keep that name.
