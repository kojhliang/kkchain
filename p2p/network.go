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
	Accept(listener net.Listener)

	// Sign message
	Sign(message []byte) ([]byte, error)

	// Verify message
	Verify(publicKey []byte, message []byte, signature []byte) bool

	CreateConnection(fd net.Conn) (Conn, error)
	CreateStream(conn Conn, protocol string) (Stream, error)

	Bootstraps() []string

	RecvMessage()
	GetConnChan() *chan Conn
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

// Notifiee is an interface for an object wishing to receive
// notifications from network. Notifiees should take care not to register other
// notifiees inside of a notification.  They should also take care to do as
// little work as possible within their notification, putting any blocking work
// out into a goroutine.
type Notifiee interface {
	Connected(Conn)      // called when a connection opened
	Disconnected(Conn)   // called when a connection closed
	OpenedStream(Stream) // called when a stream opened
	ClosedStream(Stream) // called when a stream closed
}
