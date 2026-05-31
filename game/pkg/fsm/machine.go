package fsm

import (
	"fmt"
	"maps"

	mapset "github.com/deckarep/golang-set/v2"
)

type MachineError string

func (e MachineError) Error() string {
	return string(e)
}

const (
	// StateNotFoundError operation state not found inside machine
	StateNotFoundError = MachineError("state not found")

	// MachineNotRunningError operation not allowed for stopped machine
	MachineNotRunningError = MachineError("machine not running")

	// UnknownEventError event not found in transitions table
	UnknownEventError = MachineError("unknown event")
)

// Machine Abstract finite state machine.
type Machine[E comparable, I comparable, S State[I]] interface {
	// Event changes current machine state according to given event
	Event(event E) error

	// Result if machine on end state now - returns current state and true flag, nil and false otherwise
	Result() (state S, ok bool)

	// IsRunning is machine on end state now
	IsRunning() bool

	// State returns current machine state
	State() S
}

type State[I comparable] interface {
	Identifier() I
}

type Transitions[E comparable, I comparable] map[E]I

type machine[E comparable, I comparable, S State[I]] struct {
	transitions       map[I]Transitions[E, I]
	globalTransitions Transitions[E, I]

	states       map[I]S
	finalStates  mapset.Set[I]
	currentState S
	isRunning    bool
}

func newMachine[E comparable, I comparable, S State[I]](builder *machineBuilder[E, I, S]) (Machine[E, I, S], error) {
	initialState, initialStateFound := builder.states[builder.initialState]

	if !initialStateFound {
		return nil, fmt.Errorf("initial state not set: %w", StateNotFoundError)
	}

	return &machine[E, I, S]{
		transitions:       maps.Clone(builder.transitions),
		globalTransitions: maps.Clone(builder.globalTransitions),

		states:       builder.states,
		finalStates:  builder.finalStates.Clone(),
		currentState: initialState,
		isRunning:    true,
	}, nil
}

func (m *machine[E, I, S]) Event(event E) error {
	currentState := m.currentState.Identifier()

	if transitions, ok := m.transitions[currentState]; ok {
		if nextI, ok := transitions[event]; ok {
			return m.changeState(nextI)
		}
	} else {
		return fmt.Errorf("current state not found in transitions table: %w", StateNotFoundError)
	}

	if nextI, ok := m.globalTransitions[event]; ok {
		return m.changeState(nextI)
	}

	return UnknownEventError
}

func (m *machine[E, I, S]) Result() (state S, ok bool) {
	if !m.IsRunning() {
		return m.currentState, true
	}

	var noop S
	return noop, m.IsRunning()
}

func (m *machine[E, I, S]) IsRunning() bool {
	return m.isRunning
}

func (m *machine[E, I, S]) State() S {
	return m.currentState
}

func (m *machine[E, I, S]) changeState(nextState I) error {
	if !m.isRunning {
		return MachineNotRunningError
	}

	if state, ok := m.states[nextState]; ok {
		m.currentState = state
		m.isRunning = !m.finalStates.Contains(state.Identifier())
		return nil
	}

	return StateNotFoundError
}
