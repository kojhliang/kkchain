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


// P2P represents the whole stack of p2p communication between peers
type P2P struct {
	conf p2p.Config
	host p2p.Host
	id p2p.ID
	// Node's keypair.
	keys *crypto.KeyPair
}

// New creates a new P2P instance with the specified configuration
func New(address string, conf p2p.Config) *P2P {
	keys := ed25519.RandomKeyPair()
	id := p2p.CreateID(address, keys.PublicKey)

	return &P2P{
		conf: conf,
		host: NewHost(address),
		id: id,
		keys: keys,
	}
}

// ID returns the identify of the local peer
func (p2p *P2P) ID() p2p.ID {
	return p2p.id
}

// Start kicks off the p2p stack
func (p2p *P2P) Start() error {
	// TODO: 

	return nil
}

// Conf gets configurations
func (p2p *P2P) Conf() p2p.Config {
	return p2p.conf
}

// Stop stops the p2p stack
func (p2p *P2P) Stop() {
	// TODO: 
}

// Accept connection
// FIXME: reference implementation
func (p2p *P2P) Accept(incoming net.Conn) {
	conn := NewConnection(incoming, p2p)

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

		err = p2p.dispatchMessage(conn, msg)

		if err != nil {
			glog.Error(err)
			break
		}
	}
	
	fmt.Println("exit loop for connection")
}

// dispatch message according to protocol
func (p2p *P2P) dispatchMessage(conn p2p.Conn, msg *protobuf.Message) error {
	handler, err := p2p.host.GetStreamHandler(msg.Protocol)
	if err != nil {
		return err 
	}

	var ptr types.DynamicAny
	if err = types.UnmarshalAny(msg.Message, &ptr); err != nil {
		glog.Error(err)
		return err
	}
	
	stream := NewStream(conn, msg.Protocol)
	handler(stream, ptr.Message)
	// TODO: 
	return nil
}


