package gamemachine

import (
	"context"
	"engine/utils/fsm"
	"game/internal/gameplay/world"
	"time"
)

// Event enumerated type to describe engine FSM events
type Event int

type StateIdentifier string

// State represents a single game state in the FSM. Each state is tied to a
// gameplay World through lifecycle hooks: OnEnter sets up the world,
// Update advances the simulation, and OnExit tears down what OnEnter created.
type State interface {
	fsm.State[StateIdentifier]

	// OnEnter is called once when this state becomes the active state.
	// It receives the shared World so the state can set up gameplay
	// entities or subscribe to events.
	OnEnter(ctx context.Context, world *world.World) error

	// OnExit is called once when this state is about to be replaced by
	// another state. It receives the shared World so the state can
	// clean up (remove handlers, unsubscribe, etc.).
	OnExit(ctx context.Context, world *world.World) error

	// Update advances the state by one tick. The state should drive the
	// World update and may return an Event to trigger an FSM transition.
	// Returning NoEvent keeps the current state.
	Update(ctx context.Context, world *world.World, dt time.Duration) (Event, error)
}

// Machine Alias for engine specific FSM
type Machine = fsm.Machine[Event, StateIdentifier, State]

// MachineBuilder Alias for engine specific FSM builder
type MachineBuilder = fsm.MachineBuilder[Event, StateIdentifier, State]

// Transitions Alias for engine specific transitions
type Transitions = fsm.Transitions[Event, StateIdentifier]

func NewMachineBuilder() MachineBuilder {
	return fsm.NewBuilder[Event, StateIdentifier, State]()
}
