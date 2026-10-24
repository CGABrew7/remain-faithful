package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"remain-faithful/backend/internal/apns"
	"remain-faithful/backend/internal/fcm"
)

type recordFCM struct {
	msgs []*fcm.Message
}

func (r *recordFCM) Send(_ context.Context, msg *fcm.Message) error {
	r.msgs = append(r.msgs, msg)
	return nil
}

func (r *recordFCM) IsNoop() bool { return false }

type recordAPNS struct {
	notes []*apns.Notification
}

func (r *recordAPNS) Send(_ context.Context, n *apns.Notification) error {
	copied := *n
	r.notes = append(r.notes, &copied)
	return nil
}

func (r *recordAPNS) IsNoop() bool        { return false }
func (r *recordAPNS) Environment() string { return "sandbox" }

func TestAndroidTokenUsesFCMAndLeavesAPNs(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 7, name: "Ada", email: "ada@example.com", verified: true})
	db.mu.Lock()
	db.tokens = append(db.tokens,
		memToken{userID: 7, token: "ios-token-aaa", platform: "ios", environment: "sandbox", active: true},
		memToken{userID: 7, token: "android-token-bbb", platform: "android", environment: "sandbox", active: true},
	)
	db.mu.Unlock()

	apnsRec := &recordAPNS{}
	fcmRec := &recordFCM{}
	h := &H{DB: db.open(t), APNS: apnsRec, FCM: fcmRec}
	note := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: "event-4",
		Payload: map[string]any{
			"aps": map[string]any{
				"alert": map[string]string{"title": "App Usage Alert", "body": "Ada opened an app they asked to block."},
			},
			"notification_type": "CONTENT_FLAGGED",
			"sender_name":       "Ada",
			"summary":           "Blocked app opened",
			"category":          "app_usage",
		},
	}
	h.notifyPartnerByID(context.Background(), 7, "Ada", note)

	if len(apnsRec.notes) != 1 || apnsRec.notes[0].DeviceToken != "ios-token-aaa" {
		t.Fatalf("apns tokens = %#v", apnsRec.notes)
	}
	if len(fcmRec.msgs) != 1 {
		t.Fatalf("fcm count = %d", len(fcmRec.msgs))
	}
	msg := fcmRec.msgs[0]
	if msg.Token != "android-token-bbb" {
		t.Fatalf("fcm token = %s", msg.Token)
	}
	if msg.Data["notification_type"] != "CONTENT_FLAGGED" || msg.Data["summary"] != "Blocked app opened" {
		t.Fatalf("data = %#v", msg.Data)
	}
	if msg.Title != "App Usage Alert" {
		t.Fatal(msg.Title)
	}
}

func TestRegisterAcceptsAndroidPlatform(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	h := &H{DB: db.open(t)}
	req := authed(http.MethodPost, "/users/device-token",
		`{"token":"fcm-token","platform":"Android","environment":"production"}`, 1)
	rr := httptest.NewRecorder()
	h.RegisterDeviceToken(rr, req)
	if rr.Code != 200 {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if len(db.tokens) != 1 || db.tokens[0].platform != "android" || db.tokens[0].environment != "production" {
		t.Fatalf("stored %#v", db.tokens)
	}
}
