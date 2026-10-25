package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// memDB is an in-memory fixture for handler tests. It applies the SQL text the
// handlers send, including ON CONFLICT and deletes, so a test can show what a
// statement leaves behind. It is not a Postgres server.
type memDB struct {
	mu        sync.Mutex
	users     map[int64]*memUser
	tokens    []memToken
	alerts    []memAlert
	rels      []memRel
	pInvites  []memPInvite
	gInvites  []memGInvite
	members   []memMember
	groups    map[int64]bool
	resets    []memReset
	devices   []memDevice
	removals  []memRemoval
	waitlist  []memWaitlist
	nextAlert int64
	nextRel   int64
	sqlLog    []string
}

type memUser struct {
	id             int64
	name           string
	email          string
	password       string
	verified       bool
	tokenVersion   int
	pendingEmail   string
	pendingHash    string
	pendingExpires time.Time
	appleID        string
	googleID       string
	screen         string
	heartbeat      bool
	verifyCount    int
	verifyWindow   time.Time
	screenTime     *bool
	appCount       int
	appVersion     string
	supportTaps    int
}

type memToken struct {
	userID      int64
	token       string
	platform    string
	environment string
	active      bool
}

type memAlert struct {
	id        int64
	recipient int64
	kind      string
	alertType string
	summary   string
	subject   int64
}

type memRel struct {
	id        int64
	userID    int64
	partnerID int64
	typ       string
	status    string
	pinHash   string
}

type memPInvite struct {
	id        int64
	inviterID int64
	email     string
	token     string
	status    string
}

type memGInvite struct {
	id        int64
	inviterID int64
	groupID   int64
	email     string
	token     string
	status    string
}

type memMember struct {
	groupID int64
	userID  int64
}

type memReset struct {
	userID  int64
	hash    string
	expires time.Time
}

func newMemDB() *memDB {
	return &memDB{
		users:     map[int64]*memUser{},
		groups:    map[int64]bool{},
		nextAlert: 1,
		nextRel:   1,
	}
}

func (m *memDB) addUser(u memUser) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := u
	if cp.tokenVersion == 0 {
		cp.tokenVersion = 1
	}
	m.users[cp.id] = &cp
}

func (m *memDB) open(t *testing.T) *sql.DB {
	t.Helper()
	db := sql.OpenDB(memConnector{m})
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func (m *memDB) user(id int64) (memUser, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[id]
	if !ok {
		return memUser{}, false
	}
	return *u, true
}

func (m *memDB) activeTokenUsers(token string) []int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ids []int64
	for _, row := range m.tokens {
		if row.token == token && row.active {
			ids = append(ids, row.userID)
		}
	}
	return ids
}

func (m *memDB) alertRows() []memAlert {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]memAlert, len(m.alerts))
	copy(out, m.alerts)
	return out
}

func (m *memDB) relRows() []memRel {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]memRel, len(m.rels))
	copy(out, m.rels)
	return out
}

func (m *memDB) memberRows() []memMember {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]memMember, len(m.members))
	copy(out, m.members)
	return out
}

type memConnector struct{ m *memDB }

func (c memConnector) Connect(context.Context) (driver.Conn, error) {
	return &memConn{m: c.m}, nil
}
func (c memConnector) Driver() driver.Driver { return memDriver{} }

type memDriver struct{}

func (memDriver) Open(string) (driver.Conn, error) { return nil, driver.ErrSkip }

type memConn struct{ m *memDB }

func (c *memConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c *memConn) Close() error                        { return nil }
func (c *memConn) Begin() (driver.Tx, error)           { return memTx{}, nil }
func (c *memConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return memTx{}, nil
}

func (c *memConn) QueryContext(_ context.Context, q string, nv []driver.NamedValue) (driver.Rows, error) {
	cols, rows := c.m.apply(q, namedArgs(nv))
	if cols == nil {
		cols = []string{"x"}
	}
	return &memRows{cols: cols, vals: rows}, nil
}

func (c *memConn) ExecContext(_ context.Context, q string, nv []driver.NamedValue) (driver.Result, error) {
	c.m.apply(q, namedArgs(nv))
	return stubResult{n: 1}, nil
}

type memTx struct{}

func (memTx) Commit() error   { return nil }
func (memTx) Rollback() error { return nil }

type memRows struct {
	cols []string
	vals [][]driver.Value
	i    int
}

func (r *memRows) Columns() []string { return r.cols }
func (r *memRows) Close() error      { return nil }
func (r *memRows) Next(dest []driver.Value) error {
	if r.i >= len(r.vals) {
		return io.EOF
	}
	copy(dest, r.vals[r.i])
	r.i++
	return nil
}

func namedArgs(nv []driver.NamedValue) []any {
	args := make([]any, len(nv))
	for i, v := range nv {
		args[i] = v.Value
	}
	return args
}

func asInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int32:
		return int64(t)
	default:
		return 0
	}
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

func (m *memDB) apply(q string, args []any) ([]string, [][]driver.Value) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sqlLog = append(m.sqlLog, q)
	if cols, rows, ok := m.applyDevice(q, args); ok {
		return cols, rows
	}
	if cols, rows, ok := m.applyWaitlist(q, args); ok {
		return cols, rows
	}

	switch {
	case containsAll(q, "INSERT INTO device_tokens"):
		m.upsertToken(asInt64(args[0]), asString(args[1]), asString(args[2]), asString(args[3]))
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE device_tokens", "user_id <>"):
		token, userID := asString(args[0]), asInt64(args[1])
		for i := range m.tokens {
			if m.tokens[i].token == token && m.tokens[i].userID != userID {
				m.tokens[i].active = false
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE device_tokens", "WHERE user_id"):
		userID := asInt64(args[0])
		for i := range m.tokens {
			if m.tokens[i].userID == userID {
				m.tokens[i].active = false
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE device_tokens", "WHERE token"):
		token := asString(args[0])
		for i := range m.tokens {
			if m.tokens[i].token == token {
				m.tokens[i].active = false
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "DELETE FROM alerts", "subject_user_id"):
		subject := asInt64(args[0])
		kept := m.alerts[:0]
		for _, a := range m.alerts {
			if a.subject != subject {
				kept = append(kept, a)
			}
		}
		m.alerts = kept
		return []string{"ok"}, nil

	case containsAll(q, "DELETE FROM users"):
		id := asInt64(args[0])
		delete(m.users, id)
		for i := range m.alerts {
			if m.alerts[i].subject == id {
				m.alerts[i].subject = 0
			}
		}
		kept := m.alerts[:0]
		for _, a := range m.alerts {
			if a.recipient != id {
				kept = append(kept, a)
			}
		}
		m.alerts = kept
		var tokens []memToken
		for _, row := range m.tokens {
			if row.userID != id {
				tokens = append(tokens, row)
			}
		}
		m.tokens = tokens
		return []string{"ok"}, nil

	case containsAll(q, "DELETE FROM password_reset_tokens", "RETURNING user_id"):
		hash := asString(args[0])
		now := time.Now()
		var kept []memReset
		var userID int64
		found := false
		for _, row := range m.resets {
			if !found && row.hash == hash && row.expires.After(now) {
				userID = row.userID
				found = true
				continue
			}
			kept = append(kept, row)
		}
		m.resets = kept
		if !found {
			return []string{"user_id"}, nil
		}
		return []string{"user_id"}, [][]driver.Value{{userID}}

	case containsAll(q, "DELETE FROM password_reset_tokens"):
		id := asInt64(args[0])
		var kept []memReset
		for _, row := range m.resets {
			if row.userID != id {
				kept = append(kept, row)
			}
		}
		m.resets = kept
		return []string{"ok"}, nil

	case containsAll(q, "INSERT INTO alerts"):
		m.alerts = append(m.alerts, memAlert{
			id:        m.nextAlert,
			recipient: asInt64(args[0]),
			alertType: asString(args[1]),
			summary:   asString(args[2]),
			subject:   asInt64(args[3]),
			kind:      "protection",
		})
		m.nextAlert++
		return []string{"ok"}, nil

	case containsAll(q, "INSERT INTO users", "RETURNING id, name, email"):
		id := int64(1)
		for existing := range m.users {
			if existing >= id {
				id = existing + 1
			}
		}
		verified := false
		if len(args) > 3 {
			if b, ok := args[3].(bool); ok {
				verified = b
			}
		}
		u := &memUser{
			id: id, name: asString(args[0]), email: asString(args[1]),
			verified: verified, tokenVersion: 1,
		}
		if containsAll(q, "apple_id") {
			u.appleID = asString(args[2])
		}
		if containsAll(q, "google_id") {
			u.googleID = asString(args[2])
		}
		m.users[id] = u
		return []string{"id", "name", "email"}, [][]driver.Value{{u.id, u.name, u.email}}

	case containsAll(q, "INSERT INTO relationships"):
		return m.insertRel(q, args)

	case containsAll(q, "UPDATE relationships", "pin_hash = $1"):
		id := asInt64(args[1])
		for i := range m.rels {
			if m.rels[i].id == id {
				m.rels[i].pinHash = asString(args[0])
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE relationship_invites", "RETURNING inviter_id"):
		next := "accepted"
		var match func(memPInvite) bool
		if containsAll(q, "WHERE token") {
			token := asString(args[0])
			email := asString(args[1])
			if containsAll(q, "declined") {
				next = "declined"
			}
			match = func(inv memPInvite) bool {
				return inv.token == token && strings.EqualFold(inv.email, email) && inv.status == "pending"
			}
		} else {
			id := asInt64(args[0])
			email := ""
			if len(args) > 1 {
				email = asString(args[1])
			}
			if len(args) > 2 {
				next = asString(args[2])
			}
			match = func(inv memPInvite) bool {
				return inv.id == id && inv.status == "pending" && (email == "" || strings.EqualFold(inv.email, email))
			}
		}
		for i := range m.pInvites {
			if match(m.pInvites[i]) {
				m.pInvites[i].status = next
				return []string{"inviter_id"}, [][]driver.Value{{m.pInvites[i].inviterID}}
			}
		}
		return []string{"inviter_id"}, nil

	case containsAll(q, "UPDATE relationship_invites"):
		id := asInt64(args[0])
		status := "accepted"
		if containsAll(q, "declined") {
			status = "declined"
		}
		for i := range m.pInvites {
			if m.pInvites[i].id == id || (len(args) > 0 && m.pInvites[i].token == asString(args[0]) && containsAll(q, "WHERE token")) {
				m.pInvites[i].status = status
			}
		}
		if containsAll(q, "WHERE id") {
			for i := range m.pInvites {
				if m.pInvites[i].id == id {
					m.pInvites[i].status = status
				}
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE group_invites"):
		id := asInt64(args[0])
		status := "accepted"
		if containsAll(q, "'declined'") {
			status = "declined"
		}
		var found bool
		var groupID int64
		for i := range m.gInvites {
			if m.gInvites[i].id != id || m.gInvites[i].status != "pending" {
				continue
			}
			if containsAll(q, "invitee_email") && len(args) > 1 && !strings.EqualFold(m.gInvites[i].email, asString(args[1])) {
				continue
			}
			m.gInvites[i].status = status
			groupID = m.gInvites[i].groupID
			found = true
			id = m.gInvites[i].id
			break
		}
		if containsAll(q, "RETURNING") {
			if !found {
				return []string{"id", "group_id"}, nil
			}
			return []string{"id", "group_id"}, [][]driver.Value{{id, groupID}}
		}
		return []string{"ok"}, nil

	case containsAll(q, "INSERT INTO groups"):
		id := int64(len(m.groups) + 1)
		m.groups[id] = true
		return []string{"id", "created_at"}, [][]driver.Value{{id, "2026-01-01T00:00:00Z"}}

	case containsAll(q, "INSERT INTO group_members"):
		m.members = append(m.members, memMember{groupID: asInt64(args[0]), userID: asInt64(args[1])})
		return []string{"joined_at"}, [][]driver.Value{{"2026-01-01T00:00:00Z"}}

	case containsAll(q, "UPDATE users", "SET email = $1", "email_verified = TRUE"):
		u := m.users[asInt64(args[1])]
		if u != nil && !u.verified && strings.EqualFold(u.email, asString(args[2])) {
			u.email = asString(args[0])
			u.verified = true
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "SET email_verified = TRUE", "email_verified = FALSE"):
		u := m.users[asInt64(args[0])]
		if u != nil && !u.verified && strings.EqualFold(u.email, asString(args[1])) {
			u.verified = true
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "verify_mail_count"):
		u := m.users[asInt64(args[2])]
		if u != nil {
			u.verifyCount = int(asInt64(args[0]))
			if t, ok := args[1].(time.Time); ok {
				u.verifyWindow = t
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "screen_time_granted"):
		u := m.users[asInt64(args[0])]
		if u != nil {
			if b, ok := args[1].(bool); ok {
				u.screenTime = &b
			}
			u.appCount = int(asInt64(args[2]))
			if s := asString(args[3]); s != "" {
				u.appVersion = s
			}
			if taps := int(asInt64(args[4])); taps > u.supportTaps {
				u.supportTaps = taps
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "last_heartbeat_at"):
		u := m.users[asInt64(args[0])]
		if u != nil {
			u.heartbeat = true
			u.screen = asString(args[1])
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "SET email = pending_email"):
		u := m.users[asInt64(args[0])]
		if u != nil && u.pendingEmail == asString(args[1]) {
			u.email = u.pendingEmail
			u.verified = true
			u.pendingEmail = ""
			u.pendingHash = ""
			u.pendingExpires = time.Time{}
			u.tokenVersion++
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "CASE WHEN email_verified"):
		m.linkSocial(q, args, true)
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "password_hash", "token_version"):
		u := m.users[asInt64(args[1])]
		if u != nil {
			u.password = asString(args[0])
			u.tokenVersion++
			if containsAll(q, "pending_email = NULL") {
				u.pendingEmail = ""
				u.pendingHash = ""
				u.pendingExpires = time.Time{}
			}
		}
		return []string{"ok"}, nil

	case containsAll(q, "UPDATE users", "email_verified = TRUE", "WHERE id"):
		m.linkSocial(q, args, false)
		return []string{"ok"}, nil

	case containsAll(q, "SELECT token FROM device_tokens", "environment"):
		userID := asInt64(args[0])
		env := asString(args[1])
		var rows [][]driver.Value
		for _, row := range m.tokens {
			if row.userID == userID && row.active && row.environment == env {
				rows = append(rows, []driver.Value{row.token})
			}
		}
		return []string{"token"}, rows

	case containsAll(q, "FROM device_tokens dt", "dt.token"):
		token := asString(args[0])
		var rows [][]driver.Value
		for _, row := range m.tokens {
			if row.token != token {
				continue
			}
			u := m.users[row.userID]
			name := ""
			if u != nil {
				name = u.name
			}
			rows = append(rows, []driver.Value{row.userID, name})
		}
		return []string{"id", "name"}, rows

	case containsAll(q, "SELECT name FROM users WHERE id"):
		u := m.users[asInt64(args[0])]
		if u == nil {
			return []string{"name"}, nil
		}
		return []string{"name"}, [][]driver.Value{{u.name}}

	case containsAll(q, "verify_mail_count", "lower(email)"):
		email := strings.ToLower(asString(args[0]))
		cols := []string{"id", "name", "email_verified", "apple_id", "verify_mail_count", "verify_mail_window"}
		for _, u := range m.users {
			if strings.EqualFold(u.email, email) {
				var window driver.Value
				if !u.verifyWindow.IsZero() {
					window = u.verifyWindow
				}
				return cols, [][]driver.Value{{u.id, u.name, u.verified, u.appleID, int64(u.verifyCount), window}}
			}
		}
		return cols, nil

	case containsAll(q, "verify_mail_count", "FROM users"):
		u := m.users[asInt64(args[0])]
		cols := []string{"name", "email", "email_verified", "apple_id", "verify_mail_count", "verify_mail_window"}
		if u == nil {
			return cols, nil
		}
		var window driver.Value
		if !u.verifyWindow.IsZero() {
			window = u.verifyWindow
		}
		return cols, [][]driver.Value{{u.name, u.email, u.verified, u.appleID, int64(u.verifyCount), window}}

	case containsAll(q, "email_verified", "COALESCE(google_id", "WHERE id"):
		u := m.users[asInt64(args[0])]
		cols := []string{"email", "email_verified", "apple_id", "google_id"}
		if u == nil {
			return cols, nil
		}
		return cols, [][]driver.Value{{u.email, u.verified, u.appleID, u.googleID}}

	case containsAll(q, "SELECT email, email_verified, COALESCE(apple_id, '')"):
		u := m.users[asInt64(args[0])]
		if u == nil {
			return []string{"email", "email_verified", "apple_id"}, nil
		}
		return []string{"email", "email_verified", "apple_id"}, [][]driver.Value{{u.email, u.verified, u.appleID}}

	case containsAll(q, "SELECT email, email_verified FROM users"):
		u := m.users[asInt64(args[0])]
		if u == nil {
			return []string{"email", "email_verified"}, nil
		}
		return []string{"email", "email_verified"}, [][]driver.Value{{u.email, u.verified}}

	case containsAll(q, "SELECT email FROM users WHERE id"):
		u := m.users[asInt64(args[0])]
		if u == nil {
			return []string{"email"}, nil
		}
		return []string{"email"}, [][]driver.Value{{u.email}}

	case containsAll(q, "SELECT token_version FROM users"):
		u := m.users[asInt64(args[0])]
		if u == nil {
			return []string{"token_version"}, nil
		}
		return []string{"token_version"}, [][]driver.Value{{int64(u.tokenVersion)}}

	case containsAll(q, "WHERE google_id"), containsAll(q, "WHERE apple_id"):
		idVal := asString(args[0])
		for _, u := range m.users {
			if (containsAll(q, "google_id") && u.googleID == idVal) || (containsAll(q, "apple_id") && u.appleID == idVal) {
				return []string{"id", "name", "email", "email_verified"}, [][]driver.Value{{u.id, u.name, u.email, u.verified}}
			}
		}
		return []string{"id", "name", "email", "email_verified"}, nil

	case containsAll(q, "COALESCE(apple_id", "WHERE email"):
		email := asString(args[0])
		for _, u := range m.users {
			if u.email == email {
				return []string{"id", "name", "email", "apple_id", "google_id"}, [][]driver.Value{{
					u.id, u.name, u.email, u.appleID, u.googleID,
				}}
			}
		}
		return []string{"id", "name", "email", "apple_id", "google_id"}, nil

	case containsAll(q, "SELECT id FROM users WHERE email"):
		email := asString(args[0])
		for _, u := range m.users {
			if u.email == email {
				return []string{"id"}, [][]driver.Value{{u.id}}
			}
		}
		return []string{"id"}, nil

	case containsAll(q, "pending_email_token_hash"):
		sum := asString(args[0])
		now := time.Now()
		for _, u := range m.users {
			if u.pendingHash == sum && u.pendingEmail != "" && u.pendingExpires.After(now) {
				return []string{"id", "pending_email"}, [][]driver.Value{{u.id, u.pendingEmail}}
			}
		}
		return []string{"id", "pending_email"}, nil

	case containsAll(q, "FROM email_verifications"):
		return []string{"id", "email"}, nil

	case containsAll(q, "SELECT u.email", "password_reset_tokens"):
		hash := asString(args[0])
		now := time.Now()
		for _, row := range m.resets {
			if row.hash == hash && row.expires.After(now) {
				if u := m.users[row.userID]; u != nil {
					return []string{"email"}, [][]driver.Value{{u.email}}
				}
			}
		}
		return []string{"email"}, nil

	case containsAll(q, "FROM password_reset_tokens"):
		hash := asString(args[0])
		for _, row := range m.resets {
			if row.hash == hash {
				return []string{"user_id", "expires_at"}, [][]driver.Value{{row.userID, row.expires}}
			}
		}
		return []string{"user_id", "expires_at"}, nil

	case containsAll(q, "SELECT user_id, status FROM relationships"):
		id, partnerID := asInt64(args[0]), asInt64(args[1])
		for _, rel := range m.rels {
			if rel.id == id && rel.partnerID == partnerID {
				return []string{"user_id", "status"}, [][]driver.Value{{rel.userID, rel.status}}
			}
		}
		return []string{"user_id", "status"}, nil

	case containsAll(q, "SELECT status FROM relationships"):
		userID, partnerID := asInt64(args[0]), asInt64(args[1])
		for _, rel := range m.rels {
			if rel.userID == userID && rel.partnerID == partnerID {
				return []string{"status"}, [][]driver.Value{{rel.status}}
			}
		}
		return []string{"status"}, nil

	case containsAll(q, "SELECT COUNT(*) FROM relationships"):
		userID := asInt64(args[0])
		var n int64
		for _, rel := range m.rels {
			if rel.userID == userID {
				n++
			}
		}
		return []string{"count"}, [][]driver.Value{{n}}

	case containsAll(q, "SELECT COUNT(*) FROM group_members"):
		groupID := asInt64(args[0])
		var n int64
		for _, row := range m.members {
			if row.groupID == groupID {
				n++
			}
		}
		return []string{"count"}, [][]driver.Value{{n}}

	case containsAll(q, "FROM relationship_invites", "FROM group_invites", "COUNT"):
		return []string{"count"}, [][]driver.Value{{int64(0)}}

	case containsAll(q, "SELECT partner_id FROM relationships"):
		userID := asInt64(args[0])
		var rows [][]driver.Value
		for _, rel := range m.rels {
			if rel.userID == userID && rel.status == "accepted" {
				rows = append(rows, []driver.Value{rel.partnerID})
			}
		}
		return []string{"partner_id"}, rows

	case containsAll(q, "SELECT inviter_id, token FROM relationship_invites"):
		email := asString(args[0])
		var rows [][]driver.Value
		for _, inv := range m.pInvites {
			if inv.email == email && inv.status == "pending" {
				rows = append(rows, []driver.Value{inv.inviterID, inv.token})
			}
		}
		return []string{"inviter_id", "token"}, rows

	case containsAll(q, "SELECT token, inviter_id FROM relationship_invites"):
		id := asInt64(args[0])
		email := asString(args[1])
		for _, inv := range m.pInvites {
			if inv.id == id && inv.email == email && inv.status == "pending" {
				return []string{"token", "inviter_id"}, [][]driver.Value{{inv.token, inv.inviterID}}
			}
		}
		return []string{"token", "inviter_id"}, nil

	case containsAll(q, "SELECT inviter_id, invitee_email FROM relationship_invites"):
		token := asString(args[0])
		for _, inv := range m.pInvites {
			if inv.token == token && inv.status == "pending" {
				return []string{"inviter_id", "invitee_email"}, [][]driver.Value{{inv.inviterID, inv.email}}
			}
		}
		return []string{"inviter_id", "invitee_email"}, nil

	case containsAll(q, "SELECT i.id, i.inviter_id"):
		email := asString(args[0])
		var rows [][]driver.Value
		for _, inv := range m.pInvites {
			if inv.email == email && inv.status == "pending" {
				name := ""
				if u := m.users[inv.inviterID]; u != nil {
					name = u.name
				}
				rows = append(rows, []driver.Value{inv.id, inv.inviterID, name, "2026-01-01T00:00:00Z"})
			}
		}
		return []string{"id", "inviter_id", "name", "created_at"}, rows

	case containsAll(q, "SELECT id, group_id FROM group_invites", "token = $1"):
		token := asString(args[0])
		email := asString(args[1])
		for _, inv := range m.gInvites {
			if inv.token == token && inv.email == email && inv.status == "pending" {
				return []string{"id", "group_id"}, [][]driver.Value{{inv.id, inv.groupID}}
			}
		}
		return []string{"id", "group_id"}, nil

	case containsAll(q, "SELECT id, group_id FROM group_invites"):
		email := asString(args[0])
		var rows [][]driver.Value
		for _, inv := range m.gInvites {
			if inv.email == email && inv.status == "pending" {
				rows = append(rows, []driver.Value{inv.id, inv.groupID})
			}
		}
		return []string{"id", "group_id"}, rows

	case containsAll(q, "SELECT group_id FROM group_invites"):
		id := asInt64(args[0])
		email := asString(args[1])
		for _, inv := range m.gInvites {
			if inv.id == id && inv.email == email && inv.status == "pending" {
				return []string{"group_id"}, [][]driver.Value{{inv.groupID}}
			}
		}
		return []string{"group_id"}, nil

	case containsAll(q, "SELECT i.id, i.group_id"):
		email := asString(args[0])
		var rows [][]driver.Value
		for _, inv := range m.gInvites {
			if inv.email == email && inv.status == "pending" {
				rows = append(rows, []driver.Value{inv.id, inv.groupID, "Group", inv.inviterID, "Inviter", "2026-01-01T00:00:00Z"})
			}
		}
		return []string{"id", "group_id", "name", "inviter_id", "name", "created_at"}, rows

	case containsAll(q, "SELECT id FROM groups"):
		id := asInt64(args[0])
		if m.groups[id] {
			return []string{"id"}, [][]driver.Value{{id}}
		}
		return []string{"id"}, nil

	case containsAll(q, "FROM group_members", "role"):
		return []string{"group_id"}, nil
	}
	return []string{"x"}, nil
}

func (m *memDB) upsertToken(userID int64, token, platform, environment string) {
	for i := range m.tokens {
		if m.tokens[i].userID == userID && m.tokens[i].token == token {
			m.tokens[i].active = true
			m.tokens[i].platform = platform
			m.tokens[i].environment = environment
			return
		}
	}
	m.tokens = append(m.tokens, memToken{
		userID: userID, token: token, platform: platform, environment: environment, active: true,
	})
}

func (m *memDB) insertRel(q string, args []any) ([]string, [][]driver.Value) {
	userID, partnerID := asInt64(args[0]), asInt64(args[1])
	typ := "partner"
	status := "pending"
	if len(args) > 2 {
		if s := asString(args[2]); s != "" && !containsAll(q, "'partner'") {
			typ = s
		}
	}
	if containsAll(q, "'accepted'") || containsAll(q, "SET status = 'accepted'") {
		status = "accepted"
	}
	for i := range m.rels {
		if m.rels[i].userID == userID && m.rels[i].partnerID == partnerID {
			if containsAll(q, "DO NOTHING") || containsAll(q, "relationships.type") {
				return m.relReturning(q, m.rels[i])
			}
			if containsAll(q, "SET type = EXCLUDED.type") {
				m.rels[i].typ = typ
			}
			if containsAll(q, "SET status = 'accepted'") {
				m.rels[i].status = "accepted"
			}
			return m.relReturning(q, m.rels[i])
		}
	}
	if containsAll(q, "DO NOTHING") {
		// No existing row: insert below.
	}
	row := memRel{id: m.nextRel, userID: userID, partnerID: partnerID, typ: typ, status: status}
	m.nextRel++
	m.rels = append(m.rels, row)
	return m.relReturning(q, row)
}

func (m *memDB) relReturning(q string, row memRel) ([]string, [][]driver.Value) {
	if containsAll(q, "RETURNING id, status") {
		return []string{"id", "status", "created_at"}, [][]driver.Value{{row.id, row.status, "2026-01-01T00:00:00Z"}}
	}
	if containsAll(q, "RETURNING id") {
		return []string{"id"}, [][]driver.Value{{row.id}}
	}
	return []string{"ok"}, nil
}

func (m *memDB) linkSocial(q string, args []any, clearUnverifiedPassword bool) {
	idVal := asString(args[0])
	userID := asInt64(args[1])
	u := m.users[userID]
	if u == nil {
		return
	}
	if containsAll(q, "google_id") {
		if u.googleID != "" && u.googleID != idVal {
			return
		}
		u.googleID = idVal
	}
	if containsAll(q, "apple_id") {
		if u.appleID != "" && u.appleID != idVal {
			return
		}
		u.appleID = idVal
	}
	if clearUnverifiedPassword && !u.verified {
		u.password = ""
		u.tokenVersion++
	}
	u.verified = true
}

func containsAll(q string, parts ...string) bool {
	for _, p := range parts {
		if !stringContains(q, p) {
			return false
		}
	}
	return true
}

func stringContains(q, part string) bool {
	return len(part) == 0 || (len(q) >= len(part) && searchString(q, part))
}

func searchString(q, part string) bool {
	for i := 0; i+len(part) <= len(q); i++ {
		if q[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
