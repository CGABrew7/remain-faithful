package handler

import (
	"database/sql/driver"
	"testing"
)

func TestTokenVersionAllowed(t *testing.T) {
	cases := []struct {
		name    string
		stored  int64
		claim   int
		missing bool
		want    bool
	}{
		{name: "legacy token before any reset", stored: 1, claim: 0, want: true},
		{name: "legacy token after a reset", stored: 2, claim: 0, want: false},
		{name: "current token", stored: 3, claim: 3, want: true},
		{name: "stale token", stored: 3, claim: 2, want: false},
		{name: "database error fails closed", missing: true, claim: 1, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &inviteFakeDB{query: func(string, []any) ([]string, [][]driver.Value) {
				if tc.missing {
					return none("token_version")
				}
				return one("token_version", tc.stored)
			}}
			h := &H{DB: f.open(t)}
			if got := h.TokenVersionAllowed(9, tc.claim); got != tc.want {
				t.Fatalf("allowed = %v, want %v", got, tc.want)
			}
		})
	}
}
