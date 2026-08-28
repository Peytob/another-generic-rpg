package ecsentity

import (
	"engine/utils/ecs"
	"game/internal/gameplay/tilemap"
	"game/internal/gamestate/ecscomponent"
)

func NewTilemap(w ecs.World, tilemap *tilemap.Tilemap) ecs.Entity {
	e := w.NewEntity()

	w.RegisterComponent(e, ecscomponent.Tilemap{
		Tilemap: tilemap,
	})

	return e
}
