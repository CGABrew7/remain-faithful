package core

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClientLoginAndMetadataEvent(t *testing.T) {
	var gotEvent string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch r.URL.Path {
		case "/auth/login":
			if strings.Contains(string(body), "secret-pass") && r.Header.Get("Authorization") != "" {
				t.Fatal("login must not send a bearer token")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"aaa.bbb.ccc","user":{"id":7,"name":"Ada","email":"ada@example.com","email_verified":true}}`))
		case "/events":
			if r.Header.Get("Authorization") != "Bearer token-1" {
				t.Fatalf("auth %s", r.Header.Get("Authorization"))
			}
			gotEvent = string(body)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1}`))
		case "/heartbeat":
			if string(body) != `{"screen":"monitoring"}` {
				t.Fatalf("heartbeat %s", body)
			}
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		case "/users/me/signals":
			if strings.Contains(string(body), "safari") || strings.Contains(string(body), "true") {
				t.Fatalf("signals %s", body)
			}
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	token := "token-1"
	c := NewClient(srv.URL, func() string { return token })
	resp, err := c.Login(context.Background(), "ada@example.com", "secret-pass")
	if err != nil || resp.User.ID != 7 || !resp.User.EmailVerified {
		t.Fatalf("%v %+v", err, resp)
	}
	event, err := UsageEvent(time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC), SummaryApp)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.CreateEvent(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(gotEvent), "chrome") || strings.Contains(gotEvent, "http") {
		t.Fatal(gotEvent)
	}
	if err := c.HeartbeatMonitoring(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := c.SendSignals(context.Background(), 2, 0); err != nil {
		t.Fatal(err)
	}
	bad, _ := UsageEvent(time.Now(), SummaryApp)
	bad.Summary = "https://secret.example/path"
	if err := c.CreateEvent(context.Background(), bad); err == nil {
		t.Fatal("client posted a page")
	}
}

func TestPINAndGroupCapMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/protection/pin/verify":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"success":false,"reason":"wrong_pin"}`))
		case strings.Contains(r.URL.Path, "/email-invite"):
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"error":"group has reached the maximum of 12 members"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c := NewClient(srv.URL, func() string { return "t" })
	pin, err := c.VerifyPIN(context.Background(), "1234")
	if err != nil || pin.Success || pin.Message != "Incorrect PIN. Your partner is the one who knows it." {
		t.Fatalf("%v %+v", err, pin)
	}
	_, err = c.InviteGroupEmail(context.Background(), 3, "a@b.co")
	if err == nil || !strings.Contains(err.Error(), "maximum of 12") {
		t.Fatal(err)
	}
}

func TestStoreKeepsTokenOffTheStateFile(t *testing.T) {
	dir := t.TempDir()
	s, err := OpenStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetSession("super-secret-token", User{ID: 4, Name: "Ada", Email: "ada@example.com"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Update(func(st *LocalState) error {
		st.OnboardingComplete = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	state, err := os.ReadFile(dir + "/state.json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(state), "super-secret-token") {
		t.Fatal("token landed in state")
	}
	if err := s.ClearSession(); err != nil {
		t.Fatal(err)
	}
	if s.SignedIn() || s.Token() != "" {
		t.Fatal("session remains")
	}
}
