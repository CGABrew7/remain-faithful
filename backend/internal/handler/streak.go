package handler

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/lib/pq"
)

const (
	dayBefore  = "before_account"
	dayClean   = "clean"
	dayFlagged = "flagged"
	dayOff     = "off"
)

// streakMeta is the account clock a streak is counted in.
type streakMeta struct {
	Created time.Time
	Zone    string
}

// streakResult is the home streak and the week row.
// States keeps the words shipped iOS builds already understand, so a day
// with protection off stays out of the number without turning that dot into
// a flag on the build currently in review. Coverage is the accurate row.
type streakResult struct {
	Days     int
	Best     int
	Week     []bool
	States   []string
	Coverage []string
}

// MigrateStreak stores each account's time zone and the local days protection
// checked in. The latest heartbeat, if any, is recorded once so the current
// day is not dropped on the deploy that adds the table.
func MigrateStreak(db *sql.DB) error {
	_, err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS timezone TEXT NOT NULL DEFAULT 'UTC';
		CREATE TABLE IF NOT EXISTS protection_days (
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			day     DATE   NOT NULL,
			PRIMARY KEY (user_id, day)
		);
		INSERT INTO protection_days (user_id, day)
		SELECT id, (last_heartbeat_at AT TIME ZONE 'UTC')::date
		FROM users
		WHERE last_heartbeat_at IS NOT NULL
		ON CONFLICT DO NOTHING;
	`)
	return err
}

func canonicalZone(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || len(raw) > 64 || strings.EqualFold(raw, "Local") {
		return "", false
	}
	loc, err := time.LoadLocation(raw)
	if err != nil {
		return "", false
	}
	name := loc.String()
	if name == "" || name == "Local" {
		return "", false
	}
	return name, true
}

func zoneLocation(raw string) *time.Location {
	if name, ok := canonicalZone(raw); ok {
		loc, err := time.LoadLocation(name)
		if err == nil {
			return loc
		}
	}
	return time.UTC
}

// protectionDay is the calendar date of at in the user's time zone.
func protectionDay(at time.Time, zone string) string {
	return at.In(zoneLocation(zone)).Format("2006-01-02")
}

// noteProtectionDay records that protection checked in, and stores a time zone
// when the device sent one. A failure is logged and does not fail the heartbeat.
func (h *H) noteProtectionDay(ctx context.Context, userID int64, rawZone string, at time.Time) {
	zone := "UTC"
	if name, ok := canonicalZone(rawZone); ok {
		zone = name
		if _, err := h.DB.ExecContext(ctx,
			`UPDATE users SET timezone = $2 WHERE id = $1`, userID, zone,
		); err != nil {
			log.Printf("[streak] timezone user=%d: %v", userID, err)
		}
	} else if h.DB != nil {
		var stored string
		err := h.DB.QueryRowContext(ctx,
			`SELECT COALESCE(timezone, 'UTC') FROM users WHERE id = $1`, userID,
		).Scan(&stored)
		if err == nil {
			if name, ok := canonicalZone(stored); ok {
				zone = name
			}
		}
	}
	if h.DB == nil {
		return
	}
	if _, err := h.DB.ExecContext(ctx, `
		INSERT INTO protection_days (user_id, day)
		VALUES ($1, $2::date)
		ON CONFLICT DO NOTHING
	`, userID, protectionDay(at, zone)); err != nil {
		log.Printf("[streak] protection day user=%d: %v", userID, err)
	}
}

func (h *H) streaksAt(ctx context.Context, metas map[int64]streakMeta, now time.Time) (map[int64]streakResult, error) {
	out := make(map[int64]streakResult, len(metas))
	if len(metas) == 0 {
		return out, nil
	}
	ids := make([]int64, 0, len(metas))
	for id := range metas {
		ids = append(ids, id)
	}
	protected, flagged, err := h.protectionSets(ctx, ids)
	if err != nil {
		return nil, err
	}
	for id, meta := range metas {
		out[id] = computeStreak(meta, protected[id], flagged[id], now)
	}
	return out, nil
}

func (h *H) protectionSets(ctx context.Context, ids []int64) (protected, flagged map[int64]map[string]struct{}, err error) {
	protected = map[int64]map[string]struct{}{}
	flagged = map[int64]map[string]struct{}{}
	rows, err := h.DB.QueryContext(ctx, `
		SELECT user_id, to_char(day, 'YYYY-MM-DD')
		FROM protection_days
		WHERE user_id = ANY($1)
	`, pq.Array(ids))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var day string
		if scanErr := rows.Scan(&id, &day); scanErr != nil {
			continue
		}
		addDay(protected, id, day)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	rows, err = h.DB.QueryContext(ctx, `
		SELECT e.user_id,
		       to_char((e.timestamp AT TIME ZONE COALESCE(NULLIF(u.timezone, ''), 'UTC')), 'YYYY-MM-DD')
		FROM events e
		JOIN users u ON u.id = e.user_id
		WHERE e.user_id = ANY($1)
		  AND e.category <> 'clean'
		  AND e.timestamp >= NOW() - INTERVAL '4000 days'
	`, pq.Array(ids))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var day string
		if scanErr := rows.Scan(&id, &day); scanErr != nil {
			continue
		}
		addDay(flagged, id, day)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	return protected, flagged, nil
}

func addDay(sets map[int64]map[string]struct{}, id int64, day string) {
	day = strings.TrimSpace(day)
	if len(day) >= 10 {
		day = day[:10]
	}
	if day == "" {
		return
	}
	if sets[id] == nil {
		sets[id] = map[string]struct{}{}
	}
	sets[id][day] = struct{}{}
}

// computeStreak counts consecutive clean days. A past day counts only when
// protection checked in and nothing was flagged. Today, before the first
// check-in, is left out of the run so a morning does not wipe the streak.
func computeStreak(meta streakMeta, protected, flagged map[string]struct{}, now time.Time) streakResult {
	loc := zoneLocation(meta.Zone)
	localNow := now.In(loc)
	today := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	createdLocal := meta.Created.In(loc)
	createdDay := time.Date(createdLocal.Year(), createdLocal.Month(), createdLocal.Day(), 0, 0, 0, 0, loc)

	stateOf := func(day time.Time) string {
		if day.Before(createdDay) {
			return dayBefore
		}
		key := day.Format("2006-01-02")
		if _, ok := flagged[key]; ok {
			return dayFlagged
		}
		if _, ok := protected[key]; ok {
			return dayClean
		}
		return dayOff
	}

	start := today
	todayPending := stateOf(today) == dayOff
	if todayPending {
		start = today.AddDate(0, 0, -1)
	}

	streak := 0
	for d := start; !d.Before(createdDay) && streak <= 3650; d = d.AddDate(0, 0, -1) {
		if stateOf(d) != dayClean {
			break
		}
		streak++
	}

	window := 365
	if streak+1 > window {
		window = streak + 1
	}
	best, run := 0, 0
	for i := 0; i < window; i++ {
		day := today.AddDate(0, 0, -i)
		if day.Before(createdDay) {
			break
		}
		if day.Equal(today) && todayPending {
			continue
		}
		if stateOf(day) == dayClean {
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

	week := make([]bool, 7)
	states := make([]string, 7)
	coverage := make([]string, 7)
	for i := 0; i < 7; i++ {
		day := today.AddDate(0, 0, -(6 - i))
		state := stateOf(day)
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
			// Shipped builds treat an unknown week word as a flag. Keep the
			// older word here. Coverage carries "off" for the next build.
			states[i] = dayClean
		}
	}
	return streakResult{Days: streak, Best: best, Week: week, States: states, Coverage: coverage}
}
