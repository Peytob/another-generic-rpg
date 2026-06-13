package client

import (
	"context"
	"fmt"
	"game/internal/engine/fsm"
	"game/pkg/engine"
	"game/pkg/graphic"
	"game/pkg/graphic/resource"
	"game/pkg/math/shape"
	"game/pkg/utils/logger"
	"game/pkg/window"
	"log/slog"
)

type Client struct {
	fsm     engine.Machine
	window  *window.Window
	graphic *graphic.Graphic
}

func (c *Client) Run(ctx context.Context) error {
	c.dumpRunningInfo(ctx)

	sprite := resource.NewSprite(shape.NewRect(0.2, 0.2), shape.NewRect(0, 0))
	vertex, err := graphic.BuildShader(ctx, `
#version 330 core
layout (location = 0) in vec3 aPos;

void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
`, graphic.VertexShader)
	if err != nil {
		return err
	}

	fragment, err := graphic.BuildShader(ctx, `
#version 330 core
out vec4 FragColor;

void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
`, graphic.FragmentShader)
	if err != nil {
		return err
	}

	spb := graphic.ShaderProgramBuilder{
		Vertex:   vertex,
		Fragment: fragment,
	}

	program, err := graphic.BuildShaderProgram(ctx, spb)
	if err != nil {
		return err
	}

	vertex.Delete()
	fragment.Delete()

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

		/* test */

		canvas := graphic.NewCanvas()

		canvas.Draw(sprite, graphic.DefaultCanvasOpts())
		err = c.graphic.Renderer().Render(ctx, canvas, graphic.RenderOpts{
			Shader: program,
		})
		if err != nil {
			return fmt.Errorf("failed to render: %w", err)
		}

		/* end test */

		c.window.Clear()
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
