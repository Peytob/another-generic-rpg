package state

import (
	"context"
	"game/internal/engine/fsm"
)

const InitialStateIdentifier = fsm.StateIdentifier("initial")

type initialState struct {
}

func NewInitialState() fsm.State {
	return &initialState{}
}

func (s initialState) Update(_ context.Context) (fsm.Event, error) {
	return fsm.NoEvent, nil
}

func (s initialState) Identifier() fsm.StateIdentifier {
	return InitialStateIdentifier
}
