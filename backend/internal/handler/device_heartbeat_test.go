package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func newDeviceHarness(t *testing.T) (*memDB, *H) {
	t.Helper()
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.rels = append(db.rels, memRel{id: 1, userID: 1, partnerID: 2, status: "accepted"})
	return db, &H{DB: db.open(t), SyncPush: true, APNS: scriptedAPNS{}}
}

func postHeartbeat(t *testing.T, h *H, body string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	h.Heartbeat(rr, authed(http.MethodPost, "/heartbeat", body, 1))
	return rr
}

func TestDeviceLabelsAreTypesNotAppNames(t *testing.T) {
	cases := map[string]string{
		"iphone":  "iPhone",
		"mac":     "Mac",
		"windows": "Windows PC",
		"android": "Android phone",
	}
	for kind, label := range cases {
		got, ok := deviceKindLabel(kind)
		if !ok || got != label {
			t.Fatalf("%s label = %q ok=%v", kind, got, ok)
		}
	}
	for _, kind := range []string{"", "Safari", "chrome.exe", "com.apple.mobilesafari", "Windows PC", "phone"} {
		if _, ok := deviceKindLabel(kind); ok {
			t.Fatalf("accepted kind %q", kind)
		}
	}
	if !validDeviceID("abc_DEF-012") || validDeviceID("chrome.exe") || validDeviceID("com.apple.mobilesafari") || validDeviceID(strings.Repeat("a", 65)) {
		t.Fatal("device id charset")
	}
	if protectionBody("Alex", "monitoring_disabled", "") != "Alex turned off activity monitoring" {
		t.Fatal("legacy monitoring sentence changed")
	}
	if got := protectionBody("Alex", "monitoring_disabled", "Windows PC"); got != "Alex turned off activity monitoring on their Windows PC" {
		t.Fatalf("stop sentence = %q", got)
	}
	if got := protectionBody("Alex", "signed_out", "Mac"); got != "Alex signed out on their Mac" {
		t.Fatalf("sign-out sentence = %q", got)
	}
	if got := deviceSilenceBody("Alex", "iPhone"); got != "Alex's iPhone has stopped sending heartbeats" {
		t.Fatalf("silence = %q", got)
	}
	if got := deviceSilenceBody("Alex", "Android phone"); got != "Alex's Android phone has stopped sending heartbeats" {
		t.Fatalf("android silence = %q", got)
	}
	if protectionAlertBody("Alex", "heartbeat_silence", "Safari") != "Alex's device has stopped sending heartbeats" {
		t.Fatal("legacy silence sentence changed")
	}
}

func TestLegacyHeartbeatLeavesDevicesAlone(t *testing.T) {
	db, h := newDeviceHarness(t)
	rr := postHeartbeat(t, h, `{"screen":"monitoring"}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	u, _ := db.user(1)
	if !u.heartbeat || u.screen != "monitoring" {
		t.Fatalf("account heartbeat = %+v", u)
	}
	if db.sawSQL("device_heartbeats") || db.sawSQL("device_removals") {
		t.Fatal("legacy heartbeat wrote a device row")
	}
}

func TestDeviceHeartbeatDoesNotTouchAccountTimestamp(t *testing.T) {
	db, h := newDeviceHarness(t)
	rr := postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	u, _ := db.user(1)
	if u.heartbeat {
		t.Fatal("a computer heartbeat refreshed the account timestamp")
	}
	if db.sawSQL("last_heartbeat_at") {
		t.Fatal("device heartbeat updated users.last_heartbeat_at")
	}
	items := listDevices(t, h)
	if len(items) != 1 || items[0].Label != "Windows PC" || !items[0].Monitoring || items[0].LastCheckIn == "" {
		t.Fatalf("list = %+v", items)
	}
	if items[0].DeviceID != "pc-1" {
		t.Fatalf("owner list id = %q", items[0].DeviceID)
	}
}

func TestDeviceHeartbeatRejectsAppNames(t *testing.T) {
	db, h := newDeviceHarness(t)
	cases := []string{
		`{"screen":"monitoring","device_id":"chrome.exe","device_kind":"windows"}`,
		`{"screen":"monitoring","device_id":"pc-1","device_kind":"chrome.exe"}`,
		`{"screen":"monitoring","device_id":"pc-1","device_kind":"Safari"}`,
		`{"screen":"monitoring","device_id":"pc-1"}`,
		`{"screen":"monitoring","device_kind":"windows"}`,
		`{"screen":"monitoring","device_id":"com.apple.mobilesafari","device_kind":"iphone"}`,
	}
	for _, body := range cases {
		rr := postHeartbeat(t, h, body)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s status %d", body, rr.Code)
		}
	}
	u, _ := db.user(1)
	if u.heartbeat || db.sawSQL("last_heartbeat_at") || db.sawSQL("INSERT INTO device_heartbeats") {
		t.Fatal("rejected heartbeat still wrote state")
	}
}

func TestDeviceListAndRemovalNotifiesLikeMonitoringOff(t *testing.T) {
	db, h := newDeviceHarness(t)
	if rr := postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`); rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}
	rr := httptest.NewRecorder()
	h.ReportDevice(rr, authed(http.MethodPost, "/users/me/devices", `{"device_id":"phone-1","device_kind":"iphone","monitoring":false}`, 1))
	if rr.Code != http.StatusOK {
		t.Fatalf("report %d %s", rr.Code, rr.Body.String())
	}
	items := listDevices(t, h)
	if len(items) != 2 {
		t.Fatalf("list len %d", len(items))
	}
	var phone deviceListItem
	for _, item := range items {
		if item.DeviceKind == "iphone" {
			phone = item
		}
		if strings.Contains(item.Label, ".exe") || strings.Contains(item.Label, "Safari") {
			t.Fatalf("label %q", item.Label)
		}
	}
	if phone.Label != "iPhone" || phone.Monitoring || phone.LastCheckIn != "" {
		t.Fatalf("phone row = %+v", phone)
	}

	req := authed(http.MethodDelete, "/users/me/devices/pc-1", "", 1)
	req = mux.SetURLVars(req, map[string]string{"deviceID": "pc-1"})
	rr = httptest.NewRecorder()
	h.RemoveDevice(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("remove %d %s", rr.Code, rr.Body.String())
	}
	summaries := db.alertSummaries()
	if len(summaries) != 1 || summaries[0] != "Ada turned off activity monitoring on their Windows PC" {
		t.Fatalf("alerts = %#v", summaries)
	}
	if strings.Contains(summaries[0], "pc-1") || strings.Contains(summaries[0], ".exe") {
		t.Fatalf("alert leaked identity %q", summaries[0])
	}
	left := listDevices(t, h)
	if len(left) != 1 || left[0].Label != "iPhone" {
		t.Fatalf("after remove = %+v", left)
	}

	req = authed(http.MethodDelete, "/users/me/devices/pc-1", "", 1)
	req = mux.SetURLVars(req, map[string]string{"deviceID": "pc-1"})
	rr = httptest.NewRecorder()
	h.RemoveDevice(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("second remove %d", rr.Code)
	}
	if got := db.alertSummaries(); len(got) != 1 {
		t.Fatalf("second remove alerted again: %#v", got)
	}
}

func TestRemovedDeviceCannotHeartbeatUntilReclaim(t *testing.T) {
	db, h := newDeviceHarness(t)
	postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`)
	req := authed(http.MethodDelete, "/users/me/devices/pc-1", "", 1)
	req = mux.SetURLVars(req, map[string]string{"deviceID": "pc-1"})
	rr := httptest.NewRecorder()
	h.RemoveDevice(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	rr = postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`)
	if rr.Code != http.StatusConflict || !strings.Contains(rr.Body.String(), `"code":"device_removed"`) {
		t.Fatalf("removed heartbeat %d %s", rr.Code, rr.Body.String())
	}
	u, _ := db.user(1)
	if u.heartbeat || len(listDevices(t, h)) != 0 {
		t.Fatal("removed device heartbeated back onto the account")
	}

	rr = httptest.NewRecorder()
	h.ReclaimDevice(rr, authed(http.MethodPost, "/users/me/devices/reclaim", `{"device_id":"pc-1","device_kind":"windows"}`, 1))
	if rr.Code != http.StatusOK {
		t.Fatalf("reclaim %d %s", rr.Code, rr.Body.String())
	}
	rr = postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("heartbeat after reclaim %d %s", rr.Code, rr.Body.String())
	}
	items := listDevices(t, h)
	if len(items) != 1 || items[0].Label != "Windows PC" || !items[0].Monitoring {
		t.Fatalf("reclaimed list = %+v", items)
	}
}

func TestSilentPhoneIsNotMaskedByALiveComputer(t *testing.T) {
	db, h := newDeviceHarness(t)
	postHeartbeat(t, h, `{"screen":"monitoring","device_id":"phone-1","device_kind":"iphone"}`)
	postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`)
	db.ageDevice(1, "phone-1", time.Now().Add(-2*time.Hour))
	// The shipped phone path still refreshes the account timestamp. That must
	// not clear the computer, and the computer's heartbeat must not clear the phone.
	if rr := postHeartbeat(t, h, `{"screen":"active"}`); rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}
	h.sweepDeviceHeartbeatSilence(context.Background(), 30)
	summaries := db.alertSummaries()
	if len(summaries) != 1 || summaries[0] != "Ada's iPhone has stopped sending heartbeats" {
		t.Fatalf("silence = %#v", summaries)
	}
	if strings.Contains(summaries[0], "phone-1") || strings.Contains(summaries[0], "Windows") {
		t.Fatalf("silence named the wrong thing %q", summaries[0])
	}
	h.sweepDeviceHeartbeatSilence(context.Background(), 30)
	if got := db.alertSummaries(); len(got) != 1 {
		t.Fatalf("repeated silence %#v", got)
	}

	db.ageDevice(1, "pc-1", time.Now().Add(-2*time.Hour))
	h.sweepDeviceHeartbeatSilence(context.Background(), 30)
	got := db.alertSummaries()
	if len(got) != 2 || got[1] != "Ada's Windows PC has stopped sending heartbeats" {
		t.Fatalf("second device = %#v", got)
	}
}

func TestProtectionAlertNamesDeviceTypeOnly(t *testing.T) {
	db, h := newDeviceHarness(t)
	rr := httptest.NewRecorder()
	h.SendProtectionAlert(rr, authed(http.MethodPost, "/protection/alerts",
		`{"type":"monitoring_disabled","detail":"chrome.exe was open","device_id":"Safari","device_kind":"windows"}`, 1))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	got := db.alertSummaries()
	if len(got) != 1 || got[0] != "Ada turned off activity monitoring on their Windows PC" {
		t.Fatalf("body = %#v", got)
	}
	if strings.Contains(got[0], "Safari") || strings.Contains(got[0], "chrome") || strings.Contains(got[0], "http") {
		t.Fatalf("leaked %q", got[0])
	}

	rr = httptest.NewRecorder()
	h.SendProtectionAlert(rr, authed(http.MethodPost, "/protection/alerts", `{"type":"signed_out","detail":"ignored"}`, 1))
	got = db.alertSummaries()
	if len(got) != 2 || got[1] != "Ada signed out" {
		t.Fatalf("legacy sign-out = %#v", got)
	}

	rr = httptest.NewRecorder()
	h.SendProtectionAlert(rr, authed(http.MethodPost, "/protection/alerts",
		`{"type":"signed_out","device_id":"mac-1","device_kind":"chrome.exe"}`, 1))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("bad kind %d", rr.Code)
	}
	if len(db.alertSummaries()) != 2 {
		t.Fatal("rejected kind still alerted")
	}
}

func TestReportDoesNotMoveLastCheckIn(t *testing.T) {
	_, h := newDeviceHarness(t)
	postHeartbeat(t, h, `{"screen":"monitoring","device_id":"pc-1","device_kind":"windows"}`)
	first := listDevices(t, h)
	if len(first) != 1 || first[0].LastCheckIn == "" || !first[0].Monitoring {
		t.Fatalf("first = %+v", first)
	}
	rr := httptest.NewRecorder()
	h.ReportDevice(rr, authed(http.MethodPost, "/users/me/devices", `{"device_id":"pc-1","device_kind":"windows","monitoring":false}`, 1))
	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}
	second := listDevices(t, h)
	if len(second) != 1 || second[0].Monitoring || second[0].LastCheckIn != first[0].LastCheckIn {
		t.Fatalf("report changed check-in: %+v -> %+v", first[0], second[0])
	}
}

func listDevices(t *testing.T, h *H) []deviceListItem {
	t.Helper()
	rr := httptest.NewRecorder()
	h.ListDevices(rr, authed(http.MethodGet, "/users/me/devices", "", 1))
	if rr.Code != http.StatusOK {
		t.Fatalf("list %d %s", rr.Code, rr.Body.String())
	}
	var body struct {
		Devices []deviceListItem `json:"devices"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.Devices
}
