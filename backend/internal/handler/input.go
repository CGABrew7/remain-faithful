package handler

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxPersonName     = 80
	maxPasswordBytes  = 72
	maxDeviceToken    = 256
	maxContactName    = 80
	maxContactSubject = 120
	maxContactMessage = 4000
)

func clipPushName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Someone"
	}
	if utf8.RuneCountInString(name) <= maxPersonName {
		return name
	}
	return string([]rune(name)[:maxPersonName])
}

func containsDisallowedText(raw string, allowNewlines bool) bool {
	for _, r := range raw {
		if r == '\n' || r == '\r' || r == '\t' {
			if allowNewlines {
				continue
			}
			return true
		}
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func stripFormatRunes(raw string) string {
	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, raw)
}

// cleanPersonName trims a display name, drops invisible format characters,
// and rejects control characters and names longer than 80.
func cleanPersonName(raw string) (string, string) {
	if containsDisallowedText(raw, false) {
		return "", "Name contains characters that are not allowed."
	}
	name := strings.TrimSpace(stripFormatRunes(raw))
	if name == "" {
		return "", "Name is required."
	}
	if utf8.RuneCountInString(name) > maxPersonName {
		return "", "Name must be 80 characters or fewer."
	}
	return name, ""
}

func cleanGroupName(raw string) (string, string) {
	if containsDisallowedText(raw, false) {
		return "", "Name contains characters that are not allowed."
	}
	name := strings.TrimSpace(stripFormatRunes(raw))
	if name == "" || utf8.RuneCountInString(name) > maxGroupNameLen {
		return "", "Name must be 1 to 80 characters."
	}
	return name, ""
}

func cleanCovenant(raw string) (string, string) {
	if containsDisallowedText(raw, true) {
		return "", "Covenant contains characters that are not allowed."
	}
	text := strings.TrimSpace(stripFormatRunes(raw))
	if utf8.RuneCountInString(text) > maxGroupCovenantLen {
		return "", "Covenant is too long."
	}
	return text, ""
}

func passwordProblem(password, email string) string {
	if strings.TrimSpace(password) == "" {
		return "Choose a password that is not only spaces."
	}
	if len(password) < 8 {
		return "Password must be at least 8 characters."
	}
	if len(password) > maxPasswordBytes {
		return "That password is too long. Use 72 characters or fewer."
	}
	if email != "" && strings.EqualFold(strings.TrimSpace(password), strings.TrimSpace(email)) {
		return "Choose a password that is not your email address."
	}
	return ""
}

func allowedRelationshipType(raw string) (string, bool) {
	kind := strings.ToLower(strings.TrimSpace(raw))
	if kind == "" {
		return "partner", true
	}
	switch kind {
	case "partner", "group":
		return kind, true
	default:
		return "", false
	}
}

func cleanDeviceToken(raw string) (string, string) {
	token := strings.TrimSpace(raw)
	if token == "" || len(token) > maxDeviceToken || !utf8.ValidString(token) || containsDisallowedText(token, false) {
		return "", "Device token is not valid."
	}
	return token, ""
}

func cleanDevicePlatform(raw string) (string, string) {
	platform := strings.ToLower(strings.TrimSpace(raw))
	if platform == "" {
		platform = "ios"
	}
	switch platform {
	case "ios", "android":
		return platform, ""
	default:
		return "", "Platform is not recognized."
	}
}

func isBidiRune(r rune) bool {
	switch r {
	case '\u200e', '\u200f', '\u202a', '\u202b', '\u202c', '\u202d', '\u202e',
		'\u2066', '\u2067', '\u2068', '\u2069':
		return true
	default:
		return false
	}
}

// cleanCampaignValue keeps a short label for the owner dashboard.
// Control characters and bidi overrides are removed.
func cleanCampaignValue(raw string, maxRunes int) string {
	if maxRunes <= 0 {
		maxRunes = 80
	}
	var b strings.Builder
	n := 0
	for _, r := range raw {
		if r < 32 || r == 127 || isBidiRune(r) || unicode.Is(unicode.Cf, r) {
			continue
		}
		if n >= maxRunes {
			break
		}
		b.WriteRune(r)
		n++
	}
	return strings.TrimSpace(b.String())
}
