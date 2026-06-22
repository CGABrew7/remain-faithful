# Automated Fixes

Applied: 2026-06-21

## Fix 1 — BE-14: Add `heartbeat_silence` case to `protectionAlertBody`

**File:** `backend/internal/handler/protection.go`  
**Line:** after `case "family_controls_revoked":` (around L169)

**Before:**
```go
case "family_controls_revoked":
    return name + " revoked Screen Time authorization"
default:
```

**After:**
```go
case "family_controls_revoked":
    return name + " revoked Screen Time authorization"
case "heartbeat_silence":
    return name + "'s device has stopped sending heartbeats"
default:
```

**Why:** The `protectionAlertBody` switch had no explicit case for the `heartbeat_silence` alert type. If the type is ever sent, the default branch would produce a generic "changed a protection setting" message, which is misleading for a heartbeat silence event.

**Scope of fix:** Wires the notification message only. See `MANUAL_FIXES_NEEDED.md` for the missing detection logic that would actually trigger this alert.
