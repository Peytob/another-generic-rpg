package gamemachine

import (
	"context"
	"game/internal/gameplay/world"
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

func (s resourcesLoadingState) OnEnter(_ context.Context, _ *world.World) error {
	return nil
}

func (s resourcesLoadingState) OnExit(_ context.Context, _ *world.World) error {
	return nil
}

func (s resourcesLoadingState) Update(_ context.Context, _ *world.World, _ time.Duration) (Event, error) {
	return NoEvent, nil
}
