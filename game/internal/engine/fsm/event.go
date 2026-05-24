package fsm

// Event enumerated type to describe engine FSM events
type Event int

const (
	NoEvent = Event(iota)
	StoppedEvent
)
