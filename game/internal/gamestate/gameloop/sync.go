package gameloop

import (
	"context"
)

// ServerMessage is a placeholder for messages received from the game server.
type ServerMessage struct{}

// InputPacket is a placeholder for client input sent to the game server.
type InputPacket struct{}

// Sync abstracts the client-server communication channel. Implementations
// must be non-blocking: Drain returns immediately with whatever arrived,
// Enqueue only buffers packets for the network goroutine.
type Sync interface {
	Drain() []ServerMessage
	Enqueue(p InputPacket)
}

// NoopSync is a Sync that does nothing; used until real networking appears.
type NoopSync struct{}

func (NoopSync) Drain() []ServerMessage { return nil }
func (NoopSync) Enqueue(InputPacket)    {}

// SyncStage returns a Before stage that drains pending server messages into
// Frame.ServerMessages for the simulation ticks to apply.
func SyncStage[E comparable, A comparable](s Sync) Stage[E, A] {
	return func(_ context.Context, f *Frame[E, A]) error {
		f.ServerMessages = s.Drain()
		return nil
	}
}
