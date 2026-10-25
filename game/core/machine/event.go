package machine

type Event int

const (
	NoEvent = Event(iota)
	StoppedEvent

	InitializedEvent
	ResourcesLoadedEvent
)
