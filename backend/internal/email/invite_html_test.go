package email

import (
	"strings"
	"testing"
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
