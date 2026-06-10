package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

// InvitePartner sends a partner invitation email.
// POST /relationships/invite
// Body: { "email": "..." }
// If the invitee already has an account the relationship is created immediately.
// Otherwise a pending invite is stored and an email is sent with a deep link.
func (h *H) InvitePartner(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		writeError(w, http.StatusBadRequest, "valid email is required")
		return
	}

	var inviterName string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&inviterName); err != nil {
		writeError(w, http.StatusInternalServerError, "could not find your account")
		return
	}

	// If the invitee already has an account, create the relationship directly.
	var inviteeID int64
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM users WHERE email = $1`, req.Email,
	).Scan(&inviteeID); err == nil {
		if inviteeID == userID {
			writeError(w, http.StatusBadRequest, "cannot add yourself as a partner")
			return
		}
		var relID int64
		var status, createdAt string
		if err := h.DB.QueryRowContext(r.Context(),
			`INSERT INTO relationships (user_id, partner_id, type)
			 VALUES ($1, $2, 'partner')
			 ON CONFLICT (user_id, partner_id) DO UPDATE SET type = EXCLUDED.type
			 RETURNING id, status, created_at`,
			userID, inviteeID,
		).Scan(&relID, &status, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create relationship")
			return
		}
		token, _ := inviteToken()
		acceptURL := siteBase() + "/accept-invite?token=" + token + "&type=partner"
		_ = h.Email.SendPartnerInvite(req.Email, inviterName, acceptURL)
		writeJSON(w, http.StatusCreated, map[string]any{
			"status":          "connected",
			"relationship_id": relID,
		})
		return
	}

	// Invitee has no account — store a pending invite and send an email.
	token, err := inviteToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate invite token")
		return
	}
	if _, err := h.DB.ExecContext(r.Context(),
		`INSERT INTO relationship_invites (inviter_id, invitee_email, token)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (inviter_id, invitee_email) DO UPDATE
		   SET token = EXCLUDED.token, created_at = NOW(), status = 'pending'`,
		userID, req.Email, token,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store invite")
		return
	}
	acceptURL := siteBase() + "/invite?token=" + token + "&type=partner"
	if sendErr := h.Email.SendPartnerInvite(req.Email, inviterName, acceptURL); sendErr != nil {
		fmt.Printf("[invite] email error: %v\n", sendErr)
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"status": "invited",
		"email":  req.Email,
	})
}

// AcceptPartnerInvite accepts a pending partner invite by token.
// POST /relationships/accept-invite
// Body: { "token": "..." }
func (h *H) AcceptPartnerInvite(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	var inviterID int64
	var inviteeEmail string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT inviter_id, invitee_email FROM relationship_invites
		 WHERE token = $1 AND status = 'pending'
		   AND created_at > NOW() - INTERVAL '7 days'`,
		req.Token,
	).Scan(&inviterID, &inviteeEmail); err != nil {
		writeError(w, http.StatusNotFound, "invite not found or expired")
		return
	}

	var acceptingEmail string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT email FROM users WHERE id = $1`, userID,
	).Scan(&acceptingEmail); err != nil {
		writeError(w, http.StatusInternalServerError, "could not verify your account")
		return
	}
	if acceptingEmail != inviteeEmail {
		writeError(w, http.StatusForbidden, "this invite was sent to a different email address")
		return
	}

	var relID int64
	if err := h.DB.QueryRowContext(r.Context(),
		`INSERT INTO relationships (user_id, partner_id, type, status)
		 VALUES ($1, $2, 'partner', 'accepted')
		 ON CONFLICT (user_id, partner_id) DO UPDATE SET status = 'accepted'
		 RETURNING id`,
		inviterID, userID,
	).Scan(&relID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create relationship")
		return
	}
	_, _ = h.DB.ExecContext(r.Context(),
		`UPDATE relationship_invites SET status = 'accepted' WHERE token = $1`, req.Token)

	writeJSON(w, http.StatusCreated, map[string]any{
		"relationship_id": relID,
		"status":          "accepted",
	})
}

// GroupEmailInvite sends an email invitation to join a group to someone
// who may not yet have an account.
// POST /groups/{id}/email-invite
// Body: { "email": "..." }
func (h *H) GroupEmailInvite(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	groupID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	// Verify caller is a member (any role).
	var callerRole string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	).Scan(&callerRole); err != nil {
		writeError(w, http.StatusForbidden, "you must be a member of this group to invite others")
		return
	}

	// Enforce group size cap.
	var memberCount int
	_ = h.DB.QueryRowContext(r.Context(),
		`SELECT COUNT(*) FROM group_members WHERE group_id = $1`, groupID,
	).Scan(&memberCount)
	if memberCount >= 12 {
		writeError(w, http.StatusUnprocessableEntity, "group has reached the maximum of 12 members")
		return
	}

	var inviterName, groupName string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&inviterName); err != nil {
		writeError(w, http.StatusInternalServerError, "could not find your account")
		return
	}
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM groups WHERE id = $1`, groupID,
	).Scan(&groupName); err != nil {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}

	// If the invitee already has an account, add them directly.
	var inviteeID int64
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM users WHERE email = $1`, req.Email,
	).Scan(&inviteeID); err == nil {
		var joinedAt string
		if err := h.DB.QueryRowContext(r.Context(),
			`INSERT INTO group_members (group_id, user_id, role)
			 VALUES ($1, $2, 'member')
			 ON CONFLICT (group_id, user_id) DO NOTHING
			 RETURNING joined_at`,
			groupID, inviteeID,
		).Scan(&joinedAt); err == nil {
			acceptURL := fmt.Sprintf("%s/groups/%d", siteBase(), groupID)
			_ = h.Email.SendGroupInvite(req.Email, inviterName, groupName, acceptURL)
			writeJSON(w, http.StatusCreated, map[string]any{
				"status":  "added",
				"user_id": inviteeID,
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "already_member"})
		return
	}

	// No account yet — store a pending invite and send an email.
	token, err := inviteToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate invite token")
		return
	}
	if _, err := h.DB.ExecContext(r.Context(),
		`INSERT INTO group_invites (inviter_id, group_id, invitee_email, token)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (group_id, invitee_email) DO UPDATE
		   SET token = EXCLUDED.token, created_at = NOW(), status = 'pending'`,
		userID, groupID, req.Email, token,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store group invite")
		return
	}
	acceptURL := fmt.Sprintf("%s/invite?token=%s&type=group", siteBase(), token)
	if sendErr := h.Email.SendGroupInvite(req.Email, inviterName, groupName, acceptURL); sendErr != nil {
		fmt.Printf("[invite] group email error: %v\n", sendErr)
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"status": "invited",
		"email":  req.Email,
	})
}

// inviteToken generates a cryptographically random 32-byte hex token.
func inviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// siteBase returns the website base URL from the SITE_URL env var.
func siteBase() string {
	if u := os.Getenv("SITE_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://remainfaithful.app"
}
