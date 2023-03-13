package logging

import (
	"context"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

var loggerCtxKey = struct{}{}

func AddToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return utils.AddToContext(ctx, loggerCtxKey, logger)
}

func GetFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := utils.GetFromContext(ctx, loggerCtxKey).(*zap.Logger); ok {
		return logger
	}

	return zap.NewNop()
}
