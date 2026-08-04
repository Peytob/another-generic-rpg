package gamemachine

// Event enumerated type to describe engine FSM events
const (
	NoEvent = Event(iota)
	StoppedEvent
	ResourcesLoadedEvent
)
