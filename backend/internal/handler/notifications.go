package handler

import (
	"context"
	"database/sql"
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
// Body: {"token":"<hex>","platform":"ios","environment":"sandbox"|"production"}
func (h *H) RegisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Token       string `json:"token"`
		Platform    string `json:"platform"`
		Environment string `json:"environment"`
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
	if req.Environment != "sandbox" && req.Environment != "production" {
		req.Environment = "sandbox"
	}

	_, err := h.DB.ExecContext(r.Context(), `
		INSERT INTO device_tokens (user_id, token, platform, environment, is_active, updated_at)
		VALUES ($1, $2, $3, $4, TRUE, NOW())
		ON CONFLICT (user_id, token) DO UPDATE
			SET is_active   = TRUE,
			    platform    = EXCLUDED.platform,
			    environment = EXCLUDED.environment,
			    updated_at  = NOW()
	`, userID, req.Token, req.Platform, req.Environment)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to register device token")
		return
	}

	log.Printf("[push] registered token %.8s... user=%d env=%s", req.Token, userID, req.Environment)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// SendTestPush sends a test push notification to all active tokens of the
// authenticated user. Used to verify APNs configuration end-to-end.
// POST /debug/test-push
func (h *H) SendTestPush(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	env := h.APNS.Environment()
	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT token, environment FROM device_tokens
		WHERE user_id = $1 AND is_active = TRUE
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query tokens")
		return
	}
	defer rows.Close()

	type tokenRow struct {
		token       string
		environment string
	}
	var tokens []tokenRow
	for rows.Next() {
		var t tokenRow
		if rows.Scan(&t.token, &t.environment) == nil {
			tokens = append(tokens, t)
		}
	}

	type result struct {
		Token       string `json:"token"`
		Environment string `json:"environment"`
		Sent        bool   `json:"sent"`
		Error       string `json:"error,omitempty"`
	}
	var results []result

	for _, t := range tokens {
		res := result{Token: t.token[:min(8, len(t.token))] + "...", Environment: t.environment}
		if t.environment != env {
			res.Error = fmt.Sprintf("token env=%s but server env=%s — skipped", t.environment, env)
			results = append(results, res)
			continue
		}
		payload := map[string]any{
			"aps": map[string]any{
				"alert": map[string]string{
					"title": "APNs Test",
					"body":  "Push notifications are working ✓",
				},
				"sound": "default",
			},
			"notification_type": "TEST",
		}
		n := &apns.Notification{
			DeviceToken: t.token,
			PushType:    "alert",
			Priority:    10,
			Payload:     payload,
		}
		if sendErr := h.APNS.Send(r.Context(), n); sendErr != nil {
			res.Error = sendErr.Error()
		} else {
			res.Sent = true
		}
		results = append(results, res)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"apns_configured": !h.APNS.IsNoop(),
		"server_env":      env,
		"token_count":     len(tokens),
		"results":         results,
	})
}


// SendPanicAlert sends a time-sensitive push notification to the caller's
// designated primary partner (or most recently added accepted partner as a
// fallback). Returns 400 if no accepted partner exists.
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

	// Prefer the primary partner; fall back to the most recently added accepted one.
	var partnerID int64
	err := h.DB.QueryRowContext(r.Context(), `
		SELECT partner_id FROM relationships
		WHERE user_id = $1 AND status = 'accepted' AND is_primary = TRUE
		LIMIT 1
	`, userID).Scan(&partnerID)
	if errors.Is(err, sql.ErrNoRows) {
		err = h.DB.QueryRowContext(r.Context(), `
			SELECT partner_id FROM relationships
			WHERE user_id = $1 AND status = 'accepted'
			ORDER BY created_at DESC
			LIMIT 1
		`, userID).Scan(&partnerID)
	}
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusBadRequest, "no accountability partner set — add a partner first")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up partner")
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

	// Respond immediately; fire push in the background.
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})

	go h.notifyPartnerByID(context.Background(), partnerID, callerName, n)
}

// notifyPartners queries all active device tokens belonging to accepted partners
// of userID and calls h.APNS.Send for each. On ErrInvalidToken it marks the
// token inactive and sends an APP_DELETED notification to those partners.
func (h *H) notifyPartners(ctx context.Context, userID int64, callerName string, n *apns.Notification) {
	env := h.APNS.Environment()
	rows, err := h.DB.QueryContext(ctx, `
		SELECT dt.token
		FROM   relationships r
		JOIN   device_tokens dt ON dt.user_id = r.partner_id
		WHERE  r.user_id    = $1
		  AND  r.status     = 'accepted'
		  AND  dt.is_active = TRUE
		  AND  dt.environment = $2
	`, userID, env)
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

// notifyPartnerByID sends n to all active device tokens belonging to the given
// partnerID that match the server's configured APNs environment. On
// ErrInvalidToken it marks the token inactive.
func (h *H) notifyPartnerByID(ctx context.Context, partnerID int64, callerName string, n *apns.Notification) {
	env := h.APNS.Environment()
	rows, err := h.DB.QueryContext(ctx, `
		SELECT token FROM device_tokens
		WHERE user_id = $1 AND is_active = TRUE AND environment = $2
	`, partnerID, env)
	if err != nil {
		log.Printf("[push] notifyPartnerByID: query tokens for user=%d env=%s: %v", partnerID, env, err)
		return
	}
	defer rows.Close()

	sent := 0
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
				log.Printf("[push] notifyPartnerByID: marked inactive token %.8s... for user=%d", token, partnerID)
				continue
			}
			log.Printf("[push] notifyPartnerByID: send to user=%d token=%.8s...: %v", partnerID, token, err)
		} else {
			sent++
		}
	}
	log.Printf("[push] notifyPartnerByID: user=%d env=%s sent=%d", partnerID, env, sent)
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
