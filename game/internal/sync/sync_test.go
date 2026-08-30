package sync

import (
	"testing"
)

func TestNoopSync(t *testing.T) {
	t.Parallel()

	var s Sync = NoopSync{}

	if msgs := s.Drain(); msgs != nil {
		t.Errorf("NoopSync.Drain should return nil, got %v", msgs)
	}
	s.Enqueue(InputPacket{})
}
