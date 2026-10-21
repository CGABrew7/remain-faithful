package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

	// If the invitee already has an account, keep the relationship pending
	// until they accept. Never report "connected" for a pending row.
	var inviteeID int64
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM users WHERE email = $1`, req.Email,
	).Scan(&inviteeID); err == nil {
		if inviteeID == userID {
			writeError(w, http.StatusBadRequest, "cannot add yourself as a partner")
			return
		}
		if status, msg := h.preparePartnerInvite(r.Context(), userID, inviteeID); status != 0 {
			writeError(w, status, msg)
			return
		}
		var relID int64
		var status string
		if err := h.DB.QueryRowContext(r.Context(),
			`INSERT INTO relationships (user_id, partner_id, type, status)
			 VALUES ($1, $2, 'partner', 'pending')
			 ON CONFLICT (user_id, partner_id) DO UPDATE
			   SET type = EXCLUDED.type
			 RETURNING id, status`,
			userID, inviteeID,
		).Scan(&relID, &status); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create relationship")
			return
		}
		if status == "accepted" {
			writeJSON(w, http.StatusOK, map[string]any{
				"status":          "accepted",
				"relationship_id": relID,
			})
			return
		}
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
		acceptURL := siteBase() + "/accept-invite?token=" + token + "&type=partner"
		if h.Email != nil {
			if sendErr := h.Email.SendPartnerInvite(req.Email, inviterName, acceptURL); sendErr != nil {
				log.Printf("[invite] partner email error user=%d: %v", userID, sendErr)
			}
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"status":          "pending",
			"relationship_id": relID,
		})
		return
	}

	if status, msg := h.enforceInviteBudget(r.Context(), userID); status != 0 {
		writeError(w, status, msg)
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
	if h.Email != nil {
		if sendErr := h.Email.SendPartnerInvite(req.Email, inviterName, acceptURL); sendErr != nil {
			log.Printf("[invite] partner email error user=%d: %v", userID, sendErr)
		}
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
	userID, acceptingEmail, ok := h.requireVerifiedMailbox(w, r)
	if !ok {
		return
	}

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

	if status, msg := h.checkCanInviteToGroup(r.Context(), groupID, userID); status != 0 {
		writeError(w, status, msg)
		return
	}

	res, ierr := h.inviteEmailToGroup(r.Context(), userID, groupID, req.Email)
	if ierr != nil {
		writeError(w, ierr.status, ierr.msg)
		return
	}
	switch res.outcome {
	case groupInviteAdded:
		writeJSON(w, http.StatusCreated, map[string]any{
			"status":  "added",
			"user_id": res.inviteeID,
		})
	case groupInviteAlreadyMember:
		writeJSON(w, http.StatusOK, map[string]any{"status": "already_member"})
	default:
		writeJSON(w, http.StatusCreated, map[string]any{
			"status": "invited",
			"email":  req.Email,
		})
	}
}

// maxGroupMembers is the hard cap on accountability group size.
const maxGroupMembers = 12

type groupInviteOutcome string

const (
	groupInviteAdded         groupInviteOutcome = "added"
	groupInviteAlreadyMember groupInviteOutcome = "already_member"
	groupInviteInvited       groupInviteOutcome = "invited"
)

type groupInviteResult struct {
	outcome   groupInviteOutcome
	inviteeID int64
	joinedAt  string
}

type handlerError struct {
	status int
	msg    string
}

// checkCanInviteToGroup verifies the caller is a member of the group (any
// role) and that the group is below the member cap. It returns a zero status
// when the invite may proceed.
func (h *H) checkCanInviteToGroup(ctx context.Context, groupID, userID int64) (int, string) {
	var callerRole string
	if err := h.DB.QueryRowContext(ctx,
		`SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	).Scan(&callerRole); err != nil {
		return http.StatusForbidden, "you must be a member of this group to invite others"
	}

	var memberCount int
	_ = h.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM group_members WHERE group_id = $1`, groupID,
	).Scan(&memberCount)
	if memberCount >= maxGroupMembers {
		return http.StatusUnprocessableEntity, "group has reached the maximum of 12 members"
	}
	return 0, ""
}

// inviteEmailToGroup adds an existing user to the group, or stores a pending
// group invite for an email with no account yet. Either way the invitee gets
// an email. The caller must already have run checkCanInviteToGroup.
// email must already be trimmed and lowercased.
func (h *H) inviteEmailToGroup(ctx context.Context, inviterID, groupID int64, email string) (groupInviteResult, *handlerError) {
	var inviterName, groupName string
	if err := h.DB.QueryRowContext(ctx,
		`SELECT name FROM users WHERE id = $1`, inviterID,
	).Scan(&inviterName); err != nil {
		return groupInviteResult{}, &handlerError{http.StatusInternalServerError, "could not find your account"}
	}
	if err := h.DB.QueryRowContext(ctx,
		`SELECT name FROM groups WHERE id = $1`, groupID,
	).Scan(&groupName); err != nil {
		return groupInviteResult{}, &handlerError{http.StatusNotFound, "group not found"}
	}

	// If the invitee already has an account, store a pending invite.
	// They are not a member, and they do not see the roster, until they accept.
	var inviteeID int64
	if err := h.DB.QueryRowContext(ctx,
		`SELECT id FROM users WHERE email = $1`, email,
	).Scan(&inviteeID); err == nil {
		var existing int64
		err := h.DB.QueryRowContext(ctx,
			`SELECT user_id FROM group_members WHERE group_id = $1 AND user_id = $2`,
			groupID, inviteeID,
		).Scan(&existing)
		if err == nil {
			return groupInviteResult{outcome: groupInviteAlreadyMember, inviteeID: inviteeID}, nil
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return groupInviteResult{}, &handlerError{http.StatusInternalServerError, "failed to check membership"}
		}
		if status, msg := h.enforceInviteBudget(ctx, inviterID); status != 0 {
			return groupInviteResult{}, &handlerError{status, msg}
		}
		token, err := inviteToken()
		if err != nil {
			return groupInviteResult{}, &handlerError{http.StatusInternalServerError, "failed to generate invite token"}
		}
		if _, err := h.DB.ExecContext(ctx,
			`INSERT INTO group_invites (inviter_id, group_id, invitee_email, token)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (group_id, invitee_email) DO UPDATE
			   SET token = EXCLUDED.token, inviter_id = EXCLUDED.inviter_id, created_at = NOW(), status = 'pending'`,
			inviterID, groupID, email, token,
		); err != nil {
			return groupInviteResult{}, &handlerError{http.StatusInternalServerError, "failed to store group invite"}
		}
		acceptURL := fmt.Sprintf("%s/accept-invite?token=%s&type=group", siteBase(), token)
		if h.Email != nil {
			if sendErr := h.Email.SendGroupInvite(email, inviterName, groupName, acceptURL); sendErr != nil {
				log.Printf("[invite] group email error group=%d: %v", groupID, sendErr)
			}
		}
		return groupInviteResult{outcome: groupInviteInvited, inviteeID: inviteeID}, nil
	}

	if status, msg := h.enforceInviteBudget(ctx, inviterID); status != 0 {
		return groupInviteResult{}, &handlerError{status, msg}
	}

	// No account yet — store a pending invite and send an email. The invite is
	// redeemed when this email is verified (acceptPendingGroupInvites).
	token, err := inviteToken()
	if err != nil {
		return groupInviteResult{}, &handlerError{http.StatusInternalServerError, "failed to generate invite token"}
	}
	if _, err := h.DB.ExecContext(ctx,
		`INSERT INTO group_invites (inviter_id, group_id, invitee_email, token)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (group_id, invitee_email) DO UPDATE
		   SET token = EXCLUDED.token, created_at = NOW(), status = 'pending'`,
		inviterID, groupID, email, token,
	); err != nil {
		return groupInviteResult{}, &handlerError{http.StatusInternalServerError, "failed to store group invite"}
	}
	acceptURL := fmt.Sprintf("%s/invite?token=%s&type=group", siteBase(), token)
	if h.Email != nil {
		if sendErr := h.Email.SendGroupInvite(email, inviterName, groupName, acceptURL); sendErr != nil {
			log.Printf("[invite] group email error group=%d: %v", groupID, sendErr)
		}
	}
	return groupInviteResult{outcome: groupInviteInvited}, nil
}

// preparePartnerInvite enforces the partner cap for a new row and the daily
// invite email cap. An already-accepted partnership is left alone.
func (h *H) preparePartnerInvite(ctx context.Context, userID, inviteeID int64) (int, string) {
	var status string
	err := h.DB.QueryRowContext(ctx,
		`SELECT status FROM relationships WHERE user_id = $1 AND partner_id = $2`,
		userID, inviteeID,
	).Scan(&status)
	if err == nil && status == "accepted" {
		return 0, ""
	}
	if status != "pending" {
		n, cerr := h.partnerCount(ctx, userID)
		if cerr != nil {
			return http.StatusInternalServerError, "could not check partner limit"
		}
		if n >= maxPartnersPerUser {
			return http.StatusUnprocessableEntity, "you can have up to 10 partners"
		}
	}
	return h.enforceInviteBudget(ctx, userID)
}

// ListPendingPartnerInvites returns partner invites addressed to the signed-in email.
// GET /relationships/invites/pending
func (h *H) ListPendingPartnerInvites(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	email, verified, appleID, err := h.accountMailbox(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not verify your account")
		return
	}
	if !verified {
		writeUnverifiedMailbox(w, appleID, email)
		return
	}
	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT i.id, i.inviter_id, u.name, i.created_at
		FROM relationship_invites i
		JOIN users u ON u.id = i.inviter_id
		WHERE i.invitee_email = $1 AND i.status = 'pending'
		  AND i.created_at > NOW() - INTERVAL '7 days'
		ORDER BY i.created_at DESC
	`, email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list invites")
		return
	}
	defer rows.Close()
	type item struct {
		ID          int64  `json:"id"`
		InviterID   int64  `json:"inviter_id"`
		InviterName string `json:"inviter_name"`
		CreatedAt   string `json:"created_at"`
		Type        string `json:"type"`
	}
	out := []item{}
	for rows.Next() {
		var it item
		it.Type = "partner"
		if rows.Scan(&it.ID, &it.InviterID, &it.InviterName, &it.CreatedAt) == nil {
			out = append(out, it)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

// ListPendingGroupInvites returns group invites addressed to the signed-in email.
// GET /groups/invites/pending
func (h *H) ListPendingGroupInvites(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	email, verified, appleID, err := h.accountMailbox(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not verify your account")
		return
	}
	if !verified {
		writeUnverifiedMailbox(w, appleID, email)
		return
	}
	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT i.id, i.group_id, g.name, i.inviter_id, u.name, i.created_at
		FROM group_invites i
		JOIN groups g ON g.id = i.group_id
		JOIN users u ON u.id = i.inviter_id
		WHERE i.invitee_email = $1 AND i.status = 'pending'
		  AND i.created_at > NOW() - INTERVAL '7 days'
		ORDER BY i.created_at DESC
	`, email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list invites")
		return
	}
	defer rows.Close()
	type item struct {
		ID          int64  `json:"id"`
		GroupID     int64  `json:"group_id"`
		GroupName   string `json:"group_name"`
		InviterID   int64  `json:"inviter_id"`
		InviterName string `json:"inviter_name"`
		CreatedAt   string `json:"created_at"`
		Type        string `json:"type"`
	}
	out := []item{}
	for rows.Next() {
		var it item
		it.Type = "group"
		if rows.Scan(&it.ID, &it.GroupID, &it.GroupName, &it.InviterID, &it.InviterName, &it.CreatedAt) == nil {
			out = append(out, it)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *H) userEmail(ctx context.Context, userID int64) (string, error) {
	email, _, err := h.userMailbox(ctx, userID)
	return email, err
}

func (h *H) userMailbox(ctx context.Context, userID int64) (string, bool, error) {
	email, verified, _, err := h.accountMailbox(ctx, userID)
	return email, verified, err
}

func (h *H) accountMailbox(ctx context.Context, userID int64) (string, bool, string, error) {
	var email, appleID string
	var verified bool
	err := h.DB.QueryRowContext(ctx,
		`SELECT email, email_verified, COALESCE(apple_id, '') FROM users WHERE id = $1`,
		userID,
	).Scan(&email, &verified, &appleID)
	return email, verified, appleID, err
}

// requireVerifiedMailbox loads the caller's address and refuses the request
// until that address is confirmed. The JSON code tells the app whether to
// offer a confirmation resend or to ask the person to add an email.
func (h *H) requireVerifiedMailbox(w http.ResponseWriter, r *http.Request) (int64, string, bool) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	email, verified, appleID, err := h.accountMailbox(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not verify your account")
		return 0, "", false
	}
	if !verified {
		writeUnverifiedMailbox(w, appleID, email)
		return 0, "", false
	}
	return userID, email, true
}

// writeUnverifiedMailbox explains why invites are unavailable.
// A sign-in that never received an email address cannot confirm one.
func writeUnverifiedMailbox(w http.ResponseWriter, appleID, email string) {
	if placeholderAppleEmail(appleID, email) {
		writeCodedError(w, http.StatusForbidden, "email_required",
			"Add an email address in your account and confirm it before you can accept an invite.")
		return
	}
	writeCodedError(w, http.StatusForbidden, "email_not_verified",
		"Confirm your email, then you can view and accept invites. You can resend the confirmation email.")
}

// AcceptPartnerInviteByID accepts a pending partner invite addressed to the caller.
// POST /relationships/invites/{id}/accept
func (h *H) AcceptPartnerInviteByID(w http.ResponseWriter, r *http.Request) {
	h.decidePartnerInvite(w, r, true)
}

// DeclinePartnerInviteByID declines a pending partner invite.
// POST /relationships/invites/{id}/decline
func (h *H) DeclinePartnerInviteByID(w http.ResponseWriter, r *http.Request) {
	h.decidePartnerInvite(w, r, false)
}

func (h *H) decidePartnerInvite(w http.ResponseWriter, r *http.Request, accept bool) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	inviteID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid invite id")
		return
	}
	var email string
	if accept {
		var ok bool
		userID, email, ok = h.requireVerifiedMailbox(w, r)
		if !ok {
			return
		}
	} else {
		var err error
		email, err = h.userEmail(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not verify your account")
			return
		}
	}
	var token string
	var inviterID int64
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT token, inviter_id FROM relationship_invites
		WHERE id = $1 AND invitee_email = $2 AND status = 'pending'
		  AND created_at > NOW() - INTERVAL '7 days'
	`, inviteID, email).Scan(&token, &inviterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "invite not found or expired")
		return
	}
	if !accept {
		_, _ = h.DB.ExecContext(r.Context(),
			`UPDATE relationship_invites SET status = 'declined' WHERE id = $1`, inviteID)
		_, _ = h.DB.ExecContext(r.Context(),
			`DELETE FROM relationships WHERE user_id = $1 AND partner_id = $2 AND status = 'pending'`,
			inviterID, userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "declined"})
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
		writeError(w, http.StatusInternalServerError, "failed to accept invite")
		return
	}
	_, _ = h.DB.ExecContext(r.Context(),
		`UPDATE relationship_invites SET status = 'accepted' WHERE id = $1`, inviteID)
	writeJSON(w, http.StatusOK, map[string]any{"status": "accepted", "relationship_id": relID})
}

// DeclinePartnerInvite declines by the email token.
// POST /relationships/decline-invite
// Body: { "token": "..." }
func (h *H) DeclinePartnerInvite(w http.ResponseWriter, r *http.Request) {
	h.declineByToken(w, r, false)
}

// AcceptGroupInvite accepts a group invite token for the signed-in email.
// POST /groups/accept-invite
// Body: { "token": "..." }
func (h *H) AcceptGroupInvite(w http.ResponseWriter, r *http.Request) {
	h.decideGroupByToken(w, r, true)
}

// DeclineGroupInvite declines a group invite token.
// POST /groups/decline-invite
func (h *H) DeclineGroupInvite(w http.ResponseWriter, r *http.Request) {
	h.decideGroupByToken(w, r, false)
}

// AcceptGroupInviteByID accepts a listed group invite.
// POST /groups/invites/{id}/accept
func (h *H) AcceptGroupInviteByID(w http.ResponseWriter, r *http.Request) {
	h.decideGroupByID(w, r, true)
}

// DeclineGroupInviteByID declines a listed group invite.
// POST /groups/invites/{id}/decline
func (h *H) DeclineGroupInviteByID(w http.ResponseWriter, r *http.Request) {
	h.decideGroupByID(w, r, false)
}

func (h *H) decideGroupByID(w http.ResponseWriter, r *http.Request, accept bool) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	inviteID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid invite id")
		return
	}
	var email string
	if accept {
		var ok bool
		userID, email, ok = h.requireVerifiedMailbox(w, r)
		if !ok {
			return
		}
	} else {
		email, err = h.userEmail(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not verify your account")
			return
		}
	}
	var groupID int64
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT group_id FROM group_invites
		WHERE id = $1 AND invitee_email = $2 AND status = 'pending'
		  AND created_at > NOW() - INTERVAL '7 days'
	`, inviteID, email).Scan(&groupID)
	if err != nil {
		writeError(w, http.StatusNotFound, "invite not found or expired")
		return
	}
	if !accept {
		_, _ = h.DB.ExecContext(r.Context(), `UPDATE group_invites SET status = 'declined' WHERE id = $1`, inviteID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "declined"})
		return
	}
	if err := h.joinGroupFromInvite(r.Context(), userID, inviteID, groupID); err != nil {
		if errors.Is(err, errGroupFull) {
			writeError(w, http.StatusUnprocessableEntity, "group has reached the maximum of 12 members")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to join group")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "accepted", "group_id": groupID})
}

func (h *H) decideGroupByToken(w http.ResponseWriter, r *http.Request, accept bool) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	var email string
	var err error
	if accept {
		var ok bool
		userID, email, ok = h.requireVerifiedMailbox(w, r)
		if !ok {
			return
		}
	} else {
		email, err = h.userEmail(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not verify your account")
			return
		}
	}
	var inviteID, groupID int64
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT id, group_id FROM group_invites
		WHERE token = $1 AND invitee_email = $2 AND status = 'pending'
		  AND created_at > NOW() - INTERVAL '7 days'
	`, req.Token, email).Scan(&inviteID, &groupID)
	if err != nil {
		writeError(w, http.StatusNotFound, "invite not found or expired")
		return
	}
	if !accept {
		_, _ = h.DB.ExecContext(r.Context(), `UPDATE group_invites SET status = 'declined' WHERE id = $1`, inviteID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "declined"})
		return
	}
	if err := h.joinGroupFromInvite(r.Context(), userID, inviteID, groupID); err != nil {
		if errors.Is(err, errGroupFull) {
			writeError(w, http.StatusUnprocessableEntity, "group has reached the maximum of 12 members")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to join group")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "accepted", "group_id": groupID})
}

func (h *H) declineByToken(w http.ResponseWriter, r *http.Request, _ bool) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	email, err := h.userEmail(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not verify your account")
		return
	}
	var inviteID, inviterID int64
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT id, inviter_id FROM relationship_invites
		WHERE token = $1 AND invitee_email = $2 AND status = 'pending'
		  AND created_at > NOW() - INTERVAL '7 days'
	`, req.Token, email).Scan(&inviteID, &inviterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "invite not found or expired")
		return
	}
	_, _ = h.DB.ExecContext(r.Context(), `UPDATE relationship_invites SET status = 'declined' WHERE id = $1`, inviteID)
	_, _ = h.DB.ExecContext(r.Context(),
		`DELETE FROM relationships WHERE user_id = $1 AND partner_id = $2 AND status = 'pending'`,
		inviterID, userID)
	writeJSON(w, http.StatusOK, map[string]string{"status": "declined"})
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
	return "https://www.remainfaithful.com"
}
