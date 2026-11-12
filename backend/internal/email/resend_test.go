package email

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendUsesResendWhenKeySet(t *testing.T) {
	var got map[string]any
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"id":"x"}`))
	}))
	defer srv.Close()
	c := &Client{resendKey: "re_test", sendgridKey: "sg_test", from: "noreply@remainfaithful.com",
		resendURL: srv.URL, sendgridURL: "http://127.0.0.1:1", httpClient: srv.Client()}
	if err := c.SendPasswordReset("a@example.com", "Ann \"A\" Lee", "https://x/reset"); err != nil {
		t.Fatal(err)
	}
	if auth != "Bearer re_test" {
		t.Fatalf("auth=%q", auth)
	}
	if got["from"] != `"Remain Faithful" <noreply@remainfaithful.com>` {
		t.Fatalf("from=%v", got["from"])
	}
	to := got["to"].([]any)
	if len(to) != 1 || to[0] != `"Ann \"A\" Lee" <a@example.com>` {
		t.Fatalf("to=%v", to)
	}
	if got["subject"] != "Reset your Remain Faithful password" || got["text"] == "" || got["html"] == "" {
		t.Fatalf("payload=%v", got)
	}
	if got["reply_to"] != "support@remainfaithful.com" {
		t.Fatalf("reply_to=%v", got["reply_to"])
	}
}

func TestContactReplyGoesToSubmitter(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		w.WriteHeader(200)
	}))
	defer srv.Close()
	c := &Client{resendKey: "re_test", from: "noreply@remainfaithful.com", resendURL: srv.URL, httpClient: srv.Client()}
	if err := c.SendContact("person@example.com", "Pat", "Hello", "A question", "support@remainfaithful.com"); err != nil {
		t.Fatal(err)
	}
	if got["reply_to"] != "person@example.com" {
		t.Fatalf("reply_to=%v", got["reply_to"])
	}
}

func TestSendResendErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		_, _ = w.Write([]byte(`{"message":"domain not verified"}`))
	}))
	defer srv.Close()
	c := &Client{resendKey: "re_test", from: "noreply@remainfaithful.com", resendURL: srv.URL, httpClient: srv.Client()}
	if err := c.SendContact("u@example.com", "U", "Hi", "msg", "support@remainfaithful.com"); err == nil {
		t.Fatal("expected error on 422")
	}
}

func TestProviderSelection(t *testing.T) {
	if (&Client{resendKey: "a", sendgridKey: "b"}).Provider() != "resend" {
		t.Fatal("resend should win")
	}
	if (&Client{sendgridKey: "b"}).Provider() != "sendgrid" {
		t.Fatal("sendgrid fallback")
	}
	if (&Client{}).Provider() != "log" {
		t.Fatal("log")
	}
}

func TestFormatAddressNoName(t *testing.T) {
	if formatAddress("", "a@example.com") != "a@example.com" {
		t.Fatal("bare address")
	}
}
