package clientaddr

import (
	"crypto/subtle"
	"net"
	"net/http"
	"os"
	"strings"
)

// Header names the website sends when it calls the API on a visitor's behalf.
// The API uses the address only when the shared secret matches.
const (
	SecretHeader = "X-RF-Proxy-Secret"
	IPHeader     = "X-RF-Client-IP"
)

// FromRequest is the address used for rate limits.
// When SITE_PROXY_SECRET is set and the matching header is present, the
// forwarded visitor address is used. Otherwise the previous rules apply:
// Fly-Client-IP, and outside production the first X-Forwarded-For hop.
func FromRequest(r *http.Request) string {
	if ip := trustedVisitor(r); ip != "" {
		return ip
	}
	if ip := strings.TrimSpace(r.Header.Get("Fly-Client-IP")); ip != "" {
		return ip
	}
	if os.Getenv("APP_ENV") != "production" {
		if hop := firstHop(r.Header.Get("X-Forwarded-For")); hop != "" {
			return hop
		}
	}
	host := r.RemoteAddr
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

func trustedVisitor(r *http.Request) string {
	secret := os.Getenv("SITE_PROXY_SECRET")
	if secret == "" {
		return ""
	}
	got := r.Header.Get(SecretHeader)
	if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
		return ""
	}
	return firstHop(r.Header.Get(IPHeader))
}

func firstHop(raw string) string {
	raw = strings.TrimSpace(strings.SplitN(raw, ",", 2)[0])
	if raw == "" || len(raw) > 64 {
		return ""
	}
	if net.ParseIP(raw) == nil {
		return ""
	}
	return raw
}
