package tilemap

import (
	"errors"
	"testing"
)

const (
	grassTileID = "grass"
	stoneTileID = "stone"
)

func TestTile(t *testing.T) {
	t.Run("ID returns the tile id", func(t *testing.T) {
		tile := &Tile{id: grassTileID}

		if got := tile.ID(); got != grassTileID {
			t.Errorf("ID() = %q, want %q", got, grassTileID)
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

	t.Run("SetTile replaces tile at coordinates", func(t *testing.T) {
		layer := Layer{
			tiles:  []*Tile{{id: stoneTileID}, {id: stoneTileID}, {id: stoneTileID}, {id: stoneTileID}},
			width:  2,
			height: 2,
		}
		newTile := &Tile{id: "new"}

		if err := layer.SetTile(1, 1, newTile); err != nil {
			t.Fatalf("SetTile(1, 1) unexpected error: %v", err)
		}

		got, err := layer.Tile(1, 1)
		if err != nil {
			t.Fatalf("Tile(1, 1) unexpected error: %v", err)
		}
		if got != newTile {
			t.Errorf("SetTile did not replace tile, got %p, want %p", got, newTile)
		}
	})

	t.Run("SetTile out of bounds returns error", func(t *testing.T) {
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
				err := layer.SetTile(c.x, c.y, &Tile{id: "x"})
				if !errors.Is(err, ErrOutOfRange) {
					t.Fatalf("SetTile(%d, %d) expected ErrOutOfRange, got %v", c.x, c.y, err)
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

	t.Run("Layer returns pointer so SetTile persists on tilemap", func(t *testing.T) {
		newTile := &Tile{id: "replaced"}

		layer, err := tm.Layer(0)
		if err != nil {
			t.Fatalf("Layer(0) unexpected error: %v", err)
		}
		if err := layer.SetTile(0, 0, newTile); err != nil {
			t.Fatalf("SetTile(0, 0) unexpected error: %v", err)
		}

		layer, err = tm.Layer(0)
		if err != nil {
			t.Fatalf("Layer(0) unexpected error: %v", err)
		}
		got, err := layer.Tile(0, 0)
		if err != nil {
			t.Fatalf("Tile(0, 0) unexpected error: %v", err)
		}
		if got != newTile {
			t.Errorf("SetTile did not persist on tilemap, got %p, want %p", got, newTile)
		}
	})
}

func TestNewTilemap(t *testing.T) {
	t.Run("initializes layers and tile arrays", func(t *testing.T) {
		tm, err := NewTilemap("map", 3, 4, 2)
		if err != nil {
			t.Fatalf("NewTilemap unexpected error: %v", err)
		}

		if tm.LayersCount() != 3 {
			t.Fatalf("LayersCount() = %d, want 3", tm.LayersCount())
		}
		if tm.Width() != 4 {
			t.Errorf("Width() = %d, want 4", tm.Width())
		}
		if tm.Height() != 2 {
			t.Errorf("Height() = %d, want 2", tm.Height())
		}

		for i := range tm.LayersCount() {
			layer, err := tm.Layer(i)
			if err != nil {
				t.Fatalf("Layer(%d) unexpected error: %v", i, err)
			}
			if layer.Width() != 4 {
				t.Errorf("layer %d Width() = %d, want 4", i, layer.Width())
			}
			if layer.Height() != 2 {
				t.Errorf("layer %d Height() = %d, want 2", i, layer.Height())
			}
			for y := range 2 {
				for x := range 4 {
					got, err := layer.Tile(x, y)
					if err != nil {
						t.Fatalf("Layer(%d) Tile(%d, %d) unexpected error: %v", i, x, y, err)
					}
					if got != nil {
						t.Errorf("Layer(%d) Tile(%d, %d) = %v, want nil", i, x, y, got)
					}
				}
			}
		}
	})

	t.Run("SetTile works on constructed tilemap", func(t *testing.T) {
		tm, err := NewTilemap("map", 1, 2, 2)
		if err != nil {
			t.Fatalf("NewTilemap unexpected error: %v", err)
		}
		tile := &Tile{id: "grass"}

		layer, err := tm.Layer(0)
		if err != nil {
			t.Fatalf("Layer(0) unexpected error: %v", err)
		}
		if err := layer.SetTile(1, 1, tile); err != nil {
			t.Fatalf("SetTile(1, 1) unexpected error: %v", err)
		}

		got, err := layer.Tile(1, 1)
		if err != nil {
			t.Fatalf("Tile(1, 1) unexpected error: %v", err)
		}
		if got != tile {
			t.Errorf("Tile(1, 1) = %p, want %p", got, tile)
		}
	})

	t.Run("negative dimensions return ErrInvalidDimensions", func(t *testing.T) {
		cases := []struct {
			name                       string
			layersCount, width, height int
		}{
			{"negative layersCount", -1, 2, 2},
			{"negative width", 1, -1, 2},
			{"negative height", 1, 2, -1},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				tm, err := NewTilemap("map", c.layersCount, c.width, c.height)
				if !errors.Is(err, ErrInvalidDimensions) {
					t.Fatalf("NewTilemap(%d, %d, %d) expected ErrInvalidDimensions, got %v", c.layersCount, c.width, c.height, err)
				}
				if tm != nil {
					t.Errorf("NewTilemap(%d, %d, %d) expected nil tilemap", c.layersCount, c.width, c.height)
				}
			})
		}
	})
}

func TestMustNewTilemap(t *testing.T) {
	t.Run("returns tilemap on valid input", func(t *testing.T) {
		tm := MustNewTilemap("map", 2, 3, 3)
		if tm.LayersCount() != 2 {
			t.Errorf("LayersCount() = %d, want 2", tm.LayersCount())
		}
		if tm.Width() != 3 {
			t.Errorf("Width() = %d, want 3", tm.Width())
		}
		if tm.Height() != 3 {
			t.Errorf("Height() = %d, want 3", tm.Height())
		}
	})

	t.Run("panics on invalid input", func(t *testing.T) {
		cases := []struct {
			name                       string
			layersCount, width, height int
		}{
			{"negative layersCount", -1, 2, 2},
			{"negative width", 1, -1, 2},
			{"negative height", 1, 2, -1},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("MustNewTilemap(%d, %d, %d) expected panic", c.layersCount, c.width, c.height)
					}
				}()
				MustNewTilemap("map", c.layersCount, c.width, c.height)
			})
		}
	})
}
