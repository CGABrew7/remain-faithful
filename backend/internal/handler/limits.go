package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

const (
	maxPartnersPerUser = 10
	maxInvitesPerDay   = 20
	eventPastWindow    = 7 * 24 * time.Hour
	eventFutureSkew    = 5 * time.Minute
	resetMailPerHour   = 3
	silenceRepeatAfter = 6 * time.Hour
	eventRetention     = 365 * 24 * time.Hour
)

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// parseEventTimestamp accepts the ISO-8601 values the iOS client sends.
// Empty means "use the server clock". A timestamp may be up to 7 days old
// so an offline queue can still upload, and at most a few minutes in the future.
func parseEventTimestamp(raw string, now time.Time) (time.Time, string) {
	raw = trimSpace(raw)
	if raw == "" {
		return time.Time{}, ""
	}
	var t time.Time
	var err error
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		t, err = time.Parse(layout, raw)
		if err == nil {
			break
		}
	}
	if err != nil {
		return time.Time{}, "invalid timestamp"
	}
	if t.After(now.Add(eventFutureSkew)) {
		return time.Time{}, "timestamp is too far in the future"
	}
	if t.Before(now.Add(-eventPastWindow)) {
		return time.Time{}, "timestamp is too old"
	}
	return t.UTC(), ""
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\n' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}

func (h *H) invitesSentToday(ctx context.Context, userID int64) (int, error) {
	var n int
	err := h.DB.QueryRowContext(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM relationship_invites WHERE inviter_id = $1 AND created_at > NOW() - INTERVAL '1 day')
		  + (SELECT COUNT(*) FROM group_invites WHERE inviter_id = $1 AND created_at > NOW() - INTERVAL '1 day')
		)`, userID).Scan(&n)
	return n, err
}

func (h *H) enforceInviteBudget(ctx context.Context, userID int64) (int, string) {
	n, err := h.invitesSentToday(ctx, userID)
	if err != nil {
		return http.StatusInternalServerError, "could not check invite limit"
	}
	if n >= maxInvitesPerDay {
		return http.StatusTooManyRequests, "invite limit reached for today"
	}
	return 0, ""
}

func (h *H) partnerCount(ctx context.Context, userID int64) (int, error) {
	var n int
	err := h.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM relationships WHERE user_id = $1`, userID,
	).Scan(&n)
	return n, err
}
