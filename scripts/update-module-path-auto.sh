#!/bin/bash

# Automatic Module Path Update Script (Non-interactive)
# This script updates the Go module path without user prompts
# Usage: ./scripts/update-module-path-auto.sh [new-module-path]

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}=== Automatic Module Path Update ===${NC}"
echo ""

# Get new module path from argument or config file
if [ -n "$1" ]; then
    NEW_MODULE_PATH="$1"
    echo "Using module path from argument: $NEW_MODULE_PATH"
else
    # Read from config file
    MODULE_CONF="$PROJECT_ROOT/configs/module.conf"
    if [ ! -f "$MODULE_CONF" ]; then
        echo -e "${RED}Error: No module path provided and config file not found${NC}"
        echo "Usage: $0 <module-path>"
        exit 1
    fi
    
    NEW_MODULE_PATH=$(grep "^MODULE_PATH=" "$MODULE_CONF" | cut -d'=' -f2 | tr -d ' ')
    if [ -z "$NEW_MODULE_PATH" ]; then
        echo -e "${RED}Error: MODULE_PATH not found in config file${NC}"
        exit 1
    fi
    echo "Using module path from config: $NEW_MODULE_PATH"
fi

# Get old module path from go.mod
GO_MOD="$PROJECT_ROOT/go.mod"
if [ ! -f "$GO_MOD" ]; then
    echo -e "${RED}Error: go.mod not found${NC}"
    exit 1
fi

OLD_MODULE_PATH=$(head -n 1 "$GO_MOD" | sed 's/^module //')
if [ -z "$OLD_MODULE_PATH" ]; then
    echo -e "${RED}Error: Could not read module path from go.mod${NC}"
    exit 1
fi

echo -e "Old module path: ${YELLOW}$OLD_MODULE_PATH${NC}"
echo -e "New module path: ${GREEN}$NEW_MODULE_PATH${NC}"
echo ""

# Check if already up to date
if [ "$OLD_MODULE_PATH" = "$NEW_MODULE_PATH" ]; then
    echo -e "${GREEN}Module path is already up to date!${NC}"
    exit 0
fi

# Update go.mod
echo "Updating go.mod..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|^module $OLD_MODULE_PATH|module $NEW_MODULE_PATH|" "$GO_MOD"
else
    sed -i "s|^module $OLD_MODULE_PATH|module $NEW_MODULE_PATH|" "$GO_MOD"
fi
echo -e "${GREEN}✓${NC} go.mod updated"

# Update all .go files
echo "Updating .go files..."
GO_FILES=$(find "$PROJECT_ROOT" -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*")
FILE_COUNT=0

for file in $GO_FILES; do
    if grep -q "$OLD_MODULE_PATH" "$file"; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" "$file"
        else
            sed -i "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" "$file"
        fi
        FILE_COUNT=$((FILE_COUNT + 1))
    fi
done

echo -e "${GREEN}✓${NC} Updated $FILE_COUNT Go files"

# Run go mod tidy
echo "Running go mod tidy..."
cd "$PROJECT_ROOT"
if go mod tidy 2>&1; then
    echo -e "${GREEN}✓${NC} Dependencies updated"
else
    echo -e "${YELLOW}⚠${NC} go mod tidy had warnings (this is usually OK)"
fi

echo ""
echo -e "${GREEN}=== Update completed! ===${NC}"
echo ""
echo "Module path updated from:"
echo -e "  ${YELLOW}$OLD_MODULE_PATH${NC}"
echo "to:"
echo -e "  ${GREEN}$NEW_MODULE_PATH${NC}"
echo ""
