package ecscomponent

import (
	"engine/utils/ecs"
	"game/internal/gameplay/tilemap"
)

type TilemapComponent struct {
	Tilemap *tilemap.Tilemap
}

var TilemapComponentType = ecs.ComponentTypeOfT[TilemapComponent]()
