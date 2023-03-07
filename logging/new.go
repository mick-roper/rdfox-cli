package logging

import "go.uber.org/zap"

func New() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.Level.SetLevel(zap.InfoLevel)

	return zap.Must(cfg.Build())
}
