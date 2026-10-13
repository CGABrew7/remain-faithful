package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignIncludesTokenVersionAndParseRoundTrips(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-version")
	raw, err := Sign(7, "a@b.com", 3)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 7 || claims.Email != "a@b.com" || claims.TokenVersion != 3 {
		t.Fatalf("claims = %+v", claims)
	}
	if claims.Issuer != issuer {
		t.Fatalf("issuer = %q", claims.Issuer)
	}
}

func TestParseRejectsMissingVersionFieldAsZero(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-version")
	claims := Claims{
		UserID: 1,
		Email:  "a@b.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret())
	if err != nil {
		t.Fatal(err)
	}
	got, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.TokenVersion != 0 {
		t.Fatalf("legacy token version = %d, want 0", got.TokenVersion)
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-version")
	claims := Claims{
		UserID:       1,
		Email:        "a@b.com",
		TokenVersion: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    issuer,
		},
	}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(raw); err == nil {
		t.Fatal("expected expired token to fail")
	}
}
