package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Cloudflare's published dummy secret. It is not a real key.
const turnstileDummySecret = "1x0000000000000000000000000000000AA"

func withLoadedAt(body string) string {
	if strings.Contains(body, `"loadedAt"`) {
		return body
	}
	trimmed := strings.TrimSpace(body)
	if !strings.HasPrefix(trimmed, "{") {
		return body
	}
	return fmt.Sprintf(`{"loadedAt":%d,%s`, time.Now().Add(-10*time.Second).UnixMilli(), trimmed[1:])
}

func TestFormSubmittedInTime(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		at   int64
		ok   bool
	}{
		{name: "missing", at: 0, ok: false},
		{name: "seconds", at: now.Unix(), ok: false},
		{name: "just now", at: now.UnixMilli(), ok: false},
		{name: "one second", at: now.Add(-time.Second).UnixMilli(), ok: false},
		{name: "two seconds", at: now.Add(-2 * time.Second).UnixMilli(), ok: true},
		{name: "ten seconds", at: now.Add(-10 * time.Second).UnixMilli(), ok: true},
		{name: "too old", at: now.Add(-25 * time.Hour).UnixMilli(), ok: false},
		{name: "future", at: now.Add(time.Minute).UnixMilli(), ok: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formSubmittedInTime(tc.at, now); got != tc.ok {
				t.Fatalf("formSubmittedInTime(%d) = %v, want %v", tc.at, got, tc.ok)
			}
		})
	}
}

func TestContactDropsFastSubmit(t *testing.T) {
	t.Setenv("TURNSTILE_SECRET_KEY", "")
	t.Setenv("CONTACT_TO_EMAIL", "support@remainfaithful.com")
	em := &capturingEmail{}
	h := &H{Email: em}
	body := fmt.Sprintf(`{"name":"Ada","email":"fast@example.com","message":"Hello there","loadedAt":%d}`, time.Now().UnixMilli())
	rr := httptest.NewRecorder()
	h.Contact(rr, httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(body)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	if len(em.calls) != 0 {
		t.Fatal("fast submit sent mail")
	}
}

func TestContactEmailRateLimit(t *testing.T) {
	resetFormEmailWindows()
	t.Cleanup(resetFormEmailWindows)
	t.Setenv("TURNSTILE_SECRET_KEY", "")
	t.Setenv("CONTACT_TO_EMAIL", "support@remainfaithful.com")
	em := &capturingEmail{}
	h := &H{Email: em}
	body := withLoadedAt(`{"name":"Ada","email":"Form-Limit@example.com","message":"Hello there"}`)
	for i := 0; i < formEmailLimit; i++ {
		rr := httptest.NewRecorder()
		h.Contact(rr, httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(body)))
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d status %d %s", i+1, rr.Code, rr.Body.String())
		}
	}
	rr := httptest.NewRecorder()
	h.Contact(rr, httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(body)))
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	if len(em.calls) != formEmailLimit {
		t.Fatalf("emails = %d", len(em.calls))
	}
}

func TestWaitlistDropsFastAndHoneypot(t *testing.T) {
	t.Setenv("TURNSTILE_SECRET_KEY", "")
	db := newMemDB()
	h := &H{DB: db.open(t)}
	fast := fmt.Sprintf(`{"email":"fast-list@example.com","loadedAt":%d}`, time.Now().UnixMilli())
	rr := httptest.NewRecorder()
	h.JoinWaitlist(rr, httptest.NewRequest(http.MethodPost, "/waitlist", strings.NewReader(fast)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("fast status %d", rr.Code)
	}
	rr = httptest.NewRecorder()
	h.JoinWaitlist(rr, httptest.NewRequest(http.MethodPost, "/waitlist", strings.NewReader(
		`{"email":"bot-list@example.com","company":"filled"}`,
	)))
	if rr.Code != http.StatusOK || len(db.waitlistRows()) != 0 {
		t.Fatalf("honeypot status %d rows %d", rr.Code, len(db.waitlistRows()))
	}
}

func TestTurnstileOnContactAndWaitlist(t *testing.T) {
	t.Setenv("TURNSTILE_SECRET_KEY", turnstileDummySecret)
	t.Setenv("CONTACT_TO_EMAIL", "support@remainfaithful.com")
	var sawSecret, sawToken, sawIP bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		encoded := string(body)
		if strings.Contains(encoded, "secret="+turnstileDummySecret) {
			sawSecret = true
		}
		if strings.Contains(encoded, "response=token-ok") {
			sawToken = true
		}
		if strings.Contains(encoded, "remoteip=203.0.113.10") {
			sawIP = true
		}
		action := "contact"
		if strings.Contains(encoded, "response=token-wait") {
			action = "waitlist"
		}
		if strings.Contains(encoded, "response=token-bad") {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"success":false,"error-codes":["invalid-input-response"]}`)
			return
		}
		if strings.Contains(encoded, "response=token-mismatch") {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"success":true,"action":"other"}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success":true,"action":%q}`, action)
	}))
	defer srv.Close()
	previous := turnstileVerifyURL
	turnstileVerifyURL = srv.URL
	t.Cleanup(func() { turnstileVerifyURL = previous })

	em := &capturingEmail{}
	h := &H{Email: em}
	post := func(url, body string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(withLoadedAt(body)))
		req.RemoteAddr = "203.0.113.10:443"
		rr := httptest.NewRecorder()
		if url == "/contact" {
			h.Contact(rr, req)
		} else {
			h.JoinWaitlist(rr, req)
		}
		return rr
	}

	rr := post("/contact", `{"name":"Ada","email":"turnstile-ok@example.com","message":"Hello there","turnstileToken":"token-ok","formAction":"contact"}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("accept status %d %s", rr.Code, rr.Body.String())
	}
	if len(em.calls) != 1 || !sawSecret || !sawToken || !sawIP {
		t.Fatalf("mail %d secret %v token %v ip %v", len(em.calls), sawSecret, sawToken, sawIP)
	}

	rr = post("/contact", `{"name":"Ada","email":"turnstile-bad@example.com","message":"Hello there","turnstileToken":"token-bad","formAction":"contact"}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("reject status %d %s", rr.Code, rr.Body.String())
	}

	rr = post("/contact", `{"name":"Ada","email":"turnstile-mismatch@example.com","message":"Hello there","turnstileToken":"token-mismatch","formAction":"contact"}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("mismatch status %d %s", rr.Code, rr.Body.String())
	}

	db := newMemDB()
	h.DB = db.open(t)
	rr = post("/waitlist", `{"email":"turnstile-list@example.com","turnstileToken":"token-wait","formAction":"waitlist"}`)
	if rr.Code != http.StatusOK || len(db.waitlistRows()) != 1 {
		t.Fatalf("waitlist status %d rows %d body %s", rr.Code, len(db.waitlistRows()), rr.Body.String())
	}

	rr = post("/waitlist", `{"email":"turnstile-missing@example.com","formAction":"waitlist"}`)
	if rr.Code != http.StatusBadRequest || len(db.waitlistRows()) != 1 {
		t.Fatalf("missing token status %d rows %d", rr.Code, len(db.waitlistRows()))
	}
}
