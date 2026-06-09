package math

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Transform transformation mat3 matrix
type Transform mgl32.Mat3

func (t Transform) TransformPoint(vec2 mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Mat3(t).Mul3x1(vec2.Vec3(1)).Vec2()
}

func NoTransform() Transform {
	return Transform(mgl32.Ident3())
}

// Transformation decomposed Transform matrix that can be recalculated
type Transformation struct {
	t Transform
}

func NewTransformation() Transformation {
	return Transformation{
		t: NoTransform(),
	}
}

func (t *Transformation) Transform() Transform {
	// TRS Order
	return t.t
}
