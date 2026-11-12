package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWeeklySummaryFlagDefaultsOff(t *testing.T) {
	t.Setenv("WEEKLY_PARTNER_SUMMARY", "")
	if weeklySummaryEnabled() {
		t.Fatal("empty flag must stay off")
	}
	for _, raw := range []string{"0", "false", "off", "no", " TRUE "} {
		t.Setenv("WEEKLY_PARTNER_SUMMARY", raw)
		want := strings.TrimSpace(strings.ToLower(raw)) == "true"
		if weeklySummaryEnabled() != want {
			t.Fatalf("%q enabled=%v", raw, weeklySummaryEnabled())
		}
	}
	for _, raw := range []string{"1", "true", "YES", "on"} {
		t.Setenv("WEEKLY_PARTNER_SUMMARY", raw)
		if !weeklySummaryEnabled() {
			t.Fatalf("%q should enable the job", raw)
		}
	}
}

func TestWeeklyJobSendsNothingWhenOff(t *testing.T) {
	t.Setenv("WEEKLY_PARTNER_SUMMARY", "")
	mail := &recordingEmail{}
	h := &H{DB: newMemDB().open(t), Email: mail}
	h.sendDueWeeklySummaries(context.Background(), time.Now())
	if len(mail.weekly) != 0 {
		t.Fatalf("sent %d", len(mail.weekly))
	}
}

func TestWeeklyJobLooksUpNothingHarmfulWhenOn(t *testing.T) {
	t.Setenv("WEEKLY_PARTNER_SUMMARY", "true")
	mail := &recordingEmail{}
	h := &H{DB: newMemDB().open(t), Email: mail}
	h.sendDueWeeklySummaries(context.Background(), time.Now())
	if len(mail.weekly) != 0 {
		t.Fatalf("empty lookup sent %d", len(mail.weekly))
	}
}

func TestDeliverWeeklyLetterRespectsFlag(t *testing.T) {
	mail := &recordingEmail{}
	h := &H{Email: mail}
	letter := weeklyLetter{
		ToEmail: "bea@example.com", ToName: "Bea", SubjectName: "Alex",
		Days: 6, Alerts: 2, Streak: 4,
		UnsubURL: "https://api.remainfaithful.com/weekly-summary/unsubscribe?token=" + strings.Repeat("ab", 32),
	}
	t.Setenv("WEEKLY_PARTNER_SUMMARY", "")
	if err := h.deliverWeeklyLetter(context.Background(), letter); err != nil {
		t.Fatal(err)
	}
	if len(mail.weekly) != 0 {
		t.Fatal("flag off still sent")
	}
	t.Setenv("WEEKLY_PARTNER_SUMMARY", "true")
	if err := h.deliverWeeklyLetter(context.Background(), letter); err != nil {
		t.Fatal(err)
	}
	if len(mail.weekly) != 1 {
		t.Fatalf("sent %d", len(mail.weekly))
	}
	got := mail.weekly[0]
	if got.to != "bea@example.com" || got.subjectName != "Alex" || got.days != 6 || got.alerts != 2 || got.streak != 4 {
		t.Fatalf("%+v", got)
	}
	if !strings.Contains(got.url, "unsubscribe?token=") {
		t.Fatal(got.url)
	}
}

func TestWeeklyWindowAndStreak(t *testing.T) {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal(err)
	}
	early := time.Date(2026, 11, 9, 7, 59, 0, 0, loc)
	if _, open := weeklySendOpen(early, loc); open {
		t.Fatal("before Monday 08:00")
	}
	openAt := time.Date(2026, 11, 9, 8, 0, 0, 0, loc)
	start, open := weeklySendOpen(openAt, loc)
	if !open || start != "2026-11-09" {
		t.Fatalf("monday %s open=%v", start, open)
	}
	sunday := time.Date(2026, 11, 15, 12, 0, 0, 0, loc)
	start, open = weeklySendOpen(sunday, loc)
	if !open || start != "2026-11-09" {
		t.Fatalf("sunday %s open=%v", start, open)
	}
	from, to := lastSevenDays(time.Date(2026, 11, 12, 15, 0, 0, 0, loc), loc)
	if from != "2026-11-06" || to != "2026-11-12" {
		t.Fatalf("window %s %s", from, to)
	}
	protected := map[string]struct{}{"2026-11-10": {}, "2026-11-12": {}}
	streak := protectedStreak(time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC), loc, protected, nil, time.Date(2026, 11, 12, 15, 0, 0, 0, loc))
	if streak != 1 {
		t.Fatalf("streak %d", streak)
	}
	if countDays(protected, from, to) != 2 {
		t.Fatal(countDays(protected, from, to))
	}
}

func TestWeeklyPushAndUnsubscribePage(t *testing.T) {
	body := weeklyPushBody("Alex", 1, 1, 4)
	if body != "Alex: 1 day protected, 1 alert, streak 4." {
		t.Fatal(body)
	}
	if strings.ContainsAny(body, "\u2014\u2013") || strings.Contains(strings.ToLower(body), "honest") {
		t.Fatal(body)
	}
	h := &H{}
	rr := httptest.NewRecorder()
	h.UnsubscribeWeeklySummary(rr, httptest.NewRequest(http.MethodGet, "/weekly-summary/unsubscribe?token=nope", nil))
	if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "not valid") {
		t.Fatalf("%d %s", rr.Code, rr.Body.String())
	}
	db := newMemDB()
	h = &H{DB: db.open(t)}
	rr = httptest.NewRecorder()
	h.UnsubscribeWeeklySummary(rr, httptest.NewRequest(http.MethodGet, "/weekly-summary/unsubscribe?token="+strings.Repeat("ab", 32), nil))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("unknown token %d %s", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	h.SetWeeklySummaryPreference(rr, authed(http.MethodPost, "/users/me/weekly-summary", `{}`, 1))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("missing choice %d", rr.Code)
	}
}

func TestUnsubTokenShape(t *testing.T) {
	token, err := newUnsubToken()
	if err != nil || !validUnsubToken(token) {
		t.Fatalf("%s %v", token, err)
	}
	if validUnsubToken(strings.ToUpper(token)) || validUnsubToken(token+"aa") {
		t.Fatal("accepted a token that is not lowercase hex")
	}
}
