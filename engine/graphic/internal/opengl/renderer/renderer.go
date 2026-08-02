package renderer

import (
	"context"
	oglresource "engine/graphic/internal/opengl/resource"
	grenderer "engine/graphic/renderer"
	"engine/graphic/service"
	"engine/utils/logger"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	vao           oglresource.VertexArray
	shaderService service.Shader
}

func NewRenderer(ctx context.Context, shaderService service.Shader) *Renderer {
	// TODO make VAO repository and management

	var defaultVao uint32
	gl.GenVertexArrays(1, &defaultVao)
	logger.FromCtx(ctx).Info("created vertex array buffer", slog.Int64("id", int64(defaultVao)))

	gl.BindVertexArray(defaultVao)

	return &Renderer{
		vao:           oglresource.VertexArray(defaultVao),
		shaderService: shaderService,
	}
}

func (r *Renderer) Render(ctx context.Context, canvas grenderer.Canvas, opts grenderer.RenderOpts) error {
	oglCanvas, ok := canvas.(*Canvas)
	if !ok {
		return errors.New("canvas is not ogl-compatible")
	}

	return r.renderOgl(ctx, *oglCanvas, opts)
}

func (r *Renderer) renderOgl(_ context.Context, canvas Canvas, opts grenderer.RenderOpts) error {
	// TODO buffers reusing and management

	if canvas.Empty() {
		return nil
	}

	shaderProgram := oglresource.ShaderProgram(opts.Shader.ID)
	if shaderProgram == 0 {
		return errors.New("invalid shader program")
	}

	var vbo, ebo uint32

	gl.GenBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vbo)

	gl.GenBuffers(1, &ebo)
	defer gl.DeleteBuffers(1, &ebo)

	gl.BindVertexArray(r.vao.ID())
	{
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(canvas.positions)*4*2, gl.Ptr(canvas.positions), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(canvas.elements)*4, gl.Ptr(canvas.elements), gl.STATIC_DRAW)
	}

	gl.UseProgram(shaderProgram.ID())

	// TODO constants for uniforms

	if err := r.shaderService.UniformTransform(opts.Shader, "u_model\x00", opts.Model); err != nil {
		return fmt.Errorf("failed to set model uniform: %w", err)
	}

	if err := r.shaderService.UniformMat4(opts.Shader, "u_view\x00", opts.View); err != nil {
		return fmt.Errorf("failed to set view uniform: %w", err)
	}

	gl.BindVertexArray(r.vao.ID())
	gl.PolygonMode(gl.FRONT_AND_BACK, polygonMode(opts.Mode))
	gl.DrawElements(gl.TRIANGLES, int32(len(canvas.elements)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	return nil
}

func polygonMode(mode grenderer.DrawMode) uint32 {
	switch mode {
	case grenderer.Wireframe:
		return gl.LINE
	default:
		return gl.FILL
	}
}

func (r *Renderer) Terminate(ctx context.Context) {
	logger.FromCtx(ctx).Info("deleting vertex array", slog.Int64("id", int64(r.vao)))
	r.vao.Delete()
}
