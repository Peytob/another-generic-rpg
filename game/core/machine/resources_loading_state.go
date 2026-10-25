package machine

import "context"

const ResourcesLoadingStateIdentifier = "resources_loading"

type ResourcesLoadingState struct {
}

func NewResourcesLoadingState() ResourcesLoadingState {
	return ResourcesLoadingState{}
}

func (i ResourcesLoadingState) Identifier() StateIdentifier {
	return ResourcesLoadingStateIdentifier
}

func (i ResourcesLoadingState) OnEnter(ctx context.Context) error {
	return nil
}

func (i ResourcesLoadingState) OnExit(ctx context.Context) error {
	return nil
}

func (i ResourcesLoadingState) Tick(ctx context.Context, frame Frame) (Event, error) {
	return ResourcesLoadedEvent, nil
}
