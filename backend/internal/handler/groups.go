package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

	var name, covenant, createdAt string
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT name, covenant, created_at FROM groups WHERE id = $1`,
		groupID,
	).Scan(&name, &covenant, &createdAt)
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
		"covenant":   covenant,
		"created_at": createdAt,
		"members":    members,
	})
}

const (
	maxGroupNameLen     = 80
	maxGroupCovenantLen = 8000
)

// UpdateGroup renames a group and/or saves its covenant. Any member may do
// this. A covenant change notifies the other members with a fixed line. The
// covenant text itself is not put in the push.
// PATCH /groups/{id}
// Body: { "name"?: "...", "covenant"?: "..." }
func (h *H) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Covenant *string `json:"covenant"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == nil && req.Covenant == nil {
		writeError(w, http.StatusBadRequest, "name or covenant is required")
		return
	}

	var role string
	if err = h.DB.QueryRowContext(r.Context(),
		`SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	).Scan(&role); err != nil {
		writeError(w, http.StatusForbidden, "you must be a member of this group")
		return
	}
	_ = role

	var currentName, currentCovenant string
	if err = h.DB.QueryRowContext(r.Context(),
		`SELECT name, covenant FROM groups WHERE id = $1`,
		groupID,
	).Scan(&currentName, &currentCovenant); err != nil {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}

	newName := currentName
	newCovenant := currentCovenant
	if req.Name != nil {
		newName = strings.TrimSpace(*req.Name)
		if newName == "" || len([]rune(newName)) > maxGroupNameLen {
			writeError(w, http.StatusBadRequest, "name must be 1–80 characters")
			return
		}
	}
	if req.Covenant != nil {
		newCovenant = strings.TrimSpace(*req.Covenant)
		if len([]rune(newCovenant)) > maxGroupCovenantLen {
			writeError(w, http.StatusBadRequest, "covenant is too long")
			return
		}
	}

	if _, err = h.DB.ExecContext(r.Context(),
		`UPDATE groups SET name = $1, covenant = $2 WHERE id = $3`,
		newName, newCovenant, groupID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update group")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":       groupID,
		"name":     newName,
		"covenant": newCovenant,
	})

	if newCovenant == currentCovenant {
		return
	}
	senderName, nameErr := h.lookupUserName(r.Context(), userID)
	if nameErr != nil {
		return
	}
	others, memErr := h.otherMemberIDs(r.Context(), groupID, userID)
	if memErr != nil || len(others) == 0 {
		return
	}
	notice := metadataPush(
		"Covenant Updated",
		covenantUpdatedBody(senderName),
		"COVENANT_UPDATED",
		senderName,
		fmt.Sprintf("covenant-%d", groupID),
	)
	h.inBackground(func() {
		h.notifyUserIDs(context.Background(), others, senderName, notice)
	})
}

// SendEncouragement pushes a fixed, metadata-only note to one other member
// of the group. The sender does not type the message.
// POST /groups/{id}/encouragement
// Body: { "user_id": 123 }
func (h *H) SendEncouragement(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == 0 {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	if req.UserID == userID {
		writeError(w, http.StatusBadRequest, "pick another member")
		return
	}

	var both int
	if err = h.DB.QueryRowContext(r.Context(), `
		SELECT COUNT(*) FROM group_members
		WHERE  group_id = $1 AND user_id IN ($2, $3)
	`, groupID, userID, req.UserID).Scan(&both); err != nil || both != 2 {
		writeError(w, http.StatusForbidden, "you can only encourage a member of this group")
		return
	}

	senderName, err := h.lookupUserName(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})

	notice := metadataPush(
		"Encouragement",
		encouragementBody(senderName),
		"ENCOURAGEMENT",
		senderName,
		fmt.Sprintf("encouragement-%d-%d", userID, req.UserID),
	)
	target := req.UserID
	h.inBackground(func() {
		h.notifyUserIDs(context.Background(), []int64{target}, senderName, notice)
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

// LeaveGroup removes the authenticated user from a specific group and notifies
// the people who remain. The notice is a fixed line with the leaver's name.
// DELETE /groups/{id}/members/me
func (h *H) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	leaverName, _ := h.lookupUserName(r.Context(), userID)

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

	others, _ := h.otherMemberIDs(r.Context(), groupID, userID)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if leaverName == "" || len(others) == 0 {
		return
	}
	notice := metadataPush(
		"Group Update",
		groupLeftBody(leaverName),
		"GROUP_MEMBER_LEFT",
		leaverName,
		fmt.Sprintf("left-group-%d-%d", groupID, userID),
	)
	h.inBackground(func() {
		h.notifyUserIDs(context.Background(), others, leaverName, notice)
	})
}

// LeaveAllGroups removes the authenticated user from every group they belong to
// and notifies the other members once each.
// POST /groups/leave-all
func (h *H) LeaveAllGroups(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	leaverName, _ := h.lookupUserName(r.Context(), userID)
	others, _ := h.coMemberIDs(r.Context(), userID)

	if _, err := h.DB.ExecContext(r.Context(),
		`DELETE FROM group_members WHERE user_id = $1`, userID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to leave groups")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if leaverName == "" || len(others) == 0 {
		return
	}
	notice := metadataPush(
		"Group Update",
		groupLeftBody(leaverName),
		"GROUP_MEMBER_LEFT",
		leaverName,
		fmt.Sprintf("left-all-groups-%d", userID),
	)
	h.inBackground(func() {
		h.notifyUserIDs(context.Background(), others, leaverName, notice)
	})
}
