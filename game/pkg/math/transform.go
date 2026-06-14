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
	t          Transform
	needUpdate bool

	translate mgl32.Vec2
	rotate    float32
}

func NewTransformation() *Transformation {
	return &Transformation{
		t:          NoTransform(),
		needUpdate: false,
	}
}

func (t *Transformation) Transform() Transform {
	// TRS Order
	if t.needUpdate {
		mat := mgl32.Ident3()

		if !isZero(t.translate) {
			translate := mgl32.Translate2D(t.translate.X(), t.translate.Y())
			mat = mat.Mul3(translate)
		}

		if t.rotate != 0 {
			rotate := mgl32.HomogRotate2D(t.rotate)
			mat = mat.Mul3(rotate)
		}

		t.t = Transform(mat)
		t.needUpdate = false
	}

	return t.t
}

func (t *Transformation) Translate(dx float32, dy float32) *Transformation {
	t.translate[0] = dx
	t.translate[1] = dy
	t.needUpdate = true
	return t
}

func (t *Transformation) Rotate(angle float32) *Transformation {
	t.rotate = angle
	t.needUpdate = true
	return t
}

func (t *Transformation) Clone() *Transformation {
	clone := *t
	return &clone
}

func isZero(translate mgl32.Vec2) bool {
	return translate.X() == 0 && translate.Y() == 0
}
