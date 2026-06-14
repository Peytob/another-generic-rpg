package renderer

import (
	"context"
	"fmt"
	"game/pkg/graphic/resource"
	"game/pkg/math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// RenderOpts optional options for renderer
type RenderOpts struct {
	// View matrix used to make ViewProjection matrix. Or in other words: camera
	View mgl32.Mat4

	// Model matrix used to compute rendering model transformations
	Model math.Transform

	// Shader used to render shader. Required
	Shader resource.ShaderProgram
	// RenderTarget target to render.
	RenderTarget RenderTarget
}

// RenderTarget describes render target
type RenderTarget struct {
	RenderBufferId int32
}

// Renderer low-level canvas rendering objects
type Renderer struct {
	proj mgl32.Mat4
}

func NewRenderer() *Renderer {
	return &Renderer{
		proj: mgl32.Ident4(),
	}
}

func (r *Renderer) UpdateProjection(proj mgl32.Mat4) {
	r.proj = proj
}

func (r *Renderer) Render(_ context.Context, canvas *Canvas, opts RenderOpts) error {
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

	if err := opts.Shader.UniformTransform("model\x00", opts.Model); err != nil {
		return fmt.Errorf("failed to set model uniform: %w", err)
	}

	if err := opts.Shader.UniformMat4("view\x00", opts.View); err != nil {
		return fmt.Errorf("failed to set projection uniform: %w", err)
	}

	if err := opts.Shader.UniformMat4("proj\x00", r.proj); err != nil {
		return fmt.Errorf("failed to set view uniform: %w", err)
	}

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(canvas.elements)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	return nil
}
