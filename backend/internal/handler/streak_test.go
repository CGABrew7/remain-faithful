package handler

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestLegacyStreakMatchesTheOldUTCCount(t *testing.T) {
	now := time.Date(2026, 11, 12, 15, 0, 0, 0, time.UTC)
	created := time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC)
	flagged := map[string]struct{}{"2026-11-11": {}}
	if got := legacyStreakUTC(created, now, flagged); got != 1 {
		t.Fatalf("streak %d, want 1 (today only)", got)
	}
	if got := legacyStreakUTC(now, now, nil); got != 0 {
		t.Fatalf("signup day %d", got)
	}
	if got := legacyStreakUTC(created, now, nil); got != 5 {
		t.Fatalf("open streak %d", got)
	}
}

func TestBaselineSurvivesLaunchDay(t *testing.T) {
	now := time.Date(2026, 11, 12, 15, 0, 0, 0, time.UTC)
	until := time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC)
	got := computeAccountStreak(streakMeta{
		Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Legacy:  40,
	}, nil, nil, nil, until, 30*time.Minute, -300, now)
	if got.Days != 40 {
		t.Fatalf("launch streak %d, want the baseline 40", got.Days)
	}
	if got.Best < 40 {
		t.Fatalf("best %d", got.Best)
	}
}

func TestFailedNewDayDoesNotKeepTheBaseline(t *testing.T) {
	now := time.Date(2026, 11, 14, 18, 0, 0, 0, time.UTC)
	until := time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC)
	got := computeAccountStreak(streakMeta{
		Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Legacy:  40,
	}, nil, nil, nil, until, 30*time.Minute, 0, now)
	if got.Days != 0 {
		t.Fatalf("streak %d, want 0 after an uncovered day", got.Days)
	}
}

func TestSparseCheckInsExtendTheBaseline(t *testing.T) {
	until := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, 11, 14, 22, 0, 0, 0, time.UTC)
	var beats []time.Time
	for day := 12; day <= 14; day++ {
		base := time.Date(2026, 11, day, 0, 0, 0, 0, time.UTC)
		beats = append(beats, base.Add(8*time.Hour), base.Add(12*time.Hour), base.Add(21*time.Hour))
	}
	got := computeAccountStreak(streakMeta{
		Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Legacy:  20,
	}, beats, nil, nil, until, 30*time.Minute, 0, now)
	if got.Days != 23 {
		t.Fatalf("streak %d, want 23", got.Days)
	}
}

func TestOneCheckInThenOffIsNotClean(t *testing.T) {
	now := time.Date(2026, 11, 14, 20, 0, 0, 0, time.UTC)
	day := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)
	got := computeAccountStreak(streakMeta{
		Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Legacy:  20,
	}, []time.Time{day.Add(8 * time.Hour)}, nil, []time.Time{day.Add(15 * time.Hour)},
		time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC), 30*time.Minute, 0, now)
	if got.Coverage[6] != dayOff {
		t.Fatalf("coverage %s, want off", got.Coverage[6])
	}
}

func TestFlaggedDayIsNotClean(t *testing.T) {
	now := time.Date(2026, 11, 14, 20, 0, 0, 0, time.UTC)
	day := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)
	got := computeAccountStreak(streakMeta{
		Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Legacy:  20,
	}, []time.Time{day.Add(8 * time.Hour)}, []time.Time{day.Add(9 * time.Hour)}, nil,
		time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC), 30*time.Minute, 0, now)
	if got.Coverage[6] != dayFlagged || got.Week[6] {
		t.Fatalf("coverage %s week %v", got.Coverage[6], got.Week[6])
	}
}

func TestOffsetChangeCannotCountADayWithNoCheckIn(t *testing.T) {
	beat := time.Date(2026, 11, 12, 2, 0, 0, 0, time.UTC)
	until := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, 11, 12, 18, 0, 0, 0, time.UTC)
	meta := streakMeta{Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Legacy: 5}
	utc := computeAccountStreak(meta, []time.Time{beat}, nil, nil, until, 30*time.Minute, 0, now)
	if utc.Coverage[6] != dayClean {
		t.Fatalf("utc today %s", utc.Coverage[6])
	}
	if utc.Coverage[5] == dayClean {
		t.Fatal("the previous UTC day has no check-in")
	}
	shifted := computeAccountStreak(meta, []time.Time{beat}, nil, nil, until, 30*time.Minute, -300, now)
	if shifted.Coverage[6] == dayClean {
		t.Fatal("a different offset must not mark a day with no check-in as clean")
	}
}

func TestOffEventVoidsTheDay(t *testing.T) {
	now := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	until := time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC)
	start := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)
	var beats []time.Time
	for t := start; !t.After(now); t = t.Add(10 * time.Minute) {
		beats = append(beats, t)
	}
	off := start.Add(2 * time.Hour)
	got := computeAccountStreak(streakMeta{
		Created: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Legacy:  10,
	}, beats, nil, []time.Time{off}, until, 30*time.Minute, 0, now)
	if got.Days != 0 {
		t.Fatalf("streak %d, an off event should end the run", got.Days)
	}
	if got.Coverage[6] != dayOff || got.States[6] != dayNone {
		t.Fatalf("coverage %s states %s", got.Coverage[6], got.States[6])
	}
	if got.Week[6] {
		t.Fatal("an off day must not be a check in the old week field")
	}
}

func TestOffsetIsNotPersistedAndOutOfRangeIsIgnored(t *testing.T) {
	req := httptest.NewRequest("GET", "/users/me/stats?offset_minutes=-300", nil)
	if requestOffsetMinutes(req) != -300 {
		t.Fatal(requestOffsetMinutes(req))
	}
	req = httptest.NewRequest("GET", "/users/me/stats?offset_minutes=9999", nil)
	if requestOffsetMinutes(req) != 0 {
		t.Fatal("out of range offset should be ignored")
	}
	req = httptest.NewRequest("GET", "/users/me/stats", nil)
	if requestOffsetMinutes(req) != 0 {
		t.Fatal("missing offset")
	}
}
