#!/bin/bash

# Batch runner script for go-image-generator
# Usage: ./run_batch.sh [max_id]
# Example: ./run_batch.sh 50

# Default maximum ID if not provided
DEFAULT_MIN_ID=1
DEFAULT_MAX_ID=1

# Get the ID from command line argument or use default
MIN_ID=${1:-$DEFAULT_MIN_ID}
MAX_ID=${2:-$DEFAULT_MAX_ID}

# Validate that MAX_ID is a positive integer
if ! [[ "$MAX_ID" =~ ^[1-9][0-9]*$ ]]; then
    echo "Error: MAX_ID must be a positive integer"
    echo "Usage: $0 [min_id] [max_id]"
    echo "Example: $0 10 50"
    exit 1
fi

echo "Running go-image-generator for IDs $MIN_ID to $MAX_ID"
echo "Template: assets/templates/template.json"
echo "----------------------------------------"

# Counter for successful runs
SUCCESS_COUNT=0
# Counter for failed runs
FAIL_COUNT=0

# Loop from MIN_ID to MAX_ID
for id in $(seq $MIN_ID $MAX_ID); do
    echo "Processing ID: $id"
    
    # Run the go command
    if go run cmd/main.go --template assets/templates/template.json --id $id; then
        echo "✓ Successfully processed ID: $id"
        ((SUCCESS_COUNT++))
    else
        echo "✗ Failed to process ID: $id"
        ((FAIL_COUNT++))
    fi
    
    echo "----------------------------------------"
done

# Print summary
echo "Batch processing completed!"
echo "Successful: $SUCCESS_COUNT"
echo "Failed: $FAIL_COUNT"
echo "Total: $((MAX_ID - MIN_ID + 1))"

# Exit with error code if any runs failed
if [ $FAIL_COUNT -gt 0 ]; then
    exit 1
fi
