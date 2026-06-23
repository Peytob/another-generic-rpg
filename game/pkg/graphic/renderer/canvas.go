package renderer

import (
	gres "game/pkg/graphic/resource"
	"game/pkg/math"
)

type Drawable interface {
	GetVertexes() gres.IndexesVertices
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
