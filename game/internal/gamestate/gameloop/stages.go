package gameloop

import (
	"context"
	"game/internal/input"
	"game/internal/sync"
)

// InputStage returns a Before stage that drains the collector into
// Frame.Input at the start of every frame.
func InputStage[E comparable, A comparable](c *input.InputCollector[A]) Stage[E, A] {
	return func(_ context.Context, f *Frame[E, A]) error {
		f.Input = c.Drain()
		return nil
	}
}

// SyncStage returns a Before stage that drains pending server messages into
// Frame.ServerMessages for the simulation ticks to apply.
func SyncStage[E comparable, A comparable](s sync.Sync) Stage[E, A] {
	return func(_ context.Context, f *Frame[E, A]) error {
		f.ServerMessages = s.Drain()
		return nil
	}
}
