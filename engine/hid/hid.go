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
}
