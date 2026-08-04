package ecs

import (
	"errors"
	"reflect"
)

// Component is the base type constraint. Any type can serve as a component
type Component any

// ComponentType identifies a specific component type via reflection
type ComponentType reflect.Type

var TooManyEntitiesFound = errors.New("found too many entities for method")

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

	// Query returns all entities that have components of all the specified types
	Query(componentTypes ...ComponentType) []Entity

	// QueryOne returns all entities that have component of specified type
	QueryOne(componentType ComponentType) []Entity

	// QuerySingle returns one entity that have components of all the specified types. If there are more than one entity returns
	// TooManyEntitiesFound
	QuerySingle(componentTypes ...ComponentType) (Entity, error)

	// QuerySingleOne returns one entity that have component of specified type. If there are more than one entity returns
	// TooManyEntitiesFound
	QuerySingleOne(componentType ComponentType) (Entity, error)
}

// ComponentTypeOf returns the ComponentType for the given component value
func ComponentTypeOf(component Component) ComponentType {
	return ComponentType(reflect.TypeOf(component))
}

// ComponentTypeOfT generic implementation of ComponentTypeOf
func ComponentTypeOfT[T Component]() ComponentType {
	var t T
	return ComponentTypeOf(t)
}

// GetComponent is a type-safe generic accessor for components
func GetComponent[C any](manager ComponentManager, entity Entity) (C, bool) {
	var sample C
	ct := ComponentType(reflect.TypeOf(&sample).Elem())
	component, ok := manager.GetComponent(entity, ct)
	if !ok {
		var zero C
		return zero, false
	}
	return component.(C), true
}
