package ecs

// Entity is a unique identifier for an entity in the ECS world.
type Entity uint64

// EntityManager manages entity lifecycle.
type EntityManager interface {
	// NewEntity creates a new entity and returns its identifier.
	NewEntity() Entity

	// RemoveEntity destroys an entity and all its associated components.
	RemoveEntity(entity Entity)

	// Entities returns all active entities.
	Entities() []Entity
}
