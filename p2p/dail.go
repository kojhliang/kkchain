package p2p

import (
	"errors"
	"fmt"
	"net"
	"time"
	"strings"
	"strconv"
)

type NodeDialer interface {
	Dial(*Node) (net.Conn, error)
}

type TCPDialer struct {
	*net.Dialer
}

func (t TCPDialer) Dial(dest *Node) (net.Conn, error) {
	addr := &net.TCPAddr{IP: dest.IP, Port: int(dest.TCP)}
	return t.Dialer.Dial("tcp", addr.String())
}

// dialstate schedules dials
type dialstate struct {
	dialing   map[NodeID]connFlag
	static    map[NodeID]*dialTask
	start     time.Time
	bootnodes []*Node
}

type task interface {
	Do(*Server)
}

type dialTask struct {
	flags connFlag
	dest  *Node
}

func newDialState(static []*Node, bootnodes []*Node) *dialstate {
	s := &dialstate{
		static:    make(map[NodeID]*dialTask),
		bootnodes: make([]*Node, len(bootnodes)),
	}
	copy(s.bootnodes, bootnodes)
	for _, n := range static {
		s.addStatic(n)
	}
	return s
}

func (s *dialstate) addStatic(n *Node) {
	s.static[n.ID] = &dialTask{flags: staticDialedConn, dest: n}
}

func (s *dialstate) removeStatic(n *Node) {
	delete(s.static, n.ID)
}

var (
	errSelf             = errors.New("is self")
	errAlreadyDialing   = errors.New("already dialing")
	errAlreadyConnected = errors.New("already connected")
	errRecentlyDialed   = errors.New("recently dialed")
	errNotWhitelisted   = errors.New("not contained in netrestrict whitelist")
)

func (s *dialstate) checkDial(n *Node, peers map[NodeID]*Peer) error {
	_, dialing := s.dialing[n.ID]
	switch {
	case dialing:
		return errAlreadyDialing
	case peers[n.ID] != nil:
		return errAlreadyConnected
	}
	return nil
}

func (t *dialTask) Do(srv *Server) {
	if t.dest.IP == nil {
		if !t.resolve(srv) {
			return
		}
	}
	err := t.dial(srv, t.dest)
	if err != nil {
		log.Error("failed to dial:", err)
	}
}

func (t *dialTask) resolve(srv *Server) bool {
	var resolved *Node
	peer := srv.peers[t.dest.ID]
	if peer != nil {
		addrStr := peer.RemoteAddr().String()
		addrArr := strings.Split(addrStr, ":")
		port, err := strconv.ParseUint(addrArr[1], 10, 16)
		if err != nil {
			log.Error("failed to parse string to uint:", err)
		}
		resolved = &Node{
			net.ParseIP(addrArr[0]),
			uint16(port),
			peer.ID(),
		}
	}
	t.dest = resolved
	return true
}

func (t *dialTask) dial(srv *Server, dest *Node) error {
	fd, err := srv.Dialer.Dial(dest)
	if err != nil {
		return err
	}
	return srv.SetupConn(fd, t.flags, dest)
}

func (t *dialTask) String() string {
	return fmt.Sprintf("%v %x %v:%d", t.flags, t.dest.ID[:8], t.dest.IP, t.dest.TCP)
}
