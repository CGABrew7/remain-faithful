package main

import (
	"errors"
	"time"

	"remain-faithful/desktop/core"
	"remain-faithful/desktop/monitor"
)

func (a *App) Home() (HomeView, error) {
	if !a.store.SignedIn() {
		return HomeView{}, errors.New("Not signed in. Please log in again.")
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	user, err := a.client.Me(ctx)
	if err != nil {
		return HomeView{}, a.soften(err)
	}
	_ = a.store.UpdateUser(user)
	st := a.store.State()
	status := core.MonitoringStatus(st.Consent, len(st.Selected), st.Alerts, st.Blocking, monitor.Ready())
	view := HomeView{
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		StatusTitle:   status.Title,
		StatusDetail:  status.Detail,
		StatusAction:  status.Action,
		AlertsOn:      st.Alerts,
		BlockingOn:    st.Blocking,
		LockoutOn:     st.Lockout,
		Consent:       st.Consent,
		SelectedCount: len(st.Selected),
		Week:          []string{},
		MyEvents:      []Row{},
		PartnerAlerts: []Row{},
	}
	if stats, err := a.client.Stats(ctx); err == nil {
		view.Streak = stats.StreakDays
		view.Best = stats.BestStreak
		view.Week = core.WeekMarks(stats)
	}
	now := time.Now()
	if events, err := a.client.ListEvents(ctx); err == nil {
		for _, event := range events {
			view.MyEvents = append(view.MyEvents, Row{
				ID:      event.ID,
				Label:   core.CategoryLabel(event.Category),
				Summary: event.Summary,
				When:    core.FormatWhen(event.Timestamp, now),
			})
		}
	}
	if alerts, err := a.client.ListAlerts(ctx); err == nil {
		for _, alert := range alerts {
			if !alert.Seen {
				view.Unread++
			}
			summary := alert.Event.Summary
			category := alert.Event.Category
			if summary == "" {
				summary = core.CategoryLabel(category)
			}
			when := alert.Event.Timestamp
			if when == "" {
				when = alert.CreatedAt
			}
			view.PartnerAlerts = append(view.PartnerAlerts, Row{
				ID:      alert.ID,
				Label:   core.CategoryLabel(category),
				Summary: summary,
				When:    core.FormatWhen(when, now),
			})
		}
		_ = a.client.MarkAlertsSeen(ctx)
	}
	if groups, err := a.client.ListGroups(ctx); err == nil {
		view.InGroup = len(groups) > 0
	}
	if pending, err := a.client.PendingPartnerInvites(ctx); err == nil {
		view.PendingPartners = len(pending)
	}
	if pending, err := a.client.PendingGroupInvites(ctx); err == nil {
		view.PendingGroups = len(pending)
	}
	return view, nil
}

func (a *App) MarkDiscussed(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.MarkDiscussed(ctx, id))
}

func (a *App) Panic() error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	if err := a.client.Panic(ctx); err != nil {
		return a.soften(err)
	}
	taps, _ := a.store.PanicTap(a.store.User().ID)
	st := a.store.State()
	ctx2, cancel2 := a.ctxTimeout()
	defer cancel2()
	_ = a.client.SendSignals(ctx2, len(st.Selected), taps)
	return nil
}

func (a *App) PanicPartner() (string, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	rows, err := a.client.ListRelationships(ctx)
	if err != nil {
		return "", a.soften(err)
	}
	me := a.store.User().ID
	var fallback string
	for _, row := range rows {
		if row.Status != "accepted" {
			continue
		}
		if fallback == "" {
			fallback = row.Partner.Name
		}
		if row.IsPrimary && row.UserID == me {
			return row.Partner.Name, nil
		}
	}
	return fallback, nil
}

func (a *App) VerifyPIN(pin string) (PINView, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	result, err := a.client.VerifyPIN(ctx, pin)
	if err != nil {
		if core.IsNetwork(err) {
			return PINView{Message: core.PINMessage(false, "unavailable")}, nil
		}
		return PINView{}, a.soften(err)
	}
	if !result.Success && result.Reason == "wrong_pin" {
		if note, nerr := core.WrongPINNotice(a.store.User().Name); nerr == nil {
			a.deliver([]core.ProtectionNotice{note})
		}
	}
	return PINView{Success: result.Success, Message: result.Message}, nil
}

func (a *App) GrantConsent() error {
	err := a.store.Update(func(st *core.LocalState) error {
		core.GrantConsent(st)
		return nil
	})
	if err != nil {
		return err
	}
	if a.watching() {
		_ = a.reclaimDevice()
	}
	a.beat()
	return nil
}

func (a *App) RevokeConsent(unlocked bool) error {
	st := a.store.State()
	reducing := st.Alerts || st.Blocking || st.Lockout
	pin, err := a.pinRequired(reducing)
	if err != nil {
		return err
	}
	var notes []core.ProtectionNotice
	err = a.store.Update(func(state *core.LocalState) error {
		var e error
		notes, e = core.RevokeConsent(state, pin, unlocked)
		return e
	})
	return a.afterWatch(notes, err)
}

func (a *App) SetAlerts(on, unlocked bool) error {
	pin, err := a.pinRequired(!on)
	if err != nil {
		return err
	}
	var notes []core.ProtectionNotice
	err = a.store.Update(func(state *core.LocalState) error {
		var e error
		notes, e = core.SetAlerts(state, on, pin, unlocked)
		return e
	})
	return a.afterWatch(notes, err)
}

func (a *App) SetBlocking(on, unlocked bool) error {
	pin, err := a.pinRequired(!on)
	if err != nil {
		return err
	}
	var notes []core.ProtectionNotice
	err = a.store.Update(func(state *core.LocalState) error {
		var e error
		notes, e = core.SetBlocking(state, on, pin, unlocked)
		return e
	})
	return a.afterWatch(notes, err)
}

func (a *App) SetLockout(on, unlocked bool) error {
	pin, err := a.pinRequired(!on)
	if err != nil {
		return err
	}
	var notes []core.ProtectionNotice
	err = a.store.Update(func(state *core.LocalState) error {
		var e error
		notes, e = core.SetLockout(state, on, pin, unlocked)
		return e
	})
	return a.afterWatch(notes, err)
}

func (a *App) SetSelection(keys []string, unlocked bool) error {
	st := a.store.State()
	clearing := len(core.FilterSelection(monitor.Platform(), keys)) == 0 && len(st.Selected) > 0 && (st.Alerts || st.Blocking || st.Lockout)
	pin, err := a.pinRequired(clearing)
	if err != nil {
		return err
	}
	var notes []core.ProtectionNotice
	err = a.store.Update(func(state *core.LocalState) error {
		var e error
		notes, e = core.SetSelection(monitor.Platform(), state, keys, pin, unlocked)
		return e
	})
	return a.afterWatch(notes, err)
}

func (a *App) Catalog() []core.AppChoice {
	st := a.store.State()
	chosen := map[string]bool{}
	for _, key := range st.Selected {
		chosen[key] = true
	}
	items := core.Catalog(monitor.Platform())
	if items == nil {
		items = []core.AppChoice{}
	}
	for i := range items {
		items[i].Selected = chosen[items[i].Key]
	}
	return items
}
