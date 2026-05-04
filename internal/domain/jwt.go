// internal/domain/jwt.go
package domain

// JWTDecoder defines the contract for JWT decode and validation operations.
type JWTDecoder interface {
	// Decode parses the JWT token without verifying the signature.
	// Returns the header and claims as maps.
	Decode(token string) (header map[string]interface{}, claims map[string]interface{}, err error)

	// Validate parses and verifies the JWT token signature using the provided secret (HMAC).
	// Returns the claims on success.
	Validate(token string, secret string) (claims map[string]interface{}, err error)
}
