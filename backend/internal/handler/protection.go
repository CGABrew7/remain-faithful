package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/apns"

	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
)

// Per-user rate-limit state for POST /protection/pin/verify.
// Allows 5 attempts per 15-minute window. Keyed by authenticated userID so
// the limit tracks the monitored user (who enters the PIN), not the IP.
var (
	pinMu      sync.Mutex
	pinBuckets = make(map[int64]*pinBucket)
)

type pinBucket struct {
	attempts int
	reset    time.Time
}

func init() {
	// Prune expired PIN rate-limit buckets every 15 minutes.
	go func() {
		t := time.NewTicker(15 * time.Minute)
		defer t.Stop()
		for range t.C {
			pinMu.Lock()
			now := time.Now()
			for uid, b := range pinBuckets {
				if now.After(b.reset) {
					delete(pinBuckets, uid)
				}
			}
			pinMu.Unlock()
		}
	}()
}

// pinVerifyRateLimit returns true (denied) when userID has exceeded 5 attempts
// in the current 15-minute window. Always increments before comparing.
func pinVerifyRateLimit(userID int64) bool {
	pinMu.Lock()
	defer pinMu.Unlock()
	now := time.Now()
	b, ok := pinBuckets[userID]
	if !ok || now.After(b.reset) {
		pinBuckets[userID] = &pinBucket{attempts: 1, reset: now.Add(15 * time.Minute)}
		return false
	}
	b.attempts++
	return b.attempts > 5
}

// SendProtectionAlert pushes a metadata-only PROTECTION_ALERT to all accepted
// partners of the authenticated user. No screen content is ever included.
// POST /protection/alerts
// Body: { "type": "...", "detail": "..." }
func (h *H) SendProtectionAlert(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Type   string `json:"type"`
		Detail string `json:"detail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Type == "" {
		writeError(w, http.StatusBadRequest, "type is required")
		return
	}

	var senderName string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&senderName); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	// Respond immediately; push happens in background.
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})

	go h.sendProtectionAlertToPartners(userID, senderName, req.Type, req.Detail)
}

// sendProtectionAlertToPartners fans out a PROTECTION_ALERT push to all
// accepted partners of userID. Only metadata is sent — never screen content.
func (h *H) sendProtectionAlertToPartners(userID int64, senderName, alertType, detail string) {
	ctx := context.Background()

	rows, err := h.DB.QueryContext(ctx, `
		SELECT partner_id FROM relationships
		WHERE user_id = $1 AND status = 'accepted'
	`, userID)
	if err != nil {
		log.Printf("[protection] query partners user=%d: %v", userID, err)
		return
	}
	defer rows.Close()

	var partnerIDs []int64
	for rows.Next() {
		var pid int64
		if rows.Scan(&pid) == nil {
			partnerIDs = append(partnerIDs, pid)
		}
	}
	if len(partnerIDs) == 0 {
		return
	}

	body := protectionAlertBody(senderName, alertType, detail)
	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Protection Update",
				"body":  body,
			},
			"sound": "default",
		},
		"notification_type": "PROTECTION_ALERT",
		"alert_type":        alertType,
		"sender_name":       senderName,
	}
	n := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: fmt.Sprintf("protection-%d-%s", userID, alertType),
		Payload:    payload,
	}
	for _, pid := range partnerIDs {
		h.notifyPartnerByID(ctx, pid, senderName, n)
	}
}

// protectionAlertBody returns a human-readable notification body for the given
// alert type. The detail parameter is a safe fallback for unknown types.
// Screen content is never included — only metadata.
func protectionAlertBody(name, alertType, detail string) string {
	switch alertType {
	case "monitoring_disabled":
		return name + " turned off activity monitoring"
	case "shielding_disabled":
		return name + " turned off app blocking"
	case "lockout_disabled":
		return name + " disabled App Lockout"
	case "lockout_broadcast_stopped":
		return name + " stopped Deep Scan — apps are re-shielded"
	case "pin_wrong_attempt":
		return name + " entered an incorrect Partner PIN"
	case "pin_removed":
		return name + "'s protection PIN was removed"
	case "pin_changed":
		return name + "'s protection PIN was changed"
	case "family_controls_revoked":
		return name + " revoked Screen Time authorization"
	case "heartbeat_silence":
		return name + "'s device has stopped sending heartbeats"
	default:
		if detail != "" {
			return name + ": " + detail
		}
		return name + " changed a protection setting"
	}
}

// GetPINStatus reports whether any accepted relationship for the authenticated
// user has a partner PIN configured.
// GET /protection/pin
// Response: { "is_set": bool }
func (h *H) GetPINStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var count int
	if err := h.DB.QueryRowContext(r.Context(), `
		SELECT COUNT(*) FROM relationships
		WHERE user_id = $1 AND status = 'accepted' AND pin_hash IS NOT NULL
	`, userID).Scan(&count); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query PIN status")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"is_set": count > 0})
}

// VerifyPIN bcrypt-compares the submitted PIN against the stored hash for the
// authenticated user's accepted relationship.
// Rate-limited: 5 attempts per 15-minute window per user.
// POST /protection/pin/verify
// Body: { "pin": "1234" }
// Response: { "success": bool }
func (h *H) VerifyPIN(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	if pinVerifyRateLimit(userID) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "900")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"error":"too many PIN attempts — try again in 15 minutes","success":false}`)
		return
	}

	var req struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PIN == "" {
		writeError(w, http.StatusBadRequest, "pin is required")
		return
	}

	// Fetch the PIN hash. LIMIT 1 — in practice one accepted relationship has a PIN.
	var pinHash string
	err := h.DB.QueryRowContext(r.Context(), `
		SELECT pin_hash FROM relationships
		WHERE user_id = $1 AND status = 'accepted' AND pin_hash IS NOT NULL
		LIMIT 1
	`, userID).Scan(&pinHash)
	if err != nil {
		// No PIN set — nothing to compare against.
		writeJSON(w, http.StatusOK, map[string]bool{"success": false})
		return
	}

	success := bcrypt.CompareHashAndPassword([]byte(pinHash), []byte(req.PIN)) == nil
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

// SetRelationshipPIN stores a bcrypt hash of the partner's chosen PIN.
// Only the accountability partner (partner_id == caller) may set it.
// POST /relationships/{id}/pin
// Body: { "pin": "1234" }
func (h *H) SetRelationshipPIN(w http.ResponseWriter, r *http.Request) {
	callerID, _ := rfauth.UserIDFromContext(r.Context())

	relID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid relationship id")
		return
	}

	var req struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// Enforce exactly 4 decimal digits — never log or store the plaintext.
	if len(req.PIN) != 4 {
		writeError(w, http.StatusBadRequest, "pin must be exactly 4 digits")
		return
	}
	for _, c := range req.PIN {
		if c < '0' || c > '9' {
			writeError(w, http.StatusBadRequest, "pin must contain only digits")
			return
		}
	}

	// Authorization: caller must be the partner_id on an accepted relationship.
	var monitoredUserID int64
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT user_id FROM relationships
		WHERE id = $1 AND partner_id = $2 AND status = 'accepted'
	`, relID, callerID).Scan(&monitoredUserID)
	if err != nil {
		writeError(w, http.StatusForbidden, "not authorized to set PIN on this relationship")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash PIN")
		return
	}

	if _, err = h.DB.ExecContext(r.Context(),
		`UPDATE relationships SET pin_hash = $1 WHERE id = $2`,
		string(hash), relID,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save PIN")
		return
	}

	log.Printf("[protection] PIN set relationship=%d partner=%d monitored=%d", relID, callerID, monitoredUserID)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DeleteRelationshipPIN removes the partner PIN from a relationship.
// Only the accountability partner (partner_id == caller) may remove it.
// DELETE /relationships/{id}/pin
func (h *H) DeleteRelationshipPIN(w http.ResponseWriter, r *http.Request) {
	callerID, _ := rfauth.UserIDFromContext(r.Context())

	relID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid relationship id")
		return
	}

	// Authorization enforced by the WHERE clause — only the partner_id can clear it.
	res, err := h.DB.ExecContext(r.Context(), `
		UPDATE relationships SET pin_hash = NULL
		WHERE id = $1 AND partner_id = $2 AND status = 'accepted'
	`, relID, callerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove PIN")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusForbidden, "not authorized or relationship not found")
		return
	}

	log.Printf("[protection] PIN removed relationship=%d partner=%d", relID, callerID)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
