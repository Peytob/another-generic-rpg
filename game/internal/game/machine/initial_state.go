package machine

import (
	"context"
	"game/internal/game/state"
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

func (s initialState) OnEnter(_ context.Context, _ *state.State) error {
	return nil
}

func (s initialState) OnExit(_ context.Context, _ *state.State) error {
	return nil
}

func (s initialState) Update(_ context.Context, _ *state.State, _ time.Duration) (state.Event, error) {
	return state.NoEvent, nil
}
