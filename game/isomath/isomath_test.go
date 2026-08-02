package isomath

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func approxEq(a, b mgl32.Vec2) bool {
	return math.Abs(float64(a.X()-b.X())) < 1e-5 && math.Abs(float64(a.Y()-b.Y())) < 1e-5
}

func TestGridToScreen(t *testing.T) {
	cases := []struct {
		name string
		c, r float32
		want mgl32.Vec2
	}{
		{"origin", 0, 0, mgl32.Vec2{0, 0}},
		{"plus col", 1, 0, mgl32.Vec2{TileWu / 2, TileHu / 2}},
		{"plus row", 0, 1, mgl32.Vec2{-TileWu / 2, TileHu / 2}},
		{"plus col and row", 1, 1, mgl32.Vec2{0, TileHu}},
		{"minus col and row", -1, -1, mgl32.Vec2{0, -TileHu}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := GridToScreen(c.c, c.r); !approxEq(got, c.want) {
				t.Errorf("GridToScreen(%v, %v) = %v, want %v", c.c, c.r, got, c.want)
			}
		})
	}
}

func TestScreenToGridRoundTrip(t *testing.T) {
	cases := []mgl32.Vec2{
		{3, 7},
		{-2, 5},
		{0, 0},
		{12.5, -4.25},
		{-9.9, -3.3},
	}

	for _, grid := range cases {
		screen := GridToScreen(grid.X(), grid.Y())
		got := ScreenToGrid(screen.X(), screen.Y())
		if !approxEq(got, grid) {
			t.Errorf("round-trip grid %v -> iso %v -> grid %v", grid, screen, got)
		}
	}
}

func TestToPixelsRoundTrip(t *testing.T) {
	iso := mgl32.Vec2{1.5, -0.25}

	px := ToPixels(iso)
	want := mgl32.Vec2{iso.X() * Unit, iso.Y() * Unit}
	if !approxEq(px, want) {
		t.Errorf("ToPixels(%v) = %v, want %v", iso, px, want)
	}

	back := ToUnits(px)
	if !approxEq(back, iso) {
		t.Errorf("ToUnits(ToPixels(%v)) = %v, want %v", iso, back, iso)
	}
}
