package ecscomponent

import (
	"engine/graphic/renderer"
	"engine/utils/ecs"
)

//go:generate go run typeid
type RenderLayers struct {
	TilemapCanvas renderer.Canvas
}

func (c RenderLayers) Type() ecs.ComponentType {
	return ecs.ComponentType(RenderLayersTypeID)
}

// RenderLayersTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.RenderLayers. DO NOT EDIT.
const RenderLayersTypeID int64 = -3261046247599486594
