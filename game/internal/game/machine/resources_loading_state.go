package machine

import (
	"context"
	"game/internal/game/state"
	"sync/atomic"
	"time"
)

const ResourcesLoadingStateIdentifier = StateIdentifier("resources_loading")

type resourcesLoadingState struct {
	loadingStage *int32
}

func NewResourcesLoadingState() State {
	return &resourcesLoadingState{
		loadingStage: new(int32),
	}
}

func (s resourcesLoadingState) Identifier() StateIdentifier {
	return ResourcesLoadingStateIdentifier
}

func (s resourcesLoadingState) OnEnter(ctx context.Context, _ *state.State) error {
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second * 2):
		}

		s.onLoad()
	}()

	return nil
}

func (s resourcesLoadingState) OnExit(_ context.Context, _ *state.State) error {
	return nil
}

func (s resourcesLoadingState) Update(_ context.Context, _ *state.State, _ time.Duration) (state.Event, error) {
	loadingStage := atomic.LoadInt32(s.loadingStage)

	returnEvent := state.NoEvent
	if loadingStage != 0 {
		returnEvent = state.ResourcesLoadedEvent
	}

	return returnEvent, nil
}

func (s resourcesLoadingState) onLoad() {
	atomic.StoreInt32(s.loadingStage, 1)
}
