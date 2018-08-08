package impl 

import (
	"fmt"
	"net"

	"github.com/invin/kkchain/crypto/ed25519"
	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/gogo/protobuf/types"
	"github.com/golang/glog"
)


// Network represents the whole stack of p2p communication between peers
type Network struct {
	conf p2p.Config
	host p2p.Host
	id p2p.ID
	// Node's keypair.
	keys *crypto.KeyPair
}

// New creates a new Network instance with the specified configuration
func New(address string, conf p2p.Config) *Network {
	keys := ed25519.RandomKeyPair()
	id := p2p.CreateID(address, keys.PublicKey)

	return &Network{
		conf: conf,
		host: NewHost(address),
		id: id,
		keys: keys,
	}
}

// ID returns the identify of the local peer
func (n *Network) ID() p2p.ID {
	return n.id
}

// Start kicks off the p2p stack
func (n *Network) Start() error {
	// TODO: 

	return nil
}

// Conf gets configurations
func (n *Network) Conf() p2p.Config {
	return n.conf
}

// Stop stops the p2p stack
func (n *Network) Stop() {
	// TODO: 
}

// Accept connection
// FIXME: reference implementation
func (n *Network) Accept(incoming net.Conn) {
	conn := NewConnection(incoming, n)

	defer func() {
		if incoming != nil {
			incoming.Close()
		}
	}()

	for {
		msg, err := conn.ReadMessage()
		if err != nil {
			glog.Error(err)
			break
		}

		err = n.dispatchMessage(conn, msg)

		if err != nil {
			glog.Error(err)
			break
		}
	}
	
	fmt.Println("exit loop for connection")
}

// dispatch message according to protocol
func (n *Network) dispatchMessage(conn p2p.Conn, msg *protobuf.Message) error {
	// get stream handler
	handler, err := n.host.GetStreamHandler(msg.Protocol)
	if err != nil {
		return err 
	}

	// unmarshal message
	var ptr types.DynamicAny
	if err = types.UnmarshalAny(msg.Message, &ptr); err != nil {
		glog.Error(err)
		return err
	}
	
	// handle message
	stream := NewStream(conn, msg.Protocol)
	handler(stream, ptr.Message)
	// TODO: 
	return nil
}

// Sign signs a message
func (n *Network) Sign(message []byte) ([]byte, error) {
	return n.keys.Sign(n.conf.SignaturePolicy, n.conf.HashPolicy, message)
}

// Verify verifies the message
func (n *Network) Verify(publicKey []byte, message []byte, signature []byte) bool {
	return crypto.Verify(n.conf.SignaturePolicy, n.conf.HashPolicy, publicKey, message, signature)
}
