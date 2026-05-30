package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         id,
		"name":       req.Name,
		"email":      req.Email,
		"created_at": createdAt,
	})
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
