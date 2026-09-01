package gameloop

import (
	"context"
	"game/internal/rendering"
)

// WorldRender returns an After stage that draws the game world.
func WorldRender(r rendering.Rendering) Stage {
	return func(_ context.Context, _ *Frame) error {
		// todo draw scene via r.Drawers() using f.Alpha interpolation
		return nil
	}
}
