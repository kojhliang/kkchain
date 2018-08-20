package impl

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"

	"github.com/gogo/protobuf/types"
	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/chain"
	"github.com/invin/kkchain/p2p/dht"
	"github.com/invin/kkchain/p2p/handshake"
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	errServerStopped = errors.New("server stopped")
	log              = logrus.New()
)

// Network represents the whole stack of p2p communication between peers
type Network struct {
	conf p2p.Config
	host p2p.Host
	// Node's keypair.
	keys *crypto.KeyPair
	dht  *dht.DHT
	//BootstrapNodes []*Node
	BootstrapNodes []string
	listenAddr     string
	running        bool
	connChan       chan p2p.Conn
	quit           chan struct{}
	lock           sync.Mutex
	loopWG         sync.WaitGroup
}

// NewNetwork creates a new Network instance with the specified configuration
func NewNetwork(privateKeyPath, address string, conf p2p.Config) *Network {
	keys, _ := p2p.LoadNodeKeyFromFileOrCreateNew(privateKeyPath)
	id := p2p.CreateID(address, keys.PublicKey)

	return &Network{
		conf:       conf,
		host:       NewHost(id),
		keys:       keys,
		listenAddr: address,
		connChan:   make(chan p2p.Conn),
		quit:       make(chan struct{}),
	}
}

func (n *Network) GetConnChan() *chan p2p.Conn {
	return &n.connChan
}

// Start kicks off the p2p stack
func (n *Network) Start() error {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.running {
		return errors.New("server already running")
	}
	n.running = true

	log.Info("start P2P network")

	if n.keys == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	// init handshake msg handler
	handshake.NewHandshake(n.host)
	// init chain msg handler
	chain.NewChain(n.host)
	// set host to handle dht msg,and run dht
	config := dht.DefaultConfig()
	//config.Listen = n.listenAddr
	n.dht = dht.NewDHT(config, n, n.host)

	// do dht
	go func() {
		n.dht.Start()
		select {
		case <-n.quit:
			n.dht.Stop()
			n.host.RemoveAllConnection()
		}
	}()

	// TODO:process chain sync

	// listen
	if n.listenAddr != "" {
		if err := n.startListening(); err != nil {
			return err
		}
	} else {
		log.Warn("P2P server will be useless, not listening")
	}

	// dail
	n.loopWG.Add(1)
	go n.run()

	n.loopWG.Add(1)
	go n.RecvMessage()

	return nil
}

// Conf gets configurations
func (n *Network) Conf() p2p.Config {
	return n.conf
}

func (n *Network) Bootstraps() []string {
	return n.BootstrapNodes
}

// Stop stops the p2p stack
func (n *Network) Stop() {
	n.lock.Lock()
	defer n.lock.Unlock()
	if !n.running {
		return
	}
	n.running = false
	close(n.quit)
	n.loopWG.Wait()
}

func (n *Network) startListening() error {
	addr, err := dht.ToNetAddr(n.listenAddr)
	if err != nil {
		return err
	}

	listener, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	n.listenAddr = laddr.String()
	n.loopWG.Add(1)
	go n.Accept(listener)
	return nil
}

func (n *Network) run() {
	defer n.loopWG.Done()

	// connect boostnode
	for _, node := range n.BootstrapNodes {
		peer, err := dht.ParsePeerAddr(node)
		if err != nil {
			continue
		}
		conn, _ := n.host.GetConnection(peer.ID)
		if conn != nil {
			continue
		}
		go func() {
			conn, err := n.host.Connect(peer.Address, n)
			if err != nil {
				log.WithFields(logrus.Fields{
					"address": peer.Address,
					"nodeID":  hex.EncodeToString(peer.ID.PublicKey),
				}).Error("failed to connect boost node")
			} else {
				log.WithFields(logrus.Fields{
					"address": peer.Address,
					"nodeID":  hex.EncodeToString(peer.ID.PublicKey),
				}).Info("success to connect boost node")
				msg := handshake.NewMessage(handshake.Message_HELLO)
				handshake.BuildHandshake(msg)
				stream, err := n.CreateStream(conn, "/kkchain/p2p/handshake/1.0.0")
				if err != nil {
					log.Error(err)
				} else {
					err := stream.Write(msg)
					if err != nil {
						log.Error(err)
					}
					n.host.AddConnection(peer.ID, conn)
					n.connChan <- conn
				}
			}
		}()
	}
}

// Accept connection
// FIXME: reference implementation
func (n *Network) Accept(listener net.Listener) {
	defer n.loopWG.Done()
	n.lock.Lock()
	running := n.running
	n.lock.Unlock()
	if !running {
		log.Error(errServerStopped)
		return
	}

	for {
		var (
			fd  net.Conn
			err error
		)

		fd, err = listener.Accept()
		if err != nil {
			log.Error("failed to listen:", err)
			break
		}
		fmt.Println("listener accepted...")
		conn := NewConnection(fd, n, n.host)
		if conn == nil {
			log.Error(failedNewConnection)
			continue
		}

		n.connChan <- conn
	}
}

func (n *Network) RecvMessage() {
	defer n.loopWG.Done()
	for {
		select {
		case conn := <-n.connChan:
			go func() {
				for {
					msg, err := conn.ReadMessage()
					if err != nil {
						continue
					}
					fmt.Println("\n接受的消息：", msg.Sender, msg.Message.TypeUrl, "\n")
					err = n.dispatchMessage(conn, msg)
					if err != nil {
						continue
					}
				}
			}()
		}
	}
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
		log.Error(err)
		return err
	}

	// handle message
	stream := NewStream(conn, msg.Protocol)
	handler(stream, ptr.Message)

	return nil
}

// Sign signs a message
// TODO: move to another package??
func (n *Network) Sign(message []byte) ([]byte, error) {
	return n.keys.Sign(n.conf.SignaturePolicy, n.conf.HashPolicy, message)
}

// Verify verifies the message
// TODO: move to another package??
func (n *Network) Verify(publicKey []byte, message []byte, signature []byte) bool {
	return crypto.Verify(n.conf.SignaturePolicy, n.conf.HashPolicy, publicKey, message, signature)
}

// create connection
func (n *Network) CreateConnection(fd net.Conn) (p2p.Conn, error) {
	conn := NewConnection(fd, n, n.host)
	if conn == nil {
		return nil, failedNewConnection
	}
	return conn, nil
}

// create stream
func (n *Network) CreateStream(conn p2p.Conn, protocol string) (p2p.Stream, error) {
	stream := NewStream(conn, protocol)
	if stream == nil {
		return nil, failedCreateStream
	}
	return stream, nil
}
