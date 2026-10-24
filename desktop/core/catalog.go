package core

import "strings"

// AppChoice is a local label for an app the person can choose.
// Key is a bundle id or process name used only on this computer.
type AppChoice struct {
	Key      string `json:"key"`
	Display  string `json:"display"`
	Kind     string `json:"kind"`
	Selected bool   `json:"selected"`
}

type catalogEntry struct {
	Platform string
	Key      string
	Display  string
	Kind     string
}

// Known apps the person can choose. Matching happens on the computer.
// Nothing in this list is uploaded.
var catalog = []catalogEntry{
	{Platform: "mac", Key: "com.apple.Safari", Display: "Safari", Kind: "browser"},
	{Platform: "mac", Key: "com.google.Chrome", Display: "Google Chrome", Kind: "browser"},
	{Platform: "mac", Key: "com.microsoft.edgemac", Display: "Microsoft Edge", Kind: "browser"},
	{Platform: "mac", Key: "org.mozilla.firefox", Display: "Firefox", Kind: "browser"},
	{Platform: "mac", Key: "com.brave.Browser", Display: "Brave", Kind: "browser"},
	{Platform: "mac", Key: "com.operasoftware.Opera", Display: "Opera", Kind: "browser"},
	{Platform: "mac", Key: "com.hnc.Discord", Display: "Discord", Kind: "app"},
	{Platform: "mac", Key: "com.tinyspeck.slackmacgap", Display: "Slack", Kind: "app"},
	{Platform: "mac", Key: "net.whatsapp.WhatsApp", Display: "WhatsApp", Kind: "app"},
	{Platform: "mac", Key: "ru.keepcoder.Telegram", Display: "Telegram", Kind: "app"},
	{Platform: "mac", Key: "com.spotify.client", Display: "Spotify", Kind: "app"},
	{Platform: "mac", Key: "com.valvesoftware.steam", Display: "Steam", Kind: "app"},
	{Platform: "mac", Key: "com.apple.TV", Display: "Apple TV", Kind: "app"},

	{Platform: "windows", Key: "chrome.exe", Display: "Google Chrome", Kind: "browser"},
	{Platform: "windows", Key: "msedge.exe", Display: "Microsoft Edge", Kind: "browser"},
	{Platform: "windows", Key: "firefox.exe", Display: "Firefox", Kind: "browser"},
	{Platform: "windows", Key: "brave.exe", Display: "Brave", Kind: "browser"},
	{Platform: "windows", Key: "opera.exe", Display: "Opera", Kind: "browser"},
	{Platform: "windows", Key: "discord.exe", Display: "Discord", Kind: "app"},
	{Platform: "windows", Key: "slack.exe", Display: "Slack", Kind: "app"},
	{Platform: "windows", Key: "whatsapp.exe", Display: "WhatsApp", Kind: "app"},
	{Platform: "windows", Key: "telegram.exe", Display: "Telegram", Kind: "app"},
	{Platform: "windows", Key: "spotify.exe", Display: "Spotify", Kind: "app"},
	{Platform: "windows", Key: "steam.exe", Display: "Steam", Kind: "app"},
}

func Catalog(platform string) []AppChoice {
	var out []AppChoice
	for _, e := range catalog {
		if e.Platform == platform {
			out = append(out, AppChoice{Key: e.Key, Display: e.Display, Kind: e.Kind})
		}
	}
	return out
}

func Lookup(platform, key string) (AppChoice, bool) {
	key = normalizeKey(platform, key)
	for _, e := range catalog {
		if e.Platform == platform && e.Key == key {
			return AppChoice{Key: e.Key, Display: e.Display, Kind: e.Kind}, true
		}
	}
	return AppChoice{}, false
}

func AllowQuit(platform, key string) bool {
	_, ok := Lookup(platform, key)
	return ok
}

func IsSelf(platform, key string) bool {
	key = normalizeKey(platform, key)
	switch platform {
	case "mac":
		return key == "com.remainfaithful.app.mac"
	case "windows":
		return key == "remainfaithful.exe"
	default:
		return false
	}
}

func normalizeKey(platform, key string) string {
	key = strings.TrimSpace(key)
	if platform == "windows" {
		return strings.ToLower(key)
	}
	return key
}

func FilterSelection(platform string, keys []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, key := range keys {
		choice, ok := Lookup(platform, key)
		if !ok {
			continue
		}
		if _, dup := seen[choice.Key]; dup {
			continue
		}
		seen[choice.Key] = struct{}{}
		out = append(out, choice.Key)
	}
	return out
}
