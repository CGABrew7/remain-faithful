package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	rfauth "remain-faithful/backend/internal/auth"
	"remain-faithful/backend/internal/anthropic"
	"remain-faithful/backend/internal/apns"
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

	if os.Getenv("JWT_SECRET") == "" {
		log.Println("WARNING: JWT_SECRET not set — using insecure dev default. Set it before production use.")
	}

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

	emailClient   := email.New()
	claudeClient  := anthropic.New()
	stripeClient  := payment.New()

	h := &handler.H{DB: db, APNS: apnsClient, Email: emailClient, Claude: claudeClient, Stripe: stripeClient}
	srv := &http.Server{
		Addr:         ":" + port(),
		Handler:      routes(h),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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

	origins := strings.Split(getenv("ALLOWED_ORIGINS", "https://remainfaithful.com"), ",")
	r.Use(corsMiddleware(origins))

	// Health check (unauthenticated)
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	}).Methods(http.MethodGet)

	// Tier 3 classification (unauthenticated — text only, no PII)
	r.HandleFunc("/classify", h.Classify).Methods(http.MethodPost)

	// Public auth routes
	r.HandleFunc("/auth/register",       h.Register).Methods(http.MethodPost)
	r.HandleFunc("/auth/login",          h.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/apple",          h.AppleSignIn).Methods(http.MethodPost)
	r.HandleFunc("/auth/google",         h.GoogleSignIn).Methods(http.MethodPost)
	r.HandleFunc("/auth/forgot-password", h.ForgotPassword).Methods(http.MethodPost)
	r.HandleFunc("/auth/reset-password", h.ResetPassword).Methods(http.MethodPost)

	// Protected routes — JWT required
	api := r.NewRoute().Subrouter()
	api.Use(rfauth.Middleware)

	api.HandleFunc("/users/me",             h.GetMe).Methods(http.MethodGet)
	api.HandleFunc("/relationships",               h.CreateRelationship).Methods(http.MethodPost)
	api.HandleFunc("/relationships",               h.ListRelationships).Methods(http.MethodGet)
	api.HandleFunc("/relationships/invite",        h.InvitePartner).Methods(http.MethodPost)
	api.HandleFunc("/relationships/accept-invite", h.AcceptPartnerInvite).Methods(http.MethodPost)
	api.HandleFunc("/groups",                    h.CreateGroup).Methods(http.MethodPost)
	api.HandleFunc("/groups/{id}",               h.GetGroup).Methods(http.MethodGet)
	api.HandleFunc("/groups/{id}/invite",        h.InviteMember).Methods(http.MethodPost)
	api.HandleFunc("/groups/{id}/email-invite",  h.GroupEmailInvite).Methods(http.MethodPost)
	api.HandleFunc("/events",               h.CreateEvent).Methods(http.MethodPost)
	api.HandleFunc("/events",               h.ListEvents).Methods(http.MethodGet)
	api.HandleFunc("/alerts",               h.ListAlerts).Methods(http.MethodGet)
	api.HandleFunc("/alerts/count",         h.AlertUnreadCount).Methods(http.MethodGet)
	api.HandleFunc("/alerts/mark-seen",     h.MarkAlertsSeen).Methods(http.MethodPost)
	api.HandleFunc("/users/device-token",                  h.RegisterDeviceToken).Methods(http.MethodPost)
	api.HandleFunc("/panic",                               h.SendPanicAlert).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh",                        h.RefreshToken).Methods(http.MethodPost)
	api.HandleFunc("/donations/create-checkout-session",   h.CreateCheckoutSession).Methods(http.MethodPost)

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
		id         BIGSERIAL   PRIMARY KEY,
		user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token      TEXT        NOT NULL,
		platform   TEXT        NOT NULL DEFAULT 'ios',
		is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE (user_id, token)
	);
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
	`)
	return err
}
