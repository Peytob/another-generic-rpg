package ecscomponent

import (
	"engine/utils/ecs"
	rtypes "game/internal/rendering/types"
)

type CameraComponent struct {
	Camera *rtypes.Camera
}

var CameraComponentType = ecs.ComponentTypeOfT[CameraComponent]()
