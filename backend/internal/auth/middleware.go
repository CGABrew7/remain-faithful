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

// Middleware validates the Bearer token on every protected request and injects
// the authenticated user's ID and email into the request context.
func Middleware(next http.Handler) http.Handler {
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
		ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
		ctx = context.WithValue(ctx, ctxEmail, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext extracts the authenticated user's ID injected by Middleware.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ctxUserID).(int64)
	return id, ok
}
