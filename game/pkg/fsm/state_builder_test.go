package fsm

import (
	"testing"
)

func TestStateBuilder_Transition(t *testing.T) {
	t.Parallel()

	t.Run("should add transition to state", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Build()

		machine, err := m.
			RegisterState(runningState, Transitions[string, string]{}).
			InitialState(initialState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := machine.Event(initializedEvent); err != nil {
			t.Errorf("unexpected error on event: %v", err)
		}

		if machine.State().Identifier() != runningState.Identifier() {
			t.Errorf("expected state %s, got %s", runningState.Identifier(), machine.State().Identifier())
		}
	})

	t.Run("should overwrite transition when event already exists", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			Transition(initializedEvent, loadingState.Identifier()).
			Transition(initializedEvent, runningState.Identifier()).
			Build()

		machine, err := m.
			RegisterState(runningState, Transitions[string, string]{}).
			RegisterState(loadingState, Transitions[string, string]{}).
			InitialState(initialState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := machine.Event(initializedEvent); err != nil {
			t.Errorf("unexpected error on event: %v", err)
		}

		if machine.State().Identifier() != runningState.Identifier() {
			t.Errorf("expected overwritten transition to %s, got %s",
				runningState.Identifier(), machine.State().Identifier())
		}
	})

	t.Run("should preserve all transitions added via Transition calls", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Transition(failedEvent, failedState.Identifier()).
			Build()

		mb := m.(*machineBuilder[string, string, TestState])

		transitions, ok := mb.transitions[initialState.Identifier()]
		if !ok {
			t.Fatal("state not found in builder transitions")
		}

		if target, ok := transitions[initializedEvent]; !ok || target != runningState.Identifier() {
			t.Errorf("initializedEvent should transition to %s", runningState.Identifier())
		}

		if target, ok := transitions[failedEvent]; !ok || target != failedState.Identifier() {
			t.Errorf("failedEvent should transition to %s", failedState.Identifier())
		}
	})

	t.Run("TransitionS should add transition using State directly", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			TransitionS(initializedEvent, runningState).
			Build()

		machine, err := m.
			RegisterState(runningState, Transitions[string, string]{}).
			InitialState(initialState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := machine.Event(initializedEvent); err != nil {
			t.Errorf("unexpected error on event: %v", err)
		}

		if machine.State().Identifier() != runningState.Identifier() {
			t.Errorf("expected state %s, got %s", runningState.Identifier(), machine.State().Identifier())
		}
	})

	t.Run("TransitionS should be equivalent to Transition with Identifier", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			TransitionS(initializedEvent, runningState).
			TransitionS(exitedEvent, exitedState).
			Build()

		mb := m.(*machineBuilder[string, string, TestState])

		transitions, ok := mb.transitions[initialState.Identifier()]
		if !ok {
			t.Fatal("state not found in builder transitions")
		}

		if target, ok := transitions[initializedEvent]; !ok || target != runningState.Identifier() {
			t.Errorf("initializedEvent should transition to %s", runningState.Identifier())
		}

		if target, ok := transitions[exitedEvent]; !ok || target != exitedState.Identifier() {
			t.Errorf("exitedEvent should transition to %s", exitedState.Identifier())
		}
	})
}

func TestStateBuilder_Build(t *testing.T) {
	t.Parallel()

	t.Run("should return the parent machine builder", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		result := m.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Build()

		if result != m {
			t.Error("Build should return the parent machine builder")
		}
	})

	t.Run("should register state and transitions in parent builder via RegisterState", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		// Build same state twice with different transitions - RegisterState should merge them
		m.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Build()

		m.BuildState(initialState).
			Transition(failedEvent, failedState.Identifier()).
			Build()

		mb := m.(*machineBuilder[string, string, TestState])

		transitions, ok := mb.transitions[initialState.Identifier()]
		if !ok {
			t.Fatal("state not found in builder transitions")
		}

		if _, ok := transitions[initializedEvent]; !ok {
			t.Error("initializedEvent transition should be present after merge")
		}

		if _, ok := transitions[failedEvent]; !ok {
			t.Error("failedEvent transition should be present after merge")
		}
	})
}

func TestMachineBuilder_BuildState(t *testing.T) {
	t.Parallel()

	t.Run("should create state builder for given state", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Build()

		machine, err := m.
			RegisterState(runningState, Transitions[string, string]{}).
			InitialState(initialState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := machine.Event(initializedEvent); err != nil {
			t.Errorf("unexpected error on event: %v", err)
		}

		if machine.State().Identifier() != runningState.Identifier() {
			t.Errorf("expected state %s, got %s",
				runningState.Identifier(), machine.State().Identifier())
		}
	})

	t.Run("should build full machine using only builder pattern", func(t *testing.T) {
		t.Parallel()

		m := NewBuilder[string, string, TestState]()

		m.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Build()

		m.BuildState(runningState).
			Transition(exitedEvent, exitedState.Identifier()).
			Transition(levelChangedEvent, loadingState.Identifier()).
			Build()

		m.BuildState(loadingState).
			Transition(loadedEvent, runningState.Identifier()).
			Build()

		m.BuildState(exitedState).Build()

		machine, err := m.
			InitialState(initialState.Identifier()).
			FinalStates(exitedState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := machine.Event(initializedEvent); err != nil {
			t.Errorf("unexpected error on initialized event: %v", err)
		}

		if machine.State().Identifier() != runningState.Identifier() {
			t.Errorf("expected state %s, got %s",
				runningState.Identifier(), machine.State().Identifier())
		}

		if err := machine.Event(levelChangedEvent); err != nil {
			t.Errorf("unexpected error on levelChanged event: %v", err)
		}

		if machine.State().Identifier() != loadingState.Identifier() {
			t.Errorf("expected state %s, got %s",
				loadingState.Identifier(), machine.State().Identifier())
		}

		if err := machine.Event(loadedEvent); err != nil {
			t.Errorf("unexpected error on loaded event: %v", err)
		}

		if machine.State().Identifier() != runningState.Identifier() {
			t.Errorf("expected state %s, got %s",
				runningState.Identifier(), machine.State().Identifier())
		}

		if err := machine.Event(exitedEvent); err != nil {
			t.Errorf("unexpected error on exited event: %v", err)
		}

		if machine.State().Identifier() != exitedState.Identifier() {
			t.Errorf("expected state %s, got %s",
				exitedState.Identifier(), machine.State().Identifier())
		}

		if machine.IsRunning() {
			t.Error("machine should not be running after final state")
		}
	})

	t.Run("should produce machine equivalent to direct RegisterState usage", func(t *testing.T) {
		t.Parallel()

		// Build machine A using BuildState/StateBuilder pattern
		builderA := NewBuilder[string, string, TestState]()

		builderA.BuildState(initialState).
			Transition(initializedEvent, runningState.Identifier()).
			Transition(failedEvent, failedState.Identifier()).
			Build()

		builderA.BuildState(runningState).
			Transition(exitedEvent, exitedState.Identifier()).
			Build()

		builderA.BuildState(failedState).Build()
		builderA.BuildState(exitedState).Build()

		machineA, err := builderA.
			GlobalTransitions(Transitions[string, string]{
				failedEvent: failedState.Identifier(),
			}).
			InitialState(initialState.Identifier()).
			FinalStates(failedState.Identifier(), exitedState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error building machine A: %v", err)
		}

		// Build machine B using direct RegisterState calls (same configuration)
		machineB, err := NewBuilder[string, string, TestState]().
			RegisterState(initialState, Transitions[string, string]{
				initializedEvent: runningState.Identifier(),
				failedEvent:      failedState.Identifier(),
			}).
			RegisterState(runningState, Transitions[string, string]{
				exitedEvent: exitedState.Identifier(),
			}).
			RegisterState(failedState, Transitions[string, string]{}).
			RegisterState(exitedState, Transitions[string, string]{}).
			GlobalTransitions(Transitions[string, string]{
				failedEvent: failedState.Identifier(),
			}).
			InitialState(initialState.Identifier()).
			FinalStates(failedState.Identifier(), exitedState.Identifier()).
			Build()

		if err != nil {
			t.Fatalf("unexpected error building machine B: %v", err)
		}

		// Drive both machines through same events and compare states
		events := []string{initializedEvent, exitedEvent}

		for _, event := range events {
			errA := machineA.Event(event)
			errB := machineB.Event(event)

			if errA != nil || errB != nil {
				t.Fatalf("event %s: machineA err=%v, machineB err=%v", event, errA, errB)
			}

			if machineA.State().Identifier() != machineB.State().Identifier() {
				t.Errorf("after event %s: machineA=%s, machineB=%s",
					event, machineA.State().Identifier(), machineB.State().Identifier())
			}
		}

		if machineA.IsRunning() != machineB.IsRunning() {
			t.Error("machines should have same running state")
		}
	})
}
