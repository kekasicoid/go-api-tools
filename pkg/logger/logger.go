package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	var err error
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout", "dev-out.log"}
	cfg.ErrorOutputPaths = []string{"stderr", "dev-err.log"}

	if os.Getenv("APP_ENV") == "production" {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout", "prod-out.log"}
		cfg.ErrorOutputPaths = []string{"stderr", "prod-err.log"}
	}

	Log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func Sync() {
	_ = Log.Sync()
}
