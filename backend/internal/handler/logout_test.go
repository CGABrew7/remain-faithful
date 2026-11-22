package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutBumpsSessionVersion(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 8, name: "Ada", email: "ada@example.com", verified: true, tokenVersion: 1})
	h := &H{DB: db.open(t)}

	if !h.TokenVersionAllowed(8, 1) {
		t.Fatal("session should be current before sign-out")
	}
	rr := httptest.NewRecorder()
	h.Logout(rr, authed(http.MethodPost, "/auth/logout", "", 8))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	if h.TokenVersionAllowed(8, 1) {
		t.Fatal("signed-out session should be refused")
	}
	if !h.TokenVersionAllowed(8, 2) {
		t.Fatal("a new sign-in version should be allowed")
	}
}
