package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func assertCodedError(t *testing.T, rr *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if rr.Code != status {
		t.Fatalf("status %d body %s, want %d", rr.Code, rr.Body.String(), status)
	}
	body := decodeBody(t, rr)
	if body["code"] != code {
		t.Fatalf("code = %v, body %s", body["code"], rr.Body.String())
	}
	msg, _ := body["error"].(string)
	if strings.TrimSpace(msg) == "" {
		t.Fatalf("missing error message: %s", rr.Body.String())
	}
}

func TestSetRelationshipPINPendingAndAccepted(t *testing.T) {
	postPIN := func(h *H, relID, userID int64) *httptest.ResponseRecorder {
		t.Helper()
		req := authed(http.MethodPost, "/relationships/4/pin", `{"pin":"1234"}`, userID)
		req = mux.SetURLVars(req, map[string]string{"id": "4"})
		rr := httptest.NewRecorder()
		h.SetRelationshipPIN(rr, req)
		return rr
	}

	t.Run("pending invite", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
		db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: false})
		db.mu.Lock()
		db.rels = append(db.rels, memRel{id: 4, userID: 1, partnerID: 2, status: "pending"})
		db.mu.Unlock()
		h := &H{DB: db.open(t)}
		rr := postPIN(h, 4, 2)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		if strings.Contains(rr.Body.String(), "relationship_not_accepted") {
			t.Fatalf("pending pin used a distinct code: %s", rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "not authorized to set PIN on this relationship") {
			t.Fatalf("body %s", rr.Body.String())
		}
		if db.relRows()[0].pinHash != "" {
			t.Fatal("pending partnership stored a PIN")
		}
	})

	t.Run("accepted partnership", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
		db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
		db.mu.Lock()
		db.rels = append(db.rels, memRel{id: 4, userID: 1, partnerID: 2, status: "accepted"})
		db.mu.Unlock()
		h := &H{DB: db.open(t)}
		rr := postPIN(h, 4, 2)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		if db.relRows()[0].pinHash == "" {
			t.Fatal("accepted partnership did not store a PIN")
		}
	})

	t.Run("unknown or someone else's partnership", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
		db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
		db.addUser(memUser{id: 3, name: "Cam", email: "cam@example.com", verified: true})
		db.mu.Lock()
		db.rels = append(db.rels, memRel{id: 4, userID: 1, partnerID: 2, status: "accepted"})
		db.mu.Unlock()
		h := &H{DB: db.open(t)}
		rr := postPIN(h, 4, 3)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		if strings.Contains(rr.Body.String(), "relationship_not_accepted") {
			t.Fatalf("stranger got the pending message: %s", rr.Body.String())
		}
		missing := authed(http.MethodPost, "/relationships/99/pin", `{"pin":"1234"}`, 2)
		missing = mux.SetURLVars(missing, map[string]string{"id": "99"})
		rr = httptest.NewRecorder()
		h.SetRelationshipPIN(rr, missing)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("missing status %d body %s", rr.Code, rr.Body.String())
		}
	})
}

func TestAppleVerifiedEmailMarksAccountVerified(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	key := withAppleKey(t)
	now := time.Now()
	signIn := func(t *testing.T, h *H, email string, verified bool) *httptest.ResponseRecorder {
		t.Helper()
		raw := signApple(t, key, "com.remainfaithful.app", now.Add(time.Hour), email, verified)
		body, err := json.Marshal(map[string]string{"identity_token": raw})
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		h.AppleSignIn(rr, httptest.NewRequest(http.MethodPost, "/auth/apple", strings.NewReader(string(body))))
		return rr
	}

	t.Run("new private relay address", func(t *testing.T) {
		db := newMemDB()
		h := &H{DB: db.open(t)}
		rr := signIn(t, h, "Hide.Me@privaterelay.appleid.com", true)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		u := onlyMemUser(t, db)
		if u.email != "hide.me@privaterelay.appleid.com" || !u.verified {
			t.Fatalf("user = %+v", u)
		}
		if placeholderAppleEmail(u.appleID, u.email) {
			t.Fatal("private relay address was stored as a stand-in")
		}
	})

	t.Run("later sign-in confirms an existing address", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{
			id: 8, name: "Bea", email: "hide.me@privaterelay.appleid.com",
			appleID: "apple-sub-1", verified: false,
		})
		db.mu.Lock()
		db.pInvites = append(db.pInvites, memPInvite{
			id: 3, inviterID: 1, email: "hide.me@privaterelay.appleid.com",
			token: "invite-token", status: "pending",
		})
		db.mu.Unlock()
		h := &H{DB: db.open(t)}
		rr := signIn(t, h, "hide.me@privaterelay.appleid.com", true)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		u, ok := db.user(8)
		if !ok || !u.verified {
			t.Fatalf("user = %+v", u)
		}
		rels := db.relRows()
		if len(rels) != 1 || rels[0].status != "accepted" || rels[0].partnerID != 8 {
			t.Fatalf("partnership = %+v", rels)
		}
	})

	t.Run("sign-in without an email stays unconfirmed", func(t *testing.T) {
		db := newMemDB()
		h := &H{DB: db.open(t)}
		rr := signIn(t, h, "", false)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		u := onlyMemUser(t, db)
		if u.verified || !placeholderAppleEmail(u.appleID, u.email) {
			t.Fatalf("user = %+v", u)
		}
	})

	t.Run("later sign-in stores a newly shared address", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{
			id: 8, name: "Bea", email: placeholderAppleAddress("apple-sub-1"),
			appleID: "apple-sub-1", verified: false,
		})
		h := &H{DB: db.open(t)}
		rr := signIn(t, h, "Person@Example.com", true)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		u, _ := db.user(8)
		if u.email != "person@example.com" || !u.verified {
			t.Fatalf("user = %+v", u)
		}
	})

	t.Run("later sign-in leaves another account's address alone", func(t *testing.T) {
		db := newMemDB()
		db.addUser(memUser{
			id: 8, name: "Bea", email: placeholderAppleAddress("apple-sub-1"),
			appleID: "apple-sub-1", verified: false,
		})
		db.addUser(memUser{id: 9, name: "Ada", email: "person@example.com", verified: true})
		h := &H{DB: db.open(t)}
		rr := signIn(t, h, "person@example.com", true)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
		}
		u, _ := db.user(8)
		if u.email != placeholderAppleAddress("apple-sub-1") || u.verified {
			t.Fatalf("stand-in account changed: %+v", u)
		}
		other, _ := db.user(9)
		if other.email != "person@example.com" || other.appleID != "" {
			t.Fatalf("other account changed: %+v", other)
		}
	})
}

func TestPlaceholderAppleAccountCannotAcceptInvite(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	standIn := placeholderAppleAddress("apple-sub-1")
	db.addUser(memUser{id: 2, name: "Bea", email: standIn, appleID: "apple-sub-1", verified: false})
	db.mu.Lock()
	db.pInvites = append(db.pInvites, memPInvite{
		id: 7, inviterID: 1, email: "bea@example.com", token: "raw-token", status: "pending",
	})
	db.mu.Unlock()
	h := &H{DB: db.open(t)}

	listed := httptest.NewRecorder()
	h.ListPendingPartnerInvites(listed, authed(http.MethodGet, "/relationships/invites/pending", "", 2))
	assertCodedError(t, listed, http.StatusForbidden, "email_required")

	req := authed(http.MethodPost, "/relationships/accept-invite", `{"token":"raw-token"}`, 2)
	rr := httptest.NewRecorder()
	h.AcceptPartnerInvite(rr, req)
	assertCodedError(t, rr, http.StatusForbidden, "email_required")
	if len(db.relRows()) != 0 {
		t.Fatalf("relationship created: %+v", db.relRows())
	}
}

func TestResendVerificationForSignedInUser(t *testing.T) {
	mail := &recordingEmail{}
	db := newMemDB()
	db.addUser(memUser{id: 4, name: "Bea", email: "bea@example.com", verified: false})
	db.addUser(memUser{id: 5, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{
		id: 6, name: "Cam", email: placeholderAppleAddress("apple-sub-9"),
		appleID: "apple-sub-9", verified: false,
	})
	h := &H{DB: db.open(t), Email: mail}

	rr := httptest.NewRecorder()
	h.ResendVerification(rr, httptest.NewRequest(http.MethodPost, "/auth/resend-verification", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("signed-out status %d", rr.Code)
	}

	for i := 0; i < verifyMailPerHour; i++ {
		rr = httptest.NewRecorder()
		h.ResendVerification(rr, authed(http.MethodPost, "/auth/resend-verification", "", 4))
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "sent") {
			t.Fatalf("send %d status %d body %s", i+1, rr.Code, rr.Body.String())
		}
	}
	if mail.verify != verifyMailPerHour {
		t.Fatalf("emails sent = %d", mail.verify)
	}
	rr = httptest.NewRecorder()
	h.ResendVerification(rr, authed(http.MethodPost, "/auth/resend-verification", "", 4))
	assertCodedError(t, rr, http.StatusTooManyRequests, "rate_limited")
	if mail.verify != verifyMailPerHour {
		t.Fatalf("limited request still sent mail: %d", mail.verify)
	}
	if rr.Header().Get("Retry-After") == "" {
		t.Fatal("missing Retry-After")
	}

	u, _ := db.user(4)
	u.verifyWindow = time.Now().Add(-2 * time.Hour)
	db.mu.Lock()
	db.users[4].verifyWindow = u.verifyWindow
	db.mu.Unlock()
	rr = httptest.NewRecorder()
	h.ResendVerification(rr, authed(http.MethodPost, "/auth/resend-verification", "", 4))
	if rr.Code != http.StatusOK {
		t.Fatalf("after the hour status %d body %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ResendVerification(rr, authed(http.MethodPost, "/auth/resend-verification", "", 5))
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "already_verified") {
		t.Fatalf("verified status %d body %s", rr.Code, rr.Body.String())
	}

	before := mail.verify
	rr = httptest.NewRecorder()
	h.ResendVerification(rr, authed(http.MethodPost, "/auth/resend-verification", "", 6))
	assertCodedError(t, rr, http.StatusConflict, "email_required")
	if mail.verify != before {
		t.Fatal("stand-in account was sent a confirmation email")
	}
}

func onlyMemUser(t *testing.T, db *memDB) memUser {
	t.Helper()
	db.mu.Lock()
	defer db.mu.Unlock()
	if len(db.users) != 1 {
		t.Fatalf("user count = %d", len(db.users))
	}
	for _, u := range db.users {
		return *u
	}
	return memUser{}
}
