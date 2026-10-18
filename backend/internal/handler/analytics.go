package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"remain-faithful/backend/internal/admin"
)

const (
	purgePageviewsSQL = `DELETE FROM pageviews WHERE created_at < NOW() - INTERVAL '13 months'`
	insertPageviewSQL = `
INSERT INTO pageviews (
    event_name, path, referrer_domain, utm_source, utm_medium, utm_campaign,
    country, region, device_type, visitor_hash
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
)

var (
	visitorHashPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)
	countryPattern     = regexp.MustCompile(`^[A-Z]{2}$`)
	pathPattern        = regexp.MustCompile(`^/(?:[A-Za-z0-9_-]+/?)*$`)
	collectMu          sync.Mutex
	collectBuckets     = map[string]*loginWindow{}
)

type loginWindow struct {
	count int
	start time.Time
}

// CollectAnalytics stores one cookieless website event.
// Country and region are supplied by the website from Vercel headers.
// The raw IP address is never read into this row.
// POST /analytics/collect
func (h *H) CollectAnalytics(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.DB == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Header.Get("DNT") == "1" || r.Header.Get("Sec-GPC") == "1" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	var req struct {
		Event       string `json:"event"`
		Path        string `json:"path"`
		Referrer    string `json:"referrer"`
		UTMSource   string `json:"utm_source"`
		UTMMedium   string `json:"utm_medium"`
		UTMCampaign string `json:"utm_campaign"`
		Country     string `json:"country"`
		Region      string `json:"region"`
		Device      string `json:"device"`
		VisitorHash string `json:"visitor_hash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	switch req.Event {
	case "pageview", "app_store_click", "android_signup":
	default:
		writeError(w, http.StatusBadRequest, "unknown event")
		return
	}
	path := cleanTrackPath(req.Path)
	if path == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	hash := strings.ToLower(strings.TrimSpace(req.VisitorHash))
	if !visitorHashPattern.MatchString(hash) {
		writeError(w, http.StatusBadRequest, "invalid visitor")
		return
	}
	if !allowCollect(hash, 30) || !allowCollect("ip:"+admin.ClientIP(r), 120) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	country := strings.ToUpper(strings.TrimSpace(req.Country))
	if !countryPattern.MatchString(country) {
		country = ""
	}
	region := ""
	if country == "US" && countryPattern.MatchString(strings.ToUpper(strings.TrimSpace(req.Region))) {
		region = strings.ToUpper(strings.TrimSpace(req.Region))
	}
	device := req.Device
	if device != "mobile" && device != "desktop" {
		device = "desktop"
	}
	if _, err := h.DB.ExecContext(r.Context(), insertPageviewSQL,
		req.Event,
		path,
		referrerDomain(req.Referrer),
		cleanUTM(req.UTMSource),
		cleanUTM(req.UTMMedium),
		cleanUTM(req.UTMCampaign),
		country,
		region,
		device,
		hash,
	); err != nil {
		log.Printf("analytics collect: %v", err)
		writeError(w, http.StatusInternalServerError, "could not record")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func cleanTrackPath(path string) string {
	path = strings.TrimSpace(path)
	if i := strings.IndexAny(path, "?#"); i >= 0 {
		path = path[:i]
	}
	if len(path) > 120 || !strings.HasPrefix(path, "/") {
		return ""
	}
	if path == "/admin" || strings.HasPrefix(path, "/admin/") {
		return ""
	}
	if path != "/" && !pathPattern.MatchString(path) {
		return ""
	}
	return path
}

func referrerDomain(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	switch host {
	case "remainfaithful.com", "remainfaithful.app":
		return ""
	}
	if len(host) > 253 {
		host = host[:253]
	}
	return host
}

func cleanUTM(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range s {
		if r < 32 || r == 127 {
			continue
		}
		b.WriteRune(r)
		if b.Len() >= 80 {
			break
		}
	}
	return b.String()
}

func allowCollect(key string, limit int) bool {
	collectMu.Lock()
	defer collectMu.Unlock()
	now := time.Now()
	if len(collectBuckets) > 2000 {
		for k, b := range collectBuckets {
			if now.Sub(b.start) >= time.Minute {
				delete(collectBuckets, k)
			}
		}
	}
	b, ok := collectBuckets[key]
	if !ok || now.Sub(b.start) >= time.Minute {
		collectBuckets[key] = &loginWindow{count: 1, start: now}
		return true
	}
	if b.count >= limit {
		return false
	}
	b.count++
	return true
}

// PurgeOldPageviews deletes analytics rows older than 13 months.
func (h *H) PurgeOldPageviews(ctx context.Context) error {
	if h == nil || h.DB == nil {
		return nil
	}
	_, err := h.DB.ExecContext(ctx, purgePageviewsSQL)
	return err
}

// StartPageviewRetention repeats the 13 month delete once a day.
func (h *H) StartPageviewRetention(ctx context.Context) {
	t := time.NewTicker(24 * time.Hour)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := h.PurgeOldPageviews(c); err != nil {
				log.Printf("pageview retention: %v", err)
			}
			cancel()
		}
	}
}
