package gameloop

import (
	"context"
	"game/internal/gamestate/event"
	"game/internal/input"
)

// GlobalInput returns a Before stage that maps frame input to state
// machine transitions.
func GlobalInput() Stage {
	return func(_ context.Context, f *Frame) error {
		if f.Input.Pressed(input.Exit) {
			f.RequestTransition(event.StoppedEvent)
		}
		return nil
	}
}
