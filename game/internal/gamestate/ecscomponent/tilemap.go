package ecscomponent

import (
	"game/internal/gameplay/tilemap"
)

//go:generate go run typeid
type TilemapComponent struct {
	Tilemap *tilemap.Tilemap
}

// TilemapComponentTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.TilemapComponent. DO NOT EDIT.
const TilemapComponentTypeID int64 = 5452769344686206493
