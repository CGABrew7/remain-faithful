package main

import (
	"errors"
	"time"

	"remain-faithful/desktop/core"
	"remain-faithful/desktop/monitor"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) watching() bool {
	st := a.store.State()
	return core.HeartbeatScreen(st.Consent, st.Alerts, st.Blocking, st.Lockout) != ""
}

func (a *App) deviceIdentity() (string, string, bool) {
	kind, ok := core.DeviceKind(monitor.Platform())
	if !ok {
		return "", "", false
	}
	return a.store.DeviceID(), kind, true
}

func (a *App) reclaimDevice() error {
	id, kind, ok := a.deviceIdentity()
	if !ok || !a.store.SignedIn() {
		return nil
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	err := a.client.ReclaimDevice(ctx, id, kind)
	if err == nil || core.IsNotFound(err) {
		return nil
	}
	return err
}

func (a *App) parkRemoved() {
	_ = a.store.Update(func(st *core.LocalState) error {
		core.StopWatching(st)
		return nil
	})
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "device-removed")
	}
}

func (a *App) deviceRows() ([]DeviceRow, bool) {
	if !a.store.SignedIn() {
		return []DeviceRow{}, false
	}
	id, _, known := a.deviceIdentity()
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	remote, err := a.client.ListDevices(ctx)
	if err != nil {
		return []DeviceRow{}, true
	}
	now := time.Now()
	rows := make([]DeviceRow, 0, len(remote))
	for _, d := range remote {
		label, ok := core.DeviceLabel(d.DeviceKind)
		if !ok || !core.ValidDeviceID(d.DeviceID) {
			continue
		}
		rows = append(rows, DeviceRow{
			DeviceID:    d.DeviceID,
			Label:       label,
			Monitoring:  d.Monitoring,
			LastCheckIn: core.FormatCheckIn(d.LastCheckIn, now),
			Current:     known && d.DeviceID == id,
		})
	}
	return rows, false
}

// RemoveDevice drops a device from the account. Partners hear the same sentence
// as when activity monitoring is turned off. The Partner PIN is required when one is set.
func (a *App) RemoveDevice(deviceID string, unlocked bool) error {
	if !a.store.SignedIn() {
		return errors.New("Not signed in. Please log in again.")
	}
	if !core.ValidDeviceID(deviceID) {
		return errors.New("device id is not recognized")
	}
	pin, err := a.pinRequired(true)
	if err != nil {
		return err
	}
	if pin && !unlocked {
		return core.ErrPINRequired
	}
	mine, _, ok := a.deviceIdentity()
	current := ok && deviceID == mine
	wasWatching := current && a.watching()
	if current {
		_ = a.store.Update(func(st *core.LocalState) error {
			core.StopWatching(st)
			return nil
		})
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	err = a.client.RemoveDevice(ctx, deviceID)
	if err == nil {
		return nil
	}
	if wasWatching {
		note := core.ProtectionNotice{Type: "monitoring_disabled", Detail: "Alerts were turned off."}
		if core.ValidateProtection(note) == nil {
			a.deliver([]core.ProtectionNotice{note})
		}
	}
	if core.IsNotFound(err) {
		return errors.New("Device list is not available yet.")
	}
	return a.soften(err)
}
