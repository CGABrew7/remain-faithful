package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	rfauth "remain-faithful/backend/internal/auth"
)

// CreateRelationship sends a partner request to another user by email.
// POST /relationships
// Body: { "partner_email": "...", "type": "partner" }
func (h *H) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		PartnerEmail string `json:"partner_email"`
		Type         string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.PartnerEmail = strings.ToLower(strings.TrimSpace(req.PartnerEmail))
	if req.PartnerEmail == "" {
		writeError(w, http.StatusBadRequest, "partner_email is required")
		return
	}
	if req.Type == "" {
		req.Type = "partner"
	}

	var partnerID int64
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM users WHERE email = $1`,
		req.PartnerEmail,
	).Scan(&partnerID)
	if err != nil {
		writeError(w, http.StatusNotFound, "no account found for that email")
		return
	}
	if partnerID == userID {
		writeError(w, http.StatusBadRequest, "cannot add yourself as a partner")
		return
	}

	var id int64
	var status, createdAt string
	err = h.DB.QueryRowContext(r.Context(),
		`INSERT INTO relationships (user_id, partner_id, type)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, partner_id) DO UPDATE SET type = EXCLUDED.type
		 RETURNING id, status, created_at`,
		userID, partnerID, req.Type,
	).Scan(&id, &status, &createdAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create relationship")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         id,
		"user_id":    userID,
		"partner_id": partnerID,
		"type":       req.Type,
		"status":     status,
		"created_at": createdAt,
	})
}

// ListRelationships returns all relationships the authenticated user has.
// GET /relationships
func (h *H) ListRelationships(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT r.id, r.user_id, r.partner_id, r.type, r.status, r.created_at,
		       r.is_primary, u.name, u.email
		FROM   relationships r
		JOIN   users u ON u.id = r.partner_id
		WHERE  r.user_id = $1
		ORDER  BY r.created_at DESC
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list relationships")
		return
	}
	defer rows.Close()

	type partnerInfo struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type rel struct {
		ID        int64       `json:"id"`
		UserID    int64       `json:"user_id"`
		PartnerID int64       `json:"partner_id"`
		Type      string      `json:"type"`
		Status    string      `json:"status"`
		CreatedAt string      `json:"created_at"`
		IsPrimary bool        `json:"is_primary"`
		Partner   partnerInfo `json:"partner"`
	}

	result := []rel{}
	for rows.Next() {
		var item rel
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.PartnerID,
			&item.Type, &item.Status, &item.CreatedAt,
			&item.IsPrimary,
			&item.Partner.Name, &item.Partner.Email,
		); err != nil {
			continue
		}
		item.Partner.ID = item.PartnerID
		result = append(result, item)
	}
	writeJSON(w, http.StatusOK, result)
}

// SetPrimaryPartner marks one relationship as the user's primary partner
// (only one may be primary at a time) and clears all others.
// PUT /relationships/{id}/primary
func (h *H) SetPrimaryPartner(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid relationship id")
		return
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to start transaction")
		return
	}
	defer tx.Rollback()

	// Clear all existing primary flags for this user.
	if _, err = tx.ExecContext(r.Context(),
		`UPDATE relationships SET is_primary = FALSE WHERE user_id = $1`, userID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update relationships")
		return
	}

	// Set the target relationship as primary.
	res, err := tx.ExecContext(r.Context(),
		`UPDATE relationships SET is_primary = TRUE WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to set primary partner")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(w, http.StatusNotFound, "relationship not found")
		return
	}

	if err = tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to commit")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DeleteRelationship removes an accountability partnership.
// DELETE /relationships/{id}
func (h *H) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid relationship id")
		return
	}

	res, err := h.DB.ExecContext(r.Context(),
		`DELETE FROM relationships WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove partner")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "relationship not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
