-- Admin audit log and cookieless website analytics.
-- Applied on server startup from the embedded migrate() in cmd/server.
-- Pageview rows are also deleted after 13 months by the retention job.

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

CREATE INDEX IF NOT EXISTS idx_pageviews_created_at
    ON pageviews (created_at);

CREATE INDEX IF NOT EXISTS idx_pageviews_event_created
    ON pageviews (event_name, created_at);

CREATE TABLE IF NOT EXISTS admin_audit (
    id         BIGSERIAL   PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    action     TEXT        NOT NULL,
    detail     TEXT        NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_admin_audit_created_at
    ON admin_audit (created_at);

-- Retention. The API also runs this on a daily timer.
DELETE FROM pageviews WHERE created_at < NOW() - INTERVAL '13 months';
