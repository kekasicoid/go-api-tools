package model

// JWTDecodeRequest is the request body for the JWT decode endpoint.
type JWTDecodeRequest struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"`
}

// JWTDecodeResponse is the response for the JWT decode endpoint.
type JWTDecodeResponse struct {
	Header map[string]interface{} `json:"header"`
	Claims map[string]interface{} `json:"claims"`
}

// JWTDecodeResponseSwag is used for Swagger documentation only.
type JWTDecodeResponseSwag struct {
	Code   string            `json:"response_code" example:"200"`
	Refnum string            `json:"response_refnum" example:""`
	ID     string            `json:"response_id" example:"abc-123"`
	Desc   string            `json:"response_desc" example:"success"`
	Data   JWTDecodeResponse `json:"response_data"`
}

// JWTValidateRequest is the request body for the JWT validate endpoint.
type JWTValidateRequest struct {
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"`
	Secret string `json:"secret" example:"my-secret-key"`
}

// JWTValidateResponse is the response for the JWT validate endpoint.
type JWTValidateResponse struct {
	Valid  bool                   `json:"valid"`
	Claims map[string]interface{} `json:"claims,omitempty"`
}

// JWTValidateResponseSwag is used for Swagger documentation only.
type JWTValidateResponseSwag struct {
	Code   string              `json:"response_code" example:"200"`
	Refnum string              `json:"response_refnum" example:""`
	ID     string              `json:"response_id" example:"abc-123"`
	Desc   string              `json:"response_desc" example:"success"`
	Data   JWTValidateResponse `json:"response_data"`
}
