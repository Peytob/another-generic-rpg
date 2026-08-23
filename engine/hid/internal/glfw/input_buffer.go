package glfw

import "engine/hid"

type keyEvent struct {
	key      hid.Key
	scancode int32
	action   hid.Action
	mods     hid.Modifier
}

type mouseEvent struct {
	button hid.MouseButton
	action hid.Action
	mods   hid.Modifier
}

// InputBuffer accumulates input events pushed from GLFW callbacks until
// they are drained by Hid.Dispatch. GLFW contract: all access must happen
// on the main thread only, no synchronization is required.
type InputBuffer struct {
	keyEvents   []keyEvent
	mouseEvents []mouseEvent
}

func (b *InputBuffer) pushKey(event keyEvent) {
	b.keyEvents = append(b.keyEvents, event)
}

func (b *InputBuffer) pushMouse(event mouseEvent) {
	b.mouseEvents = append(b.mouseEvents, event)
}

func (b *InputBuffer) drainKeys() []keyEvent {
	events := b.keyEvents
	b.keyEvents = nil
	return events
}

func (b *InputBuffer) drainMouse() []mouseEvent {
	events := b.mouseEvents
	b.mouseEvents = nil
	return events
}
