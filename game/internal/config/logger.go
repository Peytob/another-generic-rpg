package config

import (
	"context"
	"log/slog"
)

type LoggerOpts struct {
	Enabled bool
	Handler slog.Handler
}

func ClientLogger(_ context.Context, opts LoggerOpts) (*slog.Logger, error) {
	if !opts.Enabled {
		return slog.New(slog.DiscardHandler), nil
	}

	return slog.New(opts.Handler), nil
}
