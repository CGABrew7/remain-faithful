package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAnyPartnerPINUnlocks(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.addUser(memUser{id: 2, name: "Bea", email: "bea@example.com", verified: true})
	db.addUser(memUser{id: 3, name: "Cam", email: "cam@example.com", verified: true})
	hashA, err := bcrypt.GenerateFromPassword([]byte("1111"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	hashB, err := bcrypt.GenerateFromPassword([]byte("2222"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	db.mu.Lock()
	db.rels = append(db.rels,
		memRel{id: 1, userID: 1, partnerID: 2, status: "accepted", pinHash: string(hashA)},
		memRel{id: 2, userID: 1, partnerID: 3, status: "accepted", pinHash: string(hashB)},
	)
	db.mu.Unlock()

	h := &H{DB: db.open(t), APNS: scriptedAPNS{}, SyncPush: true}
	check := func(pin string) (int, map[string]any) {
		t.Helper()
		rr := httptest.NewRecorder()
		h.VerifyPIN(rr, authed(http.MethodPost, "/protection/pin/verify", `{"pin":"`+pin+`"}`, 1))
		var body map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &body)
		return rr.Code, body
	}

	if code, body := check("1111"); code != http.StatusOK || body["success"] != true {
		t.Fatalf("first partner PIN: %d %v", code, body)
	}
	if code, body := check("2222"); code != http.StatusOK || body["success"] != true {
		t.Fatalf("second partner PIN: %d %v", code, body)
	}
	code, body := check("0000")
	if code != http.StatusOK || body["success"] != false || body["reason"] != "wrong_pin" {
		t.Fatalf("wrong PIN: %d %v", code, body)
	}
}

func TestDeleteUserDropsDeviceTokens(t *testing.T) {
	db := newMemDB()
	db.addUser(memUser{id: 1, name: "Ada", email: "ada@example.com", verified: true})
	db.mu.Lock()
	db.tokens = append(db.tokens, memToken{userID: 1, token: "tok-ada", platform: "ios", environment: "sandbox", active: true})
	db.mu.Unlock()

	h := &H{DB: db.open(t)}
	if err := h.deleteUserByID(context.Background(), 1, "Ada"); err != nil {
		t.Fatal(err)
	}
	if active := db.activeTokenUsers("tok-ada"); len(active) != 0 {
		t.Fatalf("token still active for %v", active)
	}
}
