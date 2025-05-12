# Armur Code Scanner CLI

A powerful command-line interface for the Armur Code Scanner, designed to help developers identify security vulnerabilities, code quality issues, and potential bugs in their codebase.

## Features

- 🔍 **Code Scanning**: Scan repositories or local files for security vulnerabilities
- 📊 **Advanced Analysis**: Enable advanced scanning for deeper code analysis
- 📝 **Report Generation**: Generate OWASP or SANS compliant reports
- 🔄 **Interactive Mode**: User-friendly interactive interface
- 📈 **Scan History**: Track and review past scans
- ⚙️ **Flexible Configuration**: Easy configuration of API and Redis endpoints
- 🔐 **API Audit**: Perform detailed code audits through the API

## Installation

```bash
# Clone the repository
git clone https://github.com/Armur-AI/Armur-Code-Scanner.git

# Navigate to the CLI directory
cd Armur-Code-Scanner/cli

# Build the CLI
go build -o armur-cli

# Move to your PATH (optional)
sudo mv armur-cli /usr/local/bin/
```

## Usage

### Configuration

The CLI can be configured in three ways:

1. **Interactive Mode**:
```bash
armur-cli config
```
This will launch an interactive prompt where you can:
- Select the configuration key to set (API URL, Redis URL, or API Key)
- Enter the new value

2. **Command Line**:
```bash
# Set API URL
armur-cli config api_url https://api.armur.ai

# Set Redis URL
armur-cli config redis_url redis://localhost:6379

# Set API Key
armur-cli config api_key your-api-key-here
```

3. **Environment Variables**:
```bash
export ARMUR_API_URL=https://api.armur.ai
export ARMUR_REDIS_URL=redis://localhost:6379
export ARMUR_API_KEY=your-api-key-here
```

### API Audit

The API audit feature allows you to perform detailed code analysis through an interactive interface:

```bash
armur-cli api
```

This will launch an interactive prompt where you can:
1. Select the audit type:
   - Vulnerability
   - Audit
   - Optimization
   - Codefix
   - Documentation
2. Enter the code/content to analyze
3. Specify token count (e.g., 200)
4. Set temperature (0.0 to 1.0)

The audit results will be displayed in the terminal.

### Code Scanning

```bash
# Scan a repository
armur-cli scan https://github.com/user/repo --language python

# Scan a local file or directory
armur-cli scan /path/to/your/code

# Enable advanced scanning
armur-cli scan <target> --advanced

# Output in JSON format
armur-cli scan <target> --output json
```

### Report Generation

```bash
# Generate an OWASP report
armur-cli report <task-id> --type owasp

# Generate a SANS report
armur-cli report <task-id> --type sans
```

### View Scan Status

```bash
armur-cli status <task-id>
```

## Examples

### Basic Repository Scan

```bash
armur-cli scan https://github.com/user/repo --language python
```

### Local Directory Scan with Advanced Analysis

```bash
armur-cli scan ./my-project --advanced
```

### API Audit Example

```bash
# Start an interactive API audit
armur-cli api

# The tool will prompt you for:
# 1. Audit type (Vulnerability/Audit/Optimization/Codefix/Documentation)
# 2. Code content to analyze
# 3. Token count
# 4. Temperature setting
```

### Generate OWASP Report

```bash
armur-cli report abc123 --type owasp
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please:

1. Check the [documentation](docs/)
2. Open an issue on GitHub
3. Contact our support team at support@armur.com
