package ecscomponent

import (
	"engine/utils/ecs"
)

//go:generate go run typeid
type PlayerFlag struct {
}

func (p PlayerFlag) Type() ecs.ComponentType {
	return ecs.ComponentType(RenderLayersTypeID)
}

// PlayerFlagTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.PlayerFlag. DO NOT EDIT.
const PlayerFlagTypeID int64 = -9171615975225738789
