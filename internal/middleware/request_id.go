package middleware

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kekasicoid/go-api-tools/internal/model"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9]{1,50}$`)

func maskRequestID(requestID string) string {
	if len(requestID) <= 6 {
		return requestID
	}

	return requestID[:3] + "***" + requestID[len(requestID)-3:]
}

func GetRequestIDTTL() time.Duration {
	if val := os.Getenv("REQUEST_ID_TTL_HOURS"); val != "" {
		if hours, err := strconv.Atoi(val); err == nil && hours > 0 {
			return time.Duration(hours) * time.Hour
		}
	}
	return 24 * time.Hour
}

func ValidateRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/health" || strings.HasPrefix(path, "/swagger") {
			c.Next()
			return
		}

		requestID := c.GetHeader(model.HeadRequestIDKey)
		if requestID == "" {
			logger.Log.Info(
				"request-id validation failed: missing header",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)
			model.RespBadRequest(c, "request-id header is required")
			c.Abort()
			return
		}

		if !requestIDPattern.MatchString(requestID) {
			logger.Log.Info(
				"request-id validation failed: invalid format",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("requestID", maskRequestID(requestID)),
			)
			model.RespBadRequest(c, "request-id must be alphanumeric, without spaces, and max length 50")
			c.Abort()
			return
		}

		// Ensure Redis is initialized before use.
		if rdb == nil {
			initRedis()
		}

		redisKey := "request_id:" + requestID
		set, err := rdb.SetNX(ctx, redisKey, 1, GetRequestIDTTL()).Result()
		if err != nil && err != redis.Nil {
			logger.Log.Error(
				"request-id validation failed: redis error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("requestID", maskRequestID(requestID)),
				zap.Error(err),
			)
			model.RespInternalServerError(c, "failed to validate request-id")
			c.Abort()
			return
		}
		if !set {
			logger.Log.Info(
				"request-id validation failed: duplicate",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("requestID", maskRequestID(requestID)),
			)
			model.RespBadRequest(c, "request-id has already been used within the last 24 hours")
			c.Abort()
			return
		}

		logger.Log.Info(
			"request-id validation passed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("requestID", maskRequestID(requestID)),
		)

		c.Set(model.HeadRequestIDKey, requestID)

		c.Next()
	}
}
