package fsm

import (
	"context"
	"game/pkg/engine"
)

const ResourcesLoadingStateIdentifier = engine.StateIdentifier("resources_loading")

type resourcesLoadingState struct {
}

func NewResourcesLoadingState() engine.State {
	return &resourcesLoadingState{}
}

func (s resourcesLoadingState) Update(_ context.Context) (engine.Event, error) {
	return NoEvent, nil
}

func (s resourcesLoadingState) Identifier() engine.StateIdentifier {
	return ResourcesLoadingStateIdentifier
}
