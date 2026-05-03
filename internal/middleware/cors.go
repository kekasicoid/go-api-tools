// internal/middleware/cors.go
package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	corsOrigin := strings.TrimSpace(os.Getenv("CORS_ORIGIN"))
	corsOrigin = strings.Trim(corsOrigin, "\"'")
	allowAllOrigins := false
	allowedOrigins := map[string]struct{}{}

	if corsOrigin == "" || corsOrigin == "*" {
		allowAllOrigins = true // Default to allow all origins if not set
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
	}

	return func(c *gin.Context) {
		requestOrigin := strings.TrimRight(strings.TrimSpace(c.GetHeader("Origin")), "/")

		if allowAllOrigins {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if _, ok := allowedOrigins[requestOrigin]; ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
			c.Writer.Header().Set("Vary", "Origin")
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
