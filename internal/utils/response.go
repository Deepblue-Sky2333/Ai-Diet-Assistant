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

// 错误码常量
const (
	CodeSuccess          = 0
	CodeInvalidParams    = 40001
	CodeUnauthorized     = 40101
	CodeForbidden        = 40301
	CodeNotFound         = 40401
	CodeConflict         = 40901
	CodeTooManyRequests  = 42901
	CodeInternalError    = 50001
	CodeDatabaseError    = 50002
	CodeAIServiceError   = 50003
	CodeEncryptionError  = 50004
	CodeValidationError  = 40002
)

// 错误消息映射
var errorMessages = map[int]string{
	CodeSuccess:          "success",
	CodeInvalidParams:    "invalid parameters",
	CodeUnauthorized:     "unauthorized",
	CodeForbidden:        "forbidden",
	CodeNotFound:         "resource not found",
	CodeConflict:         "resource conflict",
	CodeTooManyRequests:  "too many requests",
	CodeInternalError:    "internal server error",
	CodeDatabaseError:    "database error",
	CodeAIServiceError:   "AI service error",
	CodeEncryptionError:  "encryption error",
	CodeValidationError:  "validation error",
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
func Error(c *gin.Context, appErr *AppError) {
	statusCode := getHTTPStatusCode(appErr.Code)
	
	// 检查是否为生产环境（release模式）
	isProduction := gin.Mode() == gin.ReleaseMode
	
	response := Response{
		Code:      appErr.Code,
		Message:   sanitizeErrorMessage(appErr.Message),
		Timestamp: time.Now().Unix(),
	}
	
	// 记录详细错误到日志
	if appErr.Err != nil {
		// 获取logger（如果存在）
		if logger, exists := c.Get("logger"); exists {
			if zapLogger, ok := logger.(*zap.Logger); ok {
				zapLogger.Error("Request error",
					zap.Int("code", appErr.Code),
					zap.String("message", appErr.Message),
					zap.Error(appErr.Err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("client_ip", c.ClientIP()),
				)
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
	switch {
	case code >= 40000 && code < 41000:
		return http.StatusBadRequest
	case code >= 40100 && code < 40200:
		return http.StatusUnauthorized
	case code >= 40300 && code < 40400:
		return http.StatusForbidden
	case code >= 40400 && code < 40500:
		return http.StatusNotFound
	case code >= 40900 && code < 41000:
		return http.StatusConflict
	case code >= 42900 && code < 43000:
		return http.StatusTooManyRequests
	case code >= 50000:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
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
func getGenericErrorMessage(code int) string {
	switch {
	case code >= 40000 && code < 41000:
		return "The request could not be processed due to invalid input"
	case code >= 40100 && code < 40200:
		return "Authentication is required to access this resource"
	case code >= 40300 && code < 40400:
		return "You do not have permission to access this resource"
	case code >= 40400 && code < 40500:
		return "The requested resource was not found"
	case code >= 40900 && code < 41000:
		return "The request conflicts with the current state of the resource"
	case code >= 42900 && code < 43000:
		return "Too many requests, please try again later"
	case code >= 50000 && code < 50100:
		return "An internal error occurred, please try again later"
	case code >= 50100 && code < 50200:
		return "A database error occurred, please try again later"
	case code >= 50200 && code < 50300:
		return "An external service error occurred, please try again later"
	default:
		return "An unexpected error occurred"
	}
}

// SanitizeError 公开的错误清理函数
// 可在其他包中使用，用于清理错误消息中的敏感信息
func SanitizeError(errMsg string) string {
	return sanitizeError(errMsg)
}
