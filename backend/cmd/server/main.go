package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"remain-faithful/backend/internal/apns"
	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/email"
	"remain-faithful/backend/internal/handler"
	"remain-faithful/backend/internal/payment"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	if err := migrate(db); err != nil {
		log.Fatalf("migration: %v", err)
	}
	log.Println("database ready")

	apnsClient, err := apns.New(
		os.Getenv("APNS_KEY_ID"),
		os.Getenv("APNS_TEAM_ID"),
		getenv("APNS_BUNDLE_ID", "com.remainfaithful.app"),
		os.Getenv("APNS_PRIVATE_KEY"),
		os.Getenv("APNS_PRODUCTION") == "true",
	)
	if err != nil {
		log.Fatalf("apns: %v", err)
	}
	if apnsClient.IsNoop() {
		log.Println("apns: WARNING — running as no-op, push notifications will not be delivered. Set APNS_KEY_ID, APNS_TEAM_ID, APNS_BUNDLE_ID, APNS_PRIVATE_KEY.")
	} else {
		log.Printf("apns: configured for %s", apnsClient.Environment())
	}

	emailClient := email.New()
	stripeClient := payment.New()

	h := &handler.H{DB: db, APNS: apnsClient, Email: emailClient, Stripe: stripeClient}

	srv := &http.Server{
		Addr:         ":" + port(),
		Handler:      routes(h),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	silenceMinutes := heartbeatSilenceMinutes()
	go h.StartHeartbeatWatcher(ctx)
	go h.StartHeartbeatSweep(ctx, silenceMinutes)
	if err := h.PurgeOldPageviews(ctx); err != nil {
		log.Printf("pageview retention: %v", err)
	}
	go h.StartPageviewRetention(ctx)

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
	log.Println("server stopped")
}

// loggingMiddleware logs method, path, status code, and duration for every request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lw.status, time.Since(start).Round(time.Millisecond))
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

// authRateLimiter limits unauthenticated endpoints to 10 requests per
// minute per IP (fixed window). Used for /auth/* and POST /contact.
// Uses only stdlib; no external dependency.
func authRateLimiter() mux.MiddlewareFunc {
	return fixedWindowRateLimiter(10, clientIP)
}

// authenticatedRateLimiter limits authenticated spam surfaces to `limit`
// requests per minute. Keyed by user ID when JWT middleware has already
// run, otherwise by client IP. Heartbeat uses a higher limit so the
// broadcast extension's 2-minute cadence is never throttled.
func authenticatedRateLimiter(limit int) mux.MiddlewareFunc {
	return fixedWindowRateLimiter(limit, authenticatedRateLimitKey)
}

func authenticatedRateLimitKey(r *http.Request) string {
	if uid, ok := rfauth.UserIDFromContext(r.Context()); ok {
		return "user:" + strconv.FormatInt(uid, 10)
	}
	return "ip:" + clientIP(r)
}

// fixedWindowRateLimiter allows `limit` requests per minute per key.
// Uses only stdlib; no external dependency.
func fixedWindowRateLimiter(limit int, keyFn func(*http.Request) string) mux.MiddlewareFunc {
	type window struct {
		count int
		start time.Time
	}
	var mu sync.Mutex
	windows := make(map[string]*window)

	// Prune windows older than 2 minutes every 5 minutes.
	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			mu.Lock()
			cutoff := time.Now().Add(-2 * time.Minute)
			for key, w := range windows {
				if w.start.Before(cutoff) {
					delete(windows, key)
				}
			}
			mu.Unlock()
		}
	}()

	allow := func(key string) bool {
		mu.Lock()
		defer mu.Unlock()
		now := time.Now()
		w, ok := windows[key]
		if !ok || now.Sub(w.start) >= time.Minute {
			windows[key] = &window{count: 1, start: now}
			return true
		}
		w.count++
		return w.count <= limit
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allow(keyFn(r)) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprint(w, `{"error":"too many requests — try again in a minute"}`)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// debugPushEnabled reports whether POST /debug/test-push should be registered.
// Off when APP_ENV=production (fly.toml). Unset / development / staging stay on.
func debugPushEnabled() bool {
	return os.Getenv("APP_ENV") != "production"
}

// clientIP extracts the real client IP, preferring the Fly-Client-IP header
// set by Fly.io's edge proxy (which cannot be spoofed by the client).
// In production a missing Fly header does not fall back to X-Forwarded-For.
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("Fly-Client-IP"); ip != "" {
		return ip
	}
	if os.Getenv("APP_ENV") != "production" {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
		}
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i != -1 {
		return host[:i]
	}
	return host
}

// corsMiddleware sets CORS headers. Only origins listed in ALLOWED_ORIGINS (comma-separated) receive
// the Allow-Origin echo; preflight OPTIONS requests are short-circuited with 204.
func corsMiddleware(allowedOrigins []string) mux.MiddlewareFunc {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if s := strings.TrimSpace(o); s != "" {
			allowed[s] = true
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Add("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// routes wires all HTTP endpoints onto a gorilla/mux router.
func routes(h *handler.H) http.Handler {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	origins := strings.Split(getenv("ALLOWED_ORIGINS", "https://remainfaithful.com"), ",")
	r.Use(corsMiddleware(origins))
	// CORS middleware answers preflight only after a route matches.
	// This catch-all lets every OPTIONS request reach that middleware.
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	// Health check (unauthenticated)
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	}).Methods(http.MethodGet)

	// Public auth routes — rate-limited to 10 requests/minute per IP.
	auth := r.PathPrefix("/auth").Subrouter()
	auth.Use(authRateLimiter())
	auth.HandleFunc("/register", h.Register).Methods(http.MethodPost)
	auth.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	auth.HandleFunc("/apple", h.AppleSignIn).Methods(http.MethodPost)
	auth.HandleFunc("/google", h.GoogleSignIn).Methods(http.MethodPost)
	auth.HandleFunc("/forgot-password", h.ForgotPassword).Methods(http.MethodPost)
	auth.HandleFunc("/reset-password", h.ResetPassword).Methods(http.MethodPost)
	auth.HandleFunc("/verify-email", h.VerifyEmail).Methods(http.MethodPost)

	// Protected routes — JWT required, and the session version must still match.
	api := r.NewRoute().Subrouter()
	api.Use(rfauth.MiddlewareWithVersion(func(userID int64, claimVersion int) bool {
		if h == nil {
			return false
		}
		return h.TokenVersionAllowed(userID, claimVersion)
	}))

	api.HandleFunc("/users/me", h.GetMe).Methods(http.MethodGet)
	api.HandleFunc("/users/me", h.UpdateMe).Methods(http.MethodPut)
	api.HandleFunc("/users/me", h.DeleteMe).Methods(http.MethodDelete)
	api.HandleFunc("/users/me/stats", h.GetUserStats).Methods(http.MethodGet)
	api.Handle("/relationships", authenticatedRateLimiter(10)(http.HandlerFunc(h.CreateRelationship))).Methods(http.MethodPost)
	api.HandleFunc("/relationships", h.ListRelationships).Methods(http.MethodGet)
	api.Handle("/relationships/invite", authenticatedRateLimiter(10)(http.HandlerFunc(h.InvitePartner))).Methods(http.MethodPost)
	api.HandleFunc("/relationships/accept-invite", h.AcceptPartnerInvite).Methods(http.MethodPost)
	api.HandleFunc("/relationships/decline-invite", h.DeclinePartnerInvite).Methods(http.MethodPost)
	api.HandleFunc("/relationships/invites/pending", h.ListPendingPartnerInvites).Methods(http.MethodGet)
	api.HandleFunc("/relationships/invites/{id}/accept", h.AcceptPartnerInviteByID).Methods(http.MethodPost)
	api.HandleFunc("/relationships/invites/{id}/decline", h.DeclinePartnerInviteByID).Methods(http.MethodPost)
	api.HandleFunc("/relationships/{id}/primary", h.SetPrimaryPartner).Methods(http.MethodPut)
	api.HandleFunc("/relationships/{id}/pin", h.SetRelationshipPIN).Methods(http.MethodPost)
	api.HandleFunc("/relationships/{id}/pin", h.DeleteRelationshipPIN).Methods(http.MethodDelete)
	api.HandleFunc("/relationships/{id}", h.DeleteRelationship).Methods(http.MethodDelete)
	// Authenticated write surfaces that fan out partner pushes are
	// fixed-window limited per user (see authenticatedRateLimiter).
	// Heartbeat is more generous so the 2-minute broadcast cadence is safe.
	api.Handle("/protection/alerts", authenticatedRateLimiter(20)(http.HandlerFunc(h.SendProtectionAlert))).Methods(http.MethodPost)
	api.HandleFunc("/protection/pin", h.GetPINStatus).Methods(http.MethodGet)
	api.HandleFunc("/protection/pin/verify", h.VerifyPIN).Methods(http.MethodPost)
	api.HandleFunc("/groups/leave-all", h.LeaveAllGroups).Methods(http.MethodPost)
	api.HandleFunc("/groups", h.ListMyGroups).Methods(http.MethodGet)
	api.HandleFunc("/groups", h.CreateGroup).Methods(http.MethodPost)
	api.HandleFunc("/groups/invites/pending", h.ListPendingGroupInvites).Methods(http.MethodGet)
	api.HandleFunc("/groups/accept-invite", h.AcceptGroupInvite).Methods(http.MethodPost)
	api.HandleFunc("/groups/decline-invite", h.DeclineGroupInvite).Methods(http.MethodPost)
	api.HandleFunc("/groups/invites/{id}/accept", h.AcceptGroupInviteByID).Methods(http.MethodPost)
	api.HandleFunc("/groups/invites/{id}/decline", h.DeclineGroupInviteByID).Methods(http.MethodPost)
	api.HandleFunc("/groups/{id}", h.GetGroup).Methods(http.MethodGet)
	api.HandleFunc("/groups/{id}", h.UpdateGroup).Methods(http.MethodPatch)
	api.Handle("/groups/{id}/encouragement", authenticatedRateLimiter(10)(http.HandlerFunc(h.SendEncouragement))).Methods(http.MethodPost)
	api.HandleFunc("/groups/{id}/members/me", h.LeaveGroup).Methods(http.MethodDelete)
	api.Handle("/groups/{id}/invite", authenticatedRateLimiter(10)(http.HandlerFunc(h.InviteMember))).Methods(http.MethodPost)
	api.Handle("/groups/{id}/email-invite", authenticatedRateLimiter(10)(http.HandlerFunc(h.GroupEmailInvite))).Methods(http.MethodPost)
	api.Handle("/events", authenticatedRateLimiter(30)(http.HandlerFunc(h.CreateEvent))).Methods(http.MethodPost)
	api.HandleFunc("/events", h.ListEvents).Methods(http.MethodGet)
	api.HandleFunc("/alerts", h.ListAlerts).Methods(http.MethodGet)
	api.HandleFunc("/alerts/count", h.AlertUnreadCount).Methods(http.MethodGet)
	api.HandleFunc("/alerts/mark-seen", h.MarkAlertsSeen).Methods(http.MethodPost)
	api.HandleFunc("/alerts/{id}/discussed", h.MarkAlertDiscussed).Methods(http.MethodPatch)
	api.Handle("/heartbeat", authenticatedRateLimiter(60)(http.HandlerFunc(h.Heartbeat))).Methods(http.MethodPost)
	api.HandleFunc("/users/device-token", h.RegisterDeviceToken).Methods(http.MethodPost)
	api.Handle("/panic", authenticatedRateLimiter(20)(http.HandlerFunc(h.SendPanicAlert))).Methods(http.MethodPost)
	api.Handle("/feedback", authenticatedRateLimiter(5)(http.HandlerFunc(h.SubmitFeedback))).Methods(http.MethodPost)
	if debugPushEnabled() {
		api.HandleFunc("/debug/test-push", h.SendTestPush).Methods(http.MethodPost)
	}
	api.HandleFunc("/auth/refresh", h.RefreshToken).Methods(http.MethodPost)
	api.HandleFunc("/auth/resend-verification", h.ResendVerification).Methods(http.MethodPost)
	api.HandleFunc("/donations/create-checkout-session", h.CreateCheckoutSession).Methods(http.MethodPost)
	api.HandleFunc("/donations/billing-portal", h.CreateBillingPortalSession).Methods(http.MethodPost)

	// Owner dashboard. Disabled (404) unless ADMIN_PASSWORD_HASH is set.
	// These routes are not behind the user JWT.
	owner := r.PathPrefix("/admin").Subrouter()
	owner.HandleFunc("/login", h.AdminLogin).Methods(http.MethodPost)
	owner.HandleFunc("/logout", h.AdminLogout).Methods(http.MethodPost)
	owner.HandleFunc("/overview", h.AdminOverview).Methods(http.MethodGet)
	owner.HandleFunc("/health", h.AdminHealth).Methods(http.MethodGet)
	owner.HandleFunc("/account", h.AdminLookupAccount).Methods(http.MethodGet)
	owner.HandleFunc("/account/resend-verification", h.AdminResendVerification).Methods(http.MethodPost)
	owner.HandleFunc("/account/password-reset", h.AdminSendPasswordReset).Methods(http.MethodPost)
	owner.HandleFunc("/account/delete", h.AdminDeleteAccount).Methods(http.MethodPost)

	// Cookieless website analytics. The website adds country from Vercel headers.
	r.HandleFunc("/analytics/collect", h.CollectAnalytics).Methods(http.MethodPost)

	// Contact form — unauthenticated, rate-limited to 10 requests/minute per IP
	// (same fixed-window limiter as /auth).
	contact := r.NewRoute().Subrouter()
	contact.Use(authRateLimiter())
	contact.HandleFunc("/contact", h.Contact).Methods(http.MethodPost)

	// Stripe webhook — unauthenticated, verified by Stripe-Signature header
	r.HandleFunc("/donations/webhook", h.DonationWebhook).Methods(http.MethodPost)

	return r
}

// openDB connects to PostgreSQL with a retry loop so the server survives Docker
// startup races where postgres isn't yet accepting connections.
func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	const maxAttempts = 10
	for i := 1; i <= maxAttempts; i++ {
		if err = db.Ping(); err == nil {
			return db, nil
		}
		log.Printf("waiting for database (%d/%d): %v", i, maxAttempts, err)
		time.Sleep(time.Duration(i) * time.Second)
	}
	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxAttempts, err)
}

func dsn() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", "postgres"),
		getenv("DB_NAME", "remainfaithful"),
	)
}

func port() string { return getenv("PORT", "8080") }

// heartbeatSilenceMinutes returns the number of minutes without a heartbeat
// before a heartbeat_silence alert is sent to the user's partners.
// Configured via HEARTBEAT_SILENCE_MINUTES (default 30).
func heartbeatSilenceMinutes() int {
	v := os.Getenv("HEARTBEAT_SILENCE_MINUTES")
	if v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
		log.Printf("HEARTBEAT_SILENCE_MINUTES=%q is not a positive integer, using default 30", v)
	}
	return 30
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// migrate applies the schema idempotently. Safe to call on every startup.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id            BIGSERIAL    PRIMARY KEY,
		name          TEXT         NOT NULL,
		email         TEXT         NOT NULL UNIQUE,
		password_hash TEXT         NOT NULL DEFAULT '',
		apple_id      TEXT         UNIQUE,
		google_id     TEXT         UNIQUE,
		created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS relationships (
		id         BIGSERIAL   PRIMARY KEY,
		user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		partner_id BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type       TEXT        NOT NULL DEFAULT 'partner',
		status     TEXT        NOT NULL DEFAULT 'pending',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE (user_id, partner_id)
	);

	CREATE TABLE IF NOT EXISTS groups (
		id         BIGSERIAL   PRIMARY KEY,
		name       TEXT        NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS group_members (
		group_id  BIGINT      NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		user_id   BIGINT      NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
		role      TEXT        NOT NULL DEFAULT 'member',
		joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY (group_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS events (
		id        BIGSERIAL   PRIMARY KEY,
		user_id   BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		category  TEXT        NOT NULL,
		severity  TEXT        NOT NULL,
		summary   TEXT        NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS alerts (
		id              BIGSERIAL   PRIMARY KEY,
		event_id        BIGINT      NOT NULL REFERENCES events(id)        ON DELETE CASCADE,
		relationship_id BIGINT               REFERENCES relationships(id) ON DELETE SET NULL,
		seen            BOOLEAN     NOT NULL DEFAULT FALSE,
		created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_events_user_id        ON events(user_id);
	CREATE INDEX IF NOT EXISTS idx_alerts_relationship   ON alerts(relationship_id);
	CREATE INDEX IF NOT EXISTS idx_relationships_user    ON relationships(user_id);
	CREATE INDEX IF NOT EXISTS idx_relationships_partner ON relationships(partner_id);

	CREATE TABLE IF NOT EXISTS device_tokens (
		id          BIGSERIAL   PRIMARY KEY,
		user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token       TEXT        NOT NULL,
		platform    TEXT        NOT NULL DEFAULT 'ios',
		environment TEXT        NOT NULL DEFAULT 'sandbox',
		is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE (user_id, token)
	);
	ALTER TABLE device_tokens ADD COLUMN IF NOT EXISTS environment TEXT NOT NULL DEFAULT 'sandbox';
	CREATE INDEX IF NOT EXISTS idx_device_tokens_user ON device_tokens(user_id);

	CREATE TABLE IF NOT EXISTS password_reset_tokens (
		id         BIGSERIAL   PRIMARY KEY,
		user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token      TEXT        NOT NULL UNIQUE,
		expires_at TIMESTAMPTZ NOT NULL,
		UNIQUE (user_id)
	);

	CREATE TABLE IF NOT EXISTS donations (
		id                BIGSERIAL    PRIMARY KEY,
		user_id           BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		stripe_session_id TEXT         NOT NULL UNIQUE,
		amount_cents      BIGINT       NOT NULL,
		monthly           BOOLEAN      NOT NULL DEFAULT FALSE,
		created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_donations_user ON donations(user_id);

	ALTER TABLE donations ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'completed';
	ALTER TABLE donations ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;
	ALTER TABLE donations ADD COLUMN IF NOT EXISTS stripe_subscription_id TEXT;
	ALTER TABLE donations ADD COLUMN IF NOT EXISTS stripe_payment_intent_id TEXT;
	CREATE INDEX IF NOT EXISTS idx_donations_customer ON donations(stripe_customer_id);
	CREATE INDEX IF NOT EXISTS idx_donations_subscription ON donations(stripe_subscription_id);

	ALTER TABLE users ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;

	-- Idempotent Stripe webhook log. Business logic runs only when the insert is new.
	CREATE TABLE IF NOT EXISTS donation_events (
		id              BIGSERIAL    PRIMARY KEY,
		stripe_event_id TEXT         NOT NULL UNIQUE,
		type            TEXT         NOT NULL,
		payload         JSONB,
		created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_donation_events_type ON donation_events(type);

	CREATE TABLE IF NOT EXISTS stripe_subscriptions (
		id                     BIGSERIAL    PRIMARY KEY,
		user_id                BIGINT       REFERENCES users(id) ON DELETE CASCADE,
		stripe_customer_id     TEXT         NOT NULL DEFAULT '',
		stripe_subscription_id TEXT         NOT NULL UNIQUE,
		status                 TEXT         NOT NULL,
		updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_stripe_subscriptions_user ON stripe_subscriptions(user_id);
	CREATE INDEX IF NOT EXISTS idx_stripe_subscriptions_customer ON stripe_subscriptions(stripe_customer_id);

	-- Idempotent column additions for existing databases.
	ALTER TABLE users ADD COLUMN IF NOT EXISTS apple_id  TEXT UNIQUE;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id TEXT UNIQUE;
	ALTER TABLE users ALTER COLUMN password_hash SET DEFAULT '';

	CREATE TABLE IF NOT EXISTS relationship_invites (
		id            BIGSERIAL   PRIMARY KEY,
		inviter_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		invitee_email TEXT        NOT NULL,
		token         TEXT        NOT NULL UNIQUE,
		status        TEXT        NOT NULL DEFAULT 'pending',
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE (inviter_id, invitee_email)
	);
	CREATE INDEX IF NOT EXISTS idx_invites_email  ON relationship_invites(invitee_email);
	CREATE INDEX IF NOT EXISTS idx_invites_token  ON relationship_invites(token);

	CREATE TABLE IF NOT EXISTS group_invites (
		id            BIGSERIAL   PRIMARY KEY,
		inviter_id    BIGINT      NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
		group_id      BIGINT      NOT NULL REFERENCES groups(id)  ON DELETE CASCADE,
		invitee_email TEXT        NOT NULL,
		token         TEXT        NOT NULL UNIQUE,
		status        TEXT        NOT NULL DEFAULT 'pending',
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE (group_id, invitee_email)
	);
	CREATE INDEX IF NOT EXISTS idx_group_invites_email ON group_invites(invitee_email);
	CREATE INDEX IF NOT EXISTS idx_group_invites_token ON group_invites(token);

	ALTER TABLE groups ADD COLUMN IF NOT EXISTS covenant TEXT NOT NULL DEFAULT '';

	ALTER TABLE relationships ADD COLUMN IF NOT EXISTS is_primary BOOLEAN NOT NULL DEFAULT FALSE;

	ALTER TABLE users ADD COLUMN IF NOT EXISTS last_heartbeat_at    TIMESTAMPTZ;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS heartbeat_screen_state TEXT;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS heartbeat_notified_at TIMESTAMPTZ;

	ALTER TABLE alerts ADD COLUMN IF NOT EXISTS recipient_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE;
	CREATE INDEX IF NOT EXISTS idx_alerts_recipient ON alerts(recipient_user_id);

	ALTER TABLE alerts ADD COLUMN IF NOT EXISTS discussed BOOLEAN NOT NULL DEFAULT FALSE;

	-- Protection features: partner PIN stored as bcrypt hash, never plaintext.
	ALTER TABLE relationships ADD COLUMN IF NOT EXISTS pin_hash TEXT;

	-- Heartbeat silence sweep: tracks the last time a silence alert was sent
	-- for each (monitored-user, partner) pair so the sweep doesn't re-alert
	-- within the same silence window. Reset implicitly when the user resumes
	-- heartbeats (last_heartbeat_at advances past last_silence_alert_at).
	ALTER TABLE relationships ADD COLUMN IF NOT EXISTS last_silence_alert_at TIMESTAMPTZ;

	-- Existing accounts stay verified. New password signups insert FALSE explicitly.
	ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN NOT NULL DEFAULT TRUE;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS token_version INT NOT NULL DEFAULT 1;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS pending_email TEXT;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS pending_email_token_hash TEXT;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS pending_email_expires_at TIMESTAMPTZ;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_mail_count INT NOT NULL DEFAULT 0;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_mail_window TIMESTAMPTZ;

	CREATE TABLE IF NOT EXISTS email_verifications (
		user_id    BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		token_hash TEXT        NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	-- Protection and tamper rows have no events row. Build 9 still receives a nested event object.
	ALTER TABLE alerts ALTER COLUMN event_id DROP NOT NULL;
	ALTER TABLE alerts ADD COLUMN IF NOT EXISTS kind TEXT NOT NULL DEFAULT 'event';
	ALTER TABLE alerts ADD COLUMN IF NOT EXISTS alert_type TEXT NOT NULL DEFAULT '';
	ALTER TABLE alerts ADD COLUMN IF NOT EXISTS summary TEXT NOT NULL DEFAULT '';
	ALTER TABLE alerts ADD COLUMN IF NOT EXISTS subject_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
	CREATE INDEX IF NOT EXISTS idx_alerts_subject ON alerts(subject_user_id);

	CREATE TABLE IF NOT EXISTS pin_verify_attempts (
		user_id      BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		failures     INT         NOT NULL DEFAULT 0,
		window_start TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_events_user_timestamp ON events(user_id, timestamp);
	CREATE INDEX IF NOT EXISTS idx_users_last_heartbeat ON users(last_heartbeat_at);

	CREATE TABLE IF NOT EXISTS pageviews (
		id              BIGSERIAL   PRIMARY KEY,
		created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		event_name      TEXT        NOT NULL,
		path            TEXT        NOT NULL DEFAULT '',
		referrer_domain TEXT        NOT NULL DEFAULT '',
		utm_source      TEXT        NOT NULL DEFAULT '',
		utm_medium      TEXT        NOT NULL DEFAULT '',
		utm_campaign    TEXT        NOT NULL DEFAULT '',
		country         TEXT        NOT NULL DEFAULT '',
		region          TEXT        NOT NULL DEFAULT '',
		device_type     TEXT        NOT NULL DEFAULT '',
		visitor_hash    TEXT        NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_pageviews_created_at ON pageviews (created_at);
	CREATE INDEX IF NOT EXISTS idx_pageviews_event_created ON pageviews (event_name, created_at);

	CREATE TABLE IF NOT EXISTS admin_audit (
		id         BIGSERIAL   PRIMARY KEY,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		action     TEXT        NOT NULL,
		detail     TEXT        NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_admin_audit_created_at ON admin_audit (created_at);
	`)
	return err
}
