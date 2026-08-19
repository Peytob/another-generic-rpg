package repositories

import "game/internal/gameplay/tilemap"

type TilemapRepository struct {
	idIndex map[string]*tilemap.Tilemap
}

func NewTilemapRepository() *TilemapRepository {
	return &TilemapRepository{
		idIndex: make(map[string]*tilemap.Tilemap),
	}
}

func (r *TilemapRepository) ByID(id string) (*tilemap.Tilemap, bool) {
	tilemap, ok := r.idIndex[id]
	return tilemap, ok
}

func (r *TilemapRepository) Put(id string, tilemap *tilemap.Tilemap) bool {
	if _, exists := r.idIndex[id]; exists {
		return false
	}

	r.idIndex[id] = tilemap
	return true
}

func (r *TilemapRepository) Delete(id string) bool {
	if _, exists := r.idIndex[id]; !exists {
		return false
	}

	delete(r.idIndex, id)
	return true
}
