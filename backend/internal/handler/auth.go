package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	rfauth "remain-faithful/backend/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

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
		`INSERT INTO users (name, email, password_hash)
		 VALUES ($1, $2, $3)
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

	// Auto-accept any pending partner invites for this email address.
	go h.acceptPendingInvites(id, req.Email)

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         id,
		"name":       req.Name,
		"email":      req.Email,
		"created_at": createdAt,
	})
}

// acceptPendingInvites runs in a goroutine after registration to create
// relationships for any pending partner invites sent to this email.
func (h *H) acceptPendingInvites(newUserID int64, email string) {
	rows, err := h.DB.Query(
		`SELECT inviter_id, token FROM relationship_invites
		 WHERE invitee_email = $1 AND status = 'pending'
		   AND created_at > NOW() - INTERVAL '7 days'`,
		email,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var inviterID int64
		var token string
		if err := rows.Scan(&inviterID, &token); err != nil {
			continue
		}
		h.DB.Exec( //nolint
			`INSERT INTO relationships (user_id, partner_id, type, status)
			 VALUES ($1, $2, 'partner', 'accepted')
			 ON CONFLICT (user_id, partner_id) DO UPDATE SET status = 'accepted'`,
			inviterID, newUserID,
		)
		h.DB.Exec( //nolint
			`UPDATE relationship_invites SET status = 'accepted' WHERE token = $1`, token)
	}
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
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name, password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&id, &name, &hash)
	if err != nil {
		// Constant-time response to prevent user enumeration.
		bcrypt.CompareHashAndPassword([]byte("$2a$10$placeholder"), []byte(req.Password)) //nolint
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := rfauth.Sign(id, req.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"id":    id,
			"name":  name,
			"email": req.Email,
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
	email, _ := rfauth.EmailFromContext(r.Context())
	token, err := rfauth.Sign(userID, email)
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
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var userID int64
	var name string
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name FROM users WHERE email = $1`, req.Email,
	).Scan(&userID, &name)
	// Always respond 200 to prevent email enumeration.
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	// Generate a 32-byte random hex token.
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	token := hex.EncodeToString(buf)
	expiresAt := time.Now().Add(time.Hour)

	_, err = h.DB.ExecContext(r.Context(), `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
			SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at
	`, userID, token, expiresAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store reset token")
		return
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s",
		getAppBaseURL(), token)
	if h.Email != nil {
		go h.Email.SendPasswordReset(req.Email, name, resetURL) //nolint
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
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
		`SELECT user_id, expires_at FROM password_reset_tokens WHERE token = $1`, req.Token,
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
	if _, err = tx.ExecContext(r.Context(),
		`UPDATE users SET password_hash = $1 WHERE id = $2`, string(hash), userID,
	); err != nil {
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

func getAppBaseURL() string {
	if u := os.Getenv("APP_BASE_URL"); u != "" {
		return u
	}
	return "https://app.remainfaithful.app"
}
