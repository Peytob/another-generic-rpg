package types

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	position mgl32.Vec2
	area     mgl32.Vec2
	maxArea  mgl32.Vec2
}

func NewCamera(position mgl32.Vec2, maxArea mgl32.Vec2) *Camera {
	return &Camera{
		position: position,
		maxArea:  maxArea,
	}
}

func (c *Camera) Position(position mgl32.Vec2) {
	c.position = position
}

func (c *Camera) GetPosition() mgl32.Vec2 {
	return c.position
}

func (c *Camera) Move(delta mgl32.Vec2) {
	c.position = c.position.Add(delta)
}

func (c *Camera) Area(w, h int) {
	c.area[0] = min(max(0, float32(w)), c.maxArea[0])
	c.area[1] = min(max(0, float32(h)), c.maxArea[1])
}

func (c *Camera) GetArea() mgl32.Vec2 {
	return c.area
}

func (c *Camera) GetMaxArea() mgl32.Vec2 {
	return c.maxArea
}

func (c *Camera) ViewMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(-c.position.X(), -c.position.Y(), 0)
}
