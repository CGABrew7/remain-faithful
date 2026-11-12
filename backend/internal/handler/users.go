package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/fcm"
)

// GetMe returns the authenticated user's profile.
// GET /users/me
func (h *H) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var id int64
	var name, email, createdAt, zone string
	var emailVerified bool
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name, email, created_at, email_verified, COALESCE(timezone, 'UTC') FROM users WHERE id = $1`,
		userID,
	).Scan(&id, &name, &email, &createdAt, &emailVerified, &zone)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":             id,
		"name":           name,
		"email":          email,
		"created_at":     createdAt,
		"email_verified": emailVerified,
		"timezone":       zone,
	})
}

// UpdateMe updates the authenticated user's name and/or email.
// PUT /users/me
// Body: { "name": "...", "email": "..." }
func (h *H) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	name, nameErr := cleanPersonName(req.Name)
	if nameErr != "" {
		writeError(w, http.StatusBadRequest, nameErr)
		return
	}
	email, valid := normalizeWaitlistEmail(req.Email)
	if !valid {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address.")
		return
	}
	req.Name = name
	req.Email = email
	var current string
	if err := h.DB.QueryRowContext(r.Context(), `SELECT email FROM users WHERE id = $1`, userID).Scan(&current); err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	var id int64
	var savedName, savedEmail, createdAt string
	pending := ""
	if req.Email == current {
		err := h.DB.QueryRowContext(r.Context(),
			`UPDATE users SET name = $1 WHERE id = $2
			 RETURNING id, name, email, created_at`,
			req.Name, userID,
		).Scan(&id, &savedName, &savedEmail, &createdAt)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update profile")
			return
		}
	} else {
		var other int64
		err := h.DB.QueryRowContext(r.Context(),
			`SELECT id FROM users WHERE email = $1 AND id <> $2`, req.Email, userID,
		).Scan(&other)
		if err == nil {
			writeError(w, http.StatusConflict, "that email cannot be used")
			return
		}
		token, err := inviteToken()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update profile")
			return
		}
		err = h.DB.QueryRowContext(r.Context(), `
			UPDATE users
			SET name = $1,
			    pending_email = $2,
			    pending_email_token_hash = $3,
			    pending_email_expires_at = NOW() + INTERVAL '24 hours'
			WHERE id = $4
			RETURNING id, name, email, created_at
		`, req.Name, req.Email, hashToken(token), userID).Scan(&id, &savedName, &savedEmail, &createdAt)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update profile")
			return
		}
		pending = req.Email
		if h.Email != nil {
			link := getAppBaseURL() + "/verify-email?token=" + token + "&purpose=change"
			_ = h.Email.SendEmailChange(req.Email, savedName, link)
		}
	}

	resp := map[string]any{
		"id":         id,
		"name":       savedName,
		"email":      savedEmail,
		"created_at": createdAt,
	}
	if pending != "" {
		resp["pending_email"] = pending
		resp["email_change"] = "pending"
	}
	writeJSON(w, http.StatusOK, resp)
}

// DeleteMe permanently deletes the authenticated user and all their data.
// Before deletion it promotes a new admin for any group where the user is the
// sole admin, then fires departure push notifications to all partners and
// group co-members.
// DELETE /users/me
func (h *H) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var userName string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&userName); err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	if err := h.deleteUserByID(r.Context(), userID, userName); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete account")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// deleteUserByID removes an account the same way DELETE /users/me does:
// promote a remaining group admin, tell partners, then delete the row.
func (h *H) deleteUserByID(ctx context.Context, userID int64, userName string) error {
	if err := h.promoteAdminsBeforeDeletion(ctx, userID); err != nil {
		log.Printf("[delete] promote admins user=%d: %v", userID, err)
	}

	var tokens []string
	var androidTokens []string
	if h.APNS != nil {
		tokens = h.collectDepartureTokens(ctx, userID)
	}
	if h.FCM != nil {
		androidTokens = h.collectAndroidDepartureTokens(ctx, userID)
	}

	if _, err := h.DB.ExecContext(ctx, `DELETE FROM alerts WHERE subject_user_id = $1`, userID); err != nil {
		return err
	}
	if _, err := h.DB.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		return err
	}
	if len(tokens) > 0 {
		go h.sendDepartureNotifications(context.Background(), userName, tokens)
	}
	if len(androidTokens) > 0 {
		go h.sendAndroidDeparture(context.Background(), userName, androidTokens)
	}
	return nil
}

// promoteAdminsBeforeDeletion promotes the longest-tenured member to admin in
// every group where userID is the sole admin but other members remain.
func (h *H) promoteAdminsBeforeDeletion(ctx context.Context, userID int64) error {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT gm.group_id
		FROM   group_members gm
		WHERE  gm.user_id = $1
		  AND  gm.role    = 'admin'
		  AND  NOT EXISTS (
		      SELECT 1 FROM group_members x
		      WHERE  x.group_id = gm.group_id
		        AND  x.user_id != $1
		        AND  x.role    = 'admin'
		  )
		  AND  EXISTS (
		      SELECT 1 FROM group_members y
		      WHERE  y.group_id = gm.group_id
		        AND  y.user_id != $1
		  )
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var groupID int64
		if err := rows.Scan(&groupID); err != nil {
			continue
		}
		if _, err := h.DB.ExecContext(ctx, `
			UPDATE group_members
			SET    role = 'admin'
			WHERE  group_id = $1
			  AND  user_id  = (
			      SELECT user_id FROM group_members
			      WHERE  group_id = $1 AND user_id != $2
			      ORDER  BY joined_at ASC
			      LIMIT  1
			  )
		`, groupID, userID); err != nil {
			log.Printf("[delete] promote admin group=%d: %v", groupID, err)
		}
	}
	return rows.Err()
}

// collectDepartureTokens returns active APNs device tokens belonging to the
// departing user's accepted partners and group co-members.
func (h *H) collectDepartureTokens(ctx context.Context, userID int64) []string {
	env := h.APNS.Environment()
	rows, err := h.DB.QueryContext(ctx, `
		SELECT DISTINCT dt.token
		FROM   device_tokens dt
		WHERE  dt.is_active   = TRUE
		  AND  dt.environment = $2
		  AND  dt.platform <> 'android'
		  AND  dt.user_id IN (
		      SELECT CASE WHEN r.user_id = $1 THEN r.partner_id ELSE r.user_id END
		      FROM   relationships r
		      WHERE  (r.user_id = $1 OR r.partner_id = $1)
		        AND  r.status = 'accepted'
		      UNION
		      SELECT gm2.user_id
		      FROM   group_members gm1
		      JOIN   group_members gm2
		             ON gm2.group_id = gm1.group_id AND gm2.user_id != $1
		      WHERE  gm1.user_id = $1
		  )
	`, userID, env)
	if err != nil {
		log.Printf("[delete] collectDepartureTokens user=%d: %v", userID, err)
		return nil
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var t string
		if rows.Scan(&t) == nil {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

// sendDepartureNotifications pushes a departure notice to every provided token.
func (h *H) sendDepartureNotifications(ctx context.Context, userName string, tokens []string) {
	userName = clipPushName(userName)
	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Partner Update",
				"body":  userName + " has left Remain Faithful and is no longer being monitored.",
			},
			"sound": "default",
		},
		"notification_type": "PARTNER_LEFT",
		"sender_name":       userName,
	}
	sent := 0
	for _, token := range tokens {
		n := &apns.Notification{
			DeviceToken: token,
			PushType:    "alert",
			Priority:    10,
			Payload:     payload,
		}
		if err := h.APNS.Send(ctx, n); err != nil {
			var invalidErr *apns.ErrInvalidToken
			if errors.As(err, &invalidErr) {
				h.markTokenInactive(ctx, token)
				continue
			}
			log.Printf("[delete] push to %.8s...: %v", token, err)
		} else {
			sent++
		}
	}
	log.Printf("[delete] departure notifications sent=%d/%d user=%s", sent, len(tokens), userName)
}

// collectAndroidDepartureTokens returns active FCM tokens for partners and
// group co-members. APNs environment is not applied.
func (h *H) collectAndroidDepartureTokens(ctx context.Context, userID int64) []string {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT token FROM device_tokens
		WHERE is_active = TRUE
		  AND platform = 'android'
		  AND user_id IN (
		      SELECT CASE WHEN user_id = $1 THEN partner_id ELSE user_id END
		      FROM relationships
		      WHERE (user_id = $1 OR partner_id = $1) AND status = 'accepted'
		      UNION
		      SELECT gm2.user_id
		      FROM group_members gm1
		      JOIN group_members gm2
		        ON gm2.group_id = gm1.group_id AND gm2.user_id != $1
		      WHERE gm1.user_id = $1
		  )
	`, userID)
	if err != nil {
		log.Printf("[delete] collectAndroidDepartureTokens user=%d: %v", userID, err)
		return nil
	}
	defer rows.Close()
	var tokens []string
	for rows.Next() {
		var t string
		if rows.Scan(&t) == nil && t != "" {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func (h *H) sendAndroidDeparture(ctx context.Context, userName string, tokens []string) {
	if h.FCM == nil {
		return
	}
	userName = clipPushName(userName)
	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Partner Update",
				"body":  userName + " has left Remain Faithful and is no longer being monitored.",
			},
			"sound": "default",
		},
		"notification_type": "PARTNER_LEFT",
		"sender_name":       userName,
	}
	for _, token := range tokens {
		msg := fcm.FromAPNS(token, "", payload)
		if err := h.FCM.Send(ctx, msg); err != nil {
			var invalid *fcm.ErrInvalidToken
			if errors.As(err, &invalid) {
				h.markTokenInactive(ctx, token)
				continue
			}
			log.Printf("[delete] fcm to %.8s...: %v", token, err)
		}
	}
}

// GetUserStats returns the authenticated user's streak.
// A day counts only when protection checked in that day and nothing was flagged.
// Days are the account time zone, which the device sends on heartbeat.
// GET /users/me/stats
// week and week_states keep the words older apps already show.
// week_coverage uses off for a day protection did not check in.
// Index 0 is 6 days ago and index 6 is today.
func (h *H) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var createdAt time.Time
	var zone string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT created_at, COALESCE(timezone, 'UTC') FROM users WHERE id = $1`, userID,
	).Scan(&createdAt, &zone); err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	results, err := h.streaksAt(r.Context(), map[int64]streakMeta{
		userID: {Created: createdAt, Zone: zone},
	}, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load stats")
		return
	}
	stats := results[userID]
	writeJSON(w, http.StatusOK, map[string]any{
		"streak_days":   stats.Days,
		"best_streak":   stats.Best,
		"week":          stats.Week,
		"week_states":   stats.States,
		"week_coverage": stats.Coverage,
		"timezone":      zoneLocation(zone).String(),
	})
}
