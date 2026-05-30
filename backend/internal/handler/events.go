package handler

import (
	"encoding/json"
	"net/http"

	rfauth "remain-faithful/backend/internal/auth"
)

// CreateEvent records a flagged monitoring event from a device and fans out
// alerts to all accepted accountability partners.
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

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to begin transaction")
		return
	}
	defer tx.Rollback()

	var eventID int64
	var timestamp string

	if req.Timestamp != "" {
		err = tx.QueryRowContext(r.Context(),
			`INSERT INTO events (user_id, category, severity, summary, timestamp)
			 VALUES ($1, $2, $3, $4, $5::timestamptz)
			 RETURNING id, timestamp`,
			userID, req.Category, req.Severity, req.Summary, req.Timestamp,
		).Scan(&eventID, &timestamp)
	} else {
		err = tx.QueryRowContext(r.Context(),
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

	// Fan out: create one alert per accepted relationship for this user.
	relRows, err := tx.QueryContext(r.Context(),
		`SELECT id FROM relationships WHERE user_id = $1 AND status = 'accepted'`,
		userID,
	)
	if err == nil {
		defer relRows.Close()
		for relRows.Next() {
			var relID int64
			if relRows.Scan(&relID) == nil {
				_, _ = tx.ExecContext(r.Context(),
					`INSERT INTO alerts (event_id, relationship_id) VALUES ($1, $2)`,
					eventID, relID,
				)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":        eventID,
		"user_id":   userID,
		"category":  req.Category,
		"severity":  req.Severity,
		"summary":   req.Summary,
		"timestamp": timestamp,
	})
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
