package client

import (
	"context"
	"errors"
	"fmt"
	"game/core/machine"
	"game/graphic"
	"game/input"
	"game/pkg/fsm"
	"game/pkg/logger"
	"game/window"
	"log/slog"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Client struct {
	window  *window.Window
	graphic *graphic.Graphic
	machine machine.Machine
}

func (c *Client) Run(ctx context.Context) error {
	c.dumpRunningInfo(ctx)
	c.window.OnSizeChanged(func(width int, height int) {
		c.onWindowSizeChanged(ctx, width, height)
	})
	width, height := c.window.Size()
	c.onWindowSizeChanged(ctx, width, height) // Initial window size update

	var lastTime time.Time
	for {
		now := time.Now()
		dt := now.Sub(lastTime)
		if lastTime.IsZero() {
			dt = 0
		}
		lastTime = now

		select {
		case <-ctx.Done():
			// Mark window as should close
			logger.FromCtx(ctx).Info("root context done, closing window")
			c.window.Close()
		default:
			// nothing, keep running
		}

		c.window.PoolEvents()

		if c.window.ShouldClose() {
			logger.FromCtx(ctx).Info("stopping game machine")
			if err := c.transition(ctx, machine.StoppedEvent); err != nil && !errors.Is(err, fsm.ErrMachineNotRunning) {
				return fmt.Errorf("failed to stop machine on window closing: %w", err)
			}
			break
		}

		c.window.Clear()
		if err := c.tick(ctx, dt); err != nil {
			return err
		}
		c.window.Show()

		if !c.machine.IsRunning() {
			logger.FromCtx(ctx).Info("machine running complete, closing window")
			c.window.Close()
		}
	}

	return nil
}

func (c *Client) Terminate(ctx context.Context) error {
	ctx = logger.WithLogValue(ctx, slog.String("svc", "client_termination"))
	logger.FromCtx(ctx).Info("terminating client")
	return nil
}

func (c *Client) dumpRunningInfo(ctx context.Context) {
	aboutGraphic := c.graphic.Telemetry().About()
	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "running with graphic",
		slog.String("name", aboutGraphic.Name),
		slog.String("version", aboutGraphic.Version),
		slog.String("renderer", aboutGraphic.Renderer))
}

func (c *Client) onWindowSizeChanged(ctx context.Context, width int, height int) {
	proj := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	if err := c.graphic.Uniform().Uniforms().ProjectionMatrix(proj); err != nil {
		logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelError, "failed to set uniform projection", logger.Error(err))
	}
}

func (c *Client) tick(ctx context.Context, dt time.Duration) error {
	if !c.machine.IsRunning() {
		return nil
	}

	frame := machine.Frame{
		Dt:    dt,
		Input: input.Frame{},
	}

	ev, err := c.machine.State().Tick(ctx, frame)
	if err != nil {
		return fmt.Errorf("state %q update: %w", c.machine.State().Identifier(), err)
	}

	if ev != machine.NoEvent {
		return c.transition(ctx, ev)
	}

	return nil
}

func (c *Client) transition(ctx context.Context, ev machine.Event) error {
	current := c.machine.State()

	logger.FromCtx(ctx).LogAttrs(
		ctx,
		slog.LevelDebug,
		"transitioning machine",
		slog.Int("event", int(ev)),
		slog.String("current_state", string(current.Identifier())))

	if err := c.machine.Event(ev); err != nil {
		return fmt.Errorf("fsm event %d: %w", ev, err)
	}

	if err := current.OnExit(ctx); err != nil {
		return fmt.Errorf("state %q onExit: %w", current.Identifier(), err)
	}

	if c.machine.IsRunning() {
		next := c.machine.State()

		if err := next.OnEnter(ctx); err != nil {
			return fmt.Errorf("state %q onEnter: %w", next.Identifier(), err)
		}
	}

	logger.FromCtx(ctx).LogAttrs(
		ctx,
		slog.LevelInfo,
		"transition performed",
		slog.Int("event", int(ev)),
		slog.String("previous_state", string(current.Identifier())),
		slog.String("new_state", string(current.Identifier())),
		slog.Bool("is_machine_running", c.machine.IsRunning()))

	return nil
}
