package payment

import (
	"net/url"
	"strings"
	"testing"
)

func TestProductionCheckoutReturnURLs(t *testing.T) {
	if ProductionSuccessURL != "https://www.remainfaithful.com/?donated=true" {
		t.Fatalf("success URL = %q", ProductionSuccessURL)
	}
	if ProductionCancelURL != "https://www.remainfaithful.com/#donate" {
		t.Fatalf("cancel URL = %q", ProductionCancelURL)
	}

	for _, raw := range []string{ProductionSuccessURL, ProductionCancelURL} {
		if strings.Contains(raw, "remainfaithful.app") {
			t.Fatalf("production return URL must not use remainfaithful.app: %q", raw)
		}
		if strings.Contains(raw, "vercel.app") {
			t.Fatalf("production return URL must not use a Vercel host: %q", raw)
		}
		if strings.Contains(raw, "/thank-you") {
			t.Fatalf("production return URL must use the homepage donated banner, not /thank-you: %q", raw)
		}

		parsed, err := url.Parse(raw)
		if err != nil {
			t.Fatalf("parse %q: %v", raw, err)
		}
		if parsed.Hostname() != "www.remainfaithful.com" {
			t.Fatalf("host = %q, want www.remainfaithful.com", parsed.Hostname())
		}
	}
}
