package ecscomponent

import (
	"engine/graphic/renderer"
	"engine/utils/ecs"
)

type RenderLayersComponent struct {
	TilemapCanvas renderer.Canvas
}

var RenderLayersComponentType = ecs.ComponentTypeOfT[RenderLayersComponent]()
