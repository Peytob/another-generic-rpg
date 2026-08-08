package ecs

import (
	"context"
	"iter"
	"reflect"
	"time"
)

// World is the central ECS container that manages entities, components,
// systems, and events.
type World interface {
	EntityManager
	ComponentManager
	Query
	SystemManager
	EventBus

	// Update advances the simulation by one tick, executing all registered
	// systems in registration order.
	Update(ctx context.Context, dt time.Duration) error
}

type world struct {
	nextEntityID         Entity
	entities             map[Entity]struct{}
	components           map[Entity][]Component
	componentsQueryIndex map[ComponentType][]Entity
	eventHandlers        map[EventType][]*eventSubscription
	systems              []System
}

func NewWorld() World {
	// todo memory allocation configuration
	return &world{
		nextEntityID:         1,
		entities:             make(map[Entity]struct{}, 32),
		components:           make(map[Entity][]Component, 128),
		componentsQueryIndex: make(map[ComponentType][]Entity, 64),
		eventHandlers:        make(map[EventType][]*eventSubscription, 16),
		systems:              make([]System, 0, 32),
	}
}

func (w *world) NewEntity() Entity {
	id := w.nextEntityID
	w.nextEntityID++
	w.entities[id] = struct{}{}
	return id
}

func (w *world) RemoveEntity(entity Entity) {
	delete(w.entities, entity)
	delete(w.components, entity)
	for componentType, entityList := range w.componentsQueryIndex {
		w.componentsQueryIndex[componentType] = filterOutEntity(entityList, entity)
	}
}

func (w *world) Entities() []Entity {
	entities := make([]Entity, 0, len(w.entities))
	for entity := range w.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (w *world) RegisterComponent(entity Entity, component Component) {
	if _, ok := w.entities[entity]; !ok {
		return
	}

	componentType := ComponentTypeOf(component)

	components := w.components[entity]
	for i, c := range components {
		if ComponentTypeOf(c) == componentType {
			components[i] = component
			return
		}
	}

	w.components[entity] = append(components, component)
	w.componentsQueryIndex[componentType] = append(w.componentsQueryIndex[componentType], entity)
}

func (w *world) UnregisterComponent(entity Entity, componentType ComponentType) {
	if _, ok := w.entities[entity]; !ok {
		return
	}

	components := w.components[entity]
	for i, c := range components {
		if ComponentTypeOf(c) == componentType {
			w.components[entity] = append(components[:i], components[i+1:]...)
			break
		}
	}
	if list, ok := w.componentsQueryIndex[componentType]; ok {
		w.componentsQueryIndex[componentType] = filterOutEntity(list, entity)
	}
}

func (w *world) GetComponent(entity Entity, componentType ComponentType) (Component, bool) {
	if _, ok := w.entities[entity]; !ok {
		return nil, false
	}

	for _, c := range w.components[entity] {
		if ComponentTypeOf(c) == componentType {
			return c, true
		}
	}
	return nil, false
}

func (w *world) HasComponent(entity Entity, componentType ComponentType) bool {
	_, ok := w.GetComponent(entity, componentType)
	return ok
}

func (w *world) Query(componentTypes ...ComponentType) []Entity {
	if len(componentTypes) == 0 {
		return nil
	}

	candidates := w.componentsQueryIndex[componentTypes[0]]
	for _, ct := range componentTypes[1:] {
		if list := w.componentsQueryIndex[ct]; len(list) < len(candidates) {
			candidates = list
		}
	}

	result := make([]Entity, 0, len(candidates))
	for _, entity := range candidates {
		if w.hasAllComponents(entity, componentTypes) {
			result = append(result, entity)
		}
	}
	return result
}

func (w *world) QueryOne(componentType ComponentType) []Entity {
	entities := w.componentsQueryIndex[componentType]
	result := make([]Entity, len(entities))
	copy(result, entities)
	return result
}

func (w *world) QuerySingle(componentTypes ...ComponentType) (Entity, error) {
	result := w.Query(componentTypes...)
	switch len(result) {
	case 0:
		return InvalidEntity, nil
	case 1:
		return result[0], nil
	default:
		return InvalidEntity, TooManyEntitiesFoundErr
	}
}

func (w *world) QuerySingleOne(componentType ComponentType) (Entity, error) {
	entities := w.componentsQueryIndex[componentType]
	switch len(entities) {
	case 0:
		return InvalidEntity, nil
	case 1:
		return entities[0], nil
	default:
		return InvalidEntity, TooManyEntitiesFoundErr
	}
}

func (w *world) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *world) Systems() []System {
	systems := make([]System, len(w.systems))
	copy(systems, w.systems)
	return systems
}

func (w *world) SystemsIter() iter.Seq[System] {
	return func(yield func(System) bool) {
		for _, system := range w.systems {
			if !yield(system) {
				return
			}
		}
	}
}

func (w *world) EmitEvent(ctx context.Context, event Event) error {
	eventType := EventType(reflect.TypeOf(event))
	subs := w.eventHandlers[eventType]

	hasInactive := false
	var firstErr error
	for _, sub := range subs {
		if !sub.active {
			hasInactive = true
			continue
		}
		if err := sub.handler(ctx, event); err != nil {
			firstErr = err
			break
		}
	}

	if hasInactive {
		w.compactSubscriptions(eventType)
	}
	return firstErr
}

func (w *world) Subscribe(eventType EventType, handler EventHandler) Subscription {
	sub := &eventSubscription{handler: handler, active: true}
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], sub)
	return sub
}

func (w *world) Update(ctx context.Context, dt time.Duration) error {
	for system := range w.SystemsIter() {
		if err := system.Execute(ctx, w, dt); err != nil {
			return err
		}
	}
	return nil
}

// compactSubscriptions rebuilds the handler slice for an event type without
// inactive (unsubscribed) entries. A fresh slice is allocated so the backing
// array of any in-flight EmitEvent range is never mutated, keeping iteration
// safe against reentrancy.
func (w *world) compactSubscriptions(eventType EventType) {
	subs := w.eventHandlers[eventType]
	filtered := make([]*eventSubscription, 0, len(subs))
	for _, sub := range subs {
		if sub.active {
			filtered = append(filtered, sub)
		}
	}
	if len(filtered) == 0 {
		delete(w.eventHandlers, eventType)
	} else {
		w.eventHandlers[eventType] = filtered
	}
}

func (w *world) hasAllComponents(entity Entity, componentTypes []ComponentType) bool {
	for _, ct := range componentTypes {
		if !w.HasComponent(entity, ct) {
			return false
		}
	}
	return true
}

func filterOutEntity(list []Entity, entity Entity) []Entity {
	filtered := list[:0]
	for _, e := range list {
		if e != entity {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
