// internal/usecase/jwt_usecase.go
package usecase

import (
	"github.com/kekasicoid/go-api-tools/internal/domain"
)

// JWTUsecase handles JWT decode and validation business logic.
type JWTUsecase struct {
	decoder domain.JWTDecoder
}

// NewJWTUsecase creates a new JWTUsecase.
func NewJWTUsecase(d domain.JWTDecoder) *JWTUsecase {
	return &JWTUsecase{decoder: d}
}

// DecodeJWT parses the JWT without verifying its signature.
func (u *JWTUsecase) DecodeJWT(token string) (header map[string]interface{}, claims map[string]interface{}, err error) {
	return u.decoder.Decode(token)
}

// ValidateJWT verifies the JWT signature with the provided HMAC secret.
func (u *JWTUsecase) ValidateJWT(token string, secret string) (claims map[string]interface{}, err error) {
	return u.decoder.Validate(token, secret)
}
