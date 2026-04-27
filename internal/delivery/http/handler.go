// internal/delivery/http/handler.go
package http

import (
	"net/http"

	"github.com/kekasicoid/go-api-tools/internal/usecase"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.usecase.FormatJSON(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"formatted": result})
}
