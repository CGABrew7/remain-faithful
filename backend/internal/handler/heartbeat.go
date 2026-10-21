package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"
)

// RecordPhoneSignals stores aggregate phone facts for the owner dashboard.
// Screen Time is yes or no, apps are a count, the version is a string, and
// support taps are a running total. App names and messages stay on the phone.
// POST /users/me/signals
func (h *H) RecordPhoneSignals(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		ScreenTimeGranted *bool  `json:"screen_time_granted"`
		MonitoredAppCount *int   `json:"monitored_app_count"`
		AppVersion        string `json:"app_version"`
		SupportTaps       *int   `json:"support_taps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ScreenTimeGranted == nil || req.MonitoredAppCount == nil || req.SupportTaps == nil {
		writeError(w, http.StatusBadRequest, "screen time, app count, and support taps are required")
		return
	}
	if *req.MonitoredAppCount < 0 || *req.MonitoredAppCount > 500 {
		writeError(w, http.StatusBadRequest, "app count is out of range")
		return
	}
	if *req.SupportTaps < 0 || *req.SupportTaps > 100000 {
		writeError(w, http.StatusBadRequest, "support tap count is out of range")
		return
	}
	version := trimVersion(req.AppVersion)
	if version == "" && req.AppVersion != "" {
		writeError(w, http.StatusBadRequest, "app version is not recognized")
		return
	}

	_, err := h.DB.ExecContext(r.Context(), `
		UPDATE users
		SET screen_time_granted = $2,
		    monitored_app_count = $3,
		    app_version = CASE WHEN $4 = '' THEN app_version ELSE $4 END,
		    support_tap_count = GREATEST(support_tap_count, $5)
		WHERE id = $1
	`, userID, *req.ScreenTimeGranted, *req.MonitoredAppCount, version, *req.SupportTaps)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save phone signals")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func trimVersion(raw string) string {
	out := make([]rune, 0, len(raw))
	for _, c := range raw {
		if c == ' ' || c == '\t' {
			continue
		}
		switch {
		case c >= '0' && c <= '9', c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c == '.':
			out = append(out, c)
		default:
			return ""
		}
		if len(out) > 32 {
			return ""
		}
	}
	return string(out)
}

// Heartbeat updates the authenticated user's last_heartbeat_at and screen state.
// POST /heartbeat
// Body: {"screen":"active"|"idle"|"monitoring"}
// "monitoring" is the value the app sends while Alerts or Blocking are on,
// including when Deep Scan is off. It refreshes last_heartbeat_at the same way
// as active and idle. The shorter "active screen stopped" nudge still applies
// only when the stored state is "active".
func (h *H) Heartbeat(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Screen string `json:"screen"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	switch req.Screen {
	case "active", "idle", "monitoring":
	default:
		writeError(w, http.StatusBadRequest, `screen must be "active", "idle", or "monitoring"`)
		return
	}

	_, err := h.DB.ExecContext(r.Context(), `
		UPDATE users
		SET last_heartbeat_at = NOW(), heartbeat_screen_state = $2
		WHERE id = $1
	`, userID, req.Screen)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update heartbeat")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// StartHeartbeatWatcher checks every 5 minutes for stale heartbeats and sends
// push notifications to users and their accountability partners as needed.
// Stops when ctx is cancelled.
func (h *H) StartHeartbeatWatcher(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.checkHeartbeats(ctx)
		}
	}
}

func (h *H) checkHeartbeats(ctx context.Context) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT id, name, last_heartbeat_at, heartbeat_screen_state
		FROM   users
		WHERE  last_heartbeat_at IS NOT NULL
		  AND  (heartbeat_notified_at IS NULL OR heartbeat_notified_at < last_heartbeat_at)
	`)
	if err != nil {
		log.Printf("heartbeat watcher: %v", err)
		return
	}
	defer rows.Close()

	now := time.Now()
	type entry struct {
		userID int64
		name   string
		gap    time.Duration
		screen string
	}
	var users []entry
	for rows.Next() {
		var (
			userID      int64
			name        string
			lastBeat    time.Time
			screenState sql.NullString
		)
		if err := rows.Scan(&userID, &name, &lastBeat, &screenState); err != nil {
			continue
		}
		users = append(users, entry{
			userID: userID,
			name:   name,
			gap:    now.Sub(lastBeat),
			screen: screenState.String,
		})
	}

	for _, u := range users {
		// Skip idle users under 30 min — expected inactivity.
		if u.screen == "idle" && u.gap < 30*time.Minute {
			continue
		}

		shouldNotify := (u.gap >= 5*time.Minute && u.screen == "active") ||
			u.gap >= 30*time.Minute
		if !shouldNotify {
			continue
		}

		body := "Your Remain Faithful monitoring has stopped. Tap to restart."
		if u.gap >= 30*time.Minute {
			body = "Your Remain Faithful monitoring has been inactive for 30+ minutes. Tap to restart."
		}
		h.notifyUserDirect(ctx, u.userID, "monitoring_stopped", "Monitoring Stopped", body)

		// Partner silence alerts are handled by StartHeartbeatSweep in
		// heartbeat_sweep.go, which uses the protection pipeline with the
		// heartbeat_silence alert type. No partner push here.

		// Mark notified so we don't spam until a new heartbeat resets the window.
		_, _ = h.DB.ExecContext(ctx,
			`UPDATE users SET heartbeat_notified_at = NOW() WHERE id = $1`, u.userID,
		)
	}
}

// notifyUserDirect sends a push to all active device tokens belonging to userID.
func (h *H) notifyUserDirect(ctx context.Context, userID int64, collapseID, title, body string) {
	rows, err := h.DB.QueryContext(ctx,
		`SELECT token FROM device_tokens WHERE user_id = $1 AND is_active = TRUE`, userID)
	if err != nil {
		log.Printf("notifyUserDirect %d: %v", userID, err)
		return
	}
	defer rows.Close()

	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{"title": title, "body": body},
			"sound": "default",
		},
		"notification_type": "MONITORING_STOPPED",
	}

	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			continue
		}
		n := &apns.Notification{
			DeviceToken: token,
			PushType:    "alert",
			Priority:    10,
			CollapseID:  fmt.Sprintf("%s-%d", collapseID, userID),
			Payload:     payload,
		}
		if err := h.APNS.Send(ctx, n); err != nil {
			var inv *apns.ErrInvalidToken
			if errors.As(err, &inv) {
				h.markTokenInactive(ctx, token)
				continue
			}
			log.Printf("notifyUserDirect %d: %v", userID, err)
		}
	}
}
