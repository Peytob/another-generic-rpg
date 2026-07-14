package renderer

import (
	gres "engine/graphic/resource"
	"engine/math"
)

type Drawable interface {
	Vertices() gres.VertexData
}

type CanvasOpts struct {
	Transform math.Transform
}

type Canvas interface {
	Empty() bool
	Draw(drawable Drawable, opts *CanvasOpts)
}

func DefaultCanvasOpts() *CanvasOpts {
	return &CanvasOpts{
		Transform: math.NoTransform(),
	}
}
