package ecs

import (
	"context"
	"iter"
	"time"
)

// System processes entities matching a component filter each tick
type System interface {
	Execute(ctx context.Context, world World, dt time.Duration) error
}

// SystemManager manages system registration and execution order
type SystemManager interface {
	// AddSystem registers a system for processing
	AddSystem(system System)

	// Systems returns all registered systems in execution order
	Systems() []System

	// SystemsIter returns an iterator over registered systems in execution order.
	SystemsIter() iter.Seq[System]
}
