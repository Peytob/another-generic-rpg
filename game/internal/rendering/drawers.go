package rendering

import (
	"engine/graphic"
	"game/internal/gamestate/repositories"
	"game/internal/rendering/draw"
)

type Drawers struct {
	Tilemap draw.TilemapDrawer
}

func NewDrawers(_ graphic.Graphic, repo repositories.Repositories) *Drawers {
	return &Drawers{
		Tilemap: draw.NewTilemapDrawer(repo.TileRepository),
	}
}
