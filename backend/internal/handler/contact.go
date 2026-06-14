package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// Contact handles website contact-form submissions.
// POST /contact (unauthenticated)
// Body: { "name": "...", "email": "...", "subject": "...", "message": "..." }
func (h *H) Contact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Name    = strings.TrimSpace(req.Name)
	req.Email   = strings.TrimSpace(req.Email)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" || req.Email == "" || req.Message == "" {
		writeError(w, http.StatusBadRequest, "name, email, and message are required")
		return
	}
	if req.Subject == "" {
		req.Subject = "General"
	}

	toEmail := os.Getenv("CONTACT_TO_EMAIL")
	if toEmail == "" {
		toEmail = "jeff@hanokventures.co"
	}

	if err := h.Email.SendContact(req.Email, req.Name, req.Subject, req.Message, toEmail); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send message — please email us directly")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
