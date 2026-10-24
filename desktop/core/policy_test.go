package core

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestUsageCrossesOneMinuteOncePerDay(t *testing.T) {
	now := time.Date(2026, 10, 24, 9, 0, 0, 0, time.Local)
	cfg := WatchConfig{Consent: true, Alerts: true, Selected: []string{"chrome.exe"}}
	var ledger Ledger
	var fired int
	for i := 0; i < 5; i++ {
		var out Outcome
		ledger, out = Evaluate(now, "windows", cfg, ledger, "chrome.exe", 20*time.Second)
		if out.Event != nil {
			fired++
			if out.Event.Summary != SummaryCategory {
				t.Fatalf("browser summary %q", out.Event.Summary)
			}
			raw, err := out.Event.Body()
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(strings.ToLower(string(raw)), "chrome") {
				t.Fatalf("payload named the app: %s", raw)
			}
		}
		if out.Block {
			t.Fatal("alerts alone do not block")
		}
	}
	if fired != 1 {
		t.Fatalf("fired %d", fired)
	}
	if ledger.BrowserSeconds != 100 {
		t.Fatalf("seconds %d", ledger.BrowserSeconds)
	}
}

func TestAppsAndBrowsersAreSeparateSets(t *testing.T) {
	now := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	cfg := WatchConfig{Consent: true, Alerts: true, Selected: []string{"chrome.exe", "discord.exe"}}
	ledger, out := Evaluate(now, "windows", cfg, Ledger{}, "discord.exe", 30*time.Second)
	ledger, out = Evaluate(now, "windows", cfg, ledger, "discord.exe", 30*time.Second)
	if out.Event == nil || out.Event.Summary != SummaryApp {
		t.Fatalf("app event %+v", out.Event)
	}
	ledger, out = Evaluate(now, "windows", cfg, ledger, "chrome.exe", 30*time.Second)
	_, out = Evaluate(now, "windows", cfg, ledger, "chrome.exe", 30*time.Second)
	if out.Event == nil || out.Event.Summary != SummaryCategory {
		t.Fatalf("browser event %+v", out.Event)
	}
}

func TestBlockingDoesNotCreateAnEvent(t *testing.T) {
	now := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	cfg := WatchConfig{Consent: true, Blocking: true, Selected: []string{"com.apple.Safari"}}
	_, out := Evaluate(now, "mac", cfg, Ledger{}, "com.apple.Safari", time.Minute)
	if !out.Block || out.BlockKey != "com.apple.Safari" {
		t.Fatalf("block %+v", out)
	}
	if out.Event != nil {
		t.Fatal("blocking does not send a content alert")
	}
}

func TestNoConsentNoWatch(t *testing.T) {
	cfg := WatchConfig{Consent: false, Alerts: true, Blocking: true, Selected: []string{"chrome.exe"}}
	_, out := Evaluate(time.Now(), "windows", cfg, Ledger{}, "chrome.exe", time.Minute)
	if out.Event != nil || out.Block {
		t.Fatal("consent gates watching")
	}
}

func TestSelfAndUnknownAreIgnored(t *testing.T) {
	cfg := WatchConfig{Consent: true, Alerts: true, Blocking: true, Selected: []string{"remainfaithful.exe"}}
	_, out := Evaluate(time.Now(), "windows", cfg, Ledger{}, "remainfaithful.exe", time.Minute)
	if out.Event != nil || out.Block {
		t.Fatal("the app must not watch itself")
	}
	cfg.Selected = []string{"chrome.exe"}
	_, out = Evaluate(time.Now(), "windows", cfg, Ledger{}, "notepad.exe", time.Minute)
	if out.Event != nil || out.Block {
		t.Fatal("unknown apps are ignored")
	}
}

func TestDayResetsTheDailyAlert(t *testing.T) {
	start := time.Date(2026, 10, 24, 23, 0, 0, 0, time.UTC)
	cfg := WatchConfig{Consent: true, Alerts: true, Selected: []string{"slack.exe"}}
	ledger, out := Evaluate(start, "windows", cfg, Ledger{}, "slack.exe", 30*time.Second)
	ledger, out = Evaluate(start, "windows", cfg, ledger, "slack.exe", 30*time.Second)
	if out.Event == nil {
		t.Fatal("expected first alert")
	}
	ledger, out = Evaluate(start.Add(24*time.Hour), "windows", cfg, ledger, "slack.exe", 30*time.Second)
	_, out = Evaluate(start.Add(24*time.Hour), "windows", cfg, ledger, "slack.exe", 30*time.Second)
	if out.Event == nil {
		t.Fatal("a new day can alert again")
	}
}

func TestHeartbeatOnlyWhileArmed(t *testing.T) {
	if HeartbeatScreen(true, false, false, false) != "" {
		t.Fatal("idle desktop must not heartbeat")
	}
	if HeartbeatScreen(false, true, true, true) != "" {
		t.Fatal("consent is required")
	}
	if HeartbeatScreen(true, true, false, false) != "monitoring" {
		t.Fatal("alerts arm the heartbeat")
	}
	if HeartbeatScreen(true, false, true, false) != "monitoring" {
		t.Fatal("blocking arms the heartbeat")
	}
}

func TestPayloadsOmitCatalogIdentity(t *testing.T) {
	event, err := UsageEvent(time.Now(), SummaryApp)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := event.Body()
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ToLower(string(raw))
	for _, entry := range catalog {
		if strings.Contains(text, strings.ToLower(entry.Key)) || strings.Contains(text, strings.ToLower(entry.Display)) {
			t.Fatalf("event contains %s: %s", entry.Display, raw)
		}
	}
	signals := DesktopSignals(3, 1)
	blob, _ := json.Marshal(signals)
	low := strings.ToLower(string(blob))
	if strings.Contains(low, "safari") || strings.Contains(low, ".exe") || signals["screen_time_granted"] != false {
		t.Fatalf("signals %s", blob)
	}
}
