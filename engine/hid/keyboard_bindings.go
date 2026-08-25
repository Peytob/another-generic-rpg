package hid

import "fmt"

type KeyboardCallback func(key Key, scancode int, action Action, mods Modifier)

func emptyKeyboardCallback(_ Key, _ int, _ Action, _ Modifier) {
}

type KeyboardBinding struct {
	Scancode int32
	Action   Action
	Mods     Modifier
}

type KeyboardBindings struct {
	bindings[KeyboardBinding, KeyboardCallback]
}

func NewKeyboardBindings(name string) KeyboardBindings {
	return KeyboardBindings{
		bindings: newBindings[KeyboardBinding, KeyboardCallback](name),
	}
}

func (kb KeyboardBindings) AddBinding(binding KeyboardBinding, callback KeyboardCallback) error {
	if callback == nil {
		return fmt.Errorf("callback is nil: keyboard bindings %q", kb.Name())
	}
	return kb.add(binding, callback)
}

func (kb KeyboardBindings) GetBindingsFor(scancode int32) []KeyboardBinding {
	result := make([]KeyboardBinding, 0)
	for key := range kb.keys {
		if key.Scancode == scancode {
			result = append(result, key)
		}
	}
	return result
}

func (kb KeyboardBindings) GetCallback(scancode int32, action Action, mods Modifier) KeyboardCallback {
	if callback, ok := kb.lookup(KeyboardBinding{Scancode: scancode, Action: action, Mods: mods}); ok {
		return callback
	}
	return emptyKeyboardCallback
}
