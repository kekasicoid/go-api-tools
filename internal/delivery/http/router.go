// internal/delivery/http/router.go
package http

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	r.POST("/tools/json/format", h.FormatJSON)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
