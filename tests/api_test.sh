#!/bin/bash

# API Integration Test Script
# Tests all 33 API endpoints to verify they are implemented

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="${API_BASE_URL:-http://localhost:9090}"
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Test results
declare -a FAILED_ENDPOINTS

# Function to test endpoint existence
test_endpoint() {
    local method=$1
    local path=$2
    local description=$3
    local expected_not_404=$4
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -n "Testing: $description ... "
    
    # Make request and get status code
    status_code=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$API_BASE_URL$path" \
        -H "Content-Type: application/json" 2>/dev/null || echo "000")
    
    # Check if endpoint exists (not 404)
    if [ "$status_code" != "404" ]; then
        echo -e "${GREEN}PASS${NC} (HTTP $status_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}FAIL${NC} (HTTP $status_code - endpoint not found)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        FAILED_ENDPOINTS+=("$method $path")
    fi
}

echo "========================================="
echo "API Integration Test"
echo "========================================="
echo "Base URL: $API_BASE_URL"
echo ""

# Test health check endpoint
echo "--- Health Check ---"
test_endpoint "GET" "/health" "Health check endpoint"

echo ""
echo "--- Authentication Module (4 endpoints) ---"
test_endpoint "POST" "/api/v1/auth/login" "User login"
test_endpoint "POST" "/api/v1/auth/refresh" "Refresh token"
test_endpoint "POST" "/api/v1/auth/logout" "User logout"
test_endpoint "PUT" "/api/v1/auth/password" "Change password"

echo ""
echo "--- Food Management Module (6 endpoints) ---"
test_endpoint "POST" "/api/v1/foods" "Create food"
test_endpoint "GET" "/api/v1/foods" "List foods"
test_endpoint "GET" "/api/v1/foods/1" "Get food by ID"
test_endpoint "PUT" "/api/v1/foods/1" "Update food"
test_endpoint "DELETE" "/api/v1/foods/1" "Delete food"
test_endpoint "POST" "/api/v1/foods/batch" "Batch import foods"

echo ""
echo "--- Meal Records Module (5 endpoints) ---"
test_endpoint "POST" "/api/v1/meals" "Create meal"
test_endpoint "GET" "/api/v1/meals" "List meals"
test_endpoint "GET" "/api/v1/meals/1" "Get meal by ID"
test_endpoint "PUT" "/api/v1/meals/1" "Update meal"
test_endpoint "DELETE" "/api/v1/meals/1" "Delete meal"

echo ""
echo "--- Diet Plans Module (6 endpoints) ---"
test_endpoint "POST" "/api/v1/plans/generate" "Generate AI diet plan"
test_endpoint "GET" "/api/v1/plans" "List diet plans"
test_endpoint "GET" "/api/v1/plans/1" "Get plan by ID"
test_endpoint "PUT" "/api/v1/plans/1" "Update plan"
test_endpoint "DELETE" "/api/v1/plans/1" "Delete plan"
test_endpoint "POST" "/api/v1/plans/1/complete" "Complete plan"

echo ""
echo "--- AI Services Module (3 endpoints) ---"
test_endpoint "POST" "/api/v1/ai/chat" "AI chat"
test_endpoint "POST" "/api/v1/ai/suggest" "AI meal suggestion"
test_endpoint "GET" "/api/v1/ai/history" "Get chat history"

echo ""
echo "--- Nutrition Analysis Module (3 endpoints) ---"
test_endpoint "GET" "/api/v1/nutrition/daily/2024-01-15" "Get daily nutrition"
test_endpoint "GET" "/api/v1/nutrition/monthly" "Get monthly nutrition"
test_endpoint "GET" "/api/v1/nutrition/compare" "Compare nutrition"

echo ""
echo "--- Dashboard Module (1 endpoint) ---"
test_endpoint "GET" "/api/v1/dashboard" "Get dashboard data"

echo ""
echo "--- Settings Management Module (5 endpoints) ---"
test_endpoint "GET" "/api/v1/settings" "Get all settings"
test_endpoint "PUT" "/api/v1/settings/ai" "Update AI settings"
test_endpoint "GET" "/api/v1/settings/ai/test" "Test AI connection"
test_endpoint "GET" "/api/v1/user/profile" "Get user profile"
test_endpoint "PUT" "/api/v1/user/preferences" "Update user preferences"

echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -gt 0 ]; then
    echo ""
    echo "Failed Endpoints:"
    for endpoint in "${FAILED_ENDPOINTS[@]}"; do
        echo -e "  ${RED}✗${NC} $endpoint"
    done
    echo ""
    exit 1
else
    echo ""
    echo -e "${GREEN}All tests passed!${NC}"
    echo ""
    exit 0
fi
