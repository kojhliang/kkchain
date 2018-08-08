package p2p

import (
	"github.com/gogo/protobuf/proto"
)


// Stream defines interface for reading and writing pb messages
type Stream interface {
	// write message
	Write(message proto.Message) error

	// returns peer ID
	GetPeerID() ID

	// return underlying connection
	Conn() Conn

	// reset stream
	Reset() error
}

// StreamHandler is the type of function used to listen for
// streams opened by the remote side.
type StreamHandler func(Stream, proto.Message)





