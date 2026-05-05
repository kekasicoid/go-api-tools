// internal/domain/instagram.go
package domain

// InstagramMedia represents a single downloadable media item.
type InstagramMedia struct {
	MediaType string // "photo" or "video"
	MediaURL  string
	ThumbURL  string
}

// InstagramMediaInfo holds all extracted media from an Instagram URL.
type InstagramMediaInfo struct {
	PostType string // "post", "reel", "tv", "story"
	Items    []InstagramMedia
}

// InstagramDownloader defines the contract for extracting media from Instagram.
type InstagramDownloader interface {
	Download(url string) (InstagramMediaInfo, error)
}
