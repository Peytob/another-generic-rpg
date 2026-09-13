// Package input provides per-frame input snapshots collected from hid
// binding callbacks.
package input

import (
	"maps"
)

// Action enumerates game input actions recognized by gameplay states.
type Action int

const (
	MoveUp Action = iota
	MoveDown
	MoveLeft
	MoveRight
	Exit
)

// InputFrame is an immutable snapshot of the input events collected during
// the current frame.
type InputFrame struct {
	pressed  map[Action]struct{}
	released map[Action]struct{}
	held     map[Action]struct{}
}

// Pressed reports whether the action was pressed during this frame.
func (f InputFrame) Pressed(a Action) bool {
	_, ok := f.pressed[a]
	return ok
}

// Released reports whether the action was released during this frame.
func (f InputFrame) Released(a Action) bool {
	_, ok := f.released[a]
	return ok
}

// Held reports whether the action is currently held down.
func (f InputFrame) Held(a Action) bool {
	_, ok := f.held[a]
	return ok
}

// Empty reports whether the snapshot contains no events at all.
func (f InputFrame) Empty() bool {
	return len(f.pressed) == 0 && len(f.released) == 0 && len(f.held) == 0
}

// InputCollector buffers raw input events delivered by hid binding
// callbacks and converts them into per-frame snapshots.
//
// It is not safe for concurrent use: binding callbacks and Drain must be
// invoked from the same goroutine (hid Dispatch and the game loop both run
// on the main game goroutine).
type InputCollector struct {
	pressed  []Action
	released []Action
	held     map[Action]struct{}
}

// NewInputCollector creates an empty InputCollector.
func NewInputCollector() *InputCollector {
	return &InputCollector{held: make(map[Action]struct{})}
}

// Press records a press event and marks the action as held.
func (c *InputCollector) Press(a Action) {
	c.pressed = append(c.pressed, a)
	c.held[a] = struct{}{}
}

// Release records a release event and drops the action from held.
func (c *InputCollector) Release(a Action) {
	c.released = append(c.released, a)
	delete(c.held, a)
}

// Drain converts buffered events into an InputFrame and clears the buffer.
// Press and release events are per-frame: each appears in exactly one
// snapshot; held state persists across frames until release.
func (c *InputCollector) Drain() InputFrame {
	pressed := make(map[Action]struct{}, len(c.pressed))
	for _, a := range c.pressed {
		pressed[a] = struct{}{}
	}

	released := make(map[Action]struct{}, len(c.released))
	for _, a := range c.released {
		released[a] = struct{}{}
	}

	held := make(map[Action]struct{}, len(c.held))
	maps.Copy(held, c.held)

	frame := InputFrame{pressed: pressed, released: released, held: held}

	c.pressed = nil
	c.released = nil

	return frame
}
