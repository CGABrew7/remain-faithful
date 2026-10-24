package monitor

import "testing"

func TestQuitRefusesUnchosenApps(t *testing.T) {
	if err := Quit("explorer.exe"); err == nil {
		t.Fatal("expected a refusal before any process is touched")
	}
	if err := Quit("com.apple.finder"); err == nil {
		t.Fatal("finder is not a chosen app")
	}
}
