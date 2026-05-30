package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"remain-faithful/backend/internal/anthropic"
	"remain-faithful/backend/internal/apns"
	"remain-faithful/backend/internal/payment"
)

// APNSSender is the interface used by handlers to send push notifications.
// *apns.Client satisfies this interface.
type APNSSender interface {
	Send(ctx context.Context, n *apns.Notification) error
}

// EmailSender is the interface used by handlers to send transactional email.
type EmailSender interface {
	SendPasswordReset(toEmail, toName, resetURL string) error
}

// Classifier is the interface used by handlers to classify content via Claude.
type Classifier interface {
	Classify(ctx context.Context, text string) (*anthropic.ClassificationResult, error)
}

// H holds shared dependencies for all HTTP handlers.
type H struct {
	DB     *sql.DB
	APNS   APNSSender
	Email  EmailSender
	Claude Classifier
	Stripe *payment.Client
}

// writeJSON serialises v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError sends a JSON error body.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
