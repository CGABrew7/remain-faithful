package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"
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
// Registered only when APP_ENV is not "production" (see debugPushEnabled).
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

// SendPanicAlert sends a time-sensitive push to every accepted partner and
// every other member of the caller's groups. Returns 400 when that audience
// is empty. The push is metadata only: the caller's name and a fixed line.
// POST /panic
func (h *H) SendPanicAlert(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	callerName, err := h.lookupUserName(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	audience, err := h.supportAudienceIDs(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up who to notify")
		return
	}
	if len(audience) == 0 {
		writeError(w, http.StatusBadRequest, "Add a partner or join a group before sending an alert")
		return
	}

	n := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: fmt.Sprintf("panic-%d", userID),
		Payload: map[string]any{
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
		},
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})

	h.inBackground(func() {
		h.notifyUserIDs(context.Background(), audience, callerName, n)
	})
}

// inBackground runs fn after the response. Tests set H.SyncPush to run it inline.
func (h *H) inBackground(fn func()) {
	if h.SyncPush {
		fn()
		return
	}
	go fn()
}

func (h *H) lookupUserName(ctx context.Context, userID int64) (string, error) {
	var name string
	err := h.DB.QueryRowContext(ctx, `SELECT name FROM users WHERE id = $1`, userID).Scan(&name)
	return name, err
}

// supportAudienceIDs is every accepted partner (either direction) and every
// other member of the user's groups. The user is never included.
func (h *H) supportAudienceIDs(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT DISTINCT other_id FROM (
			SELECT CASE WHEN user_id = $1 THEN partner_id ELSE user_id END AS other_id
			FROM   relationships
			WHERE  (user_id = $1 OR partner_id = $1)
			  AND  status = 'accepted'
			UNION
			SELECT gm2.user_id AS other_id
			FROM   group_members gm1
			JOIN   group_members gm2
			       ON gm2.group_id = gm1.group_id AND gm2.user_id <> $1
			WHERE  gm1.user_id = $1
		) audience
		WHERE other_id <> $1 AND other_id IS NOT NULL
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIDs(rows)
}

// coMemberIDs returns the other members of every group userID belongs to.
func (h *H) coMemberIDs(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT DISTINCT gm2.user_id
		FROM   group_members gm1
		JOIN   group_members gm2
		       ON gm2.group_id = gm1.group_id AND gm2.user_id <> $1
		WHERE  gm1.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIDs(rows)
}

// otherMemberIDs returns everyone in groupID except userID.
func (h *H) otherMemberIDs(ctx context.Context, groupID, userID int64) ([]int64, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT user_id FROM group_members
		WHERE  group_id = $1 AND user_id <> $2
	`, groupID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIDs(rows)
}

func scanIDs(rows *sql.Rows) ([]int64, error) {
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (h *H) notifyUserIDs(ctx context.Context, userIDs []int64, senderName string, n *apns.Notification) {
	seen := make(map[int64]struct{}, len(userIDs))
	for _, id := range userIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		note := *n
		h.notifyPartnerByID(ctx, id, senderName, &note)
	}
}

func metadataPush(title, body, notifType, senderName, collapse string) *apns.Notification {
	return &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: collapse,
		Payload: map[string]any{
			"aps": map[string]any{
				"alert": map[string]string{
					"title": title,
					"body":  body,
				},
				"sound": "default",
			},
			"notification_type": notifType,
			"sender_name":       senderName,
		},
	}
}

func encouragementBody(name string) string {
	return name + " sent you encouragement."
}

func covenantUpdatedBody(name string) string {
	return name + " updated the group covenant."
}

func groupLeftBody(name string) string {
	return name + " left the group."
}

func partnershipEndedBody(name string) string {
	return name + " ended the accountability partnership."
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
