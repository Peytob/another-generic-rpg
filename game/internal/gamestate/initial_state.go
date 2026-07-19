package gamestate

import (
	"context"
)

const InitialStateIdentifier = StateIdentifier("initial")

type initialState struct {
}

func NewInitialState() State {
	return &initialState{}
}

func (s initialState) Update(_ context.Context) (Event, error) {
	return NoEvent, nil
}

func (s initialState) Identifier() StateIdentifier {
	return InitialStateIdentifier
}
