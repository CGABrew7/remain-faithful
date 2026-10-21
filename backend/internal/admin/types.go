package admin

import (
	"encoding/json"
	"math"
	"time"
)

// Count is a people-safe number. Fewer than 3 people in the bucket
// serializes as the string "<3" so a small cell cannot point at someone.
// N is the figure to show when People is at least 3. People is the
// headcount that decides whether N may be shown.
type Count struct {
	N      int
	People int
}

// People returns a count where the number and the headcount are the same.
func People(n int) Count {
	if n < 0 {
		n = 0
	}
	return Count{N: n, People: n}
}

// Show returns n when people is at least 3, otherwise "<3".
// Use this when the displayed number is not itself a headcount
// (page views, for example) but a small audience must still hide it.
func Show(n, people int) Count {
	if n < 0 {
		n = 0
	}
	if people < 0 {
		people = 0
	}
	return Count{N: n, People: people}
}

// Zero is a real zero, safe to show because it is not a small headcount.
func Zero() Count {
	return Count{N: 0, People: 3}
}

func (c Count) MarshalJSON() ([]byte, error) {
	if c.People < 3 {
		return []byte(`"<3"`), nil
	}
	return json.Marshal(c.N)
}

// Percent is a rate. It becomes "<3" when the audience is under 3 people,
// and null when there is nothing to divide.
type Percent struct {
	People      int
	Numerator   int
	Denominator int
}

// RateAudience hides a percentage when either side of the split is a
// group of 1 or 2 people. A zero side is safe to show.
func RateAudience(part, whole int) int {
	if whole < 3 {
		return whole
	}
	other := whole - part
	if part > 0 && part < 3 {
		return 0
	}
	if other > 0 && other < 3 {
		return 0
	}
	return whole
}

func (p Percent) MarshalJSON() ([]byte, error) {
	if p.People < 3 {
		return []byte(`"<3"`), nil
	}
	if p.Denominator <= 0 {
		return []byte("null"), nil
	}
	pct := math.Round(float64(p.Numerator)/float64(p.Denominator)*1000) / 10
	return json.Marshal(pct)
}

// Average is a mean that stays hidden for a group of fewer than 3 people
// or fewer than 3 groups.
type Average struct {
	People int
	Groups int
	Sum    float64
}

func (a Average) MarshalJSON() ([]byte, error) {
	if a.People < 3 || a.Groups < 3 {
		return []byte(`"<3"`), nil
	}
	if a.Groups == 0 {
		return []byte("null"), nil
	}
	avg := math.Round((a.Sum/float64(a.Groups))*10) / 10
	return json.Marshal(avg)
}

// Gap marks a metric the server cannot compute yet.
type Gap struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}

func NotTracked(note string) Gap {
	return Gap{Status: "not_tracked_yet", Note: note}
}

type PageRow struct {
	Path     string `json:"path"`
	Visits   Count  `json:"visits"`
	Visitors Count  `json:"visitors"`
}

type NamedCount struct {
	Label string `json:"label"`
	Count Count  `json:"count"`
}

type CampaignRow struct {
	Source   string `json:"source"`
	Medium   string `json:"medium"`
	Campaign string `json:"campaign"`
	Visitors Count  `json:"visitors"`
}

type TrendPoint struct {
	Day      string `json:"day"`
	Visits   Count  `json:"visits"`
	Visitors Count  `json:"visitors"`
}

type DayCount struct {
	Day   string `json:"day"`
	Count Count  `json:"count"`
}

type Website struct {
	Visits          Count         `json:"visits"`
	Visitors        Count         `json:"visitors"`
	TopPages        []PageRow     `json:"top_pages"`
	TopReferrers    []NamedCount  `json:"top_referrers"`
	Campaigns       []CampaignRow `json:"campaigns"`
	Countries       []NamedCount  `json:"countries"`
	USStates        []NamedCount  `json:"us_states"`
	MobileVisitors  Count         `json:"mobile_visitors"`
	DesktopVisitors Count         `json:"desktop_visitors"`
	AppStoreClicks  Count         `json:"app_store_clicks"`
	AppStoreCTR     Percent       `json:"app_store_ctr_percent"`
	AndroidSignups  Count         `json:"android_signups"`
	Trend           []TrendPoint  `json:"trend"`
}

type MethodCount struct {
	Method string `json:"method"`
	Count  Count  `json:"count"`
}

type Growth struct {
	TotalAccounts   Count         `json:"total_accounts"`
	NewAccounts     Count         `json:"new_accounts"`
	ByDay           []DayCount    `json:"by_day"`
	ByWeek          []DayCount    `json:"by_week"`
	SignIn          []MethodCount `json:"sign_in"`
	EmailVerified   Percent       `json:"email_verified_percent"`
	DeletedAccounts Gap           `json:"deleted_accounts"`
}

type FunnelStep struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Status  string   `json:"status"`
	Count   *Count   `json:"count,omitempty"`
	Dropoff *Percent `json:"dropoff_percent,omitempty"`
	Note    string   `json:"note,omitempty"`
}

type Active struct {
	Basis           string `json:"basis"`
	DAU             Count  `json:"dau"`
	WAU             Count  `json:"wau"`
	MAU             Count  `json:"mau"`
	MonitoringOn    Count  `json:"monitoring_on"`
	MonitoringOff   Count  `json:"monitoring_off"`
	HeartbeatSilent Count  `json:"heartbeat_silent"`
	NeverCheckedIn  Count         `json:"never_checked_in"`
	MonitoredApps   Count         `json:"monitored_apps"`
	DailyHistory    Gap           `json:"daily_history"`
	AppVersions     VersionRollup `json:"app_versions"`
}

// VersionRollup is a count of phones per app version. No account is listed.
type VersionRollup struct {
	Status   string       `json:"status"`
	Versions []NamedCount `json:"versions"`
}

// CountRollup is one aggregate total, such as support taps.
type CountRollup struct {
	Status string `json:"status"`
	Count  Count  `json:"count"`
}

type Partnership struct {
	Status string `json:"status"`
	Count  Count  `json:"count"`
}

type GroupSize struct {
	Size   string `json:"size"`
	Groups Count  `json:"groups"`
}

type Accountability struct {
	Partnerships     []Partnership `json:"partnerships"`
	InvitesSent      Count         `json:"invites_sent"`
	Acceptance       Percent       `json:"acceptance_percent"`
	MedianAcceptance Gap           `json:"median_acceptance"`
	Groups           Count         `json:"groups"`
	AverageSize      Average       `json:"average_group_size"`
	GroupSizes       []GroupSize   `json:"group_sizes"`
	ZeroPartners     Count         `json:"users_with_zero_partners"`
}

type Alerts struct {
	PartnerAlertsByDay    []DayCount `json:"partner_alerts_by_day"`
	MonitoringPausedByDay []DayCount `json:"monitoring_paused_by_day"`
	SupportTaps           CountRollup `json:"support_taps"`
}

type EmailStat struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Count  *Count `json:"count,omitempty"`
	Note   string `json:"note,omitempty"`
}

type Notifications struct {
	PushTokens Count       `json:"push_tokens"`
	PushPeople Count       `json:"push_people"`
	APNs       Gap         `json:"apns"`
	Emails     []EmailStat `json:"emails"`
}

// Overview is the owner dashboard. It is aggregates and gaps only.
type Overview struct {
	Range          string         `json:"range"`
	GeneratedAt    time.Time      `json:"generated_at"`
	Website        Website        `json:"website"`
	Growth         Growth         `json:"growth"`
	Funnel         []FunnelStep   `json:"funnel"`
	Active         Active         `json:"active"`
	Accountability Accountability `json:"accountability"`
	Alerts         Alerts         `json:"alerts"`
	Notifications  Notifications  `json:"notifications"`
	Feedback       Gap            `json:"feedback"`
}

// AccountStatus is the only per-account view. No alerts, apps, or partners by name.
type AccountStatus struct {
	Email         string     `json:"email"`
	CreatedAt     time.Time  `json:"created_at"`
	EmailVerified bool       `json:"email_verified"`
	SignInMethod  string     `json:"sign_in_method"`
	LastActiveAt  *time.Time `json:"last_active_at"`
	PartnerCount  int        `json:"partner_count"`
}

type TableCount struct {
	Table string `json:"table"`
	Rows  Count  `json:"rows"`
}

type Health struct {
	API          string       `json:"api"`
	Database     string       `json:"database"`
	DBBytes      int64        `json:"db_bytes"`
	Tables       []TableCount `json:"tables"`
	ServerErrors Gap          `json:"server_errors"`
}
