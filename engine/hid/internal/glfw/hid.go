package glfw

import (
	"errors"

	"engine/hid"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Hid struct {
	buffer   *InputBuffer
	keyboard *Keyboard
	mouse    *Mouse
}

func NewHid(win *glfw.Window) (*Hid, error) {
	if win == nil {
		return nil, errors.New("glfw window is nil")
	}

	h := newHid()
	win.SetKeyCallback(h.onKeyCallback)
	win.SetMouseButtonCallback(h.onMouseButtonCallback)

	return h, nil
}

func newHid() *Hid {
	h := &Hid{buffer: &InputBuffer{}}
	h.keyboard = newKeyboard()
	h.mouse = newMouse()
	return h
}

func (h *Hid) Keyboard() hid.Keyboard {
	return h.keyboard
}

func (h *Hid) Mouse() hid.Mouse {
	return h.mouse
}

// Dispatch drains buffered input events and invokes callbacks bound via
// current bindings. Key events are dispatched before mouse events, relative
// order between the two channels is not preserved.
func (h *Hid) Dispatch() {
	for _, e := range h.buffer.drainKeys() {
		h.keyboard.bindings.GetCallback(e.scancode, e.action, e.mods)(e.key, int(e.scancode), e.action, e.mods)
	}
	for _, e := range h.buffer.drainMouse() {
		h.mouse.bindings.GetCallback(e.button, e.action, e.mods)(e.button, e.action, e.mods)
	}
}

func (h *Hid) onKeyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	h.buffer.pushKey(keyEvent{
		key:      fromGlfwKey(key),
		scancode: int32(scancode),
		action:   fromGlfwAction(action),
		mods:     fromGlfwMods(mods),
	})
}

func (h *Hid) onMouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	h.buffer.pushMouse(mouseEvent{
		button: fromGlfwMouseButton(button),
		action: fromGlfwAction(action),
		mods:   fromGlfwMods(mods),
	})
}
