package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	rfauth "remain-faithful/backend/internal/auth"
)

// SubmitFeedback emails the signed-in user's idea, improvement, or problem to
// the same inbox as the website form. Name and email come from the account.
// App version, device model, and OS version are included only when the client
// sends them (the app's toggle is off by default).
// POST /feedback
// Body: { "type", "summary", "details", "app_version"?, "device_model"?, "os_version"?, "company"? }
func (h *H) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	userID, _ := rfauth.UserIDFromContext(r.Context())

	var req struct {
		Type        string `json:"type"`
		Summary     string `json:"summary"`
		Details     string `json:"details"`
		AppVersion  string `json:"app_version"`
		DeviceModel string `json:"device_model"`
		OSVersion   string `json:"os_version"`
		Company     string `json:"company"` // honeypot
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Company) != "" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	subject, message, errMsg := normalizeFeedback(
		req.Type, req.Summary, req.Details, req.AppVersion, req.DeviceModel, req.OSVersion,
	)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	var name, email string
	if err := h.DB.QueryRowContext(r.Context(),
		`SELECT name, email FROM users WHERE id = $1`, userID,
	).Scan(&name, &email); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	toEmail := os.Getenv("CONTACT_TO_EMAIL")
	if toEmail == "" {
		toEmail = "jeff@hanokventures.co"
	}
	if err := h.Email.SendContact(email, name, subject, message, toEmail); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send message — please email us directly")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func normalizeFeedback(typ, summary, details, appVersion, deviceModel, osVersion string) (subject, message, errMsg string) {
	typ = strings.TrimSpace(typ)
	summary = strings.TrimSpace(summary)
	details = strings.TrimSpace(details)
	switch typ {
	case "Idea", "Improvement", "Problem":
	default:
		return "", "", "type must be Idea, Improvement, or Problem"
	}
	if summary == "" || len([]rune(summary)) > 120 {
		return "", "", "summary is required (120 characters max)"
	}
	if details == "" || len([]rune(details)) > 4000 {
		return "", "", "details are required (4000 characters max)"
	}

	subject = "[" + typ + "] " + summary
	var b strings.Builder
	b.WriteString("Type: " + typ + "\n\n")
	b.WriteString(details)
	appVersion = clip(strings.TrimSpace(appVersion), 80)
	deviceModel = clip(strings.TrimSpace(deviceModel), 80)
	osVersion = clip(strings.TrimSpace(osVersion), 40)
	if appVersion != "" || deviceModel != "" || osVersion != "" {
		b.WriteString("\n\n—\n")
		if appVersion != "" {
			b.WriteString("App: " + appVersion + "\n")
		}
		if deviceModel != "" {
			b.WriteString("Device: " + deviceModel + "\n")
		}
		if osVersion != "" {
			b.WriteString("OS: " + osVersion + "\n")
		}
	}
	return subject, b.String(), ""
}

func clip(s string, maxRunes int) string {
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	return string(r[:maxRunes])
}
