# Remain Faithful — Adversarial Review, Round 1 (Claude)

**Reviewer:** Claude (Opus 5)
**Date:** 2026-10-13
**Commit reviewed:** `784998f`
**Scope:** iOS app (`RemainFaithful*`), Go backend (`backend/`), Next.js website (`website/`)
**Method:** Read-only source review plus test execution. No product code was modified.

---

## Verification performed

| Check | Command | Result |
|---|---|---|
| Go vet | `go vet ./...` | **Pass**, exit 0, no diagnostics |
| Go tests | `go test ./...` | **Pass**, all 6 packages (2 have no test files) |
| Website build | `npm run build` | **Pass**, 35 routes generated |
| Website lint | `npm run lint` | **Not configured** — drops into the interactive ESLint setup prompt (see F-21) |
| Website tests | `npm test` | **Pass**, 6/6 |
| Secret scan | regex sweep for Stripe/SendGrid/AWS/GitHub/PEM keys | **Clean** — only `website/.env.example` placeholders |

Two findings were confirmed by executing isolated reproductions rather than by reading alone; those are marked **[reproduced]**.

---

## Severity key

| | Meaning |
|---|---|
| **P0** | Breaks the product's core safety promise, or allows account takeover / unconsented disclosure of the most sensitive data. Fix before further promotion. |
| **P1** | Serious security, privacy, or silent-monitoring-failure defect. |
| **P2** | Real defect with material user, cost, or trust impact. |
| **P3** | Hygiene, polish, minor cost. |

---

# P0 findings

## F-01 — Unauthenticated account takeover of any user via forged Apple identity token [reproduced]

**Severity: P0** · `backend/internal/handler/social_auth.go:153-181`, `:119-127`

`decodeAppleIdentityToken` base64-decodes the JWT payload and reads `sub` and `email`. It **never verifies the signature**, the issuer, the audience, or expiry. The code comments on line 151-152 acknowledge this is missing.

```go
// social_auth.go:151-153
// NOTE: For production, also verify the JWT signature using Apple's JWK endpoint at
// https://appleid.apple.com/auth/keys and validate iss, aud, and exp claims.
func decodeAppleIdentityToken(idToken string) (sub, email string, err error) {
```

`findOrCreateSocialUser` then falls back to **matching on email** when the `apple_id` is unknown:

```go
// social_auth.go:119-127
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

**Failure scenario.** An attacker who knows only a victim's email address builds the token
`base64url({"alg":"RS256"}) + "." + base64url({"sub":"anything","email":"victim@example.com"}) + "." + base64url("garbage")`
and POSTs it to `/auth/apple`. The server decodes it, finds no user with that `apple_id`, matches the victim by email, **permanently attaches the attacker's chosen `apple_id` to the victim's account**, and returns a valid 24-hour session JWT for the victim. The attacker now reads the victim's alert history, partner list, group membership, and every future flagged event — and can re-authenticate at will even after the victim changes their password, because the attacker-chosen `apple_id` persists.

The `/auth` rate limit of 10/min/IP (`main.go:117`) is no mitigation; the attack needs one request.

**Reproduction.** A verbatim copy of `decodeAppleIdentityToken` was compiled and run against a forged token with a garbage signature:

```
decodeAppleIdentityToken -> sub="attacker-controlled-id-12345" email="victim@example.com" err=<nil>
RESULT: ACCEPTED. No signature verification performed.
```

**Fix.**
1. Fetch Apple's JWKS from `https://appleid.apple.com/auth/keys` (cache by `kid`), verify the RS256 signature with `jwt.ParseWithClaims`, and assert `iss == "https://appleid.apple.com"`, `aud == <your bundle ID>`, and `exp` in the future. Reject anything that fails.
2. Remove the email-fallback account linking in `findOrCreateSocialUser` (`social_auth.go:119-127`). Linking a social identity to an existing password account must require proving control of that account (sign in first, then link) — never an assertion carried in the unverified token itself.
3. Apply the same `aud` check to Google (see F-07).

---

## F-02 — Protection alerts are fire-and-forget: airplane mode silently defeats all monitoring

**Severity: P0** · `RemainFaithful/SettingsView.swift:310,853`, `AppLockoutManager.swift:95`, `ActivitySelectionManager.swift:69`, `RemainFaithfulApp.swift:99`, `PartnerPINManager.swift:65`, `PINEntryView.swift:393`

Every one of the seven tamper-notification call sites uses `try? await` and discards the error. There is **no retry queue for protection alerts** — a repo-wide search for `pendingProtection` / `retryProtection` / `protectionQueue` returns nothing.

```swift
// SettingsView.swift:308-314
lockoutManager.setEnabled(false)
Task {
    try? await APIClient.shared.sendProtectionAlert(
        type: "lockout_disabled",
        detail: "\(userName) disabled the App Lockout feature."
    )
}
```

The local setting is changed **first and unconditionally**; the notification is a detached, unawaited, error-swallowing side effect.

Note the asymmetry: *content events* do have durable retry (`SampleHandler.writeEventToRetryQueue`, `SampleHandler.swift:~585`). The **tamper signals — the security-critical ones — do not.**

**Failure scenario.** Enable Airplane Mode. Turn off monitoring, turn off shielding, disable App Lockout, or have a partner's PIN removed. Each `sendProtectionAlert` throws; `try?` discards it; nothing is queued. Turn Airplane Mode off. Monitoring is now disabled and **no record of the change exists anywhere** — not on the device, not on the server, not in any partner's app. This requires no jailbreak, no modified client, and no technical skill.

**Fix.** Invert the order and make the signal durable:
1. Persist the protection alert to App Group storage **before** applying the local setting change, using the same pattern as `writeEventToRetryQueue`.
2. Drain that queue on app foreground, on network reachability change, and from the `BGAppRefreshTask`.
3. Have the server acknowledge; only clear the queue entry on a 2xx.
4. Also persist server-side — see F-03.

---

## F-03 — Protection and heartbeat-silence alerts have no durable server-side record

**Severity: P0** · `backend/internal/handler/protection.go:142-193`, `heartbeat_sweep.go:108-131`, `alerts.go:85-86`

`sendProtectionAlertToPartners` and `sendSilenceAlertToPartner` send an APNs push and **write nothing to the database**. `ListAlerts` can only ever return content events, because it inner-joins `events`:

```sql
-- alerts.go:85-86
FROM   alerts a
JOIN   events e ON e.id = a.event_id
WHERE  a.recipient_user_id = $1
```

and the schema forbids an alert without an event:

```sql
-- main.go:426
event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
```

**Failure scenario.** A user disables monitoring. The push fires once. The partner's phone is off, in Focus mode, has notifications disabled for the app, or the push is simply dropped — APNs is explicitly best-effort delivery. The partner opens the app later and sees **nothing**: no badge, no history entry, no indication the event ever happened. Every tamper signal in the product — `monitoring_disabled`, `shielding_disabled`, `lockout_disabled`, `family_controls_revoked`, `pin_removed`, `pin_changed`, `wrong_pin_attempt`, `deep_scan_stopped`, `heartbeat_silence` — is a single ephemeral push with no durable backing.

Combined with F-02, the tamper-detection layer has no reliable delivery path on either end.

**Fix.** Make `alerts.event_id` nullable, add `alert_type TEXT` and `kind TEXT NOT NULL DEFAULT 'event'`, and insert a row for every protection and silence alert. Change `ListAlerts` to a `LEFT JOIN` and return protection alerts in the feed. Include them in `AlertUnreadCount` so the partner sees a badge even when the push is lost.

---

## F-04 — Inviting an existing user as a partner leaves the relationship permanently `pending`, and pending partners receive no alerts

**Severity: P0** · `backend/internal/handler/invites.go:62-81`, `events.go:139`, `main.go:396`

This is the root cause of the known "invite stays pending forever" bug, and the consequence is worse than a stuck UI state.

When the invitee already has an account, `InvitePartner` inserts the relationship **without specifying `status`**, so it takes the schema default `'pending'` (`main.go:396`):

```go
// invites.go:62-67
`INSERT INTO relationships (user_id, partner_id, type)
 VALUES ($1, $2, 'partner')
 ON CONFLICT (user_id, partner_id) DO UPDATE SET type = EXCLUDED.type
 RETURNING id, status, created_at`
```

It then generates an accept token and embeds it in the email — but **never stores it**:

```go
// invites.go:72-73
token, _ := inviteToken()
acceptURL := siteBase() + "/accept-invite?token=" + token + "&type=partner"
```

There is no `INSERT INTO relationship_invites` on this path. `AcceptPartnerInvite` looks the token up in that table (`invites.go:126-131`), finds nothing, and returns `404 "invite not found or expired"` — **always**. No code path can ever move this relationship to `accepted`.

Meanwhile alert fan-out filters on exactly that status:

```sql
-- events.go:136-139
SELECT DISTINCT partner_id
FROM   relationships
WHERE  user_id = $1 AND status = 'accepted'
```

**Failure scenario.** A user invites their accountability partner, who already has an account. The API returns `{"status":"connected"}` (`invites.go:78`) and the partner appears in `ListRelationships`. **The partner never receives a single alert, ever.** Both people believe monitoring is active. This also silently disables the Partner PIN feature for that pair, since `SetRelationshipPIN` requires `status = 'accepted'` (`protection.go:335`), and excludes them from heartbeat-silence sweeps (`heartbeat_sweep.go:65`) and protection alerts (`protection.go:147`).

The misleading `"connected"` response is what makes this P0 rather than P1: the product actively tells the user they are protected when they are not.

**Fix.** On the existing-user branch, persist the invite so it is acceptable:
```go
// after creating the relationship row
_, err := h.DB.ExecContext(ctx,
    `INSERT INTO relationship_invites (inviter_id, invitee_email, token)
     VALUES ($1, $2, $3)
     ON CONFLICT (inviter_id, invitee_email) DO UPDATE
       SET token = EXCLUDED.token, created_at = NOW(), status = 'pending'`,
    userID, req.Email, token)
```
and change the response `"status"` from `"connected"` to `"invited"` so the UI stops claiming a connection that does not exist. Separately, add a backfill migration for relationships already stuck in `pending`, and surface pending state explicitly in the iOS partner list.

---

## F-05 — Any user can add any other user to a group by email, with no consent, and immediately receive their explicit-content alerts

**Severity: P0** · `backend/internal/handler/invites.go:283-305`, `events.go:141-145`, `groups.go:143`

`inviteEmailToGroup` adds an existing account to a group **immediately and unconditionally** — there is no acceptance step for group membership at all:

```go
// invites.go:283-293
if err := h.DB.QueryRowContext(ctx,
    `SELECT id FROM users WHERE email = $1`, email,
).Scan(&inviteeID); err == nil {
    err := h.DB.QueryRowContext(ctx,
        `INSERT INTO group_members (group_id, user_id, role)
         VALUES ($1, $2, 'member')
         ON CONFLICT (group_id, user_id) DO NOTHING
         RETURNING joined_at`,
        groupID, inviteeID,
    ).Scan(&joinedAt)
```

Group co-members are an alert audience with **no status check whatsoever** (unlike relationships):

```sql
-- events.go:141-145
UNION
SELECT DISTINCT gm2.user_id
FROM   group_members gm1
JOIN   group_members gm2
       ON gm2.group_id = gm1.group_id AND gm2.user_id != $1
WHERE  gm1.user_id = $1
```

**Failure scenario.** An attacker registers an account, creates a group, and calls `POST /groups/{id}/invite` with `victim@example.com`. The victim is now a member. From that moment the attacker receives every `"Explicit content detected"` / `"Self-harm content detected"` alert the victim's device generates, including category, severity, timestamp, and summary. The attacker also calls `GET /groups/{id}` and reads the victim's **name, email address, 30-day flag count, and streak** (`groups.go:143`). The victim's only signal is an invitation email saying they were invited — after they were already added — and they must notice it and manually leave.

This exposes the single most sensitive data the product holds to anyone who knows an email address. For this user population, that is a disclosure with real-world consequences.

**Fix.** Require explicit acceptance for group membership exactly as the partner flow intends to:
1. On the existing-user branch, insert a `group_invites` row with status `pending` and send the email with an accept token — do **not** insert into `group_members`.
2. Add an authenticated `POST /groups/accept-invite` that verifies the caller's email matches `invitee_email`, then inserts the membership under the existing `FOR UPDATE` cap check in `joinGroupFromInvite` (`auth.go:181-217`).
3. Until that ships, as a stopgap, exclude group co-members from the `CreateEvent` fan-out unless membership was self-initiated.

---

## F-06 — Password reset is completely broken: the link targets a dead domain, and no reset page exists on any domain

**Severity: P0** · `backend/internal/handler/auth.go:327-328,395-399`; `website/src/app/` (route absent)

The reset email is built from `getAppBaseURL()`:

```go
// auth.go:327-328
resetURL := fmt.Sprintf("%s/reset-password?token=%s", getAppBaseURL(), token)

// auth.go:395-399
func getAppBaseURL() string {
	if u := os.Getenv("APP_BASE_URL"); u != "" {
		return u
	}
	return "https://app.remainfaithful.app"
}
```

`APP_BASE_URL` is **not set in `backend/fly.toml`** (the `[env]` block contains only `PORT` and `APP_ENV`), so unless it exists as an unlisted Fly secret, the default applies — and `app.remainfaithful.app` is the dead domain.

It fails either way, because **there is no `/reset-password` route on the website**. The full route list from `next build` is: `/`, `/about`, `/accept-invite`, `/app-guide`, `/blog`, `/blog/[slug]`, `/christian-accountability-app`, `/church-accountability`, `/compare/*`, `/feedback`, `/free-accountability-app`, `/group-setup-guide`, `/how-it-works`, `/invite`, `/join/[id]`, `/partners`, `/privacy`, `/privacy-architecture`, `/terms`, `/what-is-accountability-software`, plus `/api/*`. A grep for `reset-password` across `website/src` returns nothing.

Meanwhile the app ships a fully wired Forgot Password UI (`LoginView.swift:63,78-79` → `APIClient.forgotPassword`, `APIClient.swift:363`).

**Failure scenario.** A user with a password account forgets their password and taps "Forgot Password". They receive an email. The link either fails DNS resolution (`app.remainfaithful.app`) or 404s (`www.remainfaithful.com/reset-password`). **No password account on the service can be recovered.** For this user base, "email support about my porn-accountability account" is a step many will not take; the practical outcome is account abandonment.

**Fix.**
1. Add `website/src/app/reset-password/page.tsx` that reads `?token=`, collects a new password, and POSTs to `{API}/auth/reset-password`.
2. Set `APP_BASE_URL=https://www.remainfaithful.com` in Fly (and add it to `fly.toml`'s `[env]` so it is version-controlled and cannot silently fall back).
3. Change the `getAppBaseURL()` default from the dead domain to `https://www.remainfaithful.com`, so a missing env var degrades to a working URL rather than a broken one.
4. Add a test asserting the generated reset URL's host is a domain you control.

---

# P1 findings

## F-07 — Google Sign-In does not validate the `aud` claim

**Severity: P1** · `backend/internal/handler/social_auth.go:185-209`

`verifyGoogleIDToken` calls Google's tokeninfo endpoint, which does verify the signature — but the response's `aud` (the OAuth client the token was issued *for*) is never read or checked, and neither is `email_verified`:

```go
var info struct {
    Sub   string `json:"sub"`
    Email string `json:"email"`
    Name  string `json:"name"`
}
```

**Failure scenario.** Any Google ID token issued to **any OAuth client in the world** is accepted. An attacker runs any app with Google Sign-In, has the victim sign into it (or uses a token harvested from any other app the victim uses), and replays that token to `/auth/google`. Combined with the email-fallback linking at `social_auth.go:119-127`, this takes over the victim's Remain Faithful account.

**Fix.** Parse `aud` from the tokeninfo response and reject unless it equals your iOS client ID (and web client ID, if used). Reject unless `email_verified == "true"`. Prefer verifying Google's JWKS locally rather than the tokeninfo round-trip. Also give that `http.Get` a timeout — `http.DefaultClient` has none, so a hung Google connection pins a goroutine.

## F-08 — Login user-enumeration via bcrypt timing; the mitigation does not work [reproduced]

**Severity: P1** · `backend/internal/handler/auth.go:239-243`

```go
if err != nil {
    // Constant-time response to prevent user enumeration.
    bcrypt.CompareHashAndPassword([]byte("$2a$10$placeholder"), []byte(req.Password)) //nolint
    writeError(w, http.StatusUnauthorized, "invalid email or password")
    return
}
```

`"$2a$10$placeholder"` is not a valid bcrypt hash. `CompareHashAndPassword` rejects it immediately with `hashedSecret too short to be a bcrypted password` **without performing any key derivation**, so the intended time-padding never happens.

**Reproduction** (20 iterations each, same machine):

```
placeholder CompareHashAndPassword err = crypto/bcrypt: hashedSecret too short to be a bcrypted password
KNOWN-EMAIL path  (real bcrypt hash): avg 43.892585ms
UNKNOWN-EMAIL path (placeholder)    : avg 101ns
ratio: 430742x
```

A 43 ms versus 101 ns difference is trivially measurable over the network.

**Failure scenario.** An attacker submits `POST /auth/login` with a candidate email and junk password and times the response. ~44 ms means the address has an account; ~0 ms means it does not. At the allowed 10 req/min/IP this confirms 14,400 addresses/day per IP, and the limiter is per-IP so a small proxy pool scales it linearly. For this product, confirming that an address belongs to a Remain Faithful user discloses that the person uses a pornography-accountability app — a sensitive inference about a named individual.

**Fix.** Generate a real bcrypt hash of a dummy password **once at startup** at the same cost factor and compare against that constant:

```go
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password-for-timing"), bcrypt.DefaultCost)
// in the not-found branch:
_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
```

## F-09 — `CreateRelationship` is a direct user-enumeration oracle and leaks names without consent

**Severity: P1** · `backend/internal/handler/relationships.go:41-48`, `:88-101`

```go
err := h.DB.QueryRowContext(r.Context(),
    `SELECT id FROM users WHERE email = $1`, req.PartnerEmail).Scan(&partnerID)
if err != nil {
    writeError(w, http.StatusNotFound, "no account found for that email")
    return
}
```

**Failure scenario.** Any authenticated user submits an arbitrary email and reads the status code: `404` = no account, `201` = account exists. Unlike F-08 this needs no timing analysis — the server states it outright. The attacker then calls `GET /relationships`, which joins `users` and returns the target's **real name and email** (`relationships.go:90`), for a person who never consented to any connection. The unconsented `relationships` row is also created as a side effect.

**Fix.** Return `201` with `{"status":"invited"}` regardless of whether the address resolves (mirroring the `ForgotPassword` approach at `auth.go:301-305`), and route unknown addresses into `relationship_invites` instead. Do not return the counterparty's name or email in `ListRelationships` until the relationship is `accepted`.

## F-10 — No email verification anywhere; registering with an invited address hijacks that person's accountability feed

**Severity: P1** · `backend/internal/handler/auth.go:71,84-139,145-177`

`Register` performs no ownership proof on the email address, and immediately auto-redeems every pending invite addressed to it:

```go
// auth.go:71
h.redeemPendingInvites(r.Context(), id, req.Email)
```

which accepts partner invites (`status = 'accepted'`, `auth.go:124-132`) and joins groups (`auth.go:172-176`).

**Failure scenario.** An attacker learns that `victim@example.com` was invited to a group or as a partner — from a forwarded email, a shared inbox, a guessable work address, or simply by inviting the address themselves from a second account. They register a Remain Faithful account using that address. Registration succeeds (no verification), and every pending invite is auto-accepted. The attacker is now the inviter's accountability partner and/or a member of their group, receiving their flagged-content alerts. The legitimate invitee can never register, because the address is taken.

**Fix.** Send a verification email on registration and gate `redeemPendingInvites` on verified status. At minimum, require the user to click the specific invite token rather than auto-redeeming on an email string match.

## F-11 — PIN brute force is invisible to partners and the limiter resets on every deploy

**Severity: P1** · `backend/internal/handler/protection.go:259-297`, `:23-26`, `:87-98`; `RemainFaithful/PartnerPINManager.swift:50-68`

Two compounding problems.

**The wrong-PIN alert is client-side only.** The server's `VerifyPIN` returns `{"success": false}` and notifies nobody. The partner alert is sent *by the app, about itself*:

```swift
// PartnerPINManager.swift:50-68
let success = try await APIClient.shared.verifyPartnerPIN(pin: pin)
if !success { pinLogger.notice("Wrong PIN attempt") }
...
func notifyWrongAttempt(userName: String) async {
    try? await APIClient.shared.sendProtectionAlert(type: "wrong_pin_attempt", ...)
}
```

Anyone calling `POST /protection/pin/verify` directly — with the token read straight out of App Group UserDefaults (F-13) — brute-forces the PIN with **zero partner visibility**. And per F-02 even the cooperative path is `try?`-swallowed.

**The rate limiter is per-process and in-memory** (`protection.go:23-26`), so it resets on every deploy, crash, and Fly machine restart, and multiplies by machine count. With `auto_stop_machines`/`auto_start_machines` (F-12) restarts are routine. A 4-digit PIN is 10,000 combinations.

**Fix.** Send the `wrong_pin_attempt` protection alert from `VerifyPIN` on the server, on failure, before returning. Move the attempt counter into Postgres keyed on `relationships.id` with a persisted lockout timestamp. Add exponential backoff and notify the partner when the lockout threshold is reached.

## F-12 — `min_machines_running = 0` puts the tamper-detection sweep to sleep exactly when it is needed

**Severity: P1** · `backend/fly.toml:14-16`, `backend/cmd/server/main.go:71-73`

```toml
auto_stop_machines = "stop"
auto_start_machines = true
min_machines_running = 0
```

Both watchdogs are in-process goroutines on 5-minute tickers inside the HTTP server:

```go
// main.go:71-73
silenceMinutes := heartbeatSilenceMinutes()
go h.StartHeartbeatWatcher(ctx)
go h.StartHeartbeatSweep(ctx, silenceMinutes)
```

**Failure scenario.** The heartbeat-silence sweep exists to detect the *absence* of traffic from a user. But the machine that runs it is configured to stop when there is no *inbound* traffic. With a small user base — especially overnight, which is precisely the high-risk window — the last active device goes quiet, no other requests arrive, Fly stops the machine, and the sweep that was supposed to notice the silence never runs. When traffic eventually resumes the machine cold-starts with a fresh 5-minute ticker and empty rate-limit and PIN-attempt maps (F-11).

**Fix.** Set `min_machines_running = 1` and remove `auto_stop_machines` for this app — the watchdog is not a request handler and cannot be scale-to-zero. Better still, move both sweeps out of the web process into a Fly scheduled machine or an external cron hitting an authenticated endpoint, so detection does not depend on web-tier liveness.

## F-13 — Token refresh desyncs the extensions, silently stopping event upload and heartbeats

**Severity: P1** · `RemainFaithful/APIClient.swift:186-193`, `AuthState.swift:48-49`, `SampleHandler.swift:543,700`

`AuthState.setSession` writes the token to **both** stores:

```swift
// AuthState.swift:48-49
keychain.set(token, for: "authToken")
sharedDefaults?.set(token, forKey: "authToken")
```

but `refreshTokenIfNeeded` writes to the Keychain **only**:

```swift
// APIClient.swift:186-193
func refreshTokenIfNeeded() async throws {
    guard let expiry = AuthState.shared.tokenExpiresAt,
          expiry.timeIntervalSinceNow < 3600 else { return }
    let resp: [String: String] = try await post("/auth/refresh", body: [String: String]())
    if let newToken = resp["token"] {
        KeychainHelper.shared.set(newToken, for: "authToken")   // sharedDefaults NOT updated
    }
}
```

All three extensions read from `sharedDefaults`: `SampleHandler.swift:543` (event upload), `:700` (heartbeat), `:321` (lockout alert), `DeviceActivityMonitorExtension.swift:81`, `ActivitySelectionManager.swift:163`.

**Failure scenario.** The app refreshes the token. App Group UserDefaults keeps the old one. 24 hours later the old token expires. From then on: event uploads get `401` and silently divert to a retry queue capped at 50 entries; heartbeats get `401` which `sendHeartbeat` **ignores entirely** — it logs the status code and returns:

```swift
// SampleHandler.swift:710-713
if let (_, resp) = try? await URLSession.shared.data(for: req),
   let http = resp as? HTTPURLResponse {
    logger.info("Heartbeat response: \(http.statusCode)")
}
```

The device believes it is heartbeating. The server records nothing. The sweep fires one silence alert per relationship (then never again, per F-14), and the user is permanently unmonitored with no further signal.

**Fix.** Write to both stores in `refreshTokenIfNeeded` — ideally by routing it through `AuthState.setSession` so there is one code path. Treat `401` in `sendHeartbeat` and `uploadEvent` as a hard error that triggers re-auth and surfaces a user-visible "signed out — monitoring stopped" state.

## F-14 — Heartbeat silence alerts fire at most once, then never again

**Severity: P1** · `backend/internal/handler/heartbeat_sweep.go:59-68,96-101`

The sweep's re-alert guard is:

```sql
AND (r.last_silence_alert_at IS NULL
  OR r.last_silence_alert_at < u.last_heartbeat_at)
```

After an alert, `last_silence_alert_at = NOW()` (`:96-99`). The re-arm condition requires `last_heartbeat_at` to advance past it — which only happens if the user **resumes** heartbeating. A user who goes dark and stays dark fails the condition forever.

**Failure scenario.** A user stops monitoring permanently. Thirty minutes later their partner gets exactly one push. If that push is missed — phone off, Focus mode, swiped away, or lost (and per F-03 nothing is stored, so there is no in-app trace) — there is **no further signal for the rest of time**. A user who wants to go dark can trigger the single alert at 3 a.m. and be confident it will be missed.

**Fix.** Re-alert on a schedule while silence persists — e.g. change the guard to `r.last_silence_alert_at < NOW() - INTERVAL '24 hours'` so partners get a daily reminder, with escalating copy ("X has been dark for 3 days"). Persist each one per F-03 so a missed push still shows in-app.

## F-15 — Users who never send a first heartbeat are never flagged as unmonitored

**Severity: P1** · `backend/internal/handler/heartbeat_sweep.go:63`, `heartbeat.go:68`

Both watchdogs require `last_heartbeat_at IS NOT NULL`.

**Failure scenario.** A user signs up, connects a partner, and never grants Family Controls, never starts Deep Scan, or deletes the app before the first heartbeat. `last_heartbeat_at` stays `NULL`, so they are excluded from both the sweep and the watcher. The partner sees a connected partnership and silence, which is indistinguishable from "they're doing great." The highest-risk state — never actually onboarded — is the one state the safety net cannot see.

**Fix.** Treat `last_heartbeat_at IS NULL AND users.created_at < NOW() - INTERVAL '1 hour'` as a distinct `never_started_monitoring` alert to partners, and surface setup completion in the partner-facing UI.

## F-16 — Unvalidated client `detail` string is written to server logs, carrying user names today

**Severity: P1** · `backend/internal/handler/protection.go:107-110,166-168`; iOS call sites

`SendProtectionAlert` strictly allowlists `Type` (`protection.go:119-122`) but accepts `Detail` as **arbitrary unbounded free text** and logs it:

```go
// protection.go:166-168
if detail != "" {
    log.Printf("[protection] alert user=%d type=%s detail=%q", userID, alertType, truncateForLog(detail, 80))
}
```

The code comment two lines later states "Only metadata is sent — never screen content." That is true of the APNs payload (`protectionAlertBody` ignores `detail`, `:198`) but **not** of the log.

Today's clients already put PII in it:

```swift
// SettingsView.swift:312
detail: "\(userName) disabled the App Lockout feature."
// PartnerPINManager.swift:67
detail: "\(userName) entered a wrong Protection PIN."
```

**Failure scenario.** Every protection alert writes a user's real name into Fly logs. The field is unbounded and unvalidated, so any future client change — or any modified client — can push screen-derived text through it straight into log retention, breaking the literal privacy promise through a channel nobody is auditing. The unbounded length is also a log-flooding cost vector.

**Fix.** Since `protectionAlertBody` already ignores `detail`, stop accepting it: remove it from the request struct, or cap it at ~200 chars, restrict it to an allowlist like the `Type` field, and remove the log line entirely. If a debug breadcrumb is needed, log only `user_id` and `type`.

## F-17 — Recipient email addresses and, in a misconfiguration, live password-reset links are written to logs

**Severity: P1** · `backend/internal/email/email.go:128-131,158,170`; `handler/invites.go:75,102`

```go
// email.go:128-131
if c.apiKey == "" {
    log.Printf("[sendgrid] SENDGRID_API_KEY not set — logging email instead of sending to=%s subject=%q", toEmail, subject)
    fmt.Printf("[email] to=%s subject=%q\n%s\n", toEmail, subject, plain)
    return nil
}
```

`plain` for a password reset **contains the live reset URL with its token** (`email.go:113-117`).

```go
// email.go:158, 170
log.Printf("[sendgrid] sending request to=%s from=%s subject=%q", toEmail, c.from, subject)
log.Printf("[sendgrid] response status=%d to=%s body=%s", resp.StatusCode, toEmail, string(respBody))
// invites.go:75
log.Printf("[invite] partner email error (existing user) to=%s: %v", req.Email, sendErr)
```

**Failure scenario.** Normal operation writes the email address of every user who resets a password, every invitee, and every contact-form sender into Fly logs — for a service where "has an account here" is itself sensitive. If `SENDGRID_API_KEY` is ever unset or mis-rotated, the service **silently stops sending mail and instead prints password-reset tokens to stdout**, returning `nil` so the caller believes delivery succeeded. Anyone with log access can then reset any account whose user triggered a reset.

**Fix.** Never log full addresses — log a hash or `user_id`. Never log `plain` or the SendGrid response body. Make the missing-API-key path return an error rather than `nil`, and fail startup in production if `SENDGRID_API_KEY` is unset.

## F-18 — Stolen tokens are irrevocable and infinitely refreshable; password reset does not log sessions out

**Severity: P1** · `backend/internal/auth/jwt.go:20-33`, `handler/auth.go:266-281,339-393`

JWTs are stateless with a 24-hour lifetime and no `jti`, no server-side session table, and no deny-list. `/auth/refresh` (`main.go:303`) re-signs for any caller holding a currently valid token. `ResetPassword` updates `password_hash` and deletes the reset token but **does not invalidate outstanding JWTs**.

**Failure scenario.** An attacker obtains a token — via F-01, via the plaintext copy in App Group UserDefaults (F-19), or from a device backup. They call `/auth/refresh` every 24 hours indefinitely. The victim notices something wrong and changes their password; **the attacker's session survives**, because nothing ties the token to the credential. There is no "sign out all devices" control anywhere in the product.

**Fix.** Add `users.token_version INT NOT NULL DEFAULT 0`, include it in the JWT claims, and verify it in `auth.Middleware`. Increment on password reset, on password change, and on an explicit "sign out everywhere" action. Refuse `/auth/refresh` when the version is stale.

---

# P2 findings

## F-19 — The JWT is stored in App Group UserDefaults in plaintext, where it is captured in device backups

**Severity: P2** · `RemainFaithful/AuthState.swift:48-49`, `KeychainHelper.swift:5-11`

The token is deliberately duplicated out of the Keychain into the App Group plist so the extensions can read it (documented at `KeychainHelper.swift:9`). A shared-container plist is included in unencrypted Finder/iTunes backups, where a long-lived, infinitely-refreshable bearer token (F-18) is recoverable in cleartext — unlike the Keychain copy, which is protected by the keychain backup key.

**Fix.** Use a Keychain **access group** shared between the app and its extensions (`kSecAttrAccessGroup` with the `$(AppIdentifierPrefix)` entitlement) and read the token from there in all targets. The code comment indicates this was tried and abandoned over entitlement setup; it is worth retrying, as it removes the plaintext copy entirely. If it must stay, reduce token lifetime sharply and add F-18's revocation.

## F-20 — Partner push throttle silently drops severe alerts, including self-harm

**Severity: P2** · `backend/internal/handler/events.go:192-202`

```go
SELECT COUNT(*) FROM alerts
WHERE  recipient_user_id = $1
  AND  created_at > NOW() - INTERVAL '5 minutes'
  AND  id < $2
...
if recentCount > 0 { continue }
```

The throttle is keyed on the **recipient** across *all* monitored users and ignores severity.

**Failure scenario.** A partner supports three people. Person A triggers an `informational` app-usage alert at 10:00. At 10:02 Person B's device flags `self_harm` at `severe` — whose copy reads *"needs you right now … Reach out immediately"* (`events.go:264`). The push is dropped. The alert row exists in the database, but no notification is delivered. The highest-severity signal in the product is suppressed by the lowest.

**Fix.** Scope the throttle to `(recipient_user_id, event.user_id)` and exempt `severity = 'severe'` entirely. Consider per-category throttling so a noisy `app_usage` stream cannot mask `adult_content`.

## F-21 — `npm run lint` is unconfigured and blocks on an interactive prompt; CI runs no website checks

**Severity: P2** · `website/package.json:8`; `.github/workflows/backend-ci.yml`

`npm run lint` emits:

```
? How would you like to configure ESLint?
❯  Strict (recommended)
   Base
   Cancel
```

There is no `.eslintrc*` in `website/`. The command never lints — it waits for input, and would hang a CI job indefinitely. `next build` therefore also skips ESLint.

Separately, `backend-ci.yml` runs only `go test ./...` — no `go vet`, no website build, no website tests, and **no secret scanning**. The `.githooks/pre-commit` scanner exists and is well written, but it is opt-in per clone (`git config core.hooksPath .githooks`), which undercuts the site's claim of scanning "every contribution" (F-25).

**Fix.** Commit an `.eslintrc.json` with `{"extends": "next/core-web-vitals"}`. Add a `website-ci.yml` running `npm ci && npm run lint && npm run build && npm test`. Add `go vet ./...` to the backend job and a `gitleaks`/`trufflehog` step to both.

## F-22 — Group member cap can be exceeded by concurrent invites

**Severity: P2** · `backend/internal/handler/invites.go:254-260,288-293`

`checkCanInviteToGroup` counts members, then `inviteEmailToGroup` inserts with **no lock and no re-check** — a classic TOCTOU gap. The registration path gets this right with `SELECT id FROM groups WHERE id = $1 FOR UPDATE` (`auth.go:189-202`); the direct-add path does not.

**Failure scenario.** Two admins invite simultaneously into an 11-member group. Both count 11, both insert, the group reaches 13. The documented 12-member cap — part of the product's stated design — is not enforced.

**Fix.** Move the cap check inside the same transaction as the insert, behind `SELECT ... FOR UPDATE` on the `groups` row, reusing `joinGroupFromInvite`.

## F-23 — Partner-side rows are undeletable, producing a dead button and no way to stop being watched

**Severity: P2** · `backend/internal/handler/relationships.go:200-237`, `:88-101`

`ListRelationships` deliberately returns both directions — rows where the caller is `user_id` and rows where the caller is `partner_id`. But `DeleteRelationship` scopes to one:

```go
`DELETE FROM relationships WHERE id = $1 AND user_id = $2`
```

**Failure scenario.** A user is listed as someone else's accountability partner (including via F-09's unconsented creation). The row appears in their list with an ID; the UI offers a remove action; the call returns `404 "relationship not found"`. They cannot decline the role, and they keep receiving that person's alerts.

**Fix.** Allow deletion from either side: `WHERE id = $1 AND (user_id = $2 OR partner_id = $2)`, notifying the counterparty either way. Keep the current one-sided behaviour only if the UI hides the control for partner-direction rows — but a partner should be able to withdraw.

## F-24 — `VerifyPIN` picks an arbitrary PIN when more than one partner has set one

**Severity: P2** · `backend/internal/handler/protection.go:284-288`

```sql
SELECT pin_hash FROM relationships
WHERE user_id = $1 AND status = 'accepted' AND pin_hash IS NOT NULL
LIMIT 1
```

`LIMIT 1` with no `ORDER BY` returns an unspecified row. `GetPINStatus` meanwhile reports `is_set` if *any* relationship has a PIN (`:242-245`).

**Failure scenario.** Two partners each set a PIN. The app shows "PIN is set." Only one of the two PINs opens the gate, and which one can change between queries as Postgres changes plan or row order. The other partner's PIN appears broken, and failures are indistinguishable from "no PIN set" because the not-found branch also returns `{"success": false}` (`:289-293`).

**Fix.** Decide the model explicitly. If one PIN per user: enforce it with a partial unique index and reject a second. If per-relationship: have the client send `relationship_id` and verify against that row. Return a distinct response for "no PIN configured" versus "wrong PIN".

## F-25 — "Data at rest is AES-256 encrypted" is not backed by any code

**Severity: P2** · `website/src/app/privacy-architecture/page.tsx:266`

A grep for `crypto/aes`, `AES`, and `encrypt` across `backend/` returns **no application-level encryption**. The database stores alert metadata, names, and emails as plaintext columns; the only encryption is whatever Fly applies to the underlying volume.

The overwhelming majority of the site's privacy claims **do** hold up — "screenshots, OCR text, and raw screen content never leave your device" is accurate (`SampleHandler.uploadEvent` sends only category/severity/summary/timestamp, and `summary` is server-side allowlisted at `events.go:22-41`), and the DRM limitation is stated plainly and correctly. That makes the two unsupported claims worth tightening, because the page's credibility rests on being literally true.

The second is `privacy-architecture/page.tsx:297`: "We run pre-commit secret scanning on every contribution." The hook is real and good, but it is opt-in per clone and absent from CI (F-21).

**Fix.** Reword to "Our hosting provider encrypts database volumes at rest" — or implement column-level encryption for the sensitive fields and keep the stronger claim. For the scanning claim, either add the CI step (making it true) or soften to "contributors run pre-commit secret scanning."

## F-26 — `/api/donate` is unauthenticated and unthrottled

**Severity: P2** · `website/src/app/api/donate/route.ts:4-45`

No rate limit, no bot check, no honeypot — unlike `/contact`, which has both a honeypot and a server-side limiter. Each call creates a real Stripe Checkout Session.

**Failure scenario.** A script POSTs continuously, creating unbounded Stripe sessions and Vercel function invocations. Cost accrues on both, and the Stripe dashboard fills with abandoned sessions. `/api/waitlist` and `/api/contact` on the website have the same exposure.

Credit where due: `checkoutReturnUrls` is correctly allowlisted — non-localhost origins are forced to `PRODUCTION_SITE_ORIGIN`, so there is **no open redirect** here, and the behaviour is covered by tests (4 passing).

**Fix.** Add Vercel's built-in rate limiting or a KV-backed limiter keyed by IP, plus a honeypot field mirroring `/contact`.

## F-27 — No cap on partner count, enabling event-fanout amplification

**Severity: P2** · `backend/internal/handler/invites.go:27-108`, `relationships.go:20-76`, `events.go:162-172`

Groups are capped at 12 (`invites.go:221`). Relationships have **no cap at all**. `CreateEvent` performs one `INSERT INTO alerts` and one push fan-out per partner, in a loop, inside an unbounded goroutine.

**Failure scenario.** A user adds thousands of partners (each add is one API call at 10/min/IP, or unthrottled on the authenticated `/relationships` route — note `POST /relationships` has **no** `authenticatedRateLimiter`, unlike `/events` and `/panic`). A single event then triggers thousands of row inserts and APNs requests from a 256 MB machine. This is a self-inflicted amplification DoS and a direct Postgres/APNs cost multiplier.

**Fix.** Cap partners per user (10 is generous for the use case) and enforce it in both `InvitePartner` and `CreateRelationship`. Add `authenticatedRateLimiter` to `POST /relationships` and `/relationships/invite`. Batch the alert inserts into a single multi-row statement.

## F-28 — Password-reset email bombing

**Severity: P2** · `backend/internal/handler/auth.go:286-334`, `main.go:117`

`/auth/forgot-password` is limited to 10/min/IP with **no per-account cap**.

**Failure scenario.** An attacker sends 14,400 reset requests/day for one victim address from a single IP — more from a proxy pool. The victim is buried in email, each request invalidates the previous token via `ON CONFLICT (user_id) DO UPDATE` (`auth.go:316-321`), and SendGrid volume and sender reputation both suffer.

**Fix.** Add a per-user cap (e.g. 3 reset emails/hour) checked against `password_reset_tokens`, independent of IP. Suppress the email if an unexpired token was issued within the last 5 minutes.

## F-29 — `UpdateMe` accepts any email with no validation, no uniqueness handling, and no re-verification

**Severity: P2** · `backend/internal/handler/users.go:47-87`

No format check, and a collision with another account surfaces as `500 "failed to update profile"` rather than `409`. Changing the email also leaves the JWT carrying the stale address, which `RefreshToken` then re-signs forward indefinitely (`auth.go:274-275`).

**Fix.** Validate the format, map the unique-violation to `409`, and require verification of the new address before it takes effect. Re-issue the JWT on change.

## F-30 — `APP_DELETED` notification is sent to the token that was just reported dead

**Severity: P2** · `backend/internal/handler/notifications.go:343-368`

On `ErrInvalidToken` the code marks the token inactive and then sends the `APP_DELETED` notice to **that same dead token**:

```go
h.markTokenInactive(ctx, token)
deleted := &apns.Notification{
    DeviceToken: token,   // the token APNs just rejected with 410
    ...
}
```

The send can only fail. The feature — telling partners that someone removed the app — never works.

**Fix.** Resolve the *partner's* `user_id` from the dead token, then send the notice to that user's other active tokens via `notifyPartnerByID`. Skip if they have none.

## F-31 — Group streak counts clean events, contradicting the personal stats screen

**Severity: P2** · `backend/internal/handler/groups.go:147-156` vs `users.go:273-278`

`GetGroup` computes `streak_days` from `MAX(e.timestamp)` over **all** events. `GetUserStats` correctly excludes `category NOT IN ('clean')`. But clean events are a normal product output — `"Content reviewed — no concerns"` is an allowlisted summary (`events.go:27`).

**Failure scenario.** An actively monitored, fully clean user generates `clean` events. Their personal dashboard shows a growing streak; the group screen shows **0 days** for the same person on the same day. The user whose behaviour is best looks worst to their group.

**Fix.** Add `AND e.category <> 'clean'` to the LATERAL subquery, matching `users.go:277`. Better, extract one shared SQL helper so the two screens cannot drift.

---

# P3 findings

| ID | Finding | Evidence | Fix |
|---|---|---|---|
| F-32 | `ListEvents`/`ListAlerts` hard-code `LIMIT 100` with no pagination — older history is permanently unreachable. | `events.go:289`, `alerts.go:89` | Add cursor pagination on `(timestamp, id)`. |
| F-33 | No index supporting `ORDER BY timestamp DESC` per user; only `idx_events_user_id` exists. Sorts every row for heavy users. | `main.go:432` | `CREATE INDEX idx_events_user_ts ON events(user_id, timestamp DESC);` |
| F-34 | No index on `users.last_heartbeat_at`, which the 5-minute sweep filters on. | `main.go:533`, `heartbeat_sweep.go:63` | Add a partial index `WHERE last_heartbeat_at IS NOT NULL`. |
| F-35 | Client-supplied `timestamp` on `POST /events` is unvalidated; a modified client can backdate events so they sort to the bottom of the partner's list. Invalid strings yield `500`, not `400`. | `events.go:107-113` | Reject timestamps outside `[now-24h, now+5m]`; parse before the insert. |
| F-36 | Password-reset tokens stored in plaintext. | `main.go:451-457`, `auth.go:313` | Store `sha256(token)`; compare hashes. |
| F-37 | `markTokenInactive` is globally scoped by token with no `user_id`, and `RegisterDeviceToken` accepts any token string without proving device ownership. | `notifications.go:413-415`, `:42-50` | Scope the update by `user_id`; validate token format/length. |
| F-38 | Per-push `fmt.Printf`/`log.Printf` on every APNs send and every SendGrid call — significant log volume and cost at scale. | `apns.go:190,202,205`, `notifications.go:409` | Drop to debug level or sample. |
| F-39 | Any group member can rename the group and rewrite the 8,000-char covenant; the role is read and then explicitly discarded (`_ = role`). | `groups.go:232-240` | Restrict destructive edits to `role = 'admin'`. |
| F-40 | `LeaveGroup` does not promote a replacement admin, unlike `DeleteMe`; the sole admin leaving orphans the group. | `groups.go:430-469` vs `users.go:134-177` | Reuse `promoteAdminsBeforeDeletion`. |
| F-41 | 7 pages lack a canonical URL: `/privacy`, `/terms`, `/privacy-architecture`, `/group-setup-guide`, `/invite`, `/accept-invite`, `/join/[id]`. | `website/src/app/*/page.tsx` | Add `alternates: { canonical }` to each (and `robots: { index: false }` for the invite routes). |
| F-42 | `verifyGoogleIDToken` uses `http.Get` on `http.DefaultClient`, which has no timeout; the ID token is string-concatenated into the URL unescaped. | `social_auth.go:186-188` | Use a client with a 5 s timeout and `url.QueryEscape`. |
| F-43 | Column names interpolated into SQL via `fmt.Sprintf`. Not exploitable today (call sites pass literals) but a trap for future edits. | `social_auth.go:112,125,136` | Switch on the two cases with fully literal query strings. |
| F-44 | `clientIP` falls back to the client-controlled `X-Forwarded-For` when `Fly-Client-IP` is absent, allowing rate-limit evasion if the app is ever reachable off the Fly edge. | `main.go:196-208` | Trust only `Fly-Client-IP` in production. |
| F-45 | *Withdrawn.* A 12 MB compiled `backend/server` binary sits in the working tree, but it is already covered by `.gitignore:18` and cannot enter history. No action needed. | `.gitignore:18` | — |

---

# Cost at scale

## Where the money goes

The dominant driver is the **2-minute heartbeat** from the broadcast extension:

```swift
// SampleHandler.swift:687-693
try? await Task.sleep(nanoseconds: 2 * 60 * 1_000_000_000)
```

720 requests/day per device with Deep Scan continuously active; ~60/day at a realistic 2 h/day session. Each one performs a write:

```sql
-- heartbeat.go:35-39
UPDATE users SET last_heartbeat_at = NOW(), heartbeat_screen_state = $2 WHERE id = $1
```

**This is the single most important scaling defect.** Updating the `users` row every two minutes generates one dead tuple per heartbeat on the table that every authentication, invite lookup, and name resolution reads. At 10k users that is roughly 660k dead tuples/day on `users`, driving continuous autovacuum pressure and index bloat on the hottest table in the schema.

Assumptions below: 50% of registered users active on a given day; those users average 2 h/day of Deep Scan (≈60 heartbeats); 10% of active users run it near-continuously (≈720); ~2 flagged events/active user/day; ~20 foreground API calls/active user/day. Fly and SendGrid list prices as of this review.

## Estimates

### 100 users — **$5–30/month**

| Line | Cost | Notes |
|---|---|---|
| Fly app | ~$2 | 1 × shared-cpu-1x 256 MB. Should be `min_machines_running = 1` (F-12), which is ~$1.94/mo rather than free. |
| Fly Postgres | ~$3 | shared-cpu-1x 256 MB + 1 GB volume. |
| SendGrid | $0 | Free tier (100/day) covers ~50 invites/month. |
| Vercel | $0–20 | Hobby suffices; Pro if commercial. |
| APNs | $0 | Free. |

≈5,000 requests/day (0.06 req/s). Nothing strains.

### 1,000 users — **$50–90/month**

| Line | Cost | Notes |
|---|---|---|
| Fly app | $8–15 | 1–2 × shared-cpu-1x 512 MB. 256 MB is tight given unbounded per-request goroutines in `CreateEvent`. |
| Fly Postgres | $15–30 | 1–2 GB RAM, 10 GB volume. ~60k events + ~120k alerts/month. |
| SendGrid | $20 | Essentials — past the 100/day free tier. |
| Vercel | $20 | Pro. |
| Logs | ~$5 | Per-push logging (F-38) starts to matter. |

≈50,000 requests/day (0.6 req/s). Comfortable, but `users`-table bloat begins to show without tuned autovacuum.

### 10,000 users — **$150–350/month**

| Line | Cost | Notes |
|---|---|---|
| Fly app | $30–60 | 2–4 × shared-cpu-2x 1 GB. |
| Fly Postgres | $60–150 | **The real cost.** Dedicated instance, 4–8 GB RAM, 50–100 GB. |
| SendGrid | $20–90 | Essentials → Pro depending on invite velocity. |
| Vercel | $20 | Pro, plus overage if `/api/donate` is abused (F-26). |
| Logs | $10–30 | Per-push `fmt.Printf` at this volume (F-38). |

≈1M requests/day (~12 req/s average, 30+ at peak). Postgres dominates:

- **~660k `UPDATE`s/day on `users`** — bloat and autovacuum pressure, as above.
- **Unbounded retention.** Nothing in the codebase ever deletes an `events`, `alerts`, or `donation_events` row. ~600k events + ~1.2M alerts/month, growing forever.
- **Unindexed sweep.** `sweepHeartbeatSilence` runs every 5 minutes with no `LIMIT`, joining `users` to `relationships` with no index on `last_heartbeat_at` (F-34).
- **`checkHeartbeats`** likewise scans all users with a non-null heartbeat, unbounded, every 5 minutes.

### The three changes that most reduce the 10k bill

1. **Move heartbeats off the `users` table** into a narrow `user_heartbeats(user_id PK, last_seen, screen_state)` — ideally unlogged or Redis-backed. Removes bloat from the hottest table. Raise the cadence from 2 to 5 minutes (still well inside the 30-minute silence threshold) for a 60% write reduction at no detection cost.
2. **Add retention.** Prune `events`/`alerts` older than 12 months and `donation_events` older than 90 days via a nightly job. Without this, storage grows forever.
3. **Batch alert inserts and APNs sends** in `CreateEvent` (`events.go:162-172`) — currently one round trip per partner — and cap partner count (F-27).

### Uncapped exposures

`/api/donate`, `/api/waitlist`, `/api/contact` have no rate limiting (F-26); `/auth/forgot-password` has no per-account cap (F-28); relationship count is uncapped (F-27). Each is a cost vector with no ceiling, independent of user count.

---

# App Store risk

| Risk | Assessment |
|---|---|
| **3.1.1 in-app purchase** | **Clear.** A full grep of all Swift sources for `donat`/`Donat`/`stripe`/`Stripe`/`tip` returns only one unrelated hit (a `"starvation tips"` keyword in the classifier at `SampleHandler.swift:94`). The backend retains `/donations/*` endpoints and the website has `/api/donate`, but the app neither calls nor links to them. `SettingsView.swift:1674` links to `/app-guide`; **keep that page free of donation CTAs** — a link-out to purchase from an app-referenced page is the realistic 3.1.1 exposure. |
| **5.1.1 / 5.1.2 privacy** | The literal privacy promise holds in code (F-25 context). The exposure is F-16 and F-17 — PII written to server logs — which contradicts the posture claimed on `/privacy-architecture` and would be hard to defend in a privacy review. |
| **5.1.1(v) account deletion** | **Satisfied.** `DELETE /users/me` (`users.go:94-130`) hard-deletes with cascade and is reachable from Settings. |
| **Family Controls entitlement** | The strongest latent risk. Apple grants this for self-monitoring, and the app is correctly positioned that way. But F-05 (adding anyone to a group by email with no consent, then receiving their content alerts) makes the product function as **covert third-party surveillance**. If demonstrated, that is a plausible entitlement-revocation or removal trigger, independent of any guideline citation. Fix F-05 on that basis alone. |
| **4.2 minimum functionality** | Not a concern — the app is substantial. |
| **Time-sensitive interruption level** | `SendPanicAlert` sets `"interruption-level": "time-sensitive"` (`notifications.go:171`). *Suspected*: requires the Time Sensitive Notifications entitlement; if absent, APNs rejects the push and panic alerts silently fail. Verify against the provisioning profile. |

---

# What is solidly built

Worth stating, because it is a meaningful share of the codebase:

- **The privacy promise is real in code.** `uploadEvent` transmits only `category`, `severity`, `summary`, `timestamp` — and `summary` is constrained server-side to a 12-entry allowlist plus one anchored regex (`events.go:22-41`). The `detail` field is deliberately excluded from APNs bodies (`protection.go:198`). The architecture does what the website says.
- **Event summary validation** is genuinely well designed — a modified client cannot push arbitrary text to a partner's lock screen.
- **JWT verification** is correct: HMAC method pinned, issuer checked, `WithExpirationRequired()`, and the process refuses to start without `JWT_SECRET` (`jwt.go:37-65`). The `alg=none` and algorithm-confusion classes are closed.
- **SQL** is parameterized throughout; no injection was found.
- **Stripe webhooks** verify signatures and are idempotent via a `UNIQUE` event-ID log (`donations.go:23-28`).
- **Email templates** HTML-escape all interpolated values and strip CRLF from headers (`email.go:89-109`).
- **`checkoutReturnUrls`** correctly allowlists origins — no open redirect, and it is tested.
- **No secrets in the repo**, and `.githooks/pre-commit` is a well-constructed scanner.
- **Group member cap** is correctly enforced under `SELECT ... FOR UPDATE` on the registration path (`auth.go:181-217`) — the model to copy for F-22.
- `go vet` clean, `go test` green, website builds and tests pass.

---

# Recommended order of work

1. **F-01** Apple token signature verification + remove email-fallback linking — unauthenticated takeover of any account.
2. **F-05** Group-invite consent — unconsented disclosure of flagged-content alerts; also the Family Controls entitlement risk.
3. **F-06** Ship `/reset-password` and set `APP_BASE_URL` — every password account is currently unrecoverable.
4. **F-04** Persist the partner invite token — silent total monitoring failure with a misleading "connected" response.
5. **F-02 + F-03** Durable protection alerts, client queue and server rows — closes the airplane-mode bypass.
6. **F-07** Google `aud` validation.
7. **F-12** `min_machines_running = 1` — one line; restores the watchdog.
8. **F-13** Write the refreshed token to both stores — one line; stops silent extension death.
9. **F-16 + F-17** Stop logging `detail`, email addresses, and email bodies.
10. **F-14 + F-15** Re-alert on continued silence; detect never-started users.

Items 7 and 8 are one-line changes with outsized impact and should go in the first patch regardless of ordering.
