package gamemachine

import (
	"context"
	"engine/utils/ecs"
	"time"
)

const InitialStateIdentifier = StateIdentifier("initial")

type initialState struct{}

func NewInitialState() State {
	return &initialState{}
}

func (s initialState) Identifier() StateIdentifier {
	return InitialStateIdentifier
}

func (s initialState) OnEnter(_ context.Context, _ ecs.World) error {
	return nil
}

func (s initialState) OnExit(_ context.Context, _ ecs.World) error {
	return nil
}

func (s initialState) Update(_ context.Context, _ ecs.World, _ time.Duration) (Event, error) {
	return NoEvent, nil
}
