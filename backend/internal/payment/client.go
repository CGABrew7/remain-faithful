package payment

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v76"
	portalsession "github.com/stripe/stripe-go/v76/billingportal/session"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
)

// Production Stripe Checkout return URLs. remainfaithful.app does not resolve,
// and /thank-you is a 404 on remainfaithful.com. The homepage already shows
// DonationSuccessBanner when ?donated=true is present.
const (
	ProductionSuccessURL = "https://www.remainfaithful.com/?donated=true"
	ProductionCancelURL  = "https://www.remainfaithful.com/#donate"
	// ProductionPortalReturnURL is where Billing Portal sends monthly donors
	// after they manage or cancel a subscription.
	ProductionPortalReturnURL = ProductionCancelURL
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

func (c *Client) Enabled() bool { return c != nil && c.enabled }

type CheckoutParams struct {
	AmountCents int64
	Monthly     bool
	UserID      int64
	Email       string
	CustomerID  string
}

// CheckoutIdempotencyKey is a Stripe Idempotency-Key for Checkout Session
// create. It is a SHA-256 of user_id + amount_cents + monthly + UTC day, so
// a double-tap of the same gift on the same UTC day reuses Stripe's 24h
// idempotency window instead of opening a second session.
func CheckoutIdempotencyKey(userID, amountCents int64, monthly bool, at time.Time) string {
	freq := "once"
	if monthly {
		freq = "monthly"
	}
	raw := fmt.Sprintf("rf-donate:%d:%d:%s:%s", userID, amountCents, freq, at.UTC().Format("2006-01-02"))
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func checkoutSessionParams(p CheckoutParams, now time.Time) *stripe.CheckoutSessionParams {
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
		SuccessURL: stripe.String(ProductionSuccessURL),
		CancelURL:  stripe.String(ProductionCancelURL),
	}
	params.AddMetadata("user_id", strconv.FormatInt(p.UserID, 10))
	if p.CustomerID != "" {
		params.Customer = stripe.String(p.CustomerID)
	} else if p.Email != "" {
		params.CustomerEmail = stripe.String(p.Email)
	}
	if p.Monthly {
		params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{}
		params.SubscriptionData.AddMetadata("user_id", strconv.FormatInt(p.UserID, 10))
	} else {
		params.PaymentIntentData = &stripe.CheckoutSessionPaymentIntentDataParams{}
		params.PaymentIntentData.AddMetadata("user_id", strconv.FormatInt(p.UserID, 10))
	}
	params.SetIdempotencyKey(CheckoutIdempotencyKey(p.UserID, p.AmountCents, p.Monthly, now))
	return params
}

func (c *Client) CreateCheckoutSession(ctx context.Context, p CheckoutParams) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("stripe not configured: set STRIPE_SECRET_KEY")
	}

	s, err := session.New(checkoutSessionParams(p, time.Now()))
	if err != nil {
		return "", fmt.Errorf("stripe checkout: %w", err)
	}
	return s.URL, nil
}

func (c *Client) CreateBillingPortalSession(ctx context.Context, customerID string) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("stripe not configured: set STRIPE_SECRET_KEY")
	}
	if customerID == "" {
		return "", fmt.Errorf("stripe customer id required")
	}
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(ProductionPortalReturnURL),
	}
	s, err := portalsession.New(params)
	if err != nil {
		return "", fmt.Errorf("stripe billing portal: %w", err)
	}
	return s.URL, nil
}

// FindOrCreateCustomer returns an existing Stripe customer for email, or
// creates one tagged with user_id metadata.
func (c *Client) FindOrCreateCustomer(ctx context.Context, email string, userID int64) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("stripe not configured: set STRIPE_SECRET_KEY")
	}
	if email == "" {
		return "", fmt.Errorf("email required to find or create a Stripe customer")
	}

	list := &stripe.CustomerListParams{Email: stripe.String(email)}
	list.Limit = stripe.Int64(1)
	iter := customer.List(list)
	if iter.Next() {
		if cust := iter.Customer(); cust != nil && cust.ID != "" {
			return cust.ID, nil
		}
	}
	if err := iter.Err(); err != nil {
		return "", fmt.Errorf("stripe customer list: %w", err)
	}

	create := &stripe.CustomerParams{Email: stripe.String(email)}
	create.AddMetadata("user_id", strconv.FormatInt(userID, 10))
	cust, err := customer.New(create)
	if err != nil {
		return "", fmt.Errorf("stripe customer create: %w", err)
	}
	return cust.ID, nil
}
