#!/bin/bash

set -e

echo "🔧 Fixing Wire and Go tools dependency issues..."

# Step 0: Check Go installation
echo "🔍 Step 0: Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    exit 1
fi

# Use system Go (avoid toolchain download issues)
export GOTOOLCHAIN=local
export GO111MODULE=on

echo "   ✅ Using Go version: $(go version)"

# Step 1: Clean Go cache (fixes x/tools version mismatch)
echo "🧹 Step 1: Cleaning Go cache..."
go clean -cache -modcache || echo "   ⚠️  Warning: Some cache cleanup failed"

# Step 2: Remove go.sum to force fresh download
echo "🗑️  Step 2: Removing go.sum..."
rm -f go.sum

# Step 3: Update dependencies
echo "📦 Step 3: Updating dependencies..."
GOTOOLCHAIN=local go get github.com/google/wire@v0.6.0
GOTOOLCHAIN=local go get golang.org/x/tools@v0.29.0
GOTOOLCHAIN=local go get golang.org/x/sync@latest

# Step 4: Tidy and download
echo "📥 Step 4: Tidying and downloading all dependencies..."
GOTOOLCHAIN=local go mod tidy
GOTOOLCHAIN=local go mod download

# Step 5: Install Wire CLI
echo "🔨 Step 5: Installing Wire CLI v0.6.0..."
GOTOOLCHAIN=local go install github.com/google/wire/cmd/wire@v0.6.0

# Step 6: Verify Wire version
echo "✅ Step 6: Verifying Wire installation..."
if command -v wire &> /dev/null; then
    wire version || echo "   Wire installed successfully"
    WIRE_CMD="wire"
else
    echo "   ⚠️  Wire not in PATH, using full path"
    WIRE_CMD="$(go env GOPATH)/bin/wire"
fi

# Step 7: Clean old generated files
echo "🧹 Step 7: Cleaning old generated files..."
rm -f cmd/http/wire_gen.go
rm -f cmd/cli/migration/wire_gen.go

# Step 8: Generate wire code
echo "⚙️  Step 8: Generating wire code for HTTP server..."
cd cmd/http && $WIRE_CMD
echo "✅ HTTP wire generated successfully!"
cd ../..

echo "⚙️  Step 9: Generating wire code for migration..."
cd cmd/cli/migration && $WIRE_CMD
echo "✅ Migration wire generated successfully!"
cd ../../..

echo ""
echo "✨ Done! All issues fixed successfully!"
echo ""
echo "You can now run:"
echo "  make run      # Run HTTP server"
echo "  make migrate  # Run migration"
