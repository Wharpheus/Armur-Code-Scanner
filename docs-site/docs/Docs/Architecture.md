# Armur Code Scanner - Technical Documentation

## Overview

Armur Code Scanner is a powerful static code analysis tool built with Go, designed to identify security vulnerabilities and code quality issues across multiple programming languages. The project combines various open-source static analysis tools into a seamless pipeline, providing comprehensive security scanning capabilities.

## Technology Stack

### Core Technologies

- **Go 1.22**: Primary programming language
- **Docker & Docker Compose**: Containerization and orchestration
- **Redis**: Task queue and result storage
- **Asynq**: Background task processing

### Key Dependencies

- **gin-gonic/gin**: Web framework
- **go-git/go-git**: Git operations
- **go-redis/redis**: Redis client
- **hibiken/asynq**: Task queue implementation
- **swaggo/swag**: API documentation
- **swaggo/gin-swagger**: Swagger UI integration

## Architecture

### System Components

1. **API Server**

   - Built with Gin framework
   - RESTful API endpoints
   - Swagger documentation
   - Handles scan requests and status checks

2. **Task Queue System**

   - Asynq for task management
   - Redis for queue storage
   - Asynchronous processing of scan requests

3. **Scanner Engine**

   - Multi-language support (Go, Python, JavaScript)
   - Integration with multiple security tools
   - Configurable scanning rules

4. **Result Storage**
   - Redis-based storage
   - 24-hour TTL for scan results
   - Structured result format

### System Workflow

The Armur Code Scanner operates through a well-defined workflow that ensures efficient and reliable code analysis:

1. **Request Handling**

   - User submits a scan request through one of the API endpoints
   - API server validates the request parameters
   - A unique task ID is generated for tracking
   - Request is enqueued in Redis using Asynq

2. **Task Processing**

   - Background worker picks up the task from the queue
   - For repository scans:
     - Clones the target repository to a temporary directory
     - Analyzes repository structure to determine language
   - For local scans:
     - Accesses the mounted volume containing the code
     - Validates file paths and permissions

3. **Analysis Pipeline**

   - Scanner engine determines required security tools based on language
   - Tools are executed in parallel where possible
   - Results from each tool are collected and normalized
   - Findings are categorized by severity and type
   - Duplicate findings are deduplicated

4. **Result Processing**

   - Processed results are stored in Redis with the task ID
   - Results include:
     - Vulnerability details
     - Code locations
     - Severity levels
     - Remediation suggestions
   - Results are formatted according to OWASP and SANS standards

5. **Response Flow**
   - User can check task status using the task ID
   - When complete, results are retrieved from Redis
   - Results are formatted and returned to the user
   - Temporary files and cloned repositories are cleaned up

### Running with Docker Compose

The project includes a comprehensive Docker Compose setup for easy deployment and development:

1. **Configuration**

   ```yaml
   version: "3.8"
   services:
     app:
       build: .
    container_name: api_service
    restart: always
    command: go run ./cmd/server/main.go
    volumes:
      - ./shared_tmp:/armur/repos
    ports:
      - "${APP_PORT}:4500"
    environment:
      - PYTHONDONTWRITEBYTECODE=1
      - PYTHONUNBUFFERED=1
      - APP_PORT=${APP_PORT}
    env_file:
      - .env
    depends_on:
      - redis_service

     redis_service:
    image: "redis:alpine"
    container_name: redis_service
    restart: always
    ports:
      - "${REDIS_PORT}:6379"
   ```

2. **Starting the Services**

```bash
# Build and start all services
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

3. **Service Access**

   - API Server: http://localhost:4500
   - Swagger UI: http://localhost:4500/swagger/index.html
   - Redis: localhost:6379

4. **Volume Management**

   - Local repositories should be placed in the `shared_tmp` directory
   - This directory is mounted to `/armur/repos` in the container
   - Temporary scan results are stored in Redis

5. **Environment Variables**
   - Copy `.env.example` to `.env`
   - Configure required variables:
     ```
     REDIS_HOST=redis
     REDIS_PORT=6379
     API_PORT=4500
     ```

## Project Structure

View full file structure at [Structure](./ProjectStructure.md)

```bash
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── api/              # API handlers and routes
│   ├── redis/            # Redis client connection
│   ├── tasks/            # Task implementation
│   ├── tools/            # Security tool integrations
│   └── worker/           # Background worker
├── pkg/                  # Public libraries
│   ├── common/           # Shared utilities
│   └── utils.go          # Utility functions
├── docs/                 # Swagger documentation
├── postman/             # API testing collections
├── rule_config/         # Security tool configurations
└── shared_tmp/          # Temporary storage
```

## Security Scanning Capabilities

### Supported Vulnerability Types

1. **Input Validation Issues**

   - CWE-20: Improper Input Validation
   - CWE-78: OS Command Injection
   - CWE-79: Cross-site Scripting

2. **Authentication & Authorization**

   - CWE-287: Improper Authentication
   - CWE-306: Missing Authentication
   - CWE-302: Improper Authorization

3. **Cryptographic Issues**

   - CWE-327: Broken Cryptographic Algorithms
   - CWE-330: Insufficient Random Values
   - CWE-338: Weak PRNG

4. **Data Exposure**
   - CWE-200: Sensitive Information Exposure
   - CWE-201: Sensitive Data in Sent Data
   - CWE-209: Error Message Information Leak

### Integrated Security Tools

- **Semgrep**: Pattern-based vulnerability detection
- **Gosec**: Go-specific security analysis
- **Bandit**: Python security scanning
- **ESLint**: JavaScript code analysis
- **OSV-Scanner**: Software composition analysis
- **Trufflehog**: Secret detection
- **Checkov**: Infrastructure as Code scanning
- **Trivy**: Container and infrastructure scanning

## API Endpoints

### Scan Operations

1. **Repository Scan**

   - `POST /api/v1/scan/repo`
   - Scans entire Git repositories

2. **Advanced Scan**

   - `POST /api/v1/advanced-scan/repo`
   - Comprehensive security analysis

3. **Local Scan**

   - `POST /api/v1/scan/local`
   - Scans local codebases

4. **File Scan**
   - `POST /api/v1/scan/file`
   - Individual file analysis

### Status Operations

- `GET /api/v1/status/:task_id`
  - Check scan progress and results

## Development Setup

### Prerequisites

- Docker and Docker Compose
- Go 1.22 or later
- Git

### Local Development

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Run `make docker-up` or `docker-compose up --build -d`
4. Access the application at `http://localhost:4500`
5. View API documentation at `http://localhost:4500/swagger/index.html`

## Deployment

### Docker Deployment

The application is containerized using Docker:

- Multi-stage build process
- Optimized for production
- Environment variable configuration
- Volume mounting for local repository access

### Environment Configuration

Required environment variables:

- Redis connection details
- API configuration
- Security tool settings
- Storage paths

## Contributing

### Development Guidelines

1. Follow Go best practices
2. Write comprehensive tests
3. Update documentation
4. Follow semantic versioning

## License

MIT License - See LICENSE file for details

## Support

- Discord community
- GitHub issues
- Documentation updates
- Security reporting
