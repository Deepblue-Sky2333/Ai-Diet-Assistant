# Middleware Package

This package contains HTTP middleware for the AI Diet Assistant application.

## Components

### auth.go
- **AuthMiddleware**: JWT authentication middleware
  - Validates Bearer tokens from Authorization header
  - Extracts user information and injects into context
  - Handles token expiration and invalid tokens
- **Helper Functions**:
  - `GetUserID()`: Retrieves user ID from context
  - `GetUsername()`: Retrieves username from context
  - `MustGetUserID()`: Retrieves user ID or panics

### cors.go
- **CORSMiddleware**: Cross-Origin Resource Sharing middleware
  - Validates origin against allowed list
  - Sets CORS headers (Allow-Origin, Allow-Methods, etc.)
  - Handles OPTIONS preflight requests
  - Supports wildcard domains (*.example.com)

### ratelimit.go
- **RateLimitMiddleware**: Request rate limiting middleware
  - Limits requests per user (100 requests/minute by default)
  - Supports both Redis (distributed) and memory (single-instance) storage
  - Automatic fallback to memory storage if Redis is unavailable
  - Uses sliding window algorithm for accurate rate limiting
  - Falls back to IP-based limiting for unauthenticated requests
  - Returns 429 status when limit exceeded

### logger.go
- **LoggerMiddleware**: Structured logging middleware
  - Logs all HTTP requests with structured fields
  - Records method, path, status, duration, user info, IP
  - Sanitizes sensitive information (passwords, tokens, API keys)
  - Uses different log levels based on status code

### recovery.go
- **RecoveryMiddleware**: Panic recovery middleware
  - Catches panics and prevents server crashes
  - Logs full stack trace with context information
  - Returns 500 error response to client
  - Ensures service continues running

### upload.go
- **FileValidationMiddleware**: File upload validation middleware
  - Validates file MIME types against allowed list
  - Checks file extensions to prevent malicious uploads
  - Validates file content to prevent extension spoofing
  - Enforces maximum file size limits
  - Supports custom allowed types and extensions
- **RequestSizeLimitMiddleware**: Request body size limiting middleware
  - Limits total request body size (default: 10MB globally)
  - Prevents memory exhaustion from large uploads
  - Can be configured per-route for different limits
  - Returns clear error messages when limit exceeded

### sanitize.go
- **SanitizeMiddleware**: Input sanitization middleware
  - Cleans HTML tags and entities to prevent XSS attacks
  - Removes SQL special characters as additional protection layer
  - Normalizes Unicode characters to prevent confusion attacks
  - Removes control characters (except common whitespace)
  - Processes JSON request bodies recursively
  - Only applies to POST, PUT, and PATCH requests with JSON content
- **SanitizeQueryParams**: Query parameter sanitization middleware
  - Cleans query parameters using the same sanitization rules
  - Can be applied to specific routes that need query param sanitization

## Middleware Chain Order

Recommended order for middleware registration:

```go
router := gin.New()
router.Use(RecoveryMiddleware(logger))  // 1. Catch panics first
router.Use(LoggerMiddleware(logger))    // 2. Log all requests
router.Use(CORSMiddleware(&corsConfig)) // 3. Handle CORS
router.Use(RateLimitMiddleware(&rateLimitConfig, &redisConfig, logger)) // 4. Rate limiting
router.Use(RequestSizeLimitMiddleware(10*1024*1024, logger)) // 5. Global request size limit (10MB)
router.Use(SanitizeMiddleware())        // 6. Input sanitization (XSS/SQL injection prevention)

// Protected routes
protected := router.Group("/api/v1")
protected.Use(AuthMiddleware(jwtService)) // 7. Authentication for protected routes

// File upload routes (with specific validation)
uploadConfig := FileValidationConfig{
    MaxFileSize: 5 * 1024 * 1024, // 5MB
    ValidateContent: true,
}
protected.POST("/upload", FileValidationMiddleware(uploadConfig, logger), uploadHandler)
```

## Usage Examples

### Authentication
```go
// Apply to protected routes
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(jwtService))

// In handler, get user ID
userID := middleware.MustGetUserID(c)
```

### CORS
```go
corsConfig := &config.CORSConfig{
    AllowedOrigins: []string{"http://localhost:3000", "https://example.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
}
router.Use(middleware.CORSMiddleware(corsConfig))
```

### Rate Limiting
```go
rateLimitConfig := &config.RateLimitConfig{
    Enabled: true,
    RequestsPerMinute: 100,
    StorageType: "redis", // or "memory"
}
redisConfig := &config.RedisConfig{
    Enabled: true,
    Host: "localhost",
    Port: 6379,
}
router.Use(middleware.RateLimitMiddleware(rateLimitConfig, redisConfig, logger))
```

### File Upload Validation
```go
// Global request size limit (applied to all routes)
router.Use(middleware.RequestSizeLimitMiddleware(10*1024*1024, logger)) // 10MB

// File upload with validation (applied to specific routes)
uploadConfig := middleware.FileValidationConfig{
    MaxFileSize: 5 * 1024 * 1024, // 5MB per file
    ValidateContent: true, // Validate file content matches extension
    // Optional: custom allowed types
    AllowedMimeTypes: []string{"image/jpeg", "image/png"},
    AllowedExtensions: []string{".jpg", ".jpeg", ".png"},
}
router.POST("/api/v1/upload", 
    middleware.FileValidationMiddleware(uploadConfig, logger),
    uploadHandler,
)

// Use default allowed types (images, documents, JSON)
defaultConfig := middleware.FileValidationConfig{
    MaxFileSize: 10 * 1024 * 1024, // 10MB
    ValidateContent: true,
}
router.POST("/api/v1/import", 
    middleware.FileValidationMiddleware(defaultConfig, logger),
    importHandler,
)
```

### Request Size Limiting for Specific Routes
```go
// Different size limits for different endpoints
router.POST("/api/v1/small-upload", 
    middleware.RequestSizeLimitMiddleware(1*1024*1024, logger), // 1MB
    smallUploadHandler,
)

router.POST("/api/v1/large-upload", 
    middleware.RequestSizeLimitMiddleware(50*1024*1024, logger), // 50MB
    largeUploadHandler,
)
```

### Input Sanitization
```go
// Global sanitization (applied to all JSON POST/PUT/PATCH requests)
router.Use(middleware.SanitizeMiddleware())

// Query parameter sanitization for specific routes
router.GET("/api/v1/search", 
    middleware.SanitizeQueryParams(),
    searchHandler,
)

// Example of what gets sanitized:
// Input:  {"name": "<script>alert('xss')</script>Hello", "note": "Test--comment"}
// Output: {"name": "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;Hello", "note": "Testcomment"}
```

## Security Considerations

### Input Sanitization
1. **XSS Prevention**: The sanitize middleware removes HTML tags, scripts, and dangerous protocols
2. **SQL Injection**: While sanitization provides an additional layer, always use parameterized queries
3. **Unicode Normalization**: Prevents attacks using visually similar Unicode characters
4. **Automatic Processing**: Applied globally to all JSON POST/PUT/PATCH requests
5. **Query Parameters**: Use `SanitizeQueryParams()` for routes that need query string sanitization

### File Upload Security
1. **Always validate file content**: Set `ValidateContent: true` to prevent extension spoofing
2. **Limit file sizes**: Set appropriate `MaxFileSize` to prevent DoS attacks
3. **Restrict file types**: Only allow necessary file types for your use case
4. **Store files securely**: Never store uploaded files in web-accessible directories
5. **Scan for malware**: Consider integrating virus scanning for uploaded files

### Request Size Limiting
1. **Global limits**: Apply a reasonable global limit (10MB) to all routes
2. **Route-specific limits**: Use stricter limits for non-upload routes
3. **Monitor usage**: Log rejected requests to detect potential attacks
4. **Adjust as needed**: Tune limits based on your application's requirements
