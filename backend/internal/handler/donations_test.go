package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/payment"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

const testWebhookSecret = "whsec_test_remain_faithful"

type stubResult struct{ n int64 }

func (s stubResult) LastInsertId() (int64, error) { return 1, nil }
func (s stubResult) RowsAffected() (int64, error) { return s.n, nil }

type idempotentEventExec struct {
	calls     int
	lastQuery string
	lastID    string
}

func (e *idempotentEventExec) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	e.calls++
	e.lastQuery = query
	if len(args) > 0 {
		if id, ok := args[0].(string); ok {
			e.lastID = id
		}
	}
	if e.calls == 1 {
		return stubResult{n: 1}, nil
	}
	return stubResult{n: 0}, nil
}

func signedWebhook(t *testing.T, secret, eventJSON string) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	signed := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload: []byte(eventJSON),
		Secret:  secret,
	})
	req := httptest.NewRequest(http.MethodPost, "/donations/webhook", strings.NewReader(eventJSON))
	req.Header.Set("Stripe-Signature", signed.Header)
	req.Header.Set("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func TestDonationWebhookRejectsInvalidSignature(t *testing.T) {
	h := &H{Stripe: &payment.Client{WebhookSecret: testWebhookSecret}}
	req := httptest.NewRequest(http.MethodPost, "/donations/webhook", strings.NewReader(`{"id":"evt_1"}`))
	req.Header.Set("Stripe-Signature", "t=1,v1=not-a-real-signature")
	rr := httptest.NewRecorder()

	h.DonationWebhook(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "invalid signature") {
		t.Fatalf("body = %q, want invalid signature", rr.Body.String())
	}
}

func TestInsertDonationEventIdempotent(t *testing.T) {
	fake := &idempotentEventExec{}
	ctx := context.Background()
	payload := []byte(`{"id":"evt_abc","type":"invoice.payment_failed"}`)

	inserted, err := insertDonationEvent(ctx, fake, "evt_abc", "invoice.payment_failed", payload)
	if err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if !inserted {
		t.Fatal("first insert should be new")
	}

	inserted, err = insertDonationEvent(ctx, fake, "evt_abc", "invoice.payment_failed", payload)
	if err != nil {
		t.Fatalf("second insert: %v", err)
	}
	if inserted {
		t.Fatal("duplicate stripe_event_id should ON CONFLICT DO NOTHING")
	}
	if fake.calls != 2 {
		t.Fatalf("exec calls = %d, want 2", fake.calls)
	}
	if !strings.Contains(fake.lastQuery, "ON CONFLICT (stripe_event_id) DO NOTHING") {
		t.Fatalf("SQL missing idempotent conflict clause: %s", fake.lastQuery)
	}
	if fake.lastID != "evt_abc" {
		t.Fatalf("event id = %q", fake.lastID)
	}
}

func TestDonationWebhookInvoicePaymentFailedDoesNot500(t *testing.T) {
	h := &H{Stripe: &payment.Client{WebhookSecret: testWebhookSecret}}
	eventJSON := `{
		"id":"evt_pay_fail",
		"object":"event",
		"type":"invoice.payment_failed",
		"data":{"object":{"id":"in_123","object":"invoice","customer":"cus_x","subscription":"sub_x","amount_due":500}}
	}`
	req, rr := signedWebhook(t, testWebhookSecret, eventJSON)

	h.DonationWebhook(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (nil DB / missing user must not 500)", rr.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["received"] != "true" {
		t.Fatalf("body = %#v", body)
	}
}

func TestCreateBillingPortalRequiresJWT(t *testing.T) {
	h := &H{Stripe: &payment.Client{}}
	req := httptest.NewRequest(http.MethodPost, "/donations/billing-portal", http.NoBody)
	rr := httptest.NewRecorder()

	h.CreateBillingPortalSession(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestCreateBillingPortalStripeDisabled(t *testing.T) {
	h := &H{Stripe: &payment.Client{}}
	req := httptest.NewRequest(http.MethodPost, "/donations/billing-portal", http.NoBody)
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 42))
	rr := httptest.NewRecorder()

	h.CreateBillingPortalSession(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), "payments not configured") {
		t.Fatalf("body = %q", body)
	}
}

func TestSummarizeEventPayloadOmitsCardFields(t *testing.T) {
	raw := []byte(`{
		"id":"ch_1",
		"object":"charge",
		"amount":500,
		"customer":"cus_1",
		"payment_method_details":{"card":{"number":"4242424242424242","cvc":"123"}},
		"source":{"number":"4242424242424242"}
	}`)
	event := stripe.Event{
		ID:   "evt_1",
		Type: stripe.EventTypeChargeRefunded,
		Data: &stripe.EventData{Raw: raw},
	}
	out := string(summarizeEventPayload(event))
	for _, banned := range []string{"4242424242424242", "cvc", "payment_method_details", "number"} {
		if strings.Contains(out, banned) {
			t.Fatalf("summary leaked %q: %s", banned, out)
		}
	}
	if !strings.Contains(out, "ch_1") || !strings.Contains(out, "cus_1") {
		t.Fatalf("summary missing safe ids: %s", out)
	}
}
