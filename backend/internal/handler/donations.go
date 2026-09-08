package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/payment"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// insertDonationEventSQL logs a Stripe event once. ON CONFLICT DO NOTHING makes
// webhook retries idempotent: only a new stripe_event_id runs business logic.
const insertDonationEventSQL = `
	INSERT INTO donation_events (stripe_event_id, type, payload, created_at)
	VALUES ($1, $2, $3, NOW())
	ON CONFLICT (stripe_event_id) DO NOTHING
`

type dbExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// CreateCheckoutSession creates a Stripe Checkout session for a donation.
// Donations do not unlock features — the app stays free.
// POST /donations/create-checkout-session
func (h *H) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())
	email, _ := rfauth.EmailFromContext(r.Context())

	var body struct {
		AmountDollars int  `json:"amount_dollars"`
		Monthly       bool `json:"monthly"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AmountDollars < 1 || body.AmountDollars > 10000 {
		writeError(w, http.StatusBadRequest, "invalid amount")
		return
	}

	url, err := h.Stripe.CreateCheckoutSession(r.Context(), payment.CheckoutParams{
		AmountCents: int64(body.AmountDollars) * 100,
		Monthly:     body.Monthly,
		UserID:      userID,
		Email:       email,
		CustomerID:  h.lookupStripeCustomerID(r.Context(), userID),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create checkout session")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

// CreateBillingPortalSession opens Stripe Billing Portal so a monthly donor
// can update payment methods or cancel. Returns {url}.
// POST /donations/billing-portal
func (h *H) CreateBillingPortalSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := rfauth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if h.Stripe == nil || !h.Stripe.Enabled() {
		writeError(w, http.StatusServiceUnavailable, "payments not configured")
		return
	}

	customerID := h.lookupStripeCustomerID(r.Context(), userID)
	if customerID == "" {
		email, _ := rfauth.EmailFromContext(r.Context())
		if email == "" {
			email = h.lookupUserEmail(r.Context(), userID)
		}
		var err error
		customerID, err = h.Stripe.FindOrCreateCustomer(r.Context(), email, userID)
		if err != nil {
			log.Printf("CreateBillingPortalSession: find/create customer user=%d: %v", userID, err)
			writeError(w, http.StatusBadGateway, "could not open billing portal")
			return
		}
		h.persistUserCustomerID(r.Context(), userID, customerID)
	}

	url, err := h.Stripe.CreateBillingPortalSession(r.Context(), customerID)
	if err != nil {
		log.Printf("CreateBillingPortalSession: portal user=%d: %v", userID, err)
		writeError(w, http.StatusBadGateway, "could not open billing portal")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

// DonationWebhook handles Stripe webhook events.
// POST /donations/webhook  (unauthenticated — verified by Stripe signature)
func (h *H) DonationWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read body")
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	secret := ""
	if h.Stripe != nil {
		secret = h.Stripe.WebhookSecret
	}
	if secret == "" {
		writeError(w, http.StatusServiceUnavailable, "webhook not configured")
		return
	}

	event, err := webhook.ConstructEventWithOptions(payload, sig, secret, webhook.ConstructEventOptions{
		// Dashboard endpoints may pin a different Stripe API version than
		// stripe-go v76. Signature check still applies; we only read ids and
		// metadata from the object.
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid signature")
		return
	}

	// Signature is valid. Always 200 so a handler bug is not retried forever
	// as a 5xx. Stripe retries on non-2xx.
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("DonationWebhook: recovered type=%s id=%s: %v", event.Type, event.ID, rec)
			}
		}()
		h.dispatchDonationEvent(r.Context(), event)
	}()
	writeJSON(w, http.StatusOK, map[string]string{"received": "true"})
}

func (h *H) dispatchDonationEvent(ctx context.Context, event stripe.Event) {
	summary := summarizeEventPayload(event)
	inserted, err := insertDonationEvent(ctx, h.eventExecer(), event.ID, string(event.Type), summary)
	if err != nil {
		log.Printf("DonationWebhook: event log type=%s id=%s: %v", event.Type, event.ID, err)
		// Fail-open: still run handlers; they are themselves idempotent.
	} else if !inserted {
		log.Printf("DonationWebhook: duplicate event id=%s type=%s", event.ID, event.Type)
		return
	}

	log.Printf("DonationWebhook: type=%s id=%s", event.Type, event.ID)

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var s stripe.CheckoutSession
		if err := unmarshalEventObject(event, &s); err != nil {
			log.Printf("DonationWebhook: checkout.session.completed unmarshal: %v", err)
			return
		}
		h.handleCompletedCheckout(ctx, &s)
	case stripe.EventTypeInvoicePaymentFailed:
		h.handleInvoicePaymentFailed(ctx, event)
	case stripe.EventTypeChargeRefunded:
		h.handleChargeRefunded(ctx, event)
	case stripe.EventTypeChargeRefundUpdated, stripe.EventTypeRefundCreated, stripe.EventTypeRefundUpdated:
		h.handleRefundObject(ctx, event)
	case stripe.EventTypeCustomerSubscriptionDeleted, stripe.EventTypeCustomerSubscriptionUpdated:
		h.handleSubscriptionChange(ctx, event)
	default:
		// Logged above; no further action.
	}
}

func (h *H) handleCompletedCheckout(ctx context.Context, s *stripe.CheckoutSession) {
	userID, err := strconv.ParseInt(s.Metadata["user_id"], 10, 64)
	if err != nil {
		return
	}

	amountCents := s.AmountTotal
	monthly := s.Mode == stripe.CheckoutSessionModeSubscription
	customerID := customerIDOf(s.Customer)
	subscriptionID := subscriptionIDOf(s.Subscription)
	paymentIntentID := paymentIntentIDOf(s.PaymentIntent)

	if h.DB != nil {
		_, err = h.DB.ExecContext(ctx, `
			INSERT INTO donations (
				user_id, stripe_session_id, amount_cents, monthly, created_at,
				status, stripe_customer_id, stripe_subscription_id, stripe_payment_intent_id
			)
			VALUES ($1, $2, $3, $4, $5, 'completed', $6, NULLIF($7, ''), NULLIF($8, ''))
			ON CONFLICT (stripe_session_id) DO NOTHING
		`, userID, s.ID, amountCents, monthly, time.Now(), nullIfEmpty(customerID), subscriptionID, paymentIntentID)
		if err != nil {
			log.Printf("handleCompletedCheckout: insert donation user=%d: %v", userID, err)
		}
		h.persistUserCustomerID(ctx, userID, customerID)
		if monthly && subscriptionID != "" {
			h.upsertStripeSubscription(ctx, userID, customerID, subscriptionID, "active")
		}
	}

	amountDollars := amountCents / 100
	freq := "one-time"
	if monthly {
		freq = "monthly"
	}
	alertBody := "Your " + freq + " gift of $" + strconv.FormatInt(amountDollars, 10) +
		" helps keep RF free for everyone. Thank you!"
	h.pushDonationAlert(ctx, userID, "Thank you for your support!", alertBody, "DONATION_THANKS")
}

func (h *H) handleInvoicePaymentFailed(ctx context.Context, event stripe.Event) {
	var inv stripe.Invoice
	if err := unmarshalEventObject(event, &inv); err != nil {
		log.Printf("DonationWebhook: invoice.payment_failed unmarshal: %v", err)
		return
	}

	userID := parseUserID(inv.Metadata)
	cid := customerIDOf(inv.Customer)
	subID := subscriptionIDOf(inv.Subscription)
	if userID == 0 {
		userID = h.lookupUserIDByCustomer(ctx, cid)
	}
	if userID == 0 {
		userID = h.lookupUserIDBySubscription(ctx, subID)
	}

	log.Printf("DonationWebhook: invoice.payment_failed customer=%s subscription=%s user=%d", cid, subID, userID)

	if subID != "" && h.DB != nil {
		if _, err := h.DB.ExecContext(ctx, `
			UPDATE donations SET status = 'past_due'
			WHERE stripe_subscription_id = $1
		`, subID); err != nil {
			log.Printf("DonationWebhook: mark past_due: %v", err)
		}
		h.upsertStripeSubscription(ctx, userID, cid, subID, "past_due")
	}

	if userID > 0 {
		h.pushDonationAlert(ctx, userID,
			"Donation payment unsuccessful",
			"We could not process your monthly gift. You can update your payment method anytime — Remain Faithful stays free.",
			"DONATION_PAYMENT_FAILED")
	}
}

func (h *H) handleChargeRefunded(ctx context.Context, event stripe.Event) {
	var ch stripe.Charge
	if err := unmarshalEventObject(event, &ch); err != nil {
		log.Printf("DonationWebhook: charge.refunded unmarshal: %v", err)
		return
	}
	h.markDonationRefunded(ctx, paymentIntentIDOf(ch.PaymentIntent), customerIDOf(ch.Customer), parseUserID(ch.Metadata))
}

func (h *H) handleRefundObject(ctx context.Context, event stripe.Event) {
	var rf stripe.Refund
	if err := unmarshalEventObject(event, &rf); err != nil {
		log.Printf("DonationWebhook: %s unmarshal: %v", event.Type, err)
		return
	}
	pi := paymentIntentIDOf(rf.PaymentIntent)
	var cid string
	var userID int64
	if rf.Charge != nil {
		if pi == "" {
			pi = paymentIntentIDOf(rf.Charge.PaymentIntent)
		}
		cid = customerIDOf(rf.Charge.Customer)
		userID = parseUserID(rf.Charge.Metadata)
	}
	h.markDonationRefunded(ctx, pi, cid, userID)
}

func (h *H) handleSubscriptionChange(ctx context.Context, event stripe.Event) {
	var sub stripe.Subscription
	if err := unmarshalEventObject(event, &sub); err != nil {
		log.Printf("DonationWebhook: %s unmarshal: %v", event.Type, err)
		return
	}
	userID := parseUserID(sub.Metadata)
	cid := customerIDOf(sub.Customer)
	if userID == 0 {
		userID = h.lookupUserIDByCustomer(ctx, cid)
	}
	if userID == 0 {
		userID = h.lookupUserIDBySubscription(ctx, sub.ID)
	}
	status := string(sub.Status)
	if event.Type == stripe.EventTypeCustomerSubscriptionDeleted && status == "" {
		status = "canceled"
	}
	log.Printf("DonationWebhook: %s subscription=%s status=%s user=%d", event.Type, sub.ID, status, userID)
	h.upsertStripeSubscription(ctx, userID, cid, sub.ID, status)
	if donationStatus := donationStatusFromSubscription(status); donationStatus != "" && h.DB != nil && sub.ID != "" {
		if _, err := h.DB.ExecContext(ctx, `
			UPDATE donations SET status = $1 WHERE stripe_subscription_id = $2
		`, donationStatus, sub.ID); err != nil {
			log.Printf("DonationWebhook: sync donation status: %v", err)
		}
	}
}

func (h *H) markDonationRefunded(ctx context.Context, paymentIntentID, customerID string, userID int64) {
	if h.DB == nil {
		return
	}
	if paymentIntentID != "" {
		res, err := h.DB.ExecContext(ctx, `
			UPDATE donations SET status = 'refunded' WHERE stripe_payment_intent_id = $1
		`, paymentIntentID)
		if err != nil {
			log.Printf("DonationWebhook: refund by payment_intent: %v", err)
			return
		}
		if n, _ := res.RowsAffected(); n > 0 {
			return
		}
	}
	if customerID != "" {
		if _, err := h.DB.ExecContext(ctx, `
			UPDATE donations SET status = 'refunded'
			WHERE id = (
				SELECT id FROM donations
				WHERE stripe_customer_id = $1 AND status <> 'refunded'
				ORDER BY created_at DESC LIMIT 1
			)
		`, customerID); err != nil {
			log.Printf("DonationWebhook: refund by customer: %v", err)
		}
		return
	}
	if userID > 0 {
		if _, err := h.DB.ExecContext(ctx, `
			UPDATE donations SET status = 'refunded'
			WHERE id = (
				SELECT id FROM donations
				WHERE user_id = $1 AND status <> 'refunded'
				ORDER BY created_at DESC LIMIT 1
			)
		`, userID); err != nil {
			log.Printf("DonationWebhook: refund by user: %v", err)
		}
	}
}

func (h *H) upsertStripeSubscription(ctx context.Context, userID int64, customerID, subscriptionID, status string) {
	if h.DB == nil || subscriptionID == "" {
		return
	}
	var userArg any
	if userID > 0 {
		userArg = userID
	}
	if status == "" {
		status = "unknown"
	}
	if _, err := h.DB.ExecContext(ctx, `
		INSERT INTO stripe_subscriptions (user_id, stripe_customer_id, stripe_subscription_id, status, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (stripe_subscription_id) DO UPDATE SET
			status = EXCLUDED.status,
			user_id = COALESCE(EXCLUDED.user_id, stripe_subscriptions.user_id),
			stripe_customer_id = COALESCE(NULLIF(EXCLUDED.stripe_customer_id, ''), stripe_subscriptions.stripe_customer_id),
			updated_at = NOW()
	`, userArg, customerID, subscriptionID, status); err != nil {
		log.Printf("upsertStripeSubscription: %v", err)
	}
}

func (h *H) pushDonationAlert(ctx context.Context, userID int64, title, body, notifType string) {
	if h.APNS == nil || h.DB == nil || userID == 0 {
		return
	}
	rows, err := h.DB.QueryContext(ctx,
		`SELECT token FROM device_tokens WHERE user_id = $1 AND is_active = TRUE`, userID)
	if err != nil {
		log.Printf("pushDonationAlert query user=%d: %v", userID, err)
		return
	}
	defer rows.Close()

	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": title,
				"body":  body,
			},
			"sound": "default",
		},
		"notification_type": notifType,
	}

	for rows.Next() {
		var token string
		if rows.Scan(&token) != nil {
			continue
		}
		if err := h.APNS.Send(ctx, &apns.Notification{
			DeviceToken: token,
			PushType:    "alert",
			Priority:    10,
			Payload:     payload,
		}); err != nil {
			var inv *apns.ErrInvalidToken
			if errors.As(err, &inv) {
				h.markTokenInactive(ctx, token)
				continue
			}
			log.Printf("pushDonationAlert send %d: %v", userID, err)
		}
	}
}

func (h *H) lookupStripeCustomerID(ctx context.Context, userID int64) string {
	if h.DB == nil || userID == 0 {
		return ""
	}
	var id sql.NullString
	if err := h.DB.QueryRowContext(ctx,
		`SELECT stripe_customer_id FROM users WHERE id = $1`, userID,
	).Scan(&id); err == nil && id.Valid && id.String != "" {
		return id.String
	}
	if err := h.DB.QueryRowContext(ctx, `
		SELECT stripe_customer_id FROM donations
		WHERE user_id = $1 AND stripe_customer_id IS NOT NULL AND stripe_customer_id <> ''
		ORDER BY created_at DESC LIMIT 1
	`, userID).Scan(&id); err == nil && id.Valid {
		return id.String
	}
	return ""
}

func (h *H) lookupUserEmail(ctx context.Context, userID int64) string {
	if h.DB == nil || userID == 0 {
		return ""
	}
	var email string
	if err := h.DB.QueryRowContext(ctx, `SELECT email FROM users WHERE id = $1`, userID).Scan(&email); err != nil {
		return ""
	}
	return email
}

func (h *H) persistUserCustomerID(ctx context.Context, userID int64, customerID string) {
	if h.DB == nil || userID == 0 || customerID == "" {
		return
	}
	if _, err := h.DB.ExecContext(ctx, `
		UPDATE users SET stripe_customer_id = $1
		WHERE id = $2 AND (stripe_customer_id IS NULL OR stripe_customer_id = '')
	`, customerID, userID); err != nil {
		log.Printf("persistUserCustomerID: %v", err)
	}
}

func (h *H) lookupUserIDByCustomer(ctx context.Context, customerID string) int64 {
	if h.DB == nil || customerID == "" {
		return 0
	}
	var id int64
	if err := h.DB.QueryRowContext(ctx,
		`SELECT id FROM users WHERE stripe_customer_id = $1`, customerID,
	).Scan(&id); err == nil {
		return id
	}
	if err := h.DB.QueryRowContext(ctx, `
		SELECT user_id FROM donations
		WHERE stripe_customer_id = $1
		ORDER BY created_at DESC LIMIT 1
	`, customerID).Scan(&id); err == nil {
		return id
	}
	return 0
}

func (h *H) lookupUserIDBySubscription(ctx context.Context, subscriptionID string) int64 {
	if h.DB == nil || subscriptionID == "" {
		return 0
	}
	var id int64
	if err := h.DB.QueryRowContext(ctx, `
		SELECT user_id FROM stripe_subscriptions
		WHERE stripe_subscription_id = $1 AND user_id IS NOT NULL
	`, subscriptionID).Scan(&id); err == nil {
		return id
	}
	if err := h.DB.QueryRowContext(ctx, `
		SELECT user_id FROM donations
		WHERE stripe_subscription_id = $1
		ORDER BY created_at DESC LIMIT 1
	`, subscriptionID).Scan(&id); err == nil {
		return id
	}
	return 0
}

func (h *H) eventExecer() dbExecer {
	if h == nil || h.DB == nil {
		return nil
	}
	return h.DB
}

func insertDonationEvent(ctx context.Context, db dbExecer, eventID, typ string, payload []byte) (bool, error) {
	if db == nil {
		return false, errors.New("db not configured")
	}
	if eventID == "" {
		return false, errors.New("missing stripe event id")
	}
	if len(payload) == 0 {
		payload = []byte("{}")
	}
	res, err := db.ExecContext(ctx, insertDonationEventSQL, eventID, typ, string(payload))
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func summarizeEventPayload(event stripe.Event) []byte {
	summary := map[string]any{
		"id":   event.ID,
		"type": event.Type,
	}
	if event.Data != nil && len(event.Data.Raw) > 0 {
		var obj map[string]any
		if err := json.Unmarshal(event.Data.Raw, &obj); err == nil {
			for _, key := range []string{
				"id", "object", "customer", "subscription", "payment_intent",
				"amount", "amount_total", "amount_due", "amount_refunded",
				"currency", "status", "mode", "refunded",
			} {
				if v, ok := obj[key]; ok {
					summary[key] = compactStripeRef(v)
				}
			}
			if meta, ok := obj["metadata"].(map[string]any); ok {
				if uid, ok := meta["user_id"]; ok {
					summary["user_id"] = uid
				}
			}
		}
	}
	b, err := json.Marshal(summary)
	if err != nil {
		return []byte(`{"type":"unserializable"}`)
	}
	return b
}

func compactStripeRef(v any) any {
	switch t := v.(type) {
	case string:
		return t
	case map[string]any:
		if id, ok := t["id"]; ok {
			return id
		}
	}
	return v
}

func unmarshalEventObject(event stripe.Event, dest any) error {
	if event.Data == nil {
		return errors.New("missing event data")
	}
	return json.Unmarshal(event.Data.Raw, dest)
}

func parseUserID(meta map[string]string) int64 {
	if meta == nil {
		return 0
	}
	id, err := strconv.ParseInt(meta["user_id"], 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func customerIDOf(c *stripe.Customer) string {
	if c == nil {
		return ""
	}
	return c.ID
}

func subscriptionIDOf(s *stripe.Subscription) string {
	if s == nil {
		return ""
	}
	return s.ID
}

func paymentIntentIDOf(p *stripe.PaymentIntent) string {
	if p == nil {
		return ""
	}
	return p.ID
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func donationStatusFromSubscription(status string) string {
	switch status {
	case "canceled", "unpaid", "incomplete_expired":
		return "canceled"
	case "past_due":
		return "past_due"
	case "active", "trialing":
		return "completed"
	default:
		return ""
	}
}
