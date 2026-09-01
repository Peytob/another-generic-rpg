package gameloop

import (
	"context"
	"game/internal/sync"
)

// SyncStage returns a Before stage that drains pending server messages into
// Frame.ServerMessages for the simulation ticks to apply.
func SyncStage(s sync.Sync) Stage {
	return func(_ context.Context, f *Frame) error {
		f.ServerMessages = s.Drain()
		return nil
	}
}
