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

func (i StoppedState) OnEnter(ctx context.Context) error {
	return nil
}

func (i StoppedState) OnExit(ctx context.Context) error {
	return nil
}

func (i StoppedState) Tick(ctx context.Context, frame Frame) (Event, error) {
	return 0, nil
}
