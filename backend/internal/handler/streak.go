package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

const (
	dayBefore  = "before_account"
	dayClean   = "clean"
	dayFlagged = "flagged"
	dayOff     = "off"
	dayNone    = "none"
)

// protectionOffTypes are partner notices that protection was turned off.
// A day that contains one is not clean under the post-deploy rule.
var protectionOffTypes = []string{
	"monitoring_disabled",
	"shielding_disabled",
	"lockout_disabled",
	"deep_scan_stopped",
	"signed_out",
	"selection_cleared",
	"family_controls_revoked",
}

type streakMeta struct {
	Created time.Time
	Legacy  int
}

type streakResult struct {
	Days     int
	Best     int
	Week     []bool
	States   []string
	Coverage []string
}

// MigrateStreak keeps each account's current streak as a baseline and
// records UTC heartbeat instants from then on. It does not store a time zone
// or a pre-labeled calendar day.
func MigrateStreak(db *sql.DB) error {
	if _, err := db.Exec(`
		ALTER TABLE users DROP COLUMN IF EXISTS timezone;
		DROP TABLE IF EXISTS protection_days;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS legacy_streak INT NOT NULL DEFAULT 0;
		CREATE TABLE IF NOT EXISTS streak_rule (
			id           SMALLINT PRIMARY KEY,
			legacy_until TIMESTAMPTZ NOT NULL
		);
		CREATE TABLE IF NOT EXISTS protection_heartbeats (
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			at      TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_protection_heartbeats_user_at
			ON protection_heartbeats (user_id, at);
	`); err != nil {
		return err
	}
	var existing int
	if err := db.QueryRow(`SELECT COUNT(*) FROM streak_rule WHERE id = 1`).Scan(&existing); err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}
	now := time.Now().UTC()
	if err := snapshotLegacyStreaks(db, now); err != nil {
		return err
	}
	until := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	_, err := db.Exec(`
		INSERT INTO streak_rule (id, legacy_until) VALUES (1, $1)
		ON CONFLICT (id) DO NOTHING
	`, until)
	return err
}

func snapshotLegacyStreaks(db *sql.DB, now time.Time) error {
	rows, err := db.Query(`SELECT id, created_at FROM users`)
	if err != nil {
		return err
	}
	defer rows.Close()
	type account struct {
		id      int64
		created time.Time
	}
	var accounts []account
	for rows.Next() {
		var a account
		if err := rows.Scan(&a.id, &a.created); err != nil {
			return err
		}
		accounts = append(accounts, a)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	flagRows, err := db.Query(`
		SELECT user_id, to_char((timestamp AT TIME ZONE 'UTC'), 'YYYY-MM-DD')
		FROM events
		WHERE category <> 'clean'
	`)
	if err != nil {
		return err
	}
	defer flagRows.Close()
	flagged := map[int64]map[string]struct{}{}
	for flagRows.Next() {
		var id int64
		var day string
		if err := flagRows.Scan(&id, &day); err != nil {
			return err
		}
		if len(day) >= 10 {
			day = day[:10]
		}
		if flagged[id] == nil {
			flagged[id] = map[string]struct{}{}
		}
		flagged[id][day] = struct{}{}
	}
	if err := flagRows.Err(); err != nil {
		return err
	}
	for _, a := range accounts {
		n := legacyStreakUTC(a.created, now, flagged[a.id])
		if _, err := db.Exec(`UPDATE users SET legacy_streak = $2 WHERE id = $1`, a.id, n); err != nil {
			return err
		}
	}
	return nil
}

// legacyStreakUTC is the streak the previous server returned: UTC days with
// no flag, including today, and 0 when the account was created today.
func legacyStreakUTC(created, now time.Time, flagged map[string]struct{}) int {
	today := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
	created = created.UTC()
	createdDay := time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, time.UTC)
	if createdDay.Equal(today) {
		return 0
	}
	streak := 0
	for d := today; !d.Before(createdDay) && streak <= 3650; d = d.AddDate(0, 0, -1) {
		if _, ok := flagged[d.Format("2006-01-02")]; ok {
			break
		}
		streak++
	}
	return streak
}

func heartbeatGap() time.Duration {
	raw := strings.TrimSpace(os.Getenv("HEARTBEAT_SILENCE_MINUTES"))
	if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 24*60 {
		return time.Duration(n) * time.Minute
	}
	return 30 * time.Minute
}

func requestOffsetMinutes(r *http.Request) int {
	raw := strings.TrimSpace(r.URL.Query().Get("offset_minutes"))
	if raw == "" {
		return 0
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < -14*60 || n > 14*60 {
		return 0
	}
	return n
}

func offsetLocation(minutes int) *time.Location {
	return time.FixedZone("", minutes*60)
}

func localMidnight(instant time.Time, offsetMin int) time.Time {
	loc := offsetLocation(offsetMin)
	n := instant.In(loc)
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, loc)
}

// recordProtectionHeartbeat stores the UTC instant of a check-in.
// The request offset is not accepted and nothing about the zone is saved.
func (h *H) recordProtectionHeartbeat(ctx context.Context, userID int64, at time.Time) {
	if h == nil || h.DB == nil {
		return
	}
	if _, err := h.DB.ExecContext(ctx, `
		INSERT INTO protection_heartbeats (user_id, at) VALUES ($1, $2)
	`, userID, at.UTC()); err != nil {
		log.Printf("[streak] heartbeat user=%d: %v", userID, err)
	}
}

func (h *H) streaksAt(ctx context.Context, metas map[int64]streakMeta, now time.Time, offsetMin int) (map[int64]streakResult, error) {
	out := make(map[int64]streakResult, len(metas))
	if len(metas) == 0 {
		return out, nil
	}
	var legacyUntil time.Time
	err := h.DB.QueryRowContext(ctx, `SELECT legacy_until FROM streak_rule WHERE id = 1`).Scan(&legacyUntil)
	if err == sql.ErrNoRows {
		legacyUntil = time.Time{}
		err = nil
	}
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(metas))
	for id := range metas {
		ids = append(ids, id)
	}
	gap := heartbeatGap()
	since := now.Add(-400 * 24 * time.Hour)
	if !legacyUntil.IsZero() && legacyUntil.Add(-gap).Before(since) {
		since = legacyUntil.Add(-gap)
	}
	beats, err := h.instants(ctx, `
		SELECT user_id, at FROM protection_heartbeats
		WHERE user_id = ANY($1) AND at >= $2
	`, ids, since)
	if err != nil {
		return nil, err
	}
	flags, err := h.instants(ctx, `
		SELECT user_id, timestamp FROM events
		WHERE user_id = ANY($1) AND category <> 'clean' AND timestamp >= $2
	`, ids, now.Add(-400*24*time.Hour))
	if err != nil {
		return nil, err
	}
	offs, err := h.instants(ctx, `
		SELECT subject_user_id, created_at FROM alerts
		WHERE subject_user_id = ANY($1)
		  AND alert_type = ANY($3)
		  AND created_at >= $2
	`, ids, since, pq.Array(protectionOffTypes))
	if err != nil {
		return nil, err
	}
	for id, meta := range metas {
		out[id] = computeAccountStreak(meta, beats[id], flags[id], offs[id], legacyUntil, gap, offsetMin, now)
	}
	return out, nil
}

func (h *H) instants(ctx context.Context, query string, ids []int64, since time.Time, extra ...any) (map[int64][]time.Time, error) {
	args := []any{pq.Array(ids), since}
	args = append(args, extra...)
	rows, err := h.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64][]time.Time{}
	for rows.Next() {
		var id int64
		var at time.Time
		if err := rows.Scan(&id, &at); err != nil {
			continue
		}
		out[id] = append(out[id], at.UTC())
	}
	return out, rows.Err()
}

func computeAccountStreak(meta streakMeta, beats, flags, offs []time.Time, legacyUntil time.Time, gap time.Duration, offsetMin int, now time.Time) streakResult {
	loc := offsetLocation(offsetMin)
	today := localMidnight(now, offsetMin)
	created := meta.Created.In(loc)
	createdDay := time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, loc)

	classOf := func(day time.Time) (string, bool) {
		if day.Before(createdDay) {
			return dayBefore, false
		}
		end := day.Add(24 * time.Hour)
		if instantIn(flags, day, end) {
			return dayFlagged, false
		}
		if legacyUntil.IsZero() || day.Before(legacyUntil) {
			return dayClean, false
		}
		if instantIn(offs, day, end) {
			return dayOff, false
		}
		clean, pending := dayCovered(beats, offs, day, end, now, gap)
		if pending {
			return dayOff, true
		}
		if clean {
			return dayClean, false
		}
		return dayOff, false
	}

	streak := 0
	broke := false
	for d, i := today, 0; i < 4000; i++ {
		if d.Before(createdDay) {
			break
		}
		if !legacyUntil.IsZero() && d.Before(legacyUntil) {
			break
		}
		state, pending := classOf(d)
		if pending && d.Equal(today) {
			d = d.AddDate(0, 0, -1)
			continue
		}
		if state != dayClean {
			broke = true
			break
		}
		streak++
		d = d.AddDate(0, 0, -1)
	}
	if !broke && !legacyUntil.IsZero() {
		streak += meta.Legacy
	}

	best, run := 0, 0
	for i := 0; i < 366; i++ {
		day := today.AddDate(0, 0, -i)
		if day.Before(createdDay) {
			break
		}
		state, pending := classOf(day)
		if pending && day.Equal(today) {
			continue
		}
		if state == dayClean {
			run++
			if run > best {
				best = run
			}
			continue
		}
		run = 0
	}
	if streak > best {
		best = streak
	}
	if meta.Legacy > best && !legacyUntil.IsZero() {
		best = meta.Legacy
	}

	week := make([]bool, 7)
	states := make([]string, 7)
	coverage := make([]string, 7)
	for i := 0; i < 7; i++ {
		day := today.AddDate(0, 0, -(6 - i))
		state, _ := classOf(day)
		coverage[i] = state
		switch state {
		case dayClean:
			week[i] = true
			states[i] = dayClean
		case dayBefore:
			states[i] = dayBefore
		case dayFlagged:
			states[i] = dayFlagged
		default:
			// Older builds render "none" as an empty mark, not a check.
			states[i] = dayNone
		}
	}
	return streakResult{Days: streak, Best: best, Week: week, States: states, Coverage: coverage}
}

func instantIn(times []time.Time, start, end time.Time) bool {
	for _, t := range times {
		if !t.Before(start) && t.Before(end) {
			return true
		}
	}
	return false
}

// dayCovered is true when heartbeats chain across the whole local day with
// no gap longer than the normal heartbeat interval. A check-in covers only
// the instants within that interval of its UTC timestamp.
func dayCovered(beats, offs []time.Time, start, end, now time.Time, gap time.Duration) (clean, pending bool) {
	requiredEnd := end
	inProgress := !now.Before(start) && now.Before(end)
	if inProgress {
		requiredEnd = now
	}
	pts := make([]time.Time, 0, len(beats)+1)
	if carriedIn(beats, offs, start, gap) {
		pts = append(pts, start)
	}
	for _, b := range beats {
		if !b.Before(start) && !b.After(requiredEnd) {
			pts = append(pts, b)
		}
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].Before(pts[j]) })
	if len(pts) == 0 {
		if inProgress && !now.After(start.Add(gap)) {
			return false, true
		}
		return false, false
	}
	if pts[0].After(start.Add(gap)) {
		return false, false
	}
	for i := 1; i < len(pts); i++ {
		if pts[i].Sub(pts[i-1]) > gap {
			return false, false
		}
	}
	if requiredEnd.Sub(pts[len(pts)-1]) > gap {
		return false, false
	}
	return true, false
}

func carriedIn(beats, offs []time.Time, start time.Time, gap time.Duration) bool {
	var last time.Time
	for _, b := range beats {
		if !b.Before(start.Add(-gap)) && b.Before(start) && b.After(last) {
			last = b
		}
	}
	if last.IsZero() {
		return false
	}
	for _, o := range offs {
		if o.After(last) && o.Before(start) {
			return false
		}
	}
	return true
}
