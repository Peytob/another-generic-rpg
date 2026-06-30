package resource

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position           mgl32.Vec2
	TextureCoordinates mgl32.Vec2
}

type VertexData struct {
	// Vertices array of vertices to draw. Should be in local coordinates with local transformation
	Vertices []Vertex

	// Indexes EBO data to draw. Should starts from zero (i.e. 0,1,2, 1,2,3, 2,3,4...)
	Indexes []uint32
}
