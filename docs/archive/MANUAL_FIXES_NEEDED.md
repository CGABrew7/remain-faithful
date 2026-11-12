# Manual Fixes Needed

Updated 2026-09-22. Older “MF-01 not implemented” text is obsolete.

---

## MF-01 — Heartbeat silence detection — IMPLEMENTED

**Status:** Shipped in source.

- `backend/internal/handler/heartbeat.go` writes `users.last_heartbeat_at` on `POST /heartbeat`.
- `backend/internal/handler/heartbeat_sweep.go` runs every 5 minutes from `main.go` via `StartHeartbeatSweep`.
- Threshold: `HEARTBEAT_SILENCE_MINUTES` (default 30).
- Dedup: `relationships.last_silence_alert_at` so each partner is alerted once per silence window.
- `heartbeat_silence` is a **server-only** alert type. HTTP clients cannot forge it.

**Still Jeff / ops:**

1. Confirm Fly.io secret `HEARTBEAT_SILENCE_MINUTES` is set (or accept default 30).
2. Confirm production Postgres actually has `users.last_heartbeat_at` and `relationships.last_silence_alert_at` (added by `migrate()` in `main.go`, not by the stale `001_schema.sql` file).
3. On a real device, stop heartbeats for 30+ minutes and verify the partner receives “[Name]’s device has stopped sending heartbeats.”

`001_schema.sql` is documentation-only and lags `migrate()`. Do not treat that file as the live schema.

---

## MF-02 — Family Controls distribution entitlement

**Severity:** Blocker for App Store / external TestFlight.

Apple must approve Family Controls for Team ID production distribution. Development entitlement is not enough. Jeff requests this in the Apple Developer portal if the App ID still shows Development only.

---

## MF-03 — App Store Connect upload

**Severity:** Blocker for “submitted by EOW.”

Requires Jeff’s Mac, paid Apple Developer account, certificates, and App Store Connect. Agents cannot submit the binary. Listing copy and review notes: `docs/APP_STORE_SUBMISSION.md`.

---

## MF-04 — Woodfield Stripe legal name + EIN on receipts

Payment Links are live. Confirm Checkout shows Woodfield Foundation (not a personal name) and that tax-deductible receipts include the 501(c)(3) legal name and EIN.

---

## MF-05 — Waitlist persistence

Website `/api/waitlist` currently logs and returns success without storing the email. Wire HubSpot or email before treating the homepage form as a launch channel.
