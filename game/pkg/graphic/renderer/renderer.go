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

	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	defer gl.DeleteVertexArrays(1, &vao)

	gl.GenBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vbo)

	gl.GenBuffers(1, &ebo)
	defer gl.DeleteBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	{
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(canvas.positions)*4*2, gl.Ptr(canvas.positions), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(canvas.elements)*4*2, gl.Ptr(canvas.elements), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(0)
	}
	gl.UseProgram(opts.Shader.Id())
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(canvas.elements)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	return nil
}
