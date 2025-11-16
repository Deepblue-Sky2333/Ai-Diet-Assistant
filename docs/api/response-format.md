# API Response Format

This document describes the standardized response format used by all API endpoints in the AI Diet Assistant backend.

## Success Response Format

All successful API responses follow this structure:

```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "timestamp": 1700000000
}
```

### Fields

- `code` (integer): Always `0` for successful responses
- `message` (string): Human-readable success message (default: "success")
- `data` (object/array/null): The actual response data
- `timestamp` (integer): Unix timestamp when the response was generated

### Example

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Apple",
    "category": "fruit",
    "calories": 52
  },
  "timestamp": 1700000000
}
```

## Paginated Response Format

Endpoints that return lists with pagination use this structure:

```json
{
  "code": 0,
  "message": "success",
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1700000000
}
```

### Pagination Fields

- `page` (integer): Current page number (1-indexed)
- `page_size` (integer): Number of items per page
- `total` (integer): Total number of items
- `total_pages` (integer): Total number of pages

### Pagination Defaults

- Default `page`: 1
- Default `page_size`: 20
- Maximum `page_size`: 100

### Example

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Apple",
      "category": "fruit"
    },
    {
      "id": 2,
      "name": "Banana",
      "category": "fruit"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1700000000
}
```

## Error Response Format

All error responses follow this structure:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "detailed error message",
  "timestamp": 1700000000
}
```

### Fields

- `code` (integer): Error code (see Error Codes section)
- `message` (string): Human-readable error message
- `error` (string, optional): Detailed error information (only in development mode)
- `timestamp` (integer): Unix timestamp when the error occurred

### Example

```json
{
  "code": 40001,
  "message": "invalid request parameters",
  "error": "field 'name' is required",
  "timestamp": 1700000000
}
```

## Error Codes

The API uses standardized error codes to indicate different types of errors:

### Client Errors (4xxxx)

| Code  | HTTP Status | Description                    |
|-------|-------------|--------------------------------|
| 40001 | 400         | Invalid parameters             |
| 40002 | 400         | Validation error               |
| 40101 | 401         | Unauthorized                   |
| 40301 | 403         | Forbidden                      |
| 40401 | 404         | Resource not found             |
| 40901 | 409         | Resource conflict              |
| 42901 | 429         | Too many requests (rate limit) |

### Server Errors (5xxxx)

| Code  | HTTP Status | Description           |
|-------|-------------|-----------------------|
| 50001 | 500         | Internal server error |
| 50002 | 500         | Database error        |
| 50003 | 500         | AI service error      |
| 50004 | 500         | Encryption error      |

## HTTP Status Codes

The API uses standard HTTP status codes in addition to the custom error codes:

- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters or validation error
- `401 Unauthorized`: Authentication required or failed
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server-side error

## Security Considerations

### Error Messages in Production

In production mode (when `GIN_MODE=release`):
- Detailed error messages are NOT included in the response
- Only generic error messages are returned to prevent information leakage
- Detailed errors are logged server-side for debugging

### Sensitive Data Sanitization

The API automatically sanitizes error messages to remove:
- Database connection strings and passwords
- File paths
- SQL queries
- Other sensitive information

## Usage Examples

### Successful Request

```bash
curl -X GET "http://localhost:9090/api/v1/foods/1" \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Apple",
    "category": "fruit",
    "price": 2.5,
    "unit": "piece",
    "protein": 0.3,
    "carbs": 14.0,
    "fat": 0.2,
    "fiber": 2.4,
    "calories": 52.0,
    "available": true
  },
  "timestamp": 1700000000
}
```

### Paginated Request

```bash
curl -X GET "http://localhost:9090/api/v1/foods?page=1&page_size=20" \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Apple",
      "category": "fruit"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1700000000
}
```

### Error Request

```bash
curl -X GET "http://localhost:9090/api/v1/foods/999" \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "code": 40401,
  "message": "food not found",
  "timestamp": 1700000000
}
```

## Implementation Notes

### For Developers

All handlers should use the utility functions from `internal/utils/response.go`:

- `Success(c, data)` - Standard success response
- `SuccessWithMessage(c, message, data)` - Success with custom message
- `SuccessPaginated(c, data, pagination)` - Paginated response
- `Error(c, appErr)` - Error response
- `ErrorWithMessage(c, code, message)` - Error with custom message

### Example Handler Code

```go
func (h *FoodHandler) GetFood(c *gin.Context) {
    foodID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid food id", err))
        return
    }

    food, err := h.foodService.GetFood(userID, foodID)
    if err != nil {
        utils.Error(c, utils.NewAppError(utils.CodeNotFound, "food not found", err))
        return
    }

    utils.Success(c, food)
}
```

## Validation

All API responses are validated to ensure they conform to this format. The test suite includes comprehensive tests for:

- Success response structure
- Paginated response structure
- Error response structure
- HTTP status code mapping
- Timestamp generation
- Error message sanitization

Run tests with:
```bash
go test ./internal/utils -v
```
