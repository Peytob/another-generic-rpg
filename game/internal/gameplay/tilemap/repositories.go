package tilemap

type TileRepository struct {
	idIndex map[string]*Tile
}

func NewTileRepository() *TileRepository {
	return &TileRepository{
		idIndex: make(map[string]*Tile),
	}
}

func (r *TileRepository) ByID(id string) (*Tile, bool) {
	tile, ok := r.idIndex[id]
	return tile, ok
}

func (r *TileRepository) Put(id string, tile *Tile) bool {
	if _, exists := r.idIndex[id]; exists {
		return false
	}

	r.idIndex[id] = tile
	return true
}

func (r *TileRepository) Delete(id string) bool {
	if _, exists := r.idIndex[id]; !exists {
		return false
	}

	delete(r.idIndex, id)
	return true
}

type Repository struct {
	idIndex map[string]*Tilemap
}

func NewRepository() *Repository {
	return &Repository{
		idIndex: make(map[string]*Tilemap),
	}
}

func (r *Repository) ByID(id string) (*Tilemap, bool) {
	tilemap, ok := r.idIndex[id]
	return tilemap, ok
}

func (r *Repository) Put(id string, tilemap *Tilemap) bool {
	if _, exists := r.idIndex[id]; exists {
		return false
	}

	r.idIndex[id] = tilemap
	return true
}

func (r *Repository) Delete(id string) bool {
	if _, exists := r.idIndex[id]; !exists {
		return false
	}

	delete(r.idIndex, id)
	return true
}
