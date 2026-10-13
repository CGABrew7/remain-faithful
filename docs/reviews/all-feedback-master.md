# All-feedback master list

Deduplicated list from every pre-launch review of commit `784998f` (2026-10-13).

Sources:

| Tag | Source |
|---|---|
| grok | `docs/reviews/round1-grok-4.7.md` on `cursor/review-round1-grok-4.7-4792` |
| claude | `docs/reviews/round1-claude.md` on `review/round1-claude` |
| master | `docs/reviews/round2-master.md` on `cursor/round2-review-80a8` (re-ranked the round-1 reports) |
| ux | `docs/reviews/round2-ux.md` on `cursor/round2-review-80a8` |
| ux-claude | `docs/reviews/round2-ux-claude.md` on `review/round2-ux-claude` |
| chatgpt | Jeff's pasted ChatGPT round 1 (no source access) |
| chatgpt-ux | Independent ChatGPT UX review with no code access |

Priority follows the round-2 master when it re-ranked a security item. UX-only items keep the priority from the UX reviews. **Owner** `backend` or `website` is this push. **Owner** `ios` is the other agent (build 10). Do not edit the iOS folder here.

Fixed product choices are not defects: Family Controls plus a partner notice, DeviceActivity, opt-in Deep Scan via ReplayKit (it cannot see DRM streaming), and App Lockout. Where the words describe something the code does not do, the defect is the mismatch.

---

## P0

| ID | Priority | Owner | Reviewers | Issue |
|---|---|---|---|---|
| A01 | P0 | backend | grok 1, claude F-01, master M1 | Sign in with Apple decodes the identity token and does not check the signature, `iss`, `aud`, or `exp`. A forged payload with a victim email can attach an Apple id and take the account. |
| A02 | P0 | backend | grok 2, claude F-04, master M2, ux U-09, ux-claude U-09, chatgpt 8 | Inviting an existing user as a partner inserts `pending`, returns `"connected"`, and emails a token that is never stored. Alerts require `accepted`, so the partner never hears anything. |
| A03 | P0 | backend | grok 7, claude F-05, master M3, chatgpt 8 | Adding an existing user to a group inserts them immediately. Their later events fan out to every co-member, and the roster exposes name, email, 30-day flag count, and streak. |

## P1

| ID | Priority | Owner | Reviewers | Issue |
|---|---|---|---|---|
| A04 | P1 | backend, website | grok 4, claude F-06, master M4, ux quit 4, ux-claude quit 5 | Password reset links default to `https://app.remainfaithful.app/reset-password`. No page redeems the token. |
| A05 | P1 | ios | grok 3, claude, master M5, ux, ux-claude U-03 | The app drops `violence`, `self_harm`, `app_usage`, and `clean`. One-to-one partners have no inbox. Opening Home marks alerts seen. Alert detail invents a browsing story. |
| A06 | P1 | backend, ios | grok 5, claude F-02, claude F-03, master M6, ux U-15, ux-claude U-15, chatgpt 9 | Protection and silence alerts are a single push with no saved row. Turning monitoring off, clearing the selection, and signing out do not notify partners. Offline `try?` drops the push. |
| A07 | P1 | backend, ios | grok 5, claude F-12, master M7, chatgpt 3 | Heartbeats exist only during Deep Scan. The silence sweep sleeps with `min_machines_running = 0` and does not include group members. Do not alert on a null heartbeat. |
| A08 | P1 | website, ios | grok 6, master M8, ux, ux-claude | The site says partners are notified when a blocked category is attempted, and that time windows exist. There is no ShieldAction extension. The only schedule is all day. Home says monitoring is paused whenever Deep Scan is off. |
| A09 | P1 | ios | grok 9, claude F-13, master M9 | Token refresh writes the Keychain only. The in-memory session and the broadcast extension's app-group copy stay stale. |
| A10 | P1 | ios | grok 11, master M10 | Tier 1 labels gambling, gore, and self-harm hits as explicit sexual content. |
| A11 | P1 | backend, website | grok 12, claude F-10, claude F-29, master M11 | No email verification. Signup auto-accepts invites. `UpdateMe` changes the email immediately. |
| A12 | P1 | backend | grok 13, claude F-07, master M12 | Google sign-in checks `sub` only. No `aud`, `azp`, or `email_verified`. The token is placed in a query string. |
| A13 | P1 | website | ux, ux-claude W-01, chatgpt-ux (install handoff) | No App Store link on the hero, the download section, or invite landing pages. |
| A14 | P1 | ios | master missed 1, ux, ux-claude quit 1, chatgpt-ux RF-01 | Home status is Deep Scan only. "No flags yet, great work" shows when setup never finished. |
| A15 | P1 | ios | master missed 2, ux, ux-claude quit 2, chatgpt-ux RF-02, RF-03 | Onboarding asks for Screen Time permission and never opens the app picker. |
| A16 | P1 | website | ux, ux-claude W-02, master M2 | `/accept-invite` ignores the token and tells an existing user that signing in is enough. |
| A17 | P1 | ios | master missed 3, ux | The support button can work for a group-only user while the screen says a partner is required. Protection alerts do not fan out to group members (server half is A06/A07). |
| A18 | P1 | ios | ux-claude U-01, ux U-01, chatgpt-ux | Terms and Privacy on signup are not links. Password rules disappear. No 18+ checkbox. |
| A19 | P1 | ios | ux, ux-claude quit 3 | "Continue Without Monitoring" is the same gold button as the real action. |

## P2

| ID | Priority | Owner | Reviewers | Issue |
|---|---|---|---|---|
| A20 | P2 | backend | grok 8, claude (email escape claim was wrong), master M13 | Invite and reset HTML is not escaped. Invite routes are not rate-limited. A missing SendGrid key logs the reset link and returns success. |
| A21 | P2 | backend, website | claude F-18, master M14, grok 9 (privacy sentence) | Password reset does not revoke JWTs. The privacy page says the JWT is not in UserDefaults. The app writes it to the app group. |
| A22 | P2 | backend | claude F-11, claude F-24, master M15, ux U-19 | PIN failures live in memory. `VerifyPIN` uses `LIMIT 1` and returns the same `success: false` for no PIN and a wrong PIN. |
| A23 | P2 | backend | claude F-23, master M16 | The partner being asked to watch cannot delete the relationship. The privacy page says either person can end it. |
| A24 | P2 | backend | grok 17, claude F-30, master M17 | "App was removed" is pushed at the device token Apple just rejected. |
| A25 | P2 | backend | claude F-31, master M18 | Group streak and `flags_last_30` count `clean` events. |
| A26 | P2 | website, backend | grok 16, master M19, chatgpt 2, ux | Privacy copy: "never disclosed", 30-day deletion vs immediate delete, donation rows cascade, Do Not Track, TLS 1.3, AES-256, end-to-end. About page says about $4 per person and also about $80 a month. Logs include protection detail and recipient addresses. |
| A27 | P2 | backend, website | grok 15, grok 18, claude F-26, claude F-27, claude F-28, claude F-20, claude F-35, master M20 | No partner cap. Client timestamps are unbounded. No retention. `/api/donate` is unthrottled. Reset mail is 10/min/IP only. An informational alert can suppress a severe push. Groups are capped at 12 but the check can race. |
| A28 | P2 | ios | grok 5, master M21, master missed 4 | Sign-out leaves the previous person's Screen Time selection, shield flags, and queued events on the phone. |
| A29 | P2 | backend, website | grok 10, master M22 | `next@14.2.5` and `golang-jwt/jwt/v5 v5.2.1` (GO-2025-3553). API image is `golang:1.22-alpine` on `alpine:3.19`. |
| A30 | P2 | backend | claude F-08, master (reranked from P1) | Login's dummy bcrypt string is not a real hash, so the timing pad does not run. Registration already returns 409 for a known email. |
| A31 | P2 | backend | claude F-09, master | `POST /relationships` returns "no account found for that email". |
| A32 | P2 | website | ux, ux-claude W-03 | Waitlist inputs use placeholders as their only labels. |
| A33 | P2 | ios | ux U-05, U-06, U-07, U-10, U-13, U-14, U-17, U-18, U-20, U-21, U-23, ux-claude same ids, chatgpt-ux RF-11, RF-12, RF-15, RF-16, RF-17, RF-18 | Plain language, empty states, sign-out confirmation, invite error text, streak day letters, pre-signup red Xs (server can add a tri-state), lockout jargon, tour skip, feedback truncation, contrast, 44 pt targets, VoiceOver. Pre-signup week values are a backend field the app must render. |
| A34 | P2 | ios | ux U-22, ux-claude U-22 | `support@remainfaithful.com` may not be the mailbox that actually receives mail. Suspected. |
| A35 | P2 | website | grok 19, master P3 canonicals, ux-claude W-04 (called P3), ux-claude W-05 (called P3) | Included in this push because the brief asks for them: form errors use the same gold as buttons, several pages lack canonical URLs, contact labels are not tied to fields. |

## P3 (not in this push unless noted)

| ID | Owner | Reviewers | Issue |
|---|---|---|---|
| A36 | website, ci | claude F-21, master P3 | ESLint is not configured. CI runs Go tests only. |
| A37 | backend | claude F-32, F-33, F-34 | History is the latest 100 rows. Missing `(user_id, timestamp)` and heartbeat indexes. Indexes are pulled into A27 because the sweep and retention need them. Pagination stays P3. |
| A38 | backend | claude F-36, master P3 | Reset tokens stored in plaintext. Pulled into A04/A20 for this push (store a hash). |
| A39 | backend | claude F-39, F-40 | Any member can rename the group and rewrite the covenant. Sole admin can leave with no replacement. |
| A40 | backend | claude F-37, F-38, F-43, F-44 | `markTokenInactive` is global by token string. Per-push logging. SQL column name is interpolated from a fixed caller. `X-Forwarded-For` fallback. The production IP trust is done with A27. |
| A41 | website | grok 19, chatgpt 5, master P3 | Waitlist returns success and stores nothing. |
| A42 | ios | grok 19, master P3, ux, chatgpt-ux | No in-app 18+ attestation. Policy already says 18+. |
| A43 | ios | chatgpt-ux RF-04 through RF-10, RF-20, RF-21, RF-22 | Partner expectations, group empty guidance, alert next step, Deep Scan scope, lockout wording, tour replay, feedback categories. Overlaps A05, A14, A15, A33. |
| A44 | ios | master missed 5 | Sign-out does not clear the next account's inherited shield. Same as A28. |

## Withdrawn or stale

| Claim | Reviewers | Ruling |
|---|---|---|
| Privacy page ends with "Fort Wayne" | chatgpt 1 | Stale. Not in the privacy page. |
| Family Controls distribution entitlement missing | chatgpt 4 | Stale as a code blocker. The entitlement is in the app, shield, and Device Activity targets. |
| 96 heartbeats per user per day | chatgpt 10 | Wrong. Heartbeats run only during Deep Scan. |
| Every email template HTML-escapes its values | claude "what is solid" | Wrong. Only the contact email escapes. See A20. |
| Compiled binary in the tree | claude F-45 | Withdrawn by that review. |
| Apple takeover is unconfirmed | none | Confirmed by grok, claude, and master. |

## Checked, and still not a finding

- No city, state, named church, or street address on the site. Jeff's name stays.
- No donation UI in the iOS target. Website Stripe donations stay. Do not put donation links on surfaces the app shows.
- Deep Scan uploads category, severity, summary, and timestamp. Frames stay on the device. DRM limits on the site should stay.
- `heartbeat_silence` cannot be forged through the protection HTTP handler.
- Password hashing is bcrypt. App JWTs are HS256 with issuer and expiry checks. That does not fix A01.
- Force-quitting the app does not clear a ManagedSettings shield.
