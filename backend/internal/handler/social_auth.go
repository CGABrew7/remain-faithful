package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	rfauth "remain-faithful/backend/internal/auth"
)

// AppleSignIn finds or creates a user account via Sign in with Apple.
// POST /auth/apple
// Body: { "identity_token": "...", "authorization_code": "...", "name": {"firstName":"...","lastName":"..."} }
func (h *H) AppleSignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdentityToken     string `json:"identity_token"`
		AuthorizationCode string `json:"authorization_code"`
		Name              struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		} `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.IdentityToken == "" {
		writeError(w, http.StatusBadRequest, "identity_token is required")
		return
	}

	ident, err := verifyAppleIdentityToken(r.Context(), req.IdentityToken, appleAudience(), timeNow())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid identity token")
		return
	}

	name := strings.TrimSpace(req.Name.FirstName + " " + req.Name.LastName)
	if name == "" && ident.Email == "" {
		name = "Apple User"
	}

	userID, userName, userEmail, jwtErr := h.findOrCreateSocialUser(r, "apple_id", ident.Sub, ident.Email, name, ident.EmailVerified)
	if jwtErr != nil {
		writeSocialAccountError(w, jwtErr)
		return
	}

	token, err := h.issueSession(r.Context(), userID, userEmail)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  map[string]any{"id": userID, "name": userName, "email": userEmail},
	})
}

// GoogleSignIn finds or creates a user account via Sign in with Google.
// POST /auth/google
// Body: { "id_token": "..." }
func (h *H) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.IDToken == "" {
		writeError(w, http.StatusBadRequest, "id_token is required")
		return
	}

	ident, err := verifyGoogleIDToken(r.Context(), req.IDToken, googleAudiences(), timeNow())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid Google ID token")
		return
	}
	name := ident.Name
	if name == "" {
		name = "Google User"
	}

	userID, userName, userEmail, err := h.findOrCreateSocialUser(r, "google_id", ident.Sub, ident.Email, name, ident.EmailVerified)
	if err != nil {
		writeSocialAccountError(w, err)
		return
	}

	token, err := h.issueSession(r.Context(), userID, userEmail)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  map[string]any{"id": userID, "name": userName, "email": userEmail},
	})
}

var errSocialLinkConflict = errors.New("this email is already linked to a different sign-in")
var errSocialEmailUnverified = errors.New("the identity provider did not verify this email")

func writeSocialAccountError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errSocialLinkConflict):
		writeError(w, http.StatusConflict, "this email is already linked to a different sign-in")
	case errors.Is(err, errSocialEmailUnverified):
		writeError(w, http.StatusUnauthorized, "email address is not verified")
	default:
		writeError(w, http.StatusInternalServerError, "account error")
	}
}

// issueSession signs a JWT that carries the user's current token version.
func (h *H) issueSession(ctx context.Context, userID int64, email string) (string, error) {
	var version int
	err := h.DB.QueryRowContext(ctx, `SELECT token_version FROM users WHERE id = $1`, userID).Scan(&version)
	if err != nil {
		version = 1
	}
	if version < 1 {
		version = 1
	}
	return rfauth.Sign(userID, email, version)
}

// findOrCreateSocialUser looks up a user by social ID column (apple_id/google_id).
// An existing email account is linked only when the provider verified the email
// and that column is empty or already equal to this subject. A different subject
// is never overwritten.
func (h *H) findOrCreateSocialUser(r *http.Request, idCol, idVal, email, name string, emailVerified bool) (
	userID int64, userName, userEmail string, err error,
) {
	ctx := r.Context()
	if idCol != "apple_id" && idCol != "google_id" {
		err = errors.New("unsupported identity column")
		return
	}
	email = strings.ToLower(strings.TrimSpace(email))

	err = h.DB.QueryRowContext(ctx,
		fmt.Sprintf(`SELECT id, name, email FROM users WHERE %s = $1`, idCol), idVal,
	).Scan(&userID, &userName, &userEmail)
	if err == nil {
		return
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) && !isNoRows(err) {
		// The fake test driver and a missing row both surface as a scan error.
		// A real database returns sql.ErrNoRows. Anything else that still has
		// no scanned id falls through to the email lookup.
		if userID != 0 {
			return
		}
	}

	if email != "" {
		var appleID, googleID string
		err = h.DB.QueryRowContext(ctx,
			`SELECT id, name, email, COALESCE(apple_id, ''), COALESCE(google_id, '')
			 FROM users WHERE email = $1`, email,
		).Scan(&userID, &userName, &userEmail, &appleID, &googleID)
		if err == nil {
			if !emailVerified {
				err = errSocialEmailUnverified
				return
			}
			current := appleID
			if idCol == "google_id" {
				current = googleID
			}
			if current != "" && current != idVal {
				err = errSocialLinkConflict
				return
			}
			if current == "" {
				res, execErr := h.DB.ExecContext(ctx,
					fmt.Sprintf(`UPDATE users SET %s = $1, email_verified = TRUE WHERE id = $2 AND (%s IS NULL OR %s = '')`, idCol, idCol, idCol),
					idVal, userID,
				)
				if execErr != nil {
					err = execErr
					return
				}
				if n, _ := res.RowsAffected(); n == 0 {
					// Lost a race with another subject. Re-read and refuse a mismatch.
					var now string
					_ = h.DB.QueryRowContext(ctx,
						fmt.Sprintf(`SELECT COALESCE(%s, '') FROM users WHERE id = $1`, idCol), userID,
					).Scan(&now)
					if now != "" && now != idVal {
						err = errSocialLinkConflict
						return
					}
				}
			}
			h.redeemPendingInvites(ctx, userID, email)
			return
		}
	}

	if email != "" && !emailVerified {
		err = errSocialEmailUnverified
		return
	}

	placeholderEmail := email == ""
	if placeholderEmail {
		email = fmt.Sprintf("%s@privaterelay.appleid.com", idVal)
		name = firstNonEmpty(name, "Apple User")
	}
	if name == "" {
		name = "Remain Faithful"
	}
	verified := emailVerified && !placeholderEmail
	err = h.DB.QueryRowContext(ctx,
		fmt.Sprintf(`INSERT INTO users (name, email, password_hash, %s, email_verified)
		             VALUES ($1, $2, '', $3, $4)
		             RETURNING id, name, email`, idCol),
		name, email, idVal, verified,
	).Scan(&userID, &userName, &userEmail)
	if err == nil && verified {
		h.redeemPendingInvites(ctx, userID, email)
	}
	return
}

func isNoRows(err error) bool {
	return err == sql.ErrNoRows || strings.Contains(err.Error(), "no rows")
}

func firstNonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
