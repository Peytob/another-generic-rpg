package shape

import "github.com/go-gl/mathgl/mgl32"

type Rect struct {
	//width, height float32
	//x, y          float32
	points []mgl32.Vec2
}

func NewRect(width, height float32) Rect {
	return Rect{
		//width:  width,
		//height: height,
		//x:      0,
		//y:      0,
		points: mgl32.Rect(width, height),
	}
}

func NewPositionedRect(position mgl32.Vec2, width, height float32) Rect {
	points := mgl32.Rect(width, height)
	for _, v := range points {
		v.Add(position)
	}

	return Rect{
		//width:  width,
		//height: height,
		//x:      position.X(),
		//y:      position.Y(),
		points: points,
	}
}

func (r Rect) GetPoints() []mgl32.Vec2 {
	return r.points
}
