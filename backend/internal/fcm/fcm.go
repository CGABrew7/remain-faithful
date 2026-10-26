// Package fcm sends Android pushes through the Firebase Cloud Messaging
// HTTP v1 API. It does not read or use the APNs key. When credentials are
// unset, New returns a no-op client so iOS keeps working unchanged.
package fcm

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultTokenURL = "https://oauth2.googleapis.com/token"
	defaultSendHost = "https://fcm.googleapis.com"
	messagingScope  = "https://www.googleapis.com/auth/firebase.messaging"
	tokenTTL        = 50 * time.Minute
)

// ErrInvalidToken means FCM rejected the registration token.
type ErrInvalidToken struct {
	Reason string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("fcm: invalid device token: %s", e.Reason)
}

// Message is one Android push. Data values are strings, matching FCM rules.
// Title and body are display text. Data carries the same routing keys iOS
// reads beside the aps dictionary.
type Message struct {
	Token       string
	Title       string
	Body        string
	CollapseKey string
	Data        map[string]string
	ChannelID   string
}

// Client talks to FCM. A no-op client returns nil from Send.
type Client struct {
	noop        bool
	projectID   string
	clientEmail string
	key         *rsa.PrivateKey
	http        *http.Client
	tokenURL    string
	sendHost    string

	mu         sync.Mutex
	token      string
	tokenMined time.Time
}

// New builds a client from a Firebase service account.
// Empty projectID, clientEmail, or privateKeyPEM yields a no-op client.
func New(projectID, clientEmail, privateKeyPEM string) (*Client, error) {
	projectID = strings.TrimSpace(projectID)
	clientEmail = strings.TrimSpace(clientEmail)
	privateKeyPEM = normalizePEM(privateKeyPEM)
	if projectID == "" || clientEmail == "" || privateKeyPEM == "" {
		fmt.Println("[fcm] WARNING: FCM credentials are not set. Android push is disabled (no-op client)")
		return &Client{noop: true}, nil
	}
	key, err := parseRSAKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[fcm] client ready: project=%s\n", projectID)
	return &Client{
		projectID:   projectID,
		clientEmail: clientEmail,
		key:         key,
		http:        &http.Client{Timeout: 10 * time.Second},
		tokenURL:    defaultTokenURL,
		sendHost:    defaultSendHost,
	}, nil
}

// IsNoop reports whether Send will skip the network.
func (c *Client) IsNoop() bool {
	if c == nil {
		return true
	}
	return c.noop
}

// Send delivers msg. It is safe to call from multiple goroutines.
func (c *Client) Send(ctx context.Context, msg *Message) error {
	if c == nil || c.noop {
		token := ""
		if msg != nil {
			token = msg.Token
		}
		fmt.Printf("[fcm] noop: would send to token %.8s... (FCM not configured)\n", token)
		return nil
	}
	if msg == nil || msg.Token == "" {
		return fmt.Errorf("fcm: token is required")
	}
	access, err := c.accessToken(ctx)
	if err != nil {
		return err
	}
	body, err := json.Marshal(buildRequest(msg))
	if err != nil {
		return fmt.Errorf("fcm: marshal: %w", err)
	}
	url := fmt.Sprintf("%s/v1/projects/%s/messages:send", c.sendHost, c.projectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("fcm: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+access)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("fcm: send: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("[fcm] delivered to token %.8s...\n", msg.Token)
		return nil
	}
	reason := fcmErrorCode(respBody)
	if reason == "UNREGISTERED" || reason == "INVALID_ARGUMENT" || resp.StatusCode == http.StatusNotFound {
		fmt.Printf("[fcm] invalid token %.8s... reason=%s\n", msg.Token, reason)
		return &ErrInvalidToken{Reason: reason}
	}
	fmt.Printf("[fcm] status=%d token=%.8s... reason=%s\n", resp.StatusCode, msg.Token, reason)
	return fmt.Errorf("fcm: unexpected status %d", resp.StatusCode)
}

func buildRequest(msg *Message) map[string]any {
	data := map[string]string{}
	for k, v := range msg.Data {
		if k == "" {
			continue
		}
		data[k] = v
	}
	if msg.Title != "" {
		data["title"] = msg.Title
	}
	if msg.Body != "" {
		data["body"] = msg.Body
	}
	channel := msg.ChannelID
	if channel == "" {
		channel = "alerts"
	}
	note := map[string]any{
		"title": msg.Title,
		"body":  msg.Body,
	}
	androidNote := map[string]any{
		"channel_id": channel,
	}
	message := map[string]any{
		"token":        msg.Token,
		"notification": note,
		"data":         data,
		"android": map[string]any{
			"priority":     "HIGH",
			"notification": androidNote,
		},
	}
	if msg.CollapseKey != "" {
		message["android"].(map[string]any)["collapse_key"] = msg.CollapseKey
	}
	return map[string]any{"message": message}
}

func (c *Client) accessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	if c.token != "" && time.Since(c.tokenMined) < tokenTTL {
		tok := c.token
		c.mu.Unlock()
		return tok, nil
	}
	c.mu.Unlock()

	assertion, err := c.signedJWT()
	if err != nil {
		return "", err
	}
	form := "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&assertion=" + assertion
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL, strings.NewReader(form))
	if err != nil {
		return "", fmt.Errorf("fcm: token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("fcm: token exchange: %w", err)
	}
	defer resp.Body.Close()
	var parsed struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("fcm: token decode: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || parsed.AccessToken == "" {
		return "", fmt.Errorf("fcm: token endpoint status %d", resp.StatusCode)
	}
	c.mu.Lock()
	c.token = parsed.AccessToken
	c.tokenMined = time.Now()
	c.mu.Unlock()
	return parsed.AccessToken, nil
}

func (c *Client) signedJWT() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   c.clientEmail,
		"scope": messagingScope,
		"aud":   defaultTokenURL,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := t.SignedString(c.key)
	if err != nil {
		return "", fmt.Errorf("fcm: sign jwt: %w", err)
	}
	return signed, nil
}

func parseRSAKey(pemText string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemText))
	if block == nil {
		return nil, fmt.Errorf("fcm: failed to decode PEM block")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("fcm: expected RSA private key")
		}
		return rsaKey, nil
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fcm: parse private key: %w", err)
	}
	return key, nil
}

func normalizePEM(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, `\n`, "\n")
	return raw
}

func fcmErrorCode(body []byte) string {
	var parsed struct {
		Error struct {
			Status  string `json:"status"`
			Details []struct {
				ErrorCode string `json:"errorCode"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return ""
	}
	for _, d := range parsed.Error.Details {
		if d.ErrorCode != "" {
			return d.ErrorCode
		}
	}
	return parsed.Error.Status
}
