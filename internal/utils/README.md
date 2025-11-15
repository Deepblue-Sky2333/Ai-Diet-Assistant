# Utils Package

This package contains utility functions and services for the AI Diet Assistant application.

## Components

### crypto.go
- **CryptoService**: Handles encryption and password hashing
  - `EncryptAES()`: AES-256-GCM encryption for sensitive data (API keys)
  - `DecryptAES()`: AES-256-GCM decryption
  - `HashPassword()`: bcrypt password hashing (cost=12)
  - `VerifyPassword()`: Password verification

### jwt.go
- **JWTService**: JWT token management
  - `GenerateTokenPair()`: Creates Access Token (15min) and Refresh Token (7days)
  - `ValidateToken()`: Validates and parses JWT tokens
  - `RefreshAccessToken()`: Generates new Access Token from Refresh Token
  - **Claims**: JWT claims structure with UserID and Username

### response.go
- **Response Structures**: Unified API response formats
  - `Response`: Standard API response
  - `PaginatedResponse`: Response with pagination metadata
  - `AppError`: Application error structure
- **Helper Functions**:
  - `Success()`, `SuccessWithMessage()`: Success responses
  - `SuccessPaginated()`: Paginated success response
  - `Error()`, `ErrorWithMessage()`: Error responses
  - `CalculatePagination()`: Pagination calculation
- **Error Codes**: Standardized error codes (40001-50004)

## Usage Examples

### Encryption
```go
cryptoService, _ := utils.NewCryptoService("your-32-byte-encryption-key-here")
encrypted, _ := cryptoService.EncryptAES("sensitive-data")
decrypted, _ := cryptoService.DecryptAES(encrypted)
```

### JWT
```go
jwtService := utils.NewJWTService("secret", 15, 168) // 15h access, 168h refresh
tokens, _ := jwtService.GenerateTokenPair(userID, username)
claims, _ := jwtService.ValidateToken(tokens.AccessToken)
```

### Response
```go
utils.Success(c, data)
utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "unauthorized", nil))
```
