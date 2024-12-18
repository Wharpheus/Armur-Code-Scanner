# Armur Code Scanner

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Discord](https://img.shields.io/discord/1021371417134125106)](https://discord.gg/PEycrqvd)

This is the official repository for the Armur static code scanner, built with Go. It leverages the best open-source static analysis tools, combining them into a seamless pipeline for identifying security vulnerabilities and code quality issues. This tool is designed to be efficient and effective, particularly for languages like Go, Python, and JavaScript.

This project was born from conversations with hundreds of red teamers, bug bounty hunters, security researchers, and developers. It aims to automate a significant portion of the work involved in finding security flaws in code.

Visit [armur.ai](https://armur.ai) to use the cloud-based version of this tool, which includes proprietary AI agents powered by LLMs to provide extremely detailed code security reports for Go, Python, JavaScript, and smart contract code for Solidity, Move, and Solana (Rust).

## Key Features

*   **Multi-Language Support**: Scans Go, Python, and JavaScript code.
*   **Comprehensive Vulnerability Detection**: Identifies a broad spectrum of vulnerabilities using various static analysis tools.
*   **Code Quality Analysis**: Performs checks for code style issues, complex functions, and dead code.
*   **OWASP and SANS Compliance**: Generates reports based on OWASP and SANS guidelines.
*   **Advanced Security Scans**: Detects duplicate code, exposed secrets, infrastructure misconfigurations, and SCA vulnerabilities.
*   **Individual File Scanning**: Allows for quick analysis of individual source files.
*   **Asynchronous Task Processing**: Leverages Asynq for background task processing and Redis for result storage.
*   **API Documentation**: Comprehensive API documentation using Swagger/OpenAPI.
*   **Easy Development**: Makefile provided for easy building, testing, and running.

## Supported Vulnerabilities

Armur Code Scanner is capable of detecting the following types of vulnerabilities and coding weaknesses, based on the Common Weakness Enumeration (CWE):

| CWE ID   | Vulnerability Name                                                                                   |
| :------- | :--------------------------------------------------------------------------------------------------- |
| CWE-20   | Improper Input Validation                                                                         |
| CWE-78   | Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')         |
| CWE-79   | Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting') |
| CWE-89   | Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')        |
| CWE-90   | Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')            |
| CWE-94   | Improper Control of Generation of Code ('Code Injection')                                          |
| CWE-400  | Uncontrolled Resource Consumption ('Resource Exhaustion')                                       |
| CWE-287  | Improper Authentication                                                                              |
| CWE-306  | Missing Authentication for Critical Function                                                  |
| CWE-302  | Improper Authorization  |
| CWE-269  | Improper Privilege Management  |
| CWE-284  | Improper Access Control  |
| CWE-922  | Insecure Storage of Sensitive Information  |
| CWE-384  | Session Fixation |
| CWE-613  | Insufficient Session Expiration |
| CWE-327  | Use of a Broken or Risky Cryptographic Algorithm |
| CWE-330  | Use of Insufficiently Random Values |
| CWE-338  | Use of Cryptographically Weak PRNG |
| CWE-325  | Missing Required Cryptographic Step |
| CWE-200  | Exposure of Sensitive Information to an Unauthorized Actor  |
| CWE-201  | Insertion of Sensitive Information into Sent Data    |
| CWE-598  | Information Exposure Through Query Strings in URL    |
| CWE-209  | Generation of Error Message Containing Sensitive Information   |
| CWE-310  | Cryptographic Issues    |
| CWE-502  | Deserialization of Untrusted Data |
| CWE-917  | Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection') |
| CWE-829  | Inclusion of Functionality from Untrusted Control Sphere |
| CWE-434  | Unrestricted Upload of File with Dangerous Type |
| CWE-494  | Download of Code Without Integrity Check  |
| CWE-611  | Improper Restriction of XML External Entity Reference |
| CWE-918  | Server-Side Request Forgery (SSRF) |
| CWE-862  | Missing Authorization |
| CWE-22   | Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')|
| CWE-73   | External Control of File Name or Path|
| CWE-552  | Unsafe Handling of File Uploads |
| CWE-119  | Improper Restriction of Operations within the Bounds of a Memory Buffer |
| CWE-416  | Use After Free |
| CWE-476  | NULL Pointer Dereference |
| CWE-787  | Out-of-bounds Write |
| CWE-259  | Use of Hard-coded Password |
| CWE-798  | Use of Hard-coded Credentials |
| CWE-352  | Cross-Site Request Forgery (CSRF) |
| CWE-601  | URL Redirection to Untrusted Site ('Open Redirect') |
| MANY MORE   |

### Additional Vulnerability Information

In addition to these, Armur Code Scanner also leverages the power of the following open source tools:

*   **Semgrep:** For detecting various coding patterns and security vulnerabilities.
*   **Gosec:** For Go-specific security issues.
*   **Bandit:** For Python-specific security vulnerabilities.
*   **ESLint:** For detecting JavaScript security issues and code quality problems.
*   **OSV-Scanner:** For identifying Software Composition Analysis (SCA) issues.
*   **Trufflehog:** For identifying exposed secrets in your codebase.
*   **Checkov:** For identifying Infrastructure as code misconfigurations.
*   **Trivy:** For identifying infrastructure and container vulnerabilities, and secrets.
*   **JSCPD:** For finding duplicated code.
*   **Pydocstyle, Radon, Pylint:** For Python specific code quality issues.
*   **Golint, Govet, Staticcheck, Gocyclo:** For GO specific code quality issues.
*   **Vulture:** For identifying dead code in Python projects.

## How Armur Code Scanner Works

1.  **API Request:** Initiate a scan using the API, providing a Git repository URL or a file path.
2.  **Task Enqueue:** The API enqueues a scan task using Asynq, including repository URL, language, scan type, and a unique task ID.
3.  **Repository Cloning:** If a repository URL is provided, the tool clones it to a temporary directory.
4.  **Scan Execution:** Asynq worker processes the tasks using relevant static analysis tools like Semgrep, gosec, bandit and eslint.
5.  **Result Storage:** Scan results are stored in Redis with a TTL of 24 hours, using the task ID as the key.
6.  **Status Check:** Query the scan results using the Task Status API with the unique task ID.
7.  **Report Generation:** Generate OWASP and SANS reports by fetching and reformatting the scan results.

## Getting Started

### Prerequisites

*   [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed on your system.
*   [Go](https://golang.org/dl/) installed on your system (for development purposes)

### Running Locally (Development)

1.  **Clone the Repository:**

    ```bash
    git clone <your-github-repo-link>
    cd armur-code-scanner
    ```

2.  **Start the Development Environment:**

    ```bash
    make docker-up
    ```
    OR

    ```bash
    docker-compose up --build -d
    ```

    This command does the following:
    *  Builds the application image based on `Dockerfile`.
    *  Starts the application and Redis containers in development mode using `docker-compose.yml`.
    *  Generates the swagger documentation
    *  After running this, the application will be available at `http://localhost:4500`.
    *  Swagger documentation will be available here `http://localhost:4500/swagger/index.html`

### Testing with Postman

#### A Postman collection is included in the /postman directory for easy API testing.

You can use Postman to send requests to the API endpoints. Here's how:

1.  **API Endpoints:**
    *   **`POST /api/v1/scan/repo`:**
        *   Body:
            ```json
            {
              "repository_url": "https://github.com/go-git/go-git",
              "language": "go"
            }
            ```
        *   Returns a `task_id` upon successful submission.

    *   **`POST /api/v1/advanced-scan/repo`:**
        *   Body:
             ```json
            {
              "repository_url": "https://github.com/go-git/go-git",
              "language": "go"
            }
            ```
         *   Returns a `task_id` upon successful submission.
    *   **`POST /api/v1/scan/file`:**
        *   Select `form-data` and upload the file.
        *   Returns a `task_id` upon successful submission.

    *   **`GET /api/v1/status/:task_id`:**
        *   Replace `:task_id` with the ID from a previous request.
        *   Returns the status of the task or the scan results.

    *   **`GET /api/v1/reports/owasp/:task_id`:**
        *    Replace `:task_id` with the ID from a previous request.
        *  Returns the Owasp report.

    *   **`GET /api/v1/reports/sans/:task_id`:**
        *  Replace `:task_id` with the ID from a previous request.
        *   Returns the SANS report.

### License

**This project is licensed under the MIT License - see the LICENSE file for details.**