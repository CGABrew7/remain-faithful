package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ctxUserID contextKey = "user_id"
	ctxEmail  contextKey = "email"
)

// VersionAllowed reports whether claimVersion is still valid for userID.
// A zero claimVersion is a session issued before revocation existed.
type VersionAllowed func(userID int64, claimVersion int) bool

// Middleware validates the Bearer token on every protected request and injects
// the authenticated user's ID and email into the request context.
// It does not check token revocation. Use MiddlewareWithVersion in production.
func Middleware(next http.Handler) http.Handler {
	return middleware(nil)(next)
}

// MiddlewareWithVersion is Middleware plus a session-revocation check.
func MiddlewareWithVersion(allowed VersionAllowed) func(http.Handler) http.Handler {
	return middleware(allowed)
}

func middleware(allowed VersionAllowed) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing or invalid Authorization header"}`, http.StatusUnauthorized)
				return
			}
			claims, err := Parse(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			if allowed != nil && !allowed(claims.UserID, claims.TokenVersion) {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxEmail, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ContextWithUser injects userID the same way Middleware does. Used by tests.
func ContextWithUser(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, ctxUserID, userID)
}

// UserIDFromContext extracts the authenticated user's ID injected by Middleware.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ctxUserID).(int64)
	return id, ok
}

// EmailFromContext extracts the authenticated user's email injected by Middleware.
func EmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ctxEmail).(string)
	return email, ok
}
