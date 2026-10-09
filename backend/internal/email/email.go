package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

// SendPartnerInvite sends an accountability partner invitation email.
func (c *Client) SendPartnerInvite(toEmail, inviterName, acceptURL string) error {
	subject := inviterName + " invited you to Remain Faithful"
	plain := fmt.Sprintf(
		"Hi,\n\n%s has invited you to join them as an accountability partner on Remain Faithful.\n\n"+
			"Click the link below to accept the invitation, create your account, and get started:\n\n%s\n\n"+
			"Remain Faithful is a faith-based accountability app that helps you and a trusted partner\n"+
			"walk together in integrity.\n\n"+
			"If you weren't expecting this, you can safely ignore this email.\n\n— Remain Faithful",
		inviterName, acceptURL,
	)
	html := fmt.Sprintf(
		`<p>Hi,</p>`+
			`<p><strong>%s</strong> has invited you to join them as an accountability partner on Remain Faithful.</p>`+
			`<p><a href="%s" style="background:#D2AB4C;color:#0F2050;padding:12px 24px;border-radius:8px;`+
			`text-decoration:none;font-weight:bold;display:inline-block;">Accept Invitation</a></p>`+
			`<p>Remain Faithful is a faith-based accountability app that helps you and a trusted partner `+
			`walk together in integrity.</p>`+
			`<p style="color:#888;">If you weren't expecting this, you can safely ignore this email.</p>`+
			`<p>— Remain Faithful</p>`,
		inviterName, acceptURL,
	)
	return c.send(toEmail, "", subject, plain, html)
}

// SendGroupInvite sends a group member invitation email.
func (c *Client) SendGroupInvite(toEmail, inviterName, groupName, acceptURL string) error {
	subject := inviterName + " invited you to join " + groupName + " on Remain Faithful"
	plain := fmt.Sprintf(
		"Hi,\n\n%s has invited you to join the accountability group \"%s\" on Remain Faithful.\n\n"+
			"Click the link below to accept the invitation and join the group:\n\n%s\n\n"+
			"If you weren't expecting this, you can safely ignore this email.\n\n— Remain Faithful",
		inviterName, groupName, acceptURL,
	)
	html := fmt.Sprintf(
		`<p>Hi,</p>`+
			`<p><strong>%s</strong> has invited you to join the accountability group `+
			`<strong>"%s"</strong> on Remain Faithful.</p>`+
			`<p><a href="%s" style="background:#D2AB4C;color:#0F2050;padding:12px 24px;border-radius:8px;`+
			`text-decoration:none;font-weight:bold;display:inline-block;">Join Group</a></p>`+
			`<p style="color:#888;">If you weren't expecting this, you can safely ignore this email.</p>`+
			`<p>— Remain Faithful</p>`,
		inviterName, groupName, acceptURL,
	)
	return c.send(toEmail, "", subject, plain, html)
}

// SendPasswordReset sends a reset-link email to the given address.
// SendContact forwards a contact-form submission to the app's support inbox.
// Name, email, subject, and message are HTML-escaped before they are placed
// in the HTML part. The subject header has newlines stripped so a crafted
// subject cannot inject extra mail headers.
func (c *Client) SendContact(fromEmail, fromName, subject, message, toEmail string) error {
	subjectLine, plain, htmlBody := contactParts(fromEmail, fromName, subject, message)
	return c.send(toEmail, "", subjectLine, plain, htmlBody)
}

// contactParts builds the subject, plain text, and HTML for a contact email.
func contactParts(fromEmail, fromName, subject, message string) (subjectLine, plain, htmlBody string) {
	cleanSubject := sanitizeMailHeader(subject)
	message = strings.ReplaceAll(message, "\r\n", "\n")
	message = strings.ReplaceAll(message, "\r", "\n")
	subjectLine = "[Remain Faithful Contact] " + cleanSubject
	plain = fmt.Sprintf("From: %s <%s>\nSubject: %s\n\n%s", fromName, fromEmail, cleanSubject, message)
	htmlBody = fmt.Sprintf(
		`<p><strong>From:</strong> %s &lt;%s&gt;</p><p><strong>Subject:</strong> %s</p><p>%s</p>`,
		html.EscapeString(fromName),
		html.EscapeString(fromEmail),
		html.EscapeString(cleanSubject),
		strings.ReplaceAll(html.EscapeString(message), "\n", "<br>"),
	)
	return subjectLine, plain, htmlBody
}

func sanitizeMailHeader(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

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
		log.Printf("[sendgrid] SENDGRID_API_KEY not set — logging email instead of sending to=%s subject=%q", toEmail, subject)
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
		log.Printf("[sendgrid] error marshaling payload to=%s: %v", toEmail, err)
		return fmt.Errorf("email: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		log.Printf("[sendgrid] error building request to=%s: %v", toEmail, err)
		return fmt.Errorf("email: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[sendgrid] sending request to=%s from=%s subject=%q", toEmail, c.from, subject)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[sendgrid] request error to=%s: %v", toEmail, err)
		return fmt.Errorf("email: send: %w", err)
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("[sendgrid] response status=%d to=%s but failed to read body: %v", resp.StatusCode, toEmail, readErr)
	} else {
		log.Printf("[sendgrid] response status=%d to=%s body=%s", resp.StatusCode, toEmail, string(respBody))
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("email: SendGrid returned %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
