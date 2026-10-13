# Remain Faithful — Complete Review Feedback (Claude)

**Reviewer:** Claude (Opus 5) · **Date:** 2026-10-13 · **Commit:** `784998f`
**Scope:** iOS app (SwiftUI), Go backend (Fly.io), Next.js website (Vercel)
**Method:** Read-only source review plus test execution. No product code was modified.

Consolidates both review passes into one list:
- **Round 1 (F-01 … F-45)** — security, privacy, monitoring integrity, cost at scale, App Store risk → `docs/reviews/round1-claude.md`
- **Round 2 (U-01 … U-23, W-01 … W-05)** — UX, UI, friction, accessibility → `docs/reviews/round2-ux-claude.md`

**Totals:** 72 open findings + 1 withdrawn. 8 P0 · 13 P1 · 25 P2 · 27 P3.

---

## Verification run

| Check | Command | Result |
|---|---|---|
| Go vet | `go vet ./...` | **Pass** — exit 0, no diagnostics |
| Go tests | `go test ./...` | **Pass** — all packages (apns, auth have no test files) |
| Website build | `npm run build` | **Pass** — 35 routes |
| Website lint | `npm run lint` | **Not configured** — interactive ESLint prompt (F-21) |
| Website tests | `npm test` | **Pass** — 6/6 |
| Secret scan | regex sweep, Stripe/SendGrid/AWS/GitHub/PEM | **Clean** — only `.env.example` placeholders |

Two findings were confirmed by running isolated reproductions: **F-01** (forged Apple token accepted) and **F-08** (bcrypt timing, 43.9 ms vs 101 ns).

---

## Severity key

| | Meaning |
|---|---|
| **P0** | Breaks the core safety promise, or allows account takeover / unconsented disclosure of the most sensitive data |
| **P1** | Serious security, privacy, or silent-monitoring-failure defect |
| **P2** | Real defect with material user, cost, or trust impact |
| **P3** | Hygiene, polish, minor cost |

---

# Master list — all findings

## P0 (8)

| ID | Finding | Location |
|---|---|---|
| **F-01** | Apple Sign In verifies no JWT signature and links accounts by email → unauthenticated takeover of any account from an email address alone | `social_auth.go:153-181`, `:119-127` |
| **F-02** | All 7 protection-alert call sites are `try?` fire-and-forget with no retry queue → airplane mode defeats monitoring with no trace | `SettingsView.swift:310,853`, `AppLockoutManager.swift:95`, +4 |
| **F-03** | Protection and heartbeat-silence alerts are never persisted server-side; `ListAlerts` inner-joins `events` so they can't appear | `protection.go:142-193`, `alerts.go:85-86`, `main.go:426` |
| **F-04** | Partner invite to an existing user never leaves `pending` — the accept token is generated but never stored — yet the API reports `"connected"` | `invites.go:62-81`, `events.go:139` |
| **F-05** | Any user can add any other user to a group by email with no consent, immediately receiving their flagged-content alerts | `invites.go:283-305`, `events.go:141-145` |
| **F-06** | Password reset is unusable: link targets the dead `app.remainfaithful.app`, and no `/reset-password` route exists on any domain | `auth.go:327-328,395-399`; website route absent |
| **U-03** | iOS models 4 event categories; backend sends 6, only 2 overlap. `app_usage`, `self_harm`, `violence`, `clean` are silently `compactMap`-dropped | `DashboardView.swift:5-9,28,372` |
| **U-06** | Sign Out stops all monitoring in 2 taps — no confirmation, no PIN gate, no partner notification | `SettingsView.swift:403-407`, `AuthState.swift:87` |

## P1 (13)

| ID | Finding | Location |
|---|---|---|
| **F-07** | Google Sign-In never validates `aud` or `email_verified` → any Google ID token from any OAuth client is accepted | `social_auth.go:185-209` |
| **F-08** | Login enumeration: the "constant-time" path uses an invalid bcrypt hash that errors before any KDF work (43.9 ms vs 101 ns) | `auth.go:239-243` |
| **F-09** | `CreateRelationship` states outright whether an email has an account, then leaks that person's name without consent | `relationships.go:41-48`, `:88-101` |
| **F-10** | No email verification anywhere; registering with a known-invited address auto-redeems that person's invites | `auth.go:71,84-139` |
| **F-11** | PIN brute force sends no partner alert (it's client-side only) and the limiter is in-memory, resetting on every deploy | `protection.go:259-297`, `PartnerPINManager.swift:50-68` |
| **F-12** | `min_machines_running = 0` lets Fly stop the machine hosting the in-process tamper-detection sweep | `fly.toml:14-16`, `main.go:71-73` |
| **F-13** | `refreshTokenIfNeeded` writes only to Keychain, not App Group → extensions 401 silently after 24 h; heartbeat ignores the status code | `APIClient.swift:186-193`, `SampleHandler.swift:710-713` |
| **F-14** | Heartbeat silence alerts fire once and never re-arm for a user who stays dark | `heartbeat_sweep.go:59-68,96-101` |
| **F-15** | Users who never send a first heartbeat (`IS NOT NULL`) are invisible to both watchdogs | `heartbeat_sweep.go:63`, `heartbeat.go:68` |
| **F-16** | Unvalidated, unbounded client `detail` string is logged; today it carries real user names | `protection.go:107-110,166-168` |
| **F-17** | Recipient emails logged on every send; with no SendGrid key, full bodies incl. live reset tokens print to stdout and return `nil` | `email.go:128-131,158,170` |
| **F-18** | JWTs are irrevocable and infinitely refreshable; password reset does not invalidate sessions | `jwt.go:20-33`, `auth.go:266-281,339-393` |
| **U-15** | Tamper-alert pushes route to an empty dashboard — nothing stored, nothing shown | `NotificationService.swift:134-139` |

## P2 (25)

| ID | Finding | Location |
|---|---|---|
| **F-19** | JWT duplicated into App Group UserDefaults plaintext, captured in unencrypted device backups | `AuthState.swift:48-49` |
| **F-20** | Push throttle is keyed on recipient across all monitored users and ignores severity → drops `severe` self-harm alerts | `events.go:192-202` |
| **F-21** | `npm run lint` unconfigured (hangs on interactive prompt); CI runs no website checks and no secret scanning | `website/package.json:8`, `.github/workflows/` |
| **F-22** | Group member cap check is TOCTOU on the direct-add path (the registration path locks correctly) | `invites.go:254-260,288-293` |
| **F-23** | `DeleteRelationship` is one-directional → partner-side rows show a delete button that always 404s | `relationships.go:200-237` |
| **F-24** | `VerifyPIN` uses `LIMIT 1` with no `ORDER BY` → arbitrary PIN wins when two partners each set one | `protection.go:284-288` |
| **F-25** | "Data at rest is AES-256 encrypted" has no backing code; "secret scanning on every contribution" is opt-in per clone | `privacy-architecture/page.tsx:266,297` |
| **F-26** | `/api/donate`, `/api/waitlist`, `/api/contact` have no rate limiting; each donate call creates a real Stripe session | `api/donate/route.ts:4-45` |
| **F-27** | No cap on partner count (groups cap at 12); `POST /relationships` has no rate limiter → event-fanout amplification | `invites.go:27-108`, `events.go:162-172` |
| **F-28** | `/auth/forgot-password` has no per-account cap → 14,400 emails/day to one victim from one IP | `auth.go:286-334` |
| **F-29** | `UpdateMe` accepts any email with no format check, no uniqueness mapping (500 not 409), no re-verification | `users.go:47-87` |
| **F-30** | `APP_DELETED` is sent to the very token APNs just rejected with 410 — the feature can never work | `notifications.go:343-368` |
| **F-31** | Group streak counts `clean` events; personal stats excludes them → same user shows a streak on one screen, 0 on the other | `groups.go:147-156` vs `users.go:273-278` |
| **U-02** | Disabled submit button never says why; the 8-char rule lives only in a placeholder that vanishes on typing. Sign-in enforces it too | `OnboardingView.swift:192-197`, `LoginView.swift:30-32` |
| **U-05** | The only Deep Scan instruction in onboarding is 11 pt at 30 % opacity (2.64:1), and assumes Screen Recording is already in Control Center | `OnboardingView.swift:500-505` |
| **U-07** | Invite errors are replaced with "check the email", discarding the server's accurate message | `OnboardingView.swift:424-429` |
| **U-10** | Group invite collapses `added` / `already_member` / `invited` into one "sent" state | `GroupView.swift:1327-1330` |
| **U-13** | Streak grid hardcodes Sun–Sat labels against a rolling 7-day window → letters are wrong unless today is Saturday | `DashboardView.swift:550`, `users.go:255-256` |
| **U-14** | Days before account creation render as red X marks → a new user's first screen shows six failures | `users.go:319`, `DashboardView.swift:586-589` |
| **U-17** | "App Lockout (Deep Scan)" / "Off — optional hardening" — stacked jargon | `SettingsView.swift:256-262` |
| **U-18** | Six names for two concepts; Settings sections MONITORING / SCREEN TIME / PROTECTION overlap | throughout |
| **U-19** | Wrong PIN and no-PIN-configured are indistinguishable to the user | `protection.go:289-293` |
| **U-20** | One tap on tour Skip dismisses it permanently | `ContentView.swift:55` |
| **U-21** | Feedback silently truncates at 120 / 4000 chars with no counter | `FeedbackSheet.swift:196-197,215-216` |
| **U-22** | *Suspected* — `support@remainfaithful.com` may not be a real mailbox; backend routes to `jeff@hanokventures.co`, sends from a third domain | `FeedbackSheet.swift:243,272` |

## P3 (26)

| ID | Finding | Location |
|---|---|---|
| **F-32** | `ListEvents`/`ListAlerts` hard-code `LIMIT 100`, no pagination | `events.go:289`, `alerts.go:89` |
| **F-33** | No index for `ORDER BY timestamp DESC` per user | `main.go:432` |
| **F-34** | No index on `users.last_heartbeat_at`, filtered every 5 min | `main.go:533` |
| **F-35** | Client `timestamp` unvalidated; invalid strings 500 instead of 400 | `events.go:107-113` |
| **F-36** | Password-reset tokens stored in plaintext | `main.go:451-457` |
| **F-37** | `markTokenInactive` not scoped by `user_id`; device tokens accepted unverified | `notifications.go:413-415,42-50` |
| **F-38** | Per-push stdout logging on every APNs send and SendGrid call | `apns.go:190,202,205` |
| **F-39** | Any group member can rewrite the covenant; role is read then discarded (`_ = role`) | `groups.go:232-240` |
| **F-40** | `LeaveGroup` doesn't promote a replacement admin (unlike `DeleteMe`) | `groups.go:430-469` |
| **F-41** | 7 pages lack canonical URLs | `website/src/app/*/page.tsx` |
| **F-42** | `verifyGoogleIDToken` uses `http.DefaultClient` (no timeout); token unescaped in URL | `social_auth.go:186-188` |
| **F-43** | SQL column names via `fmt.Sprintf` — not exploitable today, a trap later | `social_auth.go:112,125,136` |
| **F-44** | `clientIP` falls back to spoofable `X-Forwarded-For` | `main.go:196-208` |
| **U-01** | Terms and Privacy are plain text, not links, at 2.64:1 contrast | `OnboardingView.swift:275-280` |
| **U-04** | "Two layers" contradicts the website's five | `OnboardingView.swift:457` |
| **U-08** | "Invitation Sent!" is secretly the Continue button | `OnboardingView.swift:389-394` |
| **U-11** | "No Group Yet" advertises a join-by-link path that doesn't exist | `GroupView.swift:244-260` |
| **U-12** | Raw `localizedDescription` strings reach users | `GroupView.swift:1333` |
| **U-16** | Foreground alert banner auto-dismisses in 6 s with no history | `DashboardView.swift:328-334` |
| **U-23** | Tour calls the card "Deep Scan status" — the tour is right, the card overclaims | `AppTour.swift:40-41` |
| **W-03** | Waitlist inputs use placeholder as the only label (WCAG 3.3.2) | `WaitlistForm.tsx:57-61,80-84,110-122` |
| **W-04** | Errors rendered in brand gold (`text-wax`), indistinguishable from accents | `WaitlistForm.tsx:72,93` |
| **A-1** | Dynamic Type effectively unsupported — 361 fixed fonts vs 13 scaling styles | app-wide |
| **A-2** | Low-opacity text fails WCAG AA (2.46–3.63:1) | palette-wide |
| **A-3** | Sheet close buttons ~30 pt, below the 44 pt minimum, on 8 sheets | `ScreenTimeMonitoringView.swift:41-46` +7 |
| **A-4** | VoiceOver: 13 labels, 2 hints, 0 values across 9,868 lines | app-wide |

### Also P0/P1, listed above under their round-2 IDs

| ID | Finding | Location |
|---|---|---|
| **U-09** | Until F-04 ships, the app confidently reports a partnership that delivers nothing | `invites.go:78` |
| **W-01** | **No App Store link anywhere on the website** — hero CTA scrolls to a "Get started" section with no download | `page.tsx:104,139-168` |
| **W-02** | `/accept-invite` instructions cannot work for existing users; the token is deliberately unused | `InviteLanding.tsx:26-44`, `invite/page.tsx:12-13` |

### Withdrawn

| ID | Note |
|---|---|
| **F-45** | Flagged the 12 MB `backend/server` binary, then confirmed it's already covered by `.gitignore:18`. No action needed. |

---

# Cross-cutting

## Top 5 places a new user quits

1. **The dashboard says "Monitoring Paused" when monitoring is on.** `StatusCard` is driven solely by `isBroadcasting` — the optional Deep Scan layer (`DashboardView.swift:454-456`). Wrong in both directions. *Fix: composite state; Deep Scan becomes a subtitle, never the headline.*
2. **Nothing ever tells him to pick apps.** `ScreenTimeMonitoringView` has exactly one call site (`SettingsView.swift:94`), never linked from onboarding. Both toggles default `false`; zero apps selected. *Fix: app-picker step in onboarding + persistent "No apps selected" banner.*
3. **"Continue Without Monitoring" is the same gold button as "Enable."** `RFButton` has no secondary style. *Fix: secondary style + name the consequence in the label.*
4. **The invite link tells him to search the App Store by hand.** No `apps.apple.com` link exists anywhere on the site. *Fix: add the link — hero, `#download`, `InviteLanding`.*
5. **Forgot Password shows a green checkmark for a dead flow** (F-06). *Fix: ship `/reset-password`, set `APP_BASE_URL`.*

**Tap counts:** 20 taps before anything is monitored; 2 taps to turn it all off.

## Accessibility, measured

```
fixed-point fonts  .font(.system(size: N))   361
Dynamic Type styles .font(.body/.headline/…)  13
scaled custom fonts .font(..., relativeTo:)    2
accessibilityLabel 13 · Hint 2 · Value 0 · Hidden 3
```

Contrast computed from `rfNavy (18,33,76)`, `rfNavyMid (28,48,102)`, `rfGold (209,171,76)`:

| Element | Opacity | vs rfNavy | vs rfNavyMid | Verdict |
|---|---:|---:|---:|---|
| Terms / Privacy line | 0.30 | 2.64 | 2.50 | **FAIL** |
| Deep Scan hint | 0.30 | 2.64 | 2.50 | **FAIL** |
| Disabled button label | 0.35 | 2.46 | — | **FAIL** |
| Field placeholders | 0.38 | 3.41 | 3.15 | large text only |
| "Skip for now" | 0.40 | 3.63 | 3.33 | large text only |
| tierRow detail | 0.50 | 4.88 | 4.34 | borderline on gradient |
| Subtitles | 0.55 | 5.61 | 4.93 | PASS |
| **Gold button text** | — | **7.13** | — | **PASS** |

The brand pairing is strong. The failures come from stacking low opacity on small type. Floor body text at 0.60, secondary at 0.55, and test against `rfNavyMid` (the lighter gradient end), not `rfNavy`.

## Cost at scale

Dominant driver: the **2-minute heartbeat** (`SampleHandler.swift:687-693`) performing `UPDATE users SET last_heartbeat_at = NOW()` (`heartbeat.go:35-39`) — one dead tuple per heartbeat on the hottest table in the schema (~660k/day at 10k users).

| Users | Monthly | Notes |
|---|---|---|
| **100** | **$5–30** | ~5k req/day. Fly app ~$2 + Postgres ~$3; SendGrid free tier |
| **1,000** | **$50–90** | ~50k req/day. 512 MB app, SendGrid Essentials $20, Vercel Pro $20 |
| **10,000** | **$150–350** | ~1M req/day. Postgres dominates at $60–150 |

Assumptions: 50 % daily-active; active users average 2 h/day Deep Scan (~60 heartbeats), 10 % near-continuous (~720); ~2 flagged events and ~20 foreground calls per active user per day.

**Three changes that most reduce the 10k bill:** move heartbeats to a narrow `user_heartbeats` table (or Redis) and raise the cadence 2 → 5 min; add retention pruning (nothing ever deletes `events`, `alerts`, or `donation_events`); batch alert inserts and APNs sends in `CreateEvent`.

**Uncapped exposures:** F-26 (donate/waitlist/contact), F-27 (partner count), F-28 (reset emails).

## App Store risk

| Risk | Assessment |
|---|---|
| **3.1.1 in-app purchase** | **Clear.** No donation code in any Swift source. Keep `/app-guide` (linked from `SettingsView.swift:1674`) free of donation CTAs — that's the realistic exposure. |
| **5.1.1 / 5.1.2 privacy** | The literal promise holds in code. Exposure is F-16 / F-17 — PII in server logs — which contradicts the posture claimed on `/privacy-architecture`. |
| **5.1.1(v) account deletion** | **Satisfied.** `DELETE /users/me` hard-deletes with cascade, reachable from Settings. |
| **Family Controls entitlement** | **Strongest latent risk.** F-05 makes the product function as covert third-party surveillance. Fix on that basis alone. |
| **Sign in with Apple button** | Custom-drawn white button styled to match Google rather than `SignInWithAppleButton` (`LoginView.swift:~219`). HIG deviation worth checking. |
| **Time-sensitive interruption level** | *Suspected* — `SendPanicAlert` sets `interruption-level: time-sensitive` (`notifications.go:171`); verify the entitlement exists or panic pushes fail silently. |

---

# Recommended order of work

### Security and integrity first

1. **F-01** — Apple token signature verification + remove email-fallback linking
2. **F-05** — Group-invite consent (also the entitlement risk)
3. **F-06** — Ship `/reset-password`, set `APP_BASE_URL`
4. **F-04** — Persist the partner invite token
5. **F-02 + F-03** — Durable protection alerts: client queue + server rows
6. **F-07** — Google `aud` validation
7. **F-12** — `min_machines_running = 1` *(one line)*
8. **F-13** — Write refreshed token to both stores *(one line)*
9. **F-16 + F-17** — Stop logging `detail`, emails, and bodies
10. **F-14 + F-15** — Re-alert on continued silence; detect never-started users

Items 7 and 8 are one-line changes with outsized impact — ship them in the first patch regardless of ordering.

### UX, by impact ÷ effort

**Trivial effort, large payoff**

| # | Fix | ID |
|---|---|---|
| 1 | Add the App Store link (hero, `#download`, `InviteLanding`) | W-01 |
| 2 | Make "Continue Without Monitoring" visually secondary; name the consequence | quit #3 |
| 3 | Make Terms / Privacy real links; raise contrast to 0.6 | U-01 |
| 4 | Generate streak day letters from today | U-13 |
| 5 | Render pre-signup days as dimmed circles, not red X marks | U-14 |
| 6 | `.frame(width: 44, height: 44)` on the shared close button | A-3 |
| 7 | Surface the server's actual invite error text | U-07 |
| 8 | Raise the Deep Scan instruction to 13 pt / 0.6 opacity | U-05 |

**Low effort, high impact**

| # | Fix | ID |
|---|---|---|
| 9 | **Rewrite `StatusCard` to a composite monitoring state** | quit #1 |
| 10 | **Add the four missing `EventCategory` cases and a `default`** | U-03 |
| 11 | **PIN-gate Sign Out, confirm it, notify the partner** | U-06 |
| 12 | Persistent "No apps selected" banner → picker | quit #2 |
| 13 | Password rules as live helper text; drop the length gate on sign-in | U-02 |
| 14 | Decode and display the three group-invite outcomes | U-10 |
| 15 | Labels on waitlist inputs | W-03 |
| 16 | Forgot Password detent `.medium`; fixed heights → `minHeight` | A-1 |

**Medium effort**

| # | Fix | ID |
|---|---|---|
| 17 | Add an app-picker step to onboarding | quit #2 |
| 18 | Settle vocabulary on **Alerts** and **Blocking** | U-18 |
| 19 | Protection-alert history in the alerts list (needs F-03) | U-15 |
| 20 | VoiceOver labels for StatusCard, StreakCard, severity | A-4 |
| 21 | Ship `/reset-password`; fix the success copy | F-06 |
| 22 | Migrate text to Dynamic Type, onboarding first | A-1 |

---

# What already works well

**Backend and security**
- **The privacy promise is real in code.** `uploadEvent` transmits only category, severity, summary, timestamp — and summary is constrained server-side to a 12-entry allowlist plus one anchored regex (`events.go:22-41`). `detail` is deliberately excluded from APNs bodies (`protection.go:198`).
- **Event summary validation** is well designed — a modified client cannot push arbitrary text to a partner's lock screen.
- **JWT verification** is correct: HMAC pinned, issuer checked, `WithExpirationRequired()`, refuses to start without `JWT_SECRET` (`jwt.go:37-65`). The `alg=none` and algorithm-confusion classes are closed.
- **SQL** is parameterized throughout; no injection found.
- **Stripe webhooks** verify signatures and are idempotent via a `UNIQUE` event-ID log (`donations.go:23-28`).
- **Email templates** HTML-escape all interpolated values and strip CRLF from headers (`email.go:89-109`).
- **`checkoutReturnUrls`** correctly allowlists origins — no open redirect, and it's tested.
- **No secrets in the repo**; `.githooks/pre-commit` is a well-constructed scanner.
- **Group member cap** is correctly enforced under `SELECT … FOR UPDATE` on the registration path (`auth.go:181-217`) — the model to copy for F-22.

**Product and UX**
- **The app tour** is the best-written surface in the product: short, concrete, plain-spoken, correctly skips the invite step when there's no group, and is the only view handling Dynamic Type and Reduce Motion properly (`AppTour.swift:118`).
- **Delete Account** applies friction exactly where it belongs (`SettingsView.swift:564`).
- **Conversation starters** in `AlertDetailView` are thoughtful and non-shaming — the clearest expression of the product's philosophy in the codebase.
- **`PauseNotice`** tells the truth about whether pausing notifies partners, conditioned on real lockout state (`AppTour.swift:70-77`).
- **The `StatusCard` comment** records why a previous local-only toggle was removed for contradicting real state — right instinct, wrong source of truth.
- **Website accessibility fundamentals** are ahead of the app's: real skip link (`layout.tsx:106`), global `:focus-visible` gold outline (`globals.css:98-100`), proper focus replacement on inputs rather than a bare `outline: none` (`:157-161`), exactly one `<h1>` per page across all 24 pages, `min-h-11` touch targets, `aria-hidden` on ornaments.
- **Privacy copy** on `/privacy-architecture` and `/how-it-works` is precise and matches the code.

---

## The one-paragraph summary

The architecture is sound and the privacy promise is genuinely kept in code — that part is better than most apps in this category. What's broken is the loop between two people. A new user finishes setup with nothing monitored and a dashboard that says "Paused"; his partner clicks an invite link with no App Store link behind it; if the partner already has an account the connection never completes; the always-on layer's events are dropped by the client; and every tamper signal is an ephemeral push with no durable record on either end. Each failure is invisible to the person experiencing it. That invisibility — not any single bug — is the core finding of both reviews.
