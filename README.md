# Armur Code Scanner

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Discord](https://img.shields.io/discord/1021371417134125106)](https://discord.gg/PEycrqvd)

This is the official repository for the Armur static code scanner, built with Go. It uses the best open-source static analysis tools, combining them into a seamless pipeline for identifying security vulnerabilities and code quality issues. This tool is designed to be efficient and effective, particularly for languages like Go, Python, and JavaScript.

This project was born from conversations with hundreds of red teamers, bug bounty hunters, security researchers, and developers. It aims to automate a significant portion of the work involved in finding security flaws in code.

Visit [armur.ai](https://armur.ai) to use the cloud-based version of this tool, which includes proprietary AI agents powered by LLMs to provide extremely detailed code security reports for Go, Python, JavaScript, and smart contract code for Solidity, Move, and Solana (Rust).

## Table of Contents

1. [Key Features](#key-features)
2. [System Requirements](#system-requirements)
3. [Windows Compatibility](#windows-compatibility)
4. [Using the CLI tool](#using-the-armur-cli)
5. [How It Works](#how-armur-code-scanner-works)
6. [Getting Started](#getting-started)
7. [Detailed Running Instructions](#detailed-running-instructions)
8. [API Functions and Endpoints](#api-functions-and-endpoints)
9. [Project Structure](#project-structure)
10. [Supported Languages and Tools](#supported-languages-and-tools)
11. [Supported Vulnerabilities](#supported-vulnerabilities)
12. [Additional Vulnerabilities Information](#additional-vulnerability-information)
13. [Scanning Local Repositories via Mounted Volume](#scanning-local-repositories-via-mounted-volume)
14. [Testing with Postman](#testing-with-postman)
15. [Troubleshooting](#troubleshooting)
16. [License](#license)

## Using the Armur CLI

For a more streamlined experience, you can use the Armur CLI tool. The CLI provides a command-line interface to interact with the Armur Code Scanner, making it easier to scan your codebases directly from your terminal. Following are the steps to quickly get started using it, and for detailed instructions on installation and usage of the CLI, please refer to our [CLI Documentation](/cli/README.md).

### Installing the CLI

These instructions will allow you to quickly setup and install the cli tool for usage.

```bash
# Clone the repository
git clone https://github.com/Armur-AI/Armur-Code-Scanner.git

# Navigate to the CLI directory
cd Armur-Code-Scanner/cli

# Build the CLI
go build -o armur-cli
# verify that the cli has been installed
./armur-cli

# Move to your PATH (optional). This allows you to run armur anywhere in you machine
sudo mv armur-cli /usr/local/bin/
```

### After installing the tool

You can now run a couple of commands that allow you to do code scanning functionality.

If the armur code scanner project is not running on your machine, you can either run that manually using

```bash
1. make docker-up
## or
2. armur-cli docker
```

The you can start by running a scan on a github repository by running the following command.

```bash
armur-cli scan <github-url> --language go
```

## Key Features

- **Multi-Language Support**: Scans Go, Python, JavaScript, and Solidity smart contracts.
- **Comprehensive Vulnerability Detection**: Identifies a broad spectrum of vulnerabilities using various static analysis tools.
- **Code Quality Analysis**: Performs checks for code style issues, complex functions, and dead code.
- **OWASP and SANS Compliance**: Generates reports based on OWASP and SANS guidelines.
- **Advanced Security Scans**: Detects duplicate code, exposed secrets, infrastructure misconfigurations, and SCA vulnerabilities.
- **Individual File Scanning**: Allows for quick analysis of individual source files.
- **Asynchronous Task Processing**: Leverages Asynq for background task processing and Redis for result storage.
- **API Documentation**: Comprehensive API documentation using Swagger/OpenAPI.
- **Easy Development**: Makefile provided for easy building, testing, and running.
- **Docker-Based**: Runs in containers for consistent environments across platforms.

## System Requirements

- **Docker** and **Docker Compose** (recommended for easy setup)
- **Go 1.22+** (for development)
- **Git** (for cloning repositories)
- **4GB RAM minimum** (8GB recommended for large scans)
- **Linux/macOS/Windows with WSL2** (see Windows section below)

## Windows Compatibility

Yes, Armur Code Scanner runs on Windows! However, due to the underlying static analysis tools being primarily Linux-based, we recommend using Docker or WSL2 for the best experience.

### Running on Windows

1. **Using Docker (Recommended):**
   - Install Docker Desktop for Windows
   - Follow the Docker setup instructions below
   - All tools run in Linux containers, so no additional setup needed

2. **Using WSL2:**
   - Install WSL2 and Ubuntu distribution
   - Install Go and Docker in WSL2
   - Run the scanner from within the WSL2 environment

3. **Native Windows (Limited):**
   - Go application runs natively
   - Some tools (bandit, eslint, etc.) may require additional setup or may not work
   - Not recommended for full functionality

## How Armur Code Scanner Works

1. **API Request:** Initiate a scan using the API, providing a Git repository URL or a file path.
2. **Task Enqueue:** The API enqueues a scan task using Asynq, including repository URL, language, scan type, and a unique task ID.
3. **Repository Cloning:** If a repository URL is provided, the tool clones it to a temporary directory.
4. **Scan Execution:** Asynq worker processes the tasks using relevant static analysis tools like Semgrep, gosec, bandit and eslint based on the programming language specified.
5. **Result Storage:** Scan results are stored in Redis with a TTL of 24 hours, using the task ID as the key.
6. **Status Check:** Query the scan results using the Task Status API with the unique task ID.
7. **Report Generation:** Generate OWASP and SANS reports by fetching and reformatting the scan results.

## Getting Started

### Documentation

You can follow out detailed documentation by doing the following steps

```bash
git clone https://github.com/Armur-Ai/Armur-Code-Scanner.git
cd Armur-Code-Scanner/docs-site

npm i --legacy-peer-deps
```

This should open our developer logs at `localhost:3000`, and you can open these from your browser

## Detailed Running Instructions

### Prerequisites

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed on your system.
- [Go](https://golang.org/dl/) installed on your system (for development purposes)

### Running Locally (Development)

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/Armur-Ai/Armur-Code-Scanner.git
    cd Armur-Code-Scanner
    ```

2. **Start the Development Environment:**

    First copy/rename `.env.example` file to `.env` and run the following commands

    ```bash
    make docker-up
    ```

    OR

    ```bash
    docker-compose up --build -d
    ```

This command does the following:

- Builds the application image based on `Dockerfile`.
- Starts the application and Redis containers in development mode using `docker-compose.yml`.
- Generates the swagger documentation
- After running this, the application will be available at `http://localhost:4500`.
- Swagger documentation will be available here `http://localhost:4500/swagger/index.html`

### Alternative Running Methods

#### Using Makefile Commands

```bash
# Build the application
make build

# Run tests
make test

# Run security scans
make sec-scan

# Clean build artifacts
make clean

# Run all CI checks
make ci
```

#### Running Without Docker (Advanced)

If you prefer not to use Docker, you can run the application directly:

1. **Install Dependencies:**

   ```bash
   go mod download
   ```

2. **Start Redis:**

   ```bash
   redis-server
   ```

3. **Run the Application:**

   ```bash
   go run cmd/server/main.go
   ```

**Note:** Running without Docker requires all static analysis tools (gosec, bandit, eslint, etc.) to be installed on your system, which can be complex and platform-dependent.

#### Using the CLI

See the [CLI Documentation](/cli/README.md) for detailed CLI usage instructions.

## API Functions and Endpoints

The Armur Code Scanner provides a REST API with the following endpoints:

### Core Scanning Endpoints

#### `POST /api/v1/scan/repo`

- **Function:** Initiates a basic security scan on a Git repository
- **Parameters:**
  - `repository_url` (string): Git repository URL
  - `language` (string, optional): Programming language (go, py, js, solidity)
- **Returns:** Task ID for tracking scan progress
- **Use Case:** Standard security and code quality analysis

#### `POST /api/v1/advanced-scan/repo`

- **Function:** Runs comprehensive advanced scans including SCA, secrets detection, and infrastructure checks
- **Parameters:** Same as basic scan
- **Returns:** Task ID
- **Use Case:** Deep analysis for production codebases

#### `POST /api/v1/scan/file`

- **Function:** Scans an individual file uploaded via form-data
- **Parameters:**
  - `file` (multipart/form-data): File to scan
- **Returns:** Task ID
- **Use Case:** Quick analysis of single files

#### `POST /api/v1/scan/local`

- **Function:** Scans a local directory mounted in the container
- **Parameters:**
  - `local_path` (string): Path to local directory (e.g., `/armur/repos/my-project`)
  - `language` (string, optional): Programming language
- **Returns:** Task ID
- **Use Case:** Scanning local codebases without Git

### Status and Results

#### `GET /api/v1/status/{task_id}`

- **Function:** Retrieves scan status and results
- **Returns:**
  - `status`: "pending" or "success"
  - `data`: Scan results (when complete)
  - `task_id`: Original task identifier

#### `GET /api/v1/reports/owasp/{task_id}`

- **Function:** Generates OWASP-compliant security report
- **Returns:** Structured report with CWE mappings and remediation advice

#### `GET /api/v1/reports/sans/{task_id}`

- **Function:** Generates SANS Top 25 compliant report
- **Returns:** Report focused on most critical security issues

### Scan Result Categories

All scans return results in these categories:

- **security_issues**: Security vulnerabilities with CWE mappings
- **antipatterns_bugs**: Code quality issues and potential bugs
- **complex_functions**: Functions that may be too complex
- **docstring_absent**: Missing documentation
- **gas_issues**: Solidity gas optimization opportunities
- **build_issues**: Compilation and build problems

Advanced scans additionally include:

- **dead_code**: Unused code detection
- **duplicate_code**: Code duplication analysis
- **secret_detection**: Exposed secrets and credentials
- **infra_security**: Infrastructure misconfigurations
- **sca**: Software Composition Analysis vulnerabilities

## Supported Languages and Tools

Armur Code Scanner supports multiple programming languages, each with specialized static analysis tools:

### Go

- **gosec**: Security issues specific to Go
- **golint**: Code style and linting
- **govet**: Suspicious code constructs
- **staticcheck**: Advanced static analysis
- **gocyclo**: Cyclomatic complexity analysis
- **godeadcode**: Dead code detection

### Python

- **bandit**: Security vulnerabilities
- **pydocstyle**: Docstring style checking
- **radon**: Code complexity and maintainability
- **pylint**: Comprehensive linting
- **vulture**: Dead code detection

### JavaScript

- **ESLint**: Security and code quality rules
- **Custom security rules**: XSS, injection, etc.

### Solidity (Smart Contracts)

- **slither**: Comprehensive static analysis
- **mythril**: Symbolic execution and security analysis
- **oyente**: Symbolic execution
- **securify**: Property verification
- **smartcheck**: Static analysis
- **solc**: Compiler diagnostics and AST analysis

### Cross-Language Tools

- **Semgrep**: Custom security rules and pattern matching
- **OSV-Scanner**: Software Composition Analysis (SCA)
- **Trufflehog**: Secret detection
- **Checkov**: Infrastructure as Code security
- **Trivy**: Container and infrastructure scanning
- **JSCPD**: Code duplication detection

### Project Structure

The codebase is organized into the following main directories:

```bash
├── cmd/                    # Application entry points
│   └── server/main.go      # Main server application
├── internal/               # Private application code
│   ├── api/                # API handlers and routes
│   ├── redis/              # Redis client connection
│   ├── tasks/              # Task implementation for queueing, running and storing results on redis
│   ├── tools/              # Integration with security tools for the different programming languages
│   └── worker/             # Background worker implementation that adds tasks to the queue
├── pkg/                    # Public libraries that can be used by external applications
│   ├── common/             # Shared utilities and constants
│   └── utils.go            # General utility functions
├── docs/                   # Documentation files using swagger
├── postman/                # Postman collection for API testing
├── rule_config/            # Configuration files for security tools
├── shared_tmp/             # Temporary files directory
├── docker-compose.yml      # Docker composition for development
├── Dockerfile              # Dockerfile for the code scanner
├── .env.example            # env example file
└── README.md               # This file

```

### Supported Vulnerabilities

Armur Code Scanner is capable of detecting the following types of vulnerabilities and coding weaknesses, based on the Common Weakness Enumeration (CWE):

| CWE ID    | Vulnerability Name                                                                                                     |
| :-------- | :--------------------------------------------------------------------------------------------------------------------- |
| CWE-20    | Improper Input Validation                                                                                              |
| CWE-78    | Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')                             |
| CWE-79    | Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')                                   |
| CWE-89    | Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')                                   |
| CWE-90    | Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')                                   |
| CWE-94    | Improper Control of Generation of Code ('Code Injection')                                                              |
| CWE-400   | Uncontrolled Resource Consumption ('Resource Exhaustion')                                                              |
| CWE-287   | Improper Authentication                                                                                                |
| CWE-306   | Missing Authentication for Critical Function                                                                           |
| CWE-302   | Improper Authorization                                                                                                 |
| CWE-269   | Improper Privilege Management                                                                                          |
| CWE-284   | Improper Access Control                                                                                                |
| CWE-922   | Insecure Storage of Sensitive Information                                                                              |
| CWE-384   | Session Fixation                                                                                                       |
| CWE-613   | Insufficient Session Expiration                                                                                        |
| CWE-327   | Use of a Broken or Risky Cryptographic Algorithm                                                                       |
| CWE-330   | Use of Insufficiently Random Values                                                                                    |
| CWE-338   | Use of Cryptographically Weak PRNG                                                                                     |
| CWE-325   | Missing Required Cryptographic Step                                                                                    |
| CWE-200   | Exposure of Sensitive Information to an Unauthorized Actor                                                             |
| CWE-201   | Insertion of Sensitive Information into Sent Data                                                                      |
| CWE-598   | Information Exposure Through Query Strings in URL                                                                      |
| CWE-209   | Generation of Error Message Containing Sensitive Information                                                           |
| CWE-310   | Cryptographic Issues                                                                                                   |
| CWE-502   | Deserialization of Untrusted Data                                                                                      |
| CWE-917   | Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection') |
| CWE-829   | Inclusion of Functionality from Untrusted Control Sphere                                                               |
| CWE-434   | Unrestricted Upload of File with Dangerous Type                                                                        |
| CWE-494   | Download of Code Without Integrity Check                                                                               |
| CWE-611   | Improper Restriction of XML External Entity Reference                                                                  |
| CWE-918   | Server-Side Request Forgery (SSRF)                                                                                     |
| CWE-862   | Missing Authorization                                                                                                  |
| CWE-22    | Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')                                         |
| CWE-73    | External Control of File Name or Path                                                                                  |
| CWE-552   | Unsafe Handling of File Uploads                                                                                        |
| CWE-119   | Improper Restriction of Operations within the Bounds of a Memory Buffer                                                |
| CWE-416   | Use After Free                                                                                                         |
| CWE-476   | NULL Pointer Dereference                                                                                               |
| CWE-787   | Out-of-bounds Write                                                                                                    |
| CWE-259   | Use of Hard-coded Password                                                                                             |
| CWE-798   | Use of Hard-coded Credentials                                                                                          |
| CWE-352   | Cross-Site Request Forgery (CSRF)                                                                                      |
| CWE-601   | URL Redirection to Untrusted Site ('Open Redirect')                                                                    |
|          | MANY MORE                                                                                                              |

### Additional Vulnerability Information

In addition to these, Armur Code Scanner also leverages the power of the following open source tools:

- **Semgrep:** For detecting various coding patterns and security vulnerabilities.
- **Gosec:** For Go-specific security issues.
- **Bandit:** For Python-specific security vulnerabilities.
- **ESLint:** For detecting JavaScript security issues and code quality problems.
- **OSV-Scanner:** For identifying Software Composition Analysis (SCA) issues.
- **Trufflehog:** For identifying exposed secrets in your codebase.
- **Checkov:** For identifying Infrastructure as code misconfigurations.
- **Trivy:** For identifying infrastructure and container vulnerabilities, and secrets.
- **JSCPD:** For finding duplicated code.
- **Pydocstyle, Radon, Pylint:** For Python specific code quality issues.
- **Golint, Govet, Staticcheck, Gocyclo:** For GO specific code quality issues.
- **Vulture:** For identifying dead code in Python projects.

### Scanning Local Repositories via Mounted Volume

To enable scanning of local repositories, the Docker Compose configuration mounts a volume from your host machine into the container:

```yaml
volumes:
  - ./shared_tmp:/armur/repos
```

This line tells Docker to take the `shared_tmp` directory from your local machine (host) and mount it inside the container at the path `/armur/repos`. This means any folders you place inside `shared_tmp` will be accessible to the scanner within the container at `/armur/repos`.

To scan a local codebase:

1. Create or copy your project folder into the `shared_tmp` directory (e.g., `shared_tmp/my-local-project`).
2. When making a request to the `/api/v1/scan/local` endpoint, use the container path in your request body:

```json
{
  "local_path": "/armur/repos/my-local-project",
  "language": "go"
}
```

This allows you to scan any local project without needing to push it to a remote repository.

### Testing with Postman

A Postman collection is included in the `/postman` directory for easy API testing.

You can use Postman to send requests to the API endpoints. Here's how:

1. **API Endpoints:**

- **`POST /api/v1/scan/repo`:**

  - Body:

    ```json
    {
      "repository_url": "https://github.com/go-git/go-git",
      "language": "go"
    }
    ```

  - Returns a `task_id` upon successful submission.

- **`POST /api/v1/advanced-scan/repo`:**

  - Body:

    ```json
    {
      "repository_url": "https://github.com/go-git/go-git",
      "language": "go"
    }
    ```

  - Returns a `task_id` upon successful submission.

- **`POST /api/v1/scan/local`:**

  - Body:

    ```json
    {
      "local_path": "/armur/repos/<repo_name>",
      "language": "go"
    }
    ```

  - Returns a `task_id` upon successful submission. Note, place your repo copy under the `shared_tmp` directory.

- **`POST /api/v1/scan/file`:**

  - Select `form-data` and upload the file.
  - Returns a `task_id` upon successful submission.

- **`GET /api/v1/status/:task_id`:**

  - Replace `:task_id` with the ID from a previous request.
  - Returns the status of the task or the scan results.

- **`GET /api/v1/reports/owasp/:task_id`:**

  - Replace `:task_id` with the ID from a previous request.
  - Returns the Owasp report.

- **`GET /api/v1/reports/sans/:task_id`:**
  - Replace `:task_id` with the ID from a previous request.
  - Returns the SANS report.

## Troubleshooting

### Common Issues

#### Docker Issues

- **Port 4500 already in use:** Change the port in `docker-compose.yml` or stop the conflicting service
- **Redis connection failed:** Ensure Redis container is running: `docker ps | grep redis`
- **Permission denied:** On Linux, you may need to run `sudo chmod 666 /var/run/docker.sock`

#### Scan Issues

- **Language not detected:** Specify the language explicitly in API requests
- **Large repository timeout:** Increase timeout in `internal/tasks/queue_task.go`
- **Out of memory:** Increase Docker memory limit or scan smaller codebases

#### Tool-Specific Issues

- **Bandit/Pylint not found:** Ensure Python tools are installed in the container
- **ESLint config errors:** Check if `rule_config/eslint/` files exist
- **Solidity compilation fails:** Verify Solidity version compatibility

### Logs and Debugging

#### View Application Logs

```bash
# Docker logs
docker-compose logs -f app

# Worker logs
docker-compose logs -f
```

#### Enable Debug Mode

Set environment variable: `LOG_LEVEL=debug`

#### Check Redis Data

```bash
# Connect to Redis container
docker exec -it redis redis-cli

# List keys
KEYS *

# Get task result
GET <task_id>
```

### Performance Tuning

- **Concurrent workers:** Adjust `Concurrency` in `cmd/server/main.go`
- **Redis TTL:** Modify TTL in `internal/tasks/result_store.go`
- **Memory limits:** Update Docker Compose memory settings

### Getting Help

- Check existing [GitHub Issues](https://github.com/Armur-AI/Armur-Code-Scanner/issues)
- Join our [Discord](https://discord.gg/PEycrqvd) for community support
- Contact <support@armur.ai> for enterprise support

## License

**This project is licensed under the MIT License - see the LICENSE file for details.**
