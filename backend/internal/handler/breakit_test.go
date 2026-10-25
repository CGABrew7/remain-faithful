package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

func TestRegisterRejectsLongWhitespaceAndEmailPasswords(t *testing.T) {
	h := &H{DB: newMemDB().open(t)}
	cases := []struct {
		body string
		want string
	}{
		{`{"name":"Ada","email":"ada@example.com","password":"` + strings.Repeat("a", 73) + `"}`, "too long"},
		{`{"name":"Ada","email":"ada@example.com","password":"        "}`, "not only spaces"},
		{`{"name":"Ada","email":"ada@example.com","password":"Ada@Example.com"}`, "not your email"},
		{`{"name":"Ada\u0000","email":"ada@example.com","password":"longenough1"}`, "not allowed"},
		{`{"name":"` + strings.Repeat("A", 81) + `","email":"ada@example.com","password":"longenough1"}`, "80"},
		{`{"name":"Ada","email":"not-an-email","password":"longenough1"}`, "valid email"},
	}
	for _, tc := range cases {
		rr := httptest.NewRecorder()
		h.Register(rr, httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tc.body)))
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), tc.want) {
			t.Fatalf("body %s status %d response %s, want 400 containing %q", tc.body, rr.Code, rr.Body.String(), tc.want)
		}
	}
}

func TestResetPasswordRejectsShapeBeforeConsumingToken(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 5, name: "Ada", email: "ada@example.com", verified: true})
	db.mu.Lock()
	db.resets = append(db.resets, memReset{userID: 5, hash: hashToken("reset-token"), expires: time.Now().Add(time.Hour)})
	db.mu.Unlock()
	h := &H{DB: db.open(t)}

	tooLong := httptest.NewRecorder()
	h.ResetPassword(tooLong, httptest.NewRequest(http.MethodPost, "/auth/reset-password",
		strings.NewReader(`{"token":"reset-token","password":"`+strings.Repeat("b", 80)+`"}`)))
	if tooLong.Code != http.StatusBadRequest || !strings.Contains(tooLong.Body.String(), "too long") {
		t.Fatalf("long password %d %s", tooLong.Code, tooLong.Body.String())
	}

	same := httptest.NewRecorder()
	h.ResetPassword(same, httptest.NewRequest(http.MethodPost, "/auth/reset-password",
		strings.NewReader(`{"token":"reset-token","password":"ada@example.com"}`)))
	if same.Code != http.StatusBadRequest || !strings.Contains(same.Body.String(), "not your email") {
		t.Fatalf("email password %d %s", same.Code, same.Body.String())
	}
	if len(db.resets) != 1 {
		t.Fatal("rejected password consumed the reset token")
	}

	ok := httptest.NewRecorder()
	h.ResetPassword(ok, httptest.NewRequest(http.MethodPost, "/auth/reset-password",
		strings.NewReader(`{"token":"reset-token","password":"newpassword1"}`)))
	if ok.Code != http.StatusOK {
		t.Fatalf("reset %d %s", ok.Code, ok.Body.String())
	}
	again := httptest.NewRecorder()
	h.ResetPassword(again, httptest.NewRequest(http.MethodPost, "/auth/reset-password",
		strings.NewReader(`{"token":"reset-token","password":"newpassword1"}`)))
	if again.Code != http.StatusBadRequest || !strings.Contains(again.Body.String(), "invalid or expired") {
		t.Fatalf("second reset %d %s", again.Code, again.Body.String())
	}
}

func TestPartnerInviteDecisionIsSingleUse(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.mu.Lock()
	db.pInvites = append(db.pInvites, memPInvite{id: 7, inviterID: 1, email: "bea@example.com", token: "tok", status: "pending"})
	db.mu.Unlock()
	h := &H{DB: db.open(t)}

	first := httptest.NewRecorder()
	req := authed(http.MethodPost, "/relationships/invites/7/accept", "", 2)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	h.AcceptPartnerInviteByID(first, req)
	if first.Code != http.StatusOK {
		t.Fatalf("first %d %s", first.Code, first.Body.String())
	}
	second := httptest.NewRecorder()
	req = authed(http.MethodPost, "/relationships/invites/7/decline", "", 2)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	h.DeclinePartnerInviteByID(second, req)
	if second.Code != http.StatusNotFound {
		t.Fatalf("second %d %s", second.Code, second.Body.String())
	}
}

func TestConfirmedAccountGate(t *testing.T) {
	postGroup := func(h *H, userID int64) *httptest.ResponseRecorder {
		rr := httptest.NewRecorder()
		h.CreateGroup(rr, authed(http.MethodPost, "/groups", `{"name":"Home group","covenant":"We show up."}`, userID))
		return rr
	}

	t.Run("unverified password account", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 3, name: "Cam", email: "cam@example.com", verified: false})
		rr := postGroup(&H{DB: db.open(t)}, 3)
		if rr.Code != http.StatusForbidden || !strings.Contains(rr.Body.String(), "email_not_verified") {
			t.Fatalf("status %d %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("apple sign in", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 4, name: "Ada", email: "ada@example.com", verified: false, appleID: "apple-sub"})
		rr := postGroup(&H{DB: db.open(t)}, 4)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status %d %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("verified password account", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 5, name: "Bea", email: "bea@example.com", verified: true})
		rr := postGroup(&H{DB: db.open(t)}, 5)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status %d %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("group name limits", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 5, name: "Bea", email: "bea@example.com", verified: true})
		h := &H{DB: db.open(t)}
		rr := httptest.NewRecorder()
		h.CreateGroup(rr, authed(http.MethodPost, "/groups", `{"name":"Bad\u0000name","covenant":"ok"}`, 5))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("nul %d %s", rr.Code, rr.Body.String())
		}
		rr = httptest.NewRecorder()
		h.CreateGroup(rr, authed(http.MethodPost, "/groups", `{"name":"`+strings.Repeat("G", 81)+`","covenant":"ok"}`, 5))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("long %d %s", rr.Code, rr.Body.String())
		}
	})
}

func TestDeviceTokenAndCampaignLimits(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	h := &H{DB: db.open(t)}
	rr := httptest.NewRecorder()
	h.RegisterDeviceToken(rr, authed(http.MethodPost, "/users/device-token",
		`{"token":"`+strings.Repeat("a", 300)+`","platform":"ios"}`, 1))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("token %d %s", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	h.RegisterDeviceToken(rr, authed(http.MethodPost, "/users/device-token",
		`{"token":"abc","platform":"desktop"}`, 1))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("platform %d %s", rr.Code, rr.Body.String())
	}

	if got := cleanUTM("hello\u202eworld" + strings.Repeat("x", 100)); strings.ContainsRune(got, '\u202e') || utf8.RuneCountInString(got) > 80 {
		t.Fatalf("utm = %q", got)
	}
	if got := referrerDomain("https://example.com/\u202ehidden"); strings.ContainsRune(got, '\u202e') {
		t.Fatalf("referrer = %q", got)
	}
}

func TestContactInboxRequired(t *testing.T) {
	t.Setenv("CONTACT_TO_EMAIL", "")
	em := &capturingEmail{}
	h := &H{Email: em}
	rr := httptest.NewRecorder()
	h.Contact(rr, httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(
		`{"name":"Ada","email":"ada@example.com","message":"Hello there"}`,
	)))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	if len(em.calls) != 0 {
		t.Fatal("message was sent without a configured inbox")
	}
}

func TestResendConfirmationIsUniform(t *testing.T) {
	h := &H{DB: newMemDB().open(t)}
	rr := httptest.NewRecorder()
	h.ResendConfirmationByEmail(rr, httptest.NewRequest(http.MethodPost, "/auth/resend-confirmation",
		strings.NewReader(`{"email":"missing@example.com"}`)))
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"status":"ok"`) {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
}

func TestPushNameIsClipped(t *testing.T) {
	long := strings.Repeat("A", 200)
	body := alertBody(long, "adult_content", "severe", "Explicit content detected")
	if utf8.RuneCountInString(body) > 400 || !strings.Contains(body, strings.Repeat("A", 80)) {
		t.Fatalf("body = %q", body)
	}
	if strings.Contains(body, strings.Repeat("A", 81)) {
		t.Fatal("push body kept more than 80 letters of the name")
	}
	if got := protectionAlertBody(long, "signed_out", "ignored"); strings.Contains(got, strings.Repeat("A", 81)) {
		t.Fatalf("protection body = %q", got)
	}
}
