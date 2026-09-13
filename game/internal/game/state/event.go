package state

type Event int

const (
	NoEvent = Event(iota)
	StoppedEvent
	ResourcesLoadedEvent
)
