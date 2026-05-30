package fsm

import "maps"

type TransitionsBuilder[E comparable, I comparable] interface {
	AddTransition(event E, identifier I) TransitionsBuilder[E, I]
	RemoveTransition(event E) TransitionsBuilder[E, I]
	Build() Transitions[E, I]
}

type transitionsBuilder[E comparable, I comparable] struct {
	transitions Transitions[E, I]
}

func NewTransitionsBuilder[E comparable, I comparable]() TransitionsBuilder[E, I] {
	return &transitionsBuilder[E, I]{
		transitions: make(Transitions[E, I]),
	}
}

func (t transitionsBuilder[E, I]) AddTransition(event E, identifier I) TransitionsBuilder[E, I] {
	t.transitions[event] = identifier
	return t
}

func (t transitionsBuilder[E, I]) RemoveTransition(event E) TransitionsBuilder[E, I] {
	delete(t.transitions, event)
	return t
}

func (t transitionsBuilder[E, I]) Build() Transitions[E, I] {
	return maps.Clone(t.transitions)
}
