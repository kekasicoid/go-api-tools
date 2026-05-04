package model

// JWTDecodeValidationRequest is the request body for the JWT decode-validation endpoint.
// Secret is optional; when provided the token signature will also be verified.
type JWTDecodeValidationRequest struct {
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"`
	Secret string `json:"secret,omitempty" example:"my-secret-key"`
}

// JWTDecodeValidationResponse is the response for the JWT decode-validation endpoint.
type JWTDecodeValidationResponse struct {
	Header map[string]interface{} `json:"header"`
	Claims map[string]interface{} `json:"claims"`
	Valid  bool                   `json:"valid"`
}

// JWTDecodeValidationResponseSwag is used for Swagger documentation only.
type JWTDecodeValidationResponseSwag struct {
	Code   string                      `json:"response_code" example:"200"`
	Refnum string                      `json:"response_refnum" example:""`
	ID     string                      `json:"response_id" example:"abc-123"`
	Desc   string                      `json:"response_desc" example:"success"`
	Data   JWTDecodeValidationResponse `json:"response_data"`
}
