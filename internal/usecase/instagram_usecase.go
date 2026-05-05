// internal/usecase/instagram_usecase.go
package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kekasicoid/go-api-tools/internal/domain"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"

	"time"
)

const igDefaultCacheTTL = 30 * time.Minute

// InstagramUsecase handles Instagram media download business logic.
type InstagramUsecase struct {
	downloader  domain.InstagramDownloader
	redisClient *redis.Client
	cacheTTL    time.Duration
}

// NewInstagramUsecase creates a new InstagramUsecase.
func NewInstagramUsecase(d domain.InstagramDownloader, rdb *redis.Client) *InstagramUsecase {
	return &InstagramUsecase{
		downloader:  d,
		redisClient: rdb,
		cacheTTL:    igDefaultCacheTTL,
	}
}

// Download fetches media info for the given Instagram URL, using Redis cache when available.
func (u *InstagramUsecase) Download(url string) (domain.InstagramMediaInfo, error) {
	if u.redisClient != nil {
		cacheKey := igCacheKey(url)

		cached, err := u.redisClient.Get(context.Background(), cacheKey).Bytes()
		if err == nil {
			var info domain.InstagramMediaInfo
			if jsonErr := json.Unmarshal(cached, &info); jsonErr == nil {
				return info, nil
			}
		} else if err != redis.Nil {
			logger.Log.Warn("failed to read instagram cache", zap.Error(err))
		}
	}

	info, err := u.downloader.Download(url)
	if err != nil {
		return domain.InstagramMediaInfo{}, err
	}

	if u.redisClient != nil {
		cacheKey := igCacheKey(url)
		if data, jsonErr := json.Marshal(info); jsonErr == nil {
			if setErr := u.redisClient.Set(context.Background(), cacheKey, data, u.cacheTTL).Err(); setErr != nil {
				logger.Log.Warn("failed to write instagram cache", zap.Error(setErr))
			}
		}
	}

	return info, nil
}

func igCacheKey(url string) string {
	h := sha256.Sum256([]byte(url))
	return fmt.Sprintf("instagram_download:%s", hex.EncodeToString(h[:]))
}
