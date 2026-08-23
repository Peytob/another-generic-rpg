package glfw

import "engine/hid"

type Mouse struct {
	bindings hid.MouseBindings
}

func newMouse() *Mouse {
	return &Mouse{
		bindings: hid.NewMouseBindings(defaultBindingsName),
	}
}

func (m *Mouse) SetCurrentBindings(bindings hid.MouseBindings) {
	m.bindings = bindings
}

func (m *Mouse) GetCurrentBindings() hid.MouseBindings {
	return m.bindings
}
