# syntax=docker/dockerfile:1
#
# Production image — build from the repo root:
#   docker build -f docker/prod/Dockerfile -t haves:latest .
# (or `just docker-build`).

###############################################################################
# Builder — compile a static binary
###############################################################################
FROM golang:1.25 AS build
WORKDIR /app
# Cache deps separately from source for faster rebuilds.
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /haves .

###############################################################################
# Runtime — minimal, non-root distroless image
###############################################################################
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /haves /haves
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/haves"]
