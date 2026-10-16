# Review round 1 — grok-4.7

Read-only review of Remain Faithful before friends-and-family use. No product code was changed. Findings below are from the source in this checkout, plus the test commands at the end. Items that depend on production env vars, the live App Store listing, or Fly secrets are marked **suspected**.

Severity: **P0** breaks the privacy promise or lets someone take an account or silently miss the accountability the product exists to provide. **P1** is a serious bypass, a false claim a partner will rely on, or a broken recovery path. **P2** is a real bug with a narrower blast radius. **P3** is hygiene.

Fixed product choices were not treated as defects: Family Controls filtering plus partner notify, DeviceActivity monitoring, time-window shielding as an architecture, opt-in Deep Scan via ReplayKit (it cannot see DRM streaming), and App Lockout via the Darwin notification `com.remainfaithful.broadcastStateChanged`. Where the copy describes behavior the code does not implement, the finding is the mismatch.

## Top 10

1. **P0** Sign in with Apple does not verify the identity token. A forged JWT with the victim’s email takes the account.
2. **P0** Inviting someone who already has an account as a partner inserts a `pending` row, tells the app `"connected"`, and emails a token that is never stored. Alerts and silence detection require `accepted`. The partner list does not show status.
3. **P1** The app drops `violence`, `self_harm`, `app_usage`, and `clean` from every in-app list. One-to-one partners have no alert inbox at all. Opening Home marks alerts seen.
4. **P1** Password reset emails default to `https://app.remainfaithful.app/reset-password`. Nothing on the site or in the app redeems the token.
5. **P1** Turning off Monitor App Usage, clearing the app selection, and signing out do not notify partners. Heartbeats exist only inside Deep Scan. Silence alerts require a prior heartbeat and an accepted relationship. Fly scale-to-zero stops the sweep while the machine is asleep.
6. **P1** The website and the App Store paste say partners are notified when a blocked category is attempted, and that time-window shielding exists. There is no ShieldAction extension. The only schedule is 00:00–23:59. Alert detail invents a browsing story the event did not record.
7. **P1** Adding an existing user to a group inserts them immediately. Their later events fan out to every co-member, and `GetGroup` exposes name, email, 30-day flag count, and streak.
8. **P1** Invite and password-reset emails interpolate names and URLs into HTML without escaping, and the invite routes are not rate-limited.
9. **P1** Token refresh writes the new JWT only to the Keychain. The in-memory session and the broadcast extension’s app-group copy stay stale, so Deep Scan uploads and heartbeats fail after about 24 hours. The privacy policy says the JWT is not in UserDefaults.
10. **P1** `next@14.2.5` is inside a large set of unpatched Next.js advisories (`npm audit`: 1 critical, 9 high, 2 moderate). The API image is `golang:1.22-alpine` with `golang-jwt/jwt/v5 v5.2.1` (GO-2025-3553).

Just under the line: Tier 1 keyword and URL hits are all labeled explicit sexual content, including gambling, gore, and suicide phrases (`RemainFaithfulBroadcast/SampleHandler.swift:132-142`, `:457-466`). A self-harm phrase can be stored as adult content, and a correctly categorized `self_harm` event from Tier 2 is then dropped by the in-app mapper in finding 3.

## 1. P0 — Apple identity tokens are decoded, not verified

`decodeAppleIdentityToken` base64-decodes the JWT payload and returns `sub` and `email`. It does not check the signature, `iss`, `aud`, or `exp`. The comment on the function says production must verify Apple’s JWKs.

```149:180:backend/internal/handler/social_auth.go
// decodeAppleIdentityToken decodes the JWT payload from Apple's identity token
// and returns the Apple user ID (sub) and email.
// NOTE: For production, also verify the JWT signature using Apple's JWK endpoint at
// https://appleid.apple.com/auth/keys and validate iss, aud, and exp claims.
func decodeAppleIdentityToken(idToken string) (sub, email string, err error) {
```

`AppleSignIn` (`social_auth.go:36-60`) then calls `findOrCreateSocialUser`. If the email matches an existing account, the handler attaches the attacker’s `apple_id`, overwriting any previous value, and issues a 24-hour JWT:

```118:126:backend/internal/handler/social_auth.go
	if email != "" {
		err = h.DB.QueryRowContext(ctx,
			`SELECT id, name, email FROM users WHERE email = $1`, email,
		).Scan(&userID, &userName, &userEmail)
		if err == nil {
			_, _ = h.DB.ExecContext(ctx,
				fmt.Sprintf(`UPDATE users SET %s = $1 WHERE id = $2`, idCol), idVal, userID)
			return
		}
	}
```

`idCol` is hardcoded to `apple_id` or `google_id` by the callers, so this is not SQL injection. The takeover does not need a real Apple signature.

**Fix:** Verify the token against `https://appleid.apple.com/auth/keys`. Require `iss` of `https://appleid.apple.com`, `aud` equal to the app’s bundle id, and a valid `exp`. Attach `apple_id` to an existing email account only when Apple’s verified email matches and the row’s `apple_id` is empty or already equal to that `sub`. Do not overwrite a different `apple_id`.

## 2. P0 — An existing user invited as a partner stays pending forever

`InvitePartner` documents that an existing account is connected immediately. The insert does not set `status`. The table default is `pending` (`backend/cmd/server/main.go:396`). The handler still returns `"status": "connected"`.

```51:80:backend/internal/handler/invites.go
	// If the invitee already has an account, create the relationship directly.
	// ...
		if err := h.DB.QueryRowContext(r.Context(),
			`INSERT INTO relationships (user_id, partner_id, type)
			 VALUES ($1, $2, 'partner')
			 ON CONFLICT (user_id, partner_id) DO UPDATE SET type = EXCLUDED.type
			 RETURNING id, status, created_at`,
			userID, inviteeID,
		).Scan(&relID, &status, &createdAt); err != nil {
			// ...
		}
		token, _ := inviteToken()
		acceptURL := siteBase() + "/accept-invite?token=" + token + "&type=partner"
		// token is never written to relationship_invites
		writeJSON(w, http.StatusCreated, map[string]any{
			"status":          "connected",
			"relationship_id": relID,
		})
```

`AcceptPartnerInvite` (`invites.go:126-131`) looks up `relationship_invites` by that token. The existing-user path never inserts a row there, so the link cannot be accepted. `website/src/app/accept-invite/page.tsx:12-14` ignores the token and only renders `InviteLanding`.

New accounts do get `status = 'accepted'` in `acceptPendingInvites` (`backend/internal/handler/auth.go:124-128`). Event fan-out (`events.go:136-139`) and the silence sweep (`heartbeat_sweep.go:63-65`) both require `status = 'accepted'`. A pending partner receives neither flags nor silence alerts.

The iOS list loads every relationship and does not render `status`. Pending partners appear under “CURRENT PARTNERS” (`RemainFaithful/SettingsView.swift:898-923`, loaded at `:1046-1062`). The invite call is `invitePartner` from onboarding and Settings. `POST /relationships` (`relationships.go:56-74`) has the same pending insert and, unlike the invite handler, returns the real `pending` status. It also answers `"no account found for that email"` (`relationships.go:45-47`), which confirms whether an address is registered. The app’s `createRelationship` client method is unused by the UI; the live path is the invite handler.

**Fix:** Keep the row `pending` until the invitee accepts, return `"pending"` (not `"connected"`), store the token in `relationship_invites`, and add an in-app accept action that calls `POST /relationships/accept-invite`. Show pending versus accepted in Manage Partners. Do not send events or silence alerts until `accepted`. Apply the same accept step to `CreateRelationship`, and stop returning “no account found.”

## 3. P1 — In-app history drops categories, and one-to-one partners have no inbox

`EventCategory.init?(apiString:)` maps only `adult_content`, `gambling`, `social_media`, and `gaming`. Anything else returns nil, and `ActivityEvent.from(remote:)` drops the event.

```22:29:RemainFaithful/DashboardView.swift
    init?(apiString: String) {
        switch apiString.lowercased().replacingOccurrences(of: " ", with: "_") {
        case "adult_content": self = .adultContent
        case "gambling":      self = .gambling
        case "social_media":  self = .socialMedia
        case "gaming":        self = .gaming
        default:              return nil
        }
    }
```

The server accepts `violence`, `self_harm`, `clean`, and `app_usage` (`events.go:43-47`). Callers that drop them: `DashboardView` history via `from(remote:)`, `SettingsView` activity log, and `GroupView.swift:816-820` (`listAlerts`, the only partner history, and only inside a group-member sheet). `NotificationService.routeNotification` (`NotificationService.swift:124-128`) sends a nil payload to the dashboard instead of the detail screen.

A severe `self_harm` push body is still honest (`events.go:263-264`: “their device flagged self-harm content”). After the banner is gone, the partner cannot open it. A severe `violence` push says “Open the app to see what was flagged” (`events.go:270`) and the app cannot show it. `social_media` and `gaming` are in the iOS enum and are not categories the server accepts, so those rows never arrive from the API.

`listAlerts` is not used for a one-to-one partner. `ContentView.swift:82-95` and `:166-168` call `markAlertsSeen` when the Home tab is selected or the app becomes active on Home, which clears the badge without opening a partner alert.

**Fix:** Map `violence`, `self_harm`, and `app_usage` to their own cases. Add a partner inbox backed by `GET /alerts` for one-to-one partners, not only the group member sheet. Mark an alert seen when that partner opens it, not when Home appears.

## 4. P1 — Password reset links point at a dead host, and nothing redeems them

```395:399:backend/internal/handler/auth.go
func getAppBaseURL() string {
	if u := os.Getenv("APP_BASE_URL"); u != "" {
		return u
	}
	return "https://app.remainfaithful.app"
}
```

`ForgotPassword` (`auth.go:327-328`) emails `{base}/reset-password?token=`. `APP_BASE_URL` is not in `backend/fly.toml`. Whether production overrides it is **suspected**. There is no `website/src/app/reset-password` page. The iOS client has `forgotPassword` and `resetPassword` API methods and a “Forgot Password?” button, and no screen that accepts the token or the URL. The app guide still says the link works for one hour (`website/src/app/app-guide/guide-content.ts:109-110`).

The token is stored in plaintext (`auth.go:316-321`), one row per user, valid for one hour. `ResetPassword` (`auth.go:375-386`) changes the password and deletes the reset row. It does not revoke JWTs. `Sign` issues a 24-hour token (`backend/internal/auth/jwt.go:20-32`) with no denylist.

Login and forgot-password responses do not reveal whether the email exists. That part is fine.

**Fix:** Send the link to a real page on remainfaithful.com that `POST`s the token to `/auth/reset-password`, or handle the URL in the app. Hash reset tokens at rest. On a successful reset, invalidate outstanding JWTs for that user.

## 5. P1 — Monitoring can be turned off quietly, and silence detection barely arms

`setMonitoringEnabled(false)` stops DeviceActivity and does not send a protection alert (`RemainFaithful/ActivitySelectionManager.swift:96-98`, `stopMonitoring` at `:146-150`). The server allowlist includes `monitoring_disabled` and the comment says SettingsView sends it (`backend/internal/handler/protection.go:58-59`). No Swift file sends that string. The Screen Time screen wires the toggle straight through (`ScreenTimeMonitoringView.swift:169-171`).

`selectionDidChange` (`ActivitySelectionManager.swift:49-56`) restarts monitoring. If the selection is emptied, `startMonitoring` sets monitoring off (`:100-105`) with no alert. Turning shielding off does alert (`shielding_disabled`, `:60-74`). Opening App Restrictions is PIN-gated (`SettingsView.swift:207`), and the toggles inside that screen are not gated again.

Sign out calls `AuthState.clearSession()` with no partner notice and no PIN (`SettingsView.swift:403-407`). `clearSession` (`AuthState.swift:83-92`) removes `authToken` from the Keychain and the app group. Delete account does notify partners. DeviceActivity events are queued before POST. `drainPendingRawEvents` (`ActivitySelectionManager.swift:158-195`) uploads that queue with whatever `authToken` is in the app group. After a later sign-in, those events post as the new account.

Family Controls revocation alerts only on an in-session change from `.approved` to something else (`RemainFaithfulApp.swift:92-104`). Revoking Screen Time while the app is not running does not fire `onChange` on the next launch.

Deep Scan stop notifies partners only when App Lockout was on (`SampleHandler.broadcastFinished`, `SampleHandler.swift:291-304`). Lockout defaults off (`AppLockoutManager.swift:31`, `bool` defaults to false). The privacy policy states this (`website/src/app/privacy/page.tsx:135`). Stopping Deep Scan with lockout off is still a quiet bypass of the opt-in scanner.

Heartbeats are not sent by `APIClient.heartbeat` (defined at `APIClient.swift:292` and never called). The only sender is the broadcast extension, which sleeps two minutes first and then posts every two minutes while a broadcast is running (`SampleHandler.swift:686-709`, body is `{screen: active|idle}`). The sweep ignores users with `last_heartbeat_at IS NULL` and relationships that are not `accepted` (`heartbeat_sweep.go:59-68`). Someone who never runs Deep Scan for two minutes is never covered by silence detection. Cycling Deep Scan off for less than the silence window (default 30 minutes) resets the clock because a new heartbeat moves `last_heartbeat_at` forward.

`backend/fly.toml:14-16` sets `auto_stop_machines = "stop"` and `min_machines_running = 0`. The sweep is an in-process ticker started in `main` (`backend/cmd/server/main.go:71-73`). It does not run while the machine is stopped. The first tick is five minutes after start (`heartbeat_sweep.go:28-38`). Fly HTTP health checks do not keep a stopped machine awake.

The Screen Time shield can be overridden with the device Screen Time passcode. The app says so (`ScreenTimeMonitoringView.swift:219`). There is no ShieldAction target to report the attempt.

**Fix:** Send `monitoring_disabled` through the existing protection pipeline when monitoring is turned off or the selection becomes empty. Send a protection alert on sign-out, and do not upload another account’s `pendingRawEvents` (stamp the queue with the user id, or drop it on sign-out). On launch, if Family Controls is no longer approved and it previously was, send `family_controls_revoked`. Send heartbeats from the main app or the DeviceActivity extension while monitoring is on, not only from Deep Scan, so silence can arm. Set `min_machines_running = 1` so the sweep process stays up. That is a few dollars a month (see cost). A ShieldAction (or equivalent) should emit the existing alert when a shield is hit, which is the notify step the architecture already describes.

## 6. P1 — Public copy overstates blocking, time windows, and what an alert means

There is no ShieldAction extension in the repo. `RemainFaithfulShieldConfig` only customizes the shield UI. Shielding and monitoring both default off. The DeviceActivity schedule is a single repeating window from 00:00 to 23:59 (`ActivitySelectionManager.swift:108-112`). The threshold is one minute of cumulative use per day (`:119-130`). There is no UI for hours the user sets.

Pages that say partners are notified when a blocked category is attempted include:

- `website/src/app/page.tsx:121` and `:197`
- `website/src/app/about/page.tsx:198`
- `website/src/app/privacy-architecture/page.tsx:69`
- `website/src/app/how-it-works/page.tsx:30` and `:175`
- `website/src/app/group-setup-guide/page.tsx:52` and `:158`
- `website/src/lib/structured-data.ts:144`
- `website/src/app/privacy/page.tsx:49` (“When always-on filtering flags a blocked category”)

Time-window copy: `how-it-works/page.tsx:222` (“You can restrict chosen apps to time windows”), `how-it-works/page.tsx:11` and `:30`, and `structured-data.ts:144`. The App Store paste in `docs/APP_STORE_SUBMISSION.md:89-105` says “If a blocked app is opened…”, “Optional time-window shielding”, and “Alerts if Screen Time permission is revoked or the device goes silent”. Whether that exact text is what App Store Connect shows today is **suspected**.

`ActivityEvent.extendedDescription` (`DashboardView.swift:141-154`), shown by `AlertDetailView.swift:123`, rewrites the stored summary. Adult content that does not contain the word “search” becomes “detected during a web browsing session.” Gambling always becomes “a gambling or sports-betting website was visited.” Real summaries are generic (“Explicit content detected”, “Blocked domain detected”, “Monitored app used for 1+ minute”). The sample events at `DashboardView.swift:219-225` are demo-mode only.

DRM limits are stated honestly on the privacy-architecture and about pages. That part should stay.

**Fix:** Implement the user-chosen time windows the architecture already describes, and emit a partner alert when a shield is presented or bypassed. Until that ships, change the website and the App Store text so they match the code: a shield blocks chosen apps and categories, usage is a once-per-day one-minute DeviceActivity event, and Deep Scan is opt-in. `extendedDescription` should show the server summary, not a browsing story.

## 7. P1 — An existing user is added to a group with no accept step

```281:305:backend/internal/handler/invites.go
	// If the invitee already has an account, add them directly.
	// ...
			`INSERT INTO group_members (group_id, user_id, role)
			 VALUES ($1, $2, 'member')
			 ON CONFLICT (group_id, user_id) DO NOTHING
			 RETURNING joined_at`,
```

Any member may invite, not only an admin (`checkCanInviteToGroup`, `invites.go:242-261`; `groups.go:364-366`). `CreateEvent` fans out to every co-member (`events.go:136-145`). `GetGroup` returns each member’s name, email, 30-day flag count, and streak to any member (`groups.go:122-158`). A person who knows an email can put that account in a group and then receive their future accountability events.

The cap check runs before the insert, outside a transaction (`invites.go:254-260` then `:287-292`). Concurrent invites can exceed 12. New accounts join through `acceptPendingGroupInvites`, which does lock the group (`auth.go:141-176`).

The emailed URL for an existing user has no token (`invites.go:301`). `website/src/app/join/[id]/page.tsx` ignores the id and tells the visitor to ask for an email invite. That page is not an authorization bug.

**Fix:** Insert a pending group invite for existing users and add them only after they accept, using the same locked cap check as signup. Until they accept, do not fan events to them and do not show them the roster.

## 8. P1 — Invite mail is HTML-injectable and unbounded

`SendPartnerInvite`, `SendGroupInvite`, and `SendPasswordReset` put `inviterName`, `groupName`, `toName`, and URLs into HTML with `fmt.Sprintf` and no `html.EscapeString` (`backend/internal/email/email.go:42-52`, `:65-74`, `:118-123`). The contact path does escape (`email.go:95-100`). Registration only `TrimSpace`s the name (`auth.go:35`), so a name can contain tags or newlines. The subject line is also the raw name (`email.go:33`, `:58`).

`POST /relationships/invite`, `/groups/{id}/invite`, and `/groups/{id}/email-invite` are authenticated and are not wrapped in `authenticatedRateLimiter` (`backend/cmd/server/main.go:269`, `:288-289`). `/auth/*` and `/contact` are limited to 10 requests per minute per IP (`main.go:113-117`, `:307-311`). One logged-in user can send as many invite emails as SendGrid will accept.

The default From address is `noreply@remainfaithful.app` if `SENDGRID_FROM_EMAIL` is unset (`email.go:27`). `backend/internal/payment/client.go` notes that `remainfaithful.app` does not resolve. Whether the live From domain is verified at SendGrid is **suspected**. If it is not, invite and reset mail will not be delivered.

**Fix:** HTML-escape every interpolated value, strip newlines from names before they touch a header, and put the invite routes on the existing per-user limiter (a small daily cap is enough for a group of 12). Set `SENDGRID_FROM_EMAIL` to a domain SendGrid has verified.

## 9. P1 — Refresh does not update the token the app and the extension use

```186:192:RemainFaithful/APIClient.swift
    func refreshTokenIfNeeded() async throws {
        guard let expiry = AuthState.shared.tokenExpiresAt,
              expiry.timeIntervalSinceNow < 3600 else { return }
        let resp: [String: String] = try await post("/auth/refresh", body: [String: String]())
        if let newToken = resp["token"] {
            KeychainHelper.shared.set(newToken, for: "authToken")
        }
    }
```

`AuthState.token` prefers `inMemoryToken` (`AuthState.swift:41`). `setSession` is the only writer of the app-group `authToken` (`AuthState.swift:48-49`). The broadcast extension reads that app-group value (`SampleHandler.swift:541-544` and `:699-700`). After refresh, API calls from a running app keep using the expiring in-memory token. After a restart, the app reads the keychain token and the extension still has the old app-group token. Deep Scan uploads and heartbeats then 401. A user who had armed silence detection can false-trigger it.

The same `setSession` writes the JWT into shared `UserDefaults`. The privacy policy says JWTs live in the Keychain and not in UserDefaults (`website/src/app/privacy/page.tsx:152-154`). The extension needs a token it can read; the policy sentence is still false.

**Fix:** On refresh, update the in-memory token, the Keychain, and the app-group token together (the same writes as `setSession`). Change the privacy sentence so it matches that.

## 10. P1 — Next.js and the API base image are on unpatched releases

`website/package.json` pins `next@14.2.5`. `npm audit` on this lockfile reported 12 issues: 1 critical, 9 high, 2 moderate. The critical package is `next` itself (`0.9.9 - 16.3.0-preview.10` in the audit range), covering a long list of advisories. There is no `middleware.ts`, so the middleware authorization-bypass advisory is unlikely to apply to this app. The Windows image-optimization RCE does not apply on Vercel. `next/image` is used for local screenshots (`website/src/components/PhoneFrame.tsx:1`) and `images.remotePatterns` is empty (`website/next.config.js:52`), which lowers the image-optimizer exposure. The App Router Server Components and Server Actions denial-of-service advisories still apply to this site. This review did not confirm an unauthenticated remote code execution.

`backend/Dockerfile:2` builds with `golang:1.22-alpine` and runs on `alpine:3.19` (`Dockerfile:18`). `backend/go.mod` requires `github.com/golang-jwt/jwt/v5 v5.2.1`. `govulncheck` reports GO-2025-3553 (unbounded memory while parsing a JWT header), fixed in v5.2.2, reachable from `internal/auth/jwt.go:38`. The same run reports many standard-library issues in go1.22.2, including GO-2026-6617 (HTTP/2 HPACK encoder crash, fixed in go1.26.9). The traces are noisy because any `net/http` use is listed. The actionable fact is the toolchain, not each id.

**Fix:** Upgrade Next.js to a current patched release and re-run `npm audit`. Rebuild the API with a current Go release and bump jwt to at least v5.2.2. Rebuild the Alpine runtime base with it.

## 11. P1 — Tier 1 labels every hit as explicit sexual content

`Tier1Rules.keywordPatterns` includes gambling, gore, and suicide phrases (`SampleHandler.swift:132-142`). `matchesKeyword` posts one event with category `.explicitSexual` and summary “Explicit keyword detected in screen text” (`:463-466`). URL matches do the same for every blocklist domain (`:457-460`), and the list includes betting and gore hosts (`:127-129`). Tier 2 can still emit a second, correctly categorized event. The in-app mapper in finding 3 then drops `self_harm` and `violence`, so the event the partner can actually open is the mislabeled one.

**Fix:** Choose the category from the pattern or domain that matched. Keep self-harm distinct from explicit sexual content.

## 12. P1 — Email can be changed with no verification

`UpdateMe` writes the new email immediately (`backend/internal/handler/users.go:71-75`). A unique violation becomes a generic 500 (`:76-78`). The signed-in user can move the account onto an address that has a pending invite, and `redeemPendingInvites` is not what runs here — the invite accept path compares the account email (`invites.go:143-146`), so changing to the invited address and then accepting is enough. The old inbox loses the account. Checkout uses the JWT email, which goes stale until the next token (`donations.go` reads the authenticated user; the stored email is whatever `UpdateMe` last wrote).

**Fix:** Keep the old email until the new one is confirmed with a single-use link. Return a conflict when the address is taken, without revealing anything else.

## 13. P1 — Google ID tokens are checked for `sub` only

`verifyGoogleIDToken` (`social_auth.go:183-208`) calls `https://oauth2.googleapis.com/tokeninfo?id_token=` and accepts any tokeninfo result with a `sub`. It does not check `aud`, `azp`, or `email_verified`. The token is placed in the query string. Google still signs the token, so this is not the unauthenticated forgery in finding 1. It is token confusion: a Google ID token minted for a different client can be replayed here, and an unverified email can be attached to an account through `findOrCreateSocialUser`.

**Fix:** Require `aud` (or `azp`) to be this app’s Google client id and require `email_verified`. Verify the token locally, or POST it to tokeninfo, so it does not land in a query string and in access logs.

## 14. P2 — Server logs and the unified log hold more than the privacy sentence allows

These do not leave the device as an upload. They are not a P0 content leak.

- OCR text, first 80 characters, is passed to `Logger` (`SampleHandler.swift:452`). A matched domain is logged at `:458`. A keyword match logs the first 80 characters of OCR text at `:464`. Swift `Logger` string interpolations are redacted unless a debugger is attached. Frames are not uploaded. `uploadEvent` sends category, severity, summary, and timestamp (`SampleHandler.swift:555-556`). The temp JPEG `rf_sca_*.jpg` (`:475-477`) is local and removed.
- Protection alerts log `detail` up to 80 runes (`protection.go:166-168`). Several iOS callers put the user’s name in that string. The APNs body does not use `detail` (`protection.go:170-171`), which is the right split. The log line is not.
- Invite handlers log the recipient email (`invites.go:75`, `:102`).
- SendGrid logging prints the recipient and the response body (`email.go:158-170`).

**Fix:** Stop logging OCR text, matched domains, protection `detail`, and recipient addresses. Log an event id or a user id.

## 15. P2 — Account enumeration, PIN attempts, and event timestamps

- `Register` returns 409 “email already registered” (`auth.go:60-62`). Login and forgot-password do not. The auth limiter is 10 requests per minute per IP.
- `clientIP` trusts `X-Forwarded-For` when `Fly-Client-IP` is absent (`main.go` around the rate-limit helpers). Spoofing that header by reaching the VM outside the Fly proxy is **suspected**.
- PIN verify allows 5 attempts per 15 minutes in memory, per process (`protection.go:20-26`). A 4-digit PIN can be walked slowly, and the bucket resets on process restart and is not shared across machines.
- `CreateEvent` stores a client-supplied timestamp with no window (`events.go:107-113`). That can move streaks. The write limit is 30 requests per minute per user (`main.go:290`), which is about 43,000 events per user per day. There is no retention job. `events` and `alerts` grow until the account is deleted (`ON DELETE CASCADE`).

**Fix:** Use the same generic response for duplicate registration as for login. Key the auth limiter on Fly’s client IP only. Count PIN failures in Postgres. Reject event timestamps outside a short skew of server time, and delete or archive events past a stated retention window.

## 16. P2 — Privacy policy does not match deletion, analytics, or crypto claims

- “Account data … deleted within 30 days of account deletion” and “Donation records: retained for 7 years” (`privacy/page.tsx:105-112`). `DeleteMe` deletes the user immediately. `donations.user_id` is `ON DELETE CASCADE` (`main.go:461`). Donation rows go with the user. The rights section also says deletion is immediate and confirmation follows within 30 days (`privacy/page.tsx:134`). There is no deletion-confirmation email.
- GA4 loads whenever `NEXT_PUBLIC_GA_ID` is set (`website/src/app/layout.tsx:107-120`). Nothing reads Do Not Track. The privacy page says GA is optional and that Do Not Track opts out (`privacy/page.tsx:61-64`). Whether the env var is set in production is **suspected**. The iOS app does not ship an analytics SDK (`NSPrivacyTracking` is false). CSP allows `'unsafe-inline'` scripts (`website/next.config.js:5`).
- “All data in transit is encrypted using TLS 1.3” and “We do not log sensitive data” (`privacy/page.tsx:151-154`). `force_https` is on in `fly.toml`. The app does not pin TLS 1.3. Finding 14 is the logging mismatch.
- The architecture page says data at rest is AES-256 and that transit is end-to-end (`privacy-architecture/page.tsx:266` and `:289`). The server sees plaintext JSON. Disk encryption, if Fly provides it, is not configured in this repo. Treat AES-256 and “end-to-end” as **suspected** infrastructure claims, not app behavior.

**Fix:** Say that account deletion is immediate and that donation rows are deleted with the account, or stop cascading them and keep the financial record the policy promises. Load GA only after a real opt-in, and remove the Do Not Track sentence until the code honors it. Describe TLS as HTTPS to the API, not end-to-end encryption, and drop AES-256 unless the database volume is actually configured that way.

## 17. P2 — “App removed” is pushed at the dead device token

When APNs says a device token is invalid, `notifyPartners` marks it inactive and then sends the “app was removed” payload to that same token (`backend/internal/handler/notifications.go:343-367`). Partners do not get `APP_DELETED`.

**Fix:** Send that payload to the user’s partners’ tokens, not back to the token Apple just rejected.

## 18. P2 — Cost is not $4 per person, and a few unbounded writes are the real bill

Public list prices, used as ranges rather than a quote (checked for this review, October 2026):

- One Fly `shared-cpu-1x` machine with 256 MB is a few dollars a month. `ord` costs more than the cheapest US regions. `backend/fly.toml:25-28` is that size.
- Fly’s unmanaged Postgres docs put a single small node near a couple of dollars a month and a three-node production cluster on the order of $80–$160 a month, plus volume. `DATABASE_URL` is a secret. Which plan is running is **unknown**. Do not treat the about-page “~$80/mo” box as proof of the SKU.
- APNs is $0.
- SendGrid’s trial is on the order of 100 emails a day. Essentials is about $20 a month for 50,000 emails.
- Vercel Hobby is $0 inside the included quota. Pro is about $20 per seat. The marketing site does not scale with phones. Homepage donations are Stripe Payment Links, not the legacy `/api/donate` route.

Heartbeats are not the cost driver. They run only during Deep Scan, after a two-minute delay, then every two minutes (`SampleHandler.swift:686-709`). If 20% of users scan 30 minutes a day, that is about 15 heartbeats each, or about 3 per average user per day.

| Users | Heartbeats / day (that assumption) | Honest events | What actually costs money |
| --- | --- | --- | --- |
| 100 | ~300 | a few hundred a day if each person flags a handful of times | One small machine if it stays up, a small Postgres, maybe a SendGrid plan. On the order of $10–40/month. |
| 1,000 | ~3,000 (~0.03 req/s) | a few thousand a day | Still tens of dollars if Postgres stays a single small node and invites are capped. |
| 10,000 | ~30,000 (~0.35 req/s) | ~20,000/day at 2 flags/user, ~600k event rows/month plus alert rows, on the order of a couple of hundred MB | Still well under a few hundred dollars a month unless Postgres is a HA cluster, invite email is abused, or the 30-events/minute cap runs for a long time with no retention. |

The about page says server costs are “approximately $4 per month per participant” and, in the next grid, “~$80/mo” (`website/src/app/about/page.tsx:213-218`). Those two figures match each other only around 20 people. $80 a month is a plausible padded small-stack bill. It is not a per-user law. The marginal cost of this architecture is not $4 per user.

**Fix:** Rate-limit the invite routes (finding 8). Cap or expire `events` and `alerts` (finding 15). Set `min_machines_running = 1` so silence detection works; that is cheaper than the false alerts from a sleeping sweeper. Add an index on `users.last_heartbeat_at` when the sweep is expected to matter. Change the about-page sentence so a donor is not told that each person costs $4.

`GetUserStats` loads distinct flagged dates for one user (`users.go` around the stats query). That is fine at 10,000 users. `GetGroup`’s per-member count uses the events-by-user index. There is no index on `last_heartbeat_at`. The waitlist route logs the email and returns success without storing it (`website/src/app/api/waitlist/route.ts`, already noted as MF-05). The contact route logs name and email when `BACKEND_URL` is unset.

## 19. P3 — Smaller notes

- `npm run lint` (`next lint`) exits asking how to configure ESLint. There is no ESLint config, so lint cannot run non-interactively. This review did not add one.
- `website/public/sitemap.xml` is stale (lastmod 2025-05-29) and omits pages that `website/src/app/sitemap.ts` includes. `next build` emits `/sitemap.xml` as an app route, so the generated file should win. **Suspected** only if a CDN served the file in `public/`. `public/robots.txt` is weaker than `robots.ts` (it does not disallow `/api/`).
- Canonicals use `https://remainfaithful.com`. Stripe return URLs use `https://www.remainfaithful.com`. Duplicate-host indexing is **suspected**, not verified.
- Contact form labels are not wired with `htmlFor` / `id`. The FAQ accordion sets `aria-expanded` while collapsed answers stay in the DOM. A skip link exists (`layout.tsx:106`).
- No in-app age attestation. Registration collects name, email, and password. The privacy policy says the product is not for under 18 (`privacy/page.tsx:145-148`). The store rating is an App Store Connect setting. The about page says the founder serves in “youth and high school ministry at my local church” (`about/page.tsx:55`). It does not name a church, city, state, or street. That sentence is not a claim that the app is for minors. Other pages say 18+ and not parental controls.
- `POST /debug/test-push` is registered only when debug push is enabled, and `APP_ENV=production` turns it off (`main.go` and `fly.toml`).
- CORS echoes `ALLOWED_ORIGINS` (default `https://remainfaithful.com`, not `www`). The iOS app does not need CORS.
- `001_schema.sql` lags `migrate()` in `main.go`. The live schema is the `migrate()` function.
- `JWT_SECRET` empty panics on the first `secret()` call (`jwt.go:60-65`), not at process start. No production secrets were found in the tree. `backend/README.md` shows a dev secret example only.

## Checked, and not a finding

- No founder city, state, named church, or street address on the website.
- No donation UI in the iOS target. Website Stripe donations stay. The API still has `POST /donations/create-checkout-session` and a billing portal (`main.go:304-305`). That is not an App Store 3.1.1 issue in the binary. The webhook checks the Stripe signature, and the stored payload is a summary rather than the card object (`donations.go`, `summarizeEventPayload`).
- Deep Scan uploads category, severity, summary, and timestamp. Clean frames are not uploaded. The event summary allowlist matches the iOS strings, including the continued-activity pattern (`events.go:19-60`).
- `heartbeat_silence` cannot be forged through the protection HTTP handler (`protection.go:69-73`).
- Group roster PII requires membership (`groups.go:122`). Partner PIN set and delete check `partner_id`.
- The Family Controls entitlement is present on the app, shield, and Device Activity targets. The app was just approved, so the distribution-versus-development entitlement is not reopened here.
- Password hashing is bcrypt. App JWTs are HS256 with issuer and expiry checks (`jwt.go:35-48`). That does not fix finding 1, which never calls `Parse` on the Apple token.

## What was run

From `backend/`, with Go 1.22.2:

- `go test ./...` passed (`cmd/server`, `internal/email`, `internal/handler`, `internal/payment`). `internal/apns` and `internal/auth` have no tests.
- `go vet ./...` exited 0.
- `govulncheck ./...` exited 3. See finding 10. The standard-library traces are not each a separate product bug.

From `website/`:

- `npm test` passed (6 tests: checkout URLs and donate links).
- `npm run lint` failed because Next.js tried to start the interactive ESLint setup. No config was written.
- `npm run build` succeeded (Next 14.2.5, 35 static pages). One warning: using the edge runtime on a page disables static generation for that page. This review did not pin which page.
- `npm audit`: 1 critical, 9 high, 2 moderate. See finding 10. Transitive highs in `braces`, `browserslist`, `postcss`, `nanoid`, and Tailwind’s build tools are P3 unless they ship to browsers.

Not run: a device build, a live call against production, or a check of App Store Connect’s current description. Those gaps are why a few items above stay **suspected**.
