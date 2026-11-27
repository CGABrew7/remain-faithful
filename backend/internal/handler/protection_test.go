package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	rfauth "remain-faithful/backend/internal/auth"
)

func TestIsClientProtectionAlertType(t *testing.T) {
	iosTypes := []string{
		"monitoring_disabled",
		"shielding_disabled",
		"lockout_disabled",
		"deep_scan_stopped",
		"wrong_pin_attempt",
		"pin_removed",
		"pin_changed",
		"family_controls_revoked",
		"signed_out",
		"selection_cleared",
		"coach_nudge",
	}
	for _, typ := range iosTypes {
		if !isClientProtectionAlertType(typ) {
			t.Errorf("expected client type %q", typ)
		}
		if isServerOnlyProtectionAlertType(typ) {
			t.Errorf("client type %q must not be server-only", typ)
		}
	}

	for _, typ := range []string{"", "not_a_real_type", "monitoring_off", "spam", "heartbeat_silence", "lockout_broadcast_stopped", "pin_wrong_attempt"} {
		if isClientProtectionAlertType(typ) {
			t.Errorf("expected type %q to be rejected on client POST", typ)
		}
	}
}

func TestIsServerOnlyProtectionAlertType(t *testing.T) {
	if !isServerOnlyProtectionAlertType("heartbeat_silence") {
		t.Fatal("expected heartbeat_silence to be server-only")
	}
	if isServerOnlyProtectionAlertType("monitoring_disabled") {
		t.Fatal("client type must not be server-only")
	}
}

func TestProtectionAlertBodyKnownTypesIgnoreDetail(t *testing.T) {
	const name = "Alex"
	injected := "UNBOUNDED CLIENT DETAIL " + strings.Repeat("x", 200)

	cases := []struct {
		alertType string
		want      string
	}{
		{"deep_scan_stopped", "Alex stopped Deep Scan"},
		{"lockout_broadcast_stopped", "Alex stopped Deep Scan"},
		{"selection_cleared", "Alex removed every chosen app"},
		{"wrong_pin_attempt", "Alex entered an incorrect Partner PIN"},
		{"pin_wrong_attempt", "Alex entered an incorrect Partner PIN"},
		{"monitoring_disabled", "Alex turned off activity monitoring"},
		{"shielding_disabled", "Alex turned off app blocking"},
		{"lockout_disabled", "Alex turned off Session Shield"},
		{"pin_removed", "Alex's protection PIN was removed"},
		{"pin_changed", "Alex's protection PIN was changed"},
		{"family_controls_revoked", "Alex revoked Family Controls authorization"},
		{"heartbeat_silence", "Alex's device has stopped sending heartbeats"},
		{"coach_nudge", "Your partner would welcome a word from you today."},
	}

	for _, tc := range cases {
		got := protectionAlertBody(name, tc.alertType, injected)
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.alertType, got, tc.want)
		}
		if strings.Contains(got, "UNBOUNDED") || strings.Contains(got, "xxxx") {
			t.Errorf("%s: body leaked client detail: %q", tc.alertType, got)
		}
	}
}

func TestListedAlertCategoryUsesProtectionType(t *testing.T) {
	if got := listedAlertCategory("", "selection_cleared"); got != "selection_cleared" {
		t.Fatalf("selection_cleared category = %q", got)
	}
	if got := listedAlertCategory("", "deep_scan_stopped"); got != "deep_scan_stopped" {
		t.Fatalf("deep_scan_stopped category = %q", got)
	}
	if got := listedAlertCategory("adult_content", "selection_cleared"); got != "adult_content" {
		t.Fatalf("content category should win, got %q", got)
	}
	if got := listedAlertCategory("", ""); got != "protection" {
		t.Fatalf("empty category = %q", got)
	}
}

func TestProtectionAlertBodyDefaultIgnoresDetail(t *testing.T) {
	got := protectionAlertBody("Alex", "not_a_real_type", "should never appear in APNs")
	if got != "Alex changed a protection setting" {
		t.Fatalf("default body = %q", got)
	}
	if strings.Contains(got, "should never appear") {
		t.Fatalf("default body leaked detail: %q", got)
	}
}

func TestSendProtectionAlertRejectsUnknownType(t *testing.T) {
	h := &H{}
	req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(`{"type":"partner_spam"}`))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.SendProtectionAlert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "unknown protection alert type") {
		t.Fatalf("body = %q, want unknown-type error", rr.Body.String())
	}
}

func TestSendProtectionAlertRejectsServerOnlyTypes(t *testing.T) {
	h := &H{}
	for _, typ := range []string{"heartbeat_silence"} {
		req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(`{"type":"`+typ+`","detail":"forged silence"}`))
		req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
		rr := httptest.NewRecorder()

		h.SendProtectionAlert(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want 400", typ, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "unknown protection alert type") {
			t.Fatalf("%s: body = %q, want unknown-type error", typ, rr.Body.String())
		}
	}
}

func TestSendProtectionAlertRequiresType(t *testing.T) {
	h := &H{}
	req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(`{"type":""}`))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.SendProtectionAlert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

func TestSendProtectionAlertAcceptsClientType(t *testing.T) {
	db := sql.OpenDB(nameLookupConnector{name: "Alex"})
	t.Cleanup(func() { _ = db.Close() })

	h := &H{DB: db}
	req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(`{"type":"monitoring_disabled","detail":"should be ignored"}`))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.SendProtectionAlert(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %q", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"status":"ok"`) {
		t.Fatalf("body = %q, want status ok", rr.Body.String())
	}
}

func TestCoachNudgeRejectsExtraFieldsAndFreeDetail(t *testing.T) {
	h := &H{}
	cases := []string{
		`{"type":"coach_nudge","detail":"opened Instagram"}`,
		`{"type":"coach_nudge","detail":"Your partner would welcome a word from you today.","app":"Instagram"}`,
		`{"type":"coach_nudge","detail":"Your partner would welcome a word from you today.","pattern":"sunday"}`,
		`{"type":"coach_nudge","detail":"Your partner would welcome a word from you today.","device_id":"phone-1","device_kind":"android"}`,
		`{"type":"signed_out","detail":"Signed out.","score":0.9}`,
	}
	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(body))
		req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
		rr := httptest.NewRecorder()
		h.SendProtectionAlert(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d, want 400 (%s)", body, rr.Code, rr.Body.String())
		}
	}
}

func TestCoachNudgeOnePerLocalDay(t *testing.T) {
	resetCoachCaps()
	t.Cleanup(resetCoachCaps)

	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Alex", email: "alex@example.com"})
	h := &H{DB: db.open(t)}
	const body = `{"type":"coach_nudge","detail":"Your partner would welcome a word from you today."}`

	first := httptest.NewRequest(http.MethodPost, "/protection/alerts?offset_minutes=-300", strings.NewReader(body))
	first = first.WithContext(rfauth.ContextWithUser(first.Context(), 1))
	rr := httptest.NewRecorder()
	h.SendProtectionAlert(rr, first)
	if rr.Code != http.StatusOK {
		t.Fatalf("first status = %d body %s", rr.Code, rr.Body.String())
	}

	second := httptest.NewRequest(http.MethodPost, "/protection/alerts?offset_minutes=-300", strings.NewReader(body))
	second = second.WithContext(rfauth.ContextWithUser(second.Context(), 1))
	again := httptest.NewRecorder()
	h.SendProtectionAlert(again, second)
	if again.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d body %s", again.Code, again.Body.String())
	}
	if strings.Contains(again.Body.String(), "Instagram") || strings.Contains(again.Body.String(), "Alex") {
		t.Fatalf("rate limit body leaked detail: %s", again.Body.String())
	}
}

func TestCoachNudgeStoredRowBlocksAnotherDayCount(t *testing.T) {
	db := newMemDB()
	db.alerts = append(db.alerts, memAlert{
		subject:   4,
		alertType: "coach_nudge",
		summary:   coachNudgeDetail,
		created:   time.Now(),
	})
	h := &H{DB: db.open(t)}
	used, err := h.coachNudgeStored(t.Context(), 4, time.Now().Add(-time.Hour))
	if err != nil || !used {
		t.Fatalf("stored=%v err=%v", used, err)
	}
	other, err := h.coachNudgeStored(t.Context(), 9, time.Now().Add(-time.Hour))
	if err != nil || other {
		t.Fatalf("other user stored=%v err=%v", other, err)
	}
}

func TestProtectionBodyKeepsCoachSentenceFixed(t *testing.T) {
	got := protectionBody("Alex", "coach_nudge", "Android phone")
	if got != "Your partner would welcome a word from you today." {
		t.Fatalf("body = %q", got)
	}
}

func TestTruncateForLog(t *testing.T) {
	if got := truncateForLog("short", 80); got != "short" {
		t.Fatalf("short = %q", got)
	}
	long := strings.Repeat("a", 100)
	got := truncateForLog(long, 80)
	if got != strings.Repeat("a", 80)+"…" {
		t.Fatalf("truncated = %q", got)
	}
}

// nameLookupConnector is a test-only sql.Connector that answers
// SELECT name FROM users and returns no partners.
type nameLookupConnector struct {
	name string
}

func (c nameLookupConnector) Connect(context.Context) (driver.Conn, error) {
	return &nameLookupConn{name: c.name}, nil
}

func (c nameLookupConnector) Driver() driver.Driver { return nameLookupDriver{} }

type nameLookupDriver struct{}

func (nameLookupDriver) Open(string) (driver.Conn, error) {
	return &nameLookupConn{}, nil
}

type nameLookupConn struct {
	name string
}

func (c *nameLookupConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c *nameLookupConn) Close() error                        { return nil }
func (c *nameLookupConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }

func (c *nameLookupConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "email_verified") {
		return &nameLookupRows{cols: []string{"email", "email_verified", "apple_id", "google_id"}, vals: [][]driver.Value{{"ada@example.com", true, "", ""}}}, nil
	}
	if strings.Contains(query, "FROM users") {
		return &nameLookupRows{cols: []string{"name"}, vals: [][]driver.Value{{c.name}}}, nil
	}
	return &nameLookupRows{cols: []string{"partner_id"}}, nil
}

type nameLookupRows struct {
	cols []string
	vals [][]driver.Value
	i    int
}

func (r *nameLookupRows) Columns() []string { return r.cols }
func (r *nameLookupRows) Close() error      { return nil }
func (r *nameLookupRows) Next(dest []driver.Value) error {
	if r.i >= len(r.vals) {
		return io.EOF
	}
	copy(dest, r.vals[r.i])
	r.i++
	return nil
}
