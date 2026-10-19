package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	rfauth "remain-faithful/backend/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

// dummyPasswordHash is a real bcrypt hash so a missing account takes about
// as long as a wrong password. Generated once at startup.
var dummyPasswordHash []byte

func init() {
	h, err := bcrypt.GenerateFromPassword([]byte("remain-faithful-timing-pad"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	dummyPasswordHash = h
}

// Register creates a new user account.
// POST /auth/register
// Body: { "name": "...", "email": "...", "password": "..." }
func (h *H) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "name, email, and password are required")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	var id int64
	var createdAt string
	err = h.DB.QueryRowContext(r.Context(),
		`INSERT INTO users (name, email, password_hash, email_verified)
		 VALUES ($1, $2, $3, FALSE)
		 RETURNING id, created_at`,
		req.Name, req.Email, string(hash),
	).Scan(&id, &createdAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	// Invites stay pending until this address is confirmed. The account can
	// still sign in (build 9 does that immediately after register).
	h.sendEmailVerification(r.Context(), id, req.Email, req.Name)

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":             id,
		"name":           req.Name,
		"email":          req.Email,
		"created_at":     createdAt,
		"email_verified": false,
	})
}

// redeemPendingInvites accepts pending partner and group invites for a newly
// created account. It is bounded by a short timeout and survives client
// disconnects so a half-finished sign-up still lands in the right group.
func (h *H) redeemPendingInvites(parent context.Context, newUserID int64, email string) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), 5*time.Second)
	defer cancel()
	h.acceptPendingInvites(ctx, newUserID, email)
	h.acceptPendingGroupInvites(ctx, newUserID, email)
}

// acceptPendingInvites creates relationships for any pending partner invites
// sent to this email.
func (h *H) acceptPendingInvites(ctx context.Context, newUserID int64, email string) {
	rows, err := h.DB.QueryContext(ctx,
		`SELECT inviter_id, token FROM relationship_invites
		 WHERE invitee_email = $1 AND status = 'pending'
		   AND created_at > NOW() - INTERVAL '7 days'`,
		email,
	)
	if err != nil {
		log.Printf("[acceptPendingInvites] query error for user=%d: %v", newUserID, err)
		return
	}
	type pending struct {
		inviterID int64
		token     string
	}
	var invites []pending
	for rows.Next() {
		var p pending
		if err := rows.Scan(&p.inviterID, &p.token); err != nil {
			log.Printf("[acceptPendingInvites] scan error: %v", err)
			continue
		}
		invites = append(invites, p)
	}
	rows.Close()

	for _, p := range invites {
		if _, err := h.DB.ExecContext(ctx,
			`INSERT INTO relationships (user_id, partner_id, type, status)
			 VALUES ($1, $2, 'partner', 'accepted')
			 ON CONFLICT (user_id, partner_id) DO UPDATE SET status = 'accepted'`,
			p.inviterID, newUserID,
		); err != nil {
			log.Printf("[acceptPendingInvites] insert relationship error: %v", err)
			continue
		}
		if _, err := h.DB.ExecContext(ctx,
			`UPDATE relationship_invites SET status = 'accepted' WHERE token = $1`, p.token,
		); err != nil {
			log.Printf("[acceptPendingInvites] update invite error: %v", err)
		}
	}
}

// acceptPendingGroupInvites adds a newly created account to every group it was
// invited to by email (pending, sent within the last 7 days). The 12-member
// cap is enforced under a row lock on the group; if the group is full the
// invite is left pending.
func (h *H) acceptPendingGroupInvites(ctx context.Context, newUserID int64, email string) {
	rows, err := h.DB.QueryContext(ctx,
		`SELECT id, group_id FROM group_invites
		 WHERE invitee_email = $1 AND status = 'pending'
		   AND created_at > NOW() - INTERVAL '7 days'
		 ORDER BY created_at ASC`,
		email,
	)
	if err != nil {
		log.Printf("[acceptPendingGroupInvites] query error for user=%d: %v", newUserID, err)
		return
	}
	type pending struct {
		inviteID int64
		groupID  int64
	}
	var invites []pending
	for rows.Next() {
		var p pending
		if err := rows.Scan(&p.inviteID, &p.groupID); err != nil {
			log.Printf("[acceptPendingGroupInvites] scan error: %v", err)
			continue
		}
		invites = append(invites, p)
	}
	rows.Close()

	for _, p := range invites {
		if err := h.joinGroupFromInvite(ctx, newUserID, p.inviteID, p.groupID); err != nil {
			log.Printf("[acceptPendingGroupInvites] invite=%d group=%d: %v", p.inviteID, p.groupID, err)
		}
	}
}

var errGroupFull = errors.New("group is full; invite left pending")

func (h *H) joinGroupFromInvite(ctx context.Context, userID, inviteID, groupID int64) error {
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var lockedID int64
	if err := tx.QueryRowContext(ctx,
		`SELECT id FROM groups WHERE id = $1 FOR UPDATE`, groupID,
	).Scan(&lockedID); err != nil {
		return fmt.Errorf("lock group: %w", err)
	}
	var memberCount int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM group_members WHERE group_id = $1`, groupID,
	).Scan(&memberCount); err != nil {
		return fmt.Errorf("count members: %w", err)
	}
	if memberCount >= maxGroupMembers {
		return errGroupFull
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO group_members (group_id, user_id, role)
		 VALUES ($1, $2, 'member')
		 ON CONFLICT (group_id, user_id) DO NOTHING`,
		groupID, userID,
	); err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE group_invites SET status = 'accepted' WHERE id = $1`, inviteID,
	); err != nil {
		return fmt.Errorf("mark accepted: %w", err)
	}
	return tx.Commit()
}

// Login authenticates a user and returns a JWT.
// POST /auth/login
// Body: { "email": "...", "password": "..." }
func (h *H) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var id int64
	var name, hash string
	var tokenVersion int
	var emailVerified bool
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name, password_hash, token_version, email_verified FROM users WHERE email = $1`,
		req.Email,
	).Scan(&id, &name, &hash, &tokenVersion, &emailVerified)
	if err != nil {
		bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(req.Password)) //nolint
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := rfauth.Sign(id, req.Email, tokenVersion)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"id":             id,
			"name":           name,
			"email":          req.Email,
			"email_verified": emailVerified,
		},
	})
}

// RefreshToken issues a fresh JWT to an already-authenticated caller.
// POST /auth/refresh
func (h *H) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var email string
	var tokenVersion int
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT email, token_version FROM users WHERE id = $1`, userID,
	).Scan(&email, &tokenVersion); err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	token, err := rfauth.Sign(userID, email, tokenVersion)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// ForgotPassword generates a reset token and emails it to the user.
// POST /auth/forgot-password
// Body: { "email": "..." }
func (h *H) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// Always respond 200 to prevent email enumeration.
	_ = h.issuePasswordReset(r.Context(), req.Email)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// issuePasswordReset stores a reset token and emails the link.
// Results are "sent", "missing", "limited", or "failed".
func (h *H) issuePasswordReset(ctx context.Context, email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || h == nil || h.DB == nil {
		return "missing"
	}
	var userID int64
	var name string
	var sentCount int
	var windowStart sql.NullTime
	err := h.DB.QueryRowContext(ctx,
		`SELECT id, name, reset_mail_count, reset_mail_window FROM users WHERE email = $1`, email,
	).Scan(&userID, &name, &sentCount, &windowStart)
	if errors.Is(err, sql.ErrNoRows) {
		return "missing"
	}
	if err != nil {
		return "failed"
	}
	now := time.Now()
	if !windowStart.Valid || now.Sub(windowStart.Time) >= time.Hour {
		sentCount = 0
		windowStart.Time = now
		windowStart.Valid = true
	}
	if sentCount >= resetMailPerHour {
		return "limited"
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "failed"
	}
	token := hex.EncodeToString(buf)
	expiresAt := now.Add(time.Hour)

	_, err = h.DB.ExecContext(ctx, `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
			SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at
	`, userID, hashToken(token), expiresAt)
	if err != nil {
		return "failed"
	}
	_, _ = h.DB.ExecContext(ctx,
		`UPDATE users SET reset_mail_count = $1, reset_mail_window = $2 WHERE id = $3`,
		sentCount+1, windowStart.Time, userID,
	)

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", getAppBaseURL(), token)
	if h.Email != nil {
		if err := h.Email.SendPasswordReset(email, name, resetURL); err != nil {
			return "failed"
		}
	}
	return "sent"
}

// ResetPassword validates a reset token and updates the password.
// POST /auth/reset-password
// Body: { "token": "...", "password": "..." }
func (h *H) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	var userID int64
	var expiresAt time.Time
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT user_id, expires_at FROM password_reset_tokens WHERE token = $1`, hashToken(req.Token),
	).Scan(&userID, &expiresAt)
	if err != nil || time.Now().After(expiresAt) {
		writeError(w, http.StatusBadRequest, "invalid or expired reset token")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback() //nolint
	if _, err = tx.ExecContext(r.Context(), `
		UPDATE users
		SET password_hash = $1,
		    token_version = token_version + 1,
		    pending_email = NULL,
		    pending_email_token_hash = NULL,
		    pending_email_expires_at = NULL
		WHERE id = $2
	`, string(hash), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update password")
		return
	}
	if _, err = tx.ExecContext(r.Context(),
		`DELETE FROM password_reset_tokens WHERE user_id = $1`, userID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to clear reset token")
		return
	}
	if err = tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to commit")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// sendEmailVerification stores a hashed one-time token and emails the link.
func (h *H) sendEmailVerification(ctx context.Context, userID int64, email, name string) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		log.Printf("[verify] token user=%d: %v", userID, err)
		return
	}
	token := hex.EncodeToString(buf)
	expires := time.Now().Add(24 * time.Hour)
	if _, err := h.DB.ExecContext(ctx, `
		INSERT INTO email_verifications (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		  SET token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at, created_at = NOW()
	`, userID, hashToken(token), expires); err != nil {
		log.Printf("[verify] store user=%d: %v", userID, err)
		return
	}
	if h.Email == nil {
		return
	}
	link := fmt.Sprintf("%s/verify-email?token=%s", getAppBaseURL(), token)
	if err := h.Email.SendVerifyEmail(email, name, link); err != nil {
		log.Printf("[verify] send user=%d: %v", userID, err)
	}
}

// VerifyEmail confirms a signup address, or applies a pending email change.
// POST /auth/verify-email
// Body: { "token": "..." }
func (h *H) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Token) == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	sum := hashToken(strings.TrimSpace(req.Token))

	var userID int64
	var email string
	err := h.DB.QueryRowContext(r.Context(), `
		SELECT u.id, u.email
		FROM email_verifications v
		JOIN users u ON u.id = v.user_id
		WHERE v.token_hash = $1 AND v.expires_at > NOW()
	`, sum).Scan(&userID, &email)
	if err == nil {
		if _, err = h.DB.ExecContext(r.Context(), `UPDATE users SET email_verified = TRUE WHERE id = $1`, userID); err != nil {
			writeError(w, http.StatusInternalServerError, "could not confirm email")
			return
		}
		_, _ = h.DB.ExecContext(r.Context(), `DELETE FROM email_verifications WHERE user_id = $1`, userID)
		h.redeemPendingInvites(r.Context(), userID, email)
		writeJSON(w, http.StatusOK, map[string]any{"status": "verified", "email": email})
		return
	}

	var pending string
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT id, pending_email FROM users
		WHERE pending_email_token_hash = $1
		  AND pending_email_expires_at > NOW()
		  AND pending_email IS NOT NULL
	`, sum).Scan(&userID, &pending)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid or expired link")
		return
	}
	res, err := h.DB.ExecContext(r.Context(), `
		UPDATE users
		SET email = pending_email,
		    email_verified = TRUE,
		    pending_email = NULL,
		    pending_email_token_hash = NULL,
		    pending_email_expires_at = NULL,
		    token_version = token_version + 1
		WHERE id = $1 AND pending_email = $2
	`, userID, pending)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			writeError(w, http.StatusConflict, "that email cannot be used")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not update email")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusBadRequest, "invalid or expired link")
		return
	}
	h.redeemPendingInvites(r.Context(), userID, pending)
	writeJSON(w, http.StatusOK, map[string]any{"status": "email_updated", "email": pending})
}

// ResendVerification emails a fresh confirmation link to the signed-in user.
// POST /auth/resend-verification
func (h *H) ResendVerification(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var name, email string
	var verified bool
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name, email, email_verified FROM users WHERE id = $1`, userID,
	).Scan(&name, &email, &verified); err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if verified {
		writeJSON(w, http.StatusOK, map[string]string{"status": "already_verified"})
		return
	}
	h.sendEmailVerification(r.Context(), userID, email, name)
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func getAppBaseURL() string {
	if u := os.Getenv("APP_BASE_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	if u := os.Getenv("SITE_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://remainfaithful.com"
}
