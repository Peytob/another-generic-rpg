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

func (i ResourcesLoadingState) OnEnter(_ context.Context) error {
	return nil
}

func (i ResourcesLoadingState) OnExit(_ context.Context) error {
	return nil
}

func (i ResourcesLoadingState) Tick(_ context.Context, _ Frame) (Event, error) {
	return ResourcesLoadedEvent, nil
}
