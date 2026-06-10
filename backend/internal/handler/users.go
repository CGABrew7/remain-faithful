package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
	req.Name  = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	var id int64
	var name, email, createdAt string
	err := h.DB.QueryRowContext(r.Context(),
		`UPDATE users SET name = $1, email = $2 WHERE id = $3
		 RETURNING id, name, email, created_at`,
		req.Name, req.Email, userID,
	).Scan(&id, &name, &email, &createdAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         id,
		"name":       name,
		"email":      email,
		"created_at": createdAt,
	})
}
