// internal/delivery/http/instagram_handler.go
package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/internal/model"
	"github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

// InstagramHandler handles Instagram media download requests.
type InstagramHandler struct {
	usecase *usecase.InstagramUsecase
}

// NewInstagramHandler creates a new InstagramHandler.
func NewInstagramHandler(u *usecase.InstagramUsecase) *InstagramHandler {
	return &InstagramHandler{usecase: u}
}

// DownloadMedia godoc
// @Summary      Download Instagram media
// @Description  Extract direct download URLs for photos, videos, and reels from a public Instagram post URL. Stories are only supported when the account is public.
// @Tags         tools
// @Accept       json
// @Produce      json
// @Param        request-id  header    string                                true  "Request ID (alphanumeric, max 50 chars)"
// @Param        request     body      model.InstagramDownloadRequest        true  "Instagram post URL"
// @Success      200         {object}  model.InstagramDownloadResponseSwag
// @Failure      400         {object}  model.SwaggRespError
// @Failure      500         {object}  model.SwaggRespError
// @Router       /tools/instagram/download [post]
func (h *InstagramHandler) DownloadMedia(c *gin.Context) {
	var req model.InstagramDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request", zap.Error(err))
		model.RespBadRequest(c, "invalid request")
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		model.RespBadRequest(c, "url is required")
		return
	}

	if !strings.Contains(req.URL, "instagram.com") {
		model.RespBadRequest(c, "url must be an Instagram URL")
		return
	}

	info, err := h.usecase.Download(req.URL)
	if err != nil {
		logger.Log.Error("instagram download failed", zap.String("url", req.URL), zap.Error(err))
		model.RespInternalServerError(c, err.Error())
		return
	}

	resp := model.InstagramDownloadResponse{
		PostType: info.PostType,
		Items:    make([]model.InstagramMediaItem, 0, len(info.Items)),
	}
	for _, item := range info.Items {
		resp.Items = append(resp.Items, model.InstagramMediaItem{
			MediaType: item.MediaType,
			MediaURL:  item.MediaURL,
			ThumbURL:  item.ThumbURL,
		})
	}

	model.RespSuccess(c, "success", resp)
}
