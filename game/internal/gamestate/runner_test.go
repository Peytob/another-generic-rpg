package gamestate

import (
	"context"
	"errors"
	"testing"
	"time"

	"engine/utils/ecs"
)

// --- test helpers ---------------------------------------------------------

// trackingState records lifecycle calls for assertions.
type trackingState struct {
	id        StateIdentifier
	enterCnt  int
	exitCnt   int
	updateCnt int
	exitErr   error
	enterErr  error
	event     Event
}

func (s *trackingState) Identifier() StateIdentifier { return s.id }

func (s *trackingState) OnEnter(_ context.Context, _ ecs.World) error {
	s.enterCnt++
	return s.enterErr
}

func (s *trackingState) OnExit(_ context.Context, _ ecs.World) error {
	s.exitCnt++
	return s.exitErr
}

func (s *trackingState) Update(_ context.Context, _ ecs.World, _ time.Duration) (Event, error) {
	s.updateCnt++
	return s.event, nil
}

var errUpdateSentinel = errors.New("update boom")

type errorUpdateState struct {
	id StateIdentifier
}

func (s *errorUpdateState) Identifier() StateIdentifier                  { return s.id }
func (s *errorUpdateState) OnEnter(_ context.Context, _ ecs.World) error { return nil }
func (s *errorUpdateState) OnExit(_ context.Context, _ ecs.World) error  { return nil }
func (s *errorUpdateState) Update(_ context.Context, _ ecs.World, _ time.Duration) (Event, error) {
	return NoEvent, errUpdateSentinel
}

func buildMachine(initial State, states ...State) Machine {
	b := NewMachineBuilder().
		InitialState(initial.Identifier()).
		RegisterState(initial, make(Transitions))

	for _, st := range states {
		b = b.RegisterState(st, make(Transitions))
	}

	return b.MustBuild()
}

// --- tests ----------------------------------------------------------------

func TestRunnerStartCallsOnEnter(t *testing.T) {
	t.Parallel()

	st := &trackingState{id: "s1"}
	m := buildMachine(st)

	r := NewRunner(m)

	if st.enterCnt != 0 {
		t.Fatalf("OnEnter should not be called before Start, got %d", st.enterCnt)
	}

	if err := r.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if st.enterCnt != 1 {
		t.Errorf("OnEnter should be called once, got %d", st.enterCnt)
	}
}

func TestRunnerStartIsIdempotent(t *testing.T) {
	t.Parallel()

	st := &trackingState{id: "s1"}
	m := buildMachine(st)
	r := NewRunner(m)

	_ = r.Start(context.Background())
	_ = r.Start(context.Background())

	if st.enterCnt != 1 {
		t.Errorf("OnEnter should be called once, got %d", st.enterCnt)
	}
}

func TestRunnerUpdateAutoStarts(t *testing.T) {
	t.Parallel()

	st := &trackingState{id: "s1"}
	m := buildMachine(st)
	r := NewRunner(m)

	if err := r.Update(context.Background(), time.Millisecond); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if st.enterCnt != 1 {
		t.Errorf("OnEnter should be called once, got %d", st.enterCnt)
	}
	if st.updateCnt != 1 {
		t.Errorf("Update should be called once, got %d", st.updateCnt)
	}
}

func TestRunnerUpdateTriggersTransition(t *testing.T) {
	t.Parallel()

	s1 := &trackingState{id: "s1", event: ResourcesLoadedEvent}
	s2 := &trackingState{id: "s2"}

	m := NewMachineBuilder().
		InitialState(s1.Identifier()).
		BuildState(s1).
		Transition(ResourcesLoadedEvent, s2.Identifier()).
		Build().
		BuildState(s2).
		Build().
		MustBuild()

	r := NewRunner(m)
	_ = r.Start(context.Background())

	// This Update should: call s1.Update (returns event), then
	// s1.OnExit, FSM transition, s2.OnEnter.
	if err := r.Update(context.Background(), time.Millisecond); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if s1.updateCnt != 1 {
		t.Errorf("s1.Update should be called once, got %d", s1.updateCnt)
	}
	if s1.exitCnt != 1 {
		t.Errorf("s1.OnExit should be called once, got %d", s1.exitCnt)
	}
	if s2.enterCnt != 1 {
		t.Errorf("s2.OnEnter should be called once, got %d", s2.enterCnt)
	}
	if r.State().Identifier() != s2.Identifier() {
		t.Errorf("current state should be s2")
	}
}

func TestRunnerRecreatesWorldOnTransition(t *testing.T) {
	t.Parallel()

	s1 := &trackingState{id: "s1", event: ResourcesLoadedEvent}
	s2 := &trackingState{id: "s2"}

	m := NewMachineBuilder().
		InitialState(s1.Identifier()).
		BuildState(s1).
		Transition(ResourcesLoadedEvent, s2.Identifier()).
		Build().
		BuildState(s2).
		Build().
		MustBuild()

	r := NewRunner(m)
	_ = r.Start(context.Background())

	worldBefore := r.World()
	worldBefore.NewEntity()

	_ = r.Update(context.Background(), time.Millisecond)

	worldAfter := r.World()
	if worldAfter == worldBefore {
		t.Fatal("World should be recreated on transition")
	}
	if got := len(worldAfter.Entities()); got != 0 {
		t.Errorf("new World should be empty, got %d entities", got)
	}
}

func TestRunnerEventTriggersTransition(t *testing.T) {
	t.Parallel()

	s1 := &trackingState{id: "s1"}
	s2 := &trackingState{id: "s2"}

	m := NewMachineBuilder().
		InitialState(s1.Identifier()).
		BuildState(s1).
		Transition(StoppedEvent, s2.Identifier()).
		Build().
		BuildState(s2).
		Build().
		FinalStates(s2.Identifier()).
		MustBuild()

	r := NewRunner(m)
	_ = r.Start(context.Background())

	if err := r.Event(context.Background(), StoppedEvent); err != nil {
		t.Fatalf("Event: %v", err)
	}

	if s1.exitCnt != 1 {
		t.Errorf("s1.OnExit should be called once, got %d", s1.exitCnt)
	}
	// s2 is a final state, so OnEnter is NOT called after reaching it.
	if s2.enterCnt != 0 {
		t.Errorf("s2.OnEnter should not be called for final state, got %d", s2.enterCnt)
	}
	if r.IsRunning() {
		t.Error("Runner should not be running after reaching final state")
	}
}

func TestRunnerEventOnFinalStateOmitsOnEnter(t *testing.T) {
	t.Parallel()

	s1 := &trackingState{id: "s1"}
	stopped := &trackingState{id: StoppedStateIdentifier}

	m := NewMachineBuilder().
		InitialState(s1.Identifier()).
		BuildState(s1).
		Build().
		RegisterState(stopped, make(Transitions)).
		GlobalTransitions(Transitions{
			StoppedEvent: stopped.Identifier(),
		}).
		FinalStates(stopped.Identifier()).
		MustBuild()

	r := NewRunner(m)
	_ = r.Start(context.Background())

	if err := r.Event(context.Background(), StoppedEvent); err != nil {
		t.Fatalf("Event: %v", err)
	}

	if stopped.enterCnt != 0 {
		t.Errorf("final state OnEnter should not be called, got %d", stopped.enterCnt)
	}
}

func TestRunnerOnExitErrorPropagates(t *testing.T) {
	t.Parallel()

	errSentinel := errors.New("exit boom")
	s1 := &trackingState{id: "s1", event: ResourcesLoadedEvent, exitErr: errSentinel}
	s2 := &trackingState{id: "s2"}

	m := NewMachineBuilder().
		InitialState(s1.Identifier()).
		BuildState(s1).
		Transition(ResourcesLoadedEvent, s2.Identifier()).
		Build().
		BuildState(s2).
		Build().
		MustBuild()

	r := NewRunner(m)
	_ = r.Start(context.Background())

	err := r.Update(context.Background(), time.Millisecond)
	if !errors.Is(err, errSentinel) {
		t.Errorf("expected errSentinel, got %v", err)
	}
	// Transition should not have completed
	if r.State().Identifier() != s1.Identifier() {
		t.Errorf("state should not have changed on OnExit error")
	}
}

func TestRunnerOnEnterErrorPropagates(t *testing.T) {
	t.Parallel()

	errSentinel := errors.New("enter boom")
	s1 := &trackingState{id: "s1", event: ResourcesLoadedEvent}
	s2 := &trackingState{id: "s2", enterErr: errSentinel}

	m := NewMachineBuilder().
		InitialState(s1.Identifier()).
		BuildState(s1).
		Transition(ResourcesLoadedEvent, s2.Identifier()).
		Build().
		BuildState(s2).
		Build().
		MustBuild()

	r := NewRunner(m)
	_ = r.Start(context.Background())

	err := r.Update(context.Background(), time.Millisecond)
	if !errors.Is(err, errSentinel) {
		t.Errorf("expected errSentinel, got %v", err)
	}
}

func TestRunnerWorldAccessible(t *testing.T) {
	t.Parallel()

	st := &trackingState{id: "s1"}
	m := buildMachine(st)
	r := NewRunner(m)
	_ = r.Start(context.Background())

	if r.World() == nil {
		t.Fatal("World should not be nil after Start")
	}

	e := r.World().NewEntity()
	if e == ecs.InvalidEntity {
		t.Fatal("World should create valid entities")
	}
}

func TestRunnerUpdateErrorPropagates(t *testing.T) {
	t.Parallel()

	errState := &errorUpdateState{id: "err"}
	m := buildMachine(errState)
	r := NewRunner(m)
	_ = r.Start(context.Background())

	err := r.Update(context.Background(), time.Millisecond)
	if !errors.Is(err, errUpdateSentinel) {
		t.Errorf("expected errUpdateSentinel, got %v", err)
	}
}

func TestRunnerMultipleUpdates(t *testing.T) {
	t.Parallel()

	st := &trackingState{id: "s1"}
	m := buildMachine(st)
	r := NewRunner(m)
	_ = r.Start(context.Background())

	for i := 0; i < 5; i++ {
		if err := r.Update(context.Background(), time.Millisecond); err != nil {
			t.Fatalf("Update %d: %v", i, err)
		}
	}

	if st.updateCnt != 5 {
		t.Errorf("Update should be called 5 times, got %d", st.updateCnt)
	}
	if st.enterCnt != 1 {
		t.Errorf("OnEnter should still be called once, got %d", st.enterCnt)
	}
}
