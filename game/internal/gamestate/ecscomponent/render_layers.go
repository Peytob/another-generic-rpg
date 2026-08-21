package ecscomponent

import (
	"engine/graphic/renderer"
)

//go:generate go run typeid
type RenderLayersComponent struct {
	TilemapCanvas renderer.Canvas
}

// RenderLayersComponentTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.RenderLayersComponent. DO NOT EDIT.
const RenderLayersComponentTypeID int64 = -8825906260855309331
