package core

import "testing"

func TestWeekMarksPreferCoverage(t *testing.T) {
	got := WeekMarks(Stats{
		WeekStates:   []string{"clean", "clean"},
		WeekCoverage: []string{"off", "flagged", "clean"},
	})
	if len(got) != 3 || got[0] != "off" || got[1] != "flagged" || got[2] != "clean" {
		t.Fatalf("%v", got)
	}
	older := WeekMarks(Stats{WeekStates: []string{"before_account", "clean"}})
	if len(older) != 2 || older[0] != "before_account" || older[1] != "clean" {
		t.Fatalf("%v", older)
	}
}
