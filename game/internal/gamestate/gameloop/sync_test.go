package gameloop

import (
	"context"
	"testing"
)

type fakeSync struct {
	msgs   []ServerMessage
	queued []InputPacket
}

func (f *fakeSync) Drain() []ServerMessage { return f.msgs }
func (f *fakeSync) Enqueue(p InputPacket)  { f.queued = append(f.queued, p) }

func TestNoopSync(t *testing.T) {
	t.Parallel()

	var s Sync = NoopSync{}

	if msgs := s.Drain(); msgs != nil {
		t.Errorf("NoopSync.Drain should return nil, got %v", msgs)
	}
	s.Enqueue(InputPacket{})
}

func TestSyncStageDrainsMessages(t *testing.T) {
	t.Parallel()

	msgs := []ServerMessage{{}, {}}
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

	stage := SyncStage[int, testAction](NoopSync{})

	f := &Frame[int, testAction]{}
	if err := stage(context.Background(), f); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if f.ServerMessages != nil {
		t.Errorf("expected no server messages, got %v", f.ServerMessages)
	}
}
