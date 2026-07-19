package gamestate

import (
	"context"
)

const StoppedStateIdentifier = StateIdentifier("stopped")

type stoppedState struct {
}

func NewStoppedState() State {
	return &stoppedState{}
}

func (s stoppedState) Update(_ context.Context) (Event, error) {
	return NoEvent, nil
}

func (s stoppedState) Identifier() StateIdentifier {
	return StoppedStateIdentifier
}
