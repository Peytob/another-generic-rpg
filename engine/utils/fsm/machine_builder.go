package fsm

import (
	"maps"

	mapset "github.com/deckarep/golang-set/v2"
)

type MachineBuilder[E comparable, I comparable, S State[I]] interface {
	// RegisterState add State and all transitions to this Machine. If state already exists in
	// machine, then transitions will be merged.
	RegisterState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S]

	// BuildState returns state builder, that can be used to build state
	BuildState(state S) StateBuilder[E, I, S]

	// WriteState registers state if state not exists. Rewrites transitions if state already exists.
	WriteState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S]

	// InitialState state that be used as initial on Machine. If initial state already exists it will be rewritten.
	InitialState(state I) MachineBuilder[E, I, S]

	// GlobalTransitions adds global transitions to Machine. If a transition for an event is already
	// specified, then the old transition will be overwritten
	GlobalTransitions(transitions Transitions[E, I]) MachineBuilder[E, I, S]

	// FinalStates states that be used as final on Machine. If not set, then machine will work indefinitely.
	FinalStates(states ...I) MachineBuilder[E, I, S]

	// Build returns built machine. Possible error causes is:
	// State used in transition but not registered
	// If initial state not set
	// No one state is registered
	Build() (Machine[E, I, S], error)

	// MustBuild Build with panic on error
	MustBuild() Machine[E, I, S]
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

func (m *machineBuilder[E, I, S]) BuildState(state S) StateBuilder[E, I, S] {
	return newStateBuilder(m, state)
}

func (m *machineBuilder[E, I, S]) WriteState(state S, transitions Transitions[E, I]) MachineBuilder[E, I, S] {
	m.transitions[state.Identifier()] = maps.Clone(transitions)
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

func (m *machineBuilder[E, I, S]) FinalStates(states ...I) MachineBuilder[E, I, S] {
	for _, state := range states {
		m.finalStates.Add(state)
	}

	return m
}

func (m *machineBuilder[E, I, S]) Build() (Machine[E, I, S], error) {
	if len(m.states) == 0 {
		return nil, ErrNoStates
	}

	for leftState := range m.transitions {
		if !m.containsState(leftState) {
			return nil, ErrTransitionLeftStateNotFound
		}

		for event := range m.transitions[leftState] {
			rightState := m.transitions[leftState][event]

			if !m.containsState(rightState) {
				return nil, ErrTransitionRightStateNotFound
			}
		}
	}

	for event := range m.globalTransitions {
		rightState := m.globalTransitions[event]

		if !m.containsState(rightState) {
			return nil, ErrTransitionRightStateNotFound
		}
	}

	for finalState := range m.finalStates.Iter() {
		if !m.containsState(finalState) {
			return nil, ErrFinalStateNotFound
		}
	}

	if !m.initialStateInitialized {
		return nil, ErrInitialStateNotInitialized
	}

	if !m.containsState(m.initialState) {
		return nil, ErrInitialStateNotRegistered
	}

	return newMachine(m)
}

func (m *machineBuilder[E, I, S]) MustBuild() Machine[E, I, S] {
	res, err := m.Build()

	if err != nil {
		panic("failed to build machine: " + err.Error())
	}

	return res
}

func (m *machineBuilder[E, I, S]) containsState(identifier I) bool {
	_, ok := m.states[identifier]
	return ok
}
