package model

// SwaggRespSuccess adalah struct khusus untuk dokumentasi Swagger response 200.
type SwaggRespSuccess struct {
	Code   string      `json:"response_code" example:"200"`
	Refnum string      `json:"response_refnum" example:""`
	ID     string      `json:"response_id" example:"abc-123"`
	Desc   string      `json:"response_desc" example:"success"`
	Data   interface{} `json:"response_data"`
}

// SwaggRespError adalah struct khusus untuk dokumentasi Swagger response 400.
type SwaggRespError struct {
	Code   string `json:"response_code" example:"400"`
	Refnum string `json:"response_refnum" example:""`
	ID     string `json:"response_id" example:"abc-123"`
	Desc   string `json:"response_desc" example:"data is required"`
	Data   any    `json:"response_data"`
}
