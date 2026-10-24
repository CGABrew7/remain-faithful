package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"remain-faithful/backend/internal/apns"
	"remain-faithful/backend/internal/fcm"
	"remain-faithful/backend/internal/payment"
)

// APNSSender is the interface used by handlers to send push notifications.
// *apns.Client satisfies this interface.
type APNSSender interface {
	Send(ctx context.Context, n *apns.Notification) error
	IsNoop() bool
	Environment() string
}

// FCMSender delivers the same alert to Android tokens. *fcm.Client satisfies it.
// Nil means this process has not wired FCM. A non-nil no-op client is used
// when credentials are unset so iOS sends are unchanged.
type FCMSender interface {
	Send(ctx context.Context, msg *fcm.Message) error
	IsNoop() bool
}

// EmailSender is the interface used by handlers to send transactional email.
type EmailSender interface {
	SendPasswordReset(toEmail, toName, resetURL string) error
	SendPartnerInvite(toEmail, inviterName, acceptURL string) error
	SendGroupInvite(toEmail, inviterName, groupName, acceptURL string) error
	SendContact(fromEmail, fromName, subject, message, toEmail string) error
	SendVerifyEmail(toEmail, toName, verifyURL string) error
	SendEmailChange(toEmail, toName, verifyURL string) error
}

// H holds shared dependencies for all HTTP handlers.
type H struct {
	DB     *sql.DB
	APNS   APNSSender
	FCM    FCMSender
	Email  EmailSender
	Stripe *payment.Client
	// SyncPush runs notification side effects on the request goroutine.
	// Production leaves it false. Tests set it so they can assert sends.
	SyncPush bool
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

// writeCodedError sends a JSON error the app can branch on.
// Both fields are strings so existing clients that read error still work.
func writeCodedError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]string{"error": msg, "code": code})
}
