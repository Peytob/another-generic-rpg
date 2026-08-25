package glfw

import (
	"engine/hid"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const defaultBindingsName = "default"

type Keyboard struct {
	bindings hid.KeyboardBindings
}

func newKeyboard() *Keyboard {
	return &Keyboard{
		bindings: hid.NewKeyboardBindings(defaultBindingsName),
	}
}

func (k *Keyboard) GetScancode(key hid.Key) int32 {
	return int32(glfw.GetKeyScancode(toGlfwKey(key)))
}

func (k *Keyboard) GetKeyName(scancode int32) string {
	return glfw.GetKeyName(glfw.KeyUnknown, int(scancode))
}

func (k *Keyboard) SetCurrentBindings(bindings hid.KeyboardBindings) {
	k.bindings = bindings
}

func (k *Keyboard) GetCurrentBindings() hid.KeyboardBindings {
	return k.bindings
}
