package input

import (
	"testing"
)

func TestCollectorPressReleaseSnapshot(t *testing.T) {
	t.Parallel()

	c := NewInputCollector()

	frame := c.Drain()
	if !frame.Empty() {
		t.Fatal("empty collector should drain an empty frame")
	}

	c.Press(MoveUp)
	c.Press(MoveDown)
	c.Release(MoveDown)

	frame = c.Drain()
	if !frame.Pressed(MoveUp) || !frame.Pressed(MoveDown) {
		t.Error("both pressed actions should be in the snapshot")
	}
	if !frame.Released(MoveDown) {
		t.Error("released action should be in the snapshot")
	}
	if !frame.Held(MoveUp) {
		t.Error("MoveUp should be held while never released")
	}
	if frame.Held(MoveDown) {
		t.Error("MoveDown should not be held after release")
	}

	// Press and release events are per-frame and must not leak into the
	// next snapshot; held state persists.
	next := c.Drain()
	if next.Pressed(MoveUp) || next.Pressed(MoveDown) || next.Released(MoveDown) {
		t.Errorf("press/release events leaked into next frame: %+v", next)
	}
	if !next.Held(MoveUp) {
		t.Error("held state should persist into the next frame")
	}
}

func TestCollectorHeldPersistsUntilRelease(t *testing.T) {
	t.Parallel()

	c := NewInputCollector()

	c.Press(MoveLeft)
	_ = c.Drain()

	second := c.Drain()
	if !second.Held(MoveLeft) {
		t.Error("held state should persist across frames")
	}

	c.Release(MoveLeft)
	third := c.Drain()
	if third.Held(MoveLeft) {
		t.Error("held state should be dropped after release")
	}
	if !third.Released(MoveLeft) {
		t.Error("release should be reported in the release frame")
	}
}

func TestCollectorDuplicatePressesDeduplicated(t *testing.T) {
	t.Parallel()

	c := NewInputCollector()

	c.Press(MoveRight)
	c.Press(MoveRight)

	frame := c.Drain()
	if !frame.Pressed(MoveRight) {
		t.Error("MoveRight should be pressed")
	}
	if !frame.Held(MoveRight) {
		t.Error("MoveRight should be held")
	}
}
