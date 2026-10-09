package email

import (
	"strings"
	"testing"
)

func TestContactPartsEscapesHTML(t *testing.T) {
	subject, plain, htmlBody := contactParts(
		`eve@example.com"><script>`,
		`Ann <img src=x onerror=alert(1)>`,
		"Hello\r\nBcc: spam@evil.test",
		"line1\n<script>alert(1)</script>",
	)

	if strings.Contains(subject, "\n") || strings.Contains(subject, "\r") {
		t.Fatalf("subject header still has a newline: %q", subject)
	}
	if strings.Contains(subject, "Bcc:") && strings.ContainsAny(subject, "\r\n") {
		t.Fatalf("subject kept a header break: %q", subject)
	}
	if !strings.Contains(subject, "[Remain Faithful Contact] Hello") {
		t.Fatalf("subject = %q", subject)
	}

	for _, raw := range []string{"<script>", "<img", `onerror=`} {
		if strings.Contains(htmlBody, raw) {
			t.Fatalf("html still contains %q:\n%s", raw, htmlBody)
		}
	}
	if !strings.Contains(htmlBody, "&lt;script&gt;") {
		t.Fatalf("message was not escaped:\n%s", htmlBody)
	}
	if !strings.Contains(htmlBody, "&lt;img") {
		t.Fatalf("name was not escaped:\n%s", htmlBody)
	}
	if !strings.Contains(htmlBody, "&lt;/script&gt;") && !strings.Contains(htmlBody, "&lt;script&gt;alert") {
		t.Fatalf("script tag was not escaped:\n%s", htmlBody)
	}
	if !strings.Contains(htmlBody, "line1<br>") {
		t.Fatalf("newline was not turned into a break:\n%s", htmlBody)
	}
	if !strings.Contains(plain, "<script>alert(1)</script>") {
		t.Fatalf("plain text should keep the original characters:\n%s", plain)
	}
	if !strings.Contains(plain, "Ann <img") {
		t.Fatalf("plain text should keep the name:\n%s", plain)
	}
}
