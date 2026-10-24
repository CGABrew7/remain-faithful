package main

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"remain-faithful/desktop/core"
	"remain-faithful/desktop/monitor"
)

func (a *App) Bootstrap() Bootstrap {
	st := a.store.State()
	return Bootstrap{
		SignedIn:           a.store.SignedIn(),
		OnboardingComplete: st.OnboardingComplete,
		ShowTour:           core.ShouldShowTour(st.TourVersion, a.store.SignedIn()),
		Platform:           monitor.Platform(),
		Version:            core.VersionLabel(),
		Consent:            st.Consent,
		HelperReady:        monitor.Ready(),
	}
}

func (a *App) PasswordHint(password string) string {
	return core.PasswordHint(password)
}

func (a *App) Register(name, email, password string, adult bool) error {
	if err := core.ValidateNewAccount(name, email, password, adult); err != nil {
		return err
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	if _, err := a.client.Register(ctx, strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email)), password); err != nil {
		return err
	}
	return a.login(email, password, false)
}

func (a *App) Login(email, password string) error {
	if strings.TrimSpace(email) == "" || password == "" {
		return errors.New("Enter your email and password.")
	}
	return a.login(email, password, true)
}

func (a *App) login(email, password string, returning bool) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	resp, err := a.client.Login(ctx, strings.ToLower(strings.TrimSpace(email)), password)
	if err != nil {
		return err
	}
	if err := a.store.SetSession(resp.Token, resp.User); err != nil {
		return err
	}
	_ = a.store.Update(func(st *core.LocalState) error {
		if returning {
			st.OnboardingComplete = true
		} else {
			st.OnboardingComplete = false
			st.TourVersion = 0
			st.Consent = false
		}
		return nil
	})
	go a.signals()
	go a.deliver(a.store.TakeQueue(resp.User.ID))
	return nil
}

func (a *App) ForgotPassword(email string) (string, error) {
	if err := core.ValidateEmail(email); err != nil {
		return "", err
	}
	email = strings.ToLower(strings.TrimSpace(email))
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	if err := a.client.ForgotPassword(ctx, email); err != nil {
		return "", err
	}
	return core.ForgotPasswordMessage(email), nil
}

func (a *App) ResendVerification() (string, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	status, err := a.client.ResendVerification(ctx)
	if err != nil {
		return "", a.soften(err)
	}
	if status == "already_verified" {
		user := a.store.User()
		user.EmailVerified = true
		_ = a.store.UpdateUser(user)
	}
	return core.ResendMessage(status), nil
}

func (a *App) CompleteOnboarding() error {
	return a.store.Update(func(st *core.LocalState) error {
		st.OnboardingComplete = true
		return nil
	})
}

func (a *App) FinishTour() error {
	return a.store.Update(func(st *core.LocalState) error {
		st.TourVersion = core.TourVersion
		return nil
	})
}

func (a *App) TourSteps(inGroup bool) []core.TourStep {
	return core.TourSteps(inGroup)
}

func (a *App) SignOut(unlocked bool) error {
	st := a.store.State()
	reducing := st.Alerts || st.Blocking || st.Lockout
	pin, err := a.pinRequired(reducing)
	if err != nil {
		return err
	}
	if reducing && pin && !unlocked {
		return core.ErrPINRequired
	}
	note, err := core.SignOutNotice()
	if err != nil {
		return err
	}
	a.deliver([]core.ProtectionNotice{note})
	return a.store.ClearSession()
}

func (a *App) DeleteAccount(confirm string) error {
	if confirm != "DELETE" {
		return errors.New("Type DELETE to confirm.")
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	if err := a.client.DeleteMe(ctx); err != nil {
		return a.soften(err)
	}
	return a.store.ClearSession()
}

func (a *App) UpdateProfile(name, email string) (string, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if name == "" || core.ValidateEmail(email) != nil {
		return "", errors.New("name and email are required")
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	user, err := a.client.UpdateMe(ctx, name, email)
	if err != nil {
		return "", a.soften(err)
	}
	_ = a.store.UpdateUser(user)
	if user.PendingEmail != "" {
		return "Check that inbox to confirm the new email.", nil
	}
	return "Saved.", nil
}

func (a *App) DeviceCaption() string {
	return core.VersionLabel() + ", " + core.DeviceModel(monitor.Platform()) + ", " + osCaption()
}

func (a *App) SendFeedback(kind, summary, details string, includeDevice bool) error {
	if err := core.ValidateFeedback(kind, summary, details); err != nil {
		return err
	}
	fields := map[string]string{
		"type":    strings.TrimSpace(kind),
		"summary": strings.TrimSpace(summary),
		"details": strings.TrimSpace(details),
	}
	if includeDevice {
		fields["app_version"] = core.VersionLabel()
		fields["device_model"] = core.DeviceModel(monitor.Platform())
		fields["os_version"] = osCaption()
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.Feedback(ctx, fields))
}

func osCaption() string {
	switch monitor.Platform() {
	case "mac":
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "sw_vers", "-productVersion").Output()
		if err != nil {
			return "macOS"
		}
		return "macOS " + strings.TrimSpace(string(out))
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}
