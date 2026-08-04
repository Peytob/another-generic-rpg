package gamemachine

import (
	"context"
	"engine/utils/ecs"
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

func (s resourcesLoadingState) OnEnter(_ context.Context, _ ecs.World) error {
	return nil
}

func (s resourcesLoadingState) OnExit(_ context.Context, _ ecs.World) error {
	return nil
}

func (s resourcesLoadingState) Update(_ context.Context, _ ecs.World, _ time.Duration) (Event, error) {
	return NoEvent, nil
}
