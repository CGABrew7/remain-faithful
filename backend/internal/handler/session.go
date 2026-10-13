package handler

import (
	"context"
	"time"
)

// TokenVersionAllowed reports whether a JWT's version still matches the user.
// Tokens issued before revocation existed omit the version and decode as 0.
// Those stay valid only while the account has never been bumped past 1, so
// deploying this check does not sign everyone out. A password reset or email
// change increments token_version and invalidates older sessions.
// A database error fails closed.
func (h *H) TokenVersionAllowed(userID int64, claimVersion int) bool {
	if h == nil || h.DB == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var v int
	err := h.DB.QueryRowContext(ctx, `SELECT token_version FROM users WHERE id = $1`, userID).Scan(&v)
	if err != nil {
		return false
	}
	if claimVersion == 0 {
		return v <= 1
	}
	return claimVersion == v
}
