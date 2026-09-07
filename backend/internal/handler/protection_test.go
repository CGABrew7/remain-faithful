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
		{"deep_scan_stopped", "Alex stopped Deep Scan — apps are re-shielded"},
		{"lockout_broadcast_stopped", "Alex stopped Deep Scan — apps are re-shielded"},
		{"wrong_pin_attempt", "Alex entered an incorrect Partner PIN"},
		{"pin_wrong_attempt", "Alex entered an incorrect Partner PIN"},
		{"monitoring_disabled", "Alex turned off activity monitoring"},
		{"shielding_disabled", "Alex turned off app blocking"},
		{"lockout_disabled", "Alex disabled App Lockout"},
		{"pin_removed", "Alex's protection PIN was removed"},
		{"pin_changed", "Alex's protection PIN was changed"},
		{"family_controls_revoked", "Alex revoked Family Controls authorization"},
		{"heartbeat_silence", "Alex's device has stopped sending heartbeats"},
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
