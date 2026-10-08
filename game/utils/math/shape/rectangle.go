package shape

import "github.com/go-gl/mathgl/mgl32"

type Rect struct {
	//width, height float32
	//x, y          float32
	points []mgl32.Vec2
}

func NewRect(width, height float32) Rect {
	return Rect{
		points: []mgl32.Vec2{
			{0.0, 0.0},
			{0.0, height},
			{width, height},

			{0.0, 0.0},
			{width, height},
			{width, 0.0},
		},
	}
}

func NewZeroRect() Rect {
	return NewRect(0, 0)
}

func (r Rect) GetPoints() []mgl32.Vec2 {
	return r.points
}
