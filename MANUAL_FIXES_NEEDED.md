# Manual Fixes Needed

Items that require human action, Apple entitlements, credentials, or architectural decisions not automatable by this test loop.

---

## MF-01 — Heartbeat Silence Detection Not Implemented

**Severity:** Medium  
**Related test:** BE-14

**What's missing:** There is no logic anywhere in the codebase (iOS or backend) that detects when a user's device stops sending heartbeats and triggers a `heartbeat_silence` alert to their partners.

**What was fixed automatically:** The `protectionAlertBody` switch in `backend/internal/handler/protection.go` now has an explicit `case "heartbeat_silence":` so the notification message is correct if/when this alert type is ever sent.

**What still needs building:**

Choose one of these approaches:

**Option A — Backend cron job (recommended):**
Add a scheduled job (e.g., every 5 minutes) that queries for users whose `last_heartbeat_at` timestamp is older than a configurable threshold (e.g., 30 minutes) and calls `sendProtectionAlertToPartners` with `alert_type = "heartbeat_silence"`. Requires adding a `last_heartbeat_at` column to the `users` table (or a separate `heartbeats` table) and a cron trigger on Fly.io.

**Option B — iOS-side detection:**
Have the iOS app detect that it hasn't sent a heartbeat in X minutes (e.g., because the app was force-quit or the extension was killed) and send a `heartbeat_silence` alert on next foreground. This is harder because a killed app cannot alert.

**Option C — Partner-side polling:**
The partner app periodically polls for the last known heartbeat timestamp and shows a local in-app warning (no backend push). Less reliable for async accountability.

No credentials or entitlements are blocked here — this is a feature build.
