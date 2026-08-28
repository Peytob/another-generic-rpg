package ecsentity

import (
	"engine/utils/ecs"
)

func NewMob(w ecs.World) ecs.Entity {
	e := w.NewEntity()

	return e
}
