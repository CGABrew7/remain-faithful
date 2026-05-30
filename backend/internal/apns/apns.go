// Package apns provides a thread-safe Apple Push Notification service client
// using token-based (ES256 JWT) HTTP/2 authentication. It relies only on
// golang-jwt/jwt/v5 and the Go standard library; Go's net/http automatically
// negotiates HTTP/2 via ALPN when dialling TLS, so no additional library is
// needed.
package apns

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	hostSandbox    = "https://api.sandbox.push.apple.com"
	hostProduction = "https://api.push.apple.com"

	// tokenTTL is how long a generated JWT is reused before a fresh one is minted.
	tokenTTL = 50 * time.Minute
)

// ErrInvalidToken is returned when APNs responds with HTTP 410 (Unregistered
// or BadDeviceToken), signalling that the device token is no longer valid.
type ErrInvalidToken struct {
	Reason string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("apns: invalid device token: %s", e.Reason)
}

// Notification carries everything needed to build a single APNs request.
type Notification struct {
	// DeviceToken is the hex-encoded APNs device token.
	DeviceToken string
	// Topic is the APNs topic (usually the app's bundle ID). If empty, the
	// client's configured bundleID is used.
	Topic string
	// PushType should be "alert" or "background".
	PushType string
	// Priority is 10 (immediate) or 5 (power-conscious).
	Priority int
	// CollapseID groups notifications so only the most recent is delivered.
	CollapseID string
	// Payload is serialised to JSON and used as the request body.
	Payload any
}

// Client is a thread-safe APNs provider client. A zero-value Client returned
// when keyPEM is empty acts as a no-op: Send always returns nil.
type Client struct {
	noop bool

	keyID    string
	teamID   string
	bundleID string
	key      *ecdsa.PrivateKey
	host     string
	http     *http.Client

	mu         sync.Mutex
	token      string
	tokenMined time.Time
}

// New constructs an APNs client.
//
// If keyPEM is empty a no-op client is returned without error, which is useful
// when APNs credentials are not configured in a given environment.
func New(keyID, teamID, bundleID, keyPEM string, production bool) (*Client, error) {
	if keyPEM == "" {
		return &Client{noop: true}, nil
	}

	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, fmt.Errorf("apns: failed to decode PEM block")
	}

	raw, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("apns: parse PKCS8 key: %w", err)
	}
	ecKey, ok := raw.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("apns: expected EC private key, got %T", raw)
	}

	host := hostSandbox
	if production {
		host = hostProduction
	}

	return &Client{
		keyID:    keyID,
		teamID:   teamID,
		bundleID: bundleID,
		key:      ecKey,
		host:     host,
		http:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers n to APNs. It is safe to call from multiple goroutines.
//
// Returns *ErrInvalidToken when APNs signals the device token is no longer
// valid (HTTP 410).
func (c *Client) Send(ctx context.Context, n *Notification) error {
	if c.noop {
		return nil
	}

	tok, err := c.bearerToken()
	if err != nil {
		return fmt.Errorf("apns: generate token: %w", err)
	}

	body, err := json.Marshal(n.Payload)
	if err != nil {
		return fmt.Errorf("apns: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/3/device/%s", c.host, n.DeviceToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("apns: build request: %w", err)
	}

	topic := n.Topic
	if topic == "" {
		topic = c.bundleID
	}
	pushType := n.PushType
	if pushType == "" {
		pushType = "alert"
	}
	priority := n.Priority
	if priority == 0 {
		priority = 10
	}

	req.Header.Set("Authorization", "bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apns-push-type", pushType)
	req.Header.Set("apns-topic", topic)
	req.Header.Set("apns-priority", fmt.Sprintf("%d", priority))
	if n.CollapseID != "" {
		req.Header.Set("apns-collapse-id", n.CollapseID)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("apns: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// Parse APNs error response.
	var apnsErr struct {
		Reason    string `json:"reason"`
		Timestamp int64  `json:"timestamp"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&apnsErr)

	if resp.StatusCode == http.StatusGone {
		return &ErrInvalidToken{Reason: apnsErr.Reason}
	}
	return fmt.Errorf("apns: unexpected status %d: %s", resp.StatusCode, apnsErr.Reason)
}

// bearerToken returns a cached JWT, regenerating it if it is older than tokenTTL.
func (c *Client) bearerToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Since(c.tokenMined) < tokenTTL {
		return c.token, nil
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": c.teamID,
		"iat": now.Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	t.Header["kid"] = c.keyID

	signed, err := t.SignedString(c.key)
	if err != nil {
		return "", err
	}

	c.token = signed
	c.tokenMined = now
	return signed, nil
}
