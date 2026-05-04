// cmd/server/main.go
package main

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	docs "github.com/kekasicoid/go-api-tools/docs"
	httpDelivery "github.com/kekasicoid/go-api-tools/internal/delivery/http"
	"github.com/kekasicoid/go-api-tools/internal/middleware"
	appusecase "github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/jsonutil"
	"github.com/kekasicoid/go-api-tools/pkg/jwtutil"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
	"go.uber.org/zap"
)

// @title           Go API Tools by Arditya Kekasi
// @version         1.0
// @description     A collection of API tools in Go.

// @contact.name   Official support
// @contact.url    https://kekasi.co.id/en/contact/
// @contact.email  arditya@kekasi.co.id

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:7070
// @BasePath  /

func main() {
	// Try common locations so env works when running from root or cmd/server.
	_ = godotenv.Load(".env", "../../.env")

	// init logger
	logger.Init()
	defer logger.Sync()

	logger.Log.Info("starting server")

	// init dependency
	formatter := jsonutil.NewJSONFormatter()
	rdb := initRedisClient()
	if rdb != nil {
		defer rdb.Close()
	}
	formatterUsecase := appusecase.NewFormatterUsecase(formatter, rdb)
	handler := httpDelivery.NewHandler(formatterUsecase)

	jwtDecoder := jwtutil.NewJWTDecoder()
	jwtUsecase := appusecase.NewJWTUsecase(jwtDecoder)
	jwtHandler := httpDelivery.NewJWTHandler(jwtUsecase)

	// router
	r := httpDelivery.NewRouter()

	// middleware
	r.Use(middleware.CORS())
	r.Use(middleware.ValidateRequestID())
	r.Use(middleware.RateLimit())
	// r.Use(middleware.RequestLogger())

	// routes
	httpDelivery.RegisterRoutes(r, handler, jwtHandler)

	// run
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	// Allow Swagger host to be configured dynamically at runtime.
	swaggerHost := strings.TrimSpace(os.Getenv("SWAGGER_HOST"))
	if swaggerHost == "" {
		swaggerHost = "localhost:" + port
	}

	if strings.Contains(swaggerHost, "://") {
		if parsedURL, err := url.Parse(swaggerHost); err == nil {
			if parsedURL.Host != "" {
				docs.SwaggerInfo.Host = parsedURL.Host
			} else {
				docs.SwaggerInfo.Host = swaggerHost
			}
			if parsedURL.Scheme != "" {
				docs.SwaggerInfo.Schemes = []string{parsedURL.Scheme}
			}
		} else {
			docs.SwaggerInfo.Host = swaggerHost
		}
	} else {
		docs.SwaggerInfo.Host = swaggerHost
	}
	docs.SwaggerInfo.BasePath = "/"

	r.Run(":" + port)

}

func initRedisClient() *redis.Client {
	addr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	if addr == "" {
		logger.Log.Warn("REDIS_ADDR not set, redis cache disabled")
		return nil
	}

	db := 0
	if rawDB := strings.TrimSpace(os.Getenv("REDIS_DB")); rawDB != "" {
		parsedDB, err := strconv.Atoi(rawDB)
		if err != nil {
			logger.Log.Warn("invalid REDIS_DB, using default 0", zap.Error(err), zap.String("redis_db", rawDB))
		} else {
			db = parsedDB
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Log.Warn("redis ping failed, cache disabled", zap.Error(err))
		return nil
	}

	logger.Log.Info("redis connected for cache")
	return rdb
}
