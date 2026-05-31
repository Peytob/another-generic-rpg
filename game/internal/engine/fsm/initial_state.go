package fsm

import (
	"context"
	"game/pkg/engine"
)

const InitialStateIdentifier = engine.StateIdentifier("initial")

type initialState struct {
}

func NewInitialState() engine.State {
	return &initialState{}
}

func (s initialState) Update(_ context.Context) (engine.Event, error) {
	return NoEvent, nil
}

func (s initialState) Identifier() engine.StateIdentifier {
	return InitialStateIdentifier
}
