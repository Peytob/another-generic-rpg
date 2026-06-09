package graphic

import (
	gres "game/pkg/graphic/resource"
	"game/pkg/math"
)

type Drawable interface {
	GetVertexes() []gres.Vertex
}

type CanvasOpts struct {
	Transformation math.Transform
}

// Canvas base type to accumulate all renderable objects
type Canvas interface {
	Draw(drawable Drawable, opts *CanvasOpts)
}

func DefaultCanvasOpts() *CanvasOpts {
	return &CanvasOpts{
		Transformation: math.NoTransform(),
	}
}
