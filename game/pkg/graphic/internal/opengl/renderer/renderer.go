package renderer

import (
	"context"
	"errors"
	"fmt"
	"game/pkg/graphic/internal/opengl/resource"
	"game/pkg/graphic/renderer"
	"game/pkg/utils/logger"
	"log/slog"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Renderer low-level canvas rendering objects
type Renderer struct {
	uniformBlocks *resource.UniformBlocks

	vao resource.VertexArray // default for now
}

func NewRenderer(ctx context.Context, uniformBlocks *resource.UniformBlocks) *Renderer {
	var tilemapVao uint32
	gl.GenVertexArrays(1, &tilemapVao)
	logger.FromCtx(ctx).Info("created vertex array buffer", slog.Int64("id", int64(tilemapVao)))

	gl.BindVertexArray(tilemapVao)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)

	return &Renderer{
		uniformBlocks: uniformBlocks,
		vao:           resource.VertexArray(tilemapVao),
	}
}

func (r *Renderer) UpdateProjection(proj mgl32.Mat4) {
	r.uniformBlocks.Proj.SetProjectionMatrix(proj)
}

func (r *Renderer) Render(ctx context.Context, canvas renderer.Canvas, opts renderer.RenderOpts) error {
	oglCanvas, ok := canvas.(*Canvas)
	if !ok {
		return errors.New("canvas is not ogl-compatible")
	}

	return r.renderOgl(ctx, *oglCanvas, opts)
}

func (r *Renderer) renderOgl(_ context.Context, canvas Canvas, opts renderer.RenderOpts) error {
	// todo buffers reusing

	if canvas.Empty() {
		return nil
	}

	var vbo, ebo uint32

	gl.GenBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vbo)

	gl.GenBuffers(1, &ebo)
	defer gl.DeleteBuffers(1, &ebo)

	gl.BindVertexArray(r.vao.Id())
	{
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(canvas.positions)*4*2, gl.Ptr(canvas.positions), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(canvas.elements)*4, gl.Ptr(canvas.elements), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(0)
	}

	gl.UseProgram(opts.Shader.Id())

	if err := resource.UniformTransform(opts.Shader, "u_model\x00", opts.Model); err != nil {
		return fmt.Errorf("failed to set model uniform: %w", err)
	}

	if err := resource.UniformMat4(opts.Shader, "u_view\x00", opts.View); err != nil {
		return fmt.Errorf("failed to set view uniform: %w", err)
	}

	gl.BindVertexArray(r.vao.Id())
	gl.DrawElements(gl.TRIANGLES, int32(len(canvas.elements)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	return nil
}

func (r *Renderer) Terminate(ctx context.Context) {
	logger.FromCtx(ctx).Info("deleting vertex array", slog.Int64("id", int64(r.vao)))
	r.vao.Delete()
}
