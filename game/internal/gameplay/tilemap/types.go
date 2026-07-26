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

func (l *Layer) Tile(x, y int) (*Tile, error) {
	if x < 0 || y < 0 || x >= l.width || y >= l.height {
		return nil, ErrOutOfRange
	}

	return l.tiles[l.width*y+x], nil
}

func (l *Layer) SetTile(x, y int, tile *Tile) error {
	if x < 0 || y < 0 || x >= l.width || y >= l.height {
		return ErrOutOfRange
	}

	l.tiles[l.width*y+x] = tile
	return nil
}

func (l *Layer) Width() int {
	return l.width
}

func (l *Layer) Height() int {
	return l.height
}

type Tilemap struct {
	id     string
	layers []Layer
	width  int
	height int
}

func NewTilemap(id string, layersCount, width, height int) (*Tilemap, error) {
	if layersCount < 0 || width < 0 || height < 0 {
		return nil, ErrInvalidDimensions
	}

	layers := make([]Layer, layersCount)
	for i := range layers {
		layers[i] = Layer{
			tiles:  make([]*Tile, width*height),
			width:  width,
			height: height,
		}
	}

	return &Tilemap{
		id:     id,
		layers: layers,
		width:  width,
		height: height,
	}, nil
}

func MustNewTilemap(id string, layersCount, width, height int) *Tilemap {
	tm, err := NewTilemap(id, layersCount, width, height)
	if err != nil {
		panic(err)
	}
	return tm
}

func (t *Tilemap) LayersCount() int {
	return len(t.layers)
}

func (t *Tilemap) Layer(i int) (*Layer, error) {
	if i < 0 || i >= len(t.layers) {
		return nil, ErrOutOfRange
	}

	return &t.layers[i], nil
}

func (t *Tilemap) Width() int {
	return t.width
}

func (t *Tilemap) Height() int {
	return t.height
}
