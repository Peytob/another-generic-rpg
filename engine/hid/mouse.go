package hid

type Mouse interface {
	SetCurrentBindings(bindings MouseBindings)
	GetCurrentBindings() MouseBindings
}
