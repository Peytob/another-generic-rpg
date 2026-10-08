package client

import (
	"context"
	"game/app"
	"game/graphic"
	"game/utils/logger"
	"game/window"
	"log/slog"
	"time"
)

type Client struct {
	window  *window.Window
	graphic *graphic.Graphic
}

func (c *Client) Run(ctx context.Context) error {
	appCtx := app.WrapCtx(ctx)
	//var err error

	c.dumpRunningInfo(appCtx)
	c.window.OnSizeChanged(c.onWindowSizeChanged)
	c.onWindowSizeChanged(c.window.Size()) // Initial window size update

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
			println(dt)
		default:
			// nothing, keep running
		}

		c.window.PoolEvents()

		if c.window.ShouldClose() {
			logger.FromCtx(ctx).Debug("stopping game machine")
			//err = c.runner.Event(ctx, state.StoppedEvent)
			//if err != nil {
			//	return fmt.Errorf("failed to stop runner on window close: %w", err)
			//}
			break
		}

		c.window.Clear()
		c.window.Show()
	}

	return nil
}

func (c *Client) Terminate(ctx context.Context) error {
	ctx = logger.WithLogValue(ctx, slog.String("svc", "client_termination"))
	logger.FromCtx(ctx).Info("terminating client")
	return nil
}

func (c *Client) dumpRunningInfo(ctx app.Context) {
	aboutGraphic := c.graphic.Telemetry().About()
	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "running with graphic",
		slog.String("name", aboutGraphic.Name),
		slog.String("version", aboutGraphic.Version),
		slog.String("renderer", aboutGraphic.Renderer))
}

func (c *Client) onWindowSizeChanged(width int, height int) {
	//proj := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	//
	//ub, ok := c.rendering.Services().Uniform.GetByName(rendering.ProjViewUniformBlock)
	//if !ok {
	//	// todo log error via contex
	//	slog.Default().Error("failed to find uniform block", "name", rendering.ProjViewUniformBlock)
	//	return
	//}
	//
	//err := c.rendering.Services().Uniform.SetUniformVariableMat4(ub, rendering.ProjUniform, proj)
	//if err != nil {
	//	// todo log error via context
	//	slog.Default().Error("failed to set uniform block variable",
	//		"name", rendering.ProjViewUniformBlock,
	//		"variable", rendering.ProjUniform)
	//	return
	//}
}
