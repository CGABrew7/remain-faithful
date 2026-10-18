package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"remain-faithful/backend/internal/admin"
)

const overviewCacheTTL = 60 * time.Second

type overviewCacheEntry struct {
	at   time.Time
	body []byte
}

var (
	overviewMu    sync.Mutex
	overviewCache = map[string]overviewCacheEntry{}
	healthMu      sync.Mutex
	healthCache   overviewCacheEntry
)

func (h *H) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Cache-Control", "no-store")
	if !admin.Enabled() {
		writeError(w, http.StatusNotFound, "not found")
		return false
	}
	if !admin.Authorize(r) {
		writeError(w, http.StatusUnauthorized, "sign in required")
		return false
	}
	return true
}

func (h *H) auditAdmin(ctx context.Context, action, detail string) {
	if h == nil || h.DB == nil {
		return
	}
	action = clip(action, 80)
	detail = clip(detail, 240)
	if _, err := h.DB.ExecContext(ctx,
		`INSERT INTO admin_audit (action, detail) VALUES ($1, $2)`, action, detail,
	); err != nil {
		log.Printf("admin audit %s: %v", action, err)
	}
}

// AdminLogin checks the owner password and starts an 8 hour session.
// POST /admin/login
func (h *H) AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if !admin.Enabled() {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !admin.AllowLogin(admin.ClientIP(r), time.Now()) {
		h.auditAdmin(r.Context(), "login_blocked", "rate limit")
		w.Header().Set("Retry-After", "900")
		writeError(w, http.StatusTooManyRequests, "too many attempts")
		return
	}
	if !admin.CheckPassword(req.Password) {
		h.auditAdmin(r.Context(), "login_failure", "")
		writeError(w, http.StatusUnauthorized, "wrong password")
		return
	}
	token, err := admin.SignSession()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not start session")
		return
	}
	admin.SetCookie(w, r, token)
	h.auditAdmin(r.Context(), "login_success", "")
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "token": token})
}

// AdminLogout ends the admin session.
// POST /admin/logout
func (h *H) AdminLogout(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	admin.ClearCookie(w, r)
	h.auditAdmin(r.Context(), "logout", "")
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// AdminOverview returns cached aggregate metrics.
// GET /admin/overview?range=7d|30d|90d|all
func (h *H) AdminOverview(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	label, from := rangeStart(r.URL.Query().Get("range"), time.Now())
	overviewMu.Lock()
	if hit, ok := overviewCache[label]; ok && time.Since(hit.at) < overviewCacheTTL {
		body := hit.body
		overviewMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
		return
	}
	overviewMu.Unlock()

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	ov, err := h.loadOverview(ctx, label, from, time.Now())
	if err != nil {
		log.Printf("admin overview: %v", err)
		writeError(w, http.StatusInternalServerError, "could not load dashboard")
		return
	}
	body, err := json.Marshal(ov)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load dashboard")
		return
	}
	overviewMu.Lock()
	overviewCache[label] = overviewCacheEntry{at: time.Now(), body: body}
	overviewMu.Unlock()
	h.auditAdmin(ctx, "view_dashboard", label)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// AdminHealth returns process and database size figures.
// GET /admin/health
func (h *H) AdminHealth(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	healthMu.Lock()
	if healthCache.body != nil && time.Since(healthCache.at) < overviewCacheTTL {
		body := healthCache.body
		healthMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
		return
	}
	healthMu.Unlock()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.DB.PingContext(ctx); err != nil {
		writeJSON(w, http.StatusOK, &admin.Health{
			API: "ok", Database: "down", Tables: []admin.TableCount{},
			ServerErrors: admin.NotTracked("The API does not keep a counter of recent server errors."),
		})
		return
	}
	health, err := h.loadHealth(ctx)
	if err != nil {
		log.Printf("admin health: %v", err)
		writeError(w, http.StatusInternalServerError, "could not load health")
		return
	}
	body, err := json.Marshal(health)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load health")
		return
	}
	healthMu.Lock()
	healthCache = overviewCacheEntry{at: time.Now(), body: body}
	healthMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
