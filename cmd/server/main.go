// cmd/server/main.go
package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	docs "github.com/kekasicoid/go-api-tools/docs"
	httpDelivery "github.com/kekasicoid/go-api-tools/internal/delivery/http"
	"github.com/kekasicoid/go-api-tools/internal/middleware"
	"github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/jsonutil"
	"github.com/kekasicoid/go-api-tools/pkg/logger"
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
	usecase := usecase.NewFormatterUsecase(formatter)
	handler := httpDelivery.NewHandler(usecase)

	// router
	r := httpDelivery.NewRouter()

	// middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ValidateRequestID())

	// routes
	httpDelivery.RegisterRoutes(r, handler)

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
