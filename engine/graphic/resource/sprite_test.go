package resource

import (
	"engine/math/shape"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestSprite_Vertices(t *testing.T) {
	t.Run("should match all rectangle points positions", func(t *testing.T) {
		rect := shape.NewRect(10, 20)
		sprite := NewSprite(rect, shape.NewRect(0, 0))

		data := sprite.Vertices()

		wantPoints := rect.GetPoints()
		if len(data.Vertices) != len(wantPoints) {
			t.Fatalf("vertices count = %d, want %d", len(data.Vertices), len(wantPoints))
		}

		for i, v := range data.Vertices {
			if v.Position != wantPoints[i] {
				t.Errorf("vertex %d position = %v, want %v", i, v.Position, wantPoints[i])
			}
		}
	})

	t.Run("should produce six vertices for rectangle", func(t *testing.T) {
		sprite := NewSprite(shape.NewRect(100, 100), shape.NewRect(0, 0))

		data := sprite.Vertices()

		want := []mgl32.Vec2{
			{0, 0}, {0, 100}, {100, 100},
			{0, 0}, {100, 100}, {100, 0},
		}
		if len(data.Vertices) != len(want) {
			t.Fatalf("vertices count = %d, want %d", len(data.Vertices), len(want))
		}
		for i, v := range data.Vertices {
			if v.Position != want[i] {
				t.Errorf("vertex %d position = %v, want %v", i, v.Position, want[i])
			}
		}
	})

	t.Run("should expose two-triangle indexes", func(t *testing.T) {
		sprite := NewSprite(shape.NewRect(1, 1), shape.NewRect(0, 0))

		data := sprite.Vertices()

		wantIndexes := []uint32{0, 1, 2, 3, 4, 5}
		if len(data.Indexes) != len(wantIndexes) {
			t.Fatalf("indexes count = %d, want %d", len(data.Indexes), len(wantIndexes))
		}
		for i, idx := range data.Indexes {
			if idx != wantIndexes[i] {
				t.Errorf("index %d = %d, want %d", i, idx, wantIndexes[i])
			}
		}
	})
}
