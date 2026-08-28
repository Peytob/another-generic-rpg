package ecscomponent

import (
	"engine/utils/ecs"
	rtypes "game/internal/rendering/types"
)

//go:generate go run typeid
type Camera struct {
	Camera *rtypes.Camera
}

func (c Camera) Type() ecs.ComponentType {
	return ecs.ComponentType(CameraTypeID)
}

// CameraTypeID is the generated type ID (FNV-1a 64-bit) of
// game/internal/gamestate/ecscomponent.Camera. DO NOT EDIT.
const CameraTypeID int64 = 9177087747556155579
