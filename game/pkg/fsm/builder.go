package fsm

import (
	"errors"
	"maps"

	mapset "github.com/deckarep/golang-set/v2"
)

type MachineBuilder[E comparable, I comparable, S State[I]] interface {
	// RegisterState add State and all transitions to this Machine. If state already exists in
	// machine, then transitions will be merged.
	RegisterState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S]

	// WriteState registers state if state not exists. Rewrites transitions if state already exists.
	WriteState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S]

	// InitialState state that be used as initial on Machine. If initial state already exists it will be rewritten.
	InitialState(state I) MachineBuilder[E, I, S]

	// GlobalTransitions Add new global transactions to Machine. If global transitions already specified, then
	// new transitions will be rewritten
	GlobalTransitions(transitions Transitions[E, I]) MachineBuilder[E, I, S]

	// FinalStates states that be used as final on Machine. If not set, then machine will work indefinitely.
	FinalStates(states []I) MachineBuilder[E, I, S]

	// Build returns built machine. Possible error causes is:
	// State used in transition but not registered
	// If initial state not set
	// No one state is registered
	Build() (Machine[E, I, S], error)

	// ShouldBuild Build with panic on error
	ShouldBuild() Machine[E, I, S]
}

type machineBuilder[E comparable, I comparable, S State[I]] struct {
	transitions             map[I]Transitions[E, I]
	globalTransitions       Transitions[E, I]
	states                  map[I]S
	finalStates             mapset.Set[I]
	initialState            I
	initialStateInitialized bool
}

// NewBuilder Creates new builder
func NewBuilder[E comparable, I comparable, S State[I]]() MachineBuilder[E, I, S] {
	m := &machineBuilder[E, I, S]{
		transitions:             make(map[I]Transitions[E, I]),
		globalTransitions:       make(Transitions[E, I]),
		states:                  make(map[I]S),
		finalStates:             mapset.NewThreadUnsafeSet[I](),
		initialStateInitialized: false,
	}

	return m
}

func (m *machineBuilder[E, I, S]) RegisterState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S] {
	if existsTransitions, ok := m.transitions[state.Identifier()]; ok {
		maps.Copy(existsTransitions, transitions)
	} else {
		m.WriteState(state, transitions)
	}

	return m
}

func (m *machineBuilder[E, I, S]) WriteState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S] {
	m.transitions[state.Identifier()] = transitions
	m.states[state.Identifier()] = state
	return m
}

func (m *machineBuilder[E, I, S]) InitialState(state I) MachineBuilder[E, I, S] {
	m.initialState = state
	m.initialStateInitialized = true
	return m
}

func (m *machineBuilder[E, I, S]) GlobalTransitions(transitions Transitions[E, I]) MachineBuilder[E, I, S] {
	maps.Copy(m.globalTransitions, transitions)
	return m
}

func (m *machineBuilder[E, I, S]) FinalStates(states []I) MachineBuilder[E, I, S] {
	for _, state := range states {
		if !m.finalStates.Contains(state) {
			m.finalStates.Add(state)
		}
	}

	return m
}

func (m *machineBuilder[E, I, S]) Build() (Machine[E, I, S], error) {
	for leftState := range m.transitions {
		if !m.containsState(leftState) {
			return nil, errors.New("transition left finalState not found")
		}

		for event := range m.transitions[leftState] {
			rightState := m.transitions[leftState][event]

			if !m.containsState(rightState) {
				return nil, errors.New("transition right finalState not found")
			}
		}
	}

	for event := range m.globalTransitions {
		rightState := m.globalTransitions[event]

		if !m.containsState(rightState) {
			return nil, errors.New("transition right finalState not found")
		}
	}

	for finalState := range m.finalStates.Iter() {
		if !m.containsState(finalState) {
			return nil, errors.New("final finalState not found")
		}
	}

	if !m.initialStateInitialized {
		return nil, errors.New("initial state not initialized")
	}

	if !m.containsState(m.initialState) {
		return nil, errors.New("initial state not registered")
	}

	return newMachine(m)
}

func (m *machineBuilder[E, I, S]) ShouldBuild() Machine[E, I, S] {
	res, err := m.Build()

	if err != nil {
		panic("failed ot build machine: " + err.Error())
	}

	return res
}

func (m *machineBuilder[E, I, S]) containsState(identifier I) bool {
	_, ok := m.states[identifier]
	return ok
}
