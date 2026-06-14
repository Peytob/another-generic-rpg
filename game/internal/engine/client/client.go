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
	sprite := resource.NewSprite(shape.NewRect(100, 100), shape.NewRect(0, 0))
	sprite.Transformation().Translate(100, 100)

	w, h := c.window.Size()
	proj := mgl32.Ortho2D(0, float32(w), float32(h), 0)
	rotationTest := 0.0

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

		rotationTest += 0.001

		canvas.Draw(sprite, &renderer.CanvasOpts{
			//Transform: math.NoTransform(),
			Transform: math.NewTransformation().Rotate(float32(rotationTest * 0.3)).Transform(),
		})
		err = c.graphic.Renderer().Render(ctx, canvas, renderer.RenderOpts{
			ViewProj: proj.Mul4(mgl32.Translate3D(250, 250, 0)),
			Model:    math.NewTransformation().Rotate(float32(rotationTest)).Transform(),
			Shader:   c.graphic.Shaders().World,
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
