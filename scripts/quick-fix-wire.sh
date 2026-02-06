#!/bin/bash

# Quick fix script - run this NOW
set -e

echo "⚡ QUICK FIX - Running now..."

# Step 0: Check Go installation
echo "0. Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    exit 1
fi

# Use system Go (avoid toolchain issues)
export GOTOOLCHAIN=local
export GO111MODULE=on

echo "   Using Go version: $(go version)"

# Clean cache (fixes x/tools v0.17.0 issue)
echo "1. Cleaning cache..."
go clean -cache -modcache || echo "   Warning: Some cache cleanup failed (might need sudo)"

# Remove lock file
echo "2. Removing go.sum..."
rm -f go.sum

# Update critical packages
echo "3. Updating packages..."
GOTOOLCHAIN=local go get github.com/google/wire@v0.6.0
GOTOOLCHAIN=local go get golang.org/x/tools@v0.29.0
GOTOOLCHAIN=local go get golang.org/x/sync@latest

# Tidy
echo "4. Tidying..."
GOTOOLCHAIN=local go mod tidy

# Install wire
echo "5. Installing wire..."
GOTOOLCHAIN=local go install github.com/google/wire/cmd/wire@v0.6.0

# Verify wire is in PATH
if ! command -v wire &> /dev/null; then
    echo "   Warning: wire not in PATH. Try: export PATH=\$PATH:\$(go env GOPATH)/bin"
    WIRE_CMD="$(go env GOPATH)/bin/wire"
else
    WIRE_CMD="wire"
fi

# Clean generated files
echo "6. Cleaning old wire files..."
rm -f cmd/http/wire_gen.go cmd/cli/migration/wire_gen.go

# Generate
echo "7. Generating wire code..."
cd cmd/http && $WIRE_CMD && echo "✅ HTTP OK" && cd ../..
cd cmd/cli/migration && $WIRE_CMD && echo "✅ Migration OK" && cd ../../..

echo ""
echo "✅ DONE! Try: make run"
