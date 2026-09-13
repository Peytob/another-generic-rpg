package machine

import (
	"context"
	"engine/utils/fsm"
	"game/internal/game/state"
	"time"
)

type StateIdentifier string

type State interface {
	fsm.State[StateIdentifier]

	OnEnter(ctx context.Context, state *state.State) error
	OnExit(ctx context.Context, state *state.State) error
	Update(ctx context.Context, state *state.State, dt time.Duration) (state.Event, error)
}

type Machine = fsm.Machine[state.Event, StateIdentifier, State]
type Builder = fsm.MachineBuilder[state.Event, StateIdentifier, State]
type Transitions = fsm.Transitions[state.Event, StateIdentifier]

func NewBuilder() Builder {
	return fsm.NewBuilder[state.Event, StateIdentifier, State]()
}
