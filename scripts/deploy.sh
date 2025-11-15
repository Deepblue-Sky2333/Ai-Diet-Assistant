#!/bin/bash

# AI Diet Assistant Deployment Script
# This script helps deploy the application to a production server

set -e

echo "==================================="
echo "AI Diet Assistant Deployment Script"
echo "==================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    echo "Please run the installation script first:"
    echo "  ./scripts/install.sh"
    echo ""
    echo "Or manually copy and configure:"
    echo "  cp .env.example .env"
    echo "  nano .env"
    exit 1
fi

# Load environment variables
source .env

echo -e "${GREEN}✓ Environment variables loaded${NC}"

# ============================================
# Configuration Validation
# ============================================
echo "Validating configuration..."

VALIDATION_FAILED=false

# Function to check if a value contains placeholder text
is_placeholder() {
    local value="$1"
    if [[ "$value" == *"REPLACE"* ]] || [[ "$value" == *"example"* ]] || \
       [[ "$value" == *"change"* ]] || [[ "$value" == *"your-"* ]] || \
       [[ "$value" == *"<"* ]] || [[ "$value" == *">"* ]]; then
        return 0
    fi
    return 1
}

# Validate JWT Secret
if [ -z "$JWT_SECRET" ]; then
    echo -e "${RED}✗ JWT_SECRET is not set${NC}"
    VALIDATION_FAILED=true
elif [ ${#JWT_SECRET} -lt 32 ]; then
    echo -e "${RED}✗ JWT_SECRET is too short (must be at least 32 characters)${NC}"
    VALIDATION_FAILED=true
elif is_placeholder "$JWT_SECRET"; then
    echo -e "${RED}✗ JWT_SECRET appears to be a placeholder or example value${NC}"
    VALIDATION_FAILED=true
else
    echo -e "${GREEN}✓ JWT_SECRET is valid${NC}"
fi

# Validate Encryption Key
if [ -z "$ENCRYPTION_KEY" ]; then
    echo -e "${RED}✗ ENCRYPTION_KEY is not set${NC}"
    VALIDATION_FAILED=true
elif [ ${#ENCRYPTION_KEY} -ne 32 ]; then
    echo -e "${RED}✗ ENCRYPTION_KEY must be exactly 32 bytes (current: ${#ENCRYPTION_KEY})${NC}"
    VALIDATION_FAILED=true
elif is_placeholder "$ENCRYPTION_KEY"; then
    echo -e "${RED}✗ ENCRYPTION_KEY appears to be a placeholder or example value${NC}"
    VALIDATION_FAILED=true
else
    echo -e "${GREEN}✓ ENCRYPTION_KEY is valid${NC}"
fi

# Validate Database Password
if [ -z "$DB_PASSWORD" ]; then
    echo -e "${RED}✗ DB_PASSWORD is not set${NC}"
    VALIDATION_FAILED=true
elif [ ${#DB_PASSWORD} -lt 12 ]; then
    echo -e "${RED}✗ DB_PASSWORD is too weak (must be at least 12 characters)${NC}"
    VALIDATION_FAILED=true
elif is_placeholder "$DB_PASSWORD"; then
    echo -e "${RED}✗ DB_PASSWORD appears to be a placeholder or example value${NC}"
    VALIDATION_FAILED=true
else
    echo -e "${GREEN}✓ DB_PASSWORD is set${NC}"
fi

# Validate CORS Configuration
if [ -z "$CORS_ALLOWED_ORIGINS" ]; then
    echo -e "${RED}✗ CORS_ALLOWED_ORIGINS is not set${NC}"
    VALIDATION_FAILED=true
elif [[ "$CORS_ALLOWED_ORIGINS" == *"*"* ]]; then
    echo -e "${RED}✗ CORS_ALLOWED_ORIGINS contains wildcard (*) - not allowed in production${NC}"
    VALIDATION_FAILED=true
elif is_placeholder "$CORS_ALLOWED_ORIGINS"; then
    echo -e "${RED}✗ CORS_ALLOWED_ORIGINS appears to be a placeholder or example value${NC}"
    VALIDATION_FAILED=true
else
    echo -e "${GREEN}✓ CORS_ALLOWED_ORIGINS is valid${NC}"
fi

# Check if config.yaml exists and validate it
if [ -f configs/config.yaml ]; then
    echo -e "${GREEN}✓ configs/config.yaml exists${NC}"
    
    # Check for placeholder values in config.yaml
    if grep -q "REPLACE\|<.*>" configs/config.yaml; then
        echo -e "${RED}✗ configs/config.yaml contains placeholder values${NC}"
        VALIDATION_FAILED=true
    else
        echo -e "${GREEN}✓ configs/config.yaml has no obvious placeholders${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Warning: configs/config.yaml not found${NC}"
fi

# Check module path in go.mod
if [ -f go.mod ]; then
    MODULE_PATH=$(head -n 1 go.mod | awk '{print $2}')
    if [[ "$MODULE_PATH" == *"yourusername"* ]] || [[ "$MODULE_PATH" == *"example"* ]]; then
        echo -e "${RED}✗ Module path in go.mod contains placeholder: $MODULE_PATH${NC}"
        echo "  Please update it with: scripts/update-module-path.sh"
        VALIDATION_FAILED=true
    else
        echo -e "${GREEN}✓ Module path is properly configured: $MODULE_PATH${NC}"
    fi
fi

# If validation failed, stop deployment
if [ "$VALIDATION_FAILED" = true ]; then
    echo ""
    echo -e "${RED}==================================="
    echo "Configuration validation FAILED!"
    echo "===================================${NC}"
    echo ""
    echo "Please fix the errors above before deploying."
    echo "Run './scripts/install.sh' to generate secure configuration."
    echo ""
    exit 1
fi

echo -e "${GREEN}✓ All configuration checks passed${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed!${NC}"
    echo "Please install Go 1.21 or higher"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"

# Check if MySQL is accessible
echo "Checking database connection..."
if command -v mysql &> /dev/null; then
    if mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -e "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✓ Database connection successful${NC}"
    else
        echo -e "${YELLOW}⚠ Warning: Cannot connect to database${NC}"
        echo "Please ensure MySQL is running and credentials are correct"
    fi
else
    echo -e "${YELLOW}⚠ Warning: mysql client not found, skipping database check${NC}"
fi

# Build the application
echo "Building application..."
make build || go build -o bin/diet-assistant cmd/server/main.go

if [ -f bin/diet-assistant ]; then
    echo -e "${GREEN}✓ Application built successfully${NC}"
else
    echo -e "${RED}Error: Build failed!${NC}"
    exit 1
fi

# Run database migrations
echo "Running database migrations..."
if [ -f bin/diet-assistant ]; then
    # Check if migrate command exists in the binary
    ./bin/diet-assistant migrate up 2>/dev/null || echo -e "${YELLOW}⚠ Migration command not available or already up to date${NC}"
fi

# Create necessary directories
echo "Creating necessary directories..."
mkdir -p logs
mkdir -p uploads
echo -e "${GREEN}✓ Directories created${NC}"

# Set permissions
echo "Setting permissions..."
chmod +x bin/diet-assistant
chmod +x scripts/*.sh
echo -e "${GREEN}✓ Permissions set${NC}"

# Check if systemd service should be installed
read -p "Do you want to install systemd service? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -f scripts/diet-assistant.service ]; then
        sudo cp scripts/diet-assistant.service /etc/systemd/system/
        sudo systemctl daemon-reload
        sudo systemctl enable diet-assistant
        echo -e "${GREEN}✓ Systemd service installed${NC}"
        echo "You can now start the service with: sudo systemctl start diet-assistant"
    else
        echo -e "${YELLOW}⚠ Service file not found at scripts/diet-assistant.service${NC}"
    fi
fi

# Check if nginx configuration should be installed
read -p "Do you want to install nginx configuration? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -f scripts/nginx.conf ]; then
        sudo cp scripts/nginx.conf /etc/nginx/sites-available/diet-assistant
        sudo ln -sf /etc/nginx/sites-available/diet-assistant /etc/nginx/sites-enabled/
        sudo nginx -t && sudo systemctl reload nginx
        echo -e "${GREEN}✓ Nginx configuration installed${NC}"
    else
        echo -e "${YELLOW}⚠ Nginx config file not found at scripts/nginx.conf${NC}"
    fi
fi

echo ""
echo -e "${GREEN}==================================="
echo "Deployment completed successfully!"
echo "===================================${NC}"
echo ""
echo "Next steps:"
echo "1. Review your .env configuration"
echo "2. Start the application:"
echo "   - Manually: ./scripts/start.sh"
echo "   - With systemd: sudo systemctl start diet-assistant"
echo "3. Access the application at http://localhost:${SERVER_PORT}"
echo ""
