package main

import (
	"context"
	"game/cmd/di"
	"game/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func init() {
	// Client main method should be executed in main application thread due
	// to C libraries restrictions
	runtime.LockOSThread()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfiguration[config.ClientConfiguration]("./config/client.yaml")
	if err != nil {
		panic(err)
	}

	_, err = di.LoadClientLogger(ctx, di.LoggerOpts{
		Enabled: cfg.Log.Enabled,
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Log.Level,
		}),
	})
	if err != nil {
		panic("failed to initialize logger:" + err.Error())
	}
}
