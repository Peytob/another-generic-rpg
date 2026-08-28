package ecsentity

import (
	"engine/utils/ecs"
	"game/internal/gamestate/ecscomponent"
)

func NewPlayer(w ecs.World) ecs.Entity {
	e := NewMob(w)

	w.RegisterComponent(e, ecscomponent.PlayerFlag{})

	return e
}
