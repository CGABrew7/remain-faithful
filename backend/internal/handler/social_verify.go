package handler

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	appleIssuer  = "https://appleid.apple.com"
	appleKeysURL = "https://appleid.apple.com/auth/keys"
)

// timeNow is replaceable in tests.
var timeNow = time.Now

// appleKeySource loads Apple's current signing keys, keyed by kid.
// Tests replace it. Production fetches and caches the JWKS.
type appleKeySource func(ctx context.Context) (map[string]*rsa.PublicKey, error)

var loadAppleKeys appleKeySource = cachedAppleKeys

var appleHTTP = &http.Client{Timeout: 5 * time.Second}

var (
	appleKeysMu     sync.Mutex
	appleKeysCached map[string]*rsa.PublicKey
	appleKeysUntil  time.Time
)

func cachedAppleKeys(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	appleKeysMu.Lock()
	defer appleKeysMu.Unlock()
	if time.Now().Before(appleKeysUntil) && len(appleKeysCached) > 0 {
		return appleKeysCached, nil
	}
	keys, err := fetchAppleKeys(ctx)
	if err != nil {
		if len(appleKeysCached) > 0 {
			return appleKeysCached, nil
		}
		return nil, err
	}
	appleKeysCached = keys
	appleKeysUntil = time.Now().Add(time.Hour)
	return keys, nil
}

func fetchAppleKeys(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appleKeysURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := appleHTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple jwks status %d", resp.StatusCode)
	}
	return parseAppleJWKS(body)
}

type appleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func parseAppleJWKS(raw []byte) (map[string]*rsa.PublicKey, error) {
	var doc struct {
		Keys []appleJWK `json:"keys"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	out := make(map[string]*rsa.PublicKey, len(doc.Keys))
	for _, k := range doc.Keys {
		if k.Kty != "RSA" || k.Kid == "" || k.N == "" || k.E == "" {
			continue
		}
		pub, err := jwkToRSA(k.N, k.E)
		if err != nil {
			continue
		}
		out[k.Kid] = pub
	}
	if len(out) == 0 {
		return nil, errors.New("apple jwks contained no rsa keys")
	}
	return out, nil
}

func jwkToRSA(nB64, eB64 string) (*rsa.PublicKey, error) {
	nb, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, err
	}
	eb, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nb)
	e := new(big.Int).SetBytes(eb)
	if !e.IsInt64() || e.Int64() <= 0 || e.Int64() > int64(^uint(0)>>1) {
		return nil, errors.New("invalid rsa exponent")
	}
	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}

func appleAudience() string {
	for _, key := range []string{"APPLE_AUDIENCE", "APPLE_BUNDLE_ID", "APNS_BUNDLE_ID"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return "com.remainfaithful.app"
}

type appleIdentity struct {
	Sub           string
	Email         string
	EmailVerified bool
}

type appleTokenClaims struct {
	jwt.RegisteredClaims
	Email         any `json:"email"`
	EmailVerified any `json:"email_verified"`
}

// verifyAppleIdentityToken checks the signature against Apple's JWKS and
// requires iss, aud, and exp. email may be empty on later sign-ins.
func verifyAppleIdentityToken(ctx context.Context, raw, audience string, now time.Time) (appleIdentity, error) {
	if strings.TrimSpace(raw) == "" {
		return appleIdentity{}, errors.New("missing identity token")
	}
	keys, err := loadAppleKeys(ctx)
	if err != nil {
		return appleIdentity{}, fmt.Errorf("apple keys: %w", err)
	}
	claims := &appleTokenClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(appleIssuer),
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	tok, err := parser.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		pub, ok := keys[kid]
		if !ok {
			return nil, errors.New("unknown kid")
		}
		return pub, nil
	})
	if err != nil || tok == nil || !tok.Valid {
		if err == nil {
			err = errors.New("invalid identity token")
		}
		return appleIdentity{}, err
	}
	if claims.Subject == "" {
		return appleIdentity{}, errors.New("missing sub claim")
	}
	return appleIdentity{
		Sub:           claims.Subject,
		Email:         stringClaim(claims.Email),
		EmailVerified: boolish(claims.EmailVerified),
	}, nil
}

func stringClaim(v any) string {
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func boolish(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(t, "true")
	default:
		return false
	}
}

// googleTokenInfo posts an ID token to Google's tokeninfo endpoint.
// Tests replace it. The token is a form body, not a query string.
type googleTokenInfo func(ctx context.Context, idToken string) (map[string]any, error)

var postGoogleTokenInfo googleTokenInfo = liveGoogleTokenInfo

func liveGoogleTokenInfo(ctx context.Context, idToken string) (map[string]any, error) {
	form := url.Values{}
	form.Set("id_token", idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/tokeninfo", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := appleHTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tokeninfo status %d", resp.StatusCode)
	}
	var info map[string]any
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return info, nil
}

func googleAudiences() []string {
	raw := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_IDS"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	}
	if raw == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(part); s != "" {
			out = append(out, s)
		}
	}
	return out
}

type googleIdentity struct {
	Sub           string
	Email         string
	Name          string
	EmailVerified bool
}

func verifyGoogleIDToken(ctx context.Context, idToken string, audiences []string, now time.Time) (googleIdentity, error) {
	if strings.TrimSpace(idToken) == "" {
		return googleIdentity{}, errors.New("missing id token")
	}
	if len(audiences) == 0 {
		return googleIdentity{}, errors.New("google client id is not configured")
	}
	info, err := postGoogleTokenInfo(ctx, idToken)
	if err != nil {
		return googleIdentity{}, err
	}
	sub := stringClaim(info["sub"])
	if sub == "" {
		return googleIdentity{}, errors.New("missing sub")
	}
	aud := stringClaim(info["aud"])
	azp := stringClaim(info["azp"])
	if !audienceMatches(audiences, aud) && !audienceMatches(audiences, azp) {
		return googleIdentity{}, errors.New("token audience does not match this app")
	}
	if exp := stringClaim(info["exp"]); exp != "" {
		var sec int64
		if _, err := fmt.Sscan(exp, &sec); err == nil {
			if time.Unix(sec, 0).Before(now) {
				return googleIdentity{}, errors.New("token expired")
			}
		}
	}
	verified := boolish(info["email_verified"])
	email := strings.ToLower(stringClaim(info["email"]))
	if email != "" && !verified {
		return googleIdentity{}, errors.New("google email is not verified")
	}
	return googleIdentity{
		Sub:           sub,
		Email:         email,
		Name:          stringClaim(info["name"]),
		EmailVerified: verified,
	}, nil
}

func audienceMatches(want []string, got string) bool {
	if got == "" {
		return false
	}
	for _, w := range want {
		if w == got {
			return true
		}
	}
	return false
}
