package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"go.uber.org/zap"
)

// CORSMiddleware CORS 中间件
func CORSMiddleware(cfg *config.CORSConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// 检查来源是否在允许列表中
		if origin != "" {
			if isOriginAllowed(origin, cfg.AllowedOrigins) {
				// 设置 CORS 响应头
				c.Header("Access-Control-Allow-Origin", origin)
				
				if cfg.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
				
				// 设置允许的方法
				if len(cfg.AllowedMethods) > 0 {
					c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
				}
				
				// 设置允许的头
				if len(cfg.AllowedHeaders) > 0 {
					c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
				}
				
				// 设置暴露的头
				if len(cfg.ExposeHeaders) > 0 {
					c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
				}
				
				// 设置预检请求缓存时间
				if cfg.MaxAge > 0 {
					c.Header("Access-Control-Max-Age", cfg.MaxAge.String())
				}
			} else {
				// 记录被拒绝的来源
				logger.Warn("CORS request rejected: origin not in allowed list",
					zap.String("origin", origin),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("remote_addr", c.ClientIP()),
				)
			}
		}
		
		// 处理 OPTIONS 预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// isOriginAllowed 检查来源是否在允许列表中
// 注意：不再支持通配符 * 以增强安全性
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		// 精确匹配
		if allowed == origin {
			return true
		}
		// 支持通配符子域名匹配，例如 *.example.com
		if strings.HasPrefix(allowed, "*.") {
			domain := allowed[2:]
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}
	return false
}
