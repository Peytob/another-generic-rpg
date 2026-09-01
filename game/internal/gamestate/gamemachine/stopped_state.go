package gamemachine

import (
	"context"
	"game/internal/gameplay/world"
	"game/internal/gamestate/event"
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

func (s stoppedState) OnEnter(_ context.Context, _ *world.World) error {
	return nil
}

func (s stoppedState) OnExit(_ context.Context, _ *world.World) error {
	return nil
}

func (s stoppedState) Update(_ context.Context, _ *world.World, _ time.Duration) (event.Event, error) {
	return event.NoEvent, nil
}
