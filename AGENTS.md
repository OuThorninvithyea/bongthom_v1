# AGENTS.md — BongThom v1

## Run the project

```bash
# Start PostgreSQL (brew services has a launchd bug — use pg_ctl)
pg_ctl -D /opt/homebrew/var/postgresql@18 start

# Start Redis
redis-server --daemonize yes

# Run (hot reload)
cd admin/kaifin-api && air

# Or manual build + run
cd admin/kaifin-api && go build -o ./tmp/main . && ./tmp/main
```

Server starts on the port defined in `.env` `API_PORT` (currently 9000).

## Layer rules — follow these or it won't compile

Every domain module is 5 files in `internal/mobile/<name>/`. Build order: **model → repo → service → handler → router**.

```
HANDLER    HTTP only. Bind, validate, translate, respond. NEVER: SQL, Redis, business logic.
SERVICE    Business logic + orchestration. Calls repo + Redis + jwt pkg. NEVER: HTTP, SQL.
REPO       Database only. SQL queries. NEVER: HTTP, business logic, Redis, fiber.
```

Handler error responses **must** use the translate pattern:
```go
msg, e_msg := translate.TranslateWithError(c, e.MessageID)
if e_msg != nil { /* handle translate failure */ }
response.NewResponseError(msg, constants.Generic_error, e.Err)
```

## i18n — add keys to all 3 files or it panics

Every `MessageID` used in `NewErrorResponse("key", err)` must exist in:
- `pkg/i18n/localize/en.yaml`
- `pkg/i18n/localize/km.yaml`
- `pkg/i18n/localize/zh.yaml`

Check with: `grep "key" pkg/i18n/localize/*.yaml`

## Database

Two tables: `tbl_users` (21 columns, full audit) and `auth_users` (unused, for refresh tokens later).

Index on `tbl_users(user_name)` is UNIQUE. Trigger `trg_tbl_users_updated_at` auto-stamps `updated_at` on any UPDATE.

Connect: `psql "postgresql://outhorninvuth:postgres@127.0.0.1:5432/postgres?sslmode=disable"`

**.env uses `127.0.0.1` not `localhost`** — macOS resolves `localhost` to IPv6 and PostgreSQL rejects it.

## Redis

Single key: `session:<userID>` = `<UUID>`. Set on login, checked on every request via middleware. No TTL. redis-client.go panics if Redis is unavailable — fail-fast by design.

Connect: `redis-cli -h 127.0.0.1`

## Auth

JWT HS256, 15-minute expiry, self-validating (no DB on requests). Claims include `user_id`, `user_name`, `login_session`, `role_id`.

- Login password check is **plaintext** (SELECT WHERE password = $1). User creation uses bcrypt.
- Middleware protects all routes except `POST /api/v1/admin/auth/login`.
- Token sent as `Authorization: Bearer <token>` header.

## Wired dependencies (main.go order matters)

```
middlewares.NewJwtMiddleware(app, db, rdb)   ← BEFORE routes
handler.NewServiceHandlers(app, db, rdb, ws) ← AFTER middleware
```

Middleware must register first or routes bypass protection.

## Type conventions

| PostgreSQL | Go | Notes |
|-----------|-----|-------|
| bigint/bigserial | int64 | — |
| integer | int | — |
| NULLable columns | Pointer (`*string`, `*time.Time`, `*int`) | — |
| password | `json:"-"` | Never exposed in JSON |
| audit fields (created_by, updated_by, deleted_by) | `json:"-"` | Never exposed |
| Acronyms | UPPERCASE (`ID`, `URL`, `JWT`) | Go convention |

## Interface pattern (every module)

```go
type XxxRepo interface { Method(params) (*Return, *error_responses.ErrorResponse) }
type XxxRepoImpl struct { db *sqlx.DB }

type XxxService interface { Method(params) (*Return, *error_responses.ErrorResponse) }
type XxxServiceImpl struct { Repo XxxRepo; Redis *redis.Client }
```

Repo method names match router verbs: `users.Get("/")` → `repo.GetAll()`, `users.Get("/:id")` → `repo.GetByID()`.

## Response format

```json
// Success
{"success":true, "message":"...", "status_code":2000, "data":{...}}
// Error
{"success":false, "message":"...", "status_code":5000, "data":{"error":"..."}}
// Paginated
{"success":true, "message":"...", "status_code":2000, "data":[...], "page":1, "per_page":20, "total":100}
```

## Architecture reference

Full docs: `.opencode/architecture.md`
