# Round 2 master review — Remain Faithful

Read-only cross-review of the three round-1 reports against commit `784998f`. No product code was changed.

Reviewer method: every claim below was checked in this checkout. Line numbers are from that tree. Items that depend on Fly secrets, the live App Store listing, or a physical iPhone are marked **suspected**. This machine cannot run the iOS app. `govulncheck` and `npm audit` were not re-run; version pins were re-read, and advisory counts carried from round 1 are labeled as such.

Severity used for the re-ranked list:

| | Meaning |
|---|---|
| **P0** | Someone can take an account, or the product tells both people they are protected when alerts will never arrive, or a stranger can receive another adult's flags. |
| **P1** | A serious hole in recovery, in what a partner can see, or in the "we will notice if monitoring stops" promise. |
| **P2** | Real, narrower. Fix it, but it does not by itself make the accountability promise false for a normal pair. |
| **P3** | Hygiene. |

Fixed product choices are still not defects: Family Controls plus a partner notice, DeviceActivity, opt-in Deep Scan via ReplayKit (it cannot see DRM video), and App Lockout via the Darwin notification. Where the words describe something the code does not do, the defect is the mismatch.

---

## What changed in the ranking

Round 1 was right about the three worst holes. A few P0 and P1 labels were higher than the code supports.

- Apple sign-in forgery, the stuck partner invite, and adding someone to a group with no consent stay **P0**.
- Password reset is completely broken in this repo. It is a **P1** launch blocker for anyone who forgets a password. It is not account takeover, so it is not P0.
- Airplane mode dropping a protection alert is real. It is **P1**, not P0. The setting changes on the phone first, and the notice is a best-effort push with no saved copy. A determined person can also skip the notice with a modified app. The honest-user failure is still serious.
- "Users who never heartbeat are invisible" is only half right. Heartbeats exist only during Deep Scan. Treating a missing heartbeat as "monitoring is off" would cry wolf for everyone who only uses Screen Time.
- The login timing oracle is real and weaker than the report says. Registration already answers "email already registered" in plain text.
- ChatGPT's "Fort Wayne" line is **stale**. It is not in the privacy page.
- Claude's closing claim that every email template HTML-escapes its values is **wrong**. Only the contact email does.

---

## Rulings — Grok

| # | Ruling | Why |
|---|---|---|
| 1 Apple token | **CONFIRMED** P0 | `decodeAppleIdentityToken` base64-decodes the payload and returns `sub` and `email`. No signature, issuer, audience, or expiry check. `backend/internal/handler/social_auth.go:149-180`. A matching email overwrites `apple_id` and a normal JWT is issued (`:118-126`, `:47-60`). |
| 2 Partner invite stuck pending | **CONFIRMED** P0 | Existing-user insert omits `status`, so the default `pending` applies (`invites.go:62-66`, `main.go:396`). The handler still returns `"connected"` (`invites.go:77-79`). The accept token is never stored (`:72-73`). Fan-out requires `accepted` (`events.go:136-139`). The partner list does not render `status` (`SettingsView.swift:898-923`, loaded at `:1052-1061`). |
| 3 Categories dropped; no 1:1 inbox | **CONFIRMED** P1 | `EventCategory` maps only adult content, gambling, social media, and gaming (`DashboardView.swift:22-29`). `from(remote:)` drops anything else (`:165-166`). The server accepts `violence`, `self_harm`, `clean`, and `app_usage` (`events.go:43-47`). `listAlerts` is used only inside a group member sheet (`GroupView.swift:816-820`). Home marks every alert seen (`ContentView.swift:81-95`, `:166-168`). Social media and gaming are in the app enum and are not server categories, so those rows never arrive. |
| 4 Password reset dead | **CONFIRMED** P1 | Default link host is `https://app.remainfaithful.app` (`auth.go:395-399`). `APP_BASE_URL` is not in `backend/fly.toml`. Whether production sets it is **suspected**. There is no `website/src/app/reset-password` page either way. The app has `forgotPassword` and `resetPassword` and no screen that accepts the token (`LoginView.swift:467-471`, `APIClient.swift:363-369`). The app guide still says the link works for one hour (`website/src/app/app-guide/guide-content.ts:109-110`). Reset does not revoke JWTs (`auth.go:375-386`). |
| 5 Monitoring can be turned off quietly | **CONFIRMED** P1, one clause partly true | `setMonitoringEnabled(false)` does not send `monitoring_disabled` (`ActivitySelectionManager.swift:96-150`). No Swift file sends that string. Emptying the selection turns monitoring off inside `startMonitoring` (`:100-105`) with no alert. Sign-out is `clearSession()` with no notice (`SettingsView.swift:403-407`). `pendingRawEvents` is not cleared, and the next account's token uploads them (`ActivitySelectionManager.swift:158-163`, `AuthState.swift:83-92`). Family Controls revocation notifies only on an in-session change (`RemainFaithfulApp.swift:92-104`). Heartbeats are sent only by the broadcast extension, after a two-minute sleep, then every two minutes (`SampleHandler.swift:686-713`). `APIClient.heartbeat` is never called. The sweep ignores `last_heartbeat_at IS NULL` and non-accepted rows (`heartbeat_sweep.go:63-65`). `min_machines_running = 0` (`fly.toml:14-16`) stops the in-process sweep (`main.go:71-73`). The Screen Time passcode override is already disclosed in the app (`ScreenTimeMonitoringView.swift:219`). There is no ShieldAction target. |
| 6 Public copy overstates blocking and time windows | **CONFIRMED** P1 | No ShieldAction extension in the project. Schedule is 00:00–23:59 (`ActivitySelectionManager.swift:108-112`) with a one-minute daily threshold (`:119-130`). No hour picker. "Partners are notified when a blocked category is attempted" appears at `website/src/app/page.tsx:121` and `:197`, `about/page.tsx:198`, `privacy-architecture/page.tsx:69`, `how-it-works/page.tsx:30` and `:175`, `group-setup-guide/page.tsx:52` and `:158`, `structured-data.ts:144`, and `privacy/page.tsx:49`. Time windows: `how-it-works/page.tsx:11`, `:30`, `:222`, and `structured-data.ts:144`. `extendedDescription` invents a browsing story (`DashboardView.swift:141-154`, shown at `AlertDetailView.swift:123`). Whether App Store Connect still shows the paste in `docs/APP_STORE_SUBMISSION.md` is **suspected**. DRM limits on the site are accurate and should stay. |
| 7 Group add with no accept | **CONFIRMED** P0 | `inviteEmailToGroup` inserts into `group_members` immediately (`invites.go:281-305`). Any member may invite (`:242-251`). `CreateEvent` fans out to every co-member with no status check (`events.go:141-145`). `GetGroup` returns name, email, 30-day count, and streak to any member (`groups.go:122-158`). The cap check is outside a transaction (`invites.go:254-260` then `:287-292`). The emailed URL for an existing user has no token (`:301`). |
| 8 Invite mail HTML injection and no rate limit | **CONFIRMED** P2 (Grok called it P1) | Partner, group, and reset HTML are `fmt.Sprintf` with no `html.EscapeString` (`email.go:42-52`, `:65-74`, `:118-123`). Contact mail does escape (`:95-100`). Names are only `TrimSpace` (`auth.go:35`). Subjects use the raw name (`email.go:33`, `:58`). Invite routes are not wrapped in `authenticatedRateLimiter` (`main.go:269`, `:288-289`). `/auth` is 10/min/IP (`main.go:113-117`). SendGrid From falls back to `noreply@remainfaithful.app` (`email.go:27`). Whether that domain is verified at SendGrid is **suspected**. Severity is P2: this is injection into an email a partner opens, plus an unbounded mail bill, not account takeover. |
| 9 Refresh leaves the extension token stale | **CONFIRMED** P1 | `refreshTokenIfNeeded` writes the Keychain only (`APIClient.swift:186-192`). `AuthState.token` prefers memory (`AuthState.swift:41`). Only `setSession` writes the app-group token (`:48-49`). The extension reads that copy (`SampleHandler.swift:541-544`, `:699-700`). The privacy page says the JWT is not in UserDefaults (`privacy/page.tsx:152-154`). `setSession` does write it there. |
| 10 Old Next.js and Go | **PARTLY TRUE** P2 | `next@14.2.5` is pinned (`website/package.json:13`). `golang:1.22-alpine` and `alpine:3.19` are the image (`backend/Dockerfile:2`, `:18`). `jwt/v5 v5.2.1` is pinned (`backend/go.mod:5`). There is no `middleware.ts`, so a middleware auth-bypass advisory does not apply to a file this app does not have. `next/image` remote patterns are empty. The advisory counts and GO-2025-3553 were not re-run here. Treat "critical RCE" as not confirmed. The actionable fact is the old toolchain. |
| 11 Tier 1 labels every hit as sexual | **CONFIRMED** P1 | Keyword list includes gambling, gore, and suicide phrases (`SampleHandler.swift:132-142`). URL list includes betting and gore hosts (`:127-129`). Both paths post `.explicitSexual` (`:457-466`), which is the string `adult_content` (`:27`). Tier 2 can still emit the right category, and finding 3 then hides that one in the app. |
| 12 Email change with no verification | **CONFIRMED** P1 | `UpdateMe` writes the new email immediately (`users.go:71-75`). A unique clash becomes a generic 500 (`:76-78`). No confirmation link. |
| 13 Google token checked for `sub` only | **CONFIRMED** P1 | `verifyGoogleIDToken` reads `sub`, `email`, and `name` only (`social_auth.go:197-208`). No `aud`, `azp`, or `email_verified`. The token is concatenated into the query string (`:186-188`). `http.Get` has no timeout. This is not the unsigned forgery in finding 1. It is a token minted for some other Google client. |
| 14 Logs hold more than the privacy sentence | **CONFIRMED** P2 | OCR prefix is logged (`SampleHandler.swift:452`). A matched domain is logged (`:458`). A keyword match logs 80 characters of OCR (`:464`). Swift `Logger` interpolations are redacted unless a debugger is attached, so this is not a server upload. Protection `detail` is logged up to 80 runes (`protection.go:166-168`). Invite handlers log the recipient (`invites.go:75`, `:102`). SendGrid logs the recipient and the response body (`email.go:158-170`). These do not leave the phone as screen content. |
| 15 Enumeration, PIN attempts, timestamps | **CONFIRMED** P2 | Register returns 409 "email already registered" (`auth.go:60-62`). `clientIP` trusts `X-Forwarded-For` when `Fly-Client-IP` is absent (`main.go:196-202`). Spoofing that by reaching the VM outside Fly is **suspected**. PIN verify is 5 attempts per 15 minutes in memory (`protection.go:20-26`, `:87-98`, `:259-267`). `CreateEvent` stores a client timestamp with no window (`events.go:107-113`). Write limit is 30/min/user (`main.go:290`). No retention job. |
| 16 Privacy policy does not match deletion, analytics, or crypto | **CONFIRMED** P2 | "Deleted within 30 days" and "donation records retained 7 years" (`privacy/page.tsx:105-112`). `DeleteMe` deletes the user immediately (`users.go:118-119`). `donations.user_id` is `ON DELETE CASCADE` (`main.go:461`). The rights section also says deletion is immediate (`privacy/page.tsx:134`). There is no deletion-confirmation email. GA4 loads whenever `NEXT_PUBLIC_GA_ID` is set (`layout.tsx:107-120`). Nothing reads Do Not Track. Whether the env var is set in production is **suspected**. "TLS 1.3" and "we do not log sensitive data" (`privacy/page.tsx:151-154`). `force_https` is on. The app does not pin TLS 1.3. AES-256 and "end-to-end" are at `privacy-architecture/page.tsx:266` and `:289`. The server sees plaintext JSON. Disk encryption is not configured in this repo. |
| 17 "App removed" goes to the dead token | **CONFIRMED** P2 | On invalid token, `notifyPartners` marks it inactive and then sends `APP_DELETED` to that same token (`notifications.go:343-364`). The partner path `notifyPartnerByID` only marks the token inactive (`:397-402`) and does not notify anyone else. |
| 18 Cost is not $4 per person | **CONFIRMED** on the copy, **PARTLY TRUE** on the bill | The about page says about $4 per participant and also "~$80/mo" (`about/page.tsx:213-218`). Those two lines only agree near 20 people. Heartbeats are not 720/day for a normal user: they run only during Deep Scan, first after two minutes (`SampleHandler.swift:686-692`). Grok's dollar ranges and Claude's ranges are both models, not invoices. The code facts that drive cost are: no retention, uncapped partners, uncapped invite mail, and a 256 MB machine that is allowed to sleep. |
| 19 Smaller notes | See the P3 table | Lint really has no ESLint config (confirmed: no `.eslintrc*` under `website/`). `001_schema.sql` lagging `migrate()` was not re-opened line by line; `migrate()` in `main.go` is the live schema. Age: registration collects name, email, and password only (`auth.go:25-38`). The about page ministry sentence was re-read; it does not name a church, city, state, or street. |

---

## Rulings — Claude

| ID | Ruling | Why |
|---|---|---|
| F-01 Apple takeover | **CONFIRMED** P0 | Same code as Grok 1. The forged-token reproduction was Claude's, on a copied function. This review re-read the function. It does not verify a signature. |
| F-02 Airplane mode drops protection alerts | **CONFIRMED**, severity **overstated** as P0. Rank **P1**. | The seven call sites use `try?` and do not queue: `SettingsView.swift:310` and `:853`, `AppLockoutManager.swift:95`, `ActivitySelectionManager.swift:69`, `RemainFaithfulApp.swift:99`, `PartnerPINManager.swift:65`, `PINEntryView.swift:393`. Content events do queue (`SampleHandler` retry around `:585`). The local change happens first. This is a real gap for an honest user who is offline. It is not a silent failure of every session. |
| F-03 No durable protection or silence row | **CONFIRMED** P1 | `sendProtectionAlertToPartners` pushes and writes no alert row (`protection.go:142-193`). Silence does the same (`heartbeat_sweep.go:108-131`). `ListAlerts` inner-joins `events` (`alerts.go:85-86`). `event_id` is `NOT NULL` (`main.go:426` area). A missed push is gone. |
| F-04 Partner invite pending forever | **CONFIRMED** P0 | Same as Grok 2. `"connected"` is what makes it P0. |
| F-05 Group add with no consent | **CONFIRMED** P0 | Same as Grok 7. This is the sensitive-data disclosure. |
| F-06 Password reset | **CONFIRMED**, rank **P1** not P0 | Same as Grok 4. Recovery is dead in this repo. The account still works if the person remembers the password. |
| F-07 Google `aud` | **CONFIRMED** P1 | Same as Grok 13. |
| F-08 bcrypt timing oracle | **CONFIRMED**, severity **overstated** as P1. Rank **P2**. | `"$2a$10$placeholder"` is not a bcrypt hash (`auth.go:241`). Compare returns immediately. Claude's millisecond ratio was not re-measured here. Registration already returns 409 for a known email (`auth.go:60-62`), which is a louder oracle. |
| F-09 `CreateRelationship` enumeration | **CONFIRMED** P2 | 404 "no account found for that email" (`relationships.go:45-47`). The insert that follows is `pending`, so it does not by itself start alert fan-out. `ListRelationships` still returns the other person's name and email for that pending row (`relationships.go:88-101`). The iOS UI calls `invitePartner`, not this route. The route is still live (`main.go:267`) and has no per-user limiter. |
| F-10 No email verification; signup auto-accepts invites | **CONFIRMED** P1 | `Register` calls `redeemPendingInvites` (`auth.go:68-71`), which sets partner rows to `accepted` (`:124-128`) and joins groups under the cap lock (`:172-176`). No verification email exists. |
| F-11 PIN brute force | **PARTLY TRUE** P2 | The server returns `success: false` and notifies nobody (`protection.go:295-296`). The app sends `wrong_pin_attempt` itself, with `try?` (`PartnerPINManager.swift:65`). The limiter is in memory and resets on restart (`protection.go:23-26`). A 4-digit PIN at 5 tries per 15 minutes is slow while one process stays up. Scale-to-zero makes restarts ordinary, which weakens the limiter. Knowing the PIN does not disable monitoring on the phone unless the person is in the app. The PIN is a speed bump, not a lock. |
| F-12 Machines sleep | **CONFIRMED** P1 | `fly.toml:14-16`. Sweeps start in `main.go:71-73`. The silence sweep cannot notice silence while the machine is stopped. |
| F-13 Refresh desyncs extensions | **CONFIRMED** P1 | Same as Grok 9. `sendHeartbeat` logs a 401 and continues (`SampleHandler.swift:710-713`). |
| F-14 Silence alerts once | **PARTLY TRUE** P2 | The SQL and the comment say one alert per silence episode until a new heartbeat (`heartbeat_sweep.go:16-26`, `:59-68`, `:96-101`). That is the written design. The defect is the design plus F-03: if the one push is missed, nothing is saved and nothing repeats. |
| F-15 Never-started users are invisible | **PARTLY TRUE**, not P1 | Both watchdogs require `last_heartbeat_at IS NOT NULL` (`heartbeat_sweep.go:63`, `heartbeat.go:68`). Heartbeats are Deep Scan only. A Screen Time user who never starts a broadcast would look "never started" if this fired. The real gap is that a partner cannot see whether Screen Time monitoring is on. Do not alert on a null heartbeat. |
| F-16 `detail` logged | **CONFIRMED** P2 | Same as Grok 14. Names are in today's client strings (`SettingsView.swift:312`). The APNs body ignores `detail` (`protection.go:198`). |
| F-17 Email addresses and reset links in logs | **CONFIRMED** P2 | `email.go:128-131` prints the full plain body, which includes the reset URL, when `SENDGRID_API_KEY` is empty, and returns `nil`. Normal sends log the address and the response body (`:158`, `:170`). Whether production has the key set is **suspected**. The nil return on a missing key is the sharp edge. |
| F-18 Tokens not revoked on password reset | **CONFIRMED** P2 | Stateless 24-hour JWT (`jwt.go:20-32`). Refresh re-signs any valid token (`auth.go:266-280`). Reset updates the password and deletes the reset row only (`auth.go:375-386`). No "sign out everywhere". |
| F-19 JWT in App Group UserDefaults | **CONFIRMED** P2 | `AuthState.swift:48-49`. Extensions need a token they can read. The privacy sentence is still false. A Keychain access group is the right fix and needs an app build. |
| F-20 Throttle can drop a severe push | **CONFIRMED** P2 | The count is per recipient, any monitored person, any severity (`events.go:192-202`). The alert row is still inserted. The push is what gets skipped. Combined with finding 3, a `self_harm` row can exist and still be invisible in the app. |
| F-21 Lint unconfigured; CI is Go tests only | **CONFIRMED** P3 | No ESLint config under `website/`. `package.json` script is `next lint`. `.github/workflows/backend-ci.yml` runs `go test ./...` only. The pre-commit hook is opt-in per clone. |
| F-22 Group cap race | **CONFIRMED** P3 | Count then insert with no lock (`invites.go:254-260`, `:287-292`). Signup path locks (`auth.go:189-202`). Exceeding 12 by a race is real and narrow. |
| F-23 Partner-side row cannot be deleted | **CONFIRMED** P2 | `DeleteRelationship` requires `user_id = caller` (`relationships.go:210-212`). `ListRelationships` returns both directions (`:88-101`). The Remove button is on every row (`SettingsView.swift:953-962`) and surfaces the error in `errorMsg` (`:834-835`). A partner cannot step out. The privacy page says either person can end it (`privacy/page.tsx:110`). |
| F-24 `VerifyPIN` uses `LIMIT 1` | **CONFIRMED** P2 | `protection.go:284-288`. `GetPINStatus` is true if any row has a hash (`:242-245`). Two PINs: one of them works, and which one is unspecified. |
| F-25 AES-256 claim | **CONFIRMED** P2 | `privacy-architecture/page.tsx:266`. No application encryption in the Go tree. The pre-commit scanning sentence is only true for clones that opted into the hook. |
| F-26 `/api/donate` unthrottled | **CONFIRMED** P2 | `website/src/app/api/donate/route.ts` creates a Stripe session with no limiter and no honeypot. Return URLs are allowlisted. `/api/waitlist` returns success and stores nothing (`waitlist/route.ts:9-14`). `/api/contact` has a honeypot; if `BACKEND_URL` is unset it logs the name and email and still returns success (`contact/route.ts:15-22`). |
| F-27 No partner cap | **CONFIRMED** P2 | Groups cap at 12. Relationships do not. `POST /relationships` and `/relationships/invite` are not on the authenticated limiter (`main.go:267-269`). |
| F-28 Reset email bombing | **CONFIRMED** P2 | Forgot-password is 10/min/IP only (`main.go:256`, limiter at `:113-117`). Each request replaces the previous token (`auth.go:316-321`). |
| F-29 `UpdateMe` email | **CONFIRMED** P1 | Same as Grok 12. Format is not checked. Refresh keeps the old email inside the JWT (`auth.go:274-275`). |
| F-30 `APP_DELETED` to the dead token | **CONFIRMED** P2 | Same as Grok 17. |
| F-31 Group streak counts clean events | **CONFIRMED** P2 | `GetGroup` uses `MAX(e.timestamp)` over all events (`groups.go:147-156`). `GetUserStats` excludes `clean` (`users.go:277`). `flags_last_30` is `COUNT(*)` with no category filter (`groups.go:149`), so clean events also inflate the flag count. That second half was easy to miss. |
| F-32 Limit 100, no pagination | **CONFIRMED** P3 | `events.go:289`, `alerts.go:89`. |
| F-33 No `(user_id, timestamp)` index | **CONFIRMED** P3 | `idx_events_user_id` only (`main.go:432`). |
| F-34 No heartbeat index | **CONFIRMED** P3 | Sweep filters `last_heartbeat_at` (`heartbeat_sweep.go:63`). No such index in `migrate()`. |
| F-35 Client timestamp | **CONFIRMED** P2 | Same as Grok 15. Bad strings become 500 because the insert fails. |
| F-36 Reset tokens plaintext | **CONFIRMED** P3 | `auth.go:313-321`, table at `main.go:451-457`. |
| F-37 `markTokenInactive` by token only | **CONFIRMED** P3 | `notifications.go:413-415`. `RegisterDeviceToken` stores the posted string. |
| F-38 Per-push logging | **CONFIRMED** P3 | `notifications.go:409` logs every partner send. |
| F-39 Any member can rename and rewrite the covenant | **CONFIRMED** P3 | Role is read and discarded (`groups.go:232-240`). |
| F-40 Sole admin leaving does not promote | **CONFIRMED** P3 | `LeaveGroup` deletes the row (`groups.go:441-443`). `DeleteMe` does promote (`users.go:132-176`). |
| F-41 Missing canonicals | **Not re-checked page by page.** | Treat as a lead, not a confirmed count of seven. |
| F-42 Google `http.Get` and query string | **CONFIRMED** P3, folded into the Google P1 | `social_auth.go:186-188`. |
| F-43 Column name via `fmt.Sprintf` | **CONFIRMED** P3 | `social_auth.go:112`, `:125`, `:136`. Callers pass `apple_id` or `google_id` only. Not injectable today. |
| F-44 `X-Forwarded-For` fallback | **CONFIRMED** P3 | `main.go:196-202`. Off-proxy spoofing is **suspected**. |
| F-45 Compiled binary | **STALE / withdrawn** | Claude already withdrew it. Not re-opened. |
| "What is solid" — email templates escape everything | **WRONG** | Only `SendContact` escapes. Invites and password reset do not. |
| "What is solid" — process refuses to start without `JWT_SECRET` | **PARTLY TRUE** | `secret()` panics on first use (`jwt.go:60-65`), not at process start. The comment says "refusing to start". |
| App Store 3.1.1 clear | **CONFIRMED** | No donation UI in the Swift targets. Keep `/app-guide` free of donate buttons. The API still has checkout routes (`main.go:304-305`). That is not in the binary. |
| Time Sensitive entitlement | **PARTLY TRUE / suspected effect** | Panic sets `interruption-level: time-sensitive` (`notifications.go:171`). `RemainFaithful.entitlements` has no time-sensitive key. Whether APNs drops or downgrades that push was not tested on a device. **Suspected.** |

The privacy promise Claude praised is still true for screen content. `uploadEvent` sends category, severity, summary, and timestamp (`SampleHandler.swift:555-556`). Summaries are allowlisted (`events.go:22-41`). That should stay.

---

## Rulings — ChatGPT

ChatGPT could not read the code. These were leads.

| # | Ruling | Why |
|---|---|---|
| 1 Privacy page ends with "Fort Wayne" | **STALE** | The contact block is email and GitHub only (`privacy/page.tsx:163-168`). No city, church, or street on that page. |
| 2 "Never disclosed to any third party" vs APNs, GA, Stripe | **CONFIRMED** P2 | The opening sentence is absolute (`privacy/page.tsx:20-22`). The same page then lists Apple, Stripe, and Google Analytics (`:115-126`). It is a contradiction in the copy, not a hidden sale of data. P2, not P1. |
| 3 Heartbeat path needs an end-to-end check | **CONFIRMED** as the code path; live delivery **not tested** | Sweep, watcher, and protection type exist. They only arm after a Deep Scan heartbeat, and the machine may be asleep. See master items on silence. |
| 4 MF-02 Family Controls distribution entitlement | **STALE as a code blocker** | `MANUAL_FIXES_NEEDED.md` still lists it. The entitlement key is in the app, shield, and Device Activity targets. This review did not open Apple's portal. Do not treat the markdown checklist as proof the entitlement is missing. |
| 5 Waitlist returns success without saving | **CONFIRMED** P3 | `website/src/app/api/waitlist/route.ts:9-14`. It logs the email and returns success. |
| 6 Outbound data boundary | **PARTLY TRUE** | Spot-checked: Deep Scan upload and DeviceActivity post send category, severity, summary, and timestamp (`SampleHandler.swift:555-556`, `DeviceActivityMonitorExtension.swift:62-70`). No screenshot upload was found on those paths. A line-by-line audit of every `URLSession` in every target was not repeated beyond those. |
| 7 IDOR on every lookup | **PARTLY TRUE** | `GetGroup` requires membership (`groups.go:122-129`). `ListEvents` is the caller's rows. `ListAlerts` is the recipient. The real holes are group force-add and one-sided delete, already listed. There is not a general missing-ownership bug on every route. |
| 8 Invite token lifecycle and the 12 cap | **CONFIRMED** as several separate bugs | New-user partner tokens last 7 days and are one-time via `status = pending` (`invites.go:126-131`). Existing-user partner tokens are never stored. Group direct-add has no token. The cap is transactional only on the signup path. Invite routes are not rate-limited. |
| 9 Silent monitoring defeat | **CONFIRMED** P1 | Same family as Grok 5 and Claude F-02, F-12, F-13. Force-quit does not lift a ManagedSettings shield. That part of the website is right. Turning the in-app toggle off, signing out, or losing the network does not reliably tell a partner. |
| 10 Scale math of 96 heartbeats/user/day | **WRONG** | The code sleeps two minutes, then posts, and only while a broadcast is running (`SampleHandler.swift:686-692`). It is not a 5-minute heartbeat for 8 hours. Missing indexes and unbounded retention are real (F-33, F-34, Grok 18). |

---

## Missed by all three

1. **The home screen tells the user monitoring is paused whenever Deep Scan is off.** `BroadcastPausedBanner` says "Monitoring is paused" and "Tap to resume accountability" (`DashboardView.swift:928-931`). The only instructions are the four Control Center broadcast steps (`:961-966`). Screen Time blocking can be on the whole time. The empty activity state says "No flags yet — great work" (`:786-793`) even when no apps were chosen and monitoring never started. That is a false calm, and it is the first thing a new user sees.

2. **Onboarding never asks which apps to block.** `MonitoringSetupStep` requests Family Controls permission and then finishes (`OnboardingView.swift:478-498`). There is no `FamilyActivityPicker`. Monitoring cannot start with an empty selection (`ActivitySelectionManager.swift:100-105`). Settings then says "Authorized", which sounds finished (`SettingsView.swift:231-237`).

3. **Protection and silence alerts do not go to group members.** Content events union group co-members (`events.go:141-145`). Protection alerts query `relationships` only (`protection.go:145-148`). A person whose only accountability is a group never gets "monitoring turned off" or "device went quiet", even after those alerts are made reliable.

4. **Sign-out leaves the previous person's Screen Time selection and shield flags in the app group.** `clearSession` removes the token and the selected group id (`AuthState.swift:83-92`). It does not clear `screenTimeSelection`, `deviceActivityMonitoringEnabled`, `screenTimeShieldingEnabled`, `appLockoutEnabled`, or `pendingRawEvents`. The next account on that phone inherits the shield and can upload the queued events. Grok noted the queued events. The leftover shield was not called out.

5. **The support button can succeed for a group and the screen can still say you need a partner.** `PanicView` shows "Add an accountability partner" when no accepted relationship is loaded (`PanicView.swift:185-187`). `SendPanicAlert` notifies accepted partners and group members (`notifications.go:200-215`) and only errors when both are empty (`:155-157`). A group-only user sees a dead-end sentence above a button that can work.

---

## Master list

Ship means what has to change before the fix is real for a user. "Server" is the Fly API. "Website" is the Next.js site. "App Store" is a new iOS binary.

### P0

**M1. Sign in with Apple does not check the identity token.**
Evidence: `social_auth.go:149-180`, `:118-126`.
A forged payload with the victim's email attaches the attacker's Apple id and returns a 24-hour session. Password reset does not remove that Apple id.
Fix: verify Apple's signature against `https://appleid.apple.com/auth/keys`. Require issuer `https://appleid.apple.com`, audience equal to the app's bundle id, and a valid expiry. Link to an existing email only when that person is already signed in, or when the row's `apple_id` is empty and the email is verified by Apple. Never overwrite a different `apple_id`.
Ship: **server**.

**M2. Inviting someone who already has an account does not connect them, and the app says it did.**
Evidence: `invites.go:62-80`, `main.go:396`, `events.go:136-139`, `SettingsView.swift:898-923`.
The row stays `pending`. Alerts, silence checks, and partner PIN all require `accepted`. The accept link in the email points at a token that was never saved. The website page ignores the token (`accept-invite/page.tsx:12-14`).
Fix: keep the row pending, return `"pending"`, store the token, and accept it only when that person taps Accept while signed in as that email. Show "Waiting for them" versus "Connected" in Manage Partners. Backfill rows that are stuck pending: do not flip them to accepted without that person accepting.
Ship: **server now** (stop returning connected, store the token). **App Store** to show the status and an Accept button. A website Accept page only helps if it can sign the person in, which it cannot today.

**M3. Anyone who knows an email can put that account in a group and then receive their flags.**
Evidence: `invites.go:281-305`, `events.go:141-145`, `groups.go:142-158`.
Fix: write a pending group invite for an existing account. Add the member only after they accept, inside the same locked cap check signup already uses (`auth.go:181-217`). Until they accept, do not fan events to them and do not show them the roster.
Ship: **server** stops the leak. **App Store** for an Accept control. The website invite page should say "you have not joined until you accept", not "sign in and you are in".

### P1

**M4. Password reset links go nowhere.**
Evidence: `auth.go:327-328`, `:395-399`. No reset page in `website/src`.
Fix: add `/reset-password` on the real site. It reads the token, asks for a new password, and posts to `/auth/reset-password`. Set `APP_BASE_URL` to that site and change the code default so a missing env var cannot fall back to `app.remainfaithful.app`. On success, invalidate outstanding sessions (see M18).
Ship: **website and server**. No app build required.

**M5. The app hides violence, self-harm, and app-use alerts, and a one-to-one partner has no inbox.**
Evidence: `DashboardView.swift:22-29`, `:165-166`. `GroupView.swift:816-820`. `ContentView.swift:166-168`. `NotificationService.swift:124-128` sends an unmapped push to Home.
Fix: map `violence`, `self_harm`, and `app_usage` to their own rows. Add a partner inbox on `GET /alerts` for one-to-one partners. Mark an alert seen when that person opens it, not when Home appears. Show the server summary. Delete `extendedDescription`'s invented browsing story.
Ship: **App Store**. The rows are already stored.

**M6. Turning protection off can be missed, and there is no record if the push is missed.**
Evidence: the `try?` sites listed under F-02. `protection.go:142-193`. `alerts.go:85-86`. `ActivitySelectionManager.swift:96-150`. `SettingsView.swift:403-407`. `RemainFaithfulApp.swift:92-104`.
Fix: write the protection alert to the app group before the setting changes, and retry until the server returns success. On the server, store a row for every protection and silence alert and return it from `GET /alerts`. Send `monitoring_disabled` when monitoring is turned off or the selection becomes empty. Send it on sign-out. On launch, if Family Controls was approved and is no longer approved, send `family_controls_revoked`. Drop `pendingRawEvents` on sign-out, or stamp them with the user id.
Ship: **App Store and server**.

**M7. Silence detection does not cover the way people actually use the app, and the sweeper can be asleep.**
Evidence: `SampleHandler.swift:686-713`. `heartbeat_sweep.go:59-68`. `fly.toml:14-16`. `protection.go:145-148` (1:1 only).
Fix: set `min_machines_running = 1` so the process stays up. That is a small monthly bill, not a new product. Send a lightweight heartbeat from the main app or the Device Activity extension while Screen Time monitoring is on, so silence means "monitoring stopped", not "Deep Scan stopped". Re-alert on a slow schedule while silence continues, and store each one (M6). Include group members, not only 1:1 partners. Do not alert merely because `last_heartbeat_at` is null until a heartbeat has been armed on purpose.
Ship: **server** for the machine and the fan-out. **App Store** if heartbeats move out of Deep Scan.

**M8. The website, and the home screen, describe a product the code does not implement.**
Evidence: the pages in Grok 6. `DashboardView.swift:928-931` and `:786-793`. `ActivitySelectionManager.swift:108-112`.
Fix until ShieldAction and real time windows exist: say that a shield blocks chosen apps, that usage is one alert after one minute per day, and that Deep Scan is opt-in and separate. On the home screen, say whether apps are blocked, whether the daily usage alert is on, and whether Deep Scan is running. Never say "great work" when setup is unfinished.
Ship: **website** for the pages. **App Store** for the home screen. App Store Connect text is a listing edit, not a binary, if that paste is what is live (**suspected**).

**M9. After about a day, Deep Scan uploads and heartbeats can die quietly.**
Evidence: `APIClient.swift:186-192` versus `AuthState.swift:43-49` and `SampleHandler.swift:541-544`, `:699-713`.
Fix: on refresh, update memory, Keychain, and the app-group token together, through `setSession`. Treat HTTP 401 from the extension as "signed out, monitoring has stopped" and show that on Home.
Ship: **App Store**.

**M10. A self-harm or gambling phrase found on screen is stored as adult content.**
Evidence: `SampleHandler.swift:132-142`, `:457-466`, `:27`.
Fix: pick the category from the pattern or domain that matched. Keep self-harm distinct.
Ship: **App Store**.

**M11. There is no proof the email address belongs to the person who signed up.**
Evidence: `auth.go:68-71`, `:124-128`. `users.go:71-75`.
Fix: send a confirmation link before `redeemPendingInvites` runs, and before a new email replaces the old one. Return a conflict when the address is taken, without extra detail. Changing email must not be enough to accept someone else's invite.
Ship: **server and website**.

**M12. Google sign-in accepts a token minted for a different app.**
Evidence: `social_auth.go:185-208`, combined with `:118-126`.
Fix: require `aud` (or `azp`) to be this app's Google client id, and require `email_verified`. Give the HTTP client a timeout. Do not put the token in a query string.
Ship: **server**.

### P2

**M13. Invite and reset emails can be stuffed with HTML, and invite mail is unlimited.**
Evidence: `email.go:42-52`, `:65-74`, `:118-123`. `main.go:269`, `:288-289`.
Fix: HTML-escape every interpolated value. Strip newlines from names before they touch a subject. Put invite routes on a small daily cap. If `SENDGRID_API_KEY` is missing in production, fail the send instead of printing the reset link and returning success (`email.go:128-131`).
Ship: **server**.

**M14. A stolen session keeps working after a password change.**
Evidence: `jwt.go:20-32`, `auth.go:266-280`, `:375-386`. The token is also in App Group UserDefaults (`AuthState.swift:48-49`), which contradicts `privacy/page.tsx:152-154`.
Fix: add a token version on the user, put it in the JWT, and bump it on password reset, password change, and "sign out everywhere". Prefer a Keychain access group over UserDefaults for the extension token.
Ship: **server** for revocation. **App Store** for the Keychain move. **Website** for the privacy sentence, which should be fixed even before the Keychain move so it matches today's app.

**M15. The partner PIN can be walked, and two partners' PINs collide.**
Evidence: `protection.go:259-297`, `:284-288`. `PartnerPINManager.swift:65`.
Fix: count failures in Postgres. Send `wrong_pin_attempt` from the server on failure. If one PIN per person, reject a second. If one PIN per partner, the app must send the relationship id. Return "no PIN" and "wrong PIN" as different results.
Ship: **server**, plus a small **App Store** change if the client must send the relationship id.

**M16. The person being asked to watch cannot leave.**
Evidence: `relationships.go:210-212`. `SettingsView.swift:953-962`. `privacy/page.tsx:110`.
Fix: allow delete when the caller is either `user_id` or `partner_id`, and notify the other person.
Ship: **server**. The button already exists.

**M17. "App was removed" is sent to the token Apple just rejected.**
Evidence: `notifications.go:343-364`.
Fix: resolve the partner from the dead token's user, and send the notice to that partner's other tokens. Do not send it back to the dead token.
Ship: **server**.

**M18. The group screen punishes a clean streak.**
Evidence: `groups.go:147-156` versus `users.go:273-278`.
Fix: exclude `category = 'clean'` from both the streak and `flags_last_30`. Use one SQL helper so the two screens cannot drift.
Ship: **server**.

**M19. Privacy copy promises things the code does not do.**
Evidence: `privacy/page.tsx:20-22`, `:61-64`, `:105-112`, `:151-154`. `privacy-architecture/page.tsx:266`, `:289`. `about/page.tsx:213-218`. `layout.tsx:107-120`.
Fix: delete the absolute "never disclosed" sentence, or name the processors in that same sentence. Say account deletion is immediate and that donation rows currently delete with the account, or stop cascading them and keep the 7-year record. Load GA only after a real opt-in, and remove the Do Not Track sentence until the code honors it. Say HTTPS to the API, not end-to-end, and drop AES-256 unless the database volume is actually configured that way. Replace "$4 per participant" with the real stack cost. Stop logging OCR text, matched domains, protection `detail`, and full email addresses (Grok 14).
Ship: **website** for the sentences. **Server** for the logs and, if you keep donation rows, the cascade.

**M20. A few writes have no ceiling.**
Evidence: uncapped partners (`invites.go` has no count). Client timestamps (`events.go:107-113`). No retention job. `/api/donate` (`donate/route.ts`). Forgot-password has no per-account cap (`auth.go:286-334`). Push throttle ignores severity (`events.go:192-202`).
Fix: cap partners (10 is enough). Reject timestamps outside a short skew of server time. Delete or archive events and alerts past a stated window. Rate-limit donate and add a honeypot. Cap reset mail per account (3 an hour is enough). Do not let an informational alert suppress a severe one; exempt `severe`, and scope the throttle to one monitored person.
Ship: **server** and **website**.

**M21. Sign-out does not reset monitoring on a shared phone.**
Evidence: `AuthState.swift:83-92`. Keys in `ActivitySelectionManager.swift:35-36` and `:158-163`. `AppLockoutManager.swift:31`.
Fix: on sign-out, stop monitoring, clear the shield, and delete the app-group selection, lockout flag, and pending event queue.
Ship: **App Store**.

**M22. Old Next.js and old Go.**
Evidence: `website/package.json:13`, `backend/Dockerfile:2` and `:18`, `backend/go.mod:5`.
Fix: upgrade Next.js to a current patched release and re-run `npm audit`. Rebuild the API with a current Go release and jwt at least v5.2.2. Do not treat the round-1 "1 critical" count as re-confirmed until that audit is run again.
Ship: **website and server**. No app build.

### P3

| Item | Evidence | Fix | Ship |
|---|---|---|---|
| ESLint is not configured, so `npm run lint` cannot run by itself. CI does not build the site. | `website/package.json:8`. No `.eslintrc*`. `backend-ci.yml` tests Go only. | Commit `{"extends":"next/core-web-vitals"}`. Add a website job: install, lint, build, test. Add `go vet` to the Go job. | Website / CI |
| History is the latest 100 rows. | `events.go:289`, `alerts.go:89` | Cursor pagination. | Server, then app |
| Missing indexes for events-by-time and heartbeat sweep. | `main.go:432`, `heartbeat_sweep.go:63` | Add them when the sweep matters. | Server |
| Reset tokens stored in plaintext. | `auth.go:313` | Store a hash. Compare hashes. | Server |
| Group cap can be raced past 12. | `invites.go:254-260` | Same row lock as signup. | Server |
| Any member can rename the group and rewrite the covenant. Sole admin can leave the group with no replacement. | `groups.go:232-240`, `:441-443` | Admin-only edits. Promote on leave, same as delete-account. | Server |
| `markTokenInactive` is global by token string. | `notifications.go:413-415` | Scope by user id. Check token length. | Server |
| SQL column name is interpolated. Safe today. | `social_auth.go:112` | Two literal queries. | Server |
| Rate limit can fall back to `X-Forwarded-For`. | `main.go:196-202` | In production, trust only `Fly-Client-IP`. | Server |
| Waitlist says success and saves nothing. | `waitlist/route.ts:9-14` | Store it, or remove the form. | Website |
| No in-app 18+ attestation. Policy says 18+. | `auth.go:25-38`, `privacy/page.tsx:145-148` | A checkbox before account creation. Not a substitute for the store rating. | App Store |
| Login timing padding does not run. | `auth.go:241` | Compare against a real dummy hash created once at startup. Also stop returning "email already registered" (M20's louder oracle). | Server |
| Canonicals, sitemap, and contact `htmlFor` from Grok's P3 list. | Not all re-counted. Contact form has no `htmlFor` match in `ContactForm.tsx`. | Fix labels. Do not block launch on canonicals. | Website |

---

## Checked, and still not a finding

- No founder city, state, named church, or street address on the pages read for this review, including the privacy contact block and the about ministry sentence.
- No donation UI in the iOS target. Website Stripe donations stay.
- Deep Scan does not upload frames. The summary allowlist holds.
- `heartbeat_silence` cannot be forged through the protection HTTP handler (`protection.go:69-73`).
- Partner PIN set and delete check `partner_id` (`protection.go:331-338`, `:372-376`).
- Password hashing is bcrypt. App JWTs pin HMAC and check issuer and expiry (`jwt.go:35-48`). That does not fix M1.
- Force-quitting the app does not clear a ManagedSettings shield. The website is right about that part.
- Group roster PII requires membership. The bug is how membership is granted (M3), not the read itself.

## What this round did not do

No iPhone build, no call to production, no `npm audit`, no `govulncheck`. Those are why a handful of items stay **suspected**.
