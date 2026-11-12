package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"

	"github.com/lib/pq"
)

// MigrateWeeklySummary adds the weekly partner summary tables.
// The job stays off until WEEKLY_PARTNER_SUMMARY is set to true.
func MigrateWeeklySummary(db *sql.DB) error {
	_, err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS weekly_summary_opt_out BOOLEAN NOT NULL DEFAULT FALSE;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS weekly_summary_token TEXT;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_weekly_summary_token
			ON users (weekly_summary_token) WHERE weekly_summary_token IS NOT NULL;
		CREATE TABLE IF NOT EXISTS weekly_summary_sends (
			partner_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			subject_user_id BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			week_start      DATE        NOT NULL,
			sent_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (partner_id, subject_user_id, week_start)
		);
	`)
	return err
}

func weeklySummaryEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("WEEKLY_PARTNER_SUMMARY"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func apiPublicURL() string {
	if u := strings.TrimSpace(os.Getenv("API_PUBLIC_URL")); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://api.remainfaithful.com"
}

// StartWeeklySummaries checks once an hour. It sends nothing unless the flag is on.
func (h *H) StartWeeklySummaries(ctx context.Context) {
	h.sendDueWeeklySummaries(ctx, time.Now())
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			h.sendDueWeeklySummaries(ctx, now)
		}
	}
}

func (h *H) sendDueWeeklySummaries(ctx context.Context, now time.Time) {
	if !weeklySummaryEnabled() || h == nil || h.DB == nil {
		return
	}
	subjects, err := h.weeklySubjects(ctx)
	if err != nil {
		log.Printf("[weekly] subjects: %v", err)
		return
	}
	sent := 0
	for _, subject := range subjects {
		n, err := h.sendSubjectWeek(ctx, subject, now)
		if err != nil {
			log.Printf("[weekly] subject=%d: %v", subject.id, err)
			continue
		}
		sent += n
	}
	if sent > 0 {
		log.Printf("[weekly] sent %d summaries", sent)
	}
}

type weeklySubject struct {
	id      int64
	name    string
	created time.Time
}

type weeklyRecipient struct {
	id     int64
	name   string
	email  string
	optOut bool
	token  string
}

func (h *H) weeklySubjects(ctx context.Context) ([]weeklySubject, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT DISTINCT u.id, u.name, u.created_at
		FROM users u
		WHERE EXISTS (
			SELECT 1 FROM relationships r
			WHERE r.user_id = u.id AND r.status = 'accepted'
		)
		OR EXISTS (
			SELECT 1
			FROM group_members gm1
			JOIN group_members gm2
			  ON gm2.group_id = gm1.group_id AND gm2.user_id <> gm1.user_id
			WHERE gm1.user_id = u.id
		)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []weeklySubject
	for rows.Next() {
		var s weeklySubject
		if err := rows.Scan(&s.id, &s.name, &s.created); err != nil {
			continue
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (h *H) sendSubjectWeek(ctx context.Context, subject weeklySubject, now time.Time) (int, error) {
	loc := time.UTC
	weekStart, open := weeklySendOpen(now, loc)
	if !open {
		return 0, nil
	}
	recipients, err := h.weeklyRecipients(ctx, subject.id)
	if err != nil {
		return 0, err
	}
	facts, err := h.weeklyFacts(ctx, subject, now, loc)
	if err != nil {
		return 0, err
	}
	sent := 0
	for _, rcpt := range recipients {
		if rcpt.optOut || rcpt.id == subject.id || strings.TrimSpace(rcpt.email) == "" {
			continue
		}
		already, err := h.weeklyAlreadySent(ctx, rcpt.id, subject.id, weekStart)
		if err != nil || already {
			continue
		}
		token, err := h.ensureWeeklyToken(ctx, rcpt)
		if err != nil || token == "" {
			continue
		}
		letter := weeklyLetter{
			ToID:        rcpt.id,
			ToEmail:     rcpt.email,
			ToName:      rcpt.name,
			SubjectName: subject.name,
			Days:        facts.days,
			Alerts:      facts.alerts,
			Streak:      facts.streak,
			UnsubURL:    apiPublicURL() + "/weekly-summary/unsubscribe?token=" + token,
		}
		if err := h.deliverWeeklyLetter(ctx, letter); err != nil {
			log.Printf("[weekly] deliver partner=%d subject=%d: %v", rcpt.id, subject.id, err)
			continue
		}
		if _, err := h.DB.ExecContext(ctx, `
			INSERT INTO weekly_summary_sends (partner_id, subject_user_id, week_start)
			VALUES ($1, $2, $3::date)
			ON CONFLICT DO NOTHING
		`, rcpt.id, subject.id, weekStart); err != nil {
			log.Printf("[weekly] mark partner=%d subject=%d: %v", rcpt.id, subject.id, err)
		}
		sent++
	}
	return sent, nil
}

func (h *H) weeklyRecipients(ctx context.Context, subjectID int64) ([]weeklyRecipient, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT u.id, u.name, u.email, u.weekly_summary_opt_out, COALESCE(u.weekly_summary_token, '')
		FROM users u
		WHERE u.id IN (
			SELECT partner_id FROM relationships
			WHERE user_id = $1 AND status = 'accepted'
			UNION
			SELECT gm2.user_id
			FROM group_members gm1
			JOIN group_members gm2
			  ON gm2.group_id = gm1.group_id AND gm2.user_id <> $1
			WHERE gm1.user_id = $1
		)
	`, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []weeklyRecipient
	for rows.Next() {
		var rcpt weeklyRecipient
		if err := rows.Scan(&rcpt.id, &rcpt.name, &rcpt.email, &rcpt.optOut, &rcpt.token); err != nil {
			continue
		}
		out = append(out, rcpt)
	}
	return out, rows.Err()
}

type weekFacts struct {
	days, alerts, streak int
}

func (h *H) weeklyFacts(ctx context.Context, subject weeklySubject, now time.Time, loc *time.Location) (weekFacts, error) {
	start, end := lastSevenDays(now, loc)
	beats, err := h.daySet(ctx, `
		SELECT DISTINCT to_char(at AT TIME ZONE 'UTC', 'YYYY-MM-DD')
		FROM protection_heartbeats
		WHERE user_id = $1
		  AND at >= NOW() - INTERVAL '400 days'
	`, subject.id)
	if err != nil {
		return weekFacts{}, err
	}
	flagged, err := h.daySet(ctx, `
		SELECT DISTINCT to_char(timestamp AT TIME ZONE 'UTC', 'YYYY-MM-DD')
		FROM events
		WHERE user_id = $1
		  AND category <> 'clean'
		  AND timestamp >= NOW() - INTERVAL '4000 days'
	`, subject.id)
	if err != nil {
		return weekFacts{}, err
	}
	offs, err := h.offDays(ctx, subject.id)
	if err != nil {
		return weekFacts{}, err
	}
	clean := map[string]struct{}{}
	for day := range beats {
		if _, bad := flagged[day]; bad {
			continue
		}
		if _, bad := offs[day]; bad {
			continue
		}
		clean[day] = struct{}{}
	}
	var eventCount int
	err = h.DB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM events
		WHERE user_id = $1
		  AND category <> 'clean'
		  AND (timestamp AT TIME ZONE 'UTC')::date BETWEEN $2::date AND $3::date
	`, subject.id, start, end).Scan(&eventCount)
	if err != nil {
		return weekFacts{}, err
	}
	var legacy int
	if err := h.DB.QueryRowContext(ctx, `SELECT COALESCE(legacy_streak, 0) FROM users WHERE id = $1`, subject.id).Scan(&legacy); err != nil {
		return weekFacts{}, err
	}
	got, err := h.streaksAt(ctx, map[int64]streakMeta{
		subject.id: {Created: subject.created, Legacy: legacy},
	}, now, 0)
	if err != nil {
		return weekFacts{}, err
	}
	return weekFacts{
		days:   countDays(clean, start, end),
		alerts: eventCount,
		streak: got[subject.id].Days,
	}, nil
}

func (h *H) offDays(ctx context.Context, userID int64) (map[string]struct{}, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT DISTINCT to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD')
		FROM alerts
		WHERE subject_user_id = $1
		  AND alert_type = ANY($2)
		  AND created_at >= NOW() - INTERVAL '4000 days'
	`, userID, pq.Array(protectionOffTypes))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]struct{}{}
	for rows.Next() {
		var day string
		if rows.Scan(&day) == nil && len(day) >= 10 {
			out[day[:10]] = struct{}{}
		}
	}
	return out, rows.Err()
}

func (h *H) daySet(ctx context.Context, query string, userID int64) (map[string]struct{}, error) {
	rows, err := h.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]struct{}{}
	for rows.Next() {
		var day string
		if rows.Scan(&day) == nil && len(day) >= 10 {
			out[day[:10]] = struct{}{}
		}
	}
	return out, rows.Err()
}

func (h *H) weeklyAlreadySent(ctx context.Context, partnerID, subjectID int64, weekStart string) (bool, error) {
	var one int
	err := h.DB.QueryRowContext(ctx, `
		SELECT 1 FROM weekly_summary_sends
		WHERE partner_id = $1 AND subject_user_id = $2 AND week_start = $3::date
	`, partnerID, subjectID, weekStart).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *H) ensureWeeklyToken(ctx context.Context, rcpt weeklyRecipient) (string, error) {
	if validUnsubToken(rcpt.token) {
		return rcpt.token, nil
	}
	token, err := newUnsubToken()
	if err != nil {
		return "", err
	}
	if _, err := h.DB.ExecContext(ctx, `
		UPDATE users
		SET weekly_summary_token = $2
		WHERE id = $1 AND (weekly_summary_token IS NULL OR weekly_summary_token = '')
	`, rcpt.id, token); err != nil {
		return "", err
	}
	var stored string
	if err := h.DB.QueryRowContext(ctx, `
		SELECT COALESCE(weekly_summary_token, '') FROM users WHERE id = $1
	`, rcpt.id).Scan(&stored); err != nil {
		return "", err
	}
	if !validUnsubToken(stored) {
		return token, nil
	}
	return stored, nil
}

type weeklyLetter struct {
	ToID                 int64
	ToEmail, ToName      string
	SubjectName          string
	Days, Alerts, Streak int
	UnsubURL             string
}

func (h *H) deliverWeeklyLetter(ctx context.Context, letter weeklyLetter) error {
	if !weeklySummaryEnabled() {
		return nil
	}
	if h.Email != nil {
		if err := h.Email.SendWeeklySummary(letter.ToEmail, letter.ToName, letter.SubjectName, letter.Days, letter.Alerts, letter.Streak, letter.UnsubURL); err != nil {
			return err
		}
	}
	if h.APNS != nil && letter.ToID != 0 {
		body := weeklyPushBody(letter.SubjectName, letter.Days, letter.Alerts, letter.Streak)
		payload := map[string]any{
			"aps": map[string]any{
				"alert": map[string]string{"title": "Weekly summary", "body": body},
				"sound": "default",
			},
			"notification_type": "WEEKLY_SUMMARY",
		}
		h.notifyPartnerByID(ctx, letter.ToID, letter.SubjectName, &apns.Notification{
			PushType:   "alert",
			Priority:   10,
			CollapseID: "weekly-summary-" + strconv.FormatInt(letter.ToID, 10),
			Payload:    payload,
		})
	}
	return nil
}

func weeklyPushBody(name string, days, alerts, streak int) string {
	return clipPushName(name) + ": " +
		countPhrase(days, "day protected", "days protected") + ", " +
		countPhrase(alerts, "alert", "alerts") + ", streak " +
		strconv.Itoa(streak) + "."
}

func countPhrase(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return strconv.Itoa(n) + " " + plural
}

// weeklySendOpen is true from Monday 08:00 in loc through the rest of that week.
// The scheduled job passes UTC. It does not store a time zone.
func weeklySendOpen(now time.Time, loc *time.Location) (string, bool) {
	if loc == nil {
		loc = time.UTC
	}
	local := now.In(loc)
	weekday := int(local.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := time.Date(local.Year(), local.Month(), local.Day(), 8, 0, 0, 0, loc).AddDate(0, 0, -(weekday - 1))
	start := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, loc)
	return start.Format("2006-01-02"), !local.Before(monday)
}

func lastSevenDays(now time.Time, loc *time.Location) (string, string) {
	if loc == nil {
		loc = time.UTC
	}
	local := now.In(loc)
	end := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	start := end.AddDate(0, 0, -6)
	return start.Format("2006-01-02"), end.Format("2006-01-02")
}

func countDays(days map[string]struct{}, start, end string) int {
	n := 0
	for day := range days {
		if day >= start && day <= end {
			n++
		}
	}
	return n
}

func protectedStreak(created time.Time, loc *time.Location, protected, flagged map[string]struct{}, now time.Time) int {
	if loc == nil {
		loc = time.UTC
	}
	local := now.In(loc)
	today := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	createdLocal := created.In(loc)
	createdDay := time.Date(createdLocal.Year(), createdLocal.Month(), createdLocal.Day(), 0, 0, 0, 0, loc)
	key := func(day time.Time) string { return day.Format("2006-01-02") }
	clean := func(day time.Time) bool {
		if day.Before(createdDay) {
			return false
		}
		if _, ok := flagged[key(day)]; ok {
			return false
		}
		_, ok := protected[key(day)]
		return ok
	}
	start := today
	if _, flaggedToday := flagged[key(today)]; !flaggedToday {
		if _, ok := protected[key(today)]; !ok {
			start = today.AddDate(0, 0, -1)
		}
	}
	streak := 0
	for day := start; !day.Before(createdDay) && streak <= 3650; day = day.AddDate(0, 0, -1) {
		if !clean(day) {
			break
		}
		streak++
	}
	return streak
}

func validUnsubToken(token string) bool {
	if len(token) != 64 {
		return false
	}
	for _, c := range token {
		switch {
		case c >= '0' && c <= '9', c >= 'a' && c <= 'f':
		default:
			return false
		}
	}
	return true
}

func newUnsubToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// UnsubscribeWeeklySummary stops weekly emails for the partner who opened the link.
// GET /weekly-summary/unsubscribe?token=...
func (h *H) UnsubscribeWeeklySummary(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if !validUnsubToken(token) {
		writeWeeklyPage(w, http.StatusBadRequest, "This link is not valid.")
		return
	}
	var id int64
	err := h.DB.QueryRowContext(r.Context(), `
		SELECT id FROM users WHERE weekly_summary_token = $1
	`, token).Scan(&id)
	if err != nil {
		writeWeeklyPage(w, http.StatusBadRequest, "This link is not valid.")
		return
	}
	if _, err := h.DB.ExecContext(r.Context(), `
		UPDATE users SET weekly_summary_opt_out = TRUE WHERE id = $1
	`, id); err != nil {
		writeWeeklyPage(w, http.StatusInternalServerError, "Could not save that choice. Try the link again.")
		return
	}
	writeWeeklyPage(w, http.StatusOK, "You will no longer get weekly partner summaries.")
}

func writeWeeklyPage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("<!doctype html><meta charset=\"utf-8\"><title>Weekly summary</title><p>" + htmlEscape(message) + "</p>"))
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return replacer.Replace(s)
}

// SetWeeklySummaryPreference lets a signed-in partner opt out or back in.
// POST /users/me/weekly-summary
// Body: {"opt_out": true}
func (h *H) SetWeeklySummaryPreference(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		OptOut *bool `json:"opt_out"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OptOut == nil {
		writeError(w, http.StatusBadRequest, "opt_out is required")
		return
	}
	if _, err := h.DB.ExecContext(r.Context(), `
		UPDATE users SET weekly_summary_opt_out = $2 WHERE id = $1
	`, userID, *req.OptOut); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save preference")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"opt_out": *req.OptOut})
}
