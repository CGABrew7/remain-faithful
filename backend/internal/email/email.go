package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Client sends transactional email via SendGrid's v3 mail/send API.
// When SENDGRID_API_KEY is not set the client logs the email to stdout instead
// (useful for local development without a SendGrid account).
type Client struct {
	apiKey string
	from   string
}

// New returns a Client using SENDGRID_API_KEY and SENDGRID_FROM_EMAIL env vars.
func New() *Client {
	return &Client{
		apiKey: os.Getenv("SENDGRID_API_KEY"),
		from:   getenv("SENDGRID_FROM_EMAIL", "noreply@remainfaithful.app"),
	}
}

// SendPasswordReset sends a reset-link email to the given address.
func (c *Client) SendPasswordReset(toEmail, toName, resetURL string) error {
	subject := "Reset your Remain Faithful password"
	plain := fmt.Sprintf(
		"Hi %s,\n\nClick the link below to reset your password. It expires in 1 hour.\n\n%s\n\n"+
			"If you didn't request this, you can safely ignore this email.\n\n— Remain Faithful",
		toName, resetURL,
	)
	html := fmt.Sprintf(
		`<p>Hi %s,</p><p>Click the link below to reset your password. It expires in 1 hour.</p>`+
			`<p><a href="%s">Reset Password</a></p>`+
			`<p>If you didn't request this, you can safely ignore this email.</p><p>— Remain Faithful</p>`,
		toName, resetURL,
	)
	return c.send(toEmail, toName, subject, plain, html)
}

func (c *Client) send(toEmail, toName, subject, plain, html string) error {
	if c.apiKey == "" {
		fmt.Printf("[email] to=%s subject=%q\n%s\n", toEmail, subject, plain)
		return nil
	}
	payload := map[string]any{
		"personalizations": []map[string]any{{
			"to": []map[string]string{{"email": toEmail, "name": toName}},
		}},
		"from":    map[string]string{"email": c.from, "name": "Remain Faithful"},
		"subject": subject,
		"content": []map[string]string{
			{"type": "text/plain", "value": plain},
			{"type": "text/html", "value": html},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("email: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("email: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("email: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("email: SendGrid returned %d", resp.StatusCode)
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
