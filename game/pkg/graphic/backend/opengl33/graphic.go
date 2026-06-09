package opengl33

import (
	"context"
	"game/pkg/graphic"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Graphic struct {
	renderer Renderer
	factory  Factory
}

func NewGraphic() (*Graphic, error) {
	g := &Graphic{}

	g.factory = NewFactory()
	g.renderer = NewRenderer(g.factory.RenderTarget().WindowRenderTarget().(RenderTarget))

	return g, nil
}

func (g Graphic) NewCanvas() graphic.Canvas {
	return Canvas{}
}

func (g Graphic) Renderer() graphic.Renderer {
	return g.renderer
}

func (g Graphic) Factory() graphic.Factory {
	return g.factory
}

func (g Graphic) Telemetry() graphic.Telemetry {
	return graphic.Telemetry{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g Graphic) Terminate(_ context.Context) {
}
