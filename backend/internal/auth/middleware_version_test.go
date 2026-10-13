package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareWithVersionRejectsStaleSession(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-middleware")
	raw, err := Sign(4, "a@b.com", 2)
	if err != nil {
		t.Fatal(err)
	}
	called := false
	h := MiddlewareWithVersion(func(userID int64, claimVersion int) bool {
		if userID != 4 || claimVersion != 2 {
			t.Fatalf("callback got user %d version %d", userID, claimVersion)
		}
		return false
	})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
	if called {
		t.Fatal("handler ran for a revoked session")
	}
}

func TestMiddlewareWithVersionAllowsCurrentSession(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-middleware")
	raw, err := Sign(4, "a@b.com", 1)
	if err != nil {
		t.Fatal(err)
	}
	h := MiddlewareWithVersion(func(int64, int) bool { return true })(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rr.Code)
	}
}
