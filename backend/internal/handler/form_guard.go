package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"remain-faithful/backend/internal/clientaddr"
)

const (
	turnstileSiteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	formMinSubmit          = 2 * time.Second
	formMaxSubmit          = 24 * time.Hour
	formEmailWindow        = time.Minute
	formTooFastMessage     = "Please wait a moment and try again."
	formRetryMessage       = "Please try again."
	formEmailLimitMessage  = "Too many requests. Try again in a minute."
)

// turnstileVerifyURL is the siteverify endpoint. Tests point it at a local server.
var turnstileVerifyURL = turnstileSiteverifyURL

var turnstileHTTP = &http.Client{Timeout: 10 * time.Second}

// formEmailLimit matches the per-IP window on public form routes (10 per minute).
var formEmailLimit = 10

type publicForm struct {
	Company        string
	Email          string
	TurnstileToken string
	FormAction     string
	LoadedAt       int64
}

type emailWindow struct {
	count int
	start time.Time
}

var (
	formEmailMu      sync.Mutex
	formEmailWindows = map[string]*emailWindow{}

	turnstileSkipMu     sync.Mutex
	turnstileSkipLogged bool
)

func logTurnstileSkip() {
	turnstileSkipMu.Lock()
	defer turnstileSkipMu.Unlock()
	if turnstileSkipLogged {
		return
	}
	turnstileSkipLogged = true
	log.Printf("[turnstile] TURNSTILE_SECRET_KEY is not set; skipping verification")
}

func resetFormEmailWindows() {
	formEmailMu.Lock()
	formEmailWindows = map[string]*emailWindow{}
	formEmailMu.Unlock()
}

func formEmailKey(raw string) string {
	if email, ok := normalizeWaitlistEmail(raw); ok {
		return email
	}
	key := strings.ToLower(strings.TrimSpace(raw))
	if len(key) > 320 {
		key = key[:320]
	}
	return key
}

// allowFormEmail is a fixed window of formEmailLimit submissions per minute
// for one address. Empty addresses are not counted here; the IP limiter covers them.
func allowFormEmail(raw string) bool {
	key := formEmailKey(raw)
	if key == "" {
		return true
	}
	formEmailMu.Lock()
	defer formEmailMu.Unlock()
	now := time.Now()
	if len(formEmailWindows) > 4096 {
		cutoff := now.Add(-formEmailWindow)
		for existing, window := range formEmailWindows {
			if window.start.Before(cutoff) {
				delete(formEmailWindows, existing)
			}
		}
	}
	window := formEmailWindows[key]
	if window == nil || now.Sub(window.start) >= formEmailWindow {
		formEmailWindows[key] = &emailWindow{count: 1, start: now}
		return true
	}
	window.count++
	return window.count <= formEmailLimit
}

// formSubmittedInTime reports whether loadedAt is a millisecond timestamp
// at least formMinSubmit before now and not older than formMaxSubmit.
func formSubmittedInTime(loadedAt int64, now time.Time) bool {
	if loadedAt < 1_000_000_000_000 {
		return false
	}
	elapsed := now.Sub(time.UnixMilli(loadedAt))
	if elapsed < formMinSubmit || elapsed > formMaxSubmit {
		return false
	}
	return true
}

func actionAllowed(action string, allowed []string) bool {
	if action == "" {
		return false
	}
	for _, candidate := range allowed {
		if action == candidate {
			return true
		}
	}
	return false
}

// turnstileOK checks the widget token when TURNSTILE_SECRET_KEY is set.
// An unset secret skips the check so forms keep working before keys are added.
func turnstileOK(ctx context.Context, token, action, remoteIP string, allowed []string) bool {
	secret := strings.TrimSpace(os.Getenv("TURNSTILE_SECRET_KEY"))
	if secret == "" {
		logTurnstileSkip()
		return true
	}
	if token == "" || len(token) > 2048 || !actionAllowed(action, allowed) {
		return false
	}
	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", token)
	if ip := net.ParseIP(strings.TrimSpace(remoteIP)); ip != nil {
		form.Set("remoteip", ip.String())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, turnstileVerifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := turnstileHTTP.Do(req)
	if err != nil {
		log.Printf("[turnstile] siteverify request failed")
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("[turnstile] siteverify status %d", resp.StatusCode)
		return false
	}
	var parsed struct {
		Success bool   `json:"success"`
		Action  string `json:"action"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return false
	}
	return parsed.Success && parsed.Action == action
}

// blockPublicForm applies the honeypot, submit-time, per-email, and Turnstile
// checks. It writes the response and returns true when the handler should stop.
func blockPublicForm(w http.ResponseWriter, r *http.Request, in publicForm, allowed ...string) bool {
	if strings.TrimSpace(in.Company) != "" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return true
	}
	if !formSubmittedInTime(in.LoadedAt, time.Now()) {
		writeError(w, http.StatusBadRequest, formTooFastMessage)
		return true
	}
	if !allowFormEmail(in.Email) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, formEmailLimitMessage)
		return true
	}
	if !turnstileOK(r.Context(), strings.TrimSpace(in.TurnstileToken), strings.TrimSpace(in.FormAction), clientaddr.FromRequest(r), allowed) {
		writeError(w, http.StatusBadRequest, formRetryMessage)
		return true
	}
	return false
}
