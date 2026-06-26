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

## Production image

A dedicated multi-stage build (`docker/prod/Dockerfile`) on distroless (~24 MB):

```bash
just docker-build      # build haves:latest from docker/prod/Dockerfile
just docker-run        # run it on :8080
```

## Endpoints

| Method | Path             | Description     |
|--------|------------------|-----------------|
| GET    | `/health`        | Health check    |
| GET    | `/api/v1/ping`   | Returns `pong`  |

```bash
curl localhost:8080/health
curl localhost:8080/api/v1/ping
```

## Layout

```
main.go                       # entrypoint
internal/server/              # router & handlers
docker/dev/Dockerfile         # dev image: base → haves-dev (air hot reload)
docker/dev/docker-compose.yml # dev stack (bind mount + air hot reload)
docker/prod/Dockerfile        # prod image: build → distroless (static, non-root)
.air.toml                     # hot-reload config
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
