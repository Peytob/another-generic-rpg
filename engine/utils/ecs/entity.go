package ecs

// Entity is a unique identifier for an entity in the ECS world.
type Entity uint64

// InvalidEntity is the zero value of Entity and never refers to a real entity.
// Use it as a sentinel for "no entity".
const InvalidEntity Entity = 0

// EntityManager manages entity lifecycle.
type EntityManager interface {
	// NewEntity creates a new entity and returns its identifier.
	NewEntity() Entity

	// RemoveEntity destroys an entity and all its associated components.
	RemoveEntity(entity Entity)

	// Entities returns all active entities.
	Entities() []Entity
}
