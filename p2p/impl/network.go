package impl

import (
	"fmt"
	"net"
	"sync"

	"time"

	"strings"

	"encoding/hex"

	"github.com/gogo/protobuf/types"
	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/crypto/ed25519"
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

const (
	defaultDialTimeout = 15 * time.Second
	defaultDBPath      = "./nodedb"
)

type connFlag int

const (
	outboundConn connFlag = iota
	inboundConn
)

// Network represents the whole stack of p2p communication between peers
type Network struct {
	conf p2p.Config
	host p2p.Host
	// Node's keypair.
	keys *crypto.KeyPair

	dht            *dht.DHT
	BootstrapNodes []*Node
	listenAddr     string
	running        bool
	quit           chan struct{}
	lock           sync.Mutex
	loopWG         sync.WaitGroup
}

// NewNetwork creates a new Network instance with the specified configuration
func NewNetwork(address string, conf p2p.Config) *Network {
	keys := ed25519.RandomKeyPair()
	id := p2p.CreateID(address, keys.PublicKey)

	return &Network{
		conf:       conf,
		host:       NewHost(id),
		keys:       keys,
		listenAddr: address,
	}
}

func (n *Network) Self() *Node {
	n.lock.Lock()
	defer n.lock.Unlock()
	if !n.running {
		return nil
	}
	return n.makeSelf(n.listenAddr)
}

func (n *Network) makeSelf(listenAddr string) *Node {
	pubkey := n.keys.PublicKey
	if listenAddr == "" {
		return &Node{IP: "0.0.0.0", ID: p2p.CreateID("0.0.0.0", pubkey)}
	}
	addr := strings.Split(listenAddr, ":")
	return &Node{
		ID:      p2p.CreateID(listenAddr, pubkey),
		IP:      addr[0],
		TCPPort: addr[1],
	}
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

	n.quit = make(chan struct{})

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

	n.loopWG.Add(1)

	// dail
	go n.run()
	n.running = true

	return nil
}

// Conf gets configurations
func (n *Network) Conf() p2p.Config {
	return n.conf
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
	listener, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	n.listenAddr = laddr.String()
	n.loopWG.Add(1)
	go n.listenLoop(listener)
	return nil
}

func (n *Network) run() {
	defer n.loopWG.Done()

	for {

		// connect boostnode
		for _, node := range n.BootstrapNodes {
			conn, _ := n.host.GetConnection(node.ID)
			if conn != nil {
				continue
			}
			go func() {
				fd, err := n.host.Connect(node.Addr())
				if err != nil {
					log.WithFields(logrus.Fields{
						"address": node.Addr(),
						"nodeID":  hex.EncodeToString(node.ID.PublicKey),
					}).Error("failed to connect boost node")
				} else {
					log.WithFields(logrus.Fields{
						"address": node.Addr(),
						"nodeID":  hex.EncodeToString(node.ID.PublicKey),
					}).Info("success to connect boost node")
					msg := handshake.NewMessage(handshake.Message_HELLO)
					handshake.BuildHandshake(msg)
					conn, err = n.CreateConnection(fd)
					if err != nil {
						log.Error(err)
					} else {
						stream, err := n.CreateStream(conn, "/kkchain/p2p/handshake/1.0.0")
						if err != nil {
							log.Error(err)
						} else {
							err := stream.Write(msg)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}()
		}
		select {
		case <-n.quit:
			break
		}
	}

	n.host.RemoveAllConnection()
}

func (n *Network) listenLoop(listener net.Listener) {
	defer n.loopWG.Done()
	for {
		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = listener.Accept()
			if err != nil {
				log.Error("failed to accept:", err)
				return
			}
			break
		}
		go func() {
			n.SetupConn(fd, inboundConn, nil)
		}()
	}
}

// create connection
func (n *Network) SetupConn(fd net.Conn, flag connFlag, dialDest *Node) error {
	self := n.Self()
	if self == nil {
		return errors.New("shutdown")
	}
	err := n.setupConn(fd, flag, dialDest)
	if err != nil {
		log.WithFields(logrus.Fields{
			"address": fd.RemoteAddr().String(),
			"error":   err,
		}).Error("failed to set up connection")
		return err
	}
	return nil
}

func (n *Network) setupConn(fd net.Conn, flag connFlag, dialDest *Node) error {
	n.lock.Lock()
	running := n.running
	n.lock.Unlock()
	if !running {
		return errServerStopped
	}

	if flag == inboundConn {
		conn := NewConnection(fd, n, n.host)
		if conn == nil {
			return failedNewConnection
		}

		err := n.Accept(conn)
		if err != nil {
			return err
		}

		existConn, err := n.host.GetConnection(conn.remotePeer)
		if err != nil {
			return err
		}
		if conn == existConn {

			// when success to accept conn,notify dht the remote peer ID
			n.dht.GetRecvchan() <- conn.remotePeer
			log.WithFields(logrus.Fields{
				"addr":      fd.RemoteAddr().String(),
				"conn_flag": "inbound",
			}).Info("accept connection")
		}
	} else {

		// outbound conn

	}
	return nil
}

// Accept connection
// FIXME: reference implementation
func (n *Network) Accept(conn p2p.Conn) error {
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	msg, err := conn.ReadMessage()
	if err != nil {
		log.Error(err)
		return err
	}

	err = n.dispatchMessage(conn, msg)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
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
