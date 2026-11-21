package machine

import (
	"context"
)

const StoppedStateIdentifier = "stopped"

type StoppedState struct {
}

func NewStoppedState() StoppedState {
	return StoppedState{}
}

func (i StoppedState) Identifier() StateIdentifier {
	return StoppedStateIdentifier
}

func (i StoppedState) OnEnter(_ context.Context) error {
	return nil
}

func (i StoppedState) OnExit(_ context.Context) error {
	return nil
}

func (i StoppedState) Tick(_ context.Context, _ Frame) (Event, error) {
	return 0, nil
}
