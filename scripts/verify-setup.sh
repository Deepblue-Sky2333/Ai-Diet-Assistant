#!/bin/bash

# AI Diet Assistant Setup Verification Script
# This script verifies that the installation is correct

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================="
echo "AI Diet Assistant Setup Verification"
echo "===================================${NC}"
echo ""

ERRORS=0
WARNINGS=0

# Function to check file exists
check_file() {
    local file=$1
    local description=$2
    
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $description: $file"
        return 0
    else
        echo -e "${RED}✗${NC} $description: $file (NOT FOUND)"
        ((ERRORS++))
        return 1
    fi
}

# Function to check directory exists
check_dir() {
    local dir=$1
    local description=$2
    
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✓${NC} $description: $dir"
        return 0
    else
        echo -e "${RED}✗${NC} $description: $dir (NOT FOUND)"
        ((ERRORS++))
        return 1
    fi
}

# Function to check command exists
check_command() {
    local cmd=$1
    local description=$2
    
    if command -v $cmd &> /dev/null; then
        local version=$($cmd --version 2>&1 | head -n 1)
        echo -e "${GREEN}✓${NC} $description: $cmd ($version)"
        return 0
    else
        echo -e "${RED}✗${NC} $description: $cmd (NOT INSTALLED)"
        ((ERRORS++))
        return 1
    fi
}

# Check required commands
echo -e "${BLUE}Checking required commands...${NC}"
check_command "go" "Go"
check_command "mysql" "MySQL"
check_command "node" "Node.js"
check_command "npm" "npm"
echo ""

# Check backend files
echo -e "${BLUE}Checking backend files...${NC}"
check_file ".env" "Backend environment file"
check_file "configs/config.yaml" "Backend config file"
check_file "go.mod" "Go module file"
check_dir "internal" "Internal directory"
check_dir "cmd" "Command directory"
echo ""

# Check frontend files
echo -e "${BLUE}Checking frontend files...${NC}"
check_dir "web/frontend" "Frontend directory"
check_file "web/frontend/package.json" "Frontend package.json"
check_file "web/frontend/.env.local.example" "Frontend env example"

if [ -f "web/frontend/.env.local" ]; then
    echo -e "${GREEN}✓${NC} Frontend environment file: web/frontend/.env.local"
else
    echo -e "${YELLOW}⚠${NC} Frontend environment file: web/frontend/.env.local (NOT FOUND)"
    echo "  Run ./scripts/install.sh to create it"
    ((WARNINGS++))
fi

if [ -d "web/frontend/node_modules" ]; then
    echo -e "${GREEN}✓${NC} Frontend dependencies installed"
else
    echo -e "${YELLOW}⚠${NC} Frontend dependencies not installed"
    echo "  Run: cd web/frontend && npm install"
    ((WARNINGS++))
fi
echo ""

# Check scripts
echo -e "${BLUE}Checking scripts...${NC}"
check_file "scripts/install.sh" "Installation script"
check_file "scripts/start.sh" "Backend start script"
check_file "scripts/start-frontend.sh" "Frontend start script"
check_file "scripts/start-all.sh" "Unified start script"
check_file "scripts/stop.sh" "Stop script"
echo ""

# Check backend configuration
echo -e "${BLUE}Checking backend configuration...${NC}"
if [ -f ".env" ]; then
    # Check JWT secret
    JWT_SECRET=$(grep "^JWT_SECRET=" .env | cut -d'=' -f2)
    if [ -n "$JWT_SECRET" ] && [ ${#JWT_SECRET} -ge 32 ]; then
        echo -e "${GREEN}✓${NC} JWT secret is configured (${#JWT_SECRET} chars)"
    else
        echo -e "${RED}✗${NC} JWT secret is too short or not configured"
        ((ERRORS++))
    fi
    
    # Check encryption key
    ENCRYPTION_KEY=$(grep "^ENCRYPTION_KEY=" .env | cut -d'=' -f2)
    if [ -n "$ENCRYPTION_KEY" ] && [ ${#ENCRYPTION_KEY} -eq 32 ]; then
        echo -e "${GREEN}✓${NC} Encryption key is configured (32 bytes)"
    else
        echo -e "${RED}✗${NC} Encryption key is not 32 bytes or not configured"
        ((ERRORS++))
    fi
    
    # Check database configuration
    DB_HOST=$(grep "^DB_HOST=" .env | cut -d'=' -f2)
    DB_NAME=$(grep "^DB_NAME=" .env | cut -d'=' -f2)
    if [ -n "$DB_HOST" ] && [ -n "$DB_NAME" ]; then
        echo -e "${GREEN}✓${NC} Database configuration found"
    else
        echo -e "${RED}✗${NC} Database configuration incomplete"
        ((ERRORS++))
    fi
    
    # Check CORS configuration
    CORS_ORIGINS=$(grep "^CORS_ALLOWED_ORIGINS=" .env | cut -d'=' -f2)
    if [ -n "$CORS_ORIGINS" ]; then
        if [[ "$CORS_ORIGINS" == *"*"* ]]; then
            echo -e "${RED}✗${NC} CORS configuration uses wildcard (*)"
            ((ERRORS++))
        else
            echo -e "${GREEN}✓${NC} CORS configuration found"
        fi
    else
        echo -e "${YELLOW}⚠${NC} CORS configuration not found"
        ((WARNINGS++))
    fi
fi
echo ""

# Check frontend configuration
echo -e "${BLUE}Checking frontend configuration...${NC}"
if [ -f "web/frontend/.env.local" ]; then
    # Check demo mode
    DEMO_MODE=$(grep "^NEXT_PUBLIC_DEMO_MODE=" web/frontend/.env.local | cut -d'=' -f2)
    if [ -n "$DEMO_MODE" ]; then
        echo -e "${GREEN}✓${NC} Demo mode configured: $DEMO_MODE"
    else
        echo -e "${YELLOW}⚠${NC} Demo mode not configured"
        ((WARNINGS++))
    fi
    
    # Check API URL
    API_URL=$(grep "^NEXT_PUBLIC_API_URL=" web/frontend/.env.local | cut -d'=' -f2)
    if [ -n "$API_URL" ]; then
        echo -e "${GREEN}✓${NC} API URL configured: $API_URL"
    else
        echo -e "${RED}✗${NC} API URL not configured"
        ((ERRORS++))
    fi
fi
echo ""

# Check database connection
echo -e "${BLUE}Checking database connection...${NC}"
if [ -f ".env" ]; then
    DB_HOST=$(grep "^DB_HOST=" .env | cut -d'=' -f2)
    DB_PORT=$(grep "^DB_PORT=" .env | cut -d'=' -f2)
    DB_USER=$(grep "^DB_USER=" .env | cut -d'=' -f2)
    DB_PASSWORD=$(grep "^DB_PASSWORD=" .env | cut -d'=' -f2)
    DB_NAME=$(grep "^DB_NAME=" .env | cut -d'=' -f2)
    
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Database connection successful"
    else
        echo -e "${RED}✗${NC} Database connection failed"
        echo "  Please check your database configuration and ensure MySQL is running"
        ((ERRORS++))
    fi
fi
echo ""

# Check if backend is running
echo -e "${BLUE}Checking running services...${NC}"
if lsof -Pi :9090 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Backend is running on port 9090"
else
    echo -e "${YELLOW}⚠${NC} Backend is not running"
    echo "  Start with: ./scripts/start.sh"
fi

if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Frontend is running on port 3000"
else
    echo -e "${YELLOW}⚠${NC} Frontend is not running"
    echo "  Start with: ./scripts/start-frontend.sh"
fi
echo ""

# Summary
echo -e "${BLUE}==================================="
echo "Verification Summary"
echo "===================================${NC}"
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed!${NC}"
    echo ""
    echo "Your setup is complete. You can now:"
    echo "  1. Start all services: ./scripts/start-all.sh"
    echo "  2. Or start separately:"
    echo "     - Backend:  ./scripts/start.sh"
    echo "     - Frontend: ./scripts/start-frontend.sh"
    echo ""
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠ Setup complete with $WARNINGS warning(s)${NC}"
    echo ""
    echo "Please review the warnings above."
    echo "You can still start the application, but some features may not work correctly."
    echo ""
    exit 0
else
    echo -e "${RED}✗ Setup incomplete: $ERRORS error(s), $WARNINGS warning(s)${NC}"
    echo ""
    echo "Please fix the errors above before starting the application."
    echo ""
    echo "Common fixes:"
    echo "  - Run: ./scripts/install.sh"
    echo "  - Install dependencies: cd web/frontend && npm install"
    echo "  - Check database configuration"
    echo ""
    exit 1
fi
