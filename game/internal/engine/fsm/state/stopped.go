package state

import (
	"context"
	"game/internal/engine/fsm"
)

const StoppedStateIdentifier = fsm.StateIdentifier("stopped")

type stoppedState struct {
}

func NewStoppedState() fsm.State {
	return &stoppedState{}
}

func (s stoppedState) Update(_ context.Context) (fsm.Event, error) {
	return fsm.NoEvent, nil
}

func (s stoppedState) Identifier() fsm.StateIdentifier {
	return StoppedStateIdentifier
}
