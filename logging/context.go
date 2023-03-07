package logging

import (
	"context"

	"go.uber.org/zap"
)

var loggerCtxKey = struct{}{}

func AddToContext(ctx context.Context, logger *zap.Logger) context.Context {
	if ctx == nil || logger == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerCtxKey, logger)
}

func GetFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerCtxKey).(*zap.Logger); ok {
		return logger
	}

	return zap.NewNop()
}
