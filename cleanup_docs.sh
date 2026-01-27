#!/bin/bash

# Script to clean up redundant documentation and script files

echo "Cleaning up redundant markdown files..."

# Remove old documentation files
rm -f DEVELOPMENT.md
rm -f QUICK_REFERENCE.md
rm -f SETUP.md
rm -f START_HERE.md

echo "✅ Removed redundant documentation files"
echo ""

echo "Cleaning up redundant script files (duplicates of Makefile)..."

# Remove redundant scripts (all functionality is in Makefile now)
rm -f scripts/check-env.sh
rm -f scripts/clean-docs.sh
rm -f scripts/generate.sh
rm -f scripts/quick-setup.sh
rm -f scripts/README.md

# Remove empty scripts directory if it exists
rmdir scripts 2>/dev/null || true

echo "✅ Removed redundant script files"
echo ""
echo "📚 Remaining documentation:"
echo "  📄 README.md           - Main project README"
echo "  📄 HUONG_DAN_SETUP.md  - Complete setup guide (Vietnamese)"
echo "  📄 AUTH_FLOW.md        - Authentication flow details (English)"
echo "  📄 SUMMARY.md          - Work summary"
echo ""
echo "🔧 All functionality moved to Makefile:"
echo "  make install-tools"
echo "  make wire"
echo "  make swagger"
echo "  make generate"
echo "  make migrate"
echo "  make seed"
echo "  make dev"
echo ""
echo "✅ Cleanup completed!"


