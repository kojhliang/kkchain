package p2p

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"strings"

	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/p2p/dht"
	"github.com/invin/kkchain/p2p/impl"
	"github.com/sirupsen/logrus"
)

const (
	defaultDialTimeout = 15 * time.Second
	maxActiveDialTasks = 16
)

var (
	errServerStopped = errors.New("server stopped")
	log              = logrus.New()
)

type ServerConfig struct {
	KeyPair         *crypto.KeyPair
	MaxPeers        int
	MaxPendingPeers int
	BootstrapNodes  []*Node
	ListenAddr      string
	DBPath          string
}

// Server manages all peer connections.
type Server struct {
	ServerConfig
	dht     *dht.DHT
	network *impl.Network
	host    *impl.Host
	peers   map[ID]*Peer
	dialer  NodeDialer
	running bool
	quit    chan struct{}
	lock    sync.Mutex
	loopWG  sync.WaitGroup
}

type connFlag int

const (
	outboundConn connFlag = iota
	inboundConn
)

func (srv *Server) Self() *Node {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return nil
	}
	return srv.makeSelf(srv.ListenAddr)
}

func (srv *Server) makeSelf(listenAddr string) *Node {
	pubkey := srv.KeyPair.PublicKey
	if listenAddr == "" {
		return &Node{IP: "0.0.0.0", ID: dht.CreateID("0.0.0.0", pubkey)}
	}
	addr := strings.Split(listenAddr, ":")
	return &Node{
		ID:      dht.CreateID(listenAddr, pubkey),
		IP:      addr[0],
		TCPPort: addr[1],
	}
}

func (srv *Server) Stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	close(srv.quit)
	srv.loopWG.Wait()
}

func (srv *Server) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true

	log.Info("start P2P network")

	if srv.KeyPair == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	if srv.dialer == nil {
		srv.dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}

	srv.quit = make(chan struct{})
	srv.network = impl.NewNetwork(srv.ListenAddr, Config{})
	ID := CreateID(srv.ListenAddr, srv.KeyPair.PublicKey)
	srv.host = impl.NewHost(ID)
	selfID := dht.CreateID(srv.ListenAddr, srv.KeyPair.PublicKey)
	srv.dht = dht.NewDHT(srv.host, srv.DBPath, selfID)

	dialer := newDialState(srv.BootstrapNodes, msg_dht)

	// listen
	if srv.ListenAddr != "" {
		if err := srv.startListening(); err != nil {
			return err
		}
	} else {
		log.Warn("P2P server will be useless, not listening")
	}

	srv.loopWG.Add(1)

	// dail
	go srv.run(dialer)
	srv.running = true
	return nil
}

func (srv *Server) startListening() error {
	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	srv.ListenAddr = laddr.String()
	srv.loopWG.Add(1)
	go srv.listenLoop(listener)
	return nil
}

func (srv *Server) run(dialstate *dialstate) {
	defer srv.loopWG.Done()

	startTask := func() {
		for _, node := range srv.BootstrapNodes {
			err := dialstate.checkDial(node)
			if err != nil {
				log.Warn("wrong dial task:", err)
				continue
			}
			t := dialstate.task[node.ID]
			go t.Do(srv)
		}
	}

running:
	for {
		startTask()
		select {
		case <-srv.quit:
			break running
		}
	}

	// when server quit , close all connection
	for _, p := range srv.peers {
		p.Disconnect(DiscQuitting)
	}
}

func (srv *Server) listenLoop(listener net.Listener) {
	defer srv.loopWG.Done()
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
			srv.SetupConn(fd, inboundConn, -1, nil)
		}()
	}
}

// create connection
func (srv *Server) SetupConn(fd net.Conn, flag connFlag, msgType MsgType, dialDest *Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}
	err := srv.setupConn(fd, flag, msgType, dialDest)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":    fd.RemoteAddr().String(),
			"error": err,
		}).Error("failed to set up connection")
		return err
	}
	return nil
}

// todo:certification
func (srv *Server) setupConn(fd net.Conn, flag connFlag, msgType MsgType, dialDest *Node) error {
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}

	if flag == inboundConn {

		// inbound conn
		srv.network.Accept(fd)
		peer := newPeer(impl.NewConnection(fd, srv.network, srv.host))
		srv.dht.AddPeer(peer.ID)
	} else {

		// outbound conn
		switch msgType {
		case msg_handshake:

		case msg_dht:

		case msg_chain:

		}
	}

	log.WithFields(logrus.Fields{
		"addr": fd.RemoteAddr().String(),
		"conn": flag,
	}).Info("connection set up")
	return nil
}
