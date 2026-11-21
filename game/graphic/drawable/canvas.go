package drawable

import (
	"game/pkg/math"

	"github.com/go-gl/mathgl/mgl32"
)

// Canvas base low-level type to accumulate renderable objects into one buffer
type Canvas struct {
	// vertices count
	vertices uint32

	// elements EBO data (just 1,2,3,4,5 ... for now)
	elements []uint32

	// positions vertex position data (drawable local coordinates, ie after model and drawable transformations)
	positions []mgl32.Vec2

	// textureCoordinates local texture coordinates (normalized coordinates in target texture)
	textureCoordinates []mgl32.Vec2
}

func NewCanvas() *Canvas {
	return NewCanvasCustomBuffer(1024)
}

func NewCanvasCustomBuffer(bufferSizeVertexes uint32) *Canvas {
	return &Canvas{
		vertices:           0,
		elements:           make([]uint32, 0, bufferSizeVertexes*3),
		positions:          make([]mgl32.Vec2, 0, bufferSizeVertexes),
		textureCoordinates: make([]mgl32.Vec2, 0, bufferSizeVertexes),
	}
}

type CanvasOpts struct {
	Transform math.Transform
}

func DefaultCanvasOpts() CanvasOpts {
	return CanvasOpts{
		Transform: math.NoTransform(),
	}
}

func (c *Canvas) Empty() bool {
	return c.vertices == 0
}

func (c *Canvas) Draw(drawable Drawable, opts CanvasOpts) {
	d := drawable.Vertices()
	startIndex := c.vertices

	for _, vertex := range d.Vertices {
		// TODO append transformation
		c.vertices++

		transformedPosition := opts.Transform.TransformPoint(vertex.Position)
		c.positions = append(c.positions, transformedPosition)
		c.textureCoordinates = append(c.textureCoordinates, vertex.TextureCoordinates)
	}

	for _, index := range d.Indexes {
		c.elements = append(c.elements, startIndex+index)
	}
}
