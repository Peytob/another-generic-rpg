package renderer

import (
	"context"
	"engine/graphic/resource"
	"engine/math"

	"github.com/go-gl/mathgl/mgl32"
)

// RenderOpts optional options for renderer
type RenderOpts struct {
	// View matrix used to make ViewProjection matrix. Or in other words: camera
	View mgl32.Mat4

	// Model matrix used to compute rendering model transformations
	Model math.Transform

	// Shader used to render shader. Required
	Shader resource.Shader

	// RenderTarget target to render.
	RenderTarget RenderTarget

	// Mode primitive rasterization mode (fill or wireframe).
	Mode DrawMode
}

// RenderTarget describes render target
type RenderTarget struct {
	RenderBufferID int32
}

type DrawMode uint8

const (
	Fill DrawMode = iota
	Wireframe
)

type Renderer interface {
	Render(ctx context.Context, canvas Canvas, opts RenderOpts) error
	Terminate(ctx context.Context)
}
