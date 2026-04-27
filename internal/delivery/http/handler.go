// internal/delivery/http/handler.go
package http

import (
	"github.com/kekasicoid/go-api-tools/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *usecase.FormatterUsecase
}

func NewHandler(u *usecase.FormatterUsecase) *Handler {
	return &Handler{usecase: u}
}

type formatRequest struct {
	Data string `json:"data"`
}

func (h *Handler) FormatJSON(c *gin.Context) {
	var req formatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request",
			zap.Error(err),
		)

		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.usecase.FormatJSON(req.Data)
	if err != nil {
		logger.Log.Error("format failed",
			zap.Error(err),
		)

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("json formatted successfully")

	c.JSON(200, gin.H{"formatted": result})
}
