package hid

type MouseCallback func(button MouseButton, action Action, mods Modifier)

func emptyMouseCallback(_ MouseButton, _ Action, _ Modifier) {
}

type MouseBinding struct {
	Button   MouseButton
	Action   Action
	Mods     Modifier
	Callback MouseCallback
}

type MouseBindings struct {
	name     string
	bindings map[MouseButton]MouseBinding
}

func NewMouseBindings() MouseBindings {
	return MouseBindings{
		bindings: make(map[MouseButton]MouseBinding),
	}
}

func (kb MouseBindings) AddBinding(binding MouseBinding) error {
	return nil
}

func (kb MouseBindings) GetBindingsFor(scancode int32) (MouseBinding, error) {
	return MouseBinding{}, nil
}

func (kb MouseBindings) GetCallback(scancode int32, action Action, mods Modifier) MouseCallback {
	return emptyMouseCallback
}

func (kb MouseBindings) Name() string {
	return kb.name
}
