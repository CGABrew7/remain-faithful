package admin

import (
	"net/http"
	"sync"
	"time"

	"remain-faithful/backend/internal/clientaddr"
)

const (
	loginAttempts = 5
	loginWindow   = 15 * time.Minute
)

type loginBucket struct {
	count int
	start time.Time
}

var (
	loginMu      sync.Mutex
	loginBuckets = map[string]*loginBucket{}
)

// AllowLogin counts one attempt. It returns false once an IP has made
// 5 attempts inside 15 minutes. The counter includes failed and successful tries.
func AllowLogin(ip string, now time.Time) bool {
	if ip == "" {
		ip = "unknown"
	}
	loginMu.Lock()
	defer loginMu.Unlock()
	if len(loginBuckets) > 500 {
		for key, b := range loginBuckets {
			if now.Sub(b.start) >= loginWindow {
				delete(loginBuckets, key)
			}
		}
	}
	b, ok := loginBuckets[ip]
	if !ok || now.Sub(b.start) >= loginWindow {
		loginBuckets[ip] = &loginBucket{count: 1, start: now}
		return true
	}
	if b.count >= loginAttempts {
		return false
	}
	b.count++
	return true
}

// ResetLoginLimits clears in-memory login counters. Tests use it.
func ResetLoginLimits() {
	loginMu.Lock()
	loginBuckets = map[string]*loginBucket{}
	loginMu.Unlock()
}

// ClientIP is the address used for owner sign-in lockout.
// A website call can pass the visitor address when the shared proxy secret matches.
func ClientIP(r *http.Request) string {
	return clientaddr.FromRequest(r)
}
