# API changes for build 10

This is the contract for the iOS agent. Build 9 keeps working: old paths stay, and new JSON fields are additive. Swift `Codable` ignores unknown keys.

Base URL is the Fly app (`https://api.remainfaithful.com` unless `BACKEND_URL` says otherwise). Protected routes need `Authorization: Bearer <jwt>`.

Sessions now carry `tv` (token version). Tokens issued before this field existed decode as `tv = 0` and stay valid until the account's `token_version` is bumped past 1. A password reset or a confirmed email change bumps the version and the old token returns 401 `invalid or expired token`. `POST /auth/refresh` signs a new token with the current version. Call it after login and after a password change. Do not keep using a token that just got 401.

## Auth

### `POST /auth/apple` and `POST /auth/google`

Unchanged request bodies (`identity_token` / `id_token`, plus optional `name`).

The server now verifies the token before it creates or links an account.

- Apple: RS256 signature against `https://appleid.apple.com/auth/keys`, `iss` is `https://appleid.apple.com`, `aud` is the bundle id (`APPLE_AUDIENCE`, else `APPLE_BUNDLE_ID`, else `APNS_BUNDLE_ID`, else `com.remainfaithful.app`), and `exp` is in the future.
- Google: the token is POSTed to Google's tokeninfo endpoint. `aud` or `azp` must match `GOOGLE_CLIENT_IDS` (comma-separated) or `GOOGLE_CLIENT_ID`. If an email is present it must be verified.

A provider-verified email can attach to an existing password account only when that account's Apple or Google id is empty. A different subject on the same email is `409`. An unverified Google email is `401`. Placeholder emails do not redeem invites.

New verified social accounts redeem pending invites. Password signups do not, until email confirmation.

### `POST /auth/register`

Still `201` with `id`, `name`, `email`, `created_at`. Adds `email_verified: false`. Still does not return a JWT. Duplicate email is still `409` `email already registered` so build 9 does not treat a duplicate as success.

The account can sign in immediately. Partner and group invites for that email stay pending until `POST /auth/verify-email`.

### `POST /auth/login`

Same body. Response `token` plus `user` now includes `email_verified`.

### `POST /auth/verify-email` (public)

Body: `{ "token": "<from the email link>" }`.

- `200` `{ "status": "verified", "email": "..." }` and then pending invites for that email are accepted.
- The same route confirms an email change (`purpose=change` on the website link). Success moves `pending_email` onto `email`, sets `email_verified` true, and bumps `token_version`.
- `400` `token is required` or an invalid or expired token.

### `POST /auth/resend-verification` (JWT)

No body. `200` `{ "status": "sent" }` or `{ "status": "already_verified" }`.

### `POST /auth/forgot-password` and `POST /auth/reset-password`

Forgot still always returns `200` `{ "status": "ok" }`. At most 3 reset emails per account per hour; extra requests still return 200 and send nothing.

Reset body is unchanged: `{ "token", "password" }`. Password minimum is still 8 characters. Success bumps `token_version`, so every previous JWT for that user fails. The link is `https://remainfaithful.com/reset-password?token=...` unless `APP_BASE_URL` or `SITE_URL` is set. It works for one hour.

### `PUT /users/me`

Name-only updates still apply immediately when the email is unchanged.

A different email does **not** replace `email` yet. The response keeps the current `email` and adds `pending_email`. A confirmation link goes to the new address. `409` `that email cannot be used` on a unique clash.

### `GET /users/me`

Adds `email_verified` (boolean). Existing accounts default to `true`.

## Partner invites

`POST /relationships/invite` body is unchanged: `{ "email": "..." }`.

- No account yet: `201` `{ "status": "invited", "email" }`. Stored until that person confirms the email.
- Account exists and is not already accepted: `201` `{ "status": "pending", "relationship_id" }`. They are **not** connected. Alerts still require `accepted`.
- Already accepted: `200` `{ "status": "accepted", "relationship_id" }`.

Build 9 mostly ignores `status`. Do not treat `pending` as connected.

Caps: 10 outgoing relationships per user (`422` `you can have up to 10 partners`). 20 invite emails per user per day across partner and group invites (`429` `invite limit reached for today`). Both invite routes are also limited to 10 requests per minute per user.

### New partner routes (JWT)

| Method and path | Body | Success |
|---|---|---|
| `GET /relationships/invites/pending` | | array of `{ id, inviter_id, inviter_name, created_at, type: "partner" }` |
| `POST /relationships/invites/{id}/accept` | | `{ "status": "accepted", "relationship_id" }` |
| `POST /relationships/invites/{id}/decline` | | `{ "status": "declined" }` |
| `POST /relationships/accept-invite` | `{ "token" }` | `201` `{ "status": "accepted", "relationship_id" }` |
| `POST /relationships/decline-invite` | `{ "token" }` | `{ "status": "declined" }` |

Accept and decline require the signed-in email to match `invitee_email`. Anyone else gets `403` `this invite was sent to a different email address`. Missing or older than 7 days is `404` `invite not found or expired`.

`POST /relationships` (unused by build 9) no longer returns `404` `no account found for that email`. An unknown email is `201` `{ "status": "invited" }`. A known email inserts `pending` (subject to the same caps).

`DELETE /relationships/{id}` may be called by either person in the row. The other person is notified.

## Group invites

`POST /groups/{id}/invite` and `POST /groups/{id}/email-invite` no longer insert `group_members` for an existing user.

- Not a member: `201` `{ "status": "invited", "email" }`. They must accept.
- Already a member: `200` `{ "status": "already_member" }`.
- No account yet: `201` `{ "status": "invited", "email" }`. Redeemed when that email is confirmed, and the group must still be under 12 members at redeem time.

The old `"added"` value is not returned anymore. Build 9 treats several outcomes as sent.

Groups are hard-capped at 12 members inside the accept transaction. A full group returns `422` `group has reached the maximum of 12 members`.

`https://www.remainfaithful.com/join/{id}` does not add anyone. Membership still requires an email invite plus accept (or email confirmation for a brand-new account).

### New group routes (JWT)

| Method and path | Body | Success |
|---|---|---|
| `GET /groups/invites/pending` | | array of `{ id, group_id, group_name, inviter_id, inviter_name, created_at, type: "group" }` |
| `POST /groups/invites/{id}/accept` | | `{ "status": "accepted", "group_id" }` |
| `POST /groups/invites/{id}/decline` | | `{ "status": "declined" }` |
| `POST /groups/accept-invite` | `{ "token" }` | `{ "status": "accepted", "group_id" }` |
| `POST /groups/decline-invite` | `{ "token" }` | `{ "status": "declined" }` |

Pending lists only include invites from the last 7 days addressed to the signed-in email.

## Alerts and protection

`POST /protection/alerts` body is unchanged: `{ "alert_type": "...", "detail": "..." }`. `detail` is not copied into the push body.

Client types the server accepts:

`monitoring_disabled`, `shielding_disabled`, `lockout_disabled`, `deep_scan_stopped`, `wrong_pin_attempt`, `pin_removed`, `pin_changed`, `family_controls_revoked`, `signed_out`.

`signed_out` is new. Send it when the user signs out, before the token is cleared, so partners get a stored alert.

These are stored and are **not** accepted on POST (the server writes them itself):

- `heartbeat_silence` when `last_heartbeat_at` is set and older than `HEARTBEAT_SILENCE_MINUTES` (default 30). A null heartbeat is not an alert. The same recipient is not alerted again for 6 hours.
- `app_removed` when Apple rejects that device token. The dead token is not pushed again.

Each stored protection alert is fanned out to accepted partners and to other members of the user's groups. `POST` returns `200` after the row is stored. A failed insert is `500`.

`GET /alerts` still returns an array. Every item still has a nested `event` object so build 9 can decode the list:

```json
{
  "id": 1,
  "event_id": 0,
  "seen": false,
  "created_at": "...",
  "kind": "protection",
  "alert_type": "signed_out",
  "event": {
    "id": 0,
    "user_id": 42,
    "category": "signed_out",
    "severity": "concerning",
    "summary": "Alex signed out",
    "timestamp": "..."
  }
}
```

`kind` is `event` for normal events and `protection` for protection, tamper, silence, and app-removed rows. `event_id` is `0` when there is no events row. Build 9 drops unknown `event.category` values; that is safe. Build 10 should read `kind` and `alert_type` and show protection rows in the partner inbox.

`POST /events` still creates the event and partner alerts. Timestamps may be RFC3339 or a zoneless `2006-01-02T15:04:05`. Empty means now. More than 5 minutes in the future or more than 7 days in the past is `400`. Informational pushes for one monitored person are throttled to one per 5 minutes. `severity: "severe"` is never throttled by that window.

`GET /users/me/stats` still returns `streak_days`, `best_streak`, and `week` (7 bools, index 0 is six days ago). It adds `week_states` with the same indexes: `before_account`, `clean`, or `flagged`. Use `before_account` instead of painting a red X for days before the account existed. `week` bools are unchanged for build 9. Group streak and `flags_last_30` no longer treat category `clean` as a flag.

Events and alerts older than 365 days are deleted by the silence sweep.

## PIN

`POST /protection/pin/verify` still returns `success`. It now checks every accepted PIN on the user's relationships, not only one row. Reasons:

- `{ "success": true }`
- `{ "success": false, "reason": "no_pin" }`
- `{ "success": false, "reason": "wrong_pin" }`
- `429` `{ "success": false, "reason": "rate_limited", "error": "too many PIN attempts. Try again in 15 minutes" }` after 5 failures in 15 minutes. `Retry-After: 900`.

The server does **not** also send a `wrong_pin_attempt` push. Build 9 already posts that alert. Keep posting it from the app on a wrong PIN so partners are notified once.

## Heartbeat

`POST /heartbeat` is unchanged. Keep sending it whenever monitoring is armed, including when Deep Scan is off. The silence sweep only looks at users who have sent at least one heartbeat.

One machine stays running (`min_machines_running = 1`) so the sweep is not asleep with the app.

## What did not change on purpose

- Register still returns `409` for a known email.
- Invite tokens in email links stay plaintext. Reset and email-confirmation tokens are stored as SHA-256 hashes.
- `GET /alerts` is still the latest 100 rows. No pagination cursor in this build.
- Group rename and sole-admin leave rules are unchanged.
