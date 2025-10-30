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

# Build server binary with metadata
ARG VERSION=0.0.0-dev
ARG COMMIT=dev
ARG DATE=unknown
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-s -w -buildid= -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" -o /out/armur-codescanner ./cmd/server

########################
# Final minimal image
########################
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /out/armur-codescanner /app/armur-codescanner
EXPOSE 4500
USER nonroot:nonroot

# Default to binding on loopback; override via env at runtime
ENV APP_PORT=4500 \
    BIND_ADDR=0.0.0.0

# No shell, just the binary
# OCI labels
LABEL org.opencontainers.image.title="armur-codescanner" \
      org.opencontainers.image.description="Armur Code Scanner server" \
      org.opencontainers.image.source="https://github.com/armur-ai/Armur-Code-Scanner" \
      org.opencontainers.image.version="$VERSION" \
      org.opencontainers.image.revision="$COMMIT" \
      org.opencontainers.image.created="$DATE"

ENTRYPOINT ["/app/armur-codescanner"]
