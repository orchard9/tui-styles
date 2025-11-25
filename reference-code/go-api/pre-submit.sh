#!/bin/bash
# Pre-submit validation script for Perforce (P4) workflows
# This script runs quality checks before submitting changes to the depot
#
# Usage:
#   ./scripts/pre-submit.sh           # Run all quality checks
#   make pre-submit                   # Run via Makefile
#   p4 submit                         # (after manual validation)
#
# Performance: ~30-45 seconds for typical changes

set -e  # Exit on first error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$SERVICE_DIR"

echo "=========================================="
echo "Pre-submit Validation for Creator API"
echo "=========================================="
echo ""
echo "Running quality checks before P4 submit..."
echo ""

# Track start time
START_TIME=$(date +%s)

# Function to report step status
report_step() {
    local step_name=$1
    local step_status=$2

    if [ "$step_status" -eq 0 ]; then
        echo "✓ $step_name passed"
    else
        echo "✗ $step_name failed"
        echo ""
        echo "Please fix the issues above before submitting to P4."
        echo "To bypass (emergency only): skip validation and run 'p4 submit' directly"
        exit 1
    fi
}

# Step 1: Check code formatting
echo "[1/7] Checking code formatting..."
make fmt-check > /dev/null 2>&1
report_step "Format check" $?
echo ""

# Step 2: Run go vet
echo "[2/7] Running go vet..."
make vet > /dev/null 2>&1
report_step "Go vet" $?
echo ""

# Step 3: Run golangci-lint
echo "[3/7] Running golangci-lint..."
make lint > /dev/null 2>&1
report_step "Linting" $?
echo ""

# Step 4: Run tests
echo "[4/7] Running tests..."
make test > /dev/null 2>&1
report_step "Tests" $?
echo ""

# Step 5: Check for dead code
echo "[5/7] Checking for dead code..."
make deadcode > /dev/null 2>&1
report_step "Dead code check" $?
echo ""

# Step 6: Check code complexity
echo "[6/7] Checking code complexity..."
make complexity > /dev/null 2>&1
report_step "Complexity check" $?
echo ""

# Step 7: Run security scans
echo "[7/7] Running security scans..."
make security > /dev/null 2>&1
report_step "Security check" $?
echo ""

# Calculate elapsed time
END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))

echo "=========================================="
echo "All quality checks passed! ✓"
echo "Elapsed time: ${ELAPSED}s"
echo "=========================================="
echo ""
echo "You can now submit to P4:"
echo "  p4 submit"
echo ""
echo "To view pending changes:"
echo "  p4 opened"
echo ""

exit 0
