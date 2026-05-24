package fsm

import (
	"testing"
)

type TestState string

func (s TestState) Identifier() string {
	return string(s)
}

const (
	initialState = TestState("initializing")
	runningState = TestState("running")
	loadingState = TestState("loading")
	exitedState  = TestState("exited")
	failedState  = TestState("failed")

	initializedEvent  = "initialized"
	failedEvent       = "failed"
	levelChangedEvent = "levelChanged"
	loadedEvent       = "loaded"
	exitedEvent       = "exited"
)

func createTestMachine() Machine[string, string, TestState] {
	m, _ := NewBuilder[string, string, TestState]().
		RegisterState(initialState, Transitions[string, string]{
			initializedEvent: runningState.Identifier(),
		}).
		RegisterState(runningState, Transitions[string, string]{
			exitedEvent:       exitedState.Identifier(),
			levelChangedEvent: loadingState.Identifier(),
		}).
		RegisterState(loadingState, Transitions[string, string]{
			loadedEvent: runningState.Identifier(),
			failedEvent: initialState.Identifier(), // will be reset just for event priority testing
		}).
		RegisterState(failedState, Transitions[string, string]{}).
		RegisterState(exitedState, Transitions[string, string]{}).
		GlobalTransitions(Transitions[string, string]{
			failedEvent: failedState.Identifier(),
		}).
		InitialState(initialState.Identifier()).
		FinalStates([]string{failedState.Identifier(), exitedState.Identifier()}).
		Build()

	return m
}

func TestEventChanges(t *testing.T) {
	t.Run("should change machine state according to local transitions table", func(t *testing.T) {
		m := createTestMachine()

		if err := m.Event(initializedEvent); err != nil {
			t.Error(err)
		}

		if m.State().Identifier() != runningState.Identifier() {
			t.Errorf("wrong state after initialized event: %s", m.State().Identifier())
		}
	})

	t.Run("should change machine state according to global transitions table", func(t *testing.T) {
		m := createTestMachine()

		if err := m.Event(failedEvent); err != nil {
			t.Error(err)
		}

		if m.State().Identifier() != failedState.Identifier() {
			t.Errorf("wrong state after failed event: %s", m.State().Identifier())
		}
	})

	t.Run("should change machine with higher local transitions priority", func(t *testing.T) {
		m := createTestMachine()

		if err := m.Event(initializedEvent); err != nil {
			t.Error(err)
		}

		if err := m.Event(levelChangedEvent); err != nil {
			t.Error(err)
		}

		// Failed event on level changing state will throw machine back to initializing state
		if err := m.Event(failedEvent); err != nil {
			t.Error(err)
		}

		if m.State().Identifier() != initialState.Identifier() {
			t.Errorf("wrong state after failed event: %s", m.State().Identifier())
		}
	})
}

func TestFinalStates(t *testing.T) {
	t.Run("should not run after final state", func(t *testing.T) {
		m := createTestMachine()

		if err := m.Event(failedEvent); err != nil {
			t.Error(err)
		}

		if m.IsRunning() {
			t.Error("machine is running in final state")
		}
	})

	t.Run("should return result in final state", func(t *testing.T) {
		m := createTestMachine()

		if err := m.Event(failedEvent); err != nil {
			t.Error(err)
		}

		state, ok := m.Result()

		if !ok {
			t.Error("result flag is not ok in final state")
		}

		if state.Identifier() != failedState.Identifier() {
			t.Error("wrong state in machine result")
		}
	})

	t.Run("should run in not final state", func(t *testing.T) {
		m := createTestMachine()

		if !m.IsRunning() {
			t.Error("machine is not running in non final state")
		}
	})
}
