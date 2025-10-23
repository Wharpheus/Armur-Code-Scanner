---
description: Repository Information Overview
alwaysApply: true
---

# Armur Code Scanner Information

## Summary

Armur Code Scanner is a static code analysis tool built with Go that combines multiple open-source security tools to identify vulnerabilities and code quality issues. It supports scanning Go, Python, JavaScript, and Solidity code, detecting a wide range of security vulnerabilities based on CWE standards.

## Structure

- **cmd/**: Application entry points, including the main server
- **internal/**: Private application code (API handlers, Redis client, tasks, tools integration)
- **pkg/**: Public libraries for external applications
- **docs/**: API documentation using Swagger
- **docs-site/**: Documentation website built with Docusaurus
- **cli/**: Command-line interface tool for interacting with the scanner
- **rule_config/**: Configuration files for security tools
- **fixtures/**: Test fixtures for various languages
- **postman/**: Postman collection for API testing

## Main Repository Components

- **Server**: Main Go application that provides the API and worker functionality
- **CLI**: Separate Go module for command-line interaction
- **Documentation Site**: React-based documentation using Docusaurus

## Projects

### Server Application

**Configuration File**: go.mod

#### Language & Runtime

**Language**: Go
**Version**: 1.22
**Build System**: Go modules
**Package Manager**: Go modules

#### Server Dependencies

**Main Dependencies**:

- github.com/gin-gonic/gin v1.10.0 (Web framework)
- github.com/go-git/go-git/v5 v5.12.0 (Git operations)
- github.com/go-redis/redis/v8 v8.11.5 (Redis client)
- github.com/hibiken/asynq v0.25.0 (Task queue)
- github.com/swaggo/gin-swagger v1.6.0 (API documentation)

#### Build & Installation

```bash
make build
# or
go build -o bin/armur-codescanner ./cmd/server
```

#### Docker

**Dockerfile**: Dockerfile
**Image**: armur/codescanner:dev
**Configuration**: Multi-stage build with golang:1.22.6-alpine3.20 for building and distroless/static:nonroot for runtime

#### Testing

**Framework**: Go testing package
**Test Location**: Throughout the codebase with *_test.go files
**Run Command**:

```bash
make test
# or
go test ./... -race -count=1
```

### CLI Tool

**Configuration File**: cli/go.mod

#### CLI Language & Runtime

**Language**: Go
**Version**: 1.23.4
**Build System**: Go modules
**Package Manager**: Go modules

#### CLI Dependencies

**Main Dependencies**:

- github.com/spf13/cobra v1.8.1 (CLI framework)
- github.com/charmbracelet/huh v0.7.0 (TUI components)
- github.com/fatih/color v1.18.0 (Terminal colors)
- github.com/briandowns/spinner v1.23.1 (Terminal spinner)

#### CLI Build & Installation

```bash
cd cli
go build -o armur-cli
```

### Documentation Site

**Configuration File**: docs-site/package.json

#### Docs Site Language & Runtime

**Language**: JavaScript (React)
**Version**: Node.js >=18.0
**Build System**: npm/pnpm
**Package Manager**: npm/pnpm

#### Docs Site Dependencies

**Main Dependencies**:

- @docusaurus/core v3.9.2
- @docusaurus/preset-classic v3.9.2
- react v19.2.0
- react-dom v19.2.0

#### Docs Site Build & Installation

```bash
cd docs-site
npm i --legacy-peer-deps
npm start
```

## Development Operations

### Running Locally

```bash
# Start the development environment
make docker-up
# or
docker-compose up --build -d
```

### Security Scanning

```bash
make sec-scan
```

### CI Pipeline

```bash
make ci
```
