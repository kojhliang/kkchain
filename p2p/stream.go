package p2p

import (
	"github.com/gogo/protobuf/proto"
)

// Stream defines interface for reading and writing pb messages
type Stream interface {
	// Write message
	Write(message proto.Message) error

	// Returns remote peer ID
	RemotePeer() ID

	// Return underlying connection
	Conn() Conn

	// Reset stream
	Reset() error

	// Get protocol
	Protocol() string
}

// StreamHandler is the type of function used to listen for
// streams opened by the remote side.
type StreamHandler func(Stream, proto.Message)
