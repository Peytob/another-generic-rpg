package resource

import (
	"game/pkg/math"
	"game/pkg/math/shape"
)

// Sprite just textured rect with transformation
type Sprite struct {
	transformation *math.Transformation
	rect           shape.Rect
	textureRect    shape.Rect

	vertices []Vertex
}

func NewSprite(rect shape.Rect, textureRect shape.Rect) *Sprite {
	return &Sprite{
		transformation: math.NewTransformation(),
		rect:           rect,
		textureRect:    textureRect,
		vertices:       nil, // will be computed lazy while GetVertices call
	}
}

func (s *Sprite) Transformation() *math.Transformation {
	return s.transformation
}

func (s *Sprite) GetLocalRect() shape.Rect {
	return s.rect
}

func (s *Sprite) TextureRect() shape.Rect {
	return s.textureRect
}

func (s *Sprite) GetVertexes() IndexesVertices {
	if len(s.vertices) == 0 {
		s.updateVertices()
	}

	return IndexesVertices{
		Vertexes: s.vertices,
		Indexes: []uint32{
			0, 1, 2,
			3, 4, 5,
		},
	}
}

func (s *Sprite) updateVertices() {
	s.vertices = make([]Vertex, 6)

	rectPoints := s.rect.GetPoints()
	for i := range rectPoints {
		s.vertices[i] = Vertex{
			Position:           rectPoints[i],
			TextureCoordinates: rectPoints[i],
		}
	}
}
