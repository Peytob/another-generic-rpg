package hid

import (
	"errors"
	"testing"
)

func TestMouseBindingsAddAndGetCallback(t *testing.T) {
	var called bool
	bindings := NewMouseBindings("test")
	err := bindings.AddBinding(MouseBinding{
		Button:   MouseButtonRight,
		Action:   Pressed,
		Mods:     ModControl,
		Callback: func(_ MouseButton, _ Action, _ Modifier) { called = true },
	})
	if err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}

	bindings.GetCallback(MouseButtonRight, Pressed, ModControl)(MouseButtonRight, Pressed, ModControl)
	if !called {
		t.Fatal("expected bound callback to be invoked")
	}
}

func TestMouseBindingsDuplicateError(t *testing.T) {
	bindings := NewMouseBindings("test")
	_ = bindings.AddBinding(MouseBinding{Button: MouseButtonLeft, Action: Pressed, Callback: emptyMouseCallback})

	err := bindings.AddBinding(MouseBinding{Button: MouseButtonLeft, Action: Pressed, Callback: emptyMouseCallback})
	if !errors.Is(err, ErrBindingAlreadyExists) {
		t.Fatalf("expected ErrBindingAlreadyExists, got: %v", err)
	}
}

func TestMouseBindingsNilCallbackError(t *testing.T) {
	bindings := NewMouseBindings("test")

	err := bindings.AddBinding(MouseBinding{Button: MouseButtonLeft, Action: Pressed})
	if !errors.Is(err, ErrNilCallback) {
		t.Fatalf("expected ErrNilCallback, got: %v", err)
	}
}

func TestMouseBindingsGetBindingsFor(t *testing.T) {
	bindings := NewMouseBindings("test")
	_ = bindings.AddBinding(MouseBinding{Button: MouseButtonMiddle, Action: Pressed, Callback: emptyMouseCallback})
	_ = bindings.AddBinding(MouseBinding{Button: MouseButtonMiddle, Action: Released, Callback: emptyMouseCallback})
	_ = bindings.AddBinding(MouseBinding{Button: MouseButton8, Action: Pressed, Callback: emptyMouseCallback})

	got, err := bindings.GetBindingsFor(MouseButtonMiddle)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 bindings for middle button, got %d", len(got))
	}
}

func TestMouseBindingsName(t *testing.T) {
	if got := NewMouseBindings("camera").Name(); got != "camera" {
		t.Fatalf("expected name camera, got %q", got)
	}
}
