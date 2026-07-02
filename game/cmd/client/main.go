package main

import (
	"context"
	"game/internal/config"
	"game/internal/engine/client"
	"game/internal/engine/fsm"
	"game/internal/engine/graphic"
	"game/pkg/engine"
	"game/pkg/utils/logger"
	"game/pkg/window"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
		panic("failed to initialize logger: " + err.Error())
	}

	ctx = logger.ToCtx(ctx, l)

	l.LogAttrs(ctx, slog.LevelInfo, "initializing client")

	/* Compability init */

	err = glfw.Init()
	if err != nil {
		panic("failed to initialize GLFW: " + err.Error())
	}

	err = gl.Init()
	if err != nil {
		panic("failed to initialize OpenGl: " + err.Error())
	}

	/* Modules */

	w, err := window.Init(window.Opts{
		Width:  800,
		Height: 600,
		Title:  "another-rpg",
	})
	if err != nil {
		panic("failed to initialize window module: " + err.Error())
	}

	g, err := graphic.InitializeGraphic(ctx)
	if err != nil {
		panic("failed to initialize graphic module: " + err.Error())
	}

	/* Client */

	machine := engine.NewMachineBuilder().
		RegisterState(fsm.NewInitialState(), make(engine.Transitions)).
		InitialState(fsm.InitialStateIdentifier).
		RegisterState(fsm.NewStoppedState(), make(engine.Transitions)).
		GlobalTransitions(engine.Transitions{
			fsm.StoppedEvent: fsm.StoppedStateIdentifier,
		}).
		FinalStates(fsm.StoppedStateIdentifier).
		MustBuild()

	cl := client.NewBuilder().
		Fsm(machine).
		Window(w).
		Graphic(g).
		MustBuild()

	l.LogAttrs(ctx, slog.LevelInfo, "starting engine running")
	if err = cl.Run(ctx); err != nil {
		panic("failed engine cycle: " + err.Error())
	}

	if err = cl.Shutdown(ctx); err != nil {
		panic("failed to stop client gracefully: " + err.Error())
	}
}
