package fsm

import (
	"context"
	"game/pkg/fsm"
)

type StateIdentifier string

type State interface {
	fsm.State[StateIdentifier]

	Update(ctx context.Context) (Event, error)
}
