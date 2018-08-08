package p2p

import (
	"net"
	"io"

	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/gogo/protobuf/proto"
)

// Config defines configurations for a basic network instance
type Config struct {
	SignaturePolicy	crypto.SignaturePolicy
	HashPolicy crypto.HashPolicy	
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
	Accept(incoming net.Conn)

	// returns ID of local peer
	ID() ID

	// sign message
	Sign(message []byte) ([]byte, error)

	// verify message
	Verify(publicKey []byte, message []byte, signature []byte) bool
	
}

// Conn wraps connection related operations, such as reading and writing 
// messages
type Conn interface {
	io.Closer
	ReadMessage()(*protobuf.Message, error)
	WriteMessage(*protobuf.Message) error
	PrepareMessage(message proto.Message) (*protobuf.Message, error)
	GetPeerID() ID
}