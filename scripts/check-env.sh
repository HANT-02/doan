#!/bin/bash

# Script kiểm tra môi trường trước khi chạy

set -e

echo "🔍 Checking environment..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check Go
echo -n "Checking Go... "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go is not installed"
    echo "Install Go from: https://golang.org/dl/"
    exit 1
fi

# Check Wire
echo -n "Checking Wire... "
if command -v wire &> /dev/null; then
    echo -e "${GREEN}✓${NC} Installed"
else
    echo -e "${YELLOW}⚠${NC}  Wire not found"
    echo "Install with: make install-tools"
fi

# Check Swag
echo -n "Checking Swag... "
if command -v swag &> /dev/null; then
    echo -e "${GREEN}✓${NC} Installed"
else
    echo -e "${YELLOW}⚠${NC}  Swag not found"
    echo "Install with: make install-tools"
fi

# Check config file
echo -n "Checking config file... "
if [ -f "configs/config.yaml" ]; then
    echo -e "${GREEN}✓${NC} Found"
else
    echo -e "${YELLOW}⚠${NC}  Not found"
    echo "Copy from: cp configs/config.yaml.sample configs/config.yaml"
fi

# Check PostgreSQL connection (if config exists)
if [ -f "configs/config.yaml" ]; then
    echo -n "Checking PostgreSQL... "
    if command -v psql &> /dev/null; then
        # Extract DB info from config (simple approach)
        echo -e "${GREEN}✓${NC} psql found"
    else
        echo -e "${YELLOW}⚠${NC}  psql not found (using Docker?)"
    fi
fi

# Check Docker (optional)
echo -n "Checking Docker... "
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓${NC} Installed"
else
    echo -e "${YELLOW}⚠${NC}  Docker not found (optional)"
fi

# Check if wire_gen.go exists
echo -n "Checking generated files... "
if [ -f "cmd/http/wire_gen.go" ]; then
    echo -e "${GREEN}✓${NC} Found"
else
    echo -e "${YELLOW}⚠${NC}  Not found - run 'make generate'"
fi

echo ""
echo "🎯 Ready to start!"
echo "Run: make dev"

