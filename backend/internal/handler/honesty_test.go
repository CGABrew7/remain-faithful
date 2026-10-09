package handler

import (
	"context"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

type capturingEmail struct {
	mu    sync.Mutex
	calls []contactCall
}

type contactCall struct {
	fromEmail, fromName, subject, message, to string
}

func (e *capturingEmail) SendPasswordReset(string, string, string) error { return nil }
func (e *capturingEmail) SendPartnerInvite(string, string, string) error { return nil }
func (e *capturingEmail) SendGroupInvite(string, string, string, string) error {
	return nil
}
func (e *capturingEmail) SendContact(fromEmail, fromName, subject, message, toEmail string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.calls = append(e.calls, contactCall{fromEmail, fromName, subject, message, toEmail})
	return nil
}

type capturingAPNS struct {
	mu   sync.Mutex
	sent []apns.Notification
}

func (c *capturingAPNS) Send(_ context.Context, n *apns.Notification) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, *n)
	return nil
}
func (c *capturingAPNS) IsNoop() bool        { return false }
func (c *capturingAPNS) Environment() string { return "sandbox" }

func pushBody(n apns.Notification) string {
	payload, _ := n.Payload.(map[string]any)
	aps, _ := payload["aps"].(map[string]any)
	alert, _ := aps["alert"].(map[string]string)
	return alert["body"]
}

func withUser(req *http.Request, id int64) *http.Request {
	return req.WithContext(rfauth.ContextWithUser(req.Context(), id))
}

func TestNormalizeFeedbackOmitsDeviceUnlessProvided(t *testing.T) {
	subject, message, errMsg := normalizeFeedback("Idea", "  Shorter check-ins  ", "Please add a weekly note.", "", "", "")
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if subject != "[Idea] Shorter check-ins" {
		t.Fatalf("subject = %q", subject)
	}
	if strings.Contains(message, "App:") || strings.Contains(message, "Device:") || strings.Contains(message, "iPhone") {
		t.Fatalf("device lines leaked into message:\n%s", message)
	}
	if !strings.Contains(message, "Type: Idea") || !strings.Contains(message, "Please add a weekly note.") {
		t.Fatalf("message = %q", message)
	}

	_, withDevice, errMsg := normalizeFeedback("Problem", "Crash", "It closed.", "Remain Faithful 1.0 (8)", "iPhone15,2", "iOS 18.1")
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	for _, want := range []string{"App: Remain Faithful 1.0 (8)", "Device: iPhone15,2", "OS: iOS 18.1"} {
		if !strings.Contains(withDevice, want) {
			t.Fatalf("missing %q in\n%s", want, withDevice)
		}
	}
}

func TestNormalizeFeedbackRejectsBadInput(t *testing.T) {
	if _, _, errMsg := normalizeFeedback("Spam", "Hi", "More", "", "", ""); errMsg == "" {
		t.Fatal("expected type rejection")
	}
	if _, _, errMsg := normalizeFeedback("Idea", "", "More", "", "", ""); errMsg == "" {
		t.Fatal("expected empty summary rejection")
	}
	if _, _, errMsg := normalizeFeedback("Idea", "Hi", "", "", "", ""); errMsg == "" {
		t.Fatal("expected empty details rejection")
	}
}

func TestSubmitFeedbackUsesAccountAndSkipsDevice(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		if strings.Contains(q, "SELECT name, email FROM users") {
			return []string{"name", "email"}, [][]driver.Value{{"Jeff", "jeff@example.com"}}
		}
		return none("x")
	}}
	em := &capturingEmail{}
	h := &H{DB: f.open(t), Email: em}
	body := `{"type":"Improvement","summary":"Clearer pause text","details":"Say who gets told."}`
	req := withUser(httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body)), 4)
	rr := httptest.NewRecorder()
	h.SubmitFeedback(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
	if len(em.calls) != 1 {
		t.Fatalf("emails = %d", len(em.calls))
	}
	call := em.calls[0]
	if call.fromEmail != "jeff@example.com" || call.fromName != "Jeff" {
		t.Fatalf("from = %s <%s>", call.fromName, call.fromEmail)
	}
	if call.subject != "[Improvement] Clearer pause text" {
		t.Fatalf("subject = %q", call.subject)
	}
	if strings.Contains(call.message, "App:") || strings.Contains(call.message, "Device:") {
		t.Fatalf("device leaked: %s", call.message)
	}
}

func TestSubmitFeedbackHoneypot(t *testing.T) {
	em := &capturingEmail{}
	h := &H{Email: em}
	body := `{"type":"Idea","summary":"x","details":"y","company":"Acme Bots"}`
	req := withUser(httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body)), 4)
	rr := httptest.NewRecorder()
	h.SubmitFeedback(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if len(em.calls) != 0 {
		t.Fatal("honeypot should not send email")
	}
}

func TestContactHoneypot(t *testing.T) {
	em := &capturingEmail{}
	h := &H{Email: em}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(
		`{"name":"Bot","email":"bot@evil.test","message":"<script>","company":"filled"}`,
	))
	rr := httptest.NewRecorder()
	h.Contact(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if len(em.calls) != 0 {
		t.Fatal("honeypot should not send email")
	}
}

func TestSupportAudienceIDs(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		if strings.Contains(q, "other_id") {
			return []string{"other_id"}, [][]driver.Value{{int64(2)}, {int64(5)}, {int64(2)}}
		}
		return none("x")
	}}
	h := &H{DB: f.open(t)}
	ids, err := h.supportAudienceIDs(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 3 {
		t.Fatalf("ids = %v (SQL distinct is the database's job; the scan returns each row)", ids)
	}
}

func TestPanicNotifiesPartnersAndGroupMembers(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		switch {
		case strings.Contains(q, "SELECT name FROM users"):
			return one("name", "Alex")
		case strings.Contains(q, "other_id"):
			return []string{"other_id"}, [][]driver.Value{{int64(2)}, {int64(9)}}
		case strings.Contains(q, "SELECT token FROM device_tokens"):
			return one("token", "tokendevice")
		default:
			return none("x")
		}
	}}
	push := &capturingAPNS{}
	h := &H{DB: f.open(t), APNS: push, SyncPush: true}
	req := withUser(httptest.NewRequest(http.MethodPost, "/panic", strings.NewReader(`{}`)), 1)
	rr := httptest.NewRecorder()
	h.SendPanicAlert(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
	if len(push.sent) != 2 {
		t.Fatalf("pushes = %d, want 2 (partner and group member)", len(push.sent))
	}
	for _, n := range push.sent {
		if pushBody(n) != "Alex needs support right now" {
			t.Fatalf("body = %q", pushBody(n))
		}
		if n.DeviceToken != "tokendevice" {
			t.Fatalf("token = %q", n.DeviceToken)
		}
	}
}

func TestPanicRequiresSomeoneToNotify(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		if strings.Contains(q, "SELECT name FROM users") {
			return one("name", "Alex")
		}
		return none("other_id")
	}}
	h := &H{DB: f.open(t), APNS: &capturingAPNS{}}
	req := withUser(httptest.NewRequest(http.MethodPost, "/panic", strings.NewReader(`{}`)), 1)
	rr := httptest.NewRecorder()
	h.SendPanicAlert(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
}

func TestUpdateGroupPersistsNameAndNotifiesOnCovenant(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		switch {
		case strings.Contains(q, "SELECT role FROM group_members"):
			return one("role", "member")
		case strings.Contains(q, "SELECT name, covenant FROM groups"):
			return []string{"name", "covenant"}, [][]driver.Value{{"Old Name", "old words"}}
		case strings.Contains(q, "SELECT name FROM users"):
			return one("name", "Sam")
		case strings.Contains(q, "SELECT user_id FROM group_members"):
			return one("user_id", int64(8))
		case strings.Contains(q, "SELECT token FROM device_tokens"):
			return one("token", "tokendevice")
		default:
			return none("x")
		}
	}}
	push := &capturingAPNS{}
	h := &H{DB: f.open(t), APNS: push, SyncPush: true}
	body := `{"name":"Tuesday Men","covenant":"UNIQUE-COVENANT-TEXT"}`
	req := withUser(httptest.NewRequest(http.MethodPatch, "/groups/7", strings.NewReader(body)), 1)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()
	h.UpdateGroup(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "Tuesday Men") {
		t.Fatalf("response = %s", rr.Body.String())
	}
	updates := f.ran("UPDATE groups")
	if len(updates) != 1 {
		t.Fatalf("updates = %d", len(updates))
	}
	if len(push.sent) != 1 {
		t.Fatalf("pushes = %d", len(push.sent))
	}
	got := pushBody(push.sent[0])
	if got != covenantUpdatedBody("Sam") {
		t.Fatalf("body = %q", got)
	}
	if strings.Contains(got, "UNIQUE-COVENANT-TEXT") {
		t.Fatal("covenant text must not ride along in the push")
	}
}

func TestUpdateGroupRejectsNonMember(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		return none("role")
	}}
	h := &H{DB: f.open(t)}
	req := withUser(httptest.NewRequest(http.MethodPatch, "/groups/7", strings.NewReader(`{"name":"Nope"}`)), 1)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()
	h.UpdateGroup(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
	if len(f.ran("UPDATE groups")) != 0 {
		t.Fatal("non-member must not update the group")
	}
}

func TestSendEncouragementPushesFixedLine(t *testing.T) {
	f := &inviteFakeDB{query: func(q string, _ []any) ([]string, [][]driver.Value) {
		switch {
		case strings.Contains(q, "SELECT COUNT(*) FROM group_members"):
			return one("count", int64(2))
		case strings.Contains(q, "SELECT name FROM users"):
			return one("name", "Sam")
		case strings.Contains(q, "SELECT token FROM device_tokens"):
			return one("token", "tokendevice")
		default:
			return none("x")
		}
	}}
	push := &capturingAPNS{}
	h := &H{DB: f.open(t), APNS: push, SyncPush: true}
	req := withUser(httptest.NewRequest(http.MethodPost, "/groups/7/encouragement", strings.NewReader(`{"user_id":8}`)), 1)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()
	h.SendEncouragement(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
	if len(push.sent) != 1 {
		t.Fatalf("pushes = %d", len(push.sent))
	}
	if pushBody(push.sent[0]) != "Sam sent you encouragement." {
		t.Fatalf("body = %q", pushBody(push.sent[0]))
	}
}

func TestNoticeCopy(t *testing.T) {
	if groupLeftBody("Sam") != "Sam left the group." {
		t.Fatal(groupLeftBody("Sam"))
	}
	if partnershipEndedBody("Sam") != "Sam ended the accountability partnership." {
		t.Fatal(partnershipEndedBody("Sam"))
	}
}
