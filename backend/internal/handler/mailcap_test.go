package handler

import (
	"testing"
	"time"
)

func TestInviteMailDailyCap(t *testing.T) {
	resetInviteMailForTests()
	t.Cleanup(resetInviteMailForTests)
	inviteMailPerDay = 3
	now := time.Date(2026, 11, 22, 15, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		if !reserveInviteMail(42, now) {
			t.Fatalf("invite %d should be allowed", i+1)
		}
	}
	if reserveInviteMail(42, now) {
		t.Fatal("fourth invite should wait until tomorrow")
	}
	if !reserveInviteMail(43, now) {
		t.Fatal("another account should still have its own allowance")
	}
	next := now.Add(24 * time.Hour)
	if !reserveInviteMail(42, next) {
		t.Fatal("the next UTC day should start a new allowance")
	}
}
