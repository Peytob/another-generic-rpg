package ecscomponent

import (
	"engine/utils/ecs"
	"game/internal/rendering"
)

type CameraComponent struct {
	Camera *rendering.Camera
}

var CameraComponentType = ecs.ComponentTypeOfT[CameraComponent]()
