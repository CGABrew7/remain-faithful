package core

import "strconv"

const (
	GroupFullMessage = "Groups hold up to 12 people."
	PauseNotice      = "Pausing or turning off monitoring notifies your partners"
)

type Status struct {
	Title  string
	Detail string
	Action string
}

func MonitoringStatus(consent bool, selected int, alerts, blocking, helperOK bool) Status {
	if !helperOK {
		return Status{
			Title:  "Setup needed",
			Detail: "This copy cannot see or close other apps yet.",
		}
	}
	if !consent {
		return Status{
			Title:  "Setup needed",
			Detail: "Monitoring access is off. Alerts and Blocking stay off until you allow it.",
			Action: "Allow monitoring",
		}
	}
	if selected == 0 {
		return Status{
			Title:  "No apps selected",
			Detail: "Nothing is being watched yet. Choose the apps you want help with. You can take your time.",
			Action: "Choose apps",
		}
	}
	if !alerts && !blocking {
		return Status{
			Title:  "Monitoring off",
			Detail: "Apps are chosen. Alerts and Blocking are both off. Turn them on when you are ready.",
			Action: "Turn monitoring on",
		}
	}
	parts := ""
	if alerts {
		parts = "Alerts on"
	} else {
		parts = "Alerts off"
	}
	if blocking {
		parts += ". Blocking on"
	} else {
		parts += ". Blocking off"
	}
	return Status{
		Title:  "Monitoring on",
		Detail: parts + ". You are not doing this alone.",
	}
}

type TourStep struct {
	Tab   string `json:"tab"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func TourSteps(inGroup bool) []TourStep {
	steps := []TourStep{
		{Tab: "home", Title: "Monitoring status", Body: "This shows whether Alerts and Blocking are on."},
		{Tab: "home", Title: "Your clean streak", Body: "Days in a row without a flag. The circles are this week."},
		{Tab: "home", Title: "Recent flags", Body: "Your latest alerts. Open one to see what happened."},
		{Tab: "home", Title: "When it gets hard", Body: "One tap sends an urgent notice to your partners and your group."},
	}
	if inGroup {
		steps = append(steps,
			TourStep{Tab: "group", Title: "Your group", Body: "Everyone's streak and how they're doing. Groups hold up to 12 people."},
			TourStep{Tab: "group", Title: "Invite someone", Body: "Type an email. They join when they accept, or when they sign up with that email."},
		)
	} else {
		steps = append(steps, TourStep{Tab: "group", Title: "Your group", Body: "Start a group here. A name is enough. Groups hold up to 12 people."})
	}
	steps = append(steps,
		TourStep{Tab: "settings", Title: "Your partners", Body: "Add partners and star your primary partner."},
		TourStep{Tab: "settings", Title: "Choose apps", Body: "Pick the apps that tempt you. Blocking stops them from opening. Alerts tell your partner after one minute a day."},
		TourStep{Tab: "settings", Title: "Tell us what you think", Body: "Ideas, fixes, problems. We read every one."},
	)
	return steps
}

func ShouldShowTour(completedVersion int, signedIn bool) bool {
	return signedIn && completedVersion < TourVersion
}

func GroupInviteLink(id int64) string {
	return "https://www.remainfaithful.com/join/" + strconv.FormatInt(id, 10)
}

func PartnerShareText() string {
	return "I would like you as my accountability partner on Remain Faithful. https://www.remainfaithful.com"
}

func PartnerInviteMessage(status string) string {
	switch stringsLower(status) {
	case "pending", "connected", "invited_existing":
		return "Invite sent. Waiting for them to accept. Alerts start after they accept."
	case "invited":
		return "Invite sent. They can accept after they create an account with this email. They don't need an account before you invite them."
	case "accepted":
		return "They're connected. Alerts can start."
	default:
		return "Invite sent. They'll need to accept before alerts start."
	}
}

func GroupInviteMessage(status string) string {
	switch stringsLower(status) {
	case "added":
		return "They're in the group now."
	case "already_member", "conflict":
		return "They're already in this group."
	case "invited":
		return "We emailed them. They join after they accept. If they don't have an account yet, they can accept after they create one with that email."
	default:
		return "We emailed them. They join after they accept."
	}
}

func PINMessage(success bool, reason string) string {
	if success {
		return ""
	}
	switch reason {
	case "no_pin":
		return "No PIN is set for this partnership."
	case "rate_limited":
		return "too many PIN attempts. Try again in 15 minutes"
	case "unavailable":
		return "Couldn't check the PIN. Check your connection and try again."
	default:
		return "Incorrect PIN. Your partner is the one who knows it."
	}
}

func ForgotPasswordMessage(email string) string {
	return "If an account exists for " + email + ", you'll get an email with a link to remainfaithful.com/reset-password. The link works for one hour."
}

func ResendMessage(status string) string {
	if status == "already_verified" {
		return "Your email is already confirmed."
	}
	return "Check your email for the confirmation link."
}

func CategoryLabel(category string) string {
	switch category {
	case "adult_content":
		return "Adult Content"
	case "gambling":
		return "Gambling"
	case "violence":
		return "Violence"
	case "self_harm":
		return "Self-harm"
	case "clean":
		return "Clean"
	case "app_usage":
		return "App Use"
	case "protection", "monitoring_disabled", "shielding_disabled", "lockout_disabled",
		"signed_out", "selection_cleared", "wrong_pin_attempt", "pin_removed",
		"pin_changed", "heartbeat_silence", "family_controls_revoked", "deep_scan_stopped":
		return "Protection"
	default:
		return "Alert"
	}
}

func stringsLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}
