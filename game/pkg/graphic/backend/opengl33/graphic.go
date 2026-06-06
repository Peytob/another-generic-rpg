package opengl33

import (
	"context"
	"game/pkg/graphic"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Graphic struct {
	renderer Renderer
}

func NewGraphic() (*Graphic, error) {
	g := &Graphic{}

	return g, nil
}

func (g Graphic) NewCanvas() graphic.Canvas {
	return Canvas{}
}

func (g Graphic) Renderer() graphic.Renderer {
	return g.renderer
}

func (g Graphic) Api() graphic.Api {
	return graphic.Api{
		Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
		Version:  gl.GoStr(gl.GetString(gl.VERSION)),
		Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
	}
}

func (g Graphic) Terminate(_ context.Context) {
}
