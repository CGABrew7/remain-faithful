package core

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrPINRequired = errors.New("Enter the Partner PIN your partner set.")
	ErrNoApps      = errors.New("Choose apps first.")
	ErrConsent     = errors.New("Allow monitoring first.")
	ErrLockout     = errors.New("Turn off App Lockout first.")
)

const (
	MaxGroupMembers = 12
	TourVersion     = 1
)

type LocalState struct {
	OnboardingComplete bool               `json:"onboarding_complete"`
	TourVersion        int                `json:"tour_version"`
	Consent            bool               `json:"consent"`
	Alerts             bool               `json:"alerts"`
	Blocking           bool               `json:"blocking"`
	Lockout            bool               `json:"lockout"`
	Selected           []string           `json:"selected"`
	Ledger             Ledger             `json:"ledger"`
	Queue              []ProtectionNotice `json:"queue"`
	QueueUser          int64              `json:"queue_user"`
	PanicTaps          int                `json:"panic_taps"`
	PanicUser          int64              `json:"panic_user"`
	DeviceID           string             `json:"device_id"`
}

func (s *LocalState) normalize() {
	if s.Selected == nil {
		s.Selected = []string{}
	}
	if s.Queue == nil {
		s.Queue = []ProtectionNotice{}
	}
}

func needsPIN(reducing, pinSet, unlocked bool) error {
	if reducing && pinSet && !unlocked {
		return ErrPINRequired
	}
	return nil
}

func SetAlerts(s *LocalState, on, pinSet, unlocked bool) ([]ProtectionNotice, error) {
	s.normalize()
	if on == s.Alerts {
		return nil, nil
	}
	if on {
		if !s.Consent {
			return nil, ErrConsent
		}
		if len(s.Selected) == 0 {
			return nil, ErrNoApps
		}
		s.Alerts = true
		return nil, nil
	}
	if err := needsPIN(true, pinSet, unlocked); err != nil {
		return nil, err
	}
	s.Alerts = false
	return notices("monitoring_disabled", "Alerts were turned off.")
}

func SetBlocking(s *LocalState, on, pinSet, unlocked bool) ([]ProtectionNotice, error) {
	s.normalize()
	if on == s.Blocking {
		return nil, nil
	}
	if on {
		if !s.Consent {
			return nil, ErrConsent
		}
		if len(s.Selected) == 0 {
			return nil, ErrNoApps
		}
		s.Blocking = true
		return nil, nil
	}
	if s.Lockout {
		return nil, ErrLockout
	}
	if err := needsPIN(true, pinSet, unlocked); err != nil {
		return nil, err
	}
	s.Blocking = false
	return notices("shielding_disabled", "Blocking was turned off.")
}

func SetLockout(s *LocalState, on, pinSet, unlocked bool) ([]ProtectionNotice, error) {
	s.normalize()
	if on == s.Lockout {
		return nil, nil
	}
	if on {
		if !s.Consent {
			return nil, ErrConsent
		}
		if len(s.Selected) == 0 {
			return nil, ErrNoApps
		}
		s.Lockout = true
		s.Blocking = true
		return nil, nil
	}
	if err := needsPIN(true, pinSet, unlocked); err != nil {
		return nil, err
	}
	s.Lockout = false
	return notices("lockout_disabled", "App Lockout was turned off.")
}

func SetSelection(platform string, s *LocalState, keys []string, pinSet, unlocked bool) ([]ProtectionNotice, error) {
	s.normalize()
	next := FilterSelection(platform, keys)
	if len(next) == 0 && len(s.Selected) > 0 && (s.Alerts || s.Blocking || s.Lockout) {
		if err := needsPIN(true, pinSet, unlocked); err != nil {
			return nil, err
		}
		s.Selected = []string{}
		s.Alerts = false
		s.Blocking = false
		s.Lockout = false
		return notices("selection_cleared", "All chosen apps were removed.")
	}
	s.Selected = next
	return nil, nil
}

func RevokeConsent(s *LocalState, pinSet, unlocked bool) ([]ProtectionNotice, error) {
	s.normalize()
	reducing := s.Alerts || s.Blocking || s.Lockout
	if err := needsPIN(reducing, pinSet, unlocked); err != nil {
		return nil, err
	}
	var notes []ProtectionNotice
	if s.Alerts {
		n, _ := notices("monitoring_disabled", "Alerts were turned off.")
		notes = append(notes, n...)
	}
	if s.Blocking {
		n, _ := notices("shielding_disabled", "Blocking was turned off.")
		notes = append(notes, n...)
	}
	if s.Lockout {
		n, _ := notices("lockout_disabled", "App Lockout was turned off.")
		notes = append(notes, n...)
	}
	s.Consent = false
	s.Alerts = false
	s.Blocking = false
	s.Lockout = false
	return notes, nil
}

func GrantConsent(s *LocalState) {
	s.Consent = true
}

func SignOutNotice() (ProtectionNotice, error) {
	n := ProtectionNotice{
		Type:   "signed_out",
		Detail: "Signed out. Monitoring on this computer has stopped.",
	}
	return n, ValidateProtection(n)
}

// StopWatching turns the toggles off without a partner notice.
// Removing this computer uses it so the server can send the one alert.
func StopWatching(s *LocalState) {
	s.Alerts = false
	s.Blocking = false
	s.Lockout = false
}

// ClearMonitoring drops this account's watch settings. Sign-out already sent its notice.
// The device id stays so this computer keeps the same row after the next sign-in.
func ClearMonitoring(s *LocalState) {
	s.Alerts = false
	s.Blocking = false
	s.Lockout = false
	s.Selected = []string{}
	s.Ledger = Ledger{}
	s.Queue = []ProtectionNotice{}
	s.QueueUser = 0
}

func notices(typ, detail string) ([]ProtectionNotice, error) {
	n := ProtectionNotice{Type: typ, Detail: detail}
	if err := ValidateProtection(n); err != nil {
		return nil, err
	}
	return []ProtectionNotice{n}, nil
}

func WrongPINNotice(name string) (ProtectionNotice, error) {
	detail := strings.TrimSpace(name) + " entered a wrong Protection PIN."
	if strings.TrimSpace(name) == "" {
		detail = "A wrong Protection PIN was entered."
	}
	n := ProtectionNotice{Type: "wrong_pin_attempt", Detail: detail}
	if err := ValidateProtection(n); err != nil {
		// A name must not carry a page or an app identity. Fall back to a fixed line.
		n.Detail = "A wrong Protection PIN was entered."
		if err2 := ValidateProtection(n); err2 != nil {
			return ProtectionNotice{}, err2
		}
	}
	return n, nil
}

func PINChangedNotice(replacing bool) (ProtectionNotice, bool, error) {
	if !replacing {
		return ProtectionNotice{}, false, nil
	}
	n := ProtectionNotice{Type: "pin_changed", Detail: "Protection PIN was changed by your accountability partner."}
	return n, true, ValidateProtection(n)
}

func PINRemovedNotice() (ProtectionNotice, error) {
	n := ProtectionNotice{Type: "pin_removed", Detail: "Protection PIN removed by accountability partner."}
	return n, ValidateProtection(n)
}

func ValidateNewAccount(name, email, password string, adult bool) error {
	if !adult {
		return errors.New("Confirm that you are 18 or older and setting this up for yourself.")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("Your name is required.")
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if !PasswordMet(password) {
		return errors.New(PasswordHint(password))
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("Enter an email address.")
	}
	return nil
}

func PasswordMet(password string) bool {
	return len(password) >= 8
}

func PasswordHint(password string) string {
	if password == "" {
		return "At least 8 characters."
	}
	if PasswordMet(password) {
		return "8 or more characters."
	}
	return itoa(len(password)) + " of 8 characters."
}

func ValidatePIN(pin string) error {
	if len(pin) != 4 {
		return errors.New("pin must be exactly 4 digits")
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return errors.New("pin must contain only digits")
		}
	}
	return nil
}

func ValidateGroupName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > 80 {
		return errors.New("That name is too long.")
	}
	return nil
}

func CanInvite(memberCount int) bool {
	return memberCount < MaxGroupMembers
}

func ValidateFeedback(kind, summary, details string) error {
	switch strings.TrimSpace(kind) {
	case "Idea", "Improvement", "Problem":
	default:
		return errors.New("type must be Idea, Improvement, or Problem")
	}
	summary = strings.TrimSpace(summary)
	details = strings.TrimSpace(details)
	if summary == "" || utf8.RuneCountInString(summary) > 120 {
		return errors.New("summary is required (120 characters max)")
	}
	if details == "" || utf8.RuneCountInString(details) > 4000 {
		return errors.New("details are required (4000 characters max)")
	}
	return nil
}

func DesktopSignals(selected, taps int) map[string]any {
	if selected < 0 {
		selected = 0
	}
	if selected > 500 {
		selected = 500
	}
	if taps < 0 {
		taps = 0
	}
	return map[string]any{
		"screen_time_granted": false,
		"monitored_app_count": selected,
		"app_version":         Version,
		"support_taps":        taps,
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
