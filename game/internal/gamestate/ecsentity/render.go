package ecsentity

import (
	"engine/graphic"
	"engine/utils/ecs"
	"game/internal/gamestate/ecscomponent"
	"game/internal/rendering"

	"github.com/go-gl/mathgl/mgl32"
)

func NewRenderStateEntity(w ecs.World, graphic graphic.Graphic) ecs.Entity {
	e := w.NewEntity()

	w.RegisterComponent(e, ecscomponent.RenderLayersComponent{
		TilemapCanvas: graphic.NewCanvas(),
	})

	w.RegisterComponent(e, ecscomponent.CameraComponent{
		Camera: rendering.NewCamera(mgl32.Vec2{0, 0}, mgl32.Vec2{800, 600}),
	})

	return e
}
