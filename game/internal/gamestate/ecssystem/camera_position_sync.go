package ecssystem

import (
	"context"
	"engine/utils/ecs"
	"game/internal/gamestate/ecscomponent"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraPositionSync struct {
}

func NewCameraPositionSync() CameraPositionSync {
	return CameraPositionSync{}
}

func (s CameraPositionSync) Execute(_ context.Context, world ecs.World, _ time.Duration) error {
	camera, ok := ecs.GetSingleComponent[ecscomponent.Camera](world)
	if !ok {
		return nil
	}

	// TODO заменить на синк с позицией игрока
	camera.Camera.Move(mgl32.Vec2{-1, 0})

	return nil
}
