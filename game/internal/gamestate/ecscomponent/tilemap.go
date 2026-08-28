package ecscomponent

import (
	"engine/utils/ecs"
	"game/internal/gameplay/tilemap"
)

//go:generate go run typeid
type Tilemap struct {
	Tilemap *tilemap.Tilemap
}

func (c Tilemap) Type() ecs.ComponentType {
	return ecs.ComponentType(TilemapTypeID)
}

// TilemapTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.Tilemap. DO NOT EDIT.
const TilemapTypeID int64 = 6744433690505514958
