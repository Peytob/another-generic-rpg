package hid

import "fmt"

type KeyboardCallback func(key Key, scancode int, action Action, mods Modifier)

func emptyKeyboardCallback(_ Key, _ int, _ Action, _ Modifier) {
}

type KeyboardBinding struct {
	Scancode int32
	Action   Action
	Mods     Modifier
	Callback KeyboardCallback
}

type keyboardBindingKey struct {
	scancode int32
	action   Action
	mods     Modifier
}

type KeyboardBindings struct {
	name     string
	bindings map[keyboardBindingKey]KeyboardBinding
}

func NewKeyboardBindings(name string) KeyboardBindings {
	return KeyboardBindings{
		name:     name,
		bindings: make(map[keyboardBindingKey]KeyboardBinding),
	}
}

func (kb KeyboardBindings) AddBinding(binding KeyboardBinding) error {
	if binding.Callback == nil {
		return fmt.Errorf("%w: keyboard bindings %q", ErrNilCallback, kb.name)
	}

	key := keyboardBindingKey{scancode: binding.Scancode, action: binding.Action, mods: binding.Mods}
	if _, exists := kb.bindings[key]; exists {
		return fmt.Errorf("%w: scancode=%d action=%d mods=%d", ErrBindingAlreadyExists, binding.Scancode, binding.Action, binding.Mods)
	}

	kb.bindings[key] = binding
	return nil
}

func (kb KeyboardBindings) GetBindingsFor(scancode int32) ([]KeyboardBinding, error) {
	result := make([]KeyboardBinding, 0)
	for key, binding := range kb.bindings {
		if key.scancode == scancode {
			result = append(result, binding)
		}
	}
	return result, nil
}

func (kb KeyboardBindings) GetCallback(scancode int32, action Action, mods Modifier) KeyboardCallback {
	if binding, ok := kb.bindings[keyboardBindingKey{scancode: scancode, action: action, mods: mods}]; ok {
		return binding.Callback
	}
	return emptyKeyboardCallback
}

func (kb KeyboardBindings) Name() string {
	return kb.name
}
