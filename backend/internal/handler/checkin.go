package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"
)

const maxCheckInNote = 140

func validateCheckIn(status, note string) (string, string, string) {
	switch status {
	case "steady", "need_support", "slipped":
	default:
		return "", "", "invalid status"
	}
	note = strings.TrimSpace(note)
	note = strings.ReplaceAll(note, "\n", " ")
	note = strings.ReplaceAll(note, "\r", " ")
	if utf8.RuneCountInString(note) > maxCheckInNote {
		return "", "", "note is too long"
	}
	return status, note, ""
}

func checkInStatusLabel(status string) string {
	switch status {
	case "need_support":
		return "Asking for support"
	case "slipped":
		return "Had a slip"
	default:
		return "Steady"
	}
}

func checkInPushBody(name, status, note string) string {
	name = clipPushName(name)
	body := name + " sent a check-in: " + checkInStatusLabel(status) + "."
	if note != "" {
		body += " " + note
	}
	return body
}

// CreateCheckIn records a note the person chose to send.
// POST /check-ins
// Body: status, note, shielding_on, session_shield_on, schedule_on, deep_scan_on.
func (h *H) CreateCheckIn(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	var req struct {
		Status          string `json:"status"`
		Note            string `json:"note"`
		ShieldingOn     bool   `json:"shielding_on"`
		SessionShieldOn bool   `json:"session_shield_on"`
		ScheduleOn      bool   `json:"schedule_on"`
		DeepScanOn      bool   `json:"deep_scan_on"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	status, note, msg := validateCheckIn(req.Status, req.Note)
	if msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	var id int64
	err := h.DB.QueryRowContext(r.Context(), `
		INSERT INTO check_ins
			(user_id, status, note, shielding_on, session_shield_on, schedule_on, deep_scan_on)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		userID, status, note, req.ShieldingOn, req.SessionShieldOn, req.ScheduleOn, req.DeepScanOn,
	).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save check-in")
		return
	}

	var callerName string
	_ = h.DB.QueryRowContext(r.Context(), `SELECT name FROM users WHERE id = $1`, userID).Scan(&callerName)
	body := checkInPushBody(callerName, status, note)

	partnerRows, partnerErr := h.DB.QueryContext(r.Context(), `
		SELECT DISTINCT partner_id
		FROM   relationships
		WHERE  user_id = $1 AND status = 'accepted'
		UNION
		SELECT DISTINCT gm2.user_id
		FROM   group_members gm1
		JOIN   group_members gm2
		       ON gm2.group_id = gm1.group_id AND gm2.user_id != $1
		WHERE  gm1.user_id = $1
	`, userID)
	var partnerIDs []int64
	if partnerErr != nil {
		log.Printf("check-in %d: query partners: %v", id, partnerErr)
	} else {
		for partnerRows.Next() {
			var pid int64
			if partnerRows.Scan(&pid) == nil && pid != 0 {
				partnerIDs = append(partnerIDs, pid)
			}
		}
		partnerRows.Close()
	}

	go func() {
		ctx := context.Background()
		payload := map[string]any{
			"aps": map[string]any{
				"alert": map[string]string{
					"title": "Check-in",
					"body":  body,
				},
				"sound": "default",
			},
			"notification_type": "CHECK_IN",
			"sender_name":       clipPushName(callerName),
			"status":            status,
			"note":              note,
			"shielding_on":      req.ShieldingOn,
			"session_shield_on": req.SessionShieldOn,
			"schedule_on":       req.ScheduleOn,
			"deep_scan_on":      req.DeepScanOn,
		}
		n := &apns.Notification{
			PushType:   "alert",
			Priority:   10,
			CollapseID: "check-in",
			Payload:    payload,
		}
		for _, pid := range partnerIDs {
			h.notifyPartnerByID(ctx, pid, callerName, n)
		}
	}()

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":     id,
		"status": status,
		"note":   note,
	})
}
