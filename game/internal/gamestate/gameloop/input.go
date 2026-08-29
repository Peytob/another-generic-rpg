package gameloop

import (
	"context"
	"maps"
)

// InputFrame is an immutable snapshot of the input events collected during
// the current frame.
type InputFrame[A comparable] struct {
	pressed  map[A]struct{}
	released map[A]struct{}
	held     map[A]struct{}
}

// Pressed reports whether the action was pressed during this frame.
func (f InputFrame[A]) Pressed(a A) bool {
	_, ok := f.pressed[a]
	return ok
}

// Released reports whether the action was released during this frame.
func (f InputFrame[A]) Released(a A) bool {
	_, ok := f.released[a]
	return ok
}

// Held reports whether the action is currently held down.
func (f InputFrame[A]) Held(a A) bool {
	_, ok := f.held[a]
	return ok
}

// Empty reports whether the snapshot contains no events at all.
func (f InputFrame[A]) Empty() bool {
	return len(f.pressed) == 0 && len(f.released) == 0 && len(f.held) == 0
}

// InputCollector buffers raw input events delivered by hid binding
// callbacks and converts them into per-frame snapshots.
//
// It is not safe for concurrent use: binding callbacks and Drain must be
// invoked from the same goroutine (hid Dispatch and the game loop both run
// on the main game goroutine).
type InputCollector[A comparable] struct {
	pressed  []A
	released []A
	held     map[A]struct{}
}

// NewInputCollector creates an empty InputCollector.
func NewInputCollector[A comparable]() *InputCollector[A] {
	return &InputCollector[A]{held: make(map[A]struct{})}
}

// Press records a press event and marks the action as held.
func (c *InputCollector[A]) Press(a A) {
	c.pressed = append(c.pressed, a)
	c.held[a] = struct{}{}
}

// Release records a release event and drops the action from held.
func (c *InputCollector[A]) Release(a A) {
	c.released = append(c.released, a)
	delete(c.held, a)
}

// Drain converts buffered events into an InputFrame and clears the buffer.
// Press and release events are per-frame: each appears in exactly one
// snapshot; held state persists across frames until release.
func (c *InputCollector[A]) Drain() InputFrame[A] {
	pressed := make(map[A]struct{}, len(c.pressed))
	for _, a := range c.pressed {
		pressed[a] = struct{}{}
	}

	released := make(map[A]struct{}, len(c.released))
	for _, a := range c.released {
		released[a] = struct{}{}
	}

	held := make(map[A]struct{}, len(c.held))
	maps.Copy(held, c.held)

	frame := InputFrame[A]{pressed: pressed, released: released, held: held}

	c.pressed = nil
	c.released = nil

	return frame
}

// InputStage returns a Before stage that drains the collector into
// Frame.Input at the start of every frame.
func InputStage[E comparable, A comparable](c *InputCollector[A]) Stage[E, A] {
	return func(_ context.Context, f *Frame[E, A]) error {
		f.Input = c.Drain()
		return nil
	}
}
