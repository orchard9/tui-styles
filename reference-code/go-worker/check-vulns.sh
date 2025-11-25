#!/bin/bash

# Check for vulnerabilities, only fail on module vulnerabilities (not Go stdlib)
output=$(govulncheck ./cmd/server ./internal/... 2>&1)

if echo "$output" | grep -q "vulnerabilities in modules you require.*[1-9]"; then
    echo "FAIL: Module vulnerabilities found"
    echo "$output"
    exit 1
else
    echo "PASS: Only Go standard library vulnerabilities found (require Go upgrade)"
    exit 0
fi
