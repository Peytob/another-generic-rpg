package gameloop

import (
	"context"
)

// Simulate returns an EachTick stage that advances the gameplay
// simulation by one fixed tick.
func Simulate() Stage {
	return func(_ context.Context, _ *Frame) error {
		// todo advance world simulation by f.Dt applying f.ServerMessages
		return nil
	}
}
