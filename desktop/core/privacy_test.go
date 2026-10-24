package core

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRejectsPagesAndPicturesAndAppNames(t *testing.T) {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	bad := []OutboundEvent{
		{Category: "adult_content", Severity: "severe", Summary: "https://example.com/page", Timestamp: now.Format(time.RFC3339)},
		{Category: "app_usage", Severity: "concerning", Summary: "Safari was used", Timestamp: now.Format(time.RFC3339)},
		{Category: "app_usage", Severity: "concerning", Summary: "chrome.exe opened", Timestamp: now.Format(time.RFC3339)},
		{Category: "nope", Severity: "concerning", Summary: SummaryApp, Timestamp: now.Format(time.RFC3339)},
		{Category: "app_usage", Severity: "high", Summary: SummaryApp, Timestamp: now.Format(time.RFC3339)},
	}
	for _, e := range bad {
		if err := ValidateEvent(e); err == nil {
			t.Fatalf("accepted %+v", e)
		}
	}
	if _, err := UsageEvent(now, "Monitored app used for 1+ minute and the address was https://x"); err == nil {
		t.Fatal("usage helper accepted a page")
	}
}

func TestProtectionAllowlist(t *testing.T) {
	if err := ValidateProtection(ProtectionNotice{Type: "heartbeat_silence", Detail: "no"}); err == nil {
		t.Fatal("silence is server only")
	}
	if err := ValidateProtection(ProtectionNotice{Type: "monitoring_disabled", Detail: "Opened https://example.com"}); err == nil {
		t.Fatal("detail leaked a page")
	}
	if err := ValidateProtection(ProtectionNotice{Type: "monitoring_disabled", Detail: "Alerts were turned off."}); err != nil {
		t.Fatal(err)
	}
}

func TestEventBodyIsMetadataOnly(t *testing.T) {
	e, err := UsageEvent(time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC), SummaryCategory)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := e.Body()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if len(m) != 4 {
		t.Fatalf("keys %v", m)
	}
	for _, key := range []string{"category", "severity", "summary", "timestamp"} {
		if _, ok := m[key]; !ok {
			t.Fatalf("missing %s", key)
		}
	}
}

func TestSummaryMatchesPhone(t *testing.T) {
	if !summaryAllowed(SummaryApp) || !summaryAllowed(SummaryCategory) {
		t.Fatal("desktop summaries must match the server list")
	}
	if summaryAllowed("Flagged during browsing session " + strings.Repeat("x", 40)) {
		t.Fatal("free text is not a summary")
	}
}

func TestNetworkVersusAPI(t *testing.T) {
	if !IsNetwork(errors.New("dial failed")) {
		t.Fatal("transport errors are network")
	}
	if IsNetwork(&APIError{Status: 400, Message: "no"}) {
		t.Fatal("api errors are not retried as network")
	}
}
