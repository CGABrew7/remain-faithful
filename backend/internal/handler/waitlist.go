package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxWaitlistEmail = 254
	maxWaitlistName  = 80
)

// waitlistPlatforms are the apps people can sign up to hear about.
var waitlistPlatforms = map[string]bool{"android": true}

// MigrateWaitlist adds the waitlist table. Safe to run on every startup.
func MigrateWaitlist(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS waitlist_signups (
			id         BIGSERIAL   PRIMARY KEY,
			email      TEXT        NOT NULL,
			name       TEXT        NOT NULL DEFAULT '',
			platform   TEXT        NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE (email, platform)
		);
	`)
	return err
}

// normalizeWaitlistEmail trims and lowercases an address and reports whether
// it looks like a real mailbox: one @, a dotted domain, no spaces or hidden
// characters, and within the usual length limits.
func normalizeWaitlistEmail(raw string) (string, bool) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" || len(email) > maxWaitlistEmail || !utf8.ValidString(email) {
		return "", false
	}
	for _, r := range email {
		if r > unicode.MaxASCII || unicode.IsSpace(r) || unicode.IsControl(r) {
			return "", false
		}
	}
	if strings.Count(email, "@") != 1 {
		return "", false
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", false
	}
	at := strings.LastIndex(email, "@")
	local, domain := email[:at], email[at+1:]
	if local == "" || len(local) > 64 {
		return "", false
	}
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Contains(domain, "..") {
		return "", false
	}
	for _, label := range strings.Split(domain, ".") {
		if label == "" || len(label) > 63 || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return "", false
		}
		for _, r := range label {
			if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
				return "", false
			}
		}
	}
	return email, true
}

func cleanWaitlistName(raw string) string {
	name := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(raw))
	name = strings.TrimSpace(name)
	if utf8.RuneCountInString(name) > maxWaitlistName {
		name = string([]rune(name)[:maxWaitlistName])
	}
	return name
}

// JoinWaitlist stores a website signup for news about an app that is not out yet.
// POST /waitlist (unauthenticated, 10 req/min per IP in main.go and 10 req/min per email).
// Body: { "email", "name", "platform", "company"?, "turnstileToken"?, "formAction"?, "loadedAt"? }
// Signing up twice is fine: the second call answers the same way and changes nothing.
func (h *H) JoinWaitlist(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email          string `json:"email"`
		Name           string `json:"name"`
		Platform       string `json:"platform"`
		Company        string `json:"company"` // honeypot; people never fill this in
		TurnstileToken string `json:"turnstileToken"`
		FormAction     string `json:"formAction"`
		LoadedAt       int64  `json:"loadedAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if blockPublicForm(w, r, publicForm{
		Company:        req.Company,
		Email:          req.Email,
		TurnstileToken: req.TurnstileToken,
		FormAction:     req.FormAction,
		LoadedAt:       req.LoadedAt,
	}, "waitlist") {
		return
	}
	email, ok := normalizeWaitlistEmail(req.Email)
	if !ok {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address.")
		return
	}
	platform := strings.ToLower(strings.TrimSpace(req.Platform))
	if platform == "" {
		platform = "android"
	}
	if !waitlistPlatforms[platform] {
		writeError(w, http.StatusBadRequest, "unknown platform")
		return
	}
	name := cleanWaitlistName(req.Name)

	if _, err := h.DB.ExecContext(r.Context(),
		`INSERT INTO waitlist_signups (email, name, platform)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (email, platform) DO NOTHING`,
		email, name, platform,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "could not save your signup")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
