package scene

import (
	"game/graphic/drawable"
	"game/pkg/math"
	"game/pkg/math/shape"

	"github.com/go-gl/mathgl/mgl32"
)

const verticesPerSprite = 6

// Sprite just textured rect with transformation
type Sprite struct {
	transformation *math.Transformation
	rect           shape.Rect
	textureRect    shape.Rect

	vertices []drawable.Vertex
}

func NewSprite(rect shape.Rect, textureRect shape.Rect) *Sprite {
	return &Sprite{
		transformation: math.NewTransformation(),
		rect:           rect,
		textureRect:    textureRect,
		vertices:       nil, // will be computed lazily on first Vertices call
	}
}

func (s *Sprite) Transformation() *math.Transformation {
	return s.transformation
}

func (s *Sprite) LocalRect() shape.Rect {
	return s.rect
}

func (s *Sprite) TextureRect() shape.Rect {
	return s.textureRect
}

func (s *Sprite) Vertices() drawable.VertexData {
	if len(s.vertices) == 0 {
		s.updateVertices()
	}

	return drawable.VertexData{
		Vertices: s.vertices,
		Indexes: []uint32{
			0, 1, 2,
			3, 4, 5,
		},
	}
}

func (s *Sprite) updateVertices() {
	s.vertices = make([]drawable.Vertex, verticesPerSprite)

	rectPoints := s.rect.GetPoints()
	for i := range rectPoints {
		s.vertices[i] = drawable.Vertex{
			Position:           rectPoints[i],
			TextureCoordinates: mgl32.Vec2{},
		}
	}
}
