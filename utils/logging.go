package utils

import (
	"context"

	"go.uber.org/zap"
)

var loggerCtxKey = struct{}{}

func AddLoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return addToContext(ctx, loggerCtxKey, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := getFromContext(ctx, loggerCtxKey).(*zap.Logger); ok {
		return logger
	}

	return zap.NewNop()
}
