package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const ProductionBaseURL = "https://remain-faithful-api.fly.dev"

type User struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	CreatedAt     string `json:"created_at"`
	EmailVerified bool   `json:"email_verified"`
	PendingEmail  string `json:"pending_email"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RemoteEvent struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Category  string `json:"category"`
	Severity  string `json:"severity"`
	Summary   string `json:"summary"`
	Timestamp string `json:"timestamp"`
}

type RemoteAlert struct {
	ID        int64       `json:"id"`
	EventID   int64       `json:"event_id"`
	Seen      bool        `json:"seen"`
	CreatedAt string      `json:"created_at"`
	Kind      string      `json:"kind"`
	AlertType string      `json:"alert_type"`
	Event     RemoteEvent `json:"event"`
}

type Stats struct {
	StreakDays int      `json:"streak_days"`
	BestStreak int      `json:"best_streak"`
	Week       []bool   `json:"week"`
	WeekStates []string `json:"week_states"`
}

type Relationship struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PartnerID int64  `json:"partner_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	IsPrimary bool   `json:"is_primary"`
	PinSet    bool   `json:"pin_set"`
	Partner   User   `json:"partner"`
}

type Group struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"created_at"`
	Covenant  string   `json:"covenant"`
	Members   []Member `json:"members"`
}

type Member struct {
	UserID      int64  `json:"user_id"`
	Role        string `json:"role"`
	JoinedAt    string `json:"joined_at"`
	User        User   `json:"user"`
	FlagsLast30 int    `json:"flags_last_30"`
	StreakDays  int    `json:"streak_days"`
}

type PartnerInvite struct {
	ID          int64  `json:"id"`
	InviterID   int64  `json:"inviter_id"`
	InviterName string `json:"inviter_name"`
	CreatedAt   string `json:"created_at"`
	Type        string `json:"type"`
}

type GroupInvite struct {
	ID          int64  `json:"id"`
	GroupID     int64  `json:"group_id"`
	GroupName   string `json:"group_name"`
	InviterID   int64  `json:"inviter_id"`
	InviterName string `json:"inviter_name"`
	CreatedAt   string `json:"created_at"`
	Type        string `json:"type"`
}

type APIError struct {
	Status  int
	Message string
	Code    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "Something went wrong. Please try again."
}

func IsNetwork(err error) bool {
	var api *APIError
	if errors.As(err, &api) {
		return false
	}
	return err != nil
}

type Client struct {
	BaseURL string
	HTTP    *http.Client
	Token   func() string
}

func NewClient(baseURL string, token func() string) *Client {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = ProductionBaseURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
		Token:   token,
	}
}

func (c *Client) Register(ctx context.Context, name, email, password string) (User, error) {
	var user User
	err := c.call(ctx, http.MethodPost, "/auth/register", map[string]string{
		"name": name, "email": email, "password": password,
	}, &user, false)
	return user, err
}

func (c *Client) Login(ctx context.Context, email, password string) (AuthResponse, error) {
	var resp AuthResponse
	err := c.call(ctx, http.MethodPost, "/auth/login", map[string]string{
		"email": email, "password": password,
	}, &resp, false)
	return resp, err
}

func (c *Client) ForgotPassword(ctx context.Context, email string) error {
	return c.call(ctx, http.MethodPost, "/auth/forgot-password", map[string]string{"email": email}, nil, false)
}

func (c *Client) Refresh(ctx context.Context) (string, error) {
	var resp struct {
		Token string `json:"token"`
	}
	err := c.call(ctx, http.MethodPost, "/auth/refresh", map[string]string{}, &resp, true)
	return resp.Token, err
}

func (c *Client) ResendVerification(ctx context.Context) (string, error) {
	var resp struct {
		Status string `json:"status"`
	}
	err := c.call(ctx, http.MethodPost, "/auth/resend-verification", map[string]string{}, &resp, true)
	return resp.Status, err
}

func (c *Client) Me(ctx context.Context) (User, error) {
	var user User
	err := c.call(ctx, http.MethodGet, "/users/me", nil, &user, true)
	return user, err
}

func (c *Client) UpdateMe(ctx context.Context, name, email string) (User, error) {
	var user User
	err := c.call(ctx, http.MethodPut, "/users/me", map[string]string{"name": name, "email": email}, &user, true)
	return user, err
}

func (c *Client) DeleteMe(ctx context.Context) error {
	return c.call(ctx, http.MethodDelete, "/users/me", nil, nil, true)
}

func (c *Client) Stats(ctx context.Context) (Stats, error) {
	var stats Stats
	err := c.call(ctx, http.MethodGet, "/users/me/stats", nil, &stats, true)
	return stats, err
}

func (c *Client) ListEvents(ctx context.Context) ([]RemoteEvent, error) {
	var events []RemoteEvent
	err := c.call(ctx, http.MethodGet, "/events", nil, &events, true)
	if events == nil {
		events = []RemoteEvent{}
	}
	return events, err
}

func (c *Client) CreateEvent(ctx context.Context, event OutboundEvent) error {
	body, err := event.Body()
	if err != nil {
		return err
	}
	return c.callRaw(ctx, http.MethodPost, "/events", body, nil, true)
}

func (c *Client) ListAlerts(ctx context.Context) ([]RemoteAlert, error) {
	var alerts []RemoteAlert
	err := c.call(ctx, http.MethodGet, "/alerts", nil, &alerts, true)
	if alerts == nil {
		alerts = []RemoteAlert{}
	}
	return alerts, err
}

func (c *Client) UnreadCount(ctx context.Context) (int, error) {
	var resp struct {
		Unseen int `json:"unseen"`
	}
	err := c.call(ctx, http.MethodGet, "/alerts/count", nil, &resp, true)
	return resp.Unseen, err
}

func (c *Client) MarkAlertsSeen(ctx context.Context) error {
	return c.call(ctx, http.MethodPost, "/alerts/mark-seen", map[string]string{}, nil, true)
}

func (c *Client) MarkDiscussed(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodPatch, "/alerts/"+itoa64(id)+"/discussed", nil, nil, true)
}

func (c *Client) HeartbeatMonitoring(ctx context.Context) error {
	return c.call(ctx, http.MethodPost, "/heartbeat", map[string]string{"screen": "monitoring"}, nil, true)
}

func (c *Client) SendSignals(ctx context.Context, selected, taps int) error {
	return c.call(ctx, http.MethodPost, "/users/me/signals", DesktopSignals(selected, taps), nil, true)
}

func (c *Client) SendProtection(ctx context.Context, n ProtectionNotice) error {
	if err := ValidateProtection(n); err != nil {
		return err
	}
	return c.call(ctx, http.MethodPost, "/protection/alerts", map[string]string{
		"type": n.Type, "detail": n.Detail,
	}, nil, true)
}

func (c *Client) PINStatus(ctx context.Context) (bool, error) {
	var resp struct {
		IsSet bool `json:"is_set"`
	}
	err := c.call(ctx, http.MethodGet, "/protection/pin", nil, &resp, true)
	return resp.IsSet, err
}

type PINResult struct {
	Success bool
	Reason  string
	Message string
}

func (c *Client) VerifyPIN(ctx context.Context, pin string) (PINResult, error) {
	if err := ValidatePIN(pin); err != nil {
		return PINResult{}, err
	}
	raw, status, err := c.exchange(ctx, http.MethodPost, "/protection/pin/verify", map[string]string{"pin": pin}, true)
	if err != nil {
		return PINResult{Reason: "unavailable", Message: PINMessage(false, "unavailable")}, nil
	}
	var body struct {
		Success bool   `json:"success"`
		Reason  string `json:"reason"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(raw, &body)
	if status == http.StatusTooManyRequests || body.Reason == "rate_limited" {
		return PINResult{Success: false, Reason: "rate_limited", Message: PINMessage(false, "rate_limited")}, nil
	}
	if status < 200 || status > 299 {
		msg := body.Error
		if msg == "" {
			msg = PINMessage(false, "unavailable")
		}
		return PINResult{}, &APIError{Status: status, Message: msg}
	}
	reason := body.Reason
	if !body.Success && reason == "" {
		reason = "wrong_pin"
	}
	return PINResult{Success: body.Success, Reason: reason, Message: PINMessage(body.Success, reason)}, nil
}

func (c *Client) SetRelationshipPIN(ctx context.Context, id int64, pin string) error {
	if err := ValidatePIN(pin); err != nil {
		return err
	}
	return c.call(ctx, http.MethodPost, "/relationships/"+itoa64(id)+"/pin", map[string]string{"pin": pin}, nil, true)
}

func (c *Client) RemoveRelationshipPIN(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodDelete, "/relationships/"+itoa64(id)+"/pin", nil, nil, true)
}

func (c *Client) ListRelationships(ctx context.Context) ([]Relationship, error) {
	var rows []Relationship
	err := c.call(ctx, http.MethodGet, "/relationships", nil, &rows, true)
	if rows == nil {
		rows = []Relationship{}
	}
	return rows, err
}

func (c *Client) InvitePartner(ctx context.Context, email string) (string, error) {
	var resp struct {
		Status string `json:"status"`
	}
	err := c.call(ctx, http.MethodPost, "/relationships/invite", map[string]string{"email": email}, &resp, true)
	return resp.Status, err
}

func (c *Client) DeleteRelationship(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodDelete, "/relationships/"+itoa64(id), nil, nil, true)
}

func (c *Client) SetPrimary(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodPut, "/relationships/"+itoa64(id)+"/primary", map[string]string{}, nil, true)
}

func (c *Client) PendingPartnerInvites(ctx context.Context) ([]PartnerInvite, error) {
	var rows []PartnerInvite
	err := c.call(ctx, http.MethodGet, "/relationships/invites/pending", nil, &rows, true)
	if rows == nil {
		rows = []PartnerInvite{}
	}
	return rows, err
}

func (c *Client) AcceptPartnerInvite(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodPost, "/relationships/invites/"+itoa64(id)+"/accept", map[string]string{}, nil, true)
}

func (c *Client) DeclinePartnerInvite(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodPost, "/relationships/invites/"+itoa64(id)+"/decline", map[string]string{}, nil, true)
}

func (c *Client) ListGroups(ctx context.Context) ([]Group, error) {
	var rows []Group
	err := c.call(ctx, http.MethodGet, "/groups", nil, &rows, true)
	if rows == nil {
		rows = []Group{}
	}
	return rows, err
}

func (c *Client) CreateGroup(ctx context.Context, name string) (Group, error) {
	var g Group
	err := c.call(ctx, http.MethodPost, "/groups", map[string]string{"name": name, "covenant": ""}, &g, true)
	return g, err
}

func (c *Client) GetGroup(ctx context.Context, id int64) (Group, error) {
	var g Group
	err := c.call(ctx, http.MethodGet, "/groups/"+itoa64(id), nil, &g, true)
	return g, err
}

func (c *Client) LeaveGroup(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodDelete, "/groups/"+itoa64(id)+"/members/me", nil, nil, true)
}

func (c *Client) LeaveAllGroups(ctx context.Context) error {
	return c.call(ctx, http.MethodPost, "/groups/leave-all", map[string]string{}, nil, true)
}

func (c *Client) InviteGroupEmail(ctx context.Context, id int64, email string) (string, error) {
	var resp struct {
		Status string `json:"status"`
	}
	err := c.call(ctx, http.MethodPost, "/groups/"+itoa64(id)+"/email-invite", map[string]string{"email": email}, &resp, true)
	return resp.Status, err
}

func (c *Client) PendingGroupInvites(ctx context.Context) ([]GroupInvite, error) {
	var rows []GroupInvite
	err := c.call(ctx, http.MethodGet, "/groups/invites/pending", nil, &rows, true)
	if rows == nil {
		rows = []GroupInvite{}
	}
	return rows, err
}

func (c *Client) AcceptGroupInvite(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodPost, "/groups/invites/"+itoa64(id)+"/accept", map[string]string{}, nil, true)
}

func (c *Client) DeclineGroupInvite(ctx context.Context, id int64) error {
	return c.call(ctx, http.MethodPost, "/groups/invites/"+itoa64(id)+"/decline", map[string]string{}, nil, true)
}

func (c *Client) Encourage(ctx context.Context, groupID, userID int64) error {
	return c.call(ctx, http.MethodPost, "/groups/"+itoa64(groupID)+"/encouragement", map[string]int64{"user_id": userID}, nil, true)
}

func (c *Client) Feedback(ctx context.Context, fields map[string]string) error {
	return c.call(ctx, http.MethodPost, "/feedback", fields, nil, true)
}

func (c *Client) Panic(ctx context.Context) error {
	return c.call(ctx, http.MethodPost, "/panic", map[string]string{}, nil, true)
}

func (c *Client) call(ctx context.Context, method, path string, in any, out any, auth bool) error {
	var payload []byte
	var err error
	if in != nil {
		payload, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}
	return c.callRaw(ctx, method, path, payload, out, auth)
}

func (c *Client) callRaw(ctx context.Context, method, path string, payload []byte, out any, auth bool) error {
	raw, status, err := c.exchange(ctx, method, path, payload, auth)
	if err != nil {
		return err
	}
	if status < 200 || status > 299 {
		return apiErr(status, raw)
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) exchange(ctx context.Context, method, path string, payload any, auth bool) ([]byte, int, error) {
	var body []byte
	switch p := payload.(type) {
	case nil:
	case []byte:
		body = p
	default:
		var err error
		body, err = json.Marshal(p)
		if err != nil {
			return nil, 0, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, strings.NewReader(string(body)))
	if err != nil {
		return nil, 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "RemainFaithfulDesktop/"+Version)
	if auth {
		tok := ""
		if c.Token != nil {
			tok = c.Token()
		}
		if tok == "" {
			return nil, http.StatusUnauthorized, &APIError{Status: http.StatusUnauthorized, Message: "Not signed in. Please log in again."}
		}
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, 0, err
	}
	return raw, resp.StatusCode, nil
}

func apiErr(status int, raw []byte) error {
	var body struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	_ = json.Unmarshal(raw, &body)
	msg := body.Error
	if msg == "" {
		msg = "Something went wrong. Please try again."
	}
	return &APIError{Status: status, Message: msg, Code: body.Code}
}

func TokenExpiry(token string) (time.Time, bool) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return time.Time{}, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, false
	}
	var payload struct {
		Exp int64 `json:"exp"`
	}
	if json.Unmarshal(raw, &payload) != nil || payload.Exp == 0 {
		return time.Time{}, false
	}
	return time.Unix(payload.Exp, 0), true
}

func itoa64(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
