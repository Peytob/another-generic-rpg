package ecsentity

import (
	"engine/utils/ecs"
	"game/internal/gameplay/tilemap"
	"game/internal/gamestate/ecscomponent"
)

func NewTilemapEntity(w ecs.World, tilemap *tilemap.Tilemap) ecs.Entity {
	e := w.NewEntity()

	w.RegisterComponent(e, ecscomponent.TilemapComponent{
		Tilemap: tilemap,
	})

	return e
}
