package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/apns"
)

// CreateEvent records a flagged monitoring event from a device and fans out
// alerts to all accountability partners (via relationships AND group membership).
// POST /events
// Body: { "category": "...", "severity": "low|medium|high", "summary": "...", "timestamp": "..." }
func (h *H) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Category  string `json:"category"`
		Severity  string `json:"severity"`
		Summary   string `json:"summary"`
		Timestamp string `json:"timestamp"` // optional ISO-8601; defaults to NOW()
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Category == "" || req.Severity == "" || req.Summary == "" {
		writeError(w, http.StatusBadRequest, "category, severity, and summary are required")
		return
	}
	validCategories := map[string]bool{
		"adult_content": true, "gambling": true,
		"violence": true, "self_harm": true, "clean": true,
	}
	validSeverities := map[string]bool{
		"informational": true, "concerning": true, "severe": true,
	}
	if !validCategories[req.Category] {
		writeError(w, http.StatusBadRequest, "invalid category")
		return
	}
	if !validSeverities[req.Severity] {
		writeError(w, http.StatusBadRequest, "invalid severity")
		return
	}

	// Step 1: insert the event — single auto-committed statement, independent of alert creation.
	var eventID int64
	var timestamp string
	var err error
	if req.Timestamp != "" {
		err = h.DB.QueryRowContext(r.Context(),
			`INSERT INTO events (user_id, category, severity, summary, timestamp)
			 VALUES ($1, $2, $3, $4, $5::timestamptz)
			 RETURNING id, timestamp`,
			userID, req.Category, req.Severity, req.Summary, req.Timestamp,
		).Scan(&eventID, &timestamp)
	} else {
		err = h.DB.QueryRowContext(r.Context(),
			`INSERT INTO events (user_id, category, severity, summary)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id, timestamp`,
			userID, req.Category, req.Severity, req.Summary,
		).Scan(&eventID, &timestamp)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	// Step 2: fan out alerts to partners — best-effort after the event is committed.
	// Collect all partner IDs first (fully drain the rows before any INSERT so we
	// never attempt two queries concurrently on the same connection).
	type partnerAlert struct {
		partnerID int64
		alertID   int64
	}
	var partnerAlerts []partnerAlert

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
	if partnerErr != nil {
		log.Printf("CreateEvent %d: query partners: %v", eventID, partnerErr)
	} else {
		var partnerIDs []int64
		for partnerRows.Next() {
			var pid int64
			if partnerRows.Scan(&pid) == nil {
				partnerIDs = append(partnerIDs, pid)
			}
		}
		if rowErr := partnerRows.Err(); rowErr != nil {
			log.Printf("CreateEvent %d: iterate partners: %v", eventID, rowErr)
		}
		partnerRows.Close()

		for _, pid := range partnerIDs {
			var alertID int64
			if scanErr := h.DB.QueryRowContext(r.Context(),
				`INSERT INTO alerts (event_id, recipient_user_id) VALUES ($1, $2) RETURNING id`,
				eventID, pid,
			).Scan(&alertID); scanErr != nil {
				log.Printf("CreateEvent %d: insert alert for partner %d: %v", eventID, pid, scanErr)
			} else {
				partnerAlerts = append(partnerAlerts, partnerAlert{pid, alertID})
			}
		}
	}

	// Look up the caller's display name for the push payload.
	var callerName string
	_ = h.DB.QueryRowContext(r.Context(),
		`SELECT name FROM users WHERE id = $1`, userID,
	).Scan(&callerName)

	capturedPartnerAlerts := partnerAlerts
	capturedEventID := eventID
	capturedCategory := req.Category
	capturedSeverity := req.Severity
	capturedSummary := req.Summary
	capturedTimestamp := timestamp
	capturedCaller := callerName

	go func() {
		ctx := context.Background()
		payload := map[string]any{
			"aps": map[string]any{
				"alert": map[string]string{
					"title": alertTitle(capturedSeverity),
					"body":  alertBody(capturedCaller, capturedCategory, capturedSeverity),
				},
				"sound": "default",
			},
			"notification_type": "CONTENT_FLAGGED",
			"sender_name":       capturedCaller,
			"event_id":          capturedEventID,
			"category":          capturedCategory,
			"severity":          capturedSeverity,
			"summary":           capturedSummary,
			"timestamp":         capturedTimestamp,
		}
		n := &apns.Notification{
			PushType:   "alert",
			Priority:   10,
			CollapseID: fmt.Sprintf("event-%d", capturedEventID),
			Payload:    payload,
		}

		for _, pa := range capturedPartnerAlerts {
			// Throttle: skip push if this partner already got one in the last 5 minutes.
			var recentCount int
			_ = h.DB.QueryRowContext(ctx, `
				SELECT COUNT(*) FROM alerts
				WHERE  recipient_user_id = $1
				  AND  created_at > NOW() - INTERVAL '5 minutes'
				  AND  id < $2
			`, pa.partnerID, pa.alertID).Scan(&recentCount)
			if recentCount > 0 {
				continue
			}
			h.notifyPartnerByID(ctx, pa.partnerID, capturedCaller, n)
		}
	}()

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":        eventID,
		"user_id":   userID,
		"category":  req.Category,
		"severity":  req.Severity,
		"summary":   req.Summary,
		"timestamp": timestamp,
	})
}

// alertTitle returns a push notification title based on severity.
func alertTitle(severity string) string {
	switch severity {
	case "severe":
		return "Accountability Alert"
	case "concerning":
		return "Check In With Your Partner"
	default:
		return "Partner Activity"
	}
}

// alertBody returns a conversation-starter notification body.
func alertBody(name, category, severity string) string {
	if severity == "severe" {
		switch category {
		case "self_harm":
			return fmt.Sprintf("%s needs you right now — their device flagged self-harm content. Reach out immediately.", name)
		case "adult_content":
			return fmt.Sprintf("%s needs accountability — their device flagged explicit content. Be present and non-judgmental.", name)
		case "gambling":
			return fmt.Sprintf("%s may be struggling — their device flagged gambling content. A quick check-in goes a long way.", name)
		default:
			return fmt.Sprintf("%s needs accountability right now. Open the app to see what was flagged.", name)
		}
	}
	if severity == "concerning" {
		return fmt.Sprintf("%s may need support — their device flagged %s content. Consider checking in today.", name, category)
	}
	return fmt.Sprintf("%s's device flagged activity worth noting. Keep them in your prayers.", name)
}

// ListEvents returns the authenticated user's flagged events, newest first.
// GET /events
func (h *H) ListEvents(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT id, category, severity, summary, timestamp
		FROM   events
		WHERE  user_id = $1
		ORDER  BY timestamp DESC
		LIMIT  100
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list events")
		return
	}
	defer rows.Close()

	type event struct {
		ID        int64  `json:"id"`
		UserID    int64  `json:"user_id"`
		Category  string `json:"category"`
		Severity  string `json:"severity"`
		Summary   string `json:"summary"`
		Timestamp string `json:"timestamp"`
	}

	result := []event{}
	for rows.Next() {
		e := event{UserID: userID}
		if err := rows.Scan(&e.ID, &e.Category, &e.Severity, &e.Summary, &e.Timestamp); err != nil {
			continue
		}
		result = append(result, e)
	}
	writeJSON(w, http.StatusOK, result)
}
