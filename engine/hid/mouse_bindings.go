package hid

import "fmt"

type MouseCallback func(button MouseButton, action Action, mods Modifier)

func emptyMouseCallback(_ MouseButton, _ Action, _ Modifier) {
}

type MouseBinding struct {
	Button   MouseButton
	Action   Action
	Mods     Modifier
	Callback MouseCallback
}

type mouseBindingKey struct {
	button MouseButton
	action Action
	mods   Modifier
}

type MouseBindings struct {
	name     string
	bindings map[mouseBindingKey]MouseBinding
}

func NewMouseBindings(name string) MouseBindings {
	return MouseBindings{
		name:     name,
		bindings: make(map[mouseBindingKey]MouseBinding),
	}
}

func (mb MouseBindings) AddBinding(binding MouseBinding) error {
	if binding.Callback == nil {
		return fmt.Errorf("%w: mouse bindings %q", ErrNilCallback, mb.name)
	}

	key := mouseBindingKey{button: binding.Button, action: binding.Action, mods: binding.Mods}
	if _, exists := mb.bindings[key]; exists {
		return fmt.Errorf("%w: button=%d action=%d mods=%d", ErrBindingAlreadyExists, binding.Button, binding.Action, binding.Mods)
	}

	mb.bindings[key] = binding
	return nil
}

func (mb MouseBindings) GetBindingsFor(button MouseButton) ([]MouseBinding, error) {
	result := make([]MouseBinding, 0)
	for key, binding := range mb.bindings {
		if key.button == button {
			result = append(result, binding)
		}
	}
	return result, nil
}

func (mb MouseBindings) GetCallback(button MouseButton, action Action, mods Modifier) MouseCallback {
	if binding, ok := mb.bindings[mouseBindingKey{button: button, action: action, mods: mods}]; ok {
		return binding.Callback
	}
	return emptyMouseCallback
}

func (mb MouseBindings) Name() string {
	return mb.name
}
