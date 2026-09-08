package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rfauth "remain-faithful/backend/internal/auth"
)

func TestIsAllowedEventSummary(t *testing.T) {
	fixed := []string{
		"Explicit content detected",
		"Gambling content detected",
		"Violent content detected",
		"Self-harm content detected",
		"Content reviewed — no concerns",
		"Explicit image detected (perceptual hash match)",
		"Blocked domain detected",
		"Explicit keyword detected in screen text",
		"Monitored app used for 1+ minute",
		"Monitored app category used for 1+ minute",
		"Monitored app activity detected",
		"App activity detected",
	}
	for _, s := range fixed {
		if !isAllowedEventSummary(s) {
			t.Errorf("expected allowed fixed summary %q", s)
		}
	}

	patterns := []string{
		"Explicit content — detected 1 time in 5 min",
		"Gambling content — detected 2 times in 5 min",
		"Violent content — detected 12 times in 5 min",
		"Self-harm content — detected 3 times in 5 min",
		"Activity — detected 1 time in 5 min",
		"Activity — detected 99 times in 5 min",
	}
	for _, s := range patterns {
		if !isAllowedEventSummary(s) {
			t.Errorf("expected allowed continued-activity summary %q", s)
		}
	}

	rejected := []string{
		"",
		"he is watching porn",
		"Explicit content detected on Safari",
		"OCR text: something explicit",
		"<script>alert(1)</script>",
		"explicit content detected",
		"Content reviewed - no concerns",
		"Explicit content - detected 2 times in 5 min",
		"Explicit content — detected 0 times in 5 min",
		"Explicit content — detected 10000 times in 5 min",
		"Partner needs you — custom note",
		strings.Repeat("x", maxEventSummaryLen+1),
		"Explicit content detected" + strings.Repeat("!", 80),
	}
	for _, s := range rejected {
		if isAllowedEventSummary(s) {
			t.Errorf("expected rejected summary %q", s)
		}
	}
}

func TestValidateCreateEventKeepsCategorySeverityAllowlists(t *testing.T) {
	if msg := validateCreateEvent("adult_content", "severe", "Explicit content detected"); msg != "" {
		t.Fatalf("honest iOS payload rejected: %s", msg)
	}
	if msg := validateCreateEvent("app_usage", "concerning", "Monitored app used for 1+ minute"); msg != "" {
		t.Fatalf("device-activity payload rejected: %s", msg)
	}
	if msg := validateCreateEvent("not_a_category", "severe", "Explicit content detected"); msg != "invalid category" {
		t.Fatalf("category = %q, want invalid category", msg)
	}
	if msg := validateCreateEvent("adult_content", "high", "Explicit content detected"); msg != "invalid severity" {
		t.Fatalf("severity = %q, want invalid severity", msg)
	}
	if msg := validateCreateEvent("adult_content", "severe", "arbitrary partner text"); msg != "invalid summary" {
		t.Fatalf("summary = %q, want invalid summary", msg)
	}
	if msg := validateCreateEvent("", "severe", "Explicit content detected"); msg != "category, severity, and summary are required" {
		t.Fatalf("empty = %q, want required", msg)
	}
}

func TestCreateEventRejectsUnknownSummary(t *testing.T) {
	h := &H{}
	body := `{"category":"adult_content","severity":"severe","summary":"he is watching porn right now"}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.CreateEvent(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %q", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid summary") {
		t.Fatalf("body = %q, want invalid summary", rr.Body.String())
	}
}

func TestCreateEventRejectsOverlongSummary(t *testing.T) {
	h := &H{}
	summary := strings.Repeat("x", maxEventSummaryLen+1)
	body := `{"category":"adult_content","severity":"severe","summary":"` + summary + `"}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.CreateEvent(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %q", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid summary") {
		t.Fatalf("body = %q, want invalid summary", rr.Body.String())
	}
}

func TestCreateEventRejectsUnknownCategoryBeforeDB(t *testing.T) {
	h := &H{}
	body := `{"category":"spam","severity":"severe","summary":"Explicit content detected"}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.CreateEvent(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %q", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid category") {
		t.Fatalf("body = %q, want invalid category", rr.Body.String())
	}
}
