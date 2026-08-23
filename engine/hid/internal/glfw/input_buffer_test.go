package glfw

import (
	"testing"

	"engine/hid"
)

func TestInputBufferFIFO(t *testing.T) {
	b := &InputBuffer{}

	b.pushKey(keyEvent{key: hid.KeyA, scancode: 1, action: hid.Pressed})
	b.pushKey(keyEvent{key: hid.KeyB, scancode: 2, action: hid.Released})
	b.pushMouse(mouseEvent{button: hid.MouseButtonLeft, action: hid.Pressed})

	keys := b.drainKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 key events, got %d", len(keys))
	}
	if keys[0].key != hid.KeyA || keys[1].key != hid.KeyB {
		t.Fatalf("expected FIFO order, got %d then %d", keys[0].key, keys[1].key)
	}

	mouse := b.drainMouse()
	if len(mouse) != 1 || mouse[0].button != hid.MouseButtonLeft {
		t.Fatalf("expected 1 left button event, got %v", mouse)
	}
}

func TestInputBufferDrainClears(t *testing.T) {
	b := &InputBuffer{}
	b.pushKey(keyEvent{key: hid.KeyA, action: hid.Pressed})
	b.pushMouse(mouseEvent{button: hid.MouseButtonLeft, action: hid.Pressed})

	_ = b.drainKeys()
	_ = b.drainMouse()

	if keys := b.drainKeys(); len(keys) != 0 {
		t.Fatalf("expected drained key buffer, got %d events", len(keys))
	}
	if mouse := b.drainMouse(); len(mouse) != 0 {
		t.Fatalf("expected drained mouse buffer, got %d events", len(mouse))
	}
}
