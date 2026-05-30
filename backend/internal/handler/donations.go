package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/apns"
	"remain-faithful/backend/internal/payment"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// CreateCheckoutSession creates a Stripe Checkout session for a donation.
// POST /donations/create-checkout-session
func (h *H) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var body struct {
		AmountDollars int  `json:"amount_dollars"`
		Monthly       bool `json:"monthly"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AmountDollars < 1 {
		writeError(w, http.StatusBadRequest, "invalid amount")
		return
	}

	url, err := h.Stripe.CreateCheckoutSession(r.Context(), payment.CheckoutParams{
		AmountCents: int64(body.AmountDollars) * 100,
		Monthly:     body.Monthly,
		UserID:      userID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create checkout session")
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
	secret := h.Stripe.WebhookSecret
	if secret == "" {
		writeError(w, http.StatusServiceUnavailable, "webhook not configured")
		return
	}

	event, err := webhook.ConstructEvent(payload, sig, secret)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid signature")
		return
	}

	if event.Type == "checkout.session.completed" {
		var s stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &s); err == nil {
			h.handleCompletedCheckout(r, &s)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"received": "true"})
}

func (h *H) handleCompletedCheckout(r *http.Request, s *stripe.CheckoutSession) {
	userID, err := strconv.ParseInt(s.Metadata["user_id"], 10, 64)
	if err != nil {
		return
	}

	amountCents := s.AmountTotal
	monthly := s.Mode == stripe.CheckoutSessionModeSubscription

	_, err = h.DB.ExecContext(r.Context(), `
		INSERT INTO donations (user_id, stripe_session_id, amount_cents, monthly, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (stripe_session_id) DO NOTHING
	`, userID, s.ID, amountCents, monthly, time.Now())
	if err != nil {
		return
	}

	rows, err := h.DB.QueryContext(r.Context(),
		`SELECT token FROM device_tokens WHERE user_id = $1 AND is_active = TRUE`, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	amountDollars := amountCents / 100
	freq := "one-time"
	if monthly {
		freq = "monthly"
	}
	alertBody := "Your " + freq + " gift of $" + strconv.FormatInt(amountDollars, 10) +
		" helps keep RF free for everyone. Thank you!"

	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Thank you for your support!",
				"body":  alertBody,
			},
			"sound": "default",
		},
		"notification_type": "DONATION_THANKS",
	}

	for rows.Next() {
		var token string
		if rows.Scan(&token) != nil {
			continue
		}
		_ = h.APNS.Send(r.Context(), &apns.Notification{
			DeviceToken: token,
			PushType:    "alert",
			Priority:    10,
			Payload:     payload,
		})
	}
}
