package gamemachine

import (
	"context"
	"engine/utils/ecs"
	"time"
)

const StoppedStateIdentifier = StateIdentifier("stopped")

type stoppedState struct{}

func NewStoppedState() State {
	return &stoppedState{}
}

func (s stoppedState) Identifier() StateIdentifier {
	return StoppedStateIdentifier
}

func (s stoppedState) OnEnter(_ context.Context, _ ecs.World) error {
	return nil
}

func (s stoppedState) OnExit(_ context.Context, _ ecs.World) error {
	return nil
}

func (s stoppedState) Update(_ context.Context, _ ecs.World, _ time.Duration) (Event, error) {
	return NoEvent, nil
}
