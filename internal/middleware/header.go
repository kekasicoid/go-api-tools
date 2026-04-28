package middleware

import "github.com/gin-gonic/gin"

func GetCtx(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}
