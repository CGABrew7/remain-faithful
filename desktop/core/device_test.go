package core

import (
	"strings"
	"testing"
	"time"
)

func TestDeviceLabelsStayOnTheAllowlist(t *testing.T) {
	if kind, ok := DeviceKind("windows"); !ok || kind != KindWindows {
		t.Fatal(kind, ok)
	}
	if kind, ok := DeviceKind("mac"); !ok || kind != KindMac {
		t.Fatal(kind, ok)
	}
	if _, ok := DeviceKind("other"); ok {
		t.Fatal("unknown platform must not claim a device type")
	}
	for _, kind := range []string{KindiPhone, KindMac, KindWindows, KindAndroid} {
		label, ok := DeviceLabel(kind)
		if !ok || label == "" || strings.Contains(label, ".") || strings.Contains(strings.ToLower(label), "http") {
			t.Fatalf("%s -> %q", kind, label)
		}
	}
	for _, kind := range []string{"Safari", "chrome.exe", "com.apple.mobilesafari", "Windows PC"} {
		if _, ok := DeviceLabel(kind); ok {
			t.Fatalf("accepted %q", kind)
		}
		if ValidDeviceID(kind) && strings.Contains(kind, ".") {
			t.Fatal(kind)
		}
	}
	if !ValidDeviceID(NewDeviceID()) {
		t.Fatal("new id")
	}
}

func TestFormatCheckIn(t *testing.T) {
	now := time.Date(2026, 10, 24, 15, 4, 0, 0, time.UTC)
	if FormatCheckIn("", now) != "No check-in yet" {
		t.Fatal("empty")
	}
	if FormatCheckIn("2026-10-24T15:04:00Z", now) != "Today 3:04 PM" {
		t.Fatal(FormatCheckIn("2026-10-24T15:04:00Z", now))
	}
	if FormatCheckIn("2026-10-23T15:04:00Z", now) != "Oct 23, 3:04 PM" {
		t.Fatal(FormatCheckIn("2026-10-23T15:04:00Z", now))
	}
}
