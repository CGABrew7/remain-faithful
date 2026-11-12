package handler

import (
	"testing"
	"time"
)

// Review tests for PR #68. They describe how a real iPhone behaves.
// An iPhone cannot check in every 30 minutes around the clock. It checks in
// when the app is opened, when iOS grants a background refresh, and when a
// Screen Time event fires. A sleeping phone is silent for hours.

// reviewLaunch is legacy_until: the first day the new rule applies.
var reviewLaunch = time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC)

// A user keeps protection on all day and never turns it off. The phone checks
// in three times. That day must count as clean.
func TestReviewOrdinaryPhoneDayIsClean(t *testing.T) {
	now := time.Date(2026, 11, 15, 9, 0, 0, 0, time.UTC)
	created := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	day := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)
	// Morning: app opened. Lunch: background refresh. Evening: app opened.
	// Today: app opened ten minutes ago.
	beats := []time.Time{
		day.Add(7 * time.Hour),
		day.Add(12*time.Hour + 30*time.Minute),
		day.Add(21 * time.Hour),
		now.Add(-10 * time.Minute),
	}
	got := computeAccountStreak(streakMeta{Created: created, Legacy: 20},
		beats, nil, nil, reviewLaunch, 30*time.Minute, 0, now)
	if got.Coverage[5] != dayClean {
		t.Errorf("Nov 14 had three check-ins and no off notice: got %q, want clean", got.Coverage[5])
	}
}

// The same user had a 20-day streak before launch. Two ordinary protected
// days later, the streak must still be alive.
func TestReviewStreakSurvivesFirstDaysAfterLaunch(t *testing.T) {
	now := time.Date(2026, 11, 15, 9, 0, 0, 0, time.UTC)
	created := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	var beats []time.Time
	for _, d := range []int{13, 14} {
		day := time.Date(2026, 11, d, 0, 0, 0, 0, time.UTC)
		beats = append(beats, day.Add(7*time.Hour), day.Add(13*time.Hour), day.Add(21*time.Hour))
	}
	beats = append(beats, now.Add(-10*time.Minute))
	got := computeAccountStreak(streakMeta{Created: created, Legacy: 20},
		beats, nil, nil, reviewLaunch, 30*time.Minute, 0, now)
	if got.Days < 20 {
		t.Errorf("streak = %d after two protected days, want at least the 20-day baseline", got.Days)
	}
}

// One check-in cannot be moved onto another day by claiming a different
// offset on the stats read: the server stores the instant, not a day label.
func TestReviewOffsetCannotMoveACheckInToYesterday(t *testing.T) {
	beat := time.Date(2026, 11, 14, 11, 0, 0, 0, time.UTC) // 6 am Eastern, Nov 14
	for _, off := range []int{-12 * 60, -5 * 60, 0, 14 * 60} {
		midnight := localMidnight(beat, off)
		yesterday := midnight.AddDate(0, 0, -1)
		if instantIn([]time.Time{beat}, yesterday, midnight) {
			t.Errorf("offset %d places the check-in on the previous local day", off)
		}
	}
}

// A protection-off notice still voids the day, however many check-ins arrive.
// A protection-off notice still voids the day, however many check-ins arrive.
func TestReviewTurningProtectionOffVoidsTheDay(t *testing.T) {
	now := time.Date(2026, 11, 15, 9, 0, 0, 0, time.UTC)
	created := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	day := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)
	var beats []time.Time
	for m := 0; m < 24*60; m += 20 {
		beats = append(beats, day.Add(time.Duration(m)*time.Minute))
	}
	offs := []time.Time{day.Add(15 * time.Hour)}
	got := computeAccountStreak(streakMeta{Created: created},
		beats, nil, offs, reviewLaunch, 30*time.Minute, 0, now)
	if got.Coverage[5] != dayOff {
		t.Errorf("day with an off notice: got %q, want off", got.Coverage[5])
	}
}
