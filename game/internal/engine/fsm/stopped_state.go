package fsm

import (
	"context"
	"game/pkg/engine"
)

const StoppedStateIdentifier = engine.StateIdentifier("stopped")

type stoppedState struct {
}

func NewStoppedState() engine.State {
	return &stoppedState{}
}

func (s stoppedState) Update(_ context.Context) (engine.Event, error) {
	return NoEvent, nil
}

func (s stoppedState) Identifier() engine.StateIdentifier {
	return StoppedStateIdentifier
}
