package engine

import (
	"context"
	"game/pkg/fsm"
)

// Event enumerated type to describe engine FSM events
type Event int

type StateIdentifier string
type State interface {
	fsm.State[StateIdentifier]

	Update(ctx context.Context) (Event, error)
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
