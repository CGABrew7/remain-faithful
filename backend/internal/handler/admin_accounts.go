package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"remain-faithful/backend/internal/admin"
)

const accountLookupSQL = `
SELECT u.created_at,
       u.email_verified,
       (u.apple_id IS NOT NULL) AS apple,
       (u.google_id IS NOT NULL) AS google,
       u.last_heartbeat_at,
       (SELECT count(*) FROM relationships r
         WHERE r.status = 'accepted' AND r.type = 'partner'
           AND (r.user_id = u.id OR r.partner_id = u.id))
FROM users u
WHERE u.email = $1`

// AdminLookupAccount returns status for one exact email address.
// GET /admin/account?email=
func (h *H) AdminLookupAccount(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("email")))
	if email == "" || !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "an exact email is required")
		return
	}
	status, err := h.lookupAccountStatus(r.Context(), email)
	if errors.Is(err, sql.ErrNoRows) {
		h.auditAdmin(r.Context(), "lookup_account", "not found")
		writeError(w, http.StatusNotFound, "no account with that email")
		return
	}
	if err != nil {
		log.Printf("admin lookup: %v", err)
		writeError(w, http.StatusInternalServerError, "could not look up account")
		return
	}
	h.auditAdmin(r.Context(), "lookup_account", email)
	writeJSON(w, http.StatusOK, status)
}

func (h *H) lookupAccountStatus(ctx context.Context, email string) (*admin.AccountStatus, error) {
	var (
		created  time.Time
		verified bool
		apple    bool
		google   bool
		lastBeat sql.NullTime
		partners int64
	)
	err := h.DB.QueryRowContext(ctx, accountLookupSQL, email).Scan(
		&created, &verified, &apple, &google, &lastBeat, &partners,
	)
	if err != nil {
		return nil, err
	}
	method := "email"
	if apple {
		method = "apple"
	} else if google {
		method = "google"
	}
	status := &admin.AccountStatus{
		Email:         email,
		CreatedAt:     created,
		EmailVerified: verified,
		SignInMethod:  method,
		PartnerCount:  int(partners),
	}
	if lastBeat.Valid {
		t := lastBeat.Time
		status.LastActiveAt = &t
	}
	return status, nil
}

// AdminResendVerification sends the same confirmation email a signed-in user can request.
// POST /admin/account/resend-verification
func (h *H) AdminResendVerification(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	email, ok := adminEmailBody(w, r)
	if !ok {
		return
	}
	var id int64
	var name string
	var verified bool
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name, email_verified FROM users WHERE email = $1`, email,
	).Scan(&id, &name, &verified)
	if errors.Is(err, sql.ErrNoRows) {
		h.auditAdmin(r.Context(), "resend_verification", "not found")
		writeError(w, http.StatusNotFound, "no account with that email")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not look up account")
		return
	}
	if verified {
		h.auditAdmin(r.Context(), "resend_verification", "already verified")
		writeJSON(w, http.StatusOK, map[string]string{"status": "already_verified"})
		return
	}
	h.sendEmailVerification(r.Context(), id, email, name)
	h.auditAdmin(r.Context(), "resend_verification", email)
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

// AdminSendPasswordReset uses the same reset email as Forgot password.
// POST /admin/account/password-reset
func (h *H) AdminSendPasswordReset(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	email, ok := adminEmailBody(w, r)
	if !ok {
		return
	}
	switch h.issuePasswordReset(r.Context(), email) {
	case "sent":
		h.auditAdmin(r.Context(), "password_reset", email)
		writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
	case "limited":
		h.auditAdmin(r.Context(), "password_reset", "rate limited")
		writeError(w, http.StatusTooManyRequests, "a reset was already sent recently")
	case "missing":
		h.auditAdmin(r.Context(), "password_reset", "not found")
		writeError(w, http.StatusNotFound, "no account with that email")
	default:
		h.auditAdmin(r.Context(), "password_reset", "failed")
		writeError(w, http.StatusInternalServerError, "could not send reset")
	}
}

// AdminDeleteAccount deletes one account after the email is typed twice.
// It uses the same removal path as DELETE /users/me.
// POST /admin/account/delete
func (h *H) AdminDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	var req struct {
		Email   string `json:"email"`
		Confirm string `json:"confirm"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	confirm := strings.ToLower(strings.TrimSpace(req.Confirm))
	if email == "" || confirm == "" || email != confirm {
		writeError(w, http.StatusBadRequest, "type the email address to confirm")
		return
	}
	var id int64
	var name string
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name FROM users WHERE email = $1`, email,
	).Scan(&id, &name)
	if errors.Is(err, sql.ErrNoRows) {
		h.auditAdmin(r.Context(), "delete_account", "not found")
		writeError(w, http.StatusNotFound, "no account with that email")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not look up account")
		return
	}
	if err := h.deleteUserByID(r.Context(), id, name); err != nil {
		log.Printf("admin delete: %v", err)
		h.auditAdmin(r.Context(), "delete_account_failed", email)
		writeError(w, http.StatusInternalServerError, "could not delete account")
		return
	}
	h.auditAdmin(r.Context(), "delete_account", email)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func adminEmailBody(w http.ResponseWriter, r *http.Request) (string, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return "", false
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "an exact email is required")
		return "", false
	}
	return email, true
}
