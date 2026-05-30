-- Remain Faithful schema
-- Run manually: psql $DATABASE_URL -f migrations/001_schema.sql
-- Applied automatically on server startup via the embedded migrate() call.

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    name          TEXT         NOT NULL,
    email         TEXT         NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS relationships (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    partner_id BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       TEXT        NOT NULL DEFAULT 'partner',
    status     TEXT        NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, partner_id),
    CHECK (type   IN ('partner', 'group')),
    CHECK (status IN ('pending', 'accepted', 'declined'))
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
    PRIMARY KEY (group_id, user_id),
    CHECK (role IN ('admin', 'member'))
);

CREATE TABLE IF NOT EXISTS events (
    id        BIGSERIAL   PRIMARY KEY,
    user_id   BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category  TEXT        NOT NULL,
    severity  TEXT        NOT NULL,
    summary   TEXT        NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (severity IN ('low', 'medium', 'high'))
);

CREATE TABLE IF NOT EXISTS alerts (
    id              BIGSERIAL   PRIMARY KEY,
    event_id        BIGINT      NOT NULL REFERENCES events(id)        ON DELETE CASCADE,
    relationship_id BIGINT               REFERENCES relationships(id) ON DELETE SET NULL,
    seen            BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_user_id       ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_relationship  ON alerts(relationship_id);
CREATE INDEX IF NOT EXISTS idx_relationships_user   ON relationships(user_id);
CREATE INDEX IF NOT EXISTS idx_relationships_partner ON relationships(partner_id);
