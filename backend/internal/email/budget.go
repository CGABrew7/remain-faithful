package email

import (
	"errors"
	"sync"
	"time"
)

// ErrDailyCeiling is returned when this process has sent its mail for the UTC day.
var ErrDailyCeiling = errors.New("email: daily send ceiling reached")

// globalMailPerDay covers invite, confirmation, reset, contact, and weekly mail.
const globalMailPerDayDefault = 1000

var globalMailPerDay = globalMailPerDayDefault

var (
	globalMailMu    sync.Mutex
	globalMailDay   string
	globalMailCount int
)

func resetGlobalSendForTests() {
	globalMailMu.Lock()
	defer globalMailMu.Unlock()
	globalMailDay = ""
	globalMailCount = 0
	globalMailPerDay = globalMailPerDayDefault
}

func reserveGlobalSend(now time.Time) bool {
	globalMailMu.Lock()
	defer globalMailMu.Unlock()
	day := now.UTC().Format("2006-01-02")
	if globalMailDay != day {
		globalMailDay = day
		globalMailCount = 0
	}
	if globalMailCount >= globalMailPerDay {
		return false
	}
	globalMailCount++
	return true
}
