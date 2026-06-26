# haves

A Go web service built with the [Gin](https://github.com/gin-gonic/gin) framework.

## Requirements

- Go 1.25+ (the toolchain in `go.mod` auto-fetches if needed)
- [`just`](https://github.com/casey/just) — task runner (`brew install just`)
- Docker + Docker Compose (for the containerized workflow)

Run `just` with no arguments to list all recipes.

## Dev workflow

Two equivalent hot-reload loops — pick one:

### Docker (containerized, prod-parity)

```bash
just docker-dev        # build image + run with hot reload (foreground)
just docker-dev-d      # same, detached
just docker-dev-down   # stop
just docker-dev-logs   # tail logs
just docker-dev-sh     # shell into the container
```

Edit any `.go` file and the server rebuilds and restarts in ~3s.

### Host-native (fastest)

```bash
just install-air   # one-time: install the air reloader
just dev           # hot reload directly on the host
```

Either way the server listens on `:8080`.

## Other recipes

```bash
just run       # run once (no reload)
just build     # build static binary into ./bin
just test      # go test ./... -race
just check     # fmt + vet + tidy
```

## Configuration

Config comes from environment variables. Copy the example file and edit it:

```bash
cp .env.example .env.development.local   # gitignored
```

The Docker dev stack reads this file two ways (both in
`docker/dev/docker-compose.yml` + the justfile):

- **`env_file`** injects the variables (Postgres credentials, etc.) into the
  containers.
- **`--env-file`** (added to the `docker-dev` recipes) feeds `${...}`
  interpolation — specifically the published `PORT`, which `env_file` alone
  can't reach.

Compose overrides the container-specific values (`POSTGRES_HOST=postgres`) so
the file can keep host-oriented defaults (`POSTGRES_HOST=localhost`) for native
runs. The app itself reads its database config from the process environment via
`internal/db`.

| Variable            | Default     | Description                          |
|---------------------|-------------|--------------------------------------|
| `PORT`              | `8080`      | Published host port (container serves 8080) |
| `DATABASE_URL`      | —           | Full Postgres DSN; overrides the parts below |
| `POSTGRES_HOST`     | `localhost` | (in Docker: `postgres`)              |
| `POSTGRES_PORT`     | `5432`      |                                      |
| `POSTGRES_USER`     | `haves`     |                                      |
| `POSTGRES_PASSWORD` | `haves`     |                                      |
| `POSTGRES_DB`       | `haves`     |                                      |
| `POSTGRES_SSLMODE`  | `disable`   |                                      |

## Database

The service connects to Postgres on startup (it fails fast if the database is
unreachable). The dev compose stack runs a `postgres:17-alpine` container, and
the app waits for it to report healthy before starting. In Docker the database
is reachable at host `postgres`; from the host it's exposed on
`localhost:5432`.

## Production image

A dedicated multi-stage build (`docker/prod/Dockerfile`) on distroless (~24 MB):

```bash
just docker-build      # build haves:latest from docker/prod/Dockerfile
just docker-run        # run it on :8080
```

## Endpoints

| Method | Path             | Description                    |
|--------|------------------|--------------------------------|
| GET    | `/health`        | Liveness — process is up       |
| GET    | `/ready`         | Readiness — pings the database |
| GET    | `/api/v1/ping`   | Returns `pong`                 |

```bash
curl localhost:8080/health
curl localhost:8080/ready
curl localhost:8080/api/v1/ping
```

## Layout

```
main.go                       # entrypoint
internal/server/              # router & handlers
internal/db/                  # Postgres connection pool (pgx) + env config
docker/dev/Dockerfile         # dev image: base → haves-dev (air hot reload)
docker/dev/docker-compose.yml # dev stack (app + postgres, bind mount + air)
docker/prod/Dockerfile        # prod image: build → distroless (static, non-root)
.air.toml                     # hot-reload config
.env.example                  # config template — copy to .env for host runs
justfile                      # task runner
```

## Notes / gotchas

- **Hot reload in Docker on macOS/Windows** uses air's polling mode
  (`poll = true`) because filesystem events don't cross the Docker bind
  mount. The Go build/module caches are pre-warmed into the dev image so
  reloads only recompile the changed package (fast, low memory). After
  changing dependencies, re-run `just docker-dev` to rebuild the image.
- If a reload ever seems stuck serving an old response, check for a stale
  server still holding the port: `lsof -nP -iTCP:8080 -sTCP:LISTEN`.
