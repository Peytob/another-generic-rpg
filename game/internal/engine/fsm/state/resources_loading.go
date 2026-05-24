package state

import (
	"context"
	"game/internal/engine/fsm"
)

const ResourcesLoadingStateIdentifier = fsm.StateIdentifier("resources_loading")

type resourcesLoadingState struct {
}

func NewResourcesLoadingState() fsm.State {
	return &resourcesLoadingState{}
}

func (s resourcesLoadingState) Update(_ context.Context) (fsm.Event, error) {
	return fsm.NoEvent, nil
}

func (s resourcesLoadingState) Identifier() fsm.StateIdentifier {
	return ResourcesLoadingStateIdentifier
}
