package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rfauth "remain-faithful/backend/internal/auth"
)

func TestIsAllowedProtectionAlertType(t *testing.T) {
	// Types actually POSTed by the iOS app / broadcast extension.
	iosTypes := []string{
		"monitoring_disabled",
		"shielding_disabled",
		"lockout_disabled",
		"deep_scan_stopped",
		"wrong_pin_attempt",
		"pin_removed",
		"pin_changed",
		"family_controls_revoked",
	}
	// Types already named in backend code (switch + heartbeat sweep).
	backendTypes := []string{
		"lockout_broadcast_stopped",
		"pin_wrong_attempt",
		"heartbeat_silence",
	}

	for _, typ := range append(iosTypes, backendTypes...) {
		if !isAllowedProtectionAlertType(typ) {
			t.Errorf("expected allowed type %q", typ)
		}
	}

	for _, typ := range []string{"", "not_a_real_type", "monitoring_off", "spam"} {
		if isAllowedProtectionAlertType(typ) {
			t.Errorf("expected unknown type %q to be rejected", typ)
		}
	}
}

func TestSendProtectionAlertRejectsUnknownType(t *testing.T) {
	h := &H{}
	req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(`{"type":"partner_spam"}`))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.SendProtectionAlert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "unknown protection alert type") {
		t.Fatalf("body = %q, want unknown-type error", rr.Body.String())
	}
}

func TestSendProtectionAlertRequiresType(t *testing.T) {
	h := &H{}
	req := httptest.NewRequest(http.MethodPost, "/protection/alerts", strings.NewReader(`{"type":""}`))
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.SendProtectionAlert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}
