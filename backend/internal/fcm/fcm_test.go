package fcm

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewNoopWithoutCredentials(t *testing.T) {
	c, err := New("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !c.IsNoop() {
		t.Fatal("expected noop")
	}
	if err := c.Send(context.Background(), &Message{Token: "abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestFromAPNSCopiesRoutingKeysOnly(t *testing.T) {
	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "App Usage Alert",
				"body":  "Ada opened an app they asked to block.",
			},
			"sound": "default",
		},
		"notification_type": "CONTENT_FLAGGED",
		"sender_name":       "Ada",
		"alert_id":          int64(9),
		"category":          "app_usage",
		"severity":          "concerning",
		"summary":           "Blocked app opened",
	}
	msg := FromAPNS("tokendevice", "event-3", payload)
	if msg.Title != "App Usage Alert" || msg.Body != "Ada opened an app they asked to block." {
		t.Fatalf("title/body = %q / %q", msg.Title, msg.Body)
	}
	if msg.Data["notification_type"] != "CONTENT_FLAGGED" {
		t.Fatalf("type = %v", msg.Data)
	}
	if msg.Data["summary"] != "Blocked app opened" || msg.Data["alert_id"] != "9" {
		t.Fatalf("data = %#v", msg.Data)
	}
	if _, ok := msg.Data["aps"]; ok {
		t.Fatal("aps leaked into data")
	}
	for _, forbidden := range []string{"url", "app_name", "package", "screen"} {
		if _, ok := msg.Data[forbidden]; ok {
			t.Fatalf("unexpected key %s", forbidden)
		}
	}
	if msg.ChannelID != "alerts" {
		t.Fatal(msg.ChannelID)
	}
}

func TestFromAPNSPanicUsesUrgentChannel(t *testing.T) {
	msg := FromAPNS("t", "panic-1", map[string]any{
		"aps": map[string]any{
			"alert": map[string]any{"title": "Urgent Prayer Request", "body": "Ada needs support right now"},
		},
		"notification_type": "PANIC_ALERT",
		"sender_name":       "Ada",
	})
	if msg.ChannelID != "urgent" {
		t.Fatal(msg.ChannelID)
	}
	if msg.Title != "Urgent Prayer Request" {
		t.Fatal(msg.Title)
	}
}

func TestSendPostsMetadataAndRejectsUnregistered(t *testing.T) {
	key, pemText := testKey(t)
	var sawAuth bool
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "token") || r.URL.Path == "/" && r.Method == http.MethodPost && strings.Contains(r.Header.Get("Content-Type"), "form") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"access_token":"ya29.test"}`)
			return
		}
		sawAuth = strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":{"status":"NOT_FOUND","details":[{"errorCode":"UNREGISTERED"}]}}`)
	}))
	defer srv.Close()

	c := &Client{
		projectID:   "remain-faithful-placeholder",
		clientEmail: "push@example.iam.gserviceaccount.com",
		key:         key,
		http:        srv.Client(),
		tokenURL:    srv.URL + "/token",
		sendHost:    srv.URL,
	}
	_ = pemText
	msg := FromAPNS("device-token-value", "event-1", map[string]any{
		"aps":               map[string]any{"alert": map[string]string{"title": "Protection Update", "body": "Ada signed out"}},
		"notification_type": "PROTECTION_ALERT",
		"alert_type":        "signed_out",
		"sender_name":       "Ada",
	})
	err := c.Send(context.Background(), msg)
	var invalid *ErrInvalidToken
	if !errorAs(err, &invalid) {
		t.Fatalf("err = %v", err)
	}
	if !sawAuth {
		t.Fatal("missing bearer")
	}
	message, _ := body["message"].(map[string]any)
	data, _ := message["data"].(map[string]any)
	if data["notification_type"] != "PROTECTION_ALERT" || data["alert_type"] != "signed_out" {
		t.Fatalf("body = %#v", body)
	}
	if _, ok := data["url"]; ok {
		t.Fatal("url in fcm data")
	}
}

func errorAs(err error, target **ErrInvalidToken) bool {
	if err == nil {
		return false
	}
	inv, ok := err.(*ErrInvalidToken)
	if !ok {
		return false
	}
	*target = inv
	return true
}

func testKey(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	block := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	return key, string(block)
}
