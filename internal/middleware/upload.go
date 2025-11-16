package middleware

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AllowedFileTypes 定义允许的文件类型
var AllowedFileTypes = map[string][]string{
	// 图片类型
	"image/jpeg": {".jpg", ".jpeg"},
	"image/png":  {".png"},
	"image/gif":  {".gif"},
	"image/webp": {".webp"},
	// 文档类型
	"application/pdf":          {".pdf"},
	"application/vnd.ms-excel": {".xls"},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {".xlsx"},
	"text/csv":   {".csv"},
	"text/plain": {".txt"},
	// JSON 类型（用于批量导入）
	"application/json": {".json"},
}

// FileValidationConfig 文件验证配置
type FileValidationConfig struct {
	// MaxFileSize 最大文件大小（字节）
	MaxFileSize int64
	// AllowedMimeTypes 允许的 MIME 类型列表，如果为空则使用默认的 AllowedFileTypes
	AllowedMimeTypes []string
	// AllowedExtensions 允许的文件扩展名列表，如果为空则使用默认的 AllowedFileTypes
	AllowedExtensions []string
	// ValidateContent 是否验证文件内容（防止伪造扩展名）
	ValidateContent bool
}

// FileValidationMiddleware 文件上传验证中间件
func FileValidationMiddleware(config FileValidationConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理包含文件上传的请求
		if !strings.HasPrefix(c.ContentType(), "multipart/form-data") {
			c.Next()
			return
		}

		// 解析 multipart form
		if err := c.Request.ParseMultipartForm(config.MaxFileSize); err != nil {
			logger.Warn("Failed to parse multipart form",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("remote_addr", c.ClientIP()),
			)
			utils.Error(c, utils.NewAppError(
				utils.CodeInvalidParams,
				"failed to parse upload form, file may be too large",
				err,
			))
			c.Abort()
			return
		}

		// 验证所有上传的文件
		if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
			for fieldName, files := range c.Request.MultipartForm.File {
				for _, fileHeader := range files {
					if err := validateFile(fileHeader, config, logger); err != nil {
						logger.Warn("File validation failed",
							zap.Error(err),
							zap.String("field", fieldName),
							zap.String("filename", fileHeader.Filename),
							zap.Int64("size", fileHeader.Size),
							zap.String("remote_addr", c.ClientIP()),
						)
						utils.Error(c, utils.NewAppError(
							utils.CodeInvalidParams,
							err.Error(),
							err,
						))
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// validateFile 验证单个文件
func validateFile(fileHeader *multipart.FileHeader, config FileValidationConfig, logger *zap.Logger) error {
	// 1. 检查文件大小
	if fileHeader.Size > config.MaxFileSize {
		return utils.NewAppError(
			utils.CodeInvalidParams,
			"file size exceeds maximum allowed size",
			nil,
		)
	}

	// 2. 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isExtensionAllowed(ext, config) {
		return utils.NewAppError(
			utils.CodeInvalidParams,
			"file type not allowed",
			nil,
		)
	}

	// 3. 验证文件内容（防止伪造扩展名）
	if config.ValidateContent {
		file, err := fileHeader.Open()
		if err != nil {
			return utils.NewAppError(
				utils.CodeInternalError,
				"failed to open uploaded file",
				err,
			)
		}
		defer file.Close()

		// 读取文件头部（前 512 字节足以检测 MIME 类型）
		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return utils.NewAppError(
				utils.CodeInternalError,
				"failed to read file content",
				err,
			)
		}

		// 检测实际的 MIME 类型
		detectedMimeType := http.DetectContentType(buffer[:n])

		// 验证 MIME 类型
		if !isMimeTypeAllowed(detectedMimeType, config) {
			logger.Warn("File content does not match extension",
				zap.String("filename", fileHeader.Filename),
				zap.String("extension", ext),
				zap.String("detected_mime", detectedMimeType),
			)
			return utils.NewAppError(
				utils.CodeInvalidParams,
				"file content does not match file extension",
				nil,
			)
		}
	}

	return nil
}

// isExtensionAllowed 检查文件扩展名是否允许
func isExtensionAllowed(ext string, config FileValidationConfig) bool {
	// 如果配置了自定义扩展名列表，使用自定义列表
	if len(config.AllowedExtensions) > 0 {
		for _, allowed := range config.AllowedExtensions {
			if strings.ToLower(allowed) == ext {
				return true
			}
		}
		return false
	}

	// 否则使用默认的 AllowedFileTypes
	for _, extensions := range AllowedFileTypes {
		for _, allowedExt := range extensions {
			if allowedExt == ext {
				return true
			}
		}
	}
	return false
}

// isMimeTypeAllowed 检查 MIME 类型是否允许
func isMimeTypeAllowed(mimeType string, config FileValidationConfig) bool {
	// 标准化 MIME 类型（移除参数部分，如 charset）
	mimeType = strings.Split(mimeType, ";")[0]
	mimeType = strings.TrimSpace(mimeType)

	// 如果配置了自定义 MIME 类型列表，使用自定义列表
	if len(config.AllowedMimeTypes) > 0 {
		for _, allowed := range config.AllowedMimeTypes {
			if strings.HasPrefix(mimeType, allowed) {
				return true
			}
		}
		return false
	}

	// 否则使用默认的 AllowedFileTypes
	for allowedMime := range AllowedFileTypes {
		if strings.HasPrefix(mimeType, allowedMime) {
			return true
		}
	}
	return false
}

// RequestSizeLimitMiddleware 请求体大小限制中间件
func RequestSizeLimitMiddleware(maxSize int64, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置请求体大小限制
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		// 尝试读取一个字节来触发大小检查
		var buf bytes.Buffer
		_, err := io.CopyN(&buf, c.Request.Body, 1)

		if err != nil && err.Error() == "http: request body too large" {
			logger.Warn("Request body too large",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int64("max_size", maxSize),
				zap.String("remote_addr", c.ClientIP()),
			)
			utils.Error(c, utils.NewAppError(
				utils.CodeInvalidParams,
				"request body too large",
				err,
			))
			c.Abort()
			return
		}

		// 如果成功读取了一个字节，需要将其放回
		if buf.Len() > 0 {
			c.Request.Body = io.NopCloser(io.MultiReader(&buf, c.Request.Body))
		}

		c.Next()
	}
}
