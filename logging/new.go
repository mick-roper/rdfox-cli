package logging

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		parseLevel(level),
	)

	return zap.New(core)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug", "trace":
		return zap.DebugLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
