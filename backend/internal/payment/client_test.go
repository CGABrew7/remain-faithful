package payment

import (
	"net/url"
	"strings"
	"testing"
	"time"
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

func TestCheckoutIdempotencyKeyStableWithinUTCDay(t *testing.T) {
	day := time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC)
	a := CheckoutIdempotencyKey(7, 2500, true, day)
	b := CheckoutIdempotencyKey(7, 2500, true, day.Add(13*time.Hour))
	if a == "" || a != b {
		t.Fatalf("same-day keys differ: %q vs %q", a, b)
	}
	nextDay := CheckoutIdempotencyKey(7, 2500, true, day.Add(24*time.Hour))
	if a == nextDay {
		t.Fatal("UTC day bucket should rotate the key")
	}
	otherAmount := CheckoutIdempotencyKey(7, 5000, true, day)
	if a == otherAmount {
		t.Fatal("amount should be part of the key")
	}
	oneTime := CheckoutIdempotencyKey(7, 2500, false, day)
	if a == oneTime {
		t.Fatal("monthly flag should be part of the key")
	}
}

func TestCheckoutSessionParamsSetsIdempotencyKey(t *testing.T) {
	now := time.Date(2026, 9, 8, 15, 0, 0, 0, time.UTC)
	params := checkoutSessionParams(CheckoutParams{
		AmountCents: 2500,
		Monthly:     true,
		UserID:      9,
		Email:       "donor@example.com",
	}, now)
	want := CheckoutIdempotencyKey(9, 2500, true, now)
	if params.IdempotencyKey == nil || *params.IdempotencyKey != want {
		t.Fatalf("IdempotencyKey = %v, want %q", params.IdempotencyKey, want)
	}
	if params.CustomerEmail == nil || *params.CustomerEmail != "donor@example.com" {
		t.Fatalf("CustomerEmail = %v", params.CustomerEmail)
	}
	if params.SubscriptionData == nil || params.SubscriptionData.Metadata["user_id"] != "9" {
		t.Fatal("subscription metadata user_id missing")
	}
}

func TestBillingPortalReturnURLIsDonateSection(t *testing.T) {
	if ProductionPortalReturnURL != ProductionCancelURL {
		t.Fatalf("portal return = %q, want %q", ProductionPortalReturnURL, ProductionCancelURL)
	}
	if !strings.Contains(ProductionPortalReturnURL, "remainfaithful.com") || !strings.Contains(ProductionPortalReturnURL, "#donate") {
		t.Fatalf("portal return should be the remainfaithful.com donate section: %q", ProductionPortalReturnURL)
	}
}
