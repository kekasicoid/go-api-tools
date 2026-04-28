// internal/delivery/http/router.go
package http

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	return gin.Default()
}

func RegisterRoutes(r *gin.Engine, h *Handler) {

	tools := r.Group("/tools")
	tools.POST("/json/format", h.FormatJSON)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation only in development environment
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "development" || appEnv == "dev" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
