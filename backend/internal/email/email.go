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
	inviterName = sanitizeMailHeader(inviterName)
	subject := inviterName + " invited you to Remain Faithful"
	plain := fmt.Sprintf(
		"Hi,\n\n%s has invited you to join them as an accountability partner on Remain Faithful.\n\n"+
			"Open the link below on your iPhone to accept. You are not connected until you accept.\n\n%s\n\n"+
			"Remain Faithful is a faith-based accountability app for adults.\n\n"+
			"If you weren't expecting this, you can safely ignore this email.\n\n- Remain Faithful",
		inviterName, acceptURL,
	)
	return c.send(toEmail, "", subject, plain, partnerInviteHTML(inviterName, acceptURL))
}

func partnerInviteHTML(inviterName, acceptURL string) string {
	return fmt.Sprintf(
		`<p>Hi,</p>`+
			`<p><strong>%s</strong> has invited you to join them as an accountability partner on Remain Faithful.</p>`+
			`<p>You are not connected until you accept.</p>`+
			`<p><a href="%s" style="background:#D2AB4C;color:#0F2050;padding:12px 24px;border-radius:8px;`+
			`text-decoration:none;font-weight:bold;display:inline-block;">Accept Invitation</a></p>`+
			`<p>If you weren't expecting this, you can safely ignore this email.</p>`+
			`<p>- Remain Faithful</p>`,
		html.EscapeString(inviterName), html.EscapeString(acceptURL),
	)
}

// SendVerifyEmail asks a new account to confirm the address before invites are accepted.
func (c *Client) SendVerifyEmail(toEmail, toName, verifyURL string) error {
	toName = sanitizeMailHeader(toName)
	subject := "Confirm your Remain Faithful email"
	plain := fmt.Sprintf(
		"Hi %s,\n\nConfirm this email address so partner and group invites sent to it can be accepted.\n\n%s\n\nThe link expires in 24 hours.\n\n- Remain Faithful",
		toName, verifyURL,
	)
	htmlBody := fmt.Sprintf(
		`<p>Hi %s,</p><p>Confirm this email address so partner and group invites sent to it can be accepted.</p>`+
			`<p><a href="%s">Confirm email</a></p><p>The link expires in 24 hours.</p><p>- Remain Faithful</p>`,
		html.EscapeString(toName), html.EscapeString(verifyURL),
	)
	return c.send(toEmail, toName, subject, plain, htmlBody)
}

// SendEmailChange confirms a new address before it replaces the current one.
func (c *Client) SendEmailChange(toEmail, toName, verifyURL string) error {
	toName = sanitizeMailHeader(toName)
	subject := "Confirm your new Remain Faithful email"
	plain := fmt.Sprintf(
		"Hi %s,\n\nConfirm this new email address for your Remain Faithful account. Until you do, the account stays on the previous address.\n\n%s\n\nThe link expires in 24 hours.\n\n- Remain Faithful",
		toName, verifyURL,
	)
	htmlBody := fmt.Sprintf(
		`<p>Hi %s,</p><p>Confirm this new email address. Until you do, the account stays on the previous address.</p>`+
			`<p><a href="%s">Confirm new email</a></p><p>The link expires in 24 hours.</p><p>- Remain Faithful</p>`,
		html.EscapeString(toName), html.EscapeString(verifyURL),
	)
	return c.send(toEmail, toName, subject, plain, htmlBody)
}

// SendGroupInvite sends a group member invitation email.
func (c *Client) SendGroupInvite(toEmail, inviterName, groupName, acceptURL string) error {
	inviterName = sanitizeMailHeader(inviterName)
	groupName = sanitizeMailHeader(groupName)
	subject := inviterName + " invited you to join " + groupName + " on Remain Faithful"
	plain := fmt.Sprintf(
		"Hi,\n\n%s has invited you to join the accountability group \"%s\" on Remain Faithful.\n\n"+
			"Open the link below to accept. You are not in the group until you accept.\n\n%s\n\n"+
			"If you weren't expecting this, you can safely ignore this email.\n\n- Remain Faithful",
		inviterName, groupName, acceptURL,
	)
	htmlBody := fmt.Sprintf(
		`<p>Hi,</p>`+
			`<p><strong>%s</strong> has invited you to join the accountability group `+
			`<strong>%s</strong> on Remain Faithful.</p>`+
			`<p>You are not in the group until you accept.</p>`+
			`<p><a href="%s" style="background:#D2AB4C;color:#0F2050;padding:12px 24px;border-radius:8px;`+
			`text-decoration:none;font-weight:bold;display:inline-block;">Accept Invitation</a></p>`+
			`<p style="color:#888;">If you weren't expecting this, you can safely ignore this email.</p>`+
			`<p>- Remain Faithful</p>`,
		html.EscapeString(inviterName), html.EscapeString(groupName), html.EscapeString(acceptURL),
	)
	return c.send(toEmail, "", subject, plain, htmlBody)
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
	toName = sanitizeMailHeader(toName)
	subject := "Reset your Remain Faithful password"
	plain := fmt.Sprintf(
		"Hi %s,\n\nOpen the link below to reset your password. It expires in 1 hour.\n\n%s\n\n"+
			"If you didn't request this, you can safely ignore this email.\n\n- Remain Faithful",
		toName, resetURL,
	)
	htmlBody := fmt.Sprintf(
		`<p>Hi %s,</p><p>Open the link below to reset your password. It expires in 1 hour.</p>`+
			`<p><a href="%s">Reset Password</a></p>`+
			`<p>If you didn't request this, you can safely ignore this email.</p><p>- Remain Faithful</p>`,
		html.EscapeString(toName), html.EscapeString(resetURL),
	)
	return c.send(toEmail, toName, subject, plain, htmlBody)
}

func (c *Client) send(toEmail, toName, subject, plain, htmlBody string) error {
	if c.apiKey == "" {
		if os.Getenv("APP_ENV") == "production" {
			log.Printf("[sendgrid] SENDGRID_API_KEY is not set; refusing to send subject=%q", subject)
			return fmt.Errorf("email: SENDGRID_API_KEY is not set")
		}
		log.Printf("[sendgrid] SENDGRID_API_KEY not set; dev log subject=%q", subject)
		fmt.Printf("[email] subject=%q\n%s\n", subject, plain)
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
			{"type": "text/html", "value": htmlBody},
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

	log.Printf("[sendgrid] sending subject=%q", subject)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[sendgrid] request error to=%s: %v", toEmail, err)
		return fmt.Errorf("email: send: %w", err)
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("[sendgrid] response status=%d (body unreadable: %v)", resp.StatusCode, readErr)
	} else {
		log.Printf("[sendgrid] response status=%d", resp.StatusCode)
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
