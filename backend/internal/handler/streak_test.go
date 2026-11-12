package handler

import (
	"testing"
	"time"
)

func TestCanonicalZone(t *testing.T) {
	name, ok := canonicalZone(" America/Chicago ")
	if !ok || name != "America/Chicago" {
		t.Fatalf("zone %q ok=%v", name, ok)
	}
	for _, raw := range []string{"", "Local", "local", "Not/AZone", "../etc/passwd", "UTC;drop"} {
		if _, ok := canonicalZone(raw); ok {
			t.Fatalf("accepted %q", raw)
		}
	}
	if name, ok := canonicalZone("UTC"); !ok || name != "UTC" {
		t.Fatalf("utc %q ok=%v", name, ok)
	}
}

func TestProtectionDayUsesUserZone(t *testing.T) {
	at := time.Date(2026, 11, 12, 2, 30, 0, 0, time.UTC)
	if got := protectionDay(at, "America/New_York"); got != "2026-11-11" {
		t.Fatalf("new york day %s", got)
	}
	if got := protectionDay(at, "UTC"); got != "2026-11-12" {
		t.Fatalf("utc day %s", got)
	}
	if got := protectionDay(at, "nope"); got != "2026-11-12" {
		t.Fatalf("unknown zone fell through to %s", got)
	}
}

func TestStreakCountsOnlyProtectedDays(t *testing.T) {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 12, 15, 0, 0, 0, loc)
	created := time.Date(2026, 11, 8, 9, 0, 0, 0, loc)
	protected := days("2026-11-08", "2026-11-09", "2026-11-11", "2026-11-12")
	flagged := days("2026-11-11")

	got := computeStreak(streakMeta{Created: created, Zone: "America/Chicago"}, protected, flagged, now)
	if got.Days != 1 {
		t.Fatalf("streak %d, want 1 (today only; yesterday was flagged)", got.Days)
	}
	if got.Best != 2 {
		t.Fatalf("best %d, want 2 (Nov 8 and Nov 9)", got.Best)
	}
	if got.Coverage[6] != dayClean || got.Coverage[5] != dayFlagged || got.Coverage[4] != dayOff {
		t.Fatalf("coverage %v", got.Coverage)
	}
	if got.Coverage[0] != dayBefore || got.Coverage[1] != dayBefore {
		t.Fatalf("days before the account %v", got.Coverage)
	}
	if got.States[4] != dayClean {
		t.Fatalf("older week word for an off day is %s", got.States[4])
	}
	if got.Week[4] {
		t.Fatal("off day must not be a clean bool")
	}
	if !got.Week[6] {
		t.Fatal("today should be a clean bool")
	}
}

func TestStreakSkipsTodayUntilCheckIn(t *testing.T) {
	protected := map[string]struct{}{}
	for d := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC); d.Before(time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)); d = d.AddDate(0, 0, 1) {
		protected[d.Format("2006-01-02")] = struct{}{}
	}
	now := time.Date(2026, 11, 12, 6, 30, 0, 0, time.UTC)
	got := computeStreak(streakMeta{
		Created: time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
		Zone:    "UTC",
	}, protected, nil, now)
	if got.Days != 11 {
		t.Fatalf("streak %d, want 11 days through yesterday while today has not checked in", got.Days)
	}
	if got.Coverage[6] != dayOff {
		t.Fatalf("today coverage %s", got.Coverage[6])
	}
}

func TestStreakSignupDayCountsWhenProtected(t *testing.T) {
	now := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	created := time.Date(2026, 11, 12, 8, 0, 0, 0, time.UTC)
	got := computeStreak(streakMeta{Created: created, Zone: "UTC"}, days("2026-11-12"), nil, now)
	if got.Days != 1 || got.Best != 1 {
		t.Fatalf("signup day streak=%d best=%d", got.Days, got.Best)
	}
	quiet := computeStreak(streakMeta{Created: created, Zone: "UTC"}, nil, nil, now)
	if quiet.Days != 0 {
		t.Fatalf("signup day without a check-in counted %d", quiet.Days)
	}
}

func TestOffDayBreaksTheRun(t *testing.T) {
	now := time.Date(2026, 11, 12, 20, 0, 0, 0, time.UTC)
	created := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	protected := days("2026-11-10", "2026-11-12")
	got := computeStreak(streakMeta{Created: created, Zone: "UTC"}, protected, nil, now)
	if got.Days != 1 {
		t.Fatalf("streak %d, want 1 because Nov 11 had no check-in", got.Days)
	}
}

func days(list ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(list))
	for _, day := range list {
		out[day] = struct{}{}
	}
	return out
}
