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

// Canvas base low-level type to accumulate renderable objects into one buffer
type Canvas struct {
	vertices           uint32
	positions          []float32
	textureCoordinates []float32
}

func NewCanvas() *Canvas {
	return NewCanvasCustomBuffer(1024)
}

func NewCanvasCustomBuffer(bufferSizeVertexes int64) *Canvas {
	return &Canvas{
		vertices:           0,
		positions:          make([]float32, 0, bufferSizeVertexes*2*4), // 2 float 4 bytes each
		textureCoordinates: make([]float32, 0, bufferSizeVertexes*2*4), // 2 float 4 bytes each
	}
}

func DefaultCanvasOpts() *CanvasOpts {
	return &CanvasOpts{
		Transformation: math.NoTransform(),
	}
}

func (c *Canvas) Empty() bool {
	return c.vertices == 0
}

func (c *Canvas) Draw(drawable Drawable, opts *CanvasOpts) {
	vertexes := drawable.GetVertexes()

	for _, vertex := range vertexes {
		// todo append transformation
		c.vertices++
		c.positions = append(c.positions, vertex.Position.X(), vertex.Position.Y())
		c.textureCoordinates = append(c.textureCoordinates, vertex.TextureCoordinates.X(), vertex.TextureCoordinates.Y())
	}
}
