package ecs

import (
	"context"
	"time"
)

// System processes entities matching a component filter each tick.
type System interface {
	// Update processes all matching entities for this tick.
	Update(ctx context.Context, world World, dt time.Duration) error
}

// SystemManager manages system registration and execution order.
type SystemManager interface {
	// AddSystem registers a system for processing.
	AddSystem(system System)

	// RemoveSystem unregisters a system.
	RemoveSystem(system System)

	// Systems returns all registered systems in execution order.
	Systems() []System
}
