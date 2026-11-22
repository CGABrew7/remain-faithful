package email

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGlobalDailyCeiling(t *testing.T) {
	resetGlobalSendForTests()
	t.Cleanup(resetGlobalSendForTests)
	globalMailPerDay = 2
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"ok"}`))
	}))
	defer srv.Close()
	c := &Client{
		resendKey:  "re_test",
		from:       "noreply@remainfaithful.com",
		resendURL:  srv.URL,
		httpClient: srv.Client(),
	}
	if err := c.SendPasswordReset("a@example.com", "Ada", "https://example.com/reset"); err != nil {
		t.Fatal(err)
	}
	if err := c.SendPasswordReset("b@example.com", "Bea", "https://example.com/reset"); err != nil {
		t.Fatal(err)
	}
	err := c.SendPasswordReset("c@example.com", "Cam", "https://example.com/reset")
	if err != ErrDailyCeiling {
		t.Fatalf("third send err = %v, want ceiling", err)
	}
	globalMailDay = time.Now().UTC().Add(-48 * time.Hour).Format("2006-01-02")
	if err := c.SendPasswordReset("d@example.com", "Dee", "https://example.com/reset"); err != nil {
		t.Fatalf("next day send: %v", err)
	}
}
