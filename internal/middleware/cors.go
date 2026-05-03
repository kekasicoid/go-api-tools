// internal/middleware/cors.go
package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

func CORS() gin.HandlerFunc {
	corsOrigin := strings.TrimSpace(os.Getenv("CORS_ORIGIN"))
	corsOrigin = strings.Trim(corsOrigin, "\"'")
	allowAllOrigins := false
	allowedOrigins := map[string]struct{}{}

	if corsOrigin == "" || corsOrigin == "*" {
		allowAllOrigins = true // Default to allow all origins if not set
		logger.Log.Info("CORS: allow all origins", zap.String("env", corsOrigin))
	} else {
		for _, origin := range strings.Split(corsOrigin, ",") {
			trimmed := strings.TrimSpace(origin)
			trimmed = strings.Trim(trimmed, "\"'")
			trimmed = strings.TrimRight(trimmed, "/")
			if trimmed == "*" {
				allowAllOrigins = true
				break
			}
			if trimmed != "" {
				allowedOrigins[trimmed] = struct{}{}
			}
		}

		if !allowAllOrigins && len(allowedOrigins) == 0 {
			allowAllOrigins = true
		}
		logger.Log.Info(
			"CORS: loaded allowed origins",
			zap.Int("count", len(allowedOrigins)),
			zap.Any("allowedOrigins", allowedOrigins),
		)
	}

	return func(c *gin.Context) {
		requestOrigin := strings.TrimRight(strings.TrimSpace(c.GetHeader("Origin")), "/")
		isCORSRequest := requestOrigin != ""

		if isCORSRequest && allowAllOrigins {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if isCORSRequest {
			if _, ok := allowedOrigins[requestOrigin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
				c.Writer.Header().Set("Vary", "Origin")
			} else {
				logger.Log.Info(
					"CORS: request origin rejected",
					zap.String("requestOrigin", requestOrigin),
					zap.Any("allowedOrigins", allowedOrigins),
				)
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, request-id")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
