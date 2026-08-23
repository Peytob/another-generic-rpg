package hid

type KeyboardCallback func(key Key, scancode int, action Action, mods Modifier)

func emptyKeyboardCallback(_ Key, _ int, _ Action, _ Modifier) {
}

type KeyboardBinding struct {
	Scancode int32
	Action   Action
	Mods     Modifier
	Callback KeyboardCallback
}

type KeyboardBindings struct {
	name     string
	bindings map[int32]KeyboardBinding
}

func NewKeyboardBindings() KeyboardBindings {
	return KeyboardBindings{
		bindings: make(map[int32]KeyboardBinding),
	}
}

func (kb KeyboardBindings) AddBinding(binding KeyboardBinding) error {
	return nil
}

func (kb KeyboardBindings) GetBindingsFor(scancode int32) (KeyboardBinding, error) {
	return KeyboardBinding{}, nil
}

func (kb KeyboardBindings) GetCallback(scancode int32, action Action, mods Modifier) KeyboardCallback {
	return emptyKeyboardCallback
}

func (kb KeyboardBindings) Name() string {
	return kb.name
}
