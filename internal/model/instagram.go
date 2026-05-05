package model

// InstagramDownloadRequest is the request body for the Instagram download endpoint.
type InstagramDownloadRequest struct {
	URL string `json:"url" example:"https://www.instagram.com/arditya.kekasi/p/C2uaPTYShvn/"`
}

// InstagramMediaItem represents a single downloadable media item.
type InstagramMediaItem struct {
	MediaType string `json:"media_type" example:"photo"`
	MediaURL  string `json:"media_url"  example:"https://scontent.cdninstagram.com/..."`
	ThumbURL  string `json:"thumb_url,omitempty" example:"https://scontent.cdninstagram.com/..."`
}

// InstagramDownloadResponse is the response for the Instagram download endpoint.
type InstagramDownloadResponse struct {
	PostType string               `json:"post_type" example:"reel"`
	Items    []InstagramMediaItem `json:"items"`
}

// InstagramDownloadResponseSwag is used for Swagger documentation only.
type InstagramDownloadResponseSwag struct {
	Code   string                    `json:"response_code" example:"200"`
	Refnum string                    `json:"response_refnum" example:""`
	ID     string                    `json:"response_id" example:"abc-123"`
	Desc   string                    `json:"response_desc" example:"success"`
	Data   InstagramDownloadResponse `json:"response_data"`
}
