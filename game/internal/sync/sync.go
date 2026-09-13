// Package sync abstracts the client-server communication channel.
package sync

// Sync abstracts the client-server communication channel. Implementations
// must be non-blocking: Drain returns immediately with whatever arrived,
// Enqueue only buffers packets for the network goroutine.
type Sync interface {
	Drain() []ServerMessage
	Enqueue(p ClientMessage)
}
