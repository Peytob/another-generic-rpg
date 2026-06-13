package renderer

import (
	"context"
	"game/pkg/graphic/resource"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// RenderOpts optional options for renderer
type RenderOpts struct {
	// View contains camera data
	View View

	// Shader used to render shader. Required
	Shader resource.ShaderProgram

	// RenderTarget target to render.
	RenderTarget RenderTarget
}

// View contains data to compute ViewProj matrix
type View struct {
}

// RenderTarget describes render target
type RenderTarget struct {
	RenderBufferId int32
}

// Renderer low-level canvas rendering objects
type Renderer struct {
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Render(ctx context.Context, canvas *Canvas, opts RenderOpts) error {
	// todo buffers reusing

	if canvas.Empty() {
		return nil
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(canvas.positions)*4, gl.Ptr(canvas.positions), gl.STATIC_DRAW)

	gl.DeleteBuffers(1, &vbo)

	return nil
}
