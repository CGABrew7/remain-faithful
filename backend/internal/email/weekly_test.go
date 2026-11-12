package email

import (
	"strings"
	"testing"
)

func TestWeeklySummaryMessage(t *testing.T) {
	subject, plain, htmlBody := weeklySummaryMessage("Bea", "Alex", 6, 2, 4, "https://api.remainfaithful.com/weekly-summary/unsubscribe?token=abc")
	if subject != "Weekly summary for Alex" {
		t.Fatal(subject)
	}
	for _, part := range []string{plain, htmlBody} {
		if !strings.Contains(part, "Days protected: 6") || !strings.Contains(part, "Alerts: 2") || !strings.Contains(part, "Streak: 4") {
			t.Fatal(part)
		}
		if !strings.Contains(part, "unsubscribe?token=abc") {
			t.Fatal(part)
		}
		if strings.ContainsAny(part, "\u2014\u2013") || strings.Contains(strings.ToLower(part), "honest") {
			t.Fatal(part)
		}
	}
	if !strings.Contains(htmlBody, "Stop these weekly emails") {
		t.Fatal(htmlBody)
	}
}
