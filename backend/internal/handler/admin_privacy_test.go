package handler

import (
	"strings"
	"testing"
)

func TestAggregateSQLStaysAggregate(t *testing.T) {
	banned := []string{
		"summary", "severity", "password", "bundle", "u.email", "users.email", "u.name",
		" select token", " device_tokens.token",
	}
	for _, q := range aggregateStatements() {
		lower := strings.ToLower(q)
		for _, bad := range banned {
			if strings.Contains(lower, bad) {
				t.Errorf("aggregate SQL contains %q:\n%s", bad, q)
			}
		}
	}
}

func TestPageviewInsertDoesNotStoreIP(t *testing.T) {
	lower := strings.ToLower(insertPageviewSQL)
	if strings.Contains(lower, "ip") {
		t.Fatalf("insert mentions ip: %s", insertPageviewSQL)
	}
	if !strings.Contains(purgePageviewsSQL, "13 months") {
		t.Fatalf("retention SQL = %s", purgePageviewsSQL)
	}
}

func TestAccountLookupSelectsStatusOnly(t *testing.T) {
	lower := strings.ToLower(accountLookupSQL)
	for _, bad := range []string{"summary", "password", "u.name", "category", "severity", "token"} {
		if strings.Contains(lower, bad) {
			t.Errorf("account lookup contains %q", bad)
		}
	}
	for _, need := range []string{"created_at", "email_verified", "last_heartbeat_at", "count(*)"} {
		if !strings.Contains(lower, need) {
			t.Errorf("account lookup missing %s", need)
		}
	}
}
