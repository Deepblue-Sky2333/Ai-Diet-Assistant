package utils

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// dsnPattern 匹配数据库连接字符串的正则表达式
	// 匹配格式: user:password@tcp(host:port)/dbname?params
	dsnPattern = regexp.MustCompile(`([^:]+):([^@]+)@tcp\(([^)]+)\)/([^?]+)(\?.*)?`)

	// filePathPattern 匹配文件路径的正则表达式
	filePathPattern = regexp.MustCompile(`(/[a-zA-Z0-9_\-./]+)+`)

	// sqlPattern 匹配SQL语句的正则表达式
	sqlPattern = regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|CREATE|DROP|ALTER|TRUNCATE)\s+.+`)
)

// 错误码常量 - 标准错误码
// 根据需求文档 3.2，系统使用以下标准错误码：
// - 40001: 参数错误
// - 40101: 未授权
// - 40401: 资源不存在
// - 42901: 请求过于频繁
// - 50001: 内部错误
// - 50002: 数据库错误
// - 50003: 外部服务错误（AI服务等）
const (
	CodeSuccess         = 0     // 成功
	CodeInvalidParams   = 40001 // 参数错误（包括验证错误）
	CodeUnauthorized    = 40101 // 未授权（包括认证失败、Token过期等）
	CodeForbidden       = 40301 // 禁止访问（权限不足）
	CodeNotFound        = 40401 // 资源不存在
	CodeConflict        = 40901 // 资源冲突（如用户名已存在）
	CodeTooManyRequests = 42901 // 请求过于频繁（限流）
	CodeInternalError   = 50001 // 内部错误（包括加密错误等）
	CodeDatabaseError   = 50002 // 数据库错误
	CodeAIServiceError  = 50003 // AI服务错误（外部服务错误）
)

// 错误消息映射 - 提供默认的错误消息
var errorMessages = map[int]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "invalid parameters",
	CodeUnauthorized:    "unauthorized",
	CodeForbidden:       "forbidden",
	CodeNotFound:        "resource not found",
	CodeConflict:        "conflict",
	CodeTooManyRequests: "too many requests",
	CodeInternalError:   "internal server error",
	CodeDatabaseError:   "database error",
	CodeAIServiceError:  "external service error",
}

// AppError 应用错误结构
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, err error) *AppError {
	if message == "" {
		message = errorMessages[code]
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
	Timestamp  int64       `json:"timestamp"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// SuccessPaginated 分页成功响应
func SuccessPaginated(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Code:       CodeSuccess,
		Message:    "success",
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now().Unix(),
	})
}

// SuccessWithPagination 分页成功响应（别名）
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	SuccessPaginated(c, data, pagination)
}

// Error 错误响应
// 统一的错误处理函数，确保：
// 1. 使用标准错误码
// 2. 记录详细的错误日志
// 3. 生产环境不暴露敏感信息
// 4. 开发环境提供详细的调试信息
func Error(c *gin.Context, appErr *AppError) {
	statusCode := getHTTPStatusCode(appErr.Code)

	// 检查是否为生产环境（release模式）
	isProduction := gin.Mode() == gin.ReleaseMode

	response := Response{
		Code:      appErr.Code,
		Message:   sanitizeErrorMessage(appErr.Message),
		Timestamp: time.Now().Unix(),
	}

	// 获取用户信息用于日志记录
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	// 构建日志字段
	logFields := []zap.Field{
		zap.Int("error_code", appErr.Code),
		zap.String("error_message", appErr.Message),
		zap.Int("http_status", statusCode),
	}

	// 添加请求信息（如果存在）
	if c.Request != nil {
		logFields = append(logFields,
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)
	}

	// 添加用户信息到日志
	if userID != nil {
		if uid, ok := userID.(int64); ok {
			logFields = append(logFields, zap.Int64("user_id", uid))
		}
	}
	if username != nil {
		if uname, ok := username.(string); ok {
			logFields = append(logFields, zap.String("username", uname))
		}
	}

	// 获取logger（如果存在）
	logger, loggerExists := c.Get("logger")
	var zapLogger *zap.Logger
	if loggerExists {
		zapLogger, _ = logger.(*zap.Logger)
	}

	// 记录详细错误到日志
	if appErr.Err != nil {
		logFields = append(logFields, zap.Error(appErr.Err))

		// 如果有logger，记录日志
		if zapLogger != nil {
			// 根据错误类型选择日志级别
			switch appErr.Code {
			case CodeInvalidParams, CodeNotFound:
				// 客户端错误使用 Warn 级别
				zapLogger.Warn("Client error", logFields...)
			case CodeUnauthorized, CodeTooManyRequests:
				// 认证和限流错误使用 Warn 级别
				zapLogger.Warn("Authentication or rate limit error", logFields...)
			case CodeInternalError, CodeDatabaseError, CodeAIServiceError:
				// 服务器错误使用 Error 级别
				zapLogger.Error("Server error", logFields...)
			default:
				// 未知错误使用 Error 级别
				zapLogger.Error("Unknown error", logFields...)
			}
		}

		// 生产环境不返回详细错误信息
		if isProduction {
			// 仅返回通用错误消息
			response.Error = getGenericErrorMessage(appErr.Code)
		} else {
			// 开发环境返回清理后的错误信息
			response.Error = sanitizeError(appErr.Err.Error())
		}
	} else {
		// 没有底层错误时也记录日志
		if zapLogger != nil {
			if appErr.Code >= 50000 {
				zapLogger.Error("Server error without underlying error", logFields...)
			} else {
				zapLogger.Warn("Client error", logFields...)
			}
		}
	}

	c.JSON(statusCode, response)
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	statusCode := getHTTPStatusCode(code)
	c.JSON(statusCode, Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// getHTTPStatusCode 根据业务错误码获取 HTTP 状态码
func getHTTPStatusCode(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeInvalidParams:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInternalError:
		return http.StatusInternalServerError
	case CodeDatabaseError:
		return http.StatusInternalServerError
	case CodeAIServiceError:
		return http.StatusInternalServerError
	default:
		// 对于未知错误码，根据范围返回合适的HTTP状态码
		switch {
		case code >= 40000 && code < 41000:
			return http.StatusBadRequest
		case code >= 50000:
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
		}
	}
}

// CalculatePagination 计算分页信息
func CalculatePagination(page, pageSize, total int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// sanitizeErrorMessage 清理错误消息中的敏感信息
// 将数据库连接字符串中的密码替换为 ***
func sanitizeErrorMessage(message string) string {
	// 使用正则表达式匹配并替换DSN中的密码部分
	// 匹配格式: user:password@tcp(host:port)/dbname?params
	// 替换为: user:***@tcp(host:port)/dbname?params
	return dsnPattern.ReplaceAllString(message, "$1:***@tcp($3)/$4$5")
}

// sanitizeError 清理错误消息中的敏感信息
// 移除文件路径、SQL语句、数据库连接字符串等敏感信息
func sanitizeError(errMsg string) string {
	// 1. 清理数据库连接字符串
	sanitized := dsnPattern.ReplaceAllString(errMsg, "$1:***@tcp($3)/$4$5")

	// 2. 清理文件路径（保留文件名，移除完整路径）
	sanitized = filePathPattern.ReplaceAllStringFunc(sanitized, func(path string) string {
		// 如果路径包含敏感目录，只保留文件名
		if strings.Contains(path, "/") {
			return filepath.Base(path)
		}
		return path
	})

	// 3. 清理SQL语句（替换为通用消息）
	if sqlPattern.MatchString(sanitized) {
		sanitized = sqlPattern.ReplaceAllString(sanitized, "[SQL query]")
	}

	return sanitized
}

// getGenericErrorMessage 根据错误码返回通用错误消息
// 用于生产环境，不暴露内部实现细节
// 只返回标准错误码的通用消息
func getGenericErrorMessage(code int) string {
	switch code {
	case CodeInvalidParams:
		return "The request contains invalid parameters"
	case CodeUnauthorized:
		return "Authentication is required to access this resource"
	case CodeForbidden:
		return "You do not have permission to access this resource"
	case CodeNotFound:
		return "The requested resource was not found"
	case CodeTooManyRequests:
		return "Too many requests, please try again later"
	case CodeInternalError:
		return "An internal error occurred, please try again later"
	case CodeDatabaseError:
		return "A database error occurred, please try again later"
	case CodeAIServiceError:
		return "An external service error occurred, please try again later"
	default:
		// 对于未知错误码，根据范围返回合适的消息
		switch {
		case code >= 40000 && code < 50000:
			return "The request could not be processed"
		case code >= 50000:
			return "An internal error occurred, please try again later"
		default:
			return "An unexpected error occurred"
		}
	}
}

// SanitizeError 公开的错误清理函数
// 可在其他包中使用，用于清理错误消息中的敏感信息
func SanitizeError(errMsg string) string {
	return sanitizeError(errMsg)
}

// 便捷的错误创建函数 - 使用标准错误码

// NewInvalidParamsError 创建参数错误
func NewInvalidParamsError(message string, err error) *AppError {
	if message == "" {
		message = "invalid parameters"
	}
	return NewAppError(CodeInvalidParams, message, err)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string, err error) *AppError {
	if message == "" {
		message = "unauthorized"
	}
	return NewAppError(CodeUnauthorized, message, err)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string, err error) *AppError {
	if message == "" {
		message = "forbidden"
	}
	return NewAppError(CodeForbidden, message, err)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string, err error) *AppError {
	if message == "" {
		message = "resource not found"
	}
	return NewAppError(CodeNotFound, message, err)
}

// NewTooManyRequestsError 创建限流错误
func NewTooManyRequestsError(message string, err error) *AppError {
	if message == "" {
		message = "too many requests"
	}
	return NewAppError(CodeTooManyRequests, message, err)
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err error) *AppError {
	if message == "" {
		message = "internal server error"
	}
	return NewAppError(CodeInternalError, message, err)
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(message string, err error) *AppError {
	if message == "" {
		message = "database error"
	}
	return NewAppError(CodeDatabaseError, message, err)
}

// NewAIServiceError 创建AI服务错误
func NewAIServiceError(message string, err error) *AppError {
	if message == "" {
		message = "AI service error"
	}
	return NewAppError(CodeAIServiceError, message, err)
}
