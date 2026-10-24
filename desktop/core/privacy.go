package core

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

// Event fields the existing API accepts. Desktop posts only these.
// Screen text, addresses, screenshots, and app names are not fields on this type.
type OutboundEvent struct {
	Category  string
	Severity  string
	Summary   string
	Timestamp string
}

const (
	SummaryApp      = "Monitored app used for 1+ minute"
	SummaryCategory = "Monitored app category used for 1+ minute"
)

// Summaries the server already accepts from iOS. Desktop policy emits only
// SummaryApp and SummaryCategory. The rest stay here so a value the phone
// already uses is not treated as a new kind of payload.
var allowedEventSummaries = map[string]struct{}{
	"Explicit content detected":                       {},
	"Gambling content detected":                       {},
	"Violent content detected":                        {},
	"Self-harm content detected":                      {},
	"Content reviewed — no concerns":                  {},
	"Explicit image detected (perceptual hash match)": {},
	"Blocked domain detected":                         {},
	"Explicit keyword detected in screen text":        {},
	"Monitored app used for 1+ minute":                {},
	"Monitored app category used for 1+ minute":       {},
	"Monitored app activity detected":                 {},
	"App activity detected":                           {},
	"Blocked app opened":                              {},
	"Adult content flagged on screen":                 {},
	"Gambling flagged on screen":                      {},
	"Violence flagged on screen":                      {},
	"Self-harm content flagged on screen":             {},
	"Content reviewed":                                {},
}

// Matches the server rule for repeated on-device classifications.
var continuedActivitySummary = regexp.MustCompile(
	`^(Explicit content|Gambling content|Violent content|Self-harm content|Activity) — detected [1-9][0-9]{0,3} times? in 5 min$`,
)

var allowedCategories = map[string]struct{}{
	"adult_content": {},
	"gambling":      {},
	"violence":      {},
	"self_harm":     {},
	"clean":         {},
	"app_usage":     {},
}

var allowedSeverities = map[string]struct{}{
	"informational": {},
	"concerning":    {},
	"severe":        {},
}

// Protection types the server accepts from a client. heartbeat_silence is
// server-only and must not be posted by the app.
var allowedProtectionTypes = map[string]struct{}{
	"monitoring_disabled":     {},
	"shielding_disabled":      {},
	"lockout_disabled":        {},
	"deep_scan_stopped":       {},
	"wrong_pin_attempt":       {},
	"pin_removed":             {},
	"pin_changed":             {},
	"family_controls_revoked": {},
	"signed_out":              {},
	"selection_cleared":       {},
}

type ProtectionNotice struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

func ValidateEvent(e OutboundEvent) error {
	if _, ok := allowedCategories[e.Category]; !ok {
		return errors.New("invalid category")
	}
	if _, ok := allowedSeverities[e.Severity]; !ok {
		return errors.New("invalid severity")
	}
	if !summaryAllowed(e.Summary) {
		return errors.New("invalid summary")
	}
	if leaked(e.Category) || leaked(e.Severity) || leaked(e.Summary) || leaked(e.Timestamp) {
		return errors.New("event carries data that stays on the device")
	}
	return nil
}

func summaryAllowed(s string) bool {
	if s == "" || len(s) > 120 {
		return false
	}
	if _, ok := allowedEventSummaries[s]; ok {
		return true
	}
	return continuedActivitySummary.MatchString(s)
}

func (e OutboundEvent) Body() ([]byte, error) {
	if err := ValidateEvent(e); err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Category  string `json:"category"`
		Severity  string `json:"severity"`
		Summary   string `json:"summary"`
		Timestamp string `json:"timestamp"`
	}{
		Category:  e.Category,
		Severity:  e.Severity,
		Summary:   e.Summary,
		Timestamp: e.Timestamp,
	})
}

func UsageEvent(now time.Time, summary string) (OutboundEvent, error) {
	e := OutboundEvent{
		Category:  "app_usage",
		Severity:  "concerning",
		Summary:   summary,
		Timestamp: now.UTC().Format(time.RFC3339),
	}
	if err := ValidateEvent(e); err != nil {
		return OutboundEvent{}, err
	}
	return e, nil
}

func ValidateProtection(n ProtectionNotice) error {
	if _, ok := allowedProtectionTypes[n.Type]; !ok {
		return errors.New("unknown protection alert type")
	}
	if n.Type == "heartbeat_silence" {
		return errors.New("unknown protection alert type")
	}
	if len(n.Detail) > 200 || leaked(n.Detail) {
		return errors.New("protection detail stays on the device")
	}
	return nil
}

// leaked rejects values that would carry a page, a picture, or an app identity.
func leaked(s string) bool {
	low := strings.ToLower(s)
	if strings.Contains(low, "http://") || strings.Contains(low, "https://") || strings.Contains(low, "www.") {
		return true
	}
	if strings.Contains(low, "screenshot") || strings.Contains(low, ".exe") || strings.Contains(low, "bundleid") || strings.Contains(low, "bundle id") {
		return true
	}
	return false
}
