package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	cookieName = "rf_admin"
	issuer     = "remain-faithful-admin"
	roleOwner  = "owner"
	sessionTTL = 8 * time.Hour
)

type sessionClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Enabled reports whether admin routes should answer.
// An empty ADMIN_PASSWORD_HASH turns the whole surface off.
func Enabled() bool {
	return strings.TrimSpace(os.Getenv("ADMIN_PASSWORD_HASH")) != ""
}

// CheckPassword compares the offered password to ADMIN_PASSWORD_HASH.
// bcrypt.CompareHashAndPassword is a constant-time check against the hash.
func CheckPassword(password string) bool {
	hash := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD_HASH"))
	if hash == "" || password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// SignSession mints an 8 hour admin session.
func SignSession() (string, error) {
	return signSession(sessionTTL, time.Now())
}

// SignSessionTTL mints a session with a custom lifetime. Tests use a negative
// lifetime to prove expired tokens are rejected.
func SignSessionTTL(ttl time.Duration) (string, error) {
	return signSession(ttl, time.Now())
}

func signSession(ttl time.Duration, now time.Time) (string, error) {
	claims := sessionClaims{
		Role: roleOwner,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sessionKey())
}

// ParseSession accepts a token signed with the admin key, not a user JWT.
func ParseSession(tokenStr string) error {
	return parseSession(tokenStr, time.Now())
}

func parseSession(tokenStr string, now time.Time) error {
	if tokenStr == "" {
		return errors.New("missing session")
	}
	parser := jwt.NewParser(
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	token, err := parser.ParseWithClaims(tokenStr, &sessionClaims{}, func(t *jwt.Token) (any, error) {
		return sessionKey(), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*sessionClaims)
	if !ok || !token.Valid || claims.Role != roleOwner {
		return errors.New("invalid session")
	}
	return nil
}

// sessionKey is separate from the user JWT key. ADMIN_SESSION_SECRET wins
// when set. Otherwise the key is an HMAC of JWT_SECRET so it is not the
// same bytes as the user signing secret.
func sessionKey() []byte {
	if s := strings.TrimSpace(os.Getenv("ADMIN_SESSION_SECRET")); s != "" {
		return []byte(s)
	}
	mac := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET")))
	_, _ = mac.Write([]byte("remain-faithful-admin-session-v1"))
	return mac.Sum(nil)
}

// TokenFromRequest reads the bearer token or the admin cookie.
func TokenFromRequest(r *http.Request) string {
	if h := r.Header.Get("Authorization"); len(h) > 7 && strings.EqualFold(h[:7], "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	c, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return c.Value
}

// Authorize returns true when admin is enabled and the request carries a live session.
func Authorize(r *http.Request) bool {
	if !Enabled() {
		return false
	}
	return ParseSession(TokenFromRequest(r)) == nil
}

// SetCookie stores the session. Secure is set only on HTTPS so local
// development still works. SameSite=Strict is safe because the site and
// the API share remainfaithful.com. The website also keeps its own
// first-party cookie and sends this value as a bearer token.
func SetCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/admin",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   requestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearCookie expires the admin cookie.
func ClearCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   requestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
	})
}

func requestIsHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = r.Header.Get("Fly-Forwarded-Proto")
	}
	return strings.EqualFold(strings.TrimSpace(proto), "https")
}
