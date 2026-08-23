package hid

import "fmt"

type MouseCallback func(button MouseButton, action Action, mods Modifier)

func emptyMouseCallback(_ MouseButton, _ Action, _ Modifier) {
}

type MouseBinding struct {
	Button MouseButton
	Action Action
	Mods   Modifier
}

type MouseBindings struct {
	bindings[MouseBinding, MouseCallback]
}

func NewMouseBindings(name string) MouseBindings {
	return MouseBindings{
		bindings: newBindings[MouseBinding, MouseCallback](name),
	}
}

func (mb MouseBindings) AddBinding(binding MouseBinding, callback MouseCallback) error {
	if callback == nil {
		return fmt.Errorf("callback is nil: mouse bindings %q", mb.Name())
	}
	return mb.add(binding, callback)
}

func (mb MouseBindings) GetBindingsFor(button MouseButton) []MouseBinding {
	result := make([]MouseBinding, 0)
	for key := range mb.keys {
		if key.Button == button {
			result = append(result, key)
		}
	}
	return result
}

func (mb MouseBindings) GetCallback(button MouseButton, action Action, mods Modifier) MouseCallback {
	if callback, ok := mb.lookup(MouseBinding{Button: button, Action: action, Mods: mods}); ok {
		return callback
	}
	return emptyMouseCallback
}
