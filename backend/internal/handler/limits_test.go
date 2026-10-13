package handler

import (
	"strings"
	"testing"
	"time"
)

func TestParseEventTimestampWindow(t *testing.T) {
	now := time.Date(2026, 10, 13, 12, 0, 0, 0, time.UTC)
	if _, msg := parseEventTimestamp("", now); msg != "" {
		t.Fatalf("empty = %q", msg)
	}
	ok := now.Add(-2 * time.Hour).Format(time.RFC3339)
	got, msg := parseEventTimestamp(ok, now)
	if msg != "" || got.IsZero() {
		t.Fatalf("recent = %v %q", got, msg)
	}
	future := now.Add(time.Hour).Format(time.RFC3339)
	if _, msg := parseEventTimestamp(future, now); msg == "" {
		t.Fatal("future timestamp was accepted")
	}
	old := now.Add(-8 * 24 * time.Hour).Format(time.RFC3339)
	if _, msg := parseEventTimestamp(old, now); !strings.Contains(msg, "too old") {
		t.Fatalf("old = %q", msg)
	}
	if _, msg := parseEventTimestamp("not-a-time", now); msg != "invalid timestamp" {
		t.Fatalf("junk = %q", msg)
	}
}

func TestHashTokenIsStableAndNotPlaintext(t *testing.T) {
	a := hashToken("abc")
	b := hashToken("abc")
	if a != b || a == "abc" || len(a) != 64 {
		t.Fatalf("hash = %q", a)
	}
	if hashToken("abd") == a {
		t.Fatal("different tokens hashed the same")
	}
}

func TestCaps(t *testing.T) {
	if maxGroupMembers != 12 {
		t.Fatalf("group cap = %d", maxGroupMembers)
	}
	if maxPartnersPerUser != 10 || maxInvitesPerDay != 20 {
		t.Fatalf("partner cap %d invite cap %d", maxPartnersPerUser, maxInvitesPerDay)
	}
}
