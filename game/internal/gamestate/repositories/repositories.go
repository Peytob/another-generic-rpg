package repositories

type Repositories struct {
	TileRepository    *TileRepository
	TilemapRepository *TilemapRepository
}

func NewRepositories() Repositories {
	return Repositories{
		TileRepository:    NewTileRepository(),
		TilemapRepository: NewTilemapRepository(),
	}
}
