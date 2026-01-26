#!/bin/bash

# Script tự động generate tất cả code cần thiết

set -e

echo "🚀 Starting code generation..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if wire is installed
if ! command -v wire &> /dev/null; then
    echo -e "${YELLOW}⚠️  Wire not found. Installing...${NC}"
    go install github.com/google/wire/cmd/wire@latest
    echo -e "${GREEN}✅ Wire installed successfully${NC}"
fi

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo -e "${YELLOW}⚠️  Swag not found. Installing...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✅ Swag installed successfully${NC}"
fi

# Generate Wire code for HTTP server
echo -e "${BLUE}📦 Generating Wire code for HTTP server...${NC}"
cd cmd/http && wire && cd ../..
echo -e "${GREEN}✅ HTTP server wire generated${NC}"

# Generate Wire code for migration CLI
echo -e "${BLUE}📦 Generating Wire code for migration CLI...${NC}"
cd cmd/cli/migration && wire && cd ../../..
echo -e "${GREEN}✅ Migration wire generated${NC}"

# Generate Wire code for get-access-token CLI
echo -e "${BLUE}📦 Generating Wire code for CLI tools...${NC}"
cd cmd/cli/get-access-token && wire && cd ../../..
echo -e "${GREEN}✅ CLI tools wire generated${NC}"

# Generate Swagger documentation
echo -e "${BLUE}📚 Generating Swagger documentation...${NC}"
swag init -g cmd/http/main.go -o cmd/http/docs
echo -e "${GREEN}✅ Swagger docs generated${NC}"

echo -e "${GREEN}🎉 All code generated successfully!${NC}"
echo -e "${BLUE}💡 Run 'make run' to start the server${NC}"

