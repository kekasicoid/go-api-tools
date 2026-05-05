// pkg/instagramutil/scraper.go
package instagramutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kekasicoid/go-api-tools/internal/domain"
)

// Scraper implements domain.InstagramDownloader via HTML scraping.
type Scraper struct {
	client *http.Client
}

// NewScraper creates a new Instagram Scraper.
func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Download fetches the Instagram page and extracts downloadable media URLs.
// It supports posts (/p/), reels (/reel/), and IGTV (/tv/) URLs.
// Stories require a logged-in session and are not supported without cookies.
func (s *Scraper) Download(rawURL string) (domain.InstagramMediaInfo, error) {
	postType, shortcode, err := parseInstagramURL(rawURL)
	if err != nil {
		return domain.InstagramMediaInfo{}, err
	}

	pageURL := buildPageURL(postType, shortcode)

	body, err := s.fetchPage(pageURL)
	if err != nil {
		return domain.InstagramMediaInfo{}, err
	}

	info, err := extractFromHTML(body, shortcode, postType)
	if err != nil {
		return domain.InstagramMediaInfo{}, err
	}

	if len(info.Items) == 0 {
		return domain.InstagramMediaInfo{}, errors.New("no downloadable media found; the post may be private or Instagram's page structure has changed")
	}

	return info, nil
}

// parseInstagramURL returns (postType, shortcode, error).
// Supported paths:
//   - /p/{shortcode}
//   - /reel/{shortcode}
//   - /tv/{shortcode}
//   - /stories/{username}/{mediaID}
func parseInstagramURL(rawURL string) (postType, shortcode string, err error) {
	// Normalise: strip trailing slash and query string for matching.
	u := strings.Split(strings.TrimRight(rawURL, "/"), "?")[0]
	parts := strings.Split(u, "/")

	for i, p := range parts {
		switch p {
		case "p":
			if i+1 < len(parts) && parts[i+1] != "" {
				return "post", parts[i+1], nil
			}
		case "reel":
			if i+1 < len(parts) && parts[i+1] != "" {
				return "reel", parts[i+1], nil
			}
		case "tv":
			if i+1 < len(parts) && parts[i+1] != "" {
				return "tv", parts[i+1], nil
			}
		case "stories":
			// /stories/{username}/{mediaID}
			if i+2 < len(parts) && parts[i+2] != "" {
				return "story", parts[i+2], nil
			}
		}
	}

	return "", "", errors.New("unsupported Instagram URL; expected /p/, /reel/, /tv/, or /stories/ path")
}

func buildPageURL(postType, shortcode string) string {
	switch postType {
	case "reel":
		return "https://www.instagram.com/reel/" + shortcode + "/"
	case "tv":
		return "https://www.instagram.com/tv/" + shortcode + "/"
	default:
		return "https://www.instagram.com/p/" + shortcode + "/"
	}
}

// fetchPage fetches the Instagram page HTML with browser-like headers.
func (s *Scraper) fetchPage(pageURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Cache-Control", "max-age=0")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Instagram page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Instagram returned HTTP %d", resp.StatusCode)
	}

	const maxBody = 10 * 1024 * 1024 // 10 MB safety cap
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}

// ---- HTML parsing helpers ---------------------------------------------------

// rehydrationRe matches the __UNIVERSAL_DATA_FOR_REHYDRATION__ script tag (may span lines).
var rehydrationRe = regexp.MustCompile(`(?s)<script[^>]+id="__UNIVERSAL_DATA_FOR_REHYDRATION__"[^>]*>(.*?)</script>`)

// sharedDataRe matches the legacy window._sharedData assignment.
var sharedDataRe = regexp.MustCompile(`(?s)window\._sharedData\s*=\s*(\{.*?\});\s*</script>`)

// extractFromHTML tries each known JSON embedding strategy to find media data.
func extractFromHTML(html, shortcode, postType string) (domain.InstagramMediaInfo, error) {
	// Strategy 1: __UNIVERSAL_DATA_FOR_REHYDRATION__
	if matches := rehydrationRe.FindStringSubmatch(html); len(matches) >= 2 {
		var root interface{}
		if err := json.Unmarshal([]byte(matches[1]), &root); err == nil {
			if node := findNodeByShortcode(root, shortcode); node != nil {
				return extractMediaFromNode(node, postType), nil
			}
		}
	}

	// Strategy 2: window._sharedData
	if matches := sharedDataRe.FindStringSubmatch(html); len(matches) >= 2 {
		var root interface{}
		if err := json.Unmarshal([]byte(matches[1]), &root); err == nil {
			if node := findNodeByShortcode(root, shortcode); node != nil {
				return extractMediaFromNode(node, postType), nil
			}
		}
	}

	return domain.InstagramMediaInfo{PostType: postType}, nil
}

// findNodeByShortcode recursively searches for a JSON object whose "shortcode"
// field equals the given shortcode.
func findNodeByShortcode(v interface{}, shortcode string) map[string]interface{} {
	switch node := v.(type) {
	case map[string]interface{}:
		if sc, ok := node["shortcode"].(string); ok && sc == shortcode {
			return node
		}
		for _, val := range node {
			if found := findNodeByShortcode(val, shortcode); found != nil {
				return found
			}
		}
	case []interface{}:
		for _, item := range node {
			if found := findNodeByShortcode(item, shortcode); found != nil {
				return found
			}
		}
	}
	return nil
}

// extractMediaFromNode extracts InstagramMediaInfo from a parsed media JSON node.
// It handles both the newer API format (image_versions2 / video_versions /
// carousel_media) and the older GraphQL format (display_url / video_url /
// edge_sidecar_to_children).
func extractMediaFromNode(node map[string]interface{}, postType string) domain.InstagramMediaInfo {
	info := domain.InstagramMediaInfo{PostType: postType}

	mediaTypeNum, _ := node["media_type"].(float64)

	switch int(mediaTypeNum) {
	case 1: // photo
		info.Items = append(info.Items, singlePhotoItem(node))

	case 2: // video / reel
		info.Items = append(info.Items, singleVideoItem(node))

	case 8: // carousel / album — new API
		if items := carouselItemsNew(node); len(items) > 0 {
			info.Items = items
			return info
		}
		// fall through to old GraphQL carousel
		info.Items = carouselItemsOld(node)

	default:
		// Old GraphQL API: is_video bool
		if isVideo, ok := node["is_video"].(bool); ok {
			if isVideo {
				info.Items = append(info.Items, singleVideoItem(node))
			} else {
				info.Items = append(info.Items, singlePhotoItem(node))
			}
			return info
		}

		// Old GraphQL carousel fallback
		if items := carouselItemsOld(node); len(items) > 0 {
			info.Items = items
		}
	}

	return info
}

// singlePhotoItem builds an InstagramMedia for a photo node.
func singlePhotoItem(node map[string]interface{}) domain.InstagramMedia {
	item := domain.InstagramMedia{MediaType: "photo"}
	// New API
	item.MediaURL = firstImageCandidate(node)
	// Old GraphQL fallback
	if item.MediaURL == "" {
		item.MediaURL, _ = node["display_url"].(string)
	}
	item.ThumbURL = thumbURL(node)
	return item
}

// singleVideoItem builds an InstagramMedia for a video node.
func singleVideoItem(node map[string]interface{}) domain.InstagramMedia {
	item := domain.InstagramMedia{MediaType: "video"}
	// New API
	item.MediaURL = firstVideoCandidate(node)
	// Old GraphQL fallback
	if item.MediaURL == "" {
		item.MediaURL, _ = node["video_url"].(string)
	}
	item.ThumbURL = thumbURL(node)
	return item
}

// carouselItemsNew extracts items from the newer "carousel_media" array.
func carouselItemsNew(node map[string]interface{}) []domain.InstagramMedia {
	carouselRaw, ok := node["carousel_media"].([]interface{})
	if !ok {
		return nil
	}

	var items []domain.InstagramMedia
	for _, raw := range carouselRaw {
		child, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		mt, _ := child["media_type"].(float64)
		switch int(mt) {
		case 1:
			items = append(items, singlePhotoItem(child))
		case 2:
			items = append(items, singleVideoItem(child))
		}
	}
	return items
}

// carouselItemsOld extracts items from the older "edge_sidecar_to_children" structure.
func carouselItemsOld(node map[string]interface{}) []domain.InstagramMedia {
	sidecar, ok := node["edge_sidecar_to_children"].(map[string]interface{})
	if !ok {
		return nil
	}
	edges, ok := sidecar["edges"].([]interface{})
	if !ok {
		return nil
	}

	var items []domain.InstagramMedia
	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		child, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		if isVideo, ok := child["is_video"].(bool); ok && isVideo {
			items = append(items, singleVideoItem(child))
		} else {
			items = append(items, singlePhotoItem(child))
		}
	}
	return items
}

// firstImageCandidate returns the highest-quality image URL from image_versions2.
func firstImageCandidate(node map[string]interface{}) string {
	iv2, ok := node["image_versions2"].(map[string]interface{})
	if !ok {
		return ""
	}
	candidates, ok := iv2["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return ""
	}
	if c, ok := candidates[0].(map[string]interface{}); ok {
		url, _ := c["url"].(string)
		return url
	}
	return ""
}

// firstVideoCandidate returns the highest-quality video URL from video_versions.
func firstVideoCandidate(node map[string]interface{}) string {
	versions, ok := node["video_versions"].([]interface{})
	if !ok || len(versions) == 0 {
		return ""
	}
	if v, ok := versions[0].(map[string]interface{}); ok {
		url, _ := v["url"].(string)
		return url
	}
	return ""
}

// thumbURL returns a thumbnail URL for the node, preferring display_url.
func thumbURL(node map[string]interface{}) string {
	if v, ok := node["display_url"].(string); ok && v != "" {
		return v
	}
	if v, ok := node["thumbnail_url"].(string); ok && v != "" {
		return v
	}
	return firstImageCandidate(node)
}
