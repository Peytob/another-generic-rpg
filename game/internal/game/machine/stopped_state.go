package machine

import (
	"context"
	"game/internal/game/state"
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

func (s stoppedState) OnEnter(_ context.Context, _ *state.State) error {
	return nil
}

func (s stoppedState) OnExit(_ context.Context, _ *state.State) error {
	return nil
}

func (s stoppedState) Update(_ context.Context, _ *state.State, _ time.Duration) (state.Event, error) {
	return state.NoEvent, nil
}
