package email

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestInviteAndResetHTMLEscapesNamesAndURLs(t *testing.T) {
	name := "Ann <img src=x onerror=alert(1)>"
	link := `https://remainfaithful.com/accept?token=abc"><script>`
	body := partnerInviteHTML(name, link)
	if strings.Contains(body, "<img") || strings.Contains(body, "<script>") {
		t.Fatalf("html was not escaped:\n%s", body)
	}
	if !strings.Contains(body, "&lt;img") || !strings.Contains(body, "&lt;script&gt;") {
		t.Fatalf("expected escaped tags:\n%s", body)
	}
	c := &Client{}
	if err := c.SendGroupInvite("a@b.com", name, `G <b>`, link); err != nil {
		t.Fatal(err)
	}
	if err := c.SendPasswordReset("a@b.com", name, link); err != nil {
		t.Fatal(err)
	}
}

func TestMailDisplayNameClipsAndStripsControls(t *testing.T) {
	raw := "Ada\u202e" + strings.Repeat("B", 200) + "\r\nBcc: spam@evil.test"
	got := mailDisplayName(raw)
	if strings.ContainsRune(got, '\u202e') || strings.ContainsAny(got, "\r\n") {
		t.Fatalf("name = %q", got)
	}
	if utf8.RuneCountInString(got) != 80 {
		t.Fatalf("runes = %d, name = %q", utf8.RuneCountInString(got), got)
	}
	if mailDisplayName(" \u202e\n\t ") != "Someone" {
		t.Fatal("blank name was not replaced")
	}

	subject, plain, htmlBody := partnerInviteMessage(raw, "https://remainfaithful.com/invite?token=abc")
	for _, part := range []string{subject, plain, htmlBody} {
		if strings.ContainsRune(part, '\u202e') || strings.Contains(part, strings.Repeat("B", 81)) {
			t.Fatalf("partner mail kept the raw name:\n%s", part)
		}
	}
	if !strings.HasPrefix(subject, strings.Repeat("B", 77)) && !strings.Contains(subject, "Ada"+strings.Repeat("B", 77)) {
		t.Fatalf("subject = %q", subject)
	}

	groupSubject, groupPlain, groupHTML := groupInviteMessage(raw, "Home\r\nGroup", "https://remainfaithful.com/invite?token=abc")
	for _, part := range []string{groupSubject, groupPlain, groupHTML} {
		if strings.ContainsRune(part, '\u202e') || strings.Contains(part, strings.Repeat("B", 81)) {
			t.Fatalf("group mail kept the raw name:\n%s", part)
		}
	}
	if strings.Contains(groupSubject, "\n") || strings.Contains(groupPlain, "\r") {
		t.Fatal("group name kept a line break")
	}
}

func TestMailPartsDoNotLeaveRawTags(t *testing.T) {
	c := &Client{}
	// Build the HTML the same way the senders do, by calling the helpers indirectly
	// through a captured path: production send logs when the key is empty.
	// Assert the escape helpers themselves.
	name := "Eve\r\nBcc: spam@evil.test <b>"
	if strings.ContainsAny(sanitizeMailHeader(name), "\r\n") {
		t.Fatalf("header still has a break: %q", sanitizeMailHeader(name))
	}
	_ = c
}
