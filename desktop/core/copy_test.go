package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

func TestCopyHasNoDashesOrDonationLanguage(t *testing.T) {
	root := filepath.Clean("..")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := info.Name()
			if base == "build" || base == "wailsjs" || base == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".go", ".js", ".css", ".html", ".swift", ".plist", ".yml", ".md":
		default:
			return nil
		}
		if strings.HasSuffix(path, filepath.Join("core", "privacy.go")) || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(raw)
		if strings.ContainsRune(text, '\u2014') || strings.ContainsRune(text, '\u2013') {
			t.Errorf("%s contains a dash that copy should not use", path)
		}
		low := strings.ToLower(text)
		if strings.Contains(low, "honest") || strings.Contains(low, "donate") || strings.Contains(low, "donation") {
			t.Errorf("%s contains framing that does not belong in the app", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTourReusesPhoneSteps(t *testing.T) {
	steps := TourSteps(true)
	if len(steps) < 8 {
		t.Fatalf("%d steps", len(steps))
	}
	if steps[0].Title != "Monitoring status" {
		t.Fatal(steps[0].Title)
	}
	joined := ""
	for _, s := range steps {
		joined += s.Body + " "
	}
	if !strings.Contains(joined, "Groups hold up to 12 people.") {
		t.Fatal("missing group cap")
	}
	if ShouldShowTour(0, true) != true || ShouldShowTour(TourVersion, true) {
		t.Fatal("tour gate")
	}
	for _, r := range PauseNotice {
		if unicode.Is(unicode.Dash, r) && r != '-' {
			t.Fatal("pause notice dash")
		}
	}
}
