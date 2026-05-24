package fsm

import "game/pkg/fsm"

// Sugar for engine FSM types

type Machine = fsm.Machine[Event, StateIdentifier, State]
type Builder = fsm.MachineBuilder[Event, StateIdentifier, State]
type Transitions = fsm.Transitions[Event, StateIdentifier]

func NewBuilder() Builder {
	return fsm.NewBuilder[Event, StateIdentifier, State]()
}

func NoTransitions() Transitions {
	return make(Transitions)
}
