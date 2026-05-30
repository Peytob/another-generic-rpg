package fsm

import "game/pkg/engine"

// Event enumerated type to describe engine FSM events
const (
	NoEvent = engine.Event(iota)
	StoppedEvent
)
