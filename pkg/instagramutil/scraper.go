// pkg/instagramutil/scraper.go
package instagramutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kekasicoid/go-api-tools/internal/domain"
)

// browserUserAgent is the User-Agent header sent with every request to Instagram.
const browserUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

// balancedJSON extracts the first balanced JSON object that starts at or after
// the given offset in s. It properly skips string literals so embedded braces
// inside strings are not counted.
func balancedJSON(s string, offset int) string {
	start := strings.Index(s[offset:], "{")
	if start == -1 {
		return ""
	}
	start += offset

	depth := 0
	inString := false
	escape := false

	for i := start; i < len(s); i++ {
		c := s[i]
		if escape {
			escape = false
			continue
		}
		if c == '\\' && inString {
			escape = true
			continue
		}
		if c == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		switch c {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

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

	// Strategy 5 (last resort): try Instagram's lightweight JSON endpoint.
	// Strategies 1-4 are attempted inside extractFromHTML above.
	if len(info.Items) == 0 {
		if apiInfo, ok := s.tryAPIEndpoint(shortcode, postType); ok && len(apiInfo.Items) > 0 {
			info = apiInfo
		}
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
	// Normalize: strip trailing slash and query string for matching.
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

	req.Header.Set("User-Agent", browserUserAgent)
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

	const maxResponseBodyBytes = 10 * 1024 * 1024 // 10 MB safety cap
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}

// tryAPIEndpoint attempts to fetch media info via Instagram's lightweight JSON
// endpoint (?__a=1&__d=dis). This works for some public posts without a session
// cookie, and is used as a last resort when HTML parsing yields nothing.
// It returns (info, true) on success, or (zero, false) on any failure.
func (s *Scraper) tryAPIEndpoint(shortcode, postType string) (domain.InstagramMediaInfo, bool) {
	apiURL := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1&__d=dis", shortcode)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return domain.InstagramMediaInfo{}, false
	}

	req.Header.Set("User-Agent", browserUserAgent)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://www.instagram.com/")

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.InstagramMediaInfo{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.InstagramMediaInfo{}, false
	}

	const maxBytes = 5 * 1024 * 1024
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	if err != nil {
		return domain.InstagramMediaInfo{}, false
	}

	var root interface{}
	if err := json.Unmarshal(bodyBytes, &root); err != nil {
		return domain.InstagramMediaInfo{}, false
	}

	node := findNodeByShortcodeOrCode(root, shortcode)
	if node == nil {
		// Some responses wrap the media under "items":[{...}].
		if rootMap, ok := root.(map[string]interface{}); ok {
			if items, ok := rootMap["items"].([]interface{}); ok && len(items) > 0 {
				if first, ok := items[0].(map[string]interface{}); ok {
					node = first
				}
			}
		}
	}

	if node == nil {
		return domain.InstagramMediaInfo{}, false
	}

	info := extractMediaFromNode(node, postType)
	if len(info.Items) == 0 {
		return domain.InstagramMediaInfo{}, false
	}
	return info, true
}

// ---- HTML parsing helpers ---------------------------------------------------

// rehydrationRe matches the __UNIVERSAL_DATA_FOR_REHYDRATION__ script tag (may span lines).
var rehydrationRe = regexp.MustCompile(`(?s)<script[^>]+id="__UNIVERSAL_DATA_FOR_REHYDRATION__"[^>]*>(.*?)</script>`)

// sharedDataRe matches the legacy window._sharedData assignment.
var sharedDataRe = regexp.MustCompile(`(?s)window\._sharedData\s*=\s*(\{.*?\});\s*</script>`)

// additionalDataMarker is the JS function name used by the window.__additionalDataLoaded strategy.
const additionalDataMarker = "window.__additionalDataLoaded("

// ogPropContentRe matches <meta property="..." content="..."> (property before content).
var ogPropContentRe = regexp.MustCompile(`<meta\b[^>]*\bproperty="([^"]+)"[^>]*\bcontent="([^"]+)"`)

// ogContentPropRe matches <meta content="..." property="..."> (content before property).
var ogContentPropRe = regexp.MustCompile(`<meta\b[^>]*\bcontent="([^"]+)"[^>]*\bproperty="([^"]+)"`)

// extractFromHTML tries each known JSON embedding strategy to find media data.
func extractFromHTML(html, shortcode, postType string) (domain.InstagramMediaInfo, error) {
	// Strategy 1: __UNIVERSAL_DATA_FOR_REHYDRATION__
	if matches := rehydrationRe.FindStringSubmatch(html); len(matches) >= 2 {
		var root interface{}
		if err := json.Unmarshal([]byte(matches[1]), &root); err == nil {
			if node := findNodeByShortcodeOrCode(root, shortcode); node != nil {
				return extractMediaFromNode(node, postType), nil
			}
		}
	}

	// Strategy 2: window._sharedData
	if matches := sharedDataRe.FindStringSubmatch(html); len(matches) >= 2 {
		var root interface{}
		if err := json.Unmarshal([]byte(matches[1]), &root); err == nil {
			if node := findNodeByShortcodeOrCode(root, shortcode); node != nil {
				return extractMediaFromNode(node, postType), nil
			}
		}
	}

	// Strategy 3: window.__additionalDataLoaded
	if idx := strings.Index(html, additionalDataMarker); idx != -1 {
		jsonStr := balancedJSON(html, idx+len(additionalDataMarker))
		if jsonStr != "" {
			var root interface{}
			if err := json.Unmarshal([]byte(jsonStr), &root); err == nil {
				if node := findNodeByShortcodeOrCode(root, shortcode); node != nil {
					return extractMediaFromNode(node, postType), nil
				}
			}
		}
	}

	// Strategy 4: Open Graph meta tags — reliable for public single-photo / video posts
	// because Instagram populates OG tags for link-preview purposes.
	if info := extractFromOGMeta(html, postType); len(info.Items) > 0 {
		return info, nil
	}

	return domain.InstagramMediaInfo{PostType: postType}, nil
}

// findNodeByShortcodeOrCode recursively searches for a JSON object whose
// "shortcode" or "code" field (both used across Instagram API versions) equals
// the given shortcode.
func findNodeByShortcodeOrCode(v interface{}, shortcode string) map[string]interface{} {
	switch node := v.(type) {
	case map[string]interface{}:
		if sc, ok := node["shortcode"].(string); ok && sc == shortcode {
			return node
		}
		if sc, ok := node["code"].(string); ok && sc == shortcode {
			return node
		}
		for _, val := range node {
			if found := findNodeByShortcodeOrCode(val, shortcode); found != nil {
				return found
			}
		}
	case []interface{}:
		for _, item := range node {
			if found := findNodeByShortcodeOrCode(item, shortcode); found != nil {
				return found
			}
		}
	}
	return nil
}

// extractFromOGMeta extracts media info from Open Graph meta tags.
// Instagram always populates these for public posts so link-preview services work.
// This is a reliable fallback for single photo and video posts, but carousel
// posts only have their first item in OG tags.
func extractFromOGMeta(pageHTML, postType string) domain.InstagramMediaInfo {
	info := domain.InstagramMediaInfo{PostType: postType}

	ogProps := make(map[string]string) // property -> content

	for _, m := range ogPropContentRe.FindAllStringSubmatch(pageHTML, -1) {
		if len(m) == 3 {
			ogProps[m[1]] = html.UnescapeString(m[2])
		}
	}
	for _, m := range ogContentPropRe.FindAllStringSubmatch(pageHTML, -1) {
		if len(m) == 3 {
			ogProps[m[2]] = html.UnescapeString(m[1])
		}
	}

	videoURL := ogProps["og:video:secure_url"]
	if videoURL == "" {
		videoURL = ogProps["og:video"]
	}
	imageURL := ogProps["og:image"]

	if videoURL != "" {
		info.Items = append(info.Items, domain.InstagramMedia{
			MediaType: "video",
			MediaURL:  videoURL,
			ThumbURL:  imageURL,
		})
	} else if imageURL != "" {
		info.Items = append(info.Items, domain.InstagramMedia{
			MediaType: "photo",
			MediaURL:  imageURL,
			ThumbURL:  imageURL,
		})
	}

	return info
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
