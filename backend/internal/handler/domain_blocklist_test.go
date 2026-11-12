package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rfauth "remain-faithful/backend/internal/auth"
)

func TestDomainBlocklistMatchesBundledGroups(t *testing.T) {
	rules := domainBlocklist()
	if len(rules) != 32 {
		t.Fatalf("count %d", len(rules))
	}
	seen := map[string]string{}
	gambling, violence := 0, 0
	for _, rule := range rules {
		if rule.Host != strings.ToLower(rule.Host) || strings.Contains(rule.Host, "/") || strings.Contains(rule.Host, ":") {
			t.Fatalf("host %q", rule.Host)
		}
		if _, ok := seen[rule.Host]; ok {
			t.Fatalf("duplicate %s", rule.Host)
		}
		seen[rule.Host] = rule.Category
		switch rule.Category {
		case "adult_content":
		case "gambling":
			gambling++
		case "violence":
			violence++
		default:
			t.Fatalf("category %s", rule.Category)
		}
	}
	if gambling != 7 || violence != 2 {
		t.Fatalf("gambling %d violence %d", gambling, violence)
	}
	if seen["pornhub.com"] != "adult_content" || seen["draftkings.com"] != "gambling" || seen["liveleak.com"] != "violence" {
		t.Fatalf("%v", seen)
	}
}

func TestDomainBlocklistHandler(t *testing.T) {
	h := &H{}
	req := httptest.NewRequest(http.MethodGet, "/protection/domain-blocklist", nil)
	req = req.WithContext(rfauth.ContextWithUser(req.Context(), 1))
	rr := httptest.NewRecorder()
	h.DomainBlocklist(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var body struct {
		Domains []domainRule `json:"domains"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Domains) != 32 {
		t.Fatalf("body %d", len(body.Domains))
	}
}

func TestPerceptualHashSummaryIsNotAccepted(t *testing.T) {
	if isAllowedEventSummary("Explicit image detected (perceptual hash match)") {
		t.Fatal("image hash summary is still accepted")
	}
}
