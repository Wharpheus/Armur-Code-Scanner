SHELL := /bin/bash

# Tool versions (pin for reproducibility)
GOLANGCI_LINT_VERSION := v1.60.3
GOSEC_VERSION := v2.20.0
OSV_SCANNER_VERSION := v1.7.4
TRIVY_VERSION := 0.55.2

BIN := $(PWD)/bin
export PATH := $(BIN):$(PATH)
GO := go
GOFLAGS := -trimpath -buildvcs=false
LDFLAGS := -s -w -buildid=
PKG := ./...

# Version metadata
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.0.0-dev")
APP_NAME ?= armur-codescanner
BIN_DIR ?= dist
BUILD_MAIN ?= ./cmd/server

# Inject version info into main package if variables exist
# Adjust -X keys to match your main package vars (main.version, main.commit, main.date)
LDFLAGS := $(LDFLAGS) -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Default target
.PHONY: all
all: ci

# Dependencies and tools
.PHONY: deps
deps:
	GOPROXY=https://proxy.golang.org,direct GOSUMDB=sum.golang.org $(GO) mod download
	$(GO) mod verify

$(BIN):
	@mkdir -p $(BIN)

.PHONY: tools
tools: $(BIN)
	GOBIN=$(BIN) $(GO) install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	GOBIN=$(BIN) $(GO) install github.com/google/osv-scanner/cmd/osv-scanner@$(OSV_SCANNER_VERSION)
	curl -sSL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b $(BIN) v$(TRIVY_VERSION)

# Lint and static analysis
.PHONY: lint
lint:
	$(GO) vet $(PKG)

# Tests
.PHONY: test
test:
	$(GO) test $(PKG) -race -count=1 -coverprofile=coverage.out -covermode=atomic

# Build
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) $(BUILD_MAIN)

# Security scans
.PHONY: sec-scan
sec-scan: tools
	gosec ./...
	osv-scanner -r . || true
	trivy fs --exit-code 0 --severity HIGH,CRITICAL --skip-dirs docs-site .

# Docker
.PHONY: docker-up
# Run docker compose up dev
docker-up:
	@echo "Starting Docker containers in dev mode..."
	docker-compose up --build -d

# Stop docker containers
.PHONY: docker-down
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

.PHONY: docker-build
IMAGE ?= armur/armur-codescanner:$(VERSION)
docker-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg DATE=$(DATE) -t $(IMAGE) .

.PHONY: dockerx
# Multi-arch build using buildx
PLATFORMS ?= linux/amd64,linux/arm64
dockerx:
	docker buildx build --platform $(PLATFORMS) --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg DATE=$(DATE) -t $(IMAGE) .

# CI aggregate
.PHONY: ci
ci: deps lint test sec-scan

.PHONY: clean
clean:
	rm -rf bin coverage.* *.out

.PHONY: build run test docker-build docker-up docker-down swagger clean dev prod