#!/bin/bash

# Script setup nhanh cho lần đầu tiên

set -e

echo "🚀 Quick Setup for Doan Project"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Install tools
echo -e "${BLUE}Step 1: Installing tools...${NC}"
if ! command -v wire &> /dev/null || ! command -v swag &> /dev/null; then
    echo "Installing Wire and Swag..."
    go install github.com/google/wire/cmd/wire@latest
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✓ Tools installed${NC}"
else
    echo -e "${GREEN}✓ Tools already installed${NC}"
fi
echo ""

# Step 2: Download dependencies
echo -e "${BLUE}Step 2: Downloading dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies ready${NC}"
echo ""

# Step 3: Copy config
echo -e "${BLUE}Step 3: Setting up configuration...${NC}"
if [ ! -f "configs/config.yaml" ]; then
    cp configs/config.yaml.sample configs/config.yaml
    echo -e "${GREEN}✓ Config file created${NC}"
    echo -e "${YELLOW}⚠️  Please edit configs/config.yaml with your settings${NC}"
else
    echo -e "${GREEN}✓ Config file already exists${NC}"
fi
echo ""

# Step 4: Start Docker services (optional)
echo -e "${BLUE}Step 4: Starting Docker services (optional)...${NC}"
if command -v docker &> /dev/null; then
    read -p "Start Docker services (PostgreSQL, Redis)? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose -f tools/docker-compose.local.yml up -d
        echo -e "${GREEN}✓ Docker services started${NC}"
        echo "Waiting for services to be ready..."
        sleep 5
    else
        echo "Skipped Docker services"
    fi
else
    echo -e "${YELLOW}⚠️  Docker not found - using existing database${NC}"
fi
echo ""

# Step 5: Generate code
echo -e "${BLUE}Step 5: Generating code...${NC}"
cd cmd/http && wire && cd ../..
cd cmd/cli/migration && wire && cd ../../..
cd cmd/cli/get-access-token && wire && cd ../../..
swag init -g cmd/http/main.go -o cmd/http/docs
echo -e "${GREEN}✓ Code generated${NC}"
echo ""

# Done
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}✅ Setup complete!${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo "Next steps:"
echo "  1. Edit configs/config.yaml if needed"
echo "  2. Run: make dev"
echo ""
echo "🎉 Happy coding!"

