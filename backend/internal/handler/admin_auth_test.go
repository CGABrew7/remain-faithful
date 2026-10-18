package handler

import (
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"remain-faithful/backend/internal/admin"

	"golang.org/x/crypto/bcrypt"
)

func adminEnv(t *testing.T) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-horse"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADMIN_PASSWORD_HASH", string(hash))
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	t.Setenv("ADMIN_SESSION_SECRET", "")
	admin.ResetLoginLimits()
	t.Cleanup(admin.ResetLoginLimits)
}

func TestAdminLoginWrongPasswordAndRateLimit(t *testing.T) {
	adminEnv(t)
	h := &H{}
	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"password":"nope"}`))
		req.Header.Set("Fly-Client-IP", "198.51.100.20")
		h.AdminLogin(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("attempt %d status %d body %s", i+1, rr.Code, rr.Body.String())
		}
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"password":"correct-horse"}`))
	req.Header.Set("Fly-Client-IP", "198.51.100.20")
	h.AdminLogin(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("rate limit status %d", rr.Code)
	}
}

func TestAdminLoginSetsCookie(t *testing.T) {
	adminEnv(t)
	h := &H{}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"password":"correct-horse"}`))
	req.Header.Set("Fly-Client-IP", "198.51.100.21")
	req.Header.Set("X-Forwarded-Proto", "https")
	h.AdminLogin(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	res := rr.Result()
	defer res.Body.Close()
	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "rf_admin" {
			cookie = c
		}
	}
	if cookie == nil {
		t.Fatal("missing admin cookie")
	}
	if !cookie.HttpOnly {
		t.Fatal("cookie should be httpOnly")
	}
	if !cookie.Secure {
		t.Fatal("cookie should be secure on https")
	}
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Fatalf("SameSite = %v", cookie.SameSite)
	}
	if admin.ParseSession(cookie.Value) != nil {
		t.Fatal("cookie is not a valid admin session")
	}
}

func TestAdminExpiredSessionRejected(t *testing.T) {
	adminEnv(t)
	tok, err := admin.SignSessionTTL(-time.Second)
	if err != nil {
		t.Fatal(err)
	}
	h := &H{}
	req := httptest.NewRequest(http.MethodGet, "/admin/overview?range=7d", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	h.AdminOverview(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestCollectSkipsSecretsAndAdminPaths(t *testing.T) {
	f := &inviteFakeDB{query: func(string, []any) ([]string, [][]driver.Value) {
		return none("x")
	}}
	h := &H{DB: f.open(t)}
	hash := strings.Repeat("ab", 32)

	body := `{"event":"pageview","path":"/pricing","referrer":"https://news.ycombinator.com/item?id=1","utm_source":"hn","country":"US","region":"CA","device":"mobile","visitor_hash":"` + hash + `","ip":"203.0.113.9"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/analytics/collect", strings.NewReader(body))
	req.Header.Set("Fly-Client-IP", "203.0.113.9")
	h.CollectAnalytics(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	ins := f.ran("INSERT INTO pageviews")
	if len(ins) != 1 {
		t.Fatalf("inserts = %d", len(ins))
	}
	for _, arg := range ins[0].args {
		s, _ := arg.(string)
		if strings.Contains(s, "203.0.113.9") || strings.Contains(s, "item?id") {
			t.Fatalf("stored identifying value %q", s)
		}
	}
	if ins[0].args[1] != "/pricing" || ins[0].args[2] != "news.ycombinator.com" || ins[0].args[7] != "CA" {
		t.Fatalf("args = %#v", ins[0].args)
	}

	before := len(f.ran("INSERT INTO pageviews"))
	for _, path := range []string{"/admin", "/admin/overview"} {
		req = httptest.NewRequest(http.MethodPost, "/analytics/collect", strings.NewReader(
			`{"event":"pageview","path":"`+path+`","visitor_hash":"`+hash+`","device":"desktop"}`,
		))
		rr = httptest.NewRecorder()
		h.CollectAnalytics(rr, req)
	}
	req = httptest.NewRequest(http.MethodPost, "/analytics/collect", strings.NewReader(body))
	req.Header.Set("DNT", "1")
	rr = httptest.NewRecorder()
	h.CollectAnalytics(rr, req)
	if len(f.ran("INSERT INTO pageviews")) != before {
		t.Fatal("admin paths or DNT were stored")
	}
}

func TestAdminDeleteRequiresTypedEmail(t *testing.T) {
	adminEnv(t)
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		if strings.Contains(q, "SELECT id, name FROM users") {
			return []string{"id", "name"}, [][]driver.Value{{int64(7), "Ada"}}
		}
		return none("x")
	}}
	h := &H{DB: f.open(t)}
	tok, err := admin.SignSession()
	if err != nil {
		t.Fatal(err)
	}
	call := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/admin/account/delete", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+tok)
		rr := httptest.NewRecorder()
		h.AdminDeleteAccount(rr, req)
		return rr
	}
	if rr := call(`{"email":"ada@example.com","confirm":"nope@example.com"}`); rr.Code != http.StatusBadRequest {
		t.Fatalf("mismatch status %d", rr.Code)
	}
	if len(f.ran("DELETE FROM users")) != 0 {
		t.Fatal("mismatched confirmation deleted the account")
	}
	rr := call(`{"email":"Ada@Example.com","confirm":"ada@example.com"}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("delete status %d body %s", rr.Code, rr.Body.String())
	}
	gone := f.ran("DELETE FROM users")
	if len(gone) != 1 || gone[0].args[0] != int64(7) {
		t.Fatalf("delete args = %#v", gone)
	}
	audits := f.ran("INSERT INTO admin_audit")
	if len(audits) == 0 {
		t.Fatal("delete was not audited")
	}
}

func TestAdminLookupReturnsStatusOnly(t *testing.T) {
	adminEnv(t)
	created := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		if strings.Contains(q, "FROM users u") {
			return []string{"created_at", "email_verified", "apple", "google", "last_heartbeat_at", "count"},
				[][]driver.Value{{created, true, true, false, nil, int64(4)}}
		}
		return none("x")
	}}
	h := &H{DB: f.open(t)}
	tok, err := admin.SignSession()
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/admin/account?email=Ada@Example.com", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	h.AdminLookupAccount(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, bad := range []string{"summary", "category", "alert", "partner_id", "bundle", "Ada\""} {
		if strings.Contains(strings.ToLower(body), bad) {
			t.Fatalf("lookup body leaked %s: %s", bad, body)
		}
	}
	for _, need := range []string{"ada@example.com", "apple", "email_verified", "partner_count"} {
		if !strings.Contains(body, need) {
			t.Fatalf("lookup missing %s: %s", need, body)
		}
	}
}

func TestDisabledAdminIsNotFound(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD_HASH", "")
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	h := &H{}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"password":"x"}`))
	h.AdminLogin(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("login status %d", rr.Code)
	}
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/admin/overview", nil)
	h.AdminOverview(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("overview status %d", rr.Code)
	}
}
