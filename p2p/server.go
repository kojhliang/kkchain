package p2p

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"encoding/hex"

	"github.com/golang/protobuf/proto"
	"github.com/invin/kkchain/p2p/dht"
	"github.com/invin/kkchain/p2p/impl"
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/libp2p/go-libp2p-crypto"
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
	PrivateKey      crypto.PrivKey
	MaxPeers        int
	MaxPendingPeers int
	BootstrapNodes  []*Node
	StaticNodes     []*Node
	ListenAddr      string
	Dialer          NodeDialer
}

// Server manages all peer connections.
type Server struct {
	ServerConfig
	SendMsg      map[string]interface{}
	RecvMsg      map[string]interface{}
	peers        map[ID]*Peer
	lock         sync.Mutex
	running      bool
	listener     net.Listener
	quit         chan struct{}
	addstatic    chan *Node
	removestatic chan *Node
	addpeer      chan *conn
	delpeer      chan *conn
	loopWG       sync.WaitGroup
}

type connFlag int

const (
	dynDialedConn connFlag = 1 << iota
	staticDialedConn
	inboundConn
)

// peer connection info
type conn struct {
	conn    impl.Connection
	flags   connFlag
	connErr chan error
	id      dht.PeerID
}

func (c *conn) isThisFlag(f connFlag) bool {
	return c.flags&f != 0
}

func (srv *Server) AddPeer(node *Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

func (srv *Server) RemovePeer(node *Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

func (srv *Server) Self() *Node {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return &Node{IP: net.ParseIP("0.0.0.0")}
	}
	return srv.makeSelf(srv.listener)
}

func (srv *Server) makeSelf(listener net.Listener) *Node {
	pubkey, err := srv.PrivateKey.GetPublic().Bytes()
	if err != nil {
		log.Error("failed to get server's pubkey:", err)
		return nil
	}
	if listener == nil {
		return &Node{IP: net.ParseIP("0.0.0.0"), ID: dht.CreateID("0.0.0.0", pubkey)}
	}
	addr := listener.Addr().(*net.TCPAddr)
	return &Node{
		ID:  dht.CreateID(addr.String(), pubkey),
		IP:  addr.IP,
		TCP: uint16(addr.Port),
	}
}

func (srv *Server) Stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	if srv.listener != nil {
		srv.listener.Close()
	}
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

	if srv.PrivateKey == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *conn)
	srv.delpeer = make(chan *conn)
	srv.addstatic = make(chan *Node)
	srv.removestatic = make(chan *Node)

	dialer := newDialState(srv.StaticNodes, srv.BootstrapNodes)

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
	srv.listener = listener
	srv.loopWG.Add(1)
	go srv.listenLoop()
	return nil
}

func (srv *Server) run(dialstate *dialstate) {
	defer srv.loopWG.Done()

	// todo：dht
	startTask := func() {
		lenBoost := len(srv.BootstrapNodes)
		if lenBoost == 0 {
			return
		}
		for i := 0; lenBoost < maxActiveDialTasks && i < lenBoost; i++ {
			log.Info("new dial task")
			t := &dialTask{
				dynDialedConn,
				srv.BootstrapNodes[i],
			}
			t.Do(srv)
		}
	}

running:
	for {
		startTask()
		select {
		case <-srv.quit:
			break running
		case n := <-srv.addstatic:
			log.Debug("add static node:", n)
			// TODO:

		case n := <-srv.removestatic:
			log.Debug("remove static node:", n)
			// TODO:

		case c := <-srv.addpeer:
			p := newPeer(&c.conn)
			log.WithFields(logrus.Fields{
				"remote_addr": c.conn.RemoteAddr().String(),
				"nodeID":      hex.EncodeToString(p.ID.PublicKey),
			}).Info("add p2p peer")

			// TODO:run peer protocols

		case c := <-srv.delpeer:
			log.WithFields(logrus.Fields{
				"remote_addr": c.conn.RemoteAddr().String(),
				"nodeID":      hex.EncodeToString(c.conn.RemotePeer().PublicKey),
			}).Debug("remove p2p peer:")
			// TODO:
		}
	}

	// when server quit , close all connection
	for _, p := range srv.peers {
		p.Disconnect(DiscQuitting)
	}

	for len(srv.peers) > 0 {
		p := <-srv.delpeer
		delete(srv.peers, p.conn.RemotePeer())
	}
}

func (srv *Server) listenLoop() {
	defer srv.loopWG.Done()
	log.Info("listener up self:", srv.makeSelf(srv.listener))

	for {
		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = srv.listener.Accept()
			if err != nil {
				log.Error("failed to listen:", err)
				return
			}
			break
		}
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
		}()
	}
}

// create connection
func (srv *Server) SetupConn(fd net.Conn, flags connFlag, dialDest *Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}

	implConn := impl.Connection{
		fd,
		nil,
		nil,
		nil,
		nil,
	}
	c := &conn{
		conn:    implConn,
		flags:   flags,
		connErr: make(chan error),
	}
	err := srv.setupConn(c, dialDest)
	if err != nil {
		c.connErr <- err
		log.WithFields(logrus.Fields{
			"id":    c.id,
			"error": err,
		}).Error("failed to set up connection")
	}
	return nil
}

// todo:certification
func (srv *Server) setupConn(c *conn, dialDest *Node) error {
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}

	if c.isThisFlag(inboundConn) {

		// inbound conn , should recv peerID
		msg, err := c.conn.ReadMessage()
		if err != io.EOF {
			log.Error("failed to read data from setupConn:", err)
			return err
		}
		readBytes, err := proto.Marshal(msg)
		if err != nil {
			log.Error("failed to marshal protocol msg:", err)
			return err
		}
		err = json.Unmarshal(readBytes, &srv.RecvMsg)
		if err != nil {
			log.Error("failed to ummarshal bytes:", err)
			return err
		}
		nodeIDBytes, err := json.Marshal(srv.RecvMsg["PeerID"])
		if err != nil {
			return err
		}
		err = json.Unmarshal(nodeIDBytes, &c.id)
		if err != nil {
			return err
		}
	} else {

		// dial node , should send self peerID
		srv.SendMsg["PeerID"] = srv.Self().ID
		_, err := json.Marshal(srv.SendMsg)
		if err != nil {
			return err
		}

		// TODO:from data to protobuf msg
		msg := &protobuf.Message{}
		err = c.conn.WriteMessage(msg)
		if err != nil {
			return err
		}
	}

	if dialDest != nil && c.id.Equals(dialDest.ID) {
		log.WithFields(logrus.Fields{
			"expected_id": c.id,
			"want_id":     dialDest.ID,
		}).Error("dialed identity mismatch")
		return DiscUnexpectedIdentity
	}

	// addConn
	srv.addpeer <- c
	log.WithFields(logrus.Fields{
		"id":   c.id,
		"addr": c.conn.RemoteAddr(),
		"conn": c.flags,
	}).Info("connection set up", "inbound", dialDest == nil)
	return nil
}

type NodeInfo struct {
	ID         string
	Enode      string
	IP         string
	Port       int
	ListenAddr string
}

func (srv *Server) NodeInfo() *NodeInfo {
	node := srv.Self()
	info := &NodeInfo{
		Enode:      node.String(),
		ID:         node.ID.String(),
		IP:         node.IP.String(),
		ListenAddr: srv.ListenAddr,
	}
	info.Port = int(node.TCP)
	return info
}

func (srv *Server) PeersInfo() []*PeerInfo {
	infos := make([]*PeerInfo, 0, len(srv.peers))
	for _, peer := range srv.peers {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ID > infos[j].ID {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}
	return infos
}
