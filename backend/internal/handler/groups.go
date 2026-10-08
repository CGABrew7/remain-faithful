package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

// ListMyGroups returns all groups the authenticated user belongs to.
// GET /groups
func (h *H) ListMyGroups(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT g.id, g.name, g.created_at
		FROM   groups g
		JOIN   group_members gm ON gm.group_id = g.id
		WHERE  gm.user_id = $1
		ORDER  BY gm.joined_at ASC
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch groups")
		return
	}
	defer rows.Close()

	type group struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
	}

	groups := []group{}
	for rows.Next() {
		var g group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			continue
		}
		groups = append(groups, g)
	}

	writeJSON(w, http.StatusOK, groups)
}

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
// Caller must be in group_members; otherwise member PII is not returned.
func (h *H) GetGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	// Verify caller is a member (any role) before returning emails/names/flags/streaks.
	var callerRole string
	if err = h.DB.QueryRowContext(r.Context(),
		`SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	).Scan(&callerRole); err != nil {
		writeError(w, http.StatusForbidden, "you must be a member of this group")
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
		SELECT gm.user_id, gm.role, gm.joined_at, u.name, u.email,
		       stats.flags_last_30, stats.streak_days
		FROM   group_members gm
		JOIN   users u ON u.id = gm.user_id
		JOIN   LATERAL (
		    SELECT
		        COUNT(*) FILTER (WHERE e.timestamp > NOW() - INTERVAL '30 days') AS flags_last_30,
		        COALESCE(
		            (NOW()::date - MAX(e.timestamp)::date),
		            (NOW()::date - gm.joined_at::date)
		        )                                                                  AS streak_days
		    FROM events e
		    WHERE e.user_id = gm.user_id
		) stats ON TRUE
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
		UserID      int64    `json:"user_id"`
		Role        string   `json:"role"`
		JoinedAt    string   `json:"joined_at"`
		User        userInfo `json:"user"`
		FlagsLast30 int      `json:"flags_last_30"`
		StreakDays  int      `json:"streak_days"`
	}

	members := []member{}
	for rows.Next() {
		var m member
		if err := rows.Scan(&m.UserID, &m.Role, &m.JoinedAt, &m.User.Name, &m.User.Email,
			&m.FlagsLast30, &m.StreakDays); err != nil {
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

// InviteMember invites someone to an existing group by email. Any group
// member may invite. If the email already has an account the user is added
// immediately; otherwise a pending invite is stored and emailed, and the person
// joins automatically when they sign up with that email.
// POST /groups/:id/invite
// Body: { "user_email": "..." } (the key "email" is also accepted)
func (h *H) InviteMember(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	if status, msg := h.checkCanInviteToGroup(r.Context(), groupID, userID); status != 0 {
		writeError(w, status, msg)
		return
	}

	var req struct {
		UserEmail string `json:"user_email"`
		Email     string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	email := req.UserEmail
	if strings.TrimSpace(email) == "" {
		email = req.Email
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		writeError(w, http.StatusBadRequest, "user_email is required")
		return
	}

	res, ierr := h.inviteEmailToGroup(r.Context(), userID, groupID, email)
	if ierr != nil {
		writeError(w, ierr.status, ierr.msg)
		return
	}
	switch res.outcome {
	case groupInviteAlreadyMember:
		writeError(w, http.StatusConflict, "user is already a member of this group")
	case groupInviteAdded:
		writeJSON(w, http.StatusCreated, map[string]any{
			"status":    "added",
			"group_id":  groupID,
			"user_id":   res.inviteeID,
			"role":      "member",
			"joined_at": res.joinedAt,
		})
	default:
		writeJSON(w, http.StatusCreated, map[string]any{
			"status":   "invited",
			"group_id": groupID,
			"email":    email,
		})
	}
}

// LeaveGroup removes the authenticated user from a specific group.
// DELETE /groups/{id}/members/me
func (h *H) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	res, err := h.DB.ExecContext(r.Context(),
		`DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to leave group")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "not a member of this group")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// LeaveAllGroups removes the authenticated user from every group they belong to.
// POST /groups/leave-all
func (h *H) LeaveAllGroups(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	if _, err := h.DB.ExecContext(r.Context(),
		`DELETE FROM group_members WHERE user_id = $1`, userID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to leave groups")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
