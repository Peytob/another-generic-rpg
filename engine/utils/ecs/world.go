package ecs

import (
	"context"
	"time"
)

// World is the central ECS container that manages entities, components,
// systems, and events.
type World interface {
	EntityManager
	ComponentManager
	SystemManager
	EventBus

	// Update advances the simulation by one tick, executing all registered
	// systems in registration order.
	Update(ctx context.Context, dt time.Duration) error
}
