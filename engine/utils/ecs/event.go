package ecs

import (
	"context"
	"reflect"
)

// Event is the base type for all events. Any type can serve as an event.
type Event any

// EventType identifies a specific component type via reflection
type EventType reflect.Type

// EventHandler processes a single event.
type EventHandler func(ctx context.Context, event Event) error

// Subscription represents an active event subscription that can be cancelled.
type Subscription interface {
	Unsubscribe()
}

// eventSubscription is a single handler registration. Unsubscribe marks it
// inactive (tombstone) rather than mutating the handler slice, which keeps
// EmitEvent iteration safe against reentrancy.
type eventSubscription struct {
	handler EventHandler
	active  bool
}

func (s *eventSubscription) Unsubscribe() {
	s.active = false
}

// EventBus provides publish/subscribe functionality for events.
type EventBus interface {
	// EmitEvent publishes an event to all subscribers of its concrete type.
	EmitEvent(ctx context.Context, event Event) error

	// Subscribe registers a handler for events of the specified EventType.
	// Returns a Subscription that can be used to unsubscribe.
	Subscribe(eventType EventType, handler EventHandler) Subscription
}

// Subscribe is a type-safe generic wrapper for EventBus.Subscribe.
//
// Usage:
//
//	sub := ecs.Subscribe[*DamageEvent](world, func(ctx context.Context, e *DamageEvent) error {
//	    // handle event
//	    return nil
//	})
//	defer sub.Unsubscribe()
func Subscribe[E any](bus EventBus, handler func(ctx context.Context, event E) error) Subscription {
	var sample E
	eventType := reflect.TypeOf(&sample).Elem()

	wrapped := func(ctx context.Context, event Event) error {
		return handler(ctx, event.(E))
	}

	return bus.Subscribe(eventType, wrapped)
}
