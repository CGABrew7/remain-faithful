package admin

import (
	"os"
	"testing"
	"time"

	rfauth "remain-faithful/backend/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestEnabledRequiresHash(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD_HASH", "")
	if Enabled() {
		t.Fatal("empty hash should disable admin")
	}
	t.Setenv("ADMIN_PASSWORD_HASH", "   ")
	if Enabled() {
		t.Fatal("blank hash should disable admin")
	}
	t.Setenv("ADMIN_PASSWORD_HASH", "$2a$04$abcdefghijklmnopqrstuu012345678901234567890123456")
	if !Enabled() {
		t.Fatal("a hash should enable admin")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-horse"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADMIN_PASSWORD_HASH", string(hash))
	if CheckPassword("wrong-password") {
		t.Fatal("wrong password was accepted")
	}
	if CheckPassword("") {
		t.Fatal("empty password was accepted")
	}
	if !CheckPassword("correct-horse") {
		t.Fatal("right password was rejected")
	}

	t.Setenv("ADMIN_PASSWORD_HASH", "")
	if CheckPassword("correct-horse") {
		t.Fatal("password matched with no hash configured")
	}
}

func TestLoginRateLimit(t *testing.T) {
	ResetLoginLimits()
	t.Cleanup(ResetLoginLimits)
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	for i := 0; i < loginAttempts; i++ {
		if !AllowLogin("203.0.113.9", now) {
			t.Fatalf("attempt %d should be allowed", i+1)
		}
	}
	if AllowLogin("203.0.113.9", now.Add(time.Minute)) {
		t.Fatal("6th attempt should be blocked")
	}
	if !AllowLogin("203.0.113.10", now) {
		t.Fatal("a different IP should have its own window")
	}
	if !AllowLogin("203.0.113.9", now.Add(loginWindow+time.Second)) {
		t.Fatal("window should reset after 15 minutes")
	}
}

func TestSessionExpiryAndSeparateKey(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	t.Setenv("ADMIN_SESSION_SECRET", "")

	tok, err := signSession(sessionTTL, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := ParseSession(tok); err != nil {
		t.Fatal(err)
	}

	expired, err := signSession(-time.Second, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := ParseSession(expired); err == nil {
		t.Fatal("expired session was accepted")
	}

	userTok, err := rfauth.Sign(4, "person@example.com", 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := ParseSession(userTok); err == nil {
		t.Fatal("a user JWT was accepted as an admin session")
	}

	// A token signed with the raw JWT secret must not pass.
	claims := sessionClaims{
		Role: roleOwner,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := raw.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatal(err)
	}
	if err := ParseSession(signed); err == nil {
		t.Fatal("session signed with JWT_SECRET itself was accepted")
	}
}

func TestAdminSessionSecretOverride(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	t.Setenv("ADMIN_SESSION_SECRET", "dedicated-admin-key")
	tok, err := SignSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := ParseSession(tok); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADMIN_SESSION_SECRET", "different-key")
	if err := ParseSession(tok); err == nil {
		t.Fatal("token from another admin key was accepted")
	}
}
