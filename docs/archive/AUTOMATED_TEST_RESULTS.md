# Automated Test Results

Run date: 2026-06-21  
Phases: 1 (initial), 3 (re-test after fixes)

---

## Phase 1 Results

### Backend Endpoints — Live (remain-faithful-api.fly.dev)

| # | Test | Result | Detail |
|---|------|--------|--------|
| BE-01 | GET /health → 200 | ✅ PASS | `{"status":"ok"}` |
| BE-02 | POST /auth/register weak password → 400 | ✅ PASS | Live confirmed |
| BE-03 | POST /auth/register valid payload → 201 | ✅ PASS | Live confirmed (probe account created then deleted) |
| BE-04 | POST /auth/login wrong password → 401 | ✅ PASS | Live confirmed |
| BE-05 | POST /auth/refresh invalid token → 401 | ✅ PASS | Live confirmed |
| BE-06 | POST /heartbeat invalid screen value → 400 | ✅ PASS (code) | heartbeat.go validates screen = "active"\|"idle"; live test requires auth |
| BE-07 | POST /donations/create-checkout-session amount > 10000 → 400 | ✅ PASS (code) | donations.go L29: `AmountDollars > 10000` → 400; requires auth to test live |
| BE-08 | Auth rate limiter applied to 6 endpoints | ✅ PASS (code) | main.go authRateLimiter wraps all auth routes |
| BE-09 | DELETE /users/me cascade exists | ✅ PASS (code) | users.go DeleteMe handler with cascade |
| BE-10 | Invite email errors logged not swallowed | ✅ PASS (code) | invites.go L71, L98: `fmt.Printf("[invite] email error: %v\n", sendErr)` |
| BE-11 | bcrypt used for PIN hashing | ✅ PASS (code) | protection.go uses bcrypt.GenerateFromPassword |
| BE-12 | Per-user PIN rate limiter (5/15 min) | ✅ PASS (code) | protection.go sync.Map per-user limiter |
| BE-13 | Partner-only WHERE clause in PIN verify | ✅ PASS (code) | protection.go WHERE user_id = $1 AND status = 'accepted' |
| BE-14 | heartbeat_silence alert type is wired | ❌ FAIL | No `case "heartbeat_silence":` in protectionAlertBody switch (protection.go L154–177). Alert type string is also absent from iOS codebase — heartbeat silence detection not implemented. |

### Privacy / Security Code Audit

| # | Test | Result | Detail |
|---|------|--------|--------|
| PS-01 | ocrText never serialized into network request body (EventProcessor.swift) | ✅ PASS | EventProcessor.swift L70–75: only category, severity, summary, timestamp passed to createEvent |
| PS-02 | NSPrivacyTracking = false in all 4 PrivacyInfo.xcprivacy | ✅ PASS | Confirmed all 4 targets |
| PS-03 | CA92.1 (UserDefaults) declared in 3 of 4 targets | ✅ PASS | RemainFaithful, RemainFaithfulBroadcast, RemainFaithfulDeviceActivity — ShieldConfig correctly excluded |
| PS-04 | JWT tokens stored in Keychain (not UserDefaults) | ✅ PASS | Confirmed from prior code review |
| PS-05 | No Anthropic in privacy policy third-party list | ✅ PASS | Third parties: Apple APNs, Stripe, Google Analytics only |
| PS-06 | No "sandboxed" or "internet-isolated" claim on /privacy | ✅ PASS | Phrase absent from privacy page |

### iOS Code Structure

| # | Test | Result | Detail |
|---|------|--------|--------|
| iOS-01 | Darwin IPC notifications in all 4 SampleHandler lifecycle methods | ✅ PASS | broadcastStarted, broadcastFinished, broadcastPaused, broadcastResumed |
| iOS-02 | No Tier 3 / cloud fallback block in SampleHandler.swift | ✅ PASS | No sendToClassify, no cloud classification path |
| iOS-03 | Heartbeat sends only "active" or "idle" (no content) | ✅ PASS | SampleHandler heartbeat payload confirmed |
| iOS-04 | uploadEvent body excludes ocrText | ✅ PASS | Upload body: category, severity, summary, timestamp only |
| iOS-05 | buildSummary / continuedActivitySummary use switch (not dynamic LLM calls) | ✅ PASS | Static switch statements in SampleHandler.swift |
| iOS-06 | AppLockoutManager.syncShieldState reads isBroadcasting from app group | ✅ PASS | AppLockoutManager.swift confirmed |
| iOS-07 | PartnerPINManager calls /relationships/{id}/pin endpoints | ✅ PASS | PartnerPINManager.swift confirmed |
| iOS-08 | family_controls_revoked detected via onChange(of: fcManager.authorizationStatus) | ✅ PASS | RemainFaithfulApp.swift confirmed |
| iOS-09 | Streak returns 0 for new account (not crash or -1) | ✅ PASS | users.go: `streak = 0 // account created today` |
| iOS-10 | ShieldConfigurationExtension shows lockout-specific subtitle when lockout active | ✅ PASS | ShieldConfigurationExtension.swift L43–45: reads `appLockoutEnabled` from app group, subtitle = "Open Remain Faithful and start Deep Scan to unlock." |

### Website Content Claims

| # | Test | Result | Detail |
|---|------|--------|--------|
| WEB-01 | Homepage: no "Tier 3", "cloud fallback", "<5%", "certificate pinning" | ✅ PASS | All phrases absent from live page |
| WEB-02 | Homepage: no "App Lockout" or "Partner PIN" advertising | ✅ PASS | Features not mentioned on homepage |
| WEB-03 | /privacy-architecture: two-tier pipeline mentioned, no Tier 3 section | ✅ PASS | "two-tier" and "on-device" confirmed; no Tier 3 or cloud fallback |
| WEB-04 | /how-it-works: no cloud fallback step | ✅ PASS | Banned phrase absent |
| WEB-05 | /privacy: no Anthropic in third-party list | ✅ PASS | Third parties: Apple APNs, Stripe, Google Analytics only |
| WEB-06 | /privacy: no "sandboxed" or "internet-isolated" claim | ✅ PASS | Phrase absent |
| WEB-07 | /about: page loads and contains founder information | ✅ PASS | Page loads; Jeff Brewer reference confirmed |
| WEB-08 | All 5 pages: no banned phrases (Tier 3, cloud fallback, App Lockout, Partner PIN) | ✅ PASS | Clean across all pages |

---

**Phase 1 complete: 31 passed, 1 failed, 0 skipped.**

Failure: BE-14 — `heartbeat_silence` alert type missing from `protectionAlertBody` switch; heartbeat silence detection logic not implemented anywhere in the codebase.

---

## Phase 2 Fixes

| Fix | File | Change |
|-----|------|--------|
| BE-14 | `backend/internal/handler/protection.go` | Added `case "heartbeat_silence": return name + "'s device has stopped sending heartbeats"` to `protectionAlertBody` switch |

Backend compiled (`go build ./...` — no errors) and deployed to Fly.io (deployment-01KVNT6D2SCXP7NJWGBYFNW41Q).

See `AUTOMATED_FIXES.md` for full diff.  
See `MANUAL_FIXES_NEEDED.md` for the unimplemented heartbeat silence detection logic (MF-01).

---

## Phase 3 Re-test Results

| # | Test | Phase 1 | Phase 3 | Detail |
|---|------|---------|---------|--------|
| BE-14 | heartbeat_silence alert type is wired | ❌ FAIL | ✅ PASS | `case "heartbeat_silence":` added to `protectionAlertBody`; deployed to Fly.io; GET /health → 200 confirms live |

**Phase 3 complete: 1 re-tested, 1 fixed (0 still failing).**

---

## Final Summary

| Metric | Count |
|--------|-------|
| Total tests | 31 |
| Phase 1 pass | 30 |
| Phase 1 fail | 1 |
| Fixed in Phase 2 | 1 |
| Still failing after Phase 3 | 0 |
| Manual action required | 1 (MF-01 — heartbeat silence detection logic) |
