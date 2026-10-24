package core

import (
	"errors"
	"testing"
)

func TestAccountRules(t *testing.T) {
	if err := ValidateNewAccount("Ada", "ada@example.com", "short", true); err == nil {
		t.Fatal("short password")
	}
	if err := ValidateNewAccount("Ada", "ada@example.com", "longenough", false); err == nil {
		t.Fatal("adult confirmation is required")
	}
	if err := ValidateNewAccount("Ada", "ada@example.com", "longenough", true); err != nil {
		t.Fatal(err)
	}
	if PasswordHint("") != "At least 8 characters." {
		t.Fatal(PasswordHint(""))
	}
	if err := ValidatePIN("123"); err == nil {
		t.Fatal("pin length")
	}
	if err := ValidatePIN("12a4"); err == nil {
		t.Fatal("pin digits")
	}
	if err := ValidatePIN("1234"); err != nil {
		t.Fatal(err)
	}
	if CanInvite(12) || !CanInvite(11) {
		t.Fatal("group cap")
	}
	if GroupInviteLink(4) != "https://www.remainfaithful.com/join/4" {
		t.Fatal(GroupInviteLink(4))
	}
}

func TestReductionsRequirePINAndTellPartners(t *testing.T) {
	s := &LocalState{Consent: true, Alerts: true, Blocking: true, Selected: []string{"chrome.exe"}}
	if _, err := SetAlerts(s, false, true, false); !errors.Is(err, ErrPINRequired) {
		t.Fatal(err)
	}
	if !s.Alerts {
		t.Fatal("alerts changed without the PIN")
	}
	notes, err := SetAlerts(s, false, true, true)
	if err != nil || len(notes) != 1 || notes[0].Type != "monitoring_disabled" {
		t.Fatalf("%v %+v", err, notes)
	}
	if _, err := SetBlocking(s, false, false, false); err != nil || s.Blocking {
		t.Fatal("blocking should turn off when no PIN is set")
	}
}

func TestLockoutKeepsBlockingOn(t *testing.T) {
	s := &LocalState{Consent: true, Selected: []string{"firefox.exe"}}
	if _, err := SetLockout(s, true, false, false); err != nil {
		t.Fatal(err)
	}
	if !s.Lockout || !s.Blocking {
		t.Fatal("lockout turns blocking on")
	}
	if _, err := SetBlocking(s, false, true, true); !errors.Is(err, ErrLockout) {
		t.Fatal(err)
	}
	notes, err := SetLockout(s, false, false, false)
	if err != nil || notes[0].Type != "lockout_disabled" || !s.Blocking {
		t.Fatalf("blocking stays after lockout off: %+v %v blocking=%v", notes, err, s.Blocking)
	}
}

func TestClearingSelectionNotifiesOnce(t *testing.T) {
	s := &LocalState{Consent: true, Alerts: true, Blocking: true, Lockout: true, Selected: []string{"chrome.exe"}}
	notes, err := SetSelection("windows", s, nil, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 || notes[0].Type != "selection_cleared" {
		t.Fatalf("%+v", notes)
	}
	if s.Alerts || s.Blocking || s.Lockout || len(s.Selected) != 0 {
		t.Fatal("selection clear stops watching")
	}
	for _, n := range notes {
		if n.Type == "deep_scan_stopped" || n.Type == "family_controls_revoked" {
			t.Fatal(n.Type)
		}
	}
}

func TestSignOutNoticeAndLocalClear(t *testing.T) {
	n, err := SignOutNotice()
	if err != nil || n.Type != "signed_out" {
		t.Fatal(err)
	}
	s := &LocalState{Alerts: true, Blocking: true, Selected: []string{"chrome.exe"}, Queue: []ProtectionNotice{n}}
	ClearMonitoring(s)
	if s.Alerts || s.Blocking || len(s.Selected) != 0 || len(s.Queue) != 0 {
		t.Fatal("sign-out clears this computer")
	}
}

func TestFeedbackAndSignals(t *testing.T) {
	if err := ValidateFeedback("Idea", "", "details"); err == nil {
		t.Fatal("summary required")
	}
	if err := ValidateFeedback("Idea", "A small idea", "More detail"); err != nil {
		t.Fatal(err)
	}
	sig := DesktopSignals(3, 2)
	if sig["screen_time_granted"] != false || sig["monitored_app_count"] != 3 || sig["app_version"] != Version {
		t.Fatal(sig)
	}
}
