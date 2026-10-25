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

func (i InitState) OnEnter(ctx context.Context) error {
	return nil
}

func (i InitState) OnExit(ctx context.Context) error {
	return nil
}

func (i InitState) Tick(ctx context.Context, frame Frame) (Event, error) {
	return InitializedEvent, nil
}
