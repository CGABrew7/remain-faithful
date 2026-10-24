package main

import "remain-faithful/desktop/core"

type Bootstrap struct {
	SignedIn           bool   `json:"signedIn"`
	OnboardingComplete bool   `json:"onboardingComplete"`
	ShowTour           bool   `json:"showTour"`
	Platform           string `json:"platform"`
	Version            string `json:"version"`
	Consent            bool   `json:"consent"`
	HelperReady        bool   `json:"helperReady"`
}

type HomeView struct {
	Name            string   `json:"name"`
	EmailVerified   bool     `json:"emailVerified"`
	StatusTitle     string   `json:"statusTitle"`
	StatusDetail    string   `json:"statusDetail"`
	StatusAction    string   `json:"statusAction"`
	Streak          int      `json:"streak"`
	Best            int      `json:"best"`
	Week            []string `json:"week"`
	AlertsOn        bool     `json:"alertsOn"`
	BlockingOn      bool     `json:"blockingOn"`
	LockoutOn       bool     `json:"lockoutOn"`
	Consent         bool     `json:"consent"`
	SelectedCount   int      `json:"selectedCount"`
	Unread          int      `json:"unread"`
	InGroup         bool     `json:"inGroup"`
	MyEvents        []Row    `json:"myEvents"`
	PartnerAlerts   []Row    `json:"partnerAlerts"`
	PendingPartners int      `json:"pendingPartners"`
	PendingGroups   int      `json:"pendingGroups"`
}

type Row struct {
	ID      int64  `json:"id"`
	Label   string `json:"label"`
	Summary string `json:"summary"`
	When    string `json:"when"`
}

type SettingsView struct {
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	EmailVerified   bool             `json:"emailVerified"`
	Platform        string           `json:"platform"`
	Version         string           `json:"version"`
	DeviceCaption   string           `json:"deviceCaption"`
	Consent         bool             `json:"consent"`
	AlertsOn        bool             `json:"alertsOn"`
	BlockingOn      bool             `json:"blockingOn"`
	LockoutOn       bool             `json:"lockoutOn"`
	PinSet          bool             `json:"pinSet"`
	StatusTitle     string           `json:"statusTitle"`
	StatusDetail    string           `json:"statusDetail"`
	PauseNotice     string           `json:"pauseNotice"`
	Apps            []core.AppChoice `json:"apps"`
	Partners        []PartnerView    `json:"partners"`
	PendingPartners []InviteRow      `json:"pendingPartners"`
	PendingGroups   []InviteRow      `json:"pendingGroups"`
}

type PartnerView struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Primary   bool   `json:"primary"`
	PinSet    bool   `json:"pinSet"`
	CanSetPIN bool   `json:"canSetPin"`
	CanStar   bool   `json:"canStar"`
}

type InviteRow struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type GroupCard struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GroupDetail struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	MemberCount int          `json:"memberCount"`
	CanInvite   bool         `json:"canInvite"`
	FullMessage string       `json:"fullMessage"`
	Link        string       `json:"link"`
	Members     []MemberView `json:"members"`
}

type MemberView struct {
	UserID int64  `json:"userId"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Streak int    `json:"streak"`
	Flags  int    `json:"flags"`
	Self   bool   `json:"self"`
}

type PINView struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
