#!/bin/bash

# Script to clean up temporary documentation files

echo "🧹 Cleaning up temporary documentation files..."

# Files to DELETE (temporary docs)
rm -f CLI_MIGRATION_FIXED.md
rm -f GORM_MIGRATION_COMPLETE.md
rm -f CONSTRAINT_ERROR_FIXED.md
rm -f LOGGER_FIXED.md
rm -f IMPORT_ALIAS_FIXED.md
rm -f CHECKLIST.md
rm -f DOCS_INDEX.md
rm -f MONOREPO_FIX.md
rm -f HTTP_WIRE_FIXED.md
rm -f RESTRUCTURE_GUIDE.md
rm -f ALL_FIXED_COMPLETE.md
rm -f FINAL_FIX_COMPLETE.md
rm -f SUMMARY.md
rm -f MIGRATION_FIX.md
rm -f PROVIDER_ERROR_FIXED.md
rm -f MAKEFILE_CLEANED.md
rm -f GORM_MIGRATION.md

echo "✅ Deleted temporary documentation files"

# Files to KEEP (important docs)
echo ""
echo "📚 Keeping important files:"
echo "  ✅ README.md - Project overview"
echo "  ✅ START_HERE.md - Quick start guide"
echo "  ✅ SETUP.md - Setup instructions"
echo "  ✅ QUICK_REFERENCE.md - Command reference"
echo "  ✅ DEVELOPMENT.md - Development guide"
echo ""
echo "🎉 Documentation cleanup complete!"

