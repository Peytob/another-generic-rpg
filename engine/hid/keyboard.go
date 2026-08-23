package hid

type ScancodeMapper interface {
	GetScancode(key Key) int32
	GetKeyName(scancode int32)
}

type Keyboard interface {
	ScancodeMapper

	SetCurrentBindings(bindings KeyboardBindings)
	GetCurrentBindings() KeyboardBindings
}
