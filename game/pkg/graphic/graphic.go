package graphic

import (
	"context"

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
	renderer *Renderer
}

func NewGraphic() (*Graphic, error) {
	return &Graphic{
		renderer: NewRenderer(),
	}, nil
}

// Renderer returns configured and ready renderer object. It can be used to render
func (g *Graphic) Renderer() *Renderer {
	return g.renderer
}

func (g *Graphic) Telemetry() Telemetry {
	return Telemetry{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g *Graphic) Terminate(_ context.Context) {
}
