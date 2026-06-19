package client

import (
	"context"
	"fmt"
	"game/internal/engine/fsm"
	"game/pkg/engine"
	"game/pkg/graphic"
	"game/pkg/graphic/renderer"
	"game/pkg/graphic/resource"
	"game/pkg/math"
	"game/pkg/math/shape"
	"game/pkg/utils/logger"
	"game/pkg/window"
	"log/slog"

	"github.com/go-gl/mathgl/mgl32"
)

type Client struct {
	fsm     engine.Machine
	window  *window.Window
	graphic *graphic.Graphic
}

func (c *Client) Run(ctx context.Context) error {
	c.dumpRunningInfo(ctx)

	// test
	//w, h := c.window.Size()
	sprite := resource.NewSprite(shape.NewRect(100, 100), shape.NewRect(0, 0))
	sprite.Transformation().Translate(-50, 50)

	c.window.OnSizeChanged(c.onWindowSizeChanged)
	c.onWindowSizeChanged(c.window.Size()) // Initial window size update

	rotate := float32(0.0)

	for {
		var err error

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
			err = c.fsm.Event(fsm.StoppedEvent)
			if err != nil {
				return fmt.Errorf("failed to stop machine on window close: %w", err)
			}
			break
		}

		event, err := c.fsm.State().Update(ctx)
		if err != nil {
			return fmt.Errorf("error while executing FSM update: %w", err)
		}

		if event != fsm.NoEvent {
			err = c.fsm.Event(event)
			if err != nil {
				return fmt.Errorf("error while changing FSM state: %w", err)
			}
		}

		c.window.Clear()

		/* test */

		canvas := renderer.NewCanvas()
		canvas.Draw(sprite, &renderer.CanvasOpts{
			Transform: sprite.Transformation().Transform(),
		})
		sprite.Transformation().Rotate(rotate)
		rotate += 0.0002

		err = c.graphic.Renderer().Render(ctx, canvas, renderer.RenderOpts{
			View: mgl32.Ident4(), // todo camera
			Model: math.NewTransformation().
				//Translate(float32(w)/2, float32(h)/2).
				//Rotate(rotate).
				Transform(),
			Shader: c.graphic.Shaders().World,
		})
		if err != nil {
			return fmt.Errorf("failed to render: %w", err)
		}

		/* end test */

		c.window.Show()

		if !c.fsm.IsRunning() {
			logger.FromCtx(ctx).Info("game FSM is completed, closing window")
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
	graphicApiInfo := c.graphic.Telemetry()
	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "running with graphic",
		slog.String("name", graphicApiInfo.Name),
		slog.String("version", graphicApiInfo.Version),
		slog.String("renderer", graphicApiInfo.Renderer))
}

func (c *Client) onWindowSizeChanged(width int, height int) {
	proj := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	c.graphic.Renderer().UpdateProjection(proj)
}
