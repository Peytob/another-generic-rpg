package loop

import (
	"context"
	"game/internal/game/state"
	"game/internal/input"
)

// GlobalInput returns a Before stage that maps frame input to state
// machine transitions.
func GlobalInput() Stage {
	return func(_ context.Context, f *Frame) error {
		if f.Input.Pressed(input.Exit) {
			f.RequestTransition(state.StoppedEvent)
		}
		return nil
	}
}
