package machine

import (
	"context"
)

const InitStateIdentifier = "init"

type InitState struct{}

func NewInitState() InitState {
	return InitState{}
}

func (i InitState) Identifier() StateIdentifier {
	return InitStateIdentifier
}

func (i InitState) OnEnter(_ context.Context) error {
	return nil
}

func (i InitState) OnExit(_ context.Context) error {
	return nil
}

func (i InitState) Tick(_ context.Context, _ Frame) (Event, error) {
	return InitializedEvent, nil
}
