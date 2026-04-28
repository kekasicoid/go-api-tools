package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/internal/model"
)

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9]{1,50}$`)

func ValidateRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/health" || strings.HasPrefix(path, "/swagger") {
			c.Next()
			return
		}

		requestID := c.GetHeader(model.HeadRequestIDKey)
		if requestID == "" {
			model.RespBadRequest(c, "request-id header is required")
			c.Abort()
			return
		}

		if !requestIDPattern.MatchString(requestID) {
			model.RespBadRequest(c, "request-id must be alphanumeric, without spaces, and max length 50")
			c.Abort()
			return
		}

		c.Set(model.HeadRequestIDKey, requestID)

		c.Next()
	}
}
