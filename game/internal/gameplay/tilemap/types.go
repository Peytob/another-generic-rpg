package tilemap

type Tile struct {
	id string
}

func (t Tile) ID() string {
	return t.id
}

type Layer struct {
	tiles  []*Tile
	width  int
	height int
}

func (l Layer) Tile(x, y int) (*Tile, error) {
	if x < 0 || y < 0 || x >= l.width || y >= l.height {
		return nil, ErrOutOfRange
	}

	return l.tiles[l.width*y+x], nil
}

type Tilemap struct {
	id     string
	layers []Layer
	width  int
	height int
}

func (t Tilemap) LayersCount() int {
	return len(t.layers)
}

func (t Tilemap) Layer(i int) (Layer, error) {
	if i < 0 || i >= len(t.layers) {
		return Layer{}, ErrOutOfRange
	}

	return t.layers[i], nil
}

func (t Tilemap) Width() int {
	return t.width
}

func (t Tilemap) Height() int {
	return t.height
}
