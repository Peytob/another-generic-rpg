package repositories

import "game/internal/gameplay/tilemap"

type TileRepository struct {
	idIndex map[string]*tilemap.Tile
}

func NewTileRepository() *TileRepository {
	return &TileRepository{
		idIndex: make(map[string]*tilemap.Tile),
	}
}

func (r *TileRepository) ByID(id string) (*tilemap.Tile, bool) {
	tile, ok := r.idIndex[id]
	return tile, ok
}

func (r *TileRepository) Put(id string, tile *tilemap.Tile) bool {
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
