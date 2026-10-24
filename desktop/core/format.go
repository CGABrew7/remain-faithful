package core

import (
	"strings"
	"time"
)

func FormatWhen(raw string, now time.Time) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var t time.Time
	var err error
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	} {
		t, err = time.Parse(layout, raw)
		if err == nil {
			break
		}
	}
	if err != nil {
		return raw
	}
	mins := int(now.Sub(t).Minutes())
	if mins < 0 {
		mins = 0
	}
	switch {
	case mins < 60:
		return itoa(mins) + "m ago"
	case mins < 1440:
		return itoa(mins/60) + "h ago"
	default:
		return itoa(mins/1440) + "d ago"
	}
}

func WeekMarks(stats Stats) []string {
	if len(stats.WeekStates) > 0 {
		out := make([]string, len(stats.WeekStates))
		for i, s := range stats.WeekStates {
			switch strings.ToLower(s) {
			case "clean", "true", "1":
				out[i] = "clean"
			case "before_signup", "before-signup", "before_account", "not_applicable", "na", "none":
				out[i] = "before_account"
			default:
				out[i] = "flagged"
			}
		}
		return out
	}
	out := make([]string, len(stats.Week))
	for i, clean := range stats.Week {
		if clean {
			out[i] = "clean"
		} else {
			out[i] = "flagged"
		}
	}
	return out
}

func DeviceModel(platform string) string {
	switch platform {
	case "mac":
		return "Mac"
	case "windows":
		return "Windows PC"
	default:
		return "Desktop"
	}
}
