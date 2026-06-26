# syntax=docker/dockerfile:1

###############################################################################
# Base — shared Go toolchain
###############################################################################
FROM golang:1.25 AS base
WORKDIR /app
# Cache deps separately from source for faster rebuilds.
COPY go.mod go.sum ./
RUN go mod download

###############################################################################
# Dev — hot reload via air; source is bind-mounted at runtime (compose)
###############################################################################
FROM base AS haves-dev
RUN go install github.com/air-verse/air@latest
COPY . .
# Pre-warm GOCACHE so hot-reload rebuilds only recompile the changed package
# (the heavy deps stay cached) — keeps reloads fast and low-memory.
RUN go build -o /dev/null ./...
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
