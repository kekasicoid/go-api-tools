// internal/delivery/http/handler.go
package http

import (
	"github.com/kekasicoid/go-api-tools/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/internal/model"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *usecase.FormatterUsecase
}

func NewHandler(u *usecase.FormatterUsecase) *Handler {
	return &Handler{usecase: u}
}

// FormatJSON godoc
// @Summary      Format JSON data
// @Description  Format a raw JSON string into pretty JSON
// @Tags         tools
// @Accept       json
// @Produce      json
// @Param        request-id  header    string                 true  "Request ID (alphanumeric, max 50 chars)"
// @Param        request     body      model.FormatJsonRequest  true  "JSON data to format"
// @Success      200      {object}  model.FormatJsonResponseSwag
// @Failure      400      {object}  model.SwaggRespError
// @Router       /tools/json/format [post]
func (h *Handler) FormatJSON(c *gin.Context) {
	var req model.FormatJsonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request", zap.Error(err))
		model.RespBadRequest(c, "invalid request")
		return
	}

	if req.Data == "" {
		model.RespBadRequest(c, "data is required")
		return
	}

	result, err := h.usecase.FormatJSON(req.Data)
	if err != nil {
		logger.Log.Error("format failed", zap.Error(err))
		model.RespBadRequest(c, err.Error())
		return
	}

	model.RespSuccess(c, "success", model.FormatJsonResponse{Formatted: result})
}
