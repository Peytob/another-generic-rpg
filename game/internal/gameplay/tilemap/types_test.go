package tilemap

import (
	"errors"
	"testing"
)

func TestTile(t *testing.T) {
	t.Run("ID returns the tile id", func(t *testing.T) {
		tile := &Tile{id: "grass"}

		if got := tile.ID(); got != "grass" {
			t.Errorf("ID() = %q, want %q", got, "grass")
		}
	})
}

func TestLayer(t *testing.T) {
	tiles := []*Tile{
		{id: "0,0"}, {id: "1,0"},
		{id: "0,1"}, {id: "1,1"},
	}
	layer := Layer{tiles: tiles, width: 2, height: 2}

	t.Run("Tile returns tile by coordinates", func(t *testing.T) {
		cases := []struct {
			x, y int
			want *Tile
		}{
			{0, 0, tiles[0]},
			{1, 0, tiles[1]},
			{0, 1, tiles[2]},
			{1, 1, tiles[3]},
		}

		for _, c := range cases {
			got, err := layer.Tile(c.x, c.y)
			if err != nil {
				t.Fatalf("Tile(%d, %d) unexpected error: %v", c.x, c.y, err)
			}
			if got != c.want {
				t.Errorf("Tile(%d, %d) = %p, want %p", c.x, c.y, got, c.want)
			}
		}
	})

	t.Run("Tile out of bounds returns error", func(t *testing.T) {
		cases := []struct {
			name string
			x, y int
		}{
			{"x too large", 2, 0},
			{"y too large", 0, 2},
			{"negative x", -1, 0},
			{"negative y", 0, -1},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				got, err := layer.Tile(c.x, c.y)
				if !errors.Is(err, ErrOutOfRange) {
					t.Fatalf("Tile(%d, %d) expected ErrOutOfRange, got %v", c.x, c.y, err)
				}
				if got != nil {
					t.Errorf("Tile(%d, %d) expected nil tile, got %v", c.x, c.y, got)
				}
			})
		}
	})
}

func TestTilemap(t *testing.T) {
	layers := []Layer{
		{tiles: []*Tile{{id: "a"}}, width: 1, height: 1},
		{tiles: []*Tile{{id: "b"}}, width: 1, height: 1},
	}
	tm := Tilemap{id: "map", layers: layers, width: 1, height: 1}

	t.Run("Layer returns layer by index", func(t *testing.T) {
		got, err := tm.Layer(0)
		if err != nil {
			t.Fatalf("Layer(0) unexpected error: %v", err)
		}
		tile, err := got.Tile(0, 0)
		if err != nil {
			t.Fatalf("Tile(0, 0) unexpected error: %v", err)
		}
		if tile.ID() != "a" {
			t.Errorf("Layer(0) tile ID = %q, want %q", tile.ID(), "a")
		}

		got, err = tm.Layer(1)
		if err != nil {
			t.Fatalf("Layer(1) unexpected error: %v", err)
		}
		tile, err = got.Tile(0, 0)
		if err != nil {
			t.Fatalf("Tile(0, 0) unexpected error: %v", err)
		}
		if tile.ID() != "b" {
			t.Errorf("Layer(1) tile ID = %q, want %q", tile.ID(), "b")
		}
	})

	t.Run("Layer out of range returns error", func(t *testing.T) {
		cases := []struct {
			name string
			i    int
		}{
			{"too large", 2},
			{"negative", -1},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				_, err := tm.Layer(c.i)
				if !errors.Is(err, ErrOutOfRange) {
					t.Fatalf("Layer(%d) expected ErrOutOfRange, got %v", c.i, err)
				}
			})
		}
	})

	t.Run("LayersCount returns the number of layers", func(t *testing.T) {
		if got := tm.LayersCount(); got != 2 {
			t.Errorf("LayersCount() = %d, want 2", got)
		}
	})

	t.Run("Width and Height return dimensions", func(t *testing.T) {
		if got := tm.Width(); got != 1 {
			t.Errorf("Width() = %d, want 1", got)
		}
		if got := tm.Height(); got != 1 {
			t.Errorf("Height() = %d, want 1", got)
		}
	})
}
