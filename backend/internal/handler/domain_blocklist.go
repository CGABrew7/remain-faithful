package handler

import "net/http"

// domainRule is one host Deep Scan may match on the phone.
// The phone keeps a copy for when this list cannot be loaded.
type domainRule struct {
	Host     string `json:"host"`
	Category string `json:"category"`
}

// domainBlocklist is the same host set the iOS broadcast extension bundles.
// Categories match the on-device groups: adult, gambling, and violence.
func domainBlocklist() []domainRule {
	adult := []string{
		"pornhub.com", "xvideos.com", "xhamster.com", "xnxx.com", "redtube.com",
		"youporn.com", "tube8.com", "spankbang.com", "ixxx.com", "beeg.com",
		"brazzers.com", "realitykings.com", "bangbros.com", "onlyfans.com",
		"chaturbate.com", "livejasmin.com", "cam4.com", "stripchat.com",
		"adultfriendfinder.com", "ashleymadison.com", "benaughty.com",
		"seeking.com", "seekingarrangement.com",
	}
	gambling := []string{
		"bovada.lv", "bet365.com", "draftkings.com", "fanduel.com",
		"pokerstars.com", "betmgm.com", "barstoolsportsbook.com",
	}
	violence := []string{"bestgore.com", "liveleak.com"}
	out := make([]domainRule, 0, len(adult)+len(gambling)+len(violence))
	for _, host := range adult {
		out = append(out, domainRule{Host: host, Category: "adult_content"})
	}
	for _, host := range gambling {
		out = append(out, domainRule{Host: host, Category: "gambling"})
	}
	for _, host := range violence {
		out = append(out, domainRule{Host: host, Category: "violence"})
	}
	return out
}

// DomainBlocklist returns the host list the phone may use for Deep Scan.
// GET /protection/domain-blocklist
// The phone still has the same list bundled if this request fails.
func (h *H) DomainBlocklist(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"domains": domainBlocklist()})
}
