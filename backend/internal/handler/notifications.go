package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/apns"
)

// RegisterDeviceToken upserts an APNs device token for the authenticated user.
// POST /users/device-token
// Body: {"token":"<hex>","platform":"ios"}
func (h *H) RegisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.Platform == "" {
		req.Platform = "ios"
	}

	_, err := h.DB.ExecContext(r.Context(), `
		INSERT INTO device_tokens (user_id, token, platform, is_active, updated_at)
		VALUES ($1, $2, $3, TRUE, NOW())
		ON CONFLICT (user_id, token) DO UPDATE
			SET is_active  = TRUE,
			    platform   = EXCLUDED.platform,
			    updated_at = NOW()
	`, userID, req.Token, req.Platform)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to register device token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// SendPanicAlert sends a time-sensitive push notification to all accepted
// accountability partners of the authenticated user.
// POST /panic
func (h *H) SendPanicAlert(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var callerName string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&callerName); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Urgent Prayer Request",
				"body":  callerName + " needs support right now",
			},
			"sound":              "default",
			"interruption-level": "time-sensitive",
		},
		"notification_type": "PANIC_ALERT",
		"sender_name":       callerName,
	}

	n := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: fmt.Sprintf("panic-%d", userID),
		Payload:    payload,
	}

	// Respond immediately; fire pushes in the background.
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})

	go h.notifyPartners(context.Background(), userID, callerName, n)
}

// notifyPartners queries all active device tokens belonging to accepted partners
// of userID and calls h.APNS.Send for each. On ErrInvalidToken it marks the
// token inactive and sends an APP_DELETED notification to those partners.
func (h *H) notifyPartners(ctx context.Context, userID int64, callerName string, n *apns.Notification) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT dt.token
		FROM   relationships r
		JOIN   device_tokens dt ON dt.user_id = r.partner_id
		WHERE  r.user_id    = $1
		  AND  r.status     = 'accepted'
		  AND  dt.is_active = TRUE
	`, userID)
	if err != nil {
		log.Printf("notifyPartners: query tokens: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			continue
		}

		n.DeviceToken = token
		if err := h.APNS.Send(ctx, n); err != nil {
			var invalidErr *apns.ErrInvalidToken
			if errors.As(err, &invalidErr) {
				h.markTokenInactive(ctx, token)
				// Notify partners that the app was deleted on that device.
				deletedPayload := map[string]any{
					"aps": map[string]any{
						"alert": map[string]string{
							"title": "Partner Update",
							"body":  callerName + "'s app was removed from a device",
						},
						"sound": "default",
					},
					"notification_type": "APP_DELETED",
					"sender_name":       callerName,
				}
				deleted := &apns.Notification{
					DeviceToken: token,
					PushType:    "alert",
					Priority:    10,
					Payload:     deletedPayload,
				}
				if sendErr := h.APNS.Send(ctx, deleted); sendErr != nil {
					log.Printf("notifyPartners: send APP_DELETED: %v", sendErr)
				}
				continue
			}
			log.Printf("notifyPartners: send to %s: %v", token, err)
		}
	}
}

// markTokenInactive sets is_active=FALSE for the given device token.
func (h *H) markTokenInactive(ctx context.Context, token string) {
	_, err := h.DB.ExecContext(ctx,
		`UPDATE device_tokens SET is_active = FALSE WHERE token = $1`, token,
	)
	if err != nil {
		log.Printf("markTokenInactive: %v", err)
	}
}
