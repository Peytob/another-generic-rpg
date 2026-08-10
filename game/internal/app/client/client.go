package client

import (
	"context"
	"engine/graphic"
	"engine/utils/logger"
	"engine/window"
	"fmt"
	"game/internal/gamestate/gamemachine"
	"game/internal/rendering"
	"log/slog"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Client struct {
	runner  *gamemachine.Runner
	window  *window.Window
	graphic graphic.Graphic
}

func (c *Client) Run(ctx context.Context) error {
	c.dumpRunningInfo(ctx)

	c.window.OnSizeChanged(c.onWindowSizeChanged)
	c.onWindowSizeChanged(c.window.Size()) // Initial window size update

	if err := c.runner.Start(ctx); err != nil {
		return fmt.Errorf("failed to start runner: %w", err)
	}

	var lastTime time.Time
	for {
		var err error

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
			logger.FromCtx(ctx).Debug("stopping game machine")
			err = c.runner.Event(ctx, gamemachine.StoppedEvent)
			if err != nil {
				return fmt.Errorf("failed to stop runner on window close: %w", err)
			}
			break
		}

		c.window.Clear()
		err = c.runner.Update(ctx, dt)
		if err != nil {
			return fmt.Errorf("error while executing runner update: %w", err)
		}
		c.window.Show()

		if !c.runner.IsRunning() {
			logger.FromCtx(ctx).Info("game runner is completed, closing window")
			break
		}
	}

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	logger.FromCtx(ctx).Info("shutting down client")
	c.window.Terminate(ctx)
	c.graphic.Terminate(ctx)
	return nil
}

func (c *Client) dumpRunningInfo(ctx context.Context) {
	graphicAPIInfo := c.graphic.Telemetry()
	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "running with graphic",
		slog.String("name", graphicAPIInfo.Name),
		slog.String("version", graphicAPIInfo.Version),
		slog.String("renderer", graphicAPIInfo.Renderer))
}

func (c *Client) onWindowSizeChanged(width int, height int) {
	proj := mgl32.Ortho2D(0, float32(width), float32(height), 0)

	ub, ok := c.graphic.Repositories().Uniform.ByName(rendering.ProjViewUniformBlock)
	if !ok {
		// todo log error via contex
		slog.Default().Error("failed to find uniform block", "name", rendering.ProjViewUniformBlock)
		return
	}

	err := c.graphic.Services().Uniform.SetUniformVariableMat4(ub, rendering.ProjUniform, proj)
	if err != nil {
		// todo log error via context
		slog.Default().Error("failed to set uniform block variable", "name", rendering.ProjViewUniformBlock, "variable", rendering.ProjUniform)
		return
	}
}
