package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

type scriptedAPNS struct{ err error }

func (s scriptedAPNS) Send(context.Context, *apns.Notification) error { return s.err }
func (s scriptedAPNS) IsNoop() bool                                   { return false }
func (s scriptedAPNS) Environment() string                            { return "sandbox" }

func authed(method, path, body string, userID int64) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	return req.WithContext(rfauth.ContextWithUser(req.Context(), userID))
}

func TestRegisterDetachesDeviceTokenFromOtherUsers(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.mu.Lock()
	db.tokens = append(db.tokens, memToken{userID: 1, token: "tok-shared", platform: "ios", environment: "sandbox", active: true})
	db.mu.Unlock()

	h := &H{DB: db.open(t)}
	req := authed(http.MethodPost, "/users/device-token",
		`{"token":"tok-shared","platform":"ios","environment":"sandbox"}`, 2)
	rr := httptest.NewRecorder()
	h.RegisterDeviceToken(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("register status %d body %s", rr.Code, rr.Body.String())
	}

	active := db.activeTokenUsers("tok-shared")
	if len(active) != 1 || active[0] != 2 {
		t.Fatalf("active owners = %v, want only user 2", active)
	}
}

func TestSignedOutClearsCallerDeviceTokens(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.mu.Lock()
	db.tokens = append(db.tokens, memToken{userID: 1, token: "tok-ada", platform: "ios", environment: "sandbox", active: true})
	db.mu.Unlock()

	h := &H{DB: db.open(t), APNS: scriptedAPNS{}}
	req := authed(http.MethodPost, "/protection/alerts", `{"type":"signed_out","detail":"Signed out."}`, 1)
	rr := httptest.NewRecorder()
	h.SendProtectionAlert(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	if active := db.activeTokenUsers("tok-ada"); len(active) != 0 {
		t.Fatalf("token still active for %v", active)
	}
}

func TestRejectedTokenUsesSendTargetAsSubject(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.addUser(memUser{id: 3, name: "Cam", email: "cam@example.com", verified: true})
	db.addUser(memUser{id: 4, name: "Dee", email: "dee@example.com", verified: true})
	db.mu.Lock()
	db.tokens = append(db.tokens,
		memToken{userID: 2, token: "tok-shared", platform: "ios", environment: "sandbox", active: true},
		memToken{userID: 1, token: "tok-shared", platform: "ios", environment: "sandbox", active: true},
	)
	db.rels = append(db.rels,
		memRel{id: 1, userID: 1, partnerID: 3, typ: "partner", status: "accepted"},
		memRel{id: 2, userID: 2, partnerID: 4, typ: "partner", status: "accepted"},
	)
	db.mu.Unlock()

	h := &H{DB: db.open(t), APNS: scriptedAPNS{err: &apns.ErrInvalidToken{Reason: "Unregistered"}}, SyncPush: true}
	h.notifyPartnerByID(context.Background(), 1, "someone", &apns.Notification{})

	var subjects []int64
	var summaries []string
	for _, a := range db.alertRows() {
		if a.alertType != "app_removed" {
			continue
		}
		subjects = append(subjects, a.subject)
		summaries = append(summaries, a.summary)
	}
	if len(subjects) != 1 || subjects[0] != 1 {
		t.Fatalf("app_removed subjects = %v summaries %v, want subject 1", subjects, summaries)
	}
	if strings.Contains(summaries[0], "Bea") || !strings.Contains(summaries[0], "Ada") {
		t.Fatalf("summary = %q, want Ada and not Bea", summaries[0])
	}
}

func TestDeleteUserRemovesPartnerProtectionAlerts(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 3, name: "Cam", email: "cam@example.com", verified: true})
	db.mu.Lock()
	db.alerts = append(db.alerts, memAlert{
		id: 1, recipient: 3, kind: "protection", alertType: "signed_out",
		summary: "Ada signed out", subject: 1,
	})
	db.mu.Unlock()

	h := &H{DB: db.open(t)}
	if err := h.deleteUserByID(context.Background(), 1, "Ada"); err != nil {
		t.Fatal(err)
	}
	for _, a := range db.alertRows() {
		if a.recipient == 3 && strings.Contains(a.summary, "Ada") {
			t.Fatalf("partner alert still names the deleted user: %+v", a)
		}
	}
}

func TestUnverifiedEmailCannotAcceptInvites(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: false})
	db.mu.Lock()
	db.pInvites = append(db.pInvites,
		memPInvite{id: 7, inviterID: 1, email: "bea@example.com", token: "raw-token", status: "pending"},
		memPInvite{id: 11, inviterID: 1, email: "bea@example.com", token: "raw-token-2", status: "pending"},
	)
	db.gInvites = append(db.gInvites, memGInvite{id: 8, inviterID: 1, groupID: 4, email: "bea@example.com", token: "group-token", status: "pending"})
	db.groups[4] = true
	db.mu.Unlock()
	h := &H{DB: db.open(t)}

	t.Run("partner list asks to confirm email", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ListPendingPartnerInvites(rr, authed(http.MethodGet, "/relationships/invites/pending", "", 2))
		assertCodedError(t, rr, http.StatusForbidden, "email_not_verified")
	})

	t.Run("group list asks to confirm email", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ListPendingGroupInvites(rr, authed(http.MethodGet, "/groups/invites/pending", "", 2))
		assertCodedError(t, rr, http.StatusForbidden, "email_not_verified")
	})

	t.Run("partner accept by id", func(t *testing.T) {
		req := authed(http.MethodPost, "/relationships/invites/7/accept", "", 2)
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		rr := httptest.NewRecorder()
		h.AcceptPartnerInviteByID(rr, req)
		assertCodedError(t, rr, http.StatusForbidden, "email_not_verified")
		if len(db.relRows()) != 0 {
			t.Fatalf("relationship created: %+v", db.relRows())
		}
	})

	t.Run("partner accept by token", func(t *testing.T) {
		req := authed(http.MethodPost, "/relationships/accept-invite", `{"token":"raw-token-2"}`, 2)
		rr := httptest.NewRecorder()
		h.AcceptPartnerInvite(rr, req)
		assertCodedError(t, rr, http.StatusForbidden, "email_not_verified")
	})

	t.Run("group accept by id", func(t *testing.T) {
		req := authed(http.MethodPost, "/groups/invites/8/accept", "", 2)
		req = mux.SetURLVars(req, map[string]string{"id": "8"})
		rr := httptest.NewRecorder()
		h.AcceptGroupInviteByID(rr, req)
		assertCodedError(t, rr, http.StatusForbidden, "email_not_verified")
		if len(db.memberRows()) != 0 {
			t.Fatalf("group member created: %+v", db.memberRows())
		}
	})

	t.Run("group accept by token", func(t *testing.T) {
		req := authed(http.MethodPost, "/groups/accept-invite", `{"token":"group-token"}`, 2)
		rr := httptest.NewRecorder()
		h.AcceptGroupInvite(rr, req)
		assertCodedError(t, rr, http.StatusForbidden, "email_not_verified")
	})
}

func TestVerifiedEmailCanAcceptPartnerInvite(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.mu.Lock()
	db.pInvites = append(db.pInvites, memPInvite{id: 7, inviterID: 1, email: "bea@example.com", token: "raw-token", status: "pending"})
	db.mu.Unlock()
	h := &H{DB: db.open(t)}
	listed := httptest.NewRecorder()
	h.ListPendingPartnerInvites(listed, authed(http.MethodGet, "/relationships/invites/pending", "", 2))
	if listed.Code != http.StatusOK {
		t.Fatalf("list status %d body %s", listed.Code, listed.Body.String())
	}
	var rows []map[string]any
	if err := json.Unmarshal(listed.Body.Bytes(), &rows); err != nil {
		t.Fatalf("list body %s: %v", listed.Body.String(), err)
	}
	if len(rows) != 1 {
		t.Fatalf("verified list = %s", listed.Body.String())
	}
	req := authed(http.MethodPost, "/relationships/invites/7/accept", "", 2)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()
	h.AcceptPartnerInviteByID(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	rels := db.relRows()
	if len(rels) != 1 || rels[0].userID != 1 || rels[0].partnerID != 2 || rels[0].status != "accepted" {
		t.Fatalf("rels = %+v", rels)
	}
}

func TestSocialLinkClearsUnverifiedPasswordAndBumpsVersion(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{
		id: 5, name: "Ada", email: "ada@example.com", password: "stored-hash",
		verified: false, tokenVersion: 1,
	})
	h := &H{DB: db.open(t)}
	req := httptest.NewRequest(http.MethodPost, "/auth/google", nil)
	if _, _, _, err := h.findOrCreateSocialUser(req, "google_id", "sub-1", "ada@example.com", "Ada", true); err != nil {
		t.Fatal(err)
	}
	u, ok := db.user(5)
	if !ok {
		t.Fatal("user missing")
	}
	if u.password != "" {
		t.Fatalf("password still set: %q", u.password)
	}
	if u.tokenVersion == 1 {
		t.Fatal("token version was not bumped")
	}
	if h.TokenVersionAllowed(5, 1) {
		t.Fatal("old token version still allowed")
	}
	if !h.TokenVersionAllowed(5, u.tokenVersion) {
		t.Fatal("new token version rejected")
	}
}

func TestSocialLinkKeepsPasswordWhenEmailAlreadyVerified(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{
		id: 5, name: "Ada", email: "ada@example.com", password: "stored-hash",
		verified: true, tokenVersion: 4,
	})
	h := &H{DB: db.open(t)}
	req := httptest.NewRequest(http.MethodPost, "/auth/apple", nil)
	if _, _, _, err := h.findOrCreateSocialUser(req, "apple_id", "apple-sub", "ada@example.com", "Ada", true); err != nil {
		t.Fatal(err)
	}
	u, _ := db.user(5)
	if u.password != "stored-hash" || u.tokenVersion != 4 {
		t.Fatalf("verified link changed credentials: password %q version %d", u.password, u.tokenVersion)
	}
}

func TestPasswordResetClearsPendingEmail(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{
		id: 5, name: "Ada", email: "ada@example.com", password: "old-hash",
		verified: true, tokenVersion: 1,
		pendingEmail: "new@example.com", pendingHash: hashToken("change-token"),
		pendingExpires: time.Now().Add(time.Hour),
	})
	db.mu.Lock()
	db.resets = append(db.resets, memReset{userID: 5, hash: hashToken("reset-token"), expires: time.Now().Add(time.Hour)})
	db.mu.Unlock()
	h := &H{DB: db.open(t)}

	rr := httptest.NewRecorder()
	h.ResetPassword(rr, httptest.NewRequest(http.MethodPost, "/auth/reset-password",
		strings.NewReader(`{"token":"reset-token","password":"newpassword1"}`)))
	if rr.Code != http.StatusOK {
		t.Fatalf("reset status %d body %s", rr.Code, rr.Body.String())
	}
	u, _ := db.user(5)
	if u.pendingEmail != "" || u.pendingHash != "" {
		t.Fatalf("pending email still staged: %+v", u)
	}

	rr = httptest.NewRecorder()
	h.VerifyEmail(rr, httptest.NewRequest(http.MethodPost, "/auth/verify-email",
		strings.NewReader(`{"token":"change-token"}`)))
	if rr.Code < 400 {
		t.Fatalf("verify after reset status %d body %s", rr.Code, rr.Body.String())
	}
	u, _ = db.user(5)
	if u.email != "ada@example.com" {
		t.Fatalf("email changed to %q", u.email)
	}
}

func TestCreateRelationshipResponseDoesNotDependOnRegistration(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	h := &H{DB: db.open(t)}

	known := httptest.NewRecorder()
	h.CreateRelationship(known, authed(http.MethodPost, "/relationships",
		`{"partner_email":"bea@example.com","type":"partner"}`, 1))
	unknown := httptest.NewRecorder()
	h.CreateRelationship(unknown, authed(http.MethodPost, "/relationships",
		`{"partner_email":"nobody@example.com","type":"partner"}`, 1))
	if known.Code != unknown.Code {
		t.Fatalf("status known %d unknown %d", known.Code, unknown.Code)
	}
	if canonicalJSON(t, known.Body.Bytes()) != canonicalJSON(t, unknown.Body.Bytes()) {
		t.Fatalf("bodies differ\nknown   %s\nunknown %s", known.Body.String(), unknown.Body.String())
	}
}

func TestCreateRelationshipDoesNotRewriteExistingType(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.mu.Lock()
	db.rels = append(db.rels, memRel{id: 9, userID: 1, partnerID: 2, typ: "partner", status: "accepted"})
	db.nextRel = 10
	db.mu.Unlock()
	h := &H{DB: db.open(t)}
	rr := httptest.NewRecorder()
	h.CreateRelationship(rr, authed(http.MethodPost, "/relationships",
		`{"partner_email":"bea@example.com","type":"friend"}`, 1))
	if rr.Code >= 400 {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	rels := db.relRows()
	if len(rels) != 1 || rels[0].typ != "partner" || rels[0].status != "accepted" {
		t.Fatalf("relationship changed: %+v", rels)
	}
}

func TestHeartbeatMonitoringUpdatesTimestamp(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	h := &H{DB: db.open(t)}
	rr := httptest.NewRecorder()
	h.Heartbeat(rr, authed(http.MethodPost, "/heartbeat", `{"screen":"monitoring"}`, 1))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	u, _ := db.user(1)
	if !u.heartbeat || u.screen != "monitoring" {
		t.Fatalf("heartbeat state = %+v", u)
	}
}

func canonicalJSON(t *testing.T, raw []byte) string {
	t.Helper()
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("json %q: %v", raw, err)
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}
