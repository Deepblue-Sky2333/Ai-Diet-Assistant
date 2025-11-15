package middleware

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"golang.org/x/text/unicode/norm"
)

// SanitizeMiddleware 输入清理中间件
// 清理HTML标签、SQL特殊字符和规范化Unicode字符
func SanitizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理 POST、PUT、PATCH 请求
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
			c.Next()
			return
		}

		// 只处理 JSON 请求
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.Next()
			return
		}

		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "failed to read request body", err))
			c.Abort()
			return
		}

		// 关闭原始请求体
		c.Request.Body.Close()

		// 如果请求体为空，继续处理
		if len(body) == 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}

		// 解析 JSON
		var data interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			// 如果不是有效的 JSON，恢复原始请求体并继续
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}

		// 清理数据
		sanitized := sanitizeJSONValue(data)

		// 重新编码为 JSON
		sanitizedBody, err := json.Marshal(sanitized)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to sanitize request", err))
			c.Abort()
			return
		}

		// 替换请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
		c.Request.ContentLength = int64(len(sanitizedBody))

		c.Next()
	}
}

// sanitizeJSONValue 递归清理值
func sanitizeJSONValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return sanitizeString(v)
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = sanitizeJSONValue(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = sanitizeJSONValue(val)
		}
		return result
	default:
		return v
	}
}

// sanitizeString 清理字符串
func sanitizeString(s string) string {
	// 1. 规范化 Unicode 字符（防止 Unicode 混淆攻击）
	s = normalizeUnicode(s)

	// 2. 清理 HTML 标签和实体（防止 XSS）
	s = stripHTMLTags(s)
	s = html.UnescapeString(s)
	s = html.EscapeString(s)

	// 3. 清理 SQL 特殊字符（防止 SQL 注入）
	// 注意：这只是额外的防护层，主要防护应该在数据库层使用参数化查询
	s = sanitizeSQLChars(s)

	// 4. 移除控制字符（除了常见的空白字符）
	s = removeControlChars(s)

	// 5. 修剪空白字符
	s = strings.TrimSpace(s)

	return s
}

// normalizeUnicode 规范化 Unicode 字符
func normalizeUnicode(s string) string {
	// 使用 NFC (Canonical Decomposition, followed by Canonical Composition)
	// 这可以防止使用视觉上相似但不同的 Unicode 字符进行攻击
	return norm.NFC.String(s)
}

// stripHTMLTags 移除 HTML 标签
func stripHTMLTags(s string) string {
	// 移除 HTML 标签
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	s = htmlTagRegex.ReplaceAllString(s, "")

	// 移除 HTML 注释
	htmlCommentRegex := regexp.MustCompile(`<!--.*?-->`)
	s = htmlCommentRegex.ReplaceAllString(s, "")

	// 移除 script 标签及其内容
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	s = scriptRegex.ReplaceAllString(s, "")

	// 移除 style 标签及其内容
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	s = styleRegex.ReplaceAllString(s, "")

	// 移除 javascript: 协议
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	s = jsProtocolRegex.ReplaceAllString(s, "")

	// 移除 data: 协议（可能用于 XSS）
	dataProtocolRegex := regexp.MustCompile(`(?i)data:text/html`)
	s = dataProtocolRegex.ReplaceAllString(s, "")

	return s
}

// sanitizeSQLChars 清理潜在的 SQL 特殊字符
// 注意：这不是主要的 SQL 注入防护，主要防护应该使用参数化查询
func sanitizeSQLChars(s string) string {
	// 移除 SQL 注释标记
	s = strings.ReplaceAll(s, "--", "")
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")

	// 移除多个连续的分号（可能用于 SQL 注入）
	semicolonRegex := regexp.MustCompile(`;{2,}`)
	s = semicolonRegex.ReplaceAllString(s, ";")

	return s
}

// removeControlChars 移除控制字符
func removeControlChars(s string) string {
	var result strings.Builder
	for _, r := range s {
		// 保留常见的空白字符：空格、制表符、换行符、回车符
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
			continue
		}
		// 移除其他控制字符
		if !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SanitizeQueryParams 清理查询参数
// 可以在需要的路由上单独使用
func SanitizeQueryParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取所有查询参数
		queryParams := c.Request.URL.Query()

		// 清理每个参数
		for key, values := range queryParams {
			for i, value := range values {
				values[i] = sanitizeString(value)
			}
			queryParams[key] = values
		}

		// 更新查询参数
		c.Request.URL.RawQuery = queryParams.Encode()

		c.Next()
	}
}
