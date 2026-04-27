// cmd/server/main.go
package main

import (
	"os"

	httpDelivery "github.com/kekasicoid/go-api-tools/internal/delivery/http"
	"github.com/kekasicoid/go-api-tools/internal/middleware"
	"github.com/kekasicoid/go-api-tools/internal/usecase"
	"github.com/kekasicoid/go-api-tools/pkg/jsonutil"
)

func main() {
	// init dependency
	formatter := jsonutil.NewJSONFormatter()
	usecase := usecase.NewFormatterUsecase(formatter)
	handler := httpDelivery.NewHandler(usecase)

	// router
	r := httpDelivery.NewRouter(handler)

	// middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())

	// run
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)

}
