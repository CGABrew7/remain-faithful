package admin

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCountSuppression(t *testing.T) {
	cases := []struct {
		name string
		c    Count
		want any
	}{
		{"zero people", People(0), "<3"},
		{"one", People(1), "<3"},
		{"two", People(2), "<3"},
		{"three", People(3), float64(3)},
		{"views hidden by audience", Show(40, 2), "<3"},
		{"views shown", Show(40, 3), float64(40)},
		{"real zero", Zero(), float64(0)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := jsonDecoded(t, tc.c); got != tc.want {
				t.Fatalf("got %#v want %#v", got, tc.want)
			}
		})
	}
}

func jsonDecoded(t *testing.T, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestRateHidesSmallSide(t *testing.T) {
	if got := RateAudience(1, 10); got >= 3 {
		t.Fatalf("a single person on one side must hide the rate, audience %d", got)
	}
	if got := RateAudience(8, 10); got >= 3 {
		t.Fatalf("two people on the other side must hide the rate, audience %d", got)
	}
	if got := RateAudience(10, 10); got < 3 {
		t.Fatalf("everyone in a large group can show 100 percent, audience %d", got)
	}
	if got := RateAudience(0, 10); got < 3 {
		t.Fatalf("a zero side is safe to show, audience %d", got)
	}

	if got := jsonDecoded(t, Percent{People: RateAudience(2, 10), Numerator: 2, Denominator: 10}); got != "<3" {
		t.Fatalf("hidden rate = %#v", got)
	}
	if got := jsonDecoded(t, Percent{People: RateAudience(0, 10), Numerator: 0, Denominator: 10}); got != float64(0) {
		t.Fatalf("zero rate = %#v", got)
	}
}

func TestFunnelDoesNotRevealSmallStep(t *testing.T) {
	steps := BuildFunnel(10, 2, 2, 1, 1, 1)
	b, err := json.Marshal(steps)
	if err != nil {
		t.Fatal(err)
	}
	var decoded []struct {
		Key    string `json:"key"`
		Count  any    `json:"count"`
		Drop   any    `json:"dropoff_percent"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, step := range decoded {
		if step.Status != "tracked" {
			continue
		}
		if step.Key == "account_created" {
			if step.Count != float64(10) {
				t.Fatalf("account count = %#v", step.Count)
			}
			continue
		}
		if step.Count != "<3" {
			t.Fatalf("%s count leaked as %#v", step.Key, step.Count)
		}
		if step.Drop != nil && step.Drop != "<3" {
			t.Fatalf("%s drop-off leaked as %#v", step.Key, step.Drop)
		}
	}
	var screen, apps bool
	for _, step := range steps {
		if step.Key == "screen_time" {
			if step.Status != "tracked" || step.Count == nil {
				t.Fatalf("screen time should be a count, %+v", step)
			}
			screen = strings.Contains(step.Note, "Screen Time")
		}
		if step.Key == "apps_chosen" {
			if step.Status != "tracked" || step.Count == nil {
				t.Fatalf("apps chosen should be a count, %+v", step)
			}
			apps = strings.Contains(step.Note, "App names")
		}
	}
	if !screen || !apps {
		t.Fatal("funnel notes should describe Screen Time and chosen apps")
	}
}

func TestAverageSuppression(t *testing.T) {
	if got := jsonDecoded(t, Average{People: 4, Groups: 2, Sum: 4}); got != "<3" {
		t.Fatalf("small group average = %#v", got)
	}
	if got := jsonDecoded(t, Average{People: 12, Groups: 4, Sum: 10}); got != 2.5 {
		t.Fatalf("average = %#v", got)
	}
}
