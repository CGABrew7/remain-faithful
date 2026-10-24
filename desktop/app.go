package main

import (
	"context"
	"os"
	"sync"
	"time"

	"remain-faithful/desktop/core"
	"remain-faithful/desktop/monitor"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	store  *core.Store
	client *core.Client
	stop   chan struct{}
	once   sync.Once
}

func NewApp() *App {
	store, err := core.OpenStore("")
	if err != nil {
		store, _ = core.OpenStore(os.TempDir())
	}
	app := &App{store: store, stop: make(chan struct{})}
	app.client = core.NewClient(os.Getenv("API_BASE_URL"), func() string { return app.store.Token() })
	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.refreshIfNeeded()
	go a.deliver(a.store.TakeQueue(a.store.User().ID))
	go a.signals()
	go a.loop()
}

func (a *App) shutdown(ctx context.Context) {
	a.once.Do(func() { close(a.stop) })
}

func (a *App) loop() {
	sample := time.NewTicker(2 * time.Second)
	beat := time.NewTicker(60 * time.Second)
	poll := time.NewTicker(30 * time.Second)
	defer sample.Stop()
	defer beat.Stop()
	defer poll.Stop()
	a.beat()
	for {
		select {
		case <-a.stop:
			return
		case <-sample.C:
			a.sample()
		case <-beat.C:
			a.beat()
		case <-poll.C:
			a.noteAlerts()
		}
	}
}

func (a *App) sample() {
	if !a.store.SignedIn() {
		return
	}
	st := a.store.State()
	if core.HeartbeatScreen(st.Consent, st.Alerts, st.Blocking, st.Lockout) == "" {
		return
	}
	front, err := monitor.Sample()
	if err != nil {
		return
	}
	var event *core.OutboundEvent
	var blockKey string
	_ = a.store.Update(func(state *core.LocalState) error {
		next, out := core.Evaluate(time.Now(), monitor.Platform(), core.WatchConfig{
			Consent:  state.Consent,
			Alerts:   state.Alerts,
			Blocking: state.Blocking,
			Selected: state.Selected,
		}, state.Ledger, front.Key, 2*time.Second)
		state.Ledger = next
		event = out.Event
		blockKey = out.BlockKey
		return nil
	})
	if blockKey != "" {
		_ = monitor.Quit(blockKey)
		if a.ctx != nil {
			display := "A chosen app"
			if choice, ok := core.Lookup(monitor.Platform(), blockKey); ok {
				display = choice.Display
			}
			runtime.EventsEmit(a.ctx, "shield", display)
		}
	}
	if event == nil {
		return
	}
	ctx, cancel := a.ctxTimeout()
	err = a.client.CreateEvent(ctx, *event)
	cancel()
	if err == nil || !core.IsNetwork(err) {
		return
	}
	_ = a.store.Update(func(state *core.LocalState) error {
		if event.Summary == core.SummaryApp {
			state.Ledger.AppSent = false
			if state.Ledger.AppSeconds >= 60 {
				state.Ledger.AppSeconds = 59
			}
		} else {
			state.Ledger.BrowserSent = false
			if state.Ledger.BrowserSeconds >= 60 {
				state.Ledger.BrowserSeconds = 59
			}
		}
		return nil
	})
}

func (a *App) beat() {
	if !a.store.SignedIn() {
		return
	}
	id, kind, ok := a.deviceIdentity()
	if !ok {
		return
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	watching := a.watching()
	var err error
	if watching {
		err = a.client.HeartbeatMonitoring(ctx, id, kind)
	} else {
		err = a.client.ReportDevice(ctx, id, kind, false)
	}
	if core.IsDeviceRemoved(err) && watching {
		a.parkRemoved()
	}
}

func (a *App) signals() {
	if !a.store.SignedIn() {
		return
	}
	st := a.store.State()
	taps := 0
	if st.PanicUser == a.store.User().ID {
		taps = st.PanicTaps
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	_ = a.client.SendSignals(ctx, len(st.Selected), taps)
}

func (a *App) noteAlerts() {
	if !a.store.SignedIn() || a.ctx == nil {
		return
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	n, err := a.client.UnreadCount(ctx)
	if err != nil {
		return
	}
	runtime.EventsEmit(a.ctx, "alerts", n)
}

func (a *App) deliver(notes []core.ProtectionNotice) {
	if len(notes) == 0 || !a.store.SignedIn() {
		if len(notes) > 0 {
			_ = a.store.Enqueue(a.store.User().ID, notes)
		}
		return
	}
	userID := a.store.User().ID
	deviceID, deviceKind, _ := a.deviceIdentity()
	var keep []core.ProtectionNotice
	for _, n := range notes {
		ctx, cancel := a.ctxTimeout()
		err := a.client.SendProtection(ctx, n, deviceID, deviceKind)
		cancel()
		if err == nil {
			continue
		}
		if !core.IsNetwork(err) && n.Type == "signed_out" {
			fallback := core.ProtectionNotice{Type: "monitoring_disabled", Detail: n.Detail}
			ctx, cancel = a.ctxTimeout()
			err = a.client.SendProtection(ctx, fallback, deviceID, deviceKind)
			cancel()
		}
		if err != nil && core.IsNetwork(err) {
			keep = append(keep, n)
		}
	}
	_ = a.store.ReplaceQueue(userID, keep)
}

func (a *App) refreshIfNeeded() {
	tok := a.store.Token()
	if tok == "" {
		return
	}
	exp, ok := core.TokenExpiry(tok)
	if ok && time.Until(exp) > time.Hour {
		return
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	next, err := a.client.Refresh(ctx)
	if err != nil {
		if ok && time.Now().After(exp) {
			_ = a.store.ClearSession()
			if a.ctx != nil {
				runtime.EventsEmit(a.ctx, "signed-out", "Not signed in. Please log in again.")
			}
		}
		return
	}
	if next != "" {
		_ = a.store.UpdateToken(next)
	}
}

func (a *App) ctxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

func (a *App) soften(err error) error {
	if err == nil {
		return nil
	}
	if api, ok := err.(*core.APIError); ok && api.Status == 401 {
		_ = a.store.ClearSession()
	}
	return err
}

func (a *App) pinRequired(reducing bool) (bool, error) {
	if !reducing {
		return false, nil
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.client.PINStatus(ctx)
}

func (a *App) afterWatch(notes []core.ProtectionNotice, err error) error {
	if err != nil {
		return err
	}
	a.deliver(notes)
	if a.watching() {
		_ = a.reclaimDevice()
	}
	a.beat()
	a.signals()
	return nil
}
