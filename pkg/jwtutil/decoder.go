// pkg/jwtutil/decoder.go
package jwtutil

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTDecoder implements the domain.JWTDecoder interface.
type JWTDecoder struct{}

// NewJWTDecoder creates a new JWTDecoder.
func NewJWTDecoder() *JWTDecoder {
	return &JWTDecoder{}
}

// Decode parses a JWT token without verifying its signature.
// It returns the header and claims as generic maps.
func (d *JWTDecoder) Decode(token string) (map[string]interface{}, map[string]interface{}, error) {
	if strings.TrimSpace(token) == "" {
		return nil, nil, errors.New("token is required")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, errors.New("invalid JWT format: must have 3 parts")
	}

	// Parse without verifying signature.
	p := jwt.NewParser()
	parsed, _, err := p.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, nil, errors.New("failed to parse token: " + err.Error())
	}

	header := parsed.Header
	if header == nil {
		return nil, nil, errors.New("failed to extract header")
	}

	mapClaims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("failed to extract claims")
	}

	return header, map[string]interface{}(mapClaims), nil
}

// Validate parses and verifies the JWT token signature using the provided HMAC secret.
// Returns the verified claims on success.
func (d *JWTDecoder) Validate(token string, secret string) (map[string]interface{}, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is required")
	}

	if strings.TrimSpace(secret) == "" {
		return nil, errors.New("secret is required")
	}

	parsed, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			alg, _ := t.Header["alg"].(string)
			return nil, errors.New("unexpected signing method: " + alg)
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, errors.New("token is invalid")
	}

	mapClaims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	return map[string]interface{}(mapClaims), nil
}
