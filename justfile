# haves — task runner. Run `just` to list recipes.

set shell := ["bash", "-uc"]


# Default: show available recipes.
default:
    @just --list

# --- Local (host) -----------------------------------------------------------

# Run the server locally.
run:
    go run .

# Hot-reload locally (requires `air`; `just install-air` to get it).
dev:
    air -c .air.toml

# Build the binary into ./bin.
build:
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/haves .

# Run tests.
test:
    go test ./... -race -count=1

# Format, vet, and tidy.
fmt:
    go fmt ./...

vet:
    go vet ./...

tidy:
    go mod tidy

# Format + vet + tidy in one shot.
check: fmt vet tidy

# --- Docker dev -------------------------------------------------------------

# Dev compose file lives under docker/dev/.
dev-compose := "docker compose -f docker/dev/docker-compose.yml -p haves-dev"

# Start the dev stack with hot reload (foreground, logs attached).
docker-dev:
    {{dev-compose}} up --build

# Same, detached.
docker-dev-d:
    {{dev-compose}} up --build -d

# Stop and remove the dev stack.
docker-dev-down:
    {{dev-compose}} down

# Tail logs from the running stack.
docker-dev-logs:
    {{dev-compose}} logs -f

# Open a shell inside the running app container.
docker-dev-sh:
    {{dev-compose}} exec app bash

# --- Docker prod ------------------------------------------------------------

# Build the minimal production image.
docker-build tag="haves:latest":
    docker build -f docker/prod/Dockerfile -t {{tag}} .

# Run the production image.
docker-run tag="haves:latest":
    docker run --rm -p 8080:8080 {{tag}}
