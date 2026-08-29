package gameloop

import (
	"testing"
)

func TestCollectorPressReleaseSnapshot(t *testing.T) {
	t.Parallel()

	c := NewInputCollector[testAction]()

	frame := c.Drain()
	if !frame.Empty() {
		t.Fatal("empty collector should drain an empty frame")
	}

	c.Press(testActionA)
	c.Press(testActionB)
	c.Release(testActionB)

	frame = c.Drain()
	if !frame.Pressed(testActionA) || !frame.Pressed(testActionB) {
		t.Error("both pressed actions should be in the snapshot")
	}
	if !frame.Released(testActionB) {
		t.Error("released action should be in the snapshot")
	}
	if !frame.Held(testActionA) {
		t.Error("A should be held while never released")
	}
	if frame.Held(testActionB) {
		t.Error("B should not be held after release")
	}

	// Press and release events are per-frame and must not leak into the
	// next snapshot; held state persists.
	next := c.Drain()
	if next.Pressed(testActionA) || next.Pressed(testActionB) || next.Released(testActionB) {
		t.Errorf("press/release events leaked into next frame: %+v", next)
	}
	if !next.Held(testActionA) {
		t.Error("held state should persist into the next frame")
	}
}

func TestCollectorHeldPersistsUntilRelease(t *testing.T) {
	t.Parallel()

	c := NewInputCollector[testAction]()

	c.Press(testActionA)
	_ = c.Drain()

	second := c.Drain()
	if !second.Held(testActionA) {
		t.Error("held state should persist across frames")
	}

	c.Release(testActionA)
	third := c.Drain()
	if third.Held(testActionA) {
		t.Error("held state should be dropped after release")
	}
	if !third.Released(testActionA) {
		t.Error("release should be reported in the release frame")
	}
}

func TestCollectorDuplicatePressesDeduplicated(t *testing.T) {
	t.Parallel()

	c := NewInputCollector[testAction]()

	c.Press(testActionA)
	c.Press(testActionA)

	frame := c.Drain()
	if !frame.Pressed(testActionA) {
		t.Error("A should be pressed")
	}
	if !frame.Held(testActionA) {
		t.Error("A should be held")
	}
}

func TestInputStageDrainsCollector(t *testing.T) {
	t.Parallel()

	c := NewInputCollector[testAction]()
	stage := InputStage[int](c)

	f := &Frame[int, testAction]{}
	c.Press(testActionA)

	if err := stage(nil, f); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if !f.Input.Pressed(testActionA) {
		t.Error("input stage should drain the collector into the frame")
	}

	f2 := &Frame[int, testAction]{}
	if err := stage(nil, f2); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if f2.Input.Pressed(testActionA) || !f2.Input.Held(testActionA) {
		t.Error("second frame should not see previous press events, only held state")
	}
}
