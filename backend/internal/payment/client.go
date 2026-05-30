package payment

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

// Client wraps the Stripe SDK for checkout session creation.
// When STRIPE_SECRET_KEY is absent the client is disabled and returns errors.
type Client struct {
	WebhookSecret string
	enabled       bool
}

func New() *Client {
	key := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = key
	return &Client{
		WebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		enabled:       key != "",
	}
}

func (c *Client) Enabled() bool { return c.enabled }

type CheckoutParams struct {
	AmountCents int64
	Monthly     bool
	UserID      int64
}

func (c *Client) CreateCheckoutSession(ctx context.Context, p CheckoutParams) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("stripe not configured: set STRIPE_SECRET_KEY")
	}

	priceData := &stripe.CheckoutSessionLineItemPriceDataParams{
		Currency: stripe.String("usd"),
		ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
			Name:        stripe.String("Remain Faithful — Support the Mission"),
			Description: stripe.String("Helping keep RF free for men pursuing sexual integrity"),
		},
		UnitAmount: stripe.Int64(p.AmountCents),
	}

	mode := stripe.String(string(stripe.CheckoutSessionModePayment))
	if p.Monthly {
		mode = stripe.String(string(stripe.CheckoutSessionModeSubscription))
		priceData.Recurring = &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
			Interval: stripe.String(string(stripe.PriceRecurringIntervalMonth)),
		}
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{PriceData: priceData, Quantity: stripe.Int64(1)},
		},
		Mode:       mode,
		SuccessURL: stripe.String("https://remainfaithful.app/thank-you"),
		CancelURL:  stripe.String("https://remainfaithful.app/"),
	}
	params.AddMetadata("user_id", strconv.FormatInt(p.UserID, 10))

	s, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("stripe checkout: %w", err)
	}
	return s.URL, nil
}
