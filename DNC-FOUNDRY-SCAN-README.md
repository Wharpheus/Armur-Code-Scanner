# DNC Foundry Solidity Project Scan Setup

This setup allows you to scan the external dnc-foundry Solidity project located at `/home/wharpheus/tools/domain-ninja/dnc-foundry/src` using the Armur Code Scanner API server.

## Files Created

1. **`dnc-foundry-scan-config.json`** - Configuration file containing scan parameters
2. **`scan-dnc-foundry.sh`** - Executable script to initiate the scan

## Prerequisites

1. **API Server Running**: Ensure the Armur Code Scanner API server is running on the configured endpoint (default: `http://localhost:8080`).

   ```bash
   # Navigate to the Armur-Code-Scanner directory and start the server
   cd /home/wharpheus/tools/Armur-Code-Scanner
   go run cmd/server/main.go
   ```

2. **Path Accessibility**: The target directory `/home/wharpheus/tools/domain-ninja/dnc-foundry/src` must be accessible to the API server process.

3. **Dependencies**: Install `jq` for optimal JSON parsing (optional but recommended):

   ```bash
   sudo apt-get update && sudo apt-get install jq
   ```

## Configuration Details

The configuration file (`dnc-foundry-scan-config.json`) specifies:

- **API Endpoint**: `http://localhost:8080/api/v1`
- **Scan Type**: Local directory scan
- **Target Path**: `/home/wharpheus/tools/domain-ninja/dnc-foundry/src`
- **Language**: Solidity
- **Expected Tools**: Slither, Mythril, Oyente, Securify, SmartCheck, GasOptimizer

## Usage

### Method 1: Using the Automated Script

```bash
./scan-dnc-foundry.sh
```

The script will:

- Load configuration from `dnc-foundry-scan-config.json`
- Validate the target path exists
- Send a POST request to `/api/v1/scan/local`
- Return the task ID for tracking scan progress

### Method 2: Manual API Call

If you prefer to make the API call directly:

```bash
curl -X POST http://localhost:8080/api/v1/scan/local \
  -H "Content-Type: application/json" \
  -d '{
    "local_path": "/home/wharpheus/tools/domain-ninja/dnc-foundry/src",
    "language": "solidity"
  }'
```

## Monitoring Scan Progress

Once the scan is initiated, you'll receive a `task_id`. Use it to check the scan status:

### Get Status

```bash
curl http://localhost:8080/api/v1/status/TASK_ID_HERE
```

### Get OWASP Report

```bash
curl http://localhost:8080/api/v1/reports/owasp/TASK_ID_HERE
```

### Get SANS Report

```bash
curl http://localhost:8080/api/v1/reports/sans/TASK_ID_HERE
```

## Expected Scan Results

The scan will run multiple Solidity security analysis tools and return results categorized as:

- **Security Issues**: High, medium, and low severity vulnerabilities
- **Antipatterns/Bugs**: Code patterns that may cause issues
- **Complex Functions**: Functions with high cyclomatic complexity
- **Missing Documentation**: Functions without proper docstrings

## Troubleshooting

1. **API Server Not Running**: Make sure the server is started before initiating scans.
2. **Path Not Accessible**: Verify the target directory exists and has appropriate permissions.
3. **Network Issues**: Ensure the API endpoint in the config matches your server configuration.
4. **Port Conflicts**: Default port is 4500 (check `cmd/server/main.go`).

## Alternative: Using CLI Instead of API

If the API server is unavailable, you can scan the directory using the CLI:

```bash
cd /home/wharpheus/tools/Armur-Code-Scanner/cli
go run main.go scan /home/wharpheus/tools/domain-ninja/dnc-foundry/src --language solidity --advanced
```

## File Structure of Target Directory

The target directory contains Solidity smart contract files:

- `DomainNinjaCoin.ultra-optimized.sol` - Main contract
- `SOVCLINECoin.sol` - Contract implementation
- `TokenOwnershipManager.sol` - Token management
- Various supporting files and documentation
