package engine

import (
	"context"
)

// Engine Launcher for all inner-game logic
type Engine interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
