package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	rfauth "remain-faithful/backend/internal/auth"
)

// GetMe returns the authenticated user's profile.
// GET /users/me
func (h *H) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var id int64
	var name, email, createdAt string
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT id, name, email, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&id, &name, &email, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         id,
		"name":       name,
		"email":      email,
		"created_at": createdAt,
	})
}

// UpdateMe updates the authenticated user's name and/or email.
// PUT /users/me
// Body: { "name": "...", "email": "..." }
func (h *H) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Name  = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	var id int64
	var name, email, createdAt string
	err := h.DB.QueryRowContext(r.Context(),
		`UPDATE users SET name = $1, email = $2 WHERE id = $3
		 RETURNING id, name, email, created_at`,
		req.Name, req.Email, userID,
	).Scan(&id, &name, &email, &createdAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         id,
		"name":       name,
		"email":      email,
		"created_at": createdAt,
	})
}

// GetUserStats returns the authenticated user's clean-streak statistics computed
// from the full event history in the database.
// GET /users/me/stats
// Response: { "streak_days": int, "best_streak": int, "week": [bool×7] }
// week[0] = 6 days ago, week[6] = today (oldest → newest).
func (h *H) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var createdAt time.Time
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT created_at FROM users WHERE id = $1`, userID,
	).Scan(&createdAt); err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	// Fetch all distinct UTC dates on which at least one non-clean event was recorded.
	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT DISTINCT DATE(timestamp AT TIME ZONE 'UTC')
		FROM   events
		WHERE  user_id  = $1
		  AND  category NOT IN ('clean')
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load stats")
		return
	}
	defer rows.Close()

	flagged := map[string]bool{}
	for rows.Next() {
		var d time.Time
		if rows.Scan(&d) == nil {
			flagged[d.UTC().Format("2006-01-02")] = true
		}
	}

	now        := time.Now().UTC()
	today      := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	createdDay := time.Date(createdAt.Year(), createdAt.Month(), createdAt.Day(), 0, 0, 0, 0, time.UTC)

	// Current streak: consecutive clean days from today backward to account creation.
	streak := 0
	d := today
	for !d.Before(createdDay) {
		if flagged[d.Format("2006-01-02")] {
			break
		}
		streak++
		d = d.AddDate(0, 0, -1)
		if streak > 3650 {
			break
		}
	}
	if createdDay.Equal(today) {
		streak = 0 // account created today: no completed clean day yet
	}

	// 7-day week view (index 0 = 6 days ago, index 6 = today).
	week := make([]bool, 7)
	for i := range week {
		day    := today.AddDate(0, 0, -(6 - i))
		dayStr := day.Format("2006-01-02")
		week[i] = !day.Before(createdDay) && !flagged[dayStr]
	}

	// Best streak: longest clean run in the window max(365, streak+1) days back.
	windowDays := max(365, streak+1)
	best, run := 0, 0
	for i := range windowDays {
		day := today.AddDate(0, 0, -i)
		if day.Before(createdDay) {
			break
		}
		if flagged[day.Format("2006-01-02")] {
			run = 0
		} else {
			run++
			if run > best {
				best = run
			}
		}
	}
	best = max(best, streak)

	writeJSON(w, http.StatusOK, map[string]any{
		"streak_days": streak,
		"best_streak": best,
		"week":        week,
	})
}
