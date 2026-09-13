package machine

import (
	"context"
	"game/internal/game/state"
	"time"
)

const ResourcesLoadingStateIdentifier = StateIdentifier("resources_loading")

type resourcesLoadingState struct{}

func NewResourcesLoadingState() State {
	return &resourcesLoadingState{}
}

func (s resourcesLoadingState) Identifier() StateIdentifier {
	return ResourcesLoadingStateIdentifier
}

func (s resourcesLoadingState) OnEnter(_ context.Context, _ *state.State) error {
	return nil
}

func (s resourcesLoadingState) OnExit(_ context.Context, _ *state.State) error {
	return nil
}

func (s resourcesLoadingState) Update(_ context.Context, _ *state.State, _ time.Duration) (state.Event, error) {
	return state.ResourcesLoadedEvent, nil
}
