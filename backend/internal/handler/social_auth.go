package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	appleID, email, err := decodeAppleIdentityToken(req.IdentityToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identity token: "+err.Error())
		return
	}

	name := strings.TrimSpace(req.Name.FirstName + " " + req.Name.LastName)
	if name == "" {
		name = "Apple User"
	}

	userID, userName, userEmail, jwtErr := h.findOrCreateSocialUser(r, "apple_id", appleID, email, name)
	if jwtErr != nil {
		writeError(w, http.StatusInternalServerError, "account error: "+jwtErr.Error())
		return
	}

	token, err := rfauth.Sign(userID, userEmail)
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

	googleID, email, name, err := verifyGoogleIDToken(req.IDToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid Google ID token: "+err.Error())
		return
	}

	userID, userName, userEmail, err := h.findOrCreateSocialUser(r, "google_id", googleID, email, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "account error: "+err.Error())
		return
	}

	token, err := rfauth.Sign(userID, userEmail)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  map[string]any{"id": userID, "name": userName, "email": userEmail},
	})
}

// findOrCreateSocialUser looks up a user by social ID column (apple_id/google_id),
// falling back to email match, and creates a new user if none exists.
func (h *H) findOrCreateSocialUser(r *http.Request, idCol, idVal, email, name string) (
	userID int64, userName, userEmail string, err error,
) {
	ctx := r.Context()

	// 1. Try to find by social ID.
	err = h.DB.QueryRowContext(ctx,
		fmt.Sprintf(`SELECT id, name, email FROM users WHERE %s = $1`, idCol), idVal,
	).Scan(&userID, &userName, &userEmail)
	if err == nil {
		return // found
	}

	// 2. Try to find by email and attach the social ID.
	if email != "" {
		err = h.DB.QueryRowContext(ctx,
			`SELECT id, name, email FROM users WHERE email = $1`, email,
		).Scan(&userID, &userName, &userEmail)
		if err == nil {
			_, _ = h.DB.ExecContext(ctx,
				fmt.Sprintf(`UPDATE users SET %s = $1 WHERE id = $2`, idCol), idVal, userID)
			return
		}
	}

	// 3. Create new user (no password, social-only account).
	if email == "" {
		email = fmt.Sprintf("%s@privaterelay.appleid.com", idVal)
	}
	err = h.DB.QueryRowContext(ctx,
		fmt.Sprintf(`INSERT INTO users (name, email, password_hash, %s)
		             VALUES ($1, $2, '', $3)
		             RETURNING id, name, email`, idCol),
		name, email, idVal,
	).Scan(&userID, &userName, &userEmail)
	return
}

// decodeAppleIdentityToken decodes the JWT payload from Apple's identity token
// and returns the Apple user ID (sub) and email.
// NOTE: For production, also verify the JWT signature using Apple's JWK endpoint at
// https://appleid.apple.com/auth/keys and validate iss, aud, and exp claims.
func decodeAppleIdentityToken(idToken string) (sub, email string, err error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", "", errors.New("malformed JWT")
	}
	payload := parts[1]
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}
	payload = strings.NewReplacer("-", "+", "_", "/").Replace(payload)
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", "", fmt.Errorf("base64 decode: %w", err)
	}
	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(data, &claims); err != nil {
		return "", "", fmt.Errorf("unmarshal claims: %w", err)
	}
	if claims.Sub == "" {
		return "", "", errors.New("missing sub claim")
	}
	return claims.Sub, claims.Email, nil
}

// verifyGoogleIDToken calls Google's tokeninfo endpoint to validate an ID token
// and extract the user's Google ID, email, and name.
func verifyGoogleIDToken(idToken string) (googleID, email, name string, err error) {
	resp, err := http.Get(
		"https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken,
	)
	if err != nil {
		return "", "", "", fmt.Errorf("tokeninfo request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("tokeninfo returned %d: %s", resp.StatusCode, body)
	}
	var info struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", "", "", fmt.Errorf("parse tokeninfo: %w", err)
	}
	if info.Sub == "" {
		return "", "", "", errors.New("missing sub in tokeninfo")
	}
	return info.Sub, info.Email, info.Name, nil
}
