package admin

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// behavioralJSON is rejected anywhere on the aggregate dashboard.
// These names would mean a person's alerts, apps, or monitoring history.
var behavioralJSON = map[string]bool{
	"user_id":           true,
	"partner_id":        true,
	"subject_user_id":   true,
	"recipient_user_id": true,
	"email":             true,
	"name":              true,
	"summary":           true,
	"category":          true,
	"severity":          true,
	"token":             true,
	"password":          true,
	"password_hash":     true,
	"visitor_hash":      true,
	"ip":                true,
	"bundle_id":         true,
	"app_name":          true,
	"members":           true,
	"screen":            true,
	"alert_type":        true,
	"event_id":          true,
}

func TestOverviewJSONKeysAreAggregate(t *testing.T) {
	walkJSONNames(t, reflect.TypeOf(Overview{}), behavioralJSON)

	c := People(2)
	drop := Percent{People: 0, Numerator: 8, Denominator: 10}
	sample := Overview{
		Range:       "7d",
		GeneratedAt: time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC),
		Website: Website{
			Visits: Show(9, 2), Visitors: People(2),
			TopPages:     []PageRow{{Path: "/privacy", Visits: Show(4, 1), Visitors: People(1)}},
			TopReferrers: []NamedCount{{Label: "example.com", Count: People(1)}},
			Campaigns:    []CampaignRow{{Source: "newsletter", Visitors: People(2)}},
			Countries:    []NamedCount{{Label: "US", Count: People(2)}},
			USStates:     []NamedCount{{Label: "CA", Count: People(1)}},
			Trend:        []TrendPoint{{Day: "2026-10-18", Visits: Show(3, 1), Visitors: People(1)}},
		},
		Growth: Growth{
			TotalAccounts: People(2), NewAccounts: People(1),
			ByDay:           []DayCount{{Day: "2026-10-18", Count: People(1)}},
			SignIn:          []MethodCount{{Method: "email", Count: People(2)}},
			DeletedAccounts: NotTracked("not kept"),
		},
		Funnel: []FunnelStep{{
			Key: "invited_partner", Label: "Invited a partner", Status: "tracked",
			Count: &c, Dropoff: &drop,
		}},
		Active: Active{
			DAU: People(1), WAU: People(2), MAU: People(2),
			DailyHistory: NotTracked("no daily heartbeat log"),
			AppVersions:  NotTracked("version is not sent"),
		},
		Accountability: Accountability{
			Partnerships:     []Partnership{{Status: "accepted", Count: People(2)}},
			GroupSizes:       []GroupSize{{Size: "2", Groups: People(1)}},
			MedianAcceptance: NotTracked("no accepted_at"),
		},
		Alerts: Alerts{
			PartnerAlertsByDay:    []DayCount{{Day: "2026-10-18", Count: Show(5, 1)}},
			MonitoringPausedByDay: []DayCount{{Day: "2026-10-18", Count: Show(4, 2)}},
			SupportTaps:           NotTracked("not saved"),
		},
		Notifications: Notifications{
			PushTokens: Show(2, 1), PushPeople: People(1),
			APNs:   NotTracked("not logged"),
			Emails: []EmailStat{{Type: "verification", Status: "not_tracked_yet", Note: "not logged"}},
		},
		Feedback: NotTracked("emailed, not stored"),
	}
	b, err := json.Marshal(sample)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "example.com/secret") {
		t.Fatal("unexpected path detail")
	}
	assertNoBehavioralValues(t, b)
}

func TestAccountStatusIsStatusOnly(t *testing.T) {
	allowed := map[string]bool{
		"email": true, "created_at": true, "email_verified": true,
		"sign_in_method": true, "last_active_at": true, "partner_count": true,
	}
	typ := reflect.TypeOf(AccountStatus{})
	if typ.NumField() != len(allowed) {
		t.Fatalf("account status has %d fields, want %d", typ.NumField(), len(allowed))
	}
	for i := 0; i < typ.NumField(); i++ {
		name := strings.Split(typ.Field(i).Tag.Get("json"), ",")[0]
		if !allowed[name] {
			t.Fatalf("unexpected account field %s", name)
		}
	}
	when := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	b, err := json.Marshal(AccountStatus{
		Email: "person@example.com", CreatedAt: when, EmailVerified: true,
		SignInMethod: "apple", LastActiveAt: &when, PartnerCount: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	for key := range raw {
		if !allowed[key] {
			t.Fatalf("response key %s is not account status", key)
		}
	}
}

func walkJSONNames(t *testing.T, typ reflect.Type, banned map[string]bool) {
	t.Helper()
	seen := map[reflect.Type]bool{}
	var walk func(reflect.Type)
	walk = func(typ reflect.Type) {
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			return
		}
		if seen[typ] {
			return
		}
		seen[typ] = true
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			tag := f.Tag.Get("json")
			name := strings.Split(tag, ",")[0]
			if name != "" && name != "-" && banned[name] {
				t.Errorf("%s.%s uses banned json name %q", typ.Name(), f.Name, name)
			}
			ft := f.Type
			for ft.Kind() == reflect.Pointer || ft.Kind() == reflect.Slice {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct && ft.PkgPath() != "time" {
				walk(ft)
			}
		}
	}
	walk(typ)
}

func assertNoBehavioralValues(t *testing.T, body []byte) {
	t.Helper()
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatal(err)
	}
	var walk func(any)
	walk = func(node any) {
		switch n := node.(type) {
		case map[string]any:
			for k, child := range n {
				if behavioralJSON[k] {
					t.Errorf("dashboard JSON contained %q", k)
				}
				walk(child)
			}
		case []any:
			for _, child := range n {
				walk(child)
			}
		case string:
			lower := strings.ToLower(n)
			for _, bad := range []string{"adult_content", "bundle", "password_hash", "explicit content"} {
				if strings.Contains(lower, bad) {
					t.Errorf("dashboard JSON contained %q", bad)
				}
			}
		}
	}
	walk(v)
}
