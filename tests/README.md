# Testing Documentation

This directory contains tests for the AI Diet Assistant backend API.

## Test Structure

### Unit Tests

Unit tests are located alongside the code they test:

- `internal/utils/date_test.go` - Tests for date parsing and formatting functions
- `internal/utils/response_test.go` - Tests for API response format functions

### Integration Tests

- `tests/integration_test.go` - Basic integration tests for response formats
- `tests/api_test.sh` - Shell script to test all 33 API endpoints against a running server

## Running Tests

### Run All Unit Tests

```bash
go test ./...
```

### Run Specific Package Tests

```bash
# Test utils package
go test -v ./internal/utils

# Test integration tests
go test -v ./tests
```

### Run API Integration Tests

The API integration test script tests all 33 endpoints against a running server:

```bash
# Start the server first
./bin/diet-assistant

# In another terminal, run the API tests
./tests/api_test.sh

# Or specify a custom API URL
API_BASE_URL=http://localhost:8080 ./tests/api_test.sh
```

## Test Coverage

### Unit Tests Coverage

The following functionality is covered by unit tests:

#### Date Handling (`internal/utils/date_test.go`)
- ✅ Parse date in YYYY-MM-DD format
- ✅ Parse date in ISO 8601 format
- ✅ Parse date to start of day (00:00:00)
- ✅ Parse date to end of day (23:59:59)
- ✅ Format date to YYYY-MM-DD
- ✅ Format date to ISO 8601
- ✅ Date parsing consistency across formats

#### Response Format (`internal/utils/response_test.go`)
- ✅ Success response format (code, message, data, timestamp)
- ✅ Success response with custom message
- ✅ Paginated response format (includes pagination object)
- ✅ Error response format (code, message, error, timestamp)
- ✅ Error response with underlying error
- ✅ HTTP status code mapping for all error codes
- ✅ Response timestamp generation
- ✅ Pagination calculation (page, page_size, total, total_pages)
- ✅ Sensitive data sanitization in error messages

### API Endpoints Coverage

The API integration test script (`tests/api_test.sh`) verifies all 33 endpoints exist:

#### Authentication Module (4 endpoints)
1. ✅ POST `/api/v1/auth/login` - User login
2. ✅ POST `/api/v1/auth/refresh` - Refresh token
3. ✅ POST `/api/v1/auth/logout` - User logout
4. ✅ PUT `/api/v1/auth/password` - Change password

#### Food Management Module (6 endpoints)
5. ✅ POST `/api/v1/foods` - Create food
6. ✅ GET `/api/v1/foods` - List foods
7. ✅ GET `/api/v1/foods/:id` - Get food by ID
8. ✅ PUT `/api/v1/foods/:id` - Update food
9. ✅ DELETE `/api/v1/foods/:id` - Delete food
10. ✅ POST `/api/v1/foods/batch` - Batch import foods

#### Meal Records Module (5 endpoints)
11. ✅ POST `/api/v1/meals` - Create meal
12. ✅ GET `/api/v1/meals` - List meals
13. ✅ GET `/api/v1/meals/:id` - Get meal by ID
14. ✅ PUT `/api/v1/meals/:id` - Update meal
15. ✅ DELETE `/api/v1/meals/:id` - Delete meal

#### Diet Plans Module (6 endpoints)
16. ✅ POST `/api/v1/plans/generate` - Generate AI diet plan
17. ✅ GET `/api/v1/plans` - List diet plans
18. ✅ GET `/api/v1/plans/:id` - Get plan by ID
19. ✅ PUT `/api/v1/plans/:id` - Update plan
20. ✅ DELETE `/api/v1/plans/:id` - Delete plan
21. ✅ POST `/api/v1/plans/:id/complete` - Complete plan

#### AI Services Module (3 endpoints)
22. ✅ POST `/api/v1/ai/chat` - AI chat
23. ✅ POST `/api/v1/ai/suggest` - AI meal suggestion
24. ✅ GET `/api/v1/ai/history` - Get chat history

#### Nutrition Analysis Module (3 endpoints)
25. ✅ GET `/api/v1/nutrition/daily/:date` - Get daily nutrition
26. ✅ GET `/api/v1/nutrition/monthly` - Get monthly nutrition
27. ✅ GET `/api/v1/nutrition/compare` - Compare nutrition

#### Dashboard Module (1 endpoint)
28. ✅ GET `/api/v1/dashboard` - Get dashboard data

#### Settings Management Module (5 endpoints)
29. ✅ GET `/api/v1/settings` - Get all settings
30. ✅ PUT `/api/v1/settings/ai` - Update AI settings
31. ✅ GET `/api/v1/settings/ai/test` - Test AI connection
32. ✅ GET `/api/v1/user/profile` - Get user profile
33. ✅ PUT `/api/v1/user/preferences` - Update user preferences

### Additional Tests

- ✅ Health check endpoint (`/health`)
- ✅ 404 handler for non-existent endpoints
- ✅ Authentication requirement for protected endpoints
- ✅ Parameter validation
- ✅ Pagination validation
- ✅ Error handling and error response format

## Test Requirements Verification

This test suite verifies all requirements from the backend-cleanup spec:

### Requirement 1: Frontend Code Cleanup
- ✅ No frontend routes exist (verified by API test script)
- ✅ Only API endpoints are accessible

### Requirement 2: API Completeness
- ✅ All 33 API endpoints are implemented and accessible

### Requirement 3: API Defect Fixes
- ✅ Response format compliance (tested in response_test.go)
- ✅ Error code standardization (tested in response_test.go)
- ✅ Parameter validation (tested in integration tests)
- ✅ Pagination handling (tested in response_test.go)
- ✅ Date handling (tested in date_test.go)

### Requirement 4: Code Quality
- ✅ Error handling tested
- ✅ Response format consistency tested
- ✅ Date parsing consistency tested

### Requirement 5: Router Cleanup
- ✅ Only /api/v1 routes exist (verified by API test script)
- ✅ Health check endpoint exists

### Requirement 6: Middleware Optimization
- ✅ Authentication middleware works (endpoints return 401 without auth)
- ✅ Error handling middleware works (404 responses are correct)

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run unit tests
        run: go test -v ./...
      
      - name: Build application
        run: go build -o bin/diet-assistant cmd/server/main.go
      
      - name: Start server
        run: ./bin/diet-assistant &
        env:
          DATABASE_HOST: localhost
          DATABASE_PASSWORD: postgres
      
      - name: Wait for server
        run: sleep 5
      
      - name: Run API integration tests
        run: ./tests/api_test.sh
```

## Manual Testing Checklist

For comprehensive manual testing, verify:

- [ ] All 33 API endpoints respond (not 404)
- [ ] Protected endpoints require authentication (return 401 without token)
- [ ] Parameter validation works (returns 400 with invalid params)
- [ ] Pagination works correctly (page, page_size, total, total_pages)
- [ ] Date formats are accepted (YYYY-MM-DD and ISO 8601)
- [ ] Error responses include all required fields
- [ ] Success responses include all required fields
- [ ] Health check endpoint returns 200 OK

## Troubleshooting

### Tests Fail to Connect to Database

If unit tests fail due to database connection issues:

```bash
# Skip tests that require database
go test -v ./internal/utils
go test -v ./tests
```

### API Test Script Fails

If the API test script fails:

1. Ensure the server is running: `./bin/diet-assistant`
2. Check the server is listening on the correct port (default: 9090)
3. Verify database and Redis are accessible
4. Check server logs for errors

### Server Won't Start

If the server won't start for testing:

1. Check database connection settings in `configs/config.yaml`
2. Ensure PostgreSQL is running
3. Run migrations: `./scripts/run-migrations.sh`
4. Check logs for specific errors

## Future Improvements

- Add end-to-end tests with real database transactions
- Add performance/load testing
- Add security testing (SQL injection, XSS, etc.)
- Add API contract testing with OpenAPI spec
- Add mutation testing to verify test quality
