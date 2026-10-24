package core

import "time"

const (
	alertAfterSeconds = 60
	maxSampleSeconds  = 30
)

// Ledger counts foreground time on this computer for the current local day.
// It is not uploaded. AppSent and BrowserSent remember the one daily alert
// for each set, matching the phone: chosen apps are one set, browsers are the other.
type Ledger struct {
	Day            string `json:"day"`
	AppSeconds     int    `json:"app_seconds"`
	BrowserSeconds int    `json:"browser_seconds"`
	AppSent        bool   `json:"app_sent"`
	BrowserSent    bool   `json:"browser_sent"`
}

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
		if ledger.BrowserSeconds >= alertAfterSeconds && !ledger.BrowserSent {
			event, err := UsageEvent(now, SummaryCategory)
			if err == nil {
				ledger.BrowserSent = true
				out.Event = &event
			}
		}
	} else {
		ledger.AppSeconds += secs
		if ledger.AppSeconds >= alertAfterSeconds && !ledger.AppSent {
			event, err := UsageEvent(now, SummaryApp)
			if err == nil {
				ledger.AppSent = true
				out.Event = &event
			}
		}
	}
	return ledger, out
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
