package ecs

import (
	"context"
	"reflect"
)

// Event is the base type for all events. Any type can serve as an event.
type Event any

// EventHandler processes a single event.
type EventHandler func(ctx context.Context, event Event) error

// Subscription represents an active event subscription that can be cancelled.
type Subscription interface {
	Unsubscribe()
}

// EventBus provides publish/subscribe functionality for events.
type EventBus interface {
	// Emit publishes an event to all subscribers of its concrete type.
	Emit(ctx context.Context, event Event) error

	// Subscribe registers a handler for events of the specified reflect.Type.
	// Returns a Subscription that can be used to unsubscribe.
	Subscribe(eventType reflect.Type, handler EventHandler) Subscription
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
