package handler

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

// Per-device heartbeats.
//
// Shipped iOS 1.0.2 posts {"screen"} only. That path stays on
// users.last_heartbeat_at and never writes a device row. Partners of those
// phones still get "Name's device has stopped sending heartbeats".
//
// Newer clients send device_id plus device_kind. Each device has its own
// last_seen. A heartbeat for one device does not refresh another device, and
// does not refresh users.last_heartbeat_at.
//
// device_kind is an allowlist. The partner sentence uses the label
// (iPhone, Mac, Windows PC, Android phone). The device id is an opaque token
// for the owner. It is never copied into a partner alert.
//
// Next iOS build (kind "iphone") and the Android build (kind "android") should
// add Settings > My devices:
//   - GET /users/me/devices shows label, monitoring on/off, and last check-in.
//   - Mark the row whose device_id matches the id stored on that phone.
//   - POST /users/me/devices with monitoring true or false when protection changes.
//     This does not move last_seen. Only POST /heartbeat does.
//   - POST /heartbeat includes device_id, device_kind, and screen.
//   - Turning monitoring on calls POST /users/me/devices/reclaim first, then heartbeats.
//   - DELETE /users/me/devices/{id} tells partners the same way turning activity
//     monitoring off does, then the row is gone. A later heartbeat from that id
//     returns 409 code device_removed until the person turns monitoring on again.
//   - Protection posts may include the same device_id and device_kind.
//   - Store the id in app prefs. Do not reuse the push token. Do not send app
//     names, page addresses, or screen contents.
//   - iOS 1.0.2 will not appear in the list until that next build ships.

// deviceKindLabels is the only text that may follow "on their" or "Name's".
// Keep it aligned with desktop/core/device.go.
var deviceKindLabels = map[string]string{
	"iphone":  "iPhone",
	"mac":     "Mac",
	"windows": "Windows PC",
	"android": "Android phone",
}

func deviceKindLabel(kind string) (string, bool) {
	label, ok := deviceKindLabels[kind]
	return label, ok
}

func validDeviceID(id string) bool {
	if id == "" || len(id) > 64 {
		return false
	}
	for _, c := range id {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '_', c == '-':
		default:
			return false
		}
	}
	return true
}

func deviceCollapseKey(deviceID string) string {
	sum := sha256.Sum256([]byte(deviceID))
	return hex.EncodeToString(sum[:4])
}

// protectionBody is the partner sentence. An empty label keeps the legacy line.
// The label must already be an allowlisted device type.
func protectionBody(name, alertType, deviceLabel string) string {
	body := protectionAlertBody(name, alertType, "")
	if deviceLabel == "" {
		return body
	}
	return body + " on their " + deviceLabel
}

func deviceSilenceBody(name, deviceLabel string) string {
	return name + "'s " + deviceLabel + " has stopped sending heartbeats"
}

// MigrateDeviceHeartbeats adds the per-device tables. Safe to run on every startup.
func MigrateDeviceHeartbeats(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS device_heartbeats (
			user_id            BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			device_id          TEXT        NOT NULL,
			device_kind        TEXT        NOT NULL,
			monitoring         BOOLEAN     NOT NULL DEFAULT FALSE,
			last_seen_at       TIMESTAMPTZ,
			screen             TEXT        NOT NULL DEFAULT '',
			silence_alerted_at TIMESTAMPTZ,
			created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, device_id)
		);
		CREATE INDEX IF NOT EXISTS idx_device_heartbeats_seen ON device_heartbeats (last_seen_at);

		CREATE TABLE IF NOT EXISTS device_removals (
			user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			device_id   TEXT        NOT NULL,
			device_kind TEXT        NOT NULL,
			removed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, device_id)
		);
	`)
	return err
}

func deviceSilenceWindow() time.Duration {
	minutes := 30
	if raw := strings.TrimSpace(os.Getenv("HEARTBEAT_SILENCE_MINUTES")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			minutes = n
		}
	}
	return time.Duration(minutes) * time.Minute
}

func (h *H) deviceIsRemoved(ctx context.Context, userID int64, deviceID string) (bool, error) {
	var one int
	err := h.DB.QueryRowContext(ctx, `
		SELECT 1 FROM device_removals WHERE user_id = $1 AND device_id = $2
	`, userID, deviceID).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *H) upsertDeviceHeartbeat(ctx context.Context, userID int64, deviceID, kind, screen string) error {
	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO device_heartbeats (user_id, device_id, device_kind, monitoring, last_seen_at, screen)
		VALUES ($1, $2, $3, TRUE, NOW(), $4)
		ON CONFLICT (user_id, device_id) DO UPDATE
		SET device_kind = EXCLUDED.device_kind,
		    monitoring = TRUE,
		    last_seen_at = NOW(),
		    screen = EXCLUDED.screen
	`, userID, deviceID, kind, screen)
	return err
}

func (h *H) upsertDeviceReport(ctx context.Context, userID int64, deviceID, kind string, monitoring bool) error {
	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO device_heartbeats (user_id, device_id, device_kind, monitoring)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, device_id) DO UPDATE
		SET device_kind = EXCLUDED.device_kind,
		    monitoring = EXCLUDED.monitoring
	`, userID, deviceID, kind, monitoring)
	return err
}

type deviceReportRequest struct {
	DeviceID   string `json:"device_id"`
	DeviceKind string `json:"device_kind"`
	Monitoring *bool  `json:"monitoring"`
}

func decodeDeviceReport(r *http.Request, requireMonitoring bool) (deviceReportRequest, string, bool) {
	var req deviceReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, "invalid request body", false
	}
	if !validDeviceID(req.DeviceID) {
		return req, "device id is not recognized", false
	}
	if _, ok := deviceKindLabel(req.DeviceKind); !ok {
		return req, "device type is not recognized", false
	}
	if requireMonitoring && req.Monitoring == nil {
		return req, "monitoring is required", false
	}
	return req, "", true
}

// ReportDevice registers this device and whether it is currently watching.
// It does not move last_seen. POST /users/me/devices
func (h *H) ReportDevice(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	req, msg, ok := decodeDeviceReport(r, true)
	if !ok {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	removed, err := h.deviceIsRemoved(r.Context(), userID, req.DeviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save device")
		return
	}
	if removed {
		writeCodedError(w, http.StatusConflict, "device_removed", "this device was removed from the account.")
		return
	}
	if err := h.upsertDeviceReport(r.Context(), userID, req.DeviceID, req.DeviceKind, *req.Monitoring); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save device")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ReclaimDevice lets the person add a removed device back by turning monitoring on.
// POST /users/me/devices/reclaim
func (h *H) ReclaimDevice(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	req, msg, ok := decodeDeviceReport(r, false)
	if !ok {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	if _, err := h.DB.ExecContext(r.Context(), `
		DELETE FROM device_removals WHERE user_id = $1 AND device_id = $2
	`, userID, req.DeviceID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save device")
		return
	}
	if err := h.upsertDeviceReport(r.Context(), userID, req.DeviceID, req.DeviceKind, true); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save device")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type deviceListItem struct {
	DeviceID    string `json:"device_id"`
	DeviceKind  string `json:"device_kind"`
	Label       string `json:"label"`
	Monitoring  bool   `json:"monitoring"`
	LastCheckIn string `json:"last_check_in"`
}

// ListDevices returns the owner's devices. Partner clients never call this.
// GET /users/me/devices
func (h *H) ListDevices(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT device_id, device_kind, monitoring, last_seen_at
		FROM device_heartbeats
		WHERE user_id = $1
		ORDER BY created_at, device_id
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}
	defer rows.Close()

	cutoff := time.Now().Add(-deviceSilenceWindow())
	items := []deviceListItem{}
	for rows.Next() {
		var (
			id, kind   string
			monitoring bool
			seen       sql.NullTime
		)
		if err := rows.Scan(&id, &kind, &monitoring, &seen); err != nil {
			continue
		}
		label, ok := deviceKindLabel(kind)
		if !ok || !validDeviceID(id) {
			continue
		}
		item := deviceListItem{
			DeviceID:   id,
			DeviceKind: kind,
			Label:      label,
			Monitoring: monitoring,
		}
		if seen.Valid {
			item.LastCheckIn = seen.Time.UTC().Format(time.RFC3339)
			if seen.Time.Before(cutoff) {
				item.Monitoring = false
			}
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"devices": items})
}

// RemoveDevice drops one device and tells partners that activity monitoring
// was turned off on that device type. DELETE /users/me/devices/{deviceID}
func (h *H) RemoveDevice(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	deviceID := mux.Vars(r)["deviceID"]
	if !validDeviceID(deviceID) {
		writeError(w, http.StatusBadRequest, "device id is not recognized")
		return
	}

	var kind string
	err := h.DB.QueryRowContext(r.Context(), `
		SELECT device_kind FROM device_heartbeats WHERE user_id = $1 AND device_id = $2
	`, userID, deviceID).Scan(&kind)
	if errors.Is(err, sql.ErrNoRows) {
		removed, rerr := h.deviceIsRemoved(r.Context(), userID, deviceID)
		if rerr != nil {
			writeError(w, http.StatusInternalServerError, "failed to remove device")
			return
		}
		if removed {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
			return
		}
		writeError(w, http.StatusNotFound, "device not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove device")
		return
	}
	label, ok := deviceKindLabel(kind)
	if !ok {
		writeError(w, http.StatusInternalServerError, "failed to remove device")
		return
	}

	var senderName string
	if err := h.DB.QueryRowContext(r.Context(), `SELECT name FROM users WHERE id = $1`, userID).Scan(&senderName); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove device")
		return
	}

	// Notify before the row disappears. A failed send leaves the device in place
	// so a retry can still tell partners. The sentence matches monitoring_disabled.
	body := protectionBody(senderName, "monitoring_disabled", label)
	if err := h.storeAndPushProtectionText(userID, senderName, "monitoring_disabled", body, false, deviceCollapseKey(deviceID)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove device")
		return
	}
	if _, err := h.DB.ExecContext(r.Context(), `
		INSERT INTO device_removals (user_id, device_id, device_kind)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, device_id) DO NOTHING
	`, userID, deviceID, kind); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove device")
		return
	}
	if _, err := h.DB.ExecContext(r.Context(), `
		DELETE FROM device_heartbeats WHERE user_id = $1 AND device_id = $2
	`, userID, deviceID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove device")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type silentDevice struct {
	userID int64
	name   string
	id     string
	kind   string
}

func (h *H) sweepDeviceHeartbeatSilence(ctx context.Context, silenceMinutes int) {
	if silenceMinutes <= 0 {
		silenceMinutes = 30
	}
	rows, err := h.DB.QueryContext(ctx, `
		SELECT d.user_id, u.name, d.device_id, d.device_kind
		FROM device_heartbeats d
		JOIN users u ON u.id = d.user_id
		WHERE d.last_seen_at IS NOT NULL
		  AND d.last_seen_at < NOW() - INTERVAL '1 minute' * $1
		  AND (d.silence_alerted_at IS NULL OR d.silence_alerted_at < d.last_seen_at)
	`, silenceMinutes)
	if err != nil {
		log.Printf("[device-heartbeat] query error: %v", err)
		return
	}
	defer rows.Close()

	var targets []silentDevice
	for rows.Next() {
		var row silentDevice
		if err := rows.Scan(&row.userID, &row.name, &row.id, &row.kind); err != nil {
			continue
		}
		targets = append(targets, row)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[device-heartbeat] query error: %v", err)
		return
	}
	if len(targets) == 0 {
		return
	}
	log.Printf("[device-heartbeat] %d silent device(s)", len(targets))

	for _, row := range targets {
		label, ok := deviceKindLabel(row.kind)
		if !ok {
			continue
		}
		body := deviceSilenceBody(row.name, label)
		// silenceDedup stays false. The account-level 6-hour skip would hide a
		// second device that went quiet in the same window.
		if err := h.storeAndPushProtectionText(row.userID, row.name, "heartbeat_silence", body, false, deviceCollapseKey(row.id)); err != nil {
			log.Printf("[device-heartbeat] alert user=%d: %v", row.userID, err)
			continue
		}
		if _, err := h.DB.ExecContext(ctx, `
			UPDATE device_heartbeats
			SET silence_alerted_at = NOW()
			WHERE user_id = $1 AND device_id = $2
		`, row.userID, row.id); err != nil {
			log.Printf("[device-heartbeat] mark user=%d: %v", row.userID, err)
		}
	}
}
