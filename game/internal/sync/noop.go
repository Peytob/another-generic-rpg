package sync

// NoopSync Stub for disabled network
type NoopSync struct {
}

func (NoopSync) Drain() []ServerMessage {
	return nil
}

func (NoopSync) Enqueue(ClientMessage) {
}
