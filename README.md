# haves

A Go web service built with the [Gin](https://github.com/gin-gonic/gin) framework.

## Requirements

- Go 1.25+ (the toolchain in `go.mod` will auto-fetch if needed)

## Run

```bash
go run .
```

The server listens on `:8080`.

## Endpoints

| Method | Path             | Description        |
|--------|------------------|--------------------|
| GET    | `/health`        | Health check       |
| GET    | `/api/v1/ping`   | Returns `pong`     |

```bash
curl localhost:8080/health
curl localhost:8080/api/v1/ping
```

## Layout

```
main.go                  # entrypoint
internal/server/         # router & handlers
```
