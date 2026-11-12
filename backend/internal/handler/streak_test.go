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

func TestWholeDayNeedsAChainAndIgnoresASingleCheckIn(t *testing.T) {
	gap := 30 * time.Minute
	// 02:30 UTC is 21:30 on the previous calendar day at UTC-5.
	beat := time.Date(2026, 11, 12, 2, 30, 0, 0, time.UTC)
	start := time.Date(2026, 11, 11, 5, 0, 0, 0, time.UTC) // local midnight, offset -300
	end := start.Add(24 * time.Hour)
	now := end.Add(time.Hour)
	clean, pending := dayCovered([]time.Time{beat}, nil, start, end, now, gap)
	if clean || pending {
		t.Fatal("one evening check-in must not cover the whole day")
	}
	nextClean, _ := dayCovered([]time.Time{beat}, nil, end, end.Add(24*time.Hour), now, gap)
	if nextClean {
		t.Fatal("that check-in must not cover the next day")
	}
}

func TestCarryAndGaps(t *testing.T) {
	gap := 30 * time.Minute
	start := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	now := end.Add(time.Hour)
	var beats []time.Time
	for t := start.Add(-10 * time.Minute); !t.After(end); t = t.Add(20 * time.Minute) {
		beats = append(beats, t)
	}
	clean, pending := dayCovered(beats, nil, start, end, now, gap)
	if !clean || pending {
		t.Fatal("a 20 minute chain carried from before midnight should cover the day")
	}
	beats = beats[:len(beats)-1]
	// Drop the middle of the day.
	var gapped []time.Time
	for _, b := range beats {
		if b.Hour() == 12 {
			continue
		}
		gapped = append(gapped, b)
	}
	clean, _ = dayCovered(gapped, nil, start, end, now, gap)
	if clean {
		t.Fatal("a gap longer than the heartbeat interval is not a clean day")
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
