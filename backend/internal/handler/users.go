package handler

import (
	"net/http"

	rfauth "remain-faithful/backend/internal/auth"
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
	var name, email, createdAt string
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name, email, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&id, &name, &email, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         id,
		"name":       name,
		"email":      email,
		"created_at": createdAt,
	})
}
