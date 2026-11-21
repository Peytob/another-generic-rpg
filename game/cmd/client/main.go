package main

import (
	"context"
	"fmt"
	"game/app/client"
	"game/config"
	"game/graphic"
	"game/pkg/logger"
	"game/window"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	glMajorVersion = 3
	glMinorVersion = 3
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

	log, err := config.ClientLogger(ctx, config.LoggerOpts{
		Enabled: cfg.Log.Enabled,
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Log.Level,
		}),
	})
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	ctx = logger.ToCtx(ctx, log)

	err = run(ctx, cfg)
	if err != nil {
		panic("failed while running game: " + err.Error())
	}
}

func run(ctx context.Context, cfg config.ClientConfiguration) error {
	log := logger.FromCtx(ctx)

	log.LogAttrs(ctx, slog.LevelInfo, "initializing glfw")
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw init: %w", err)
	}
	defer func() {
		log.LogAttrs(ctx, slog.LevelInfo, "terminating glfw")
		glfw.Terminate()
	}()

	log.LogAttrs(ctx, slog.LevelInfo, "initializing window")
	glfwWindow, err := window.Init(window.Opts{
		Width:   cfg.Window.Width,
		Height:  cfg.Window.Height,
		Title:   "window",
		GLMajor: glMajorVersion,
		GLMinor: glMinorVersion,
		Visible: true,
	})
	if err != nil {
		return fmt.Errorf("window init: %w", err)
	}
	defer func() {
		if err := glfwWindow.Terminate(ctx); err != nil {
			log.LogAttrs(ctx, slog.LevelError, "terminate glfw window", logger.Error(err))
		}
	}()

	log.LogAttrs(ctx, slog.LevelInfo, "initializing opengl")
	if err := gl.Init(); err != nil {
		return fmt.Errorf("opengl init: %w", err)
	}

	openglGraphic, err := graphic.Initialize(ctx, graphic.Opts{})
	if err != nil {
		return fmt.Errorf("graphic init: %w", err)
	}
	defer func() {
		if err := openglGraphic.Terminate(ctx); err != nil {
			log.LogAttrs(ctx, slog.LevelError, "terminate graphic", logger.Error(err))
		}
	}()

	return runClient(ctx, glfwWindow, openglGraphic)
}

func runClient(ctx context.Context, glfwWindow *window.Window, openglGraphic *graphic.Graphic) error {
	log := logger.FromCtx(ctx)

	c := client.NewBuilder().
		Window(glfwWindow).
		Graphic(openglGraphic).
		MustBuild()

	defer func() {
		if err := c.Terminate(ctx); err != nil {
			log.LogAttrs(ctx, slog.LevelError, "terminate client", logger.Error(err))
		}
	}()

	if err := c.Run(ctx); err != nil {
		return fmt.Errorf("client run: %w", err)
	}

	return nil
}
