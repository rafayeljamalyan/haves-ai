# CLAUDE.md

## What this is

`haves` is a Go web service on [Gin](https://github.com/gin-gonic/gin) backed by Postgres via [pgx](https://github.com/jackc/pgx) v5. Module path is `haves` (internal imports: `haves/internal/...`).

See `README.md` for commands, configuration, the env-file split, Docker hot-reload gotchas, endpoints, and layout. Tasks run through `just` (`just` with no args lists recipes). Single test: `go test ./internal/server -run TestName -race`.

## Architecture & conventions

- **Dependency flow** (`main.go`): open the DB pool (`db.New` / `db.DSNFromEnv`), **fail fast** if Postgres is unreachable, then `server.New(database)`. The DB handle is injected, not global.
- **`internal/db`** — `DB` wraps a concurrency-safe `*pgxpool.Pool`; `New` pings on open. `DSNFromEnv` prefers `DATABASE_URL`, else composes from `POSTGRES_*` vars.
- **`internal/server`** — `New(*db.DB)` builds the engine. DB-dependent handlers hang off the `handler` struct; standalone ones are plain funcs. New API routes go under the `/api/v1` group.
- **Response envelope** — `respond` and `respondError` (`response.go`) are the **only** ways any handler may write a response, including the `/health` and `/ready` probes. Never call `c.JSON` (or `c.String`/`c.Data`/etc.) directly — that bypasses the envelope. Every response has the shape `{"meta": {"error": ...}, "data": ...}`:
  - `respond(c, status, data)` — success: payload in `data`, `meta.error` is `null`.
  - `respondError(c, status, code, msg)` — failure: `data` is `null`, `meta.error` is `{"message": msg, "code": code}` where `code` is an app-level `int` distinct from the HTTP `status`.
