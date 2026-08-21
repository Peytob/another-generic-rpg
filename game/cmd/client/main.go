package main

import (
	"context"
	"engine/utils/logger"
	"engine/window"
	"fmt"
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

// Window and OpenGL context parameters.
const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "another-rpg"
	glMajor      = 3
	glMinor      = 3
)

func main() {
	// Client main method should be executed in main application thread due
	// to C libraries restrictions
	runtime.LockOSThread()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfiguration[config.ClientConfiguration]("./config/client.yaml")
	if err != nil {
		panic(err)
	}

	l, err := initLogger(ctx, cfg)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	ctx = logger.ToCtx(ctx, l)

	l.LogAttrs(ctx, slog.LevelInfo, "initializing client")

	/* MachineModules */

	repo := repositories.NewRepositories()

	w, err := initWindow()
	if err != nil {
		panic("failed to initialize window module: " + err.Error())
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

func initLogger(ctx context.Context, cfg *config.ClientConfiguration) (*slog.Logger, error) {
	l, err := config.ClientLogger(ctx, config.LoggerOpts{
		Enabled: cfg.Log.Enabled,
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Log.Level,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("init client logger: %w", err)
	}
	return l, nil
}

// initWindow initializes GLFW, creates the window and loads the OpenGL
// bindings; all C library calls must run on the locked main thread.
func initWindow() (*window.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("glfw init: %w", err)
	}

	w, err := window.Init(window.Opts{
		Width:   windowWidth,
		Height:  windowHeight,
		Title:   windowTitle,
		GLMajor: glMajor,
		GLMinor: glMinor,
		Visible: true,
	})
	if err != nil {
		return nil, fmt.Errorf("window init: %w", err)
	}

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("opengl init: %w", err)
	}

	return w, nil
}
