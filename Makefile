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
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/armur-codescanner ./cmd/server

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
docker-build:
	docker build --pull --no-cache -t armur/codescanner:dev .

# CI aggregate
.PHONY: ci
ci: deps lint test sec-scan

.PHONY: clean
clean:
	rm -rf bin coverage.* *.out

.PHONY: build run test docker-build docker-up docker-down swagger clean dev prod