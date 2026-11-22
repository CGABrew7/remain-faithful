package handler

import (
	"sync"
	"time"
)

// inviteMailPerDay is enough for a full group of 12, with room to correct a few addresses.
const inviteMailPerDayDefault = 50

var inviteMailPerDay = inviteMailPerDayDefault

var (
	inviteMailMu    sync.Mutex
	inviteMailDay   string
	inviteMailCount = map[int64]int{}
)

func resetInviteMailForTests() {
	inviteMailMu.Lock()
	defer inviteMailMu.Unlock()
	inviteMailDay = ""
	inviteMailCount = map[int64]int{}
	inviteMailPerDay = inviteMailPerDayDefault
}

// reserveInviteMail counts one invite email for this account on the UTC day.
func reserveInviteMail(userID int64, now time.Time) bool {
	inviteMailMu.Lock()
	defer inviteMailMu.Unlock()
	day := now.UTC().Format("2006-01-02")
	if inviteMailDay != day {
		inviteMailDay = day
		inviteMailCount = map[int64]int{}
	}
	if inviteMailCount[userID] >= inviteMailPerDay {
		return false
	}
	inviteMailCount[userID]++
	return true
}
