package clientaddr

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnsetSecretIgnoresProxyHeaders(t *testing.T) {
	t.Setenv("SITE_PROXY_SECRET", "")
	t.Setenv("APP_ENV", "production")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.8:1234"
	req.Header.Set(IPHeader, "203.0.113.5")
	req.Header.Set(SecretHeader, "anything")
	req.Header.Set("X-Forwarded-For", "198.51.100.9")
	if got := FromRequest(req); got == "203.0.113.5" || got == "198.51.100.9" {
		t.Fatalf("untrusted headers were used: %q", got)
	}
}

func TestMatchingSecretUsesVisitorAddress(t *testing.T) {
	t.Setenv("SITE_PROXY_SECRET", "shared-secret")
	t.Setenv("APP_ENV", "production")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.8:1234"
	req.Header.Set("Fly-Client-IP", "9.9.9.9")
	req.Header.Set(SecretHeader, "shared-secret")
	req.Header.Set(IPHeader, "203.0.113.5, 10.0.0.1")
	if got := FromRequest(req); got != "203.0.113.5" {
		t.Fatalf("visitor address = %q", got)
	}
}

func TestWrongSecretKeepsFlyAddress(t *testing.T) {
	t.Setenv("SITE_PROXY_SECRET", "shared-secret")
	t.Setenv("APP_ENV", "production")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Fly-Client-IP", "9.9.9.9")
	req.Header.Set(SecretHeader, "nope")
	req.Header.Set(IPHeader, "203.0.113.5")
	if got := FromRequest(req); got != "9.9.9.9" {
		t.Fatalf("address = %q", got)
	}
}

func TestBadVisitorAddressFallsThrough(t *testing.T) {
	t.Setenv("SITE_PROXY_SECRET", "shared-secret")
	t.Setenv("APP_ENV", "production")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Fly-Client-IP", "9.9.9.9")
	req.Header.Set(SecretHeader, "shared-secret")
	req.Header.Set(IPHeader, "not an ip")
	if got := FromRequest(req); got != "9.9.9.9" {
		t.Fatalf("address = %q", got)
	}
}
