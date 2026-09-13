package machine

import (
	"context"
	"fmt"
	"game/internal/game/repositories"
	"game/internal/game/state"
	"game/internal/gameplay/world"
	"time"
)

type Runner struct {
	machine      Machine
	repositories repositories.Repositories

	state   *state.State
	started bool
}

// NewRunner creates a Runner. The World is created lazily on Start.
func NewRunner(machine Machine, repositories repositories.Repositories) *Runner {
	return &Runner{
		machine:      machine,
		repositories: repositories,
	}
}

func (r *Runner) Start(ctx context.Context) error {
	if r.started {
		return nil
	}

	r.started = true

	if r.state == nil {
		r.state = &state.State{
			World:        &world.World{},
			Repositories: r.repositories,
		}
	}

	return r.machine.State().OnEnter(ctx, r.state)
}

// Update advances the simulation by one tick. If the current state returns a
// non-NoEvent event, the transition is performed automatically.
func (r *Runner) Update(ctx context.Context, dt time.Duration) error {
	if !r.started {
		if err := r.Start(ctx); err != nil {
			return err
		}
	}

	if !r.machine.IsRunning() {
		return nil
	}

	ev, err := r.machine.State().Update(ctx, r.state, dt)
	if err != nil {
		return fmt.Errorf("state %q update: %w", r.machine.State().Identifier(), err)
	}

	if ev != state.NoEvent {
		return r.transition(ctx, ev)
	}

	return nil
}

// Event triggers an FSM transition from outside the state update loop
// (e.g. a window-close signal). It calls OnExit on the current state, performs
// the transition, creates a fresh World, and calls OnEnter on the new state.
func (r *Runner) Event(ctx context.Context, ev state.Event) error {
	if !r.started || !r.machine.IsRunning() {
		return r.machine.Event(ev)
	}
	return r.transition(ctx, ev)
}

// State returns the currently active state.
func (r *Runner) State() State {
	return r.machine.State()
}

// IsRunning reports whether the FSM has not yet reached a final state.
func (r *Runner) IsRunning() bool {
	return r.machine.IsRunning()
}

// transition performs the full lifecycle: OnExit -> FSM.Event -> new World ->
// OnEnter. The old World is discarded after OnExit completes.
func (r *Runner) transition(ctx context.Context, ev state.Event) error {
	current := r.machine.State()

	if err := current.OnExit(ctx, r.state); err != nil {
		return fmt.Errorf("state %q onExit: %w", current.Identifier(), err)
	}

	if err := r.machine.Event(ev); err != nil {
		return fmt.Errorf("fsm event %d: %w", ev, err)
	}

	r.state.World = &world.World{}

	if r.machine.IsRunning() {
		next := r.machine.State()
		if err := next.OnEnter(ctx, r.state); err != nil {
			return fmt.Errorf("state %q onEnter: %w", next.Identifier(), err)
		}
	}

	return nil
}
