package handler

import (
	"net/http"

	rfauth "remain-faithful/backend/internal/auth"
)

// AlertUnreadCount returns the number of unseen alerts for the authenticated user.
// GET /alerts/count
func (h *H) AlertUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	var count int
	err := h.DB.QueryRowContext(r.Context(), `
		SELECT COUNT(*)
		FROM   alerts
		WHERE  recipient_user_id = $1 AND seen = FALSE
	`, userID).Scan(&count)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to count alerts")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"unseen": count})
}

// MarkAlertsSeen marks all unseen alerts for the authenticated user as seen.
// POST /alerts/mark-seen
func (h *H) MarkAlertsSeen(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	_, err := h.DB.ExecContext(r.Context(), `
		UPDATE alerts
		SET    seen = TRUE
		WHERE  recipient_user_id = $1 AND seen = FALSE
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to mark alerts seen")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ListAlerts returns all alerts sent to the authenticated user.
// GET /alerts
func (h *H) ListAlerts(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT a.id,
		       a.event_id,
		       a.seen,
		       a.created_at,
		       e.user_id,
		       e.category,
		       e.severity,
		       e.summary,
		       e.timestamp
		FROM   alerts a
		JOIN   events e ON e.id = a.event_id
		WHERE  a.recipient_user_id = $1
		ORDER  BY a.created_at DESC
		LIMIT  100
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list alerts")
		return
	}
	defer rows.Close()

	type eventInfo struct {
		ID        int64  `json:"id"`
		UserID    int64  `json:"user_id"`
		Category  string `json:"category"`
		Severity  string `json:"severity"`
		Summary   string `json:"summary"`
		Timestamp string `json:"timestamp"`
	}
	type alert struct {
		ID        int64     `json:"id"`
		EventID   int64     `json:"event_id"`
		Seen      bool      `json:"seen"`
		CreatedAt string    `json:"created_at"`
		Event     eventInfo `json:"event"`
	}

	result := []alert{}
	for rows.Next() {
		var a alert
		if err := rows.Scan(
			&a.ID, &a.EventID, &a.Seen, &a.CreatedAt,
			&a.Event.UserID, &a.Event.Category, &a.Event.Severity,
			&a.Event.Summary, &a.Event.Timestamp,
		); err != nil {
			continue
		}
		a.Event.ID = a.EventID
		result = append(result, a)
	}
	writeJSON(w, http.StatusOK, result)
}
