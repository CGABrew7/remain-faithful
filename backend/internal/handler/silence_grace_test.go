package handler

import (
	"context"
	"testing"
	"time"
)

func TestSilenceGraceElapsed(t *testing.T) {
	boot := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	grace := 30 * time.Minute
	if silenceGraceElapsed(boot, boot.Add(29*time.Minute), grace) {
		t.Fatal("inside the restart grace")
	}
	if !silenceGraceElapsed(boot, boot.Add(30*time.Minute), grace) {
		t.Fatal("grace should be finished at the threshold")
	}
	if !silenceGraceElapsed(time.Time{}, boot, grace) {
		t.Fatal("an unset boot time is not a grace window")
	}
}

func TestRestartGraceHoldsDeviceSilenceAlerts(t *testing.T) {
	db, h := newDeviceHarness(t)
	postHeartbeat(t, h, `{"screen":"monitoring","device_id":"phone-1","device_kind":"iphone"}`)
	db.ageDevice(1, "phone-1", time.Now().Add(-2*time.Hour))

	h.bootAt = time.Now()
	h.sweepDeviceHeartbeatSilence(context.Background(), 30)
	if got := db.alertSummaries(); len(got) != 0 {
		t.Fatalf("grace sent %#v", got)
	}

	h.bootAt = time.Now().Add(-31 * time.Minute)
	h.sweepDeviceHeartbeatSilence(context.Background(), 30)
	got := db.alertSummaries()
	if len(got) != 1 || got[0] != "Ada's iPhone has stopped sending heartbeats" {
		t.Fatalf("after grace %#v", got)
	}
}

func TestRestartGraceSkipsAccountSilenceQuery(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t), bootAt: time.Now()}
	h.sweepHeartbeatSilence(context.Background(), 30)
	if db.sawSQL("last_heartbeat_at <") {
		t.Fatal("queried account silence during the restart grace")
	}

	h.bootAt = time.Now().Add(-time.Hour)
	h.sweepHeartbeatSilence(context.Background(), 30)
	if !db.sawSQL("last_heartbeat_at <") {
		t.Fatal("expected the silence query after the grace")
	}
}
