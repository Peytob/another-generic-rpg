package logger

import (
	"context"
	"log/slog"
)

type loggerKeyT struct{}

var loggerKey = loggerKeyT{}

func ToCtx(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}
