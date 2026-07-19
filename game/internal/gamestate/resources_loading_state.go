package gamestate

import (
	"context"
)

const ResourcesLoadingStateIdentifier = StateIdentifier("resources_loading")

type resourcesLoadingState struct {
}

func NewResourcesLoadingState() State {
	return &resourcesLoadingState{}
}

func (s resourcesLoadingState) Update(_ context.Context) (Event, error) {
	return NoEvent, nil
}

func (s resourcesLoadingState) Identifier() StateIdentifier {
	return ResourcesLoadingStateIdentifier
}
