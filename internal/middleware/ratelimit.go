// internal/middleware/ratelimit.go
package middleware

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kekasicoid/go-api-tools/internal/model"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func initRedis() {
	addr := os.Getenv("REDIS_ADDR")
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				return 0
			}
			return db
		}(),
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis at %s: %v", addr, err)
	} else {
		log.Printf("Successfully connected to Redis at %s", addr)
	}
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
			model.RespInternalServerError(c, "internal server error")
			c.Abort()
			return
		}

		// Set expiration for the key if it's new
		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if count > 100 {
			model.RespTooManyRequests(c, "too many requests")
			c.Abort()
			return
		}

		c.Next()

	}
}
