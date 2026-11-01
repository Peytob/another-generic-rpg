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

func WithLogValue(ctx context.Context, value slog.Attr) context.Context {
	l := FromCtx(ctx)
	l = l.With(value)
	return ToCtx(ctx, l)
}

func Uint32(key string, id uint32) slog.Attr {
	return slog.Uint64(key, uint64(id))
}

func Error(err error) slog.Attr {
	return slog.Any("err", err)
}
