// internal/middleware/ratelimit.go
package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				return 0
			}
			return db
		}(),
	})
}

func RateLimit() gin.HandlerFunc {
	if rdb == nil {
		initRedis()
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "ratelimit:" + ip

		// Increment the request count for the IP
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// Set expiration for the key if it's new
		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if count > 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
