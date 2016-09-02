// Package hh provides a hinted handoff service for Pub.
package hh

type Service interface {

	// Start the hinted handoff service.
	Start() error

	// Stop the hinted handoff service.
	Stop()

	// Append add key/value byte slice to end of the buffer.
	Append(cluster, topic string, key, value []byte) error

	// Empty returns whether the buffer has no inflight entries.
	Empty(cluster, topic string) bool

	// FlushInflights flush all inflight entries inside buffer to final message storage.
	FlushInflights()
}

var Default Service
