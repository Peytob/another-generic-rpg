package gameloop

import (
	"context"
	"testing"

	"game/internal/input"
	"game/internal/sync"
)

type fakeSync struct {
	msgs   []sync.ServerMessage
	queued []sync.InputPacket
}

func (f *fakeSync) Drain() []sync.ServerMessage { return f.msgs }
func (f *fakeSync) Enqueue(p sync.InputPacket)  { f.queued = append(f.queued, p) }

func TestInputStageDrainsCollector(t *testing.T) {
	t.Parallel()

	c := input.NewInputCollector[input.Action]()
	stage := InputStage[int](c)

	f := &Frame[int, input.Action]{}
	c.Press(input.MoveUp)

	if err := stage(context.Background(), f); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if !f.Input.Pressed(input.MoveUp) {
		t.Error("input stage should drain the collector into the frame")
	}

	f2 := &Frame[int, input.Action]{}
	if err := stage(context.Background(), f2); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if f2.Input.Pressed(input.MoveUp) || !f2.Input.Held(input.MoveUp) {
		t.Error("second frame should not see previous press events, only held state")
	}
}

func TestSyncStageDrainsMessages(t *testing.T) {
	t.Parallel()

	msgs := []sync.ServerMessage{{}, {}}
	stage := SyncStage[int, testAction](&fakeSync{msgs: msgs})

	f := &Frame[int, testAction]{}
	if err := stage(context.Background(), f); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if len(f.ServerMessages) != 2 {
		t.Errorf("expected 2 server messages in frame, got %d", len(f.ServerMessages))
	}
}

func TestSyncStageNoopLeavesFrameEmpty(t *testing.T) {
	t.Parallel()

	stage := SyncStage[int, testAction](sync.NoopSync{})

	f := &Frame[int, testAction]{}
	if err := stage(context.Background(), f); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if f.ServerMessages != nil {
		t.Errorf("expected no server messages, got %v", f.ServerMessages)
	}
}
