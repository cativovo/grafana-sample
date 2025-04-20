package ctxvalue

import (
	"context"
	"log/slog"
)

const ctxKeyLogger ctxKey = "logger"

func Logger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxKeyLogger).(*slog.Logger)
	if !ok {
		panic("ctx doesn't have logger")
	}

	return logger
}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, l)
}
