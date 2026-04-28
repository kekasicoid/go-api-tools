package model

type FormatJsonRequest struct {
	Data string `json:"data,omitempty" example:"{\"name\":\"kekasi\",\"age\":25}"`
}

type FormatJsonResponse struct {
	Formatted string `json:"formatted,omitempty" example:"{\n  \"name\": \"kekasi\",\n  \"age\": 25\n}"`
}

type FormatJsonResponseSwag struct {
	Code   string             `json:"response_code" example:"200"`
	Refnum string             `json:"response_refnum" example:""`
	ID     string             `json:"response_id" example:"abc-123"`
	Desc   string             `json:"response_desc" example:"success"`
	Data   FormatJsonResponse `json:"response_data"`
}
