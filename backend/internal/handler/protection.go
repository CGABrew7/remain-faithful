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

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

// clientProtectionAlertTypes are the type strings the iOS app / broadcast
// extension actually POST. Unknown and server-only types are rejected so a
// captured token cannot spam partners or forge heartbeat_silence.
//
// lockout_broadcast_stopped / pin_wrong_attempt are legacy protectionAlertBody
// names only — iOS sends deep_scan_stopped / wrong_pin_attempt. They are not
// accepted on HTTP POST.
var clientProtectionAlertTypes = map[string]struct{}{
	"monitoring_disabled":     {}, // SettingsView
	"shielding_disabled":      {}, // ActivitySelectionManager
	"lockout_disabled":        {}, // SettingsView
	"deep_scan_stopped":       {}, // AppLockoutManager, broadcast SampleHandler
	"wrong_pin_attempt":       {}, // PartnerPINManager
	"pin_removed":             {}, // SettingsView
	"pin_changed":             {}, // PINEntryView
	"family_controls_revoked": {}, // RemainFaithfulApp
	"signed_out":              {}, // SettingsView, build 10
	"selection_cleared":       {}, // ActivitySelectionManager, all chosen apps removed
}

// serverOnlyProtectionAlertTypes are generated internally (heartbeat sweep).
// sendProtectionAlertToPartners may still send them; HTTP POST must not.
var serverOnlyProtectionAlertTypes = map[string]struct{}{
	"heartbeat_silence": {}, // heartbeat_sweep
}

func isClientProtectionAlertType(t string) bool {
	_, ok := clientProtectionAlertTypes[t]
	return ok
}

func isServerOnlyProtectionAlertType(t string) bool {
	_, ok := serverOnlyProtectionAlertTypes[t]
	return ok
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
		Type       string `json:"type"`
		Detail     string `json:"detail"`
		DeviceID   string `json:"device_id"`
		DeviceKind string `json:"device_kind"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Type == "" {
		writeError(w, http.StatusBadRequest, "type is required")
		return
	}
	if isServerOnlyProtectionAlertType(req.Type) || !isClientProtectionAlertType(req.Type) {
		writeError(w, http.StatusBadRequest, "unknown protection alert type")
		return
	}
	// Shipped iOS 1.0.2 sends type and detail only. Detail is still ignored.
	// A newer client may name the device type so the sentence can say
	// "on their Windows PC". Free text is rejected so an app name cannot land
	// in the partner notification.
	label := ""
	collapseExtra := ""
	if req.DeviceID != "" || req.DeviceKind != "" {
		if !validDeviceID(req.DeviceID) {
			writeError(w, http.StatusBadRequest, "device id is not recognized")
			return
		}
		var ok bool
		label, ok = deviceKindLabel(req.DeviceKind)
		if !ok {
			writeError(w, http.StatusBadRequest, "device type is not recognized")
			return
		}
		collapseExtra = deviceCollapseKey(req.DeviceID)
	}
	if _, ok := h.requireConfirmedAccount(w, r); !ok {
		return
	}

	var senderName string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&senderName); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	// Sign-out is the app's server signal that this account left the device.
	// Clear this account's tokens before the alert is stored.
	if req.Type == "signed_out" {
		if _, err := h.DB.ExecContext(r.Context(), `
			UPDATE device_tokens
			SET is_active = FALSE, updated_at = NOW()
			WHERE user_id = $1
		`, userID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to record protection alert")
			return
		}
	}

	// Store the alert before acknowledging so a missed push still has a history row.
	body := protectionBody(senderName, req.Type, label)
	if err := h.storeAndPushProtectionText(userID, senderName, req.Type, body, false, collapseExtra); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record protection alert")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// sendProtectionAlertToPartners fans out a PROTECTION_ALERT push to all
// accepted partners of userID. It does not apply the HTTP client allowlist, so
// internal callers (heartbeat sweep) may send server-only types.
// Only metadata is sent — never screen content.
// storeAndPushProtection writes one alert row per accepted partner and group
// co-member, then pushes. detail from the client is ignored.
// silenceDedup skips a recipient who already has this silence alert in the last 6 hours.
func (h *H) storeAndPushProtection(userID int64, senderName, alertType string, silenceDedup bool) error {
	return h.storeAndPushProtectionText(userID, senderName, alertType, protectionAlertBody(senderName, alertType, ""), silenceDedup, "")
}

// storeAndPushProtectionText is storeAndPushProtection with a finished sentence.
// collapseExtra keeps two devices from replacing each other's push. It is a
// short hash, never the device id and never part of the sentence.
func (h *H) storeAndPushProtectionText(userID int64, senderName, alertType, body string, silenceDedup bool, collapseExtra string) error {
	ctx := context.Background()
	ids, err := h.protectionRecipients(ctx, userID)
	if err != nil {
		log.Printf("[protection] query recipients user=%d: %v", userID, err)
		return err
	}
	var deliver []int64
	for _, pid := range ids {
		if silenceDedup && h.recentSilenceAlert(ctx, userID, pid) {
			continue
		}
		if _, err := h.DB.ExecContext(ctx, `
			INSERT INTO alerts (event_id, recipient_user_id, kind, alert_type, summary, subject_user_id)
			VALUES (NULL, $1, 'protection', $2, $3, $4)
		`, pid, alertType, body, userID); err != nil {
			log.Printf("[protection] insert alert user=%d recipient=%d: %v", userID, pid, err)
			return err
		}
		deliver = append(deliver, pid)
	}
	if len(deliver) == 0 {
		return nil
	}
	h.inBackground(func() {
		h.pushProtection(userID, senderName, alertType, body, deliver, collapseExtra)
	})
	return nil
}

func (h *H) protectionRecipients(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT partner_id FROM relationships
		WHERE user_id = $1 AND status = 'accepted'
		UNION
		SELECT DISTINCT gm2.user_id
		FROM group_members gm1
		JOIN group_members gm2
		  ON gm2.group_id = gm1.group_id AND gm2.user_id != $1
		WHERE gm1.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	seen := map[int64]bool{}
	for rows.Next() {
		var pid int64
		if rows.Scan(&pid) == nil && pid != userID && !seen[pid] {
			seen[pid] = true
			ids = append(ids, pid)
		}
	}
	return ids, rows.Err()
}

func (h *H) recentSilenceAlert(ctx context.Context, subject, recipient int64) bool {
	var n int
	_ = h.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM alerts
		WHERE subject_user_id = $1 AND recipient_user_id = $2
		  AND alert_type = 'heartbeat_silence'
		  AND created_at > NOW() - INTERVAL '6 hours'
	`, subject, recipient).Scan(&n)
	return n > 0
}

func (h *H) pushProtection(userID int64, senderName, alertType, body string, partnerIDs []int64, collapseExtra string) {
	ctx := context.Background()
	// detail is never interpolated into the APNs body. Only fixed copy per type,
	// plus an allowlisted device label when the client sent one.
	collapse := fmt.Sprintf("protection-%d-%s", userID, alertType)
	if collapseExtra != "" {
		collapse += "-" + collapseExtra
	}
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
		"sender_name":       clipPushName(senderName),
	}
	n := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: collapse,
		Payload:    payload,
	}
	for _, pid := range partnerIDs {
		h.notifyPartnerByID(ctx, pid, senderName, n)
	}
}

// sendProtectionAlertToPartners is the internal entry used by the silence sweep.
func (h *H) sendProtectionAlertToPartners(userID int64, senderName, alertType, _ string) {
	_ = h.storeAndPushProtection(userID, senderName, alertType, alertType == "heartbeat_silence")
}

// protectionAlertBody returns fixed notification copy for the given alert type.
// detail is ignored so client-supplied text cannot inflate the APNs body.
// Screen content is never included — only metadata.
func protectionAlertBody(name, alertType, _ string) string {
	name = clipPushName(name)
	switch alertType {
	case "monitoring_disabled":
		return name + " turned off activity monitoring"
	case "shielding_disabled":
		return name + " turned off app blocking"
	case "lockout_disabled":
		return name + " disabled App Lockout"
	case "lockout_broadcast_stopped", "deep_scan_stopped":
		return name + " stopped Deep Scan"
	case "selection_cleared":
		return name + " removed every chosen app"
	case "pin_wrong_attempt", "wrong_pin_attempt":
		return name + " entered an incorrect Partner PIN"
	case "pin_removed":
		return name + "'s protection PIN was removed"
	case "pin_changed":
		return name + "'s protection PIN was changed"
	case "family_controls_revoked":
		return name + " revoked Family Controls authorization"
	case "signed_out":
		return name + " signed out"
	case "app_removed":
		return name + "'s app was removed from a device"
	case "heartbeat_silence":
		return name + "'s device has stopped sending heartbeats"
	default:
		return name + " changed a protection setting"
	}
}

func truncateForLog(s string, maxRunes int) string {
	if maxRunes <= 0 || s == "" {
		return ""
	}
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	return string(r[:maxRunes]) + "…"
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

	if denied, err := h.pinAttemptDenied(r.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "could not check PIN attempts")
		return
	} else if denied {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "900")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"error":"too many PIN attempts. Try again in 15 minutes","success":false,"reason":"rate_limited"}`)
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

	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT pin_hash FROM relationships
		WHERE user_id = $1 AND status = 'accepted' AND pin_hash IS NOT NULL AND pin_hash <> ''
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up PIN")
		return
	}
	defer rows.Close()
	var hashes []string
	for rows.Next() {
		var hash string
		if rows.Scan(&hash) == nil && hash != "" {
			hashes = append(hashes, hash)
		}
	}
	if len(hashes) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"success": false, "reason": "no_pin"})
		return
	}
	for _, hash := range hashes {
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.PIN)) == nil {
			writeJSON(w, http.StatusOK, map[string]any{"success": true})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": false, "reason": "wrong_pin"})
}

func (h *H) pinAttemptDenied(ctx context.Context, userID int64) (bool, error) {
	var failures int
	err := h.DB.QueryRowContext(ctx, `
		INSERT INTO pin_verify_attempts (user_id, failures, window_start)
		VALUES ($1, 1, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
		  failures = CASE
		    WHEN pin_verify_attempts.window_start < NOW() - INTERVAL '15 minutes' THEN 1
		    ELSE pin_verify_attempts.failures + 1
		  END,
		  window_start = CASE
		    WHEN pin_verify_attempts.window_start < NOW() - INTERVAL '15 minutes' THEN NOW()
		    ELSE pin_verify_attempts.window_start
		  END
		RETURNING failures
	`, userID).Scan(&failures)
	if err != nil {
		return false, err
	}
	return failures > 5, nil
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

	// A protection PIN can be saved only after the partnership is accepted.
	// A row that belongs to this partner but is still waiting says so.
	// Any other row keeps the same not-authorized response as before.
	var monitoredUserID int64
	var status string
	err = h.DB.QueryRowContext(r.Context(), `
		SELECT user_id, status FROM relationships
		WHERE id = $1 AND partner_id = $2
	`, relID, callerID).Scan(&monitoredUserID, &status)
	if err != nil {
		writeError(w, http.StatusForbidden, "not authorized to set PIN on this relationship")
		return
	}
	if status != "accepted" {
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
