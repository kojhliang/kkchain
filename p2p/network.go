package p2p

import (
	"io"
	"net"

	"github.com/gogo/protobuf/proto"
	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/p2p/protobuf"
)

//// Config defines configurations for a basic network instance
type Config struct {
	SignaturePolicy crypto.SignaturePolicy
	HashPolicy      crypto.HashPolicy
}

// Network defines the interface for network communication
type Network interface {
	// Start kicks off the network stack
	Start() error

	// Get configurations
	Conf() Config

	// Stop the network stack
	Stop()

	// Accept connection
	Accept(incoming Conn) error

	// Sign message
	Sign(message []byte) ([]byte, error)

	// Verify message
	Verify(publicKey []byte, message []byte, signature []byte) bool

	CreateConnection(fd net.Conn) (Conn, error)
	CreateStream(conn Conn, protocol string) (Stream, error)
}

// Conn wraps connection related operations, such as reading and writing
// messages
type Conn interface {
	io.Closer

	// Read message
	ReadMessage() (*protobuf.Message, error)

	// Write message
	WriteMessage(*protobuf.Message) error

	// Prepare message for further process
	PrepareMessage(message proto.Message) (*protobuf.Message, error)

	// Returns remote peer
	RemotePeer() ID
}
