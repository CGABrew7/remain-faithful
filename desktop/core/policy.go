package core

import "time"

const (
	maxSampleSeconds = 30
)

// Ledger counts foreground time on this computer for the current local day.
// It is not uploaded. AppSent and BrowserSent remember the one daily alert
// for each set, matching the phone: chosen apps are one set, browsers are the other.
type Ledger struct {
	Day                string `json:"day"`
	AppSeconds         int    `json:"app_seconds"`
	BrowserSeconds     int    `json:"browser_seconds"`
	AppSent            bool   `json:"app_sent"`
	BrowserSent        bool   `json:"browser_sent"`
	AppSentMinutes     []int  `json:"app_sent_minutes,omitempty"`
	BrowserSentMinutes []int  `json:"browser_sent_minutes,omitempty"`
}

// usageMinutes are the partner alerts for one local day. Each fires once.
var usageMinutes = []int{1, 5, 15, 30}

type WatchConfig struct {
	Consent  bool
	Alerts   bool
	Blocking bool
	Selected []string
}

type Outcome struct {
	Block    bool
	BlockKey string
	Event    *OutboundEvent
}

// Evaluate decides what a foreground sample means.
// key is a local bundle id or process name. It is never copied into Outcome.Event.
func Evaluate(now time.Time, platform string, cfg WatchConfig, ledger Ledger, key string, delta time.Duration) (Ledger, Outcome) {
	if !cfg.Consent || (!cfg.Alerts && !cfg.Blocking) {
		return ledger, Outcome{}
	}
	day := now.Format("2006-01-02")
	if ledger.Day != day {
		ledger = Ledger{Day: day}
	}
	key = normalizeKey(platform, key)
	if key == "" || IsSelf(platform, key) {
		return ledger, Outcome{}
	}
	choice, ok := Lookup(platform, key)
	if !ok || !selected(cfg.Selected, choice.Key) {
		return ledger, Outcome{}
	}

	var out Outcome
	if cfg.Blocking {
		out.Block = true
		out.BlockKey = choice.Key
	}
	if !cfg.Alerts {
		return ledger, out
	}

	secs := int(delta / time.Second)
	if secs < 0 {
		secs = 0
	}
	if secs > maxSampleSeconds {
		secs = maxSampleSeconds
	}
	if choice.Kind == "browser" {
		ledger.BrowserSeconds += secs
		if minutes, ok := nextUsageMinutes(ledger.BrowserSeconds, ledger.BrowserSentMinutes, ledger.BrowserSent); ok {
			event, err := UsageEvent(now, usageSummary("category", minutes))
			if err == nil {
				ledger.BrowserSentMinutes = append(ledger.BrowserSentMinutes, minutes)
				if minutes == 1 {
					ledger.BrowserSent = true
				}
				out.Event = &event
			}
		}
	} else {
		ledger.AppSeconds += secs
		if minutes, ok := nextUsageMinutes(ledger.AppSeconds, ledger.AppSentMinutes, ledger.AppSent); ok {
			event, err := UsageEvent(now, usageSummary("app", minutes))
			if err == nil {
				ledger.AppSentMinutes = append(ledger.AppSentMinutes, minutes)
				if minutes == 1 {
					ledger.AppSent = true
				}
				out.Event = &event
			}
		}
	}
	return ledger, out
}

func nextUsageMinutes(seconds int, sent []int, legacyOne bool) (int, bool) {
	for _, minutes := range usageMinutes {
		if seconds < minutes*60 {
			return 0, false
		}
		if containsMinutes(sent, minutes) || (minutes == 1 && legacyOne) {
			continue
		}
		return minutes, true
	}
	return 0, false
}

func containsMinutes(sent []int, minutes int) bool {
	for _, n := range sent {
		if n == minutes {
			return true
		}
	}
	return false
}

func usageSummary(kind string, minutes int) string {
	unit := "minutes"
	if minutes == 1 {
		unit = "minute"
	}
	if kind == "category" {
		return "Monitored app category used for " + itoa(minutes) + "+ " + unit
	}
	return "Monitored app used for " + itoa(minutes) + "+ " + unit
}

// RewindUsage undoes one usage alert that failed to upload, so the next sample can send it again.
func RewindUsage(ledger Ledger, summary string) Ledger {
	minutes, kind := usageMinutesForSummary(summary)
	if minutes == 0 {
		return ledger
	}
	floor := minutes*60 - 1
	if kind == "category" {
		ledger.BrowserSentMinutes = removeMinutes(ledger.BrowserSentMinutes, minutes)
		if minutes == 1 {
			ledger.BrowserSent = false
		}
		if ledger.BrowserSeconds >= minutes*60 {
			ledger.BrowserSeconds = floor
		}
		return ledger
	}
	ledger.AppSentMinutes = removeMinutes(ledger.AppSentMinutes, minutes)
	if minutes == 1 {
		ledger.AppSent = false
	}
	if ledger.AppSeconds >= minutes*60 {
		ledger.AppSeconds = floor
	}
	return ledger
}

func usageMinutesForSummary(summary string) (int, string) {
	for _, minutes := range usageMinutes {
		if summary == usageSummary("app", minutes) {
			return minutes, "app"
		}
		if summary == usageSummary("category", minutes) {
			return minutes, "category"
		}
	}
	return 0, ""
}

func removeMinutes(sent []int, minutes int) []int {
	out := sent[:0]
	for _, n := range sent {
		if n != minutes {
			out = append(out, n)
		}
	}
	return out
}

func selected(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

// HeartbeatScreen is "monitoring" only while this computer is actually watching.
// An empty result means do not post a heartbeat. Desktop heartbeats name this
// computer and do not refresh the account timestamp shared with a phone.
func HeartbeatScreen(consent, alerts, blocking, lockout bool) string {
	if consent && (alerts || blocking || lockout) {
		return "monitoring"
	}
	return ""
}
