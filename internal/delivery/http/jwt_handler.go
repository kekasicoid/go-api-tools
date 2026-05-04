// internal/delivery/http/jwt_handler.go
package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/internal/model"
	"github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

// JWTHandler handles JWT decode and validation requests.
type JWTHandler struct {
	usecase *usecase.JWTUsecase
}

// NewJWTHandler creates a new JWTHandler.
func NewJWTHandler(u *usecase.JWTUsecase) *JWTHandler {
	return &JWTHandler{usecase: u}
}

// DecodeJWT godoc
// @Summary      Decode JWT token
// @Description  Parse a JWT token and return its header and claims without verifying the signature
// @Tags         tools
// @Accept       json
// @Produce      json
// @Param        request-id  header    string                    true  "Request ID (alphanumeric, max 50 chars)"
// @Param        request     body      model.JWTDecodeRequest    true  "JWT token to decode"
// @Success      200         {object}  model.JWTDecodeResponseSwag
// @Failure      400         {object}  model.SwaggRespError
// @Router       /tools/jwt/decode [post]
func (h *JWTHandler) DecodeJWT(c *gin.Context) {
	var req model.JWTDecodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request", zap.Error(err))
		model.RespBadRequest(c, "invalid request")
		return
	}

	if req.Token == "" {
		model.RespBadRequest(c, "token is required")
		return
	}

	header, claims, err := h.usecase.DecodeJWT(req.Token)
	if err != nil {
		logger.Log.Error("jwt decode failed", zap.Error(err))
		model.RespBadRequest(c, err.Error())
		return
	}

	model.RespSuccess(c, "success", model.JWTDecodeResponse{
		Header: header,
		Claims: claims,
	})
}

// ValidateJWT godoc
// @Summary      Validate JWT token
// @Description  Verify a JWT token signature using the provided HMAC secret and return its claims
// @Tags         tools
// @Accept       json
// @Produce      json
// @Param        request-id  header    string                      true  "Request ID (alphanumeric, max 50 chars)"
// @Param        request     body      model.JWTValidateRequest    true  "JWT token and HMAC secret"
// @Success      200         {object}  model.JWTValidateResponseSwag
// @Failure      400         {object}  model.SwaggRespError
// @Router       /tools/jwt/validate [post]
func (h *JWTHandler) ValidateJWT(c *gin.Context) {
	var req model.JWTValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("invalid request", zap.Error(err))
		model.RespBadRequest(c, "invalid request")
		return
	}

	if req.Token == "" {
		model.RespBadRequest(c, "token is required")
		return
	}

	if req.Secret == "" {
		model.RespBadRequest(c, "secret is required")
		return
	}

	claims, err := h.usecase.ValidateJWT(req.Token, req.Secret)
	if err != nil {
		logger.Log.Warn("jwt validation failed", zap.Error(err))
		model.RespSuccess(c, "success", model.JWTValidateResponse{Valid: false})
		return
	}

	model.RespSuccess(c, "success", model.JWTValidateResponse{
		Valid:  true,
		Claims: claims,
	})
}
