package rendering

import (
	"engine/graphic"
	"game/internal/rendering/draw"
)

type Drawers struct {
	Tilemap draw.TilemapDrawer
}

func NewDrawers(_ graphic.Graphic) *Drawers {
	return &Drawers{
		Tilemap: draw.NewTilemapDrawer(nil),
	}
}
