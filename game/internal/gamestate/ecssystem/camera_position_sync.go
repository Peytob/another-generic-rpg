package ecssystem

import (
	"context"
	"engine/utils/ecs"
	"game/internal/gamestate/ecscomponent"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraPositionSyncSystem struct {
}

func NewCameraPositionSyncSystem() CameraPositionSyncSystem {
	return CameraPositionSyncSystem{}
}

func (s CameraPositionSyncSystem) Execute(_ context.Context, world ecs.World, _ time.Duration) error {
	camera, ok := ecs.GetSingleComponent[ecscomponent.CameraComponent](world)
	if !ok {
		return nil
	}

	camera.Camera.Move(mgl32.Vec2{-1, 0})

	return nil
}
