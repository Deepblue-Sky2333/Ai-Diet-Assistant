#!/bin/bash

# Module Path Update Script
# This script updates the Go module path throughout the entire project.
# It reads the module path from configs/module.conf and updates:
# - go.mod module declaration
# - All import statements in .go files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the project root directory (parent of scripts directory)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}=== Module Path Update Script ===${NC}"
echo ""

# Check if module.conf exists
MODULE_CONF="$PROJECT_ROOT/configs/module.conf"
if [ ! -f "$MODULE_CONF" ]; then
    echo -e "${RED}Error: Configuration file not found: $MODULE_CONF${NC}"
    echo "Please create the file with MODULE_PATH=your/module/path"
    exit 1
fi

# Read the module path from config file
echo "Reading module path from $MODULE_CONF..."
MODULE_PATH=$(grep "^MODULE_PATH=" "$MODULE_CONF" | cut -d'=' -f2 | tr -d ' ')

if [ -z "$MODULE_PATH" ]; then
    echo -e "${RED}Error: MODULE_PATH not found in $MODULE_CONF${NC}"
    echo "Please add: MODULE_PATH=your/module/path"
    exit 1
fi

echo -e "New module path: ${GREEN}$MODULE_PATH${NC}"
echo ""

# Validate module path format
if [[ ! "$MODULE_PATH" =~ ^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$ ]]; then
    echo -e "${YELLOW}Warning: Module path format may be invalid: $MODULE_PATH${NC}"
    echo "Expected format: domain.com/username/repository"
    
    if [ -t 0 ]; then
        # Interactive mode
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted."
            exit 1
        fi
    else
        # Non-interactive mode
        echo "Running in non-interactive mode, continuing..."
        read -t 1 REPLY || REPLY="y"
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted."
            exit 1
        fi
    fi
fi

# Check for placeholder values
if [[ "$MODULE_PATH" == *"yourusername"* ]] || [[ "$MODULE_PATH" == *"example"* ]]; then
    echo -e "${RED}Error: Module path contains placeholder values: $MODULE_PATH${NC}"
    echo "Please update configs/module.conf with your actual repository path."
    exit 1
fi

# Get the old module path from go.mod
GO_MOD="$PROJECT_ROOT/go.mod"
if [ ! -f "$GO_MOD" ]; then
    echo -e "${RED}Error: go.mod not found at $GO_MOD${NC}"
    exit 1
fi

OLD_MODULE_PATH=$(head -n 1 "$GO_MOD" | sed 's/^module //')

if [ -z "$OLD_MODULE_PATH" ]; then
    echo -e "${RED}Error: Could not read module path from go.mod${NC}"
    exit 1
fi

echo -e "Old module path: ${YELLOW}$OLD_MODULE_PATH${NC}"
echo ""

# Check if paths are the same
if [ "$OLD_MODULE_PATH" = "$MODULE_PATH" ]; then
    echo -e "${GREEN}Module path is already up to date!${NC}"
    exit 0
fi

# Confirm with user (unless running in non-interactive mode)
if [ -t 0 ]; then
    # Interactive mode
    echo -e "${YELLOW}This will update all Go files in the project.${NC}"
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
else
    # Non-interactive mode (piped input)
    echo -e "${YELLOW}Running in non-interactive mode, proceeding with update...${NC}"
    read -t 1 REPLY || REPLY="y"
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

echo ""
echo "Starting update process..."
echo ""

# Step 1: Update go.mod
echo "1. Updating go.mod..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS 需要提供备份扩展名
    sed -i '.bak' "s|^module $OLD_MODULE_PATH|module $MODULE_PATH|" "$GO_MOD"
else
    # Linux
    sed -i.bak "s|^module $OLD_MODULE_PATH|module $MODULE_PATH|" "$GO_MOD"
fi
rm -f "$GO_MOD.bak"
echo -e "   ${GREEN}✓${NC} go.mod updated"

# Step 2: Update all .go files
echo "2. Updating import statements in .go files..."
GO_FILES=$(find "$PROJECT_ROOT" -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*")
FILE_COUNT=0

for file in $GO_FILES; do
    if grep -q "$OLD_MODULE_PATH" "$file"; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '.bak' "s|$OLD_MODULE_PATH|$MODULE_PATH|g" "$file"
        else
            # Linux
            sed -i.bak "s|$OLD_MODULE_PATH|$MODULE_PATH|g" "$file"
        fi
        rm -f "$file.bak"
        FILE_COUNT=$((FILE_COUNT + 1))
    fi
done

echo -e "   ${GREEN}✓${NC} Updated $FILE_COUNT Go files"

# Step 3: Run go mod tidy
echo "3. Running go mod tidy..."
cd "$PROJECT_ROOT"
if go mod tidy; then
    echo -e "   ${GREEN}✓${NC} Dependencies updated"
else
    echo -e "   ${RED}✗${NC} go mod tidy failed"
    exit 1
fi

# Step 4: Verify by compiling
echo "4. Verifying update (compiling project)..."
if go build -o /dev/null ./cmd/server; then
    echo -e "   ${GREEN}✓${NC} Compilation successful"
else
    echo -e "   ${RED}✗${NC} Compilation failed"
    echo ""
    echo -e "${RED}Error: The project does not compile after the update.${NC}"
    echo "Please check the errors above and fix them manually."
    exit 1
fi

echo ""
echo -e "${GREEN}=== Update completed successfully! ===${NC}"
echo ""
echo "Module path updated from:"
echo -e "  ${YELLOW}$OLD_MODULE_PATH${NC}"
echo "to:"
echo -e "  ${GREEN}$MODULE_PATH${NC}"
echo ""
echo "All Go files have been updated and the project compiles successfully."
