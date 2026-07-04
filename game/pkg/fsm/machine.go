package fsm

import (
	"fmt"
	"maps"

	mapset "github.com/deckarep/golang-set/v2"
)

// Machine Abstract finite state machine. Only for single-gorutine use, add mutex if you want use it from
// many gorutines
type Machine[E comparable, I comparable, S State[I]] interface {
	// Event changes current machine state according to given event
	Event(event E) error

	// Result if machine on final state - returns current state and true flag, zero value and false otherwise
	Result() (state S, ok bool)

	// IsRunning returns whether the machine is still running (has not reached a final state)
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
		return nil, fmt.Errorf("initial state not set: %w", ErrStateNotFound)
	}

	transitions := make(map[I]Transitions[E, I], len(builder.transitions))
	for id, table := range builder.transitions {
		transitions[id] = maps.Clone(table)
	}

	return &machine[E, I, S]{
		transitions:       transitions,
		globalTransitions: maps.Clone(builder.globalTransitions),

		states:       maps.Clone(builder.states),
		finalStates:  builder.finalStates.Clone(),
		currentState: initialState,
		isRunning:    true,
	}, nil
}

func (m *machine[E, I, S]) Event(event E) error {
	if !m.isRunning {
		return ErrMachineNotRunning
	}

	currentState := m.currentState.Identifier()

	if transitions, ok := m.transitions[currentState]; ok {
		if nextI, ok := transitions[event]; ok {
			return m.changeState(nextI)
		}
	}

	if nextI, ok := m.globalTransitions[event]; ok {
		return m.changeState(nextI)
	}

	return ErrUnknownEvent
}

func (m *machine[E, I, S]) Result() (S, bool) {
	if !m.IsRunning() {
		return m.currentState, true
	}

	var zero S
	return zero, false
}

func (m *machine[E, I, S]) IsRunning() bool {
	return m.isRunning
}

func (m *machine[E, I, S]) State() S {
	return m.currentState
}

func (m *machine[E, I, S]) changeState(nextState I) error {
	if state, ok := m.states[nextState]; ok {
		m.currentState = state
		m.isRunning = !m.finalStates.Contains(state.Identifier())
		return nil
	}

	return ErrStateNotFound
}
