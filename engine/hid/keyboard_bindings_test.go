package hid

import (
	"testing"
)

func TestKeyboardBindingsAddAndGetCallback(t *testing.T) {
	var called bool
	bindings := NewKeyboardBindings("test")
	err := bindings.AddBinding(KeyboardBinding{
		Scancode: 42,
		Action:   Pressed,
		Mods:     ModShift,
	}, func(_ Key, _ int, _ Action, _ Modifier) { called = true })
	if err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}

	callback := bindings.GetCallback(42, Pressed, ModShift)
	callback(KeyA, 42, Pressed, ModShift)
	if !called {
		t.Fatal("expected bound callback to be invoked")
	}
}

func TestKeyboardBindingsSameScancodeDifferentAction(t *testing.T) {
	var pressed, released bool
	bindings := NewKeyboardBindings("test")
	_ = bindings.AddBinding(KeyboardBinding{
		Scancode: 42,
		Action:   Pressed,
	}, func(_ Key, _ int, _ Action, _ Modifier) { pressed = true })
	if err := bindings.AddBinding(KeyboardBinding{
		Scancode: 42,
		Action:   Released,
	}, func(_ Key, _ int, _ Action, _ Modifier) { released = true }); err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}

	bindings.GetCallback(42, Pressed, ModNone)(KeyA, 42, Pressed, ModNone)
	bindings.GetCallback(42, Released, ModNone)(KeyA, 42, Released, ModNone)

	if !pressed || !released {
		t.Fatalf("expected both actions to have own callbacks: pressed=%v released=%v", pressed, released)
	}
}

func TestKeyboardBindingsDuplicateError(t *testing.T) {
	bindings := NewKeyboardBindings("test")
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed}, emptyKeyboardCallback)

	err := bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed}, emptyKeyboardCallback)
	if err == nil {
		t.Fatal("expected duplicate binding error")
	}
}

func TestKeyboardBindingsNilCallbackError(t *testing.T) {
	bindings := NewKeyboardBindings("test")

	err := bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed}, nil)
	if err == nil {
		t.Fatal("expected nil callback error")
	}
}

func TestKeyboardBindingsNoMatchReturnsEmptyCallback(t *testing.T) {
	bindings := NewKeyboardBindings("test")
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed, Mods: ModShift}, emptyKeyboardCallback)

	if got := bindings.GetCallback(1, Pressed, ModNone); got == nil {
		t.Fatal("expected non-nil callback")
	}
	if got := bindings.GetCallback(2, Pressed, ModNone); got == nil {
		t.Fatal("expected non-nil callback")
	}

	bindings.GetCallback(1, Pressed, ModNone)(KeyA, 1, Pressed, ModNone)
	bindings.GetCallback(2, Pressed, ModNone)(KeyA, 2, Pressed, ModNone)
}

func TestKeyboardBindingsGetBindingsFor(t *testing.T) {
	bindings := NewKeyboardBindings("test")
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 42, Action: Pressed}, emptyKeyboardCallback)
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 42, Action: Released}, emptyKeyboardCallback)
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 7, Action: Pressed}, emptyKeyboardCallback)

	got, err := bindings.GetBindingsFor(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 bindings for scancode 42, got %d", len(got))
	}
}

func TestKeyboardBindingsName(t *testing.T) {
	if got := NewKeyboardBindings("movement").Name(); got != "movement" {
		t.Fatalf("expected name movement, got %q", got)
	}
}
