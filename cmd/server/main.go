// cmd/server/main.go
package main

import (
	"os"

	"github.com/joho/godotenv"
	httpDelivery "github.com/kekasicoid/go-api-tools/internal/delivery/http"
	"github.com/kekasicoid/go-api-tools/internal/middleware"
	"github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/jsonutil"
	"github.com/kekasicoid/go-api-tools/pkg/logger"

	_ "github.com/kekasicoid/go-api-tools/docs"
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
	r := httpDelivery.NewRouter(handler)

	// middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(middleware.RequestLogger())

	// run
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)

}
