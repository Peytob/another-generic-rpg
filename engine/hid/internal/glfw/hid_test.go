package glfw

import (
	"testing"

	"engine/hid"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestNewHidNilWindowError(t *testing.T) {
	if _, err := NewHid(nil); err == nil {
		t.Fatal("expected error for nil window")
	}
}

func TestHidImplementsHidInterface(t *testing.T) {
	var _ hid.Hid = newHid()
}

func TestDispatchKeyboardEvent(t *testing.T) {
	h := newHid()

	type call struct {
		key      hid.Key
		scancode int
		action   hid.Action
		mods     hid.Modifier
	}
	var calls []call

	bindings := hid.NewKeyboardBindings("test")
	err := bindings.AddBinding(hid.KeyboardBinding{
		Scancode: 42,
		Action:   hid.Pressed,
		Mods:     hid.ModShift,
	}, func(key hid.Key, scancode int, action hid.Action, mods hid.Modifier) {
		calls = append(calls, call{key, scancode, action, mods})
	})
	if err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}
	h.Keyboard().SetCurrentBindings(bindings)

	h.onKeyCallback(nil, glfw.KeyA, 42, glfw.Press, glfw.ModShift)
	h.Dispatch()

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	got := calls[0]
	if got.key != hid.KeyA || got.scancode != 42 || got.action != hid.Pressed || got.mods != hid.ModShift {
		t.Fatalf("unexpected call args: %+v", got)
	}

	h.Dispatch()
	if len(calls) != 1 {
		t.Fatal("expected buffer to be drained by previous Dispatch")
	}
}

func TestDispatchMouseEvent(t *testing.T) {
	h := newHid()

	type call struct {
		button hid.MouseButton
		action hid.Action
		mods   hid.Modifier
	}
	var calls []call

	bindings := hid.NewMouseBindings("test")
	err := bindings.AddBinding(hid.MouseBinding{
		Button: hid.MouseButtonRight,
		Action: hid.Released,
	}, func(button hid.MouseButton, action hid.Action, mods hid.Modifier) {
		calls = append(calls, call{button, action, mods})
	})
	if err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}
	h.Mouse().SetCurrentBindings(bindings)

	h.onMouseButtonCallback(nil, glfw.MouseButtonRight, glfw.Release, 0)
	h.Dispatch()

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	got := calls[0]
	if got.button != hid.MouseButtonRight || got.action != hid.Released || got.mods != hid.ModNone {
		t.Fatalf("unexpected call args: %+v", got)
	}
}

func TestDispatchUnmatchedEvent(t *testing.T) {
	h := newHid()

	h.onKeyCallback(nil, glfw.KeyZ, 999, glfw.Repeat, glfw.ModAlt)
	h.onMouseButtonCallback(nil, glfw.MouseButton8, glfw.Press, 0)

	h.Dispatch()
}
