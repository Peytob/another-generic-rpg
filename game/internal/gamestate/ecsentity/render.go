package ecsentity

import (
	"engine/graphic/renderer"
	"engine/utils/ecs"
	"game/internal/gamestate/ecscomponent"
	rtypes "game/internal/rendering/types"

	"github.com/go-gl/mathgl/mgl32"
)

func NewRenderState(w ecs.World, tilemapCanvas renderer.Canvas) ecs.Entity {
	e := w.NewEntity()

	w.RegisterComponent(e, ecscomponent.RenderLayers{
		TilemapCanvas: tilemapCanvas,
	})

	w.RegisterComponent(e, ecscomponent.Camera{
		Camera: rtypes.NewCamera(mgl32.Vec2{0, 0}, mgl32.Vec2{800, 600}),
	})

	return e
}
