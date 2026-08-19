package main

import (
	"context"
	"engine/utils/logger"
	"engine/window"
	"game/internal/app/client"
	"game/internal/config"
	"game/internal/gamestate/repositories"
	"game/internal/rendering"
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

	/* Compatibility init */

	err = glfw.Init()
	if err != nil {
		panic("failed to initialize GLFW: " + err.Error())
	}

	/* MachineModules */

	repo := repositories.NewRepositories()

	w, err := window.Init(window.Opts{
		Width:   800,
		Height:  600,
		Title:   "another-rpg",
		GLMajor: 3,
		GLMinor: 3,
		Visible: true,
	})
	if err != nil {
		panic("failed to initialize window module: " + err.Error())
	}

	err = gl.Init()
	if err != nil {
		panic("failed to initialize OpenGl: " + err.Error())
	}

	r, err := rendering.NewRendering(ctx, repo)
	if err != nil {
		panic("failed to initialize graphic module: " + err.Error())
	}

	/* Client */

	cl := client.NewBuilder().
		Window(w).
		Rendering(r).
		Repositories(repo).
		MustBuild()

	l.LogAttrs(ctx, slog.LevelInfo, "starting engine running")
	if err = cl.Run(ctx); err != nil {
		panic("failed engine cycle: " + err.Error())
	}

	if err = cl.Shutdown(ctx); err != nil {
		panic("failed to stop client gracefully: " + err.Error())
	}
}
