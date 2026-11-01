package handler

import (
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// memWaitlist is the in-memory stand-in for waitlist_signups.
type memWaitlist struct {
	email    string
	name     string
	platform string
}

// applyWaitlist handles the waitlist statements. apply already holds m.mu.
func (m *memDB) applyWaitlist(q string, args []any) ([]string, [][]driver.Value, bool) {
	if !containsAll(q, "INSERT INTO waitlist_signups", "ON CONFLICT (email, platform) DO NOTHING") {
		return nil, nil, false
	}
	email, name, platform := asString(args[0]), asString(args[1]), asString(args[2])
	for _, row := range m.waitlist {
		if row.email == email && row.platform == platform {
			return []string{"ok"}, nil, true
		}
	}
	m.waitlist = append(m.waitlist, memWaitlist{email: email, name: name, platform: platform})
	return []string{"ok"}, nil, true
}

func (m *memDB) waitlistRows() []memWaitlist {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]memWaitlist, len(m.waitlist))
	copy(out, m.waitlist)
	return out
}

func postWaitlist(t *testing.T, h *H, body string) *httptest.ResponseRecorder {
	t.Helper()
	t.Setenv("TURNSTILE_SECRET_KEY", "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/waitlist", strings.NewReader(withLoadedAt(body)))
	req.Header.Set("Content-Type", "application/json")
	h.JoinWaitlist(rr, req)
	return rr
}

func TestJoinWaitlistStoresSignup(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t)}

	rr := postWaitlist(t, h, `{"name":"  Ada 🙏 ","email":"  Ada+Android@Example.COM "}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rr.Code, rr.Body.String())
	}
	rows := db.waitlistRows()
	if len(rows) != 1 {
		t.Fatalf("rows = %d", len(rows))
	}
	if rows[0].email != "ada+android@example.com" || rows[0].name != "Ada 🙏" || rows[0].platform != "android" {
		t.Fatalf("row = %+v", rows[0])
	}
}

func TestJoinWaitlistDuplicateIsIdempotent(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t)}
	for i := 0; i < 3; i++ {
		rr := postWaitlist(t, h, `{"email":"dup@example.com","platform":"android"}`)
		if rr.Code != http.StatusOK {
			t.Fatalf("try %d status = %d body = %s", i, rr.Code, rr.Body.String())
		}
	}
	rr := postWaitlist(t, h, `{"email":"DUP@example.com "}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("case variant status = %d", rr.Code)
	}
	if n := len(db.waitlistRows()); n != 1 {
		t.Fatalf("rows = %d, want 1", n)
	}
}

func TestJoinWaitlistRejectsBadInput(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t)}
	long := strings.Repeat("a", 250) + "@example.com"
	bad := []string{
		`{"email":""}`,
		`{"email":"@"}`,
		`{"email":"bob@x"}`,
		`{"email":"not-an-email"}`,
		`{"email":"a@@b.com"}`,
		`{"email":"a b@example.com"}`,
		`{"email":"a@example..com"}`,
		`{"email":"a@-example.com"}`,
		`{"email":"zero\u200bwidth@example.com"}`,
		`{"email":"full@examplｅ.com"}`,
		`{"email":"x@example.com\r\nBcc: y@example.com"}`,
		`{"email":"Name <x@example.com>"}`,
		`{"email":"` + long + `"}`,
		`{"email":"ok@example.com","platform":"<script>"}`,
		`{"email":123}`,
		`{bad`,
	}
	for _, body := range bad {
		rr := postWaitlist(t, h, body)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want 400", body, rr.Code)
		}
	}
	if n := len(db.waitlistRows()); n != 0 {
		t.Fatalf("stored %d bad rows", n)
	}
}

func TestJoinWaitlistCapsNameAndDropsHoneypot(t *testing.T) {
	db := newMemDB()
	h := &H{DB: db.open(t)}
	rr := postWaitlist(t, h, `{"email":"long@example.com","name":"`+strings.Repeat("N", 5000)+`\u0000"}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	rows := db.waitlistRows()
	if len(rows) != 1 || len([]rune(rows[0].name)) != maxWaitlistName {
		t.Fatalf("rows = %+v", rows)
	}
	rr = postWaitlist(t, h, `{"email":"bot@example.com","company":"Acme"}`)
	if rr.Code != http.StatusOK || len(db.waitlistRows()) != 1 {
		t.Fatalf("honeypot stored a row: %d %d", rr.Code, len(db.waitlistRows()))
	}
}
