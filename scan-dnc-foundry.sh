#!/bin/bash

# Script to initiate a scan of the dnc-foundry Solidity project using the Armur Code Scanner API
# Loads configuration from dnc-foundry-scan-config.json

# Configuration file
CONFIG_FILE="dnc-foundry-scan-config.json"

# Check if config file exists
if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "Error: Configuration file '$CONFIG_FILE' not found."
    exit 1
fi

# Load configuration using jq or parse manually
if command -v jq &> /dev/null; then
    API_ENDPOINT=$(jq -r '.api_endpoint' "$CONFIG_FILE")
    LOCAL_PATH=$(jq -r '.target.local_path' "$CONFIG_FILE")
    LANGUAGE=$(jq -r '.target.language' "$CONFIG_FILE")
    SCAN_MODE=$(jq -r '.scan_mode' "$CONFIG_FILE")
else
    echo "Warning: jq not found, using grep. JSON parsing may be limited."
    # Basic parsing without jq
    API_ENDPOINT=$(grep '"api_endpoint"' "$CONFIG_FILE" | cut -d'"' -f4)
    LOCAL_PATH=$(grep '"local_path"' "$CONFIG_FILE" | cut -d'"' -f4)
    LANGUAGE=$(grep '"language"' "$CONFIG_FILE" | sed 's/.*"language":\s*"\([^"]*\)".*/\1/')
    SCAN_MODE=$(grep '"scan_mode"' "$CONFIG_FILE" | sed 's/.*"scan_mode":\s*"\([^"]*\)".*/\1/')
fi

echo "Configuration loaded:"
echo "API Endpoint: $API_ENDPOINT"
echo "Local Path: $LOCAL_PATH"
echo "Language: $LANGUAGE"
echo "Scan Mode: $SCAN_MODE"
echo ""

# Check if path is accessible
if [[ ! -d "$LOCAL_PATH" ]]; then
    echo "Error: Target path '$LOCAL_PATH' does not exist or is not a directory."
    exit 1
fi

echo "Target directory exists and is accessible."
echo ""

# Prepare API request
SCAN_ENDPOINT="$API_ENDPOINT/scan/local"

# Prepare JSON payload
PAYLOAD=$(cat <<EOF
{
  "local_path": "$LOCAL_PATH",
  "language": "$LANGUAGE"
}
EOF
)

echo "Initiating scan..."
echo "Endpoint: $SCAN_ENDPOINT"
echo "Payload: $PAYLOAD"
echo ""

# Make API request
RESPONSE=$(curl -s -X POST "$SCAN_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD")

# Check if request was successful
if [[ $? -ne 0 ]]; then
    echo "Error: Failed to make API request. Is the API server running?"
    echo "Make sure to start the Armur Code Scanner API server on the configured endpoint."
    exit 1
fi

# Extract task ID from response
TASK_ID=$(echo "$RESPONSE" | grep -o '"task_id":"[^"]*"' | cut -d'"' -f4)

if [[ -n "$TASK_ID" ]]; then
    echo "Scan initiated successfully!"
    echo "Task ID: $TASK_ID"
    echo ""
    echo "To check the status of your scan, use:"
    echo "curl $API_ENDPOINT/status/$TASK_ID"
    echo ""
    echo "To get OWASP report:"
    echo "curl $API_ENDPOINT/reports/owasp/$TASK_ID"
else
    echo "Error: Failed to extract task ID from response."
    echo "API Response: $RESPONSE"
    exit 1
fi
