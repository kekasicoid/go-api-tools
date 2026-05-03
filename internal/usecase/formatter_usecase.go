// internal/usecase/formatter_usecase.go
package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kekasicoid/go-api-tools/internal/domain"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

const defaultCacheTTL = 1 * time.Hour
const cacheTTLEnvName = "JSON_FORMATTER_TTL_HOURS"

type FormatterUsecase struct {
	formatter   domain.Formatter
	redisClient *redis.Client
	cacheTTL    time.Duration
}

func NewFormatterUsecase(f domain.Formatter, rdb *redis.Client) *FormatterUsecase {
	cacheTTL := defaultCacheTTL
	if rawTTL := strings.TrimSpace(os.Getenv(cacheTTLEnvName)); rawTTL != "" {
		ttlHours, err := strconv.Atoi(rawTTL)
		if err != nil {
			logger.Log.Warn("invalid JSON_FORMATTER_TTL_HOURS, using default 1 hour", zap.String("json_formatter_ttl_hours", rawTTL), zap.Error(err))
		} else if ttlHours <= 0 {
			logger.Log.Warn("JSON_FORMATTER_TTL_HOURS must be > 0, using default 1 hour", zap.Int("json_formatter_ttl_hours", ttlHours))
		} else {
			cacheTTL = time.Duration(ttlHours) * time.Hour
		}
	}

	return &FormatterUsecase{
		formatter:   f,
		redisClient: rdb,
		cacheTTL:    cacheTTL,
	}
}

func (u *FormatterUsecase) FormatJSON(input string) (string, error) {
	if u.redisClient != nil {
		cacheKey := buildFormatJSONCacheKey(input)

		cachedValue, err := u.redisClient.Get(context.Background(), cacheKey).Result()
		if err == nil {
			return cachedValue, nil
		}

		if err != redis.Nil {
			logger.Log.Warn("failed to read cached formatted json", zap.Error(err))
		}

		result, formatErr := u.formatter.Format(input)
		if formatErr != nil {
			return "", formatErr
		}

		if err := u.redisClient.Set(context.Background(), cacheKey, result, u.cacheTTL).Err(); err != nil {
			logger.Log.Warn("failed to write formatted json cache", zap.Error(err))
		}

		return result, nil
	}

	return u.formatter.Format(input)
}

func buildFormatJSONCacheKey(input string) string {
	h := sha256.Sum256([]byte(input))
	return fmt.Sprintf("format_json:%s", hex.EncodeToString(h[:]))
}
