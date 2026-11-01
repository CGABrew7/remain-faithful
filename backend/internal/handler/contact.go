package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

// Contact handles website contact-form submissions.
// POST /contact (unauthenticated; 10 req/min/IP via authRateLimiter in main.go,
// plus 10 req/min for the same email).
// Body: { "name", "email", "subject", "message", "company"?, "turnstileToken"?, "formAction"?, "loadedAt"? }
func (h *H) Contact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		Subject        string `json:"subject"`
		Message        string `json:"message"`
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
	}, "contact", "church-pilot", "feedback") {
		return
	}
	name, nameErr := cleanPersonName(req.Name)
	if nameErr != "" {
		writeError(w, http.StatusBadRequest, nameErr)
		return
	}
	email, valid := normalizeWaitlistEmail(req.Email)
	if !valid {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address.")
		return
	}
	subject := strings.TrimSpace(req.Subject)
	message := strings.TrimSpace(req.Message)
	if containsDisallowedText(subject, false) || containsDisallowedText(message, true) {
		writeError(w, http.StatusBadRequest, "That message contains characters that are not allowed.")
		return
	}
	if message == "" {
		writeError(w, http.StatusBadRequest, "name, email, and message are required")
		return
	}
	if utf8.RuneCountInString(subject) > maxContactSubject {
		writeError(w, http.StatusBadRequest, "Subject must be 120 characters or fewer.")
		return
	}
	if utf8.RuneCountInString(message) > maxContactMessage {
		writeError(w, http.StatusBadRequest, "Message must be 4000 characters or fewer.")
		return
	}
	if subject == "" {
		subject = "General"
	}

	toEmail := strings.TrimSpace(os.Getenv("CONTACT_TO_EMAIL"))
	if toEmail == "" {
		log.Printf("[contact] CONTACT_TO_EMAIL is not set")
		writeError(w, http.StatusInternalServerError, "We could not send your message. Please email support@remainfaithful.com.")
		return
	}

	if err := h.Email.SendContact(email, name, subject, message, toEmail); err != nil {
		writeError(w, http.StatusInternalServerError, "We could not send your message. Please email support@remainfaithful.com.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
