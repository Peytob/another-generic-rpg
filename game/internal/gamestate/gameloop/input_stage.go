package gameloop

import (
	"context"
	"game/internal/input"
)

// InputStage returns a Before stage that drains the collector into
// Frame.Input at the start of every frame.
func InputStage(c *input.InputCollector) Stage {
	return func(_ context.Context, f *Frame) error {
		f.Input = c.Drain()
		return nil
	}
}
