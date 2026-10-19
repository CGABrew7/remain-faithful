package admin

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

// ClientIP prefers Fly-Client-IP. Fly's proxy adds that header as the address
// it accepted the connection from and tells apps to prefer it over
// X-Forwarded-For when Fly is the edge:
// https://fly.io/docs/networking/request-headers/
// Production does not fall back to X-Forwarded-For.
func ClientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("Fly-Client-IP")); ip != "" {
		return ip
	}
	if os.Getenv("APP_ENV") != "production" {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
		}
	}
	host := r.RemoteAddr
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}
