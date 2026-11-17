package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateCheckIn(t *testing.T) {
	status, note, msg := validateCheckIn("steady", "  Walking today  ")
	if msg != "" || status != "steady" || note != "Walking today" {
		t.Fatalf("got %q %q %q", status, note, msg)
	}
	if _, _, msg := validateCheckIn("watching", "hi"); msg == "" {
		t.Fatal("expected invalid status")
	}
	if _, _, msg := validateCheckIn("slipped", strings.Repeat("a", maxCheckInNote+1)); msg == "" {
		t.Fatal("expected a long note to fail")
	}
	_, note, msg = validateCheckIn("need_support", "line\nbreak")
	if msg != "" || strings.Contains(note, "\n") {
		t.Fatalf("note %q msg %q", note, msg)
	}
}

func TestCheckInPushBody(t *testing.T) {
	got := checkInPushBody("Alex", "steady", "Home tonight")
	if got != "Alex sent a check-in: Steady. Home tonight" {
		t.Fatalf("body = %q", got)
	}
	if got := checkInPushBody("Alex", "slipped", ""); got != "Alex sent a check-in: Had a slip." {
		t.Fatalf("empty note body = %q", got)
	}
}

func TestCreateCheckInStoresAUserNote(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t)}
	rr := httptest.NewRecorder()
	body := `{"status":"steady","note":"Home tonight","shielding_on":true,"session_shield_on":true,"schedule_on":false,"deep_scan_on":false}`
	h.CreateCheckIn(rr, authed(http.MethodPost, "/check-ins", body, 1))
	if rr.Code != http.StatusCreated {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"note":"Home tonight"`) {
		t.Fatalf("body %s", rr.Body.String())
	}
}

func TestCreateCheckInRejectsUnknownStatus(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t)}
	rr := httptest.NewRecorder()
	h.CreateCheckIn(rr, authed(http.MethodPost, "/check-ins", `{"status":"secret"}`, 1))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
}

func TestNormalizeEventSignals(t *testing.T) {
	tier := 2
	gotTier, bucket, msg := normalizeEventSignals(&tier, "High")
	if msg != "" || bucket != "high" || gotTier == nil || *gotTier != 2 {
		t.Fatalf("tier %v bucket %q msg %q", gotTier, bucket, msg)
	}
	if _, _, msg := normalizeEventSignals(nil, "exact"); msg == "" {
		t.Fatal("expected invalid bucket")
	}
	bad := 9
	if _, _, msg := normalizeEventSignals(&bad, ""); msg == "" {
		t.Fatal("expected invalid tier")
	}
	if _, bucket, msg := normalizeEventSignals(nil, ""); msg != "" || bucket != "" {
		t.Fatalf("empty signals: %q %q", bucket, msg)
	}
}

func TestAlertBodyKeepsConfidenceSeparateFromSummary(t *testing.T) {
	got := alertBodyWithConfidence("Ada", "adult_content", "severe", "Explicit content detected", "high")
	if !strings.Contains(got, "Confidence: high.") {
		t.Fatalf("body = %q", got)
	}
	plain := alertBodyWithConfidence("Ada", "adult_content", "severe", "Explicit content detected", "")
	if strings.Contains(plain, "Confidence") {
		t.Fatalf("empty bucket changed the body: %q", plain)
	}
}
