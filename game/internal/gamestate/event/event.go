// Package event defines the game state machine event vocabulary shared by
// the FSM wiring and the frame loop.
package event

// Event enumerated type to describe game state FSM events
type Event int

const (
	NoEvent = Event(iota)
	StoppedEvent
	ResourcesLoadedEvent
)
