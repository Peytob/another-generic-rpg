package fsm

type StateBuilder[E comparable, I comparable, S State[I]] interface {
	// Transition adds transition to state. If transition for this event already exists it will be rewritten
	Transition(event E, nextState I) StateBuilder[E, I, S]

	// TransitionS adds transition to state using State directly. If transition for this event already exists it will be rewritten
	TransitionS(event E, nextState S) StateBuilder[E, I, S]

	// Build calls 'RegisterState' from parent MachineBuilder and returns it.
	Build() MachineBuilder[E, I, S]
}

func newStateBuilder[E comparable, I comparable, S State[I]](
	m *machineBuilder[E, I, S],
	state S,
) StateBuilder[E, I, S] {
	return stateBuilder[E, I, S]{
		parentMachine:      m,
		state:              state,
		transitionsBuilder: NewTransitionsBuilder[E, I](),
	}
}

type stateBuilder[E comparable, I comparable, S State[I]] struct {
	parentMachine      *machineBuilder[E, I, S]
	state              S
	transitionsBuilder TransitionsBuilder[E, I]
}

func (s stateBuilder[E, I, S]) Transition(event E, nextState I) StateBuilder[E, I, S] {
	s.transitionsBuilder.AddTransition(event, nextState)
	return s
}

func (s stateBuilder[E, I, S]) TransitionS(event E, nextState S) StateBuilder[E, I, S] {
	return s.Transition(event, nextState.Identifier())
}

func (s stateBuilder[E, I, S]) Build() MachineBuilder[E, I, S] {
	return s.parentMachine.RegisterState(s.state, s.transitionsBuilder.Build())
}
