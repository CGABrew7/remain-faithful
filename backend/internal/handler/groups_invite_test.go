package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	rfauth "remain-faithful/backend/internal/auth"

	"github.com/gorilla/mux"
)

// ── scripted fake database ──────────────────────────────────────────────────
//
// inviteFakeDB answers queries by substring match and records every statement
// so tests can assert which writes happened. It implements just enough of
// database/sql/driver for QueryRowContext, ExecContext and BeginTx.

type inviteFakeDB struct {
	mu    sync.Mutex
	log   []fakeStmt
	query func(q string, args []any) (cols []string, rows [][]driver.Value)
}

type fakeStmt struct {
	sql  string
	args []any
}

func (f *inviteFakeDB) record(q string, nv []driver.NamedValue) []any {
	args := make([]any, len(nv))
	for i, v := range nv {
		args[i] = v.Value
	}
	f.mu.Lock()
	f.log = append(f.log, fakeStmt{sql: q, args: args})
	f.mu.Unlock()
	return args
}

// ran returns the statements whose SQL contains substr.
func (f *inviteFakeDB) ran(substr string) []fakeStmt {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []fakeStmt
	for _, s := range f.log {
		if strings.Contains(s.sql, substr) {
			out = append(out, s)
		}
	}
	return out
}

func (f *inviteFakeDB) open(t *testing.T) *sql.DB {
	t.Helper()
	db := sql.OpenDB(inviteFakeConnector{f})
	t.Cleanup(func() { _ = db.Close() })
	return db
}

type inviteFakeConnector struct{ f *inviteFakeDB }

func (c inviteFakeConnector) Connect(context.Context) (driver.Conn, error) {
	return &inviteFakeConn{f: c.f}, nil
}
func (c inviteFakeConnector) Driver() driver.Driver { return inviteFakeDriver{} }

type inviteFakeDriver struct{}

func (inviteFakeDriver) Open(string) (driver.Conn, error) { return nil, driver.ErrSkip }

type inviteFakeConn struct{ f *inviteFakeDB }

func (c *inviteFakeConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c *inviteFakeConn) Close() error                        { return nil }
func (c *inviteFakeConn) Begin() (driver.Tx, error)           { return inviteFakeTx{c.f}, nil }
func (c *inviteFakeConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.f.record("BEGIN", nil)
	return inviteFakeTx{c.f}, nil
}

func (c *inviteFakeConn) QueryContext(_ context.Context, q string, nv []driver.NamedValue) (driver.Rows, error) {
	args := c.f.record(q, nv)
	cols, rows := c.f.query(q, args)
	if cols == nil {
		cols = []string{"x"}
	}
	return &inviteFakeRows{cols: cols, vals: rows}, nil
}

func (c *inviteFakeConn) ExecContext(_ context.Context, q string, nv []driver.NamedValue) (driver.Result, error) {
	c.f.record(q, nv)
	return stubResult{n: 1}, nil
}

type inviteFakeTx struct{ f *inviteFakeDB }

func (t inviteFakeTx) Commit() error   { t.f.record("COMMIT", nil); return nil }
func (t inviteFakeTx) Rollback() error { return nil }

type inviteFakeRows struct {
	cols []string
	vals [][]driver.Value
	i    int
}

func (r *inviteFakeRows) Columns() []string { return r.cols }
func (r *inviteFakeRows) Close() error      { return nil }
func (r *inviteFakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.vals) {
		return io.EOF
	}
	copy(dest, r.vals[r.i])
	r.i++
	return nil
}

func one(col string, v driver.Value) ([]string, [][]driver.Value) {
	return []string{col}, [][]driver.Value{{v}}
}

func none(col string) ([]string, [][]driver.Value) { return []string{col}, nil }

// ── recording email sender ──────────────────────────────────────────────────

type recordingEmail struct {
	mu          sync.Mutex
	groupInvite []string // acceptURL per call
	groupTo     []string
}

func (e *recordingEmail) SendPasswordReset(string, string, string) error           { return nil }
func (e *recordingEmail) SendPartnerInvite(string, string, string) error           { return nil }
func (e *recordingEmail) SendContact(string, string, string, string, string) error { return nil }
func (e *recordingEmail) SendGroupInvite(to, _, _, acceptURL string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.groupTo = append(e.groupTo, to)
	e.groupInvite = append(e.groupInvite, acceptURL)
	return nil
}

// ── scenario helpers ────────────────────────────────────────────────────────

type inviteScenario struct {
	isMember      bool
	memberCount   int64
	inviteeExists bool
	alreadyMember bool
}

func (s inviteScenario) router() func(q string, args []any) ([]string, [][]driver.Value) {
	return func(q string, _ []any) ([]string, [][]driver.Value) {
		switch {
		case strings.Contains(q, "SELECT role FROM group_members"):
			if s.isMember {
				return one("role", "admin")
			}
			return none("role")
		case strings.Contains(q, "SELECT COUNT(*) FROM group_members"):
			return one("count", s.memberCount)
		case strings.Contains(q, "SELECT name FROM users"):
			return one("name", "Jeff")
		case strings.Contains(q, "SELECT name FROM groups"):
			return one("name", "App Review Group")
		case strings.Contains(q, "SELECT id FROM users WHERE email"):
			if s.inviteeExists {
				return one("id", int64(42))
			}
			return none("id")
		case strings.Contains(q, "INSERT INTO group_members"):
			if s.alreadyMember {
				return none("joined_at")
			}
			return one("joined_at", "2026-10-08T20:00:00Z")
		}
		return none("x")
	}
}

func doInvite(t *testing.T, handler func(*H, http.ResponseWriter, *http.Request), sc inviteScenario, body string) (*httptest.ResponseRecorder, *inviteFakeDB, *recordingEmail) {
	t.Helper()
	t.Setenv("SITE_URL", "")
	f := &inviteFakeDB{query: sc.router()}
	em := &recordingEmail{}
	h := &H{DB: f.open(t), Email: em}
	req := httptest.NewRequest(http.MethodPost, "/groups/7/invite", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()
	handler(h, rr, req)
	return rr, f, em
}

func inviteMember(h *H, w http.ResponseWriter, r *http.Request)     { h.InviteMember(w, r) }
func groupEmailInvite(h *H, w http.ResponseWriter, r *http.Request) { h.GroupEmailInvite(w, r) }

func decodeBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
		t.Fatalf("decode body %q: %v", rr.Body.String(), err)
	}
	return m
}

// ── InviteMember (Group tab) ────────────────────────────────────────────────

func TestInviteMemberNonMemberForbidden(t *testing.T) {
	rr, f, em := doInvite(t, inviteMember, inviteScenario{isMember: false}, `{"user_email":"a@b.com"}`)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "you must be a member of this group to invite others") {
		t.Fatalf("body = %q", rr.Body.String())
	}
	if len(f.ran("INSERT")) != 0 || len(em.groupInvite) != 0 {
		t.Fatalf("non-member must not write or email")
	}
}

func TestInviteMemberUnknownEmailCreatesPendingInvite(t *testing.T) {
	rr, f, em := doInvite(t, inviteMember, inviteScenario{isMember: true, memberCount: 2},
		`{"user_email":"  New.Person@Example.com "}`)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %q", rr.Code, rr.Body.String())
	}
	if got := decodeBody(t, rr)["status"]; got != "invited" {
		t.Fatalf("status field = %v, want invited", got)
	}
	ins := f.ran("INSERT INTO group_invites")
	if len(ins) != 1 {
		t.Fatalf("group_invites inserts = %d, want 1", len(ins))
	}
	if ins[0].args[0] != int64(1) || ins[0].args[1] != int64(7) || ins[0].args[2] != "new.person@example.com" {
		t.Fatalf("insert args = %v, want inviter 1, group 7, lowercased email", ins[0].args)
	}
	if len(f.ran("INSERT INTO group_members")) != 0 {
		t.Fatalf("unknown email must not be added as a member")
	}
	if len(em.groupInvite) != 1 {
		t.Fatalf("group invite emails = %d, want 1", len(em.groupInvite))
	}
	url := em.groupInvite[0]
	if !strings.HasPrefix(url, "https://www.remainfaithful.com/invite?token=") || !strings.HasSuffix(url, "&type=group") {
		t.Fatalf("accept URL = %q", url)
	}
}

func TestInviteMemberExistingUserAdded(t *testing.T) {
	rr, f, em := doInvite(t, inviteMember, inviteScenario{isMember: true, memberCount: 2, inviteeExists: true},
		`{"user_email":"partner@example.com"}`)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %q", rr.Code, rr.Body.String())
	}
	body := decodeBody(t, rr)
	if body["status"] != "added" || body["user_id"] != float64(42) {
		t.Fatalf("body = %v", body)
	}
	if len(f.ran("INSERT INTO group_members")) != 1 {
		t.Fatalf("expected one group_members insert")
	}
	if len(f.ran("INSERT INTO group_invites")) != 0 {
		t.Fatalf("existing user must not get a pending invite")
	}
	if len(em.groupInvite) != 1 || em.groupInvite[0] != "https://www.remainfaithful.com/invite?type=group" {
		t.Fatalf("emails = %v", em.groupInvite)
	}
}

func TestInviteMemberAlreadyMemberConflict(t *testing.T) {
	rr, _, em := doInvite(t, inviteMember, inviteScenario{isMember: true, memberCount: 2, inviteeExists: true, alreadyMember: true},
		`{"user_email":"partner@example.com"}`)
	if rr.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rr.Code)
	}
	if len(em.groupInvite) != 0 {
		t.Fatalf("already-member must not be emailed")
	}
}

func TestInviteMemberCapTwelve(t *testing.T) {
	rr, f, em := doInvite(t, inviteMember, inviteScenario{isMember: true, memberCount: 12},
		`{"user_email":"thirteenth@example.com"}`)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "maximum of 12 members") {
		t.Fatalf("body = %q", rr.Body.String())
	}
	if len(f.ran("INSERT")) != 0 || len(em.groupInvite) != 0 {
		t.Fatalf("full group must not write or email")
	}
}

func TestInviteMemberAcceptsEmailKey(t *testing.T) {
	rr, f, _ := doInvite(t, inviteMember, inviteScenario{isMember: true, memberCount: 1},
		`{"email":"Someone@Example.com"}`)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %q", rr.Code, rr.Body.String())
	}
	ins := f.ran("INSERT INTO group_invites")
	if len(ins) != 1 || ins[0].args[2] != "someone@example.com" {
		t.Fatalf("inserts = %v", ins)
	}
}

func TestInviteMemberRequiresEmail(t *testing.T) {
	rr, _, _ := doInvite(t, inviteMember, inviteScenario{isMember: true, memberCount: 1}, `{"user_email":"  "}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

// ── GroupEmailInvite (Settings) keeps its responses ─────────────────────────

func TestGroupEmailInviteUnknownEmailStillInvited(t *testing.T) {
	rr, f, _ := doInvite(t, groupEmailInvite, inviteScenario{isMember: true, memberCount: 3}, `{"email":"x@y.com"}`)
	if rr.Code != http.StatusCreated || decodeBody(t, rr)["status"] != "invited" {
		t.Fatalf("status = %d body = %q", rr.Code, rr.Body.String())
	}
	if len(f.ran("INSERT INTO group_invites")) != 1 {
		t.Fatalf("expected pending invite insert")
	}
}

func TestGroupEmailInviteAlreadyMemberOK(t *testing.T) {
	rr, _, _ := doInvite(t, groupEmailInvite, inviteScenario{isMember: true, memberCount: 3, inviteeExists: true, alreadyMember: true}, `{"email":"x@y.com"}`)
	if rr.Code != http.StatusOK || decodeBody(t, rr)["status"] != "already_member" {
		t.Fatalf("status = %d body = %q", rr.Code, rr.Body.String())
	}
}

func TestGroupEmailInviteCapTwelve(t *testing.T) {
	rr, _, _ := doInvite(t, groupEmailInvite, inviteScenario{isMember: true, memberCount: 12}, `{"email":"x@y.com"}`)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rr.Code)
	}
}

// ── redemption on sign-up ─────────────────────────────────────────────────

func redeemRouter(memberCount int64, groupInvites [][]driver.Value) func(string, []any) ([]string, [][]driver.Value) {
	return func(q string, _ []any) ([]string, [][]driver.Value) {
		switch {
		case strings.Contains(q, "FROM relationship_invites"):
			return []string{"inviter_id", "token"}, nil
		case strings.Contains(q, "FROM group_invites"):
			return []string{"id", "group_id"}, groupInvites
		case strings.Contains(q, "SELECT id FROM groups WHERE id = $1 FOR UPDATE"):
			return one("id", int64(7))
		case strings.Contains(q, "SELECT COUNT(*) FROM group_members"):
			return one("count", memberCount)
		case strings.Contains(q, "INSERT INTO users"):
			return []string{"id", "created_at"}, [][]driver.Value{{int64(99), "2026-10-08T20:00:00Z"}}
		}
		return none("x")
	}
}

func TestAcceptPendingGroupInvitesJoinsAndMarksAccepted(t *testing.T) {
	f := &inviteFakeDB{query: redeemRouter(2, [][]driver.Value{{int64(5), int64(7)}})}
	h := &H{DB: f.open(t)}
	h.acceptPendingGroupInvites(context.Background(), 99, "new@example.com")

	if q := f.ran("FROM group_invites"); len(q) != 1 || q[0].args[0] != "new@example.com" {
		t.Fatalf("pending lookup = %v", q)
	}
	ins := f.ran("INSERT INTO group_members")
	if len(ins) != 1 || ins[0].args[0] != int64(7) || ins[0].args[1] != int64(99) {
		t.Fatalf("member insert = %v", ins)
	}
	upd := f.ran("UPDATE group_invites SET status = 'accepted'")
	if len(upd) != 1 || upd[0].args[0] != int64(5) {
		t.Fatalf("invite update = %v", upd)
	}
	if len(f.ran("COMMIT")) != 1 {
		t.Fatalf("expected commit")
	}
}

func TestAcceptPendingGroupInvitesSkipsWhenGroupFull(t *testing.T) {
	f := &inviteFakeDB{query: redeemRouter(12, [][]driver.Value{{int64(5), int64(7)}})}
	h := &H{DB: f.open(t)}
	h.acceptPendingGroupInvites(context.Background(), 99, "new@example.com")

	if len(f.ran("INSERT INTO group_members")) != 0 {
		t.Fatalf("full group must not gain a member")
	}
	if len(f.ran("UPDATE group_invites")) != 0 {
		t.Fatalf("invite must stay pending when the group is full")
	}
	if len(f.ran("COMMIT")) != 0 {
		t.Fatalf("full group must not commit")
	}
}

func TestRegisterRedeemsGroupInviteBeforeResponding(t *testing.T) {
	f := &inviteFakeDB{query: redeemRouter(1, [][]driver.Value{{int64(5), int64(7)}})}
	h := &H{DB: f.open(t)}
	req := httptest.NewRequest(http.MethodPost, "/auth/register",
		strings.NewReader(`{"name":"New","email":"New@Example.com","password":"longenough1"}`))
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %q", rr.Code, rr.Body.String())
	}
	// No goroutine: by the time Register returns the membership must exist.
	ins := f.ran("INSERT INTO group_members")
	if len(ins) != 1 || ins[0].args[1] != int64(99) {
		t.Fatalf("member insert = %v", ins)
	}
	if q := f.ran("FROM group_invites"); len(q) != 1 || q[0].args[0] != "new@example.com" {
		t.Fatalf("lookup must use the normalized email: %v", q)
	}
}

func TestSocialSignupRedeemsOnlyRealEmails(t *testing.T) {
	newUser := func(q string, a []any) ([]string, [][]driver.Value) {
		switch {
		case strings.Contains(q, "WHERE apple_id"), strings.Contains(q, "WHERE email"):
			return []string{"id", "name", "email"}, nil
		case strings.Contains(q, "INSERT INTO users"):
			return []string{"id", "name", "email"}, [][]driver.Value{{int64(99), "A", a[1]}}
		}
		return redeemRouter(1, [][]driver.Value{{int64(5), int64(7)}})(q, a)
	}

	t.Run("placeholder email skips redemption", func(t *testing.T) {
		f := &inviteFakeDB{query: newUser}
		h := &H{DB: f.open(t)}
		req := httptest.NewRequest(http.MethodPost, "/auth/apple", nil)
		if _, _, _, err := h.findOrCreateSocialUser(req, "apple_id", "sub123", "", "A"); err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(f.ran("FROM group_invites")) != 0 || len(f.ran("FROM relationship_invites")) != 0 {
			t.Fatalf("placeholder account must not redeem invites")
		}
	})

	t.Run("real email redeems", func(t *testing.T) {
		f := &inviteFakeDB{query: newUser}
		h := &H{DB: f.open(t)}
		req := httptest.NewRequest(http.MethodPost, "/auth/google", nil)
		if _, _, _, err := h.findOrCreateSocialUser(req, "google_id", "g1", "Real@Example.com", "A"); err != nil {
			t.Fatalf("err = %v", err)
		}
		q := f.ran("FROM group_invites")
		if len(q) != 1 || q[0].args[0] != "real@example.com" {
			t.Fatalf("group invite lookup = %v", q)
		}
		if len(f.ran("INSERT INTO group_members")) != 1 {
			t.Fatalf("expected group join")
		}
	})
}

func TestSiteBaseDefaultsToLiveSite(t *testing.T) {
	t.Setenv("SITE_URL", "")
	if got := siteBase(); got != "https://www.remainfaithful.com" {
		t.Fatalf("siteBase() = %q", got)
	}
	t.Setenv("SITE_URL", "https://example.org/")
	if got := siteBase(); got != "https://example.org" {
		t.Fatalf("siteBase() override = %q", got)
	}
}
