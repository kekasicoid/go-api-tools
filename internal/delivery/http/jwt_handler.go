// internal/delivery/http/jwt_handler.go
package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kekasicoid/go-api-tools/internal/model"
	"github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

// JWTHandler handles JWT decode-validation requests.
type JWTHandler struct {
	usecase *usecase.JWTUsecase
}

// NewJWTHandler creates a new JWTHandler.
func NewJWTHandler(u *usecase.JWTUsecase) *JWTHandler {
	return &JWTHandler{usecase: u}
}

// DecodeValidateJWT godoc
// @Summary      Decode and validate JWT token
// @Description  Parse a JWT token and return its header and claims. When a secret is provided the signature is also verified and the valid field reflects the result.
// @Tags         tools
// @Accept       json
// @Produce      json
// @Param        request-id  header    string                               true  "Request ID (alphanumeric, max 50 chars)"
// @Param        request     body      model.JWTDecodeValidationRequest     true  "JWT token and optional HMAC secret"
// @Success      200         {object}  model.JWTDecodeValidationResponseSwag
// @Failure      400         {object}  model.SwaggRespError
// @Router       /tools/jwt/decode-validation [post]
func (h *JWTHandler) DecodeValidateJWT(c *gin.Context) {
	var req model.JWTDecodeValidationRequest
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

	resp := model.JWTDecodeValidationResponse{
		Header: header,
		Claims: claims,
		Valid:  false,
	}

	if req.Secret != "" {
		_, validationErr := h.usecase.ValidateJWT(req.Token, req.Secret)
		resp.Valid = validationErr == nil
		if validationErr != nil {
			logger.Log.Warn("jwt validation failed", zap.Error(validationErr))
		}
	}

	model.RespSuccess(c, "success", resp)
}
