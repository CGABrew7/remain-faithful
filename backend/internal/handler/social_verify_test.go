package handler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func withAppleKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	prev := loadAppleKeys
	loadAppleKeys = func(context.Context) (map[string]*rsa.PublicKey, error) {
		return map[string]*rsa.PublicKey{"test-kid": &key.PublicKey}, nil
	}
	t.Cleanup(func() { loadAppleKeys = prev })
	return key
}

func signApple(t *testing.T, key *rsa.PrivateKey, aud string, exp time.Time, email string, verified bool) string {
	t.Helper()
	claims := appleTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "apple-sub-1",
			Issuer:    appleIssuer,
			Audience:  jwt.ClaimStrings{aud},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(exp.Add(-time.Minute)),
		},
		Email:         email,
		EmailVerified: verified,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = "test-kid"
	raw, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestVerifyAppleRejectsUnsignedPayload(t *testing.T) {
	withAppleKey(t)
	forged := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2lkIn0.eyJzdWIiOiJhdHRhY2tlciIsImVtYWlsIjoidmljdGltQGV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cHM6Ly9hcHBsZWlkLmFwcGxlLmNvbSIsImF1ZCI6ImNvbS5yZW1haW5mYWl0aGZ1bC5hcHAiLCJleHAiOjk5OTk5OTk5OTl9.c2ln"
	_, err := verifyAppleIdentityToken(context.Background(), forged, "com.remainfaithful.app", time.Now())
	if err == nil {
		t.Fatal("forged apple token was accepted")
	}
}

func TestVerifyAppleAcceptsSignedToken(t *testing.T) {
	key := withAppleKey(t)
	now := time.Now()
	raw := signApple(t, key, "com.remainfaithful.app", now.Add(time.Hour), "person@example.com", true)
	got, err := verifyAppleIdentityToken(context.Background(), raw, "com.remainfaithful.app", now)
	if err != nil {
		t.Fatal(err)
	}
	if got.Sub != "apple-sub-1" || got.Email != "person@example.com" || !got.EmailVerified {
		t.Fatalf("identity = %+v", got)
	}
}

func TestVerifyAppleRejectsWrongAudienceAndExpiry(t *testing.T) {
	key := withAppleKey(t)
	now := time.Now()
	wrongAud := signApple(t, key, "com.other.app", now.Add(time.Hour), "person@example.com", true)
	if _, err := verifyAppleIdentityToken(context.Background(), wrongAud, "com.remainfaithful.app", now); err == nil {
		t.Fatal("wrong audience was accepted")
	}
	expired := signApple(t, key, "com.remainfaithful.app", now.Add(-time.Minute), "person@example.com", true)
	if _, err := verifyAppleIdentityToken(context.Background(), expired, "com.remainfaithful.app", now); err == nil {
		t.Fatal("expired token was accepted")
	}
}

func TestVerifyGoogleRequiresAudienceAndVerifiedEmail(t *testing.T) {
	prev := postGoogleTokenInfo
	t.Cleanup(func() { postGoogleTokenInfo = prev })
	now := time.Now()
	postGoogleTokenInfo = func(context.Context, string) (map[string]any, error) {
		return map[string]any{
			"sub":            "g-1",
			"aud":            "other-client",
			"azp":            "other-client",
			"email":          "person@example.com",
			"email_verified": "true",
			"exp":            "9999999999",
			"name":           "Pat",
		}, nil
	}
	if _, err := verifyGoogleIDToken(context.Background(), "tok", []string{"our-client"}, now); err == nil {
		t.Fatal("foreign audience was accepted")
	}

	postGoogleTokenInfo = func(context.Context, string) (map[string]any, error) {
		return map[string]any{
			"sub":            "g-1",
			"aud":            "our-client",
			"email":          "person@example.com",
			"email_verified": false,
			"exp":            "9999999999",
		}, nil
	}
	if _, err := verifyGoogleIDToken(context.Background(), "tok", []string{"our-client"}, now); err == nil {
		t.Fatal("unverified email was accepted")
	}

	postGoogleTokenInfo = func(context.Context, string) (map[string]any, error) {
		return map[string]any{
			"sub":            "g-1",
			"aud":            "our-client",
			"email":          "Person@Example.com",
			"email_verified": "true",
			"exp":            "9999999999",
			"name":           "Pat",
		}, nil
	}
	got, err := verifyGoogleIDToken(context.Background(), "tok", []string{"our-client"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if got.Sub != "g-1" || got.Email != "person@example.com" || !got.EmailVerified || got.Name != "Pat" {
		t.Fatalf("identity = %+v", got)
	}
}

func TestVerifyGoogleRejectsExpiredToken(t *testing.T) {
	prev := postGoogleTokenInfo
	t.Cleanup(func() { postGoogleTokenInfo = prev })
	postGoogleTokenInfo = func(context.Context, string) (map[string]any, error) {
		return map[string]any{
			"sub":            "g-1",
			"aud":            "our-client",
			"email":          "a@b.com",
			"email_verified": true,
			"exp":            "10",
		}, nil
	}
	if _, err := verifyGoogleIDToken(context.Background(), "tok", []string{"our-client"}, time.Unix(100, 0)); err == nil {
		t.Fatal("expired google token was accepted")
	}
}

func TestBoolish(t *testing.T) {
	if !boolish(true) || !boolish("true") || !boolish("TRUE") {
		t.Fatal("expected true")
	}
	if boolish(false) || boolish("false") || boolish(nil) || boolish(1) {
		t.Fatal("expected false")
	}
}
