# Remain Faithful — Backend API

REST API for the Remain Faithful iOS accountability app. Built with Go, gorilla/mux, lib/pq, and PostgreSQL.

## Quick Start (Docker)

```bash
cd backend
docker-compose up --build
```

The server starts on **http://localhost:8080**. The database schema is applied automatically on first boot; no manual migration step needed.

To stop and remove volumes:

```bash
docker-compose down -v
```

## Local Development (no Docker)

Prerequisites: Go 1.22+, PostgreSQL 14+.

```bash
# 1. Create the database
createdb remainfaithful

# 2. Configure environment
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=remainfaithful
export JWT_SECRET=dev-secret-change-me
export PORT=8080

# 3. Download dependencies and run
cd backend
go mod download
go run ./cmd/server
```

## Environment Variables

| Variable        | Default                           | Description                                     |
|-----------------|-----------------------------------|-------------------------------------------------|
| `PORT`          | `8080`                            | HTTP listen port                                |
| `DATABASE_URL`  | *(unset)*                         | Full Postgres DSN — overrides all `DB_*` vars   |
| `DB_HOST`       | `localhost`                       | PostgreSQL host                                 |
| `DB_PORT`       | `5432`                            | PostgreSQL port                                 |
| `DB_USER`       | `postgres`                        | PostgreSQL user                                 |
| `DB_PASSWORD`   | `postgres`                        | PostgreSQL password                             |
| `DB_NAME`       | `remainfaithful`                  | PostgreSQL database                             |
| `JWT_SECRET`    | `dev-secret-change-in-production` | HS256 signing secret — **always override in production** |

## API Reference

All routes except `/health`, `/auth/register`, and `/auth/login` require:

```
Authorization: Bearer <token>
```

---

### Health

```
GET /health
```
Returns `{"status":"ok"}`. Useful for load-balancer probes.

---

### Auth

#### Register
```
POST /auth/register
```
```json
{ "name": "Jeff Brewer", "email": "jeff@example.com", "password": "secure123" }
```
Returns `201` with the created user (no token — call `/auth/login` next).

#### Login
```
POST /auth/login
```
```json
{ "email": "jeff@example.com", "password": "secure123" }
```
Returns:
```json
{
  "token": "<jwt>",
  "user": { "id": 1, "name": "Jeff Brewer", "email": "jeff@example.com" }
}
```
Tokens expire after **24 hours**.

---

### Users

#### Get my profile
```
GET /users/me
```

---

### Relationships

#### Create a partner relationship
```
POST /relationships
```
```json
{ "partner_email": "partner@example.com", "type": "partner" }
```
`type` can be `"partner"` or `"group"`. New relationships start with `status: "pending"`.

#### List my relationships
```
GET /relationships
```
Returns all relationships where the current user is `user_id`, joined with the partner's profile.

---

### Groups

#### Create a group
```
POST /groups
```
```json
{ "name": "Iron Brotherhood" }
```
The creator is automatically added as `admin`.

#### Get group with members
```
GET /groups/:id
```

#### Invite a member (admin only)
```
POST /groups/:id/invite
```
```json
{ "user_email": "newmember@example.com" }
```

---

### Events

Events are flagged monitoring detections submitted by the iOS app.

#### Report an event
```
POST /events
```
```json
{
  "category":  "adult_content",
  "severity":  "high",
  "summary":   "Flagged during browsing session",
  "timestamp": "2024-05-01T14:30:00Z"
}
```
`timestamp` is optional (defaults to `NOW()`). On creation, an **alert** is automatically
fanned out to every accepted accountability partner.

Valid `severity` values: `low`, `medium`, `high`.

#### List my events
```
GET /events
```
Returns up to 100 most-recent events for the authenticated user.

---

### Alerts

#### List alerts I've received
```
GET /alerts
```
Returns alerts sent to the current user by their partners — i.e., events created by people
the current user holds accountable. Each alert embeds the full event payload.

---

## Database Schema

```
users            id · name · email · password_hash · created_at
relationships    id · user_id → users · partner_id → users · type · status · created_at
groups           id · name · created_at
group_members    group_id → groups · user_id → users · role · joined_at
events           id · user_id → users · category · severity · summary · timestamp
alerts           id · event_id → events · relationship_id → relationships · seen · created_at
```

The schema is applied idempotently (all `CREATE TABLE IF NOT EXISTS`) so restarting the server against an existing database is safe.

## Running Tests

```bash
cd backend
go test ./...
```

(Unit/integration tests are not included in this scaffold — add them in `internal/handler/*_test.go`.)

## Project Layout

```
backend/
├── cmd/server/main.go          # server entry-point, router, DB init
├── internal/
│   ├── auth/
│   │   ├── jwt.go              # token sign / parse
│   │   └── middleware.go       # HTTP auth middleware
│   └── handler/
│       ├── handler.go          # shared H struct + JSON helpers
│       ├── auth.go             # /auth/register, /auth/login
│       ├── users.go            # /users/me
│       ├── relationships.go    # /relationships
│       ├── groups.go           # /groups, /groups/:id, /groups/:id/invite
│       ├── events.go           # /events
│       └── alerts.go           # /alerts
├── migrations/001_schema.sql   # standalone SQL reference
├── Dockerfile
├── docker-compose.yml
└── README.md
```
