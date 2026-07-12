package math

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Transform transformation mat4 matrix
type Transform mgl32.Mat4

func NoTransform() Transform {
	return Transform(mgl32.Ident4())
}

func (t Transform) TransformPoint(vec2 mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Mat4(t).Mul4x1(vec2.Vec4(0, 1)).Vec2()
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
		mat := mgl32.Ident4()

		if !isZero(t.translate) {
			translate := mgl32.Translate3D(t.translate.X(), t.translate.Y(), 0)
			mat = mat.Mul4(translate)
		}

		if t.rotate != 0 {
			rotate := mgl32.HomogRotate3DZ(t.rotate)
			mat = mat.Mul4(rotate)
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
