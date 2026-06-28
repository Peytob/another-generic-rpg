package opengl

import (
	"context"
	"fmt"
	"game/pkg/graphic"
	"game/pkg/graphic/internal/opengl/renderer"
	"game/pkg/graphic/internal/opengl/resource"
	grenderer "game/pkg/graphic/renderer"
	gresource "game/pkg/graphic/resource"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Graphic struct {
	renderer *renderer.Renderer

	shaders       *gresource.Shaders
	uniformBlocks *resource.UniformBlocks
}

func NewGraphic(ctx context.Context) (*Graphic, error) {
	shaders, err := loadShaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load shaders: %w", err)
	}

	uniformBlocks, err := loadUniformBlocks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load uniform blocks: %w", err)
	}

	err = bindUniformBlocks(shaders, uniformBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to bind uniform blocks: %w", err)
	}

	r := renderer.NewRenderer(ctx, uniformBlocks)

	g := &Graphic{
		renderer: r,

		shaders:       shaders,
		uniformBlocks: uniformBlocks,
	}

	return g, nil
}

// Renderer returns configured and ready renderer object. It can be used to render
func (g *Graphic) Renderer() grenderer.Renderer {
	return g.renderer
}

func (g *Graphic) NewCanvas() grenderer.Canvas {
	return renderer.NewCanvas()
}

func (g *Graphic) Shaders() *gresource.Shaders {
	return g.shaders
}

func (g *Graphic) Telemetry() graphic.Telemetry {
	return graphic.Telemetry{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g *Graphic) Terminate(ctx context.Context) {
	terminateShaders(ctx, g.shaders)
	terminateUniformBlocks(ctx, g.uniformBlocks)
	g.renderer.Terminate(ctx)
}
