package hid

import (
	"errors"
	"testing"
)

func TestKeyboardBindingsAddAndGetCallback(t *testing.T) {
	var called bool
	bindings := NewKeyboardBindings("test")
	err := bindings.AddBinding(KeyboardBinding{
		Scancode: 42,
		Action:   Pressed,
		Mods:     ModShift,
		Callback: func(_ Key, _ int, _ Action, _ Modifier) { called = true },
	})
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
		Callback: func(_ Key, _ int, _ Action, _ Modifier) { pressed = true },
	})
	if err := bindings.AddBinding(KeyboardBinding{
		Scancode: 42,
		Action:   Released,
		Callback: func(_ Key, _ int, _ Action, _ Modifier) { released = true },
	}); err != nil {
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
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed, Callback: emptyKeyboardCallback})

	err := bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed, Callback: emptyKeyboardCallback})
	if !errors.Is(err, ErrBindingAlreadyExists) {
		t.Fatalf("expected ErrBindingAlreadyExists, got: %v", err)
	}
}

func TestKeyboardBindingsNilCallbackError(t *testing.T) {
	bindings := NewKeyboardBindings("test")

	err := bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed})
	if !errors.Is(err, ErrNilCallback) {
		t.Fatalf("expected ErrNilCallback, got: %v", err)
	}
}

func TestKeyboardBindingsNoMatchReturnsEmptyCallback(t *testing.T) {
	bindings := NewKeyboardBindings("test")
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 1, Action: Pressed, Mods: ModShift, Callback: emptyKeyboardCallback})

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
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 42, Action: Pressed, Callback: emptyKeyboardCallback})
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 42, Action: Released, Callback: emptyKeyboardCallback})
	_ = bindings.AddBinding(KeyboardBinding{Scancode: 7, Action: Pressed, Callback: emptyKeyboardCallback})

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
