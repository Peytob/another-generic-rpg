package di

import (
	"context"
	"io"
	"log/slog"
)

type LoggerOpts struct {
	Enabled bool
	Handler slog.Handler
}

func ClientLogger(_ context.Context, opts LoggerOpts) (*slog.Logger, error) {
	if !opts.Enabled {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	}

	return slog.New(opts.Handler), nil
}
