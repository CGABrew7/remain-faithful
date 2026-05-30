package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Classify runs Tier 3 cloud content classification on a text snippet extracted
// from the device screen. It never receives raw screenshots — only text.
// POST /classify
// Body: { "text": "...", "context": "..." }
// Response: { "category": "...", "confidence": 0.0-1.0, "severity": "..." }
func (h *H) Classify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text    string `json:"text"`
		Context string `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	if len(req.Text) > 2000 {
		req.Text = req.Text[:2000]
	}

	if h.Claude == nil {
		writeError(w, http.StatusServiceUnavailable, "classification service not configured")
		return
	}

	result, err := h.Claude.Classify(r.Context(), req.Text)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "classification failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"category":   result.Category,
		"confidence": result.Confidence,
		"severity":   result.Severity,
	})
}
