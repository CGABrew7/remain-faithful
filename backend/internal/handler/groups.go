package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

// CreateGroup creates a new accountability group and adds the creator as admin.
// POST /groups
// Body: { "name": "...", "covenant": "..." }
func (h *H) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Name     string `json:"name"`
		Covenant string `json:"covenant"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to begin transaction")
		return
	}
	defer tx.Rollback()

	var groupID int64
	var createdAt string
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO groups (name, covenant) VALUES ($1, $2) RETURNING id, created_at`,
		req.Name, req.Covenant,
	).Scan(&groupID, &createdAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create group")
		return
	}

	if _, err = tx.ExecContext(r.Context(),
		`INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'admin')`,
		groupID, userID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add creator as admin")
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         groupID,
		"name":       req.Name,
		"covenant":   req.Covenant,
		"created_at": createdAt,
	})
}

// GetGroup returns a group and its full member list.
// GET /groups/:id
func (h *H) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	var name, createdAt string
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT name, created_at FROM groups WHERE id = $1`,
		groupID,
	).Scan(&name, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}

	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT gm.user_id, gm.role, gm.joined_at, u.name, u.email
		FROM   group_members gm
		JOIN   users u ON u.id = gm.user_id
		WHERE  gm.group_id = $1
		ORDER  BY gm.joined_at ASC
	`, groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch members")
		return
	}
	defer rows.Close()

	type userInfo struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type member struct {
		UserID   int64    `json:"user_id"`
		Role     string   `json:"role"`
		JoinedAt string   `json:"joined_at"`
		User     userInfo `json:"user"`
	}

	members := []member{}
	for rows.Next() {
		var m member
		if err := rows.Scan(&m.UserID, &m.Role, &m.JoinedAt, &m.User.Name, &m.User.Email); err != nil {
			continue
		}
		m.User.ID = m.UserID
		members = append(members, m)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         groupID,
		"name":       name,
		"created_at": createdAt,
		"members":    members,
	})
}

// InviteMember adds a user to an existing group. Caller must be a group admin.
// POST /groups/:id/invite
// Body: { "user_email": "..." }
func (h *H) InviteMember(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	var callerRole string
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	).Scan(&callerRole)
	if err != nil || callerRole != "admin" {
		writeError(w, http.StatusForbidden, "only group admins can invite members")
		return
	}

	var req struct {
		UserEmail string `json:"user_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.UserEmail = strings.ToLower(strings.TrimSpace(req.UserEmail))
	if req.UserEmail == "" {
		writeError(w, http.StatusBadRequest, "user_email is required")
		return
	}

	var inviteeID int64
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM users WHERE email = $1`,
		req.UserEmail,
	).Scan(&inviteeID)
	if err != nil {
		writeError(w, http.StatusNotFound, "no account found for that email")
		return
	}

	var joinedAt string
	err = h.DB.QueryRowContext(r.Context(),
		`INSERT INTO group_members (group_id, user_id, role)
		 VALUES ($1, $2, 'member')
		 ON CONFLICT (group_id, user_id) DO NOTHING
		 RETURNING joined_at`,
		groupID, inviteeID,
	).Scan(&joinedAt)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusConflict, "user is already a member of this group")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add member")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"group_id":  groupID,
		"user_id":   inviteeID,
		"role":      "member",
		"joined_at": joinedAt,
	})
}
