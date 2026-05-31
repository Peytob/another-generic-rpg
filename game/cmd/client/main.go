package main

import (
	"context"
	"game/internal/config"
	"game/internal/engine/client"
	"game/internal/engine/fsm"
	"game/pkg/engine"
	"game/pkg/utils/logger"
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

	l, err := config.ClientLogger(ctx, config.LoggerOpts{
		Enabled: cfg.Log.Enabled,
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Log.Level,
		}),
	})
	if err != nil {
		panic("failed to initialize logger:" + err.Error())
	}

	ctx = logger.ToCtx(ctx, l)

	l.LogAttrs(ctx, slog.LevelInfo, "initializing client")

	machine := engine.NewMachineBuilder().
		RegisterState(fsm.NewInitialState(), make(engine.Transitions)).
		InitialState(fsm.InitialStateIdentifier).
		RegisterState(fsm.NewStoppedState(), make(engine.Transitions)).
		GlobalTransitions(engine.Transitions{
			fsm.StoppedEvent: fsm.StoppedStateIdentifier,
		}).
		FinalStates([]engine.StateIdentifier{
			fsm.StoppedStateIdentifier,
		}).
		MustBuild()

	cl := client.NewBuilder().
		Fsm(machine).
		MustBuild()

	l.LogAttrs(ctx, slog.LevelInfo, "starting engine running")
	if err = cl.Run(ctx); err != nil {
		panic("failed engine cycle: " + err.Error())
	}
}
