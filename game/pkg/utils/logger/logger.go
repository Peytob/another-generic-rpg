package logger

import (
	"context"
	"log/slog"
)

type loggerKeyT struct{}

func ToCtx(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKeyT{}, l)
}

func FromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKeyT{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}
