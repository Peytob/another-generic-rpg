package hid

type Action int32

const (
	Pressed = Action(iota)
	Released
	Repeat
)

type Hid interface {
	Keyboard() Keyboard
	Mouse() Mouse

	// Dispatch drains buffered input events and invokes callbacks
	// bound to them via current keyboard/mouse bindings.
	Dispatch()
}
