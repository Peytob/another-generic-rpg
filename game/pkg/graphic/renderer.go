package graphic

import "context"

// RenderOpts optional options for renderer
type RenderOpts struct {
	// RenderTarget target to render. Renderer should use default target if nil
	RenderTarget RenderTarget
}

// View contains data to compute ViewProj matrix
type View struct {
}

// RenderTarget describes render target
type RenderTarget interface {
}

// Renderer Used to render canvas to
type Renderer interface {
	Render(ctx context.Context, canvas Canvas, view View, opts RenderOpts) error
}
