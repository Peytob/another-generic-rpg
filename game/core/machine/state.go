package machine

import (
	"context"
	"game/pkg/fsm"
)

type StateIdentifier string

type State interface {
	fsm.State[StateIdentifier]

	OnEnter(ctx context.Context) error
	OnExit(ctx context.Context) error
	Tick(ctx context.Context, frame Frame) (Event, error)
}

type Machine = fsm.Machine[Event, StateIdentifier, State]
type Builder = fsm.MachineBuilder[Event, StateIdentifier, State]
type Transitions = fsm.Transitions[Event, StateIdentifier]

func NewBuilder() Builder {
	return fsm.NewBuilder[Event, StateIdentifier, State]()
}
