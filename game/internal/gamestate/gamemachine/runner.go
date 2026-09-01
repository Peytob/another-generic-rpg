package gamemachine

import (
	"context"
	"fmt"
	"game/internal/gameplay/world"
	"game/internal/gamestate/event"
	"time"
)

// Runner connects a game state FSM with a gameplay World. On every state
// transition a fresh World is created, giving each state an isolated
// simulation context. The Runner ensures that state lifecycle hooks
// (OnEnter / OnExit) are called at the correct moments:
//
//   - Start creates the initial World and calls OnEnter on the initial state.
//   - Update calls the current state's Update; if the state returns an event,
//     the Runner performs the transition (OnExit -> FSM.Event -> new World ->
//     OnEnter).
//   - Event performs an externally triggered transition with the same
//     lifecycle guarantees.
//
// Runner is not safe for concurrent use; it must be called from a single
// goroutine (the main game loop).
type Runner struct {
	machine Machine
	world   *world.World
	started bool
}

// NewRunner creates a Runner. The World is created lazily on Start.
func NewRunner(machine Machine) *Runner {
	return &Runner{
		machine: machine,
	}
}

// NewRunnerWithWorld creates a Runner with a pre-configured initial World.
// The World is still recreated on every subsequent state transition.
func NewRunnerWithWorld(machine Machine, w *world.World) *Runner {
	return &Runner{
		machine: machine,
		world:   w,
	}
}

// Start activates the Runner by creating the initial World (unless one was
// provided via NewRunnerWithWorld) and calling OnEnter on the initial state.
// Calling Start more than once is a no-op.
func (r *Runner) Start(ctx context.Context) error {
	if r.started {
		return nil
	}
	r.started = true
	if r.world == nil {
		r.world = &world.World{}
	}
	return r.machine.State().OnEnter(ctx, r.world)
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

	ev, err := r.machine.State().Update(ctx, r.world, dt)
	if err != nil {
		return fmt.Errorf("state %q update: %w", r.machine.State().Identifier(), err)
	}

	if ev != event.NoEvent {
		return r.transition(ctx, ev)
	}

	return nil
}

// Event triggers an FSM transition from outside the state update loop
// (e.g. a window-close signal). It calls OnExit on the current state, performs
// the transition, creates a fresh World, and calls OnEnter on the new state.
func (r *Runner) Event(ctx context.Context, ev event.Event) error {
	if !r.started || !r.machine.IsRunning() {
		return r.machine.Event(ev)
	}
	return r.transition(ctx, ev)
}

// State returns the currently active state.
func (r *Runner) State() State {
	return r.machine.State()
}

// World returns the gameplay World for the currently active state.
func (r *Runner) World() *world.World {
	return r.world
}

// IsRunning reports whether the FSM has not yet reached a final state.
func (r *Runner) IsRunning() bool {
	return r.machine.IsRunning()
}

// transition performs the full lifecycle: OnExit -> FSM.Event -> new World ->
// OnEnter. The old World is discarded after OnExit completes.
func (r *Runner) transition(ctx context.Context, ev event.Event) error {
	current := r.machine.State()

	if err := current.OnExit(ctx, r.world); err != nil {
		return fmt.Errorf("state %q onExit: %w", current.Identifier(), err)
	}

	if err := r.machine.Event(ev); err != nil {
		return fmt.Errorf("fsm event %d: %w", ev, err)
	}

	r.world = &world.World{}

	if r.machine.IsRunning() {
		next := r.machine.State()
		if err := next.OnEnter(ctx, r.world); err != nil {
			return fmt.Errorf("state %q onEnter: %w", next.Identifier(), err)
		}
	}

	return nil
}
