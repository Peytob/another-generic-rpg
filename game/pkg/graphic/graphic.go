package graphic

import (
	"context"
	"fmt"
	"game/pkg/graphic/renderer"
	gresource "game/pkg/graphic/resource"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Telemetry struct {
	// Graphic backend provider (Vulkan / OGL / etc)
	Name string

	// Graphic backend provider version
	Version string

	// Usually backend hardware name (or other renderer data)
	Renderer string
}

type Graphic struct {
	renderer *renderer.Renderer
	shaders  *gresource.Shaders
}

func NewGraphic(ctx context.Context) (*Graphic, error) {
	shaders, err := LoadShaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load shaders: %w", err)
	}

	return &Graphic{
		renderer: renderer.NewRenderer(),
		shaders:  shaders,
	}, nil
}

// Renderer returns configured and ready renderer object. It can be used to render
func (g *Graphic) Renderer() *renderer.Renderer {
	return g.renderer
}

func (g *Graphic) Shaders() *gresource.Shaders {
	return g.shaders
}

func (g *Graphic) Telemetry() Telemetry {
	return Telemetry{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g *Graphic) Terminate(ctx context.Context) {
	g.shaders.Terminate(ctx)
}
