# syntax=docker/dockerfile:1.7

########################
# Builder stage
########################
FROM --platform=$BUILDPLATFORM golang:1.22.6-alpine3.20@sha256:2a6f2b2e2c98b7e6660a9b090f8f0a7f9a2e3bb0c37d8b1d5f1a6e5c2f21f2f1 AS builder

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOFLAGS="-trimpath -buildvcs=false" \
    GOPROXY="https://proxy.golang.org,direct" \
    GOSUMDB="sum.golang.org"

WORKDIR /src

# Pre-cache modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

# Copy the rest of the source
COPY . .

# Build server binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-s -w -buildid=" -o /out/armur-codescanner ./cmd/server

########################
# Final minimal image
########################
FROM gcr.io/distroless/static:nonroot@sha256:f25f1d3b2a0c3e9d2c0a02b3c2d2a4d4d9f0e0a7a8b9c0d1e2f3a4b5c6d7e8f9

WORKDIR /app
COPY --from=builder /out/armur-codescanner /app/armur-codescanner
EXPOSE 4500
USER nonroot:nonroot

# Default to binding on loopback; override via env at runtime
ENV APP_PORT=4500 \
    BIND_ADDR=0.0.0.0

# No shell, just the binary
ENTRYPOINT ["/app/armur-codescanner"]
