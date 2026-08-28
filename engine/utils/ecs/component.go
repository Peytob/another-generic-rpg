package ecs

import (
	"errors"
)

// Component is the base type constraint for components.
//
// Implementations must provide a Type method returning their unique
// ComponentType. Type must be cheap, deterministic, and must not
// dereference the receiver (a value receiver is recommended) so it can be
// called on zero values.
type Component interface {
	Type() ComponentType
}

// ComponentType identifies a specific component type.
//
// Values must be unique per component type; they are assigned by the
// component definition rather than allocated by the ecs package.
type ComponentType int64

var TooManyEntitiesFoundErr = errors.New("found too many entities for method")

// ComponentManager manages component storage, retrieval, and querying
type ComponentManager interface {
	// RegisterComponent associates a component with an entity,
	// replacing any existing component of the same type
	RegisterComponent(entity Entity, component Component)

	// UnregisterComponent removes the component of the specified type from the entity
	UnregisterComponent(entity Entity, componentType ComponentType)

	// GetComponent retrieves the component of the specified type from the entity
	// Returns the component and true if found, or nil and false otherwise
	GetComponent(entity Entity, componentType ComponentType) (Component, bool)

	// HasComponent reports whether the entity has a component of the specified type
	HasComponent(entity Entity, componentType ComponentType) bool
}

type Query interface {
	// Query returns all entities that have components of all the specified types
	Query(componentTypes ...ComponentType) []Entity

	// QueryOne returns all entities that have component of specified type
	QueryOne(componentType ComponentType) []Entity

	// QuerySingle returns one entity that have components of all the specified types. If there are more than one entity returns
	// TooManyEntitiesFoundErr. If no entities found returns InvalidEntity
	QuerySingle(componentTypes ...ComponentType) (Entity, error)

	// QuerySingleOne returns one entity that have component of specified type. If there are more than one entity returns
	// TooManyEntitiesFoundErr. If no entities found returns InvalidEntity
	QuerySingleOne(componentType ComponentType) (Entity, error)
}

// GetComponent is a type-safe generic accessor for components
func GetComponent[C Component](manager ComponentManager, entity Entity) (C, bool) {
	var zero C
	ct := zero.Type()
	component, ok := manager.GetComponent(entity, ct)
	if !ok {
		return zero, false
	}
	return component.(C), true
}

// GetSingleComponent is a type-safe generic accessor for single components
func GetSingleComponent[C Component](manager interface {
	Query
	ComponentManager
}) (C, bool) {
	var zero C
	ct := zero.Type()

	e, err := manager.QuerySingleOne(ct)
	if err != nil {
		// todo log error here
		return zero, false
	}
	if e == InvalidEntity {
		return zero, false
	}

	component, ok := manager.GetComponent(e, ct)
	if !ok {
		return zero, false
	}

	return component.(C), true
}
