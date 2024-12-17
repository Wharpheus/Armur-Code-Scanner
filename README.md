# Armur Code Scanner

This is the official repository for the Armur static scanner built with GO. It uses the best open source tools for static scanning and combines them in a pipeline. It works amazingly for languages like GO, Rust and Python.

This was built after talking to hundreds of Red Teamers, Bug Bounty hunters, security researchers and most importantly, developers to build a solution that automates a huge part of the work required to find security issues and vulnerabilities in code.

Visit armur.ai to use the cloud based version of this same tool.

On the Armur platform, you can use our proprietary AI tool that uses Agents powered by LLMs for security tooling that help provide extremely detailed reports on code security for GO, Python and Rust code and smart contracts code for Solidity, Move and Solana (Rust) - on blockchains

## Features

*   **Multi-Language Support:** Supports scanning of Go, Python, and JavaScript code.
*   **Vulnerability Detection:** Identifies a wide range of vulnerabilities using static analysis tools.
*   **Code Quality Checks:** Performs checks for code style issues and complex functions.
*   **OWASP and SANS Reports:** Generates reports based on OWASP and SANS guidelines.
*   **Advanced Scans:** Detects duplicate code, secrets, infrastructure issues, and SCA issues.
*   **File Scan Support:** Ability to scan individual files for quick analysis.
*   **Task Management:** Uses Asynq for background task processing, and Redis as a task result store.

## Supported Vulnerabilities

Armur Code Scanner detects the following types of vulnerabilities and coding weaknesses, based on the Common Weakness Enumeration (CWE):

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
| CWE-302  |  Improper Authorization  |
| CWE-269 |  Improper Privilege Management  |
| CWE-284 | Improper Access Control  |
| CWE-922 |  Insecure Storage of Sensitive Information  |
| CWE-384  | Session Fixation |
| CWE-613 |  Insufficient Session Expiration |
| CWE-327 | Use of a Broken or Risky Cryptographic Algorithm |
| CWE-330 | Use of Insufficiently Random Values |
| CWE-338 | Use of Cryptographically Weak PRNG |
| CWE-325 | Missing Required Cryptographic Step |
| CWE-200 |  Exposure of Sensitive Information to an Unauthorized Actor  |
| CWE-201  |  Insertion of Sensitive Information into Sent Data    |
| CWE-598   |  Information Exposure Through Query Strings in URL    |
| CWE-209 |   Generation of Error Message Containing Sensitive Information   |
| CWE-310  |  Cryptographic Issues    |
| CWE-502  | Deserialization of Untrusted Data |
| CWE-917  | Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection') |
| CWE-829  | Inclusion of Functionality from Untrusted Control Sphere |
| CWE-434  | Unrestricted Upload of File with Dangerous Type |
| CWE-494 |  Download of Code Without Integrity Check  |
| CWE-611 | Improper Restriction of XML External Entity Reference |
| CWE-918 | Server-Side Request Forgery (SSRF) |
| CWE-862 | Missing Authorization |
| CWE-22 | Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')|
| CWE-73 | External Control of File Name or Path|
| CWE-552 | Unsafe Handling of File Uploads |
| CWE-119 | Improper Restriction of Operations within the Bounds of a Memory Buffer |
| CWE-416 | Use After Free |
| CWE-476 | NULL Pointer Dereference |
| CWE-787 |  Out-of-bounds Write |
| CWE-259 | Use of Hard-coded Password |
| CWE-798 | Use of Hard-coded Credentials |
| CWE-352 | Cross-Site Request Forgery (CSRF) |
| CWE-601 | URL Redirection to Untrusted Site ('Open Redirect') |

| AND MANY MORE |

### Additional Vulnerability Information

In addition to these, Armur Code Scanner also leverages tools such as:

*   **Semgrep:** For detecting various coding patterns and security vulnerabilities.
*   **Gosec:** For Go-specific security issues.
*   **Bandit:** For Python-specific security vulnerabilities.
*   **ESLint:** For detecting JavaScript security issues and code quality problems.
*   **OSV-Scanner:** For identifying Software Composition Analysis issues.
*   **Trufflehog:** For identifying exposed secrets in your codebase.
*   **Checkov:** For identifying Infrastructure as code misconfigurations.
*   **Trivy:** For identifying infrastructure and container vulnerabilities, and secrets.
*   **JSCPD:** For finding duplicated code.
*   **Pydocstyle, Radon, Pylint:** For Python specific code quality issues.
*  **Golint, Govet, Staticcheck, Gocyclo:** For GO specific code quality issues.
* **Vulture:** For identifying dead code in python projects.

## How Armur Code Scanner Works

1.  **API Request:** A scan is triggered via an API request which can be a POST request, containing a git repository URL or a file.
2.  **Task Enqueue:** The API enqueues a scan task using Asynq, with the information about repository URL, language, and scan type in the task payload, and task id.
3. **Repository Cloning:** If the payload contains repo url, it is cloned into a temporary directory.
4. **Scan Execution:** The Asynq worker processes these tasks. Based on task type it will run either a simple scan or an advanced scan, on either the repository directory or the file path.
5.  **Scanning:** The tool will then use appropriate static analysis tools, such as Semgrep, gosec, bandit, eslint etc., and generate scan results.
6.  **Result Storage:** Scan results are stored in a Redis database with a TTL of 24 hours, using task id as a key.
7.  **Status Check:** The scan results can be queried via a task status API using the task id.
8.  **Report Generation:** OWASP and SANS reports can be generated by fetching and reformatting the stored results.

## Getting Started

### Prerequisites

*   Docker and Docker Compose installed on your system.
*   Go installed on your system (for development)

### Running Locally (Development)

1.  **Clone the Repository:**

    ```bash
    git clone <github-repo-link>
    cd <cloned-directory>
    ```

2.  **Start the Development Environment:**

    ```bash
     docker-compose -f docker-compose.dev.yml up --build -d
    ```
    This command does the following:
       * Builds the `armur-tools` image (including installing all go tools).
        * Builds the application image based on `Dockerfile`.
        * Starts the application and Redis containers.

    After running this the application will be available at `http://localhost:4500`.

### Testing with Postman

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

