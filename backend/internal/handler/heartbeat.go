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

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/apns"
)

// Heartbeat updates the authenticated user's last_heartbeat_at and screen state.
// POST /heartbeat
// Body: {"screen":"active"|"idle"}
func (h *H) Heartbeat(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Screen string `json:"screen"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Screen != "active" && req.Screen != "idle" {
		writeError(w, http.StatusBadRequest, `screen must be "active" or "idle"`)
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

		if u.gap >= 30*time.Minute && u.screen == "active" {
			h.notifyPartnersOfStoppedMonitoring(ctx, u.userID, u.name)
		}

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

func (h *H) notifyPartnersOfStoppedMonitoring(ctx context.Context, userID int64, name string) {
	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Partner Check-In",
				"body":  name + "'s monitoring has been inactive. Consider reaching out.",
			},
			"sound": "default",
		},
		"notification_type": "PARTNER_MONITORING_STOPPED",
		"sender_name":       name,
	}
	n := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: fmt.Sprintf("partner_monitoring-%d", userID),
		Payload:    payload,
	}
	h.notifyPartners(ctx, userID, name, n)
}
