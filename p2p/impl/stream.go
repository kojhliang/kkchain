package impl

import (
	"github.com/gogo/protobuf/proto"
	"github.com/invin/kkchain/p2p"
	"github.com/pkg/errors"
)

var (
	errEmptyMsg       = errors.New("received an empty message from peer")
	errInvalidMessage = errors.New("invalid message")
	errVerifySign     = errors.New("invalid signature")
)

// Stream implements the Stream interface and provides basic routiens
// for sending and writing pb messages
type Stream struct {
	conn     p2p.Conn
	protocol string
}

// NewStream creates a new stream
func NewStream(conn p2p.Conn, protocol string) *Stream {
	return &Stream{conn, protocol}
}

// Write wraps a pb message and send it
func (s *Stream) Write(message proto.Message) error {
	// sign the message
	signed, err := s.conn.PrepareMessage(message)
	if err != nil {
		return err
	}

	// set protocol
	signed.Protocol = s.protocol

	return s.conn.WriteMessage(signed)
}

// RemotePeer returns remote peer ID
func (s *Stream) RemotePeer() p2p.ID {
	return s.conn.RemotePeer()
}

// Conn returns the connection to which stream attaches
func (s *Stream) Conn() p2p.Conn {
	return s.conn
}

// Reset resets current stream
func (s *Stream) Reset() error {
	return s.conn.Close()
}

// Protocol returns protocl of this stream
func (s *Stream) Protocol() string {
	return s.protocol
}
