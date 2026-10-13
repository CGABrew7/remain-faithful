package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rfauth "remain-faithful/backend/internal/auth"
)

func TestDebugPushEnabled(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	if debugPushEnabled() {
		t.Fatal("debug push must be off when APP_ENV=production")
	}

	t.Setenv("APP_ENV", "staging")
	if !debugPushEnabled() {
		t.Fatal("debug push must be on when APP_ENV=staging")
	}

	t.Setenv("APP_ENV", "development")
	if !debugPushEnabled() {
		t.Fatal("debug push must be on when APP_ENV=development")
	}

	t.Setenv("APP_ENV", "")
	if !debugPushEnabled() {
		t.Fatal("debug push must be on when APP_ENV is empty")
	}
}

func TestDebugPushRouteGatedByAppEnv(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	t.Setenv("APP_ENV", "production")
	router := routes(nil)
	req := httptest.NewRequest(http.MethodPost, "/debug/test-push", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound && rr.Code != http.StatusMethodNotAllowed {
		// gorilla/mux returns 404 for unmatched routes.
		t.Fatalf("production: POST /debug/test-push status = %d, want 404", rr.Code)
	}

	t.Setenv("APP_ENV", "development")
	router = routes(nil)
	req = httptest.NewRequest(http.MethodPost, "/debug/test-push", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	// Route is registered; JWT middleware rejects the request before the handler.
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("development: POST /debug/test-push status = %d, want 401 (route mounted)", rr.Code)
	}
}

func TestBillingPortalRouteRequiresJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	router := routes(nil)
	req := httptest.NewRequest(http.MethodPost, "/donations/billing-portal", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("POST /donations/billing-portal status = %d, want 401", rr.Code)
	}
}

func TestFixedWindowRateLimiter(t *testing.T) {
	lim := fixedWindowRateLimiter(3, func(*http.Request) string { return "test" })
	h := lim(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for i := 0; i < 3; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("req %d: status %d, want 204", i+1, rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("4th req: status %d, want 429", rr.Code)
	}
	if rr.Header().Get("Retry-After") != "60" {
		t.Fatalf("Retry-After = %q, want 60", rr.Header().Get("Retry-After"))
	}
}

func TestAuthenticatedRateLimiterKeysByUser(t *testing.T) {
	lim := authenticatedRateLimiter(1)
	h := lim(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	reqFor := func(uid int64) *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		return req.WithContext(rfauth.ContextWithUser(req.Context(), uid))
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, reqFor(1))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("user 1 first: status %d, want 204", rr.Code)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, reqFor(1))
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("user 1 second: status %d, want 429", rr.Code)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, reqFor(2))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("user 2 first: status %d, want 204 (separate bucket)", rr.Code)
	}
}

func TestClientIPIgnoresForwardedForInProduction(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.RemoteAddr = "10.0.0.8:1234"
	req.Header.Set("X-Forwarded-For", "1.2.3.4")

	t.Setenv("APP_ENV", "production")
	if got := clientIP(req); got == "1.2.3.4" {
		t.Fatalf("production trusted X-Forwarded-For: %q", got)
	}
	req.Header.Set("Fly-Client-IP", "9.9.9.9")
	if got := clientIP(req); got != "9.9.9.9" {
		t.Fatalf("Fly-Client-IP = %q", got)
	}

	t.Setenv("APP_ENV", "development")
	req.Header.Del("Fly-Client-IP")
	if got := clientIP(req); got != "1.2.3.4" {
		t.Fatalf("development X-Forwarded-For = %q", got)
	}
}
