package core

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Device kinds the server will name in a partner alert.
// Keep this aligned with backend deviceKindLabels.
const (
	KindiPhone  = "iphone"
	KindMac     = "mac"
	KindWindows = "windows"
	KindAndroid = "android"
)

func DeviceKind(platform string) (string, bool) {
	switch platform {
	case KindMac, KindWindows:
		return platform, true
	default:
		return "", false
	}
}

func DeviceLabel(kind string) (string, bool) {
	switch kind {
	case KindiPhone:
		return "iPhone", true
	case KindMac:
		return "Mac", true
	case KindWindows:
		return "Windows PC", true
	case KindAndroid:
		return "Android phone", true
	default:
		return "", false
	}
}

// ValidDeviceID accepts the opaque token stored on this computer.
// Dots and slashes are rejected so an app id cannot be used as the token.
func ValidDeviceID(id string) bool {
	if id == "" || len(id) > 64 {
		return false
	}
	for _, c := range id {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '_', c == '-':
		default:
			return false
		}
	}
	return true
}

func NewDeviceID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		t := time.Now().UnixNano()
		for i := 0; i < 8; i++ {
			b[i] = byte(t >> (uint(i) * 8))
		}
	}
	return hex.EncodeToString(b[:])
}

func FormatCheckIn(raw string, now time.Time) string {
	raw = trimSpace(raw)
	if raw == "" {
		return "No check-in yet"
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return "No check-in yet"
	}
	local := t.In(now.Location())
	clock := local.Format("3:04 PM")
	if local.Year() == now.Year() && local.YearDay() == now.YearDay() {
		return "Today " + clock
	}
	return local.Format("Jan 2") + ", " + clock
}

func trimSpace(s string) string {
	out := make([]rune, 0, len(s))
	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		out = append(out, c)
	}
	return string(out)
}
